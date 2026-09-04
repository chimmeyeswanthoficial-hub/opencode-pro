# ⚡ OpenCode vs. OpenCode Pro: The Definitive Comparison Guide

> **A clear, highlighted breakdown showing why OpenCode Pro is 10x more powerful, autonomous, and beginner-friendly than original OpenCode.**

---

## 🏆 Quick Summary: At a Glance

| Feature Category | Original OpenCode ⚪ | OpenCode Pro ⚡ (Ours) |
| :--- | :--- | :--- |
| **Prompt Engineering** | ❌ Manual only (Fails with vague user prompts) | ✅ **Automated Thinking & Golden Prompt Synthesizer** |
| **Context Mentions** | ⚠️ `@` Files only | ✅ **Universal `@` (Files, Tools, MCPs) + `/` (Skills) + `#` (Rules, Git, Specs)** |
| **TUI Workflow** | ⚪ Standard Text Box | ✅ **Interactive Optimizer Bar + Live Reasoning Inspector (`Ctrl+P`, `Tab`)** |
| **Error Recovery** | ❌ None (Crashes or stops on test failure) | ✅ **Autonomous Self-Healing Diagnostics & Auto-Repair Loop** |
| **API & Headless Mode** | ❌ Terminal TUI only | ✅ **REST & JSON-RPC Daemon (`opencode daemon --port 8080`)** |
| **Command Centre** | ❌ None | ✅ **Complete Obsidian Command Centre Vault with HITL Checkbox Gate** |
| **Multi-Agent Orchestration** | ❌ Single session only | ✅ **LangGraph StateGraph Bridge (PM ➔ Architect ➔ Worker ➔ QA)** |
| **Startup Packs** | ❌ None | ✅ **Pre-built Skills (`/startup:api-scaffold`, `/startup:soc2-audit`, etc.)** |
| **Telemetry & Cost Tracking** | ❌ Local text logs only | ✅ **Native Langfuse Observability & Token Cost Monitoring** |

---

## 🔄 The Workflow Transformation: Before vs. After

### ⚪ In Original OpenCode (The Old Way)
```text
[Non-Engineer User] ➔ Types "fix the auth bug" ➔ Agent gets confused / hallucinates ➔ Fails silently
```

### ⚡ In OpenCode Pro (The New Way)
```text
[Non-Engineer User] ➔ Types "fix auth bug @internal/auth/jwt.go /test #rules:security"
                                   │
                                   ▼
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                        ⚡ AUTOMATED PROMPT OPTIMIZER (Ours)                            │
│  - Introspects Go/SQLC stack + Git diff + Target Files + Security Rules                │
│  - Synthesizes Structured Golden Prompt with Verification Steps                        │
└──────────────────────────────────┬─────────────────────────────────────────────────────┘
                                   │
                                   ▼
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                       ⚡ AUTONOMOUS AGENT + SELF-HEALING                               │
│  - Executes edits with edit/patch/grep tools                                           │
│  - Runs `go test ./...` ➔ If failed, self-heals compiler diagnostics automatically     │
│  - Syncs report and git diff to Obsidian Command Centre & Langfuse Telemetry           │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 🔍 Deep-Dive: 7 Major Superpowers in OpenCode Pro

### 1. 🧠 Smart Prompt Optimizer for Non-Prompt Engineers
- **Original OpenCode:** Requires the user to write detailed, technical prompts with explicit context constraints. Non-engineers struggle to get useful results.
- **OpenCode Pro:** Includes an **Automated Thinking & Reasoning Engine** (`internal/optimizer/`). When you type a rough prompt, it inspects your repository, infers the tech stack, resolves target files, detects missing test requirements, and expands it into an expert **Golden Prompt**.

### 2. 🎯 Universal Context Mentions (`@`, `/`, `#`)
- **Original OpenCode:** Only supported `@` for file autocompletion.
- **OpenCode Pro:** Multiplexes 3 distinct context layers:
  - **`@` Mention:** Project files, folders, LSP symbols, built-in tools (`@tool:edit`, `@tool:grep`), and MCP servers (`@mcp:github`, `@mcp:postgres`).
  - **`/` Skills:** 10+ standard and startup skill workflows (`/test`, `/refactor`, `/security-audit`, `/startup:api-scaffold`, `/startup:pitch-to-spec`).
  - **`#` Rules & Specs:** Project coding conventions (`#rules`, `#rules:strict-types`, `#rules:clean-architecture`), live git state (`#git:diff`, `#git:branch`), and architecture specs (`#spec:auth`, `#spec:db`).

