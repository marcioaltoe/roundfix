---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/cli/cli.go
line: 2632
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIpP,comment:PRRC_kwDOS0qyts7aUVDO
review_hash: 0248ec62b0308dd5df70123790fd35b47bb6f8587340108d89709433f6b5d82f
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:18Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Duplicate `daemon.Dependencies` literal.**

This engine construction is byte-for-byte identical to the one in `executeResolveCycle` (lines 2407-2417). Extract a single constructor so a future dependency addition cannot be applied to only one of the two call sites.



<details>
<summary>♻️ Proposed helper</summary>

```go
func newResolveEngine(collaborators engineCollaborators, runStore *store.Store, ui *runUI) (*daemon.Engine, error) {
	return daemon.NewEngine(daemon.Dependencies{
		Runner:    collaborators.runner,
		Verifier:  collaborators.verifier,
		Committer: collaborators.committer,
		Pusher:    collaborators.pusher,
		Source:    collaborators.source,
		Runs:      runStore,
		Worktree:  collaborators.worktree,
		Sink:      ui.sink,
		Progress:  ui.progress,
	})
}
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 2619 - 2632, Extract the duplicated
daemon.Dependencies construction from executeResolveCycle and the current engine
creation block into a shared newResolveEngine helper accepting
engineCollaborators, the run store, and run UI. Replace both call sites with
this helper while preserving their existing error propagation and engine
behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:refactor_suggestion -->

<!-- cr-comment:v1:acc6d4cff97185259decda88 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Extracted `newResolveEngine` and reused it from resolve and watch, removing the duplicated dependency assembly. Focused `internal/cli` tests passed.
