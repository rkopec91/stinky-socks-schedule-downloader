package schedule

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchScheduleWeek(t *testing.T) {
	mockData := `[{"LocalTeamName":"Puck Luck","VisitorTeamName":"Team B","Date":"2025-11-07","StartTime":"14:00:00","EndTime":"15:00:00","SportCenterName":"Rink 1"}]`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(mockData))
	}))
	defer server.Close()

	oldAPIURL := apiURL
	apiURL = func(startDate, endDate string) string { return server.URL }
	defer func() { apiURL = oldAPIURL }()

	items, err := FetchScheduleWeek("2025-11-07", "2025-11-07", "fake-cookie")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("Expected 1 schedule item, got %d", len(items))
	}
}
