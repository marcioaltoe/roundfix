---
source: coderabbit
pr: "143"
round: 1
round_created_at: "2026-08-09T01:14:50Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/specs-0082-0083
head_sha: 49ad4407ca050772b13a812b9625cf45d940e19e
file: internal/agent/selection_capabilities.go
line: 549
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XiApL,comment:PRRC_kwDOS0qyts7fC8Re
review_hash: 9e1c94948fd97c670eb5ad9c5094c475a2067778eaeb8005959cd17f6e3fffeb
duplicate_of: ""
source_review_id: "4890173306"
source_review_submitted_at: "2026-08-09T00:28:49Z"
---

# Issue 045: _ Functional Correctness_ _ Major_ _ Quick win_

## Review Comment

_🎯 Functional Correctness_ | _🟠 Major_ | _⚡ Quick win_

**Retain the exact requested model-effort variant first.**

`bindsRequestedModel` retains generic canonical matches in advertised order. If more than 64 variants share one canonical model, the map can fill before `<model>[<reasoning_effort>]` is visited. `PlanSelectionAssignment` then rejects an advertised configured selection.

Prioritize the exact requested variant before generic canonical matches. Add a regression case with more than 64 variants and the requested variant at the end.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/agent/selection_capabilities.go` around lines 527 - 549, Update the
capability-retention logic around bindsRequestedModel so the exact requested
model-effort variant is retained before generic canonical matches, even when
more than maxRetainedCapabilityValues variants are advertised. Preserve the
existing retention behavior for currentValue and canonical matches, and add a
regression case covering over 64 variants with the requested variant last.
```

</details>

<!-- fingerprinting:phantom:medusa:tapir -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:40f492671ace2e10ff130b22 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `resolved`
- Notes: Confirmed real before fixing. `retainAdvertisedValues` already
  prioritised `currentValue`, the requested effort, and every
  `bindsRequestedModel` match over the generic fill, but it treated all
  canonical matches alike. When more than `maxRetainedCapabilityValues`
  advertised values bind to the requested canonical model, the map fills on
  sibling `<model>[<effort>]` variants and the exact requested value is dropped,
  so `PlanSelectionAssignment` rejects a selection the runtime does advertise.
  Fixed by an exact-match pass that retains `currentValue`, the requested
  effort, and the requested model by identity before canonical matching runs; it
  adds at most three entries and cannot itself exhaust the bound. Regression
  added as `TestRetentionKeepsExactRequestedModelAmongItsOwnVariants`, with 128
  variants of one canonical model and the exact value advertised last. Proven
  non-vacuous: with the exact-match pass disabled it fails reporting
  `retained 64 models`, and passes with it restored. The sibling test
  `TestRetentionKeepsRequestedModelPastTheBound` missed this because its fixture
  advertises distinct vendors, so only one value ever bound.
