package server

type Message struct {
	ID   string `json:"id"`
	Role string `json:"role"`
	Text string `json:"text"`
	Time string `json:"time"`
}

type SuggestionCard struct {
	Icon          string `json:"icon"`
	TextMain      string `json:"textMain"`
	TextSecondary string `json:"textSecondary"`
}

type ChartPoint struct {
	Label string `json:"label"`
	Value int    `json:"value"`
	X     int    `json:"x"`
}

type CalendarDay struct {
	Day      int    `json:"day"`
	Mood     string `json:"mood"`
	IsFuture bool   `json:"isFuture"`
}

type Metrics struct {
	AverageScore float64       `json:"averageScore"`
	AverageDelta float64       `json:"averageDelta"`
	StreakDays   int           `json:"streakDays"`
	CurrentMonth string        `json:"currentMonth"`
	ChartData    []ChartPoint  `json:"chartData"`
	CalendarDays []CalendarDay `json:"calendarDays"`
}

type chatRequest struct {
	Message string `json:"message"`
}

type chatResponse struct {
	Message  Message         `json:"message"`
	Metrics  Metrics         `json:"metrics"`
	Analysis *AnalysisResult `json:"analysis,omitempty"`
}

type ScoresDelta struct {
	Anxiety      *float64 `json:"anxiety"`
	Energy       *float64 `json:"energy"`
	Mood         *float64 `json:"mood"`
	Openness     *float64 `json:"openness"`
	Focus        *float64 `json:"focus"`
	Irritability *float64 `json:"irritability"`
}

type AnalysisExcerpt struct {
	Quote      *string `json:"quote"`
	ContextTag *string `json:"context_tag"`
}

type AnalysisResult struct {
	ScoresDelta   ScoresDelta     `json:"scores_delta"`
	Excerpt       AnalysisExcerpt `json:"excerpt"`
	Triggers      []string        `json:"triggers"`
	SignalStrength string          `json:"signal_strength"`
}
