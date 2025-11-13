package main

import (
	"context"
	"flag"
	"log/slog"
	"log_analyzer/dbwriter"
	"log_analyzer/pkg/segmenter"
	"os"

	"github.com/jackc/pgx/v5"
)

func main() {
	logPath := flag.String("path", "/home/fizza/goProjects/log_analyzer/logs", "Path to the log directory")
	flag.Parse()
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("Unable to connect to database", "error", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)
	logStore, err := segmenter.ParseLogSegments(*logPath)
	if err != nil {
		slog.Error("Failed to parse logs", "error", err)

	}

	err = dbwriter.InsertLogs(ctx, conn, logStore)
	if err != nil {
		slog.Error("Failed to insert logs", "error", err)
	} else {
		slog.Info("All logs inserted successfully")
	}
}
