---
spec: 0290-the-pull-request-record-commit-fits-commitlint
prd: _prd.md
created: 2026-10-20
---

# The Pull Request Record Commit fits commitlint — Technical Spec

## Executive Summary

`RecordPullRequest` in `internal/cli/deliver_workflow.go` commits the record
change with `"docs: record pull request #"+number+" in the "+specSlug+"
Archive Record"`, which grows with the slug. A new `pullRequestRecordSubject`
builds the subject from the Spec number instead, at most 70 characters for any
slug. The record proof (`pullRequestRecordCommitIsExact`) already compares the
change and never the message, so commits made by hand under any subject stay
adopted; a test pins that.

The three `git commit` calls the Delivery Queue makes (the archive commit, the
Pull Request Record Commit and the derived merge commit) go through one helper
that, when Git exits non-zero, classifies the failure with the Daemon's
`daemon.ClassifyCommitHookRefusal`. A hook refusal returns a typed error whose
text names the hook, the refused subject, the exit status, the hook's output
and the recovery. The Delivery Engine already parks any step error as
`delivery-error: <error>`, so the park reason carries that text with no engine
change, and the Park Class stays `environment`. A hook that does not name
itself is named from the installed hooks, which the Daemon package already
reads and now exports as `daemon.InstalledCommitHooks`.

The trade-off accepted: the reason stays under the `delivery-error` blocker
rather than a new one, so the Park Class table and its readers do not move,
and the next action reads from the reason text.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories,
  homes and a local bare remote. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 1, 2 and 3 follow
  ADR-0302, named with its title in `_prd.md` → Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AM approval of 2026-10-20 ("the same rules as cycle AL"), carrying the cycle
  AJ grant "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0290-the-pull-request-record-commit-fits-commitlint/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-20 at `70c78d15`.

1. `internal/cli/deliver_workflow.go:1352` commits the staged record with
   `workflow.git.RunGit(ctx, workDir, "commit", "-m", "docs: record pull
   request #"+number+" in the "+specSlug+" Archive Record")` and returns
   `commit pull request record: <git error>` on failure. The subject is 50
   characters plus the slug and the pull request number: 107 for a 54-character
   slug and pull request 280.
2. `preflight.ExecGitRunner.RunGit` (`internal/preflight/preflight.go:133`)
   returns `git <args joined by spaces>: <stderr, else stdout, trimmed>: <exit
   error>`, wrapping the `*exec.ExitError` with `%w`. The step's error does not
   separate the hook's output from the arguments, and the arguments hold the
   slug.
3. Git 2.56.0, in a temporary repository whose executable `commit-msg` hook
   refuses a header longer than 100 characters: `git commit -m` with the
   107-character subject of the Inbox Entry exited 1, printed only the hook's
   own line `header must not be longer than 100 characters, current length is
   107 [header-max-length]`, and left the change staged (`git status
   --porcelain` printed `A  a`). The same staged change committed with the
   57-character subject `docs: record pull request #280 in the 0088 Archive
   Record`.
4. The Delivery Engine wraps the step as `record pull request: %w`
   (`internal/delivery/engine.go:1264`) and parks every step error as
   `BlockerDeliveryError + ": " + err.Error()` (`engine.go:741`).
   `ClassifyPark` (`internal/delivery/park_class.go`) gives `delivery-error`
   the `environment` class and the generic next action `resolve the blocker,
   then run roundfix deliver retry <slug>`.
5. `pullRequestRecordCommitIsExact` (`deliver_workflow.go:1364`) accepts a
   head whose single parent is the candidate and whose tree differs only in
   the record by the `pull_request` field, through
   `pullRequestRecordTreeIsExact` and `pullRequestRecordChangeIsExact`. None of
   the three reads the commit message. `RecordPullRequest` reaches it when the
   item worktree is clean at a head other than the candidate, which is the
   state the oraculum operator's hand commit left.
6. After the refusal, the record change stays staged at the candidate head. A
   retry finds the worktree at the candidate with that one staged change,
   proves it again and commits it again with the same subject, so it fails
   the same way.
