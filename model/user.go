package model

import "time"

type User struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	Password    string    `json:"password"`
	Email       string    `json:"email"`
	Created     time.Time `json:"created"`
	Updated     time.Time `json:"updated"`
	IsActive    bool      `json:"isActive,omitempty"`
	IsSuspended bool      `json:"isSuspended,omitempty"`
	IsVerified  bool      `json:"isVerified,omitempty"`
}
