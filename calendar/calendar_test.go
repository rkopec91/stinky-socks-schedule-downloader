package calendar

import (
	"testing"

	"schedule-downloader/models"
)

func TestSyncEvents_Empty(t *testing.T) {
	var events []models.Event

	// We just check it doesn't panic with empty events
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("SyncEvents panicked with empty events")
		}
	}()

	SyncEvents(nil, events) // srv can be nil because we're only testing safety
}
