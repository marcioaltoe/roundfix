---
task: task_02
spec: 0364-a-rerun-finishes-a-partial-npm-publication
status: pending
type: infra
complexity: medium
---

# Task 02: The Publication Preflight lets an already published coordinate through

## Overview

The rerun of a release job whose npm publication stopped part way refuses in
the Publication Preflight with `registry: <coordinate>@<version> is already
used` for every coordinate the first attempt published (runs 41750342366 and
42056741468, attempt 2). This Task gives the preflight a `published` state
for a version listed under the packument's `versions`, writes those
coordinates to `$RUNNER_TEMP/roundfix-release-published` for the next steps,
keeps `used` for a version listed only under `time`, and stops when all six
coordinates are published. It also builds the Go test harness that runs the
workflow's own step scripts with stub `curl` and `npm`, and characterizes
today's publish failure on a transparency-log 409. Task 03 adds the content
comparison and the publish skip.

## Requirements

1. MUST change the step "Publication preflight" of
   `.github/workflows/release.yml` (it stays where it is, on both triggers):
   - `CLASSIFY_JQ` becomes the classifier `_techspec.md` → Interfaces
     gives: `published` when `.versions` has the tag, else `used` when
     `.time` has it, else `cooldown` when `.time` has `unpublished`, else
     `eligible`, with the same `invalid packument` error;
   - before the loop it empties `$RUNNER_TEMP/roundfix-release-published`
     (`PUBLISHED="$RUNNER_TEMP/roundfix-release-published"; : > "$PUBLISHED"`)
     and counts the coordinates and the published ones;
   - a `published` coordinate keeps its table row and summary row with the
     state `published`, is appended to that file, prints exactly
     `::notice::registry: <coordinate>@<version> is already published; a tag run compares it with its own build`
     (`echo "::notice::registry: $coordinate@$TAG is already published; a tag run compares it with its own build"`)
     and does not block;
   - after the existing `blocked` exit, when every coordinate is
     `published`, it prints exactly
     `echo "::error::registry: every coordinate of the Release Set is already published at $TAG"`
     and exits 1;
   - every other state, line and exit keeps its bytes, including
     `::error::registry: $coordinate@$TAG is already used`.
2. MUST create `internal/verifyselect/release_workflow_test.go`, headed with
   a `// Suite:` / `// Invariant:` / `// Boundary IN:` / `// Boundary OUT:`
   comment like its neighbours, holding the harness `_techspec.md` →
   Testing Approach describes:
   - it reads `.github/workflows/release.yml` from the repository root
     (`findRepositoryRoot`) with `gopkg.in/yaml.v3`, finds a step of job
     `release` by name, writes its `run` text to a temporary file and runs
     `bash -e <file>` with `exec.CommandContext`, using the path
     `testfixture.RequireTool(t, "bash")` returns (the External Test Tool
     rule refuses a bare `"bash"`), in a temporary repository
     that holds copies of `dist/npm/platforms.json`,
     `dist/npm/roundfix/package.json` and each
     `dist/npm/packages/*/package.json`;
   - the child's environment is built, never inherited: `PATH` is an
     ABSOLUTE stub directory followed by the test process's `PATH`;
     `HOME` and `TMPDIR` are temporary; `RUNNER_TEMP`,
     `GITHUB_STEP_SUMMARY`, `TARGET_VERSION`, `GITHUB_REF_NAME`,
     `GITHUB_EVENT_NAME`, `FALLBACK_WINDOW` are set by the case, and
     `NPM_FALLBACK_TOKEN` is an obviously fake value built by
     concatenation; `NODE_AUTH_TOKEN` is never set by the harness;
   - a stub `curl` (an `sh` script) writes the packument fixture for the
     requested coordinate to the `--output` file and prints `200`, or prints
     `404` when there is none;
   - a stub `npm` (an `sh` script) serves `pack` of the package in its
     working directory and `pack <coordinate>@<version>` from fixture
     tarballs into `--pack-destination` (or the working directory), printing
     `[{"filename":"<file>","integrity":"<integrity>"}]` and exiting 1 when
     the fixture is missing; it appends each `publish` to a log with the
     directory's base name and whether `NODE_AUTH_TOKEN` is set, and for a
     configured directory prints configured output and exits 1;
   - it resolves `bash`, `jq`, `gzip` and `cmp` with
     `testfixture.RequireTool`, so a host without one skips naming it;
   - it writes each stub with a non-executable mode and then `os.Chmod`s it
     to `0o755`, as `full_contract_run_test.go` does: the Written Executable
     rule (`TestNoTestWritesAnExecutableOutsideTheResidue`) refuses a
     literal executable mode in `os.WriteFile`;
   - every top-level test and subtest calls `t.Parallel()` first, and the
     harness waits only on the step's exit or the test's context;
   - no case reaches the network, publishes, or signals a process, and every
     command ends with the test's context.
