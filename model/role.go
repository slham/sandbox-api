package model

import "time"

type Role struct {
	ID      int       `json:"id"`
	Name    string    `json:"name"`
	Created time.Time `json:"created"`
	Updated time.Time `json:"updated"`
}

type UserRole struct {
	UserID string `json:"user_id"`
	RoleID int    `json:"role_id"`
}
