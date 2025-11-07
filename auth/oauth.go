package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"golang.org/x/oauth2"
)

func TokenValid(config *oauth2.Config, tok *oauth2.Token) bool {
	ctx := context.Background()
	tokenSource := config.TokenSource(ctx, tok)
	newTok, err := tokenSource.Token()
	if err != nil {
		return false
	}
	if newTok.AccessToken != tok.AccessToken {
		SaveToken("token.json", newTok)
	}
	return true
}

func GetClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := TokenFromFile(tokenFile)
	if err != nil || !TokenValid(config, tok) {
		tok = GetTokenFromWeb(config)
		SaveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	listener := "localhost:8080"
	codeCh := make(chan string)
	srv := &http.Server{Addr: listener}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if errMsg := r.URL.Query().Get("error"); errMsg != "" {
			http.Error(w, "Authorization error: "+errMsg, http.StatusBadRequest)
			return
		}
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "No authorization code found.", http.StatusBadRequest)
			return
		}
		fmt.Fprintln(w, "✅ Authorization received! You can close this window.")
		codeCh <- code
	})

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start local server: %v", err)
		}
	}()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("\n🔗 Open the following URL in your browser (may open automatically):\n%v\n\n", authURL)
	go openBrowser(authURL)

	authCode := <-codeCh
	_ = srv.Shutdown(context.Background())

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

func openBrowser(url string) error {
	var cmd string
	var args []string
	switch os := runtime.GOOS; os {
	case "windows":
		cmd = "rundll32"
		args = []string{"url.dll,FileProtocolHandler", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default:
		cmd = "xdg-open"
		args = []string{url}
	}
	return exec.Command(cmd, args...).Start()
}

func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
