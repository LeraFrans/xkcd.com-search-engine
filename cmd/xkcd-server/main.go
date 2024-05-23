package main

import (
	"fmt"
	"log"
	"net/http"
	"task6/config"
	"task6/internal/server"
)

func main() {
	// Создаем экземпляр HTTPHandler
	handler := server.NewHTTPHandler(nil) // Здесь должен быть реализован ComicService

	// Инициализируем обработчики
	handler.Init()
	log.Println("Listening...")

	// GoCron для ежедневного обновления БД
	go server.RunCron()

	// определяем порт
	portInt := config.ReadConfig().Port
	port := fmt.Sprintf(":%d", portInt)

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(port, nil))
}
