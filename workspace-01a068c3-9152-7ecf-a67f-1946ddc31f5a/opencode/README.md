# ⌬ OpenCode Pro — Autonomous AI Coding Agent & Command Centre

<p align="center">
  <img src="https://img.shields.io/badge/Language-Go%201.24+-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white" alt="Go" />
  <img src="https://img.shields.io/badge/TUI-Bubbletea-FD971F.svg?style=for-the-badge" alt="Bubbletea" />
  <img src="https://img.shields.io/badge/Protocol-MCP%20Enabled-7952B3.svg?style=for-the-badge" alt="MCP" />
  <img src="https://img.shields.io/badge/Automation-Obsidian%20%2B%20LangGraph-7F52FF.svg?style=for-the-badge" alt="Obsidian" />
  <img src="https://img.shields.io/badge/Telemetry-Langfuse-black.svg?style=for-the-badge" alt="Langfuse" />
  <img src="https://img.shields.io/badge/License-MIT-blue.svg?style=for-the-badge" alt="License" />
</p>

> **The next-generation terminal AI agent equipped with automated prompt engineering, universal context mentions (`@`, `/`, `#`), headless daemon orchestration, and a complete Obsidian Command Centre.**

---

## ⚡ What Makes OpenCode Pro Different?

Traditional AI coding assistants fail when non-prompt engineers give vague instructions like *"fix login bug"*. OpenCode Pro solves this by introducing a **Smart Prompt Optimization & Reasoning Engine** directly inside the terminal:

1. **Smart Prompt Optimizer for Non-Engineers**: Automatically converts rough, vague ideas into structured, production-grade **Golden Prompts** with clear objectives, file anchors, tool recommendations, and verification criteria.
2. **Universal Context Mentions (`@`, `/`, `#`)**:
   - `@` : Mention Files, Folders, LSP Symbols, Built-in Tools (`edit`, `grep`, `bash`), and MCP Servers.
   - `/` : Invoke Skills & Custom Workflows (`/test`, `/refactor`, `/security-audit`, `/startup:api-scaffold`).
   - `#` : Attach Project Rules (`#rules`, `#rules:security`), Git Context (`#git:diff`, `#git:branch`), or Architecture Specs (`#spec:auth`, `#spec:db`).
3. **Below-Input Optimizer Bar & Live Inspector**:
   - Real-time status pill directly under your comment box.
   - Press **`Tab`** to instantly expand your rough text into an optimized prompt.
   - Press **`Ctrl+P`** to open the **Prompt Inspector Modal** and inspect the step-by-step reasoning chain.
   - Press **`Ctrl+O`** to toggle the auto-optimizer on or off.
4. **Autonomous Obsidian Command Centre**:
   - Drop markdown tasks into an Obsidian vault `Inbox/` folder and watch OpenCode execute them in the background, updating logs, diffs, and test reports.
5. **Headless Daemon API (`opencode daemon`)**:
   - Exposes REST & JSON-RPC endpoints for **LangGraph multi-agent loops**, CI/CD pipelines, and web dashboards.
6. **Built-in Startup Automation Packs**:
   - Pre-configured skills for bootstrapping APIs, converting pitch decks to specs, SOC2 compliance audits, and E2E testing.

---

## 🖥️ Terminal UI & Workflow

```text
┌───────────────────────────────────────────────────────────────────────────────────────────────────┐
│ > add user role verification to jwt middleware @internal/auth/jwt.go /test #rules                 │
├───────────────────────────────────────────────────────────────────────────────────────────────────┤
│ ⚡ OPTIMIZED • 📁 1 file(s) • 🛠️ 1 skill(s) • 📜 1 rule(s)   [Tab] Expand • [Ctrl+P] Inspect • [Enter] Run│
└───────────────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Quick Start

### 1. Build & Install from Source

```bash
# Clone the repository
git clone https://github.com/opencode-ai/opencode.git
cd opencode

# Build binary
go build -o opencode .

# Move to your PATH
sudo mv opencode /usr/local/bin/
```

### 2. Launch Interactive TUI

```bash
# Launch in the current repository
opencode

# Launch with debug logging
opencode -d
```

### 3. CLI One-Shot Prompt Optimization

```bash
# Optimize any prompt and view reasoning chain
opencode optimize "add stripe webhook verification @server.go /test #rules:security"

