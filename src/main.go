package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"slackbot-go/infra/di"
	"slackbot-go/slackbot"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/slack-go/slack/slackevents"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("req.Body===========%+v", req.Body)
	// challenge認証用
	parsedReq, err := slackbot.ParseRequest(req.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	if parsedReq.Challenge != "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 200,
			Body:       parsedReq.Challenge,
		}, nil
	}

	eventsAPIEvent, err := slackevents.ParseEvent(
		json.RawMessage(req.Body),
		slackevents.OptionVerifyToken(
			&slackevents.TokenComparator{
				VerificationToken: os.Getenv("SLACK_VERIFICATION_TOKEN"),
			},
		),
	)
	// パースもしくは Token の認証がうまく行かなかった場合のエラー処理
	if err != nil {
		log.Printf("ERROR: %s", err)
		return makeResponse(http.StatusInternalServerError, err.Error())
	}

	message := slackbot.AI(di.NewAI(), req)

	return makeResponse(http.StatusOK, message), nil
}

func makeResponse(code int, body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       body,
	}
}
