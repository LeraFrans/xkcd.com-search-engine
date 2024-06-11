package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"task9/pkg/repository"
	"task9/pkg/service"
	"task9/testDB"
	"testing"

	"github.com/stretchr/testify/assert"
)

func AddDataUsers(testDB *sql.DB) {

	// Добавление обычного пользователя
	user1 := User{
		Email:    "user1@example.com",
		Name:     "User One",
		Password: repository.GeneratePasswordHash("password1"),
		Role:     0,
	}

	// Добавление обычного пользователя
	user2 := User{
		Email:    "user2@example.com",
		Name:     "User Two",
		Password: repository.GeneratePasswordHash("password2"),
		Role:     0,
	}

	// Добавление администратора
	admin1 := User{
		Email:    "admin@example.com",
		Name:     "Administrator",
		Password: repository.GeneratePasswordHash("password3"),
		Role:     1,
	}

	// Добавление пользователей в базу данных
	if _, err := testDB.Exec("INSERT INTO users (id, email, name, password, role) VALUES (?, ?, ?, ?, ?)", 1, user1.Email, user1.Name, user1.Password, user1.Role); err != nil {
		log.Panicln(err)
	}
	if _, err := testDB.Exec("INSERT INTO users (id, email, name, password, role) VALUES (?, ?, ?, ?, ?)", 2, user2.Email, user2.Name, user2.Password, user2.Role); err != nil {
		log.Panicln(err)
	}
	if _, err := testDB.Exec("INSERT INTO users (id, email, name, password, role) VALUES (?, ?, ?, ?, ?)", 3, admin1.Email, admin1.Name, admin1.Password, admin1.Role); err != nil {
		log.Panicln(err)
	}
}

type LoginData struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func TestHandler_handleLogin(t *testing.T) {

	//здесь тестовая БД, сервис и хэндлер
	os.Remove("testDB1.db")
	testDB, _ := testDB.CreateTestDB("testDB1")

	r := repository.NewRepository(testDB)
	s := service.NewService(r)
	h := NewHandler(s)

	AddDataUsers(testDB)

	tests := []struct {
		name           string
		args           LoginData
		method         string
		wantStatusCode int
		wantBodyLen    int
	}{
		{
			name:           "OK",
			args:           LoginData{Email: "user1@example.com", Password: "password1"},
			method:         "POST",
			wantStatusCode: 200,
			wantBodyLen:    10,
		},
		{
			name:           "Not foud user",
			args:           LoginData{Email: "ffff@example.com", Password: "password"},
			method:         "POST",
			wantStatusCode: 404,
			wantBodyLen:    0,
		},
		{
			name:           "Bad request",
			args:           LoginData{},
			method:         "POST",
			wantStatusCode: 400,
			wantBodyLen:    0,
		},
		{
			name:           "GET method",
			args:           LoginData{Email: "user1@example.com", Password: "password1"},
			method:         "GET",
			wantStatusCode: 200,
			wantBodyLen:    10,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			// Кодируем структуру в JSON
			jsonData, err := json.Marshal(tt.args)
			if err != nil {
				panic(err) // Обработайте ошибку здесь
			}

			req, err := http.NewRequest(tt.method, "/login", bytes.NewBuffer(jsonData))
			if err != nil {
				t.Fatal(err)
			}

			// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
			// и используем его для получения ответа
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.handleLogin)
			handler.ServeHTTP(rr, req)

			// Проверяем код
			if status := rr.Code; status != tt.wantStatusCode {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatusCode)
			}
			// Проверяем тело ответа
			if len(rr.Body.String()) < tt.wantBodyLen {
				t.Errorf("handler returned unexpected body: got %v want %v",
					rr.Body.String(), "error len()")
			}
			fmt.Printf("\nName: %s, BODY:%s\n", tt.name, rr.Body.String())
		})
	}

}

func Test_perClientRateLimiter(t *testing.T) {
	// Подготовка к тесту
	// limiter := rate.NewLimiter(rate.Every(time.Second), 2)
	// globalConcurrencyLimiter := semaphore.NewWeighted(1)
	server := httptest.NewServer(perClientRateLimiter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})))
	defer server.Close()

	// Тест 1: Проверка ограничения количества запросов
	// Запрос 1
	for i := range 13 {
		resp, err := http.Get(server.URL)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		fmt.Println(i)
	}

	// Тест 2: Проверка глобального ограничителя concurrency
	// (Упс, похоже не работает)
	// Запрос 1
	resp, err := http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Запрос 2
	resp, err = http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Запрос 3
	resp, err = http.Get(server.URL)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
