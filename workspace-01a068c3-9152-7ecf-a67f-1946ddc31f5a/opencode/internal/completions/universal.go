package completions

import (
	"fmt"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/opencode-ai/opencode/internal/config"
	"github.com/opencode-ai/opencode/internal/tui/components/dialog"
)

// UniversalContextGroup multiplexes completions for @, /, and #
type UniversalContextGroup struct {
	filesProvider  dialog.CompletionProvider
	skillsProvider dialog.CompletionProvider
	rulesProvider  dialog.CompletionProvider
}

// NewUniversalContextGroup creates a multi-trigger context provider
func NewUniversalContextGroup() *UniversalContextGroup {
	return &UniversalContextGroup{
		filesProvider:  NewFileAndFolderContextGroup(),
		skillsProvider: NewSkillsContextGroup(),
		rulesProvider:  NewRulesContextGroup(),
	}
}

func (u *UniversalContextGroup) GetId() string {
	return "universal"
}

func (u *UniversalContextGroup) GetEntry() dialog.CompletionItemI {
	return dialog.NewCompletionItem(dialog.CompletionItem{
		Title: "Universal Context (@ files/tools, / skills, # rules)",
		Value: "universal",
	})
}

// BuiltinToolsList returns built-in agent tools formatted for @ mentions
func BuiltinToolsList() []string {
	return []string{
		"@tool:edit",
		"@tool:patch",
		"@tool:write",
		"@tool:bash",
		"@tool:grep",
		"@tool:glob",
		"@tool:ls",
		"@tool:view",
		"@tool:diagnostics",
		"@tool:fetch",
		"@tool:sourcegraph",
	}
}

// GetChildEntries dynamically inspects the query prefix (@, /, #) and yields relevant matches
func (u *UniversalContextGroup) GetChildEntries(query string) ([]dialog.CompletionItemI, error) {
	trimmed := strings.TrimSpace(query)

	// Case 1: Slash Command or Skill Trigger (/)
	if strings.HasPrefix(trimmed, "/") {
		return u.skillsProvider.GetChildEntries(trimmed)
	}

	// Case 2: Project Rule, Spec or Git Trigger (#)
	if strings.HasPrefix(trimmed, "#") {
		return u.rulesProvider.GetChildEntries(trimmed)
	}

	// Case 3: Mention Files, Builtin Tools, or MCP Servers (@)
	cleanQuery := strings.TrimPrefix(trimmed, "@")

	var results []dialog.CompletionItemI

	// Include Builtin Tools
	tools := BuiltinToolsList()
	matchedTools := tools
	if cleanQuery != "" {
		matchedTools = fuzzy.Find(cleanQuery, tools)
	}
	for _, t := range matchedTools {
		results = append(results, dialog.NewCompletionItem(dialog.CompletionItem{
			Title: fmt.Sprintf("%s (Built-in Tool)", t),
			Value: t,
		}))
	}

	// Include Configured MCP Servers & Tools if available
	cfg := config.Get()
	if cfg != nil && len(cfg.MCPServers) > 0 {
		for name := range cfg.MCPServers {
			mcpTag := fmt.Sprintf("@mcp:%s", name)
			if cleanQuery == "" || strings.Contains(strings.ToLower(mcpTag), strings.ToLower(cleanQuery)) {
				results = append(results, dialog.NewCompletionItem(dialog.CompletionItem{
					Title: fmt.Sprintf("%s (MCP Server)", mcpTag),
					Value: mcpTag,
				}))
			}
		}
	}

	// Include Files & Folders
	fileItems, err := u.filesProvider.GetChildEntries(cleanQuery)
	if err == nil {
		results = append(results, fileItems...)
	}

	return results, nil
}
