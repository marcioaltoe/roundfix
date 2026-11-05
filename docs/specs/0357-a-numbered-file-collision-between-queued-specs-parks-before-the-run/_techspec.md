---
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
prd: _prd.md
created: 2026-11-05
---

# A numbered file collision between queued Specs parks before the Run — Technical Spec

## Executive Summary

Four slices. First, Project Config gains `delivery.numbered_files`, a list of
path patterns with the Derived Path Declaration grammar, empty by default, and
a predicate that tells a Numbered File from any other matching file. Second,
`internal/delivery` gains a pure comparison that turns the Numbered Files an
item adds, those the default branch holds and those every earlier unmerged
item adds into a report of collisions and next free numbers, and the engine
parks a colliding item as `numbered-file-collision` before its Run, refuses a
retry of an item that has not run while it collides, and classifies the park
in the conflict Park Class. Third, the command workflow's Delivery Revalidation
reads the declaration at the default-branch commit and, only when it declares
patterns, reads the files from Git, the plans from the Spec folders and the
earlier items from the Run Database. Fourth, the documents, `CONTEXT.md` and
ADR-0367. The trade-off this design accepts is that only declared or
committed numbers are seen before a Run: a number a Run would generate without
its plan naming it is found later, for every later item, once a branch holds
it.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or stored;
  the comparison reads existing file names and Spec slugs. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the check reads Git, the Spec folders and the Run
  Database. Source: `docs/agents/go.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0312 skips Delivery Revalidation
  for an Archived Item Resume, ADR-0312: "runs neither Delivery Revalidation
  nor the Implement Command"; the check lives inside Delivery Revalidation,
  so it does not run there. The declaration is read at the default-branch
  commit as `delivery.derived_paths` is, but no merge step changes, and
  prerequisite handling is untouched. ADR-0367, written by task_03, records
  the decision below. Source: `docs/agents/domain.md`,
  `docs/agents/autonomous-work.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ batch 2 answers of 2026-11-05 in the authorization record; bounded files
  are the Roundfix skill files task_03 changes. No Makefile, workflow,
  `go.mod`, `go.sum`, lint, formatter, Baseline module or `.roundfixrc.yml`
  change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read at `ad77a797`:

1. `(*Engine).advanceItem` (`internal/delivery/engine.go`) creates or continues
   the item worktree, proves an archived item, then calls `Revalidate` with the
   merge commits of earlier merged items. It records `premise-changed` and
   `owner-older-than-main` warnings, parks `revalidation-failed: <codes>` on
   strict findings, and otherwise moves the item to `running`.
2. `(*Engine).Retry` calls `Revalidate(..., nil)` only for an item with no
   recorded Run, refuses with `retry Delivery Queue item %q: revalidation
   failed: <codes>` and leaves the item unchanged; an item with a Run is not
   revalidated.
3. `(*commandDeliveryWorkflow).Revalidate` (`internal/cli/deliver_revalidate.go`)
   runs the strict check in the item worktree and compares declared production
   Go `interface:` paths with earlier merge commits. It reads no Run Database
   row and no file name number.
4. `ClassifyPark` (`internal/delivery/park_class.go`) maps
   `pull-request-conflict` alone to the conflict Park Class;
   `TestEveryBlockerHasAParkClass` fails on an exported blocker without a case.
5. `roundconfig.DeliveryConfigAtCommit` reads Project Config from a default
   commit with User Config beneath it; `delivery.derived_paths` entries are
   validated by `validateDerivedPathPatterns` and matched by
   `matchesDerivedPath`. A Project Config with an unknown `delivery` key is
   refused.
6. `(*commandDeliveryWorkflow).CreateItemBranch` fetches the default branch
   into `refs/remotes/<remote>/<default>` and creates a new item worktree at
   it; a recorded item branch is reused as it is, so a retry runs on the base
   the branch was created from.