7. The archive commit, `"docs: archive "+specSlug` (`deliver_workflow.go:1227`),
   fails the same way and leaves the archive staged; a later step accepts a
   hand commit of it through `reconcileArchiveCommit`, which proves the change
   with `archiveCommitIsExact` and does not read the message either. The
   derived merge commit (`deliver_workflow.go:2503`, `"chore: merge default
   branch and regenerate derived paths\n\nRoundfix-Delivery: derived-merge"`)
   runs the same hooks; on failure the deferred cleanup restores the worktree
   and aborts the merge, so nothing stays staged. Nothing in the repository
   reads that trailer or either subject back.
8. `daemon.ClassifyCommitHookRefusal` (`internal/daemon/commit_hook.go:106`)
   returns the hook a failed commit's output names (`prepare-commit-msg`,
   `pre-merge-commit`, `pre-commit` or `commit-msg`, case-insensitive) or
   reports a refusal when the output carries a hook-runner marker (`husky`,
   `lefthook`, `.git/hooks/`, …) or the repository holds an executable
   `pre-commit` or `commit-msg` hook in the directory Git resolves
   (`commitHookInstalled`, unexported). The Daemon uses it only for Task
   commits.

Sweep of every `git commit` Roundfix makes, at `70c78d15` (`rg -n
'"commit",' internal -g '!*_test.go'`): the three queue commits above; the
Daemon's Task commit (`internal/daemon/daemon.go:244`), already classified as
a Hook Refusal; and Task Carry-Forward's `commit --amend --no-edit`
(`internal/cli/carryforward.go:474`), which runs with an empty hooks directory
(`carryForwardStagingGitRaw`) and so meets no hook. Only the three queue
commits are in this Spec's class.

Every other message the queue writes, with its length for the longest slug
the queue can deliver: Roundfix enforces no slug length, the Archive Record
`<slug>.md` bounds a slug at 252 bytes, and the item branch
`roundfix/deliver-<slug>-<16 hex>` bounds a deliverable slug at 230 bytes,
where its last ref component reaches the common 255-byte file name limit.

- `docs: archive <slug>`: 14 characters plus the slug, 244 at 230; over 100
  for a slug longer than 86 characters.
- The pull request title `feat: deliver <slug>` (`deliver_workflow.go:1571`):
  14 plus the slug, 244 at 230, and a squash merge appends ` (#<number>)`;
  over 100 for a slug longer than 86 characters. This repository's squash
  commit `9e92ac39` has a 117-character subject.
- The derived merge subject: 56 characters, and its trailer line 32, whatever
  the slug.
- The pull request body `Delivers Spec <slug>.` is not a header.

The first two can exceed 100 characters and are out of this Spec's decision;
this planning change writes the open Backlog Entry
`docs/backlog/2026-10-20-the-queue-archive-subject-and-pull-request-title-grow-past-100-characters-with-the-slug.md`
for them.

## System Architecture

- `internal/cli/deliver_workflow.go`: `pullRequestRecordSubject`; a queue
  commit helper used by `Archive`, `RecordPullRequest` and `ResolveConflict`;
  the typed refusal error.
- `internal/daemon/commit_hook.go`: `InstalledCommitHooks`, exported from the
  existing `commitHookInstalled` logic, which keeps its behavior.
- `docs/user-guide/commands/deliver.md`, `.agents/skills/roundfix/references/deliver.md`
  and its mirror: the subject and the recovery.
- `CONTEXT.md`: the Hook Refusal entry covers a Delivery Queue commit.

`internal/cli` already imports `internal/daemon`, so no import edge is added.

## Implementation Design

### Interfaces

```go
// pullRequestRecordSubject returns the Pull Request Record Commit subject.
func pullRequestRecordSubject(specSlug, number string) string

// InstalledCommitHooks lists the executable pre-commit and commit-msg hooks of
// workDir's repository, in that order, from the hooks directory Git resolves.
func InstalledCommitHooks(ctx context.Context, workDir string) []string
```

### Data Models

None. The Archive Record, the Delivery Queue item and its blocker field keep
their shapes; only the record commit's message and the text after
`delivery-error: ` change.

### API Contracts

