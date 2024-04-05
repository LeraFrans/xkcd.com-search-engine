// 
// Не успела пока сделать пакетную структуру и флаги -о и -n (в модулях и пакетах запуталась окончательно, постараюсь в субботу разобраться)
//
// 1) У меня тут получилось, что половина кода это записи логов после каждого шага практически. Так и должно быть? (в первый раз такое делаю просто),
// 2) Окей ли это, что мы сначала создаём мапу, записываем в неё все данные и только потом её полностью переписываем в json? Можно же ещё сначала сделать json и потом в него по одному добавлять данные о каждом комиксе сразу после получения их с сервера. Что будет эффективнее: работа с дополнительной мапой или постоянное обновление json "в реальном времени"?
// 3) Собрала url сайта простой конкатенацией строки, где меняется только номер комикса. Ну и source_url типа меняется, но он только в теории, потому что для какого-то другого сайта это всё работать не будет. Я не понимаю, для чего нам дана возможность менять source_url, если на другом сайте даже поля у json будут другие. Просто на случай изменения доменного имени сайта с коммиксами, да?

package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

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

func MakeJSONwithComicsData () {
	// Чтение конфигурационного файла
	source_url, db_file := readConfig()
	max_num := 100

	//тут будут храниться конечные данные перед записью их в json
	resultMap := make(map[int]OutputData)
	//циклом проходимся по каждому из комиксов
	for num := 1; num < max_num; num++ {
		input := getDataOfOneComic(num, source_url) //получаем данные об одном комиксе
		output := DataProcess(input)                //обрабатываем их
		resultMap[input.Num] = output               //добавляем в результирующую мапу
	}
	//делаем json из результирующей мапы
	writeResultInJSON(resultMap, db_file)
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
	if err := json.NewDecoder(resp.Body).Decode(&input); err != nil {
		log.Print("Failed decoding from json: ", err)
	}

	return input
}

// Нормализует данные об одном комиксе и оформляет их в структурку
func DataProcess(input InputData) OutputData {
	//обрабатываем нормализатором
	transcriptWithAlt := fmt.Sprint(input.Transcript, input.Alt) //описание и краткое описание сливаем в одну строку, чтобы для каждого нормализацию не делать отдельно
	normalizated, err := WorldsNormalizator(transcriptWithAlt)
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
