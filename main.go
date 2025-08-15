package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

type Event struct {
	Title    string
	Start    time.Time
	End      time.Time
	Location string
}

type ScheduleItem struct {
	LocalTeamName   string `json:"LocalTeamName"`
	VisitorTeamName string `json:"VisitorTeamName"`
	Date            string `json:"Date"`
	StartTime       string `json:"StartTime"`
	EndTime         string `json:"EndTime"`
	SportCenterName string `json:"SportCenterName"`
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

	// Headers
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
	teamFilter := "Puck Luck"
	cookie := "YOUR_SESSION_COOKIE_HERE" // Replace with your actual cookie

	// Default: today → one year from today
	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")

	fmt.Printf("Fetching schedule for team containing '%s' from %s to %s...\n", teamFilter, startDate, endDate)

	eventsJSON, err := fetchScheduleWeek(startDate, endDate, cookie)
	if err != nil {
		log.Fatalf("Error fetching schedule: %v", err)
	}

	events := []Event{}
	for _, e := range eventsJSON {
		if strings.Contains(strings.ToLower(e.LocalTeamName), strings.ToLower(teamFilter)) ||
			strings.Contains(strings.ToLower(e.VisitorTeamName), strings.ToLower(teamFilter)) {

			startTime, _ := time.Parse("2006-01-02T15:04:05", e.Date+"T"+e.StartTime)
			endTime, _ := time.Parse("2006-01-02T15:04:05", e.Date+"T"+e.EndTime)

			title := fmt.Sprintf("%s vs %s", e.LocalTeamName, e.VisitorTeamName)
			location := e.SportCenterName
			events = append(events, Event{
				Title:    title,
				Start:    startTime,
				End:      endTime,
				Location: location,
			})
		}
	}

	fmt.Println("Filtered Events:")
	for _, ev := range events {
		fmt.Printf("%s @ %s — %s to %s\n", ev.Title, ev.Location, ev.Start.Format(time.RFC1123), ev.End.Format(time.RFC1123))
	}
}
