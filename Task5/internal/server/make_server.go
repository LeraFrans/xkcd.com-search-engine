package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"task5/config"
	"task5/internal/comic"
	"task5/internal/words"
	"time"

	"github.com/go-co-op/gocron"
)

type ComicsCount struct {
	NewComics   int `json:"new_comics"`
	TotalComics int `json:"total_comics"`
}

type URL struct {
	Url []string `json:"url"`
}

func GoServer() {

	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/pics", handlePics)

	log.Println("Listening...")

	// GoCron для ежедневного обновления БД
	go runCron()

	// определяем порт
	portInt := config.ReadConfig().Port
	port := fmt.Sprintf(":%d", portInt)

	log.Fatal(http.ListenAndServe(port, nil))
}

func handleUpdate(w http.ResponseWriter, r *http.Request) {

	if r.Method != "POST" {
		fmt.Fprintf(w, "Sorry, only POST methods are supported.")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("I do UPDATE")

	w.WriteHeader(http.StatusOK)                             // Возвращаем статус OK
	newComics, totalComics := comic.MakeJSONWithComicsData() // создаём/обновляем БД
	cc := ComicsCount{newComics, totalComics}
	json.NewEncoder(w).Encode(&cc) // Кодируем структуру в JSON и отправляем клиенту

	// сюда тоже результат отпринтовываем на всякий случай
	fmt.Printf("\n Update is complete!\nNew comics: %d\nTotal comics: %d\n", cc.NewComics, cc.TotalComics)

}

func handlePics(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	searchQuery := r.URL.Query().Get("search")                 // Получаем строку поиска из URL.
	normSearchQuery, _ := words.WordsNormalizator(searchQuery) // нормализуем и стеммим её

	if !comic.IsDataBaseExist() {
		// Если база данных не существует, возвращаем ошибку 500.
		http.Error(w, "Database not found. Please, do update:\n> make test_post", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)                 // Возвращаем статус OK
	u := URL{comic.IndexSearch(normSearchQuery)} // в структуру записываем результат поиска
	json.NewEncoder(w).Encode(&u)                // Кодируем структуру в JSON и отправляем клиенту

	// сюда тоже результат отпринтовываем на всякий случай
	for _, str := range comic.IndexSearch(normSearchQuery) {
		fmt.Println(str)
	}
}

func runCron() {
	// настраиваем местное время
	localTime, err := time.LoadLocation("Europe/Moscow")
	fmt.Println(localTime)
	if err != nil {
		log.Println(err)
	}

	// инициализируем объект планировщика
	s := gocron.NewScheduler(localTime)

	// выполняем задачу раз в день в определённое время
	_, err = s.Every(1).Day().At("03:32").Do(func() {
		fmt.Println("\nThe database is currently being updated daily")
		newComics, totalComics := comic.MakeJSONWithComicsData() // создаём/обновляем БД
		fmt.Printf("\n Update is complete!\nNew comics: %d\nTotal comics: %d\n", newComics, totalComics)
	})
	if err != nil {
		log.Println(err)
	}

	// запускаем планировщик с блокировкой текущего потока
	s.StartBlocking()
}
