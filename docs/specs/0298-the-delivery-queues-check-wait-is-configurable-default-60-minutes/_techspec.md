---
spec: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
prd: _prd.md
created: 2026-10-22
---

# The Delivery Queue's check wait is configurable, default 60 minutes — Technical Spec

## Executive Summary

The Delivery Engine bounds every Check Wait with `engine.checkTimeout`, which
`NewEngine` sets from `EngineDependencies.CheckTimeout` or, when that is not
positive, from the constant `defaultCheckTimeout = 30 * time.Minute`; the
`deliver` command never sets the dependency. This Spec adds the User Config
key `delivery.check_timeout` to `internal/config` (a Go duration, default
`60m`, greater than 0 and at most `24h`, ignored with a warning in Project
Config), makes the engine's default the same constant, passes the loaded value
from `internal/cli` through an extracted dependency builder, and makes the
`checks-timeout` park text name the key and the value. The three deadline call
sites keep their shape, so each wait still starts its own deadline.

The trade-off accepted: an adopter that sets nothing now waits up to twice as
long before a `checks-timeout` park, in exchange for not parking while a
normal CI run is still going.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the key is a YAML key and the park text extends an existing
  blocker. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses fake clocks, fake forges
  and temporary homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 1 and 2 follow ADR-0236,
  ADR-0304 and ADR-0311, named with their titles in `_prd.md` → Project
  Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AO approval of 2026-10-22 ("the same rules as cycle AN"), carrying the cycle
  AJ grant "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at `34b4e246`. Line numbers move when
other Specs merge; the function names do not.

1. `internal/delivery/engine.go` declares `defaultCheckTimeout = 30 *
   time.Minute` and `defaultCheckInterval = 15 * time.Second`. `NewEngine`
   uses `dependencies.CheckTimeout` when it is positive and the constant
   otherwise.
2. `newCommandDeliveryEngine` (`internal/cli/deliver_workflow.go`) builds the
   engine for `deliver start`, `deliver resume` (both through
   `runDeliveryOwner`) and `deliver retry` (`runDeliverRetry`), passing
   `RunSlotWait: loaded.Config.Runs.SlotWait` and no `CheckTimeout`. Only
   tests set `CheckTimeout` (2, 3, 5 seconds and 10 minutes).
3. Three deadlines read `engine.clock.Now().Add(engine.checkTimeout)`:
   - `Engine.createPullRequest`: while an existing Pull Request reports an
     earlier candidate commit as its head, the engine re-reads it every check
     interval until the deadline, then fails with `find or create pull
     request: PR Head Branch is at "<sha>", expected "<sha>"`;
   - `Engine.checkCandidate`, when the item enters `checking`;
   - `Engine.checkCandidate` again, once per call (`timeoutRestarted`), when
     it re-runs a failed check (ADR-0236).
   Each is a local computed when its step starts; nothing carries a deadline
   between calls.
4. At the deadline `Engine.checkCandidate` calls `engine.park(ctx, gitRoot,
   item, BlockerChecksTimeout)`: the blocker is the bare `checks-timeout`.
   With the fake clock, an engine without `CheckTimeout` whose checks stay
   pending parked after exactly 30 minutes and 121 polls.
5. `internal/config/config.go`: `deliveryOverlay` knows `derived_paths`,
   `item_binary`, `linux_lane` and `linux_lane_image`; strict decoding refused
   a User Config or Project Config with `delivery:\n  check_timeout: 60m` as
   `yaml: unmarshal errors`. `applyConfigContent` removes `runs.slot_wait`,
   `runs.max_active`, `jev.monthly_ceiling_usd` and `store.run_retention_days`
   from a Project Config with `warnIgnoredProjectSetting`, which writes
   `config: <key> in Project Config is ignored; set <key> in User Config`.
   Every loader, `DeliveryConfigAtCommit` for the fetched default branch
   included, goes through `applyConfigContent` with its source.
   `defaults.profile_proof_setup_timeout` is validated with `must be greater
   than 0 and at most 10m`.

Sweep of every reader of the blocker, at `34b4e246` (`rg -n
'BlockerChecksTimeout|checks-timeout' internal docs .agents skills`):

- `ClassifyPark` (`internal/delivery/park_class.go`) cuts the blocker at the
  first `:` and maps `checks-timeout` to `environment`; status and the Pending
  Question print the whole blocker.
- `store` retry compares the stored blocker with the parked item's own
  blocker, whatever its text.