3. MUST add to that file
   `TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne`, a
   table over the step "Publication preflight" for version `0.132.1` with
   packuments built per coordinate, asserting exit code, the stdout lines
   and the published file:
   - all six eligible: exit 0, no `published` row, the file empty;
   - the v0.132.1 shape (five platform coordinates published, `roundfix`
     eligible): exit 0, five `published` rows and five notice lines, the
     file listing the five platform coordinates in Release Set order;
   - a version under `time` only for `roundfix` (the others eligible): exit
     1 with `::error::registry: roundfix@0.132.1 is already used` and the
     file empty;
   - all six published: exit 1 with
     `::error::registry: every coordinate of the Release Set is already published at 0.132.1`;
   - one platform coordinate published and `roundfix` in cooldown (`.time`
     has `unpublished`): exit 1 with the cooldown line, so a published
     coordinate never hides a blocking one.
4. MUST add `TestReleasePublishStopsAtATransparencyLog409WithoutATokenRetry`:
   the step "Publish to npm" after an all-eligible preflight, with
   `FALLBACK_WINDOW=1` and the stub `npm` failing `roundfix`'s publish with
   the two lines `npm error code TLOG_CREATE_ENTRY_ERROR` and
   `npm error error creating tlog entry - (409) an equivalent entry already exists in the transparency log`;
   it asserts exit 1, the line
   `::error::publish: roundfix failed for a non-identity reason`, the five
   platform publishes then exactly one `roundfix` publish, and no publish
   with `NODE_AUTH_TOKEN` set. This characterizes today's behavior, which
   this Spec keeps.
5. MUST run `go test -timeout 20m -count=1 ./internal/testfixture` after
   adding the file and keep it passing, because its repository rules
   (External Test Tools, Written Executables, Parallel Test Packages,
   Real-Time Budgets) scan the new test file.
6. MUST keep `TestMainAndReleaseRunEveryContract` passing unchanged; MUST
   NOT change the steps "Cross-compile and stage", "Publish to npm" or
   "GitHub Release", any other workflow, the Makefile, `go.mod`, `go.sum`,
   or any file under `internal/suiteguardcontract`.

## Subtasks

- [ ] Change the preflight classifier, its published list and the all-published refusal.
- [ ] Build the step-script harness with the stub `curl` and `npm`.
- [ ] Add the preflight table test and the 409 characterization test.

## Acceptance Criteria

- [ ] A coordinate already published at the version passes the preflight as `published` and is listed for the next steps.
- [ ] A version listed only under `time`, a cooldown and a fully published Release Set still stop the preflight with their lines.
- [ ] A transparency-log 409 still stops the publish step after one attempt, without the token fallback.

## Context

- instruction: `docs/adr/0374-a-rerun-completes-a-release-set-from-the-same-content.md`
- instruction: `docs/adr/0146-release-publication-is-all-or-nothing-with-a-bounded-token-fallback.md`
- instruction: `docs/user-guide/release-runbook.md`
- instruction: `internal/verifyselect/full_contract_run_test.go`
- instruction: `internal/verifyselect/main_test.go`
- instruction: `internal/testfixture/external_tool.go`
- instruction: `internal/testfixture/written_executable_test.go`
- instruction: `dist/npm/platforms.json`
- interface: `.github/workflows/release.yml`
- creates: `internal/verifyselect/release_workflow_test.go`

## Consumers

- The preflight's `published` state and `$RUNNER_TEMP/roundfix-release-published` — consumers: `./internal/verifyselect`

## Verification

- `for t in TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne TestReleasePublishStopsAtATransparencyLog409WithoutATokenRetry; do grep -q "func $t(" internal/verifyselect/release_workflow_test.go || { printf 'missing %s\n' "$t" >&2; exit 1; }; done; grep -qF -- 'roundfix-release-published' .github/workflows/release.yml || exit 1; go test -timeout 20m -count=1 ./internal/testfixture || exit 1; out="$(mktemp)"; go test -timeout 20m -count=1 -v -run '^(TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne|TestReleasePublishStopsAtATransparencyLog409WithoutATokenRetry|TestMainAndReleaseRunEveryContract)$' ./internal/verifyselect > "$out" 2>&1; status=$?; cat "$out"; passed=0; grep -q -- '--- PASS: TestReleasePreflightClassifiesAPublishedVersionApartFromAUsedOne' "$out" && grep -q -- '--- PASS: TestReleasePublishStopsAtATransparencyLog409WithoutATokenRetry' "$out" && passed=1; rm -f "$out"; test "$status" -eq 0 && test "$passed" -eq 1` — expected: exit 0. Before this Task the test file and the published list do not exist, so the first check fails; afterwards the repository test rules accept the new file, the preflight table and the 409 characterization pass and are not skipped, and the release step order still holds.

## References

- `_prd.md` → Goals; User Stories 1-2; Core Feature 1; Success Metric 4; Success Metric 5
- `_techspec.md` → Interfaces; API Contracts 1 and 3; Testing Approach; Build Order 2
- ADR-0374
