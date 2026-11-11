---
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
status: active
created: 2026-11-10
---

# A Machine Verification Capacity bounds Verifications across repositories — Technical Spec

## Executive Summary

Two repositories' Verifications that run at once on one machine compete for
its cores and fail each other's time-sensitive tests. This Spec adds the
**Machine Verification Capacity**: `runs.max_verifying` in User Config, `1`
by default, `0` off. Every Verification the Daemon starts takes its unit of
the Run's own Verification Capacity first and one unit of the machine's
second, and the Delivery Queue's repository gate takes one unit for its
repository Verification and Linux Lane together. The units live as entry
files under `verification-capacity/` in the Roundfix Home, changed only under
one advisory lock and served in ticket order; a holder is removed only when
ADR-0373's rule proves its owner gone. While a Run's Verification waits, its
Run Budget does not run, a Stop Request still ends the wait cleanly, and
`runs list`, `deliver status` and `doctor` name the holder and the waiters.
ADR-0384 records the decision and amends ADR-0056.

The trade-off accepted: with the default, a machine verifies one repository
at a time, so a Run can wait for another repository's whole suite; the
Agents keep working, the wait costs no Run Budget, and the owner can raise or
turn off the bound in one key. Entry files instead of a Run Database table
keep every older binary able to open the shared Run Database, at the price of
a directory that only Roundfix's own lock protects.

## Project Constraints

- Identifier strategy: not applicable — no domain identifier is minted; the
  ticket is a local sequence number in the Roundfix Home, and Run IDs and
  Spec slugs are the existing ones. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0384 (new) records the decision:
  "A Machine Verification Capacity bounds the Verifications of every
  repository". ADR-0056 is amended by task_04: "The capacity is owned by one
  Task cycle, not by a machine-wide service, so it does not coordinate other
  Runs" stays true of the Run's own capacity, and an outer machine capacity
  now wraps it. ADR-0373 applies unchanged: "A proven start-identity mismatch
  at a recorded owner PID means the owner is gone". ADR-0044 applies
  unchanged: "a provably dead owner makes the lock orphaned". ADR-0304 applies
  unchanged: "`runs.slot_wait` sits beside `runs.max_active`". ADR-0365 applies
  unchanged: "The turn is a file.". ADR-0342 applies unchanged: the Integrated
  Repository Verification it defines is one of the gated steps. ADR-0313,
  cited by ADR-0373, applies unchanged. Not applying: ADR-0025, ADR-0038 and
  ADR-0034, which ADR-0056 refines or cites; ADR-0311, which cites ADR-0304
  and sets the Delivery Queue's check wait; and ADR-0182, which cites
  ADR-0342 and decides the facts a Task settles on. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — the cycle BB batch 2 grant quoted in
  `_authorization.md`; the Governed Paths are the Roundfix skill files tasks
  02 and 03 change, through the skill sync and version record. No Makefile, workflow,
  Go module, lint or formatter change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0374-a-machine-verification-capacity-bounds-verifications-across-repositories/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/runs.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured at `60ac8142`.

1. `internal/daemon/task_engine.go`: `fairVerificationGate` is the Run's
   Verification Capacity, created per Task cycle by `newVerificationGate`.
   `acquireVerificationCapacity` is its only door: `verifyTaskPreWork` (the
   pre-work probe) calls it, and `runTaskVerificationRequest` calls it for
   every attempt request — `verifyTask` (Task attempts, with the repository
   Verification appended by `taskWithRequiredRepositoryVerification`, and the
   exclusive retry), `verifyRepositoryPrecondition`,
   `verifyIntegratedRepository` and `runQARepositoryVerification`. While it
   queues, `watchStopRequest` turns a Stop Request into
   `stopQueuedVerification`, which starts no command.
2. The gate is an in-process object; two Implement Runs, which are two
   processes, never see each other's attempts. ADR-0056: "The capacity is
   owned by one Task cycle, not by a machine-wide service, so it does not
   coordinate other Runs".
3. `internal/cli/deliver_workflow.go` `commandDeliveryWorkflow.Gate` trims the
   shared caches, runs `defaults.verification` through `daemon.ExecVerifier`
   in the item worktree, then the Linux Lane, in the Queue Owner process.
4. `internal/config/config.go`: `Runs` holds `MaxActive` and `SlotWait`;
   `runsOverlay` decodes `max_active` and `slot_wait`; the Project Config
   reader removes both with `warnIgnoredProjectSetting`; `validateConfig`
   refuses negatives; `defaultConfigYAML` writes `runs:` with
   `max_active` only for the user scope. A User Config with
   `runs.max_verifying` is refused as an unknown key.
