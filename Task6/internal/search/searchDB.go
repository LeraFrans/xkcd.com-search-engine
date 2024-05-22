package search

import (
	"log"
	"sort"
	"strconv"
	"strings"
	"task6/internal/repository"

	_ "github.com/mattn/go-sqlite3" // Импорт драйвера для SQLite
)

// поиск номеров релевантных комиксов в таблице с индексами   (здесь input - слайс ключевых слов для поиска)
func IndexSearch(input []string, repa *repository.SQLiteRepository) []repository.IndexResult {

	finded := repository.FindInIndex(input, repa) //все комиксы, где встречается хоть раз какое-то слово
	occurrence := findOccurrence(finded)          //теперь храним комикс и количество релевантных слов в нём

	// Сортировка слайса по полю CoutOfRevevantWords
	sort.Slice(occurrence, func(i, j int) bool {
		return occurrence[i].CoutOfRevevantWords > occurrence[j].CoutOfRevevantWords
	})

	// возвращаем только те номера комиксов, где больше всего слов (первые n из отсортированных, n определяетя в findRelevantCount())
	return occurrence[:findRelevantCount(occurrence)]
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
