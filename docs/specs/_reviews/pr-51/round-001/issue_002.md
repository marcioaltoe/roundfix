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
line: 349
severity: nitpick
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6VJIT5,comment:PRRC_kwDOS0qyts7bkbJr
review_hash: c636caf0eeffdf25db3bf9d09a08b10b7371915617bd3e6d62cf0cb780b92cd5
duplicate_of: ""
source_review_id: "4820128835"
source_review_submitted_at: "2026-07-30T14:56:10Z"
---

# Issue 002: _ Performance & Scalability_ _ Trivial_ _ Quick win_

## Review Comment

_🚀 Performance & Scalability_ | _🔵 Trivial_ | _⚡ Quick win_

**Two git subprocesses per commit on a preflight path.**

`QAReportOnlyBranch` is called from Branch Integrity Preflight for every pending Implement branch, and a branch with many commits spawns `2N` git processes before the first non-QA message aborts the loop. A single `git log --format=%H%x00%B` (or `--name-only -z`) pass over `targetHead..runHead` gives the same evidence in one invocation.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@internal/worktree/worktree.go` around lines 318 - 349, Update
QAReportOnlyBranch to gather each commit’s subject/message and changed paths in
a single git log invocation over targetHead..runHead, rather than calling
runner.Run for show and diff-tree inside the commit loop. Preserve the existing
QA message matching, non-empty path, qaDirs containment, and error-result
behavior while parsing the combined NUL-delimited output.
```

</details>

<!-- fingerprinting:phantom:medusa:komodo -->

<!-- cr-indicator-types:nitpick -->

<!-- cr-comment:v1:654d2d6720e3603d85a6c55d -->

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes:
  - The prior implementation ran `rev-list` once, then `show` and `diff-tree` for each commit; the real-Git regression failed before the fix with three calls for one QA commit and five calls for two QA commits.
  - `QAReportOnlyBranch` now reads every commit message and NUL-delimited changed path with one `git log --name-only -z` invocation. The parser preserves the exact QA message match, non-empty path requirement, active or archived QA-directory containment, empty-range behavior, and conservative false result for malformed evidence.
  - Focused evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/worktree -run 'Test(QAReportOnlyBranch|InspectTerminalRunBoundsSupersededReason|InspectTerminalRunClassifiesSupersededQAReport|ApplyTerminalRunSuperseded)' -count=1` passed.
  - Integration evidence: `rtk proxy env GOCACHE=/Users/marcio/dev/roundfix/.gocache go test ./internal/cli -run '^(TestRunReconcile|TestBranchIntegrityPreflight|TestBranchIntegrityIntegrationPlan)' -count=1` passed.
  - The Daemon owns the authoritative `make verify` run after this Agent turn.
