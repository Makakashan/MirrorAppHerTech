package server

import (
	"os"
	"strings"
	"time"
)

const demoSuggestionIcon = "https://www.figma.com/api/mcp/asset/929399b3-69ce-47ea-ac8a-5ca9f3d4a8bf"

func demoDataEnabled() bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv("MIRROR_DEMO_DATA")))
	return value != "false" && value != "0" && value != "off"
}

func initialMessages(now time.Time) []Message {
	messages := []Message{
		{
			ID:   "welcome",
			Role: "ai",
			Text: "Hi! I'm Mirror. Tell me what has been on your mind, and I will help you reflect on the patterns behind it.",
			Time: formatClock(now),
		},
	}

	if !demoDataEnabled() {
		return messages
	}

	return append(messages,
		Message{
			ID:   "demo-user-1",
			Role: "user",
			Text: "I felt anxious before presenting my idea and then avoided replying to messages.",
			Time: formatClock(now.Add(-46 * time.Minute)),
		},
		Message{
			ID:   "demo-ai-1",
			Role: "ai",
			Text: "That looks like a pressure-and-avoidance loop: public evaluation triggered anxiety, then avoidance gave short relief. Track what you feared people would notice.",
			Time: formatClock(now.Add(-45 * time.Minute)),
		},
		Message{
			ID:   "demo-user-2",
			Role: "user",
			Text: "After a walk I felt calmer and could finish the task.",
			Time: formatClock(now.Add(-18 * time.Minute)),
		},
		Message{
			ID:   "demo-ai-2",
			Role: "ai",
			Text: "That is a useful stabilizing pattern. Movement may be helping you lower activation enough to return to focused work.",
			Time: formatClock(now.Add(-17 * time.Minute)),
		},
	)
}

func initialMetrics(now time.Time) Metrics {
	if !demoDataEnabled() {
		return buildMetrics(now, nil)
	}

	return Metrics{
		AverageScore: 8.1,
		AverageDelta: 0.9,
		StreakDays:   12,
		CurrentMonth: now.Format("January"),
		ChartData: []ChartPoint{
			{Label: "Mon", Value: 68, X: 10},
			{Label: "Tue", Value: 74, X: 107},
			{Label: "Wed", Value: 71, X: 203},
			{Label: "Thu", Value: 86, X: 300},
			{Label: "Fri", Value: 82, X: 397},
			{Label: "Sat", Value: 90, X: 493},
			{Label: "Sun", Value: 78, X: 590},
		},
		CalendarDays: demoCalendar(now),
	}
}

func suggestionCards() []SuggestionCard {
	if !demoDataEnabled() {
		return []SuggestionCard{}
	}

	return []SuggestionCard{
		{
			Icon:          demoSuggestionIcon,
			TextMain:      "I felt anxious today and could not focus",
			TextSecondary: "help me understand what may have triggered it",
		},
		{
			Icon:          demoSuggestionIcon,
			TextMain:      "I keep postponing important things",
			TextSecondary: "show me the pattern behind this behavior",
		},
		{
			Icon:          demoSuggestionIcon,
			TextMain:      "I felt calmer after changing my routine",
			TextSecondary: "show what habits may support this state",
		},
	}
}

func demoCalendar(now time.Time) []CalendarDay {
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).
		AddDate(0, 0, -weekday+1-21)

	moods := []string{
		"neutral", "good", "bad", "neutral", "good", "great", "good",
		"neutral", "bad", "good", "good", "great", "great", "neutral",
		"good", "neutral", "bad", "good", "great", "good", "neutral",
		"bad", "neutral", "good", "good", "great", "good", "neutral",
	}

	days := make([]CalendarDay, 28)
	for i := range days {
		d := monday.AddDate(0, 0, i)
		days[i] = CalendarDay{Day: d.Day(), Mood: moods[i], IsFuture: d.After(now)}
	}
	return days
}