5. `taskCycleBudget` holds a deadline that `renew` moves at each Task
   settlement, `expiredAt` checks, and `watch` turns into a cycle
   cancellation when its timer fires at an unchanged deadline. No
   Verification command has a time limit of its own.
6. `internal/store`: `OwnerGone` answers ADR-0373's question;
   `acquireWriteLock`/`releaseWriteLock` (flock on Unix, `LockFileEx` on
   Windows) take an advisory lock on an open file with context cancellation.
7. `roundfix runs list`, `roundfix deliver status` and `roundfix doctor` say
   nothing about another repository's Verification; a Run waiting for the
   Run's own capacity shows `Verifying`.

## System Architecture

- `internal/store/machine_verification.go` (new): the entry files, the lock,
  the ticket sequence, `AcquireMachineVerification`, the lease and
  `ReadMachineVerification`, all platform-neutral over the existing lock
  helpers.
- `internal/config`: the key, its default, validation, Project Config warning
  and template line.
- `internal/daemon`: the `MachineVerificationGate` dependency, the second
  acquisition in `acquireVerificationCapacity`, the step label, the Run
  Budget pause, the wait event and progress line.
- `internal/cli`: the production gate for Implement Runs, the Delivery Queue
  gate's acquisition, and the capacity block of `runs list`, `deliver status`
  and `doctor`.
- Documentation: the configuration and command guides, four Roundfix
  references, `CONTEXT.md` and ADR-0056's amendment notice.

## Implementation Design

### Interfaces

```go
// internal/store
const (
	MachineVerificationHolderRun      = "run"
	MachineVerificationHolderDelivery = "delivery"
)

type MachineVerificationHolder struct {
	Kind       string
	RunID      string
	Repository string
	SpecSlug   string
	Step       string
}

type MachineVerificationEntry struct {
	Ticket        int64
	State         string // "waiting" or "holding"
	Holder        MachineVerificationHolder
	OwnerPID      int
	OwnerIdentity string
	Since         time.Time // entry written
	AcquiredAt    time.Time // zero while waiting
	Path          string
}

type MachineVerificationWait struct {
	Capacity int
	Holding  []MachineVerificationEntry
	Ahead    []MachineVerificationEntry
}

type MachineVerificationLease struct{ /* unexported */ }

func (lease *MachineVerificationLease) Release() error
func (lease *MachineVerificationLease) Waited() time.Duration

func AcquireMachineVerification(ctx context.Context, homeDir string, holder MachineVerificationHolder, capacity int, onWait func(MachineVerificationWait)) (*MachineVerificationLease, error)

type MachineVerificationView struct {
	Entries    []MachineVerificationEntry // ticket order
	Unreadable []string                   // paths
}

func ReadMachineVerification(ctx context.Context, homeDir string) (MachineVerificationView, error)
func DescribeMachineVerificationHolder(holder MachineVerificationHolder) string

// internal/daemon
type MachineVerificationGate interface {
	Acquire(ctx context.Context, holder store.MachineVerificationHolder, onWait func(store.MachineVerificationWait)) (release func(), err error)
}
// Dependencies gains: MachineVerification MachineVerificationGate
```

### Data Models

Directory `<home>/.roundfix/verification-capacity/`, mode `0700`:

- `lock`: the advisory lock file, opened per acquisition and held only while
  the directory is read or changed.
- `sequence`: the last ticket issued, decimal, rewritten under the lock.
- `<ticket>.json`, the ticket zero-padded to 20 digits, written to a
  temporary file in the directory and renamed into place:

```json
{
  "schema": "roundfix/machine-verification-entry/v1",
  "ticket": 7,
  "state": "holding",
  "kind": "run",
  "run_id": "run_20261110T215130Z_fa9f763e0a43c500",
  "repository": "/Users/marcio/dev/roundfix",
  "spec": "0371-a-results-row-after-a-blank-line-is-refused-naming-the-line",
  "step": "Task task_02 Verification attempt 1",
  "owner_pid": 81234,
  "owner_identity": "darwin:1762814490.123456",
  "since": "2026-11-10T22:41:10Z",
  "acquired_at": "2026-11-10T22:41:30Z"
}
```

`acquired_at` is absent while `state` is `waiting`. `run_id` is empty for a
`delivery` holder. No Run Database table, column or schema version changes.

### API Contracts

