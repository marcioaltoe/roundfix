---
task: task_04
spec: 0213-a-test-suite-that-does-not-flake
status: completed
type: test
complexity: medium
---

# Task 04: The Assets Sync template lives in memory

## Overview

The Assets Sync tests in `internal/baseline` copy one on-disk template that
stays in the system temporary directory for the whole package run. When the
maintainer cleared system temporaries during the 0203 gate run, every copy
failed with `open .git/objects/05: no such file or directory`. Removing the
template during a run reproduces that message in four of four consumer tests.
The template is also built inside `sync.Once` with the first caller's
`testing.T`. This Task builds the template without a `testing.T`, reads it into
memory and removes the on-disk build, so each test writes its own copy into its
own temporary directory.

## Requirements

1. MUST replace the template build with `buildAssetsSyncTemplate`, as the
   TechSpec section "The Assets Sync template" describes:
   - it takes no `testing.T` and returns its error;
   - it runs Git through an error-returning runner with the same `-c` flags;
   - it reads the target and source repositories, `.git` included, into
     memory with their file modes;
   - it removes its build directory and records the path.
2. MUST keep one build per package process through the existing `sync.Once`,
   and report a build error once, from every caller, as
   `build Assets Sync template: <error>`.
3. MUST make `newAssetsSyncTarget` and `newAssetsSyncSource` write the
   in-memory entries into `t.TempDir()`, and keep their signatures and the
   revision they return.
4. MUST remove `removeAssetsSyncTemplate` and its `TestMain` call, and keep the
   suite guard installation unchanged.
5. MUST add `TestAssetsSyncTemplateLeavesNoDirectoryBehind`. It forces the
   build, then asserts that:
   - the recorded build directory does not exist;
   - a target copy passes `git fsck` and `git rev-parse HEAD`;
   - a source copy's `HEAD` equals the recorded revision.
6. MUST NOT change production code or any other test's assertions.

## Subtasks

- [ ] Build the template without a `testing.T` into memory.
- [ ] Write per-test copies from memory.
- [ ] Add the no-directory-behind test and repeat the Assets Sync tests.

## Acceptance Criteria

- [ ] No `assets-sync-template` directory exists once the template is built,
      and per-test copies are valid Git repositories.
- [ ] The Assets Sync tests pass three repeated runs at `-cpu 1,4`.

## Context

- interface: `internal/baseline/assets_sync_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestAssetsSyncTemplateLeavesNoDirectoryBehind$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestAssetsSyncTemplateLeavesNoDirectoryBehind' || { printf '%s\n' "$out"; exit 1; }` — expected: exit 0; before this Task the test does not exist, so the command fails.
- `! grep -n 'removeAssetsSyncTemplate' internal/baseline/assets_sync_test.go || exit 1; out="$(go test -count=3 -cpu 1,4 -v -run '^(TestAssetsSync|TestBaselineAssetsSync)' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestAssetsSyncProvenanceAndPreMutationRefusals' || { printf '%s\n' "$out"; exit 1; }` — expected: exit 0 with no failed iteration; before this Task the on-disk template removal exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 4; Core Feature 5; Success Metric 5; Acceptance evidence
- [_techspec.md](_techspec.md) — The Assets Sync template; Testing Approach 4; Build Order 4
- ADR-0126

## Result

The Assets Sync template now stores both repository trees in memory, including
`.git`, directory modes and file modes. `buildAssetsSyncTemplate` takes no
`testing.T`, returns errors, and removes its recorded build directory before
returning, including on error. The existing `sync.Once` stores the build error
for every caller to report as `build Assets Sync template: <error>`. Git uses
the same configuration flags and fixture identity through an error-returning
runner.

Each target and source helper restores entries into its caller's `t.TempDir()`;
their signatures and recorded source revision are preserved. The shared
directory cleanup function and its `TestMain` call are removed. The suite guard
installation, production code and existing test assertions are unchanged.

Focused-check evidence:

- Before changing the helpers,
  `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test -count=1 -run 'TestAssetsSyncTemplateLeaves' ./internal/baseline`
  exited 1: the new regression test found that the template build directory
  remained after construction.
- After implementation,
  `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test -count=3 -cpu=1,4 -run 'AssetsSync' ./internal/baseline`
  exited 0 (`ok roundfix/internal/baseline 113.266s`). This focused selection
  includes every Assets Sync test and the new regression test.
- Acceptance criterion 1: the regression test asserts the recorded build
  directory is absent before creating copies, runs target `git fsck` and
  `git rev-parse HEAD`, and compares source `HEAD` with the recorded revision.
  It passed in the focused repeated check.
- Acceptance criterion 2: the focused selection passed three repetitions at
  each of CPU settings 1 and 4, including the provenance refusal subtests.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0 after the code
  changes; diff review found only this test helper slice and the assigned Task
  file changed.

The declared Verification commands were not run. Status remains Daemon-owned;
no commit, push or pull request was created. No follow-up work was identified.

## Carry-forward provenance

- Source Run: `run_20261002T064202Z_a22e1dd4644c4dc6`
- Source commit: `b022095fae3d3e54fcecd27f6b0e3a424fb2262b`
