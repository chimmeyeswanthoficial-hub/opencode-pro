package daemon

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/opencode-ai/opencode/internal/app"
	"github.com/opencode-ai/opencode/internal/logging"
	"github.com/opencode-ai/opencode/internal/optimizer"
	"github.com/opencode-ai/opencode/internal/version"
)

// DaemonServer provides HTTP & JSON-RPC API for OpenCode orchestration
type DaemonServer struct {
	app       *app.App
	optimizer *optimizer.PromptOptimizer
	port      int
	host      string
}

// NewDaemonServer creates a new daemon server
func NewDaemonServer(app *app.App, host string, port int) *DaemonServer {
	if host == "" {
		host = "0.0.0.0"
	}
	if port <= 0 {
		port = 8080
	}
	return &DaemonServer{
		app:       app,
		optimizer: optimizer.NewPromptOptimizer("."),
		port:      port,
		host:      host,
	}
}

// Start begins listening on the HTTP server
func (s *DaemonServer) Start(ctx context.Context) error {
	mux := http.NewServeMux()

	// 1. Health & Status
	mux.HandleFunc("/api/v1/health", s.handleHealth)

	// 2. Prompt Optimizer Endpoint
	mux.HandleFunc("/api/v1/optimize", s.handleOptimize)

	// 3. Mention Context Resolution Endpoint
	mux.HandleFunc("/api/v1/context", s.handleContext)

	// 4. Agent Execution Endpoint
	mux.HandleFunc("/api/v1/run", s.handleRun)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.host, s.port),
		Handler:      corsMiddleware(mux),
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 600 * time.Second,
	}

	logging.Info(fmt.Sprintf("OpenCode Headless Daemon listening on %s:%d", s.host, s.port))

	go func() {
		<-ctx.Done()
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdownCtx)
	}()

	err := server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *DaemonServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	info := optimizer.IntrospectProject(".")
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"status":      "ok",
		"version":     version.Version,
		"language":    info.PrimaryLanguage,
		"framework":   info.Framework,
		"testCommand": info.TestCommand,
		"gitBranch":   info.GitBranch,
	})
}

type OptimizeRequest struct {
	Prompt string `json:"prompt"`
}

func (s *DaemonServer) handleOptimize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OptimizeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	res := s.optimizer.Optimize(r.Context(), req.Prompt)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}

type ContextRequest struct {
	Prompt string `json:"prompt"`
}

func (s *DaemonServer) handleContext(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ContextRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	parsed := optimizer.ParsePromptMentions(req.Prompt)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(parsed)
}

type RunRequest struct {
	Prompt      string `json:"prompt"`
	AutoEnhance bool   `json:"auto_enhance"`
	Format      string `json:"format"`
}

func (s *DaemonServer) handleRun(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RunRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	finalPrompt := req.Prompt
	if req.AutoEnhance {
		res := s.optimizer.Optimize(r.Context(), req.Prompt)
		if res.OptimizedPrompt != "" {
			finalPrompt = res.OptimizedPrompt
		}
	}

	if req.Format == "" {
		req.Format = "json"
	}

	w.Header().Set("Content-Type", "application/json")
	err := s.app.RunNonInteractive(r.Context(), finalPrompt, req.Format, true)
	if err != nil {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"status": "error",
			"error":  err.Error(),
		})
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]any{
		"status": "success",
		"prompt": finalPrompt,
	})
}
