package main

import (
	"log"
	"net/http"
	"task9/webserver"
)

func main (){
	// Инициализируем обработчики
	webserver.Init()
	log.Println("Listening...")

	// Запускаем сервер
	log.Fatal(http.ListenAndServe(":8081", nil))
}