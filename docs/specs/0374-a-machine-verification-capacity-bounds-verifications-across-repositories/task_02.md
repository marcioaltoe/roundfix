---
task: task_02
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
status: completed
type: backend
complexity: high
---

# Task 02: Every Daemon Verification takes the Machine Verification Capacity after the Run's own

## Overview

The Daemon gains a `MachineVerificationGate` dependency. Every Verification
it starts — pre-work probe, repository precondition, Task attempt with the
settlement command, exclusive retry, Integrated Repository Verification and
the QA Task's repository Verification — takes the Run's own unit first and a
machine unit second, starts no command before both are held, and gives the
machine unit back first. While it waits the Run Budget pauses, a Stop Request
still ends the wait cleanly, and the wait is journaled and printed. The
Implement Command wires the gate to the store with `runs.max_verifying`. The
configuration guide and the Roundfix `implement` reference describe the key
and the gated steps.

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `acquireVerificationCapacity` is the only function that takes Verification
   Capacity, reached by `verifyTaskPreWork` and `runTaskVerificationRequest`;
   name each caller of `runTaskVerificationRequest`.
2. MUST add `MachineVerificationGate` and `Dependencies.MachineVerification`
   to `internal/daemon/engine.go` with the signature of `_techspec.md` →
   Interfaces, and `Step string` to `verificationAttemptRequest`, set by the
   four constructions with the exact labels of API Contract 4
   (`Task <id> pre-work probe`, `Task <id> repository precondition`,
   `Integrated Repository Verification`, `QA Task <id> repository Verification`);
   a request with an empty step is labelled
   `Task <id> Verification <identity>` from its `identity()`.
3. MUST change `acquireVerificationCapacity` in
   `internal/daemon/task_engine.go` as API Contract 4 says: keep the Run's
   capacity path as it is; then, when the dependency is not nil, acquire a
   machine unit with holder kind `run`, the Run ID, the Spec slug and the
   step, under the same Stop Request watch; on a Stop Request release the
   Run unit and return `stopQueuedVerification`; on another error release the
   Run unit and wrap the error with
   `acquire Machine Verification Capacity for run %q Task %s attempt %d retry %d`;
   return a release that gives the machine unit back before the Run unit;
   keep the context and Stop Request re-checks after both units are held.
   With a nil dependency every Verification behaves exactly as before.
   MUST NOT take a machine unit while waiting for the Run's own unit
   (Invariant 2), and the exclusive retry takes exactly one machine unit.
4. MUST, on the first `onWait` of an acquisition, pause the Run Budget and
   publish through `publishDaemonEvent` one `daemon.status` Run Event whose
   summary is
   `<step> waits for the Machine Verification Capacity: <held> of <capacity> held by <holders>.`
   (`<holders>`: each holding entry as
   `<DescribeMachineVerificationHolder>, <step>`, joined by `; `) with
   payload keys `machine_verification` (`waiting`), `step`, `capacity`,
   `held`, `waiting_ahead` and `holders`, and write the same summary as one
   line to `Progress`; when the acquisition returns after a wait, resume the
   Run Budget and publish `<step> took the Machine Verification Capacity after waiting <seconds>s.`
   with payload keys `machine_verification` (`acquired`), `step` and
   `waited_seconds`.
5. MUST add `pause` to `taskCycleBudget` as API Contract 5 says: while
   paused, `expiredAt` answers false and `watch` does not cancel; on the last
   resume an enabled, unexpired deadline moves later by the paused length and
   `changed` is signalled; a `renew` while paused restarts the paused length
   from the settlement time. Every existing budget test passes unchanged.
6. MUST wire the production gate in `internal/cli` (a new
   `internal/cli/machine_verification.go`) as API Contract 6 says: it calls
   `store.AcquireMachineVerification` with the loaded Roundfix Home,
   `Runs.MaxVerifying` and the Implement Run's checkout root as the holder's
   repository, and its release writes `roundfix: warning: <release error>` to
   standard error on failure without failing the Task; pass it to
   `daemon.Dependencies` in `executeImplementCycle` (`internal/cli/implement.go`).
   No other engine construction gains it.