7. A conflicting Pull Request is found only in the check wait; a conflict
   outside `delivery.derived_paths` parks `pull-request-conflict: <paths>`.

## System Architecture

- `internal/config`: the `delivery.numbered_files` key, its validation and
  `(Delivery).IsNumberedFile`.
- `internal/delivery`: `FindNumberedFileCollisions`, `NumberedFileReport` and
  `BlockerNumberedFileCollision`; `Revalidation.NumberedFiles`; the park in
  `advanceItem`, the refusal in `Retry` and the case in `ClassifyPark`.
- `internal/cli`: `numberedFileReport`, called by `Revalidate` after the
  premise comparison, reading the default commit, the item and earlier items
  through Git, their Spec folders through `spec.Load`, and the queue through
  the Run Database.
- Documents: `CONTEXT.md`, ADR-0367, the configuration and `deliver` guides,
  the Roundfix `deliver` reference and its mirror.

## Implementation Design

### Interfaces

```go
// internal/config
type Delivery struct {
	// existing fields …
	NumberedFiles []string
}

// IsNumberedFile reports whether a numbered_files pattern matches name and
// the base name of name starts with an ASCII digit.
func (delivery Delivery) IsNumberedFile(name string) bool

// internal/delivery
const BlockerNumberedFileCollision = "numbered-file-collision"

type NumberedFileClaims struct {
	Owner string   // the earlier item's Spec slug
	Paths []string // the Numbered Files that item adds
}

type NumberedFileCollision struct {
	Path  string // a Numbered File the item adds
	Other string // the colliding file
	Owner string // the earlier item's Spec slug, or "" for the default branch
}

type NumberedFileReport struct {
	Collisions []NumberedFileCollision
	NextFree   []string // "<directory>/<digits>", one per colliding directory
}

func FindNumberedFileCollisions(defaultFiles, subject []string, ahead []NumberedFileClaims) NumberedFileReport
func (report NumberedFileReport) Detail() string
func (report NumberedFileReport) Blocker() string

type Revalidation struct {
	// existing fields …
	NumberedFiles NumberedFileReport
}
```

### Data Models

No Run Database change. A park stores its text in the existing blocker column.

### API Contracts

1. API Contract: `delivery.numbered_files` is a YAML list of strings in User
   or Project Config, default empty. A project value replaces the user value;
   an absent project key inherits it; an explicit `[]` clears it. Each entry
   passes `validateDerivedPathPatterns("delivery.numbered_files", entries)`, so
   a refusal reads `delivery.numbered_files has unsafe path "<entry>"` or
   `delivery.numbered_files has invalid pattern "<entry>": <cause>`. The merge
   of the overlay sits beside the `derived_paths` merge in `applyConfigContent`
   and the validation beside `validateDerivedPaths` in `Validate`.
2. API Contract: `(Delivery).IsNumberedFile(name)` is true when
   `matchesDerivedPath(delivery.NumberedFiles, name)` holds and `path.Base(name)`
   starts with a byte `0` through `9`. With `[drizzle/]`,
   `drizzle/0052_a.sql` and `drizzle/meta/0052_snapshot.json` are Numbered
   Files and `drizzle/meta/_journal.json` is not; with `[drizzle/*.sql]`,
   `drizzle/x0052.sql` and `db/0052_a.sql` are not.
3. API Contract: `FindNumberedFileCollisions` reads each path's number as the
   run of ASCII digits at the start of its base name, parsed with
   `strconv.ParseUint` base 10. A path whose base name does not start with a
   digit, or whose run is longer than 18 digits, is ignored. Two paths share a
   number when `path.Dir` is equal and the values are equal, so `052_b.sql`
   and `0052_a.sql` share one. Each distinct subject path, in sorted order,
   collides with every default file and every earlier item's path of the same
   directory and number, whatever the names, the same path included.
   `Collisions` is sorted by `Path`, then `Other`, then the rendered owner. For
   each directory that holds a collision, in sorted order, `NextFree` holds the
   directory, `/`, and one plus the highest number among the default files and
   every earlier item's paths in that directory, zero-padded to the largest
   digit count among the subject's colliding paths there; the root directory
   `.` renders the digits alone. Neither input slice is modified.
