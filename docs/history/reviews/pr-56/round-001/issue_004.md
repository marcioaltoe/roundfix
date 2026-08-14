---
source: coderabbit
pr: "56"
round: 1
round_created_at: "2026-07-31T14:02:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: .agents/skills/archive-spec/SKILL.md
line: 68
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vb8yQ,comment:PRRC_kwDOS0qyts7b_18E
review_hash: abe920ae659f1ce4236ed8df2e897339c5fd6ae52a416935ac525e7c7712a67f
duplicate_of: ""
source_review_id: "4829144282"
source_review_submitted_at: "2026-07-31T14:01:53Z"
---

# Issue 004: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Validate index ownership and row uniqueness.**

The archive check verifies path existence and source absence, but it does not verify that each `owner` equals the current four-digit Spec or that each source appears once. A stale index can pass while violating the sole-owner contract.

Parse `_index.md` and enforce ownership, uniqueness, valid types, and path containment before archiving.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md around lines 63 - 68, Update the
archive validation instructions around the _index.md inspection to parse every
row before archiving and enforce that each owner matches the current four-digit
Spec, each source appears exactly once, each row has a valid type, and
referenced paths remain within the expected archive directory. Preserve the
existing path-existence and source-absence checks, reporting offending rows and
retaining command output as evidence.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:9b2a078e29af315d1e9c2533 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The former prose checked only current-path existence and original-path
  absence. It did not parse the fixed table or validate its ownership and
  uniqueness invariants.
- Fix: the documented gate now parses the fixed five-column table, requires
  the current four-digit Spec owner, accepts only `inbox` and `finding`, checks
  the source path matches its type, validates the adoption date, rejects
  duplicate `source` and `path` values, and emits the normalized rows as retained
  evidence before the existing path/source checks complete.
- Focused evidence: the extracted command rejected an invalid `report` type and
  owner `9999` for Spec `0001` (exit 1) and rejected a duplicated source (exit
  1), while the valid fixture exited 0. The focused Skill contract test and
  `skills-sync-check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
