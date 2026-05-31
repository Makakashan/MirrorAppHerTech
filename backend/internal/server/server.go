package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type App struct {
	mu       sync.Mutex
	messages []Message
	metrics  Metrics
	ai       *NvidiaClient
}

func New() *App {
	now := time.Now()
	return &App{
		messages: initialMessages(now),
		metrics:  initialMetrics(now),
		ai:       NewNvidiaClientFromEnv(),
	}
}

func (a *App) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", a.handleHealth)
	mux.HandleFunc("GET /api/messages", a.handleMessages)
	mux.HandleFunc("POST /api/chat", a.handleChat)
	mux.HandleFunc("GET /api/suggestions", a.handleSuggestions)
	mux.HandleFunc("GET /api/metrics", a.handleMetrics)
	return cors(mux)
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
		Text: a.generateReflection(r.Context(), text),
		Time: formatClock(now),
	}

	a.mu.Lock()
	a.messages = append(a.messages, userMessage, aiMessage)
	a.metrics = buildMetrics(now, collectUserTexts(a.messages))
	metrics := a.metrics
	a.mu.Unlock()

	writeJSON(w, http.StatusOK, chatResponse{Message: aiMessage, Metrics: metrics})
}

func (a *App) generateReflection(ctx context.Context, text string) string {
	if a.ai == nil {
		return buildReflection(text)
	}

	response, err := a.ai.Complete(ctx, text)
	if err != nil || strings.TrimSpace(response) == "" {
		return buildReflection(text)
	}

	return response
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
		Message nvidiaMessage `json:"message"`
	} `json:"choices"`
}

func NewNvidiaClientFromEnv() *NvidiaClient {
	apiKey := os.Getenv("NVIDIA_API_KEY")
	if apiKey == "" {
		return nil
	}

	endpoint := os.Getenv("NVIDIA_INVOKE_URL")
	if endpoint == "" {
		endpoint = "https://integrate.api.nvidia.com/v1/chat/completions"
	}

	model := os.Getenv("NVIDIA_MODEL")
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

func (c *NvidiaClient) Complete(ctx context.Context, input string) (string, error) {
	payload := nvidiaChatRequest{
		Model: c.model,
		Messages: []nvidiaMessage{
			{
				Role:    "system",
				Content: "You are Mirror, a concise reflection assistant. Help users notice emotions, context, body signals, patterns, and next reflective questions. Do not diagnose, label the user, or present yourself as a therapist. Keep replies supportive, practical, and under 120 words.",
			},
			{
				Role:    "user",
				Content: input,
			},
		},
		MaxTokens:   512,
		Temperature: 0.7,
		TopP:        0.95,
		Stream:      false,
	}

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

	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}
