---
source: coderabbit
pr: "37"
round: 1
round_created_at: "2026-07-27T01:53:02Z"
status: invalid
head_repository: marcioaltoe/roundfix
head_branch: ma/0036-doctor-skill-readiness
head_sha: 9a6b7f9433b9779afe75f38d833b780ceb2555ed
file: .agents/skills/golang-testing/SKILL.md
line: 170
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6T6oNU,comment:PRRC_kwDOS0qyts7Zy0Nt
review_hash: f30131620eec8f46f02219ff8447dc907cfb0f2c25f5e199eb665de1a3c30c2a
duplicate_of: ""
source_review_id: "4783144632"
source_review_submitted_at: "2026-07-27T01:52:02Z"
---

# Issue 002: _ Maintainability & Code Quality_ _ Major_ _ Quick win_

## Review Comment

_📐 Maintainability & Code Quality_ | _🟠 Major_ | _⚡ Quick win_

**Remove or reframe the Testify guidance for this repository.**

The repository-wide `**/*_test.go` rule requires only the standard-library `testing` package and explicitly forbids introducing testify. This newly added section presents Testify as supported, so it conflicts with the authoritative testing policy. Rewrite the examples using `t.Errorf`/`t.Fatalf`, or clearly mark this guidance as non-applicable here.

<details>
<summary>🧰 Tools</summary>

<details>
<summary>🪛 LanguageTool</summary>

[style] ~145-~145: Consider shortening this phrase to just ‘whether’, unless you mean ‘regardless of whether’.
Context: ...y hiding which case broke. This happens whether or not the subtest calls `t.Parallel()`.  ```g...

(WHETHER)

</details>

</details>

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In @.agents/skills/golang-testing/SKILL.md around lines 143 - 170, Reframe the
“Assert Scope Leaking into Subtests” section to comply with the repository’s
standard-library-only testing policy. Remove Testify references and rewrite both
examples and the verification guidance using subtest-local testing.T methods
such as t.Errorf or t.Fatalf, while preserving the requirement that failures are
attributed to the correct subtest.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:814a1484574623e06a31446b -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `INVALID`
- Notes: This finding treats a generic, externally managed skill snapshot as
  repository-authored guidance. The repository's narrower Go rules in
  `docs/agents/go.md` and `docs/agents/specific-repository.md` already require
  stdlib `testing` and forbid Testify, so those rules govern Roundfix code even
  when an upstream skill describes Testify usage for other repositories.
  Rewriting the snapshot locally would violate the upstream-skill ownership
  and immutable-hash contracts in `docs/agents/skill-dispatch.md` and
  `skills-lock.json`. No repository change is appropriate.

## Verification

- Inspected the complete upstream `golang-testing` snapshot and the
  repository's narrower Go testing rules.
- Inspected the `golang-testing` lock entry and commit `672eb5c`, which refresh
  the external skill artifact and its hash as one immutable snapshot.
