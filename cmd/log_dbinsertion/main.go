package main

import (
	"log/slog"
	"log_analyzer/pkg/database"
	"log_analyzer/pkg/parser"
	"os"
)

const dbUrl = "postgresql:///log_analyzer?host=/var/run/postgresql/"

func commandHandler(args []string) error {
	db, err := database.CreateDB(dbUrl)
	if err != nil {
		return err
	}
	switch args[0] {
	case "init":
		err := database.InitDB(db)
		if err != nil {
			return err
		}
	case "add":
		dirPath := args[1]
		if dirPath == "" {
			slog.Error("Specify directory!")
		}
		entries, err := parser.ParseLogFiles(dirPath)
		if err != nil {
			return err
		}
		for _, entry := range entries {
			dbEntry := database.Entry{
				TimeStamp: entry.Time,
				Level:     string(entry.Level),
				Component: entry.Component,
				Host:      entry.Host,
				RequestId: entry.ReqID,
				Message:   entry.Msg,
			}
			err := database.AddDB(db, dbEntry)
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func main() {
	err := commandHandler(os.Args[1:])
	if err != nil {
		slog.Error("Error in invocation", "error", err)
		os.Exit(-1)
	}
}
