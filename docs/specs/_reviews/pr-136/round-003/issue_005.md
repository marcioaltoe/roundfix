---
source: coderabbit
pr: "136"
round: 3
round_created_at: "2026-08-06T20:20:19Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: fba018672a8f31a3a4f59e6afd21d2c03c6a220f
file: docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md
line: 126
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XGgXH,comment:PRRC_kwDOS0qyts7ebLMI
review_hash: 8de69510db09c271435f1a2a6f1bbde387ce796daac011cec2656c5f180d88e0
duplicate_of: ""
source_review_id: "4877969817"
source_review_submitted_at: "2026-08-06T20:19:25Z"
---

# Issue 005: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Define the terminal-event path after `CloseJournal`.**

Lines 119–123 close the shared journal before terminal settlement, but terminal settlement must persist `daemon.outcome`. If `CloseJournal` rejects later appends, the specified sequence cannot persist that event. Define whether terminal settlement inserts `daemon.outcome` directly through `withWriteTx`, or move `CloseJournal` after the terminal transaction. Specify error propagation for that path.

As per coding guidelines, integration points and failure modes must be explicit.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md` around lines
119 - 126, Clarify the terminal-event sequence in the Run lifecycle
specification: ensure daemon.outcome is persisted through an explicitly defined
terminal transaction despite CloseJournal rejecting later appends, or move
CloseJournal after that transaction. State how errors from the terminal
transaction and CloseJournal propagate through the existing Run-failure path and
prevent requested-outcome settlement.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:0410ee1f22765a01f4ee111d -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: The TechSpec now names the terminal operation explicitly:
  `CompleteRun` bypasses the closed JournalWriter and uses `withWriteTx` to
  atomically commit terminal state, Active Run lock release, cursor allocation,
  and `daemon.outcome`. It also defines Close retryability, pre-commit rollback,
  ambiguous-commit reconciliation, Failed fallback behavior, and the rule that
  prevents a conflicting fallback before the requested outcome is disproved.
- Focused evidence: the post-change `rtk rg` contract probe found the direct
  `CompleteRun` interface, closed-writer bypass, requested-outcome refusal, and
  post-Close coverage. `rtk env
  GOCACHE=/Users/marcio/dev/roundfix/.gocache go test -count=1 -parallel=1
  ./internal/speccheck -run '^TestCheckCorpusBudget$'` passed. `rtk env
  GOCACHE=/Users/marcio/dev/roundfix/.gocache go run -buildvcs=false
  ./cmd/roundfix spec check` exited 0 with no findings for Specs 0080 and 0081;
  it separately reported the pre-existing missing Task Graph, reference-index,
  and Vocabulary Contract skips. `rtk git diff --check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
