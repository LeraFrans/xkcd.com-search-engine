package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"task8/testDB"
	"testing"
)

// нужна таблица индексов но нельзя идти в слой сервиса
// мб сделать дубликат той функции сюда? или вручную заполнить таблицу индексов?
func TestSearchSQLite_FindInIndex(t *testing.T) {
	testDB, err := testDB.CreateTestDB("testDB3.db")
	if err != nil {
		log.Fatalf("failed to initializing db: %s", err.Error())
	}
	r := NewRepository(testDB)

	addDataComic(testDB)
	r.addDataIndex()

	tests := []struct {
		name           string
		input_keywords []string
		want           []IndexFromDB
		wantErr        bool
	}{
		{
			name:           "OK",
			input_keywords: []string{"through", "flow"},
			want:           []IndexFromDB{{"through", "2 4"}, {"flow", "4"}},
			wantErr:        false,
		},
		{
			name:           "test 2",
			input_keywords: []string{},
			want:           nil,
			wantErr:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindInIndex(tt.input_keywords)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchSQLite.FindInIndex() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchSQLite.FindInIndex() = %v, want %v", got, tt.want)
			}
		})
	}
	os.Remove("testDB3.db")
}

func addDataComic(testDB *sql.DB) {

	// заполняем таблицу комиксов
	if _, err := testDB.Exec("INSERT INTO keyword (number, url, keywords) VALUES (?, ?, ?)", 1, "https://imgs.xkcd.com/comics/barrel_cropped_(1).jpg", "boy barrel is drift ocean next noth seen t els don all sit float wonder where into distanc"); err != nil {
		log.Print(err)
	}
	if _, err := testDB.Exec("INSERT INTO keyword (number, url, keywords) VALUES (?, ?, ?)", 2, "https://imgs.xkcd.com/comics/tree_cropped_(1).jpg", "sphere refer onli halfway le princ thought sketch two opposit side titl grow be about tree petit through"); err != nil {
		log.Print(err)
	}
	if _, err := testDB.Exec("INSERT INTO keyword (number, url, keywords) VALUES (?, ?, ?)", 3, "https://imgs.xkcd.com/comics/island_color.jpg", "sketch island hello"); err != nil {
		log.Print(err)
	}
	if _, err := testDB.Exec("INSERT INTO keyword (number, url, keywords) VALUES (?, ?, ?)", 4, "https://imgs.xkcd.com/comics/landscape_cropped_(1).jpg", "flow through ocean sketch landscap sun horizon river"); err != nil {
		log.Print(err)
	}
}

// копия функции  UpdateIndexTable() из пакета service
func (r *Repository) addDataIndex() string {
	// считывает данные из таблицы комиксов
	comics, err := r.GetDataForIndexTable()
	if err != nil {
		return "Error geting data from DB for index table"
	}

	// Создание инвертированного индекса
	var index = make(map[string][]string)
	for _, comic := range comics {
		for _, keyword := range strings.Split(comic.Keywords, " ") {
			index[keyword] = append(index[keyword], fmt.Sprint(comic.Num))
		}
	}

	// вставка данных в БД в таблицу words_index
	err = r.WriteResultInIndexTable(index)
	if err != nil {
		return "Error writing data in index table"
	}

	return ""
}

func TestSearchSQLite_FindURL(t *testing.T) {
	testDB, err := testDB.CreateTestDB("testDB4.db")
	if err != nil {
		log.Fatalf("failed to initializing db: %s", err.Error())
	}
	r := NewRepository(testDB)

	addDataComic(testDB)
	r.addDataIndex()

	tests := []struct {
		name    string
		indexes []IndexResult
		want    []string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "OK",
			indexes: []IndexResult{{Num: 2}},
			want:    []string{"https://imgs.xkcd.com/comics/tree_cropped_(1).jpg"},
			wantErr: false,
		},
		{
			name:    "NO ROWS",
			indexes: []IndexResult{},
			want:    []string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.FindURL(tt.indexes)
			if (err != nil) != tt.wantErr {
				t.Errorf("SearchSQLite.FindURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchSQLite.FindURL() = %v, want %v", got, tt.want)
			}
		})
	}

	os.Remove("testDB4.db")
}

func TestSearchSQLite_IsTableHasAnyRows(t *testing.T) {

	os.Remove("testDB5.db")

	// Создание новой базы данных
	testDB, _ := sql.Open("sqlite3", "testDB5.db")
	r := NewSearchSQLite(testDB)

	// Создание таблицы keyword
	stmt, _ := testDB.Prepare("CREATE TABLE IF NOT EXISTS keyword (number INTEGER, url TEXT, keywords TEXT)")
	stmt.Exec()
	t.Run("NO rows", func(t *testing.T) {
		if got := r.IsTableHasAnyRows("keyword"); got != false {
			t.Errorf("SearchSQLite.IsTableHasAnyRows() = %v, want %v", got, false)
		}
	})

	// добавление строк в таблицу
	addDataComic(testDB)

	t.Run("OK", func(t *testing.T) {
		if got := r.IsTableHasAnyRows("keyword"); got != true {
			t.Errorf("SearchSQLite.IsTableHasAnyRows() = %v, want %v", got, true)
		}
	})
	os.Remove("testDB5.db")

}
