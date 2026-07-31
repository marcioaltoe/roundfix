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
line: 37
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vb8yJ,comment:PRRC_kwDOS0qyts7b_179
review_hash: ef4446cb0f044e06d239f0d5939174715eff39d5f29d47eee2447ad4dea5db9f
duplicate_of: ""
source_review_id: "4829144282"
source_review_submitted_at: "2026-07-31T14:01:53Z"
---

# Issue 002: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Implement the documented `qa_override` path.**

The skill allows failed or missing QA when the user explicitly says `archive anyway` and records `qa_override: true`. The supplied `internal/spec/archive.go:30-75` always rejects missing or non-pass QA, and `internal/cli/archive.go:28-71` does not pass an override field. This path cannot succeed.

Add an explicit override to the request and CLI. Limit it to QA. Keep self-containment non-overridable. Otherwise, remove the override instructions.







Also applies to: 70-73

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md around lines 30 - 37, Implement the
documented QA-only override across the archive flow: add an explicit override
field to the archive request, expose and propagate an “archive anyway” option
from the CLI, and update the validation in the archive logic to allow missing or
non-passing QA only when that override is set. Preserve the existing
self-containment rejection regardless of the override, and ensure stamped
frontmatter records qa_override: true when used; otherwise remove the override
documentation.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:f9a5b8594d07ae278f91ea6d -->

---

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Use the repository QA verifier instead of a line grep.**

The exact `grep -n '^verdict: pass$'` check does not validate YAML or `rows_blocked_finding`. `internal/spec/qa.go:156-201` rejects malformed reports and rejects `verdict: pass` when `rows_blocked_finding` is nonzero. The documented preflight can disagree with `QAVerdict`.

Invoke the repository verifier or document equivalent parsing and blocked-row checks.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/archive-spec/SKILL.md around lines 30 - 37, Update the QA
preflight instructions in the archive procedure to invoke the repository QA
verifier, or explicitly document equivalent YAML parsing and
rows_blocked_finding validation, instead of relying on grep for verdict: pass.
Ensure archive eligibility matches QAVerdict, including rejection of malformed
reports and pass reports with nonzero blocked rows, while preserving the
existing override behavior.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2095fbce7d461358b792a12a -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID` for the QA-verifier mismatch; `INVALID` for the proposed
  CLI override path.
- Notes: The override premise conflates the Skill's direct stamp-and-move
  procedure with the separate `roundfix archive` command. Spec 0060 explicitly
  ships no Go code, and the Skill already records an explicitly authorized
  manual QA override in the stamped PRD. No CLI flag or ArchiveRequest field is
  part of this Spec's accepted contract. The second finding is valid: a literal
  `grep '^verdict: pass$'` can accept malformed YAML or a pass report with
  `rows_blocked_finding > 0`.
- Fix: replaced the line grep with the repository `internal/spec.QAVerdict`
  semantics: filename-contract report selection, YAML parsing, supported
  verdict validation, typed non-negative blocked counts, and rejection of a
  finding-blocked pass. The existing QA-only override and non-overridable
  self-containment boundary remain unchanged.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test ./skills -run '^TestSpecReferenceLifecycleSkillContracts$'` passed;
  `rtk make GOCACHE=/Users/marcio/dev/roundfix/.gocache skills-sync-check`
  passed; full `./skills` package tests passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
