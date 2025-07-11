package model

type Calendar struct {
  ID string `json:"id"`
	UserID       string    `json:"user_id"`
	Name string    `json:"name,omitempty"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
}
