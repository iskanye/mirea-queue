package models

type Group struct {
	Name        string `json:"fullTitle"`
	ScheduleUrl string `json:"iCalLink"`
}
