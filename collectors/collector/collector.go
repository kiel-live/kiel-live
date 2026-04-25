package collector

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kiel-live/kiel-live/pkg/client"
	"github.com/kiel-live/kiel-live/pkg/version"
)

// ExecuteFunc is the main logic of a collector, called after setup and client connection.
type ExecuteFunc func(ctx context.Context, c client.Client) error

// Options defines the configuration for a collector.
type Options struct {
	// Name of the collector.
	Name string
	// Execute is the main function of the collector.
	Execute ExecuteFunc
}

// Collector is the base collector instance.
type Collector struct {
	execute ExecuteFunc
	name    string
}

// New creates a new Collector from the given options.
func New(opt Options) *Collector {
	return &Collector{
		execute: opt.Execute,
		name:    opt.Name,
	}
}

// Run starts the collector and exits with code 1 on error.
func (c *Collector) Run() {
	if err := c.run(); err != nil {
		slog.Error("collector failed", "error", err)
		os.Exit(1)
	}
}

func (c *Collector) run() error {
	slog.Info(c.name+" collector", "version", version.Version)

	if err := godotenv.Load(); err != nil {
		slog.Debug("No .env file found")
	}

	level := slog.LevelInfo
	if os.Getenv("LOG") == "debug" {
		level = slog.LevelDebug
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level})))

	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			return fmt.Errorf("failed to load timezone from TZ environment variable: %w", err)
		}
		slog.Debug("Set time.Local to", "tz", tz)
	}

	server := os.Getenv("COLLECTOR_SERVER")
	if server == "" {
		return fmt.Errorf("please provide a server address for the collector with COLLECTOR_SERVER")
	}

	token := os.Getenv("COLLECTOR_TOKEN")
	if token == "" {
		return fmt.Errorf("please provide a token for the collector with COLLECTOR_TOKEN")
	}

	cl := client.NewClient(server, token)
	if err := cl.Connect(); err != nil {
		return err
	}
	defer func() {
		if err := cl.Disconnect(); err != nil {
			slog.Error("failed to disconnect", "error", err)
		}
	}()

	if c.execute == nil {
		panic("collector execute function is not set")
	}

	return c.execute(context.Background(), cl)
}
