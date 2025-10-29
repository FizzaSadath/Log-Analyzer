package main

import (
	"fmt"
	"log_analyzer/segmenter"
)

func main() {
	// entries, _ := parser.ParseLogFiles("../logs")
	// for _, entry := range entries {
	// 	fmt.Println(entry)
	// }
	// fmt.Println(len(entries))
	logStore, _ := segmenter.ParseLogSegments("../logs")
	fmt.Println(logStore.Segments[0])
}
