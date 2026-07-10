---
source: coderabbit
pr: "20"
round: 1
round_created_at: "2026-07-10T19:46:39Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/model-catalog-and-stream-efficiency
head_sha: f73fd40d026660e67999ceb7cbb016d7b1c039ad
file: internal/agent/acpx_runner.go
line: 291
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6P_tlq,comment:PRRC_kwDOS0qyts7USTxl
review_hash: 2e6d46be0face60e977675604c1d751f7fbb87ced13db47504a70b103c14363d
duplicate_of: ""
source_review_id: "4674496937"
source_review_submitted_at: "2026-07-10T19:45:35Z"
---

# Issue 001: _ Stability & Availability_ _ Major_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _⚡ Quick win_

**Bound selection preflight setup with a deadline.**

`codexEnvForSession`, `sessions ensure`, and `set` use a potentially unbounded caller context. A hung adapter can block setup, doctor, or operational commands indefinitely; the cleanup timeout only starts afterward.

<details>
<summary>Proposed fix</summary>

```diff
+	acpxPreflightSetupTimeout       = 30 * time.Second
 	acpxPreflightCleanupTimeout     = 5 * time.Second
```

```diff
-	codexEnv, err := runner.codexEnvForSession(ctx, runtime, sessionName)
+	setupCtx, setupCancel := context.WithTimeout(ctx, acpxPreflightSetupTimeout)
+	codexEnv, err := runner.codexEnvForSession(setupCtx, runtime, sessionName)
 	if err != nil {
+		setupCancel()
 		return err
 	}
-	setupErr := runner.applyDisposableSelection(ctx, runtime, sessionName, workDir, codexEnv)
+	setupErr := runner.applyDisposableSelection(setupCtx, runtime, sessionName, workDir, codexEnv)
+	setupCancel()
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func (runner ACPXRunner) probeSelection(ctx context.Context, runtime RuntimeSpec, workDir string) error {
	if err := validateRuntimeSelection(runtime); err != nil {
		return err
	}
	sessionName, err := disposablePreflightSessionName()
	if err != nil {
		return err
	}
	setupCtx, setupCancel := context.WithTimeout(ctx, acpxPreflightSetupTimeout)
	codexEnv, err := runner.codexEnvForSession(setupCtx, runtime, sessionName)
	if err != nil {
		setupCancel()
		return err
	}
	setupErr := runner.applyDisposableSelection(setupCtx, runtime, sessionName, workDir, codexEnv)
	setupCancel()

	cleanupCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), acpxPreflightCleanupTimeout)
	defer cancel()
	cleanupErr := runner.closeDisposableSession(cleanupCtx, runtime, sessionName, workDir)
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/acpx_runner.go` around lines 275 - 291, Bound the disposable
selection setup phase in probeSelection with its own timeout context created
from ctx before calling codexEnvForSession and applyDisposableSelection,
ensuring environment lookup, session ensure, and selection set cannot run
indefinitely. Keep cleanup using its separate acpxPreflightCleanupTimeout
context and preserve the existing error handling and cancellation behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:64e6617f90fd48aae8a25412 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes:
  - Added a dedicated `30s` setup timeout context for disposable selection preflight setup so Codex environment resolution, `sessions ensure`, and selection `set` calls are bounded independently from cleanup.
  - Added regression coverage that verifies Codex setup inspection receives a context with a deadline.
  - Verification: `make verify` passed in this session. Result: `go test ./...` reported 1046 passed in 19 packages, `roundfix skills check` passed, and `go build` completed.
