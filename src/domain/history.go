package domain

const (
	MessageFromUser = iota
	MessageFromAI
)

type History struct {
	MessageFrom int
	Message     string
}
