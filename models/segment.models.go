package models

import "time"

type Segment struct {
	FileName   string       `json:"fileName"`
	LogEntries []LogEntry   `json:"logEntries"`
	StartTime  time.Time    `json:"startTime"`
	EndTime    time.Time    `json:"endTime"`
	Index      SegmentIndex `json:"index"`
}
type SegmentIndex struct {
	ByLevel     map[string][]int `json:"byLevel"`
	ByComponent map[string][]int `json:"byComponent"`
	ByHost      map[string][]int `json:"byHost"`
	ByReqID     map[string][]int `json:"byReqID"`
}
