package search

import (
	"fmt"
	"testing"
)

func Test_GetAccessToken(t *testing.T) {

	actual, _ := EbayGetAccessToken()

	expected := EbayAccessTokenResponse{}
	expected.TokenType = "Application Access Token"

	fmt.Println(actual)
	if actual.TokenType != expected.TokenType {
		t.Fatalf(`GetAccessToken() failed`)
	}
}

func Test_EbaySearch(t *testing.T) {
	q := "cake"
	expected := "bar"
	token, _ := EbayGetAccessToken()
	actual, _ := EbaySearch(q, token.AccessToken)

	if actual != expected {
		t.Fatalf(`Search(%s) failed`, q)
	}
}
