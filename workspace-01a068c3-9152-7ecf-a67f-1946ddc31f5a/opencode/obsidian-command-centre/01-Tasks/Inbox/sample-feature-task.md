---
title: Implement Secure Webhook Verification Middleware
status: pending
priority: high
target_repo: .
tools_allowed:
  - edit
  - bash
  - tool:grep
require_human_approval: false
---

# Task: Implement Webhook Verification Middleware

## 🎯 Goal
Implement HMAC-SHA256 signature verification for incoming payment webhooks.

## 📝 Context & Mentions
- Target File: `@internal/auth/jwt.go`
- Skill Template: `/startup:api-scaffold`
- Project Rule: `#rules:security`
- Testing: `/test`

## 💡 Specifications
1. Read secret key from environment variable `WEBHOOK_SIGNING_SECRET`.
2. Compute HMAC-SHA256 against raw request payload.
3. Compare against `X-Signature-SHA256` header using constant-time comparison.
4. Add unit test suite testing valid, invalid, and expired signatures.
