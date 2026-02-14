package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kiel-live/kiel-live/gateway/api"
	"github.com/kiel-live/kiel-live/gateway/database"
	"github.com/kiel-live/kiel-live/gateway/hub"
	"github.com/kiel-live/kiel-live/gateway/search"
)

func main() {
	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "4568"
	}

	collectorToken := os.Getenv("GATEWAY_COLLECTOR_TOKEN")
	if collectorToken == "" {
		log.Fatal("GATEWAY_COLLECTOR_TOKEN is required")
	}

	db := database.NewMemoryDatabase()

	search := search.NewMemorySearch()

	hub := hub.NewHub(db)
	go hub.Run()

	apiRouter := http.NewServeMux()
	api.NewAPIServer(db, search, hub, apiRouter, collectorToken)

	router := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	fmt.Println("Backend listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
