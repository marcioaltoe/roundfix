---
source: coderabbit
pr: "130"
round: 1
round_created_at: "2026-08-06T03:34:01Z"
status: invalid
terminal_reason: "The requested Build Order backfill would rewrite archived Spec 0065."
head_repository: marcioaltoe/roundfix
head_branch: ma/0065-loop-order-and-verification-honesty
head_sha: 7d35358ba9f77ceeda86ec5c34d7c4485a7eb8f9
file: docs/specs/_archived/0065-loop-order-and-verification-honesty/_techspec.md
line: 202
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6W2v5p,comment:PRRC_kwDOS0qyts7eEK62
review_hash: 1c7a00dc3b8e625749e98528f9c5eb47440bfcecddf7bcf22aacb4e918ace1b0
duplicate_of: ""
source_review_id: "4870656161"
source_review_submitted_at: "2026-08-06T03:12:28Z"
---

# Issue 004: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Record the corrective Task and terminal QA Task in the Build Order.**

The task manifest adds `task_07` and terminal QA `task_06`, but this Build Order ends at Skill synchronisation. Task graph edges must derive from the TechSpec Build Order, and ADR-0091 requires the QA Task to be terminal. Add both steps with explicit dependencies. Use `(depends on: none)` instead of `(independent of 1)` for the independent step.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/_archived/0065-loop-order-and-verification-honesty/_techspec.md`
around lines 191 - 202, Update the Build Order section to use “(depends on:
none)” for Verification honesty, then add explicit steps for the corrective task
task_07 and terminal QA task_06 after Skill synchronisation, with dependencies
matching the task manifest and ensuring task_06 is terminal as required by
ADR-0091.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:20fd88c658b93010bd9d2995 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The archived Build Order predates the corrective `task_07` and
  terminal `task_06` entries now present in the manifest. The repository's
  archive contract preserves that historical artifact byte-for-byte; new
  requirements must be carried by a new Spec.
- Daemon Verification: `make verify` not run; Daemon-owned.
