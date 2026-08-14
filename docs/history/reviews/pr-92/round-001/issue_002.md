---
source: coderabbit
pr: "92"
round: 1
round_created_at: "2026-08-03T19:34:00Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spawn-economy
head_sha: 7765d1f6d62e59ebf68ca2e4e2e273733da58425
file: internal/baseline/skills_restore_git_test.go
line: 112
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WGXK0,comment:PRRC_kwDOS0qyts7c9hS-
review_hash: 8a684c3818cadd74ce2708db3fc232a2c450f536d8063b2c27b3198592801978
duplicate_of: ""
source_review_id: "4847882119"
source_review_submitted_at: "2026-08-03T19:33:10Z"
---

# Issue 002: _ Stability & Availability_ _ Trivial_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🔵 Trivial_ | _⚡ Quick win_

**Bind the helper subprocess to the test context.**

`exec.Command` gives the helper process no deadline and no cancellation. If the helper ever blocks before writing its partial frame, `Read` blocks until the package test timeout instead of failing this test. The production constructor already uses `exec.CommandContext`.

<details>
<summary>♻️ Proposed change</summary>

```diff
-	command := exec.Command(os.Args[0], "-test.run=^TestBatchObjectReaderProcessHelper$")
+	command := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestBatchObjectReaderProcessHelper$")
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	command := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^TestBatchObjectReaderProcessHelper$")
	command.Env = append(os.Environ(), "ROUNDFIX_BATCH_OBJECT_HELPER=die-mid-stream")
	reader, err := startBatchObjectReader(command)
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/skills_restore_git_test.go` around lines 110 - 112, Update
the helper command setup before startBatchObjectReader to use the test context
with exec.CommandContext instead of exec.Command, preserving the existing
arguments and environment so a blocked helper subprocess is canceled when the
test context ends.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:ffdcac2511746733e590375e -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The helper process used `exec.Command`, so the test could wait until the package timeout if the helper blocked before emitting its partial frame.

## Result

- Bound the helper subprocess to `t.Context()` with `exec.CommandContext` while preserving its arguments and environment.
- Focused check: `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache-run_20260803T193424Z-baseline go test ./internal/baseline -run '^TestBatchObjectReader(ReportsProcessDeathMidStream|RejectsOversizedObjectsAndPoisonsProtocolFailures)$' -count=1` — passed.
- Package check: `rtk proxy env GOCACHE=/private/tmp/roundfix-gocache-run_20260803T193424Z-baseline go test ./internal/baseline -count=1` — passed.
- The Daemon owns the configured `make verify` run.
