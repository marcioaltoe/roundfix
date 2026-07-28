---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: internal/store/agent_selection_test.go
line: 27
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USIp1,comment:PRRC_kwDOS0qyts7aUVEJ
review_hash: 1f1ea1f4374e4f5e098662ee6146b3841c2ca813348c57003d8f802267da8aa3
duplicate_of: ""
source_review_id: "4793702781"
source_review_submitted_at: "2026-07-28T04:25:19Z"
---

# Issue 012: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Use `schemaVersion` here too, matching the generalization applied in `internal/store/store_test.go`.**

The same PR replaced hard-coded version literals with `schemaVersion` across `store_test.go`; these two assertions still pin `11` and will need editing on the next bump.

<details>
<summary>♻️ Proposed change</summary>

```diff
-	if version != 11 {
-		t.Fatalf("expected schema version 11, got %d", version)
+	if version != schemaVersion {
+		t.Fatalf("expected schema version %d, got %d", schemaVersion, version)
 	}
```
</details>





Also applies to: 45-46

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/store/agent_selection_test.go` around lines 26 - 27, Update both
schema version assertions in the agent selection tests to compare against the
shared schemaVersion constant instead of the hard-coded literal 11, while
preserving the existing failure messages and version checks.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:30b7c7d235479efd38c30103 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Replaced the duplicated schema-version literals with `schemaVersion` in both agent-selection fixtures. Focused schema-version store tests passed.
