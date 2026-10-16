---
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
prd: _prd.md
created: 2026-10-16
---

# The queue reports a gone owner, and a dismissal survives a fix — Technical Spec

## Executive Summary

All three behaviors attach to code that already exists. `runDeliverStatus`
reuses the owner proof `deliver resume` and `deliver retry` run
(`store.ProcessAlive`, then `OwnerProcessController.ProveOwner` against the
recorded start-time identity) and prints one `Owner:` line, without writing
the Run Database. One helper in `internal/cli/deliver.go` builds the
missing-queue refusal for `status`, `resume`, `retry` and `stop`; it reads a
new read-only store query, `DeliveryQueueItemWorktrees`, and compares resolved
paths as ADR-0268 does. `runReviewDisposeCommand` replaces the equality check
for `--dismiss` with an ancestry check and an anchor comparison read from Git,
and the disposition gains an optional `currentHead` field.

The trade-off accepted is ADR-0279's narrowness. Status reports only a
recorded owner proven gone, so a queue whose owner exited and released its
claim still prints no `Owner:` line; and a dismissal after a fix demands the
anchored lines at the same numbers, so a fix that only shifts them still needs
the dismissal first. In exchange there is no false report of a reused PID, no
write from a read command, and no dismissal that outlives a change to the code
it judged.

## Project Constraints

- Identifier strategy: not applicable — no identifier, table or column is
  added; the lookup reads the recorded item `worktree`, and `currentHead` is a
  commit name Git reports. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call; Git and the Run Database are read locally. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0279 governs every decision below:
  "Status stays read-only and never reclaims the owner record". ADR-0044:
  "anything short of proof keeps the block". ADR-0266: "a mismatch sends
  nothing and counts as proof that the recorded process ended". ADR-0174: "The
  disposition ledger stays shared, because every disposition already names its
  repository and reviewed head". ADR-0196: "Each finding must open with a
  `path:line` anchor". ADR-0272: "The current record keeps its role". ADR-0268: "Path
  comparison resolves the longest existing prefix of a path". ADR-0256: "A
  release waits for the skills that describe a changed surface". ADR-0244: "A
  Spec declares the domain terms it introduces". ADR-0259: "A test runs in parallel unless it
  names why it cannot". ADR-0271: "A test bounds by real time only what it
  asserts". ADR-0280: "A guide Task records its own surface fingerprint".
  ADR-0281: "A test replaces an executable it rebuilds instead of writing into
  it"; not applicable, because no test here rebuilds an executable.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH grant of 2026-10-16 under the cycle AE "Concedo" and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/review.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-16 at `9e723b82` by `go test -overlay` probes in
`internal/cli` (nothing written to the repository), each in a temporary home
and repository:

- A queue with one item at `checking` and a recorded owner PID `1<<30`, which
  no process can hold (`store.ProcessAlive` returned `false`), printed exit 0,
  empty stderr and exactly
  `0300-example\tchecking\t-\t-\n`, `Limits: deadline none, retries per item none, concurrency 1, tokens none\n`
  and `Usage: no Runs recorded\n`. Nothing names the owner.
- With the queue recorded for the main checkout and an item worktree created
  by `git worktree add` and recorded through
  `RecordDeliveryQueueItemWorktree`, `deliver retry <slug>` run in the item
  worktree exited 2 with `Retry refused`, `Reason:` and
  `retry Delivery Queue item "<slug>": Delivery Queue for repository "<item worktree>" does not exist`,
  then the `No side effects:` block; `deliver status` there exited 2 with
  `Preflight failed` and `Delivery Queue for repository "<item worktree>" does not exist`.
  `loaded.GitRoot` was the unresolved temporary path (`/var/...`), not the
  path Git resolves (`/private/var/...`).
- A findings record whose finding carries the anchor `review.txt:1-1`, after a
  commit that adds `other-fix.txt` only, refused `review dispose F1 --dismiss
  --evidence "not a defect"` with exit 2 and
  `roundfix: review dispose refused: dismissal requires current HEAD "<head>" to equal reviewed head "<reviewed>"`.

