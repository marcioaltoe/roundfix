---
source: coderabbit
pr: "38"
round: 2
round_created_at: "2026-07-27T17:14:22Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-outcome-integrity
head_sha: 233964dd96ec27fc6bc709eda82b14457fc7f61b
file: internal/cli/cli.go
line: 655
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UG-Ox,comment:PRRC_kwDOS0qyts7aENB_
review_hash: 38ab7ac3d78edfeb2ec3221ebcaa7a9f5a0d7c859856542b99b0a4c88454346f
duplicate_of: ""
terminal_reason: "forceStopRun proves owner identity before session cleanup; identity-proof failure leaves Agent Sessions active"
source_review_id: "4788632386"
source_review_submitted_at: "2026-07-27T15:23:14Z"
---

# Issue 001: _ Stability & Availability_ _ Major_ _ Heavy lift_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _🏗️ Heavy lift_

**Agent Sessions are cancelled before owner identity is proven.**

`bestEffortForceStopAgentSessions` runs at Line 633, before the recorded-owner checks at Lines 634-655. When identity proof fails (PID reuse, unreadable identity), the Run is deliberately left Active with its lock retained — but its registered Agent Sessions have already been cancelled and closed, and their selections marked `closed`. The escape hatch then leaves a still-owned Active Run in a degraded state that the retry advice in the help text cannot repair.

Prove the owner first, then clean up sessions.



<details>
<summary>🛠️ Proposed reordering</summary>

```diff
-	warnings := bestEffortForceStopAgentSessions(ctx, runStore, active)
 	pid, ok := activeOwnerPID(active)
 	if !ok {
-		return stopResult{Run: active, Warnings: warnings}, forceStopOwnerError{
+		return stopResult{Run: active}, forceStopOwnerError{
 			RunID: active.ID,
 			PID:   0,
 			Step:  "validate recorded owner PID",
 			Err:   store.ErrOwnerProcessIdentityUnproven,
 		}
 	}
 	if err := ownerProcesses.TerminateAndWait(ctx, pid, active.OwnerIdentity); err != nil {
 		step := "prove owner exit"
 		var controlErr store.OwnerProcessControlError
 		if errors.As(err, &controlErr) && strings.TrimSpace(controlErr.Step) != "" {
 			step = controlErr.Step
 		}
-		return stopResult{Run: active, Warnings: warnings}, forceStopOwnerError{
+		return stopResult{Run: active}, forceStopOwnerError{
 			RunID: active.ID,
 			PID:   pid,
 			Step:  step,
 			Err:   err,
 		}
 	}
+	warnings := bestEffortForceStopAgentSessions(ctx, runStore, active)
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	pid, ok := activeOwnerPID(active)
	if !ok {
		return stopResult{Run: active}, forceStopOwnerError{
			RunID: active.ID,
			PID:   0,
			Step:  "validate recorded owner PID",
			Err:   store.ErrOwnerProcessIdentityUnproven,
		}
	}
	if err := ownerProcesses.TerminateAndWait(ctx, pid, active.OwnerIdentity); err != nil {
		step := "prove owner exit"
		var controlErr store.OwnerProcessControlError
		if errors.As(err, &controlErr) && strings.TrimSpace(controlErr.Step) != "" {
			step = controlErr.Step
		}
		return stopResult{Run: active}, forceStopOwnerError{
			RunID: active.ID,
			PID:   pid,
			Step:  step,
			Err:   err,
		}
	}
	warnings := bestEffortForceStopAgentSessions(ctx, runStore, active)
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 633 - 655, Reorder the stop flow so active
agent sessions are cleaned up only after the recorded owner identity has been
validated and TerminateAndWait succeeds. In the surrounding force-stop logic,
move bestEffortForceStopAgentSessions after the activeOwnerPID check and
ownerProcesses.TerminateAndWait call, while preserving warnings and error
results for validation or termination failures.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d9e1e07332056a348645f88c -->

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The finding does not apply to the current implementation. `forceStopRun` calls the read-only `ownerProcesses.ProveOwner` before `bestEffortForceStopAgentSessions`; an unprovable or reused owner therefore returns before any cancel or close. Moving cleanup after `TerminateAndWait` would also contradict the archived PRD and Task 03 requirement that Force Stop cancel registered Agent Sessions before terminating the recorded owner. Fresh Batch 001 evidence: `TestRunForceStopOwnerProofFailurePreservesAgentSessions` passed in `rtk proxy env GOCACHE=/tmp/roundfix-batch001-cli-cache go test ./internal/cli -run '^(TestCompletionWinnerOwnerVersusForceStopPublishesOneTerminalOutcome|TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion|TestRunForceStopOwnerProofFailurePreservesAgentSessions)$' -count=1`.
