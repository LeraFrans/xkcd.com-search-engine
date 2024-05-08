package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"task5/comic"
	"task5/words"

	"gopkg.in/yaml.v2"
)

type ComicsCount struct {
	NewComics   int `json:"new_comics"`
	TotalComics int `json:"total_comics"`
}

type URL struct {
	Url []string `json:"url"`
}

func main() {

	http.HandleFunc("/update", handleUpdate)
	http.HandleFunc("/pics", handlePics)

	log.Println("Listening...")

	_, _, _, portInt := readConfig() // читаем порт из конфига
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

	w.WriteHeader(http.StatusOK)                                                          // Возвращаем статус OK
	source_url, db_file, parallel, _ := readConfig()                                      // читаем нужные перменные из конфига
	newComics, totalComics := comic.MakeJSONWithComicsData(source_url, db_file, parallel) // создаём/обновляем БД
	cc := ComicsCount{newComics, totalComics}
	json.NewEncoder(w).Encode(&cc) // Кодируем структуру в JSON и отправляем клиенту
}

func handlePics(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		fmt.Fprintf(w, "Sorry, only GET methods are supported.")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	searchQuery := r.URL.Query().Get("search")                 // Получаем строку поиска из URL.
	normSearchQuery, _ := words.WordsNormalizator(searchQuery) // нормализуем и стеммим её

	w.WriteHeader(http.StatusOK)                 // Возвращаем статус OK
	u := URL{comic.IndexSearch(normSearchQuery)} // в структуру записываем результат поиска
	json.NewEncoder(w).Encode(&u)                // Кодируем структуру в JSON и отправляем клиенту

	// сюда тоже результат отпринтовываем на всякий случай
	for _, str := range comic.IndexSearch(normSearchQuery) {
		fmt.Println(str)
	}
}

// Возвращает переменные из "config.yaml"
func readConfig() (string, string, int, int) {
	// Чтение файла

	configName := "config.yaml"

	content, err := os.ReadFile(configName)
	if err != nil {
		panic(err)
	}

	// Разбор YAML-контента
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(content, &yamlMap)
	if err != nil {
		panic(err)
	}

	source_url, ok := yamlMap["source_url"].(string)
	if !ok {
		panic("Неверный тип данных для 'source_url'")
	}
	db_file, ok := yamlMap["db_file"].(string)
	if !ok {
		panic("Неверный тип данных для 'db_file'")
	}
	parallel, ok := yamlMap["parallel"].(int)
	if !ok {
		panic("Неверный тип данных для 'parallel'")
	}
	port, ok := yamlMap["port"].(int)
	if !ok {
		panic("Неверный тип данных для 'port'")
	}

	return source_url, db_file, parallel, port
}
