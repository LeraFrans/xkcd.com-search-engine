package words_processing

// Это всё для цели xkcd в сmd. Для сервера это не надо

import (
	"flag"
	"log"
)

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
	result, err := WordsNormalizator(parsedLine)
	if err != nil {
		log.Println("Error input normalization: ", err)
	}

	return result
}