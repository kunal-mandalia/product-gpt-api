package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
	"github.com/kunal-mandalia/product-gpt-api/chatgpt"
	"github.com/kunal-mandalia/product-gpt-api/search"
)

// request structs
type ProductRecommendationsBody struct {
	Query_request  string `json:"query_request"`
	Query_response string `json:"query_response"`
}

func requiredValue(s string, fromUser bool) (string, *events.APIGatewayProxyResponse) {
	if s == "" {
		statusCode := 0
		body := ""
		if fromUser {
			statusCode = 400
			body = "Bad request (incomplete, missing value)"
		} else {
			statusCode = 500
			body = "Server error (config)"
		}
		fmt.Println(statusCode, body)
		return "", &events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       body,
		}
	}
	return s, nil
}

func handleUpstreamResponse(res interface{}, err error) (*events.APIGatewayProxyResponse, error) {
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Server error (upstream)",
		}, nil
	}

	// send string
	b, err := json.Marshal(res)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Marshal error",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(b),
	}, nil
}

// expose textcompletion and search endpoints
func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println(request)

	// check api keys are loaded
	chatGPTApiKey, e := requiredValue(os.Getenv("CHATGPT_API_KEY"), false)
	if e != nil {
		return e, nil
	}
	searchApiKey, e := requiredValue(os.Getenv("CUSTOMSEARCH_API_KEY"), false)
	if e != nil {
		return e, nil
	}

	if strings.Contains(request.Path, "/textcompletion") {
		q, e := requiredValue(request.QueryStringParameters["q"], true)
		if e != nil {
			return e, nil
		}
		res, err := chatgpt.TextCompletion(chatGPTApiKey, q)
		return handleUpstreamResponse(res, err)
	}

	if strings.Contains(request.Path, "/product_recommendations") {
		args := ProductRecommendationsBody{}
		json.Unmarshal([]byte(request.Body), &args)
		qReq, e := requiredValue(args.Query_request, true)
		if e != nil {
			return e, nil
		}
		qRes, e := requiredValue(args.Query_response, true)
		if e != nil {
			return e, nil
		}
		query := chatgpt.ProductRecommendationsQuery(qRes, qReq)
		res, err := chatgpt.TextCompletion(chatGPTApiKey, query)
		return handleUpstreamResponse(res, err)
	}

	if strings.Contains(request.Path, "/search") {
		q, e := requiredValue(request.QueryStringParameters["q"], true)
		if e != nil {
			return e, nil
		}
		res, err := search.ProductList(searchApiKey, q)
		return handleUpstreamResponse(res, err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Bad query",
	}, nil
}

func main() {
	godotenv.Load()
	lambda.Start(handler)
}
