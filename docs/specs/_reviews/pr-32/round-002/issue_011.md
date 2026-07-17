---
source: coderabbit
pr: "32"
round: 2
round_created_at: "2026-07-17T13:23:47Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/release-plan-and-agent-selection-profiles
head_sha: d7ab1933ac9fdcf0c94d73e2f417d99d38e43fe7
file: internal/spec/spec.go
line: 239
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Ru5tC,comment:PRRC_kwDOS0qyts7Wt95e
review_hash: b757d574c9e5745d59e94504af9e4af51642a860ca77540b9403aee2bd5c25b9
duplicate_of: ""
source_review_id: "4721765481"
source_review_submitted_at: "2026-07-17T10:25:31Z"
---

# Issue 011: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Do not let an empty or malformed projection table disable validation.**

`Load` treats `len(projections) == 0` as “no table,” while the parser silently skips malformed rows. A table whose IDs are all mistyped therefore bypasses the required projection checks. Duplicate Task IDs also overwrite earlier rows.

Recognize the canonical table explicitly, return whether it exists, and reject malformed or duplicate rows before validating every graph node.






As per coding guidelines, “Do not use workarounds; fix the root cause instead.” 


Also applies to: 299-372

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/spec/spec.go` around lines 222 - 239, Update Load and the
projection-table parsing flow so table presence is tracked independently of the
number of parsed rows. Identify the canonical table explicitly, reject malformed
rows and duplicate Task IDs during parsing, and propagate those errors before
validating graph nodes. Run projection validation whenever the table exists,
including when its parsed row set is empty.
```

</details>

<!-- fingerprinting:phantom:medusa:sol -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:1e72cf32efc90ab3ee8fad26 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Fixed Task Type projection parsing to track canonical table presence, reject malformed and duplicate rows, and validate graph nodes even for empty tables. Evidence: `GOCACHE=/private/tmp/roundfix-go-build rtk go test ./internal/agent ./internal/cli ./internal/config ./internal/daemon ./internal/releaseplan ./internal/spec ./internal/store ./internal/tui` passed.
