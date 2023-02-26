package search

import (
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func setup() {
	godotenv.Load()
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func Test_GetAccessToken(t *testing.T) {
	actual, _ := EbayGetAccessToken()

	expected := EbayAccessTokenResponse{}
	expected.TokenType = "Application Access Token"

	fmt.Println(actual)
	if actual.TokenType != expected.TokenType {
		t.Fatalf(`EbayGetAccessToken() failed`)
	}
}

func Test_EbaySearch(t *testing.T) {
	q := "cake"
	expected := 3
	token, _ := EbayGetAccessToken()
	actual, _ := EbaySearch(q, "EBAY_GB", token.AccessToken)

	if len(actual.ItemSummaries) != expected {
		t.Fatalf(`Search(%s) failed`, q)
	}
}
