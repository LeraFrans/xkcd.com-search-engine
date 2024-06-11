package testDB

import (
	"database/sql"
)

func CreateTestDB(nameDB string) (*sql.DB, error) {
	// Создание новой базы данных
	db, err := sql.Open("sqlite3", nameDB)
	if err != nil {
		return nil, err
	}

	// Создание таблицы users
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS users (Id INTEGER, Email TEXT, Name TEXT, password TEXT, Role INTEGER)")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	// Создание таблицы keyword
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS keyword (number INTEGER, url TEXT, keywords TEXT)")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	// Создание таблицы words_index
	stmt, err = db.Prepare("CREATE TABLE IF NOT EXISTS words_index (word TEXT, numbers TEXT)")
	if err != nil {
		return nil, err
	}
	_, err = stmt.Exec()
	if err != nil {
		return nil, err
	}

	return db, nil
}
