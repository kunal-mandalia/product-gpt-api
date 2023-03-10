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
	expected := 10
	token, _ := EbayGetAccessToken()
	campaignId := os.Getenv("EBAY_CAMPAIGN_ID")
	actual, _ := EbaySearch(q, "EBAY_GB", 10, token.AccessToken, campaignId)

	if len(actual.ItemSummaries) != expected {
		t.Fatalf(`Search(%s) failed`, q)
	}
}
