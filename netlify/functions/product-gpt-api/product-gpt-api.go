package main

import (
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/kunal-mandalia/product-gpt-api/chatgpt"
)

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	apiKey := os.Getenv("CHATGPT_API_KEY")
	if apiKey == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Server error",
		}, nil
	}

	q := request.QueryStringParameters["q"]
	if q == "" {
		return &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Missing query",
		}, nil
	}

	res, err := chatgpt.GetTextCompletion(apiKey, q)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       res,
	}, nil
}

func main() {
	godotenv.Load()
	lambda.Start(handler)
}
