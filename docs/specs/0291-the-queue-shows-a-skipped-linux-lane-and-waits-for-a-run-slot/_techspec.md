---
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
prd: _prd.md
created: 2026-10-20
---

# The queue shows a skipped Linux Lane and waits for a Run slot — Technical Spec

## Executive Summary

Four code slices with their documentation, and no new package, command, flag,
Delivery Stage, Run Database column or schema version.

1. The Linux Lane tells an expected skip from a Lane Environment Skip, the
   Delivery Queue gate returns the environment skip's reason with its passing
   result, and the Delivery Engine records it on the item as the warning
   `linux-lane-skipped: <reason>` and adds a paragraph to the Pull Request
   body. `roundfix deliver status` already prints item warnings.
2. `roundfix doctor` gains a `linux-lane:` line, printed only while the
   checkout's effective `delivery.linux_lane` is `true`, that checks Docker and
   the lane's image with `docker info` and `docker image inspect` and nothing
   else.
3. `runs.slot_wait` joins User Config, and the Delivery Engine turns a
   ceiling-only refusal of the Implement step into a Run Slot Wait: the running
   item carries the wait as its blocker text, the engine polls the Run slots,
   starts the Run when one frees and parks at the earlier deadline exactly as
   today.
4. The queue's Implement step recognizes the ceiling-only refusal of `roundfix
   implement`, reads the Run slots the way `implement` counts them, and the
   Queue Owner's engine receives the configured wait.

The trade-offs accepted:

- The lane skip rides on the item's free-text warning and the wait on its
  blocker text, parsed back by marker. A schema change would make older item
  binaries refuse the shared Run Database (ADR-0303, ADR-0304).
- The queue recognizes the ceiling refusal from `implement`'s public preflight
  text, because every preflight refusal shares exit code `2`.
- A waiting item occupies its queue for up to `runs.slot_wait`; later items of
  the same queue wait behind it, as behind any running item.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record field
  is minted; the warning and blocker texts and `runs.slot_wait` are the only
  new values. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  recipient is added; the Pull Request body goes through the existing GitHub
  boundary, and tests use fake Docker, fake runners and temporary homes.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 and 2 follow ADR-0303
  ("A Linux Lane skipped for its environment is shown where the delivery is
  read"), ADR-0288 ("The Delivery Queue gate runs the changed Go packages in a
  Linux Lane"), ADR-0294 ("The Linux Lane runs like CI") and ADR-0297 ("The
  Linux Lane reads the repository history read-only"); Build Order 3 and 4
  follow ADR-0304 ("A queue item waits for a Run slot up to runs.slot_wait"),
  ADR-0199 ("A queue's token ceiling stops new starts and never a running
  Task") and ADR-0279 ("Status reports a Queue Owner proven gone, and a
  dismissal survives a fix that leaves its anchor"). Every step follows
  ADR-0184, ADR-0259, ADR-0271, ADR-0285, ADR-0189, ADR-0256, ADR-0280,
  ADR-0104, ADR-0240, ADR-0276, ADR-0265 and ADR-0183, each named with its
  title in `_prd.md` → Project Constraints. Set no obligation here, because
  this Spec changes none of their rules: ADR-0193, ADR-0167, ADR-0229, ADR-0233, ADR-0237, ADR-0250, ADR-0262, ADR-0281, ADR-0286, ADR-0295, ADR-0298, ADR-0197, ADR-0269, ADR-0306, ADR-0272, ADR-0282, ADR-0278, ADR-0290, ADR-0299, ADR-0092, ADR-0120, ADR-0247, ADR-0248, ADR-0254, ADR-0258, ADR-0260, ADR-0261, ADR-0277, ADR-0165, ADR-0251, ADR-0169, ADR-0196, each named in `_prd.md` →
  Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AM approval of 2026-10-20 (the same rules as cycle AL), carrying the cycle
  AJ grant "Mexer em Governed Paths declarados", and the 2026-09-30 skills
  grant ("considere autorizado a ajustar todas as skills se necessário").
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/runtime.md`,
  `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-20 at `70c78d15` (v0.100.0). Line numbers are those
of that commit.

- The Linux Lane (`internal/cli/deliver_linux_lane.go`).
  - `runLinuxLane` (`:41`) handles every skip through one closure (`:56`):
    it sets `result.Skip` and `result.Passed`, writes `roundfix: notice:
    <label>Linux Lane skipped: <reason>` to the console and the lane log, and
    returns. The reasons, in order: `no go.mod at the repository root` (`:72`),
    `go.mod declares no go version` (`:89`), `no changed Go package` (`:119`),
    the three object store reasons of `linuxLaneHistoryMount`, `go is not on
    PATH`, `docker is not on PATH`, `docker info failed: <line>`, `image
    <image> is not present locally` (`:158`) and `docker could not start the
    container: <line>` (`:230`, exit `125`). `linuxLaneResult` has no field
    that tells them apart.
  - `recordLinuxLaneSkip` (`:410`) writes the same notice when
    `linuxLaneOptIn` cannot read the fetched default.
  - `TestLinuxLaneSkipsWithANamedReason` and
    `TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead`
    (`internal/cli/deliver_linux_lane_test.go`) pin each notice and a passing
    gate result, and nothing else.
- The gate and publication (`internal/cli/deliver_workflow.go`).
  - `Gate` (`:1456`) returns `delivery.GateResult{Passed: true}` after every
    skip; `GateResult` (`internal/delivery/engine.go:101`) holds `Passed`,
    `Reason` and `Blocker` only.
  - `Engine.gateCandidate` (`engine.go:1065`) moves a passing item to
    `publishing` and records nothing else; `setStage` clears `Blocker` and
    keeps `Warning`.
  - `Publication` (`:1550`) returns the body `Delivers Spec <slug>.` (`:1572`),
    which `Engine.createPullRequest` passes to `FindOrCreatePullRequest`
    unchanged.
  - `roundfix deliver status` (`internal/cli/deliver.go:375`) prints
    `<slug>\t<stage>\t<blocker or ->\t<worktree or ->` per item (`:406`), then
    `Warning: <slug> <warning>` for each item with a warning (`:410`),
    whatever its stage.
- Doctor (`internal/cli/doctor.go`). `runDoctorCommand` (`:56`) prints node,
  acpx, adapter, profiles, recommendations, pre-pr-review, the machine
  readiness lines, skills, residue, storage and codex; no check reads
  `delivery.linux_lane` or calls Docker. Readiness findings render through
  `readinessResult` (`internal/cli/readiness.go:16`) as `<code>: <text>` with
  `next:` actions.
- The ceiling.
  - `store.ActiveImplementRunCeilingError.Error()`
    (`internal/store/store.go:281`) renders `Active Implement Run ceiling
    reached: <n> Run(s) hold <max> slot(s)`, one `- run_id=…` line per holder
    and the `runs.max_active` hint. `implement` prints it with
    `printPreflightFailure` (`internal/cli/implement.go:404`) and exit `2`
    after reclaiming holders whose owner process is gone
    (`reclaimOrphanedActiveRun`, `internal/cli/cli.go:3517`); `CreateRun`
    counts again inside its write transaction (`store.go:610-620`), which
    `TestRunCeilingIsCountedInsideTheCreateTransaction` proves.
  - The queue runs `roundfix implement --spec <slug>` as a child
    (`RunSpec`, `deliver_workflow.go:918`). `deliveryRunResult` (`:1065`)
    returns `command.failure("roundfix implement")` (`:1100`) for every exit
    that is neither `0` nor a Run that ended unresolved or over budget.
    `runCandidate` (`engine.go:957`) wraps it as `run Implement executor: …`,
    and `Engine.Run` parks the item `delivery-error: <error>` (`:741`).
  - An overlay probe on 2026-10-20 that wrote nothing to the repository
    (`git status --porcelain` empty before and after) fed
    `deliveryRunResult` the stderr of `printPreflightFailure("implement",
    store.ActiveImplementRunCeilingError{…})` with exit `2` and an unchanged
    latest Run. It returned the error `roundfix implement failed with exit code
    2: Preflight failed` followed by the `Reason:` block, `  Active Implement
    Run ceiling reached: 1 Run(s) hold 3 slot(s)`, the holder line, the
    `runs.max_active` hint, the `No side effects:` block and the `Usage:`
    block, so the park reason today is `delivery-error: run Implement executor:
    roundfix implement failed with exit code 2: Preflight failed` plus that
    whole report.
- Configuration (`internal/config/config.go`). `Runs` holds `MaxActive` only
  (`:239`); `runsOverlay` reads `max_active` (`:462`), applied only outside
  Project Config (`:1958`); a Project Config `runs.max_active` is removed with
  `config: runs.max_active in Project Config is ignored; set runs.max_active in
  User Config` (`:1623`, `:733`).

## System Architecture

- `internal/cli`: `deliver_linux_lane.go` marks expected skips and exports one
  image resolver to doctor; `deliver_workflow.go` returns the lane skip from
  `Gate`, recognizes the ceiling refusal in `deliveryRunResult`, answers
  `RunSlots`, and wires the wait into `newCommandDeliveryEngine`; `doctor.go`
  and `health.go` add the `linux-lane` check.
- `internal/delivery`: `engine.go` records the lane warning, extends the Pull
  Request body, and runs the Run Slot Wait in `runCandidate`.
- `internal/config`: `config.go` reads and validates `runs.slot_wait`.
- `internal/store`: `store.go` exports the ceiling message's fixed prefix.

## Implementation Design

### Interfaces

```go
// internal/delivery
const WarningLinuxLaneSkipped = "linux-lane-skipped"
const WaitRunSlot = "run-slot-wait"

