---
source: coderabbit
pr: "136"
round: 2
round_created_at: "2026-08-06T19:47:02Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: 2a1d4725a703a2baf5514952d9986761bc2a234d
file: docs/specs/0080-cheap-detectors-run-before-the-gate/_techspec.md
line: 106
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XE5Ya,comment:PRRC_kwDOS0qyts7eY0jE
review_hash: 860725487ac91e819a76fa9399286916e8d844b294860a5d7aa1d2e1912b2c45
duplicate_of: ""
source_review_id: "4877313912"
source_review_submitted_at: "2026-08-06T18:14:54Z"
---

# Issue 006: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Define report materialization for mechanical results and skips.**

The active PRD requires every absent-artifact skip to record the detector and missing artifact. This interface has no skip field, and it states that the stage never edits a report. When the stage withholds the Agent Session, no component is defined to write the blocked rows or recorded skips. Add the result and report-writing contract, then replace “silent skip” with the observable recorded-skip behavior.

   


Also applies to: 133-138, 182-185

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0080-cheap-detectors-run-before-the-gate/_techspec.md` around
lines 90 - 106, The MechanicalResult contract must include recorded skips with
the detector and missing artifact, and define how mechanical findings, carried
rows, blocking rows, and skips are materialized into the report when the Agent
Session is withheld. Update the MechanicalStage specification and related
sections to replace silent skips with observable recorded-skip behavior,
including the report-writing responsibility and contract.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fa8a853fa048e73ed0e46378 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: `MechanicalResult` now carries typed Findings, Carried rows, Blocked
  rows, and Skips whose fields name the detector and missing artifact. The
  TechSpec assigns Markdown ownership to a Daemon report materializer, defines
  both blocking and non-blocking materialization, and replaces silent skips
  with recorded skips. A blocking result closes a mechanical-stage-only fail
  report with zero pending rows before the Agent Session is withheld.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  run -buildvcs=false ./cmd/roundfix spec check` passed with no findings for
  Spec 0080; `rtk git diff --check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
