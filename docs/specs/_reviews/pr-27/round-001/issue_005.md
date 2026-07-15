---
source: coderabbit
pr: "27"
round: 1
round_created_at: "2026-07-15T15:54:46Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 43147cdff5f36ec1ac2bf276c3747400474d3fab
file: internal/cli/cli.go
line: 3976
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95M,comment:PRRC_kwDOS0qyts7V5tZ6
review_hash: 912fc09acf84f6f3f1fbac876cd1ca4ac283e655eae2817a72b1c6d32f69cf3a
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_005.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---




# Issue 005: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Do not attribute every uncommitted change to the Agent.**

Preflight permits existing untracked files, and Roundfix may write artifacts itself. This message can misidentify user-owned files; state only that the tracked tree was clean at preflight and ask the user to inspect current changes.

<details>
<summary>Proposed wording</summary>

```diff
-	fmt.Fprintln(stderr, "Uncommitted changes in the checkout are Agent work from this Run because Preflight started from a clean tracked tree.")
+	fmt.Fprintln(stderr, "The tracked checkout was clean at Preflight; inspect current tracked and untracked changes before retrying.")
```
</details>

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
func printAgentCheckoutChangesNotice(stderr io.Writer) {
	fmt.Fprintln(stderr, "The tracked checkout was clean at Preflight; inspect current tracked and untracked changes before retrying.")
}
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli.go` around lines 3975 - 3976, Update
printAgentCheckoutChangesNotice so it only states that the tracked tree was
clean when Preflight started and asks the user to inspect the current
uncommitted changes; remove the claim that all checkout changes are Agent work.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:81d9aecd90eda0ddcf47b4b0 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
