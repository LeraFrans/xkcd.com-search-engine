package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// GET запрос
func (h *Handler) handlePics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// обработка входящей строки
	searchQuery := r.URL.Query().Get("search")

	// поиск урлов
	urls, errorString, statusCod := h.services.Search.SearchComics(searchQuery)
	if statusCod != 200 {
		http.Error(w, errorString, statusCod)
		return
	}

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)     // Возвращаем статус OK
	json.NewEncoder(w).Encode(&urls) // Кодируем структуру в JSON и отправляем клиенту

	// сюда тоже результат отпринтовываем на всякий случай
	fmt.Println("Your comics:")
	for _, str := range urls {
		fmt.Println(str)
	}
}
