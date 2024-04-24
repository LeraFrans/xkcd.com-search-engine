package comic

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"sort"
	"task4/words"
)

type Pair struct {
	Num   int // Ключ
	Count int // Значение
}

func SimpleSearch(input []string) []Pair {
	// Считываем данные из database.json
	data := readDataBase()

	// Проходимся по базе и выписываем оттуда номера комиксов, в которых встречалось нужное слово
	var index = make(map[string][]int) // ключ - слово, значение - массив номеров комиксов, где есть слово
	for key, item := range data {
		for _, keyword := range item.Keywords {
			if containsElement(input, keyword) {
				index[keyword] = append(index[keyword], key)
			}
		}
	}

	occurrence := findOccurrence(index)              // считаем сколько совпадающих слов найдено в каждом номере комикса
	sortedOccurrence := sortByOccurrence(occurrence) // сортируем по убыванию количества релевантных слов в комиксе

	// // Выводим отсортированные пары.
	// fmt.Println("Аfter sorting:")
	// for _, pair := range sortedOccurrence {
	// 	fmt.Printf("%d: %d\n", pair.Num, pair.Count)
	// }

	return sortedOccurrence[:findRelevantCount(sortedOccurrence)]

}

func IndexSearch(input []string) []Pair {
	сreateIndex()

	finded := findInIndexJSON(input)
	occurrence := findOccurrence(finded)
	sortedOccurrence := sortByOccurrence(occurrence)

	// // Выводим отсортированные пары.
	// fmt.Println("Аfter sorting:")
	// for _, pair := range sortedOccurrence {
	// 	fmt.Printf("%d: %d\n", pair.Num, pair.Count)
	// }

	return sortedOccurrence[:findRelevantCount(sortedOccurrence)]
}

// Создание файла с индексами
func сreateIndex() {

	// Проверяем наличие файла с индексами
	indexFilePath := "database/index.json"
	if _, err := os.Stat(indexFilePath); err != nil {
		if os.IsNotExist(err) {
			createFileIfNotExist(indexFilePath)
		}
	} else {
		return
	}

	// Считываем данные из database.json
	data := readDataBase()

	// Создание инвертированного индекса
	var index = make(map[string][]int)
	for key, item := range data {
		for _, keyword := range item.Keywords {
			index[keyword] = append(index[keyword], key)
		}
	}

	// Подготовка индексов к записи в json
	dataBytes, _ := json.MarshalIndent(index, "", "	")

	// Запись индексов в json
	indexFile, _ := os.Open(indexFilePath)
	defer indexFile.Close()
	os.WriteFile(indexFilePath, dataBytes, 0)
}

// Парсинг строки из консоли
func parsArguments() string {
	error_message := "Please use the -s flag with string in double quotes"
	//pointer at start string
	pointerToParsedLine := flag.String("s", error_message, error_message)
	flag.Parse()
	return *pointerToParsedLine
}

// Нормализация слов из консоли
func InputProccessing() []string {
	parsedLine := parsArguments()
	result, err := words.WordsNormalizator(parsedLine)
	if err != nil {
		log.Println("Error input normalization: ", err)
	}

	return result
}

// Содержит ли слайс определённый элемент
func containsElement(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

// Читает файл database.json
func readDataBase() map[int]OutputData {
	// Чтение database.json
	dataBaseFilePath := "database/database.json"
	dataFile, err := os.ReadFile(dataBaseFilePath)
	if err != nil {
		panic(err)
	}

	// Декодирование JSON в слайс структур
	data := make(map[int]OutputData)
	err = json.Unmarshal(dataFile, &data)
	if err != nil {
		log.Println(err)
	}

	return data
}

// считает, сколько релевантных слов в каждом комиксе
func findOccurrence(index map[string][]int) map[int]int {
	occurrence := make(map[int]int) // номер комикса -> количество совпавших слов
	for _, numbers := range index {
		for n := range numbers {
			occurrence[numbers[n]]++
		}
	}
	return occurrence
}

// сортирует мапу из предыдущей функции по убыванию релевантных слов
func sortByOccurrence(occurrence map[int]int) []Pair {
	// Создаем слайс пар для сортировки.
	pairs := make([]Pair, 0, len(occurrence))
	for num, count := range occurrence {
		pairs = append(pairs, Pair{num, count})
	}

	// Сортируем слайс пар по значению.
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].Count > pairs[j].Count
	})

	return pairs
}

// Для функции IndexSearch(). Ищет в файле с индексами нужные нам релевантные слова
func findInIndexJSON(input []string) map[string][]int {

	// чтение
	indexFilePath := "../database/index.json"
	indexFile, err := os.ReadFile(indexFilePath)
	if err != nil {
		panic(err)
	}

	// Декодирование JSON в слайс структур
	indexFromJSON := make(map[string][]int)
	err = json.Unmarshal(indexFile, &indexFromJSON)
	if err != nil {
		log.Println(err)
	}

	// мапа, где каждому слову соответсвуют номера комиксов, в которых оно встречается
	finded := make(map[string][]int)
	for word, nums := range indexFromJSON {
		if containsElement(input, word) {
			finded[word] = nums
		}
	}

	return finded
}

// Ищем максимальное число комиксов для вывода
// Либо 10, либо меньше, если количество релевантных слов меньше 2
func findRelevantCount(all []Pair) int {
	result := 0

	for _, elem := range all {
		if elem.Count == 1 || result == 10 {
			break
		}
	}

	return result
}
