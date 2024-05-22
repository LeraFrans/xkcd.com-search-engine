package main

import (
	"fmt"
	"log"
	"task6/internal/repository"

	"task6/internal/search"
	"task6/internal/words_processing"
)

func main() {

	db, errConnect := repository.ConnectDB()
	if errConnect != nil {
		log.Print(errConnect)
	}
	defer db.Close()
	repa := repository.NewSQLiteRepository(db) //эту репу везде будем передавать как ссылку на БД

	fmt.Println(search.IndexSearch(words_processing.InputProccessing(), repa))

}
