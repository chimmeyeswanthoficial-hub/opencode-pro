package optimizer

import (
	"fmt"
	"strings"
)

// FormatGoldenPrompt builds a structured, high-precision engineering prompt
func FormatGoldenPrompt(
	rawGoal string,
	detectedIntent string,
	projectInfo ProjectInfo,
	targetFiles []string,
	tools []string,
	skills []string,
	rules []string,
	ruleContents []string,
) string {
	var sb strings.Builder

	// Header & Role
	sb.WriteString("### Role & Objective\n")
	if detectedIntent != "" {
		sb.WriteString(fmt.Sprintf("Act as a Principal %s Engineer. %s\n\n", projectInfo.PrimaryLanguage, detectedIntent))
	} else {
		sb.WriteString(fmt.Sprintf("Act as a Principal %s Engineer to execute the following task with high precision:\n%s\n\n", projectInfo.PrimaryLanguage, rawGoal))
	}

	// Target Scope & Codebase Context
	sb.WriteString("### Scope & Target Files\n")
	if len(targetFiles) > 0 {
		for _, f := range targetFiles {
			sb.WriteString(fmt.Sprintf("- `%s`\n", f))
		}
	} else {
		sb.WriteString("- Inspect workspace and relevant source files before making changes.\n")
	}
	if projectInfo.Framework != "" {
		sb.WriteString(fmt.Sprintf("- Tech Stack: %s (%s)\n", projectInfo.PrimaryLanguage, projectInfo.Framework))
	}
	sb.WriteString("\n")

	// Applied Skills & Methodology
	if len(skills) > 0 {
		sb.WriteString("### Skill Guidelines & Methodology\n")
		for _, s := range skills {
			sb.WriteString(fmt.Sprintf("- **%s**\n", s))
		}
		sb.WriteString("\n")
	}

	// Recommended Tools
	if len(tools) > 0 {
		sb.WriteString("### Recommended Tools\n")
		for _, t := range tools {
			sb.WriteString(fmt.Sprintf("- `%s`\n", t))
		}
		sb.WriteString("\n")
	}

	// Applied Rules & Constraints
	sb.WriteString("### Constraints & Quality Standards\n")
	sb.WriteString("- Write clean, production-grade code adhering to repository patterns.\n")
	sb.WriteString("- Avoid breaking existing functionality or removing exported APIs.\n")
	sb.WriteString("- Handle error scenarios gracefully with explicit context.\n")
	if len(rules) > 0 {
		for _, r := range rules {
			sb.WriteString(fmt.Sprintf("- Applied Rule: `%s`\n", r))
		}
	}
	if len(ruleContents) > 0 {
		sb.WriteString("\n**Rule Details:**\n")
		for _, rc := range ruleContents {
			sb.WriteString(fmt.Sprintf("%s\n", rc))
		}
	}
	sb.WriteString("\n")

	// Definition of Done & Verification
	sb.WriteString("### Definition of Done & Verification\n")
	sb.WriteString("1. Inspect and make necessary code edits accurately.\n")
	if projectInfo.TestCommand != "" {
		sb.WriteString(fmt.Sprintf("2. Run tests to verify correctness: `%s`\n", projectInfo.TestCommand))
	} else {
		sb.WriteString("2. Verify changes with diagnostics/tests without regressions.\n")
	}
	sb.WriteString("3. Provide a clear summary of all modifications.\n")

	return sb.String()
}
