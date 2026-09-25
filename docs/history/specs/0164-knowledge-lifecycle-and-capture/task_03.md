---
task: task_03
spec: 0164-knowledge-lifecycle-and-capture
status: completed
type: backend
complexity: medium
---

# Task 03: A terminal record closes without a Spec

## Overview

`detectArchiveLicenses` in `internal/speccheck/citations.go` demands an `absorbed_by` for every archived Finding, though the layout guide allows `closure_reason` and `closure_evidence`. `ClassifyBacklogEntry` in `internal/spec/retirement.go` retires only `declined`, and `detectBacklogPromotion` in `internal/speccheck/backlog.go` ignores a terminal Backlog Entry left in `docs/backlog/`.

This is an authorized tooling Task for `internal/speccheck/backlog.go` and `internal/speccheck/backlog_test.go`, bounded by [_authorization.md](_authorization.md). Stop before any other governed mutation.

## Requirements

1. MUST make `parseFindingFrontmatter` read `closure_reason` and `closure_evidence`, and `detectArchiveLicenses` accept an archived Finding with no `absorbed_by` when both are non-empty.
2. MUST keep reporting `SC-ARCHIVE-LICENSE` for a present, unresolvable `absorbed_by` even beside closure fields, and for a Finding with neither license, with a Fix naming both routes.
3. MUST make `ClassifyBacklogEntry` keep retiring `declined` and also retire `done`, `deprecated`, `superseded`, `closed` and `cancelled` when the entry names a `spec` or a non-empty `reason`; `open`, `promoted` and unknown statuses never retire.
4. MUST make `detectBacklogPromotion` report, under `SC-BACKLOG-UNMOVED`, a terminal Backlog Entry still in `docs/backlog/`, with a Fix naming `docs/history/backlog/` and, when it has no `spec`, the missing `reason`.
5. MUST add no detector code and keep `TestCheckCorpusGolden` passing unchanged.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] An archived Finding with closure fields and no absorber passes the archive check.
- [ ] Closure fields beside an invalid `absorbed_by` still fail.
- [ ] Each terminal Backlog status retires with a `spec` or `reason`; `open` and an unknown status do not.
- [ ] A terminal Backlog Entry in `docs/backlog/` is reported.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/backlog.go`
- interface: `internal/spec/retirement.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestArchivedFindingClosesWithReasonAndEvidence|TestArchivedFindingClosureCannotHideAnInvalidAbsorber|TestTerminalBacklogEntryLeftActiveIsUnmoved|TestClassifyBacklogEntryRetiresEveryClosedStatus)$" ./internal/speccheck ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestArchivedFindingClosesWithReasonAndEvidence TestArchivedFindingClosureCannotHideAnInvalidAbsorber TestTerminalBacklogEntryLeftActiveIsUnmoved TestClassifyBacklogEntryRetiresEveryClosedStatus; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Terminal dispositions

## Result

Implemented the two documented terminal-disposition routes without adding a detector code. Archived Findings now close through either a resolvable `absorbed_by` or non-empty `closure_reason` and `closure_evidence`; a present invalid absorber remains an `SC-ARCHIVE-LICENSE` error. Backlog retirement now recognizes every declared terminal status with its required `spec` or `reason`, and `SC-BACKLOG-UNMOVED` reports terminal entries that remain under `docs/backlog/` with the History Root destination and missing-reason guidance when applicable.

Acceptance evidence:

- An archived Finding with closure fields and no absorber passes: `TestArchivedFindingClosesWithReasonAndEvidence` passed in `./internal/speccheck`.
- Closure fields beside an invalid `absorbed_by` still fail: `TestArchivedFindingClosureCannotHideAnInvalidAbsorber` passed, including the invalid-absorber diagnostic and the missing-license Fix naming both routes.
- Every terminal Backlog status retires with a `spec` or non-empty `reason`, while `open`, `promoted`, unknown, blank-evidence, and missing-evidence cases do not: `TestClassifyBacklogEntryRetiresEveryClosedStatus` passed in `./internal/spec`.
- A terminal Backlog Entry left in `docs/backlog/` is reported: `TestTerminalBacklogEntryLeftActiveIsUnmoved` passed and asserted both `docs/history/backlog/` and the missing `reason` in the Fix.

Focused checks:

- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 ./internal/speccheck ./internal/spec` — passed.
- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -tags docscontract -run '^TestCheckCorpusGolden$' ./internal/docscontract` — passed with the existing golden unchanged.
- `make verify-incremental` — passed after rerunning with the network access required by GitHub-backed integration checks; vet, the repository Go suite, skill synchronization checks, skill validation, and the build all exited successfully.

The Task's declared `## Verification` command was not run; Daemon Verification remains authoritative. No follow-up work was discovered.

## Carry-forward provenance

- Source Run: `run_20260924T223414Z_54e6f85ae98b6550`
- Source commit: `64dee5744bd809e5aae0e53dc64e0a4eb07a3340`
