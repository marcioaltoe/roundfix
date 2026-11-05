---
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
prd: _prd.md
created: 2026-11-05
---

# Only a .git entry that Git accepts marks a Git root — Technical Spec

## Executive Summary

`findGitRoot` (`internal/config/config.go`) returns the first ancestor of the
working directory where `os.Stat(<dir>/.git)` succeeds, so any file or
directory named `.git` makes a Git root. `findSkillsGitRoot`
(`internal/cli/cli.go`) holds a byte-identical copy for the project scope of
`skills install`. On the maintainer's MacBook a `~/.git` holding only `gk/`
and `info/` made the home directory a Git root, which Git itself rejects.

The fix replaces both walks with one exported `config.FindGitRoot` that
accepts a `.git` entry only when Git would: a directory that passes Git's own
repository test, or a regular file starting with `gitdir:`. Anything else is
skipped and the walk continues to the parent. The test reads the file system
and starts no `git` process (ADR-0362).

The trade-off accepted: about 80 test fixtures faked a repository with an
empty `.git` directory and now need a marker the new rule accepts. They get one
helper, `gittest.MarkRepository`, which writes the smallest `.git` directory
Git accepts (`HEAD`, `objects/`, `refs/`). Git then answers in those
checkouts as in an empty repository instead of failing, and the one fake that
depended on Git failing, the CLI tests' `fakeWorktree`, writes the change it
reports to disk.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0362 (new) records the rule:
  "A `.git` that Git rejects is skipped and the walk goes on to the parent, as
  Git's own discovery does." It amends ADR-0357, whose Decision says "How the
  Git root is found does not change: any `.git` entry, directory or file,
  still marks a repository." ADR-0357's refusal text and Read-Only Run Command
  contract stay. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ answer of 2026-11-05 and the cycle AY rules quoted in the authorization
  record; the Governed Paths are `internal/cli/cli_test.go` and
  `internal/gittest/gittest.go`, which task_01 changes, and the Roundfix skill files task_03 changes; no Makefile,
  workflow, module or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0352-only-a-git-entry-that-git-accepts-marks-a-git-root/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/review-runs.md`, `.agents/skills/roundfix/references/runs.md`, `internal/cli/cli_test.go`, `internal/gittest/gittest.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Read on 2026-11-05 at `38cc5854`; `_prd.md` → Validation records the runs.

1. `findGitRoot(start)` walks up from `start` and returns the first directory
   where `os.Stat(filepath.Join(dir, ".git"))` succeeds, or `""`. Its callers
   are `Load` (`GitRoot`, then `ProjectConfigPath` and
   `ResolveWorktreeLocation`), `Init` for `--scope project` and
   `profileConfigTarget` (`internal/config/profile_config.go`) for `profiles
   configure --scope project`.
2. `findSkillsGitRoot` in `internal/cli/cli.go` is the same walk;
   `defaultResolveSkillsProjectRoot` returns its answer, or falls back to `git
   rev-parse --show-toplevel` when it finds nothing.
