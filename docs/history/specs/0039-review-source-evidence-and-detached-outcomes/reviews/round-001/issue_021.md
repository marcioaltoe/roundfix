---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/daemon/task_engine.go
line: 916
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Uf6Oc,comment:PRRC_kwDOS0qyts7aoLRT
review_hash: d11700affaf6a15bae0cb7367a3b81fb6ec013d22b926e83b307c0081547b3ff
duplicate_of: ""
source_review_id: "4800337236"
source_review_submitted_at: "2026-07-28T17:53:09Z"
---

# Issue 021: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Validate `retryUsed` before mutating Run state.**

`UpdateRunState(..., StateVerifying)` runs before the `retryUsed == nil` guard, so a programmer error leaves the Run in `Verifying` while the call fails. Move the guard above the state update so the argument check has no side effect.

<details>
<summary>🛠️ Proposed fix</summary>

```diff
 func (engine *Engine) verifyTask(ctx context.Context, plan TaskPlan, task spec.Task, ordinal int, attempt int, retryUsed *bool) (verificationAttemptOutcome, error) {
+	if retryUsed == nil {
+		return verificationAttemptOutcome{}, fmt.Errorf("verify run %q Task %s: temporary retry state is required", plan.RunID, task.ID)
+	}
 	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateVerifying); err != nil {
 		return verificationAttemptOutcome{}, fmt.Errorf("update run %q to state %q before Task %s verification: %w", plan.RunID, store.StateVerifying, task.ID, err)
 	}
-	if retryUsed == nil {
-		return verificationAttemptOutcome{}, fmt.Errorf("verify run %q Task %s: temporary retry state is required", plan.RunID, task.ID)
-	}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func (engine *Engine) verifyTask(ctx context.Context, plan TaskPlan, task spec.Task, ordinal int, attempt int, retryUsed *bool) (verificationAttemptOutcome, error) {
	if retryUsed == nil {
		return verificationAttemptOutcome{}, fmt.Errorf("verify run %q Task %s: temporary retry state is required", plan.RunID, task.ID)
	}
	if err := engine.deps.Runs.UpdateRunState(ctx, plan.RunID, store.StateVerifying); err != nil {
		return verificationAttemptOutcome{}, fmt.Errorf("update run %q to state %q before Task %s verification: %w", plan.RunID, store.StateVerifying, task.ID, err)
	}
	request := verificationAttemptRequest{
		RunID:                   plan.RunID,
		WorkDir:                 plan.WorkDir,
		ArtifactDir:             plan.ArtifactDir,
		BatchNumber:             ordinal,
		WorkItem:                task.ID,
		Attempt:                 attempt,
		Mode:                    verificationShared,
		Capacity:                plan.VerificationConcurrency,
		TemporaryRetryAvailable: !*retryUsed,
		Commands:                task.Verification,
		Publish: func(ctx context.Context, summary string, payload map[string]any) error {
			if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, task.ID, runevent.KindDaemonVerification, summary, payload); err != nil {
				return fmt.Errorf("publish verification event for run %q Task %s: %w", plan.RunID, task.ID, err)
			}
			return nil
		},
	}
	verification, err := engine.runTaskVerificationRequest(ctx, plan, task, request)
	if err != nil || verification.TemporaryFailure == nil || *retryUsed {
		return verification, err
	}

	*retryUsed = true
	request.Retry = 1
	request.Mode = verificationExclusive
	request.TemporaryRetryAvailable = false
	return engine.runTaskVerificationRequest(ctx, plan, task, request)
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/task_engine.go` around lines 881 - 916, Move the nil check
for retryUsed in verifyTask before the engine.deps.Runs.UpdateRunState call.
Keep the existing error and return behavior, ensuring invalid arguments are
rejected without mutating the run state; leave the verification request and
retry flow unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:8bf28cc8f8c57fa864b9411a -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Current HEAD updated the Run to `Verifying` before rejecting a nil
  temporary-retry state. Moved the argument guard ahead of the state mutation
  and added `TestVerifyTaskRejectsMissingRetryStateWithoutMutatingRun`.
  Focused evidence:
  `GOCACHE=/private/tmp/roundfix-batch-002-gocache rtk go test ./internal/daemon -run '^TestVerifyTaskRejectsMissingRetryStateWithoutMutatingRun$' -count=1`
  passed; the combined daemon, Run Event, and TUI package check passed 395
  tests. The Daemon owns the configured `make verify` run after this Agent
  turn.
