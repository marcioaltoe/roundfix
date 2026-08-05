---
source: coderabbit
pr: "113"
round: 1
round_created_at: "2026-08-05T02:12:07Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0068-implementation
head_sha: c9af2617f988bd63e1bd8f22c6178758a8e5fd40
file: internal/specaudit/audit_test.go
line: 147
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6WhZqJ,comment:PRRC_kwDOS0qyts7dk96y
review_hash: 49c1cad5d4982ffb24721cfeb7ef53e4ad7ac5307bcc49537f2d617ca4282bbf
duplicate_of: ""
source_review_id: "4860420451"
source_review_submitted_at: "2026-08-05T02:11:27Z"
---

# Issue 003: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Use `errors.Is` for the sentinel check.**

`os.IsNotExist` does not traverse wrapped errors. The coding guidelines require `errors.Is` for sentinel matching, and `errorlint` is in the enabled linter set.





<details>
<summary>♻️ Proposed fix</summary>

```diff
-	if _, err := os.Stat(worktreePath); !os.IsNotExist(err) {
+	if _, err := os.Stat(worktreePath); !errors.Is(err, os.ErrNotExist) {
 		t.Fatalf("reclaimed worktree stat error = %v, want not exists", err)
 	}
```

Add `"errors"` to the import block.
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	if _, err := os.Stat(worktreePath); !errors.Is(err, os.ErrNotExist) {
		t.Fatalf("reclaimed worktree stat error = %v, want not exists", err)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/specaudit/audit_test.go` around lines 145 - 147, Update the
reclaimed worktree existence assertion in the audit test to use errors.Is with
os.ErrNotExist instead of os.IsNotExist, and add the errors import to the
existing import block.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:eacb064aad8060e8a052c6de -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The reclaimed-worktree assertion used `os.IsNotExist`, which does not provide the required wrapped-sentinel matching contract.

## Resolution

- Replaced the check with `errors.Is(err, os.ErrNotExist)` and added the standard-library import.
- Focused evidence: `rtk env GOCACHE=/private/tmp/roundfix-review0068-audit-cache go test ./internal/specaudit -run '^(TestAuditClassifiesResidueBranch|TestAuditScratchWorktreeReclaimCommandRuns|TestAuditPreservesActiveRunSurvivors)$' -count=1` exited 0; the full affected-package command also exited 0.
- Daemon Verification: `make verify` was not run; the Daemon owns that command.