type GateResult struct {
	Passed   bool
	Reason   string
	Blocker  string
	LaneSkip string // a Lane Environment Skip's reason; "" when the lane ran, was off or had nothing to test
}

type RunSlotHolder struct{ RunID, Repository, SpecSlug string }
type RunSlots struct {
	MaxActive int
	Holders   []RunSlotHolder
}
func (slots RunSlots) Free() bool // MaxActive == 0 || len(Holders) < MaxActive

// RunSlotRefusal is the Implement step's error when the only refusal is the ceiling.
type RunSlotRefusal struct {
	Slots RunSlots
	Err   error // the failure the queue parks with today
}
func (refusal *RunSlotRefusal) Error() string // Err.Error()
func (refusal *RunSlotRefusal) Unwrap() error

type RunSlotReader interface {
	RunSlots(ctx context.Context) (RunSlots, error)
}

func RunSlotWaitText(deadline time.Time, slots RunSlots) string
func LaneSkipBody(reason string) string

// EngineDependencies gains:
//	RunSlots            RunSlotReader
//	RunSlotWait         time.Duration
//	RunSlotPollInterval time.Duration // 30s when zero

// internal/store
const ActiveImplementRunCeilingReason = "Active Implement Run ceiling reached"

// internal/config
type Runs struct {
	MaxActive int
	SlotWait  time.Duration
}
```

### Data Models

No file, table, column, stage or record changes shape. The item's existing
`Warning` gains the segment `linux-lane-skipped: <reason>`, always last; the
item's existing `Blocker` holds the Run Slot Wait text while its stage is
`running`. User Config gains `runs.slot_wait`.

### API Contracts

1. API Contract: expected skips and Lane Environment Skips.
   `linuxLaneResult` gains `SkipExpected bool`, set only for `no go.mod at the
   repository root` and `no changed Go package`. Every other reason
   `runLinuxLane` skips with is a Lane Environment Skip. The console notice,
   the lane log, the reasons, their order and `roundfix deliver lane`'s exit
   `3` on any skip do not change.
2. API Contract: the gate returns the environment skip. `Gate` returns
   `GateResult{Passed: true, LaneSkip: <reason>}` when the lane is enabled and
   skipped with `SkipExpected` false, and `GateResult{Passed: true, LaneSkip:
   <reason>}` when `linuxLaneOptIn` cannot read the fetched default. A
   disabled lane, an expected skip, a lane that ran and a failed repository
   gate leave `LaneSkip` empty. A lane failure keeps its `Reason` and
   `Blocker`.
3. API Contract: the engine records it. On a passing gate, `gateCandidate`
   removes an earlier `linux-lane-skipped: ` segment from the item's warning
   (from that marker to the end, with the `; ` before it), then appends
   `linux-lane-skipped: <LaneSkip>`, joined by `; ` to any remaining warning,
   when `LaneSkip` is not empty, and persists the item with the `publishing`
   stage. `publishCandidate` appends `"\n\n" + LaneSkipBody(<reason>)` to the
   planned body when the item's warning carries the segment, where
   `LaneSkipBody` returns `Linux Lane: not run. The Delivery Queue gate
   skipped it: <reason>. This candidate's Go packages were not tested on Linux
   before this pull request.` The engine writes no other line; a failed gate
   leaves the warning unchanged.
4. API Contract: doctor's `linux-lane` check.
   - `doctorDependencies` gains `linuxLane func(context.Context,
     roundconfig.Loaded) CheckResult`, set by `defaultDoctorDependencies`;
     `runDoctorCommand` prints its result after `storage` and before `codex`
     only when the dependency is set and `loaded.Config.Delivery.LinuxLane`
     is `true`. `HealthCheckLinuxLane = "linux-lane"`.
   - The image is `delivery.linux_lane_image` when set, else `golang:<go
     directive>` of `go.mod` at the checkout's `HEAD`, through the one image
     resolver the lane also uses. Without a Git root it reports `skipped
     (requires a Git repository)`; without `go.mod` at `HEAD`, `skipped (no
     go.mod at the repository root)`.
   - It resolves `docker` on the command environment's PATH as the lane does
     and runs only `docker info` and then `docker image inspect <image>`, each
     bounded by ten seconds, with the command environment. It never pulls,
     runs or removes anything.
   - Findings, rendered by `readinessResult`, all `warn`:
     `DR-LANE-DOCKER-MISSING: docker is not on PATH` (next `install Docker, or
     set delivery.linux_lane: false`); `DR-LANE-DOCKER-UNREACHABLE: docker info
     failed: <first non-empty output line>` (next `start Docker, then re-run
     roundfix doctor`), after which the image is not inspected;
     `DR-LANE-IMAGE-MISSING: image <image> is not present locally` (next
     `docker pull <image>`); `DR-LANE-IMAGE-UNKNOWN: go.mod declares no go
     version` (next `set delivery.linux_lane_image`). Without findings: `ok
     (docker reachable; image <image> present)`.
5. API Contract: `runs.slot_wait`. A Go duration in User Config, default
   `2h`, `0` allowed. A negative value fails validation with `runs.slot_wait
   must be greater than or equal to 0`; an unparsable value fails parsing. A
   Project Config value is removed before parsing with `config: runs.slot_wait
   in Project Config is ignored; set runs.slot_wait in User Config` and never
   applied. The init templates do not write it.
6. API Contract: the Run Slot Wait (`runCandidate`).
   - The wait applies only when `RunSpec`'s error is a `*RunSlotRefusal`
     (`errors.As`), `RunSlots` is set and `RunSlotWait` is greater than zero.
     Otherwise the error is returned as today.
   - At the first refusal the wait deadline is the clock's now plus
     `RunSlotWait`; when the queue's deadline is set and earlier, the wait
     deadline is the queue's and the queue deadline wins.
   - While now is before the deadline: the item's `Blocker` is set to
     `RunSlotWaitText(deadline, slots)` and persisted with the stage still
     `running` whenever the text changes; the engine logs `roundfix: Delivery
     Queue item <slug>: waiting for a Run slot until <deadline>; <n> Run(s)
     hold <max> slot(s)` once per wait; it sleeps for the poll interval or the
     time left, whichever is shorter, then reads `RunSlots`. When the slots
     are free and the deadline has not passed, it clears `Blocker`, persists
     the item, logs `roundfix: Delivery Queue item <slug>: a Run slot is free;
     starting the Implement Run` and calls `RunSpec` again. A new refusal
     continues the same wait toward the same deadline.
   - At the deadline: if the queue deadline won, the item parks
     `queue-deadline`; otherwise `runCandidate` returns `run Implement
     executor: <the latest refusal's Err>`, which `Engine.Run` parks as
     `delivery-error: …`, byte-identical to today's reason.
   - A sleep that returns an error clears `Blocker` through a context without
     cancellation and returns the error; a `RunSlots` error returns `wait for
     a Run slot: <error>`, which parks as `delivery-error`.
   - `RunSlotWaitText` renders `run-slot-wait: until <deadline in UTC
     RFC3339>; <n> Run(s) hold <max> slot(s)` followed, when there are
     holders, by `: ` and `<run id> in <repository> for <spec>` per holder
     joined by `, `.
