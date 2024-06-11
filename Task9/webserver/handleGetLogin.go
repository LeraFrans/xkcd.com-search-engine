package webserver

import (
	"net/http"
	"text/template"
)

func handleGetLogin(w http.ResponseWriter, r *http.Request) {
		// Проверяем метод запроса, чтобы убедиться, что это GET запрос
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	
		// Создаем шаблон для формы входа
		t := template.Must(template.ParseFiles("./webserver/logform.html"))
		

	
		// Выполняем шаблон и отправляем его в ответ
		if err := t.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
}