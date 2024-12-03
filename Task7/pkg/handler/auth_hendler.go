package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"task7/config"
	"time"

	"golang.org/x/sync/semaphore"
	"golang.org/x/time/rate"
)

// для ограничения общего количества запросов от разных клиентов в одно время
var globalConcurrencyLimiter = semaphore.NewWeighted(int64(config.ReadConfig().Server.Concurrency_limit)) // Ограничиваем до 5 одновременных запросов

type (
	// Структура данных с информацией о пользователе
	User struct {
		Id       int    `json:"-" db:"id"`
		Email    string `json:"email" binding:"required"`
		Name     string `json:"name" binding:"required"`
		password string `json:"password" binding:"required"`
		Role     int    `json:"role" binding:"required"` // 1 - admin, 0 - simple user
	}

	LoginResponse struct {
		Token string `json:"token"`
	}
)

func (h *Handler) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		fmt.Fprintf(w, "Sorry, only POST methods are supported.\n")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные от пользователя
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		fmt.Fprintf(w, "The request could not be decoded\n")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Создаём токен
	tokenString, err := h.services.Authorization.GenerateToken(credentials.Email, credentials.Password)
	if err != nil {
		log.Println(err)
		return
	}

	// Отправляем токен пользователю
	response := LoginResponse{Token: tokenString}
	json.NewEncoder(w).Encode(response)
	w.Header().Set("Authorization", tokenString)
}

// ограничитель запросов от одного клиента
func perClientRateLimiter(next http.HandlerFunc) http.HandlerFunc {

	type client struct {
		limiter  *rate.Limiter // сколько раз может вызвать
		lastSeen time.Time     // время последнего вызова
	}

	var mu sync.Mutex
	var clients = make(map[string]*client) // коллекция клиентов, с которыми работаем сейчас

	go func() {
		// если с момента последнего запроса от клиента прошло больше 5 минут,
		// то удаляем его из коллекции клиентов. И при следующем вызове он будет
		// считаться опять как новый клиент.
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 5*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr) // получение IP-адреса, с которого был запрос, будем ограничивать только его одного
		if err != nil {
			fmt.Fprintf(w, "Failed to get an IP\n")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		mu.Lock()
		// проверям, есть ли клиент с таким ip в нашей коллекции
		if _, found := clients[ip]; !found {
			// если его нет, то делаем ему новый лимитер и добавляем в коллекцию
			clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(config.ReadConfig().Rate_limit_per_second), config.ReadConfig().Rate_limit)} //(2 запроса в сек, максимум 10 от одного клиента за всё время (через 5 мин можно заново))
		}
		clients[ip].lastSeen = time.Now() // обновляем клиенту дату проследнего запроса
		if !clients[ip].limiter.Allow() { // если ограничитель не разрешает больше делать запросы, шлём статус код ошибки
			mu.Unlock()
			fmt.Fprintf(w, "You have exceeded the request limit\n")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		mu.Unlock()

		// Пропускаем запрос через глобальный ограничитель concurrency limiter
		if err := globalConcurrencyLimiter.Acquire(context.Background(), 1); err != nil {
			// Обработка ошибки, если не удалось получить разрешение на выполнение запроса
			fmt.Fprintf(w, "Failed to acquire semaphore: %v\n", err)
			return
		}
		defer globalConcurrencyLimiter.Release(1)

		// Вызываем следующую функцию-обработчик
		next(w, r)
	})
}
