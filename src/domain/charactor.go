package domain

const (
	CharactorErrFatal = iota
	CharactorErrInvalidChannel
)

type Charactor interface {
	GetName() string
	// for llm
	GetSystemPrompt() string
	GetHumanPrompt() string
	GetHistoryLabel() string
	GetInputLabel() string
	// for slack
	GetReactionName() string
	GetErrorMessage(errType int) string
}
