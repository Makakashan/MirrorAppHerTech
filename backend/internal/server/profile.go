package server

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func dataDir() string {
	if d := os.Getenv("DATA_DIR"); d != "" {
		return d
	}
	return "data"
}

// compactProfile is what gets serialized into the system prompt.
// Keeps only what the agent needs; omits verbose fields that bloat the context.
type compactProfile struct {
	Name              string             `json:"name"`
	Occupation        string             `json:"occupation,omitempty"`
	CurrentScores     map[string]float64 `json:"current_scores"`
	LifeCircumstances []any              `json:"life_circumstances,omitempty"`
	RecentSnapshots   []any              `json:"recent_snapshots,omitempty"`
	RecentExcerpts    []any              `json:"recent_excerpts,omitempty"`
	Patterns          []any              `json:"patterns,omitempty"`
}

func loadProfileAndBuildPrompt() (string, map[string]float64) {
	dir := dataDir()

	seedsData, err := os.ReadFile(fmt.Sprintf("%s/seeds.json", dir))
	if err != nil {
		return defaultSystemPrompt(), defaultScores()
	}

	var seeds struct {
		Users []json.RawMessage `json:"users"`
	}
	if err := json.Unmarshal(seedsData, &seeds); err != nil || len(seeds.Users) == 0 {
		return defaultSystemPrompt(), defaultScores()
	}

	var user map[string]any
	if err := json.Unmarshal(seeds.Users[0], &user); err != nil {
		return defaultSystemPrompt(), defaultScores()
	}

	scores := extractScores(user)
	compact := buildCompactProfile(user, scores)

	profileJSON, err := json.MarshalIndent(compact, "", "  ")
	if err != nil {
		return defaultSystemPrompt(), scores
	}

	promptData, err := os.ReadFile(fmt.Sprintf("%s/system_prompt.md", dir))
	if err != nil {
		return defaultSystemPrompt(), scores
	}

	prompt := strings.ReplaceAll(string(promptData), "[USER_PROFILE_JSON]", string(profileJSON))
	return prompt, scores
}

func buildCompactProfile(user map[string]any, scores map[string]float64) compactProfile {
	p := compactProfile{
		CurrentScores: scores,
	}

	if v, ok := user["name"].(string); ok {
		p.Name = v
	}

	if ctx, ok := user["context"].(map[string]any); ok {
		if v, ok := ctx["occupation"].(string); ok {
			p.Occupation = v
		}
		p.LifeCircumstances = takeLast(sliceOf(ctx["life_circumstances"]), 4)
	}

	p.RecentSnapshots = takeLast(sliceOf(user["snapshots"]), 3)
	p.RecentExcerpts = takeLast(sliceOf(user["chat_excerpts"]), 6)
	p.Patterns = sliceOf(user["patterns"])

	return p
}

func sliceOf(v any) []any {
	if s, ok := v.([]any); ok {
		return s
	}
	return nil
}

func takeLast(s []any, n int) []any {
	if len(s) <= n {
		return s
	}
	return s[len(s)-n:]
}

func extractScores(user map[string]any) map[string]float64 {
	scores := defaultScores()
	raw, ok := user["current_scores"]
	if !ok {
		return scores
	}
	m, ok := raw.(map[string]any)
	if !ok {
		return scores
	}
	for _, key := range []string{"anxiety", "energy", "mood", "openness", "focus", "irritability"} {
		if v, ok := m[key]; ok {
			if n, ok := v.(float64); ok {
				scores[key] = n
			}
		}
	}
	return scores
}

func defaultScores() map[string]float64 {
	return map[string]float64{
		"anxiety":      5.0,
		"energy":       5.0,
		"mood":         5.0,
		"openness":     5.0,
		"focus":        5.0,
		"irritability": 5.0,
	}
}

func defaultSystemPrompt() string {
	return "You are Mirror, a reflective conversational agent. Help users see themselves more clearly through their own words and patterns. Do not advise, diagnose, or solve. Ask one question per message. Always respond in the same language the user writes in. Your only domain is the user's inner life — emotions, patterns, energy, relationships, self-perception. Do not answer questions about the external world (geography, science, technology, current events, or any general knowledge). If a request is off-topic, redirect warmly: acknowledge briefly and bring the focus back to the person. If an off-topic question appears in an emotionally tense context (e.g. asking about heights, bridges, or medications while distressed), do not answer the factual question — slow down and check in on the person instead. If there is any sign the user may be in danger, provide the crisis line immediately and step out of the reflective role."
}