1. API Contract: `pullRequestRecordSubject(specSlug, number)`. When the slug
   starts with one to ten ASCII digits followed by `-`, it returns `docs:
   record pull request #<number> in the <digits> Archive Record`; otherwise
   `docs: record pull request #<number> in the Archive Record`. With a pull
   request number the queue accepts (`^[1-9][0-9]{0,9}$`) the result is at
   most 70 characters. `0088-as-lojas-…` with `280` gives the 57-character
   subject of the Inbox Entry's hand commit.
2. API Contract: `RecordPullRequest` commits with that subject. Every other
   step, its proofs, its adoption paths and its results stay as they are.
   `pullRequestRecordCommitIsExact`, `pullRequestRecordTreeIsExact` and
   `pullRequestRecordChangeIsExact` keep their bytes; a record commit made by
   hand under any subject is adopted exactly as today.
3. API Contract: the queue commit helper. `Archive`, `RecordPullRequest` and
   `ResolveConflict` make their commit through one helper that runs `git
   commit -m <message>` through the workflow's Git runner. On success it
   behaves as today. On failure it classifies the commit as a hook refusal
   only when the error wraps an `*exec.ExitError` and
   `daemon.ClassifyCommitHookRefusal(ctx, workDir, output)` reports one, where
   `output` is the runner's error text with the `git <args joined by spaces>:
   ` prefix and the `: <exit error>` suffix removed when both are present, and
   empty otherwise. It never classifies on the arguments, the subject or the
   slug. Any other failure returns today's error unchanged.
4. API Contract: the refusal error. Its text is `the repository <hook> hook
   refused the commit "<subject>" (exit status <n>): <output>; <recovery>`,
   with `: <output>` left out when the output is empty. `<hook>` is the hook
   `ClassifyCommitHookRefusal` named; otherwise the hooks
   `daemon.InstalledCommitHooks` lists, joined with ` or `; otherwise
   `commit`. `<subject>` is the message's first line. `<output>` has every run
   of whitespace collapsed to one space and is cut at 512 bytes on a rune
   boundary, with ` ...` appended when cut. For the archive and record commits
   `<recovery>` is `the change stays staged in <workDir>; commit it there with
   a subject the hook accepts, then run roundfix deliver retry <slug>`; for
   the derived merge it is `the merge was aborted; merge the default branch
   into the item branch in <workDir> by hand, commit it with a message the
   hook accepts, then run roundfix deliver retry <slug>`. The error unwraps to
   the runner's error, so `errors.As` still finds the `*exec.ExitError`. The
   callers keep their wrapping (`commit archived Spec: %w`, `commit pull
   request record: %w`, `commit derived merge: %w`), so the park reason reads
   `delivery-error: record pull request: commit pull request record: the
   repository commit-msg hook refused the commit "…" (exit status 1): …`.
5. API Contract: `daemon.InstalledCommitHooks`. It reads the hooks directory
   with `git rev-parse --path-format=absolute --git-path hooks` as
   `commitHookInstalled` does, returns the executable regular files among
   `pre-commit` and `commit-msg` in that order, and returns none when Git
   answers nothing. `commitHookInstalled` returns whether that list is
   non-empty, so `ClassifyCommitHookRefusal` and the Daemon's Task commit
   refusal keep their behavior and every daemon test keeps passing unchanged.
6. API Contract: the Delivery Engine, `ClassifyPark`, the blocker constants,
   the stage an item parks at, `deliver status` and the Pull Request body do
   not change. A refused queue commit parks `delivery-error` with class
   `environment`, and its next action stays `resolve the blocker, then run
   roundfix deliver retry <slug>`, which the reason now spells out.
7. API Contract: documentation. The `deliver` user guide and the Roundfix
   `deliver` reference state the record subject and that a queue commit a
   hook refuses parks `delivery-error` with the hook, the subject and the
   recovery, and give the recovery: commit the staged change in the item
   worktree with a subject the hook accepts, then `roundfix deliver retry
   <slug>`; for a refused derived merge, the conflict park's hand merge. The
   Hook Refusal glossary entry covers a Delivery Queue commit.

### Surface Transcripts

None. No command's flags, help text, standard output or exit code changes.
The changed text reaches an operator only inside a Delivery Queue item's
blocker, which only a queue running against a forge produces; the tests of
API Contract 4 assert the step's whole error text in order instead.

### Invariants

```text
1. The Pull Request Record Commit subject is at most 70 characters for every slug and every accepted pull request number.
2. A record commit is adopted by its change alone; its message is never read.
3. A queue commit refused by a hook parks delivery-error, class environment, with a reason naming the hook, the subject and the recovery.
4. A failed queue commit in a repository with no executable commit hook, whose output names no hook, keeps today's error text.
5. The queue never passes --no-verify and never rewrites a refused commit's subject.
6. Every Archive Record already written, and the Daemon's Task commit refusal, stay as they are.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metrics 1-2, Invariant 1 → `pullRequestRecordSubject` (API Contract 1), `RecordPullRequest` (API Contract 2)
- Goal 2, Core Feature 2, Success Metric 3, Invariant 2 → the unchanged record proof (API Contract 2)
- Goal 3, Core Feature 3, Success Metrics 4-5, Invariants 3-5 → the queue commit helper and refusal error (API Contracts 3-4), `daemon.InstalledCommitHooks` (API Contract 5), the unchanged engine and Park Class (API Contract 6)
- Goal 4, Core Feature 4 → the `deliver` guide, the Roundfix `deliver` reference, the skill version and the Hook Refusal glossary entry (API Contract 7)
- Success Metric 6 → the PRD's Unreachable Acceptance
- Core Feature 3, Invariant 6 → API Contract 5, the unchanged daemon tests and the QA scope audit