7. MUST create `internal/daemon/machine_verification_test.go` with the four
   tests `_techspec.md` → Testing Approach names for `internal/daemon`, using
   a fake gate that records each request and blocks until the test grants it,
   the existing Task cycle fixtures, fake verifiers and the fake budget clock
   and timer of `task_budget_renewal_test.go`:
   `TestEveryDaemonVerificationStepTakesTheMachineCapacityAfterTheRunsOwn`
   (one row per step with its exact label; the fake verifier records no
   command start before the grant; with Verification Capacity 1 and two
   ready Tasks, the second Task's machine request is never made while the
   first Task's machine unit is held),
   `TestAStopRequestWhileWaitingForTheMachineCapacityStartsNoCommand`,
   `TestTheRunBudgetDoesNotRunWhileAVerificationWaitsForTheMachineCapacity`
   (on `taskCycleBudget` itself: a one-minute budget paused 30 s after its
   start and resumed two minutes later answers `expiredAt` false during the
   pause and at 2 min 59 s after the start, and true at 3 min; and its
   `watch`, given a fake timer that fires during the pause, does not cancel)
   and
   `TestTheMachineCapacityWaitIsJournaledAndPrinted` (both summaries and
   payload keys verbatim). Create `internal/cli/machine_verification_test.go`
   with `TestImplementGivesTheEngineTheUserConfigMachineCapacity` (the gate
   built from a loaded User Config with `max_verifying: 1` writes one holding
   `run` entry under the temporary home naming the checkout root, and with
   `0` writes no directory). Every top-level test calls `t.Parallel()` first.
8. MUST keep every existing test of `./internal/daemon` and every existing
   Implement test of `./internal/cli` passing unchanged; a test that needs
   two Runs to verify at once in one temporary home sets
   `runs.max_verifying: 0` in its User Config, and the Task names each such
   test in its `## Result` (none is known at authoring time).
9. MUST add `## Machine Verification Capacity` to
   `docs/user-guide/configuration.md` right after
   `## Active Implement Run ceiling`, covering `_techspec.md` → API
   Contract 9 for task_02, and containing these sentences word for word:
   "After the update every machine verifies one repository at a time; set
   `runs.max_verifying` to a higher number, or to `0`, to change that.";
   "Waiters take a unit in arrival order."; "A holder that is still alive
   keeps its unit however long it runs; stop its Run with `roundfix stop` to
   free it."; "While a Verification waits for the Machine Verification
   Capacity, its Run's `budget.max_run_duration` does not run."; and the
   warning line
   `config: runs.max_verifying in Project Config is ignored; set runs.max_verifying in User Config`.
   The section lists the gated steps and the ungated ones exactly as
   ADR-0384's Decision does. In `## Task and Verification capacities`, the
   paragraph beginning "Verification Capacity is not a machine-wide lock"
   gains the sentence "The Machine Verification Capacity below bounds
   Verifications across repositories."
10. MUST change item 9 of the canonical
    `.agents/skills/roundfix/references/implement.md` (Task Capacity and
    Verification Capacity) to describe `runs.max_verifying`, its default,
    `0`, the order after the Run's own capacity, the exclusive retry's one
    unit, the Run Budget pause and the wait event summary, naming the
    Machine Verification Capacity; then run `make skills-sync` and raise the
    Roundfix skill version with
    `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
    Never write inside a `### QA settlement` section, never hand-edit a
    mirror, version field or record.
11. MUST, as the Task's last step, after its last guide and skill edit,
    re-record the Behavior Surface Record with
    `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
    Of the Governed Paths the Task may change only
    `.agents/skills/roundfix/SKILL.md`,
    `.agents/skills/roundfix/references/implement.md` and
    `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Characterize the callers of the Run's capacity.
- [ ] Add the gate dependency and the step labels.
- [ ] Take the machine unit after the Run's own, with Stop Request handling and the budget pause.
- [ ] Journal and print the wait.
- [ ] Wire the production gate in the Implement Command.
- [ ] Add the Daemon and CLI tests.
- [ ] Write the configuration guide section and the `implement` reference, sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] No Daemon Verification command starts before both units are held, and the machine unit is released first.
- [ ] A Stop Request during a machine wait starts no command, and the Run Budget does not run during the wait.
- [ ] The Implement Command takes its machine units from `runs.max_verifying` in the Roundfix Home.
- [ ] The configuration guide and both copies of the `implement` reference describe the capacity.

## Context

