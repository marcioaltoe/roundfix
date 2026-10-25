---
task: task_05
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
status: completed
type: backend
complexity: medium
---

# Task 05: The suite guard fingerprints the repository through Git

## Overview

The maintainer's answer of 2026-10-25, at about 11:00 BRT, keeps task_04
and also asks for the second and last corrective Task: `internal/suiteguard`
computes its repository fingerprint through `git` subprocesses instead of
walking the tree and digesting every file in the test process. The guard
keeps every detection it has today. Git lists the watched paths and hashes
their content. A mode-only `os.Lstat` pass, which runs before and after
`m.Run` and so is not a Go cache input, keeps the permission bits. The
first step characterizes the current guard with a table of 22 leak cases;
the same table must pass after the change.

## Requirements

1. MUST answer the maintainer's answer of 2026-10-25 quoted in
   `_authorization.md` (`_prd.md` → Sources 3) by implementing API Contract
   3 and Invariant 4 of `_techspec.md`.
2. MUST characterize first. Before changing `suiteguard.go`, add the new
   file `internal/suiteguard/leak_cases_test.go` and the fixture changes of
   Requirement 5, and run them against the unchanged guard. Every leak case
   and `TestGuardLeavesNoRepositoryPathInTheTestInputLog` pass on it, and only
   `TestGuardRefusesARootOutsideAGitWorkTree` fails. Record that run in the
   Task's Result. The file holds:
   - a fixture builder that runs `git init` in a temporary directory, writes
     a root `.gitignore` with `/ignored/`, `*.log` and `!keep.log`, a
     `sub/.gitignore` with `local.tmp`, the tracked files `tracked.txt`,
     `dirty.txt` and `sub/nested.txt` and a tracked symbolic link `link` to
     `tracked.txt`, commits them with a fixture identity, then leaves
     `dirty.txt` modified, `scratch.txt` untracked and an empty untracked
     directory `emptydir`;
   - `TestGuardReportsEveryLeakClass`, one parallel subtest per case, each
     on its own fixture, which runs the mutation inside the guarded
     subprocess (the existing `TestGuardFixture`, new mode `leak`, case name
     in `ROUNDFIX_SUITEGUARD_FIXTURE_CASE`) and requires a failing exit with
     each expected `suiteguard: <change>: <path>` line, or exit 0 for an
     ignored path. The 22 cases are: tracked file edited; tracked file
     deleted; tracked file emptied keeps its size zero
     (`sub/nested.txt`); dirty tracked file edited again; untracked file
     created; untracked file created in a new directory (wants
     `created: newdir` and `created: newdir/deep/created.txt`); untracked
     file created beside tracked files; pre-existing untracked file edited;
     pre-existing untracked file deleted; empty directory created;
     pre-existing empty directory removed; tracked file made executable;
     untracked file made executable; tracked symbolic link retargeted;
     tracked file replaced by a symbolic link; tracked file permission
     narrowed (0o600, wants `modified: tracked.txt`); tracked directory
     permission narrowed (`sub` to 0o700, wants `modified: sub`); new
     directory holding only ignored files (`onlylogs/x.log`, wants
     `created: onlylogs`); file matching a negated ignore rule created
     (`keep.log`); and, passing, file under a directory-only ignore rule
     created (`ignored/build.log`), file matching a root ignore pattern
     created (`sub/trace.log`) and file matching a nested ignore file
     created (`sub/local.tmp`);
   - `TestGuardLeavesNoRepositoryPathInTheTestInputLog`: a guarded fixture
     run with `-test.testlogfile` whose log is non-empty and names no `open`
     or `stat` under the fixture repository, in either its path or its
     resolved form;
   - `TestGuardRefusesARootOutsideAGitWorkTree`: a guarded run on a plain
     temporary directory exits non-zero and prints
     `suiteguard: fingerprint repository before tests:`.
