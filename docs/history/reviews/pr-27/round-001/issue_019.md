---
source: coderabbit
pr: "27"
round: 1
round_created_at: "2026-07-15T15:54:46Z"
status: duplicated
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 43147cdff5f36ec1ac2bf276c3747400474d3fab
file: skills/write-tasks/references/task-template.md
line: 95
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ96O,comment:PRRC_kwDOS0qyts7V5tbL
review_hash: 895697236d42a2fa9f1c59fcd3c3c0ee7fd5928a150100c14a7a82f824d3c464
duplicate_of: /Users/marcio/dev/roundfix/docs/specs/_reviews/pr-27/round-002/issue_019.md
terminal_reason: 'Agent failed: Agent Batch failed after acpx exited with code 1: agent/protocol error'
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---




# Issue 019: _ Maintainability & Code Quality_ _ Minor_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟡 Minor_ | _⚡ Quick win_

**Use `rg` for repository searches in generated task gates.**

This template currently propagates guidance that contradicts the repository’s required search tooling.

<details>
<summary>Proposed correction</summary>

```diff
-     Use portable shell forms: prefer grep over rg in task gates, avoid wc-pipeline shape checks, use repository build flags such as go build -buildvcs=false ./... when a build is required, and include executable checks that prove the Task's effect. -->
+     Use portable shell forms: use rg or rg --files for repository searches, avoid wc-pipeline shape checks, use repository build flags such as go build -buildvcs=false ./... when a build is required, and include executable checks that prove the Task's effect. -->
```
</details>





As per coding guidelines, “Use `rg` or `rg --files` for local code search.”

<!-- suggestion_start -->

<details>
<summary>📝 Committable suggestion</summary>

> ‼️ **IMPORTANT**
> Carefully review the code before committing. Ensure that it accurately replaces the highlighted code, contains no missing lines, and has no issues with indentation. Thoroughly test & benchmark the code to ensure it meets the requirements.

```suggestion
     Use portable shell forms: use rg or rg --files for repository searches, avoid wc-pipeline shape checks, use repository build flags such as go build -buildvcs=false ./... when a build is required, and include executable checks that prove the Task's effect. -->
```

</details>

<!-- suggestion_end -->

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@skills/write-tasks/references/task-template.md` at line 95, Update the
task-template guidance comment to require rg or rg --files for repository and
local code searches in generated task gates, replacing the instruction to prefer
grep. Preserve the existing guidance about avoiding wc-pipeline checks, using
repository build flags, and including executable effect checks.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c4142c620e6f368fd04e695d -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
