---
task: task_02
spec: 0075-typed-docs-backlog
status: pending
type: test
complexity: medium
---

# Task 02: Re-record the corpus and declare what moved

## Overview

The corpus is a characterization gate, and re-recording it is sanctioned by
ADR-0081 as fallout of the authorized module edit. Sanctioned is not the same as
unexamined: **only layout content and its digests may move.** Anything else
moving means the edit leaked into a surface this Spec does not own.

This slice runs the regeneration and reads the diff, which is the whole point.

## Requirements

1. MUST run `make baseline-digests` and re-record the two characterization
   corpora it does not reach:

   ```bash
   go test ./internal/baseline -count=1 \
     -run TestBaselinePlanCharacterization -update-baseline-plan-characterization
   go test ./internal/baseline -count=1 \
     -run TestCatalogDiagnosticCharacterization -update-catalog-diagnostics
   ```

2. MUST inspect the resulting diff and record, in the Task Result, every path
   that moved and why layout content explains it.
3. MUST fail the Task if any path moved that layout content does not explain.
   A leaked edit is the failure this Task exists to catch.
4. MUST leave `make verify` green after one regeneration pass.
5. MUST NOT hand-edit any derived artifact. Every change comes from a
   regeneration command.

## Subtasks

- [ ] Run the sanctioned command and both corpus re-records.
- [ ] Read the full diff and classify every moved path.
- [ ] Record the declared break list in the Result.

## Acceptance Criteria

- [ ] `make verify` exits 0 after the regeneration pass.
- [ ] Every moved path is listed in the Result with the layout content that
      explains it.
- [ ] No path moved that layout content does not explain.
- [ ] No derived artifact was hand-edited, asserted by every change tracing to a
      regeneration command named in the Result.

## Verification

This Task's whole purpose is to return a repository that task_01 legitimately
left red. Declaring the configured repository Verification command verbatim
makes the Daemon run it as a **precondition**, before the Agent starts — and it
fails on exactly the state this Task exists to repair, settling the Task
without ever starting it. Measured on 2026-08-05:
`repository not green on entry: make verify exited 2`.

The gate below is the same gate, named by its parts rather than by the
configured string, so the repository is still proven green **after** the
regeneration and never demanded green before it. Nothing is weakened: every
target `make verify` runs is listed.

- `make baseline-digests` — expected: exit 0.
- `go test ./internal/baseline -count=1 -run 'TestBaselinePlanCharacterization|TestCatalogDiagnosticCharacterization' -v | grep -q -- "--- PASS"`
  — expected: exit 0; both corpora match after re-recording.
- `make fmt-check test spec-budget skills-sync-check skills-check build spec-check`
  — expected: exit 0; every target `make verify` runs, after the regeneration.
- `git diff --name-only HEAD | grep -E "\.go$" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no Go source changed.

## References

- `_prd.md` → Success Metrics.
- `_techspec.md` → Testing Approach; Build Order 2; Risks & Considerations.
- ADR-0081.