1. **Config key** (`internal/config/config.go`, task_01).
   - `Runs` gains `MaxVerifying int`, built-in `defaultRunsMaxVerifying = 1`;
     `runsOverlay` gains `MaxVerifying *int` with tag `yaml:"max_verifying"`,
     applied from User Config only, beside `MaxActive`.
   - The Project Config reader removes `runs.max_verifying` and calls
     `warnIgnoredProjectSetting("runs.max_verifying")`, which prints
     `config: runs.max_verifying in Project Config is ignored; set runs.max_verifying in User Config`.
   - `validateConfig` refuses a negative value with
     `runs.max_verifying must be greater than or equal to 0`; a non-integer
     value fails strict decoding as every integer key does.
   - `defaultConfigYAML` writes, for the user scope only, right after the
     `max_active` line:

     ```yaml
       # Maximum Verifications at once across repositories; 0 disables the bound.
       max_verifying: 1
     ```

     so `roundfix init --scope user` and `roundfix setup` write it. The
     project template is unchanged.
2. **Acquisition** (`internal/store/machine_verification.go`, task_01).
   `AcquireMachineVerification`:
   - With `capacity <= 0` it returns a lease whose `Release` does nothing,
     and creates no directory and no file.
   - Otherwise it creates the directory, takes the lock, increments
     `sequence`, writes its own entry as `waiting` with its own PID,
     `OwnerProcessIdentity(ctx, os.Getpid())` (empty when it cannot be read)
     and `since`, and releases the lock.
   - Each round takes the lock and reads every `*.json` entry. An entry that
     cannot be read or parsed, or whose `schema` differs, is skipped and its
     path kept for the warning; it is neither counted nor removed. Every other
     entry but its own whose `OwnerGone(ctx, pid, identity)` answers gone is
     removed. Then, with `holding` the count of `holding` entries and `ahead`
     the `waiting` entries with a lower ticket, it rewrites its own entry as
     `holding` with `acquired_at` when `holding + ahead < capacity`, releases
     the lock and returns. Otherwise it releases the lock, calls `onWait`
     once (on the first refused round only) with the capacity, the holding
     entries and the entries ahead, and waits for the next round.
   - A round follows every 2 seconds (`machineVerificationPollInterval`), and
     at once when another acquisition in the same process releases or
     removes an entry (a package-level broadcast), so attempts of one process
     never wait for the poll.
   - When `ctx` ends while waiting, it removes its own entry under the lock,
     with a context that `ctx`'s cancellation does not cancel and that is
     bounded by 30 seconds, and returns `ctx`'s error wrapped with
     `wait for the Machine Verification Capacity`.
   - `Release` removes its own entry under the lock with the same
     uncancelled, bounded context, wakes the process's waiters, and is safe
     to call twice; its error names the entry path.
   - `Waited` is the time from the first refused round to the grant, zero for
     an immediate grant.
3. **Reader** (`internal/store/machine_verification.go`, task_01).
   `ReadMachineVerification` takes no lock and writes nothing: a missing
   directory is an empty view; it returns the readable entries in ticket
   order without those `OwnerGone` proves gone, and the unreadable paths.
   `DescribeMachineVerificationHolder` returns
   `Run <run id> of <spec> in <repository>` for a `run` holder and
   `Delivery Queue item <spec> in <repository>` for a `delivery` holder.
4. **Daemon acquisition** (`internal/daemon`, task_02).
   - `verificationAttemptRequest` gains `Step string`. The constructions set
     it: `verifyTaskPreWork` `Task <id> pre-work probe`;
     `verifyRepositoryPrecondition` `Task <id> repository precondition`;
     `verifyIntegratedRepository` `Integrated Repository Verification`;
     `runQARepositoryVerification` `QA Task <id> repository Verification`.
     `verifyTask` leaves it empty, and the acquisition uses
     `Task <id> Verification <identity>`, where `<identity>` is the request's
     `identity()` (`attempt 1`, `attempt 1 retry 1`).
   - `acquireVerificationCapacity` keeps its Run-capacity path unchanged, then,
     when `engine.deps.MachineVerification` is not nil, calls its `Acquire`
     with holder kind `run`, the Run ID, the Spec slug and the step (the
     gate fills the repository), under the same
     `watchStopRequest` wrapping the Run-capacity wait uses. When the machine
     wait ends by a Stop Request it releases the Run unit and returns
     `stopQueuedVerification`; on any other error it releases the Run unit
     and returns the error wrapped with
     `acquire Machine Verification Capacity for run %q Task %s attempt %d retry %d`.
     The returned release gives the machine unit back first, then the Run
     unit. The context and Stop Request re-checks after acquisition stay
     after both units are held.
   - The first time `onWait` runs for an attempt, the engine pauses the Run
     Budget, publishes one `daemon.status` Run Event through
     `publishDaemonEvent` with summary
     `<step> waits for the Machine Verification Capacity: <held> of <capacity> held by <holders>.`
     where `<holders>` lists `<holder description>, <step>` joined by `; `,
     and payload `{"machine_verification":"waiting","step":…,"capacity":…,"held":…,"waiting_ahead":…,"holders":[…]}`,
     and writes the same summary as one line to `Progress`. When the
     acquisition returns after a wait, the Run Budget resumes and one
     `daemon.status` Run Event says
     `<step> took the Machine Verification Capacity after waiting <seconds>s.`
     with payload `{"machine_verification":"acquired","step":…,"waited_seconds":…}`.
