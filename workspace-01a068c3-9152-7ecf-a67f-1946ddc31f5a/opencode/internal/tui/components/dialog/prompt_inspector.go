package dialog

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opencode-ai/opencode/internal/optimizer"
	"github.com/opencode-ai/opencode/internal/tui/layout"
	"github.com/opencode-ai/opencode/internal/tui/styles"
	"github.com/opencode-ai/opencode/internal/tui/theme"
	"github.com/opencode-ai/opencode/internal/tui/util"
)

// PromptInspectorCloseMsg is sent when the inspector dialog is closed
type PromptInspectorCloseMsg struct{}

// PromptInspectorApplyMsg is sent when user chooses to replace their editor input with the enhanced prompt
type PromptInspectorApplyMsg struct {
	EnhancedPrompt string
}

// PromptInspectorRunMsg is sent when user chooses to directly run the enhanced prompt
type PromptInspectorRunMsg struct {
	EnhancedPrompt string
}

type PromptInspectorDialog interface {
	tea.Model
	layout.Bindings
	SetResult(result optimizer.OptimizationResult)
	SetSize(width, height int)
}

type promptInspectorCmp struct {
	width    int
	height   int
	result   optimizer.OptimizationResult
	viewport viewport.Model
}

type promptInspectorKeyMap struct {
	Apply  key.Binding
	Run    key.Binding
	Escape key.Binding
}

var promptInspectorKeys = promptInspectorKeyMap{
	Apply: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "insert into editor"),
	),
	Run: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "run enhanced prompt"),
	),
	Escape: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "close"),
	),
}

func (p *promptInspectorCmp) Init() tea.Cmd {
	return nil
}

func (p *promptInspectorCmp) SetResult(result optimizer.OptimizationResult) {
	p.result = result
	p.updateContent()
}

func (p *promptInspectorCmp) SetSize(width, height int) {
	p.width = width
	p.height = height
	modalWidth := width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}
	p.viewport.Width = modalWidth - 4
	p.viewport.Height = height - 12
	if p.viewport.Height < 10 {
		p.viewport.Height = 10
	}
	p.updateContent()
}

func (p *promptInspectorCmp) updateContent() {
	t := theme.CurrentTheme()
	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Primary())
	subHeaderStyle := lipgloss.NewStyle().Bold(true).Foreground(t.Text())
	codeBlockStyle := lipgloss.NewStyle().
		Border(lipgloss.NormalBorder(), false, false, false, true).
		BorderForeground(t.Primary()).
		Padding(0, 1).
		Foreground(t.Text())

	var sb strings.Builder

	sb.WriteString(headerStyle.Render("⚡ ORIGINAL USER INPUT"))
	sb.WriteString("\n")
	sb.WriteString(codeBlockStyle.Render(p.result.OriginalPrompt))
	sb.WriteString("\n\n")

	sb.WriteString(headerStyle.Render("🧠 INTROSPECTION & REASONING CHAIN"))
	sb.WriteString("\n")
	if p.result.Reasoning != "" {
		sb.WriteString(codeBlockStyle.Render(p.result.Reasoning))
	} else {
		sb.WriteString(codeBlockStyle.Render("No active mentions. Prompt will be executed with default workspace context."))
	}
	sb.WriteString("\n\n")

	sb.WriteString(subHeaderStyle.Render("🎯 TARGET CONTEXT ATTACHMENTS"))
	sb.WriteString("\n")
	var attachments []string
	if len(p.result.TargetFiles) > 0 {
		attachments = append(attachments, fmt.Sprintf("• Files: %s", strings.Join(p.result.TargetFiles, ", ")))
	}
	if len(p.result.AppliedSkills) > 0 {
		attachments = append(attachments, fmt.Sprintf("• Skills: %s", strings.Join(p.result.AppliedSkills, ", ")))
	}
	if len(p.result.AppliedRules) > 0 {
		attachments = append(attachments, fmt.Sprintf("• Rules: %s", strings.Join(p.result.AppliedRules, ", ")))
	}
	if len(attachments) > 0 {
		sb.WriteString(strings.Join(attachments, "\n"))
	} else {
		sb.WriteString("• None explicitly tagged; auto-introspection active.")
	}
	sb.WriteString("\n\n")

	sb.WriteString(headerStyle.Render("✨ SYNTHESIZED GOLDEN PROMPT (FOR AGENT)"))
	sb.WriteString("\n")
	sb.WriteString(codeBlockStyle.Render(p.result.OptimizedPrompt))

	p.viewport.SetContent(sb.String())
}

func (p *promptInspectorCmp) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, promptInspectorKeys.Escape):
			return p, util.CmdHandler(PromptInspectorCloseMsg{})
		case key.Matches(msg, promptInspectorKeys.Apply):
			return p, tea.Batch(
				util.CmdHandler(PromptInspectorApplyMsg{EnhancedPrompt: p.result.OptimizedPrompt}),
				util.CmdHandler(PromptInspectorCloseMsg{}),
			)
		case key.Matches(msg, promptInspectorKeys.Run):
			return p, tea.Batch(
				util.CmdHandler(PromptInspectorRunMsg{EnhancedPrompt: p.result.OptimizedPrompt}),
				util.CmdHandler(PromptInspectorCloseMsg{}),
			)
		}
	}

	var cmd tea.Cmd
	p.viewport, cmd = p.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return p, tea.Batch(cmds...)
}

func (p *promptInspectorCmp) View() string {
	t := theme.CurrentTheme()
	baseStyle := styles.BaseStyle()

	modalWidth := p.width - 10
	if modalWidth > 90 {
		modalWidth = 90
	}

	title := baseStyle.
		Foreground(t.Primary()).
		Bold(true).
		Width(modalWidth).
		Padding(0, 1).
		Render("⚡ Smart Prompt Optimizer & Reasoning Inspector")

	controls := baseStyle.
		Foreground(t.TextMuted()).
		Width(modalWidth).
		Padding(0, 1).
		Render("[Tab] Insert into Editor • [Enter] Run Immediately • [Esc] Close • [↑/↓] Scroll")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		baseStyle.Width(modalWidth).Render(""),
		p.viewport.View(),
		baseStyle.Width(modalWidth).Render(""),
		controls,
	)

	return baseStyle.Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderBackground(t.Background()).
		BorderForeground(t.Primary()).
		Width(modalWidth + 4).
		Render(content)
}

func (p *promptInspectorCmp) BindingKeys() []key.Binding {
	return layout.KeyMapToSlice(promptInspectorKeys)
}

// NewPromptInspectorDialog creates a new prompt inspector dialog
func NewPromptInspectorDialog() PromptInspectorDialog {
	vp := viewport.New(80, 15)
	return &promptInspectorCmp{
		viewport: vp,
	}
}
