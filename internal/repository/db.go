package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB - глобальная переменная с пулом соединений, 
// через которую мы будем делать запросы из других мест программы.
var DB *pgxpool.Pool

func ConnectDB() {
	// Строка подключения (DSN - Data Source Name).
	// Формат: postgres://пользователь:пароль@хост:порт/имя_базы
	// Как договаривались, используем пароль admin.
	dsn := "postgres://postgres:admin@localhost:5432/sycinema_db"

	var err error
	// Создаем пул соединений
	DB, err = pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal("❌ Ошибка при создании пула соединений:", err)
	}

	// Делаем тестовый пинг, чтобы убедиться, что сервер БД реально работает
	err = DB.Ping(context.Background())
	if err != nil {
		log.Fatal("❌ База данных недоступна. Проверь, запущен ли PostgreSQL:", err)
	}

	fmt.Println("🔌 Успешное подключение к PostgreSQL!")
}