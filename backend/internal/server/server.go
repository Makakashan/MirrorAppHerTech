package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type App struct {
	mu              sync.Mutex
	messages        []Message
	aiHistory       []nvidiaMessage
	metrics         Metrics
	ai              *NvidiaClient
	systemPrompt    string
	scores          map[string]float64
	sessionExcerpts []AnalysisExcerpt
	sessionTriggers map[string]int // trigger → occurrence count
}

func New() *App {
	now := time.Now()
	systemPrompt, scores := loadProfileAndBuildPrompt()
	return &App{
		messages:        initialMessages(now),
		aiHistory:       []nvidiaMessage{},
		metrics:         initialMetrics(now),
		ai:              NewNvidiaClientFromEnv(),
		systemPrompt:    systemPrompt,
		scores:          scores,
		sessionExcerpts: []AnalysisExcerpt{},
		sessionTriggers: map[string]int{},
	}
}

// sessionContext builds an addendum to the system prompt with what the AI
// has learned during this session — excerpts, recurring triggers, live scores.
// Called under mu lock.
func (a *App) sessionContext() string {
	if len(a.sessionExcerpts) == 0 && len(a.sessionTriggers) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("\n\n## Learned during this session (do not quote directly)\n")

	if len(a.sessionExcerpts) > 0 {
		sb.WriteString("\nNew excerpts:\n")
		for _, e := range a.sessionExcerpts {
			if e.Quote != nil && e.ContextTag != nil {
				fmt.Fprintf(&sb, "- \"%s\" [%s]\n", *e.Quote, *e.ContextTag)
			}
		}
	}

	if len(a.sessionTriggers) > 0 {
		sb.WriteString("\nRecurring triggers this session: ")
		first := true
		for t, count := range a.sessionTriggers {
			if !first {
				sb.WriteString(", ")
			}
			fmt.Fprintf(&sb, "%s×%d", t, count)
			first = false
		}
		sb.WriteString("\n")
	}

	sb.WriteString("\nLive scores (updated from conversation):\n")
	for _, k := range []string{"anxiety", "energy", "mood", "openness", "focus", "irritability"} {
		fmt.Fprintf(&sb, "- %s: %.1f\n", k, a.scores[k])
	}

	return sb.String()
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", a.handleHealth)
	mux.HandleFunc("GET /api/messages", a.handleMessages)
	mux.HandleFunc("POST /api/chat", a.handleChat)
	mux.HandleFunc("GET /api/suggestions", a.handleSuggestions)
	mux.HandleFunc("GET /api/metrics", a.handleMetrics)
	return requestLogger(cors(mux))
}

type statusWriter struct {
	http.ResponseWriter
	status int
}

func (sw *statusWriter) WriteHeader(code int) {
	sw.status = code
	sw.ResponseWriter.WriteHeader(code)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sw := &statusWriter{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(sw, r)
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, sw.status, time.Since(start).Round(time.Millisecond))
	})
}

func (a *App) WithStatic(api http.Handler, dist string) http.Handler {
	files := http.FileServer(http.Dir(dist))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			api.ServeHTTP(w, r)
			return
		}

		path := filepath.Join(dist, filepath.Clean(r.URL.Path))
		if r.URL.Path != "/" {
			if _, err := http.Dir(dist).Open(strings.TrimPrefix(path, dist)); err == nil {
				files.ServeHTTP(w, r)
				return
			}
		}
		http.ServeFile(w, r, filepath.Join(dist, "index.html"))
	})
}

func (a *App) handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *App) handleMessages(w http.ResponseWriter, _ *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	writeJSON(w, http.StatusOK, a.messages)
}

