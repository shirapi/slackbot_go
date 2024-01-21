package main

import (
	"context"
	"slackbot-go/infra/router"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	chiadapter "github.com/awslabs/aws-lambda-go-api-proxy/chi"
)

func main() {
	lambda.Start(handler)
}

func handler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	mux := router.SetupMux()
	adapter := chiadapter.New(mux)
	return adapter.ProxyWithContext(ctx, req)
	// https://github.com/awslabs/aws-lambda-go-api-proxy/blob/master
	// var httpHandler http.Handler = http.HandlerFunc(slackbot.Talk)
	// adapter := httpadapter.New(httpHandler)
	// return adapter.ProxyWithContext(ctx, req)
}
