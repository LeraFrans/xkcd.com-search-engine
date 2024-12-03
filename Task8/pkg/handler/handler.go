package handler

import (
	"net/http"
	"task8/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) Init() {

	//все три хендлера пропускаем через функцию - ограничитель запросов от одного клиента
	http.HandleFunc("GET /pics", perClientRateLimiter(http.HandlerFunc(h.handlePics)))
	http.HandleFunc("POST /update", perClientRateLimiter(http.HandlerFunc(h.handleUpdate)))
	http.HandleFunc("POST /login", perClientRateLimiter(http.HandlerFunc(h.handleLogin)))
}
