package main

import (
	"errors"
	"flag"
	"strings"
	"unicode"

	"github.com/kljensen/snowball"
)

func parsArguments() string {
	error_message := "Please use the -s flag with string in double quotes"
	//pointer at start string
	pointerToParsedLine := flag.String("s", error_message, error_message)
	flag.Parse()
	return *pointerToParsedLine
}

func trimPunctuation(st string) []string {
	// Splyce string into words. Use lambda function as separator
	return strings.FieldsFunc(st, func(symbol rune) bool {
		// Split on any character that is not a letter or a number
		return !unicode.IsLetter(symbol) && !unicode.IsNumber(symbol)
	})
}

func deleteMostCommon(input []string) []string {
	var commonWords = map[string]struct{}{
		"a": {}, "and": {}, "be": {}, "have": {}, "i": {}, "me": {},
		"in": {}, "of": {}, "that": {}, "the": {}, "to": {},
		"m": {}, "s": {}, "ll": {}, "it": {}, "for": {}, "not": {},
		"on": {}, "with": {}, "he": {}, "as": {}, "you": {}, "am": {},
		"at": {}, "this": {}, "by": {}, "his": {}, "from": {},
		"they": {}, "we": {}, "her": {}, "she": {}, "or": {}, "are": {},
		"an": {}, "will": {}, "my": {}, "would": {}, "there": {},
		"their": {}, "what": {}, "so": {}, "if": {}, "who": {},
		"get": {}, "which": {}, "when": {}, "can": {}, "him": {},
		"your": {}, "some": {}, "them": {}, "then": {}, "its": {},
		"also": {}, "us": {},
	}
	output := make([]string, len(input))
	for _, elem := range input {
		_, ok := commonWords[elem]
		if !ok {
			output = append(output, elem)
		}
	}
	return output
}

// removes duplicates
func unique(input []string) string {
	output := ""
	set := map[string]bool{}
	for _, elem := range input {
		set[elem] = true
	}
	for key, _ := range set {
		output = output[:] + " " + key
	}
	return strings.Trim(output, " ")
}

// modification of the word form
func stemming(input []string) (string, error) {
	ss := strings.Join(input, " ")
	stemmed, err := snowball.Stem(ss, "english", true)
	if err == nil {
		return unique(strings.Fields(stemmed)), nil
	}
	return "", errors.New("Stemming error")
}
