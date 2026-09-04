package logging

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

// LangfuseConfig holds observability settings
type LangfuseConfig struct {
	Host      string
	PublicKey string
	SecretKey string
	Enabled   bool
}

// GetLangfuseConfig loads Langfuse config from environment variables
func GetLangfuseConfig() LangfuseConfig {
	host := os.Getenv("LANGFUSE_HOST")
	if host == "" {
		host = "http://localhost:3000"
	}
	pk := os.Getenv("LANGFUSE_PUBLIC_KEY")
	sk := os.Getenv("LANGFUSE_SECRET_KEY")

	return LangfuseConfig{
		Host:      host,
		PublicKey: pk,
		SecretKey: sk,
		Enabled:   pk != "" && sk != "",
	}
}

// TraceEvent represents an execution step or generation trace
type TraceEvent struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	StartTime time.Time      `json:"startTime"`
	EndTime   time.Time      `json:"endTime"`
	Input     any            `json:"input"`
	Output    any            `json:"output"`
	Metadata  map[string]any `json:"metadata"`
	Status    string         `json:"status"`
}

// LogLangfuseTrace sends an asynchronous trace event to a Langfuse instance
func LogLangfuseTrace(ctx context.Context, event TraceEvent) {
	cfg := GetLangfuseConfig()
	if !cfg.Enabled {
		return
	}

	go func() {
		defer RecoverPanic("langfuse-telemetry", nil)

		payload := map[string]any{
			"batch": []map[string]any{
				{
					"type":      "trace-create",
					"id":        event.ID,
					"name":      event.Name,
					"timestamp": event.StartTime.Format(time.RFC3339Nano),
					"input":     event.Input,
					"output":    event.Output,
					"metadata":  event.Metadata,
				},
			},
		}

		body, err := json.Marshal(payload)
		if err != nil {
			return
		}

		endpoint := fmt.Sprintf("%s/api/public/ingestion", cfg.Host)
		req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, endpoint, bytes.NewReader(body))
		if err != nil {
			return
		}

		req.SetBasicAuth(cfg.PublicKey, cfg.SecretKey)
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err == nil && resp != nil {
			_ = resp.Body.Close()
		}
	}()
}