### 3. 🖥️ Below-Input Optimizer Bar & Live Prompt Inspector Modal
- **Original OpenCode:** Simple cursor input.
- **OpenCode Pro:**
  - Dynamic **Optimizer Bar** right beneath the comment input showing attached context count (`📁 2 files • 🛠️ 1 skill • 📜 1 rule`).
  - Press **`Tab`** to instantly swap rough text with the AI-optimized Golden Prompt.
  - Press **`Ctrl+P`** to launch the **Prompt Inspector Modal** and inspect the step-by-step reasoning chain and tool bindings before running.
  - Press **`Ctrl+O`** to toggle auto-optimization on/off.

### 4. 🔄 Autonomous Self-Healing & Diagnostic Auto-Repair Loop
- **Original OpenCode:** If code compilation fails or tests break, the user has to manually copy-paste the error and re-prompt the agent.
- **OpenCode Pro:** Built-in **Self-Healing Engine** (`internal/optimizer/self_healing.go`). Parses Go, TypeScript, Python, and Rust compiler stack traces, generates targeted auto-repair prompts, and iterates until the build is green.

### 5. 🌐 Headless Daemon Mode & JSON-RPC REST API
- **Original OpenCode:** Could only run inside an interactive terminal TUI.
- **OpenCode Pro:** Start with `opencode daemon --port 8080` to turn OpenCode into a high-speed background execution service with REST and JSON-RPC endpoints (`/api/v1/optimize`, `/api/v1/run`, `/api/v1/context`, `/api/v1/health`).

### 6. 📊 Obsidian Command Centre & HITL Checkbox Gate
- **Original OpenCode:** No workspace or GUI dashboard integration.
- **OpenCode Pro:** Comes with a turnkey **Obsidian Mission Control Vault** (`obsidian-command-centre/`):
  - **Drop-in Autopilot:** Drop `.md` tasks in `01-Tasks/Inbox/` and OpenCode automatically executes them in the background.
  - **Human-in-the-Loop (HITL) Gatekeeper:** Destructive commands trigger a review note with a `- [ ] Approve Execution` checkbox in Obsidian. Checking `- [x]` in Obsidian automatically resumes execution!

### 7. 🐍 LangGraph Multi-Agent Orchestrator Kit
- **Original OpenCode:** Single-threaded execution.
- **OpenCode Pro:** Includes a ready-to-run **LangGraph StateGraph** (`langgraph-bridge/`):
  - Connects a **PM Agent ➔ Architect ➔ OpenCode Worker ➔ QA Validator ➔ Obsidian Sync** in an autonomous cyclic workflow.

---

## 📋 Concrete Example: Prompt Transformation

### What the Non-Engineer Types:
```text
> add role verification to jwt middleware @internal/auth/jwt.go /test #rules:security
```

### ⚡ What OpenCode Pro Automatically Synthesizes for the Agent:
```markdown
### Role & Objective
Act as a Principal Go Engineer. Execute /test workflow on request: 'add role verification to jwt middleware'

### Scope & Target Files
- `internal/auth/jwt.go`
- Tech Stack: Go (SQLC / Standard Library)

### Skill Guidelines & Methodology
- **/test (Testing): Analyze the selected files, detect the test framework (e.g. go test, jest, pytest), and write complete, passing unit & edge-case tests with high coverage.**

### Recommended Tools
- `edit`
- `grep`
- `view`
- `bash`

### Constraints & Quality Standards
- Write clean, production-grade code adhering to repository patterns.
- Avoid breaking existing functionality or removing exported APIs.
- Handle error scenarios gracefully with explicit context.
- Applied Rule: `#rules:security - Zero-trust input sanitization, least privilege, safe crypto`

### Definition of Done & Verification
1. Inspect and make necessary code edits accurately.
2. Run tests to verify correctness: `go test -v ./...`
3. Provide a clear summary of all modifications.
```

---

## 🎯 Which One Should You Choose?

- Choose **Original OpenCode** if you only want a basic terminal chat window and are willing to manually write comprehensive prompts every time.
- Choose **OpenCode Pro** if you want an **autonomous, self-healing software development powerhouse** that lets anyone (founders, PMs, engineers) build software 10x faster with complete Obsidian and LangGraph automation.
