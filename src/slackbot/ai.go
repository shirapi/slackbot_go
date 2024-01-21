package slackbot

// func ai(ai domain.Charactor) (string, error) {

// 	llm, err := openai.New()
// 	if err != nil {
// 		return "", err
// 	}

// 	prompt := prompts.NewChatPromptTemplate([]prompts.MessageFormatter{
// 		prompts.NewSystemMessagePromptTemplate(
// 			ai.GetSystemPrompt(),
// 			[]string{ai.GetHistoryLabel()},
// 		),
// 		prompts.NewHumanMessagePromptTemplate(
// 			ai.GetHumanPrompt(),
// 			[]string{ai.GetInputLabel()},
// 		),
// 	})

// 	result, err := prompt.Format(map[string]any{
// 		ai.GetHistoryLabel(): "",
// 		ai.GetInputLabel():   "",
// 	})
// 	if err != nil {
// 		return "", err
// 	}

// 	completion, err := llm.Call(context.Background(), result)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	return completion, nil
// }
