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
line: 105
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vce6W,comment:PRRC_kwDOS0qyts7cAnDB
review_hash: bd40826dedf5e71055551b82966e9ab1e52a8765a3fbf2f4e2435a37cee8f88b
duplicate_of: ""
source_review_id: "4829407112"
source_review_submitted_at: "2026-07-31T14:33:10Z"
---

# Issue 004: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Create the reference index before the first move, or mark adoption as incomplete.**

`archive-spec` skips self-containment when `references/_index.md` is absent to preserve legacy Specs. This workflow creates the index only after moving sources. If a move succeeds and the index step fails, the current Spec can enter the legacy archive path without provenance or link validation. Write the complete index before any `git mv`, or add an explicit adoption marker that archive validation rejects while incomplete.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-prd/SKILL.md around lines 104 - 105, The write-prd
workflow must establish reference-index state before moving adopted sources.
Update the steps around the “Index” instruction so
docs/specs/<slug>/references/_index.md is fully written and validated before any
git mv operation, ensuring archive-spec cannot treat a partially adopted Spec as
legacy; alternatively, add an explicit incomplete-adoption marker that archive
validation rejects until the index and moves are complete.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:eb7097d1895b983a02c76bef -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The workflow created `_index.md` after moving sources, while
  `archive-spec` treats only a genuinely absent index as legacy. An interrupted
  move sequence could therefore leave adopted state without the marker that
  activates self-containment validation.

## Resolution

- Reordered adoption so the complete `_index.md` is written and structurally
  validated before the first finding status update or source move.
- Documented that this early index intentionally makes an interrupted adoption
  non-legacy, so `archive-spec` rejects it until every indexed move completes.
- Added ordered Skill contract assertions that pin preflight, index, status
  update, move, and final rewrite/gate order.

## Focused evidence

- `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./skills -run TestSpecReferenceLifecycleSkillContracts -count=1` — exit 0.
- `rtk make skills-sync-check` — exit 0.
- `rtk make baseline-digests` — exit 0; the repeat run reported
  `ok:true, changed:false`.
- `rtk git diff --check` — exit 0.
- `make verify` — not run; Daemon-owned authoritative Verification.
