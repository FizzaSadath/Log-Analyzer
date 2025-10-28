package main

import (
	"fmt"
	"log_analyzer/parser"
)

func main() {
	entries, _ := parser.ParseLogFiles("../logs")
	for _, entry := range entries {
		fmt.Println(entry)
	}
	fmt.Println(len(entries))
}
