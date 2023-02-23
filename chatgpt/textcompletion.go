package chatgpt

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type TextCompletionBody struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Max_tokens  int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
}

type TextCompletionChoice struct {
	Text          string `json:"text"`
	Finish_reason string `json:"finish_reason"`
}

type TextCompletionResponse struct {
	Id      string                 `json:"id"`
	Choices []TextCompletionChoice `json:"choices"`
}

func TextCompletion(apiKey string, input string) (TextCompletionResponse, error) {
	url := "https://api.openai.com/v1/completions"
	data := TextCompletionBody{"text-davinci-003", input, 2000, 0}
	m, err := json.Marshal(data)
	b := bytes.NewBuffer(m)
	d1 := TextCompletionResponse{}

	if err != nil {
		return d1, err
	}

	r, err := http.NewRequest("POST", url, b)
	if err != nil {
		log.Fatal(err)
		return d1, err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return d1, err
	}
	if res.StatusCode != http.StatusOK {
		return d1, errors.New("upstream error")
	}

	defer res.Body.Close()

	json.NewDecoder(res.Body).Decode(&d1)
	return d1, nil
}

func ProductRecommendationsQuery(initialQuery string, initialResponse string) string {
	return "\ncontext:\n" +
		initialQuery + "\n" +
		initialResponse + "\n\n" +
		"prompt:\n" +
		"Identify and provide information on products and services found within the context above." +
		" Specifically, output an array of JSON objects which include the following:" +
		" key: name, value: name of product or service" +
		" key: link, value: a link to that product or service," +
		" key: isProduct, value: true if the link refers to a product, false if it refers to a service" +
		" key: characterRange, value: an array of start / end numbers indicating "
}

func EntityExtractionQuery(q string) string {
	return "I want you to identify and classify entities in the following passage of text until I state END_OF_PASSAGE :\n" +
		q + "END_OF_PASSAGE" +
		"\nReturn a table with the following columns:" +
		" The 1st is the entity name." +
		" The 2nd column is the entity type." +
		" The 3rd column shows if the entity is saleable or not." +
		" The 4th column is the category of the entity" +
		" The 5th column is the entity type" +
		" The 6th column is either Product, Service, or Other depending on the entity type"
}
