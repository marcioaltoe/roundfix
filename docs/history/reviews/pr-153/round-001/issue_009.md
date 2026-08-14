---
source: coderabbit
pr: "153"
round: 1
round_created_at: "2026-08-10T21:31:46Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/a-proof-that-can-refuse
head_sha: d68e5a2a65875cf1f7a5e9976514c7be60ee5d5d
file: internal/agent/selection_catalogue_characterization_test.go
line: 4
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6YBWiy,comment:PRRC_kwDOS0qyts7fswOX
review_hash: 5006fa586becb708966579e2d34e6f1ad35fb651af53bc8f8b7d2826b642a51e
duplicate_of: ""
source_review_id: "4900622643"
source_review_submitted_at: "2026-08-10T20:19:14Z"
---

# Issue 009: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**The suite header contradicts the tests below it.**

Line 2 records the invariant as "today's proof accepts an echoed claude model". `TestSelectionCatalogueCharacterizationClaudeRefusesAnUnofferedModel` and its comment at lines 89-90 state the opposite: the catalogue now refuses that model. Line 4 places "the membership verdict" out of scope, but `TestRuntimeCatalogueRecordsAContaminatedAdvertisement` and the Claude refusal test both exercise the membership verdict.

A reader uses this header to learn the recorded invariant. Update it to the post-change behavior and drop the `task_03` reference.



<details>
<summary>📝 Proposed header</summary>

```diff
-// Suite: Agent Selection catalogue characterization.
-// Invariant: today's proof accepts an echoed claude model and refuses it on codex.
-// Boundary IN: ProveExactSelection and the fake ACPX process boundary.
-// Boundary OUT: the membership verdict added by task_03.
+// Suite: Agent Selection catalogue characterization.
+// Invariant: the pre-request Runtime Catalogue refuses an unadvertised Agent
+// Model on claude, and the codex adapter refuses it first on codex (ADR-0119).
+// Boundary IN: ProveExactSelection, applySessionSelection, and the fake ACPX
+// process boundary.
+// Boundary OUT: live ACP Runtimes and adapter binaries.
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
// Suite: Agent Selection catalogue characterization.
// Invariant: the pre-request Runtime Catalogue refuses an unadvertised Agent
// Model on claude, and the codex adapter refuses it first on codex (ADR-0119).
// Boundary IN: ProveExactSelection, applySessionSelection, and the fake ACPX
// process boundary.
// Boundary OUT: live ACP Runtimes and adapter binaries.
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/selection_catalogue_characterization_test.go` around lines 1 -
4, Update the suite header in the agent selection catalogue characterization
tests to state that an echoed, unoffered Claude model is refused, and revise the
boundary description to include the membership verdict rather than marking it
out of scope. Remove the task_03 reference while preserving the existing test
behavior.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fe52fb9fc15b72cc67845a30 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `RESOLVED`
- Notes: Updated the suite header in internal/agent/selection_catalogue_characterization_test.go to the post-change invariant: the pre-request Runtime Catalogue refuses an unadvertised Agent Model on claude, and the codex adapter refuses it first on codex (ADR-0119). Boundary IN now includes `applySessionSelection`, Boundary OUT is live ACP Runtimes, and the stale task_03 reference is dropped. Focused: `go test ./internal/agent -run 'TestSelectionCatalogueCharacterization'` ok.
