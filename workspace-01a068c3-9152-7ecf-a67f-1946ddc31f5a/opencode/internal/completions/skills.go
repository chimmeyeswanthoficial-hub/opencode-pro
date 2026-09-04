package completions

import (
	"fmt"
	"strings"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/opencode-ai/opencode/internal/optimizer"
	"github.com/opencode-ai/opencode/internal/tui/components/dialog"
)

type skillsContextGroup struct {
	prefix string
	skills []optimizer.SkillDefinition
}

func NewSkillsContextGroup() dialog.CompletionProvider {
	return &skillsContextGroup{
		prefix: "/",
		skills: optimizer.DefaultBuiltinSkills(),
	}
}

func (cg *skillsContextGroup) GetId() string {
	return cg.prefix
}

func (cg *skillsContextGroup) GetEntry() dialog.CompletionItemI {
	return dialog.NewCompletionItem(dialog.CompletionItem{
		Title: "Skills & Commands",
		Value: "skills",
	})
}

func (cg *skillsContextGroup) GetChildEntries(query string) ([]dialog.CompletionItemI, error) {
	customCmds, _ := dialog.LoadCustomCommands()

	var targets []string
	lookup := make(map[string]string)

	for _, s := range cg.skills {
		val := s.Name
		targets = append(targets, val)
		lookup[val] = fmt.Sprintf("%s (%s) - %s", s.Name, s.Category, s.Description)
	}

	for _, cmd := range customCmds {
		val := "/" + cmd.ID
		targets = append(targets, val)
		lookup[val] = fmt.Sprintf("%s - %s", val, cmd.Description)
	}

	var matches []string
	cleanQuery := strings.TrimPrefix(query, "/")
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
