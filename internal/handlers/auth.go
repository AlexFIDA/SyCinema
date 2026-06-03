package handlers

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"sycinema/internal/repository" // Импортируем наш репозиторий
)

// RegisterPageHandler просто отдает HTML страницу с формой (метод GET)
func RegisterPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/register.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// RegisterPostHandler принимает данные из формы и сохраняет юзера (метод POST)
func RegisterPostHandler(w http.ResponseWriter, r *http.Request) {
	// 1. Просим Go распарсить данные из HTML-формы
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка чтения данных формы", http.StatusBadRequest)
		return
	}

	// 2. Достаем значения по их атрибуту name="" из HTML
	username := r.FormValue("username")
	email := r.FormValue("email")
	password := r.FormValue("password")

	// 3. Вызываем нашу функцию из репозитория
	err = repository.CreateUser(username, email, password)
	if err != nil {
		log.Println("Ошибка регистрации:", err)
		http.Error(w, "Ошибка при регистрации (возможно, email уже занят)", http.StatusInternalServerError)
		return
	}

	// 4. Если всё успешно, отправляем простой текст (позже заменим на редирект)
	fmt.Fprintf(w, "✅ Пользователь %s успешно зарегистрирован!", username)
}

// LoginPageHandler отдает HTML форму логина
func LoginPageHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/login.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки страницы", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

// LoginPostHandler обрабатывает попытку входа
func LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка чтения формы", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")

	// Проверяем пользователя через нашу новую функцию БД
	username, err := repository.AuthenticateUser(email, password)
	if err != nil {
		// Если пароль не подошел — ругаемся
		http.Error(w, "❌ Неверный email или пароль", http.StatusUnauthorized)
		return
	}

	// Если авторизация успешна:
	// Позже мы будем выдавать здесь JWT-токен (куки), чтобы сайт "запомнил" пользователя.
	// А пока просто радуемся успеху!
	fmt.Fprintf(w, "✅ Добро пожаловать обратно, %s!", username)
}