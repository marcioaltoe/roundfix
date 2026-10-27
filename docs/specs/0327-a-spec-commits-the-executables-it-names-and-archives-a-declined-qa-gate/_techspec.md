---
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
prd: _prd.md
created: 2026-10-26
---

# A Spec commits the executables it names, settle says what it removes, and archive honours a declined QA gate — Technical Spec

## Executive Summary

Three defects of the Run lifecycle, from Quaestor's report (`_prd.md` →
Sources). First, the Daemon's stage filter refuses every untracked executable,
so a Spec that asks for a launchd script loses its Task; the authorization
record gains `executables:`, read from the authorization the Run resolved at
start, and the Task commit keeps exactly those paths (new ADR-0340, amending
ADR-0157). Characterizing the filter found that it judges an untracked
directory as one entry, so an executable inside a new directory already rides
a Task or Batch commit; both commits now expand the directory before the
filter. Second, `roundfix settle` exits `0` after it integrates and removes the
kept Run Worktree, but says nothing about the removal; the 128 Quaestor saw is
Git's exit in a shell whose working directory was that worktree, so settle now
names each worktree it removes. Third, `roundfix archive` reads a QA Report for
every Spec with a Task Graph, so a `qa: declined` Spec can only archive with an
override; it now archives on Task completion alone, as ADR-0154 decided.

The trade-off accepted: the executables list is one more field an author must
fill before the Run, and an amendment must land on `main` before a retry
reads it. Reading it from the Task file instead would let the work widen its
own grant.

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local repository commands only.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0157, ADR-0154 and the new
  ADR-0340, as `_prd.md` → Project Constraints states; ADR-0251, ADR-0262,
  ADR-0269 and ADR-0306 cite ADR-0154 or one another for other archive and
  cancel paths (a legacy folder's failed QA, the cancelled record, ignored
  Spec files, a link leaving the repository); this Spec leaves each of them
  unchanged. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AQ batch 5 answers quoted in the authorization record grant the twelve
  Governed Paths below. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/settle.md`, `.agents/skills/roundfix/references/spec-delivery.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `internal/spec/archive.go`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Current behavior

Read and measured on 2026-10-26 at `6a2630e5`, in a prototype clone.

