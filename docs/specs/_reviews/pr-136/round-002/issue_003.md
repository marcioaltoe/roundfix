---
source: coderabbit
pr: "136"
round: 2
round_created_at: "2026-08-06T19:47:02Z"
status: resolved
head_repository: marcioaltoe/roundfix
head_branch: ma/0079-one-door-for-fleet-knowledge
head_sha: 2a1d4725a703a2baf5514952d9986761bc2a234d
file: docs/adr/0098-run-events-append-in-batches-with-durability-per-batch.md
line: 12
severity: major
author: coderabbitai[bot]
source_ref: thread:PRRT_kwDOS0qyts6XE5YE,comment:PRRC_kwDOS0qyts7eY0im
review_hash: 53bffbac80482e6dfcb257a947f62139f50a193cff9f87c965f31d9379d7d22a
duplicate_of: ""
source_review_id: "4877313912"
source_review_submitted_at: "2026-08-06T18:14:54Z"
---

# Issue 003: _ Data Integrity & Integration_ _ Major_ _ Heavy lift_

## Review Comment

_🗄️ Data Integrity & Integration_ | _🟠 Major_ | _🏗️ Heavy lift_

**Define the batch contract before implementation.**

Lines [3]-[9] promise that a crash loses at most the batch in flight, but the ADR does not define the batch boundary. Specify count, time, shutdown, and error flush triggers. Specify event ordering, commit failure handling, and retry or idempotency behavior. Without these rules, the durability and replay guarantee cannot be tested, and a retry can duplicate or silently drop raw events.

As per coding guidelines, cross-boundary contracts must define batching, ordering, retries, and replay semantics.

<details>
<summary>🤖 Prompt for AI Agents</summary>

```
Verify each finding against current code. Fix only still-valid issues, skip the
rest with a brief reason, keep changes minimal, and validate.

In `@docs/adr/0098-run-events-append-in-batches-with-durability-per-batch.md`
around lines 3 - 12, Expand the ADR’s batch contract to define count-based,
time-based, shutdown, and error-triggered flush boundaries. Specify that events
retain input order, commit failures preserve the uncommitted batch for retry,
and retries use an idempotency strategy that prevents duplicates while ensuring
no raw events are silently dropped; document the resulting replay semantics and
crash-loss bound.
```

</details>

<!-- fingerprinting:phantom:triton:caracal -->

<!-- cr-indicator-types:potential_issue -->

<!-- cr-comment:v1:e947d950fe5a556847d1c6b2 -->

_Source: Coding guidelines_

<!-- This is an auto-generated comment by CodeRabbit -->

## Triage

- Decision: `VALID`
- Notes: ADR-0098 and Spec 0081 now define count, linger, immediate-event,
  error, shutdown, Agent-teardown, and terminal-settlement boundaries. They
  preserve publisher order, retain a failed batch, reconcile an ambiguous
  commit against its assigned cursor range, retry only the same cursors, and
  fail a partial or different match instead of duplicating or dropping raw
  events. The crash bound is the open Store batch's count or linger window.
- Focused evidence: `rtk env GOCACHE=/Users/marcio/dev/roundfix/.gocache go
  test -count=1 -parallel=1 ./internal/speccheck -run
  '^TestCheckCorpusBudget$'` passed; `rtk git diff --check` passed.
- Daemon Verification: `make verify` not run; Daemon-owned.
