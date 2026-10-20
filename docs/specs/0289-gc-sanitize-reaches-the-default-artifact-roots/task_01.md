---
task: task_01
spec: 0289-gc-sanitize-reaches-the-default-artifact-roots
status: pending
type: backend
complexity: medium
---

# Task 01: One Artifact Root resolver serves the storage report and Run Retention

## Overview

The storage report and Run Retention each carry their own code to turn a Run's
recorded Artifact Directory into the root the Run wrote to, and sanitation is
about to need a third (`_prd.md` → Sources 1). This Task makes the storage
report's resolver the one exported resolver, keeps the storage report's output,
and moves Run Retention onto it. It is verifiable on its own on temporary
Roundfix Homes.

## Requirements

1. MUST implement API Contract 1 of `_techspec.md` in
   `internal/store/journal.go`: export `ResolveRunArtifactRoot`,
   `ErrUnsafeRunID` and `ErrArtifactRootUnresolved` as `_techspec.md` →
   Interfaces sketches them, add the recorded value to `ArtifactRootRun` as
   `ArtifactDir`, and remove `storageArtifactRoot`, whose rules the exported
   function keeps except for the `os.Lstat` probe.
2. MUST keep the storage report's output, as API Contract 2 states: on
   `ErrArtifactRootUnresolved` it measures the cleaned recorded root, or `.`
   for a blank one, and an unsafe Run ID still fails it with `measure Run
   artifacts: unsafe Run ID "<id>"`. Its sections, columns, groups, totals and
   reconciliation do not change.
3. MUST implement API Contract 3 in `internal/cli/run_retention.go`:
   `runRetentionArtifact` calls `store.ResolveRunArtifactRoot` and carries no
   resolution of its own, so neither `config.ResolveArtifactDirectory` nor
   `config.DefaultArtifactDirectoryForPath` is called in that file. Both
   resolver errors keep the Run as `artifact root unproven` before any further
   filesystem read; every later step stays as it is.
4. MUST state the hazard in the code's behavior: the resolver reads the Run
   Database row and the recorded repositories' Git metadata and returns a path
   to its caller; it never calls `ValidateArtifactDirectory` and creates,
   writes or removes nothing in Roundfix Home, a repository or a root.
5. MUST declare these breaks, each covered by a test below: a relative or `~/`
   recorded root is found by Run Retention under the Run's Git root or the
   home directory instead of the working directory; a Run whose root cannot be
   resolved is kept instead of failing the sweep; the storage report keeps the
   key's root when the key's `runs/<run-id>` is a symbolic link whose target
   is missing. Every existing test keeps its name and passes unchanged,
   including `TestStorageReportMeasuresTheArtifactRootARunWroteTo`,
   `TestRunStorageReportListsTheDefaultRootOfARunWithoutOne`,
   `TestGCRemovesRunsPastRunRetention` and
   `TestGCKeepsARunUnderAnUnprovenArtifactRoot`.
6. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows, using temporary homes and repositories, and never
   reading `~/.roundfix`:
   - in `internal/store/journal_test.go`,
     `TestResolveRunArtifactRootIsTheRootTheRunWroteTo`, with the subtests
     `_techspec.md` → Testing Approach lists: blank root, blank root of a
     removed linked worktree, relative root, `~/` root, absolute root, a
     directory only under the checkout's default, a symbolic link at the key's
     `runs/<run-id>`, an unsafe Run ID (`errors.Is` with `ErrUnsafeRunID`),
     blank Git and repository roots (`errors.Is` with
     `ErrArtifactRootUnresolved`), and the home and repository file trees
     identical before and after every call;
   - in `internal/cli/gc_run_retention_test.go`,
     `TestRunRetentionReadsTheRootTheSharedResolverReturns`, through `gc`: a
     terminal Run past Run Retention whose relative recorded root resolves
     under its Git root, outside Roundfix Home, with a directory there is kept
     with `artifact root unproven 1`; one whose resolved relative root holds no
     directory leaves the Run Database; a Run recording no root whose
     directory sits under the repository's default root leaves with that
     directory.
