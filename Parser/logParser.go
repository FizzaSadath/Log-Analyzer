package main

import (
	"fmt"
	"regexp"
	"time"
)

type LogEntry struct {
	raw       string
	time      time.Time
	level     string
	component string
	host      string
	reqID     string
	msg       string
}

func parseLogEntry(s string) (LogEntry, error) {
	pattern := `^(?P<time>\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d+)\s+\|\s+(?P<level>[A-Z]+)\s+\|\s+(?P<component>[\w-]+)\s+\|\s+host=(?P<host>[\w-]+)\s+\|\s+request_id=(?P<request_id>[\w-]+)\s+\|\s+msg="(?P<msg>.*)"$`
	re := regexp.MustCompile(pattern)
	match := re.FindStringSubmatch(s)
	if match == nil {
		return LogEntry{}, fmt.Errorf("invalid format")
	}
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if name != "" {
			result[name] = match[i]
		}
	}
	parsedTime, err := time.Parse("2006-01-02 15:04:05.000", result["time"])
	if err != nil {
		return LogEntry{}, fmt.Errorf("failed to parse time: %v", err)
	}
	entry := LogEntry{
		raw:       match[0],
		time:      parsedTime,
		level:     result["level"],
		component: result["component"],
		host:      result["host"],
		reqID:     result["request_id"],
		msg:       result["msg"],
	}
	return entry, nil

}
