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
line: 67
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vb8yL,comment:PRRC_kwDOS0qyts7b_18A
review_hash: 828b942d086e00fffa4c62a23528a4d7f2ca3fc7028d9de0a949dd3041242079
duplicate_of: ""
source_review_id: "4829144282"
source_review_submitted_at: "2026-07-31T14:01:53Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Reject index paths that escape `references/`.**

The documented command tests `docs/specs/<slug>/references/<path>` without requiring a normalized path under `references/`. A row such as `../../other-spec/file.md` can pass when that file exists outside this Spec.

Canonicalize each path. Reject absolute paths, `..` traversal, and symlinks that escape `references/`.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md around lines 63 - 67, Update the
index-validation instructions in the archive-spec workflow to canonicalize each
row path before testing it. Reject absolute paths, parent-directory traversal,
and symlinks resolving outside docs/specs/<slug>/references/, then run the
existence checks only for validated paths while preserving the existing
offending-path/source output and adoption guidance.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:dc970abb4005bdf95fced303 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The old `references/<path>` concatenation accepted path separators and
  followed symlinks, so it did not prove that the indexed file was the adopted
  basename inside `references/`.
- Fix: the archive gate now requires one basename relative to `_index.md`,
  rejects empty, absolute, dot, traversal, slash, and backslash forms, rejects
  a symlinked `references/` directory or indexed file, and checks the resulting
  entry is a regular file before checking source absence.
- Focused evidence: the extracted documented command passed a valid index
  (exit 0) and rejected `../../outside.md` (exit 1) plus an escaping symlink
  (exit 1), each naming the offending path; `rtk bash -n` accepted the command
  block; the focused Skill contract test and `skills-sync-check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
