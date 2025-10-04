package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kiel-live/kiel-live/hub/api"
	"github.com/kiel-live/kiel-live/hub/hub"
	"github.com/kiel-live/kiel-live/pkg/database"
)

func main() {
	port := os.Getenv("HUB_PORT")
	if port == "" {
		port = "4568"
	}

	collectorToken := os.Getenv("HUB_COLLECTOR_TOKEN")
	if collectorToken == "" {
		log.Fatal("HUB_COLLECTOR_TOKEN is required")
	}

	db := database.NewMemoryDatabase()
	hub := hub.NewHub(db)
	go hub.Run()

	apiRouter := http.NewServeMux()
	api.NewAPIServer(db, hub, apiRouter, collectorToken)

	router := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	fmt.Println("Backend listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