3. MUST replace `fingerprint` in `internal/suiteguard/suiteguard.go`. It
   takes only the root and returns one state per watched path:
   - `git ls-files -z --cached --others --exclude-per-directory=.gitignore`
     lists tracked files and untracked files that the repository's
     `.gitignore` files do not ignore; `--exclude-standard` is not used, so
     `.git/info/exclude` and the user's global excludes file change nothing,
     as today;
   - `git ls-files -z --others --directory --exclude-per-directory=.gitignore`
     supplies untracked directories, empty or not (entries ending in `/`);
   - the listed files go, NUL-separated on standard input, to
     `git update-index --add --remove --info-only -z --stdin` with
     `GIT_INDEX_FILE` set to a file in a fresh `os.MkdirTemp` directory
     outside the repository, removed afterwards, and
     `git ls-files -z --stage` on that index gives each present path's Git
     mode and object ID. `--info-only` writes no object and `--remove`
     drops a tracked path the working tree lost;
   - every parent directory of a hashed path and every untracked directory
     is a directory entry;
   - one `os.Lstat` per entry records the full `fs.FileMode`, so permission
     bits, the executable bit and the file type keep being compared;
   - every `git` call runs with `-C <root>`, an environment without any
     inherited `GIT_` variable plus `GIT_OPTIONAL_LOCKS=0`, and
     `-c core.fileMode=true -c core.symlinks=true -c core.autocrlf=false
     -c core.safecrlf=false -c core.fsmonitor=false`; a failing call
     returns an error naming the git subcommand and its standard error,
     which `run` reports as today, so a root outside a Git work tree fails
     the package.
4. MUST remove `loadIgnoreRules`, the hand-written ignore parser and its
   helpers, and `digestFile`, and keep `Main`, `run`'s messages (including
   `guard cost: measured %d paths`), `compare`, the Sanctioned regeneration
   handling and the helper residue check byte-identical in behavior.
5. MUST make every existing guarded fixture a Git work tree: in
   `internal/suiteguard/suiteguard_test.go` each `t.TempDir()` used as a
   guarded root becomes a fresh Git repository, and `runGuardedFixture`
   forwards the case variable; in `internal/suiteguard/residue_test.go` the
   helper fixture's root becomes an initialized Git work tree. Their
   assertions keep their bytes.
6. MUST NOT change `internal/suiteguardcontract`, any caller of
   `suiteguard.Main`, `internal/cli`, the Makefile, the CI workflows or any
   file task_01, task_02 or task_04 changed.

## Subtasks

- [ ] Add the leak-case table, the input-log test and the refusal test, and characterize the current guard.
- [ ] Fingerprint through Git with a mode-only Lstat pass, and remove the ignore parser.
- [ ] Make the existing fixtures Git work trees.

## Acceptance Criteria

- [ ] All 22 leak cases pass before and after the change; the refusal test fails before and passes after.
- [ ] Every guarded package still passes on this repository, and the guard still builds for Windows.

## Consumers

None. `suiteguard.Main` keeps its signature and messages; every package whose `TestMain` installs it keeps its result on a Git checkout.

## Context

- instruction: `internal/suiteguardcontract/contract.go`
- instruction: `internal/suiteguardcontract/regeneration.go`
- interface: `internal/suiteguard/suiteguard.go`
- interface: `internal/suiteguard/suiteguard_test.go`
- interface: `internal/suiteguard/residue_test.go`
- creates: `internal/suiteguard/leak_cases_test.go`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestGuardReportsEveryLeakClass|TestGuardLeavesNoRepositoryPathInTheTestInputLog|TestGuardRefusesARootOutsideAGitWorkTree|TestGuardNamesTheViolatingPath|TestGuardPassesOnAnIsolatedWrite|TestGuardIgnoresRepositoryRules|TestSanctionedRegenerationIsDeclaredInProcess|TestGuardReportsAndKillsHelperResidue|TestGuardPassesWhenEveryHelperEnded|TestGuardMarksItsTestBinary|TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning)$' ./internal/suiteguard 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGuardReportsEveryLeakClass TestGuardLeavesNoRepositoryPathInTheTestInputLog TestGuardRefusesARootOutsideAGitWorkTree TestGuardNamesTheViolatingPath TestGuardIgnoresRepositoryRules TestSanctionedRegenerationIsDeclaredInProcess; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; cases="$(printf '%s\n' "$out" | grep -c -- '--- PASS: TestGuardReportsEveryLeakClass/')"; test "$cases" -eq 22 || { printf 'leak cases passed: %s, want 22\n' "$cases" >&2; exit 1; }; if grep -q 'func loadIgnoreRules' internal/suiteguard/suiteguard.go; then echo 'the hand-written ignore parser is still present' >&2; exit 1; fi; GOOS=windows go build -buildvcs=false ./internal/suiteguard && go test -count=1 -timeout 15m -tags repocontract ./internal/suiteguard ./internal/suiteguardcontract` — expected: exit 0. Before this Task the three new tests do not exist, so the PASS check fails. With the guard unchanged and only the tests added, the refusal test fails. Removing the mode pass fails the two permission cases, dropping `--directory` fails the three directory cases, dropping `--others` fails seven untracked cases, and dropping `--remove` fails the tracked deletion case. Afterwards all named tests and all 22 cases pass, the parser is gone, the guard builds for Windows, and the repository contracts pass.

