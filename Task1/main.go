package main

import (
	"fmt"
	"strings"
)

func main() {
	parsedLine := parsArguments()
	withoutPunctuationArray := trimPunctuation(strings.ToLower(parsedLine))
	withoutCommonWordsArray := deleteMostCommon(withoutPunctuationArray)
	result, err := stemming(withoutCommonWordsArray)
	if err == nil {
		fmt.Println(result)
	}
}