- instruction: `docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md`
- instruction: `docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md`
- instruction: `internal/daemon/task_budget_renewal_test.go`
- instruction: `internal/daemon/task_engine_test.go`
- instruction: `internal/daemon/integrated_repository_verification_test.go`
- instruction: `internal/daemon/settlement_repository_verification_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/daemon/machine_verification_test.go`
- interface: `internal/cli/implement.go`
- creates: `internal/cli/machine_verification.go`
- creates: `internal/cli/machine_verification_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/implement.md`
- interface: `skills/roundfix/references/implement.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `daemon.MachineVerificationGate` and `daemon.Dependencies.MachineVerification` — consumers: `./internal/cli`, `./internal/daemon`
- The `daemon.status` wait events — consumers: `./internal/daemon`
- The configuration guide and the `implement` reference — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 30m -v -run '^(TestEveryDaemonVerificationStepTakesTheMachineCapacityAfterTheRunsOwn|TestAStopRequestWhileWaitingForTheMachineCapacityStartsNoCommand|TestTheRunBudgetDoesNotRunWhileAVerificationWaitsForTheMachineCapacity|TestTheMachineCapacityWaitIsJournaledAndPrinted)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryDaemonVerificationStepTakesTheMachineCapacityAfterTheRunsOwn TestAStopRequestWhileWaitingForTheMachineCapacityStartsNoCommand TestTheRunBudgetDoesNotRunWhileAVerificationWaitsForTheMachineCapacity TestTheMachineCapacityWaitIsJournaledAndPrinted; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'Budget|VerificationCapacity|IntegratedRepositoryVerification|PreWork|Precondition|QARepository' ./internal/daemon` — expected: exit 0. Before this Task the four tests do not exist; afterwards they pass, and the existing budget, capacity, integrated, pre-work and precondition tests pass unchanged.

- `out="$(go test -count=1 -timeout 30m -v -run '^TestImplementGivesTheEngineTheUserConfigMachineCapacity$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestImplementGivesTheEngineTheUserConfigMachineCapacity ' || { echo "missing pass: TestImplementGivesTheEngineTheUserConfigMachineCapacity" >&2; exit 1; }; go test -count=1 -timeout 40m -run 'Implement|RunCeiling' ./internal/cli` — expected: exit 0. Before this Task the test does not exist; afterwards it passes, and every existing Implement test of the CLI passes.