- Five existing tests assert the bare blocker after a timeout:
  `TestPersistentCheckReadErrorsParkAtTheDeadline` and
  `TestChecksTimeoutParksTheItem` (`internal/delivery/engine_test.go`, 2s),
  `TestABlockedMergeStateParksAtTheChecksTimeout`
  (`internal/delivery/merge_state_test.go`, 3s),
  `TestAnOldCheckAttemptCannotKeepResettingTheTimeout`
  (`internal/delivery/check_rerun_test.go`, 5s) and
  `TestStaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce`
  (`internal/delivery/stale_check_test.go`, 5s).
- `internal/cli/deliver_park_status_test.go` and
  `internal/cli/deliver_owner_status_test.go` set a bare `checks-timeout` by
  hand, the form a record made before this change keeps.
- Text: `docs/user-guide/commands/deliver.md` and the Roundfix `deliver`
  reference say "The existing checks timeout parks the item as
  `checks-timeout`" and "up to the check timeout"; neither gives the value.
  `docs/user-guide/configuration.md` has no entry. An older decision record
  quotes a past `checks-timeout` park and stays as written.

## System Architecture

- `internal/config/config.go`: `DeliveryCheckTimeoutKey`,
  `DefaultDeliveryCheckTimeout`, `MaxDeliveryCheckTimeout`,
  `Delivery.CheckTimeout`, the overlay field, the Project Config removal, the
  built-in value and the validation.
- `internal/delivery/engine.go`: `defaultCheckTimeout` reads
  `config.DefaultDeliveryCheckTimeout`; `checksTimeoutBlocker` builds the park
  text. `internal/delivery` already depends on `internal/config` through
  `internal/store`, so the new import adds no package edge to the graph.
- `internal/cli/deliver_workflow.go`: `commandDeliveryEngineDependencies`,
  extracted from `newCommandDeliveryEngine`, sets `CheckTimeout`.
- `docs/user-guide/configuration.md`, `docs/user-guide/commands/deliver.md`,
  `.agents/skills/roundfix/references/deliver.md` and its mirror: the key and
  the park text. `CONTEXT.md`: **Check Wait**.

## Implementation Design

### Interfaces

```go
// internal/config
const DeliveryCheckTimeoutKey = "delivery.check_timeout"
const DefaultDeliveryCheckTimeout = 60 * time.Minute
const MaxDeliveryCheckTimeout = 24 * time.Hour

type Delivery struct {
	CheckTimeout time.Duration
	// existing fields unchanged
}

// internal/delivery
func checksTimeoutBlocker(timeout time.Duration) string

// internal/cli
func commandDeliveryEngineDependencies(runStore *store.Store, loaded roundconfig.Loaded) delivery.EngineDependencies
```

### Data Models

`config.Delivery` gains `CheckTimeout`. The Delivery Queue item, its stages,
the Run Database schema and `delivery.EngineDependencies` keep their shapes;
the item's blocker text grows a detail after `: `.

### API Contracts

1. API Contract: the key. `delivery.check_timeout` decodes with the existing
   `durationValue` (`0` or a `time.ParseDuration` value). `Builtin()` sets
   `Delivery.CheckTimeout` to `DefaultDeliveryCheckTimeout`, `60m`. A User
   Config value replaces it. `Validate` refuses a value `<= 0` or `>
   MaxDeliveryCheckTimeout` with exactly `delivery.check_timeout must be
   greater than 0 and at most 24h`; an unparsable value fails decoding with
   `invalid duration`. No init template and no Full example changes.
2. API Contract: User Config only. `applyConfigContent` removes
   `delivery.check_timeout` from a Project Config source before decoding and
   calls `warnIgnoredProjectSetting(DeliveryCheckTimeoutKey)`, so the warning
   is `config: delivery.check_timeout in Project Config is ignored; set
   delivery.check_timeout in User Config`, written once per load; the overlay
   applies the value only for a source other than Project Config. This holds
   for the checkout's Project Config and for `DeliveryConfigAtCommit`'s
   default-branch Project Config. An invalid Project Config value is removed
   before decoding, so it is ignored, not refused.
3. API Contract: the engine. `defaultCheckTimeout` equals
   `config.DefaultDeliveryCheckTimeout`; `NewEngine` keeps using
   `dependencies.CheckTimeout` when positive. The three deadlines keep their
   code: each wait computes its deadline from its own start, and
   `Engine.checkCandidate` restarts it at most once per call after a re-run.
