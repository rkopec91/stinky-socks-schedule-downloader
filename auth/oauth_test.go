package auth

import (
	"os"
	"testing"

	"golang.org/x/oauth2"
)

func TestTokenSaveLoad(t *testing.T) {
	token := &oauth2.Token{AccessToken: "12345"}
	file := "test_token.json"

	SaveToken(file, token)
	defer os.Remove(file)

	loaded, err := TokenFromFile(file)
	if err != nil {
		t.Fatalf("Error loading token: %v", err)
	}
	if loaded.AccessToken != token.AccessToken {
		t.Errorf("Expected %s, got %s", token.AccessToken, loaded.AccessToken)
	}
}
