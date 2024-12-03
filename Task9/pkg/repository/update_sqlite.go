package repository

import (
	"database/sql"
	"log"
	"strings"
)

type UpdateSQLite struct {
	db *sql.DB
}

func NewUpdateSQLite(db *sql.DB) *UpdateSQLite {
	return &UpdateSQLite{db: db}
}

// информация об одном комиксе (уже обработанная)
type Comic struct {
	Num      int
	Url      string
	Keywords string
}

// финальная информация по индексному поиску, где каждому комиксу присваивается количество релевантных слов в нём
type IndexResult struct {
	Num                 int
	CoutOfRevevantWords int
	Url                 string
}

// сколько комиксов сейчас записано в базе
func (r *UpdateSQLite) FindOurMaxNumberOfComics() (int, error) {

	// Подготовка SQL запроса для получения максимального номера
	stmt, errPrepare := r.db.Prepare("SELECT MAX(number) FROM keyword")
	if errPrepare != nil {
		return 0, errPrepare
	}
	defer stmt.Close()

	// Выполнение запроса
	var maxNumber int
	if errQuery := stmt.QueryRow().Scan(&maxNumber); errQuery != nil {
		log.Println(errQuery)
		return 0, errQuery
	}

	return maxNumber, nil
}

// Содержит ли таблица хоть одну запись (для проверки перед поиском)
func (r *UpdateSQLite) IsTableHasAnyRows(table_name string) bool {


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

// Записывает результат в БД в таблицу комиксов
func (r *UpdateSQLite) WriteResultInKeywordTable(resultComicsSlice []Comic) error {

	// Подготовка SQL-запросов для вставки данных
	stmt, _ := r.db.Prepare("INSERT OR REPLACE INTO keyword (number, url, keywords) VALUES (?, ?, ?)")
	defer stmt.Close()

	// Запись данных из resultComicsSlice в базу данных
	for _, data := range resultComicsSlice {
		_, err := stmt.Exec(data.Num, data.Url, data.Keywords)
		if err != nil {
			return err
		}
	}

	return nil
}

// берёт инфу обо всех комиксах из таблицы с комиксами
func (r *UpdateSQLite) GetDataForIndexTable() ([]Comic, error) {
	// Запрос нужных данных
	rows, err := r.db.Query("SELECT number, keywords FROM keyword")
	if err != nil {
		return nil, err
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

	return comics, nil
}

// записывает уже обработанную инфу в таблицу с индексами
func (r *UpdateSQLite) WriteResultInIndexTable(index map[string][]string) error {
	for word, numbers := range index {
		_, err := r.db.Exec("insert or replace into words_index (word, numbers) values ($1, $2)",
			word, strings.Join(numbers, " "))
		if err != nil {
			return err
		}
	}
	return nil
}
