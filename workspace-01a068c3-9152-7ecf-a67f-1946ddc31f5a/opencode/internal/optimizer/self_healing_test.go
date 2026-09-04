package optimizer

import (
	"strings"
	"testing"
)

func TestParseExecutionDiagnostics_GoErrors(t *testing.T) {
	output := `# github.com/opencode-ai/opencode/internal/auth
internal/auth/jwt.go:42:15: undefined: GenerateToken
internal/auth/jwt.go:88:2: cannot use claims (variable of type *Claims) as type jwt.Claims
FAIL	github.com/opencode-ai/opencode/internal/auth [build failed]`

	res := ParseExecutionDiagnostics(output, 1)
	if !res.HasErrors {
		t.Errorf("expected HasErrors to be true")
	}
	if len(res.DetectedErrors) != 2 {
		t.Fatalf("expected 2 errors, got %d", len(res.DetectedErrors))
	}
	if res.DetectedErrors[0].File != "internal/auth/jwt.go" {
		t.Errorf("expected File 'internal/auth/jwt.go', got '%s'", res.DetectedErrors[0].File)
	}
	if res.DetectedErrors[0].Line != 42 {
		t.Errorf("expected Line 42, got %d", res.DetectedErrors[0].Line)
	}
	if !strings.Contains(res.AutoFixPrompt, "Autonomous Self-Healing Diagnostic Agent") {
		t.Errorf("expected auto fix prompt to contain 'Autonomous Self-Healing Diagnostic Agent'")
	}
	if !strings.Contains(res.AutoFixPrompt, "internal/auth/jwt.go:42") {
		t.Errorf("expected auto fix prompt to contain error line")
	}
}

func TestParseExecutionDiagnostics_NoError(t *testing.T) {
	output := `=== RUN   TestLsTool_Info
--- PASS: TestLsTool_Info (0.00s)
PASS
ok  	github.com/opencode-ai/opencode/internal/llm/tools	0.026s`

	res := ParseExecutionDiagnostics(output, 0)
	if res.HasErrors {
		t.Errorf("expected HasErrors to be false")
	}
}
