package repository

import (
	"database/sql"
)

type Authorization interface {
	GetUser(email, password string) (User, string, int)
	GetRole(userID int) (int, error)
}

type Update interface {
	FindOurMaxNumberOfComics() (int, error)
	WriteResultInKeywordTable(resultComicsSlice []Comic) error
	GetDataForIndexTable() ([]Comic, error)
	IsTableHasAnyRows(table_name string) bool
	WriteResultInIndexTable(index map[string][]string) error
}

type Search interface {
	FindInIndex(input_keywords []string) ([]IndexFromDB, error)
	FindURL(indexes []IndexResult) ([]string, error)
	IsTableHasAnyRows(table_name string) bool
}

type Repository struct {
	Authorization
	Update
	Search
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		Authorization: NewAuthSQLite(db),
		Update: NewUpdateSQLite(db),
		Search: NewSearchSQLite(db),
	}
}
