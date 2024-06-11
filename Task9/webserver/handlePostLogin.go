package webserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func handlePostLogin(w http.ResponseWriter, r *http.Request) {

	// Проверяем, что запрос является POST-запросом
	if r.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ошибка при чтении тела запроса: %v", err)
		return
	}

	// Парсим тело запроса
	pairs := strings.Split(string(body), "&")
	// Парсим каждую пару
	var email, password string
	for _, pair := range pairs {
		kv := strings.Split(pair, "=")
		key := kv[0]
		value := kv[1]

		switch key {
		case "email":
			email = strings.ReplaceAll(value, "%40", "@") // тут почему-то %40 вместо собачки прилетает, поэтому меняем
		case "password":
			password = value
		}
	}

	// Формируем структуру для отправки
	credentials := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}

	// Кодируем структуру в JSON
	jsonData, err := json.Marshal(credentials)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ошибка при кодировании JSON: %v", err)
		return
	}

	// Отправляем POST запрос на другой сервер
	resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Ошибка при отправке запроса: %v", err)
		return
	}

	bb, _ := ioutil.ReadAll(resp.Body) // Читаем тело ответа в байты
	responseString := string(bb)       // Преобразуем байты в строку

	// Обрабатываем ответ
	defer resp.Body.Close()

	// получаем токен в виде строки
	// Разделяем строку на части до первого двоеточия
	var token2 string
	parts := strings.Split(responseString, ":")
	if len(parts) > 1 {
		// Извлекаем часть после двоеточия и удаляем кавычки
		token1 := parts[1][1 : len(parts[1])-1]
		token2 = token1[:len(token1)-2]
	}

	cookie := http.Cookie{
		Name:     "token",
		Value:    token2,
		Expires:  time.Now().Add(24 * time.Hour), // Токен действителен в течение 24 часов
		HttpOnly: true,                           // Защищает от CSRF-атак
	}
	http.SetCookie(w, &cookie)

	// Возвращаем страницу с формой поиска
	http.Redirect(w, r, "/comics", http.StatusFound)
	fmt.Fprintf(w, "Ответ: %s", resp.Status) // Установка куки с токеном
}
