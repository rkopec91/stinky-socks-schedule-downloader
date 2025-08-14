package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type Game struct {
	LocalTeamName   string `json:"LocalTeamName"`
	VisitorTeamName string `json:"VisitorTeamName"`
	Date            string `json:"Date"`
	StartTime       string `json:"StartTime"`
	EndTime         string `json:"EndTime"`
	Note            string `json:"Note"`
}

type Event struct {
	Title string
	Start time.Time
	End   time.Time
	Note  string
}

func main() {
	fmt.Println("Schedule Downloader")

	// Example URL for a week
	url := "https://byot.kreezee-sports.com/api/v2/solutions/8031/schedule?startDate=2025-08-11%2000:00:00&endDate=2025-08-17%2023:59:59"

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Failed to access url: %v", err)
	}
	defer resp.Body.Close()

	var games []Game
	if err := json.NewDecoder(resp.Body).Decode(&games); err != nil {
		log.Fatalf("Failed to decode JSON: %v", err)
	}

	events := []Event{}

	for _, g := range games {
		localLower := strings.ToLower(g.LocalTeamName)
		visitorLower := strings.ToLower(g.VisitorTeamName)

		// Filter for any team containing "puck luck"
		if !strings.Contains(localLower, "puck luck") && !strings.Contains(visitorLower, "puck luck") {
			continue
		}

		// Parse start and end times
		startTime, err := time.Parse("2006-01-02T15:04:05", g.Date+"T"+g.StartTime)
		if err != nil {
			log.Printf("Failed to parse start time for %s vs %s: %v", g.LocalTeamName, g.VisitorTeamName, err)
			continue
		}

		endTime, err := time.Parse("2006-01-02T15:04:05", g.Date+"T"+g.EndTime)
		if err != nil {
			log.Printf("Failed to parse end time for %s vs %s: %v", g.LocalTeamName, g.VisitorTeamName, err)
			continue
		}

		title := fmt.Sprintf("%s vs %s", g.LocalTeamName, g.VisitorTeamName)

		events = append(events, Event{
			Title: title,
			Start: startTime,
			End:   endTime,
			Note:  g.Note,
		})
	}

	fmt.Println("Filtered Events:")
	for _, e := range events {
		fmt.Printf("%s | Start: %s | End: %s | Note: %s\n", e.Title, e.Start.Format(time.RFC1123), e.End.Format(time.RFC1123), e.Note)
	}
}
