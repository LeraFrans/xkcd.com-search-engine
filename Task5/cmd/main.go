package main

import (
	//"fmt"
	//"fmt"
	"os"
	"task5/comic"

	//"time"

	"gopkg.in/yaml.v2"
)

// нельзя запускать оба варианта сразу. + сделай вывод именно URL

func main() {

	//comic.MakeJSONWithComicsData(readConfig())

	// выбрать один из вариантов
	//result := comic.SimpleSearch(comic.InputProccessing())
	//fmt.Println(comic.IndexSearch(comic.InputProccessing()))

	//fmt.Println(result)

	// for i := range result {
	// 	fmt.Println(result[i])
	// }

	comic.MakeJSONWithComicsData(readConfig())

}

// Возвращает переменные из "config.yaml"
func readConfig() (string, string, int) {
	// Чтение файла

	//configName := comic.ParsArgument_C()
	configName := "config.yaml"

	content, err := os.ReadFile(configName)
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
	parallel, ok := yamlMap["parallel"].(int)
	if !ok {
		panic("Неверный тип данных для 'parallel'")
	}

	return source_url, db_file, parallel
}
