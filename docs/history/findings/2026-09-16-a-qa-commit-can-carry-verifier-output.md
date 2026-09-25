---
status: done
created_at: 2026-09-16
updated_at: 2026-09-25
absorbed_by: 0141-a-commit-that-carries-the-work-and-nothing-else
---

# QA gate — A QA commit can carry files the repository Verification wrote (2026-09-16)

The pre-PR review of Spec 0139's delivery raised this against the Daemon's new
repository Verification step. Reading the QA gate showed the exposure is older
than that step and does not reach this repository, so it is recorded rather than
repaired inside that Spec.

## 1. The QA commit stages every path that appeared after the step's snapshot

- Symptom / evidence:
  - The QA step takes its before-snapshot when it begins, and the QA Report
    commit later stages what is dirty relative to that snapshot
    (`internal/daemon/task_engine.go`).
  - Anything the repository Verification writes into the worktree after that
    snapshot is therefore a candidate for the QA Report commit. The stageable
    filter removes repository-external paths, symlink crossings and executables,
    not ordinary generated files.
  - This predates Spec 0139: the QA Agent ran the same command inside the same
    window. That Spec moved the run to the Daemon, which keeps the exposure
    rather than creating it.
  - This repository is not affected by its own command: `make verify` writes
    `/bin/` and `/.gocache/`, both ignored by Git.
- Root cause: the commit baseline is captured before a step that may write to the
  worktree, and nothing re-establishes it afterwards.
- Action / suggestion:
  - Route to a Spec that re-establishes the commit baseline after the repository
    Verification, so only the report, its evidence and the Agent's own writes
    reach the QA commit.
  - Keep the existing exclusions for external paths, symlink crossings and
    executables.
  - A repository whose Verification writes tracked or unignored files is the case
    that makes this visible.

## Addendum — 2026-09-25 — Delivered by Spec 0141

Revalidated during the 2026-09-25 triage: `commitQAReport` in
`internal/daemon/task_engine.go` subtracts the paths the repository Verification
wrote (`verificationWindowPaths`), delivered by commit `a383d4d6` (Spec 0141-a-commit-that-carries-the-work-and-nothing-else);
`TestQAReportCommitExcludesVerificationWrites` pins it. One edge remains: a path
written by both the verifier and the Agent is dropped; only the new-directory
case is tested.
