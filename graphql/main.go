package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/kiel-live/kiel-live/graphql/graph"
	"github.com/kiel-live/kiel-live/shared/database"
	"github.com/kiel-live/kiel-live/shared/hub"
	"github.com/kiel-live/kiel-live/shared/pubsub"
	"github.com/kiel-live/kiel-live/testing/usage"
)

const defaultPort = "4567"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	db := database.NewMemoryDatabase()
	if err := db.Open(); err != nil {
		log.Fatal(err)
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		Hub: &hub.Hub{
			DB:     db,
			PubSub: pubsub.NewMemory(),
		},
	}}))
	srv.AddTransport(&transport.Websocket{})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	u := usage.NewUsage("graphql")
	http.HandleFunc("/perf", func(w http.ResponseWriter, r *http.Request) {
		amountClients := r.URL.Query().Get("amountClients")
		if amountClients == "" {
			http.Error(w, "amountClients is required", http.StatusBadRequest)
			return
		}

		go u.Collect(amountClients)
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
