---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: docs/specs/0082-the-manifest-already-answered-that/_techspec.md
line: 54
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XeKK9,comment:PRRC_kwDOS0qyts7e9kiz
review_hash: 15697669edf149cab5e8b560cb672743603ea4c6f5ca4c4dc6a25f4a268df449
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4888818931"
source_review_submitted_at: "2026-08-08T12:40:11Z"
---


# Issue 020: _ Maintainability & Code Quality_ _ Major_ _ Heavy lift_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _🏗️ Heavy lift_

**Complete the active ADR obligations.**

The PRD marks ADR-0047, ADR-0066, ADR-0067, ADR-0070, and ADR-0081 as binding for this Spec. This Technical Spec does not classify them as applicable or not applicable.

Add each omitted ADR to this row. Map its concrete constraint to the design, or state why it does not apply. This keeps the Technical Spec consistent with the PRD before Tasks derive implementation work.

As per coding guidelines, `_techspec.md` must classify active ADR obligations as applicable or not applicable with reasons.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0082-the-manifest-already-answered-that/_techspec.md` around lines
41 - 54, Update the “Active ADR obligations” row in the Technical Spec to
include ADR-0047, ADR-0066, ADR-0067, ADR-0070, and ADR-0081. Classify each as
applicable with its concrete design constraint or explicitly not applicable with
a reason, while preserving the existing classifications and source reference.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:9ee703c8a06b10e54a7aa974 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added ADR-0047 (immutable fingerprinting), ADR-0066 (CLI exit codes), ADR-0067 (structured CLI output), ADR-0070 (CLI response completeness), and ADR-0081 (digested guidance pins) to the Active ADR obligations in _techspec.md with classifications.
