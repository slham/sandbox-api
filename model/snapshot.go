package model

import "time"

type Snapshot struct {
	ID         string `json:"id"`
	UserID     string
	CalendarID string    `json:"calendar_id"`
	Done       time.Time `json:"done"`
	Workout    Workout   `json:"workout"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
}
