---
status: done
created_at: 2026-07-17
updated_at: 2026-07-18
---

# Spec Runs — concurrent Task Verification exhausted the test environment (2026-07-17)

Resolved in planning by
[Spec 0042 — Verification Capacity and Daemon Task Settlement](../specs/0042-verification-capacity-and-daemon-task-settlement/_prd.md).

Roundfix implemented Vortex Spec `0006-uniformidade-uuidv7` in detached Run
`run_20260717T151124Z_1d2c757911fbf76a` with `worktree.concurrency: 2` and Worktree Bootstrap
`bun install`. The Run exposed a conflict between concurrent Agent work, project-mandated full
verification, and the documented Daemon-owned Verification contract. The implementation changes
passed their focused tests and typecheck, while concurrent full backend suites failed during test
environment setup.

## 1. Concurrent Task Worktrees ran the full repository gate at the same time

- **Symptom / evidence**: Task 02 and Task 04 ran `rtk make verify` concurrently from separate
  Task Worktrees. Both passed formatting, lint, OnionCry, and typecheck, then failed in
  `backend#test` with `Hook timed out in 30000ms` from integration-suite `beforeAll` hooks. The
  Task 02 log named Better Auth helpers, session middleware, and several Visio integrations. The
  Task 04 log named `erp-connection-resolution`, `platform-health`, `review-import-matching`, and
  `check-boundaries-cli`. The affected Task 04 tests passed when rerun in isolation. Sandboxed
  attempts also reported `listen EPERM` while SMTP tests tried to bind `127.0.0.1`.
- **Root cause**: not yet proven by a controlled sequential rerun. The evidence points to resource
  contention: `worktree.concurrency: 2` allowed two complete Vitest integration suites to create
  containers and local listeners at once, while their setup hooks retained a 30-second limit. The
  shipped Roundfix skill already warns that concurrent Tasks can run heavy Verification commands
  simultaneously and consume matching CPU and cache resources.
- **Action / suggestion**: keep Vortex's `make verify` unchanged and set its Project Config to
  `worktree.concurrency: 1` until Roundfix can serialize Verification separately. Confirm the
  hypothesis by settling the failed Tasks sequentially with the same Verification commands. Do
  not weaken the gate, remove suites, or increase timeouts merely to accommodate concurrent Runs.

## 2. One concurrency setting couples Agent work and Verification

- **Symptom / evidence**: Vortex benefits from two concurrent implementation Task Worktrees, but
  the same `worktree.concurrency` value also permits their full Verification commands to overlap.
  Project Config has no independent limit for Verification. The only current configuration that
  prevents the overlap is `worktree.concurrency: 1`, which serializes the entire Task lifecycle.
- **Root cause**: the Wave scheduler applies one concurrency limit to Task Worktrees. Verification
  has no repository-scoped semaphore or separate configured capacity.
- **Action / suggestion**: add a Verification capacity setting, for example
  `verification.concurrency: 1`, independent of `worktree.concurrency`. The Daemon can continue
  running two Agents concurrently while queueing heavy Verification per repository. The Run Event
  Stream and Live Run View must distinguish `waiting-for-verification` from an Agent that is still
  working or a Verification command that is running.

## 3. Agent-owned status prevented Daemon-owned Verification

- **Symptom / evidence**: the Run Event Stream contains Verification events only for Task 01.
  Tasks 02 and 04 ran `rtk make verify` inside their Agent Sessions, changed their task files to
  `status: failed`, and settled failed without a Daemon Verification attempt, Verification
  Feedback, or the documented repair attempt. Their dependent Tasks were then skipped even though
  focused tests and typecheck passed and the failing integration tests were unrelated or passed
  alone.
- **Root cause**: the Agent workflow can treat its own verification attempt and task-file status as
  authoritative before the Daemon reaches the Verification boundary. This conflicts with the
  Roundfix glossary and skill contract, where the Daemon runs Verification verbatim and owns the
  pass/fail decision.
- **Action / suggestion**: make the Daemon the sole writer of terminal Task status during an
  Implement Command. The Agent must hand back implementation-ready work without running the
  declared Task Verification commands or marking the Task `completed` or `failed`. The Daemon must
  run those commands through the serialized Verification queue, return Verification Feedback to
  the same Agent Session on attempt 1 failure, and settle the Task only after its final verdict.
  Project-level rules that require a full gate remain satisfied because the Daemon's Verification
  is part of Task completion.

## 4. A transient capacity failure causes dependency skips and a second Run

- **Symptom / evidence**: after Task 02 and Task 04 settled failed, Roundfix skipped dependent
  Tasks 03 and 08 while continuing independent Tasks. Recovery requires the Active Run to finish,
  sequential Settle Commands for the kept Task Worktrees, and another Implement Command for the
  skipped dependency chain.
- **Root cause**: a failed Task correctly blocks dependents, but the Run cannot distinguish a
  deterministic implementation failure from a Verification capacity failure produced by
  concurrent execution. Settle Command also refuses recovery while an Active Run owns the Spec.
- **Action / suggestion**: the serialized Verification boundary removes the primary failure mode.
  If Verification still fails before any assertion runs, classify the infrastructure condition
  separately and retry it under exclusive Verification capacity before settling the Task failed.
  Keep the current Worktree preservation and dependency blocking behavior when the exclusive retry
  also fails.

## What worked — keep

- Task Worktrees preserved every failed Task's implementation for inspection and Settle Command
  recovery.
- Focused tests and typecheck isolated the UUIDv7 behavior from the full-suite environment failure.
- The Console Log, Run Event Stream, task result sections, and retained RTK logs made the missing
  Daemon Verification and exact timeout locations observable.
- Independent Tasks continued after one dependency chain became blocked.
- `make verify` remained an explicit non-negotiable project gate; the failure was reported instead
  of being hidden through validation configuration changes.

## Recommended sequence

1. Set Vortex `worktree.concurrency` to `1` and keep Worktree Bootstrap as `bun install`.
2. Settle the retained failed Tasks sequentially and record whether unchanged Verification passes.
3. Specify independent Verification capacity and Daemon-only Task status ownership in a Roundfix
   implementation Spec.
4. Restore `worktree.concurrency: 2` only after the Daemon serializes heavy Verification.
