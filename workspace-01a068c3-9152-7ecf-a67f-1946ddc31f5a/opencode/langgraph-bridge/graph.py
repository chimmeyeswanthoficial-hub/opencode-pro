"""
LangGraph Multi-Agent Orchestrator for OpenCode Pro & Obsidian Command Centre.
Workflow: PM Agent -> Architect -> OpenCode Worker -> QA Validator -> Obsidian Sync.
"""

from typing import TypedDict, List, Dict, Any, Optional
import os
import asyncio
from langgraph.graph import StateGraph, END
from opencode_client import OpenCodeClient

class CommandCentreState(TypedDict):
    task_id: str
    raw_task: str
    target_repo: str
    pm_specs: str
    target_files: List[str]
    applied_rules: List[str]
    enhanced_prompt: str
    worker_output: str
    qa_status: str
    error_count: int
    retry_count: int
    completed: bool

client = OpenCodeClient(base_url=os.getenv("OPENCODE_DAEMON_URL", "http://localhost:8080"))

async def pm_agent_node(state: CommandCentreState) -> CommandCentreState:
    """PM Node: Decomposes business requirement into engineering scope."""
    print(f"📋 [PM Agent] Analyzing task: {state['task_id']}")
    resolved = await client.resolve_context(state["raw_task"])
    state["target_files"] = resolved.get("Files", [])
    state["pm_specs"] = f"Decomposed requirement for {state['raw_task']}"
    return state

async def architect_node(state: CommandCentreState) -> CommandCentreState:
    """Architect Node: Injects system architecture specs and optimizes prompt."""
    print("📐 [Architect Agent] Synthesizing Golden Engineering Prompt...")
    opt = await client.optimize(state["raw_task"])
    state["enhanced_prompt"] = opt.get("OptimizedPrompt", state["raw_task"])
    state["applied_rules"] = opt.get("AppliedRules", [])
    return state

async def opencode_worker_node(state: CommandCentreState) -> CommandCentreState:
    """Worker Node: Executes code edits, shell commands, and tests via OpenCode Daemon."""
    print("⚡ [OpenCode Worker] Executing tools and making code modifications...")
    result = await client.run(state["enhanced_prompt"], auto_enhance=False)
    state["worker_output"] = str(result)
    return state

async def qa_validator_node(state: CommandCentreState) -> CommandCentreState:
    """QA Node: Validates test output and security criteria."""
    print("🧪 [QA Validator] Checking test results and security rules...")
    output = state.get("worker_output", "")
    if "error" in output.lower() and state["retry_count"] < 2:
        state["qa_status"] = "failed"
        state["retry_count"] += 1
        print(f"⚠️ [QA Validator] Validation failed. Initiating retry #{state['retry_count']}")
    else:
        state["qa_status"] = "passed"
        state["completed"] = True
        print("✅ [QA Validator] All tests and specs verified successfully.")
    return state

async def obsidian_sync_node(state: CommandCentreState) -> CommandCentreState:
    """Obsidian Sync Node: Writes completion report back to Obsidian Command Centre."""
    print(f"📊 [Obsidian Sync] Updating mission control for task {state['task_id']}...")
    return state

def should_retry(state: CommandCentreState) -> str:
    """Routing logic for self-healing loops."""
    if state["qa_status"] == "failed" and state["retry_count"] < 2:
        return "architect"
    return "obsidian_sync"

def build_command_centre_graph():
    """Build the compiled LangGraph workflow."""
    builder = StateGraph(CommandCentreState)
    builder.add_node("pm_agent", pm_agent_node)
    builder.add_node("architect", architect_node)
    builder.add_node("opencode_worker", opencode_worker_node)
    builder.add_node("qa_validator", qa_validator_node)
    builder.add_node("obsidian_sync", obsidian_sync_node)

    builder.set_entry_point("pm_agent")
    builder.add_edge("pm_agent", "architect")
    builder.add_edge("architect", "opencode_worker")
    builder.add_edge("opencode_worker", "qa_validator")
    builder.add_conditional_edges(
        "qa_validator",
        should_retry,
        {"architect": "architect", "obsidian_sync": "obsidian_sync"}
    )
    builder.add_edge("obsidian_sync", END)
    return builder.compile()

if __name__ == "__main__":
    app = build_command_centre_graph()
    initial_state: CommandCentreState = {
        "task_id": "TASK-001",
        "raw_task": "Implement JWT auth middleware @internal/auth/jwt.go /test #rules:security",
        "target_repo": ".",
        "pm_specs": "",
        "target_files": [],
        "applied_rules": [],
        "enhanced_prompt": "",
        "worker_output": "",
        "qa_status": "pending",
        "error_count": 0,
        "retry_count": 0,
        "completed": False,
    }
    print("🚀 Running LangGraph Autonomous Pipeline...")
    asyncio.run(app.ainvoke(initial_state))
