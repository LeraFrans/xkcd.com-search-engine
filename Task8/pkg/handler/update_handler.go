package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

// POST запрос
func (h *Handler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Получения JWT-токена из запроса
	tokenString := r.Header.Get("Authorization")

	// Проверка роли пользователя
	isAdmin, errString, statusCode := h.services.Authorization.IsAdmin(tokenString)
	if statusCode != 200 {
		http.Error(w, errString, statusCode)
	}

	if !isAdmin {
		// если не админ, посылаем сообщение о недосаточных правах и выходим
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Access denied. You must have administrator privileges to perform this action.")
		return
	}

	// Включить миграцию, если нет БД
	// database.Migrations()
	// fmt.Println("Finish Migration")

	fmt.Println("\nI do UPDATE DB...")

	// создаём/обновляем таблицу комиксов + создаём ответ клиенту (количество новых комиксов и общее количество комиксов в базе)
	fmt.Println("Start update keyword table...")
	response, errComic := h.services.Update.UpdateComicTable()
	if errComic != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(errComic)
	}
	fmt.Println("Finish update keyword table!")

	// создаём/обновляем таблицу индексов
	fmt.Println("Start update index table...")
	errString, errIndex := h.services.Update.UpdateIndexTable()
	if errIndex != nil {
		http.Error(w, errString, http.StatusInternalServerError)
		log.Print(errIndex)
	}
	fmt.Println("Finish update index table!")

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)         // Возвращаем статус OK
	json.NewEncoder(w).Encode(&response) // Кодируем структуру в JSON и отправляем клиенту
}
