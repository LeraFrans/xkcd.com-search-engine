package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"task6/internal/repository"
	"task6/internal/search"
	"task6/internal/words_processing"
)



type HTTPHandler struct {
	comicsService repository.ComicsService
}

func NewHTTPHandler(comicsService repository.ComicsService) *HTTPHandler {
	return &HTTPHandler{
		comicsService: comicsService,
	}
}

func (h *HTTPHandler) Init() {
	http.HandleFunc("GET /pics", h.handlePics)
	http.HandleFunc("POST /update", h.handleUpdate)
}

// GET запрос
func (h *HTTPHandler) handlePics(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// обработка входящей строки
	searchQuery := r.URL.Query().Get("search")                                // Получаем строку для поиска
	normSearchQuery, errNormalization := words_processing.WordsNormalizator(searchQuery) // нормализуем и стеммим её
	if errNormalization != nil {
		log.Printf("Error normalizing search query: %v", errNormalization)
		http.Error(w, "Error normalizing search query.", http.StatusBadRequest)
		return
	}

	// Подключаемся к БД
	db, errConnect := repository.ConnectDB()
	if errConnect != nil {
		log.Print(errConnect)
		http.Error(w, "Error connecting to database.", http.StatusInternalServerError)
		return
	}
	defer db.Close()
	repa := repository.NewSQLiteRepository(db) //эту репу везде будем передавать как ссылку на БД

	// Если таблица комиксов пустая, просим обновить базу
	if !repository.IsTableHasAnyRows("keyword", repa) {
		http.Error(w, "Keyword table is exist. Please, do /update request.", http.StatusBadRequest)
	}

	// поиск комиксов
	indexes := search.IndexSearch(normSearchQuery, repa) // ищем релевантные номера комиксов в таблице индексов
	urls := repository.FindURL(indexes, repa)            // по номерам смотрим урлы нужных комиксов

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)     // Возвращаем статус OK
	json.NewEncoder(w).Encode(&urls) // Кодируем структуру в JSON и отправляем клиенту

	// сюда тоже результат отпринтовываем на всякий случай
	fmt.Println("\nYour comics:\n")
	for _, str := range urls {
		fmt.Println(str)
	}
}

// POST запрос
func (h *HTTPHandler) handleUpdate(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fmt.Println("\nI do UPDATE DB...")

	// Включить миграцию, если нет БД
	// database.Migrations()
	// fmt.Println("Finish Migration")

	// Подключаемся к БД
	db, errConnect := repository.ConnectDB()
	if errConnect != nil {
		log.Print(errConnect)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()
	repa := repository.NewSQLiteRepository(db) //эту репу везде будем передавать как ссылку на БД

	// создаём/обновляем таблицу комиксов
	fmt.Println("Start update keyword table...")
	response, errComic := repa.UpdateComicTable()
	if errComic != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(errComic)
	}
	fmt.Println("Finish update keyword table!")

	// создаём/обновляем таблицу индексов
	fmt.Println("Start update index table...")
	errIndex := repa.UpdateIndexTable()
	if errIndex != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Print(errIndex)
	}
	fmt.Println("Finish update index table!")

	// Ответ клиенту
	w.WriteHeader(http.StatusOK)         // Возвращаем статус OK
	json.NewEncoder(w).Encode(&response) // Кодируем структуру в JSON и отправляем клиенту
}
