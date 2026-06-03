package handlers

import (
	"html/template"
	"log"
	"net/http"
)

// HomeHandler обрабатывает запросы к главной странице
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	// Парсим наш HTML файл
	tmpl, err := template.ParseFiles("web/templates/index.html")
	if err != nil {
		log.Println("Ошибка загрузки шаблона:", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}

	// Отправляем сгенерированный HTML в браузер
	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Println("Ошибка рендера шаблона:", err)
	}
}