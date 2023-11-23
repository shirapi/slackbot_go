package slackbot

import (
	"context"
	"log"
	"os"
	"slackbot-go/domain"

	"github.com/slack-go/slack"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/prompts"
)

var (
	slackClient *slack.Client
)

func init() {
	slackClient = slack.New(os.Getenv("SLACK_BOT_OAUTH_TOKEN"))
}

func AI(ai domain.Charactor, req Request) string {
	llm, err := openai.New()
	if err != nil {
		log.Fatal(err)
	}

	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
		prompts.NewSystemMessagePromptTemplate(
			ai.GetSystemPrompt(),
			[]string{ai.GetHistoryLabel()},
		),
		prompts.NewHumanMessagePromptTemplate(
			ai.GetHumanPrompt(),
			[]string{ai.GetInputLabel()},
		),
	})

	result, err := prompt.Format(map[string]any{
		ai.GetHistoryLabel(): "",
		ai.GetInputLabel():   "",
	})
	if err != nil {
		log.Fatal(err)
	}

	completion, err := llm.Call(context.Background(), result)
	if err != nil {
		log.Fatal(err)
	}
	return completion
}
