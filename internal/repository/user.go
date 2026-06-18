package repository

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func CreateUser(username, email, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (username, email, password_hash) VALUES ($1, $2, $3)`
	_, err = DB.Exec(context.Background(), query, username, email, string(hash))
	return err
}

func AuthenticateUser(email, password string) (int, string, error) {
	var id int
	var username string
	var hash string

	query := `SELECT id, username, password_hash FROM users WHERE email = $1`
	err := DB.QueryRow(context.Background(), query, email).Scan(&id, &username, &hash)
	if err != nil {
		return 0, "", errors.New("неверный email или пароль")
	}

	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return 0, "", errors.New("неверный email или пароль")
	}

	return id, username, nil
}