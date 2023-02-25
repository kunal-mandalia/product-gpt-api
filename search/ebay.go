package search

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type EbayAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func EbayGetAccessToken() (EbayAccessTokenResponse, error) {
	tokenEndpoint := os.Getenv("EBAY_TOKEN_ENDPOINT")
	clientId := os.Getenv("EBAY_CLIENT_ID")
	clientSecret := os.Getenv("EBAY_CLIENT_SECRET")

	token := EbayAccessTokenResponse{}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("scope", "https://api.ebay.com/oauth/api_scope")

	creds := clientId + ":" + clientSecret
	b64Creds := base64.StdEncoding.EncodeToString([]byte(creds))

	r, _ := http.NewRequest("POST", tokenEndpoint, strings.NewReader(data.Encode()))
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Authorization", "Basic "+b64Creds)

	client := &http.Client{}
	resp, err := client.Do(r)
	if err != nil {
		return token, err
	}
	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println(string(bodyBytes))
		return token, errors.New("upstream error")
	}

	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(&token)
	return token, nil
}

func EbaySearch(q string, accessToken string) (interface{}, error) {
	data := EbaySearchResponse{}
	apiEndpoint := os.Getenv("EBAY_BROWSE_API_ENDPOINT")
	URLString := apiEndpoint + "/item_summary/search?limit=3&q=" + url.QueryEscape(q)

	r, _ := http.NewRequest("GET", URLString, nil)
	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+accessToken)
	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		bodyBytes, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyBytes))
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(res.Body)
		fmt.Println(string(bodyBytes))
		return nil, errors.New("upstream error")
	}
	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&data)
	fmt.Println(data)
	return data, nil
}
