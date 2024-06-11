package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"task9/pkg/repository"
	"task9/pkg/service"
	"task9/testDB"
)

func TestHandler_handlePics(t *testing.T) {
	// здесь тестовая БД, сервис и хэндлер
	os.Remove("testDB2.db")
	testDB, _ := testDB.CreateTestDB("testDB2.db")

	r := repository.NewRepository(testDB)
	s := service.NewService(r)
	h := NewHandler(s)

	addDataComic(testDB)
	s.UpdateIndexTable()

	tests := []struct {
		name           string
		searchString   string
		method         string
		wantStatusCode int
		wantBodyLen    int
	}{
		{
			name:           "OK",
			searchString:   "sun?river?float?into",
			method:         "GET",
			wantStatusCode: 200,
			wantBodyLen:    2,
		},

		{
			name:           "No rows",
			searchString:   "call?trying?but?it?says?number?blocked",
			method:         "GET",
			wantStatusCode: 500,
			wantBodyLen:    1,
		},
		{
			name:           "POST method",
			searchString:   "sun?river?float?into",
			method:         "POST",
			wantStatusCode: 200,
			wantBodyLen:    2,
		},
	}

	//s.UpdateIndexTable()

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			// Создаем URL с параметрами поиска
			url := fmt.Sprintf("http://localhost:8080/pics?search=%s", tt.searchString)

			// Создаем запрос GET с параметрами поиска
			req, err := http.NewRequest(tt.method, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			// Мы создаем ResponseRecorder(реализует интерфейс http.ResponseWriter)
			// и используем его для получения ответа
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(h.handlePics)
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

func addDataComic(testDB *sql.DB) {

	// заполняем таблицу комиксов
	if _, err := testDB.Exec("INSERT INTO keyword (number, url, keywords) VALUES (?, ?, ?)", 1, "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg", "boy barrel is drift ocean next noth seen t els don all sit float wonder where into distanc"); err != nil {
		log.Print(err)
	}
	if _, err := testDB.Exec("INSERT INTO keyword (number, url, keywords) VALUES (?, ?, ?)", 2, "https://imgs.xkcd.com/comics/tree_cropped_(1).jpg", "sphere refer onli halfway le princ thought sketch two opposit side titl grow be about tree petit through"); err != nil {
		log.Print(err)
	}
	if _, err := testDB.Exec("INSERT INTO keyword (number, url, keywords) VALUES (?, ?, ?)", 3, "https://imgs.xkcd.com/comics/island_color.jpg", "sketch island hello"); err != nil {
		log.Print(err)
	}
	if _, err := testDB.Exec("INSERT INTO keyword (number, url, keywords) VALUES (?, ?, ?)", 4, "https://imgs.xkcd.com/comics/landscape_cropped_(1).jpg", "flow through ocean sketch landscap sun horizon river"); err != nil {
		log.Print(err)
	}
}
