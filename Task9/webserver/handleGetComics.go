package webserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
)

func handleComics(w http.ResponseWriter, r *http.Request) {

	// Получаем токен из куки
	cookie, err := r.Cookie("token")
	if err != nil {
		// Если куки нет, возможно, пользователь не авторизован
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Проверяем, что токен не пустая строка
	if len(cookie.Value) < 5 {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Обработка GET запроса на /comics без параметров
	if r.Method == "GET" && r.URL.Query().Get("search") == "" {
		// Возвращаем страницу с формой поиска
		t := template.Must(template.ParseFiles("./webserver/searchform.html"))
		// Выполняем шаблон и отправляем его в ответ
		if err := t.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// Обработка GET запроса на /comics с параметрами
	if r.Method == "GET" && r.URL.Query().Get("search") != "" {

		//получение поисковой строки
		searchString := r.URL.Query().Get("search")

		// Создаем URL с параметрами поиска
		url := fmt.Sprintf("http://localhost:8080/pics?search=%s", searchString)

		// Создаем запрос GET с параметрами поиска
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			log.Println("Ошибка создания запроса:", err)
		}

		// клиент для отправки запроса на поиск
		client := &http.Client{}

		// Отправляем запрос и получаем ответ
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка отправки запроса:", err)
			return
		}
		defer resp.Body.Close()

		// Читаем тело ответа
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Ошибка чтения тела ответа:", err)
			return
		}

		var responseURLs []string // слайс урлов - результатов поиска
		err = json.Unmarshal(body, &responseURLs)
		if err != nil {
			fmt.Println("Ошибка десериализации тела ответа:", err)
			return
		}

		// Загрузка содержимого HTML файла
		TemplatePath := "./webserver/resultform.html"
		searchResultsTemplateContent, err := ioutil.ReadFile(TemplatePath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Создание шаблона
		searchResultsTemplate, err := template.New("searchResults").Parse(string(searchResultsTemplateContent))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Выполнение шаблона и отправка его в ответ
		err = searchResultsTemplate.Execute(w, struct {
			Urls []string
		}{Urls: responseURLs})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