1. `internal/daemon/task_engine.go` `FilterStageablePaths` keeps a tracked
   path whatever its mode and drops an untracked regular file with an execute
   bit as `executable file`, `Lost: true`. `prepareTaskCommit` calls it, and
   `commitTask`'s caller turns a lost path into the Task failure `Task commit
   lost output: <path> (executable file)`. `publishDroppedStagePath` prints
   `roundfix: refused executable file <path> (mode <mode>); build artifacts and
   deliberately executable repository files are not valid Work Item output`
   for every commit kind. ADR-0157: "An executable file that is not tracked
   stays refused, which keeps the original protection against a build artifact
   riding into the tree."
2. The filter's input comes from `git status --porcelain` without
   `--untracked-files=all` (`GitWorktreeSnapshotter.Snapshot` in
   `internal/daemon/daemon.go`), so a new directory arrives as one entry, the
   filter sees a directory and keeps it, and `prepareTaskCommit` expands it
   into files only after the filter (`expandUntrackedCommitDirectories`). A
   Batch commit (`internal/daemon/engine.go`) stages the directory entry with
   `git add -f`. Measured: a Task whose Agent writes `newdir/helper.sh` with
   mode `0755` settles `completed` and commits it; a Batch whose snapshot is
   `newdir/` stages `newdir`. The QA Report commit reads a file-granular
   snapshot (`snapshotFiles`) and has no such hole.
3. `internal/authorization/authorization.go` `parseAuthorizationFrontmatter`
   reads `status`, `granted`, `action`, `consuming`, `paths`, `operations` and
   `precondition_repairs`; other keys are ignored. The record a Run uses is
   the one `roundfix implement` resolves at start (`ReadSpecAuthorization` in
   `internal/cli/implement.go`), carried as `TaskPlan.Authorization`. The
   Project Config precedent: `authorizationBoundsProjectConfig` lets a Task
   commit keep `.roundfixrc.yml` only when a granted record's `paths` lists it.
   Readers of the record: `internal/speccheck/constraints.go` (validity, which
   reports a `malformed_record` reason as `SC-TOOLING-UNTYPED`),
   `internal/speccheck/mechanical.go` and `undeclared.go` (`paths`,
   regenerations), `internal/suiteguardcontract/regeneration.go`
   (regenerations), `internal/cli/deliver_workflow.go` and `deliver_plan.go`
   (operations), `internal/cli/settle.go` (commit authority) and
   `internal/cli/implement.go` and the Daemon (`precondition_repairs`, Project
   Config, governed mutations). None reads an unknown key.
4. `internal/cli/settle.go` `settleTaskAndCommit` stages with `git add -A`
   and filters nothing, so it commits an executable a Task commit refuses; a
   tracked file whose mode a Task changes passes the filter as tracked. Both
   stay as they are (`_prd.md` → Non-Goals; open Backlog Entry
   `docs/backlog/2026-10-26-settle-and-a-tracked-mode-change-commit-an-executable-the-task-commit-refuses.md`).
5. `runSettleCommand` integrates a kept Task Worktree
   (`integrateTaskWorktree`, which removes it), then the Run
   (`integrateSettledRun`, which fast-forwards the user branch and removes the
   Run Worktree through `cleanupCleanRunWorktree`), then prints the commit
   lines and `settled <task> completed — <sha>` and returns `0`. Nothing names
   a removal. Measured with the built binary on a kept Run: run from the user
   checkout, stdout ends `settled task_01 completed — <sha>`, the exit is `0`
   and the Run Worktree is gone. Roundfix's exit codes are `0`, `1`, `2` and
   `130`; none is `128`.
6. In a shell whose working directory was removed, `git rev-parse HEAD`
   prints `fatal: Unable to read current working directory: No such file or
   directory` and exits `128` (measured with `/bin/sh`). Quaestor's 128 came
   from such a command run after settle, in the Run Worktree it had been
   editing.
7. `internal/spec/archive.go` `Archive`, without an override, requires every
   Task completed, then `ReadQAReport`; with no report it fails `no passing QA
   verdict: no QA Report found for Spec directory …`. It never reads
   `Graph.QADeclined`, which `internal/spec/spec.go` `parseQADeclaration` sets
   only for `qa: declined` with a non-empty `qa_reason` (and no QA Task then
   exists). With `--qa-override` and a declined graph whose Tasks are
   completed, the override succeeds and records `qa_override_qa_outcome:
   missing`. `BuildArchiveRecord` gives a record without a QA Report the
   disposition `no-qa`. ADR-0154: "A normal declined-gate or qualifying
   declared-partial case retains its own eligibility rules and does not
   acquire an invented override." `docs/agents/docs-layout.md`: "Declared QA
   decline and qualifying declared-only partial evidence are evaluated under
   their own contracts, not automatically called overrides."
8. Settle run from inside the Run Worktree resolves that worktree as the user
   checkout, commits there and prints `integration pending — git merge
   --ff-only roundfix/run-<id>`, exit `1` (measured). Out of scope; open
   Backlog Entry
   `docs/backlog/2026-10-26-settle-run-inside-the-run-worktree-takes-it-for-the-user-checkout.md`.
9. The Delivery Queue archives through `roundfix archive <slug>`
   (`internal/cli/deliver_workflow.go` `Archive`), so the archive fix reaches
   the queue unchanged.

## System Architecture

### Authorized Executables

The authorization record gains `Executables []string`. The Task commit builds
a set from it only when the Run's resolution is granted, and the stage filter
keeps an untracked executable whose exact repository-relative path is in the
set. QA Report and Batch commits pass no set. Only a Task commit's refusal
line gains the hint, so QA and Batch refusals keep their bytes (API Contracts
1 to 3).

### One file at a time

`prepareTaskCommit` and the Batch commit call
`expandUntrackedCommitDirectories` before the filter instead of after it, so
the filter judges each file of a new directory (API Contract 4). An
executable that a Task writes into a directory Git has never seen is then
refused like one written beside tracked files, and the Task commit's refusal
names the file and how the Spec's authorization can list it (API Contract 3),
while the directory's regular files are still staged by name.

### Archive of a Declared QA Decline

`Archive` skips the QA Report read when the graph is declined and every Task
is completed; the override path refuses a declined graph whose Tasks are all
completed, as it refuses any Spec that qualifies (API Contract 5). Its usage,
guide, skills and glossary entries follow in the same Task (API Contract 8).

### Settle names what it removes

`runSettleCommand` prints one stderr line after each successful removal
(API Contract 6): `Task Worktree removed: <path>` after the Task Worktree
integration and `Run Worktree removed: <path>` after the Run-level
integration. Settle keeps removing what it integrated, as the Implement
Command does after a Clean outcome, and its exit code stays the settled
outcome; a caller whose shell stood in a removed worktree now reads why its
next Git command exits 128 there. Its usage, guide, skills and glossary
entries follow in the same Task (API Contract 9).

### The rules in writing

Each change carries its own written rule: the executables rule in the
`write-tasks` skill, the Roundfix `implement` reference, the `implement` guide
and `CONTEXT.md` (API Contract 7); the decline rule in the archive usage, the
Roundfix `archive` and `spec-delivery` references, the `archive-spec` skill,
the `archive` guide and `CONTEXT.md` (API Contract 8); the settle rule in the
settle usage, the Roundfix `settle` and `spec-delivery` references, the
`settle` guide and `CONTEXT.md` (API Contract 9).

task_01 builds Authorized Executables and One file at a time, task_02 the
Archive of a Declared QA Decline, and task_03 Settle names what it removes,
each with its rule in writing.
Each Task is a vertical slice through one command's code and its documents.
The three share `CONTEXT.md`, the Roundfix `SKILL.md` version pair, the
owned-skill record and the Behavior Surface record, so they run in series.

## Implementation Design

### Interfaces

```go
// internal/authorization: AuthorizationRecord gains
//	Executables []string // exact repository-relative file paths
// internal/daemon
func FilterStageablePaths(ctx context.Context, workDir string, paths []string) ([]string, []DroppedStagePath) // unchanged: no executable allowed
func filterStageablePaths(ctx context.Context, workDir string, paths []string, executables map[string]bool) ([]string, []DroppedStagePath)
func authorizedExecutables(resolution spec.AuthorizationResolution) map[string]bool // nil unless granted
```

### Data Models

The authorization record's frontmatter gains an optional key; nothing else
persists differently.

```yaml
executables:
  - scripts/auditor-continuo.sh
