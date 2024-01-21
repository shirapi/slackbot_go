package domain

type AIModel interface {
	Call(ai Charactor, histories []History, userMessage string) (string, error)
}
