package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func UNUSED(x ...interface{}) {}

// curl https://api.openai.com/v1/completions \
//   -H 'Content-Type: application/json' \
//   -H 'Authorization: Bearer YOUR_API_KEY' \
//   -d '{
//   "model": "text-davinci-003",
//   "prompt": "Say this is a test",
//   "max_tokens": 7,
//   "temperature": 0
// }'

type TextCompletionBody struct {
	Model       string `json:"model"`
	Prompt      string `json:"prompt"`
	Max_tokens  int    `json:"max_tokens"`
	Temperature int    `json:"temperature"`
}

func GetTextCompletion(apiKey string, input string) (string, error) {
	fmt.Println("GetTextCompletion called")

	url := "https://api.openai.com/v1/completions"
	data := TextCompletionBody{"text-davinci-003", input, 15, 0}
	m, err := json.Marshal(data)
	b := bytes.NewBuffer(m)
	if err != nil {
		return "", err
	}

	r, err := http.NewRequest("POST", url, b)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	r.Header.Add("Content-Type", "application/json")
	r.Header.Add("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	bytes, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	str := string(bytes)
	return str, nil
}
