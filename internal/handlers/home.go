package handlers

import (
	"html/template"
	"net/http"

	"sycinema/internal/models"
	"sycinema/internal/repository"
	"sycinema/internal/service"
)

type HomeData struct {
	Username string
	Series   []models.Series
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	data := HomeData{}

	if cookie, err := r.Cookie("token"); err == nil {
		if claims, err := service.ParseToken(cookie.Value); err == nil {
			data.Username = claims.Username
		}
	}

	data.Series, _ = repository.GetAllSeries() // Игнорируем ошибку для краткости, в проде лучше логгировать

	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}