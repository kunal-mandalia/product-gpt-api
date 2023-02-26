package main

import (
	"encoding/json"
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
		return "", &events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       body,
		}
	}
	return s, nil
}

func handleUpstreamResponse(res interface{}, err error) (*events.APIGatewayProxyResponse, error) {
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Headers"] = "Content-Type"
	headers["Access-Control-Allow-Methods"] = "GET, POST, PUT, DELETE"

	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       "Server error (upstream)",
		}, nil
	}

	// send string
	b, err := json.Marshal(res)
	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       "Marshal error",
		}, nil
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers,
		Body:       string(b),
	}, nil
}

// expose textcompletion and search endpoints
func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
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

	if strings.Contains(request.Path, "/product_recommendation") {
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
		query := chatgpt.ProductRecommendationsQuery(qReq, qRes)
		res, err := chatgpt.TextCompletion(chatGPTApiKey, query)
		return handleUpstreamResponse(res, err)
	}

	if strings.Contains(request.Path, "/entities") {
		args := ProductRecommendationsBody{}
		json.Unmarshal([]byte(request.Body), &args)
		qReq, e := requiredValue(args.Query_request, true)
		if e != nil {
			return e, nil
		}
		query := chatgpt.EntityExtractionQuery(qReq)
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

	if strings.Contains(request.Path, "/ebay_search") {
		q, e := requiredValue(request.QueryStringParameters["q"], true)
		if e != nil {
			return e, nil
		}
		m, e := requiredValue(request.QueryStringParameters["marketplace"], true)
		if e != nil {
			return e, nil
		}
		// TODO: cache token
		t, _ := search.EbayGetAccessToken()
		res, err := search.EbaySearch(q, m, t.AccessToken)
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
