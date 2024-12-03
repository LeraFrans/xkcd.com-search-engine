package service

import (
	"database/sql"
	"log"
	"task8/pkg/repository"
	"task8/testDB"
)

func TestService() (*sql.DB, *Service) {
	testDB, err := testDB.CreateTestDB("testDB.db")
	if err != nil {
		log.Fatalf("failed to initializing db: %s", err.Error())
	}
	repos := repository.NewRepository(testDB)
	services := NewService(repos)

	return testDB, services
}
