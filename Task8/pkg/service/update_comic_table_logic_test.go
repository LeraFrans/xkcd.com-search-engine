package service

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"task8/pkg/repository"
	"testing"

	"github.com/go-playground/assert/v2"
)

// // Оригинал этой функции никак не проверишь, тк нельзя установить ожидаемое значение,
// // количество комикосв на сайте постоянно увеличивается. Создаём заглушку на 10 комиксов для скорости
// func findMaxNumberOfComics() int {
// 	return 10
// }

func Test_isPageExists(t *testing.T) {
	tests := []struct {
		name     string
		page_num int
		want     bool
	}{
		{
			name:     "OK",
			page_num: 5,
			want:     true,
		},
		{
			name:     "NOT PAGE",
			page_num: 123456,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPageExists(tt.page_num); got != tt.want {
				t.Errorf("isPageExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDataProcess(t *testing.T) {

	tests := []struct {
		name string
		args DataFromXksdCom
		want repository.Comic
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			args: DataFromXksdCom{
				11,
				"[[A boy sits in a barrel which is floating in an ocean.]]\nBoy: None of the places i floated had mommies.\n{{Alt: Awww.}}",
				"Awww.",
				"https://imgs.xkcd.com/comics/barrel_mommies.jpg",
			},
			want: repository.Comic{
				Num:      11,
				Url:      "https://imgs.xkcd.com/comics/barrel_mommies.jpg",
				Keywords: "none place sit is float ocean awww boy barrel had mommi",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			got := DataProcess(tt.args)
			assert.Equal(t, got.Num, tt.want.Num)
			assert.Equal(t, got.Url, tt.want.Url)
			assert.NotEqual(t, 0, len(got.Keywords))
			assert.Equal(t, len(got.Keywords), len(tt.want.Keywords))

		})
	}
}

func Test_getDataOfOneComic(t *testing.T) {
	tests := []struct {
		name string
		num  int
		want DataFromXksdCom
	}{
		// TODO: Add test cases.
		{
			name: "OK",
			num:  11,
			want: DataFromXksdCom{
				11,
				"[[A boy sits in a barrel which is floating in an ocean.]]\nBoy: None of the places i floated had mommies.\n{{Alt: Awww.}}",
				"Awww.",
				"https://imgs.xkcd.com/comics/barrel_mommies.jpg",
			},
		},
		{
			name: "NOT FOUND PAGE",
			num:  404,
			want: DataFromXksdCom{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getDataOfOneComic(tt.num); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getDataOfOneComic() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_readCache(t *testing.T) {

	var empty []byte

	tests := []struct {
		name     string
		filename string
		want     bool
		want1    []byte
	}{
		{
			name:     "No cache",
			filename: "test_cache.txt",
			want:     false,
			want1:    empty,
		},
		{
			name:     "Have cache",
			filename: "test_cache.txt",
			want:     true,
			want1:    []byte("1 2 3"),
		},
	}

	// тут на всякий случай удаление файла, если он существует
	if _, err := os.Stat("test_cache.txt"); !os.IsNotExist(err) {
		os.Remove("test_cache.txt")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := readCache(tt.filename)
			if got != tt.want {
				t.Errorf("readCache() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("readCache() got1 = %v, want %v", got1, tt.want1)
			}
		})

		// а тут его создание и заполнение чем-то
		file, _ := os.Create("test_cache.txt")
		defer file.Close()
		file.Write([]byte("1 2 3"))
	}
}

func Test_createCache(t *testing.T) {

	//удаляем кэш если он был
	if _, err := os.Stat("written_comics.txt"); !os.IsNotExist(err) {
		os.Remove("written_comics.txt")
	}

	writen1 := []int{1, 2, 3}
	t.Run("No cache", func(t *testing.T) {
		createCache(&writen1)
		//проверяем, что файл создался
		res := 1
		if _, err := os.Stat("written_comics.txt"); os.IsNotExist(err) {
			res = 0
		}
		assert.Equal(t, 1, res)
	})

	writen2 := []int{4, 5, 6}
	t.Run("Have cache", func(t *testing.T) {
		createCache(&writen2)

		// проверяем содержимое файла

		content, _ := os.ReadFile("written_comics.txt")

		res := 1
		for num := 1; num <= 6; num++ {
			if !bytes.Contains(content, []byte(fmt.Sprintf("%d", num))) {
				res = 0
			}
		}
		assert.Equal(t, 1, res)
	})
}

func TestUpdateService_findNumbersToGet(t *testing.T) {

	os.Remove("testDB.db")
	os.Remove("written_comics.txt")
	testDB, s := TestService() // создаём тестовый сервис, использующий тестовую БД
	defer testDB.Close()

	max_num := 100

	tests := []struct {
		name    string
		numbers chan int
		max_num int
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "No DB, No cahce",
			numbers: make(chan int, max_num),
			max_num: max_num,
			want:    100,
			wantErr: false,
		},
		{
			name:    "DB, No cahce 1",
			numbers: make(chan int, max_num),
			max_num: max_num,
			want:    96,
			wantErr: false,
		},
	}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			got, err := s.findNumbersToGet(tt.numbers, tt.max_num)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateService.findNumbersToGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UpdateService.findNumbersToGet() = %v, want %v", got, tt.want)
			}
		})

		// после первого теста заполняем БД
		addDataComic(testDB)
	}

	// ещё один тест
	t.Run("DB, No cahce 2", func(t *testing.T) {
		got, err := s.findNumbersToGet(make(chan int, 4), 4)
		if (err != nil) != false {
			t.Errorf("UpdateService.findNumbersToGet() error = %v, wantErr %v", err, false)
			return
		}
		if got != 0 {
			t.Errorf("UpdateService.findNumbersToGet() = %v, want %v", got, 0)
		}
	})

	//ещё один тест 2

	// создадим кэш
	writen := []int{1, 2}
	createCache(&writen)
	t.Run("DB, cahce", func(t *testing.T) {
		got, err := s.findNumbersToGet(make(chan int, 4), 4)
		if (err != nil) != false {
			t.Errorf("UpdateService.findNumbersToGet() error = %v, wantErr %v", err, false)
			return
		}
		if got != 2 {
			t.Errorf("UpdateService.findNumbersToGet() = %v, want %v", got, 2)
		}
	})

	os.Remove("testDB.db")
	os.Remove("written_comics.txt")
}

func TestUpdateService_UpdateComicTable(t *testing.T) {
	os.Remove("testDB.db")
	os.Remove("written_comics.txt")
	testDB, s := TestService() // создаём тестовый сервис, использующий тестовую БД
	defer testDB.Close()

	max_numbers := findMaxNumberOfComics()

	tests := []struct {
		name    string
		want    UpdateResponse
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "first",
			want:    UpdateResponse{max_numbers, max_numbers},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := s.UpdateComicTable()
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateService.UpdateComicTable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UpdateService.UpdateComicTable() = %v, want %v", got, tt.want)
			}
		})
	}

	os.Remove("testDB.db")
	os.Remove("written_comics.txt")
}

// func Test_findMaxNumberOfComics(t *testing.T) {

// 	t.Run("OK", func(t *testing.T) {
// 		if got := findMaxNumberOfComics(); got < 2935 {
// 			t.Errorf("findMaxNumberOfComics() = %v, it is less then %d", got, 2935)
// 		}
// 	})
// }
