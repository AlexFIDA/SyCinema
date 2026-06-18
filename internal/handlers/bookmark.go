package handlers

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"sycinema/internal/repository"
	"sycinema/internal/service"
)

func ToggleBookmarkHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("⚠️ ToggleBookmark: Кука token не найдена:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	claims, err := service.ParseToken(cookie.Value)
	if err != nil {
		log.Println("⚠️ ToggleBookmark: Ошибка парсинга токена:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	if claims.UserID == 0 {
		log.Println("⚠️ ToggleBookmark: UserID в токене равен 0!")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	seriesID, _ := strconv.Atoi(r.FormValue("series_id"))
	repository.ToggleBookmark(claims.UserID, seriesID)

	http.Redirect(w, r, "/series?id="+r.FormValue("series_id"), http.StatusSeeOther)
}

func BookmarksPageHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("⚠️ BookmarksPage: Кука token не найдена")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	claims, err := service.ParseToken(cookie.Value)
	if err != nil {
		log.Println("⚠️ BookmarksPage: Ошибка токена:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	bookmarks, err := repository.GetBookmarks(claims.UserID)
	if err != nil {
		log.Println("❌ BookmarksPage: Ошибка БД:", err)
		http.Error(w, "Ошибка БД", http.StatusInternalServerError)
		return
	}

	data := struct {
		Username string
		Series   interface{}
	}{
		Username: claims.Username,
		Series:   bookmarks,
	}

	tmpl, err := template.ParseFiles("web/templates/bookmarks.html")
	if err != nil {
		log.Println("❌ BookmarksPage: Ошибка загрузки файла шаблона:", err)
		http.Error(w, "Шаблон не найден", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		log.Println("❌ BookmarksPage: Ошибка рендеринга:", err)
	}
}

func ToggleFavoriteEpisodeHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Println("⚠️ ToggleFavorite: Кука token не найдена:", err)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	claims, err := service.ParseToken(cookie.Value)
	if err != nil || claims.UserID == 0 {
		log.Println("⚠️ ToggleFavorite: Ошибка токена или UserID=0")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	episodeID, _ := strconv.Atoi(r.FormValue("episode_id"))
	seriesID := r.FormValue("series_id")
	favoritesFilter := r.FormValue("favorites")

	err = repository.ToggleFavoriteEpisode(claims.UserID, episodeID)
	if err != nil {
		log.Println("❌ ToggleFavorite: Ошибка БД:", err)
	}

	redirectURL := "/series?id=" + seriesID
	if favoritesFilter == "true" {
		redirectURL += "&favorites=true"
	}
	
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}