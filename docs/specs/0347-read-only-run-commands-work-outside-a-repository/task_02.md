---
task: task_02
spec: 0347-read-only-run-commands-work-outside-a-repository
status: completed
type: backend
complexity: medium
---

# Task 02: The Read-Only Run Commands load without the worktree-location refusal and warn instead

## Overview

Bare `runs`, `runs list`, `runs show`, `attach` and `events` only read
the Run Database, yet a stray `.git` in the home directory makes each of them
exit `2` with the `worktree.location` refusal. This Task gives them one
loader that sets `WorktreeLocationOptional` and prints the refusal once on
stderr as a warning, proves their output is unchanged inside a repository and
under a stray `.git` directory or file, asserts Surface Transcripts 2 and 3,
and states the contract in the `attach` and `events` guides. It answers the
maintainer report of 2026-11-03, adopted as
[references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md](references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md).

## Requirements

1. MUST characterize first: create `internal/cli/read_only_run_commands_test.go`
   with `TestReadOnlyRunCommandsKeepTheirOutputInsideARepository`
   (`_techspec.md` → Testing Approach), run it green on the unchanged
   commands, and record that run in the Task Result.
2. MUST add `loadRunReaderCommandConfig(environment commandEnvironment, stderr
   io.Writer) (roundconfig.Loaded, error)` to `internal/cli/cli.go` as API
   Contract 3 states: it builds the same `LoadOptions` as
   `loadCommandConfig`, sets `WorktreeLocationOptional`, and when
   `WorktreeLocationErr` is set writes exactly one stderr line
   `roundfix: warning: <WorktreeLocationErr>; commands that create a Run
   Worktree refuse it` before returning. `loadCommandConfig` keeps its bytes.
3. MUST switch the config load of `runRunsListCommand` (`runs.go`),
   `runRunsShowCommand` (`runs_show.go`), `runAttachCommand` (`attach.go`),
   `runEventsCommand` (`events.go`) and `runRunsBrowserCommand`
   (`runbrowser.go`) to `loadRunReaderCommandConfig`, and no other command.
   The `runs list` scope rule is task_03's; do not change it here.
4. MUST add `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory` and
   `TestStrayGitDirectoryTranscripts` to the same file as `_techspec.md` →
   Testing Approach describes. `TestStrayGitDirectoryTranscripts` asserts
   Surface Transcripts 2 and 3 whole, one subtest each: every stdout and
   stderr line in order and the exit code, with `<home>` the test's
   temporary home. The Run Browser cases use the existing interactive-input,
   `ROUNDFIX_TUI`, `withBrowserAttachCockpit` and `withRunBrowserSession`
   seams. Every top-level test calls `t.Parallel()` first; no test reads the
   real home directory or starts a process that outlives it.
5. MUST state in `docs/user-guide/commands/attach.md` and
   `docs/user-guide/commands/events.md` that the command is a Read-Only Run
   Command: it runs from any directory, and when `worktree.location` would be
   refused there it prints that refusal once on stderr as a warning, naming the
   Git root and the `.git` it found, and continues.
6. MUST re-record the Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   as the Task's last step, after the last guide edit, and never hand-edit it.
7. MUST NOT change `internal/config`, the `runs` help text, any skill, the
   Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the five commands inside a repository.
- [ ] Add the loader and switch the five commands to it.
- [ ] Add the stray-`.git` and transcript tests.
- [ ] Update the attach and events guides, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] Inside a repository the five commands print what they printed before.
- [ ] Under a stray `.git` directory or file they exit and print as inside a repository, with the warning first on stderr; `reconcile` still refuses with the Git root named.
- [ ] Surface Transcripts 2 and 3 hold whole.

## Context

