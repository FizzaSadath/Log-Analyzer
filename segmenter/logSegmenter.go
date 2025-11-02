package segmenter

import (
	"bufio"
	"fmt"
	"log_analyzer/indexer"
	"log_analyzer/models"
	"log_analyzer/parser"
	"os"
	"path/filepath"
	"sync"
	"time"
)

func ParseLogSegments(path string) (models.LogStore, error) {
	start := time.Now()
	logStore := models.LogStore{
		Segments: []models.Segment{},
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return logStore, fmt.Errorf("failed to read directory: %v", err)
	}

	var wg sync.WaitGroup
	var mu sync.Mutex // protects logStore.Segments

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		wg.Add(1)
		go func(file os.DirEntry) {
			defer wg.Done()

			filepath := filepath.Join(path, file.Name())
			f, err := os.Open(filepath)
			if err != nil {
				fmt.Printf("Skipping file %s due to error: %v\n", filepath, err)
				return
			}
			defer f.Close()

			var logEntries []models.LogEntry
			scanner := bufio.NewScanner(f)
			scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // allow 10MB lines

			for scanner.Scan() {
				line := scanner.Text()
				entry, err := parser.ParseLogEntry(line)
				if err == nil {
					logEntries = append(logEntries, *entry)
				}
			}

			if len(logEntries) == 0 {
				return
			}

			index := indexer.BuildSegmentIndex(logEntries)
			segment := models.Segment{
				FileName:   file.Name(),
				LogEntries: logEntries,
				StartTime:  logEntries[0].Time,
				EndTime:    logEntries[len(logEntries)-1].Time,
				Index:      index,
			}

			// lock before writing to shared slice
			mu.Lock()
			logStore.Segments = append(logStore.Segments, segment)
			mu.Unlock()

		}(file)
	}

	wg.Wait()

	elapsed := time.Since(start)
	fmt.Println("Segment parsing took:", elapsed)
	return logStore, nil
}