```

### API Contracts

1. API Contract: the `executables` field. An absent or null key is an empty
   list. A present key must be a YAML list of strings; each is an exact
   repository-relative file path: non-empty, no surrounding whitespace, no
   leading or trailing `/`, no `*`, `?`, `[` or `\`, equal to its
   `path.Clean` form, not `.` and not starting `../`, and not repeated. Any
   other shape returns the reason code `malformed_record` with `Field:
   "executables"`, the offending value, and a detail that starts `executables`
   (`executables entry "<value>" <why>`, `executables repeats path "<value>"`,
   or the shared `executables must be a list`), so the record grants nothing
   and `roundfix spec check` reports `SC-TOOLING-UNTYPED`. The list bounds no
   Governed Path and changes no `paths` rule.
2. API Contract: the Task commit. `prepareTaskCommit` passes
   `authorizedExecutables(plan.Authorization)`; the filter keeps an untracked
   regular file with an execute bit when its stage path, in slash form, is in
   the set, and the commit stages it with the paths the Committer already
   receives, so Git records the mode it reads (`_prd.md` → Acceptance
   evidence). An unlisted one is dropped exactly as today (`executable file`,
   its mode, `Lost: true`, the Task fails `Task commit lost output: <path>
   (executable file)`). A Run whose authorization is absent, refused or
   unresolved keeps no executable.
3. API Contract: the refusal line. For a Task commit (a non-empty Task id),
   the progress line and the Run Event summary append
   ` unless the Spec's _authorization.md lists the path under executables`
   after `not valid Work Item output`; the event payload keeps `decision`,
   `path`, `reason`, `mode` and `task`. A QA Report commit and a Batch commit
   print today's bytes and never keep an executable, listed or not.
4. API Contract: one file at a time. `prepareTaskCommit` expands untracked
   directories before the Project Config drop is checked against files and
   before the filter; the Batch commit in `internal/daemon/engine.go` expands
   them after its Project Config drop and before `FilterStageablePaths`, and
   returns the expansion error as it returns its other Git errors. A new
   directory's executable is then refused like any other; its regular files
   are staged by name.
5. API Contract: archive of a Declared QA Decline. Without an override, a
   graph with `QADeclined` archives when every Task is `completed`, reads no
   QA Report, records `unproven` as empty, and writes the record
   `BuildArchiveRecord` already writes (disposition `no-qa`, no
   `qa_override`); a pending, failed or skipped Task still refuses with the
   existing message. With `--qa-override`, a declined graph whose Tasks are
   all completed is refused with the existing `QA archive override is not
   allowed because the Spec already qualifies for normal archive`; with a
   non-completed Task the existing override rules apply. Every other graph
   keeps today's path byte for byte.
6. API Contract: settle's removal lines. After `integrateTaskWorktree`
   succeeds, settle prints `Task Worktree removed: <task worktree path>` on
   stderr; after `integrateSettledRun` returns no pending command, it prints
   `Run Worktree removed: <run worktree path>`. Paths are the recorded refs'
   paths. An integration conflict, a pending integration or a failure prints
   neither line for a worktree that was not removed. Stdout and the exit codes
   do not change.
7. API Contract: the executables rule in writing (task_01). This sentence,
   verbatim and without inline code, is added to the `write-tasks` skill under
   its tooling steps, to its task template as a Verification comment, to the
   Roundfix `implement` reference beside the refusal text, and to the
   `implement` guide: "A Task commit keeps an untracked executable file only
   when its exact path is listed under executables in the Spec's authorization
   record; every other untracked executable is refused." Each place also says
   that the list must be on the delivery target before the Run reads it, that
   QA Report and Batch commits never keep one, and that a repository with
   core.fileMode false records the file as 100644. The `implement` reference's
   refusal text gains the Task-commit suffix of API Contract 3. `CONTEXT.md`
   adds **Authorized Executable**.
8. API Contract: the decline rule in writing (task_02). This sentence,
   verbatim and without inline code, is added to the Roundfix `archive` and
   `spec-delivery` references, the `archive` guide and a new `### Declared QA
   decline` section of the `archive-spec` skill: "A Task Graph that declares
   qa: declined with a non-empty qa_reason archives once every Task is
   completed, without a QA Report and without an override, and its Archive
   Record has the disposition no-qa." Each place says that `--qa-override` is
   refused for such a Spec. `archiveUsage`'s first paragraph names the case
   ("every Task is completed and the Task Graph declares qa: declined with a
   qa_reason") and its override paragraph adds "It is also refused for a Task
   Graph that declares qa: declined once every Task is completed."
   `CONTEXT.md` adds **Declared QA Decline** and revises **Archive Command**
   with the sentence "A Declared QA Decline archives on Task completion
   alone, without an override."
9. API Contract: the settle rule in writing (task_03). This sentence,
   verbatim and without inline code, is added to the Roundfix `settle` and
   `spec-delivery` references and the `settle` guide: "After a successful
   integration, settle removes the kept Task Worktree and Run Worktree it
   integrated and names each on stderr; leave any shell that stands in one of
   them before you run settle, because Git exits 128 in a removed directory."
   Each place shows the two stderr lines and says to run settle from the user
   checkout. `settleUsage` adds "After a successful integration, settle
   removes the kept Task Worktree and Run Worktree it integrated and names
   each on stderr." `CONTEXT.md` revises **Run Worktree** ("removed after a
   Clean integrated spec outcome or after settle integrates it") and
   **Settle Command** ("names each worktree it removes on stderr").

### Surface Transcripts

None. Settle's stderr lines need a kept Run, which the QA gate cannot create
without an Agent; API Contract 6 is asserted by named tests instead.

### Declared breaks

1. An untracked executable inside a new directory is now refused in a Task
   commit (the Task fails) and dropped from a Batch commit; before, it was
   committed.
2. `--qa-override` on a declined Spec whose Tasks are completed is now
   refused; Quaestor's override of 0020 and 0021 would be.
3. A Task commit's executable refusal line gains its suffix:
   `internal/daemon/task_engine_test.go`
   `TestTaskCommitDropsExecutableFileAndCommitsRemainingPaths` changes its
   expected line.
4. Settle's stderr gains the removal lines: in
   `internal/cli/settle_test.go`,
   `TestSettleTaskStatusRetargetsKeptRunWorktreeAndCleansUpAfterIntegration`,
   `TestRunSettleRetargetsKeptTaskWorktreeAndCleansUpAfterIntegration`,
   `TestSettleAcceptsCompletedTaskFromKeptTaskWorktreeAfterHookRefusal` and
   `TestSettleCommitsDeletedAndRenamedWorkFromTaskWorktree` expect them.
5. `archive --help` and `settle --help` change: the usage pin in
   `internal/cli/archive_pull_request_test.go` (task_02) and the Behavior
   Surface record (each Task) move.

### Invariants

1. No commit keeps an untracked executable whose path a granted
   authorization of the Run does not list under `executables`.
2. QA Report and Batch commits never keep an untracked executable.
3. A declined graph never archives with a non-completed Task.
4. Settle never prints a removal line for a worktree it did not remove, and
   its exit code is the settled outcome.

## Coverage Map

- Goal 1 → Authorized Executables (API Contracts 1-2)
- Goal 2 → One file at a time (API Contracts 3-4)
- Goal 3 → Settle names what it removes (API Contract 6)
- Goal 4 → Archive of a Declared QA Decline (API Contract 5)
- Goal 5 → The rules in writing (API Contracts 7-9)
- Core Features 1-2 → API Contracts 1-4
- Core Feature 3 → API Contract 6
- Core Feature 4 → API Contract 5
- Core Feature 5 → API Contracts 7-9
- Success Metrics 1-2 → API Contracts 2-3, Invariants 1-2
- Success Metric 3 → API Contract 4
- Success Metric 4 → API Contract 1
- Success Metric 5 → API Contract 5, Invariant 3
- Success Metric 6 → API Contract 6, Invariant 4
- Success Metric 7 → API Contracts 7-9

Deviations from the PRD: None.

## Absence Evidence

1. The `executables` field — command: `grep -rn 'executables' internal/authorization internal/spec/authorization.go internal/daemon/task_engine.go`;
   result: no match at `6a2630e5`.
2. A decline branch in archive — command: `grep -n 'QADeclined' internal/spec/archive.go internal/cli/archive.go`;
   result: no match at `6a2630e5`; `grep -rn 'QADeclined' internal` finds only
   `internal/spec/spec.go`, `internal/speccheck/consumers.go` and
   `internal/speccheck/transcripts.go`.
3. Settle's removal lines — command: `grep -rn 'Worktree removed' internal`;
   result: no match at `6a2630e5`.
4. **Authorized Executable** and **Declared QA Decline** — command:
   `grep -n 'Authorized Executable\|Declared QA Decline' CONTEXT.md`; result:
   no match at `6a2630e5`.

## ADR Alternatives

- ADR-0157: its decision stands for every unlisted executable; ADR-0340
  amends only its accepted cost, and its rejected alternatives (reading the
  file's content, staging every executable) stay rejected.
- ADR-0154: implemented as decided; the rejected alternative of treating a
  declined gate as an override is the behavior this Spec removes.
- ADR-0251: unchanged; a legacy folder's failed QA keeps `failed-qa`, and a
  declined graph never has a QA Report to fail.
- ADR-0262: unchanged; a cancelled Spec keeps its own record and never
  reaches the decline branch.
- ADR-0269: unchanged; the ignored-file refusal runs before eligibility for a
  declined Spec too.
- ADR-0306: unchanged; the link refusal runs before eligibility for a
  declined Spec too.
- ADR-0340: its rejected alternatives (staging every executable, naming
  executables in the Task file, globs or directories) are not taken.

## Integration Points

The Delivery Queue reaches the archive fix through the `roundfix archive`
subprocess (Current behavior 9). The Linux Lane runs the changed
`internal/authorization`, `internal/daemon`, `internal/spec` and
`internal/cli` tests in `golang:1.26`; the mode tests read `100755` there too
(the container's Git keeps `core.fileMode: true`).

## Testing Approach

- task_01: `internal/authorization/executables_test.go` (the field and every
  malformed shape) and `internal/daemon/task_engine_executables_test.go`:
  a listed executable kept and an unlisted one refused with the hint
  (`prepareTaskCommit`), a real-Git Task Cycle committing `100755`, a real-Git
  Task Cycle refusing `newdir/helper.sh`, a Batch staging only the regular
  file of a new directory, and a QA Report commit refusing a listed path with
  today's line.
- task_02: `internal/cli/archive_declined_qa_test.go` through CLI dispatch
  (accepted decline, refused override, pending Task) and the archive usage
  pin.
- task_03: settle tests in `internal/cli/settle_test.go` asserting the full
  stderr of both removal cases and of a pending integration.
- Each Task: phrase checks on every file it writes, the owned-skill record,
  and the documentation contracts with the Behavior Surface record
  re-recorded as its last step.
- QA, each sabotage on its own `git archive` copy of the audited head (lesson
  39, the case named in parentheses distinguishes the mutation):
  - drop the set check in the filter → `TestTaskCommitKeepsAnAuthorizedExecutable`
    and `TestTaskCycleCommitsAnAuthorizedExecutableWithItsMode` fail (input:
    a listed `scripts/launch.sh` with mode `0755`);
  - read the set in QA commits too → `TestQAReportCommitRefusesAListedExecutable`
    fails (input: a listed path under `qa/evidence/`);
  - expand directories after the filter again →
    `TestTaskCycleRefusesAnUnlistedExecutableInANewDirectory` fails (input:
    `newdir/helper.sh` in a directory Git has never seen);
  - drop the Batch expansion → `TestResolveCycleRefusesAnExecutableInsideANewBatchDirectory`
    fails (input: a snapshot entry `newdir/`);
  - accept a glob → `TestAuthorizationRefusesAMalformedExecutablesEntry/glob`
    fails;
  - read the QA Report for a declined graph again →
    `TestArchiveAcceptsADeclaredQADecline` fails;
  - skip the override refusal for a declined graph →
    `TestArchiveRefusesAnOverrideForAQualifyingDeclaredQADecline` fails;
  - print `Run Worktree removed:` also when integration is pending →
    `TestSettleNamesNoRemovalWhenIntegrationIsPending` fails.

## Build Order

1. task_01: Authorized Executables and one file at a time, with their written
   rule (API Contracts 1-4 and 7).
2. task_02: archive of a Declared QA Decline with its written rule (API
   Contracts 5 and 8) (depends on: 1, for the shared `CONTEXT.md`, Roundfix
   `SKILL.md`, owned-skill record and Behavior Surface record).
3. task_03: settle's removal lines with their written rule (API Contracts 6
   and 9) (depends on: 2, for the same shared files and
   `spec-delivery.md`).
4. task_04: the QA gate (depends on: 3).

## Risks & Considerations

- A Spec already delivering an executable through a new directory (Current
  behavior 2) now fails its Task; the refusal names the file and the fix, and
  ADR-0157 always meant this.
- `executables:` must be on the Delivery Base before a Run or retry reads it,
  like every authorization amendment.
- A repository with `core.fileMode: false` commits a listed file as `100644`;
  the docs say so.
- Spec 0322 rewrites sentences of `settle.md`, `archive.md` and the
  `archive-spec` skill before this Spec; task_02 and task_03 add paragraphs
  and pin no bytes of theirs.

## Glossary

- adds: **Authorized Executable**
- adds: **Declared QA Decline**
- changes: **Run Worktree**
- changes: **Settle Command**
- changes: **Archive Command**

## Decisions

- `executables:` in the authorization frontmatter, exact paths, read from the
  Run's resolved authorization (ADR-0340).
- Reuse the `malformed_record` reason code with `Field: "executables"`, so
  `roundfix spec check` needs no change.
- Expand directories before the filter in Task and Batch commits, the class
  fix of the hole the characterization found.
- Keep settle's removal; name it. No new exit code.
- Reuse the `no-qa` disposition for a Declared QA Decline; no record schema
  change.
