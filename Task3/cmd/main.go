package main

import (
	"fmt"
	"os"
	"task3/comic"
	"time"

	"gopkg.in/yaml.v2"
)

func main() {
	start := time.Now()

	comic.MakeJSONWithComicsData(readConfig())

	fmt.Println("\nTime: ", time.Now().Sub(start).Seconds())
}


// Возвращает переменные из "config.yaml"
func readConfig() (string, string, int) {
	// Чтение файла

	configName := comic.ParsArgument_C()

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
