package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"task9/pkg/repository"
	"task9/pkg/service"
	"task9/testDB"
)

func TestHandleUpdate(t *testing.T) {

	// здесь тестовая БД, сервис и хэндлер
	os.Remove("testDB3.db")
	testDB, _ := testDB.CreateTestDB("testDB3.db")

	r := repository.NewRepository(testDB)
	s := service.NewService(r)
	h := NewHandler(s)

	addDataComic(testDB)
	s.UpdateIndexTable()
	AddDataUsers(testDB)

	t.Run("OK", func(t *testing.T) {
		// Создаем фиктивный запрос
		req := httptest.NewRequest("GET", "/update", nil)
		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(h.handleUpdate)
		handler.ServeHTTP(recorder, req)

		assert.Equal(t, http.StatusOK, recorder.Code)
	})

	// Проверяем, что был отправлен статус Unauthorized
	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/update", nil)
		recorder := httptest.NewRecorder()
		handler := http.HandlerFunc(h.handleUpdate)
		handler.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusUnauthorized, recorder.Code)
	})

	// Проверяем, что был отправлен статус OK
	t.Run("Unauthorized", func(t *testing.T) {

		// получаем токен админа (внутри функции происходит авторизация)
		token := h.getAdminToken()
		fmt.Printf("OUR TOKEN IS: %s\n", token)

		req := httptest.NewRequest("POST", "/update", nil)
		recorder := httptest.NewRecorder()

		//передаём токен в заголовок как значение ключа "Authorization"
		//но почему-то не работает
		req.Header.Set("Authorization", token)

		handler := http.HandlerFunc(h.handleUpdate)
		handler.ServeHTTP(recorder, req)
		assert.Equal(t, http.StatusOK, recorder.Code)
	})
}

func (h *Handler) getAdminToken() string {
	// Кодируем структуру в JSON
	jsonData, err := json.Marshal(LoginData{Email: "admin@example.com", Password: "password3"})
	if err != nil {
		panic(err) // Обработайте ошибку здесь
	}

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonData))

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(h.handleLogin)
	handler.ServeHTTP(rr, req)

	jsonStr := rr.Body.String()

	type TokenResponse struct {
		Token string `json:"token"`
	}
	var tokenResponse TokenResponse
	// Десериализация JSON в структуру
	err = json.Unmarshal([]byte(jsonStr), &tokenResponse)
	if err != nil {
		fmt.Println("Ошибка десериализации:", err)
		return ""
	}

	// Извлечение токена из структуры
	token := tokenResponse.Token
	//fmt.Println("Токен:", token) // Выведет только токен без фигурных скобок и названия ключа

	return token

}
