package domain

const (
	TalkChallenge = iota
	TalkNormal
)

type Response struct {
	Challenge bool
	Message   string `json:"message"`
}
