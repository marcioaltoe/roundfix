---
source: coderabbit
pr: "27"
round: 2
round_created_at: "2026-07-15T16:07:28Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/dogfood-findings-remediation
head_sha: 81ff89aabce7ee1f748504ad6e4bcf1ac52ea200
file: internal/cli/implement.go
line: 625
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6RJ95b,comment:PRRC_kwDOS0qyts7V5taP
review_hash: 93868452c680a03552887ccf4a8f474a28bfb28581c71152a95d9fd421631a02
duplicate_of: ""
source_review_id: "4705718496"
source_review_submitted_at: "2026-07-15T15:38:26Z"
---

# Issue 009: _ Functional Correctness_ _ Minor_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟡 Minor_ | _⚡ Quick win_

**Normalize reasons to one stdout line.**

`strings.TrimSpace` preserves embedded newlines in `TaskOutcome.Reason`, so a daemon error can inject additional report lines. The documented CLI contract requires a single `reason:` line.

<details>
<summary>Proposed fix</summary>

```diff
-		reason := strings.TrimSpace(outcome.Reason)
+		reason := strings.Join(strings.Fields(outcome.Reason), " ")
```
</details>

Add a multiline-reason case to the renderer test.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/cli/implement.go` around lines 624 - 625, Normalize the reason
before rendering it in the report block keyed by reasons[task.ID], replacing
embedded newlines with spacing so each TaskOutcome.Reason produces exactly one
“reason:” stdout line while preserving its content. Add a renderer test covering
a multiline reason and asserting the output remains a single line.
```

</details>

<!-- fingerprinting:phantom:poseidon:terra -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:1449d2adedcb7038474fe031 -->

_Source: Coding guidelines_

<!-- This is an auto-generated reply by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Implement task outcome reasons are normalized with `strings.Fields` so each reason renders on one stdout line; `make verify` passed.
