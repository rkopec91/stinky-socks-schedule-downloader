package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"schedule-downloader/auth"
	calendarPkg "schedule-downloader/calendar"
	"schedule-downloader/schedule"
	"time"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
)

func main() {
	ctx := context.Background()
	teamFilter := "Puck Luck"
	cookie := os.Getenv("KREEZEE_COOKIE") // safer than hardcoding

	startDate := time.Now().Format("2006-01-02")
	endDate := time.Now().AddDate(1, 0, 0).Format("2006-01-02")

	fmt.Printf("Fetching schedule for team containing '%s' from %s to %s...\n", teamFilter, startDate, endDate)

	scheduleItems, err := schedule.FetchScheduleWeek(startDate, endDate, cookie)
	if err != nil {
		log.Fatalf("Error fetching schedule: %v", err)
	}

	loc, _ := time.LoadLocation("America/New_York")
	events := schedule.FilterAndTransform(scheduleItems, teamFilter, loc)

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
	client := auth.GetClient(config)

	srv, err := calendarPkg.NewService(ctx, client)
	if err != nil {
		log.Fatalf("Unable to create Calendar service: %v", err)
	}

	calendarPkg.SyncEvents(srv, events)
}
