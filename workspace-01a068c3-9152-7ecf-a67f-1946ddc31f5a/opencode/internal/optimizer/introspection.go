package optimizer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProjectInfo contains introspected repository details
type ProjectInfo struct {
	PrimaryLanguage string
	Framework       string
	PackageManager  string
	TestCommand     string
	GitBranch       string
	ModifiedFiles   []string
	HasRulesFile    bool
	RulesSummary    string
	TopLevelFiles   []string
}

// IntrospectProject scans the current working directory for tech stack and environment details
func IntrospectProject(dir string) ProjectInfo {
	if dir == "" {
		dir = "."
	}

	info := ProjectInfo{
		PrimaryLanguage: "Unknown",
		TestCommand:     "test",
	}

	// 1. Scan Top Level Files
	entries, err := os.ReadDir(dir)
	if err == nil {
		for _, e := range entries {
			if !e.IsDir() {
				info.TopLevelFiles = append(info.TopLevelFiles, e.Name())
			}
		}
	}

	// 2. Detect Tech Stack & Test Commands
	switch {
	case fileExists(filepath.Join(dir, "go.mod")):
		info.PrimaryLanguage = "Go"
		info.PackageManager = "go mod"
		info.TestCommand = "go test -v ./..."
		if fileExists(filepath.Join(dir, "sqlc.yaml")) {
			info.Framework = "SQLC / Standard Library"
		}
	case fileExists(filepath.Join(dir, "package.json")):
		info.PrimaryLanguage = "TypeScript/JavaScript"
		info.PackageManager = "npm"
		if fileExists(filepath.Join(dir, "pnpm-lock.yaml")) {
			info.PackageManager = "pnpm"
		} else if fileExists(filepath.Join(dir, "yarn.lock")) {
			info.PackageManager = "yarn"
		} else if fileExists(filepath.Join(dir, "bun.lockb")) {
			info.PackageManager = "bun"
		}
		info.TestCommand = info.PackageManager + " test"

		// Detect UI / Framework
		if fileExists(filepath.Join(dir, "next.config.js")) || fileExists(filepath.Join(dir, "next.config.ts")) {
			info.Framework = "Next.js"
		} else if fileExists(filepath.Join(dir, "vite.config.ts")) || fileExists(filepath.Join(dir, "vite.config.js")) {
			info.Framework = "Vite / React / Vue"
		}
	case fileExists(filepath.Join(dir, "pyproject.toml")) || fileExists(filepath.Join(dir, "requirements.txt")):
		info.PrimaryLanguage = "Python"
		info.PackageManager = "pip/poetry"
		info.TestCommand = "pytest -v"
		if fileExists(filepath.Join(dir, "manage.py")) {
			info.Framework = "Django"
			info.TestCommand = "python manage.py test"
		} else if fileExists(filepath.Join(dir, "app.py")) {
			info.Framework = "Flask/FastAPI"
		}
	case fileExists(filepath.Join(dir, "Cargo.toml")):
		info.PrimaryLanguage = "Rust"
		info.PackageManager = "cargo"
		info.TestCommand = "cargo test"
	}

	// 3. Inspect Git State
	if branchOut, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output(); err == nil {
		info.GitBranch = strings.TrimSpace(string(branchOut))
	}
	if diffOut, err := exec.Command("git", "status", "--porcelain").Output(); err == nil {
		lines := strings.Split(string(diffOut), "\n")
		for _, l := range lines {
			if len(l) > 3 {
				info.ModifiedFiles = append(info.ModifiedFiles, strings.TrimSpace(l[3:]))
			}
		}
	}

	// 4. Inspect Rules Files
	for _, ruleFile := range []string{".opencode/rules.md", "CLAUDE.md", "AGENTS.md", ".rules.md"} {
		p := filepath.Join(dir, ruleFile)
		if fileExists(p) {
			info.HasRulesFile = true
			if content, err := os.ReadFile(p); err == nil {
				summary := string(content)
				if len(summary) > 200 {
					summary = summary[:200] + "..."
				}
				info.RulesSummary = summary
			}
			break
		}
	}

	return info
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
