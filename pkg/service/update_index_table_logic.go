package service

import (
	"fmt"
	"strings"
)

// для POST запроса, обновляет таблицу инвертированных индексов по таблице комиксов
// возвращаемая строка это описание, где именно произошла ошибка
func (s *UpdateService) UpdateIndexTable() (string, error) {

	// считывает данные из таблицы комиксов 
	comics, err := s.repo.GetDataForIndexTable()
	if err != nil {
		return "Error geting data from DB for index table", err
	}

	// Создание инвертированного индекса
	var index = make(map[string][]string)
	for _, comic := range comics {
		for _, keyword := range strings.Split(comic.Keywords, " ") {
			index[keyword] = append(index[keyword], fmt.Sprint(comic.Num))
		}
	}

	// вставка данных в БД в таблицу words_index
	err = s.repo.WriteResultInIndexTable(index)
	if err != nil {
		return "Error writing data in index table", err
	}

	return "", err

}
