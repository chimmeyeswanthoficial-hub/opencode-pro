package completions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/opencode-ai/opencode/internal/optimizer"
	"github.com/opencode-ai/opencode/internal/tui/components/dialog"
)

type rulesContextGroup struct {
	prefix string
	rules  []optimizer.RuleDefinition
}

func NewRulesContextGroup() dialog.CompletionProvider {
	return &rulesContextGroup{
		prefix: "#",
		rules:  optimizer.DefaultBuiltinRules(),
	}
}

func (cg *rulesContextGroup) GetId() string {
	return cg.prefix
}

func (cg *rulesContextGroup) GetEntry() dialog.CompletionItemI {
	return dialog.NewCompletionItem(dialog.CompletionItem{
		Title: "Project Rules & Specs",
		Value: "rules",
	})
}

func (cg *rulesContextGroup) GetChildEntries(query string) ([]dialog.CompletionItemI, error) {
	var targets []string
	lookup := make(map[string]string)

	for _, r := range cg.rules {
		targets = append(targets, r.Tag)
		lookup[r.Tag] = fmt.Sprintf("%s (%s) - %s", r.Tag, r.Category, r.Description)
	}

	// Dynamic: scan for any custom markdown files in .opencode/rules/
	rulesDir := filepath.Join(".", ".opencode", "rules")
	if entries, err := os.ReadDir(rulesDir); err == nil {
		for _, e := range entries {
			if strings.HasSuffix(e.Name(), ".md") {
				tag := "#rules:" + strings.TrimSuffix(e.Name(), ".md")
				targets = append(targets, tag)
				lookup[tag] = fmt.Sprintf("%s (Custom Rule) - %s", tag, e.Name())
			}
		}
	}

	var matches []string
	cleanQuery := strings.TrimPrefix(query, "#")
	if cleanQuery == "" {
		matches = targets
	} else {
		matches = fuzzy.Find(cleanQuery, targets)
	}

	items := make([]dialog.CompletionItemI, 0, len(matches))
	for _, m := range matches {
		items = append(items, dialog.NewCompletionItem(dialog.CompletionItem{
			Title: lookup[m],
			Value: m,
		}))
	}

	return items, nil
}