- `tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qE -- 'After the update every machine verifies one repository at a time; set .runs.max_verifying. to a higher number, or to .0., to change that.' || { echo 'missing: configuration sentence' >&2; exit 1; }; grep -q '^## Machine Verification Capacity$' docs/user-guide/configuration.md || { echo 'missing: configuration section' >&2; exit 1; }; for f in .agents/skills/roundfix/references/implement.md skills/roundfix/references/implement.md; do grep -qF 'runs.max_verifying' "$f" || { echo "missing key: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither the guide nor the reference names the key; afterwards both copies of the reference and the guide describe it, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record pass.

## References

- `_prd.md` → Goal 1; Goal 2; User Story 1; User Story 2; Core Feature 2; Core Feature 6; Core Feature 7; Success Metric 3; Success Metric 4; Success Metric 7
- `_techspec.md` → Current behavior 1, 2 and 5; Interfaces; API Contracts 4-6; API Contract 9; Invariants 1 and 2; Declared breaks; Testing Approach; Build Order 2
- ADR-0384
- ADR-0056

## Result

Skills loaded: implement-task, coding-guidelines, golang-concurrency, golang-context, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, golang-cli, testing-boss, agentic-cli-design, tech-writer, evidence-gate.

Starting-tree characterization (before implementation): `acquireVerificationCapacity`
is the only function that acquires Verification Capacity. Its callers are
`verifyTaskPreWork` and `runTaskVerificationRequest`. The callers of
`runTaskVerificationRequest` are `verifyIntegratedRepository` (normal and exclusive
retry), `verifyRepositoryPrecondition`, `verifyTask` (normal and exclusive retry),
and `runQARepositoryVerification`. The starting tree has no machine gate or budget
pause. Inspection: `rg` for those three acquisition/request symbols in
`internal/daemon/task_engine.go` and the request and dependency types in
`internal/daemon/engine.go`.

Initial worktree: only this Task file's Daemon-owned status differs from HEAD.
The PRD/TechSpec constraints and approved authorization bound the owned-skill
changes to the three Governed Paths this Task declares. Implementation evidence
will be appended below; Task status and declared Verification belong to the Daemon.

### Implementation and acceptance evidence

- Both units before commands, machine released first: added the optional
  `MachineVerificationGate` dependency and the four explicit step labels; ordinary
  Task attempts derive their label from `identity()`. The existing Run acquisition
  path remains unchanged in `acquireRunVerificationCapacity`, including its nil-gate
  behavior. The outer acquisition waits for the Run first, then takes one machine
  unit, rechecks cancellation/Stop Request after both grants and releases in reverse
  order. The step-table test exercises pre-work, repository precondition, Task
  settlement (including the repository command), Task exclusive retry, Integrated
  Repository Verification (including retry) and QA repository Verification. The
  fake verifier rejects any command without a granted machine lease. The fake
  machine gate checks the Run unit is already held at acquisition and still held
  at release. The two-ready-Tasks row observes the second Task queue for the Run
  before checking that it has issued no machine request while the first holds it.
- Stop Request and Run Budget: Stop/cancel/error/grant-race rows start no command,
  preserve pending Task status and return the Run unit. Wait callbacks pause the
  budget only once per acquisition. The budget counts overlapping pauses, resumes
  on acquisition success and failure, extends the deadline only at the last resume
  and restarts pause accounting at a timely settlement renewal. Tests assert a
  one-minute budget paused after 30 seconds expires exactly at three minutes after
  a two-minute pause; overlapping/duplicate resumes, renewal during pause, disabled
  budgets and a fake watchdog timer firing during a pause are also covered.
- Implement configuration: only `executeImplementCycle` receives the production
  gate. It uses the loaded Home, `Runs.MaxVerifying` and the preflight checkout
  root. The CLI test loads real User Config values `1` and `0` in temporary homes,
  observes the holding `run` entry and checkout identity at `1`, no capacity
  directory at `0`, and removal after cancellation. A separate test forces release
  failure and observes the `roundfix: warning:` diagnostic without a Task error.
- Wait observability: the journal/progress test checks both exact summaries and
  all payload keys, two holder descriptions joined by `; `, waiting-ahead count,
  first-callback deduplication and a measured twelve-second wait.
- Guide and references: added the Machine Verification Capacity section immediately
  after Active Implement Run ceiling, the required sentences and Project Config
  warning, gated/ungated steps, FIFO order, owner-death proof and per-Home scope.
  Item 9 of the canonical implement reference describes the machine bound, order,
  exclusive retry, budget pause and wait events. The generated mirror matches it.

### Focused checks

All commands below ran in this worktree. These are implementation checks, not
Daemon Verification or terminal QA.

- `go test -count=1 -race -timeout 8m ./internal/daemon ./internal/cli -run 'Test(EveryDaemonVerificationStep|AStopRequestWhileWaitingForTheMachine|TheRunBudgetDoesNotRun|TheMachineCapacityWait|NilMachineVerification|ImplementGivesTheEngine|MachineVerificationRelease)'`
  — exit 0, daemon 3.551 s and CLI 2.042 s, including the final code/test edits.
- `go test -count=1 -race -timeout 15m ./internal/daemon ./internal/cli -run '^Test(TaskBudget|TaskCycle.*Verification|PreWorkProbe|Repository.*Precondition|IntegratedRepository|QA.*Repository|Implement|RunCeiling)'`
  — exit 0, daemon 9.036 s and CLI 13.440 s. Existing tests were unchanged;
  no existing test needed a `runs.max_verifying: 0` override.
- `make skills-sync` — exit 0; both Roundfix `SKILL.md` files and both implement
  references subsequently compared byte-identical with `cmp`.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  — exit 0; the command raised Roundfix from `0.1.212` to `0.1.213` and updated
  its owned-version record. No mirror or version was hand-edited.
- `make baseline-digests` — exit 0, `changed: false`; no derived Baseline paths changed.
- `go test -count=1 -timeout 5m ./skills -run '^(TestCheckValidatesRoundfixSkillArtifacts|TestRoundfixSkillIndexNamesEveryReference|TestRoundfixSkillIndexLinksResolve|TestRoundfixSkillEntryFileStaysWithinItsBudget)$'`
  — exit 0.
- `go test -count=1 -timeout 5m -tags docscontract ./internal/docscontract -run '^TestUserGuideLinksResolve$'`
  — exit 0.
- `git diff --check` — exit 0. Changed/untracked-path inspection found only this
  Task's declared files. No sibling Task, Task Graph, status field or QA settlement
  section was edited. The only changed Governed Paths are the three this Task
  authorizes.

The sandbox initially denied the Go build cache and the canonical skill write;
the same operations succeeded through the required escalated tool calls.

Final regeneration, after all guide/skill edits and this evidence entry:
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
The tool transcript records its process outcome; only
`docs/references/behavior-surfaces.json` is its sanctioned output.

The declared Verification commands, full suites and terminal settlement remain
for the Daemon. No commit, push or Pull Request was made. No follow-up outside
this slice was required.
