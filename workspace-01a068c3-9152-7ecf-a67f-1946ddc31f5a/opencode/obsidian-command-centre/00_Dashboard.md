# 📊 Autonomous Startup Command Centre & Mission Control

> **System Status:** 🟢 All Systems Operational | **Connected Agent:** OpenCode Pro | **Telemetry:** Langfuse Active

---

## 🚀 Quick Actions & Dispatch
- Drop new task markdown files into `01-Tasks/Inbox/` to auto-execute.
- To execute via terminal: `opencode daemon --watch-obsidian .`
- Use `@` for files & MCP tools, `/` for skill templates, `#` for project rules.

---

## 📥 Active Task Queue

```dataview
TABLE file.mtime as "Modified", status as "Status", priority as "Priority"
FROM "01-Tasks/Inbox"
SORT file.mtime desc
```

---

## 🛑 Waiting for Human Approval (HITL)

```dataview
TABLE file.mtime as "Requested", risk as "Risk Level", tools_requested as "Tools"
FROM "01-Tasks/Waiting-Approval"
```

---

## ✅ Recently Completed Autonomous Tasks

```dataview
TABLE file.mtime as "Completed", status as "Status"
FROM "01-Tasks/Completed"
SORT file.mtime desc
LIMIT 10
```

---

## 🏛️ Architecture Decision Records (ADRs)

```dataview
TABLE status as "Status", deciders as "Deciders", date as "Date"
FROM "03-Decisions-ADR"
SORT date desc
```

---

## 📈 Observability & System Resources
- **Observability Host:** `http://localhost:3000` (Langfuse)
- **Local RAG Provider:** Hybrid Dense + BM25 MCP
- **Daemon API:** `http://localhost:8080/api/v1`