7. API Contract: the queue recognizes the ceiling refusal.
   `deliveryRunResult` returns `&delivery.RunSlotRefusal{Err:
   command.failure("roundfix implement")}` when the exit code is `2`, no new
   Implement Run of the Spec exists after the command, and in its standard
   error, with ANSI color sequences removed, the first non-empty line after
   the line `Reason:` starts with `store.ActiveImplementRunCeilingReason` and
   a colon. Every other outcome keeps today's result. `RunSpec` fills the
   refusal's `Slots` from `RunSlots`; when that read fails it returns the
   plain `Err`, which parks as today. `ActiveImplementRunCeilingError.Error()`
   builds its first line from the same constant, so its text is unchanged.
8. API Contract: the Run slots. `commandDeliveryWorkflow.RunSlots` returns
   `MaxActive` from the Queue Owner's `runs.max_active` and, as `Holders`,
   each Active Implement Run from the Run Database except one whose recorded
   owner PID names a process that is not alive, the holders `implement` would
   reclaim, sorted by Run ID, with the Run's Git root as `Repository`. It
   reads only. `newCommandDeliveryEngine` passes the workflow as `RunSlots`
   and `runs.slot_wait` as `RunSlotWait`.
9. API Contract: everything else stays. Command names, flags, exit codes, the
   lane's notices, skip reasons and outcomes, every park reason other than the
   two above, Delivery Retry, `roundfix implement`'s output and exit codes,
   the doctor lines of a checkout whose lane is off, the store's ceiling
   check and every existing test name keep their behavior.

