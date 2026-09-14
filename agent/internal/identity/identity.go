package identity

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

type Identity struct {
	AgentID   string    `json:"agentId"`
	CreatedAt time.Time `json:"createdAt"`
}

const filePermissions = 0o600
const dirPermissions = 0o700

func Load(path string) (Identity, error) {
	identity, err := readFromDisk(path)
	if err == nil {
		return identity, nil
	}
	if !os.IsNotExist(err) {
		return Identity{}, fmt.Errorf("identity: read existing identity: %w", err)
	}

	identity = Identity {
		AgentID: uuid.NewString(),
		CreatedAt: time.Now().UTC(),
	}

	if err := writeToDisk(path, identity); err != nil {
		return Identity{}, fmt.Errorf("identity: persist new identity: %w", err)
	}

	return identity, nil
}

func readFromDisk(path string) (Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Identity{}, err
	}

	var identity Identity
	if err := json.Unmarshal(data, &identity); err != nil {
		return Identity{}, fmt.Errorf("parse identity file: %w", err)
	}

	return identity, nil
}

func writeToDisk(path string, identity Identity) error {
	if err := os.MkdirAll(filepath.Dir(path), dirPermissions); err != nil {
		return fmt.Errorf("create identity directory: %w", err)
	}

	data, err := json.MarshalIndent(identity, "", " ")
	if err != nil {
		return fmt.Errorf("encode identity file: %w", err)
	}

	if err := os.WriteFile(path, data, filePermissions); err != nil {
		return fmt.Errorf("write identity file: %w", err)
	}

	return nil
}