4. API Contract: `Detail` joins `<Path> collides with <Other> (<owner>)`, with
   `<owner>` the Spec slug or `default branch`, by `, `, then appends
   `; next free number: ` and `NextFree` joined by `, `. `Blocker` is
   `numbered-file-collision: ` followed by `Detail`. An empty report returns
   `""` from both. Example:
   `numbered-file-collision: drizzle/0052_b.sql collides with drizzle/0052_a.sql (alpha), drizzle/meta/0052_snapshot.json collides with drizzle/meta/0052_snapshot.json (default branch); next free number: drizzle/0053, drizzle/meta/0053`.
5. API Contract: in `advanceItem`, after the `revalidation-failed` park and
   before `setStage(running)`, a report with collisions parks the item with
   `revalidation.NumberedFiles.Blocker()`; the worktree stays recorded and no
   Run starts. Strict findings take precedence, and warnings are recorded as
   today.
6. API Contract: in `Retry`, for an item with no Run, after the findings
   refusal, a report with collisions returns
   `retry Delivery Queue item "<slug>": numbered file collision: <Detail>`,
   leaving the item unchanged and carrying nothing forward.
7. API Contract: `ClassifyPark` maps `numbered-file-collision` to the conflict
   Park Class with the next action
   `in <worktree>, renumber the item's files the blocker names from its next free number, in the Task declarations on the item branch and in any file the branch holds, or deliver the earlier item it names first and merge the default branch into the item branch; then run roundfix deliver retry <slug>`.
   `roundfix deliver status` prints it as `Park: <slug> conflict: <next>`, and
   the Pending Question answers with it.
8. API Contract: `(*commandDeliveryWorkflow).Revalidate` sets
   `NumberedFiles` from `numberedFileReport(ctx, workDir, specsRoot.Path,
   specSlug, graph)` after the premise comparison, and wraps its error as
   `check numbered files: %w`. `numberedFileReport`:
   1. resolves the default branch with `preflight.DetectDefaultBranch(ctx,
      workDir, "", workflow.git)` and the remote from `watch.push_remote`
      (default `origin`); an undetermined branch or an unresolvable
      `refs/remotes/<remote>/<default>^{commit}` returns an empty report;
   2. reads the configuration with `roundconfig.DeliveryConfigAtCommit` at
      that commit; an empty `NumberedFiles` returns an empty report with no
      further Git command and no Run Database read;
   3. otherwise fetches `+refs/heads/<default>:refs/remotes/<remote>/<default>`
      from the remote, resolves the commit again and re-reads the
      configuration there (an empty list then returns an empty report);
   4. lists the Numbered Files of that commit with `git ls-tree -r --name-only -z`;
   5. takes the subject's claims: the Numbered Files that
      `git diff --name-only --no-renames --diff-filter=A -z <merge-base> HEAD`
      lists in `workDir`, where the merge base is that of the default commit
      and `HEAD`, plus the `creates:` paths of each Task of `graph` that is
      neither completed nor of type `qa`, kept when they are Numbered Files
      absent at the merge base;
   6. reads the queue with `workflow.store.DeliveryQueue(ctx,
      workflow.loaded.GitRoot)` and, for each item before the subject that is
      not `merged`, takes its claims the same way: its head is
      `refs/heads/<item branch>` when that ref resolves and the default commit
      otherwise; its plan is loaded with `spec.Load` from its recorded
      worktree when that worktree holds the Spec's `_tasks.md` at the same
      relative Specs Root, and from the subject's Specs Root otherwise; a load
      error leaves only its branch claims. A subject absent from the queue
      has no earlier items;
   7. returns `delivery.FindNumberedFileCollisions`.

