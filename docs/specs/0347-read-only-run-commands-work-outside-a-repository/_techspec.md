---
spec: 0347-read-only-run-commands-work-outside-a-repository
prd: _prd.md
created: 2026-11-04
---

# Read-only Run commands work outside a repository — Technical Spec

## Executive Summary

Every command loads configuration through `config.Load`, which ends with
`ResolveWorktreeLocation` and fails when `worktree.location` lies inside the
Git root of the working directory. `findGitRoot` takes the first ancestor that
holds any `.git`, so a stray `.git` in the home directory captures the default
`~/.roundfix/worktrees`, and commands that never create a Run Worktree fail with
a refusal that names neither the root nor the `.git`.

The fix has three parts. The refusal appends the Git root and its `.git` path
after its current text. `LoadOptions` gains `WorktreeLocationOptional`: with it
`Load` returns the refusal in `Loaded.WorktreeLocationErr`, leaves the location
empty and succeeds. The five Read-Only Run Commands (bare `runs`, `runs list`,
`runs show`, `attach`, `events`) load through one CLI helper that sets the
option and prints the refusal as one stderr warning; and `runs list` with no
Git root lists every repository, as `--all` does.

The trade-off accepted: under a stray `.git`, `runs list` without `--all`
still scopes to that root and prints `No Runs found.`; the warning above it is
what tells the user why ([ADR-0357](../../adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md)).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0357 (new) records the Read-Only
  Run Command contract: "Every other command keeps refusing." The location and
  its refusal stay for every command that creates, inspects or removes a Run
  Worktree. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AX answers of 2026-11-04 quoted in the authorization record; the Governed
  Paths are the Roundfix skill files task_04 changes, `internal/cli/cli_test.go`
  and the coverage record task_03 changes; no Makefile, workflow,
  module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0347-read-only-run-commands-work-outside-a-repository/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/events.md`, `.agents/skills/roundfix/references/review-runs.md`, `.agents/skills/roundfix/references/runs.md`, `docs/references/coverage-record.json`, `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-04 at `6a0ed904`; `_prd.md` → Validation records the runs.

1. `findGitRoot` (`internal/config/config.go`) walks up from the working
   directory and returns the first directory holding a `.git` entry of any
   kind (`os.Stat`), or `""`.
2. `Load` resolves the User and Project Config, validates, then calls
   `ResolveWorktreeLocation(location, gitRoot, homeDir)`; when the resolved
   location is inside or equal to the Git root it returns
   `worktree.location must not be inside the repository tree: "<location>"`,
   and `Load` fails with it. `worktree_build_cache.go` and
   `doctor_disk_space.go` call `ResolveWorktreeLocation` directly too.
3. `loadCommandConfig` (`internal/cli/cli.go`) is the loader of
   `runRunsListCommand`, `runRunsShowCommand`, `runAttachCommand`,
   `runEventsCommand` and `runRunsBrowserCommand`; none of them, nor the Run
   Browser loop and the cockpit they open, reads `Config.Worktree.Location`.
   `runStopCommand`, `reconcile`, `implement` and the Delivery Queue read it.
4. `runRunsListCommand` refuses a missing Git root without `--all` with
   `runs list requires a Git repository; pass --all to list every
   repository` (exit `2`); `TestRunRunsListOutsideRepositoryRequiresAll` pins
   that refusal. With `--all` it lists every repository and adds the
   repository column.
5. The `runs` help says the listing "is scoped to the current repository";
   the Roundfix `runs` reference says that outside a Git repository `runs
   list` without `--all` exits `2`.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/config/config.go`: the refusal text, the option and the field.
- `internal/config/worktree_location_test.go` (new): the refusal and the
  option.
- `internal/cli/cli.go`: `loadRunReaderCommandConfig` and the `runs` help.
- `internal/cli/runs.go`, `runs_show.go`, `attach.go`, `events.go`,
  `runbrowser.go`: the loader call; `runs.go` also the outside-repository
  scope.
- `internal/cli/read_only_run_commands_test.go` (new): the stray-`.git` tests
  and Surface Transcripts 2 and 3.
- `internal/cli/runs_list_outside_repository_test.go` (new): the
  outside-repository listing and Surface Transcript 1.
- The Roundfix `runs`, `events` and `review-runs` references, the `runs`,
  `attach`, `events` and configuration guides, and `CONTEXT.md`.

## Implementation Design

### Interfaces

```go
// internal/config/config.go
type LoadOptions struct {
	HomeDir string
	WorkDir string
	Stderr  io.Writer
	// WorktreeLocationOptional keeps a worktree.location refusal from failing
	// Load, for a Read-Only Run Command that never creates a Run Worktree.
	WorktreeLocationOptional bool
}

