package repository

import (
	"fmt"
	"log"
	"strings"
)

// для POST запроса, обновляет таблицу инвертированных индексов по таблице комиксов
func (repa *SQLiteRepository) UpdateIndexTable() error {

	// Запрос нужных данных
	rows, err := repa.db.Query("SELECT number, keywords FROM keyword")
	if err != nil {
		return err
	}
	defer rows.Close()

	// Итерация по результатам запроса и заполнение слайса комиксов
	comics := []Comic{}
	for rows.Next() {
		c := Comic{} // один комикс
		err := rows.Scan(&c.Num, &c.Keywords)
		if err != nil {
			log.Println(err)
			continue
		}
		comics = append(comics, c)
	}

	// Создание инвертированного индекса
	var index = make(map[string][]string)
	for _, comic := range comics {
		for _, keyword := range strings.Split(comic.Keywords, " ") {
			index[keyword] = append(index[keyword], fmt.Sprint(comic.Num))
		}
	}

	// вставка данных в БД в таблицу words_index
	for word, numbers := range index {
		_, err := repa.db.Exec("insert or replace into words_index (word, numbers) values ($1, $2)",
			word, strings.Join(numbers, " "))
		if err != nil {
			return err
		}
	}
	return nil
}
