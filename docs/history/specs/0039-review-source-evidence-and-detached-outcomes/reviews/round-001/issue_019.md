---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/cli/implement.go
line: 284
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Uf6OY,comment:PRRC_kwDOS0qyts7aoLRO
review_hash: 8390d4de531303f8dcef2cd66a500af81c1eceb39e922feee51069b70d6f05c7
duplicate_of: ""
source_review_id: "4800337236"
source_review_submitted_at: "2026-07-28T17:53:09Z"
---

# Issue 019: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Two adjacent `int` parameters are silently swappable.**

`taskCapacity` and `verificationCapacity` sit next to each other in a 20-parameter positional signature; transposing them at the single call site compiles and produces a wrong-but-plausible schedule. Group them in a small value type so the call site names them.

<details>
<summary>♻️ Suggested shape</summary>

```go
type implementCapacities struct {
	task         int
	verification int
}

func executeImplementCycle(..., capacities implementCapacities, ...) (daemon.TaskCycleResult, error)
```

Call site:

```go
implementCapacities{
	task:         loadedConfig.Config.Worktree.Concurrency,
	verification: loadedConfig.Config.Verification.Concurrency,
}
```
</details>





Also applies to: 620-620

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/implement.go` at line 284, Replace the adjacent positional task
and verification capacity parameters in executeImplementCycle with a named
implementCapacities value containing task and verification fields. Update the
executeImplementCycle declaration, its internal references, and both call sites
to construct implementCapacities with the corresponding Worktree.Concurrency and
Verification.Concurrency values.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:60e50984e34e33034ac53329 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Grouped task and verification concurrency in `implementCapacities`, removing the adjacent same-typed positional parameters from `executeImplementCycle`. Four focused implement-capacity tests passed.
