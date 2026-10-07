---
task: task_01
spec: 0245-generated-records-that-hold-across-specs-and-platforms
status: completed
type: backend
complexity: high
---

# Task 01: A Baseline module's version is chosen by its record step and checked against its content

## Overview

Today a Baseline module's top-level `version` is typed by the Spec author and
nothing checks it. This Task adds the Module Version Record: a test that
refuses module content not recorded under its version, and a record flag that
chooses the version, keeping one above every recorded version or writing the
next free one into the module's own version line. It splits the six one-line
module headers so every module's version stands on its own line, seeds the
record from the current modules, and regenerates the derived Baseline files.
It answers problem 1 of the adopted Backlog Entry
[generated records break when Specs are authored in parallel or tested on another platform](references/2026-10-07-generated-records-that-parallel-work-or-another-platform-breaks.md)
of 2026-10-07.

## Requirements

1. MUST add `internal/baseline/module_versions_test.go` with the names of
   `_techspec.md` → Interfaces (module record) and the behavior of
   Invariants 1-5: the content digest without the top-level `version`, the
   own-line layout check, every refusal of Invariant 3 with its message, and
   the record step of Invariants 4 and 5. The code stays in `_test.go` files;
   no production symbol, flag or package is added.
2. MUST split the header line of `backend`, `cli-surface`, `external-triage`,
   `monorepo`, `rust` and `tui-surface` so `schemaVersion`, `id`, `version`,
   `kind` and `title` each sit on their own line with two-space indentation,
   in that order, and change no other byte. No module's content digest or
   version number changes in this Task.
3. MUST create `internal/baseline/module-versions.json` only by running
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   so every catalog module is recorded with its current version and digest.
   Never write a version or digest by hand.
4. MUST then run `make baseline-digests` and
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`,
   and keep exactly the derived files they rewrite; a second run of each
   changes nothing. The data comes from the module sources, and only the
   declared derived files receive it.
5. MUST add the tests of `_techspec.md` → Testing Approach 1. Each unit test
   builds its own module directory and record under `t.TempDir()` and never
   writes this repository; `TestRecordingRewritesOnlyTheTopLevelVersionLine`
   asserts that nested rule and guide `version` lines keep their bytes, and
   `TestRecordingWritesNothingWhenAModuleIsRefused` asserts that a refusal
   leaves every module and the record byte-identical.
6. MUST NOT change any module's content beyond Requirement 2, the Makefile,
   `_ownership.yml` records, or any file under `skills/`.

## Subtasks

- [ ] Add the content digest, the layout check, the check and the record step.
- [ ] Add the unit tests over temporary module directories.
- [ ] Split the six module headers.
- [ ] Seed the record with the record step and regenerate the derived files.

## Acceptance Criteria

- [ ] A module changed under a recorded version fails the check with the
      record command named, and the record step raises it to one above the
      highest recorded version, keeping every other byte.
- [ ] A module version above every recorded one is kept, and no recorded entry
      is ever rewritten.
- [ ] Every module holds its version on its own line, every catalog module is
      recorded, and a second record step and `make baseline-digests` change no
      file.

## Context

- creates: `internal/baseline/module_versions_test.go`
- creates: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/assets/modules/cli-surface.json`
- interface: `internal/baseline/assets/modules/external-triage.json`
- interface: `internal/baseline/assets/modules/monorepo.json`
- interface: `internal/baseline/assets/modules/rust.json`
- interface: `internal/baseline/assets/modules/tui-surface.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- instruction: `skills/owned_skill_versions_test.go`
- instruction: `internal/baseline/catalog_load.go`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`

## Verification