### Surface Transcripts

None. No command, flag or output format is added: the new park is a blocker
value printed through the existing `deliver status` row, `Park:` line and
Pending Question, and `TestDeliverStatusNamesTheNumberedFileCollisionPark`
asserts those lines for a seeded queue.

## Invariants

1. Without `delivery.numbered_files`, Delivery Revalidation runs no new Git
   command, reads no Run Database row and returns the result it returns at
   `ad77a797`.
2. Only an earlier item that is not `merged` and the default branch are
   compared with the subject; a later item never parks an earlier one.
3. A collision parks before any Run, Task carry-forward or review.
4. The check writes nothing: no file, ref other than the default branch's
   remote-tracking ref, or Run Database row.
5. Strict findings keep parking as `revalidation-failed`, ahead of a collision.

## Declared breaks

No existing test is expected to change. The existing Delivery Revalidation
fixtures declare no `delivery.numbered_files`, so Invariant 1 keeps them
passing; their workflows carry no Run Database, which the guard in step 2 of
API Contract 8 never reaches. `TestParkClassesAndNextActions` gains one row,
and `checkParkActionDocumentation`
(`internal/docscontract/command_documentation_test.go`) also requires the new
action sentence and the conflict row entry, so
`TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem` fails until task_03
writes them, in the same Task.

## ADR-0367

task_03 creates
`docs/adr/0367-a-numbered-file-collision-with-an-earlier-queue-item-parks-before-the-run.md`
with exactly this content, apart from line wrapping:

```markdown
---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# A numbered file collision with an earlier queue item parks before the Run

On 2026-11-05 fiscus reported five Specs of one Delivery Queue that all
started their database migrations at number 0052. Each Spec took the next
number of the main it knew, and the queue found the clash only at merge,
after the Run, the review, the gate and the checks, so each merge made the
others renumber by hand.

A repository now declares its numbered files in Project Config under
`delivery.numbered_files`, with the path grammar of a Derived Path
Declaration. A matching file whose name starts with digits is a Numbered File,
and its directory and the value of those digits are its place in a sequence.
Delivery Revalidation, before an item's first Run and on a retry of an item
that has not run, collects the Numbered Files the item adds: those its branch
adds since its merge base with the default branch, and those its not yet
completed Tasks declare under `creates:`. One that shares a directory and a
number with a file on the default branch, or with a file an earlier item of
the same queue adds while it is not merged, is a Numbered File Collision. The
item parks as `numbered-file-collision` in the conflict Park Class, naming
each pair, the earlier item or the default branch, and the next free number,
and the operator renumbers the plan, or delivers the earlier item first and
merges the default branch, before the retry.

## Considered Options

- A default pattern list. Rejected: names that start with digits are common
  outside sequences, such as dated Backlog Entries and Findings, so any default
  would park unrelated items; the declaration is off until a repository makes
  it.
- Waiting for the earlier item, as a Prerequisite Spec does. Rejected: a plan
  that names a taken number must change even after the earlier item merges,
  and a branch created before that merge keeps its old base.
- Renumbering by the queue. Rejected for now: renaming a migration also
  rewrites tool-owned files, such as a journal and its snapshots, that only
  the repository's generator owns; the queue reports the next free number and
  leaves the change to the operator or the Agent.
- Checking at merge only, as today. Rejected: it spends the Run, the review,
  the gate and the checks before the conflict appears.

## Consequences

A number a Run generates without its plan declaring it is seen only once a
branch holds it; from then on every later item is held to it. A retry of an
item whose Run already started is not checked and still meets the conflict at
merge, where `delivery.derived_paths` and `pull-request-conflict` keep their
behavior. Queue order decides precedence: the later item renumbers.
Timestamped folders, as drizzle-kit v1 writes them, are outside the
declaration.
```

## Coverage Map

- Goal 1, Story 1, Story 5, Core Features 1–2, Success Metric 1 →
  `Delivery.NumberedFiles`, its overlay and validation,
  `(Delivery).IsNumberedFile` (task_01).