4. API Contract: the park text. `checksTimeoutBlocker(timeout)` returns
   `checks-timeout: current-head checks not settled within
   delivery.check_timeout <timeout.String()>`, and `Engine.checkCandidate`
   parks with it at the deadline. `BlockerChecksTimeout` keeps its value,
   `ClassifyPark` keeps mapping it to `environment`, and `Next` keeps its
   default text. A bare `checks-timeout` recorded before this change still
   classifies and retries.
5. API Contract: the wiring. `newCommandDeliveryEngine` returns
   `delivery.NewEngine(runStore, commandDeliveryEngineDependencies(runStore,
   loaded))`; the builder creates the workflow and the GitHub CLI exactly as
   today and sets `CheckTimeout: loaded.Config.Delivery.CheckTimeout` beside
   `RunSlotWait`. It writes no file and starts no process.
6. API Contract: documentation. `docs/user-guide/configuration.md` gains a
   `## Delivery Queue check wait` section and a key-reference row;
   `docs/user-guide/commands/deliver.md` and the Roundfix `deliver` reference
   give the park text and the key; `CONTEXT.md` defines **Check Wait**.

### Surface Transcripts

None. No command's flags, help text, standard output or exit code changes.
The park text reaches an operator through `roundfix deliver status`, which
prints the stored blocker unchanged; the tests of API Contracts 2 and 4 assert
the exact warning and blocker strings.

### Invariants

```text
1. With nothing set, every Check Wait is 60 minutes; a User Config value replaces it; nothing else does.
2. delivery.check_timeout is greater than 0 and at most 24h, or the configuration does not load.
3. A Project Config delivery.check_timeout never changes the wait and always produces the ignored-setting warning.
4. Each Check Wait starts its own deadline; a re-run restarts it at most once per wait.
5. A checks-timeout park names delivery.check_timeout and the duration it waited, and stays in the environment Park Class.
6. The deliver command's engine receives the loaded delivery.check_timeout.
```

## Coverage Map

Deviations from the PRD: None.

- Goals 1-3, Core Features 1-2, Success Metrics 3-4, Invariants 2-3 → the key, its validation and its scope (API Contracts 1-2)
- Goal 2, Core Feature 3, Success Metrics 1-2, Invariants 1 and 6 → the engine default and the wiring (API Contracts 3 and 5)
- Goal 4, Core Feature 3, Success Metric 5, Invariant 4 → the unchanged deadline call sites and their tests (API Contract 3)
- Goal 5, Core Feature 4, Success Metric 1, Invariant 5 → `checksTimeoutBlocker` (API Contract 4)
- Goal 6, Core Feature 5, Success Metric 6 → the guides, the Roundfix reference and **Check Wait** (API Contract 6)
- Success Metric 7 → the PRD's Unreachable Acceptance

## Absence Evidence

1. `DeliveryCheckTimeoutKey`, `DefaultDeliveryCheckTimeout`, `MaxDeliveryCheckTimeout` and the `check_timeout` overlay field — command: `rg -n -e check_timeout -e CheckTimeoutKey -e DefaultDeliveryCheckTimeout -e MaxDeliveryCheckTimeout internal`; result: no match at `34b4e246` (exit 1); strict decoding refuses the key in both scopes.
2. `checksTimeoutBlocker` and `commandDeliveryEngineDependencies` — command: `rg -n -e 'func commandDeliveryEngineDependencies' -e checksTimeoutBlocker internal`; result: no match at `34b4e246` (exit 1); the park uses the bare constant and the dependencies are built inline in `newCommandDeliveryEngine`.
3. A configured check wait in the `deliver` command — command: `rg -n CheckTimeout internal/cli -g '!*_test.go'`; result: no match at `34b4e246` (exit 1); the command always used the engine's constant.

## ADR Alternatives

- ADR-0236: no rejected alternative is taken; a failed check is still judged on a merge with the current default branch, and the timeout still restarts only once per wait after a re-run.
- ADR-0304: no rejected alternative is taken; this Spec adds no stage, column or schema and no Project Config key for a machine's wait, and leaves `runs.slot_wait` unchanged.
- ADR-0311: no rejected alternative is taken; the wait is not a constant, not a Project Config key, not derived from a workflow, and `0` is refused rather than given a meaning.

## Integration Points

None at runtime beyond the existing forge reads: the queue polls the same
current-head checks through the same `gh` calls. No test, Verification command
or QA row reads the real `~/.roundfix`, another repository or a forge.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary homes and the engine's fake clock and sleeper, and waits on no
real time.

