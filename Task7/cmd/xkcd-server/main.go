package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"task7/config"
	"task7/pkg/handler"
	"task7/pkg/repository"
	"task7/pkg/service"
)

func main() {

	var wg sync.WaitGroup

	db, err := repository.ConnectDB()
	if err != nil {
		log.Fatalf("failed to initializing db: %s", err.Error())
	}
	defer db.Close()

	// создание трёх слоёв: репозитория(работа с БД), сервиса(бизнес-логика) и хендлеров(связь с клиентом)
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	hand := handler.NewHandler(services)

	wg.Add(1)
	go func() {
		defer wg.Done()
		// Инициализируем обработчики
		hand.Init()
		log.Println("Listening...")

		// определяем порт
		portInt := config.ReadConfig().Port
		port := fmt.Sprintf(":%d", portInt)

		// Запускаем сервер
		log.Fatal(http.ListenAndServe(port, nil))
	}()

	wg.Wait()

}
