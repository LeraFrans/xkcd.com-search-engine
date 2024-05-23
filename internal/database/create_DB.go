package database

import (
	"errors"
	"fmt"
	"log"

	// Библиотека для миграций
	"github.com/golang-migrate/migrate/v4"
	// Драйвер для выполнения миграций SQLite 3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	// Драйвер для получения миграций из файлов
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrations() {
	var storagePath, migrationsPath, migrationsTable string

	storagePath = "internal/database/sqlite_db/xkcdDB.db"
	migrationsPath = "internal/database/sqlite_db/"
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
