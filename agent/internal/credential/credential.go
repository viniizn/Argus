package credential

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

const (
	filePermissions = 0o600
	dirPermissions  = 0o700
)

var ErrNotFound = errors.New("credential: not found")

type Secret string

func (Secret) String() string   { return "<redacted>" }
func (Secret) GoString() string { return "credential.Secret(<redacted>)" }

func (s Secret) Reveal() string { return string(s) }

type stored struct {
	Credential string `json:"credential"`
}

func Load(path string) (Secret, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("credential: read file: %w", err)
	}

	var s stored
	if err := json.Unmarshal(data, &s); err != nil {
		return "", fmt.Errorf("credential: parse file: %w", err)
	}
	if s.Credential == "" {
		return "", errors.New("credential: file is empty")
	}
	return Secret(s.Credential), nil
}

func Save(path string, secret Secret) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, dirPermissions); err != nil {
		return fmt.Errorf("credential: create directory: %w", err)
	}

	data, err := json.Marshal(stored{Credential: secret.Reveal()})
	if err != nil {
		return fmt.Errorf("credential: encode: %w", err)
	}

	tmp, err := os.CreateTemp(dir, ".credential-*")
	if err != nil {
		return fmt.Errorf("credential: create temp file: %w", err)
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)

	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return fmt.Errorf("credential: write temp file: %w", err)
	}
	if err := tmp.Chmod(filePermissions); err != nil {
		tmp.Close()
		return fmt.Errorf("credential: set permissions: %w", err)
	}
	if err := tmp.Sync(); err != nil {
		tmp.Close()
		return fmt.Errorf("credential: sync: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("credential: close temp file: %w", err)
	}

	if err := os.Rename(tmpName, path); err != nil {
		return fmt.Errorf("credential: rename: %w", err)
	}
	return nil
}