package schedule

import (
	"schedule-downloader/models"
	"strings"
	"time"
)

// FilterAndTransform filters schedule items by team name and converts them to Event structs.
func FilterAndTransform(items []models.ScheduleItem, teamFilter string, loc *time.Location) []models.Event {
	events := []models.Event{}

	for _, e := range items {
		if strings.Contains(strings.ToLower(e.LocalTeamName), strings.ToLower(teamFilter)) ||
			strings.Contains(strings.ToLower(e.VisitorTeamName), strings.ToLower(teamFilter)) {

			dateOnly := e.Date[:10] // YYYY-MM-DD
			startTime, err := time.ParseInLocation("2006-01-02T15:04:05", dateOnly+"T"+e.StartTime, loc)
			if err != nil {
				continue
			}
			endTime, err := time.ParseInLocation("2006-01-02T15:04:05", dateOnly+"T"+e.EndTime, loc)
			if err != nil {
				continue
			}

			title := e.LocalTeamName + " vs " + e.VisitorTeamName
			events = append(events, models.Event{
				Title:    title,
				Start:    startTime,
				End:      endTime,
				Location: e.SportCenterName,
				Home:     e.LocalTeamName,
				Away:     e.VisitorTeamName,
			})
		}
	}

	return events
}