- Goal 2, Story 2, Story 4, Core Features 3–4, Success Metric 2 →
  `FindNumberedFileCollisions` (task_02) and `numberedFileReport` (task_04).
- Goal 3, Story 3, Core Features 5–6, Success Metrics 3–4 →
  `NumberedFileReport.Blocker`, the park in `advanceItem`, the refusal in
  `Retry`, the `ClassifyPark` case (task_02) and the `deliver status` lines
  (task_04).
- Goal 4, Core Feature 7, Success Metric 5 → `CONTEXT.md`, ADR-0367, the
  configuration and `deliver` guides and both copies of the `deliver`
  reference (task_03).

Deviations from the PRD: None.

## Absence Evidence

1. `delivery.numbered_files` and `IsNumberedFile` — command: `grep -rn 'numbered_files\|NumberedFiles\|IsNumberedFile' internal docs/user-guide`; result: no match.
2. `FindNumberedFileCollisions`, `NumberedFileReport` and `BlockerNumberedFileCollision` — command: `grep -rn 'NumberedFile\|numbered-file-collision' internal`; result: no match.
3. `numberedFileReport` in Delivery Revalidation — command: `grep -n 'DeliveryQueue\|numbered' internal/cli/deliver_revalidate.go`; result: no match; `Revalidate` reads no Run Database row.
4. A file-name number comparison anywhere in the queue — command: `grep -rn 'ParseUint' internal/delivery internal/cli/deliver_*.go | grep -v _test`; result: only `check_rerun.go`, which reads GitHub check run ids, and `existingItemBranches` in `deliver_workflow.go`, which counts commits.

## ADR Alternatives

- ADR-0312: it rejected revalidating and rerunning an item branch that
  already holds a proven archive; the check lives inside Delivery
  Revalidation and so never runs for an Archived Item Resume.
- ADR-0367: written here; the design takes none of its rejected options (a
  default list, a wait like a Prerequisite Spec's, queue renumbering,
  merge-only detection).

## Integration Points

Git through `preflight.GitRunner` (`ls-tree`, `merge-base`, `diff`,
`rev-parse`, and one `fetch` from the configured remote when patterns are
declared); the Spec folders through `spec.Load`; the Run Database through the
workflow's store. Tests use temporary repositories with a bare remote on disk,
so no network is reached.

## Testing Approach

- task_01: `internal/config/numbered_files_test.go`, through
  `ResolveConfigProposal`.
- task_02: `internal/delivery/numbered_files_test.go` for the pure
  comparison; `internal/delivery/revalidate_test.go` with the existing
  `recordingItemRevalidator`, `newRevalidationTestEngine`,
  `newRetryRevalidationTestEngine` and `seedNoRunRetryItem` seams;
  `internal/delivery/park_class_test.go`.
- task_03: the park action contract in
  `internal/docscontract/command_documentation_test.go`, with its rejection
  test.
- task_04: `internal/cli/deliver_numbered_files_test.go` over temporary
  repositories built like `newDeliveryRevalidationGitFixture` (a bare
  `origin` on disk, item worktrees linked to the queue repository), a
  temporary Run Database from `store.Open`, and `runCLI` for `deliver status`.

Sabotages, each on a temporary copy, with the test and the input that sees it:

1. `IsNumberedFile` drops the leading-digit test:
   `TestIsNumberedFileNeedsAMatchAndALeadingDigit` fails on
   `drizzle/meta/_journal.json` under `[drizzle/]`.
2. The overlay never copies `numbered_files`:
   `TestNumberedFileDeclarationsAreReadAndValidated/project_replaces_user` fails.
3. Numbers compared as strings, not values:
   `TestFindNumberedFileCollisions/by_value` fails on `052_b.sql` against
   `0052_a.sql`.