7. MUST document the change in `docs/user-guide/commands/gc.md`, in its Run
   Retention section, with the sentence "Run Retention finds each Run's
   artifact directory under the Artifact Root the Run wrote to, resolved the
   way `roundfix storage report` resolves it." followed by the three
   resolution rules (blank, relative, absolute).
8. MUST re-record the Behavior Surface Record, because this Task edits a user
   guide page: `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   The Agent runs it itself before handing back; no record is hand-edited.

## Subtasks

- [ ] Export the resolver and its errors from the storage report's code.
- [ ] Keep the storage report's fallbacks through the typed errors.
- [ ] Move Run Retention onto the resolver and drop its copy.
- [ ] Add the resolver and Run Retention tests.
- [ ] Add the guide sentence and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] One exported function resolves a Run's Artifact Root, and the storage report and Run Retention call it.
- [ ] The storage report's Artifact Roots lines are unchanged.
- [ ] Run Retention keeps and removes the same Runs for blank and absolute roots and finds a relative root under the Git root.
- [ ] The resolver creates nothing.

## Consumers

- `store.ResolveRunArtifactRoot` and `store.ArtifactRootRun` — consumers: `./internal/store`, `./internal/cli`

## Context

- instruction: `docs/adr/0301-gc-sanitation-classifies-the-artifact-root-each-run-wrote-to.md`
- instruction: `docs/adr/0255-run-retention-removes-terminal-runs-whole-and-compacts-incrementally.md`
- instruction: `docs/adr/0296-the-storage-report-reconciles-every-page-of-one-storage-snapshot.md`
- instruction: `internal/config/config.go`
- instruction: `internal/cli/gc.go`
- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- interface: `internal/cli/run_retention.go`
- interface: `internal/cli/gc_run_retention_test.go`
- interface: `docs/user-guide/commands/gc.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestResolveRunArtifactRootIsTheRootTheRunWroteTo|TestStorageReportMeasuresTheArtifactRootARunWroteTo|TestStorageReportReconcilesMeasuredTotalsWithoutMutation|TestDiscoverArtifactRootsReturnsEveryRecordedRepository)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestResolveRunArtifactRootIsTheRootTheRunWroteTo TestStorageReportMeasuresTheArtifactRootARunWroteTo TestStorageReportReconcilesMeasuredTotalsWithoutMutation; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestRunRetentionReadsTheRootTheSharedResolverReturns|TestGCRemovesRunsPastRunRetention|TestGCKeepsARunUnderAnUnprovenArtifactRoot|TestGCDryRunReportsRunRetention|TestRunStorageReportListsTheDefaultRootOfARunWithoutOne|TestRunStorageReportOutsideGitRepository)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunRetentionReadsTheRootTheSharedResolverReturns TestGCRemovesRunsPastRunRetention TestGCKeepsARunUnderAnUnprovenArtifactRoot TestRunStorageReportListsTheDefaultRootOfARunWithoutOne; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -q 'func ResolveRunArtifactRoot(' internal/store/journal.go || { echo 'missing exported resolver' >&2; exit 1; }; ! grep -q 'func storageArtifactRoot(' internal/store/journal.go || { echo 'storage report keeps its own resolver' >&2; exit 1; }; ! grep -q 'DefaultArtifactDirectoryForPath' internal/cli/run_retention.go || { echo 'Run Retention keeps its own resolution' >&2; exit 1; }; ! grep -q 'ResolveArtifactDirectory' internal/cli/run_retention.go || { echo 'Run Retention keeps its own resolution' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/gc.md | grep -qF -- 'Run Retention finds each Run' || { echo 'missing guide sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/gc.md | grep -qF -- 'under the Artifact Root the Run wrote to' || { echo 'missing guide sentence' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/store ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the resolver test does not exist, so the command fails at its first pass check. Afterwards the resolver and Run Retention tests pass with the storage report and retention tests they must not break, the storage report and Run Retention carry no resolution of their own, the guide states the rule, both packages build for Windows, and the Behavior Surface Record is current.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Core Feature 4; Success Metric 3; Prerequisites
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 8; Invariants 1-2, 5; Interfaces; Data Models; Current behavior; Testing Approach; Build Order 1
- ADR-0301; ADR-0255; ADR-0296; ADR-0259; ADR-0271; ADR-0285; ADR-0280
