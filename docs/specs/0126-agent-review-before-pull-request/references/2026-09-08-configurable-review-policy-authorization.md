---
status: approved
granted: 2026-09-08
action: make pre-PR review configurable as codex, claude, coderabbit or none in canonical guidance
consuming: 0126-agent-review-before-pull-request
paths:
  - internal/baseline/assets/modules/core.json
  - internal/baseline/assets/modules/autonomous-work.json
  - docs/agents/agent-instructions.md
  - docs/agents/autonomous-work.md
  - docs/agents/setup-context.json
---

# Canonical configurable review policy authorization

The maintainer corrected the policy on 2026-09-08: "O review antes da PR pode
ou não existir, ou seja, pode ser code, claude, coderabbit ou none". This
session interprets "code" as Codex, preserving the previously selected default.
The maintainer also requires every canonical instruction to live in
`internal/baseline/assets/modules/` and its derived guides.

This instruction grants the named canonical policy and generated guidance:
Codex remains the default; explicit project selection may choose Codex,
Claude, CodeRabbit or none. CodeRabbit is optional. Explicit none omits review
and permits otherwise authorized publication/merge with required checks;
it records that review was disabled, never that review passed. An enabled
reviewer's failure or absence does not select none. Existing QA, exact-head
checks, permissions and execution limits remain operative.

The grant covers only the five listed files and sanctioned digest regeneration
from these source changes. It does not implement configuration parsing,
provider adapters, skip receipts or durable queue behavior, alter the current
project's selected mode, enable a provider service, change credentials or
approve paid calls. Other implementation paths remain under the proposed
Spec authorization. Historical review evidence retains its original meaning.

Commit this record separately before consuming tooling changes. The known
regeneration/verification compatibility failures remain prerequisites to
publishing the resulting tooling change as verified.
