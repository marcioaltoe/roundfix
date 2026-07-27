---
source: coderabbit
pr: "39"
round: 1
round_created_at: "2026-07-27T21:22:30Z"
status: invalid
terminal_reason: "Archived Spec 0038 makes retained-worktree guidance the sole trailing note, so emitting hidden guidance too would violate the shipped CLI contract."
head_repository: marcioaltoe/roundfix
head_branch: ma/terminal-run-worktree-reconciliation
head_sha: 44fa422cea404a2d8c47e4b8011f065c4c0481ba
file: internal/cli/runs.go
line: 127
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6UNfMU,comment:PRRC_kwDOS0qyts7aNkK2
review_hash: 1bf400f2fe9ff6c784bfa6c69c2be7ccb18d3f20e416b75dd5962a5ad277cacc
duplicate_of: ""
source_review_id: "4791610618"
source_review_submitted_at: "2026-07-27T21:21:25Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Retained-worktree note suppresses the hidden/older guidance.** When any terminal worktree is retained, the user no longer sees the `--state all` / `--limit 0` hint even though rows are hidden. Printing both notes keeps the pagination guidance intact.

<details>
<summary>🛠️ Proposed fix</summary>

```diff
-	note := runsListRetainedWorktreeNote(retainedWorktrees)
-	if note == "" {
-		note = runsListHiddenNote(opts.state, len(runs), len(matching), len(visible))
-	}
-	if note != "" {
-		fmt.Fprintln(stderr, note)
+	for _, note := range []string{
+		runsListHiddenNote(opts.state, len(runs), len(matching), len(visible)),
+		runsListRetainedWorktreeNote(retainedWorktrees),
+	} {
+		if note != "" {
+			fmt.Fprintln(stderr, note)
+		}
 	}
```
</details>

Note the added tests at `internal/cli/cli_test.go` assert exact stderr equality and would need updating.

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	for _, note := range []string{
		runsListHiddenNote(opts.state, len(runs), len(matching), len(visible)),
		runsListRetainedWorktreeNote(retainedWorktrees),
	} {
		if note != "" {
			fmt.Fprintln(stderr, note)
		}
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/runs.go` around lines 121 - 127, Update the note-rendering logic
in the runs-list flow so runsListRetainedWorktreeNote and runsListHiddenNote can
both be emitted when applicable, instead of the retained-worktree note
suppressing hidden/older guidance. Preserve the existing note formatting and
update the exact stderr expectations in the related CLI tests.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3b024e570dd1985ebafc9aa7 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The suppression is intentional. Archived Spec 0038 Task 05 requires
  the retained-worktree count to be the bounded stderr guidance and records
  hidden-row guidance as the fallback only when no retained residue exists.
  `docs/user-guide/commands.md` also states that the retained-worktree
  diagnostic is the one trailing stderr note. Printing both would change the
  accepted byte-stable CLI contract, so no code or test change is appropriate.
