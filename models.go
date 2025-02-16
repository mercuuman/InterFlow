package main

import "time"

type UserOut struct {
	UserID       int       `json:"user_id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	AvatarURL    string    `json:"avatar_url"`
	Description  string    `json:"description"`
	Respects     int       `json:"respects"`
	IsDeleted    bool      `json:"is_deleted"`
	CreatedAt    time.Time `json:"created_at"`
}

// UserIn — структура для входных данных от клиента (регистрация).
type UserIn struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserState хранит состояние пользователя, включая токен подтверждения.
type UserState struct {
	UserID     int       `json:"user_id"`
	Token      string    `json:"token"`
	IsVerified bool      `json:"is_verified"`
	CreatedAt  time.Time `json:"created_at"`
}
