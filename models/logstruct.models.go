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
	Raw       string    `json:"raw"`
	Time      time.Time `json:"time"`
	Level     LogLevel  `json:"level"`
	Component string    `json:"component"`
	Host      string    `json:"host"`
	ReqID     string    `json:"redID"`
	Msg       string    `json:"message"`
}