## System Architecture

- `internal/cli/deliver.go`: `runDeliverStatus` calls `deliveryOwnerLine`
  after the `Warning:` lines and before the `Park:` lines. A new helper,
  `missingDeliveryQueue`, builds the refusal `status`, `resume`, `stop` use in
  place of their `fmt.Errorf("Delivery Queue for repository %q does not
  exist", ...)`, and `runDeliverRetry` reads the queue once after opening the
  store and refuses through `printDeliverRetryRefusal` with the same helper
  before calling the engine.
- `internal/store/delivery.go`: `DeliveryQueueItemWorktrees` lists every
  recorded item worktree. It writes nothing and needs no migration.
- `internal/cli/review.go`: the `--dismiss` branch of
  `runReviewDisposeCommand` keeps the equality fast path and otherwise calls
  `proveDismissalAnchor`; `reviewFindingDisposition` gains `CurrentHead`.
- Nothing in `internal/delivery`, the review correction proof
  (`deliver_review_correction.go`), the Reviewer Lineage
  (`review_lineage.go`) or review reuse changes: they read a dismissal's
  `disposition`, `evidence` and `headCommit`, which keep their meaning.

## Implementation Design

### Interfaces

```go
// internal/store/delivery.go
type DeliveryQueueItemWorktree struct {
	GitRoot, SpecSlug, Worktree string
}
func (store *Store) DeliveryQueueItemWorktrees(ctx context.Context) ([]DeliveryQueueItemWorktree, error)

// internal/cli/deliver.go
func deliveryOwnerLine(ctx context.Context, queue store.DeliveryQueue, controller OwnerProcessController, stderr io.Writer) string
func missingDeliveryQueue(ctx context.Context, runStore *store.Store, gitRoot string) error

// internal/cli/review.go
type reviewFindingDisposition struct {
	// existing fields unchanged, in their order
	CurrentHead string `json:"currentHead,omitempty"`
}
func proveDismissalAnchor(ctx context.Context, git preflight.GitRunner, root, reviewed, head string, finding reviewFinding) error
```

### Data Models

No table, column or migration. The disposition ledger line gains the optional
`currentHead` key, written last and only when the current HEAD differs from
`headCommit`. Every existing line stays valid and decodes as before.

### API Contracts

1. API Contract: owner report. When `queue.OwnerPID > 0` and at least one
   item's stage is neither `merged` nor `parked` (the rule
   `ReleaseIdleDeliveryQueueOwner` applies), status proves the owner: if
   `store.ProcessAlive(pid)` is false the state is `is not running`; otherwise
   `ProveOwner(ctx, pid, identity)` decides: `nil` with the PID still alive
   adds no line, `nil` with the PID now absent is `is not running`, an error
   matching `store.ErrOwnerProcessIdentityUnproven` is `has a different
   process identity`, and any other error adds no line and writes
   `roundfix: Delivery Queue owner PID <pid> could not be checked: <error>` to
   stderr. A state prints, on stdout, after the `Warning:` lines and before the
   `Park:` lines,
   `Owner: Delivery Queue owner PID <pid> <state>; run roundfix deliver resume`.
   A missing controller is treated as an error that adds no line. Status opens
   the Run Database read-only and never calls `ReleaseDeliveryQueueOwner`; its
   exit code stays 0 and every other line keeps its bytes and order.
