package fim

import (
	"context"
	"log/slog"
	"time"

	"github.com/viniizn/Argus/agent/internal/model"
)

type Config struct {
	Paths    []string
	Interval time.Duration
}

type Runner struct {
	scanner *scanner
	cfg     Config
	logger  *slog.Logger
}

func New(agentID string, cfg Config, logger *slog.Logger) *Runner {
	s := newState()
	sc := newScanner(agentID, cfg.Paths, s, logger)

	return &Runner{
		scanner: sc,
		cfg:     cfg,
		logger:  logger,
	}
}

func (r *Runner) Run(ctx context.Context, events chan<- model.Event) {
	r.logger.Info("fim: starting",
		"paths", r.cfg.Paths,
		"interval", r.cfg.Interval.String(),
	)

	r.runScan(events)

	ticker := time.NewTicker(r.cfg.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.runScan(events)
		case <-ctx.Done():
			r.logger.Info("fim: stopped")
			return
		}
	}
}

func (r *Runner) runScan(events chan<- model.Event) {
	detected := r.scanner.scan()

	for _, event := range detected {
		r.logger.Info("fim: event detected",
			"type", event.Type,
			"severity", event.Severity,
			"path", event.Payload["path"],
		)

		select {
		case events <- event:
		default:
			r.logger.Warn("fim: event channel full, dropping event",
				"type", event.Type,
				"path", event.Payload["path"],
			)
		}
	}
}
