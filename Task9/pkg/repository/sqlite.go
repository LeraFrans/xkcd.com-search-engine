package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"task9/config"

	// Библиотека для миграций
	"github.com/golang-migrate/migrate/v4"
	// Драйвер для выполнения миграций SQLite 3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// Драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrations() {
	var storagePath, migrationsPath, migrationsTable string

	storagePath = "pkg/repository/database/xkcdDB.db"
	migrationsPath = "schema"
	migrationsTable = "migration_table"

	// Создаем объект мигратора, передав креды нашей БД
	m, err := migrate.New("file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		log.Println(err)
	}

	// Выполняем миграции до последней версии
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		log.Println(err)
	}
}

// Создание подключения к базе данных (вызывается в хендлерах)
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.ReadConfig().Dsn)
	if err != nil {
		log.Printf("Не удалось открыть базу данных: %v", err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
