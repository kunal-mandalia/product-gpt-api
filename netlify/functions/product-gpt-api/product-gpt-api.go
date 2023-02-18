package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/kunal-mandalia/product-gpt-api/chatgpt"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "Hello, World!",
	}, nil
}

func main() {
	chatgpt.GetTextCompletion("foo")
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}
