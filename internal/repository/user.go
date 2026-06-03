package repository

import (
	"context"
	"fmt"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser хеширует пароль и сохраняет нового пользователя в БД
func CreateUser(username, email, password string) error {
	// 1. Генерируем хеш пароля (стоимость 10 - оптимальный баланс скорости и безопасности)
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return fmt.Errorf("ошибка хеширования пароля: %v", err)
	}

	// 2. SQL-запрос на вставку. Используем $1, $2, $3 для защиты от SQL-иньекций
	query := `
		INSERT INTO users (username, email, password_hash) 
		VALUES ($1, $2, $3)
	`

	// 3. Выполняем запрос через наш глобальный пул соединений DB
	_, err = DB.Exec(context.Background(), query, username, email, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("ошибка при сохранении пользователя: %v", err)
	}

	return nil
}

func AuthenticateUser(email, password string) (string, error) {
	var username string
	var hash string

	// 1. Ищем пользователя по Email
	query := `SELECT username, password_hash FROM users WHERE email = $1`
	
	// QueryRow сканирует ровно одну строку из базы в наши переменные
	err := DB.QueryRow(context.Background(), query, email).Scan(&username, &hash)
	if err != nil {
		// Если ничего не нашли (или другая ошибка)
		return "", errors.New("неверный email или пароль")
	}

	// 2. Сравниваем введенный пароль с хешем из БД
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// Пароли не совпали
		return "", errors.New("неверный email или пароль")
	}

	// 3. Если всё ок, возвращаем имя пользователя
	return username, nil
}