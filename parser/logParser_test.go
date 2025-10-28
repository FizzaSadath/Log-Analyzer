package parser

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestParseLogEntry(t *testing.T) {
	line := `2025-10-23 15:17:08.636 | WARN | api-server | host=worker01 | request_id=req-4leuyy-5910 | msg="Cache cleared"`
	entry, err := ParseLogEntry(line)

	if err != nil {
		t.Errorf("Log Parsing Failed!")
	}
	expectedTime, _ := time.Parse("2006-01-02 15:04:05.000", "2025-10-23 15:17:08.636")
	if entry.Raw != line {
		t.Errorf("Expected raw to be %q but got %q", line, entry.Raw)
	}
	if !entry.Time.Equal(expectedTime) {
		t.Errorf("Expected time %v but got %v", expectedTime, entry.Time)
	}
	if entry.Level != "WARN" {
		t.Errorf("Expected WARN but got %s./n", entry.Level)
	}

	if entry.Component != "api-server" {
		t.Errorf("Expected api-server but got %s.\n", entry.Component)
	}
	if entry.Host != "worker01" {
		t.Errorf("Expected api-server but got %s.\n", entry.Host)
	}
	if entry.ReqID != "req-4leuyy-5910" {
		t.Errorf("Expected req-4leuyy-5910 but got %s.\n", entry.ReqID)
	}
	if entry.Msg != "Cache cleared" {
		t.Errorf("Expected 'Cache cleared' but got %s.\n", entry.Msg)
	}

}

func TestParseInvalidLogEntry(t *testing.T) {
	invalidLine := `invalid log line`
	_, err := ParseLogEntry(invalidLine)
	if err == nil {
		t.Errorf("Expected error for invalid format but got none")
	}
}

func TestParseLogEntryBadTime(t *testing.T) {

	badTimeLine := `2025-10-23 15:17:08.636000 | WARN | api-server | host=worker01 | request_id=req-4leuyy-5910 | msg="Cache cleared"`
	_, err := ParseLogEntry(badTimeLine)
	if err == nil || !strings.Contains(err.Error(), "failed to parse time") {
		t.Errorf("Expected time parsing error, got %v", err)
	}
}
func TestParseLogFilesBadDirectory(t *testing.T) {
	path := "../logss"
	_, err := ParseLogFiles(path)
	if err == nil {
		t.Errorf("Expected 'no such directory' error but got none.")
	}
}
func TestParseLogFilesValidFile(t *testing.T) {
	tmpDir := t.TempDir()
	logContent := `2025-10-23 15:17:08.636 | INFO | api-server | host=worker01 | request_id=req-xyz | msg="Cache cleared"`
	tmpFile := filepath.Join(tmpDir, "valid.log")

	err := os.WriteFile(tmpFile, []byte(logContent+"\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp log file: %v", err)
	}

	entries, err := ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 log entry, got %d", len(entries))
	}
	if entries[0].Host != "worker01" {
		t.Errorf("Expected host=worker01, got %s", entries[0].Host)
	}
}

func TestParseLogFilesInvalidLog(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "invalid.log")

	// invalid format line
	err := os.WriteFile(tmpFile, []byte("invalid log line\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create temp invalid log file: %v", err)
	}

	entries, err := ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// no valid entries should be parsed
	if len(entries) != 0 {
		t.Errorf("Expected 0 valid entries, got %d", len(entries))
	}
}
func TestParseLogFilesUnreadableFile(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "unreadable.log")

	err := os.WriteFile(badFile, []byte("data"), 0000)
	if err != nil {
		t.Fatalf("Failed to create unreadable file: %v", err)
	}
	defer os.Chmod(badFile, 0644) // reset permission after test

	_, err = ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
func TestParseLogFilesSkipsSubdir(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "nested")
	os.Mkdir(subDir, 0755)

	logContent := `2025-10-23 15:17:08.636 | WARN | scheduler | host=worker02 | request_id=req-abc | msg="Job delayed"`
	tmpFile := filepath.Join(tmpDir, "log1.log")
	os.WriteFile(tmpFile, []byte(logContent+"\n"), 0644)

	entries, err := ParseLogFiles(tmpDir)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if len(entries) != 1 {
		t.Errorf("Expected 1 entry (subdir skipped), got %d", len(entries))
	}
}
