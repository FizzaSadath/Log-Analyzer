package main

import (
	"log"
	"log_analyzer/pkg/database"
	"log_analyzer/pkg/web"
)

const dbUrl = "postgresql:///logsdb?host=/var/run/postgresql/"

func main() {
	// Connect DB
	db, err := database.CreateDB(dbUrl)
	if err != nil {
		log.Fatal("DB connection failed:", err)
	}

	// Build router
	r := web.SetupRouter(db)

	// Start server
	log.Println("Server running at http://localhost:8080")
	r.Run(":8080")
}
