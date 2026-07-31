---
source: coderabbit
pr: "56"
round: 2
round_created_at: "2026-07-31T14:33:59Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: .agents/skills/write-prd/SKILL.md
line: 103
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vce6R,comment:PRRC_kwDOS0qyts7cAnC7
review_hash: 372d92f31da9f77adaee63b6959a20fa36e6ea89024db2fe28e0d66b0c1f7766
duplicate_of: ""
source_review_id: "4829407112"
source_review_submitted_at: "2026-07-31T14:33:10Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Preflight basename collisions before changing adoption state.**

`git mv <source> .../<basename>` creates a flat destination in `references/`. Two relied-upon sources can have the same basename. The first move then succeeds and the second fails. Because the workflow flips finding status before the move, the Spec can be left partially adopted with `status: done` while a source remains at its old path. Detect duplicate basenames and existing destination files before the first status update or move.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-prd/SKILL.md around lines 99 - 103, The workflow
instructions must require preflight validation before any adoption status update
or source move: detect duplicate source basenames targeting the same references
directory and existing destination files, then abort before changing state when
either collision exists. Update the “Move” procedure while preserving its
one-move, byte-preserving behavior.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:f3fadca1e07b5d0b40a8ac94 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The workflow changed a finding's adoption state before `git mv` and
  had no whole-inventory check for duplicate basenames or occupied
  destinations, so a later move could fail after partial mutation.

## Resolution

- Added an explicit whole-inventory preflight before index creation, finding
  status changes, or moves. It rejects duplicate source basenames and any
  occupied destination path, including a symbolic link.
- Kept the one-move, no-copy, byte-preserving contract and added ordered Skill
  contract assertions for the new workflow sequence.

## Focused evidence

- `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./skills -run TestSpecReferenceLifecycleSkillContracts -count=1` — exit 0.
- `rtk make skills-sync-check` — exit 0.
- `rtk make baseline-digests` — exit 0; the repeat run reported
  `ok:true, changed:false`.
- `rtk git diff --check` — exit 0.
- `make verify` — not run; Daemon-owned authoritative Verification.
