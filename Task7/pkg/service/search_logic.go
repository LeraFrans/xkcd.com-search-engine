package service

import (
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"task7/pkg/repository"

	_ "github.com/mattn/go-sqlite3" // Импорт драйвера для SQLite
)

// для связи между слоями
type SearchService struct {
	repo repository.Search
}

func NewSearchService(repo repository.Search) *SearchService {
	return &SearchService{repo: repo}
}

// самая верхняя из логики поиска функция
// возвращаемая string это описание ошибки, а int это её статус код, чтобы потом в хендлере вернуть
func (s *SearchService) SearchComics(searchQuery string) ([]string, string, int) {
	// Получаем строку для поиска
	normSearchQuery, errNormalization := WordsNormalizator(searchQuery) // нормализуем и стеммим её
	if errNormalization != nil {
		return nil, "Error normalizing search query.", http.StatusBadRequest
	}

	// Если таблица комиксов пустая, просим обновить базу
	if !s.repo.IsTableHasAnyRows("keyword") {
		return nil, "Keyword table is exist. Please, do /update request.", http.StatusBadRequest
	}


	// поиск комиксов
	indexes, err := s.IndexSearch(normSearchQuery) // ищем релевантные номера комиксов в таблице индексов
	if err != nil {
		return nil, "Error Index search", http.StatusInternalServerError
	}
	urls, err := s.repo.FindURL(indexes)            // по номерам смотрим урлы нужных комиксов
	if err != nil {
		return nil, "Error findinf URL", http.StatusInternalServerError
	}

	return urls, "", http.StatusOK

} 

// поиск номеров релевантных комиксов в таблице с индексами   (здесь input - слайс ключевых слов для поиска)
func (s *SearchService) IndexSearch(input []string) ([]repository.IndexResult, error) {

	finded, err := s.repo.FindInIndex(input)  //все комиксы, где встречается хоть раз какое-то слово
	if err != nil {
		return nil, err
	}
	occurrence := findOccurrence(finded) //теперь храним комикс и количество релевантных слов в нём

	// Сортировка слайса по полю CoutOfRevevantWords
	sort.Slice(occurrence, func(i, j int) bool {
		return occurrence[i].CoutOfRevevantWords > occurrence[j].CoutOfRevevantWords
	})

	// возвращаем только те номера комиксов, где больше всего слов (первые n из отсортированных, n определяетя в findRelevantCount())
	return occurrence[:findRelevantCount(occurrence)], nil
}

// считает, сколько релевантных слов в каждом комиксе
func findOccurrence(indexes []repository.IndexFromDB) []repository.IndexResult {
	occurrenceMap := make(map[int]int) // [номер комикса] -> количество совпавших слов
	for _, index := range indexes {
		for _, n := range strToSliceOfInt(index.Numbers) {
			occurrenceMap[n]++
		}
	}

	// переводим в нужную структуру (костыль)
	var result []repository.IndexResult
	for number, count := range occurrenceMap {
		r := repository.IndexResult{
			Num:                 number,
			CoutOfRevevantWords: count,
		}
		result = append(result, r)
	}
	return result
}

func strToSliceOfInt(str string) []int {
	sliceOfString := strings.Split(str, " ")

	var sliceOfInt []int

	for _, s := range sliceOfString {
		i, err := strconv.Atoi(s)
		if err != nil {
			log.Panicln(err)
		}
		sliceOfInt = append(sliceOfInt, i)
	}

	return sliceOfInt
}

// Ищем максимальное число комиксов для вывода
// Либо 10, либо меньше, если количество релевантных слов меньше 2
func findRelevantCount(all []repository.IndexResult) int {
	result := 0

	for _, elem := range all {
		if elem.CoutOfRevevantWords == 1 || result == 10 {
			break
		}
		result++
	}
	return result
}