### Surface Transcripts

1. Surface Transcript: status shows a Lane Environment Skip. A Delivery Queue
   in a temporary repository holds one merged item `0285-example`, no
   worktree recorded, whose warning the engine recorded after the gate
   returned the reason `image golang:1.26 is not present locally`; the queue
   has no owner, limits or Runs.

   ```transcript
   $ roundfix deliver status
   stdout:
   0285-example	merged	-	-
   Warning: 0285-example linux-lane-skipped: image golang:1.26 is not present locally
   Limits: deadline none, retries per item none, concurrency 1, tokens none
   Usage: no Runs recorded
   stderr:
   exit: 0
   ```

2. Surface Transcript: status shows a Run Slot Wait. The queue holds one item
   `0291-example` in stage `running` with worktree `<worktree>`, whose blocker
   is `RunSlotWaitText` of deadline `2026-10-21T00:00:00Z` and one holder,
   Run `run_holder` in `<holder repository>` for `0001-holder`, at a ceiling
   of 1; the queue has no owner, limits or Runs of its own.

   ```transcript
   $ roundfix deliver status
   stdout:
   0291-example	running	run-slot-wait: until 2026-10-21T00:00:00Z; 1 Run(s) hold 1 slot(s): run_holder in <holder repository> for 0001-holder	<worktree>
   Limits: deadline none, retries per item none, concurrency 1, tokens none
   Usage: no Runs recorded
   stderr:
   exit: 0
   ```

