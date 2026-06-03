package main

import (
	"log"
	"net/http"

	// Импортируем наши собственные пакеты.
	// Если твой модуль в go.mod называется иначе, замени "sycinema" на свое название.
	"sycinema/internal/handlers"
	"sycinema/internal/repository"
)

func main() {
	// 1. Инициализация базы данных
	// Первым делом поднимаем соединение с PostgreSQL.
	// Если база недоступна, программа завершится с ошибкой (мы прописали log.Fatal внутри ConnectDB).
	log.Println("Инициализация базы данных...")
	repository.ConnectDB()

	// 2. Раздача статических файлов (CSS, картинки, скрипты)
	// Указываем серверу искать статику в папке web/static.
	// StripPrefix отрезает "/static/" из URL, чтобы сервер искал файл "css/style.css", а не "static/css/style.css" внутри папки.
	staticDir := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticDir))

	// 3. Раздача медиафайлов (Наши тестовые фильмы)
	// То же самое, но для папки media.
	mediaDir := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", mediaDir))

	// 4. Настройка маршрутов (Роутинг)
	// Главная страница
	http.HandleFunc("/", handlers.HomeHandler)

	// Страница регистрации
	// Так как стандартный http-роутер в Go до версии 1.22 не умеет элегантно разделять GET и POST,
	// мы делаем это вручную внутри анонимной функции.
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			// Если браузер просто запрашивает страницу — отдаем HTML форму
			handlers.RegisterPageHandler(w, r)
		case http.MethodPost:
			// Если форма отправляет данные — обрабатываем их и пишем в БД
			handlers.RegisterPostHandler(w, r)
		default:
			// Если пришел какой-то другой запрос (например, PUT или DELETE) — отказываем
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	// Страница входа
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.LoginPageHandler(w, r)
		case http.MethodPost:
			handlers.LoginPostHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	// 5. Запуск сервера
	port := ":8080"
	log.Printf("🎬 Сервер SyCinema успешно запущен на http://localhost%s", port)

	// ListenAndServe блокирует текущую горутину и начинает слушать входящие подключения.
	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Критическая ошибка при работе сервера:", err)
	}
}
