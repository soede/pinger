package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"github.com/docker/docker/client"
)

func main() {
	config := loadConfig()
	slog.SetDefault(NewLogger(&config))
	slog.Info("Launching pinger service", slog.String("configuration", config.String()))

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	go func() {
		<-sigChan
		slog.Info("Interrupt signal received, shutdown...", slog.String("module", "pinger.main"))
		cancel()
	}()

	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		slog.Error("Error creating Docker client",
			slog.String("module", "pinger.main"),
			slog.String("message", err.Error()))
	}
	ticker := time.NewTicker(config.PingInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			slog.Info("Completing the main ping cycle", slog.String("module", "pinger.main"))
			return
		case <-ticker.C:
			checkContainers(ctx, dockerClient, config)
		}
	}
}

func checkContainers(ctx context.Context, client *client.Client, config Config) {
	containers, err := fetchContainers(ctx, client)
	if err != nil {
		slog.Error("Failed to fetch containers",
			slog.String("module", "pinger.checkContainers"),
			slog.String("error", err.Error()))
		return
	}
	pingResults := eachContainers(containers, config.PingInterval)
	err = sendResult(ctx, pingResults, config)
	if err != nil {
		slog.Error("Failed to send results",
			slog.String("module", "pinger.checkContainers"),
			slog.String("error", err.Error()))
	}
}
