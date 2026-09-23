package heartbeat

import (
	"context"
	"log/slog"
	"net/http"
	"time"
	"fmt"

	"github.com/viniizn/Argus/agent/internal/model"
	"github.com/viniizn/Argus/agent/internal/retry"
	"github.com/viniizn/Argus/agent/internal/transport"
)

// Payload is the data sent to the server on each heartbeat.
type Payload struct {
	AgentID      string    `json:"agentId"`
	Hostname     string    `json:"hostname"`
	OS           string    `json:"os"`
	Architecture string    `json:"architecture"`
	IPAddress    string    `json:"ipAddress"`
	AgentVersion string    `json:"agentVersion"`
	Timestamp    time.Time `json:"timestamp"`
}

// Runner sends periodic heartbeats to the Argus server.
type Runner struct {
	client   *transport.Client
	payload  Payload
	interval time.Duration
	logger   *slog.Logger
}

// New creates a new heartbeat Runner.
func New(
	client *transport.Client,
	agentID string,
	info model.SystemInfo,
	interval time.Duration,
	logger *slog.Logger,
) *Runner {
	return &Runner{
		client: client,
		payload: Payload{
			AgentID:      agentID,
			Hostname:     info.Hostname,
			OS:           info.OS,
			Architecture: info.Architecture,
			IPAddress:    info.IPAddress,
			AgentVersion: info.AgentVersion,
		},
		interval: interval,
		logger:   logger,
	}
}

// Run starts the heartbeat loop. It sends a heartbeat immediately,
// then repeats at the configured interval. It stops when ctx is cancelled.
func (r *Runner) Run(ctx context.Context) {
	r.send(ctx)

	ticker := time.NewTicker(r.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.send(ctx)
		case <-ctx.Done():
			r.logger.Info("heartbeat stopped")
			return
		}
	}
}

func (r *Runner) send(ctx context.Context) {
	err := retry.Do(ctx, retry.DefaultConfig(), r.logger, func() error {
		r.payload.Timestamp = time.Now().UTC()

		resp, err := r.client.Post(ctx, "/api/v1/heartbeat", r.payload)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			r.logger.Info("heartbeat sent", "agentId", r.payload.AgentID)
			return nil
		}

		if resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests {
			return &retry.NonRetryableError{
				Cause: fmt.Errorf("server rejected heartbeat with status %d", resp.StatusCode),
			}
		}

		return fmt.Errorf("heartbeat failed with status %d", resp.StatusCode)
	})

	if err != nil {
		r.logger.Warn("heartbeat gave up after retries", "error", err)
	}
}