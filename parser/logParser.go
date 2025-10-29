package parser

import (
	"bufio"
	"fmt"
	"log_analyzer/models"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

var logPattern = regexp.MustCompile(`^(?P<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+\|\s+(?P<level>[A-Z]+)\s+\|\s+(?P<component>[\w-]+)\s+\|\s+host=(?P<host>[\w-]+)\s+\|\s+request_id=(?P<request_id>[\w-]+)\s+\|\s+msg="(?P<msg>.*)"$`)

// type LogEntry struct {
// 	Raw       string
// 	Time      time.Time
// 	Level     string
// 	Component string
// 	Host      string
// 	ReqID     string
// 	Msg       string
// }

func ParseLogEntry(s string) (*models.LogEntry, error) {
	//pattern := `^(?P<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+\|\s+(?P<level>[A-Z]+)\s+\|\s+(?P<component>[\w-]+)\s+\|\s+host=(?P<host>[\w-]+)\s+\|\s+request_id=(?P<request_id>[\w-]+)\s+\|\s+msg="(?P<msg>.*)"$`
	//re := regexp.MustCompile(pattern)
	match := logPattern.FindStringSubmatch(s)
	if match == nil {
		return nil, fmt.Errorf("invalid format")
	}
	result := make(map[string]string)
	for i, name := range logPattern.SubexpNames() {
		if name != "" {
			result[name] = match[i]
		}
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05.000", result["time"])
	if err != nil {
		return nil, fmt.Errorf("failed to parse time: %v", err)
	}
	entry := models.LogEntry{
		Raw:       match[0],
		Time:      parsedTime,
		Level:     models.ToLogLevel(result["level"]),
		Component: result["component"],
		Host:      result["host"],
		ReqID:     result["request_id"],
		Msg:       result["msg"],
	}

	return &entry, nil

}
func ParseLogFiles(path string) ([]models.LogEntry, error) {
	var allEntries []models.LogEntry
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory : %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filepath := filepath.Join(path, file.Name())
		f, err := os.Open(filepath)
		if err != nil {
			fmt.Printf("Skipping file %s due to error: %v\n", filepath, err)
			continue
		}
		scanner := bufio.NewScanner(f)
		scanner.Buffer(make([]byte, 0, 1024*1024), 10*1024*1024) // allow 10MB lines

		for scanner.Scan() {
			line := scanner.Text()
			entry, err := ParseLogEntry(line)
			if err == nil {
				allEntries = append(allEntries, *entry)
			}
		}
		f.Close()

	}
	return allEntries, nil
}