2. API Contract: owning checkout. `DeliveryQueueItemWorktrees` returns every
   item row with a non-empty `worktree`, ordered by `git_root` and
   `position`. When `DeliveryQueue(ctx, gitRoot)` finds nothing,
   `missingDeliveryQueue` compares `config.CanonicalizeResolvedPath(gitRoot)`
   with the canonical form of each recorded worktree whose `GitRoot` differs
   from `gitRoot`; the first match yields
   `Delivery Queue for repository "<gitRoot>" does not exist; this checkout is the item worktree of <item slug> in the Delivery Queue of "<owning git root>"; run the command in "<owning git root>"`.
   Without a match, or when the lookup fails, it yields today's
   `Delivery Queue for repository "<gitRoot>" does not exist`. `status`,
   `resume` and `stop` print it through `printDeliverFailure`; `retry` prints
   `retry Delivery Queue item "<slug>": ` followed by it through
   `printDeliverRetryRefusal`, before the engine runs, so no item changes and
   no owner starts. All four exit 2.
3. API Contract: dismissal after a fix. With `--dismiss`, a current HEAD equal
   to the reviewed head is accepted as today and writes no `currentHead`.
   Otherwise, in this order: `git merge-base --is-ancestor <reviewed> HEAD`
   exiting 1 refuses with
   `dismissal requires current HEAD "<head>" to equal or descend from reviewed head "<reviewed>"`,
   and any other failure refuses with `check dismissal ancestry: <error>`; a
   finding whose record holds no `anchor` refuses with
   `dismissal at HEAD "<head>" after reviewed head "<reviewed>" requires finding "<id>" to record an anchor`;
   an anchor whose path is missing at either commit, whose lines
   `startLine`..`endLine` (1-based, inclusive) do not all exist at both
   commits, or whose lines differ byte for byte at the same numbers refuses
   with
   `dismissal at HEAD "<head>" after reviewed head "<reviewed>" requires the anchor of finding "<id>" (<path>:<start>-<end>) to be unchanged`.
   A start line of 0 compares the whole blob. Commits are read through Git
   (`rev-parse --verify <commit>:<path>` and `cat-file blob`), never from the
   working tree. An accepted dismissal appends and prints the line with
   `headCommit` set to the reviewed head and `currentHead` set to HEAD. Every
   refusal starts stderr with `roundfix: review dispose refused: `, exits 2
   and appends nothing; every other dispose rule and message is unchanged.
4. API Contract: everything else stays. The `deliver` and `review dispose`
   help text, `deliver status` in a repository that owns its queue, the
   missing-queue refusal outside any item worktree, `--fixed-by`, the ledger's
   existing keys and the readers of the ledger keep their bytes and answers.

### Surface Transcripts

Each block uses literal tab characters between the columns of an item row.

1. Surface Transcript: status names a recorded Queue Owner that is not
   running.

   ```transcript
   $ roundfix deliver status
   stdout:
   0300-example	checking	-	-
   Owner: Delivery Queue owner PID <pid> is not running; run roundfix deliver resume
   Limits: deadline none, retries per item none, concurrency 1, tokens none
   Usage: no Runs recorded
   stderr:
   exit: 0
   ```

2. Surface Transcript: retry in an item worktree names the owning checkout.

   ```transcript
   $ roundfix deliver retry 0300-example
   stdout:
   stderr:
   Retry refused

   Reason:
     retry Delivery Queue item "0300-example": Delivery Queue for repository "<item worktree>" does not exist; this checkout is the item worktree of 0300-example in the Delivery Queue of "<checkout>"; run the command in "<checkout>"

   No side effects:
     Roundfix did not change the Delivery Queue item, start a queue owner, commit, or push.
   exit: 2
   ```

3. Surface Transcript: status in an item worktree names the owning checkout.

   ```transcript
   $ roundfix deliver status
   stdout:
   stderr:
   Preflight failed

   Reason:
     Delivery Queue for repository "<item worktree>" does not exist; this checkout is the item worktree of 0300-example in the Delivery Queue of "<checkout>"; run the command in "<checkout>"

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix deliver status --help' for usage.
   exit: 2
   ```

4. Surface Transcript: a dismissal after a fix that leaves the anchor.

   ```transcript
   $ roundfix review dispose F1 --dismiss --evidence "not a defect"
   stdout:
   {"repository":"<repository>","headCommit":"<reviewed head>","finding":"F1","text":"<finding text>","disposition":"dismissed","evidence":"not a defect","recordedAt":"<timestamp>","currentHead":"<current head>"}
   stderr:
   exit: 0
   ```