- `internal/config/delivery_check_timeout_test.go` (new):
  - `TestDeliveryCheckTimeoutIsReadFromUserConfig`: absent (60 minutes,
    `DefaultDeliveryCheckTimeout`, and the constant equals 60 minutes), `45m`
    (45 minutes) and `24h` (`MaxDeliveryCheckTimeout`).
  - `TestProjectConfigCannotSetTheDeliveryCheckTimeout`: User Config `45m`
    and Project Config `0`, `5m`, `-1h` and `invalid`; the value stays 45
    minutes and standard error is exactly the warning line.
  - `TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds`: `0`, `0s`, `-1m`
    and `24h1s` give the range error; `soon` gives `invalid duration`.
- `internal/delivery/check_wait_test.go` (new):
  - `TestTheCheckWaitDefaultsToSixtyMinutes`: an engine without
    `CheckTimeout` and checks pending forever parks after exactly
    `config.DefaultDeliveryCheckTimeout`, 60 minutes and 241 polls, with the
    blocker `checks-timeout: current-head checks not settled within
    delivery.check_timeout 1h0m0s` in the `environment` Park Class.
  - `TestTheCheckWaitUsesTheConfiguredDuration`: `CheckTimeout` 45 minutes
    parks after exactly 45 minutes with `... 45m0s`.
  - `TestEachCheckWaitStartsItsOwnDeadline`: with 45 minutes, the first call
    sees `pass` at the 161st poll, 40 minutes in, and moves to `merging`; the
    item is set back to `checking` and a second call, pending forever, parks
    exactly 45 minutes after its own start.
  - `TestThePullRequestHeadWaitUsesTheCheckWait`: `createPullRequest` against
    an existing Pull Request that keeps reporting the earlier candidate
    commit fails with `PR Head Branch is at "earlier", expected "candidate"`
    after exactly 45 minutes.
- `internal/cli/deliver_check_timeout_test.go` (new):
  `TestTheDeliveryEngineReceivesTheConfiguredCheckWait`: with
  `roundconfig.Builtin()` the builder's `CheckTimeout` is
  `DefaultDeliveryCheckTimeout`; with `Delivery.CheckTimeout` 45 minutes and
  `Runs.SlotWait` 90 minutes it carries both.
- Declared break: the five tests named under Current behavior assert the new
  park text for their timeout (`... delivery.check_timeout 2s`, `3s`, `5s`)
  instead of the bare `checks-timeout`; their inputs, polls and other
  assertions do not change.
- Unchanged and run by name: `TestCheckRerunRestartsTheTimeout`,
  `TestAConflictingPullRequestStopsTheCheckWaitAtOnce`,
  `TestRunSlotWaitIsReadFromUserConfig`,
  `TestProjectConfigCannotSetTheRunSlotWait` and
  `TestDeliverStatusPrintsAParkLinePerParkedItem`.

## Build Order

1. The configuration guide, the `deliver` guide, the Roundfix `deliver`
   reference, the skill version, the Behavior Surface Record and **Check
   Wait** (no dependency: the text states the contracts of `_techspec.md`,
   and its Verification reads only documents).
2. The key, its validation and scope, the engine default and park text, the
   wiring, and their tests (depends on: 1, which writes the guides for the
   surface this step changes, so the code and its text land in one chain).
3. QA gate (depends on: 1, 2).

No Spec is listed in `requires`. Spec 0297 delivers first and changes the
collision check; Spec 0299 delivers after this one and also changes
`internal/delivery`, for resuming an archived item. This Spec changes only the
default, the park text and the three Check Wait readers' source of the value,
and it rebases on 0297.

## Risks & Considerations

- An adopter that sets nothing waits up to 60 minutes per Check Wait, and a
  check re-run can extend one wait to two hours; a stuck check holds its queue
  that long before parking.
- A Delivery Queue that runs an item binary built from a branch older than
  this change refuses a User Config that sets `delivery.check_timeout`, as an
  older binary refuses `runs.slot_wait`; the default needs no setting.
- The park text changes for every `checks-timeout` park; a script that
  compares the blocker for equality must compare its prefix before `:`.
- A Queue Owner already running keeps the value it read; the operator stops
  and resumes it, or retries a parked item, to apply a new value.

## Glossary

- adds: **Check Wait**

## Decisions

- The wait is the User Config key `delivery.check_timeout`, 60 minutes by
  default, greater than 0 and at most 24 hours; see ADR-0311.
- Project Config cannot set it, as for `runs.slot_wait`; see ADR-0304.
- The single restart after a re-run is kept; see ADR-0236.
- The dependency builder is extracted from `newCommandDeliveryEngine` so the
  wiring has one testable place; it is production code that the command
  calls, not a test hook.
