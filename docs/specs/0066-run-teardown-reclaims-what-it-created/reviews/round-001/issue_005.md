---
source: coderabbit
pr: "116"
round: 1
round_created_at: "2026-08-05T05:03:08Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0066-implementation
head_sha: 76774c31db4c686cac700bc83af0bcc68521b6c6
file: internal/cli/cli_test.go
line: 9296
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Wi9X3,comment:PRRC_kwDOS0qyts7dnSbW
review_hash: 7e6e0989276db5799af38b6cfb3979081fe97fd3d50824f6883f75da722db77b
duplicate_of: ""
source_review_id: "4861144519"
source_review_submitted_at: "2026-08-05T05:01:44Z"
---

# Issue 005: _ Maintainability & Code Quality_ _ Trivial_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _⚡ Quick win_

**Add a case where the classification returns both `Current` and `Releasable`.**

No table case sets `Current` and `Releasable` together. That combination is the only way to reach `internal/cli/cli.go` lines 2163-2177, which disregards the current failed-cycle branch with the proof text `current QA Report %q classifies the failed-cycle set`. It also has a preserve fallback when `CurrentReport` is empty. Both sub-branches are untested. Add one case that returns `Current`, `CurrentReport`, and a non-empty `Releasable`, and assert the `proof=current QA Report` diagnostic.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/cli_test.go` around lines 9254 - 9296, Add a table case in the
classification tests whose callback returns Current, CurrentReport, and a
non-empty Releasable set, then assert the successful result includes a
diagnostic containing “proof=current QA Report”. Ensure the fixture exercises
the current-report path in the CLI classification handling while preserving the
existing cases.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:ad1a1b9deec7efb85cf65dda -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes:
  - The requested `Current` plus `Releasable` disregard diagnostic encodes the behavior corrected by issue 001: current evidence must remain actionable and must not be reported as disregarded.
  - The revised end-to-end regression covers the valid boundary instead: older proven-superseded branches receive disregard diagnostics, while the current evidence branch has no disregard diagnostic and follows normal integration.
  - Focused evidence: `TestBranchIntegrityPreflightWatchDisregardsOnlySupersededFailedQACycles` passed as part of the 1,247-test affected-package run.