5. Surface Transcript: a dismissal after a fix that changed the anchor.

   ```transcript
   $ roundfix review dispose F1 --dismiss --evidence "not a defect"
   stdout:
   stderr:
   roundfix: review dispose refused: dismissal at HEAD "<head>" after reviewed head "<reviewed head>" requires the anchor of finding "F1" (<path>:<start>-<end>) to be unchanged
   exit: 2
   ```

### Invariants

```text
1. Status never writes the Run Database and never reclaims an owner record.
2. A live PID whose recorded identity matches, or cannot be read, is never reported gone.
3. A queue without a recorded owner, or with no item that can advance, prints no Owner line.
4. Outside an item worktree the missing-queue refusal keeps the delivery base's bytes.
5. A refused dismissal appends nothing to the ledger.
6. An accepted dismissal keeps headCommit equal to the reviewed head.
7. A dismissal at the reviewed head writes the same line as before.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, User Story 1, Core Feature 1 → `deliveryOwnerLine` (API Contract 1, Surface Transcript 1)
- Goal 2, User Story 2, Core Feature 2 → `proveDismissalAnchor`, `CurrentHead` (API Contract 3, Surface Transcripts 4 and 5)
- Goal 3, User Story 3, Core Feature 3 → `missingDeliveryQueue`, `DeliveryQueueItemWorktrees` (API Contract 2, Surface Transcripts 2 and 3)
- Core Feature 4 → the `deliver` and `review` user guides, the Roundfix skill's `deliver` and `review` references, `CONTEXT.md`
- Success Metric 1 → API Contract 1, `internal/cli/deliver_owner_status_test.go` and the QA gate
- Success Metric 2 → API Contract 3, `internal/cli/review_dismiss_after_fix_test.go` and the QA gate
- Success Metric 3 → API Contracts 2 and 4, `internal/cli/deliver_item_worktree_refusal_test.go` and the QA gate

## Absence Evidence

1. Owner report in status — command: `grep -n 'ProveOwner\|ProcessAlive' internal/cli/deliver.go`; result: lines 156, 160, 162, 315, 453, 458, 460 and 512 at `9e723b82`, all in `retry`, `start`, `resume` and `stop`; none inside `runDeliverStatus` (lines 347-405).
2. `DeliveryQueueItemWorktrees` — command: `grep -n 'SELECT' internal/store/delivery.go`; result: the item rows are read only by `git_root` (`WHERE git_root = ?`); no query lists worktrees across queues.
3. Owning-checkout refusal — command: `grep -rn 'is the item worktree of' internal`; result: no match at `9e723b82`.
4. Dismissal after a fix and `currentHead` — command: `grep -rn 'to equal reviewed head\|json:"currentHead' internal/cli`; result: only the refusal at `internal/cli/review.go:1012`; no ledger key `currentHead`.

## ADR Alternatives

- ADR-0279: this Spec's decision; it rejects reporting a queue without a recorded owner, a PID-only proof, reclaiming from status, a whole-file or any-descendant dismissal rule, a `--head` flag and resolving the owning queue automatically.
- ADR-0044: no rejected alternative is taken; status reports only a proven death and reclaims nothing.
- ADR-0266: no rejected alternative is taken; status sends only the existence probe and decides reuse by the recorded start identity.
- ADR-0174: no rejected alternative is taken; the ledger stays shared and each line still names its repository and reviewed head.
- ADR-0196: no rejected alternative is taken; validation is untouched and the dismissal reads the anchor it recorded.
- ADR-0272: the rejected `review dispose --head` flag stays rejected; dispose still reads the current record.
- ADR-0268: no rejected alternative is taken; the lookup compares resolved paths and the queue records no new path.
- ADR-0256: no rejected alternative is taken; task_04 updates both covering skill references.
- ADR-0244: no rejected alternative is taken; both revised terms follow its declaration contract.
- ADR-0259: no rejected alternative is taken; every new test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0280: no rejected alternative is taken; each guide Task declares the Behavior Surface Record and names its record step.
- ADR-0281: not applicable; no test rebuilds an executable.

## Integration Points

None outside the repository. Git is read through the command's Git runner in
the checkout `review dispose` runs in, and the Run Database through the store.

## Testing Approach

New tests call `t.Parallel()` first and run in temporary homes and
repositories with the existing fixtures: `newImplementWorkspace`,
`newParkedDeliveryQueueForRetry`'s pattern, `fakeDeliveryOwnerProcesses`,
`updateCommandDependenciesForTest`, `newReviewDispositionFixture`,
`writeDispositionRecord` and `gittest`. A gone owner uses the PID `1<<30`,
which no process can hold on the supported systems, never a reaped child's
PID, because the suite wraps the PID counter (ADR-0266). A live owner uses the
test's own PID with the fake controller deciding its identity. An item
worktree is made with `git worktree add` and recorded through
`RecordDeliveryQueueItemWorktree`, and the resolved-path case records it
through a symbolic link to its parent.

No existing test is expected to change. Every existing status test records no
owner; no test pins the missing-queue refusal; and
`TestReviewDisposeRefusesADismissalAtAMovedHead` keeps passing because its
finding records no anchor, so the dismissal is now refused by the anchor rule
instead of the equality rule. A Task that finds otherwise names the test and
the contract it moves to.

- `internal/cli/deliver_owner_status_test.go`: Surface Transcript 1, the
  different identity, a proven owner, no recorded owner, no advanceable item,
  an unreadable identity, and the owner record left in place.
- `internal/store/delivery_item_worktrees_test.go` and
  `internal/cli/deliver_item_worktree_refusal_test.go`: the query, Surface
  Transcripts 2 and 3, `resume` and `stop` in an item worktree starting
  nothing, the resolved-path match, and the unchanged refusal elsewhere.
- `internal/cli/review_dismiss_after_fix_test.go`: Surface Transcripts 4 and
  5, a moved anchor, a deleted anchored file, a missing anchor, a head that
  does not descend, the whole-file anchor, the unchanged line at the reviewed
  head, and the review correction proof accepting a dismissal recorded after
  the fix.

## Build Order

1. Owner report in status, its tests and the `deliver` guide paragraph (no
   dependency).
2. Owning-checkout refusal, the store query, tests and the `deliver` guide
   paragraph (depends on: 1, which edits the same command file and guide).
3. Dismissal after a fix, `currentHead`, tests and the `review` guide
   paragraph (depends on: 2, because both re-record the Behavior Surface
   Record).
4. Roundfix skill `deliver` and `review` references, mirrors, version,
   Behavior Surface Record and glossary (depends on: 1, 2, 3).
5. QA gate (depends on: 4).

No other Spec is a prerequisite. Rebase on main after the other cycle AH plans
merge and run the strict check again.

## Risks & Considerations

- A process that exits between the existence probe and the identity proof is
  read again, as `resume` does, so it is reported gone rather than running.
- A queue whose owner was stopped with `deliver stop` prints no `Owner:` line;
  the guide says that `deliver resume` restarts it.
- A finding whose anchor names lines the fix shifted is refused; the guide
  says to dismiss before committing such a fix.
- `PID 1<<30` is above every supported system's PID ceiling (macOS 99998,
  Linux 4194304); on another platform the gone-owner tests would need another
  absent PID.

## Glossary

- changes: **Queue Owner**
- changes: **Review Finding Disposition**

## Decisions

- Status proves the owner with the existing start-time proof and never reclaims it. See ADR-0279.
- A dismissal after a fix compares the anchored lines at the same numbers at both commits. See ADR-0279.
- The item worktree refusal names the owning checkout. See ADR-0279.
- `currentHead` is optional and written last, so every ledger reader keeps its answer.
