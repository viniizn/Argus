package enroll

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/viniizn/Argus/agent/internal/credential"
	"github.com/viniizn/Argus/agent/internal/model"
	"github.com/viniizn/Argus/agent/internal/retry"
	"github.com/viniizn/Argus/agent/internal/transport"
)

type request struct {
	Token        string `json:"token"`
	AgentID      string `json:"agentId"`
	Hostname     string `json:"hostname"`
	OS           string `json:"os"`
	Architecture string `json:"architecture"`
	IPAddress    string `json:"ipAddress"`
	AgentVersion string `json:"agentVersion"`
}

type response struct {
	Credential string `json:"credential"`
}

// Enroll exchanges a one-time enrollment token for a credential.
func Enroll(ctx context.Context, client *transport.Client, logger *slog.Logger,
	token, agentID string, info model.SystemInfo) (credential.Secret, error) {

	body := request{
		Token: token, AgentID: agentID, Hostname: info.Hostname, OS: info.OS,
		Architecture: info.Architecture, IPAddress: info.IPAddress,
		AgentVersion: info.AgentVersion,
	}

	var out credential.Secret
	err := retry.Do(ctx, retry.DefaultConfig(), logger, func() error {
		resp, err := client.Post(ctx, "/api/v1/agents/enroll", body)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		switch {
		case resp.StatusCode == http.StatusCreated:
			var r response
			if err := json.NewDecoder(http.MaxBytesReader(nil, resp.Body, 4096)).Decode(&r); err != nil || r.Credential == "" {
				return &retry.NonRetryableError{Cause: fmt.Errorf("invalid enroll response")}
			}
			out = credential.Secret(r.Credential)
			return nil
		case resp.StatusCode >= 400 && resp.StatusCode < 500 && resp.StatusCode != http.StatusTooManyRequests:
			return &retry.NonRetryableError{Cause: fmt.Errorf("enroll rejected with status %d", resp.StatusCode)}
		default:
			return fmt.Errorf("enroll failed with status %d", resp.StatusCode)
		}
	})
	return out, err
}