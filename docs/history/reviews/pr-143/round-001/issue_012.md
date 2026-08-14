---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/cli/baseline_human_test.go
line: 113
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XbA2u,comment:PRRC_kwDOS0qyts7e5EBM
review_hash: 1ec43efe81ddb349840d4ba102c48061f7fe48c72c53aca7ada9da4e99c2a7f0
duplicate_of: ""
terminal_reason: 'Verification failed: command "make verify" exited with exit status 2; diagnostics: /Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20260809T011807Z_77466ac6469b46dc/verification/batch-001-attempt-2.log'
source_review_id: "4887493512"
source_review_submitted_at: "2026-08-08T00:23:24Z"
---


# Issue 012: _ Maintainability & Code Quality_ _ Trivial_ _ Low value_

## Review Comment

_📐 Maintainability & Code Quality_ | _🔵 Trivial_ | _💤 Low value_

**Name the prompts that this answer script drives.**

`changeAnswers` is a 17-element script of `"1"`, `"2"`, and `"3"` values with no mapping to prompts. A reader cannot tell which answer selects "Change Baseline Profile" or which trailing `"2"` declines the recomputed Plan. Add a short comment that names each group, or build the slice from named constants.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/baseline_human_test.go` around lines 109 - 113, Document the
answer sequence used by changeAnswers by adding a concise comment identifying
which prompt each initial answer, repeated "1" group, and trailing "2"
correspond to, including the "Change Baseline Profile" selection and recomputed
Plan decline. Keep the existing 17-answer ordering and values unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:a54943a087ba82c7e3b6086d -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Added concise comment documenting the changeAnswers sequence in baseline_human_test.go explaining which prompt each answer group corresponds to. `rtk go build ./...` passes.