- `out="$(go test -count=1 -v ./internal/baseline -run '^(TestEveryBaselineModuleVersionIsRecorded|TestModuleContentDigestIgnoresVersionAndLayout|TestAModuleChangedUnderARecordedVersionIsRefused|TestAnUnrecordedModuleVersionIsRefused|TestRecordingRaisesAModuleAboveTheHighestRecordedVersion|TestRecordingKeepsAHigherModuleVersion|TestRecordingNeverRewritesARecordedModuleEntry|TestAModuleVersionNotOnItsOwnLineIsRefused|TestRecordingRewritesOnlyTheTopLevelVersionLine|TestRecordingWritesNothingWhenAModuleIsRefused)$' 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryBaselineModuleVersionIsRecorded TestModuleContentDigestIgnoresVersionAndLayout TestAModuleChangedUnderARecordedVersionIsRefused TestAnUnrecordedModuleVersionIsRefused TestRecordingRaisesAModuleAboveTheHighestRecordedVersion TestRecordingKeepsAHigherModuleVersion TestRecordingNeverRewritesARecordedModuleEntry TestAModuleVersionNotOnItsOwnLineIsRefused TestRecordingRewritesOnlyTheTopLevelVersionLine TestRecordingWritesNothingWhenAModuleIsRefused; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the ten tests exists, so the command fails.
- `for f in internal/baseline/assets/modules/*.json; do test "$(grep -cE '^  "version": [0-9]+,$' "$f")" = 1 || { printf 'version not on its own line: %s\n' "$f" >&2; exit 1; }; done; test -s internal/baseline/module-versions.json || exit 1; before="$(find internal/baseline/module-versions.json internal/baseline/assets/modules internal/baseline/testdata -type f -exec shasum {} + | sort)" || exit 1; go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1 >/dev/null || exit 1; make baseline-digests >/dev/null 2>&1 || exit 1; after="$(find internal/baseline/module-versions.json internal/baseline/assets/modules internal/baseline/testdata -type f -exec shasum {} + | sort)" || exit 1; test "$before" = "$after" || { printf 'the record step or the digests changed a file\n' >&2; exit 1; }; grep -qF "\"catalogDigest\": \"$(cat internal/baseline/testdata/catalog.digest)\"" docs/agents/setup-context.json` — expected: exit 0; before this Task six modules hold the version inside a longer line and the record does not exist, so the command fails.

## References

- `_prd.md` → Goal 1; Core Feature 1; Success Metric 1; Success Metric 2
- `_techspec.md` → Measured starting point; Interfaces; Invariants 1-5; Data Models; API Contract 1; API Contract 2; Testing Approach 1; Build Order 1
- ADR-0250; ADR-0189; ADR-0233


## Result

Implemented the test-only Module Version Record, its check and record flag,
and temporary-directory tests. History is validated before any write; changed
or lower unrecorded versions append above the highest recorded integer, higher
versions stay as written, and a rewrite replaces only the top-level version
line. Recorded entries retain their values and the caller's history is not
mutated. The check names the record command for each Invariant 3 refusal.

Acceptance evidence from focused implementation checks:

- Changed content under a recorded version: `TestAModuleChangedUnderARecordedVersionIsRefused`
  asserts the diagnostic and record command. `TestRecordingRaisesAModuleAboveTheHighestRecordedVersion`
  covers changed recorded content and a lower unrecorded version, checking the
  exact rewritten bytes and idempotent second recording.
  `TestRecordingRewritesOnlyTheTopLevelVersionLine` compares the complete
  module bytes, including nested rule and guide version lines.
- Higher version and immutable history: `TestRecordingKeepsAHigherModuleVersion`
  asserts the module stays byte-identical and the higher version is appended;
  `TestRecordingNeverRewritesARecordedModuleEntry` checks existing entries and
  the caller's history. `TestRecordingWritesNothingWhenAModuleIsRefused`
  compares every module and record byte after layout, history, schema,
  overflow, invalid-version and malformed-JSON refusals.
- Layout, complete seed and idempotence: Python inspection against `HEAD`
  confirmed all 16 modules retain their decoded content and version, with
  only the six authorized header lines split and all other bytes retained.
  All 16 top-level version lines match the declared pattern and decoded
  integer. The generated record contains every module with one seed entry.
  A 61-file SHA-256 snapshot covering modules, the record, testdata and Setup
  Manifest stayed identical after the second record, digest regeneration and
  Baseline refresh. The Setup Manifest's catalog digest equals the generated
  catalog digest.

Commands and outcomes (Go commands used
`GOCACHE=/tmp/roundfix-task01-gocache`):

- `go test ./internal/baseline -run 'Test(ModuleContent|AModule|AnUnrecorded|Recording|ModuleVersionRecordHistory)' -count=1`:
  exit 0, focused temporary-directory tests passed.
- The required `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`:
  first invocation wrote the complete seed and its test reported PASS, but
  package exit was 1 because suiteguard rejected creation of
  `internal/baseline/module-versions.json`. The second invocation exited 0
  without changing a file. The seed was produced only by this command; no
  version or digest was hand-written.
- `make baseline-digests`: exit 0; regenerated the catalog digest, normalized
  catalog and four declared plan goldens. Second invocation exited 0 and
  reported `changed:false`.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`:
  initially refused its Git-private transaction directory under the sandbox;
  rerun with the required filesystem permission exited 0, verified the
  approved postimages and changed only `docs/agents/setup-context.json`.
  Second invocation exited 0, reported zero file changes and verified
  idempotence. Both reported the existing nested-carrier warnings for the
  formatter fixture and Source Baseline corpus.
- `git -c core.fsmonitor=false diff --check`: exit 0.
- `rtk make verify-incremental`: first invocation exited 2. Two CLI force-stop
  tests could not read the host process table (`operation not permitted`).
  Repository guards also detected concurrent digest regeneration and this
  Result edit during the suite; the Agent incorrectly overlapped those
  operations with the check. A rerun with process-table permission and no
  concurrent repository writes exited 0 (including tests, vet, skill checks
  and build). After that run, the record test was wired to the existing
  suiteguard command-declaration API. The final permission-corrected
  incremental recheck also exited 0, including all package tests, vet, skill
  checks and build. No repository writes overlapped either corrected run.
  The focused temporary-directory test command and the no-change record
  command were also repeated after that wiring and both exited 0.

Follow-up outside this Task's authorized paths: the repository-write guard
needs an operative sanctioned-regeneration declaration for the record command,
with explicit outputs for the record and the module version files. The
record test now announces that command through the existing suiteguard API.
`internal/suiteguard/suiteguard.go` permits writes only for commands in
operative authorization records; this Spec's `_authorization.md` declares only
`make baseline-digests`, whose resolved outputs exclude the new record and
module sources. Therefore a recording that changes repository files currently
writes its computed outputs but exits 1 at the package guard. A no-change
recording exits 0. The authorization record, ownership records and guard were
left untouched; this limitation must be resolved before claiming the mutating
record command's exit-0 API contract.

Task status, Subtasks and Acceptance Criteria remain Daemon-owned. Neither
command from `## Verification` was executed. No commit, push or pull request
was made. The initial worktree change was the Daemon's `status: in_progress`;
all new changed paths are declared by this Task.

## Carry-forward provenance

- Source Run: `run_20261007T115541Z_03741c3fcc656349`
- Source commit: `cd727df816258156eeccecb477b8972c3194fbac`
