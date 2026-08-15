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

func handler(ctx context.Context, req events.APIGatewayV2HTTPRequest) (events.APIGatewayV2HTTPResponse, error) {
	mux := router.SetupMux()
	adapter := chiadapter.NewV2(mux)
	return adapter.ProxyWithContextV2(ctx, req)
}
