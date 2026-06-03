package models

import "time"

// User описывает структуру нашего пользователя в базе данных
type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // минус означает, что хэш пароля никогда не улетит на фронтенд в JSON
	CreatedAt    time.Time `json:"created_at"`
}