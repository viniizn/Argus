package model

import "time"

type EventType string

const (
	EventFileCreated EventType = "FILE_CREATED"
	EventFileModified EventType = "FILE_MODIFIED"
	EventFileDeleted EventType = "FILE_DELETED"
)

type Severity string

const (
	SeverityLow Severity = "LOW"
	SeverityMedium Severity = "MEDIUM"
	SeverityHigh Severity = "HIGH"
	SeverityCritical Severity = "CRITICAL"
)

type Event struct {
	ID string `json:"id"`
	AgentID string `json:"agentId"`
	Type EventType `json:"type"`
	Severity Severity `json:"severity"`
	Timestamp time.Time `json:"timestamp"`
	Payload map[string]any `json:"payload"`
}