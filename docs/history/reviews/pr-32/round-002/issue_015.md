---
source: coderabbit
pr: "32"
round: 2
round_created_at: "2026-07-17T13:23:47Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: d7ab1933ac9fdcf0c94d73e2f417d99d38e43fe7
file: internal/tui/timeline_group.go
line: 185
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5tR,comment:PRRC_kwDOS0qyts7Wt95y
review_hash: cd2dc3a525d4c5abcc7d39503a1c542de13a0efada09fcc9cd82f36a0675ab3d
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---

# Issue 015: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Fall back to the persisted summary when projection fails.**

Returning `""` suppresses the entire timeline event when a legacy or malformed selection payload cannot be projected, despite `event.Summary` or `entry.text` being readable.

<details>
<summary>Proposed fix</summary>

```diff
-	if record, ok, err := runevent.ProjectSelectionLifecycle(entry.event); ok {
-		if err != nil {
-			return ""
-		}
+	if record, ok, err := runevent.ProjectSelectionLifecycle(entry.event); ok && err == nil {
 		return runevent.SelectionLifecycleSummary(record)
 	}
```
</details>

This also preserves the PR’s legacy Run-summary compatibility requirement.

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
	if record, ok, err := runevent.ProjectSelectionLifecycle(entry.event); ok && err == nil {
		return runevent.SelectionLifecycleSummary(record)
	}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/tui/timeline_group.go` around lines 180 - 185, Update the
selection-summary logic around ProjectSelectionLifecycle so projection errors
fall back to the persisted event summary, using event.Summary or entry.text when
available, instead of returning an empty string. Preserve the projected
SelectionLifecycleSummary result when projection succeeds and maintain legacy
run-summary compatibility.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:14b10afb7ab68af539119bbd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Fixed TUI timeline selection summaries to fall back to persisted event summaries when selection lifecycle projection fails, with malformed-payload coverage. Evidence: `GOCACHE=/private/tmp/roundfix-go-build rtk go test ./internal/agent ./internal/cli ./internal/config ./internal/daemon ./internal/releaseplan ./internal/spec ./internal/store ./internal/tui` passed.