# Output as JSON for scripts or LangGraph pipelines
opencode optimize "refactor db queries @repo.go #rules:clean-architecture" --json
```

### 4. Run Headless Daemon Mode (for Obsidian & LangGraph)

```bash
# Start background REST API on port 8080
opencode daemon --port 8080

# Start daemon with automated Obsidian task watcher
opencode daemon --port 8080 --watch-obsidian ./obsidian-command-centre
```

---

## 🎯 Universal Context Triggers (`@`, `/`, `#`)

| Trigger | Description | Examples |
| :--- | :--- | :--- |
| **`@` Mention** | Target files, directories, LSP symbols, built-in tools, and active MCP servers | `@internal/auth/jwt.go`<br>`@tool:grep`<br>`@tool:edit`<br>`@mcp:github` |
| **`/` Skills** | Execute built-in engineering workflows and custom prompt macros | `/test`<br>`/refactor`<br>`/security-audit`<br>`/startup:api-scaffold`<br>`/startup:pitch-to-spec` |
| **`#` Rules & Specs** | Enforce project conventions, inject live git state, or reference domain specs | `#rules`<br>`#rules:security`<br>`#rules:strict-types`<br>`#git:diff`<br>`#git:staged`<br>`#spec:auth` |

---

## 🧠 Startup Automation Packs

OpenCode Pro includes specialized skills designed for startups and fast-moving teams:

- **`/startup:pitch-to-spec`**: Converts high-level product briefs or pitch notes into database schemas, API routes, and task breakdowns.
- **`/startup:api-scaffold`**: Generates production-grade REST/gRPC endpoints with input validation, database models, and OpenAPI specs.
- **`/startup:soc2-audit`**: Evaluates codebase against SOC2 security criteria (rate limiting, audit logs, secret hygiene, encryption).
- **`/startup:e2e-tests`**: Generates complete end-to-end user journey test suites with framework auto-detection.

---

## 🏛️ Obsidian Command Centre Architecture

The included `obsidian-command-centre/` vault provides a visual human-in-the-loop dashboard:

```text
📁 obsidian-command-centre/
├── 📊 00_Dashboard.md               <-- Real-time Task Queue & KPI Mission Control
├── 📥 01-Tasks/
│   ├── 📥 Inbox/                    <-- Drop task markdown files here to auto-run
│   ├── ⚡ In-Progress/              <-- Currently executing tasks
│   ├── 🛑 Waiting-Approval/         <-- Human-in-the-loop approval gates
│   └── ✅ Completed/                <-- Output reports, diffs, and test results
├── 🤖 02-Agents/                    <-- Agent prompts & LangGraph workflow graphs
├── 🏛️ 03-Decisions-ADR/             <-- Architecture Decision Records (ADRs)
└── 🧠 04-Knowledge-Base/            <-- Architecture & business guidelines
```

---

## ⚙️ Configuration (`~/.config/opencode/config.yaml` or `.opencode.json`)

```json
{
  "$schema": "https://raw.githubusercontent.com/opencode-ai/opencode/main/opencode-schema.json",
  "theme": "catppuccin-mocha",
  "default_agent": "coder",
  "mcp_servers": {
    "github": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-github"]
    },
    "postgres": {
      "type": "stdio",
      "command": "npx",
      "args": ["-y", "@modelcontextprotocol/server-postgres", "postgresql://localhost/mydb"]
    }
  }
}
```

### Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `LANGFUSE_HOST` | Observability host URL for telemetry | `http://localhost:3000` |
| `LANGFUSE_PUBLIC_KEY` | Public key for Langfuse | Optional |
| `LANGFUSE_SECRET_KEY` | Secret key for Langfuse | Optional |
| `OPENAI_API_KEY` | OpenAI provider key | Optional |
| `ANTHROPIC_API_KEY` | Anthropic Claude provider key | Optional |
| `GEMINI_API_KEY` | Google Gemini provider key | Optional |

---

## 🧪 Testing & Quality Assurance

Run all unit and integration tests across the codebase:

```bash
go test -v ./...
```

Run code formatting and vet checks:

```bash
go vet ./...
go fmt ./...
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these guidelines:
1. Ensure all code adheres to `.opencode/rules.md`.
2. Write unit tests for new features.
3. Verify that `go test ./...` and `go build ./...` pass with zero warnings.

---

## 📜 License

Distributed under the **MIT License**. See `LICENSE` for details.
