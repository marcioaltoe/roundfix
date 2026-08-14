---
source: coderabbit
pr: "56"
round: 1
round_created_at: "2026-07-31T14:02:35Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: .agents/skills/write-prd/SKILL.md
line: 83
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vb8yV,comment:PRRC_kwDOS0qyts7b_18K
review_hash: df7a9da07a4badd2b184fe2f97958e62da6bd08dca9f3660fac06863bd1d73e4
duplicate_of: ""
source_review_id: "4829144282"
source_review_submitted_at: "2026-07-31T14:01:53Z"
---

# Issue 006: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Create the Spec reference directory before the adoption move.**

Lines 81-83 place adoption before the Spec folder is written. Lines 98-100 move into `docs/specs/<slug>/references/`. Line 129 creates the Spec folder later and does not explicitly create `references/`. The first `git mv` can therefore fail because its destination does not exist.

Create the numbered Spec and `references/` directory before step 5, or move adoption after directory creation.







Also applies to: 98-100

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-prd/SKILL.md around lines 81 - 83, Reorder the
write-prd workflow so the numbered Spec directory and its references
subdirectory are created before the adoption step that moves source documents
into them. Update the steps around the adoption instructions and the later
Spec-folder creation so the first git mv always targets an existing directory,
while preserving the documented adoption sequence.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:a01d3fde80d3390d1167c972 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The prior order could run `git mv` before either the numbered Spec
  directory or `references/` existed.
- Fix: write-prd step 4 now resolves the numbered slug and runs `mkdir -p
  docs/specs/<slug>/references` before adoption. Step 5 writes `_prd.md` into
  that prepared folder instead of creating the folder after the move.
- Focused evidence: the initially red
  `TestSpecReferenceLifecycleSkillContracts/PRD_adoption` case passed after the
  edit; full `./skills` tests and `skills-sync-check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
