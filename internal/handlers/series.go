package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"sycinema/internal/repository"
	"sycinema/internal/models"
	"sycinema/internal/service"
)

func SeriesHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	series, err := repository.GetSeriesByID(id)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	userID := 0
	username := ""
	isBookmarked := false

	if cookie, err := r.Cookie("token"); err == nil {
		if claims, err := service.ParseToken(cookie.Value); err == nil {
			userID = claims.UserID
			username = claims.Username
			isBookmarked = repository.IsBookmarked(userID, id)
		}
	}

	// Передаем userID (даже если он 0 для гостя, запрос отработает корректно)
	episodes, _ := repository.GetEpisodesBySeriesID(id, userID)

	// Если пришли со страницы закладок (в URL есть favorites=true), фильтруем серии
	showFavorites := r.URL.Query().Get("favorites") == "true"
	if showFavorites {
		var filtered []models.Episode
		for _, ep := range episodes {
			if ep.IsFavorite {
				filtered = append(filtered, ep)
			}
		}
		episodes = filtered
	}

	data := struct {
		Username     string
		Series       interface{}
		Episodes     interface{}
		IsBookmarked bool
		FavoritesOnly bool
	}{
		Username:     username,
		Series:       series,
		Episodes:     episodes,
		IsBookmarked: isBookmarked,
		FavoritesOnly: showFavorites,
	}

	tmpl, _ := template.ParseFiles("web/templates/series.html")
	tmpl.Execute(w, data)
}