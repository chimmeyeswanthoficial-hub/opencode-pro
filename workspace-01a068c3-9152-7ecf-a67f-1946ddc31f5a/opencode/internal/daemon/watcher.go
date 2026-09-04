package daemon

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/opencode-ai/opencode/internal/app"
	"github.com/opencode-ai/opencode/internal/logging"
	"github.com/opencode-ai/opencode/internal/optimizer"
)

// ObsidianTaskWatcher monitors an Obsidian Vault directory for new task files and approvals
type ObsidianTaskWatcher struct {
	vaultPath string
	app       *app.App
	optimizer *optimizer.PromptOptimizer
	pollRate  time.Duration
}

// NewObsidianTaskWatcher creates a task watcher instance
func NewObsidianTaskWatcher(vaultPath string, app *app.App) *ObsidianTaskWatcher {
	return &ObsidianTaskWatcher{
		vaultPath: vaultPath,
		app:       app,
		optimizer: optimizer.NewPromptOptimizer("."),
		pollRate:  2 * time.Second,
	}
}

// Start begins the watcher loop
func (w *ObsidianTaskWatcher) Start(ctx context.Context) {
	inboxDir := filepath.Join(w.vaultPath, "01-Tasks", "Inbox")
	waitingDir := filepath.Join(w.vaultPath, "01-Tasks", "Waiting-Approval")
	inProgressDir := filepath.Join(w.vaultPath, "01-Tasks", "In-Progress")
	completedDir := filepath.Join(w.vaultPath, "01-Tasks", "Completed")

	// Ensure directories exist
	_ = os.MkdirAll(inboxDir, 0755)
	_ = os.MkdirAll(waitingDir, 0755)
	_ = os.MkdirAll(inProgressDir, 0755)
	_ = os.MkdirAll(completedDir, 0755)

	logging.Info(fmt.Sprintf("Obsidian Task Watcher active on: %s", w.vaultPath))

	ticker := time.NewTicker(w.pollRate)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			w.checkInbox(ctx, inboxDir, waitingDir, inProgressDir, completedDir)
			w.checkWaitingApprovals(ctx, waitingDir, inProgressDir, completedDir)
		}
	}
}

func (w *ObsidianTaskWatcher) checkInbox(ctx context.Context, inboxDir, waitingDir, inProgressDir, completedDir string) {
	entries, err := os.ReadDir(inboxDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}

		inboxFile := filepath.Join(inboxDir, entry.Name())
		contentBytes, err := os.ReadFile(inboxFile)
		if err != nil {
			continue
		}
		rawContent := string(contentBytes)
		if strings.TrimSpace(rawContent) == "" {
			continue
		}

		// Check if Human-in-the-Loop approval is required
		if strings.Contains(rawContent, "require_human_approval: true") {
			waitingFile := filepath.Join(waitingDir, entry.Name())
			approvalCard := fmt.Sprintf("%s\n\n---\n### 🛑 Action Requires Human Approval\n- **Status:** Pending Review\n\n- [ ] Approve Execution\n", rawContent)
			_ = os.WriteFile(waitingFile, []byte(approvalCard), 0644)
			_ = os.Remove(inboxFile)
			logging.Info(fmt.Sprintf("Obsidian Task routed to Waiting-Approval: %s", entry.Name()))
			continue
		}

		// Execute directly
		w.processTask(ctx, inboxFile, inProgressDir, completedDir, entry.Name(), rawContent)
	}
}

func (w *ObsidianTaskWatcher) checkWaitingApprovals(ctx context.Context, waitingDir, inProgressDir, completedDir string) {
	entries, err := os.ReadDir(waitingDir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			continue
		}

		waitingFile := filepath.Join(waitingDir, entry.Name())
		contentBytes, err := os.ReadFile(waitingFile)
		if err != nil {
			continue
		}
		rawContent := string(contentBytes)

		// Check if checkbox has been checked by human: - [x] Approve Execution
		if strings.Contains(rawContent, "- [x] Approve Execution") || strings.Contains(rawContent, "- [X] Approve Execution") {
			logging.Info(fmt.Sprintf("Human Approval granted in Obsidian for: %s", entry.Name()))
			w.processTask(ctx, waitingFile, inProgressDir, completedDir, entry.Name(), rawContent)
		}
	}
}

func (w *ObsidianTaskWatcher) processTask(ctx context.Context, sourceFile, inProgressDir, completedDir, fileName, rawContent string) {
	inProgressFile := filepath.Join(inProgressDir, fileName)
	completedFile := filepath.Join(completedDir, fileName)

	if err := os.Rename(sourceFile, inProgressFile); err != nil {
		return
	}

	logging.Info(fmt.Sprintf("Executing Obsidian task: %s", fileName))

	// Optimize prompt
	optResult := w.optimizer.Optimize(ctx, rawContent)

	// Execute with OpenCode
	startTime := time.Now()
	execErr := w.app.RunNonInteractive(ctx, optResult.OptimizedPrompt, "text", true)
	duration := time.Since(startTime)

	// Append Execution Results
	var report strings.Builder
	report.WriteString(rawContent)
	report.WriteString("\n\n---\n")
	report.WriteString(fmt.Sprintf("## ⚡ OpenCode Execution Report\n"))
	report.WriteString(fmt.Sprintf("- **Executed At:** %s\n", time.Now().Format(time.RFC3339)))
	report.WriteString(fmt.Sprintf("- **Duration:** %s\n", duration.Round(time.Millisecond)))
	if execErr != nil {
		report.WriteString(fmt.Sprintf("- **Status:** ❌ Failed (%v)\n", execErr))
	} else {
		report.WriteString("- **Status:** ✅ Completed Successfully\n")
	}
	report.WriteString("\n### Introspected Reasoning\n")
	report.WriteString(fmt.Sprintf("```\n%s\n```\n", optResult.Reasoning))

	_ = os.WriteFile(completedFile, []byte(report.String()), 0644)
	_ = os.Remove(inProgressFile)
	logging.Info(fmt.Sprintf("Obsidian task completed: %s -> %s", fileName, completedFile))
}
