package schedule

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"schedule-downloader/models"
	"strings"
)

var apiURL = func(startDate, endDate string) string {
	return fmt.Sprintf(
		"https://byot.kreezee-sports.com/api/v2/solutions/8031/schedule?startDate=%s%%2000:00:00&endDate=%s%%2023:59:59",
		startDate, endDate,
	)
}

// FetchScheduleWeek fetches the schedule between startDate and endDate using the given cookie.
func FetchScheduleWeek(startDate, endDate, cookie string) ([]models.ScheduleItem, error) {
	url := apiURL(startDate, endDate)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

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

	var data []models.ScheduleItem
	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %v", err)
	}

	return data, nil
}
