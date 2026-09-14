package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/viniizn/Argus/agent/internal/collector"
	"github.com/viniizn/Argus/agent/internal/identity"
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

const identityPath = "./data/identity.json"

func run(ctx context.Context, logger *slog.Logger) error {
	id, err := identity.Load(identityPath)
	if err != nil {
		return fmt.Errorf("load identity: %w", err)
	}

	logger.Info("agent identity loaded",
		"agentId", id.AgentID,
		"createdAt", id.CreatedAt,	
	)
	
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