package di

import (
	"os"
	"slackbot-go/infra/aimodel"
	"slackbot-go/infra/charactor"
	"slackbot-go/slackbot"

	"github.com/slack-go/slack"
)

func NewHiyokoTalker() *slackbot.Talker {
	return &slackbot.Talker{
		Client: slack.New(os.Getenv("SLACK_BOT_OAUTH_TOKEN")),
		Model:  aimodel.NewOpenAI(),
		AI:     charactor.NewHiyoko(),
	}
}
