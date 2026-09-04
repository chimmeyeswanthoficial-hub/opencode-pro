package chat

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/optimizer"
	"github.com/opencode-ai/opencode/internal/tui/styles"
	"github.com/opencode-ai/opencode/internal/tui/theme"
)

// OptimizerBarModel renders the bottom status bar for prompt optimization
type OptimizerBarModel struct {
	Width        int
	Enabled      bool
	HasContent   bool
	Result       optimizer.OptimizationResult
	MentionCount int
	FilesCount   int
	SkillsCount  int
	RulesCount   int
	IsThinking   bool
}

func NewOptimizerBarModel() OptimizerBarModel {
	return OptimizerBarModel{
		Enabled: true,
	}
}

// UpdateState updates the optimizer bar state based on the current editor text
func (b *OptimizerBarModel) UpdateState(text string, opt *optimizer.PromptOptimizer) {
	trimmed := strings.TrimSpace(text)
	if trimmed == "" {
		b.HasContent = false
		b.Result = optimizer.OptimizationResult{}
		b.FilesCount = 0
		b.SkillsCount = 0
		b.RulesCount = 0
		return
	}

	b.HasContent = true
	// Fast local reasoning and optimization
	if opt != nil {
		b.Result = opt.Optimize(nil, text)
		b.FilesCount = len(b.Result.TargetFiles)
		b.SkillsCount = len(b.Result.AppliedSkills)
		b.RulesCount = len(b.Result.AppliedRules)
	}
}

func (b *OptimizerBarModel) View() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle().Width(b.Width)

	if !b.Enabled {
		return baseStyle.
			Foreground(t.TextMuted()).
			Background(t.Background()).
			Render(" ⚡ Optimizer: Off (Ctrl+O to enable)")
	}

	if !b.HasContent {
		hint := baseStyle.
			Foreground(t.TextMuted()).
			Background(t.Background()).
			Render(" ⚡ Pro Prompt Engine Active • Type @ for files/tools, / for skills, # for rules")
		return hint
	}

	// Status badge
	badgeStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(t.Background()).
		Background(t.Primary()).
		Padding(0, 1)

	badge := badgeStyle.Render("⚡ OPTIMIZED")

	// Context pills
	var pills []string
	if b.FilesCount > 0 {
		pills = append(pills, fmt.Sprintf("📁 %d file(s)", b.FilesCount))
	}
	if b.SkillsCount > 0 {
		pills = append(pills, fmt.Sprintf("🛠️ %d skill(s)", b.SkillsCount))
	}
	if b.RulesCount > 0 {
		pills = append(pills, fmt.Sprintf("📜 %d rule(s)", b.RulesCount))
	}

	contextSummary := strings.Join(pills, " • ")
	if contextSummary == "" {
		contextSummary = "Auto-inferring project scope"
	}

	infoStyle := lipgloss.NewStyle().
		Foreground(t.Text()).
		Background(t.Background()).
		Padding(0, 1)

	shortcutStyle := lipgloss.NewStyle().
		Foreground(t.TextMuted()).
		Background(t.Background()).
		Padding(0, 1)

	leftPart := lipgloss.JoinHorizontal(lipgloss.Center, badge, infoStyle.Render(contextSummary))
	rightPart := shortcutStyle.Render("[Tab] Expand Prompt • [Ctrl+P] Inspect • [Ctrl+O] Toggle")

	// Calculate space in between
	totalLen := lipgloss.Width(leftPart) + lipgloss.Width(rightPart)
	spaceLen := b.Width - totalLen
	if spaceLen < 1 {
		spaceLen = 1
	}
	space := strings.Repeat(" ", spaceLen)

	barContent := lipgloss.JoinHorizontal(lipgloss.Center, leftPart, space, rightPart)
	return lipgloss.NewStyle().
		Background(t.Background()).
		Width(b.Width).
		Render(barContent)
}
