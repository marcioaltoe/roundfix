---
source: coderabbit
pr: "56"
round: 1
round_created_at: "2026-07-31T14:02:35Z"
status: invalid
terminal_reason: "Spec 0060 explicitly assigns self-containment to the archive-spec Skill and ships no Go code; internal/spec.Archive retains its separate Tasks-and-QA command contract."
head_repository: marcioaltoe/roundfix
head_branch: ma/0060-spec-owned-reference-lifecycle
head_sha: 05752e266533235d41a554f01dd42584bd24d46d
file: .agents/skills/archive-spec/SKILL.md
line: 3
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vb8yE,comment:PRRC_kwDOS0qyts7b_171
review_hash: 1b0a04cda75dba87f580261f7f5e650ca90113462b55e44050013bfe3aec5e26
duplicate_of: ""
source_review_id: "4829144282"
source_review_submitted_at: "2026-07-31T14:01:53Z"
---

# Issue 001: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Enforce self-containment at the archive mutation boundary.**

This skill requires indexed-path, source-path, and forbidden-link checks. The supplied `internal/spec/archive.go:30-75` performs task and QA checks before `os.Rename`, and `internal/cli/archive.go:28-71` calls that service directly. The archive command can therefore move a Spec without the new self-containment precondition.

Add the validation to `spec.Archive` before stamping or moving, and cover it with tests.







Also applies to: 15-17, 39-73

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md at line 3, Update spec.Archive to run
indexed-path, source-path, and forbidden-link self-containment validation before
stamping archive metadata or calling os.Rename; ensure validation failure leaves
the spec unmoved and unstamped. Add tests covering these validation failures and
successful archiving, including the direct internal/cli/archive.go call path.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8e0400a0d6b51cd7ac816ad5 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: The finding conflates two archive paths. The accepted TechSpec states
  that Spec 0060 ships no Go code, deliberately keeps the new gate in the
  authorial Skill, and treats a future Go gate as follow-up only if agents skip
  the instruction (`_techspec.md`, Executive Summary and Risks &
  Considerations). `.agents/skills/archive-spec/SKILL.md` performs its own
  preflight, frontmatter stamp, and `git mv`; it does not call
  `internal/spec.Archive`. The public Archive Command retains its separately
  documented Tasks-and-QA-only contract in `CONTEXT.md` and the Roundfix Skill.
  Changing that command here would contradict the accepted implementation
  boundary rather than repair this PR's workflow.
- Focused evidence: complete reads of
  `docs/specs/0060-spec-owned-reference-lifecycle/_prd.md`, `_techspec.md`,
  `task_01.md`, `.agents/skills/archive-spec/SKILL.md`,
  `internal/spec/archive.go`, and `internal/cli/archive.go`.
- Daemon Verification: `make verify` not run; Daemon-owned.
