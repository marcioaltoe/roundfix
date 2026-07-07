---
source: coderabbit
pr: "19"
round: 1
round_created_at: "2026-07-07T22:44:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/run-browser
head_sha: f2726493dff5e63e604139d27d147973ff650cf5
file: internal/cli/attach.go
line: 70
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6PCMCY,comment:PRRC_kwDOS0qyts7S9Vv1
review_hash: fda803cc0975c8526e81749b1a746693fbeccf4268fba923ad731d834ab9f9c2
duplicate_of: ""
source_review_id: "4648487653"
source_review_submitted_at: "2026-07-07T19:59:41Z"
---

# Issue 001: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Misleading remediation when attach has no Git repository.**

The failure directs the user to `roundfix runs list`, but per this same file's sibling check in `internal/cli/runbrowser.go` (Line 90-93) and the README's `runs list` boundary description, `runs list` itself exits `2` outside a Git repository unless `--all` is passed. This message should point to `roundfix runs list --all`, matching the parallel message in `runbrowser.go`.

<details>
<summary>🩹 Proposed fix</summary>

```diff
-			printAttachFailure(validationError{message: "attach without a Run ID requires a Git repository; run 'roundfix runs list' to discover Runs"}, stderr)
+			printAttachFailure(validationError{message: "attach without a Run ID requires a Git repository; run 'roundfix runs list --all' to discover Runs"}, stderr)
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		gitRoot := strings.TrimSpace(loaded.GitRoot)
		if gitRoot == "" {
			printAttachFailure(validationError{message: "attach without a Run ID requires a Git repository; run 'roundfix runs list --all' to discover Runs"}, stderr)
			return exitPreflight
		}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/attach.go` around lines 66 - 70, The attach preflight message in
attach should not point users to roundfix runs list by itself when there is no
Git repository, because that command also fails outside a repo unless --all is
used. Update the validationError message in the attach flow to reference
roundfix runs list --all, matching the sibling logic in runbrowser.go and
keeping the remediation consistent with the existing repository-boundary
behavior.
```

</details>

<!-- fingerprinting:phantom:poseidon:beignet -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d3d8750ff355a82ae1e16b97 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `valid`
- Notes: The original no-repository preflight path is no longer present, but the remaining non-interactive `attach` missing-Run-ID message still pointed to repository-scoped `roundfix runs list`. That remediation can fail outside a Git repository, so it now points to `roundfix runs list --all`.

## Resolution

- Status: `resolved`
- Changes: Updated `internal/cli/attach.go` and the CLI regression test to require the machine-wide discovery command.
- Verification: `rtk make verify` passed in this session for the configured `make verify` gate.
