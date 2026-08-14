---
source: coderabbit
pr: "51"
round: 1
round_created_at: "2026-07-30T16:50:38Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/spec-0053-qa-gate-reachability
head_sha: 4443081a669ed7731c37bfddb2615804338503a1
file: internal/worktree/worktree.go
line: 541
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6VJIUA,comment:PRRC_kwDOS0qyts7bkbJ1
review_hash: 6b800e5e73f011c4f16e362569aba0512b9bef42a9e7a553c11827c9d0c4549a
duplicate_of: ""
source_review_id: "4820128835"
source_review_submitted_at: "2026-07-30T14:56:10Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Bound the superseded reason string.**

Every other reason is a fixed constant, but this one embeds a full report path. With an archived path plus a long Spec slug the reason exceeds the 160-character bounded-reason contract asserted in `assertTerminalRunReconciliation` (worktree_test.go Line 1916) and persisted into the reconciliation payload. Truncate or bound the composed reason.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree.go` around lines 531 - 541, The superseded
reconciliation reason assigned in the supersedingQAReport branch can exceed the
160-character contract because it includes the full report path. Bound the
composed value before assigning result.Reason, preserving the identifying prefix
and applying the existing reason-length convention if available; keep
result.State and result.SupersedingReport unchanged.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:8c6615442df9bb35ad54fd20 -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - A real-Git regression using an archived QA path and long Spec slug reproduced a superseded reason longer than the existing 160-byte reconciliation contract.
  - Superseded reasons now preserve the identifying prefix, replace embedded line breaks, truncate on a UTF-8 rune boundary to at most 160 bytes, and append a truncation marker. `SupersedingReport` retains the complete path.
  - The regression also proves `ApplyTerminalRun` persists the bounded reason before cleanup.
  - Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/worktree ./internal/daemon -count=1` passed.
  - Integration evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/cli -run '^(TestRunReconcile|TestBranchIntegrityPreflight|TestBranchIntegrityIntegrationPlan)' -count=1` passed.
  - The Daemon owns the authoritative `make verify` run after this Agent turn.
