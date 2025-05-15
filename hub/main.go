package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kiel-live/kiel-live/pkg/database"
)

func main() {
	port := "4568"

	db := database.NewMemoryDatabase()
	hub := newHub(db)
	go hub.run()

	apiRouter := http.NewServeMux()
	NewAPIServer(db, hub, apiRouter)

	router := http.NewServeMux()
	router.Handle("/api/v1/", http.StripPrefix("/api/v1", apiRouter))

	fmt.Println("Backend listening on http://localhost:" + port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
