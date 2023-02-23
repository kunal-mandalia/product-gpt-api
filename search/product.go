package search

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
)

type CustomSearchResponseItem struct {
	Title   string `json:"title"`
	Link    string `json:"link"`
	Snippet string `json:"snippet"`
}
type CustomSearchResponse struct {
	Items []CustomSearchResponseItem `json:"items"`
}

func ProductList(apiKey string, q string) (CustomSearchResponse, error) {
	data := CustomSearchResponse{}

	URLString := "https://customsearch.googleapis.com/customsearch/v1?" +
		"q=" + url.QueryEscape(q) +
		"&key=" + apiKey +
		"&cx=510f41b60d9334db2" +
		"&safe=1"

	r, err := http.NewRequest("GET", URLString, nil)
	if err != nil {
		log.Fatal(err)
		return data, err
	}

	r.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return data, err
	}
	if res.StatusCode != http.StatusOK {
		return data, errors.New("upstream error")
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&data)
	return data, nil
}
