package model

type User struct {
	ID          int    `json:"id"`
	UserName    string `json:"name"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	Created     string `json:"created"`
	Updated     string `json:"updated"`
	IsActive    bool   `json:"isActive,omitempty"`
	IsSuspended bool   `json:"isSuspended,omitempty"`
	IsVerified  bool   `json:"isVerified,omitempty"`
}