5. **Run Budget pause** (`internal/daemon/task_engine.go`, task_02).
   `taskCycleBudget` gains `pause(now) (resume func(now time.Time))` with a
   pause count and the time the count left zero. While the count is above
   zero `expiredAt` answers false and `watch` does not cancel when its timer
   fires; it waits for the next change. When the count returns to zero, a
   set and unexpired deadline moves later by the paused length and `changed`
   is signalled. A `renew` while paused sets the deadline from the settlement
   time and restarts the paused length from that time. A disabled budget
   ignores pauses.
6. **Implement wiring** (`internal/cli`, task_02). The Implement Command gives
   the engine a `MachineVerificationGate` over
   `store.AcquireMachineVerification` with the loaded Roundfix Home,
   `Runs.MaxVerifying` and, as the holder's repository, the root of the
   checkout the Implement Run was started from (the preflight's Git root),
   never the Run Worktree's path. Its release writes
   `roundfix: warning: <release error>` to standard error when `Release`
   fails and never fails the Task.
7. **Delivery Queue gate** (`internal/cli/deliver_workflow.go`, task_03).
   `commandDeliveryWorkflow.Gate`, after the shared cache trim and before the
   repository Verification, acquires holder kind `delivery` with the
   checkout's repository root, the Spec slug and step `repository gate`, and
   releases after the Linux Lane, on every return path. On its first wait it
   writes to the queue log
   `roundfix: Delivery Queue item <spec>: waiting for the Machine Verification Capacity; <held> of <capacity> held by <holders>`.
   An acquisition error is returned as the gate's error wrapped with
   `wait for the Machine Verification Capacity`.
8. **Status surfaces** (`internal/cli`, task_03). One renderer prints, from
   `ReadMachineVerification` and the reader's own `runs.max_verifying`:

   ```text
   Machine Verification Capacity: <held> of <capacity> held, <waiting> waiting
     holding: <holder description>, <step>, since <acquired_at RFC3339 UTC>
     waiting <n>: <holder description>, <step>, since <since RFC3339 UTC>
     unreadable: <path>
   ```

   with `<capacity>` replaced by `off` (`Machine Verification Capacity: off, <held> held, <waiting> waiting`)
   when the reader's value is `0`, holders in ticket order, then waiters in
   ticket order numbered from 1, then unreadable paths.
   - `roundfix runs list` prints the block on standard output after the
     listing, both formats, including after `No Runs found.`, only when the
     view has an entry or an unreadable path.
   - `roundfix deliver status` prints it right before the `Limits:` line,
     under the same condition.
   - `roundfix doctor` prints one `verification-capacity:` result after
     `disk-space:` through a new `verificationCapacity` dependency:
     `ok (runs.max_verifying <n>; <held> held, <waiting> waiting)`, or with
     `off` in place of `runs.max_verifying <n>` when it is `0`, followed by
     `; holding: <holder description>, <step>` and
     `; waiting <n>: <holder description>, <step>` for each entry; with an
     unreadable path it is `warn` with `; unreadable entry <path>` and next
     action `inspect <path> and remove it if no Roundfix process wrote it`.
     A read failure is `warn` naming the error. Doctor stays read-only.
9. **Documentation** (tasks 02 to 04). task_02: `docs/user-guide/configuration.md` gains
   `## Machine Verification Capacity` after `## Active Implement Run ceiling`:
   the key, its default, `0`, User Config only with the warning, which steps
   take it and which do not, the arrival order, the proof that frees a dead
   holder, the Run Budget pause, the per-Home scope, and the sentence
   "After the update every machine verifies one repository at a time; set
   `runs.max_verifying` to a higher number, or to `0`, to change that." The
   `## Task and Verification capacities` paragraph that says Verification
   Capacity is not a machine-wide lock names the new section; the canonical
   Roundfix reference `implement.md` says the same, synced to `skills/`.
   task_03: `docs/user-guide/commands/runs.md`, `deliver.md` and `doctor.md`
   show their block or line, and the canonical references `runs.md`,
   `deliver.md` and `setup.md` say the same, synced to `skills/`. task_04:
   `CONTEXT.md` adds **Machine Verification Capacity** and changes
   **Verification Capacity**, **Waiting for Verification** and
   **Run Budget** with the texts of API Contract 10. ADR-0056 gains its
   amendment notice.
10. **Glossary texts** (task_04), exactly:
    - **Machine Verification Capacity**: "The maximum number of Verifications
      that run at once on one machine across every repository sharing the
      Roundfix Home, configured by `runs.max_verifying` in User Config, `1`
      by default and off at `0`. It is taken by every Verification the Daemon
      starts, after its Run's Verification Capacity, and by the repository
      gate of a Delivery Queue; waiters take it in arrival order, and only a
      proof that its owner is gone frees a holder (ADR-0384)." The closing `(ADR-0384)`
      is written as the glossary's usual Markdown link to
      `docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md`.
      `_Avoid_: machine-wide test lock, core budget, Run ceiling`
    - **Verification Capacity**: keeps its first sentence and adds "An attempt
      holds it while it waits for the Machine Verification Capacity and gives
      it back last." Its `_Avoid_` line becomes
      `_Avoid_: Task Capacity, worktree concurrency, Machine Verification Capacity`.
    - **Waiting for Verification**: "before the Task acquires Verification
      Capacity" becomes "before the Task holds its Run's Verification Capacity
      and the Machine Verification Capacity".
    - **Run Budget**: adds "It does not run while a Verification of the Run
      waits for the Machine Verification Capacity."

### Surface Transcripts

The fixture is a temporary `HOME` with User Config `runs:\n  max_verifying: 1\n`
and no Run Database, run from a temporary Git repository. The test writes two
entry files in the format of Data Models with its own PID: a `delivery`
holder for `0001-holder` in `<repository one>`, step `repository gate`,
`acquired_at` `2026-11-10T22:41:30Z`; and a `run` waiter
`run_20261110T221208Z_946b774b8b696f09` of `0002-waiter` in
`<repository two>`, step `Task task_03 Verification attempt 2`, `since`
`2026-11-10T22:50:02Z`. `<repository one>` and `<repository two>` are the
test's temporary repository paths.

1. Surface Transcript: runs list names the holder and the waiter.

   ```transcript
   $ roundfix runs list --all --format plain
   stdout:
   No Runs found.
   Machine Verification Capacity: 1 of 1 held, 1 waiting
     holding: Delivery Queue item 0001-holder in <repository one>, repository gate, since 2026-11-10T22:41:30Z
     waiting 1: Run run_20261110T221208Z_946b774b8b696f09 of 0002-waiter in <repository two>, Task task_03 Verification attempt 2, since 2026-11-10T22:50:02Z
   stderr:
   exit: 0
   ```

### Invariants

1. No Daemon Verification command and no Delivery Queue repository gate
   starts while its holder has no `holding` entry, when the capacity is above
   zero.
2. No process holds a machine unit while it waits for its Run's own unit.
3. An entry is removed only by its owner or by a round that proves its owner
   gone; never by its age.
4. A waiter is granted only when `holding + ahead < capacity`.
5. `ReadMachineVerification` and the three status surfaces write nothing.
6. With `runs.max_verifying: 0` no file is written under the Roundfix Home.

### Declared breaks

- Every adopter's Verifications become one at a time per machine after the
  update; `runs.max_verifying: 0` restores today's behavior.
- A Run's Verification can wait for another repository's; its Run Budget
  pauses for that wait.
- `roundfix doctor` prints one more line, `verification-capacity:`.
- An older Roundfix binary refuses a User Config that sets
  `runs.max_verifying` as an unknown key.

## Coverage Map

- User Story 1, User Story 2, Core Feature 2, Core Feature 6, Core Feature 7,
  Success Metric 3, Success Metric 4, Goals 1 and 2 → API Contracts 4 to 6,
  Invariants 1 and 2, `TestEveryDaemonVerificationStepTakesTheMachineCapacityAfterTheRunsOwn`,
  `TestAStopRequestWhileWaitingForTheMachineCapacityStartsNoCommand`,
  `TestTheRunBudgetDoesNotRunWhileAVerificationWaitsForTheMachineCapacity`,
  `TestTheMachineCapacityWaitIsJournaledAndPrinted`,
  `TestImplementGivesTheEngineTheUserConfigMachineCapacity`, the
  configuration guide and the `implement` reference (task_02).
- Core Feature 1, User Story 4, Goal 5, Core Feature 9 → API Contract 1,
  Invariant 6, `TestRunsMaxVerifyingDefaultsToOne`,
  `TestRunsMaxVerifyingAcceptsZeroAndRefusesANegativeValue`,
  `TestRunsMaxVerifyingInProjectConfigIsIgnoredWithAWarning`,
  `TestTheUserConfigTemplateWritesRunsMaxVerifying`,
  `TestMachineVerificationCapacityOffWritesNothing` (task_01).
- Core Feature 4, Core Feature 5, User Story 5, Success Metric 1,
  Success Metric 2, Goal 3 → API Contracts 2 and 3, Invariants 3 and 4,
  `TestMachineVerificationCapacityServesWaitersInTicketOrder`,
  `TestMachineVerificationCapacityNeverHoldsMoreThanItsUnits`,
  `TestMachineVerificationCapacityRemovesOnlyAnOwnerProvenGone`,
  `TestACancelledMachineVerificationWaitLeavesNoEntry`,
  `TestReadMachineVerificationLeavesOutAnOwnerProvenGoneAndWritesNothing`
  (task_01).
- Core Feature 3, Core Feature 8, User Story 3, Success Metric 5,
  Success Metric 6, Goal 4 → API Contracts 7 and 8, Surface Transcript 1,
  Invariant 5, `TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns`,
  `TestRunsListNamesTheMachineCapacityHolderAndWaiters`,
  `TestDeliverStatusNamesTheMachineCapacityHolderAndWaiters`,
  `TestDoctorReportsTheMachineVerificationCapacity`,
  `TestRunsListMachineCapacityTranscript`, the three command guides and the
  `runs`, `deliver` and `setup` references (task_03).
- Success Metric 7 → API Contracts 9 and 10, the documentation, skill and
  glossary contracts (tasks 02 to 04).

Deviations from the PRD: None.

## Absence Evidence

1. The key — command: `grep -rn 'max_verifying' internal docs .agents skills`;
   result: no match.
2. The capacity store and its directory — command:
   `ls internal/store/machine_verification.go; grep -rn 'verification-capacity\|MachineVerification' internal`;
   result: no such file, no match.
3. A machine-wide Verification gate in the Daemon — command:
   `grep -n 'MachineVerification\|machine_verification' internal/daemon/*.go internal/cli/*.go`;
   result: no match.
4. The status block and doctor line — command:
   `grep -rn 'Machine Verification Capacity\|verification-capacity' internal docs .agents skills`;
   result: no match.
5. The glossary term — command:
   `grep -n 'Machine Verification Capacity' CONTEXT.md`; result: no match.

## ADR Alternatives

- ADR-0384: the decision this Spec records; it rejects `runs.max_active: 1`,
  a Run Database table, one advisory lock per slot, reclaiming by timeout or
  heartbeat, an exclusive retry that takes the whole machine, gating the QA
  Agent's suites and a core budget alone, and the design takes none of them.
- ADR-0056: amended; it chose a per-Task-cycle capacity over a machine-wide
  service, and the design keeps that capacity untouched and adds the machine
  capacity outside it, so the Run's own order of attempts and its exclusive
  retry stay as decided.
- ADR-0373: applies unchanged; it rejected reporting a reused owner as
  residue, walking the new holder's tree, treating an unreadable identity as
  gone and releasing from Doctor, and the design takes none of them: an
  unreadable identity keeps its entry, and the readers remove nothing.
- ADR-0044: applies unchanged; it rejected silent reclamation without proof
  of death, and the design reclaims an entry only on proof.
- ADR-0304: applies unchanged; it rejected a distinct exit code, retrying
  the child on every poll, a `waiting` stage or column, a Project Config key
  and an unbounded wait, and the design takes none of them except that the
  machine wait has no deadline of its own: it ends when the holder ends, a
  Stop Request or `deliver stop` arrives, or the owner is proven gone, because
  a Verification command has no time limit to bound it.
- ADR-0365: applies unchanged; it rejected a stage or column, deriving the
  turn from the Run Database alone, starting the next item early, a Project
  Config key and a default of `1` for reviews, and the design takes none of
  them: the capacity is a file set, read by the status commands.
- ADR-0342: applies unchanged; it rejected skipping the integrated tree
  check, and the design gates that check like any other.
- ADR-0313: applies unchanged; it rejected signalling a mismatched PID, and
  the design signals nothing.
- ADR-0025: not applicable; it set the Task Capacity, which this Spec does
  not change, so none of its rejected options is taken.
- ADR-0038: not applicable; it bounded Agent repairs of Verification, which
  this Spec does not change, so none of its rejected options is taken.
- ADR-0034: not applicable; it governs shared Worktree Bootstrap, which this
  Spec does not gate, so none of its rejected options is taken.
- ADR-0311: not applicable; it set the Delivery Queue's check wait, which
  this Spec does not change, so none of its rejected options is taken.
- ADR-0182: not applicable; it decided the facts a Task settles on, which
  this Spec does not change, so none of its rejected options is taken.

## Integration Points

The Roundfix Home's `verification-capacity/` directory is shared by every
Roundfix process of the user on the machine. Tests use temporary homes only,
never the real `~/.roundfix`, and no process outside the test is signalled.

## Testing Approach

Every Go test uses temporary homes and repositories, fake runners, fake
verifiers or harmless shell commands; none reaches the network, Docker or the
real `~/.roundfix`. Waits end on the event the assertion reads or the harness
wait budget (`internal/testwait`), never on a short sleep.

- `internal/store`, new `machine_verification_test.go`:
  `TestMachineVerificationCapacityServesWaitersInTicketOrder` (capacity 1:
  A holds; B, then C, wait, each after the previous one's `onWait` ran;
  releasing A grants B and not C; releasing B grants C),
  `TestMachineVerificationCapacityNeverHoldsMoreThanItsUnits` (capacity 2,
  eight acquisitions in goroutines; the maximum observed number of holders is
  2 and every acquisition is granted), `TestMachineVerificationCapacityOffWritesNothing`
  (capacity 0: no directory), `TestACancelledMachineVerificationWaitLeavesNoEntry`,
  and `TestReadMachineVerificationLeavesOutAnOwnerProvenGoneAndWritesNothing`
  (the directory's names, sizes and modification times are equal before and
  after).
- `internal/store`, new `machine_verification_unix_test.go` (`//go:build
  unix`): `TestMachineVerificationCapacityRemovesOnlyAnOwnerProvenGone`, a
  table over a holding entry written by the test for the parent PID with
  identities `reused` (`<GOOS>:0.0`, removed, the acquisition is granted at
  once), `absent` (a reaped child's PID from `reapedOwnerPID`, removed),
  `matching` (the parent's real identity, kept), `foreign` (`foreign:0.0`,
  kept) and `empty` (kept); a kept entry makes the acquisition wait, observed
  by `onWait`, and the test then cancels it.
- `internal/config`, new `machine_verification_config_test.go`: the four
  tests of the Coverage Map, through `Load` with temporary homes and
  `DefaultConfigYAML`/`DefaultProjectConfigYAML`.
- `internal/daemon`, new `machine_verification_test.go`, with a fake
  `MachineVerificationGate` that records each request and blocks until the
  test grants it: `TestEveryDaemonVerificationStepTakesTheMachineCapacityAfterTheRunsOwn`
  (one row per step of API Contract 4 with its exact step label, the fake
  verifier records no start before the grant, and the Run unit is taken
  before the machine request and released after the machine release),
  `TestAStopRequestWhileWaitingForTheMachineCapacityStartsNoCommand`,
  `TestTheRunBudgetDoesNotRunWhileAVerificationWaitsForTheMachineCapacity`
  (on `taskCycleBudget`: a one-minute budget paused 30 s after its start for
  two minutes is not expired during the pause nor at 2 min 59 s, is expired
  at 3 min, and its `watch` with a fake timer firing during the pause does
  not cancel) and `TestTheMachineCapacityWaitIsJournaledAndPrinted`.
- `internal/cli`, new `machine_verification_test.go` (task_02):
  `TestImplementGivesTheEngineTheUserConfigMachineCapacity` (the production
  gate built for a loaded config writes a holding entry under the temporary
  home with the configured capacity and none at `0`).
- `internal/cli`, new `machine_verification_status_test.go` (task_03), beside
  the new renderer file `machine_verification_status.go`:
  `TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns` (the configured
  Verification is a shell command that copies the capacity directory's entry
  names and states to a file; the copy shows one `holding` `delivery` entry,
  and none remains after a passing and a failing gate),
  `TestRunsListNamesTheMachineCapacityHolderAndWaiters`,
  `TestDeliverStatusNamesTheMachineCapacityHolderAndWaiters`,
  `TestDoctorReportsTheMachineVerificationCapacity` (rows: empty, holder and
  waiter, off, unreadable entry) and `TestRunsListMachineCapacityTranscript`
  (Surface Transcript 1, whole).
- Sabotages, each a one-line mutation the QA gate applies to a temporary
  copy of the audited head, with the named test that must then fail. The
  author ran each on a prototype of tasks 01 to 03 on 2026-11-10 and saw the
  named test fail:
  1. grant when `holding < capacity`, ignoring the waiters ahead →
     `TestMachineVerificationCapacityKeepsALaterWaiterBehindALiveWaiterAhead`;
  2. never remove an entry whose owner is proven gone →
     `TestMachineVerificationCapacityRemovesOnlyAnOwnerProvenGone`, rows
     `reused` and `absent`;
  3. also remove an entry whose identity is empty or foreign →
     `TestMachineVerificationCapacityRemovesOnlyAnOwnerProvenGone`, rows
     `empty` and `foreign`;
  4. `Release` keeps its entry → `TestACancelledMachineVerificationWaitLeavesNoEntry`,
     run alone;
  5. create the directory at capacity `0` →
     `TestMachineVerificationCapacityOffWritesNothing`;
  6. the reader keeps an entry whose owner is proven gone →
     `TestReadMachineVerificationLeavesOutAnOwnerProvenGoneAndWritesNothing`;
  7. the Daemon skips the machine acquisition →
     `TestEveryDaemonVerificationStepTakesTheMachineCapacityAfterTheRunsOwn`;
  8. `expiredAt` ignores the pause →
     `TestTheRunBudgetDoesNotRunWhileAVerificationWaitsForTheMachineCapacity`;
  9. the resume does not move the deadline →
     `TestTheRunBudgetDoesNotRunWhileAVerificationWaitsForTheMachineCapacity`;
  10. the repository gate passes capacity `0` to the acquisition →
      `TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns`;
  11. the repository gate never releases its lease →
      `TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns`.
- Every existing test that pins the whole standard output of `roundfix
  doctor` with the default doctor dependencies on the default branch when
  task_03 runs gains the `verification-capacity:` line; every other existing
  test of `runs list`, `deliver status`, `implement` and the Daemon passes
  unchanged.

## Build Order

1. task_01: the key, the template line, the entry store with its reader, and
   their tests; the Behavior Surface Record re-recorded last.
2. task_02: the Daemon acquisition, step labels, Run Budget pause, wait event,
   the Implement wiring, the configuration guide section and the `implement`
   reference; depends on step 1.
3. task_03: the Delivery Queue gate, the three status surfaces, their command
   guides and the `runs`, `deliver` and `setup` references; depends on step 1,
   and runs after step 2 because both raise the Roundfix skill version and
   re-record the Behavior Surface Record.
4. task_04: the glossary entries and ADR-0056's notice; depends on step 3, so
   the vocabulary lands once every guide that uses it exists.
5. task_05: the QA gate.

No other Spec is required. Spec 0373 delivers first and changes other
sections of `CONTEXT.md`, the Roundfix `SKILL.md` pair, the owned-skill
version record and the Behavior Surface Record; each Task regenerates the
records from the tree it runs on. The core budget Spec that follows reads
`runs.max_verifying` and the `MachineVerificationGate` seam; it is not
required here.

## Risks & Considerations

- A live holder whose command hangs keeps its unit; the surfaces name it and
  `roundfix stop` ends it. No age rule frees it, by decision.
- A live process that stops polling while waiting (suspended by a signal)
  keeps its place in line; the surfaces name it as `waiting 1`.
- An unreadable entry is not counted, so a corrupted file cannot stop the
  machine; it can let one more Verification run, and `doctor` warns.
- Processes with different `runs.max_verifying` values compare against their
  own value; a change to User Config applies to processes started after it.
- The QA Agent's suites and any Agent's focused tests still run alongside;
  the core budget Spec lightens them.
- A queue whose gate waits occupies its owner; later items wait behind it, as
  with any running item.

## Glossary

- adds: **Machine Verification Capacity**
- changes: **Verification Capacity**
- changes: **Waiting for Verification**
- changes: **Run Budget**

## Decisions

Autonomous route: no person present; decisions taken from `_prd.md` → Decisions, the operator's cycle BB batch 2 brief for Spec 0374 and the characterization above.

- The machine unit is taken after the Run's own unit and released before it,
  so no process waits for its own Run while holding the machine.
- The exclusive retry takes one machine unit: taking every unit would make it
  wait for every other repository; with the default it is alone anyway.
- Entry files under one advisory lock, the lock helpers the Run Database
  writer already uses, so there is no schema change and the Windows lock comes
  for free.
- Ticket order with a 2-second cross-process poll and an in-process wake, so
  attempts of one Run never pay the poll.
- An unreadable entry is skipped, never removed and never counted: removing
  what cannot be proven violates ADR-0373, and counting it could stop every
  repository on one corrupted file.
- The seam for the core budget Spec: `runs.max_verifying` and
  `MachineVerificationGate` are the inputs it divides the cores by; this Spec
  exports no environment variable.
