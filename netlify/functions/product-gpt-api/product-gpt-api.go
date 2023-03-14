package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/kunal-mandalia/product-gpt-api/openai"
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
	allowedOrigin, e := requiredValue(os.Getenv("ALLOWED_ORIGIN"), false)
	fmt.Println("allowed origin", allowedOrigin)
	if e != nil {
		sentry.CaptureException(errors.New("missing env allowed origin"))
		return e, nil
	}

	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = allowedOrigin
	headers["Access-Control-Allow-Headers"] = "Content-Type"
	headers["Access-Control-Allow-Methods"] = "GET, POST"

	if err != nil {
		sentry.CaptureException(err)
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       "Server error (upstream)",
		}, nil
	}

	b, err := json.Marshal(res)
	if err != nil {
		sentry.CaptureException(err)
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
	openAIApiKey, e := requiredValue(os.Getenv("OPENAI_API_KEY"), false)
	if e != nil {
		return e, nil
	}
	searchApiKey, e := requiredValue(os.Getenv("CUSTOMSEARCH_API_KEY"), false)
	if e != nil {
		return e, nil
	}
	ebayCampaignId, e := requiredValue(os.Getenv("EBAY_CAMPAIGN_ID"), false)
	if e != nil {
		return e, nil
	}

	if strings.Contains(request.Path, "/textcompletion") {
		sentry.CaptureMessage("api_hit: /textcompletion/q=" + request.QueryStringParameters["q"])
		q, e := requiredValue(request.QueryStringParameters["q"], true)
		if e != nil {
			return e, nil
		}
		res, err := openai.TextCompletion(openAIApiKey, q)
		if err != nil {
			sentry.CaptureException(err)
		}
		return handleUpstreamResponse(res, err)
	}

	if strings.Contains(request.Path, "/product_recommendation") {
		sentry.CaptureMessage("api_hit: /product_recommendation/")
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
		query := openai.ProductRecommendationsQuery(qReq, qRes)
		res, err := openai.TextCompletion(openAIApiKey, query)
		if err != nil {
			sentry.CaptureException(err)
		}
		return handleUpstreamResponse(res, err)
	}

	if strings.Contains(request.Path, "/entities") {
		sentry.CaptureMessage("api_hit: /entities")
		args := ProductRecommendationsBody{}
		json.Unmarshal([]byte(request.Body), &args)
		qReq, e := requiredValue(args.Query_request, true)
		if e != nil {
			return e, nil
		}
		query := openai.EntityExtractionQuery(qReq)
		res, err := openai.TextCompletion(openAIApiKey, query)
		if err != nil {
			sentry.CaptureException(err)
		}
		return handleUpstreamResponse(res, err)
	}

	if strings.Contains(request.Path, "/search") {
		sentry.CaptureMessage("api_hit: /search?q" + request.QueryStringParameters["q"])
		q, e := requiredValue(request.QueryStringParameters["q"], true)
		if e != nil {
			return e, nil
		}
		res, err := search.ProductList(searchApiKey, q)
		if err != nil {
			sentry.CaptureException(err)
		}
		return handleUpstreamResponse(res, err)
	}

	if strings.Contains(request.Path, "/ebay_search") {
		sentry.CaptureMessage("api_hit: /ebay_search?q=" + request.QueryStringParameters["q"])
		q, e := requiredValue(request.QueryStringParameters["q"], true)
		if e != nil {
			return e, nil
		}
		marketPlace, e := requiredValue(request.QueryStringParameters["marketplace"], true)
		if e != nil {
			return e, nil
		}
		limit, e := requiredValue(request.QueryStringParameters["limit"], true)
		if e != nil {
			return e, nil
		}
		nLimit, _ := strconv.Atoi(limit)

		// TODO: cache token
		t, _ := search.EbayGetAccessToken()
		res, err := search.EbaySearch(q, marketPlace, nLimit, t.AccessToken, ebayCampaignId)
		if err != nil {
			sentry.CaptureException(err)
		}
		return handleUpstreamResponse(res, err)
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 400,
		Body:       "Bad query",
	}, nil
}

func main() {
	godotenv.Load()
	sentryDsn, e := requiredValue(os.Getenv("SENTRY_DSN"), false)
	if e != nil {
		log.Fatalf("missing sentry dsn")
	}
	environment, e := requiredValue(os.Getenv("ENVIRONMENT"), false)
	if e != nil {
		log.Fatalf("missing environment")
	}
	sErr := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDsn,
		Environment: environment,
	})
	if sErr != nil {
		log.Fatalf("sentry.Init: %s", sErr)
	}
	lambda.Start(handler)
}
