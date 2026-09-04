# 🐍 LangGraph Multi-Agent Orchestrator Bridge

This package connects **LangGraph** (StateGraph multi-agent loops) with **OpenCode Pro Headless Daemon** and **Obsidian Command Centre**.

## Architecture & Workflow

```
[Obsidian Task Note] ➔ [PM Agent] ➔ [Architect Agent] ➔ [OpenCode Worker] ➔ [QA Validator] ➔ [Obsidian Sync]
                                                             │ (on test failure)
                                                             └──⮌ (Self-Healing Retry Loop)
```

## Quick Start

### 1. Install Dependencies
```bash
pip install -r requirements.txt
```

### 2. Start OpenCode Daemon
```bash
opencode daemon --port 8080 --watch-obsidian ../obsidian-command-centre
```

### 3. Run the Multi-Agent Pipeline
```bash
python graph.py
```

### 4. Or Run via Docker Compose (with Langfuse Observability)
```bash
docker-compose up -d
```
