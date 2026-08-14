---
source: coderabbit
pr: "136"
round: 2
round_created_at: "2026-08-06T19:47:02Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: 2a1d4725a703a2baf5514952d9986761bc2a234d
file: docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md
line: 69
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XE5Ye,comment:PRRC_kwDOS0qyts7eY0jL
review_hash: 76f000b1a07828bbcfe4ccffb86e9289f96344f1afd9bc2a2f077537dd237e96
duplicate_of: ""
source_review_id: "4877313912"
source_review_submitted_at: "2026-08-06T18:14:54Z"
---

# Issue 007: _ Stability & Availability_ _ Major_ _ Heavy lift_

## Review Comment

_🩺 Stability & Availability_ | _🟠 Major_ | _🏗️ Heavy lift_

**Choose one writer-concurrency design.**

The spec claims that parallel Runs will stop producing `SQLITE_BUSY`, but it does not define the mechanism. Batching alone does not prove zero busy errors or define cursor allocation and error propagation across writers. Select one writer, connection, and transaction discipline. Record rejected alternatives before implementation.

As per coding guidelines, a technical specification must present one proposal and document trade-offs and rejected alternatives.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/specs/0081-a-journal-cheap-to-write-and-keep/_techspec.md` around lines
52 - 69, Update the technical specification to select one explicit
writer-concurrency design covering the writer count, SQLite connection usage,
transaction discipline, cursor allocation, and error propagation needed to
prevent SQLITE_BUSY across parallel Runs. Add a concise trade-off section
documenting the chosen approach and rejected alternatives before implementation,
while preserving the existing batching, retention, read, cockpit, and
measurement scope.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:fcaa4951dff43b41414c2305 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: Spec 0081 now selects one writer discipline: each operational process
  opens one single-connection writer Store, and every Roundfix write
  transaction acquires the same machine-wide advisory lock before `BeginTx`
  through commit or rollback. Cursor allocation stays in the transaction and
  errors propagate to the caller. The trade-off and the rejected batching-only,
  per-process-goroutine, and per-Run-database alternatives are explicit.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  run -buildvcs=false ./cmd/roundfix spec check` passed with no findings for
  Spec 0081; `rtk git diff --check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
