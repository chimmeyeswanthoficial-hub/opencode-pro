package optimizer

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIntrospectProject(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "opencode_opt_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create fake go.mod
	err = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module testapp\n\ngo 1.22\n"), 0644)
	require.NoError(t, err)

	info := IntrospectProject(tempDir)
	assert.Equal(t, "Go", info.PrimaryLanguage)
	assert.Equal(t, "go test -v ./...", info.TestCommand)
}

func TestParsePromptMentions(t *testing.T) {
	raw := "refactor login endpoint @internal/auth/jwt.go @tool:grep /test #rules:security #git:diff"
	ctx := ParsePromptMentions(raw)

	assert.Contains(t, ctx.Files, "internal/auth/jwt.go")
	assert.Contains(t, ctx.Tools, "grep")
	assert.NotEmpty(t, ctx.Skills)
	assert.Equal(t, "/test", ctx.Skills[0].Name)
	assert.NotEmpty(t, ctx.Rules)
}

func TestPromptOptimizer(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "opencode_opt_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	_ = os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte("module myapp\n"), 0644)

	opt := NewPromptOptimizer(tempDir)
	res := opt.Optimize(context.Background(), "create payment endpoint @payment.go /startup:api-scaffold #spec:api")

	assert.True(t, res.IsEnhanced)
	assert.Contains(t, res.TargetFiles, "payment.go")
	assert.Contains(t, res.OptimizedPrompt, "Role & Objective")
	assert.Contains(t, res.OptimizedPrompt, "payment.go")
	assert.Contains(t, res.OptimizedPrompt, "Startup Pack")
	assert.NotEmpty(t, res.Reasoning)
}
