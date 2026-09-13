package model

import "time"

type SystemInfo struct {
	Hostname string
	OS string
	Architecture string
	IPAddress string
	AgentVersion string
	CollectedAt time.Time
}