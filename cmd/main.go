package main

import (
	"flag"
	"fmt"
	"log/slog"
	"log_analyzer/filter"
	"log_analyzer/segmenter"
)

func main() {
	logPath := flag.String("path", "/home/fizza/goProjects/log_analyzer/logs", "Path to the log directory (required)")
	level := flag.String("level", "", "Comma separated list of log levels")
	component := flag.String("component", "", "Comma separated list of components")
	host := flag.String("host", "", "Comma separated list of hosts")
	reqID := flag.String("reqID", "", "Comma separated list of requestIDs")
	startTimeString := flag.String("after", "", "Filter by start time [2006-01-02 15:04:05]")
	endTimeString := flag.String("before", "", "Filter by end time [2006-01-02 15:04:05]")
	flag.Parse()

	logStore, err := segmenter.ParseLogSegments(*logPath)
	if err != nil {
		slog.Error("Failed to parse logs\n")
	}

	levels := split(*level)
	components := split(*component)
	hosts := split(*host)
	reqIDs := split(*reqID)
	startTime := parseTimeFlag(*startTimeString)
	endTime := parseTimeFlag(*endTimeString)

	filteredLogs := filter.FilterLogs(logStore, levels, components, hosts, reqIDs, startTime, endTime)
	fmt.Printf("Found %d matching entries\n", len(filteredLogs))
	for _, entry := range filteredLogs {
		fmt.Println(entry.Raw)
	}
}