func (a *App) handleChat(w http.ResponseWriter, r *http.Request) {
	var req chatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	text := strings.TrimSpace(req.Message)
	if text == "" {
		writeError(w, http.StatusBadRequest, "message is required")
		return
	}
	if len(text) > 2000 {
		writeError(w, http.StatusBadRequest, "message is too long")
		return
	}

	a.mu.Lock()
	a.aiHistory = append(a.aiHistory, nvidiaMessage{Role: "user", Content: text})
	history := trimHistory(a.aiHistory, 10) // keep last 10 messages (5 turns)
	systemPrompt := a.systemPrompt + a.sessionContext()
	a.mu.Unlock()

	// Call 1 (dialogue) + Call 2 (analyzer) in parallel
	var (
		aiText   string
		analysis *AnalysisResult
		wg       sync.WaitGroup
	)
	wg.Add(2)

	go func() {
		defer wg.Done()
		if a.ai != nil {
			t := time.Now()
			var err error
			aiText, err = a.ai.Chat(r.Context(), systemPrompt, history)
			if err != nil {
				log.Printf("call1 error (%s): %v — retrying", time.Since(t).Round(time.Millisecond), err)
				aiText, err = a.ai.Chat(r.Context(), systemPrompt, history)
			}
			if err != nil || strings.TrimSpace(aiText) == "" {
				log.Printf("call1 fallback (%s)", time.Since(t).Round(time.Millisecond))
				aiText = buildReflection(text)
			} else {
				log.Printf("call1 ok (%s)", time.Since(t).Round(time.Millisecond))
			}
		} else {
			aiText = buildReflection(text)
		}
	}()

	go func() {
		defer wg.Done()
		if a.ai != nil {
			t := time.Now()
			var err error
			analysis, err = a.ai.Analyze(context.Background(), text)
			if err != nil {
				log.Printf("call2 error (%s): %v — retrying", time.Since(t).Round(time.Millisecond), err)
				analysis, err = a.ai.Analyze(context.Background(), text)
			}
			if err != nil {
				log.Printf("call2 failed (%s): %v", time.Since(t).Round(time.Millisecond), err)
			} else {
				log.Printf("call2 ok signal=%s (%s)", analysis.SignalStrength, time.Since(t).Round(time.Millisecond))
			}
		}
	}()

	wg.Wait()

	now := time.Now()
	userMessage := Message{
		ID:   fmt.Sprintf("%d-user", now.UnixNano()),
		Role: "user",
		Text: text,
		Time: formatClock(now),
	}
	aiMessage := Message{
		ID:   fmt.Sprintf("%d-ai", now.UnixNano()),
		Role: "ai",
		Text: aiText,
		Time: formatClock(now),
	}

	a.mu.Lock()
	a.aiHistory = append(a.aiHistory, nvidiaMessage{Role: "assistant", Content: aiText})
	if analysis != nil && analysis.SignalStrength != "weak" {
		a.applyScoresDelta(analysis)
		a.saveSessionData(analysis)
	}
	a.messages = append(a.messages, userMessage, aiMessage)
	a.metrics = buildMetrics(now, collectUserTexts(a.messages))
	metrics := a.metrics
	a.mu.Unlock()

	writeJSON(w, http.StatusOK, chatResponse{Message: aiMessage, Metrics: metrics, Analysis: analysis})
}

func trimHistory(history []nvidiaMessage, maxMessages int) []nvidiaMessage {
	if len(history) <= maxMessages {
		cp := make([]nvidiaMessage, len(history))
		copy(cp, history)
		return cp
	}
	cp := make([]nvidiaMessage, maxMessages)
	copy(cp, history[len(history)-maxMessages:])
	return cp
}

func (a *App) saveSessionData(analysis *AnalysisResult) {
	if analysis.Excerpt.Quote != nil && analysis.Excerpt.ContextTag != nil {
		a.sessionExcerpts = append(a.sessionExcerpts, analysis.Excerpt)
		log.Printf("session excerpt saved: [%s] %.60q", *analysis.Excerpt.ContextTag, *analysis.Excerpt.Quote)
	}
	for _, t := range analysis.Triggers {
		if t != "" {
			a.sessionTriggers[t]++
		}
	}
}

func (a *App) applyScoresDelta(analysis *AnalysisResult) {
	const chatSourceWeight = 0.5
	var signalMod float64
	switch analysis.SignalStrength {
	case "strong":
		signalMod = 0.7
	case "moderate":
		signalMod = 0.3
	default:
		return
	}
	w := chatSourceWeight * signalMod

	apply := func(key string, delta *float64) {
		if delta == nil {
			return
		}
		v := a.scores[key] + *delta*w
		if v < 0 {
			v = 0
		} else if v > 10 {
			v = 10
		}
		a.scores[key] = v
	}

	d := analysis.ScoresDelta
	apply("anxiety", d.Anxiety)
	apply("energy", d.Energy)
	apply("mood", d.Mood)
	apply("openness", d.Openness)
	apply("focus", d.Focus)
	apply("irritability", d.Irritability)
}

func (a *App) handleSuggestions(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, suggestionCards())
}

func (a *App) handleMetrics(w http.ResponseWriter, _ *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()

	writeJSON(w, http.StatusOK, a.metrics)
}

