package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/viniizn/Argus/agent/internal/collector"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.Info("agent starting")

	if err := run(ctx, logger); err != nil {
		logger.Error("agent exited with error", "error", err)
		os.Exit(1)
	}

	logger.Info("agent stopped cleanly")
}

func run(ctx context.Context, logger *slog.Logger) error {
	info, err := collector.SystemInfo()
	if err != nil {
		return fmt.Errorf("collect system info: %w", err)
	}

	logger.Info("system info collected",
		"hostname", info.Hostname,
		"os", info.OS,
		"architecture", info.Architecture,
		"ip", info.IPAddress,
		"agentVersion", info.AgentVersion,
	)

	<-ctx.Done()
	logger.Info("shutdown signal received")
	return nil
}