4. Next free taken from the default branch only:
   `TestFindNumberedFileCollisions/next_free_above_every_claim` fails
   (`drizzle/0052` rendered, `drizzle/0054` expected).
5. Directory ignored: `TestFindNumberedFileCollisions/other_directory` fails
   on `drizzle/meta/0052_snapshot.json` against `drizzle/0052_a.sql`.
6. The engine ignores `NumberedFiles`:
   `TestRevalidationParksANumberedFileCollisionBeforeTheRun` fails (stage
   `running`, a `run` event).
7. The retry skips the check:
   `TestRetryOfAnItemWithoutARunRefusesWhileANumberedFileCollides` fails.
8. `ClassifyPark` drops the case: `TestEveryBlockerHasAParkClass` and
   `TestParkClassesAndNextActions/numbered-file-collision` fail.
9. Every queue item compared, later ones included:
   `TestRevalidateReportsANumberedFileCollision/ignored` fails on the later
   item that claims the same number.
10. Merged items compared: the same `/ignored` subtest fails on the merged
    item whose file is not on the default branch.
11. Completed Tasks' `creates:` counted: `/ignored` fails on the completed
    Task's path.
12. Reading the store before the declaration:
    `TestRevalidateReportsANumberedFileCollision/off` panics on its nil store.
13. The park action contract skips the new blocker:
    `TestTheParkActionContractRequiresTheNumberedFileCollision` fails when the
    sentence is removed from the first park document.

## Build Order

1. `delivery.numbered_files` and `IsNumberedFile` (task_01).
2. The pure comparison, the report, the park, the retry refusal and the Park
   Class (task_02; independent of 1, because it takes Numbered File paths the
   caller has already filtered).
3. `CONTEXT.md`, ADR-0367, the guides, the `deliver` reference and the park
   action contract (task_03; depends on: 1, 2, for `ClassifyPark`'s next
   action and the shared Behavior Surface Record).
4. `numberedFileReport` in Delivery Revalidation and the `deliver status`
   test (task_04; depends on: 1, 2, 3; the `deliver` guide and reference that
   describe this CLI surface come first).
5. QA gate (depends on: 4).

Specs 0355 and 0356, delivered before this Spec, may change
`internal/delivery/engine.go`, `internal/delivery/park_class.go`, the
`deliver` guide and reference, the Roundfix `SKILL.md` pair, the version
record, `CONTEXT.md` and the Behavior Surface Record; Spec 0352 may change
`internal/config/config.go`. Each Task rebases onto them, keeps to the
functions it names (`advanceItem`, `Retry`, `ClassifyPark`,
`applyConfigContent`, `Validate`, `Revalidate`) and re-runs the record
commands, which recompute versions and fingerprints.

## Risks & Considerations

- A plan that does not declare its migration path is seen only once its branch
  holds the file; ADR-0367 states the limit, and the guide tells adopters to
  declare migration files under `creates:`.
- A parked earlier item that is never delivered keeps parking later colliding
  items; the next action names it, and a new queue without it releases them.
- The fetch adds one network round trip per revalidation only for repositories
  that declare patterns; a failed fetch parks the item as `delivery-error`,
  as every other revalidation error does.
- Pattern breadth: a pattern over dated files would report dated names as
  collisions; the guide says to declare only directories numbered in sequence.

## Glossary

- adds: **Numbered File**
- adds: **Numbered File Collision**
- changes: **Delivery Revalidation**

## Decisions

- The declaration is a plain list of patterns, not entries with commands:
  nothing runs, so there is nothing else to declare. See ADR-0367.
- The comparison is pure and lives in `internal/delivery`; Git and the Run
  Database stay in the command workflow, as for the premise comparison.
- The declaration is read at the default-branch commit, like
  `delivery.derived_paths`, so an item branch cannot turn the check off.
- The park reuses the conflict Park Class: the cause is a predicted pull
  request conflict.
- Strict findings park first; a collision is reported once the plan is
  otherwise valid.
