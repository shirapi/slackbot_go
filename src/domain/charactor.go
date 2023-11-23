package domain

type Charactor interface {
	GetName() string
	GetSystemPrompt() string
	GetHumanPrompt() string
	GetHistoryLabel() string
	GetInputLabel() string
}
