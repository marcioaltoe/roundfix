---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: invalid
terminal_reason: "The Waves omission is historical drift in archived Spec 0065, whose artifacts must remain byte-identical."
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: docs/specs/_archived/0065-loop-order-and-verification-honesty/_tasks.md
line: 43
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v5l,comment:PRRC_kwDOS0qyts7eEK6y
review_hash: dd004e15867ba265d626ad4ca9c2d96db5a77711ca427ed1c5760dfefee055eb
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:28Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Include `task_07` in the execution waves.**

The graph requires `task_07` before terminal QA task `task_06`, but the Waves summary jumps from `task_05` to `task_06`. Regenerate it as `3 → task_05 · 4 → task_07 · 5 → task_06` so the projection matches the graph.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0065-loop-order-and-verification-honesty/_tasks.md`
around lines 42 - 43, Update the Waves summary to include task_07 between
task_05 and task_06, using the sequence 3 → task_05 · 4 → task_07 · 5 → task_06
so it matches the execution graph.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:47ef44ea3724cd36c5508c48 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The graph includes `task_07` before terminal `task_06`, while the
  Waves prose omits it. That inconsistency is in archived Spec 0065, which the
  repository forbids backfilling or rewriting after completion.
- Daemon Verification: `make verify` not run; Daemon-owned.
