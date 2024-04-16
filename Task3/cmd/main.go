package main

import (
	"fmt"
	"task3/comic"
	"time"
)

func main() {
	start := time.Now()

	comic.MakeJSONwithComicsData()

	fmt.Println("\nTime: ", time.Now().Sub(start).Seconds())
}
