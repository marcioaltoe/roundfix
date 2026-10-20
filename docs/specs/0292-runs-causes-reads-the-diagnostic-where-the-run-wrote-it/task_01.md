---
task: task_01
spec: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
status: pending
type: backend
complexity: medium
---

# Task 01: Runs causes reads the diagnostic under the resolved Artifact Root

## Overview

`roundfix runs causes` reads a failed attempt's diagnostic only under the
Artifact Directory the Run recorded, taken literally, and every Implement Run
records none, so no Implement Run's diagnostic is ever read (`_prd.md` →
Sources 1). This Task makes the report find the Run's directory through
`store.ResolveRunArtifactRoot`, the resolver the storage report, GC
sanitation and Run Retention use, and keeps every guard on the diagnostic. It
is verifiable on its own on temporary Roundfix Homes.

## Requirements

1. MUST implement API Contract 1 of `_techspec.md` in
   `internal/runcause/report.go`: add `HomeDir` to `Request` as
   `_techspec.md` → Interfaces sketches it, and make `diagnosticTail` take the
   home and build the Run's directory as `<root>/runs/<run-id>`, where the
   root is `store.ResolveRunArtifactRoot` called with the Run's ID, Git root,
   repository root (the Git root when the recorded repository root is blank)
   and trimmed recorded Artifact Directory. Production code in `internal/runcause` MUST NOT call
   `config.ResolveArtifactDirectory`, `config.DefaultArtifactDirectoryForPath`
   or `config.ValidateArtifactDirectory`, and MUST NOT join the recorded value
   into a path itself.
2. MUST map the resolver's errors as API Contract 1 states: an error wrapping
   `store.ErrUnsafeRunID` or `store.ErrArtifactRootUnresolved` gives no
   evidence (an empty tail) without any further filesystem read and without
   failing the report; any other resolver error fails `Build` with
   `resolve Artifact Root of Run "<id>": <error>`.
3. MUST keep every guard unchanged on the resolved directory: the absolute
   diagnostic path, the `filepath.Rel` containment in `runs/<run-id>`, the
   `os.Lstat` walk of every component refusing a symbolic link or a missing
   component, the regular-file requirement, the tail bound and the existing
   error texts.
4. MUST implement API Contract 2 in `internal/cli/runs_causes.go`: set
   `Request.HomeDir` from `loaded.HomeDir`, the value that opens the Run
   Database reader. Flags, help, exit codes, text output and JSON fields do
   not change.
5. MUST state the hazard in the code's behavior: data comes from the Run
   Database row, the recorded repositories' Git metadata and the loaded home;
   the resolved root and the diagnostic text reach only the classifier; no
   diagnostic text, path or key reaches stdout or JSON; the report creates,
   writes and removes nothing.
6. MUST declare this break, covered by the tests below: an attempt of a Run
   that recorded a blank or relative Artifact Directory, whose diagnostic sits
   under the resolved root and passes the guards, is now classified from its
   diagnostic (API Contract 3). Every existing test keeps its name and body
   and passes, including `TestBuildIgnoresALogOutsideTheArtifactDirectory`,
   `TestBuildReportsAttemptsTasksAndCorrectives`,
   `TestBuildBoundsEvidenceAndExportsNoPrivateCommandText`,
   `TestRunsCausesPrintsAClassifiedWindow`, `TestRunsCausesPrintsJSON` and
   both `TestBuildLeavesTheRunDatabaseUnchanged` tests. An existing fixture
   helper may gain parameters or resolve its repository path's symbolic links
   when the existing tests that use it keep their assertions.
7. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows, using temporary homes and repositories whose paths
   have their symbolic links resolved, never reading `~/.roundfix`, and
   starting no process in `internal/runcause`:
   - in `internal/runcause/report_test.go`,
     `TestBuildReadsTheDiagnosticUnderTheResolvedArtifactRoot`, with the
     subtests `_techspec.md` → Testing Approach lists: `blank root`, `relative
     root`, `absolute root`, `blank root without home`, the six guards
     (`leaf symlink`, `parent symlink`, `outside`, `other run`, `directory`,
     `missing`) applied under the default root, `diagnostic under the recorded
     value` and `inspection error` (which skips on Windows with a reason
     naming the platform, because Windows reports a path through a file as
     missing); each subtest asserts task_01 attempt 1's class, signature and
     evidence (or `Build`'s error), and that the home's and repository's file
     hashes are identical before and after `Build`;
   - in `internal/cli/runs_causes_test.go`,
     `TestRunsCausesReadsTheDiagnosticOfARunWithoutAnArtifactRootTranscript`,
     asserting Surface Transcript 1 of `_techspec.md` whole: every stdout line
     in order with only `<run id>` and `<digest>` substituted, empty stderr and
     exit 0;
   - in `internal/cli/runs_causes_test.go`,
     `TestRunsCausesReadsTheDiagnosticUnderARelativeArtifactRoot`, through the
     command with the Run recording `rel-artifacts` and its directory under
     the repository: the same stdout as the transcript.
8. MUST add to `docs/user-guide/commands/runs.md`, in the "Why Verification
   failed" section after the paragraph that names the five classes, the
   sentence "A failed attempt's diagnostic is read under the Artifact Root
   the Run wrote to, resolved the way `roundfix storage report` resolves it:"
   followed, in the same paragraph, by the three resolution rules (a blank
   recorded Artifact Directory uses the repository's default root, a relative
   one resolves under the Run's Git root, an absolute one is used as
   recorded) and the sentence "A diagnostic outside the Run's own
   `runs/<run-id>` directory, reached through a symbolic link, or not a
   regular file gives the attempt no evidence." The paragraph does not end
   with a colon.
9. MUST re-record the Behavior Surface Record, because this Task edits a user
   guide page: `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   The Agent runs it itself before handing back; no record is hand-edited.

## Subtasks

- [ ] Add the home to the request and resolve the Run's root in the read.
- [ ] Pass the loaded home from the command.
- [ ] Add the report and command tests.
- [ ] Add the guide paragraph and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A Run that recorded no Artifact Directory, or a relative one, has its diagnostic read under the root it wrote to, and its attempts are classified as for an absolute root.
- [ ] Every guard on the diagnostic holds on the resolved directory, and an unresolvable root gives no evidence without failing the report.
- [ ] `internal/runcause` carries no resolution of its own.
- [ ] The report writes nothing.
- [ ] The `runs` guide states where diagnostics are read, the three resolution rules and the guards.

## Consumers

- `runcause.Request` and the `runs causes` report — consumers: `./internal/runcause`, `./internal/cli`

## Context

- instruction: `docs/adr/0301-gc-sanitation-classifies-the-artifact-root-each-run-wrote-to.md`
- instruction: `docs/adr/0214-why-a-verification-failed-is-classified-by-deterministic-signatures.md`
- instruction: `internal/store/journal.go`
- interface: `internal/runcause/report.go`
- interface: `internal/runcause/report_test.go`
- interface: `internal/cli/runs_causes.go`
- interface: `internal/cli/runs_causes_test.go`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestBuildReadsTheDiagnosticUnderTheResolvedArtifactRoot|TestBuildIgnoresALogOutsideTheArtifactDirectory|TestBuildReportsAttemptsTasksAndCorrectives|TestBuildBoundsEvidenceAndExportsNoPrivateCommandText|TestBuildLeavesTheRunDatabaseUnchanged)$' ./internal/runcause 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBuildReadsTheDiagnosticUnderTheResolvedArtifactRoot TestBuildIgnoresALogOutsideTheArtifactDirectory TestBuildReportsAttemptsTasksAndCorrectives TestBuildBoundsEvidenceAndExportsNoPrivateCommandText TestBuildLeavesTheRunDatabaseUnchanged; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for sub in blank_root relative_root absolute_root blank_root_without_home leaf_symlink parent_symlink outside other_run directory missing diagnostic_under_the_recorded_value inspection_error; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestBuildReadsTheDiagnosticUnderTheResolvedArtifactRoot/$sub " || { printf 'missing subtest pass: %s\n' "$sub" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestRunsCausesReadsTheDiagnosticOfARunWithoutAnArtifactRootTranscript|TestRunsCausesReadsTheDiagnosticUnderARelativeArtifactRoot|TestRunsCausesPrintsAClassifiedWindow|TestRunsCausesPrintsJSON|TestBuildLeavesTheRunDatabaseUnchanged)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunsCausesReadsTheDiagnosticOfARunWithoutAnArtifactRootTranscript TestRunsCausesReadsTheDiagnosticUnderARelativeArtifactRoot TestRunsCausesPrintsAClassifiedWindow TestRunsCausesPrintsJSON TestBuildLeavesTheRunDatabaseUnchanged; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -q 'store.ResolveRunArtifactRoot(' internal/runcause/report.go || { echo 'runs causes does not call the shared resolver' >&2; exit 1; }; ! grep -rqE --exclude='*_test.go' 'ResolveArtifactDirectory|DefaultArtifactDirectoryForPath|ValidateArtifactDirectory' internal/runcause || { echo 'runs causes carries its own resolution' >&2; exit 1; }; grep -q 'HomeDir = loaded.HomeDir' internal/cli/runs_causes.go || { echo 'the command does not pass the home' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/runcause ./internal/cli || exit 1; GOOS=windows go vet ./internal/runcause || exit 1; tr -s '[:space:]' ' ' < docs/user-guide/commands/runs.md | grep -qF -- 'diagnostic is read under the Artifact Root the Run wrote to' || { echo 'missing guide sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/runs.md | grep -qF -- 'regular file gives the attempt no evidence' || { echo 'missing guide guard sentence' >&2; exit 1; }; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new report test does not exist, so the command fails at its first pass check. Afterwards the resolved-root tests and the existing guard, transcript and read-only tests pass, the report calls the shared resolver and carries no resolution of its own, the command passes the home, both packages build for Windows with the report's tests compiling there, the guide states the rule, and the Behavior Surface Record is current.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Invariants 1-5; Interfaces; Data Models; Current behavior; Testing Approach; Build Order 1
- ADR-0301; ADR-0214; ADR-0259; ADR-0271; ADR-0285; ADR-0280
