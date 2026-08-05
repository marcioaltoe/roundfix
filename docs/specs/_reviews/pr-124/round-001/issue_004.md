---
source: coderabbit
pr: "124"
round: 1
round_created_at: "2026-08-05T16:50:26Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0077-a-green-check-is-not-a-review
head_sha: 4a03df27595a73705316edfb149bea641e3b5772
file: internal/cli/cli.go
line: 5087
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wuaz6,comment:PRRC_kwDOS0qyts7d35tQ
review_hash: 2352a7a129bdb6b64dea3fc586e8943dc7a434902c928bc6e989cf98e3b2fcad
duplicate_of: ""
source_review_id: "4866751340"
source_review_submitted_at: "2026-08-05T16:49:39Z"
---

# Issue 004: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Rewrap the new help text and fix the "resolves pending" phrasing.**

Two text issues exist in this block:

1. Lines 5084 and 5085 are much longer than the surrounding help text, which wraps near 72 characters. Line 5085 runs past 140 characters. In a narrow terminal the paragraph wraps inconsistently against the rest of the usage output.
2. "An unrecognised signal resolves pending" and "refusal resolves skipped" omit the preposition. "resolves to pending" and "resolves to skipped" read correctly.

`cli_test.go` Line 403 asserts the exact substrings "unrecognised signal resolves pending" and "explicit Review Source refusal resolves skipped". Update those expectations together with the text.



<details>
<summary>📝 Proposed rewrap and wording fix</summary>

```diff
   a current-head CodeRabbit APPROVED review is also accepted, with zero
-  unresolved CodeRabbit threads. An unrecognised signal resolves pending even
-  when its check conclusion is success: a green check is not evidence that a review ran.
-  An explicit Review Source refusal resolves skipped; watch will not merge that head or clear it for merge. If no accepted Evidence appears
-  within the grace period, watch ends CleanUnverified and exits 3. Omit all Agent Selection flags
-  to use the review profile. A one-Run override requires --agent, --model, and --reasoning-effort together.
+  unresolved CodeRabbit threads. An unrecognised signal resolves to pending
+  even when its check conclusion is success: a green check is not evidence
+  that a review ran. An explicit Review Source refusal resolves to skipped;
+  watch will not merge that head or clear it for merge. If no accepted
+  Evidence appears within the grace period, watch ends CleanUnverified and
+  exits 3. Omit all Agent Selection flags to use the review profile. A
+  one-Run override requires --agent, --model, and --reasoning-effort together.
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
  Worktree. With --until-clean, Clean requires accepted Review Source Evidence
  on the pushed head. The only check-or-status route to a verified head is a
  recognised review-completed current-head CodeRabbit check or commit status;
  a current-head CodeRabbit APPROVED review is also accepted, with zero
  unresolved CodeRabbit threads. An unrecognised signal resolves to pending
  even when its check conclusion is success: a green check is not evidence
  that a review ran. An explicit Review Source refusal resolves to skipped;
  watch will not merge that head or clear it for merge. If no accepted
  Evidence appears within the grace period, watch ends CleanUnverified and
  exits 3. Omit all Agent Selection flags to use the review profile. A
  one-Run override requires --agent, --model, and --reasoning-effort together.
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 5079 - 5087, Rewrap the help text in the
CLI help block around the existing worktree/clean guidance so it matches the
surrounding 72-character formatting, and update the phrasing in that block to
use “resolves to pending” and “resolves to skipped” in the Help/usage string
built by the relevant cli.go help text path. Also update the corresponding
cli_test.go substring assertions to match the revised wording so the help output
and test expectations stay aligned.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:251b71157167c21ba31199a7 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Watch help now uses “resolves to pending” and “resolves to skipped”
  and wraps the paragraph consistently with the surrounding usage text. The
  exact help assertions were updated with the corrected phrases.
- Focused evidence: the updated help assertion failed against the old wording,
  then `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test
  ./internal/cli -count=1 -run
  '^(TestRunCommandHelp|TestCommandUsageDocumentsProfileLedAndCompleteSelectionOverrides)$'`
  passed against the final wrapped text; every changed help line is at most 72
  characters. A broader package attempt reached an unrelated concurrent Git
  worktree test and failed with `commondir: Result too large`; authoritative
  Verification remains Daemon-owned.
- Daemon Verification: `make verify` not run; Daemon-owned.
