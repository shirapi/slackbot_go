package di

import (
	"os"
	"slackbot-go/infra/aimodel"
	"slackbot-go/infra/charactor"
	"slackbot-go/slackbot"

	"github.com/slack-go/slack"
)

func NewHiyokoTalker() (*slackbot.Talker, error) {
	// TODO 環境変数のチェック
	return &slackbot.Talker{
		Client: slack.New(os.Getenv("SLACK_BOT_OAUTH_TOKEN")),
		Model:  aimodel.NewOpenAI(os.Getenv("AI_MODEL")),
		AI:     charactor.NewHiyoko(),
	}, nil
}
