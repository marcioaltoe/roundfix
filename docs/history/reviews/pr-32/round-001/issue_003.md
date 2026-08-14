---
source: coderabbit
pr: "32"
round: 1
round_created_at: "2026-07-17T10:26:16Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: f7ff075d90b898620702e0d2c3a736020b4750d3
file: internal/cli/implement.go
line: 154
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5sT,comment:PRRC_kwDOS0qyts7Wt94m
review_hash: 38b46531716124e8533a4df77b47d9aeaf40982691b213ce766e9b5a2931daef
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-32/round-002/issue_003.md
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---


# Issue 003: _ Functional Correctness_ _ Major_ _ Heavy lift_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _🏗️ Heavy lift_

**Resolve profile categories from the same committed graph that executes.**

`graph` is loaded from the potentially dirty user checkout, while the cycle later uses `executionGraph` from the committed Run Worktree. An uncommitted task-type change can therefore validate/select the wrong profiles and fail only after creating the Run.

Load the committed graph before operational preflight and use that graph for both category selection and execution.

As per coding guidelines, “Do not use workarounds; fix the root cause instead.”

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/implement.go` around lines 137 - 154, Update the implement flow
to load the committed Run Worktree graph before calling
runProfileOperationalPreflight, and derive categories through
implementProfileCategories from that committed graph. Use the same
executionGraph for operational profile validation, runtime selection, and the
later execution cycle; do not use the potentially dirty checkout graph for these
decisions.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:dcd3488936385371c9831f29 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
