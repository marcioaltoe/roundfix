---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: docs/agents/docs-layout.md
line: 25
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XeHgI,comment:PRRC_kwDOS0qyts7e9gyP
review_hash: 99a72ee4446c083d53676a338dfbea3eb750cde65c1c3569b0bf9ddc3a4d91df
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4888802346"
source_review_submitted_at: "2026-08-08T12:29:17Z"
---


# Issue 016: _ Data Integrity & Integration_ _ Major_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _⚡ Quick win_

**Use the status that matches the adoption outcome.**

Line 25 requires `done` before moving every finding into a Spec. The same file defines `partial` for findings covered only in part and `deferred` for findings that are not implemented. Require `done` only for full adoption. Use `partial` or `deferred` for the other outcomes so the workflow does not mark incomplete findings as complete.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/agents/docs-layout.md` at line 25, Update the findings workflow
documentation near the docs/findings/ entry to require done only when a finding
is fully adopted by a Spec; require partial for findings addressed in part and
deferred for findings not implemented, while preserving the existing status
recording and move-to-references workflow.
```

</details>

<!-- fingerprinting:phantom:poseidon:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:2fa50e7e4b43578186a05519 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Updated findings workflow in docs-layout.md to require done only for full adoption, partial for findings addressed in part, deferred for not implemented.
