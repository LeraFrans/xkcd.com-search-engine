package server

import (
	"fmt"
	"log"
	"time"
	"task6/internal/repository"

	"github.com/go-co-op/gocron"
)

func RunCron() {
	// настраиваем местное время
	localTime, err := time.LoadLocation("Europe/Moscow")
	fmt.Println(localTime)
	if err != nil {
		log.Println(err)
	}

	// инициализируем объект планировщика
	s := gocron.NewScheduler(localTime)

	// выполняем задачу раз в день в определённое время
	_, err = s.Every(1).Day().At("03:32").Do(func() {

		//Тут почти полная копия того, что было в хендлере обновления

		fmt.Println("\nThe database is currently being updated daily")
		// Подключаемся к БД
		db, errConnect := repository.ConnectDB()
		if errConnect != nil {
			log.Print(errConnect)
			return
		}
		defer db.Close()
		repa := repository.NewSQLiteRepository(db) //эту репу везде будем передавать как ссылку на БД

		// создаём/обновляем таблицу комиксов
		fmt.Println("Start update keyword table...")
		response, errComic := repa.UpdateComicTable()
		if errComic != nil {
			log.Print(errComic)
		}
		fmt.Printf("Finish update keyword table!\nNew comics: %d	Total comics: %d\n", response.New, response.Total)

		// создаём/обновляем таблицу индексов
		fmt.Println("Start update index table...")
		errIndex := repa.UpdateIndexTable()
		if errIndex != nil {
			log.Print(errIndex)
		}
		fmt.Println("Finish update index table!")
		fmt.Println("\n Update is complete! ")

	})
	if err != nil {
		log.Println(err)
	}

	// запускаем планировщик с блокировкой текущего потока
	s.StartBlocking()
}