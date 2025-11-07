package calendar

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"schedule-downloader/models"
	"time"

	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func NewService(ctx context.Context, client *http.Client) (*calendar.Service, error) {
	return calendar.NewService(ctx, option.WithHTTPClient(client))
}

func SyncEvents(srv *calendar.Service, events []models.Event) {
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

		existing, err := srv.Events.List("primary").
			ShowDeleted(false).
			SingleEvents(true).
			TimeMin(e.Start.Add(-5 * time.Minute).Format(time.RFC3339)).
			TimeMax(e.End.Add(5 * time.Minute).Format(time.RFC3339)).
			Q(e.Title).
			Do()

		if err != nil {
			log.Printf("❌ Failed to query existing events for %s: %v", e.Title, err)
			continue
		}

		if len(existing.Items) > 0 {
			evID := existing.Items[0].Id
			updated, err := srv.Events.Update("primary", evID, gEvent).Do()
			if err != nil {
				log.Printf("❌ Failed to update event %s: %v", gEvent.Summary, err)
			} else {
				fmt.Printf("🔄 Event updated: %s (%s)\n", updated.Summary, updated.HtmlLink)
			}
		} else {
			created, err := srv.Events.Insert("primary", gEvent).Do()
			if err != nil {
				log.Printf("❌ Failed to create event %s: %v", gEvent.Summary, err)
			} else {
				fmt.Printf("✅ Event created: %s (%s)\n", created.Summary, created.HtmlLink)
			}
		}
	}
}
