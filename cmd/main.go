package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"schedule-downloader/auth"
	googlecalendar "schedule-downloader/calendar"
	"schedule-downloader/models"
	"schedule-downloader/schedule"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func main() {
	// CLI flags
	teamFilter := flag.String("team", "Puck Luck", "Comma-separated list of team names to filter schedule")
	startDate := flag.String("start", time.Now().Format("2006-01-02"), "Start date (YYYY-MM-DD)")
	endDate := flag.String("end", time.Now().AddDate(1, 0, 0).Format("2006-01-02"), "End date (YYYY-MM-DD)")
	cookie := flag.String("cookie", os.Getenv("KREEZEE_COOKIE"), "Kreezee session cookie")
	calendars := flag.String("calendar", "primary", "Comma-separated list of Google Calendar IDs")

	flag.Parse()

	if *cookie == "" {
		log.Fatal("Error: KREEZEE_COOKIE must be provided either as flag or environment variable")
	}

	fmt.Printf("Fetching schedule for teams '%s' from %s to %s...\n", *teamFilter, *startDate, *endDate)

	// Split multiple teams and calendars
	teams := parseCommaSeparated(*teamFilter)
	calendarIDs := parseCommaSeparated(*calendars)

	// Fetch and filter schedule
	scheduleItems, err := schedule.FetchScheduleWeek(*startDate, *endDate, *cookie)
	if err != nil {
		log.Fatalf("Error fetching schedule: %v", err)
	}

	loc, _ := time.LoadLocation("America/New_York")
	events := []models.Event{}
	for _, team := range teams {
		events = append(events, schedule.FilterAndTransform(scheduleItems, team, loc)...)
	}

	fmt.Println("Filtered Events:")
	for _, ev := range events {
		fmt.Printf("%s @ %s — %s to %s\n", ev.Title, ev.Location, ev.Start.Format(time.RFC1123), ev.End.Format(time.RFC1123))
	}

	// Google Calendar setup
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.ConfigFromJSON(b, calendar.CalendarScope)
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client %v", err)
	}
	client := auth.GetClient(config)

	srv, err := googlecalendar.NewService(context.Background(), client)
	if err != nil {
		log.Fatalf("Unable to create Calendar service: %v", err)
	}

	// Sync events to all specified calendars
	for _, calID := range calendarIDs {
		fmt.Printf("Syncing events to calendar: %s\n", calID)
		googlecalendar.SyncEvents(srv, events)
	}
}

// parseCommaSeparated splits a comma-separated string and trims spaces
func parseCommaSeparated(s string) []string {
	parts := strings.Split(s, ",")
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	}
	return parts
}
