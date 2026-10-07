---
schema: roundfix/archive-record/v1
spec: 0098-a-hook-that-cannot-outrank-the-gate
title: A hook that cannot outrank the gate
status: archived
created: "2026-08-12"
archived: "2026-08-26"
disposition: pass
source: docs/history/specs/0098-a-hook-that-cannot-outrank-the-gate
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_08
qa_report: qa-report-2026-08-25-02.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "167"
delivery_commit: 86c3ca440f661d9223727513c614df10e33af015
---

# A hook that cannot outrank the gate

The Daemon runs the authoritative Verification and then commits. When the repository's commit hook refuses, the work is reverted and the Run ends failed — and the repair loop covers a Verification failure, not a hook failure. Three Runs died this way in one Spec, every time with work that was correct and already verified, left staged in a Task Worktree. Two gates where one cannot be satisfied nor recovered is a design conflict, and the invariant that resolves it — a commit hook may never be stricter than the authoritative Verification — is written nowhere. The command built for recovery makes it worse: it refuses work that is completed but uncommitted, because its contract assumes lost work is always failed.
