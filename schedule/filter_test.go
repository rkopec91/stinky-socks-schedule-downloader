package schedule

import (
	"schedule-downloader/models"
	"testing"
	"time"
)

func TestFilterAndTransform(t *testing.T) {
	loc, _ := time.LoadLocation("America/New_York")

	items := []models.ScheduleItem{
		{LocalTeamName: "Puck Luck", VisitorTeamName: "Other Team", Date: "2025-11-07", StartTime: "14:00:00", EndTime: "15:00:00", SportCenterName: "Rink 1"},
		{LocalTeamName: "Another Team", VisitorTeamName: "Puck Luck", Date: "2025-11-08", StartTime: "16:00:00", EndTime: "17:00:00", SportCenterName: "Rink 2"},
		{LocalTeamName: "Other Team", VisitorTeamName: "Third Team", Date: "2025-11-09", StartTime: "18:00:00", EndTime: "19:00:00", SportCenterName: "Rink 3"},
	}

	events := FilterAndTransform(items, "Puck Luck", loc)
	if len(events) != 2 {
		t.Errorf("Expected 2 events, got %d", len(events))
	}

	for _, e := range events {
		if e.Home != "Puck Luck" && e.Away != "Puck Luck" {
			t.Errorf("Event does not match filter: %+v", e)
		}
		if e.Start.IsZero() || e.End.IsZero() {
			t.Errorf("Event times not parsed: %+v", e)
		}
	}
}
