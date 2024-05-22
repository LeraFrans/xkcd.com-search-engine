package repository

import (
	"database/sql"
	"log"
	"task6/config"
)

// используется в пакете server для структуры HTTPHandler
type ComicsService interface {
	List(proposal string) ([]Comic, error)
	Update() (UpdateResponse, error)
}

// просто "временная структура", тут для удобства храним считанные с сайта данные об одном комиксе до обработки
type DataFromXksdCom struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

// информация об одном комиксе (уже обработанная)
type Comic struct {
	Num      int
	Url      string
	Keywords string
}

// Эти двое используются в файле findInDB
// информация об одном индексе, взятая из БД
type IndexFromDB struct {
	Word    string
	Numbers string
}

// финальная информация по индексному поиску, где каждому комиксу присваивается количество релевантных слов в нём
type IndexResult struct {
	Num                 int
	CoutOfRevevantWords int
	Url                 string
}

// отчёт по обновлению таблицы комиксов
type UpdateResponse struct {
	Total int `json:"total"`
	New   int `json:"new"`
}

// после создания подключения к БД, во все функции, где делаем запросы к базе, передаём бд в этом виде (обычно с именем переменной repa)
type SQLiteRepository struct {
	db *sql.DB
	//client xkcd.Client
}

func NewSQLiteRepository(db *sql.DB) *SQLiteRepository {
	return &SQLiteRepository{db: db}
}

// Создание подключения к базе данных (вызывается в хендлерах)
func ConnectDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", config.ReadConfig().Dsn)
	if err != nil {
		log.Printf("Не удалось открыть базу данных: %v", err)
		return nil, err
	}

	return db, nil
}
