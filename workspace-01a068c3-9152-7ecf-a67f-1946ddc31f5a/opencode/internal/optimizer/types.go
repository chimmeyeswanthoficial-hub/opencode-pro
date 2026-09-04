package optimizer

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SkillDefinition defines an automated skill or workflow template
type SkillDefinition struct {
	ID          string
	Name        string
	Category    string
	Description string
	Template    string
}

// DefaultBuiltinSkills returns standard startup and engineering skills
func DefaultBuiltinSkills() []SkillDefinition {
	return []SkillDefinition{
		{
			ID:          "test",
			Name:        "/test",
			Category:    "Testing",
			Description: "Generate comprehensive unit, integration & table-driven tests",
			Template:    "Analyze the selected files, detect the test framework (e.g. go test, jest, pytest), and write complete, passing unit & edge-case tests with high coverage.",
		},
		{
			ID:          "refactor",
			Name:        "/refactor",
			Category:    "Engineering",
			Description: "Refactor for clean code, DRY, performance & readability",
			Template:    "Refactor the code to improve maintainability, reduce complexity, remove duplication, and adhere to idiomatic patterns without altering public behavior.",
		},
		{
			ID:          "security-audit",
			Name:        "/security-audit",
			Category:    "Security",
			Description: "Scan for vulnerabilities, OWASP Top 10, SQLi & secrets",
			Template:    "Perform a deep security audit checking for SQL injection, XSS, unvalidated inputs, insecure deserialization, plain-text secrets, and missing auth checks.",
		},
		{
			ID:          "doc",
			Name:        "/doc",
			Category:    "Documentation",
			Description: "Generate comprehensive docstrings, API markdown & README",
			Template:    "Generate detailed, idiomatic documentation, exported function comments, API request/response schemas, and usage examples.",
		},
		{
			ID:          "explain",
			Name:        "/explain",
			Category:    "Analysis",
			Description: "Step-by-step architectural and execution flow walkthrough",
			Template:    "Break down the codebase logic, control flow, data models, error handling, and component interactions in clear, plain language.",
		},
		{
			ID:          "startup:api-scaffold",
			Name:        "/startup:api-scaffold",
			Category:    "Startup Pack",
			Description: "Scaffold production REST/gRPC API with models, DB & routes",
			Template:    "Scaffold a complete production-grade API endpoint with request validation, database query models, service layer, error handling, and Swagger/OpenAPI docs.",
		},
		{
			ID:          "startup:pitch-to-spec",
			Name:        "/startup:pitch-to-spec",
			Category:    "Startup Pack",
			Description: "Convert product pitch/brief into architecture & task list",
			Template:    "Analyze the product requirements, decompose into domain entities, define database schemas, list API endpoints, and output an actionable step-by-step task breakdown.",
		},
		{
			ID:          "startup:soc2-audit",
			Name:        "/startup:soc2-audit",
			Category:    "Startup Pack",
			Description: "Audit codebase against SOC2 compliance & data privacy",
			Template:    "Evaluate codebase against SOC2 Trust Services Criteria (Security, Availability, Confidentiality): verify audit logging, rate limiting, encryption in transit/rest, and access control.",
		},
		{
			ID:          "startup:e2e-tests",
			Name:        "/startup:e2e-tests",
			Category:    "Startup Pack",
			Description: "Generate end-to-end user journey tests",
			Template:    "Write complete end-to-end integration test scenarios covering happy paths, edge cases, authentication flows, and failure recoveries.",
		},
		{
			ID:          "benchmark",
			Name:        "/benchmark",
			Category:    "Performance",
			Description: "Create performance benchmarks & memory allocation profiling",
			Template:    "Construct benchmark suites, profile hotspots, reduce allocations, and optimize critical execution loops.",
		},
	}
}

// RuleDefinition defines a project rule, git context, or spec tag
type RuleDefinition struct {
	Tag         string
	Category    string
	Description string
	Content     string
}

