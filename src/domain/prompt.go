package domain

type Prompt struct {
	HistoryLabel string
	InputLabel   string
}

var DefaultPrompt = Prompt{
	HistoryLabel: "history",
	InputLabel:   "input",
}