- instruction: `docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/runs.go`
- interface: `internal/cli/runs_show.go`
- interface: `internal/cli/attach.go`
- interface: `internal/cli/events.go`
- interface: `internal/cli/runbrowser.go`
- creates: `internal/cli/read_only_run_commands_test.go`
- interface: `docs/user-guide/commands/attach.md`
- interface: `docs/user-guide/commands/events.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The output of `runs`, `runs list`, `runs show`, `attach` and `events` — consumers: `./internal/cli`
- The Behavior Surface Record of the `attach` and `events` guides — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReadOnlyRunCommandsKeepTheirOutputInsideARepository|TestReadOnlyRunCommandsWarnOnAStrayGitDirectory|TestStrayGitDirectoryTranscripts)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReadOnlyRunCommandsKeepTheirOutputInsideARepository TestReadOnlyRunCommandsWarnOnAStrayGitDirectory TestStrayGitDirectoryTranscripts; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/cli` — expected: exit 0. Before this Task none of the three tests exists, so the first pass line is missing; afterwards they pass with the whole CLI suite.
- `for f in docs/user-guide/commands/attach.md docs/user-guide/commands/events.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Read-Only Run Command' || { echo "missing contract in $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither guide names the contract; afterwards both do and the documentation contracts, the Behavior Surface fingerprint among them, pass.

## References

- `_prd.md` → Goals; Core Feature 3; Success Metric 2; Success Metric 4
- `_techspec.md` → Current behavior 3; API Contract 3; API Contract 5; Surface Transcript 2; Surface Transcript 3; Invariants 2 and 3; Declared breaks; Testing Approach; Build Order 2


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, knowledge-workspace, evidence-gate.

Implemented the reader-only loader and switched exactly `runs list`, `runs
show`, `attach`, `events`, and the bare `runs` browser to it. It uses the same
load options as the strict loader plus `WorktreeLocationOptional`, emits the
specified warning once, and propagates configuration and warning-write errors.
`loadCommandConfig` is byte-identical; repository scoping and `runs` help are
unchanged. Both command guides now state the any-directory and warning contract.

Focused evidence by acceptance criterion:

- Inside-repository output: before changing production commands,
  `go test -count=1 -run '^TestReadOnlyRunCommandsKeepTheirOutputInsideARepository$' ./internal/cli`
  exited 0 (`ok roundfix/internal/cli`, 1.032s). The characterization exercises
  a resolved Run through show (text and JSON), events, explicit attach, the
  scoped and all-repository lists, and both browser entry points. Readers
  outside any repository match inside stdout and stderr exactly; the scoped
  list contains exactly the resolved Run's row. Browser selection, cockpit,
  refresh and cancellation use existing seams.
- Stray `.git` directory and file: before the loader change,
  `go test -count=1 -run '^TestReadOnlyRunCommandsWarnOnAStrayGitDirectory$' ./internal/cli`
  exited 1, reproducing exit 2 from every reader. After the change,
  `go test -count=1 -run '^(TestReadOnlyRunCommands|TestStrayGitDirectoryTranscripts)' ./internal/cli`
  exited 0 (0.781s). Both marker forms preserve inside-repository stdout and
  exit codes, with stderr exactly the warning followed by baseline stderr.
  A hidden-terminal-Run diagnostic proves ordering; the scoped list still
  reports no Runs for the stray root. Both browsers select the stored Run,
  return from the cockpit, and cancel with only the warning on stderr.
- Surface Transcripts 2 and 3: `TestStrayGitDirectoryTranscripts` compares the
  entire stdout, stderr and exit code in separate named subtests, substituting
  only the temporary home. The list warns and exits 0; `reconcile` retains
  its exact refusal and usage line and exits 2, naming the Git root and `.git`.
- Parallel execution and nearby contracts:
  `go test -race -count=1 -run '^(TestReadOnlyRunCommands|TestStrayGitDirectoryTranscripts|TestRunRunsListOutsideRepositoryRequiresAll$|TestAttachWithoutRunDatabaseFailsAsCLIError$|TestEventsValidationErrorsEmitNoStdout$)' ./internal/cli`
  exited 0 (2.815s). This also retains task_03's current outside-repository
  refusal and the existing missing-store and invalid-events contracts.
  Every new top-level test starts with `t.Parallel()`; fixtures use temporary
  homes and existing fake collaborators, without live Agent sessions or network.

`git diff --check` exited 0. The initial changed-path baseline contained only
this Task file's Daemon-owned status change; it is preserved. No other Task,
Task Graph, config package, skill, tooling, module or workflow file was edited.
The selected repository gates, full package suites and declared Verification
commands are reserved for the Daemon. No commit, push or Pull Request was made.

Behavior Surface Record: the sanctioned recorder is the final implementation
step after both guide edits; its outcome is recorded below.

`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0 (0.440s).
The generated diff contains only the two edited guide fingerprints. No guide
or implementation edit followed regeneration.