// DefaultBuiltinRules returns standard context rule definitions
func DefaultBuiltinRules() []RuleDefinition {
	return []RuleDefinition{
		{
			Tag:         "#rules",
			Category:    "Project Rules",
			Description: "Include repository global coding rules and constraints (.opencode/rules.md)",
			Content:     "Enforce repository coding standards, linting requirements, and architectural boundaries.",
		},
		{
			Tag:         "#rules:strict-types",
			Category:    "Type Safety",
			Description: "Strict typing: no any, explicit error handling, robust schemas",
			Content:     "Ensure all types, interfaces, and function signatures are strictly typed without placeholders or loose type casting.",
		},
		{
			Tag:         "#rules:clean-architecture",
			Category:    "Architecture",
			Description: "Domain-Driven Design, decoupled layers, interface segregation",
			Content:     "Maintain clean separation between domain logic, data persistence, and transport/API controllers.",
		},
		{
			Tag:         "#rules:security",
			Category:    "Security Policy",
			Description: "Zero-trust input sanitization, least privilege, safe crypto",
			Content:     "Follow defensive coding guidelines: sanitize all external inputs, enforce authorization checks, and avoid leaking internal errors.",
		},
		{
			Tag:         "#git:diff",
			Category:    "Git Context",
			Description: "Inject unstaged working directory git diff into context",
			Content:     "Include current unstaged changes from git diff.",
		},
		{
			Tag:         "#git:staged",
			Category:    "Git Context",
			Description: "Inject staged git changes into context",
			Content:     "Include staged git changes from git diff --cached.",
		},
		{
			Tag:         "#git:branch",
			Category:    "Git Context",
			Description: "Current branch name, upstream info, and recent commits",
			Content:     "Include current branch name and git status.",
		},
		{
			Tag:         "#spec:auth",
			Category:    "Domain Spec",
			Description: "Authentication, session management & JWT token specifications",
			Content:     "Reference authentication design specifications and security token lifecycles.",
		},
		{
			Tag:         "#spec:db",
			Category:    "Domain Spec",
			Description: "Database migrations, relations, transactions & index rules",
			Content:     "Reference database schema rules, indexing strategies, and migration conventions.",
		},
		{
			Tag:         "#spec:api",
			Category:    "Domain Spec",
			Description: "REST/gRPC API status codes, response wrapping & error schema",
			Content:     "Reference unified API JSON response wrapping, HTTP status conventions, and error codes.",
		},
		{
			Tag:         "#vault:architecture",
			Category:    "Obsidian Vault",
			Description: "System architecture and component diagrams from Obsidian",
			Content:     "Reference Obsidian architecture knowledge base notes.",
		},
		{
			Tag:         "#vault:decisions",
			Category:    "Obsidian Vault",
			Description: "Architecture Decision Records (ADRs) from Obsidian",
			Content:     "Reference existing Architectural Decision Records in Obsidian.",
		},
	}
}

// ResolveRuleContent resolves a tag into its dynamic string content
func ResolveRuleContent(tag string) string {
	switch tag {
	case "#git:diff":
		cmd := exec.Command("git", "diff")
		if out, err := cmd.Output(); err == nil && len(out) > 0 {
			return string(out)
		}
		return "Git diff: no unstaged changes."
	case "#git:staged":
		cmd := exec.Command("git", "diff", "--cached")
		if out, err := cmd.Output(); err == nil && len(out) > 0 {
			return string(out)
		}
		return "Git diff staged: no staged changes."
	case "#git:branch":
		cmd := exec.Command("git", "status", "-s", "-b")
		if out, err := cmd.Output(); err == nil {
			return string(out)
		}
		return "Git status: not a git repository."
	case "#rules":
		for _, path := range []string{".opencode/rules.md", "CLAUDE.md", "AGENTS.md", ".rules.md"} {
			if content, err := os.ReadFile(path); err == nil {
				return string(content)
			}
		}
		return "Standard project rules: follow clean code, maintain test coverage, handle errors explicitly."
	default:
		if strings.HasPrefix(tag, "#rules:") {
			ruleName := strings.TrimPrefix(tag, "#rules:")
			path := filepath.Join(".", ".opencode", "rules", ruleName+".md")
			if content, err := os.ReadFile(path); err == nil {
				return string(content)
			}
		}
		return ""
	}
}

// ResolvedContext represents all parsed tokens
type ResolvedContext struct {
	Files       []string
	Tools       []string
	McpServers  []string
	Skills      []SkillDefinition
	Rules       []RuleDefinition
	RuleContent []string
}

// ParsePromptMentions parses @, /, and # tokens in a prompt
func ParsePromptMentions(text string) ResolvedContext {
	rc := ResolvedContext{}
	tokens := strings.Fields(text)

	allSkills := DefaultBuiltinSkills()
	allRules := DefaultBuiltinRules()

	for _, token := range tokens {
		switch {
		case strings.HasPrefix(token, "@tool:"):
			toolName := strings.TrimPrefix(token, "@tool:")
			rc.Tools = append(rc.Tools, toolName)
		case strings.HasPrefix(token, "@mcp:"):
			mcpName := strings.TrimPrefix(token, "@mcp:")
			rc.McpServers = append(rc.McpServers, mcpName)
		case strings.HasPrefix(token, "@"):
			filePath := strings.TrimPrefix(token, "@")
			if filePath != "" {
				rc.Files = append(rc.Files, filePath)
			}
		case strings.HasPrefix(token, "/"):
			skillName := strings.TrimPrefix(token, "/")
			for _, s := range allSkills {
				if strings.EqualFold(s.ID, skillName) || strings.EqualFold(s.Name, token) {
					rc.Skills = append(rc.Skills, s)
				}
			}
		case strings.HasPrefix(token, "#"):
			for _, r := range allRules {
				if strings.EqualFold(r.Tag, token) {
					rc.Rules = append(rc.Rules, r)
					content := ResolveRuleContent(token)
					if content != "" {
						rc.RuleContent = append(rc.RuleContent, fmt.Sprintf("=== Rule: %s ===\n%s", token, content))
					}
				}
			}
		}
	}

	return rc
}
