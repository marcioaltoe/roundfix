---
task: task_03
spec: 0347-read-only-run-commands-work-outside-a-repository
status: completed
type: backend
complexity: medium
---

# Task 03: Runs List outside a repository lists every repository

## Overview

Outside a Git repository `roundfix runs list` exits `2` and asks for
`--all`, although there is no narrower scope than every repository. This
Task makes it list every repository there, exactly as `runs list --all`,
keeps the repository scope wherever a Git root is found, updates the `runs`
help and guide, replaces the test that pinned the old refusal, and asserts
Surface Transcript 1. It answers the maintainer report of 2026-11-03, adopted
as
[references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md](references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md).

## Requirements

1. MUST replace, in `runRunsListCommand` (`internal/cli/runs.go`), the
   refusal `runs list requires a Git repository; pass --all to list every
   repository` with API Contract 4 of `_techspec.md`: when the Git root is
   `""`, the listing runs as with `--all` (rows, repository column, trailing
   note, exit `0`); a detected Git root keeps the repository scope.
2. MUST change the scope sentence of the `runs` help in
   `internal/cli/cli.go` to add `; outside a Git repository it lists every
   repository` after "shows the 20 newest Active Runs", keeping every other
   help line.
3. MUST remove `TestRunRunsListOutsideRepositoryRequiresAll` from
   `internal/cli/cli_test.go` and add
   `TestRunRunsListOutsideRepositoryListsEveryRepository` and
   `TestRunsListOutsideARepositoryTranscript` to the new
   `internal/cli/runs_list_outside_repository_test.go` as `_techspec.md` → Testing
   Approach describes; `TestRunsListOutsideARepositoryTranscript` asserts
   Surface Transcript 1 whole (stdout, stderr and exit code). Each new
   top-level test calls `t.Parallel()` first.
4. MUST state in `docs/user-guide/commands/runs.md`, next to the default
   scope, that outside a Git repository it lists every repository, as `--all`
   does, and that `runs`, `runs list` and `runs show` are Read-Only Run
   Commands that print a `worktree.location` refusal as one stderr warning.
5. MUST re-record the coverage record only with
   `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`,
   because a recorded top-level test is removed, and then the Behavior Surface
   Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last guide and help edit. Never
   hand-edit either record.
6. MUST NOT change `internal/config`, the loader task_02 added, any skill,
   the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] List every repository when there is no Git root.
- [ ] Update the runs help and guide.
- [ ] Replace the refusal test and assert Surface Transcript 1.
- [ ] Re-record the coverage record, then the Behavior Surface Record.

## Acceptance Criteria

- [ ] Outside a repository `runs list` with any other flags equals `runs list --all` and exits `0`; inside one, its scope is unchanged.
- [ ] Surface Transcript 1 holds whole, and the help and guide state the scope.

## Context

- instruction: `docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md`
- interface: `internal/cli/runs.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- creates: `internal/cli/runs_list_outside_repository_test.go`
- interface: `docs/user-guide/commands/runs.md`
- interface: `docs/references/coverage-record.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The `runs list` scope and the `runs` help — consumers: `./internal/cli`, `./internal/docscontract`
- The coverage record — consumers: `./internal/spec`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestRunRunsListOutsideRepositoryListsEveryRepository|TestRunsListOutsideARepositoryTranscript)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunRunsListOutsideRepositoryListsEveryRepository TestRunsListOutsideARepositoryTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; ! grep -q 'func TestRunRunsListOutsideRepositoryRequiresAll' internal/cli/cli_test.go || { echo 'old refusal test still present' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/cli && go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec` — expected: exit 0. Before this Task neither new test exists, so the first pass line is missing; afterwards they pass, the old test is gone, and the CLI suite and the coverage record agree.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/runs.md | grep -qF -- 'it lists every repository' || { echo 'missing scope in the runs guide' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide does not state the outside-repository scope; afterwards it does and the documentation contracts, the Behavior Surface fingerprint of the help and guide among them, pass.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 1
- `_techspec.md` → Current behavior 4-5; API Contract 4; API Contract 5; Surface Transcript 1; Invariant 4; Declared breaks; Testing Approach; Build Order 3

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

- Acceptance criterion 1: `runRunsListCommand` sets the existing all-repository
  option when no Git root is detected, sharing the query, repository column,
  filters, bounds and trailing note with explicit `--all`. A detected root
  retains its scope. The replacement test uses Runs from two repositories
  and compares stdout, stderr and exit code for default, all-state, terminal,
  bounded, unbounded and combined flags, including explicit `--all` outside
  a repository; it also asserts that an inside-repository listing excludes
  the other repository and its repository column.
- Acceptance criterion 2: `TestRunsListOutsideARepositoryTranscript` asserts
  Surface Transcript 1 in full. Before the production change, the focused
  command `go test ./internal/cli -run '^TestRunsListOutsideARepositoryTranscript$' -count=1`
  exited 1: the command under test returned exit 2, empty stdout and the old
  refusal. After the change, it returned exit 0, exactly `No Runs found.\n`
  on stdout and empty stderr. The help adds the specified scope sentence;
  the guide states the outside-repository scope and the single stderr
  warning for `runs`, `runs list` and `runs show`.
- Focused checks: `go test ./internal/cli -run '^(TestRunRunsListOutsideRepositoryListsEveryRepository|TestRunsListOutsideARepositoryTranscript|TestRunRunsListPrintsStableColumnsNewestFirst|TestRunRunsListStateFlagFiltersAndNotes|TestRunRunsListLimitBoundsNewestMatches|TestStrayGitDirectoryTranscripts)$' -count=1 -v`
  exited 0: all six top-level tests and their subtests passed. This covers
  the new behavior, existing repository output, filtering and limits, and
  the prior Task's stray-root transcripts. `git diff --check` exited 0.
- Removed only the superseded refusal test from `cli_test.go`. The required
  sanctioned command `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
  exited 0 and regenerated the coverage record. Its whole-tree inventory
  also includes previously unrecorded tests already present in this checkout;
  the only removed test name is the superseded refusal test. No generated
  record was edited by hand.
- Final required operation, after the last source, help, guide and Result edit:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
  This command owns the Behavior Surface Record regeneration; its process
  output supplies the regeneration result.
- Scope postflight: newly changed paths are the declared CLI source/tests,
  runs guide, two generated records and this Task file. The pre-existing
  Daemon status change remains intact. No config, loader, skill, Task Graph,
  other Task, Makefile, workflow or Go module edits. No follow-up work found.
- Declared Verification and repository Verification remain Daemon-owned and
  were not run. Task status remains `in_progress`; no terminal Task verdict,
  commit, push or pull request is claimed.