func buildReflection(input string) string {
	q := strings.ToLower(input)
	switch {
	case hasAny(q, "anxious", "anxiety", "panic", "worried", "stress", "stressed"):
		return "It sounds like your nervous system is staying on high alert. A useful next step is to name the trigger, the body signal, and the thought that followed it. This is not a diagnosis, but it can make the pattern easier to observe."
	case hasAny(q, "sad", "empty", "lonely", "tired", "hopeless"):
		return "I hear a lower emotional tone in that. Try separating what happened from what you concluded about yourself. Mirror can track this over time and show whether the feeling appears around specific situations, people, or routines."
	case hasAny(q, "angry", "irritated", "frustrated", "rage"):
		return "Anger often points to a crossed boundary, unmet need, or accumulated pressure. Write down what felt unfair, what you needed in that moment, and what action would protect you without escalating the situation."
	case hasAny(q, "procrastinate", "postpon", "avoid", "can't start", "cannot start"):
		return "Avoidance usually protects you from discomfort in the short term while increasing pressure later. The pattern to test is: what emotion appears right before you switch away from the task?"
	case hasAny(q, "happy", "good", "proud", "calm", "better"):
		return "That is worth tracking too. Notice what supported this state: sleep, movement, social contact, progress, or lower pressure. Positive patterns are as useful as difficult ones."
	default:
		return "I can reflect this as a pattern: describe the situation, the emotion, the body reaction, and what you did next. Over several entries, Mirror will make recurring triggers and coping strategies easier to see."
	}
}

func buildMetrics(now time.Time, texts []string) Metrics {
	baseValues := []int{72, 85, 78, 91, 88, 76, 84}
	if len(texts) > 0 {
		score := scoreText(texts[len(texts)-1])
		baseValues[len(baseValues)-1] = score
	}

	labels := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	xs := []int{10, 107, 203, 300, 397, 493, 590}
	chart := make([]ChartPoint, len(labels))
	total := 0
	for i := range labels {
		total += baseValues[i]
		chart[i] = ChartPoint{Label: labels[i], Value: baseValues[i], X: xs[i]}
	}

	average := float64(total) / float64(len(baseValues)) / 10
	return Metrics{
		AverageScore: round1(average),
		AverageDelta: 1.2,
		StreakDays:   max(1, 14+len(texts)),
		CurrentMonth: now.Format("January"),
		ChartData:    chart,
		CalendarDays: buildCalendar(now, texts),
	}
}

func buildCalendar(now time.Time, texts []string) []CalendarDay {
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		AddDate(0, 0, -weekday+1-21)

	moods := []string{"great", "good", "neutral", "bad", "great"}
	days := make([]CalendarDay, 28)
	for i := range days {
		d := monday.AddDate(0, 0, i)
		mood := moods[d.Day()%len(moods)]
		if len(texts) > 0 && sameDate(d, now) {
			mood = moodFromScore(scoreText(texts[len(texts)-1]))
		}
		days[i] = CalendarDay{Day: d.Day(), Mood: mood, IsFuture: d.After(now)}
	}
	return days
}

func scoreText(input string) int {
	q := strings.ToLower(input)
	score := 80
	score -= 12 * countAny(q, "anxious", "anxiety", "panic", "stress", "stressed", "worried")
	score -= 10 * countAny(q, "sad", "empty", "lonely", "hopeless", "angry", "frustrated")
	score += 8 * countAny(q, "happy", "calm", "proud", "better", "grateful", "relieved")
	if score < 50 {
		return 50
	}
	if score > 98 {
		return 98
	}
	return score
}

func moodFromScore(score int) string {
	switch {
	case score >= 85:
		return "great"
	case score >= 72:
		return "good"
	case score >= 60:
		return "neutral"
	default:
		return "bad"
	}
}

func collectUserTexts(messages []Message) []string {
	texts := make([]string, 0)
	for _, msg := range messages {
		if msg.Role == "user" {
			texts = append(texts, msg.Text)
		}
	}
	return texts
}

func hasAny(input string, needles ...string) bool {
	for _, needle := range needles {
		if strings.Contains(input, needle) {
			return true
		}
	}
	return false
}

func countAny(input string, needles ...string) int {
	count := 0
	for _, needle := range needles {
		if strings.Contains(input, needle) {
			count++
		}
	}
	return count
}

func sameDate(a, b time.Time) bool {
	ay, am, ad := a.Date()
	by, bm, bd := b.Date()
	return ay == by && am == bm && ad == bd
}

func round1(value float64) float64 {
	return float64(int(value*10+0.5)) / 10
}

