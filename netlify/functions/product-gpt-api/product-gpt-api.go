package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/getsentry/sentry-go"
	"github.com/joho/godotenv"
	"github.com/kunal-mandalia/product-gpt-api/openai"
	"github.com/kunal-mandalia/product-gpt-api/search"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type ProductRecommendationsBody struct {
	Query_request  string `json:"query_request"`
	Query_response string `json:"query_response"`
}
type CachedResponse struct {
	Key   string
	Value string
}

func makeCorsHeaders() map[string]string {
	headers := make(map[string]string)
	headers["Access-Control-Allow-Origin"] = "*"
	headers["Access-Control-Allow-Headers"] = "*"
	headers["Access-Control-Allow-Methods"] = "GET, POST, OPTIONS"
	return headers
}

func requiredValue(s string, fromUser bool) (string, *events.APIGatewayProxyResponse) {
	headers := makeCorsHeaders()

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
			Headers:    headers,
			Body:       body,
		}
	}
	return s, nil
}

func handleUpstreamResponse(res interface{}, err error, cachedRes *CachedResponse) (*events.APIGatewayProxyResponse, error) {
	headers := makeCorsHeaders()

	if err != nil {
		sentry.CaptureException(err)
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Headers:    headers,
			Body:       "Server error (upstream)",
		}, nil
	}

	if cachedRes != nil {
		sentry.CaptureMessage("cache_hit/key=" + cachedRes.Key)
		return &events.APIGatewayProxyResponse{
			StatusCode: 200,
			Headers:    headers,
			Body:       cachedRes.Value,
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
func wrappedHandler(rdb *redis.Client, cacheDuration time.Duration, openAIApiKey string, ebayCampaignId string) func(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
		fmt.Println("handler executed")
		if strings.Contains(request.Path, "/textcompletion") {
			sentry.CaptureMessage("api_hit: /textcompletion/q=" + request.QueryStringParameters["q"])
			cacheKey := "/textcompletion/q=" + request.QueryStringParameters["q"]
			cachedCommand := rdb.Get(ctx, cacheKey)
			cachedValue := cachedCommand.Val()
			if cachedValue != "" {
				c := CachedResponse{cacheKey, cachedValue}
				return handleUpstreamResponse(nil, nil, &c)
			}

			q, e := requiredValue(request.QueryStringParameters["q"], true)
			if e != nil {
				return handleUpstreamResponse(nil, errors.New("missing query"), nil)
			}
			res, err := openai.TextCompletion(openAIApiKey, q)
			if err != nil {
				sentry.CaptureException(err)
				return handleUpstreamResponse(nil, err, nil)
			}
			b, err := json.Marshal(res)
			if err != nil {
				sentry.CaptureException(err)
				return handleUpstreamResponse(nil, err, nil)
			}
			rdb.Set(ctx, cacheKey, string(b), cacheDuration)
			return handleUpstreamResponse(res, err, nil)
		}

		if strings.Contains(request.Path, "/entities") {
			sentry.CaptureMessage("api_hit: /entities")
			args := ProductRecommendationsBody{}
			json.Unmarshal([]byte(request.Body), &args)
			qReq, e := requiredValue(args.Query_request, true)
			if e != nil {
				return handleUpstreamResponse(nil, errors.New("missing entities query request"), nil)
			}

			cacheKey := "/entities?body_query_request=" + qReq
			cachedCommand := rdb.Get(ctx, cacheKey)
			cachedValue := cachedCommand.Val()
			if cachedValue != "" {
				c := CachedResponse{cacheKey, cachedValue}
				return handleUpstreamResponse(openai.TextCompletionResponse{}, nil, &c)
			}

			query := openai.EntityExtractionQuery(qReq)
			res, err := openai.TextCompletion(openAIApiKey, query)
			if err != nil {
				sentry.CaptureException(err)
				return handleUpstreamResponse(nil, err, nil)
			}
			b, err := json.Marshal(res)
			if err != nil {
				sentry.CaptureException(err)
				return handleUpstreamResponse(nil, err, nil)
			}
			rdb.Set(ctx, cacheKey, string(b), cacheDuration)
			return handleUpstreamResponse(res, err, nil)
		}

		if strings.Contains(request.Path, "/ebay_search") {
			sentry.CaptureMessage("api_hit: /ebay_search?q=" + request.QueryStringParameters["q"])

			q, e := requiredValue(request.QueryStringParameters["q"], true)
			if e != nil {
				return handleUpstreamResponse(nil, errors.New("missing ebay search query"), nil)
			}
			marketPlace, e := requiredValue(request.QueryStringParameters["marketplace"], true)
			if e != nil {
				return handleUpstreamResponse(nil, errors.New("missing ebay search marketplace"), nil)
			}
			limit, e := requiredValue(request.QueryStringParameters["limit"], true)
			if e != nil {
				return handleUpstreamResponse(nil, errors.New("missing ebay search limit"), nil)
			}
			nLimit, _ := strconv.Atoi(limit)

			cacheKey := "/ebay_search?q=" + q + "&marketplace=" + marketPlace + "&limit" + limit
			cachedCommand := rdb.Get(ctx, cacheKey)
			cachedValue := cachedCommand.Val()
			if cachedValue != "" {
				c := CachedResponse{cacheKey, cachedValue}
				return handleUpstreamResponse(nil, nil, &c)
			}

			token := GetEbayAccessToken(rdb)
			if token == "" {
				return handleUpstreamResponse(nil, errors.New("no ebay access token"), nil)
			}

			res, err := search.EbaySearch(q, marketPlace, nLimit, token, ebayCampaignId)
			if err != nil {
				sentry.CaptureException(err)
				return handleUpstreamResponse(nil, err, nil)
			}

			b, err := json.Marshal(res)
			if err != nil {
				sentry.CaptureException(err)
				return handleUpstreamResponse(nil, err, nil)
			}
			rdb.Set(ctx, cacheKey, string(b), cacheDuration)
			return handleUpstreamResponse(res, err, nil)
		}

		return handleUpstreamResponse(nil, errors.New("bad request"), nil)
	}
}

func GetEbayAccessToken(rdb *redis.Client) string {
	cacheKey := "ebay_access_token"
	cachedCommand := rdb.Get(ctx, cacheKey)
	cachedValue := cachedCommand.Val()
	if cachedValue != "" {
		sentry.CaptureMessage("cache_hit: ebay_access_token")
		return cachedValue
	}
	t, err := search.EbayGetAccessToken()
	if err != nil {
		sentry.CaptureException(err)
		return ""
	}
	rdb.Set(ctx, cacheKey, t.AccessToken, time.Second*time.Duration(t.ExpiresIn))
	sentry.CaptureMessage("cache_set: ebay_access_token")
	return t.AccessToken
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
	openAIApiKey, e := requiredValue(os.Getenv("OPENAI_API_KEY"), false)
	if e != nil {
		log.Fatalf("missing openapi key")
	}
	ebayCampaignId, e := requiredValue(os.Getenv("EBAY_CAMPAIGN_ID"), false)
	if e != nil {
		log.Fatalf("missing ebay compaign id")
	}
	sErr := sentry.Init(sentry.ClientOptions{
		Dsn:         sentryDsn,
		Environment: environment,
	})
	if sErr != nil {
		log.Fatalf("sentry.Init: %s", sErr)
	}

	redisUrl, e := requiredValue(os.Getenv("REDIS_URL"), false)
	if e != nil {
		log.Fatalf("missing redis_url")
	}
	options, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("failed to parse redis_url")
	}
	cacheDuration, dErr := time.ParseDuration("240h")
	if dErr != nil {
		log.Fatalf("failed to define default cache duration")
	}

	rdb := redis.NewClient(options)
	lambda.Start(wrappedHandler(rdb, cacheDuration, openAIApiKey, ebayCampaignId))
}
