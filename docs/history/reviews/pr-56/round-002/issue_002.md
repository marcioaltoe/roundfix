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
line: 176
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vce6M,comment:PRRC_kwDOS0qyts7cAnC1
review_hash: f742d56806dff27fa71b214a5c1332224f5d491583c301f9399ae543cca8744e
duplicate_of: ""
source_review_id: "4829407112"
source_review_submitted_at: "2026-07-31T14:33:10Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Enforce the source-to-path basename invariant.**

`write-prd` moves each source to `references/<basename>` and records the current relative path. This validator only requires `path` to be an existing safe basename. A row can claim `docs/findings/a.md` with `path` set to `b.md`; if `a.md` is absent and `b.md` exists, the check passes even though the indexed file is not the moved source. Require `path` to equal the basename of `source`, or define a rename rule in both skills.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md around lines 161 - 176, The
self-containment validator in the archive-spec skill must enforce that each
indexed path names the moved source file. In the loop reading source and path,
compare path with the basename of source and reject mismatches before validating
the target file; keep the existing safe-basename, existence, and source-absence
checks unchanged.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8d59bcca1674e891b76a7677 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The row validator checked that `path` was a safe existing basename
  but did not prove it was the basename of the row's never-updated `source`.

## Resolution

- The archive loop now derives the source basename with shell parameter
  expansion and rejects a mismatched `path` before resolving or validating the
  target file.
- Added ordered Skill contract assertions for the basename comparison and its
  placement before target-file validation.

## Focused evidence

- `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./skills -run TestSpecReferenceLifecycleSkillContracts -count=1` — exit 0.
- `rtk make skills-sync-check` — exit 0.
- `rtk make baseline-digests` — exit 0; the repeat run reported
  `ok:true, changed:false`.
- `rtk git diff --check` — exit 0.
- `make verify` — not run; Daemon-owned authoritative Verification.