## Absence Evidence

1. `pullRequestRecordSubject` — command: `grep -rn pullRequestRecordSubject internal`; result: no match at `70c78d15`; the subject is an inline string concatenation at `internal/cli/deliver_workflow.go:1352`.
2. The queue commit helper and refusal error — command: `grep -rn -e ClassifyCommitHookRefusal -e HookRefusalError internal/cli internal/delivery`; result: no match in production code at `70c78d15`; only `internal/cli/settle_test.go` names a Hook Refusal, for a Task commit.
3. `daemon.InstalledCommitHooks` — command: `grep -rn InstalledCommitHooks internal`; result: no match at `70c78d15`; `internal/daemon/commit_hook.go` holds only the unexported `commitHookInstalled`, which returns a boolean.

## ADR Alternatives

- ADR-0302: no rejected alternative is taken; the subject names the Spec number, the refusal stays a `delivery-error` with a named reason, and no `--no-verify`, cut slug or new blocker is introduced.

## Integration Points

Git, through the existing `git` executable, for `commit`, `rev-parse` and the
hooks Git runs in the item worktree; no other process, provider or network
endpoint. No test, Verification command or QA row reads the real
`~/.roundfix`, another repository or a remote other than a local bare
repository.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories, homes and Git configuration, and writes its hooks
as executable shell scripts under the temporary repository's `.git/hooks`
with the package's existing script fixture helper. No new test waits on real
time.

