package main

import (
	"testing"
	"time"
)

func TestParseLogEntry(t *testing.T) {
	line := `2025-10-23 15:17:08.636 | WARN | api-server | host=worker01 | request_id=req-4leuyy-5910 | msg="Cache cleared"`
	entry, err := parseLogEntry(line)

	if err != nil {
		t.Errorf("Log Parsing Failed!")
	}
	expectedTime, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-23 15:17:08.636")
	if entry.raw != line {
		t.Errorf("Expected raw to be %q but got %q", line, entry.raw)
	}
	if !entry.time.Equal(expectedTime) {
		t.Errorf("Expected time %v but got %v", expectedTime, entry.time)
	}
	if entry.level != "WARN" {
		t.Errorf("Expected WARN but got %s./n", entry.level)
	}

	if entry.component != "api-server" {
		t.Errorf("Expected api-server but got %s.\n", entry.component)
	}
	if entry.host != "worker01" {
		t.Errorf("Expected api-server but got %s.\n", entry.component)
	}
	if entry.reqID != "req-4leuyy-5910" {
		t.Errorf("Expected req-4leuyy-5910 but got %s.\n", entry.reqID)
	}
	if entry.msg != "Cache cleared" {
		t.Errorf("Expected 'Cache cleared' but got %s.\n", entry.msg)
	}
	// if entry.time!=2025-10-23 15:17:08.636{

	// }

}
