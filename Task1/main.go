package main

import (
	"fmt"
	"strings"
)

func main() {
	parsedLine := parsArguments()
	result, err := worldsNormalizator(parsedLine)
	if err == nil {
		fmt.Println(strings.Join(result, " "))
	}
}