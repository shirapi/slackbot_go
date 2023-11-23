package di

import (
	"slackbot-go/domain"
	"slackbot-go/infra/charactor"
)

func NewAI() domain.Charactor {
	return charactor.NewHiyoko()
}
