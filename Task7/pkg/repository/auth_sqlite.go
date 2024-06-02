package repository

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"net/http"
)

// полная информация о юзере из БД (пароль в виде зашифрованного хэша)
type User struct {
	Id       int    `json:"-" db:"id"`
	Email    string `json:"email" binding:"required"`
	Name     string `json:"name" binding:"required"`
	password string `json:"password" binding:"required"`
	Role     int    `json:"role" binding:"required"` // 1 - admin, 0 - simple user
}

// для зависимостей
type AuthSQLite struct {
	db *sql.DB
}

func NewAuthSQLite(db *sql.DB) *AuthSQLite {
	return &AuthSQLite{db: db}
}

// по введёнными пользователем email и паролю ищет юзера в БД и возвращает полный данные о нём (в виде структуры User)
// возвращаемая string это описание ошибки, а int это её статус код, чтобы потом в хендлере вернуть
func (r *AuthSQLite) GetUser(email, password string) (User, string, int) {

	var user User

	// Проверяем корректность предоставленных данных
	if len(email) == 0 || len(password) == 0 {
		return user, "The entered data is incorrect\n", http.StatusBadRequest
	}

	// шифруем пароль
	password_hash := generatePasswordHash(password)

	// Проверяем наличие пользователя в базе данных
	stmt, err := r.db.Prepare("SELECT id, email, name, password, role FROM users WHERE email = ? AND password = ?")
	if err != nil {
		return user, "Error prepare request", http.StatusInternalServerError
	}
	defer stmt.Close()

	err = stmt.QueryRow(email, password_hash).Scan(&user.Id, &user.Email, &user.Name, &user.password, &user.Role)
	switch {
	case err == sql.ErrNoRows:
		return user, "Wrong email or password", http.StatusNotFound
	case err != nil:
		return user, "There is no user with this email", http.StatusInternalServerError
	}

	//если всё ок
	return user, "", 200
}

// по ID пользователя определяет его роль (1 - админ, 0 - обычный юзер)
func (r *AuthSQLite) GetRole(userID int) (int, error) {
	stmt, err := r.db.Prepare("SELECT role FROM users WHERE id = ?")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var role int
	err = stmt.QueryRow(userID).Scan(&role)
	if err != nil {
		return 0, err
	}

	return role, nil
}

// возвращает закодированный хэш пароля
func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte("hjqrhjqw124617ajfhajs")))
}

// вспомогательная функция, тк не реализована регистрация пользователей
func AddUsers() {
	// Подключаемся к БД
	db, errConnect := ConnectDB()
	if errConnect != nil {
		log.Print(errConnect)
		return
	}
	defer db.Close()

	// Добавление обычного пользователя
	user1 := User{
		Email:    "user1@example.com",
		Name:     "User One",
		password: generatePasswordHash("password1"),
		Role:     0,
	}

	// Добавление обычного пользователя
	user2 := User{
		Email:    "user2@example.com",
		Name:     "User Two",
		password: generatePasswordHash("password2"),
		Role:     0,
	}

	// Добавление администратора
	admin1 := User{
		Email:    "admin@example.com",
		Name:     "Administrator",
		password: generatePasswordHash("password3"),
		Role:     1,
	}

	// Добавление пользователей в базу данных
	if _, err := db.Exec("INSERT INTO users (email, name, password, role) VALUES (?, ?, ?, ?)", user1.Email, user1.Name, user1.password, user1.Role); err != nil {
		log.Panicln(err)
	}
	if _, err := db.Exec("INSERT INTO users (email, name, password, role) VALUES (?, ?, ?, ?)", user2.Email, user2.Name, user2.password, user2.Role); err != nil {
		log.Panicln(err)
	}
	if _, err := db.Exec("INSERT INTO users (email, name, password, role) VALUES (?, ?, ?, ?)", admin1.Email, admin1.Name, admin1.password, admin1.Role); err != nil {
		log.Panicln(err)
	}
}
