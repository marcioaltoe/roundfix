---
source: coderabbit
pr: "56"
round: 2
round_created_at: "2026-07-31T14:33:59Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: .agents/skills/archive-spec/SKILL.md
line: 84
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vce6G,comment:PRRC_kwDOS0qyts7cAnCu
review_hash: 2d0422839edf22988516782e0b42e0a109fee6d19672aebde68dbff5972e51a0
duplicate_of: ""
source_review_id: "4829407112"
source_review_submitted_at: "2026-07-31T14:33:10Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Reject a symlinked `references/_index.md`.**

The preflight rejects a symlinked `references/` directory and symlinked indexed files, but it does not reject the index file itself. `awk` follows a symlinked index and can validate provenance stored outside the Spec. Require the index to be a regular, non-symlinked file before parsing it.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md around lines 74 - 84, Update the
preflight checks around the `index` variable to reject a symlinked or
non-regular `references/_index.md` before invoking `awk`; require the index path
to identify a regular file without following symlinks, while preserving the
existing failure behavior and parsing flow for valid files.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:40f8effd5bc8d77c0f4945be -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The validator rejected a symbolic-link `references/` directory but
  invoked `awk` without first rejecting a symbolic-link or non-regular
  `references/_index.md`.

## Resolution

- Added a pre-`awk` guard that rejects `_index.md` when it is a symbolic link
  or is not a regular file. The migration boundary now also treats a broken
  index symlink as invalid indexed state rather than a legacy Spec.
- Added an ordered Skill contract assertion proving the index guard appears
  before `awk` parsing.

## Focused evidence

- `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./skills -run TestSpecReferenceLifecycleSkillContracts -count=1` — exit 0.
- `rtk make skills-sync-check` — exit 0.
- `rtk make baseline-digests` — exit 0; the repeat run reported
  `ok:true, changed:false`.
- `rtk git diff --check` — exit 0.
- `make verify` — not run; Daemon-owned authoritative Verification.
