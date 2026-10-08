---
task: task_03
spec: 0254-a-faster-make-test
status: completed
type: test
complexity: medium
---

# Task 03: The sequential residue in internal/cli shrinks and pinned history reads less

## Overview

In the authoring prototype the 43 tests left sequential in `internal/cli`
still cost 124 s under load, and 78 s of that came from six
repository-profile tests that call `t.Chdir` only so that `baseline profile
init` finds the repository (`_techspec.md` → Current behavior). The command
environment already carries a per-test working directory
(`setCommandWorkDirForTest`), which is ADR-0089's answer. Separately,
`TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive` writes the
6,643 files of the pinned archived corpus to read its 759 QA Reports. This
Task gives the six tests their per-test working directory, makes them
parallel, lowers the `internal/cli` ceiling to 40, and narrows two
pinned-history fixtures to the paths they read (ADR-0259). It is verifiable
on its own: the six tests are parallel and pass, the ceiling holds, and the
narrowed tests pass reading only QA Reports. The authoring prototype made
exactly this change to the six tests and they passed in parallel.

This is an authorized tooling Task for one Governed Path,
`internal/cli/cli_test.go`, which it may change only to replace process-wide
state with a per-test value. It may change only the paths in its Context and
this Task file.

## Requirements

1. MUST make `newBaselineRepositoryProfileFixture` set the working directory
   with `setCommandWorkDirForTest` instead of `t.Chdir`, make
   `runRepositoryProfileCommand` run the command through `runCLIContext`, which
   reads the per-test command environment, instead of `RunContext`, and make
   `TestBaselineUpdateReachesCurrentWithARepositoryProfile`,
   `TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile`,
   `TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile`,
   `TestBaselineSkillsRestoreAcceptsARepositoryProfile`,
   `TestBaselineSkillsRestoreNamesTheProfilePath` and
   `TestDoctorComparesARepositoryProfileWithItsSnapshot` call `t.Parallel()`
   first. Any other test that called the helper and was sequential only for
   it MUST become parallel too.
2. SHOULD convert any other `internal/cli` Sequential Test whose only
   process-wide state is a value the command environment already carries per
   test (the working directory, `HOME`, or an environment entry the command
   reads through its environment rather than from the process). A test whose
   state a child process reads from the process environment, or that swaps a
   package-level variable, stays sequential. No production file changes for
   this.
3. MUST lower the `internal/cli` ceiling to 40 and leave the other ceilings
   unchanged.
4. MUST make `TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive`
   and `TestArchivedQAReportCorpusRemainsReadable` pass `gittest.PinnedHistory`
   a glob pathspec for the archived QA Reports
   (`:(glob)<archive dir>/*/qa/qa-report-*.md`, built from
   `ArchiveDir(ArchiveKindSpec)`) instead of the whole archive directory. Both
   tests MUST still find and check every QA Report they found before; each
   MUST fail when its glob matches no file. Tests that read whole archived
   Spec folders keep the whole directory.
5. MUST NOT change any assertion, expected value, golden or production file
   (Invariant 1). `internal/gittest` does not change.
6. MUST record in the Result the time of `go test -count=1 ./internal/spec`
   before and after the pathspec change, on the same machine, back to back.

## Subtasks

- [ ] Replace the working-directory change with the per-test value.
- [ ] Make the six tests and any other eligible test parallel.
- [ ] Lower the `internal/cli` ceiling.
- [ ] Narrow the two pinned-history pathspecs and time `internal/spec`.

## Acceptance Criteria

- [ ] The six repository-profile tests are parallel and pass.
- [ ] The rule logs `internal/cli` within a ceiling of 40.
- [ ] The two QA Report sweeps read only QA Reports and check the same set.
- [ ] `internal/cli` and `internal/spec` pass with `-race -short`.

## Context

- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `docs/adr/0089-code-under-test-takes-its-environment-explicitly.md`
- instruction: `internal/gittest/pinned_history.go`
- creates: `internal/testfixture/parallel_tests_test.go`
- interface: `internal/cli/baseline_update_repository_profile_test.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/carryforward_hooks_test.go`
- interface: `internal/cli/deliver_item_binary_test.go`
- interface: `internal/cli/readiness_forge_test.go`
- interface: `internal/spec/qa_frontmatter_test.go`
- interface: `internal/spec/qa_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestEveryTestInAParallelTestPackageRunsInParallel$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -Eq 'internal/cli: ([0-9]|[1-3][0-9]|40) Sequential Tests, ceiling 40' || { printf 'internal/cli is not within its ceiling of 40\n' >&2; exit 1; }; ! grep -q 't[.]Chdir' internal/cli/baseline_update_repository_profile_test.go || { printf 'the repository-profile fixture still changes the process working directory\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^(TestBaselineUpdateReachesCurrentWithARepositoryProfile|TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile|TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile|TestBaselineSkillsRestoreAcceptsARepositoryProfile|TestBaselineSkillsRestoreNamesTheProfilePath|TestDoctorComparesARepositoryProfileWithItsSnapshot)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestBaselineUpdateReachesCurrentWithARepositoryProfile TestBaselineUpdatePreviewListsATrailingSkillWithARepositoryProfile TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile TestBaselineSkillsRestoreAcceptsARepositoryProfile TestBaselineSkillsRestoreNamesTheProfilePath TestDoctorComparesARepositoryProfileWithItsSnapshot; do printf '%s\n' "$out" | grep -q -- "=== PAUSE $name" || { printf '%s is not parallel\n' "$name" >&2; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the ceiling is 48 and the six tests never pause, so the command fails.
- `for file in internal/spec/qa_frontmatter_test.go internal/spec/qa_test.go; do grep -qF ':(glob)' "$file" || { printf '%s still materializes the whole archive\n' "$file" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive|TestArchivedQAReportCorpusRemainsReadable)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive TestArchivedQAReportCorpusRemainsReadable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -race -short -timeout 30m ./internal/cli ./internal/spec || exit 1` — expected: exit 0; before this Task the sweep passes the whole archive directory, so the first check fails.

## References

- [_prd.md](_prd.md) — Goals 3–4; User Stories 1–2; Core Feature 4; Success Metric 2
- [_techspec.md](_techspec.md) — Current behavior; API Contract 2; Invariants 1–3; Testing Approach; Build Order 3
- ADR-0259; ADR-0089; ADR-0090


## Result

The six repository-profile tests now call `t.Parallel()` first. Their shared
fixture supplies `setCommandWorkDirForTest(t, root)` and their command helper
uses `runCLIContext(t, t.Context(), ...)`, so profile discovery receives the
per-test working directory. Searching `internal/cli` found no additional
callers of that fixture. The CLI Sequential Test ceiling falls from 48 to 40;
the other ceilings keep their original values.

Both archived QA Report sweeps now request
`":(glob)" + ArchiveDir(ArchiveKindSpec) + "/*/qa/qa-report-*.md"` from
`gittest.PinnedHistory`. Their report discovery, assertions and empty-corpus
refusals are unchanged. Tests that read whole archived Spec folders retain
their whole-directory fixtures. No production file, expected value, golden,
or `internal/gittest` file changed.

### Focused evidence

All Go checks used `GOCACHE=/private/tmp/roundfix-task03-gocache` on this
machine with Go 1.26.5. The cache was warmed by
`go test -run '^$' ./internal/spec` before the timing pair.

| Acceptance criterion | Evidence in this turn |
| --- | --- |
| Six repository-profile tests are parallel and pass | `go test -count=1 -race -short -shuffle=on -v -run 'RepositoryProfile\|ProfilePath' ./internal/cli` exited 0 (`ok`, 5.742 s). The captured output contains `=== PAUSE` and `--- PASS` for each of the six named tests; an explicit output check confirmed all six. |
| CLI rule stays within ceiling 40 | `go test -count=1 -v ./internal/testfixture` exited 0. The repository rule logged `internal/cli: 37 Sequential Tests, ceiling 40`. |
| QA sweeps read only reports and check the same set | A read-only Python inspection of `git archive --format=tar` at `gittest.PinnedHistoryRevision` compared the whole archive and the new glob. Both contain the identical nonempty set of 352 QA Report paths across 202 Spec directories. Regular files materialized fall from 5,570 to 352. The sorted report-path set has SHA-256 `786ef852156223212383c5ad64b81c388b53b40c8ac9fc1fd3b5af94306aa19a` in both archives. Both `go test -count=1 ./internal/spec` runs passed, exercising both sweeps. |
| CLI and Spec pass with race detection and short mode | Focused CLI check above and `go test -count=1 -race -short -shuffle=on -run 'Test(QAReportReaderAgreesWithTheDerivedVerification\|ArchivedQAReportCorpus)' ./internal/spec` exited 0; Spec reported `ok`, 8.011 s. These cover the changed tests. The full-package `-race -short` check remains Daemon-owned and was not run in this turn. |

The task's two declared Verification commands, `make verify`, and
`make verify-changed` were not run. Focused output was captured at
`/private/tmp/roundfix-task03-profile-check.log`,
`/private/tmp/roundfix-task03-rule-check.log`, and
`/private/tmp/roundfix-task03-qa-check.log`.

### Back-to-back Spec timing

`/usr/bin/time -p go test -count=1 ./internal/spec` ran immediately before
and immediately after changing only the two QA sweep pathspecs, with the same
machine, cache, and environment. No other test command ran during the pair.

| Fixture | Package test time | Wall | User CPU | System CPU | Outcome |
| --- | --- | --- | --- | --- | --- |
| Whole archive directory | 17.185 s | 17.34 s | 88.74 s | 43.84 s | exit 0 |
| Archived QA Report glob | 10.341 s | 10.84 s | 31.75 s | 39.84 s | exit 0 |

This is one local pair, not a claim about CI or the full repository suite.
The pinned corpus available here has 352 reports; the authoring prototype's
759-report count is not used as the current measurement.

### Scope review

Other Sequential Tests in the allowed CLI files keep their existing reason:
carry-forward hooks, item-binary fixtures, and forge probes supply environment
values that their child processes read; the TUI and color cases explicitly
test process defaults under ADR-0089; the detached-child branch calls
`os.Unsetenv` in production. None qualifies for an additional conversion
without changing its contract or production code.

`gofmt -l` reported no formatting changes for the four changed Go test files.
`git diff --check` passed. The initial worktree change was the Daemon's
`status: in_progress` in this Task file. The implementation changes only the
four declared test paths and this Result section. Task status, checkboxes,
the Task Graph, and other Task files remain Daemon-owned. No follow-up work
was added to this diff.

## Carry-forward provenance

- Source Run: `run_20261008T215927Z_c26d3579de5392c8`
- Source commit: `0fb78b66a0d1438979c4f104970bb9f03d657521`
