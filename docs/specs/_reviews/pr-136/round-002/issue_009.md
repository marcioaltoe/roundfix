---
source: coderabbit
pr: "136"
round: 2
round_created_at: "2026-08-06T19:47:02Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: 2a1d4725a703a2baf5514952d9986761bc2a234d
file: docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md
line: 125
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XE5Yl,comment:PRRC_kwDOS0qyts7eY0jU
review_hash: 5f0a220fbeaca662dd78903ea7bc757c92722f7d4263a56d5c0b3870c5065482
duplicate_of: ""
source_review_id: "4877313912"
source_review_submitted_at: "2026-08-06T18:14:54Z"
---

# Issue 009: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Resolve the retention-shape proposal in this tech spec.**

The Data Models section defers a payload side-table, and the Build Order defers whether payload shedding or a second retention shape will exist until later measurement. This leaves the retention contract open even though the PRD makes retention behavior and payload-loss handling core requirements. Choose the supported retention behavior, or explicitly remove it from scope. Record rejected alternatives and the data-loss consequence.

As per coding guidelines, a technical specification must present one proposal and document remaining risks.    


Also applies to: 142-153, 198-201

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md` around lines
119 - 125, Resolve the retention contract in the tech spec by selecting one
supported retention behavior or explicitly removing retention behavior from
scope. Update the Data Models, Build Order, and related sections to make payload
shedding versus a second retention shape unambiguous, document rejected
alternatives and the resulting payload-loss consequence, and record any
remaining risks.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:5ebcd59532555bef5ceafc2b -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The PRD and TechSpec now select one retention contract: the existing
  terminal-only, age-based Journal Retention behavior remains unchanged, and
  this Spec only makes its eligibility query cheap. A payload side table,
  payload shedding, and a second window are explicitly out of scope and
  rejected. The existing consequence remains explicit: GC deletes a terminal
  Run's complete journal after the configured boundary; retained Runs lose no
  payload before it.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test -count=1 -parallel=1 ./internal/speccheck -run
  '^TestCheckCorpusBudget$'` passed; `rtk env
  GOCACHE=/Users/marcio/dev/roundfix/.gocache go run -buildvcs=false
  ./cmd/roundfix spec check` passed with no findings for Spec 0081.
- Daemon Verification: `make verify` not run; Daemon-owned.
