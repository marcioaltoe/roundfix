---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: internal/spec/spec.go
line: 98
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v6X,comment:PRRC_kwDOS0qyts7eEK7u
review_hash: dd654ea6b2d15cf827d3e24151d8648a1b35738d2f28066111a01215a93da743
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:28Z"
---

# Issue 022: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Preserve declaration source locations.**

`Requirements` and `RehearsalCases` discard their Markdown line numbers. `ContradictoryRequirements` then reports `requirementIndex + 1`, and `UndeclaredRehearsal` reports line `1`. The replay contract expects the actual declaration lines.

Store each declaration as text plus source line. Populate that metadata during parsing and use it when creating findings. This prevents diagnostics from linking users to unrelated lines.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/spec.go` around lines 97 - 98, Update the declarations
represented by Requirements and RehearsalCases to retain each Markdown entry’s
text and source line during parsing, then use those stored lines when creating
ContradictoryRequirements and UndeclaredRehearsal findings instead of
requirementIndex + 1 or a hardcoded line 1. Preserve the existing declaration
text and replay behavior.
```

</details>

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:153464edd90bf979ff4e3edf -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
