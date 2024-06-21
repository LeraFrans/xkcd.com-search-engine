package repository

import (
	"database/sql"
	"log"
)

// для связи между слоями
type SearchSQLite struct {
	db *sql.DB
}

func NewSearchSQLite(db *sql.DB) *SearchSQLite {
	return &SearchSQLite{db: db}
}

// информация об одном индексе, взятая из БД
type IndexFromDB struct {
	Word    string
	Numbers string
}

// Для функции IndexSearch(). Ищет в таблице индексов нужные нам релевантные слова
func (r *SearchSQLite) FindInIndex(input_keywords []string) ([]IndexFromDB, error) {

	var result []IndexFromDB

	// Я правда долго пыталась, используя плейсхолдеры и преобразования типов, засунуть это в один запрос,
	// но ничего не получилось, поэтому вот отдельный запрос для каждого ключевого слова (какой кошмар)
	for _, word := range input_keywords {

		// запрос
		row := r.db.QueryRow("SELECT word, numbers FROM words_index WHERE word = (?)", word)

		// Считывание индекса из курсора
		var y IndexFromDB
		err := row.Scan(&y.Word, &y.Numbers)
		if err != nil {
			return nil, err
		}

		// Добавление индекса в слайс
		result = append(result, y)
	}

	return result, nil
}

// по номерам комиксом ищем в бд их урлы
func (r *SearchSQLite) FindURL(indexes []IndexResult) ([]string, error) {
	result := []string{}
	for _, index := range indexes {
		row := r.db.QueryRow("SELECT url FROM keyword WHERE number = ?", index.Num)

		// Считывание урла из курсора
		var url string
		err := row.Scan(&url)
		if err != nil {
			return nil, err
		}

		// Добавление урла в слайс
		result = append(result, url)
	}

	return result, nil
}

// Содержит ли таблица хоть одну запись (для проверки перед поиском)
func (r *SearchSQLite) IsTableHasAnyRows(table_name string) bool {

	// Проверка наличия хотя бы одной записи в таблице keyword
	stmt2, err := r.db.Prepare("SELECT COUNT(*) FROM keyword")
	if err != nil {
		log.Println(err)
	}
	defer stmt2.Close()

	// Выполнение запроса и проверка результата
	row := stmt2.QueryRow()
	var count int
	if err := row.Scan(&count); err != nil {
		log.Println(err)
	}

	// Если count == 0, значит таблица пустая
	return count > 0
}
