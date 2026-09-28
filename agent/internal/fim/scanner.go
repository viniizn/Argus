package fim

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/viniizn/Argus/agent/internal/model"
)

type scanner struct {
	agentID string
	paths   []string
	state   *state
	logger  *slog.Logger
}

func newScanner(agentID string, paths []string, state *state, logger *slog.Logger) *scanner {
	return &scanner{
		agentID: agentID,
		paths:   paths,
		state:   state,
		logger:  logger,
	}
}

func (s *scanner) scan() []model.Event {
	var events []model.Event

	for _, root := range s.paths {
		info, err := os.Stat(root)
		if err != nil {
			if os.IsNotExist(err) {
				events = append(events, s.handleDeleted(root))
				continue
			}
			s.logger.Warn("fim: failed to stat path",
				"path", root,
				"error", err,
			)
			continue
		}

		if info.IsDir() {
			dirEvents := s.scanDir(root)
			events = append(events, dirEvents...)
		} else {
			if event, ok := s.scanFile(root); ok {
				events = append(events, event)
			}
		}
	}

	return events
}

func (s *scanner) scanDir(root string) []model.Event {
	var events []model.Event

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			s.logger.Warn("fim: error walking path",
				"path", path,
				"error", err,
			)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		if event, ok := s.scanFile(path); ok {
			events = append(events, event)
		}

		return nil
	})

	if err != nil {
		s.logger.Warn("fim: failed to walk directory",
			"path", root,
			"error", err,
		)
	}

	return events
}

func (s *scanner) scanFile(path string) (model.Event, bool) {
	hash, err := hashFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return s.handleDeleted(path), true
		}
		s.logger.Warn("fim: failed to hash file",
			"path", path,
			"error", err,
		)
		return model.Event{}, false
	}

	previous, known := s.state.get(path)
	s.state.set(path, hash)

	if !known {
		return s.newEvent(model.EventFileCreated, model.SeverityLow, path, map[string]any{
			"path":    path,
			"newHash": hash,
		}), true
	}

	if hash != previous {
		return s.newEvent(model.EventFileModified, model.SeverityHigh, path, map[string]any{
			"path":    path,
			"oldHash": previous,
			"newHash": hash,
		}), true
	}

	return model.Event{}, false
}

func (s *scanner) handleDeleted(path string) model.Event {
	_, wasKnown := s.state.get(path)
	s.state.delete(path)

	severity := model.SeverityHigh
	if !wasKnown {
		severity = model.SeverityLow
	}

	return s.newEvent(model.EventFileDeleted, severity, path, map[string]any{
		"path": path,
	})
}

func (s *scanner) newEvent(
	eventType model.EventType,
	severity model.Severity,
	path string,
	payload map[string]any,
) model.Event {
	return model.Event{
		ID:        uuid.NewString(),
		AgentID:   s.agentID,
		Type:      eventType,
		Severity:  severity,
		Timestamp: time.Now().UTC(),
		Payload:   payload,
	}
}

func hashFile(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hash file: %w", err)
	}

	return hex.EncodeToString(h.Sum(nil)), nil
}