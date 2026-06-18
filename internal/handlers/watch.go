package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"sycinema/internal/repository"
	"sycinema/internal/service"
)

func WatchHandler(w http.ResponseWriter, r *http.Request) {
	epStr := r.URL.Query().Get("ep")
	epID, err := strconv.Atoi(epStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	episode, err := repository.GetEpisodeByID(epID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	series, _ := repository.GetSeriesByID(episode.SeriesID)

	username := ""
	if cookie, err := r.Cookie("token"); err == nil {
		if claims, err := service.ParseToken(cookie.Value); err == nil {
			username = claims.Username
		}
	}

	data := struct {
		Username string
		Series   interface{}
		Episode  interface{}
	}{
		Username: username,
		Series:   series,
		Episode:  episode,
	}

	tmpl, err := template.ParseFiles("web/templates/watch.html")
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}