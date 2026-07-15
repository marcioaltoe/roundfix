---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: internal/daemon/task_engine.go
line: 928
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95n,comment:PRRC_kwDOS0qyts7V5tac
review_hash: 5a92561cc761ad3b033585815e90603f1589055cedef665ed2371e8327b1cf2a
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 012: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Propagate progress-writer failures.**

Line 918 ignores `fmt.Fprint`'s error, concealing failed diagnostic output and triggering `errcheck`.

<details>
<summary>Proposed fix</summary>

```diff
 	warning := fmt.Sprintf("roundfix: warning: Task %s completed with no changes outside the Spec Root (%s)\n", taskID, shape)
-	fmt.Fprint(engine.deps.Progress, warning)
+	if _, err := fmt.Fprint(engine.deps.Progress, warning); err != nil {
+		return fmt.Errorf("write no-op Task commit warning for run %q Task %s: %w", plan.RunID, taskID, err)
+	}
```
</details>

As per coding guidelines, “Wrap errors with `%w`.”

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func (engine *Engine) publishNoOpTaskCommitWarning(ctx context.Context, plan TaskPlan, taskID string, ordinal int, shape string) error {
	if shape == "" {
		return nil
	}
	warning := fmt.Sprintf("roundfix: warning: Task %s completed with no changes outside the Spec Root (%s)\n", taskID, shape)
	if _, err := fmt.Fprint(engine.deps.Progress, warning); err != nil {
		return fmt.Errorf("write no-op Task commit warning for run %q Task %s: %w", plan.RunID, taskID, err)
	}
	summary := fmt.Sprintf("Task %s completed with no changes outside the Spec Root (%s).", taskID, shape)
	payload := map[string]any{
		"decision": "warning",
		"warning":  "no_op_task_commit",
		"task":     taskID,
		"shape":    shape,
	}
	if err := engine.publishTaskEvent(ctx, plan.RunID, ordinal, taskID, runevent.KindDaemonCommit, summary, payload); err != nil {
		return fmt.Errorf("publish no-op Task commit warning for run %q Task %s: %w", plan.RunID, taskID, err)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 golangci-lint (2.12.2)</summary>

[error] 918-918: Error return value of `fmt.Fprint` is not checked

(errcheck)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/daemon/task_engine.go` around lines 913 - 928, The
publishNoOpTaskCommitWarning method ignores the error returned by fmt.Fprint
when writing to engine.deps.Progress. Capture that error, wrap it with
contextual information using %w, and return it before constructing or publishing
the task event; preserve the existing successful-output and empty-shape
behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:ec480ba9452bedba94d97f58 -->

_Sources: Coding guidelines, Linters/SAST tools_

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `publishNoOpTaskCommitWarning` now wraps and returns progress-writer failures before publishing the event; `make verify` passed.
