package models

import "time"

type LogLevel string

const (
	INFO  LogLevel = "INFO"
	WARN  LogLevel = "WARN"
	DEBUG LogLevel = "DEBUG"
	ERROR LogLevel = "ERROR"
)

func ToLogLevel(s string) LogLevel {
	return LogLevel(s)
}

type LogEntry struct {
	Raw       string
	Time      time.Time
	Level     LogLevel
	Component string
	Host      string
	ReqID     string
	Msg       string
}
