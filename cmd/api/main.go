package main

import (
	"log"
	"net/http"

	"sycinema/internal/handlers"
	"sycinema/internal/repository"
)

func main() {
	log.Println("Инициализация базы данных...")
	repository.ConnectDB()

	// Раздача статики и медиа
	staticDir := http.FileServer(http.Dir("web/static"))
	http.Handle("/static/", http.StripPrefix("/static/", staticDir))

	mediaDir := http.FileServer(http.Dir("media"))
	http.Handle("/media/", http.StripPrefix("/media/", mediaDir))

	// Роуты
	http.HandleFunc("/", handlers.HomeHandler)

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.RegisterPageHandler(w, r)
		case http.MethodPost:
			handlers.RegisterPostHandler(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

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

	http.HandleFunc("/logout", handlers.LogoutHandler)
	http.HandleFunc("/series", handlers.SeriesHandler)
	http.HandleFunc("/watch", handlers.WatchHandler)
	http.HandleFunc("/bookmark/toggle", handlers.ToggleBookmarkHandler)
	http.HandleFunc("/bookmarks", handlers.BookmarksPageHandler)
	http.HandleFunc("/episode/favorite", handlers.ToggleFavoriteEpisodeHandler)

	// Запуск
	port := ":8080"
	log.Printf("🎬 Сервер SyCinema запущен на http://localhost%s", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("Критическая ошибка:", err)
	}
}
