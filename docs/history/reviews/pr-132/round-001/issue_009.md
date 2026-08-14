---
source: coderabbit
pr: "132"
round: 1
round_created_at: "2026-08-06T09:54:40Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0073-skill-versions-decoupled-from-the-binary
head_sha: 8cde14417b3d169f259d8e0cf3ed0d6930f50f0e
file: skills/baseline_skill_contract_integration_test.go
line: 42
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W6kag,comment:PRRC_kwDOS0qyts7eJyli
review_hash: 7062e11afd9c88bfee3cecdfdec40a985a06abb9592d9dea8822cfec0f2f32cb
duplicate_of: ""
source_review_id: "4872547928"
source_review_submitted_at: "2026-08-06T08:19:10Z"
---

# Issue 009: _ Stability & Availability_ _ Minor_ _ Quick win_

## Review Comment

_🩺 Stability & Availability_ | _🟡 Minor_ | _⚡ Quick win_

**Do not discard the close error.**

If `WriteString` fails, the code discards the `file.Close()` result. Capture both results and fail the test when either operation fails.

As per coding guidelines, `**/*.go`: “Always check returned errors; never discard them with `_`.”






<details>
<summary>Proposed fix</summary>

```diff
-		if _, err := file.WriteString("\n<!-- compatibility-preserving owned skill edit -->\n"); err != nil {
-			_ = file.Close()
-			t.Fatalf("edit owned skill %s: %v", relative, err)
-		}
-		if err := file.Close(); err != nil {
-			t.Fatalf("close owned skill edit target %s: %v", relative, err)
+		_, writeErr := file.WriteString("\n<!-- compatibility-preserving owned skill edit -->\n")
+		closeErr := file.Close()
+		if writeErr != nil || closeErr != nil {
+			t.Fatalf("edit owned skill %s: write error = %v, close error = %v", relative, writeErr, closeErr)
 		}
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
		_, writeErr := file.WriteString("\n<!-- compatibility-preserving owned skill edit -->\n")
		closeErr := file.Close()
		if writeErr != nil || closeErr != nil {
			t.Fatalf("edit owned skill %s: write error = %v, close error = %v", relative, writeErr, closeErr)
		}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/baseline_skill_contract_integration_test.go` around lines 37 - 42,
Update the file-edit error path in the integration test to capture the result of
file.Close instead of discarding it with _. Ensure the test fails when either
WriteString or Close returns an error, while preserving the existing error
context for the owned skill edit target.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:d7148ac882f47c182e727050 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The write-failure branch discarded `file.Close()` and could hide a second failure. The test now captures both `writeErr` and `closeErr`, closes the file unconditionally after the write attempt, and reports both results with the owned Skill path. Focused evidence: `rtk go test ./skills -count=1 -run '^TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical$' -v` with the repository-local Go cache exited 0. Authoritative `make verify` remains Daemon-owned.
