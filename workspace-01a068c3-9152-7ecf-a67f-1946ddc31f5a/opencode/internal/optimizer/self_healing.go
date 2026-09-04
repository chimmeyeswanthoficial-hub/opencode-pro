package optimizer

import (
	"fmt"
	"regexp"
	"strings"
)

// DiagnosticError represents a compiler or test failure error
type DiagnosticError struct {
	File      string
	Line      int
	Column    int
	Message   string
	ErrorType string // "compile", "test", "lint", "runtime"
	Snippet   string
}

// SelfHealingResult contains the diagnostic analysis and auto-repair prompt
type SelfHealingResult struct {
	HasErrors      bool
	DetectedErrors []DiagnosticError
	Summary        string
	AutoFixPrompt  string
}

// Patterns for common compilers and test runners
var (
	// Go compiler error: internal/auth/jwt.go:42:15: undefined: GenerateToken
	goErrorRegex = regexp.MustCompile(`(?m)^([a-zA-Z0-9_\-\./\\]+\.go):(\d+):(?:(\d+):)?\s*(.+)$`)

	// TypeScript/JS error: src/auth.ts:15:7 - error TS2304: Cannot find name 'jwt'
	tsErrorRegex = regexp.MustCompile(`(?m)^([a-zA-Z0-9_\-\./\\]+\.tsx?):(\d+):(\d+)\s*-\s*error\s*([A-Z0-9]+):\s*(.+)$`)

	// Python pytest/traceback error: File "test_auth.py", line 25, in test_login
	pyErrorRegex = regexp.MustCompile(`(?m)File "([^"]+\.py)", line (\d+), in (.+)`)

	// Rust compiler error: --> src/main.rs:12:5
	rustErrorRegex = regexp.MustCompile(`(?m)-->\s*([a-zA-Z0-9_\-\./\\]+\.rs):(\d+):(\d+)`)
)

// ParseExecutionDiagnostics inspects execution output / test failure logs and extracts structured diagnostics
func ParseExecutionDiagnostics(rawOutput string, exitCode int) SelfHealingResult {
	if exitCode == 0 && !strings.Contains(strings.ToLower(rawOutput), "fail") && !strings.Contains(strings.ToLower(rawOutput), "error") {
		return SelfHealingResult{HasErrors: false}
	}

	var errors []DiagnosticError

	// 1. Scan for Go errors
	for _, match := range goErrorRegex.FindAllStringSubmatch(rawOutput, -1) {
		if len(match) >= 5 {
			lineNum := 0
			fmt.Sscanf(match[2], "%d", &lineNum)
			errors = append(errors, DiagnosticError{
				File:      match[1],
				Line:      lineNum,
				Message:   match[4],
				ErrorType: "compile/test",
			})
		}
	}

	// 2. Scan for TypeScript errors
	for _, match := range tsErrorRegex.FindAllStringSubmatch(rawOutput, -1) {
		if len(match) >= 6 {
			lineNum := 0
			fmt.Sscanf(match[2], "%d", &lineNum)
			errors = append(errors, DiagnosticError{
				File:      match[1],
				Line:      lineNum,
				Message:   fmt.Sprintf("[%s] %s", match[4], match[5]),
				ErrorType: "typescript",
			})
		}
	}

	// 3. Scan for Python errors
	for _, match := range pyErrorRegex.FindAllStringSubmatch(rawOutput, -1) {
		if len(match) >= 4 {
			lineNum := 0
			fmt.Sscanf(match[2], "%d", &lineNum)
			errors = append(errors, DiagnosticError{
				File:      match[1],
				Line:      lineNum,
				Message:   fmt.Sprintf("Exception in %s", match[3]),
				ErrorType: "python",
			})
		}
	}

	// 4. Scan for Rust errors
	for _, match := range rustErrorRegex.FindAllStringSubmatch(rawOutput, -1) {
		if len(match) >= 4 {
			lineNum := 0
			fmt.Sscanf(match[2], "%d", &lineNum)
			errors = append(errors, DiagnosticError{
				File:      match[1],
				Line:      lineNum,
				Message:   "Rust compiler diagnostics",
				ErrorType: "rust",
			})
		}
	}

	if len(errors) == 0 && exitCode != 0 {
		lines := strings.Split(strings.TrimSpace(rawOutput), "\n")
		lastFew := lines
		if len(lines) > 5 {
			lastFew = lines[len(lines)-5:]
		}
		errors = append(errors, DiagnosticError{
			Message:   strings.Join(lastFew, " "),
			ErrorType: "runtime/test-failure",
		})
	}

	if len(errors) == 0 {
		return SelfHealingResult{HasErrors: false}
	}

	// Synthesize Auto-Repair Golden Prompt
	var fixPrompt strings.Builder
	fixPrompt.WriteString("### Role & Objective\n")
	fixPrompt.WriteString("Act as an Autonomous Self-Healing Diagnostic Agent. The previous execution encountered errors during build or test execution. Diagnose and repair the root causes immediately.\n\n")

	fixPrompt.WriteString("### Detected Diagnostics & Failure Summary\n")
	for i, err := range errors {
		if err.File != "" {
			fixPrompt.WriteString(fmt.Sprintf("%d. `%s:%d` - %s\n", i+1, err.File, err.Line, err.Message))
		} else {
			fixPrompt.WriteString(fmt.Sprintf("%d. %s\n", i+1, err.Message))
		}
	}
	fixPrompt.WriteString("\n")

	fixPrompt.WriteString("### Execution Instructions\n")
	fixPrompt.WriteString("1. Inspect the offending files using `view` and `grep`.\n")
	fixPrompt.WriteString("2. Apply surgical fixes with `edit` or `patch` without removing necessary functionality.\n")
	fixPrompt.WriteString("3. Re-run tests via `bash` to verify that all errors are resolved.\n")

	return SelfHealingResult{
		HasErrors:      true,
		DetectedErrors: errors,
		Summary:        fmt.Sprintf("Detected %d failure diagnostics. Initiating auto-repair loop.", len(errors)),
		AutoFixPrompt:  fixPrompt.String(),
	}
}