func formatClock(t time.Time) string {
	return t.Format("15:04")
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(value); err != nil && !errors.Is(err, http.ErrHandlerTimeout) {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

type NvidiaClient struct {
	apiKey     string
	endpoint   string
	model      string
	httpClient *http.Client
}

type nvidiaChatRequest struct {
	Model       string          `json:"model"`
	Messages    []nvidiaMessage `json:"messages"`
	MaxTokens   int             `json:"max_tokens"`
	Temperature float64         `json:"temperature"`
	TopP        float64         `json:"top_p"`
	Stream      bool            `json:"stream"`
}

type nvidiaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type nvidiaChatResponse struct {
	Choices []struct {
		Message      nvidiaMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
}

func NewNvidiaClientFromEnv() *NvidiaClient {
	apiKey := os.Getenv("AI_API_KEY")
	if apiKey == "" {
		return nil
	}

	endpoint := os.Getenv("AI_INVOKE_URL")
	if endpoint == "" {
		endpoint = "https://integrate.api.nvidia.com/v1/chat/completions"
	}

	model := os.Getenv("AI_MODEL")
	if model == "" {
		model = "stepfun-ai/step-3.7-flash"
	}

	return &NvidiaClient{
		apiKey:   apiKey,
		endpoint: endpoint,
		model:    model,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

func (c *NvidiaClient) doRequest(ctx context.Context, payload nvidiaChatRequest) (string, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/mental-mirror")
	req.Header.Set("X-Title", "Mental Mirror")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		responseBody, _ := io.ReadAll(io.LimitReader(res.Body, 2048))
		return "", fmt.Errorf("nvidia api status %d: %s", res.StatusCode, strings.TrimSpace(string(responseBody)))
	}

	var parsed nvidiaChatResponse
	if err := json.NewDecoder(res.Body).Decode(&parsed); err != nil {
		return "", err
	}
	if len(parsed.Choices) == 0 {
		return "", errors.New("nvidia api returned no choices")
	}

	choice := parsed.Choices[0]
	if choice.FinishReason != "" && choice.FinishReason != "stop" {
		log.Printf("nvidia finish_reason=%q content_len=%d", choice.FinishReason, len(choice.Message.Content))
	}
	content := strings.TrimSpace(choice.Message.Content)
	if content == "" {
		return "", fmt.Errorf("nvidia api returned empty content (finish_reason=%q)", choice.FinishReason)
	}
	return content, nil
}

func (c *NvidiaClient) Chat(ctx context.Context, systemPrompt string, history []nvidiaMessage) (string, error) {
	messages := make([]nvidiaMessage, 0, len(history)+1)
	messages = append(messages, nvidiaMessage{Role: "system", Content: systemPrompt})
	messages = append(messages, history...)

	return c.doRequest(ctx, nvidiaChatRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   1024,
		Temperature: 0.7,
		TopP:        0.95,
	})
}

const analyzerSystem = `You are a silent background analyzer for a mental health reflection app. You receive a single message from a user and return a structured JSON analysis. You are never shown to the user. Be precise and conservative — only flag what is clearly present in the text.`

const analyzerUserTemplate = `Analyze the following user message and return ONLY a valid JSON object with no explanation, no markdown, no code blocks.

User message:
"__MSG__"

Return this exact structure:
{
  "scores_delta": {
    "anxiety":      <float -2.0 to 2.0 or null if not detectable>,
    "energy":       <float -2.0 to 2.0 or null if not detectable>,
    "mood":         <float -2.0 to 2.0 or null if not detectable>,
    "openness":     <float -2.0 to 2.0 or null if not detectable>,
    "focus":        <float -2.0 to 2.0 or null if not detectable>,
    "irritability": <float -2.0 to 2.0 or null if not detectable>
  },
  "excerpt": {
    "quote": <string: verbatim phrase from the message worth remembering, or null>,
    "context_tag": <string: one snake_case tag describing the pattern, or null>
  },
  "triggers": [<string>, ...],
  "signal_strength": <"strong" | "moderate" | "weak">
}

Rules:
- scores_delta: positive = increase, negative = decrease. Use null when there is no clear signal.
- excerpt.quote: only if emotionally significant or reveals a recurring pattern. Otherwise null.
- triggers: short noun phrases, max 3, empty array if none.
- signal_strength: "strong" if multiple clear signals, "moderate" if one clear signal, "weak" if short or neutral.`

func (c *NvidiaClient) Analyze(ctx context.Context, userMessage string) (*AnalysisResult, error) {
	userPrompt := strings.ReplaceAll(analyzerUserTemplate, "__MSG__", userMessage)

	raw, err := c.doRequest(ctx, nvidiaChatRequest{
		Model: c.model,
		Messages: []nvidiaMessage{
			{Role: "system", Content: analyzerSystem},
			{Role: "user", Content: userPrompt},
		},
		MaxTokens:   600,
		Temperature: 0.2,
		TopP:        0.9,
	})
	if err != nil {
		return nil, err
	}

	log.Printf("analyze raw (pre-strip, truncated 512): %.512s", raw)

	// Strip markdown fences: ```json ... ``` or ``` ... ```
	raw = strings.TrimSpace(raw)
	if strings.HasPrefix(raw, "```") {
		// remove first line (``` or ```json)
		if i := strings.Index(raw, "\n"); i != -1 {
			raw = raw[i+1:]
		}
		// remove trailing ```
		if strings.HasSuffix(strings.TrimSpace(raw), "```") {
			raw = raw[:strings.LastIndex(raw, "```")]
		}
		raw = strings.TrimSpace(raw)
	}

	var result AnalysisResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		log.Printf("analyze post-strip parse error, raw: %.512s", raw)
		return nil, fmt.Errorf("parse analysis JSON: %w", err)
	}
	return &result, nil
}