type Loaded struct {
	// ...existing fields...
	WorktreeLocationErr error // the refusal when the option was set; Location is then ""
}

// internal/cli/cli.go
func loadRunReaderCommandConfig(environment commandEnvironment, stderr io.Writer) (roundconfig.Loaded, error)
```

### Data Models

No schema changes.

### API Contracts

1. API Contract: the refusal. `ResolveWorktreeLocation` returns
   `fmt.Errorf("worktree.location must not be inside the repository tree: %q is inside the Git root %q (found %q)", resolved, gitRoot, filepath.Join(gitRoot, ".git"))`
   where it returned the first clause alone; the text before ` is inside` is
   byte-identical to today's whole message. Every other error of the function
   keeps its bytes, and every caller gets the new text.
2. API Contract: optional location. With `WorktreeLocationOptional` set, a
   `ResolveWorktreeLocation` error does not fail `Load`: `Load` returns nil, the
   error in `Loaded.WorktreeLocationErr`, `Config.Worktree.Location` `""`, and
   every other field as without the option. Without the option, or when the
   location resolves, `Load` behaves as today and `WorktreeLocationErr` is nil.
3. API Contract: Read-Only Run Commands. `runRunsListCommand`,
   `runRunsShowCommand`, `runAttachCommand`, `runEventsCommand` and
   `runRunsBrowserCommand` load through `loadRunReaderCommandConfig`, which
   sets the option and, when `WorktreeLocationErr` is set, writes once to
   stderr `roundfix: warning: <WorktreeLocationErr>; commands that create a Run
   Worktree refuse it` before the command's own output. `loadCommandConfig`
   and every other command keep refusing.
4. API Contract: Runs List scope. When the Git root is `""`, `runs list`
   lists every repository exactly as `runs list --all` with the same other
   flags: rows, repository column, trailing note and exit `0`. A detected Git
   root, stray or not, keeps the repository scope. `--all` keeps its meaning
   everywhere.
5. API Contract: the rules in writing.
   - `runs` help: the scope sentence adds `; outside a Git repository it lists
     every repository`.
   - `.agents/skills/roundfix/references/runs.md`: the sentence "Outside a Git
     repository, `runs list` without `--all` exits `2` and names `--all` as the
     alternative." becomes the outside-repository listing; a paragraph names
     the Read-Only Run Commands and the warning; the closing sentence of the
     Run Browser paragraph says `runs list` is repository-scoped inside a
     repository.
   - `references/events.md`: `events` is a Read-Only Run Command and prints the
     warning on stderr, never on stdout.
   - `references/review-runs.md`, at the `worktree.location` bullet, and
     `docs/user-guide/configuration.md`, next to the `worktree` block: a
     location inside the current Git root is refused, the refusal names the
     root and the `.git` it found (a stray one in a parent directory such as
     the home directory is found this way), and Read-Only Run Commands only
     warn.
   - `docs/user-guide/commands/runs.md`, `attach.md`, `events.md`: the same
     facts for their commands.
   - `CONTEXT.md`: adds **Read-Only Run Command**.

### Surface Transcripts

Each runs with `HOME` set to a fresh temporary directory `<home>`, no User
Config, and no Run Database.

1. Surface Transcript: `runs list` from a directory outside any Git
   repository lists every repository, here none.

   ```transcript
   $ roundfix runs list
   stdout:
   No Runs found.
   stderr:
   exit: 0
   ```

2. Surface Transcript: `runs list` from `<home>/drive` with an empty `<home>/.git`
   directory warns and keeps the repository scope.

   ```transcript
   $ roundfix runs list
   stdout:
   No Runs found.
   stderr:
   roundfix: warning: worktree.location must not be inside the repository tree: "<home>/.roundfix/worktrees" is inside the Git root "<home>" (found "<home>/.git"); commands that create a Run Worktree refuse it
   exit: 0
   ```

3. Surface Transcript: `reconcile` from the same directory still refuses,
   naming the Git root and the `.git`.

   ```transcript
   $ roundfix reconcile
   stdout:
   stderr:
   roundfix: reconcile failed: worktree.location must not be inside the repository tree: "<home>/.roundfix/worktrees" is inside the Git root "<home>" (found "<home>/.git")
   Run 'roundfix reconcile --help' for usage.
   exit: 2
   ```

### Invariants

```text
1. The refusal text at 6a0ed904 is a prefix of the refusal after this Spec, for every location inside or equal to the Git root.
2. Only the five Read-Only Run Commands load with the option; every other command that loads configuration fails on the refusal as before.
3. Inside a Git repository whose tree does not hold worktree.location, the five commands print the same stdout, stderr and exit code as at 6a0ed904.
4. With a Git root, runs list keeps the repository scope; without one it equals runs list --all.
5. findGitRoot keeps accepting a .git directory and a .git file.
```

### Declared breaks

- `runs list` outside a Git repository exits `0` and lists every repository
  instead of exiting `2` with `runs list requires a Git repository; pass --all
  to list every repository`. `TestRunRunsListOutsideRepositoryRequiresAll`
  pins the old contract and is removed, replaced by
  `TestRunRunsListOutsideRepositoryListsEveryRepository`, so the coverage
  record is re-recorded.
- The `worktree.location` refusal gains text after its old message; the only
  test reading it (`TestLoadRejectsWorktreeLocationInsideRepository`) matches
  with `strings.Contains` and keeps passing
  (`grep -rn 'inside the repository tree' internal --include=*_test.go`).
- A Read-Only Run Command whose working directory has a Git root holding
  `worktree.location` succeeds with a warning instead of exiting `2`.
- The `runs` help text changes, so the Behavior Surface Record is
  re-recorded.

## Coverage Map

- Goal 1, Core Feature 4, Success Metric 1 → API Contract 4,
  `TestRunRunsListOutsideRepositoryListsEveryRepository`,
  `TestRunsListOutsideARepositoryTranscript` (Surface Transcript 1),
  Invariant 4.
- Goal 2, Core Features 2 and 3, Success Metric 2 → API Contracts 2 and 3,
  `TestLoadWithOptionalWorktreeLocationReturnsTheRefusal`,
  `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory`,
  `TestStrayGitDirectoryTranscripts` (Surface Transcript 2), Invariants 2
  and 5.
- Goal 3, Core Feature 1, Success Metric 3 → API Contract 1,
  `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit`,
  `TestStrayGitDirectoryTranscripts` (Surface Transcript 3), Invariants 1
  and 2.
- Goal 4, Success Metric 4 →
  `TestReadOnlyRunCommandsKeepTheirOutputInsideARepository`, Invariant 3.
- Goal 5, Core Feature 5, Success Metric 5 → API Contract 5.

Deviations from the PRD: None.

## Absence Evidence

1. `WorktreeLocationOptional` and `WorktreeLocationErr` — command: `grep -rn 'WorktreeLocationOptional\|WorktreeLocationErr' internal`; result: no match.
2. A loader that tolerates a worktree-location refusal — command: `grep -n 'ResolveWorktreeLocation' internal/config/config.go internal/cli/cli.go`; result: only the declaration and the one call in `Load`, whose error fails it; `cli.go` has none.
3. `loadRunReaderCommandConfig` — command: `grep -rn 'loadRunReaderCommandConfig' internal`; result: no match.
4. The three new test files — command: `ls internal/config/worktree_location_test.go internal/cli/read_only_run_commands_test.go internal/cli/runs_list_outside_repository_test.go`; result: no such files.
5. A Read-Only Run Command term — command: `grep -n 'Read-Only Run Command' CONTEXT.md docs/adr`; result: no match.

## ADR Alternatives

- ADR-0357: the decision this Spec records; it rejects silent skipping, every
  command lenient, a stricter `.git` recognition and `runs causes` as a
  Read-Only Run Command, and the design takes none of them.

## Integration Points

None external.

## Testing Approach

All tests are hermetic: temporary homes and repositories through
`t.TempDir()`, the existing `withCLIWorkspace`, `setCommandEnvironmentForTest`,
`runResolveForAttachTest`, `seedRunsForListColumns`, Run Browser and
interactive-input seams of `internal/cli`, and no network. Every new top-level
test calls `t.Parallel()` first.

- `internal/config/worktree_location_test.go`:
  - `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit`: rows `git
    directory`, `git file` (a `.git` file holding `gitdir: ...`) and `location
    is the root`, each with the working directory below the root; the error
    equals API Contract 1 exactly and starts with the old text.
  - `TestLoadWithOptionalWorktreeLocationReturnsTheRefusal`: under a stray
    `.git` the option returns no error, `WorktreeLocationErr` equal to the
    error `Load` returns without it, an empty location, the stray root and a
    User Config value still applied; in a repository whose tree does not hold
    the location, the option returns the same location and root as without it.
- `internal/cli/read_only_run_commands_test.go` (task_02):
  - `TestReadOnlyRunCommandsKeepTheirOutputInsideARepository`
    (characterization, green at `6a0ed904`): a resolved Run in a repository;
    `runs show`, `events` and `attach` exit `0` with no stderr and print the
    same output inside the repository and outside any repository; `runs list
    --state all` inside prints exactly one row for the Run and no stderr.
  - `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory`: subtests `directory`
    and `file`; from `<home>/drive`, `runs show`, `events`, `attach` and `runs
    list --all --state all` exit and print as inside the repository with
    stderr equal to the warning line plus the inside stderr; `runs list
    --state all` prints `No Runs found.` with the warning; the Run Browser
    through `attach` and bare `runs` exits `0` with the warning as its whole
    stderr.
- `internal/cli/runs_list_outside_repository_test.go` (task_03):
  - `TestRunRunsListOutsideRepositoryListsEveryRepository`: Runs of two
    repositories; `runs list`, `runs list --state all` and `runs list --limit
    1` outside any repository equal the same command with `--all` inside a
    repository, and name the other repository.
  - `TestRunsListOutsideARepositoryTranscript`: Surface Transcript 1 whole.
  - `TestStrayGitDirectoryTranscripts`: Surface Transcripts 2 and 3 whole, one
    subtest each, with `<home>` the test's temporary home.
- Existing suites: `./internal/config`, `./internal/cli`, `./skills` and
  `-tags docscontract ./internal/docscontract`.
- Sabotages, each run by the author on a scratch copy of the prototype at
  `6a0ed904`, with the case whose input makes it observable:
  - the refusal reverted to its old text:
    `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit` fails on `git
    directory`, and `TestStrayGitDirectoryTranscripts` on its third case;
  - the option ignored (`Load` always fails), the error dropped, or the
    unresolved location kept: `TestLoadWithOptionalWorktreeLocationReturnsTheRefusal`
    fails;
  - a `.git` file no longer marks a root:
    `TestWorktreeLocationRefusalNamesTheGitRootAndItsDotGit` fails on `git
    file`;
  - `attach`, `events`, `runs show` or the Run Browser back on
    `loadCommandConfig`: `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory`
    fails on that command; `runs list` back on it:
    `TestStrayGitDirectoryTranscripts` fails on its second case;
  - the warning not printed: `TestReadOnlyRunCommandsWarnOnAStrayGitDirectory`
    fails;
  - the option set in `loadCommandConfig` for every command:
    `TestStrayGitDirectoryTranscripts` fails on `reconcile`;
  - `runs list` without a Git root left without the repository column or
    still scoped: `TestRunRunsListOutsideRepositoryListsEveryRepository`
    fails;
  - the old refusal restored for a missing Git root:
    `TestRunsListOutsideARepositoryTranscript` fails.

## Build Order

1. The refusal text and the optional location in `internal/config`, with
   their tests (task_01).
2. The Read-Only Run Command loader in the five commands, the stray-`.git`
   tests and Surface Transcripts 2 and 3, the `attach` and `events` guides,
   with the Behavior Surface Record re-recorded last (task_02) (depends on: 1).
3. `runs list` outside a repository, the `runs` help, Surface Transcript 1,
   the removal of the old refusal test, the `runs` guide, with the coverage
   record and the Behavior Surface Record re-recorded last (task_03)
   (depends on: 2).
4. The Roundfix `runs`, `events` and `review-runs` references, the
   configuration guide and `CONTEXT.md`, with the mirrors synced, the
   owned-skill versions and the Behavior Surface Record recorded last
   (task_04) (depends on: 3).
5. QA gate (task_05) (depends on: 4).

task_02 needs task_01 for the option; task_03 changes `runs.go` and `cli.go`
after task_02 and re-records the Behavior Surface Record task_02 re-records; task_04 documents what task_02 and task_03
built and shares the Behavior Surface Record with both. No Verification
depends on another Spec; 0345 and 0348 merge first and may move the Roundfix
skill version, the Behavior Surface Record and `CONTEXT.md`, which are
re-recorded or edited narrowly here.

## Risks & Considerations

- Under a stray `.git`, `runs list` without `--all` lists that root's Runs,
  usually none; the warning names the root and `.git`, and `--all` still lists
  everything.
- A home directory that really is a repository (a dotfiles checkout) gets the
  warning on every Read-Only Run Command until `worktree.location` is set
  outside it, which `implement` already requires.
- Git itself does not take an empty `.git` for a repository; Roundfix does.
  Aligning detection with Git is left for a later Spec because every command
  depends on it.

## Glossary

- adds: **Read-Only Run Command** — a command that only reads the Run Database and runs from any directory

## Decisions

- One option on `config.Load` and one CLI loader, used by exactly five
  commands; no new command or flag.
- The refusal grows after its old text, so a reader matching the old message
  keeps matching.
- The warning is printed, not suppressed, and only once per command.
- `runs list` outside a repository equals `--all`; no extra note, because the
  repository column already shows the scope.
- `findGitRoot` is unchanged: a `.git` file of a linked worktree or submodule
  must keep marking a repository, and Git's own validity test would change
  detection for every command.
- ADR-0357 records the contract; no other ADR changes.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found no page deciding this. Exa: git(1) on repository discovery
  and the Pro Git environment variables chapter.
- Four implementation Tasks in series and a QA gate.
