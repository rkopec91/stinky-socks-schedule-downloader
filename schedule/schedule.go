package schedule

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"schedule-downloader/models"
	"strings"
	"time"
)

func FetchScheduleWeek(startDate, endDate, cookie string) ([]models.ScheduleItem, error) {
	url := fmt.Sprintf(
		"https://byot.kreezee-sports.com/api/v2/solutions/8031/schedule?startDate=%s%%2000:00:00&endDate=%s%%2023:59:59",
		startDate, endDate,
	)

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
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

	var data []models.ScheduleItem
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return data, nil
}

func FilterAndTransform(items []models.ScheduleItem, teamFilter string, loc *time.Location) []models.Event {
	events := []models.Event{}
	for _, e := range items {
		if strings.Contains(strings.ToLower(e.LocalTeamName), strings.ToLower(teamFilter)) ||
			strings.Contains(strings.ToLower(e.VisitorTeamName), strings.ToLower(teamFilter)) {

			if len(e.Date) < 10 {
				log.Printf("Unexpected date format: %s", e.Date)
				continue
			}
			dateOnly := e.Date[:10]
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

			events = append(events, models.Event{
				Title:    fmt.Sprintf("%s vs %s", e.LocalTeamName, e.VisitorTeamName),
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