### Invariants

```text
1. A Linux Lane skip passes the gate; only a Lane Environment Skip of an enabled or unreadable lane reaches the item, the Pull Request body and doctor.
2. The lane's console and log notices and roundfix deliver lane's exit codes are unchanged.
3. Doctor runs no Docker command that pulls, runs or removes, and prints its linux-lane line only while delivery.linux_lane is true.
4. The queue waits only on a refusal whose sole reason is the Active Implement Run ceiling, never past runs.slot_wait or the queue deadline.
5. The store's CreateRun transaction is the only authority on the ceiling; waiting queues never hold more Active Implement Runs than runs.max_active allows.
6. At the end of a wait the item parks with today's delivery-error reason, or queue-deadline when the queue deadline came first.
7. No Delivery Stage, Run Database column or schema version is added.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1, Invariants 1-2 → `linuxLaneResult.SkipExpected`, `Gate`, `gateCandidate`, `publishCandidate` (API Contracts 1-3), Surface Transcript 1
- Goal 2, Success Metric 2 → API Contracts 1-3
- Goal 1, Core Feature 2, Success Metric 3, Invariant 3 → doctor's `linux-lane` check (API Contract 4)
- Goal 3, Core Feature 3, Success Metric 4, Invariants 4, 6-7 → `runs.slot_wait`, `runCandidate` (API Contracts 5-6), Surface Transcript 2
- Goal 4, Success Metric 5, Invariant 5 → the store's `CreateRun` transaction and API Contracts 6-8
- Core Feature 4 → `deliveryRunResult`, `RunSlots`, `newCommandDeliveryEngine` (API Contracts 7-8)
- Core Feature 5 → the `deliver`, `doctor` and configuration guides, the Roundfix `deliver`, `setup` and `runtime` references, **Lane Environment Skip** and **Run Slot Wait** in `CONTEXT.md`
- Success Metric 6 → the PRD's Unreachable Acceptance
- Goal 2, Success Metric 2, Success Metric 5, API Contract 9 → the existing lane, gate, doctor, implement, store and delivery suites and the QA scope audit

## Absence Evidence

1. The lane skip's kind and its gate field — command: `grep -n 'SkipExpected\|LaneSkip' internal/cli/deliver_linux_lane.go internal/cli/deliver_workflow.go internal/delivery/engine.go`; result: no match at `70c78d15`.
2. The item warning and Pull Request paragraph — command: `grep -rn 'linux-lane-skipped\|LaneSkipBody\|not tested on Linux' internal docs .agents`; result: no match at `70c78d15`.
3. Doctor's Linux Lane check — command: `grep -rn 'linux-lane\|DR-LANE' internal/cli/doctor.go internal/cli/health.go docs/user-guide/commands/doctor.md .agents/skills/roundfix/references/setup.md`; result: no match at `70c78d15`.
4. `runs.slot_wait` — command: `grep -rn 'slot_wait\|SlotWait' internal docs .agents`; result: no match at `70c78d15`.
5. The Run Slot Wait in the engine and its types — command: `grep -rn 'RunSlot\|run-slot-wait\|WaitRunSlot' internal`; result: no match at `70c78d15`.
6. The exported ceiling prefix — command: `grep -n 'ActiveImplementRunCeilingReason' internal/store/store.go`; result: no match at `70c78d15`; the text is a format string inside `Error()`.

## ADR Alternatives

- ADR-0303: the design is the decision; none of its rejected alternatives (parking on a skip, showing every skip, a new column, doctor reading the fetched default) is taken.
- ADR-0304: the design is the decision; none of its rejected alternatives (a new exit code, retrying the child each poll, a new stage, a Project Config key, an unbounded wait) is taken.
- ADR-0288: no rejected alternative is taken; a skip still passes the gate and the lane still covers the changed packages.
- ADR-0294: no rejected alternative is taken; the lane still never pulls and keeps its image rules.
- ADR-0297: no rejected alternative is taken; the object store mount is unchanged, and its skips become Lane Environment Skips.
- ADR-0199: no rejected alternative is taken; the token ceiling still stops only new starts, and the wait never stops a running Task.
- ADR-0279: no rejected alternative is taken; status still reports a Queue Owner proven gone beside a stale wait text.
- ADR-0184: no rejected alternative is taken; both status surfaces are stated as transcripts.
- ADR-0259: no rejected alternative is taken; every new top-level test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; waits use fake clocks and sleepers, and the race test waits on events up to the test deadline.
- ADR-0285: no rejected alternative is taken; the new tests compile for Windows, a unix-only case in its own constrained file.
- ADR-0189: no rejected alternative is taken; the Roundfix skill edits raise its recorded version.
- ADR-0256: no rejected alternative is taken; the covering `deliver` and `setup` references change with the surfaces.
- ADR-0280: no rejected alternative is taken; each guide Task records its own fingerprint.
- ADR-0104: no rejected alternative is taken; the outside-evidence rows rest on GitHub's documentation and an existing store test.
- ADR-0240: no rejected alternative is taken; the real-Docker doctor run is one Unreachable Acceptance entry.
- ADR-0276: no rejected alternative is taken; network-free rows qualify, and byte invariants compare with the Delivery Base.
- ADR-0265: no rejected alternative is taken; every new component has its absence proof.
- ADR-0183: no rejected alternative is taken; each attributed quote carries its receipt.
- ADR-0193: no rejected alternative is taken; prerequisite waits are unchanged.
- ADR-0167: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0229: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0233: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0237: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0250: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0281: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0286: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0295: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0298: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0197: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0269: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0306: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0272: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0282: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0278: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0290: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0299: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0092: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0120: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0247: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0248: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0254: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0258: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0260: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0261: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0277: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0165: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0251: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0169: no rejected alternative is taken; this Spec changes none of its rules.
- ADR-0196: no rejected alternative is taken; this Spec changes none of its rules.

## Integration Points

The Run Database, read for the queue and the Active Implement Runs and written
only through the item updates the engine already makes; the child `roundfix
implement`'s exit code and standard error; Docker, through `docker info` and
`docker image inspect` in `doctor`; the GitHub boundary, through the Pull
Request body only. No test, Verification command or QA row runs real Docker,
reaches the network, or reads the real `~/.roundfix`.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary homes and repositories, ends every child it starts in
`t.Cleanup`, and declares no real-time budget under ten seconds.

- `internal/cli/deliver_lane_skip_visible_test.go`, new:
  - `TestGateReturnsTheLaneEnvironmentSkip`, through `Gate` with the compiled
    Docker stand-in of `deliver_linux_lane_test.go`: subtests `no go`, `no
    docker`, `daemon`, `image`, `start refused`, `no version` and `unreadable
    default` return their reason in `LaneSkip`; `no module`, `no package`, a
    disabled lane, a lane that passed and a failed repository gate return it
    empty; every console notice keeps its exact text;
  - `TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript`, asserting
    Surface Transcript 1 whole.
- `internal/delivery/lane_skip_warning_test.go`, new, with the engine fakes:
  - `TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody`: an item
    with a `premise-changed` warning gets `; linux-lane-skipped: <reason>`
    and the fake boundary receives `Delivers Spec <slug>.` followed by the
    `LaneSkipBody` paragraph;
  - `TestAQuietGateLeavesNoLaneWarningAndReplacesAnEarlierOne`: an empty
    `LaneSkip` leaves the body and warning unchanged, and a retried gate with
    an empty or a different reason removes or replaces the earlier segment.
- `internal/cli/doctor_linux_lane_test.go`, new:
  - `TestDoctorLinuxLaneCheckReportsDockerAndTheImage`, calling the check
    with the Docker stand-in: subtests `ok`, `image missing`, `docker
    missing`, `daemon unreachable` (no image inspection recorded),
    `configured image`, `no go version` and `no module`, each asserting the
    rendered line and that the stand-in recorded only `info` and `image
    inspect <image>`;
  - `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`, through
    `runDoctorCommand` with fake dependencies: with the lane on, the
    `linux-lane:` line sits between `storage:` and `codex:`; with it off, the
    output equals today's.
- `internal/config/run_slot_wait_test.go`, new:
  `TestRunSlotWaitIsReadFromUserConfig` (default `2h`, `90m`, `0`),
  `TestProjectConfigCannotSetTheRunSlotWait` (the exact warning, User Config
  kept), `TestRunSlotWaitRefusesANegativeOrInvalidValue`.
- `internal/delivery/run_slot_wait_test.go`, new, with a fake clock, a fake
  sleeper and a fake `RunSlotReader`:
  - `TestARunCeilingRefusalWaitsAndStartsWhenASlotFrees` (blocker text
    persisted while waiting, cleared before the second `RunSpec`, the two log
    lines);
  - `TestARunSlotWaitParksAsTodayAtItsDeadline` (reason byte-identical to the
    no-wait park of the same refusal);
  - `TestTheQueueDeadlineEndsARunSlotWait` (`queue-deadline`);
  - `TestAZeroRunSlotWaitParksAtOnce` and
    `TestOtherRunSpecFailuresDoNotWait`;
  - `TestAStoppedOwnerClearsTheRunSlotWait` (cancelled sleeper);
  - `TestTwoQueuesWaitingForOneSlotStartOneRun`: two engines over two queues
    in one Run Database at a ceiling of 1, whose fake runners create their
    Run with `store.CreateRun` and `MaxActiveImplementRuns: 1` and return a
    `*RunSlotRefusal` on `store.ActiveImplementRunCeilingError`; after the
    holder completes, both pollers are released together; exactly one Run is
    created, Active Implement Runs never exceed 1, and the other item still
    carries its wait text. It waits on channel events through the
    `internal/testwait` helpers, bounded only by the test deadline.
- `internal/cli/deliver_run_slot_wait_test.go`, new:
  - `TestDeliveryRunResultRecognizesOnlyTheCeilingRefusal`, with the real
    `printPreflightFailure` output of a ceiling error (plain and colored), of
    an Active Run in the working tree, of a ceiling error after a new Run, and
    an exit `1`;
  - `TestRunSpecReturnsARunSlotRefusalWithItsHolders`, through `RunSpec` with
    a compiled item binary that prints a ceiling refusal and exits `2`;
  - `TestDeliverStatusShowsARunSlotWaitTranscript`, asserting Surface
    Transcript 2 whole.
- `internal/cli/deliver_run_slot_wait_unix_test.go`, new, `//go:build unix`:
  `TestRunSlotsExcludeAHolderWhoseOwnerIsGone`, with a Run whose owner PID is
  a `sh` child the test started and waited for.
