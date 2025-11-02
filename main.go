package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

type Event struct {
	Title    string
	Start    time.Time
	End      time.Time
	Location string
	Home     string
	Away     string
}

type ScheduleItem struct {
	LocalTeamName   string `json:"LocalTeamName"`
	VisitorTeamName string `json:"VisitorTeamName"`
	Date            string `json:"Date"`
	StartTime       string `json:"StartTime"`
	EndTime         string `json:"EndTime"`
	SportCenterName string `json:"SportCenterName"`
}

func tokenValid(config *oauth2.Config, tok *oauth2.Token) bool {
	ctx := context.Background()
	tokenSource := config.TokenSource(ctx, tok)

	// Try getting a new token
	newTok, err := tokenSource.Token()
	if err != nil {
		return false
	}

	// If the new token differs from the original, persist it
	if newTok.AccessToken != tok.AccessToken {
		saveToken("token.json", newTok)
	}

	return true
}

func getClient(config *oauth2.Config) *http.Client {
	tokenFile := "token.json"
	tok, err := tokenFromFile(tokenFile)
	if err != nil || !tokenValid(config, tok) {
		tok = getTokenFromWeb(config)
		saveToken(tokenFile, tok)
	}
	return config.Client(context.Background(), tok)
}

func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	// Local redirect listener (must match Google Cloud redirect URI)
	listener := "localhost:8080"
	codeCh := make(chan string)
	srv := &http.Server{Addr: listener}

	// Handle the OAuth2 callback
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

		fmt.Fprintln(w, "✅ Authorization received! You can close this window and return to the terminal.")
		codeCh <- code
	})

	// Start the server asynchronously
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start local server: %v", err)
		}
	}()

	// Open the URL in the user’s browser
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("\n🔗 Open the following URL in your browser (it may open automatically):\n%v\n\n", authURL)
	// Attempt to open browser automatically (optional)
	go func() {
		_ = openBrowser(authURL)
	}()

	// Wait for the authorization code
	authCode := <-codeCh

	// Shut down the server after we get the code
	_ = srv.Shutdown(context.Background())

	// Exchange the code for a token
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
	default: // linux, freebsd, etc.
		cmd = "xdg-open"
		args = []string{url}
	}
	return exec.Command(cmd, args...).Start()
}

func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.Create(path)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func fetchScheduleWeek(startDate, endDate, cookie string) ([]ScheduleItem, error) {
	url := fmt.Sprintf(
		"https://byot.kreezee-sports.com/api/v2/solutions/8031/schedule?startDate=%s%%2000:00:00&endDate=%s%%2023:59:59",
		startDate, endDate,
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Cookie", cookie)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if strings.HasPrefix(string(bodyBytes), "<") {
		return nil, fmt.Errorf("got HTML instead of JSON; likely requires login or headers")
	}

	var data []ScheduleItem
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return data, nil
}

func main() {
	ctx := context.Background()
	teamFilter := "Puck Luck"
	cookie := "YOUR_SESSION_COOKIE_HERE"

	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")

	fmt.Printf("Fetching schedule for team containing '%s' from %s to %s...\n", teamFilter, startDate, endDate)

	eventsJSON, err := fetchScheduleWeek(startDate, endDate, cookie)
	if err != nil {
		log.Fatalf("Error fetching schedule: %v", err)
	}

	loc, _ := time.LoadLocation("America/New_York")
	events := []Event{}
	for _, e := range eventsJSON {
		if strings.Contains(strings.ToLower(e.LocalTeamName), strings.ToLower(teamFilter)) ||
			strings.Contains(strings.ToLower(e.VisitorTeamName), strings.ToLower(teamFilter)) {

			dateOnly := e.Date[:10] // YYYY-MM-DD
			startTime, err := time.ParseInLocation("2006-01-02T15:04:05", dateOnly+"T"+e.StartTime, loc)
			if err != nil {
				log.Printf("Failed to parse start time for event %s: %v", e.LocalTeamName+" vs "+e.VisitorTeamName, err)
				continue
			}

			endTime, err := time.ParseInLocation("2006-01-02T15:04:05", dateOnly+"T"+e.EndTime, loc)
			if err != nil {
				log.Printf("Failed to parse end time for event %s: %v", e.LocalTeamName+" vs "+e.VisitorTeamName, err)
				continue
			}

			title := fmt.Sprintf("%s vs %s", e.LocalTeamName, e.VisitorTeamName)
			location := e.SportCenterName
			events = append(events, Event{
				Title:    title,
				Start:    startTime,
				End:      endTime,
				Location: location,
				Home:     e.LocalTeamName,
				Away:     e.VisitorTeamName,
			})
		}
	}

	fmt.Println("Filtered Events:")
	for _, ev := range events {
		fmt.Printf("%s @ %s — %s to %s\n", ev.Title, ev.Location, ev.Start.Format(time.RFC1123), ev.End.Format(time.RFC1123))
	}

	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client %v", err)
	}
	client := getClient(config)

	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}

	for _, e := range events {
		gEvent := &calendar.Event{
			Summary:  e.Title,
			Location: e.Location,
			Start: &calendar.EventDateTime{
				DateTime: e.Start.Format(time.RFC3339),
				TimeZone: "America/New_York",
			},
			End: &calendar.EventDateTime{
				DateTime: e.End.Format(time.RFC3339),
				TimeZone: "America/New_York",
			},
		}

		// Look for an existing event in the same time window with same title
		existing, err := srv.Events.List("primary").
			ShowDeleted(false).
			SingleEvents(true).
			TimeMin(e.Start.Add(-time.Minute * 5).Format(time.RFC3339)).
			TimeMax(e.End.Add(time.Minute * 5).Format(time.RFC3339)).
			Q(e.Title).
			Do()

		if err != nil {
			log.Printf("❌ Failed to query existing events for %s: %v", e.Title, err)
			continue
		}

		if len(existing.Items) > 0 {
			// Update the first matching event
			evID := existing.Items[0].Id
			updated, err := srv.Events.Update("primary", evID, gEvent).Do()
			if err != nil {
				log.Printf("❌ Failed to update event %s: %v", gEvent.Summary, err)
			} else {
				fmt.Printf("🔄 Event updated: %s (%s)\n", updated.Summary, updated.HtmlLink)
			}
		} else {
			// Insert new
			created, err := srv.Events.Insert("primary", gEvent).Do()
			if err != nil {
				log.Printf("❌ Failed to create event %s: %v", gEvent.Summary, err)
			} else {
				fmt.Printf("✅ Event created: %s (%s)\n", created.Summary, created.HtmlLink)
			}
		}
	}
}
