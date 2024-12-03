package service

import (
	"database/sql"
	"log"
	"os"
	"task8/pkg/repository"
	"testing"
)

func TestAuthService_GenerateToken(t *testing.T) {
	os.Remove("testDB.db")
	testDB, s := TestService() // создаём тестовый сервис, использующий тестовую БД
	defer testDB.Close()

	tests := []struct {
		name     string
		username string
		password string
		want     int
		want1    string
		wantErr  error
	}{
		// TODO: Add test cases.
		{
			name:     "User is not exist",
			username: "name",
			password: "pp",
			want:     404,
			want1:    "",
		},
		{
			name:     "Bad Request",
			username: "",
			password: "",
			want:     400,
			want1:    "",
		},
		{
			name:     "Server Error",
			username: "user1@example.com",
			password: "password1",
			want:     500,
			want1:    "",
		},
	}

	AddDataUsers(testDB)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, _ := s.GenerateToken(tt.username, tt.password)
			if got != tt.want {
				t.Errorf("AuthService.GenerateToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AuthService.GenerateToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
	os.Remove("testDB.db")
}

func AddDataUsers(testDB *sql.DB) {

	// Добавление обычного пользователя
	user1 := repository.User{
		Email:    "user1@example.com",
		Name:     "User One",
		Password: repository.GeneratePasswordHash("password1"),
		Role:     0,
	}

	// Добавление обычного пользователя
	user2 := repository.User{
		Email:    "user2@example.com",
		Name:     "User Two",
		Password: repository.GeneratePasswordHash("password2"),
		Role:     0,
	}

	// Добавление администратора
	admin1 := repository.User{
		Email:    "admin@example.com",
		Name:     "Administrator",
		Password: repository.GeneratePasswordHash("password3"),
		Role:     1,
	}

	// Добавление пользователей в базу данных
	if _, err := testDB.Exec("INSERT INTO users (email, name, password, role) VALUES (?, ?, ?, ?)", user1.Email, user1.Name, user1.Password, user1.Role); err != nil {
		log.Panicln(err)
	}
	if _, err := testDB.Exec("INSERT INTO users (email, name, password, role) VALUES (?, ?, ?, ?)", user2.Email, user2.Name, user2.Password, user2.Role); err != nil {
		log.Panicln(err)
	}
	if _, err := testDB.Exec("INSERT INTO users (email, name, password, role) VALUES (?, ?, ?, ?)", admin1.Email, admin1.Name, admin1.Password, admin1.Role); err != nil {
		log.Panicln(err)
	}
}

func TestAuthService_IsAdmin(t *testing.T) {

	os.Remove("testDB.db")
	testDB, s := TestService() // создаём тестовый сервис, использующий тестовую БД
	defer testDB.Close()

	tests := []struct {
		name        string
		tokenString string
		want        bool
		want1       string
		want2       int
	}{
		// TODO: Add test cases.
		{
			name:        "Empty row",
			tokenString: "",
			want:        false,
			want1:       "error parsing token",
			want2:       401,
		},
		{
			name:        "No user",
			tokenString: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTc5OTk1NzYsImlhdCI6MTcxNzk1NjM3NiwidXNlcl9pZCI6MX0.-Dgb7ynuvitixlIoXGNra5u8vpMVvPfx7ahTh20FVU4",
			want:        false,
			want1:       "error find user in DB",
			want2:       500,
		},

	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := s.IsAdmin(tt.tokenString)
			if got != tt.want {
				t.Errorf("AuthService.IsAdmin() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AuthService.IsAdmin() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("AuthService.IsAdmin() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
	os.Remove("testDB.db")
}
