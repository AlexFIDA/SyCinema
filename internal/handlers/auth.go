package handlers

import (
	"html/template"
	"net/http"
	"time"

	"sycinema/internal/repository"
	"sycinema/internal/service"
)

func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/register.html")
	tmpl.Execute(w, nil)
}

func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, _ := template.ParseFiles("web/templates/login.html")
	tmpl.Execute(w, nil)
}

func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	err := repository.CreateUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		http.Error(w, "Ошибка регистрации", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	
	// 1. Получаем ID и Username из базы
	id, username, err := repository.AuthenticateUser(r.FormValue("email"), r.FormValue("password"))
	if err != nil {
		http.Error(w, "Неверный логин или пароль", http.StatusUnauthorized)
		return
	}

	// 2. Обязательно передаем полученный ID в токен!
	tokenString, err := service.GenerateToken(id, username)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}