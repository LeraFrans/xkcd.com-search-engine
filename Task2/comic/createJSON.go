// Мы немножечко запутались по поводу функционала флага -n, у меня он работает
// как ограничитель количества обрабатываемых коммиксов. То есть без него в json
// будет записано 2916 штук, а с ним json перезапишется с меньшим числом (n) комиксов

package comic

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"slices"
	"strconv"
	"task2/words"

	"gopkg.in/yaml.v2"
)

// просто "временная структура", тут для удобства храним считанные с сайта данные об одном комиксе до обработки
type InputData struct {
	Num        int    `json:"num"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
}

// а тут уже обработанные данные, предназначенные для записи в json
type OutputData struct {
	Url      string   `json:"url"`
	Keywords []string `json:"keywords"`
}

// основная функция для экспорта
func MakeJSONwithComicsData() {
	// Чтение конфигурационного файла
	source_url, db_file := readConfig()

	// Парсим флаги, получаем булёвое значение флага -о и максимальное число комиксов
	oFlag, max_num := parsArguments()

	//тут будут храниться конечные данные перед записью их в json
	resultMap := make(map[int]OutputData)
	//циклом проходимся по каждому из комиксов
	for num := 1; num <= max_num; num++ {
		input := getDataOfOneComic(num, source_url) //получаем данные об одном комиксе
		output := DataProcess(input)                //обрабатываем их
		resultMap[input.Num] = output               //добавляем в результирующую мапу
		//через каждые 100 делаем запись в json
		if num%100 == 0 {
			writeResultInJSON(resultMap, db_file)
		}
	}
	//делаем финальную запись в json
	writeResultInJSON(resultMap, db_file)

	if oFlag {
		consolePrint(resultMap)
	}
}

func readConfig() (string, string) {
	// Чтение файла
	content, err := os.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	// Разбор YAML-контента
	var yamlMap map[string]interface{}
	err = yaml.Unmarshal(content, &yamlMap)
	if err != nil {
		panic(err)
	}

	source_url, ok := yamlMap["source_url"].(string)
	if !ok {
		panic("Неверный тип данных для 'source_url'")
	}
	db_file, ok := yamlMap["db_file"].(string)
	if !ok {
		panic("Неверный тип данных для 'db_file'")
	}

	return source_url, db_file
}

// Возвращает данные об одном комиксе, полученные с сервера (данные в виде временной структуры InputData)
func getDataOfOneComic(num int, source_url string) InputData {

	//собираем url (без понятия, что делать в случае изменения "/info.0.json", да и вообще не понимаю, для чего нам возможность менять source_url, ведь для другого сайта ничего работать не будет же)
	url := fmt.Sprint(source_url, "/", num, "/info.0.json")

	//получаем данные от сервера
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Failed to make a get request to the server: ", err)
	}
	defer resp.Body.Close()

	//декодируем данные во временную структуру
	var input InputData
	json.NewDecoder(resp.Body).Decode(&input)
	// пришлось убрать обработку ошибок при декодировании, тк она не давала парсить те комиксы, где
	// некоторые поля были не заполнены

	return input
}

// Нормализует данные об одном комиксе и оформляет их в структурку
func DataProcess(input InputData) OutputData {
	//обрабатываем нормализатором
	transcriptWithAlt := fmt.Sprint(input.Transcript, input.Alt) //описание и краткое описание сливаем в одну строку, чтобы для каждого нормализацию не делать отдельно
	normalizated, err := words.WorldsNormalizator(transcriptWithAlt)
	if err != nil {
		log.Fatal(err)
	}

	//формируем элемент данных об одном коммиксе и добавляем его в результирующую мапу
	output := OutputData{
		Url:      input.Img,
		Keywords: normalizated,
	}

	return output
}

// Проверяет, существует ли файл json, если нет, то создаёт его
func checkFileIsExist(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		_, err := os.Create(filename)
		if err != nil {
			return err
		}
	}
	return nil
}

// Записывает результирующую мапу в json
func writeResultInJSON(resultMap map[int]OutputData, db_file string) {
	//подготовка результирующей мапы к записи в json
	dataBytes, err := json.MarshalIndent(resultMap, "", "	")
	if err != nil {
		log.Fatal("Serialization error: ", err)
	}

	//проверяем, существует ли файл json
	err = checkFileIsExist(db_file)
	if err != nil {
		log.Fatal("File is not exixt and can't create:", err)
	}

	//запись
	err = os.WriteFile(db_file, dataBytes, 0)
	if err != nil {
		log.Fatal("Error writing to the json file: ", err)
	}
}

// Парсинг аргументов командной строки
func parsArguments() (bool, int) {
	//первичный поиск флагов -n или -о
	nFlag := flag.Bool("n", false, "Флаг -n")
	oFlag := flag.Bool("o", false, "Флаг -o")
	flag.Parse()

	max_num := 2916 // дефолтное значение

	// обрабатываем оставшиеся аргументы (ищем max_num после -n или флаг -о, если не нашли его вначале)
	args := flag.Args()
	//устанавливаем max_num
	if *nFlag {
		n_num, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Неверный формат числа: %s\n", args[0])
		}
		if max_num > 2916 {
			log.Print("The maximum comic number cannot be more than 2916\n")
		} else {
			max_num = n_num
		}
	}
	//без этого костыля не получилось найти флаг -о, если он стоит после -n
	if slices.Contains(args, "-o") {
		*oFlag = true
	}

	return *oFlag, max_num
}

// Печать результирующей мапы в консоль для флага -о
func consolePrint(resultMap map[int]OutputData) {
	for key, value := range resultMap {
		fmt.Printf("Number: %d: {\n", key)
		fmt.Printf("\tURL: %s\n", value.Url)
		fmt.Printf("\tKey Words: [")
		for i, word := range value.Keywords {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Print(word)
		}
		fmt.Print("]\n}\n")
	}
}