## References

- `_prd.md` → Sources; Goals; Core Feature 5; Success Metric 6
- `_techspec.md` → Current behavior; API Contract 3; Invariant 4; Testing Approach; Build Order 4

## Result

### Characterization before the fingerprint change

- Added the 22-case leak table, the input-log check and the non-Git-root
  refusal check, and initialized the existing guarded fixture roots with Git.
- Ran `go test -count=1 -timeout 3m -v -run 'TestGuardReportsEveryLeakClass|TestGuardLeavesNoRepositoryPathInTheTestInputLog|TestGuardRefusesARootOutsideAGitWorkTree' ./internal/suiteguard`
  with `suiteguard.go` unchanged (`git diff` for that file was empty).
  All 22 leak subtests and the input-log test passed. Only
  `TestGuardRefusesARootOutsideAGitWorkTree` failed: the original guard
  accepted the plain directory with exit 0. Overall exit 1, as expected.

### Implementation and focused evidence

- `fingerprint(root)` now lists files and untracked directories with Git's
  per-directory `.gitignore` handling. Git hashes working-tree content into
  a disposable external index with `--info-only --remove`; stage records
  supply the Git mode and object ID. Parent directories and untracked
  directories receive entries, and one `os.Lstat` per entry records its full
  mode. The hand-written ignore parser and `digestFile` were removed.
- Each Git subprocess clears inherited `GIT_` variables, sets
  `GIT_OPTIONAL_LOCKS=0`, and applies all five required core settings.
  Git errors name the subcommand and include stderr. The temporary index
  directory is removed on return.
- Existing fixture assertions remain unchanged. Their guarded roots are
  initialized Git repositories. The new fixture commits the required files
  and symbolic link, then leaves dirty, untracked and empty-directory state.
- Acceptance criterion 1: the same focused characterization command recorded
  above exited 0 after the implementation. All 22 leak subtests, the input-log
  test and the non-Git refusal test passed. Both before/after logs were
  counted explicitly: 22 passing leak subtests each. The empty-file case
  also asserts zero size inside the guarded child.
- Existing behavior check:
  `go test -count=1 -timeout 3m -run 'TestGuardNames|TestGuardPasses|TestGuardIgnores|TestSanctioned|TestGuardReportsAndKills|TestGuardMarks|TestHelperResidue' ./internal/suiteguard`
  exited 0. This covers diagnostics, ignored writes, isolated writes,
  sanctioned and refused regeneration, binary marks and helper residue.
  After isolating the residue fixture's Git initialization from inherited
  Git settings, `go test -count=1 -timeout 2m -run '^TestHelperResidueSkipsAnUnreadableProcessAndKeepsScanning$' ./internal/suiteguard`
  also exited 0.
- Acceptance criterion 2, portability evidence:
  `GOOS=windows go test -c -o /tmp/0313-suiteguard-windows.test.exe ./internal/suiteguard`
  exited 0, compiling both the package and its tests for Windows. Execution
  of every guarded package and repository-contract Verification remains
  Daemon-owned and was not run in this child turn.
- Compared `Main`, `compare`, and the sanctioned-regeneration functions with
  `HEAD`: their function text is unchanged. Helper-residue production code
  was untouched. `git diff --check` exited 0.
- Scope: only the four authorized `internal/suiteguard` files and this
  Task's Result were edited. The pre-existing Daemon status change is
  preserved. No manifest, other Task, caller, CLI, contract package, build
  configuration or workflow was changed; no commit, push or PR was made.
- Environment: the executor rejected the symlinked writable root, so shell
  commands used the authorized worktree's resolved physical path with
  sandbox escalation. No executor configuration was changed.

Declared Verification and repository-wide acceptance remain for the Daemon;
this Result records focused implementation evidence, not Task settlement.

## Carry-forward provenance

- Source Run: `run_20261025T193403Z_ade2bddbacb762cf`
- Source commit: `5807780c7cea58c41515a1e9af80b02f82175dca`
