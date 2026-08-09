---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: failed
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/baseline/plan.go
line: 566
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XeKLE,comment:PRRC_kwDOS0qyts7e9kjB
review_hash: 422d5f07fd2286aab05ab67b5d344ec7619ec552e916a096348eb99fce5f37bd
duplicate_of: ""
source_review_id: "4888818931"
source_review_submitted_at: "2026-08-08T12:40:11Z"
---

# Issue 024: _ Functional Correctness_ _ Major_ _ Heavy lift_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _🏗️ Heavy lift_

**Do not block a restorable external skill before the skills stage.**

`ResolveProfileAlignment` runs before this managed-refresh path. A missing required external skill therefore returns an action-required result before `baseline update` can reach its post-apply restoration stage. The QA report records this as F-002 for `capability.context7`.

Route restorable external-skill drift to the skills stage, then verify profile alignment before reporting success. Keep preflight blocks for divergences that the update cannot repair. Add an end-to-end test with a missing required external skill.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/baseline/plan.go` around lines 548 - 556, Update the managed-refresh
planning flow around ResolveProfileAlignment so missing required external skills
that post-apply restoration can repair are deferred to the skills stage rather
than returned as preflight action-required results. Preserve preflight blocking
for divergences the update cannot repair, then re-run profile alignment after
restoration before reporting success. Add an end-to-end test covering a missing
required external skill such as capability.context7.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:df0adc52d8cbb6e73052e939 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `failed`
- Notes: This requires a multi-step architectural change across `plan.go`, `profile_alignment.go`, and `baseline_update.go`. The current flow: `BuildPlan` → `ResolveProfileAlignment` (line 346) returns blocking divergences for missing required external skills → `BuildPlan` returns `actionOutcome` with `Plan == nil` → `baseline_update.go` returns early at line 216 before reaching the skills stage at line 270. To fix: (1) `ResolveProfileAlignment` or its caller must distinguish restorable external-skill divergences from non-restorable ones; (2) `BuildPlan` must defer restorable divergences instead of blocking at line 356; (3) `baseline_update.go` must re-run `ResolveProfileAlignment` after the skills stage completes at line 293; (4) an end-to-end test must exercise `capability.context7` restoration. The change touches the core planning contract and requires careful design of the "restorable vs non-restorable" boundary. This is correctly tagged as a "heavy lift" and warrants a separate focused PR.
