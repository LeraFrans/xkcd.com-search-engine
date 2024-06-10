package service

import (
	"errors"
	"fmt"
	"net/http"
	"task8/pkg/repository"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// для связи между слоями
type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

const (
	salt       = "hjqrhjqw124617ajfhajs"
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

// создаёт токен для аутентифицирующегося пользователя
func (s *AuthService) GenerateToken(username, password string) (int, string, error) {

	//получаем всю инфу о юзере из БД
	user, errString, status := s.repo.GetUser(username, password)
	if status != 200 {
		return status, "", errors.New(errString)
	}

	// генерируем токен
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	// возвращаем токен в виде зашифрованной строки
	rezToken, err := token.SignedString([]byte(signingKey))
	return 200, rezToken, err
}

// парсит токен, потом по ID проверяет в бд, является ли юзер админом
func (s *AuthService) IsAdmin (tokenString string) (bool, string, int) {
	fmt.Printf("\nTOKEN STRING: %s\n", tokenString)

	//получаем ID пользователя по токену
	userID, err := s.ParseToken(tokenString)
	if err != nil {
		return false, "error parsing token", http.StatusUnauthorized
	}
	fmt.Printf("USER ID: %d\n", userID)

	// идёт в бд и возвращает роль пользователя
	role, err := s.repo.GetRole(userID)
	if err != nil {
		return false, "error find user in DB", http.StatusInternalServerError
	}

	return role == 1, "", http.StatusOK
}

// по токену получает ID пользователя 
func (s *AuthService) ParseToken(accessToken string) (int, error) {

	// колбэк для проверки подписи токена
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем, что метод подписи является HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		// Возвращаем секретный ключ для проверки подписи.
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}

	// преобразуем в нужную нам структуру
	claims, _ := token.Claims.(*tokenClaims)

	// возвращаем ID пользователя
	return claims.UserId, nil
}