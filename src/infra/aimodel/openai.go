package aimodel

import (
	"context"
	"fmt"
	"slackbot-go/domain"
	"strings"

	"github.com/tmc/langchaingo/chains"
	"github.com/tmc/langchaingo/llms/openai"
	"github.com/tmc/langchaingo/memory"
	"github.com/tmc/langchaingo/outputparser"
	"github.com/tmc/langchaingo/prompts"
)

type OpenAI struct {
	model string
}

func NewOpenAI(model string) *OpenAI {
	return &OpenAI{
		model: model,
	}
}

func (a *OpenAI) Call(ai domain.Charactor, histories []domain.History, userMessage string) (string, error) {

	llm, err := openai.New()
	if err != nil {
		return "", err
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

	ctx := context.Background()
	chatHistory := a.convertHistory(ctx, histories)
	chainMemory := memory.NewConversationBuffer(
		memory.WithChatHistory(chatHistory),
		memory.WithReturnMessages(true),
		memory.WithHumanPrefix("H"),
		memory.WithAIPrefix("A"),
		memory.WithMemoryKey(ai.GetHistoryLabel()),
		memory.WithInputKey(ai.GetInputLabel()),
	)

	message, err := chains.Predict(
		ctx,
		chains.LLMChain{
			Prompt:       prompt,
			LLM:          llm,
			OutputParser: outputparser.NewSimple(),
			Memory:       chainMemory,
			// Memory:       memory.NewSimple(),
		},
		map[string]any{
			ai.GetInputLabel():   userMessage,
			ai.GetHistoryLabel(): chatHistory,
		},
		chains.WithModel(a.model),
		chains.WithTemperature(0.6),
	)
	if err != nil {
		fmt.Println("chain error:", err)
		return "", err
	}

	// 返信に名前が付加されていたら削除
	// 先頭に fmt.Sprintf("%s: ", ai.GetName()) の結果が存在したら空白に置換する
	reply := strings.Replace(message, fmt.Sprintf("%s:", ai.GetName()), "", 1)
	reply = strings.Replace(reply, fmt.Sprintf("%s：", ai.GetName()), "", 1)
	return reply, nil
}

func (a *OpenAI) convertHistory(ctx context.Context, histories []domain.History) *memory.ChatMessageHistory {
	history := memory.NewChatMessageHistory()
	for _, h := range histories {
		if h.MessageFrom == domain.MessageFromUser {
			history.AddUserMessage(ctx, h.Message)
		} else {
			history.AddAIMessage(ctx, h.Message)
		}
	}
	return history
}
