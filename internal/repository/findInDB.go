package repository

import (
	"log"
)

// Для функции IndexSearch(). Ищет в таблице индексов нужные нам релевантные слова
func FindInIndex(input_keywords []string, repa *SQLiteRepository) []IndexFromDB {


	var result []IndexFromDB

	// Я правда долго пыталась, используя плейсхолдеры и преобразования типов, засунуть это в один запрос,
	// но ничего не получилось, поэтому вот отдельный запрос для каждого ключевого слова (какой кошмар)
	for _, word := range input_keywords {

		// запрос
		row := repa.db.QueryRow("SELECT word, numbers FROM words_index WHERE word = (?)", word)

		// Считывание индекса из курсора
		var r IndexFromDB
		err := row.Scan(&r.Word, &r.Numbers)
		if err != nil {
			log.Println(err)
		}

		// Добавление индекса в слайс
		result = append(result, r)
	}

	return result
}

// По нoмерам комиксов ищем в БД их урлы
func FindURL(indexes []IndexResult, repa *SQLiteRepository) []string {

	result := []string{}
	for _, index := range indexes {
		row := repa.db.QueryRow("SELECT url FROM keyword WHERE number = (?)", index.Num)

		// Считывание урла из курсора
		var url string
		err := row.Scan(&url)
		if err != nil {
			log.Println(err)
		}

		// Добавление урла в слайс
		result = append(result, url)
	}

	return result
}
