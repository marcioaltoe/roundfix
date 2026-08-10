---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/agent/selection_assignment_test.go
line: 200
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWij,comment:PRRC_kwDOS0qyts7fswOB
review_hash: cf399d4c33d7b2423320d12258efd6918522305e67d65e2646d373a061c6a7bc
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:14Z"
---

# Issue 006: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Derive at least one table assignment from `PlanSelectionAssignment`.**

Every case hand-writes the `SelectionAssignment` and asserts `selectionStateMatches` accepts it. The defect this PR fixes (`claude/opus/high` reported `effective_selection_mismatch`) was a disagreement between the planner and the matcher. A hand-written assignment cannot detect that disagreement, so the "claude opus normalizes its echoed alias" case can pass while the real pipeline still mismatches.

`runtimeDeferredSelectionState` (line 254) already plans the assignment from the capabilities. Apply the same approach to the alias case so the planner and the matcher are proved together.



<details>
<summary>♻️ Proposed change for the alias case</summary>

```diff
 	tests := []struct {
 		name       string
 		assignment SelectionAssignment
 		state      SelectionCapabilities
+		runtime    RuntimeSpec
 	}{
```

Then, for the alias case, plan the assignment instead of literal fields:

```go
{
	name:    "claude opus normalizes its echoed alias",
	runtime: RuntimeSpec{ID: "claude", Model: "opus", ReasoningEffort: "high"},
	state:   matchingIndependentSelectionState("opus[1m]", "opus", "effort", "high"),
},
```

In the subtest, when `tt.runtime.ID != ""`, call `PlanSelectionAssignment(tt.runtime, tt.state)` and use its result as the assignment.
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/selection_assignment_test.go` around lines 137 - 200, Update
TestProofAcceptsAMatchingSelectionAmongSiblings so the Claude Opus alias case
derives its assignment through PlanSelectionAssignment instead of hard-coded
SelectionAssignment fields. Add the corresponding RuntimeSpec test-table field,
populate it for that case, and in the subtest call PlanSelectionAssignment when
RuntimeSpec.ID is set before invoking selectionStateMatches; preserve literal
assignments for the remaining cases.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:22a74a4cd365153bbd00934c -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: The "claude opus normalizes its echoed alias" case in `TestProofAcceptsAMatchingSelectionAmongSiblings` now plans its assignment through `PlanSelectionAssignment(runtime, state)` instead of hard-coded `SelectionAssignment` fields. Added a `runtime RuntimeSpec` table field populated only for that case; the subtest derives the assignment from the planner when `runtime.ID != ""`, so the planner and matcher are proved together. The remaining cases keep their literal assignments. Focused: `go test ./internal/agent -run 'TestProofAcceptsAMatchingSelectionAmongSiblings|TestProofStillRejects'` ok.
