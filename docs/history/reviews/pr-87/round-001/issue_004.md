---
source: coderabbit
pr: "87"
round: 1
round_created_at: "2026-08-03T15:34:03Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/qa-task-not-flag
head_sha: a12c1a665c5970773e04c4a145c6b9b0c5a0e686
file: internal/daemon/task_engine.go
line: 306
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WBeNt,comment:PRRC_kwDOS0qyts7c2V0R
review_hash: 100c49cf00f7637d194debd51e6a48bec4747e92f5c12c0c7d9251c77622c114
duplicate_of: ""
source_review_id: "4845660382"
source_review_submitted_at: "2026-08-03T15:14:34Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Report the withheld gate instead of ending the cycle silently.**

When a gate dependency is not completed, the condition at Line 301 skips the QA step with no progress line and no Run Event. The cycle then finishes with an empty `QAVerdict` and no explanation, and `assertNoQAStep` confirms no `daemon.qa` event is published. A user who reads the outcome cannot distinguish "no gate declared" from "gate withheld until dependencies complete".

Write a progress line, or publish a status event naming the unsettled dependencies, when a declared gate is withheld.

<details>
<summary>♻️ Suggested addition</summary>

```diff
 	if qaTask != nil && qaTask.Status != spec.StatusCompleted && taskNeedsCompleted(*qaTask, statuses) {
 		if err := engine.stopIfRequested(ctx, plan.RunID, ordinal+1); err != nil {
 			return result, fmt.Errorf("stop run %q before the QA step: %w", plan.RunID, err)
 		}
 		ordinal++
 		verdict, reportPath, err := engine.runQAGate(ctx, plan, *qaTask, ordinal)
 		if err != nil {
 			return result, err
 		}
 		result.QAVerdict = verdict
 		result.QAReportPath = reportPath
+	} else if qaTask != nil && qaTask.Status != spec.StatusCompleted {
+		fmt.Fprintf(engine.deps.Progress, "QA step withheld: Task %s still waits on its dependencies\n", qaTask.ID)
 	}
```
</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/task_engine.go` around lines 292 - 306, Update the QA
handling around runTaskScheduler and runQAGate so a declared qaTask whose
dependencies are not yet completed emits progress or a Run Event identifying the
unsettled dependencies before the cycle ends. Preserve the existing
pending/resumable behavior and continue running the gate only when
taskNeedsCompleted returns true, while ensuring the no-gate path remains
distinguishable.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:5b74392e61720cce28b9b2d5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: A declared unsettled gate now writes a progress line naming the gate and every uncompleted direct dependency while remaining pending and emitting no `daemon.qa` event. The regression reproduced the previously silent path and `go test ./internal/daemon -count=1` passed with 171 tests.
