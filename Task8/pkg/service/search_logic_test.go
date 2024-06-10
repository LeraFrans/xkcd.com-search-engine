package service

import (
	"database/sql"
	"log"
	"os"
	"reflect"
	"task8/pkg/repository"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func Test_findRelevantCount(t *testing.T) {
	type args struct {
		all []repository.IndexResult
	}

	var dat1 []repository.IndexResult
	for i := range 15 {
		dat1 = append(dat1, repository.IndexResult{
			Num:                 i,
			CoutOfRevevantWords: i + 3,
			Url:                 "https://example.com/",
		})
	}

	var dat2 []repository.IndexResult
	for i := range 15 {
		dat1 = append(dat1, repository.IndexResult{
			Num:                 i,
			CoutOfRevevantWords: 1,
			Url:                 "https://example.com/",
		})
	}

	tests := []struct {
		name string
		all  []repository.IndexResult
		want int
	}{
		// TODO: Add test cases.
		{
			name: "Ten",
			all:  dat1,
			want: 10,
		},

		{
			name: "Null",
			all:  dat2,
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findRelevantCount(tt.all); got != tt.want {
				t.Errorf("findRelevantCount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_strToSliceOfInt(t *testing.T) {

	tests := []struct {
		name string
		str  string
		want []int
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			str:  "1 2 3",
			want: []int{1, 2, 3},
		},
		{
			name: "ERR",
			str:  "avkkefh",
			want: []int{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := strToSliceOfInt(tt.str); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("strToSliceOfInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_findOccurrence(t *testing.T) {

	index1 := []repository.IndexFromDB{
		{"word1", "1 2 3"},
		{"word2", "2 3 4"},
	}

	want1 := []repository.IndexResult{
		{1, 1, ""},
		{2, 2, ""},
		{3, 2, ""},
		{4, 1, ""},
	}

	tests := []struct {
		name    string
		indexes []repository.IndexFromDB
		want    []repository.IndexResult
	}{
		// TODO: Add test cases.
		{
			name:    "OK",
			indexes: index1,
			want:    want1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := findOccurrence(tt.indexes); !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("findOccurrence() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSearchService_SearchComics(t *testing.T) {

	os.Remove("testDB.db")
	testDB, s := TestService() // создаём тестовый сервис, использующий тестовую БД
	defer testDB.Close()

	tests := []struct {
		name        string
		searchQuery string
		want        []string
		want1       string
		want2       int
	}{
		{
			name:        "Empty DB",
			searchQuery: "call?trying?but?it?says?number?blocked",
			want:        nil,
			want1:       "Keyword table is exist. Please, do /update request.",
			want2:       400,
		},
		{
			name:        "Index error",
			searchQuery: "sphere refer onli halfway le princ",
			want:        nil,
			want1:       "Error Index search",
			want2:       500,
		},
	}

	for _, tt := range tests {

		// после первого раза пустую БД заполняем данными
		if tt.name == "Index error" {
			addDataComic(testDB)
			s.UpdateIndexTable()
		}

		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := s.SearchComics(tt.searchQuery)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SearchService.SearchComics() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SearchService.SearchComics() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("SearchService.SearchComics() got2 = %v, want %v", got2, tt.want2)
			}
		})

	}
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
