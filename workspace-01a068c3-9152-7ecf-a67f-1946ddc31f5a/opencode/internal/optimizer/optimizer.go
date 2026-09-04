package optimizer

import (
	"context"
	"fmt"
	"strings"
)

// OptimizationResult contains the transformed prompt and diagnostic reasoning
type OptimizationResult struct {
	OriginalPrompt   string
	OptimizedPrompt  string
	Reasoning        string
	DetectedIntent   string
	TargetFiles      []string
	RecommendedTools []string
	AppliedSkills    []string
	AppliedRules     []string
	TestCommand      string
	IsEnhanced       bool
}

// PromptOptimizer encapsulates reasoning and prompt synthesis
type PromptOptimizer struct {
	projectInfo ProjectInfo
}

// NewPromptOptimizer creates a prompt optimizer instance
func NewPromptOptimizer(dir string) *PromptOptimizer {
	return &PromptOptimizer{
		projectInfo: IntrospectProject(dir),
	}
}

// Optimize takes a raw user prompt, parses mentions, and generates an enhanced prompt
func (p *PromptOptimizer) Optimize(ctx context.Context, rawPrompt string) OptimizationResult {
	trimmed := strings.TrimSpace(rawPrompt)
	if trimmed == "" {
		return OptimizationResult{OriginalPrompt: rawPrompt}
	}

	// 1. Parse Mentions (@, /, #)
	mentionContext := ParsePromptMentions(trimmed)

	// 2. Clean Raw Goal Text (remove @, /, # tokens for clean intent extraction)
	words := strings.Fields(trimmed)
	var goalWords []string
	for _, w := range words {
		if !strings.HasPrefix(w, "@") && !strings.HasPrefix(w, "/") && !strings.HasPrefix(w, "#") {
			goalWords = append(goalWords, w)
		}
	}
	cleanGoal := strings.Join(goalWords, " ")
	if cleanGoal == "" {
		cleanGoal = rawPrompt
	}

	// 3. Determine Intent & Missing Gaps
	var reasoningSteps []string
	reasoningSteps = append(reasoningSteps, fmt.Sprintf("1. Introspected project stack: %s (%s)", p.projectInfo.PrimaryLanguage, p.projectInfo.Framework))

	targetFiles := mentionContext.Files
	if len(targetFiles) > 0 {
		reasoningSteps = append(reasoningSteps, fmt.Sprintf("2. Resolved %d target file(s): %s", len(targetFiles), strings.Join(targetFiles, ", ")))
	} else if len(p.projectInfo.ModifiedFiles) > 0 {
		reasoningSteps = append(reasoningSteps, fmt.Sprintf("2. Inferred target from active git modifications: %s", strings.Join(p.projectInfo.ModifiedFiles, ", ")))
	} else {
		reasoningSteps = append(reasoningSteps, "2. No explicit target files mentioned; instructed agent to inspect workspace.")
	}

	var appliedSkills []string
	for _, s := range mentionContext.Skills {
		appliedSkills = append(appliedSkills, fmt.Sprintf("%s (%s): %s", s.Name, s.Category, s.Template))
		reasoningSteps = append(reasoningSteps, fmt.Sprintf("3. Applied skill workflow: %s", s.Name))
	}

	var appliedRules []string
	for _, r := range mentionContext.Rules {
		appliedRules = append(appliedRules, fmt.Sprintf("%s - %s", r.Tag, r.Description))
	}
	if p.projectInfo.HasRulesFile && len(appliedRules) == 0 {
		appliedRules = append(appliedRules, "#rules (repository conventions)")
	}
	if len(appliedRules) > 0 {
		reasoningSteps = append(reasoningSteps, fmt.Sprintf("4. Attached project constraints: %s", strings.Join(appliedRules, ", ")))
	}

	tools := mentionContext.Tools
	if len(tools) == 0 {
		// Heuristically suggest essential tools
		tools = []string{"edit", "grep", "view", "bash"}
	}

	for _, mcp := range mentionContext.McpServers {
		tools = append(tools, fmt.Sprintf("mcp:%s", mcp))
	}
	reasoningSteps = append(reasoningSteps, fmt.Sprintf("5. Selected execution tools: %s", strings.Join(tools, ", ")))

	// Formulate Detected Intent
	intent := cleanGoal
	if len(mentionContext.Skills) > 0 {
		intent = fmt.Sprintf("Execute %s workflow on request: '%s'", mentionContext.Skills[0].Name, cleanGoal)
	}

	// 4. Synthesize Golden Prompt
	goldenPrompt := FormatGoldenPrompt(
		cleanGoal,
		intent,
		p.projectInfo,
		targetFiles,
		tools,
		appliedSkills,
		appliedRules,
		mentionContext.RuleContent,
	)

	return OptimizationResult{
		OriginalPrompt:   rawPrompt,
		OptimizedPrompt:  goldenPrompt,
		Reasoning:        strings.Join(reasoningSteps, "\n"),
		DetectedIntent:   intent,
		TargetFiles:      targetFiles,
		RecommendedTools: tools,
		AppliedSkills:    appliedSkills,
		AppliedRules:     appliedRules,
		TestCommand:      p.projectInfo.TestCommand,
		IsEnhanced:       true,
	}
}
