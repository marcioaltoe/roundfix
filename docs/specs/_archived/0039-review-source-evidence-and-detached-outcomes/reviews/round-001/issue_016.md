---
source: coderabbit
pr: "40"
round: 1
round_created_at: "2026-07-28T18:11:57Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/review-evidence-and-verification-capacity
head_sha: 2f9f75357d07ae37bd94ddc54579d8e8b8d6ef0b
file: docs/specs/0039-review-source-evidence-and-detached-outcomes/qa/evidence/2026-07-28-rerun-02/make-verify.txt
line: 30
severity: minor
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6USSo7,comment:PRRC_kwDOS0qyts7aUjcA
review_hash: 54bf2536d94efebd48a3e989fdbdf2ff2e0007c826526c345acd05fe68aaca6e
duplicate_of: ""
source_review_id: "4793775214"
source_review_submitted_at: "2026-07-28T04:41:35Z"
---

# Issue 016: _ Data Integrity & Integration_ _ Minor_ _ Quick win_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟡 Minor_ | _⚡ Quick win_

**Reconcile the dirty build identity before treating this as exact-commit evidence.**

The binary records `BuildCommit=f3649fe-dirty`, while the surrounding QA artifacts describe it as built from exact commit `f3649fef83f720c51da883444038b76e9def6296`. Either rebuild from a clean worktree, or explicitly document that only QA evidence files were dirty and verify that no product paths differed; otherwise the exercised binary is not reproducibly identified.

As per coding guidelines, verification claims must prove the exact artifact being claimed.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In
`@docs/specs/0039-review-source-evidence-and-detached-outcomes/qa/evidence/2026-07-28-rerun-02/make-verify.txt`
around lines 29 - 30, Reconcile the production binary’s dirty BuildCommit
identity with the claimed exact commit f3649fef83f720c51da883444038b76e9def6296:
rebuild from a clean worktree, or document that only QA evidence files were
dirty and verify no product paths differed. Update the verification evidence so
the exercised artifact is reproducibly identified before claiming exact-commit
coverage.
```

</details>

<!-- fingerprinting:phantom:triton:luna -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:43ece88d06f83335f1f6e4ca -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Reconciled the historical QA evidence with the exact verification commit and parent, and recorded that the commit changed only Spec-local QA evidence/report paths. `git rev-list --parents -n 1 8971bda...` and `git diff-tree --name-only` confirmed the recorded identities and path scope.
