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
line: 100
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6Vb8ya,comment:PRRC_kwDOS0qyts7b_18P
review_hash: 99c2003bc21899f4da445acd9d48f922d44328d5ccc60f5cc467afb89bb375b3
duplicate_of: ""
source_review_id: "4829144282"
source_review_submitted_at: "2026-07-31T14:01:53Z"
---

# Issue 007: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Resolve the byte-preservation and link-rewrite conflict.**

Line 100 requires preserving every byte of the adopted source. Lines 118-120 require rewriting every link to the old path across the repository. If an adopted source contains such a link, both requirements cannot hold.

Define whether adopted-source links are excluded from rewriting, or permit controlled link-only edits and update the byte-preservation rule.







Also applies to: 118-120

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/write-prd/SKILL.md around lines 98 - 100, Clarify the rules
in the adoption and link-rewrite guidance so they no longer conflict: explicitly
state whether links inside the moved source are excluded from repository-wide
rewriting, or revise the byte-preservation requirement to allow only controlled
link updates. Update the instructions around the “Move” step and the
link-rewrite steps to define the chosen behavior while preserving the required
git move semantics.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:3270fd28fb3969cfe878e973 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: "preserve every byte" and "rewrite every link" were contradictory
  when an adopted source itself contained an old-path Markdown destination.
- Fix: the move remains byte-preserving. After the move, step 7 may change only
  Markdown link destinations, including destinations inside the adopted source;
  observations and all other source content remain immutable.
- Focused evidence: the focused Skill contract test, full `./skills` tests,
  `skills-sync-check`, and `git diff --check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
