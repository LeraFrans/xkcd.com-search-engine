package repository

import (
	"database/sql"
	"log"
	"os"
	"task9/testDB"
	"testing"
)

func TestUpdateSQLite_IsTableHasAnyRows(t *testing.T) {
	os.Remove("testDB6.db")

	// Создание новой базы данных
	testDB, _ := sql.Open("sqlite3", "testDB6.db")
	r := NewUpdateSQLite(testDB)

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

	os.Remove("testDB6.db")
}

func TestUpdateSQLite_FindOurMaxNumberOfComics(t *testing.T) {

	os.Remove("testDB7.db")
	testDB, err := testDB.CreateTestDB("testDB7.db")
	if err != nil {
		log.Fatalf("failed to initializing db: %s", err.Error())
	}
	r := NewRepository(testDB)

	tests := []struct {
		name    string
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "No rows",
			want:    0,
			wantErr: false,
		},
		{
			name:    "OK",
			want:    4,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := r.FindOurMaxNumberOfComics()
			if got != tt.want {
				t.Errorf("UpdateSQLite.FindOurMaxNumberOfComics() = %v, want %v", got, tt.want)
			}
		})

		addDataComic(testDB)
	}
}

func TestUpdateSQLite_WriteResultInKeywordTable(t *testing.T) {
	
	os.Remove("testDB8.db")
	testDB, err := testDB.CreateTestDB("testDB8.db")
	if err != nil {
		log.Fatalf("failed to initializing db: %s", err.Error())
	}
	r := NewUpdateSQLite(testDB)

	tests := []struct {
		name    string
		resultComicsSlice []Comic
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "No rows",
			resultComicsSlice: []Comic{},
			wantErr: false,
		},
		{
			name: "OK",
			resultComicsSlice: []Comic{{1, "URL", "word1 word2"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := r.WriteResultInKeywordTable(tt.resultComicsSlice); (err != nil) != tt.wantErr {
				t.Errorf("UpdateSQLite.WriteResultInKeywordTable() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
