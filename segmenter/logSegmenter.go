package segmenter

import (
	"bufio"
	"fmt"
	"log_analyzer/indexer"
	"log_analyzer/models"
	"log_analyzer/parser"
	"os"
	"path/filepath"
)

func ParseLogSegments(path string) (models.LogStore, error) {
	LogStore := models.LogStore{
		Segments: []models.Segment{},
	}
	files, err := os.ReadDir(path)
	if err != nil {
		return LogStore, fmt.Errorf("failed to read directory : %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filepath := filepath.Join(path, file.Name())
		f, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v", filepath, err)
		}
		var LogEntries []models.LogEntry
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // allow 10MB lines

		for scanner.Scan() {
			line := scanner.Text()
			entry, err := parser.ParseLogEntry(line)
			if err == nil {
				LogEntries = append(LogEntries, *entry)
			}
		}
		f.Close()
		if len(LogEntries) == 0 {
			continue
		}
		index := indexer.BuildSegmentIndex(LogEntries)
		segment := models.Segment{
			FileName:   file.Name(),
			LogEntries: LogEntries,
			StartTime:  LogEntries[0].Time,
			EndTime:    LogEntries[len(LogEntries)-1].Time,
			Index:      index,
		}
		LogStore.Segments = append(LogStore.Segments, segment)
	}
	return LogStore, nil

}
