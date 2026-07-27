---
source: coderabbit
pr: "38"
round: 1
round_created_at: "2026-07-27T15:34:32Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-outcome-integrity
head_sha: 9ed57622bb92f138aa3e23d4d59e260ebbff0116
file: internal/cli/cli.go
line: 655
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UG-Ox,comment:PRRC_kwDOS0qyts7aENB_
review_hash: d17e9ce472ea04a8112c7b81e638964a491a5a8f02d58d3b5ed7041247d12b70
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260727T152947Z_936cd84aa803ba5d/verification/batch-001-attempt-2.log'
source_review_id: "4788632386"
source_review_submitted_at: "2026-07-27T15:23:14Z"
---


# Issue 002: _ Stability & Availability_ _ Major_ _ Heavy lift_

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

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `forceStopRun` mutated registered Agent Sessions before validating and terminating the recorded owner, leaving an Active Run degraded when owner proof failed. Session cleanup now starts only after `TerminateAndWait` succeeds, and the regression asserts failed owner proof leaves the active selection untouched. The canonical and embedded Roundfix skill contract was synchronized with `rtk make skills-sync`. Verification Feedback identified the corresponding TypeScript/Bun setup snapshot digest as stale; that single generated contract field now matches the canonical skill. Focused evidence: `rtk proxy env GOCACHE=/tmp/roundfix-run-936cd84aa803ba5d-gocache go test ./internal/cli -run '^(TestCompletionWinnerOwnerVersusForceStopPublishesOneTerminalOutcome|TestRunForceStopOwnerFailurePreservesAgentSessions)$' -count=1` passed; `TestRunStopForceRegisteredAgentSessionCleanupTargetsActiveScopesInOrder` also passed with owner proof first; `rtk proxy env GOCACHE=/tmp/roundfix-run-936cd84aa803ba5d-gocache go test -count=1 ./skills -run '^TestAuthorialSkillSync/typescript-bun.json$'` and `rtk make skills-sync-check` passed after the digest refresh.
