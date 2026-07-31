---
source: coderabbit
pr: "56"
round: 3
round_created_at: "2026-07-31T14:59:09Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: .agents/skills/write-prd/SKILL.md
line: 99
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vc6Tn,comment:PRRC_kwDOS0qyts7cBOhg
review_hash: 888896524efbf3adecba60f921823e3cec2104eb87bba99d68be51bcb0968f77
duplicate_of: ""
source_review_id: "4829633138"
source_review_submitted_at: "2026-07-31T14:58:04Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Reserve `_index.md` during preflight.**

A source named `docs/findings/_index.md` or `docs/_inbox/_index.md` passes preflight because `references/_index.md` does not exist yet. Step 5 then creates the index at that path, and Step 7 cannot move the source. Reject this basename before writing the index or include the reserved index path in collision checks.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-prd/SKILL.md around lines 95 - 99, The preflight step
(Step 4) validates destination paths and rejects collisions but does not account
for `_index.md` as a reserved basename that Step 5 will create at
`references/_index.md`. Update the preflight validation to either reject sources
with basename `_index.md` outright before checking destination collisions, or
explicitly include the reserved index path `references/_index.md` in the
collision detection even though it may not exist yet at preflight time. This
ensures Step 7 can safely move all non-rejected sources without encountering
path conflicts created by Step 5.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:dccfb705b0b323ce91d7250e -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Step 4 checked only duplicate basenames and destinations that already
  existed. A source basename `_index.md` therefore passed preflight before step
  5 created the same reserved destination, leaving step 7 unable to move that
  source.

## Resolution

- Added `_index.md` as an explicit reserved source basename in the whole-inventory
  preflight, before index creation, finding status changes, or source moves.
- Added a focused Skill contract assertion so the reserved-basename rule cannot
  disappear while the broader collision wording remains.

## Focused evidence

- Before the fix,
  `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./skills -run TestSpecReferenceLifecycleSkillContracts -count=1`
  exited 1 and reported that the reserved `_index.md` contract was missing. The
  same command exited 0 after the fix.
- `rtk make skills-sync-check` — exit 0.
- `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./skills -count=1`
  — exit 0.
- `rtk make baseline-digests` — exit 0; the required repeat run reported
  `ok:true, changed:false`.
- `make verify` was not run; authoritative Verification is Daemon-owned.