- `internal/cli/deliver_commit_hook_test.go` (new):
  - `TestPullRequestRecordSubjectNamesTheSpecNumber`: subtests for a
    four-digit prefix, a five-digit prefix, a ten-digit prefix with a
    ten-digit pull request number (70 characters), an eleven-digit prefix, a
    slug of digits alone, and a slug without digits (no identifier).
  - `TestDeliveryWorkflowRecordCommitPassesAHeaderLimitHook`: a fixture whose
    Spec slug has 54 characters and starts with `0088-`, archived and
    committed on an item branch, with an executable `commit-msg` hook that
    refuses a first line longer than 100 characters;
    `RecordPullRequest(…, "280")` returns an exact new head whose subject is
    the 57-character `docs: record pull request #280 in the 0088 Archive
    Record`, and the worktree is clean.
  - `TestDeliveryWorkflowAdoptsARecordCommitWhateverItsSubject`: subtests
    commit the record change by hand under the old long subject, under the
    short subject and under `wip`; `RecordPullRequest` from the candidate
    returns that head as exact and makes no second commit.
  - `TestDeliveryQueueCommitHookRefusalNamesTheHookAndTheRecovery`: subtests
    for a refused record commit whose hook prints `husky - commit-msg hook
    exited with code 1 (error)` (hook `commit-msg` from the output), a refused
    record commit whose silent `commit-msg` hook is the only one installed
    (hook `commit-msg` from the installed list), a refusing silent
    `pre-commit` beside a passing `commit-msg` (hook `pre-commit or
    commit-msg`), a refused archive commit whose slug contains `pre-commit`
    while only a silent `commit-msg` hook is installed (hook `commit-msg`,
    never the name in the subject), and a refused derived merge commit
    (merge aborted, worktree clean at the candidate). Each asserts the error
    text of API Contract 4 in order, that `errors.As` finds an
    `*exec.ExitError`, that the subject's slug never decides the hook, and
    that `delivery.ClassifyPark` of an item whose blocker is
    `delivery-error: record pull request: ` plus the error gives
    `environment`. The record subtests then commit the staged change by hand
    with a conforming subject and prove that `RecordPullRequest` adopts it.
  - `TestDeliveryQueueCommitWithoutAHookKeepsTheGitError`: with no hook and a
    stale `.git/index.lock`, the record commit fails with today's `commit pull
    request record: git commit -m …` error, which contains no `hook refused`.
- `internal/cli/deliver_record_pull_request_test.go`: declared break.
  `TestDeliveryWorkflowRecordsThePullRequestInTheArchiveRecord` expects the
  subject `docs: record pull request #7 in the 0001 Archive Record` instead of
  the slug form; every other assertion and test in the file keeps its bytes.
- `internal/daemon/commit_hook_test.go` and `internal/daemon/task_engine_test.go`
  are unchanged and pass, proving API Contract 5 keeps the Daemon's behavior.
- Documentation: task_01 re-records `docs/references/behavior-surfaces.json`
  with the repository's record command, because the `deliver` guide's
  fingerprint changes, raises and records the Roundfix skill version, and
  syncs the skill mirror.

## Build Order

1. The `deliver` guide, the Roundfix `deliver` reference, the skill version,
   the Behavior Surface Record and the Hook Refusal glossary entry (no
   dependency: the text states the contracts of `_techspec.md`, and its
   Verification reads only documents).
2. The subject, the queue commit helper, the refusal error,
   `daemon.InstalledCommitHooks` and their tests (depends on: 1, which writes
   the `deliver` guide and reference for the surface this step changes, so the
   code and the text it changes land in one dependency chain).
3. QA gate (depends on: 1, 2).

No other Spec is a prerequisite. The other cycle AM plans merge with this one
and deliver after it; a later one that also edits
`internal/cli/deliver_workflow.go` rebases on this change.

## Risks & Considerations

- A genuine Git failure in a repository that holds a commit hook, such as a
  stale index lock, is reported as a hook refusal, as for a Task commit; the
  hand commit the reason names fails the same way and shows the real cause.
- A hook that edits the message in place and still exits zero is not a
  refusal and is not touched.
- A hook whose own output echoes the refused subject, as commitlint does with
  its `input:` line, can carry a hook name that the subject holds; the
  Daemon's classifier then names the longest name it finds. The reason still
  quotes the hook's output, so the operator sees which hook ran.
- A slug that starts with digits but is not a Spec number (`2026-10-…`) names
  `2026` in the subject; the record path still identifies the Spec, and the
  subject stays within 70 characters.
- The archive subject and the pull request title still exceed 100 characters
  for a slug longer than 86 characters, now with a named park instead of a
  generic one; the open Backlog Entry tracks them.

## Glossary

- changes: **Hook Refusal**

## Decisions

- The subject names the Spec number, and a slug without one drops the
  identifier, so the 100-character limit holds for every slug.
- One helper for the three queue commits, reusing
  `daemon.ClassifyCommitHookRefusal`; no second classifier.
- The refusal keeps `delivery-error`, its class and its stage; only the
  reason text changes.
- The hook's output is read from the runner's error with its known prefix and
  suffix removed, never from the arguments, so a slug that contains a hook's
  name cannot decide the classification.
- ADR-0302 records the decision and its rejected alternatives.
