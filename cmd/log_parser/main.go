package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"log_analyzer/pkg/segmenter"
	"os"
)

func main() {
	logPath := flag.String("path", "/home/fizza/goProjects/log_analyzer/logs", "Path to the log directory")
	outFile := flag.String("out", "parsed_logs.json", "Output JSON file path")
	flag.Parse()

	logStore, err := segmenter.ParseLogSegments(*logPath)
	if err != nil {
		slog.Error("Failed to parse logs", "error", err)

	}

	file, err := os.Create(*outFile)
	if err != nil {
		slog.Error("Could not create output file", "error", err)

	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(logStore); err != nil {
		slog.Error("Error writing JSON", "error", err)
	}

	fmt.Printf("Logs parsed successfully and saved to %s\n", *outFile)
}
