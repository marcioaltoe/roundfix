---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: pending
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: internal/spec/task.go
line: 190
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v6Z,comment:PRRC_kwDOS0qyts7eEK7x
review_hash: a2a58c17601c38b337dd5880c8f60ee2ed35336cdb21b67e530958c8cffffde0
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:29Z"
---

# Issue 023: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Stop extraction at every heading.**

A `###` heading inside `## Requirements` does not end `inSection`. Its heading and content can become continuation text for the prior numbered requirement. This changes the declarations that the coherence checks evaluate.

Terminate the section on any Markdown heading after the matching `##` heading. Add a parser test with a nested heading followed by content.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/task.go` around lines 171 - 190, Update taskSectionLines to
stop collecting lines when any Markdown heading is encountered after the
matching ## heading, including ### and deeper headings, while preserving the
existing section matching behavior. Add a parser test covering a nested heading
followed by content and verify that neither the nested heading nor its content
is returned.
```

</details>

<!-- fingerprinting:phantom:poseidon:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:c923e9f774df7c61c34e24cd -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `UNREVIEWED`
- Notes:
