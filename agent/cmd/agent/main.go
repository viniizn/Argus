package main

import (
    "context"
    "fmt"
    "log/slog"
    "os"
    "os/signal"
    "sync"
    "syscall"
    "time"

    "github.com/viniizn/Argus/agent/internal/collector"
    "github.com/viniizn/Argus/agent/internal/config"
    "github.com/viniizn/Argus/agent/internal/heartbeat"
    "github.com/viniizn/Argus/agent/internal/identity"
    "github.com/viniizn/Argus/agent/internal/transport"
	"github.com/viniizn/Argus/agent/internal/fim"
	"github.com/viniizn/Argus/agent/internal/model"
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
	cfg, err := config.Load("./data/config.json")
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger.Info("agent configuration loaded",
		"serverUrl", cfg.ServerURL,
		"logLevel", cfg.LogLevel,
		"heartbeatInterval", cfg.HeartbeatInterval,
	)

	id, err := identity.Load(cfg.IdentityPath)
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

	events := make(chan model.Event, 100)

	fimCfg := fim.Config{
		Paths:    []string{"./data"},
		Interval: 10 * time.Second,
	}
	fimRunner := fim.New(id.AgentID, fimCfg, logger)

	client := transport.New(cfg.ServerURL)

	hb := heartbeat.New(
		client,
		id.AgentID,
		info,
		time.Duration(cfg.HeartbeatInterval)*time.Second,
		logger,
	)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		fimRunner.Run(ctx, events)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		hb.Run(ctx)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case event := <-events:
				logger.Info("event received",
					"id", event.ID,
					"type", event.Type,
					"severity", event.Severity,
					"agentId", event.AgentID,
				)
			case <-ctx.Done():
				return
			}
		}
	}()

	<-ctx.Done()
	logger.Info("shutdown signal received, waiting for goroutines")
	wg.Wait()

	return nil
}