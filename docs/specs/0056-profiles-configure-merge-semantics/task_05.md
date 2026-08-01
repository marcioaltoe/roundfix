---
task: task_05
spec: 0056-profiles-configure-merge-semantics
status: pending
type: backend
complexity: medium
---

# Task 05: Make a refusal distinguishable by exit code

## Overview

A declined confirmation in a non-interactive context writes nothing and exits
zero, so automation reads a refusal as a successful write. This Task separates
the two by exit code alone, which is the only signal a script reliably has, and
evolves the machine result additively so an Agent can tell a refusal from an
applied change without parsing prose.

## Requirements

1. MUST exit non-zero when a write is declined, or when confirmation cannot be
   obtained non-interactively and explicit consent was not given.
2. MUST use the exit code the repository's documented contract already assigns
   to an unresolved or failed operation, and MUST NOT mint a new code.
3. MUST keep validation failures — invalid flags, a category named in both the
   fragment and a removal, a failed proof — on the existing validation exit
   code, distinct from a refusal.
4. MUST keep an applied write and an already-satisfied no-op exiting zero, and
   MUST keep `--dry-run` exiting zero without writing.
5. MUST print the refusal on the diagnostics stream, leaving the requested
   output stream for requested output.
6. MUST evolve the machine result additively, so existing fields keep their
   meaning and a refusal is machine-detectable.
7. MUST leave the merge, summary, and proof behavior from earlier slices
   unchanged.

## Subtasks

- [ ] Exit non-zero on a declined or unconfirmable write.
- [ ] Keep validation failures on their own distinct code.
- [ ] Add the refusal marker to the machine result additively.
- [ ] Route the refusal message to the diagnostics stream.

## Acceptance Criteria

- [ ] A declined confirmation exits non-zero and writes nothing.
- [ ] A non-interactive invocation without explicit consent exits with the same
      refusal code and writes nothing.
- [ ] A validation failure exits with the validation code, not the refusal code.
- [ ] An applied write, an already-satisfied no-op, and `--dry-run` each exit
      zero.
- [ ] The refusal message appears on the diagnostics stream, not the output
      stream.
- [ ] The machine result marks the refusal, and every previously present field
      keeps its meaning and type.
- [ ] `git status --porcelain` shows no path outside `internal/cli/` and this
      task file.

## Context

- instruction: `docs/agents/cli.md`
- interface: `internal/cli/profiles_configure.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/cli -run TestProfilesConfigureExitCodes -count=1` —
  expected: exit 0; refusal, validation failure, applied, no-op, and dry-run
  each assert their own code.
- `go test ./internal/cli -run TestProfilesConfigureChangeSummary -count=1` —
  expected: exit 0; the summary from task 04 is unchanged.
- `go test ./internal/config -run TestProfilesConfigWriterCharacterization -count=1`
  — expected: exit 0.
- `go test ./internal/config ./internal/cli -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 4; Core Features 5; Success Metrics (automation
  distinguishes by exit code alone).
- `_techspec.md` → API Contracts: exit codes; Build Order 5.
