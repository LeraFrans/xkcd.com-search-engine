package main

import (
	"fmt"
	"task5/internal/comic"
)

// нельзя запускать оба варианта сразу. + сделай вывод именно URL

func main() {

	fmt.Println(comic.IndexSearch(comic.InputProccessing()))

	//comic.MakeJSONWithComicsData()

}
