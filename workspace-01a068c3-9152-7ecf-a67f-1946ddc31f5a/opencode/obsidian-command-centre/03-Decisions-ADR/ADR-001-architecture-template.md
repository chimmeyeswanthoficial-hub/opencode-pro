---
title: ADR-001: Architecture Decision Record Template
status: accepted
deciders: OpenCode Agent, Core Engineering
date: 2026-09-03
---

# ADR-001: OpenCode Modular Extension Architecture

## Context and Problem Statement
How should OpenCode handle universal context mentions (`@`, `/`, `#`), prompt optimization for non-prompt engineers, and headless orchestration with Obsidian and LangGraph?

## Decision Drivers
- Zero breaking changes to existing CLI / TUI interface.
- Local-first, zero-leakage privacy for sensitive startup codebases.
- High performance sub-100ms response times.
- Frictionless workflow for non-prompt engineers.

## Considered Options
1. Cloud-only prompt rewriting API.
2. Direct integration within OpenCode TUI with local introspection, universal context parser, and headless daemon mode.

## Decision Outcome
Chosen Option: **Option 2**, because it guarantees complete privacy, works offline, and allows seamless integration with both human interactive terminals and headless automation orchestrators (Obsidian + LangGraph).
