package repository

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"task8/testDB"
	"testing"
)

func TestAuthSQLite_GetUser(t *testing.T) {

	testDB, err := testDB.CreateTestDB("testDB1.db")
	if err != nil {
		log.Fatalf("failed to initializing db: %s", err.Error())
	}
	r := NewRepository(testDB)

	AddDataUsers(testDB)

	tests := []struct {
		name     string
		email    string
		password string
		want     User
		want1    string
		want2    int
	}{
		{
			name:     "Incorrect input",
			email:    "",
			password: "",
			want:     User{},
			want1:    "The entered data is incorrect\n",
			want2:    400,
		},
		{
			name:     "Wrong password",
			email:    "eee@example.com",
			password: "eee",
			want:     User{},
			want1:    "Wrong email or password",
			want2:    404,
		},
		{
			name:     "OK",
			email:    "user1@example.com",
			password: "password1",
			want:     User{1, "user1@example.com", "User One", "686a7172686a7177313234363137616a6668616a73e38ad214943daad1d64c102faec29de4afe9da3d", 0},
			want1:    "",
			want2:    200,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := r.GetUser(tt.email, tt.password)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AuthSQLite.GetUser() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("AuthSQLite.GetUser() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("AuthSQLite.GetUser() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}

	os.Remove("testDB1.db")
}

func AddDataUsers(testDB *sql.DB) {

	// Добавление обычного пользователя
	user1 := User{
		Email:    "user1@example.com",
		Name:     "User One",
		Password: GeneratePasswordHash("password1"),
		Role:     0,
	}

	// Добавление обычного пользователя
	user2 := User{
		Email:    "user2@example.com",
		Name:     "User Two",
		Password: GeneratePasswordHash("password2"),
		Role:     0,
	}

	// Добавление администратора
	admin1 := User{
		Email:    "admin@example.com",
		Name:     "Administrator",
		Password: GeneratePasswordHash("password3"),
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

func TestAuthSQLite_GetRole(t *testing.T) {

	testDB, err := testDB.CreateTestDB("testDB2.db")
	if err != nil {
		log.Fatalf("failed to initializing db: %s", err.Error())
	}
	r := NewRepository(testDB)

	AddDataUsers(testDB)

	tests := []struct {
		name    string
		userID int
		want    int
	}{
		{
			name: "Simple user",
			userID: 1,
			want:    0,
		},
		{
			name: "Admin",
			userID: 3,
			want:    1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := r.GetRole(tt.userID)
			if got != tt.want {
				t.Errorf("AuthSQLite.GetRole() = %v, want %v", got, tt.want)
			}
		})
	}
}
