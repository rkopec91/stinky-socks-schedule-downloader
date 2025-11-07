package models

import "time"

type Event struct {
	Title    string
	Start    time.Time
	End      time.Time
	Location string
	Home     string
	Away     string
}

type ScheduleItem struct {
	LocalTeamName   string `json:"LocalTeamName"`
	VisitorTeamName string `json:"VisitorTeamName"`
	Date            string `json:"Date"`
	StartTime       string `json:"StartTime"`
	EndTime         string `json:"EndTime"`
	SportCenterName string `json:"SportCenterName"`
}