3. Every later reader of `Loaded.GitRoot` (`RepositoryRoot`, the Run
   Database's repository key, `runs list` scope, `cancel`, `upgrade`, `deliver
   lane`, `follow-up`, `runs causes`, the artifact directory) takes the root
   from `Load`, so changing the walk changes them all, and only where the walk
   changes.
4. Other `.git` readers: `findBaselineRepositoryRoot`
   (`internal/cli/baseline_profile.go`) has its own walk that also refuses a
   `.git` that is neither a file nor a directory, and is followed by Git
   inspection; `internal/store/store.go`'s repository backfill and
   `worktree_build_cache.go` read the `.git` of a recorded Run path;
   `internal/spec`, `internal/speccheck`, `internal/worktree` and
   `internal/baseline` ask `git rev-parse`. This Spec leaves them as they are.
5. Tests: `grep -rnE '(mustMkdir|os.Mkdir)\(.*"\.git"\)' internal --include=*_test.go`
   finds the empty-`.git` fixtures; a probe that recorded every call where the
   old and the new walk differ, run over `go test ./...`, found them in the
   files listed under Declared breaks and nowhere else.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/config/config.go`: `FindGitRoot` and its unexported `.git` test
  replace `findGitRoot`; `Load`, `Init` and `profileConfigTarget` call it.
- `internal/config/profile_config.go`: the call.
- `internal/gittest/gittest.go`: `MarkRepository`.
- `internal/config/git_root_test.go` (new): the walk, Git's agreement and
  `Load` under a rejected `.git`.
- Test fixtures in `internal/config`, `internal/cli` and `internal/judge`:
  `gittest.MarkRepository` where an empty `.git` directory stood; the CLI
  tests' `fakeWorktree` writes the file it reports.
- `internal/cli/cli.go`: `defaultResolveSkillsProjectRoot` calls
  `roundconfig.FindGitRoot`; `findSkillsGitRoot` is deleted.
- `internal/cli/rejected_git_directory_test.go` (new): Surface Transcripts
  1 to 3 and the skills project root.
- Documentation: the Roundfix `review-runs` and `runs` references, the
  configuration and `runs` guides, and ADR-0357's amendment notice.

## Implementation Design

### Interfaces

```go
// internal/config/config.go

// FindGitRoot returns the nearest directory at or above start whose .git
// entry Git accepts (ADR-0362), or "" when there is none.
func FindGitRoot(start string) string

// gitEntryAccepted reports whether path, a .git entry, is a directory that
// passes Git's repository test or a regular file starting with "gitdir:".
func gitEntryAccepted(path string) bool

// internal/gittest/gittest.go

// MarkRepository writes the smallest .git directory Git accepts under dir: a
// HEAD naming a branch and empty objects and refs directories.
func MarkRepository(t testing.TB, dir string)
```

### Data Models

No schema changes.

### API Contracts

1. API Contract: the `.git` test. `gitEntryAccepted(path)` follows symbolic
   links (`os.Stat`) and returns true only for:
   - a regular file whose content starts with `gitdir:`, whatever its target;
   - a directory whose `HEAD`, read with `os.Lstat`, is a symbolic link whose
     target starts with `refs/`, or a file whose content starts with `ref:`
     followed by optional whitespace and `refs/`, or with a run of exactly 40
     or 64 hexadecimal digits; and whose common directory holds an `objects`
     directory and a `refs` directory. The common directory is the path in
     `<path>/commondir`, trimmed and resolved against `<path>` when relative,
     or `<path>` when that file is absent.
   Every other entry, and every read error, returns false.
2. API Contract: the walk. `FindGitRoot(start)` cleans `start` and, from it
   up to the file-system root, returns the first directory whose `.git` entry
   passes API Contract 1, or `""`. It starts no process and reads no
   environment variable. `Load`, `Init`, `profileConfigTarget` and
   `defaultResolveSkillsProjectRoot` call it, and no other copy of the walk
   remains in `internal/config` or `internal/cli`.
3. API Contract: unchanged messages. The `worktree.location` refusal, the
   Read-Only Run Command warning, `project init requires a Git root; use
   --scope user outside a repository`, `project profiles configure requires a
   Git root; use --scope user outside a repository`, `reconcile requires
   running inside a Git repository` and `project skill install requires
   running inside a git repository or passing --dir` keep their bytes; only
   which directories reach them changes.
4. API Contract: the test marker. `gittest.MarkRepository(t, dir)` creates
   `<dir>/.git/objects/` and `<dir>/.git/refs/` with `os.MkdirAll` and writes
   `<dir>/.git/HEAD` holding `ref: refs/heads/main\n`, failing the test on any
   error. `FindGitRoot` accepts it, and `git rev-parse --show-toplevel` there
   prints `<dir>`.
5. API Contract: the rule in writing.
   - The rule sentence, the same Markdown in four files:

     ```markdown
     Only a `.git` entry that Git accepts marks a Git root: a directory with a valid `HEAD` and `objects` and `refs` directories, or a file that starts with `gitdir:`. Any other `.git`, such as a directory that holds only `gk` and `info`, is skipped as Git skips it.
     ```

   - `.agents/skills/roundfix/references/review-runs.md`, at the
     `worktree.location` bullet, and `docs/user-guide/configuration.md`, next
     to the `worktree` block: the sentence "Repository detection finds any
     `.git` entry in an ancestor, so a stray `.git` in a parent such as the
     home directory is found this way." is replaced by the rule sentence.
   - `docs/user-guide/commands/runs.md`, in the Read-Only Run Command
     paragraph: "including a stray `.git` in a parent such as the home
     directory." becomes "including one in a parent such as the home
     directory." followed by the rule sentence, and "A detected Git root,
     including a stray one, keeps" becomes "A detected Git root, including
     one in a parent directory, keeps".
   - `.agents/skills/roundfix/references/runs.md`: the rule sentence follows
     "Outside a Git repository, `runs list` lists every repository as `runs
     list --all` does and exits `0`."
   - `docs/adr/0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md`:
     after the Decision paragraph "How the Git root is found does not
     change...", a new paragraph; at the end of the Considered Option
     "Recognize a `.git` only when Git would...", one sentence; and
     `updated_at: 2026-11-05T12:00:00Z`:

     ```markdown
     **Amended by ADR-0362 (2026-11-05).** A `.git` entry marks a Git root only when Git accepts it: a directory that passes Git's repository test or a file starting with `gitdir:`. A stray `.git` that Git rejects is skipped. The refusal text and the Read-Only Run Command contract still hold.
     ```

     In the ADR file, `ADR-0362` in that bold opening is the Markdown link
     whose target is `0362-only-a-git-entry-that-git-accepts-marks-a-git-root.md`,
     the form of every amendment notice under `docs/adr/`.

     ```markdown
     ADR-0362 takes the first half of this option, recognizing a `.git` only when Git would.
     ```

   Line wrapping may differ; the words may not.

### Surface Transcripts

Each runs with `HOME` set to a fresh temporary directory `<home>` holding
`<home>/.git/gk/` and `<home>/.git/info/` and nothing else, no User Config and
no Run Database, from the working directory `<home>/drive`. After each,
`<home>/.roundfixrc.yml` does not exist.

1. Surface Transcript: `runs list` finds no Git root and prints no warning.

   ```transcript
   $ roundfix runs list
   stdout:
   No Runs found.
   stderr:
   exit: 0
   ```

2. Surface Transcript: `reconcile` no longer refuses `worktree.location`; it
   refuses the missing repository.

   ```transcript
   $ roundfix reconcile
   stdout:
   stderr:
   roundfix: reconcile failed: reconcile requires running inside a Git repository
   Run 'roundfix reconcile --help' for usage.
   exit: 2
   ```

3. Surface Transcript: `init --scope project` refuses instead of writing
   Project Config into the home directory.

   ```transcript
   $ roundfix init --scope project
   stdout:
   stderr:
   Init failed

   Reason:
     project init requires a Git root; use --scope user outside a repository

   Usage:
     Run 'roundfix init --help' for usage.
   exit: 2
   ```

### Invariants

```text
1. For every .git shape in Testing Approach where git rev-parse --show-toplevel answers with a top level or with "not a git repository", FindGitRoot answers the same directory or "".
2. A repository made by git init, a linked worktree, a submodule or any other gitdir: file, and a real repository nested under a .git that Git rejects each resolve to the same Git root as at 38cc5854, and Project Config is read from that root.
3. Every message in API Contract 3 keeps its bytes.
4. Inside a Git repository Git accepts, every command prints what it printed at 38cc5854.
```

### Declared breaks

- A `.git` directory that Git rejects no longer marks a Git root, so under one
  no Project Config is read from that directory, `runs list` prints no
  warning and lists the Runs of every repository, `reconcile`, `stop`, `gc`,
  `implement` and `deliver` no longer refuse `worktree.location` (each then
  answers as outside a repository), and `init --scope project` and `profiles
  configure --scope project` refuse instead of writing into it.
- A regular `.git` file without `gitdir:` no longer marks a Git root.
- Tests that faked a repository with an empty `.git` directory call
  `gittest.MarkRepository` instead, with no assertion changed:
  - `internal/config`: `config_test.go`, `delivery_check_timeout_test.go`,
    `jev_ceiling_test.go`, `pre_pr_review_omit_paths_test.go`,
    `profile_deviation_test.go`, `profile_proof_setup_timeout_test.go`,
    `profiles_characterization_test.go`, `review_source_scope_test.go`,
    `run_retention_test.go`, `run_slot_wait_test.go`,
    `secondbrain_test.go`, `settlement_config_test.go`,
    `worktree_caches_test.go`, `worktree_shared_caches_test.go`;
  - `internal/cli`: `cli_test.go` (including `withCLIWorkspace`),
    `doctor_test.go`, `runs_list_outside_repository_test.go`;
  - `internal/judge`: `task_acceptance_measure_test.go`,
    `task_acceptance_stats_test.go`.
- The Spec 0347 tests that put a `.git` directory at the Git root,
  `internal/config/worktree_location_test.go` (rows without `gitFile`) and
  `internal/cli/read_only_run_commands_test.go` (`directory` and
  `TestStrayGitDirectoryTranscripts`), create it with
  `gittest.MarkRepository`, a home that is a repository Git accepts such as a
  dotfiles checkout; their `.git` file rows keep `gitdir: <dir>`. Their
  assertions and expected text are unchanged.
- `fakeWorktree.Snapshot` in `internal/cli/cli_test.go` reported
  `src/agent-change.go` without creating it, and the commit boundary kept the
  path only because `git ls-files` failed against the empty `.git`
  (`indexTracksExactPath`, `internal/daemon/task_engine.go`). Under a
  `.git` Git accepts, Git answers that the path is neither on disk nor in the
  index and the path is dropped, so `TestRunResolveNoAgentConsoleSuppressesAgentDisplayOnly`,
  `TestRunResolveFinalPushRunsOnceAfterAllUnresolvedTerminal`,
  `TestRunResolveProcessesAllBatchesBeforeFinalPush`,
  `TestWatchFinalPushRunsOncePerCleanRoundThroughEngine`,
  `TestAttachRendersWatchDaemonEventsInTimeline` and
  `TestRunOperationalCommandAcceptsMVPFlags` lost their "Batch commit
  created" line in the prototype. The fake now writes the file under the work
  directory before reporting it, and those tests pass unchanged.
- No test is renamed or removed, so the coverage record is not re-recorded.
- `findSkillsGitRoot` is deleted; no test names it
  (`grep -rn findSkillsGitRoot internal`: only its definition and its one
  call).

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 2 → API Contracts 1 and 2,
  `TestFindGitRootAcceptsOnlyAGitEntryGitAccepts`,
  `TestFindGitRootAgreesWithGit`, Invariant 1.
- Goal 2, Success Metric 3 → API Contracts 1 and 2 (a `gitdir:` file and a
  directory that passes Git's test still count), Testing Approach rows `git
  init repository`, `linked worktree`, `gitdir file` and `repository under a
  rejected directory` of `TestFindGitRootAcceptsOnlyAGitEntryGitAccepts`,
  the existing `internal/config` Project Config hierarchy tests, which read
  Project Config from a `gittest.MarkRepository` root, and the Spec 0347
  tests, Invariants 2 and 4.
- Goal 3, Success Metric 1 → API Contracts 2 and 3,
  `TestLoadUnderARejectedGitDirectoryFindsNoGitRoot`,
  `TestRejectedGitDirectoryTranscripts` (Surface Transcripts 1 to 3),
  Invariant 3.
- Goal 4, Core Feature 2, Success Metric 4 → API Contract 2,
  `TestSkillsProjectRootSkipsARejectedGitDirectory`.
- Core Feature 3 → API Contract 4 and Declared breaks.
- Goal 5, Core Feature 4, Success Metric 5 → API Contract 5.

Deviations from the PRD: None.

## Absence Evidence

1. `FindGitRoot` and a `.git` test that reads `HEAD`, `objects` or `refs` —
   command: `grep -rn 'FindGitRoot\|"HEAD"\|"objects"' internal/config internal/cli/cli.go`;
   result: no `FindGitRoot`; `internal/config` reads `.git` only through
   `os.Stat` in `findGitRoot` and `repositoryRootMetadata`, and neither
   checks `HEAD`, `objects` or `refs`.
2. `gittest.MarkRepository` — command: `grep -rn 'MarkRepository' internal`;
   result: no match.
3. The two new test files — command: `ls internal/config/git_root_test.go internal/cli/rejected_git_directory_test.go`;
   result: no such files.

## ADR Alternatives

- ADR-0362: the decision this Spec records; it rejects asking `git
  rev-parse`, a `HEAD`-and-`objects` test without `refs`, stopping at any
  regular `.git` file, skipping a `gitdir:` file whose target is missing, and a
  home-directory special case, and the design takes none of them.
- ADR-0357: this Spec takes the option ADR-0357 rejected for its own scope,
  "Recognize a `.git` only when Git would", so ADR-0362 amends it and task_03
  writes the notice; its refusal text, its warning and its five Read-Only Run
  Commands are unchanged.

## Integration Points

None external. The `git` binary is a Base Test Tool; the tests use it to make
real repositories and as the oracle of `TestFindGitRootAgreesWithGit`.

## Testing Approach

All tests are hermetic: `t.TempDir()` directories, `gittest.InitRepo` and
`gittest.Run` with the isolated Git environment, the existing
`setCommandEnvironmentForTest` and `runCLI` seams of `internal/cli`, and no
network. Every new top-level test calls `t.Parallel()` first, and each new
test file opens with the `Suite`, `Invariant`, `Boundary IN` and `Boundary
OUT` comment lines the package's test files carry.

- `internal/config/git_root_test.go` (task_01). One table of `.git` shapes,
  each built under a fresh `<base>` with the working directory below the
  shape's root:
  - rejected, expected root `""`: `only gk and info`, `empty directory`,
    `HEAD and objects without refs`, `HEAD that names nothing` (`junk`),
    `file without gitdir` (`garbage`);
  - accepted, expected root the shape's directory: `symbolic ref HEAD`,
    `detached HEAD` (40 hex digits), `SHA-256 detached HEAD` (64),
    `common directory` (`HEAD` plus `commondir` naming a directory with
    `objects` and `refs`), `gitdir file` (pointing at a `git init`
    repository), `gitdir file whose target is missing`, `git init
    repository`, `linked worktree` (`git worktree add`), `repository under a
    rejected directory` (`<base>/.git` holding `gk` and `info`, the working
    directory inside `<base>/drive/repo`).
  - `TestFindGitRootAcceptsOnlyAGitEntryGitAccepts`: every row.
  - `TestFindGitRootAgreesWithGit`: every row except `file without gitdir`
    and `gitdir file whose target is missing`, where Git fails with another
    message; `git rev-parse --show-toplevel` runs in the working directory
    with `gittest.IsolatedEnv()` and `GIT_CEILING_DIRECTORIES` set to the
    parent of the physical `<base>`, and both answers are compared after
    `filepath.EvalSymlinks`.
  - `TestLoadUnderARejectedGitDirectoryFindsNoGitRoot`: `Load` from
    `<home>/drive` under `<home>/.git/{gk,info}` returns no error, an empty
    `GitRoot` and `ProjectConfigPath` and the location
    `<home>/.roundfix/worktrees`; `Init` with `--scope project` returns the
    API Contract 3 refusal and writes no `.roundfixrc.yml`;
    `profileConfigTarget` returns its refusal.
- `internal/cli/rejected_git_directory_test.go` (task_02):
  - `TestRejectedGitDirectoryTranscripts`: Surface Transcripts 1 to 3 whole,
    one subtest each, plus the absence of `<home>/.roundfixrc.yml`.
  - `TestSkillsProjectRootSkipsARejectedGitDirectory`:
    `defaultResolveSkillsProjectRoot` from `<home>/drive` returns the API
    Contract 3 refusal, and from `<home>/drive/repo/sub` in a `git init`
    repository returns `<home>/drive/repo`.
- Existing suites: `./internal/config`, `./internal/cli`, `./internal/judge`,
  `./internal/gittest`, `./skills` and `-tags docscontract
  ./internal/docscontract`.
- Sabotages, each run by the author on a scratch copy of the prototype at
  `38cc5854`, with the case whose input makes it observable:
  - `gitEntryAccepted` back to "any entry": `only gk and info` fails in both
    walk tests, and `TestRejectedGitDirectoryTranscripts` fails on Surface
    Transcript 1 (the warning returns);
  - `refs` not required: `HEAD and objects without refs` fails;
  - `HEAD` content not checked: `HEAD that names nothing` fails;
  - `commondir` ignored: `common directory` fails;
  - a regular file accepted without `gitdir:`: `file without gitdir` fails;
  - a `gitdir:` file accepted only when its target exists: `gitdir file
    whose target is missing` fails;
  - `defaultResolveSkillsProjectRoot` on the old copy of the walk:
    `TestSkillsProjectRootSkipsARejectedGitDirectory` fails on its first
    assertion;
  - `gittest.MarkRepository` writing an empty `.git` directory again:
    `TestLoadAppliesConfigPrecedence` and every other fixture test that reads
    Project Config fail;
  - `fakeWorktree.Snapshot` no longer writing its file:
    `TestRunOperationalCommandAcceptsMVPFlags` fails on `resolve`.

## Build Order

1. `FindGitRoot` and the `.git` test in `internal/config`, `MarkRepository` in
   `internal/gittest`, the fixture changes of Declared breaks,
   `internal/config/git_root_test.go`, and the rule sentence in the
   configuration guide with the Behavior Surface Record re-recorded last
   (task_01).
2. The skills project root on `FindGitRoot`, deleting `findSkillsGitRoot`,
   and `internal/cli/rejected_git_directory_test.go` with Surface Transcripts
   1 to 3 (task_02) (depends on: 1).
3. The Roundfix `review-runs` and `runs` references, the `runs` guide,
   ADR-0357's amendment notice, with the mirrors synced, the owned-skill
   versions and the Behavior Surface Record recorded last (task_03) (depends
   on: 2).
4. QA gate (task_04) (depends on: 3).

task_01 must change every fixture in the same Task as the walk, because the
repository Verification runs at its settlement and the old fixtures fail under
the new walk; it also edits the configuration guide, so the CLI files it and
task_02 change have a guide in their dependency closure. task_02 needs `FindGitRoot` and changes `internal/cli/cli.go`,
whose tests task_01 also changes. task_03 documents what both make true and re-records the Behavior Surface
Record task_01 re-records. No
Verification depends on another Spec.

## Risks & Considerations

- A user whose home directory is a real repository Git accepts (a dotfiles
  checkout) still gets ADR-0357's warning and refusal; that is a repository,
  and `worktree.location` outside it is the documented answer.
- Under `gittest.MarkRepository` Git answers in the fixture checkouts instead
  of failing. The prototype ran the whole module with it: only the six tests
  of the `fakeWorktree` break above changed, so no other test depended on Git
  failing there.
- A filesystem that denies reading `HEAD` or `commondir` makes the entry
  rejected and the walk continue, as Git's `access` checks do.
- Git's own test reads `GIT_OBJECT_DIRECTORY`; Roundfix does not, as before,
  so a repository whose objects live only there resolves at `38cc5854` and
  after alike only when its `.git` also holds `objects/`.

## Glossary

None.

## Decisions

- One exported walk in `internal/config`, called by configuration and by the
  skills project root; no new package.
- The test is Git's `is_git_directory` and `validate_headref`, read from the
  file system; no `git` process (ADR-0362).
- A `gitdir:` file marks the root without checking its target; a file
  without `gitdir:` is skipped (ADR-0362).
- Fixtures use the smallest `.git` directory Git accepts, one helper for
  every package; a `gitdir:` file with a missing target was tried and broke
  19 tests that later run `git init` there, write inside `.git` or read Git's
  "not a git repository" text.
- `findBaselineRepositoryRoot`, the store backfill and the `git
  rev-parse` callers are out of scope (Current behavior 4).
- ADR-0362 records the rule and amends ADR-0357.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation found no page deciding this. Exa: gitrepository-layout(5) and
  gitcore-tutorial(7). Git's `setup.c` read from github.com/git/git.
- Three implementation Tasks in series and a QA gate.
