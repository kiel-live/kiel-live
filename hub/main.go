package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/hub/server"
	"github.com/kiel-live/kiel-live/hub/store"
	"github.com/kiel-live/kiel-live/pkg/version"
)

func main() {
	level := slog.LevelInfo
	if os.Getenv("LOG") == "debug" {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	if err := godotenv.Load(); err != nil {
		slog.Debug("No .env file found")
	}

	slog.Info("kiel-live hub", "version", version.Version)

	if tz := os.Getenv("TZ"); tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			slog.Error("failed to load timezone", "tz", tz, "error", err)
			os.Exit(1)
		}
		time.Local = loc
	}

	token := os.Getenv("COLLECTOR_TOKEN")
	if token == "" {
		slog.Error("please provide COLLECTOR_TOKEN")
		os.Exit(1)
	}

	listen := os.Getenv("HUB_LISTEN")
	if listen == "" {
		listen = ":4568"
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	st := store.New()
	st.Start(ctx)
	srv := server.New(st, token)

	mux := http.NewServeMux()
	srv.RegisterRoutes(mux)

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	slog.Info("hub listening", "addr", listen)
	if err := http.ListenAndServe(listen, mux); err != nil {
		slog.Error("hub server failed", "error", err)
		os.Exit(1)
	}
}