- Unchanged and run as consumers: the lane, `deliver lane`, doctor, implement
  ceiling, store ceiling and delivery engine suites.

## Build Order

1. Lane Environment Skip on the item: `SkipExpected`, `GateResult.LaneSkip`,
   the warning and body, their tests and the `deliver` guide's Linux Lane
   text (no dependency).
2. Doctor's `linux-lane` check, the shared image resolver, the `doctor`
   guide, the Roundfix `deliver`, `setup` and `runtime` references for both lane
   surfaces, the skill version and **Lane Environment Skip** in `CONTEXT.md`
   (depends on: 1, for the lane file, the references' lane text and the
   Behavior Surface Record).
3. `runs.slot_wait` and the Run Slot Wait in the engine, their tests and the
   configuration guide (depends on: 2, because both change
   `internal/delivery/engine.go` after 1 and the Behavior Surface Record).
4. The queue's refusal recognition, `RunSlots`, the engine wiring, the status
   transcript, the `deliver` guide's wait section, the Roundfix `deliver`
   reference, the skill version and **Run Slot Wait** in `CONTEXT.md`
   (depends on: 3, for the engine types).
5. QA gate (depends on: 1, 2, 3, 4).

## Risks & Considerations

- An older item binary prints the same ceiling text, so recognition holds for
  it; a future change to that text would stop the wait, which
  `TestDeliveryRunResultRecognizesOnlyTheCeilingRefusal` catches because it
  renders the real error.
- A holder Run with no recorded owner PID counts as holding a slot, as
  `implement` counts it.
- A killed owner can leave the wait text on a running item; status then also
  prints the owner proven gone, and `roundfix deliver resume` starts a new
  wait.
- An older Roundfix binary refuses a User Config that sets `runs.slot_wait`;
  the configuration guide says so, and the default needs no setting.

## Glossary

- adds: **Lane Environment Skip**
- adds: **Run Slot Wait**

## Decisions

- Mark expected skips in the lane result instead of matching reason strings in
  the gate, so a new skip reason is an environment skip by default. See
  ADR-0303.
- Carry the skip on the item's warning and the wait on its blocker text,
  because both are durable and already printed by status. See ADR-0303 and
  ADR-0304.
- Recognize the refusal from `implement`'s rendered `Reason:` line through an
  exported prefix that the store's error also uses, so the two cannot drift.
- Poll the Run Database, not the child, and let `CreateRun` decide.
- Keep the wait's deadline from the first refusal across lost races, and
  start a new wait only on resume or retry.
