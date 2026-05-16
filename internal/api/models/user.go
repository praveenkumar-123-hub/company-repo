package models

import "time"

type Role string

const (
	RoleAdmin Role = "Admin"
	RoleUser  Role = "User"
)

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"` // Never export password
	Role      Role      `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

type SignupRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     Role   `json:"role"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}