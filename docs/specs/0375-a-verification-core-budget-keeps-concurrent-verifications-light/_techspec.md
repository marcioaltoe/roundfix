---
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: active
created: 2026-11-10
---

# A Verification Core Budget keeps concurrent Verifications light — Technical Spec

## Executive Summary

A Go Verification takes every core because `GOMAXPROCS` and `go test -p`
default to the CPU count, and Spec 0374's Machine Verification Capacity
cannot gate an Agent's own suites. This Spec adds the
**Verification Core Budget**: a new `internal/corebudget` package computes
half of the logical CPUs, at least 1, divided by `runs.max_verifying`, or by
`runs.max_active` while the capacity is off, never below 1, unless a positive
whole-number `GOMAXPROCS` in the process environment overrides it. The
Implement Command computes it once and hands it to the Daemon, whose verifier
sets `GOMAXPROCS` and `ROUNDFIX_VERIFICATION_CORES` on every Verification
command and whose Agent Sessions carry both beside the Run Marker. The
Delivery Queue's repository gate sets both on its Verification and passes
both to the Linux Lane container. `roundfix doctor` prints a
`verification-cores:` line against the active Implement Runs. ADR-0385
records the decision.

The trade-off accepted: the budget is static per process, so with the
default a lone Verification uses only half of an idle machine; in exchange
the budgets that may run at once always fit in that half, and an Agent's
suite or a second Verification still finds free cores. A test runner that
reads neither variable is unaffected until its repository opts in.

## Project Constraints

- Identifier strategy: not applicable — no domain identifier is minted; the
  budget is a number computed per process and never stored. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0385 (new) records the decision:
  "Roundfix computes one number and exports it under two names." ADR-0384
  applies unchanged: "`runs.max_verifying` is the Machine Verification
  Capacity", which the budget divides by. Not applying: ADR-0304, ADR-0365
  and ADR-0373, which ADR-0384 cites and which govern the Run Slot Wait, the
  pre-PR review's capacity and the proof that an owner is gone; and ADR-0056,
  whose amendment by the capacity Spec cites ADR-0384 and which sets the per-Run
  Verification Capacity this Spec does not change. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — the cycle BB batch 2 grant quoted in
  `_authorization.md`; the Governed Paths are the Roundfix skill files tasks
  02 and 03 change, through the skill sync and version record. No Makefile,
  workflow, Go module, lint, formatter, test-runner or Baseline change.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0375-a-verification-core-budget-keeps-concurrent-verifications-light/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/setup.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured at `b233761e`, with Spec 0374's plan read as the tree this Spec
lands on.

1. `internal/daemon/daemon.go` `ExecVerifier.Verify` runs `sh -c <command>`
   with `judge.WithoutKeyVariables(os.Environ())` plus the Run Marker when
   `VerifyRequest.RunID` is set. Nothing sets `GOMAXPROCS`.
2. The Daemon reaches that verifier at two sites: `verificationAttempt` in
   `internal/daemon/engine.go`, which every Task attempt, exclusive retry,
   repository precondition, Integrated Repository Verification and QA
   repository Verification goes through, and `probeCommandsForRun` in
   `internal/daemon/verification_probe.go`, which the pre-work probe
   (`verifyTaskPreWork`) and `roundfix spec check --run-verification`
   (`ProbeCommands`) share.
3. Every Agent Session request goes through `runAgentSession` in
   `internal/daemon/agent_session_owner.go`. In `internal/agent/acpx_runner.go`
   the session-preparation and prompt paths append the Run Marker to the
   per-session overrides that `acpxCommandEnv` applies over the inherited
   environment.
4. `internal/cli/deliver_workflow.go` `commandDeliveryWorkflow.Gate` runs
   `defaults.verification` through `daemon.ExecVerifier` in the Queue Owner,
   then `runLinuxLane`, whose `docker run` argument list fixes the container
   environment (`HOME`, `GOCACHE`, `GOMODCACHE`, `GOFLAGS`, `GOPROXY`,
   `GOTOOLCHAIN`, `GOSUMDB`). `roundfix deliver lane` runs the same builder.
5. This repository's Makefile runs `go test -parallel 16` without `-p`; `go
   help build` says of `-p`: "The default is GOMAXPROCS, normally the number
   of CPUs available."
6. `roundfix doctor` prints no line about cores. Spec 0374 adds a
   `verification-capacity:` line after `disk-space:` and the key
   `runs.max_verifying`, `1` by default.

## System Architecture

- `internal/corebudget` (new package): the computation, its description and
  the environment rewrite; pure functions over their inputs.
- `internal/daemon`: `VerifyRequest.Cores`, applied by `ExecVerifier`;
  `Dependencies.VerificationCores`, passed at both verifier sites and set on
  every Agent Session request.
- `internal/agent`: `ExecuteRequest.Cores`, appended beside the Run Marker.
- `internal/cli`: the budget helper over the command dependencies, the
  Implement wiring, the Delivery Queue gate and Linux Lane, and the doctor
  line.
- Documentation: the configuration guide, the doctor and deliver command
  guides, three Roundfix references, `CONTEXT.md`.

## Implementation Design

### Interfaces

```go
// internal/corebudget
const (
	Variable   = "ROUNDFIX_VERIFICATION_CORES"
	GoVariable = "GOMAXPROCS"
)
const (
	SourceEnvironment  = "environment"
	SourceMaxVerifying = "runs.max_verifying"
	SourceMaxActive    = "runs.max_active"
	SourceUndivided    = "undivided"
)
type Inputs struct{ LogicalCPUs, MaxVerifying, MaxActive int; Environ []string }
type Budget struct {
	Cores, LogicalCPUs, Share, Divisor int
	Source                             string
	MaxVerifying, MaxActive            int
}
func Compute(in Inputs) Budget
func (budget Budget) Describe() string
func (budget Budget) Concurrent(activeRuns int) int
func Environment(environ []string, cores int) []string
```

`internal/daemon` gains `VerifyRequest.Cores int` and
`Dependencies.VerificationCores int`; `internal/agent` gains
`ExecuteRequest.Cores int`; `internal/cli` gains
`commandDependencies.logicalCPUs func() int` (default `runtime.NumCPU`),
`implementCapacities.cores`, `linuxLaneRequest.Cores`,
`doctorDependencies.verificationCores` and the check name
`HealthCheckVerificationCores = "verification-cores"`.

### Data Models

None — nothing is stored; the budget lives in the process that computed it
and in the environments it sets.

### API Contracts

1. **Computation** (`internal/corebudget`, task_01). `Compute`:
   - `LogicalCPUs` below 1 counts as 1; `Share` is `LogicalCPUs / 2`, at
     least 1.
   - When `Environ` holds `GOMAXPROCS` (the last entry wins) whose value is a
     non-empty string of ASCII digits naming an integer from 1 to 2147483647,
     `Cores` is that integer, `Source` is `environment` and `Divisor` is 1.
     Any other value, including `0`, a sign, letters or an empty value, is
     ignored. `ROUNDFIX_VERIFICATION_CORES` in `Environ` is never read.
   - Otherwise `Divisor` and `Source` are `MaxVerifying` and
     `runs.max_verifying` when `MaxVerifying > 0`, else `MaxActive` and
     `runs.max_active` when `MaxActive > 0`, else 1 and `undivided`, and
     `Cores` is `Share / Divisor`, at least 1.
   - `Describe` returns, by source: `GOMAXPROCS <cores> in the environment`;
     `half of <cpus> logical CPUs divided by runs.max_verifying <divisor>`;
     `half of <cpus> logical CPUs divided by runs.max_active <divisor> while runs.max_verifying is 0`;
     `half of <cpus> logical CPUs, undivided while runs.max_verifying and runs.max_active are 0`.
   - `Concurrent(activeRuns)` is `MaxVerifying` when above 0, else
     `MaxActive` when above 0, else `activeRuns`, at least 1.
   - `Environment(environ, cores)` returns a new slice: with `cores < 1` an
     unchanged copy; otherwise every `GOMAXPROCS` and
     `ROUNDFIX_VERIFICATION_CORES` entry removed, other entries in order, then
     `GOMAXPROCS=<cores>` and `ROUNDFIX_VERIFICATION_CORES=<cores>`. It never
     changes its input.
2. **Verifier** (`internal/daemon/daemon.go`, task_01). `ExecVerifier.Verify`
   sets the command environment to
   `corebudget.Environment(<today's environment>, req.Cores)` after the Run
   Marker is applied. With `Cores` 0 the environment is byte-identical to
   today's.
3. **Implement Run** (task_02).
   - `engine.go` `verificationAttempt` passes
     `Cores: engine.deps.VerificationCores`; `probeCommandsForRun` gains a
     `cores int` parameter after `verifier`, `verifyTaskPreWork` passes
     `engine.deps.VerificationCores` and `ProbeCommands` passes 0, so `spec
     check --run-verification` exports nothing.
   - `runAgentSession` sets `req.Cores = engine.deps.VerificationCores` before
     either branch.
   - `acpx_runner.go`: wherever the Run Marker is appended to the session
     overrides (session preparation and prompt), `coreBudgetOverrides(req.Cores)`
     is appended after it: nothing for `Cores < 1`, otherwise
     `GOMAXPROCS=<n>` and `ROUNDFIX_VERIFICATION_CORES=<n>`, which
     `acpxCommandEnv` puts in place of the inherited copies.
   - `internal/cli/verification_cores.go` (new): `verificationCoreBudget(ctx,
     config)` calls `Compute` with the command dependencies' `logicalCPUs`
     (`runtime.NumCPU` when nil), `config.Runs.MaxVerifying`,
     `config.Runs.MaxActive` and the dependencies' `environ()`.
   - `runImplementCommand` sets `implementCapacities.cores` from
     `verificationCoreBudget(ctx, loadedConfig.Config).Cores`, and
     `executeImplementCycle` passes it as `Dependencies.VerificationCores`.
     No other engine construction sets it, so review Runs export nothing.
4. **Delivery Queue** (task_03). `commandDeliveryWorkflow.Gate` computes
   `cores := verificationCoreBudget(ctx, workflow.loaded.Config).Cores` once
   and passes it as `VerifyRequest.Cores` to the repository Verification and
   as `linuxLaneRequest.Cores` to the lane; `roundfix deliver lane` passes its
   own budget too. `runLinuxLane` inserts, when `Cores > 0`, the four
   arguments `-e`, `GOMAXPROCS=<n>`, `-e`, `ROUNDFIX_VERIFICATION_CORES=<n>`
   right after `-e GOSUMDB=off` and before the first `-v`. The item binary
   build and derived-file regeneration keep `Cores` 0.
5. **Doctor** (task_03). `doctorDependencies.verificationCores`, default
   `defaultDoctorVerificationCoresResult` in the new
   `internal/cli/doctor_verification_cores.go`, prints one result right after the
   `verification-capacity:` line (after `disk-space:` when that line is
   absent), only when the dependency is not nil. It computes the budget from
   the loaded User Config and the command dependencies, counts
   `store.ActiveImplementRuns` in a read-only open of the Run Database (no
   Run Database is 0 active Runs), and formats:
   - `ok` detail `<cores> per Verification: <Describe>; <k> active Implement Runs`;
   - `warn` when `Concurrent(k) * cores > LogicalCPUs`, the detail above
     followed by `; <c> at once may take <c×cores> of <cpus> logical CPUs`,
     next action `lower GOMAXPROCS in the environment or set runs.max_verifying in User Config`;
   - `warn` when the Run Database cannot be read:
     `<cores> per Verification: <Describe>; active Implement Runs unknown: <error>`,
     next action `inspect the Run Database with roundfix runs list`.
   Doctor stays read-only and the line never changes its exit code.
6. **Documentation** (tasks 02 to 04). task_02:
   `docs/user-guide/configuration.md` gains `## Verification Core Budget`
   right after `## Machine Verification Capacity`, and item 9 of the
   canonical `implement` reference gains one paragraph; task_03: the doctor
   and deliver command guides and the canonical `deliver` and `setup`
   references; task_04: **Verification Core Budget** in `CONTEXT.md` with
   the text of API Contract 7. Every required sentence is listed in its
   Task.
7. **Glossary text** (task_04), right after the entry
   Machine Verification Capacity, exactly: "The number of cores Roundfix
   gives each Verification it starts for an Implement Run or for the
   repository gate of a Delivery Queue, and each Agent Session of an
   Implement Run, as `GOMAXPROCS` and `ROUNDFIX_VERIFICATION_CORES`. It is
   half of the logical CPUs, at least 1, divided by `runs.max_verifying`, or
   by `runs.max_active` while that is `0`, and never below 1; a positive
   whole-number `GOMAXPROCS` in the environment of the Roundfix process
   takes its place (ADR-0385)." The closing `(ADR-0385)` is the glossary's
   usual Markdown link to
   `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md`.
   `_Avoid_: core limit, CPU quota, Machine Verification Capacity`

### Surface Transcripts

None. The doctor line is the only new command output, and the real
`roundfix doctor` proves Agent Selection Profiles by starting ACP adapters,
which no Verification or QA row may do; `TestDoctorReportsTheVerificationCoreBudget`
asserts each prescribed line whole.

### Invariants

1. With a budget of `n`, every Verification command of an Implement Run, its
   Agent Sessions, the repository gate's command and the Linux Lane container
   see exactly one `GOMAXPROCS` entry and one `ROUNDFIX_VERIFICATION_CORES`
   entry, both `n`.
2. With a budget of 0 the verifier's environment and the acpx overrides are
   what they are today.
3. The budget is at least 1 and, unless `GOMAXPROCS` overrides it, the
   budget times its divisor is at most the share whenever the share is at
   least the divisor.
4. `Compute` reads only its inputs; doctor writes nothing.

### Declared breaks

- Every Verification of an Implement Run and of the Delivery Queue's gate,
  and every Agent Session of an Implement Run, now runs with `GOMAXPROCS`
  and `ROUNDFIX_VERIFICATION_CORES` set; a Go suite on a twelve-core machine
  with the default runs at most six package binaries at once.
- `roundfix doctor` prints one more line, `verification-cores:`.
- The Linux Lane's `docker run` argument list gains four arguments; the test
  that pins it, `TestLinuxLaneRunsTheChangedPackagesInAContainer`, changes
  with task_03.

## Coverage Map

- Core Feature 1, Core Feature 2, User Story 4, Goals 2 and 4,
  Success Metric 1 → API Contract 1, Invariant 3,
  `TestTheVerificationCoreBudgetIsHalfTheLogicalCPUsDividedByTheMachineCapacity`,
  `TestAnExplicitGOMAXPROCSWinsTheVerificationCoreBudget`,
  `TestTheCoreBudgetEnvironmentReplacesBothVariables` (task_01).
- Core Feature 3, Core Feature 4, User Story 1, User Story 2, Goal 1,
  Success Metric 2 → API Contracts 2 and 3, Invariants 1 and 2,
  `TestTheVerificationCommandReceivesTheCoreBudget` (task_01),
  `TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget`,
  `TestAnAgentSessionReceivesTheCoreBudget`,
  `TestImplementGivesTheEngineTheVerificationCoreBudget` (task_02).
- Core Feature 5, User Story 3, Goal 3, Success Metric 3 → API Contract 4,
  `TestTheRepositoryGateCommandReceivesTheCoreBudget`,
  `TestLinuxLaneRunsTheChangedPackagesInAContainer` (task_03), and the
  configuration guide's opt-in example (task_02).
- Core Feature 6, User Story 5, Goal 5, Success Metric 4 → API Contract 5,
  Invariant 4, `TestDoctorReportsTheVerificationCoreBudget` (task_03).
- Core Feature 7 → API Contracts 3 and 4: only the Implement Command and the
  queue gate set a budget, and no key is added (tasks 02 and 03).
- Success Metric 5 → API Contracts 6 and 7 (tasks 02 to 04).

Deviations from the PRD: None.

## Absence Evidence

1. `internal/corebudget` and the variable — command:
   `ls internal/corebudget; grep -rn 'ROUNDFIX_VERIFICATION_CORES\|corebudget' internal docs .agents skills`;
   result: no such directory, no match.
2. A budget on Verification commands or Agent Sessions — command:
   `grep -rn 'GOMAXPROCS\|NumCPU' internal --include=*.go`;
   result: no match outside tests (the Makefile comment on `-parallel` is the
   only mention in the repository's own tooling).
3. The doctor line — command:
   `grep -rn 'verification-cores' internal docs .agents skills`;
   result: no match.
4. The glossary term — command:
   `grep -n 'Verification Core Budget' CONTEXT.md`; result: no match.

## ADR Alternatives

- ADR-0385: the decision this Spec records; it rejects a fixed number,
  dividing by the current holders, a share key, overriding an explicit
  `GOMAXPROCS`, leaving Agent Sessions out and a Baseline clause, and the
  design takes none of them.
- ADR-0384: applies unchanged; it rejected a core budget alone as a
  substitute for the capacity, and the design keeps the capacity and only
  divides by its value, taking none of its rejected options.
- ADR-0304: not applicable; it set the Run Slot Wait, which this Spec does
  not change, so none of its rejected options is taken.
- ADR-0365: not applicable; it set the pre-PR review's own capacity, which
  this Spec does not change, so none of its rejected options is taken.
- ADR-0373: not applicable; it decided when a recorded owner is gone, which
  this Spec does not read, so none of its rejected options is taken.
- ADR-0056: not applicable; it set the per-Run Verification Capacity, which
  this Spec does not change, and the capacity Spec's amendment makes it
  cite ADR-0384, so none of its rejected options is taken.

## Integration Points

`GOMAXPROCS` is read by the Go runtime and, through it, by `go test -p`; the
pkg.go.dev runtime documentation says "If the GOMAXPROCS environment variable
is set to a positive whole number, GOMAXPROCS defaults to that value."
`ROUNDFIX_VERIFICATION_CORES` is read only by a test runner whose repository
configures it. Docker receives both through `-e`. No test reaches Docker, the
network or a real ACP adapter.

## Testing Approach

Every Go test uses temporary homes and repositories, fake runners, fake
verifiers, the existing fake acpx and fake docker, and harmless shell
commands. A test that computes a budget from command dependencies sets
`logicalCPUs` and an environment without `GOMAXPROCS` or
`ROUNDFIX_VERIFICATION_CORES` unless the row supplies one, because the suite
itself runs inside a Verification that carries a budget.

- `internal/corebudget/corebudget_test.go` (new, task_01): the three
  `Compute`/`Environment` tests of the Coverage Map, with rows twelve CPUs at
  capacities 1 and 2, thirteen CPUs, capacity off with `max_active` 4, both
  off, one CPU, four CPUs at capacity 5; explicit values `4`, `10`, a last
  duplicate, `0`, `-2`, `+4`, `four`, empty, and the Roundfix variable alone.
- `internal/daemon/verification_cores_test.go` (new): task_01
  `TestTheVerificationCommandReceivesTheCoreBudget` (an `ExecVerifier` command
  writes `${GOMAXPROCS-unset}|${ROUNDFIX_VERIFICATION_CORES-unset}|<count of GOMAXPROCS entries>`
  to a file: `4|4|1` at budget 4, the inherited values at 0).
- `internal/daemon/run_core_budget_test.go` (new, task_02):
  `TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget` (a Task
  cycle with a recording verifier and the repository Verification at
  settlement: the pre-work probe, the settlement command and every Agent
  Session request carry 5 at budget 5 and 0 at 0).
- `internal/agent/core_budget_test.go` (new, task_02), with a recorder hooked
  into the fake acpx of `acpx_runner_test.go`:
  `TestAnAgentSessionReceivesTheCoreBudget` (an inherited `GOMAXPROCS=12`;
  `sessions ensure` and `prompt` see `3` and `3` at budget 3, and `12` with no
  Roundfix variable at 0).
- `internal/cli/verification_cores_test.go` (new, task_02): the hermetic
  budget-input helper (twelve logical CPUs, no inherited variable) and
  `TestImplementGivesTheEngineTheVerificationCoreBudget` (`roundfix
  implement` with `runs.max_verifying: 2` gives every Verification and Agent
  Session request 3, and 7 with `GOMAXPROCS=7`).
- `internal/cli/deliver_core_budget_test.go` (new, task_03):
  `TestTheRepositoryGateCommandReceivesTheCoreBudget` (`3|3`, and `5|5` with
  `GOMAXPROCS=5`).
- `internal/cli/doctor_verification_cores_test.go` (new, task_03):
  `TestDoctorReportsTheVerificationCoreBudget` (the six lines below, whole,
  plus the default dependency in a home without a Run Database).
- `internal/cli/deliver_linux_lane_test.go` (task_03): the fixture context
  gives eight logical CPUs and the pinned argument list gains
  `-e GOMAXPROCS=4 -e ROUNDFIX_VERIFICATION_CORES=4` after `-e GOSUMDB=off`.

The doctor rows, whole lines on twelve logical CPUs:

```text
verification-cores: ok (6 per Verification: half of 12 logical CPUs divided by runs.max_verifying 1; 2 active Implement Runs)
verification-cores: ok (2 per Verification: half of 12 logical CPUs divided by runs.max_active 3 while runs.max_verifying is 0; 0 active Implement Runs)
verification-cores: warn (6 per Verification: half of 12 logical CPUs, undivided while runs.max_verifying and runs.max_active are 0; 3 active Implement Runs; 3 at once may take 18 of 12 logical CPUs; next: lower GOMAXPROCS in the environment or set runs.max_verifying in User Config)
verification-cores: ok (12 per Verification: GOMAXPROCS 12 in the environment; 1 active Implement Runs)
verification-cores: warn (8 per Verification: GOMAXPROCS 8 in the environment; 1 active Implement Runs; 2 at once may take 16 of 12 logical CPUs; next: lower GOMAXPROCS in the environment or set runs.max_verifying in User Config)
verification-cores: warn (6 per Verification: half of 12 logical CPUs divided by runs.max_verifying 1; active Implement Runs unknown: open Run Database: denied; next: inspect the Run Database with roundfix runs list)
```

Sabotages, each a one-line mutation the QA gate applies to a temporary copy
of the audited head, with the named test that must then fail. The author ran
each on a prototype of tasks 01 to 03 on 2026-11-10 and saw the named test
fail and pass without the mutation:

1. `Compute` never reads an explicit `GOMAXPROCS` →
   `TestAnExplicitGOMAXPROCSWinsTheVerificationCoreBudget`, row
   `positive whole number`;
2. a leading `+` counts as a digit →
   `TestAnExplicitGOMAXPROCSWinsTheVerificationCoreBudget`, row
   `signed is ignored`;
3. `runs.max_verifying` is never the divisor →
   `TestTheVerificationCoreBudgetIsHalfTheLogicalCPUsDividedByTheMachineCapacity`,
   row `two at once`;
4. the share is every logical CPU →
   `TestTheVerificationCoreBudgetIsHalfTheLogicalCPUsDividedByTheMachineCapacity`,
   row `default on twelve`;
5. `Environment` keeps the inherited `GOMAXPROCS` →
   `TestTheCoreBudgetEnvironmentReplacesBothVariables`;
6. `ExecVerifier` applies a budget of 0 →
   `TestTheVerificationCommandReceivesTheCoreBudget`, row `budget four`;
7. `verificationAttempt` sends `Cores` 0 →
   `TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget`;
8. the pre-work probe sends `Cores` 0 →
   `TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget`;
9. `runAgentSession` leaves `Cores` unset →
   `TestEveryRunVerificationAndAgentSessionReceivesTheCoreBudget`;
10. `coreBudgetOverrides` returns nothing →
    `TestAnAgentSessionReceivesTheCoreBudget`, row `budget three`;
11. the Implement Command leaves `implementCapacities.cores` unset →
    `TestImplementGivesTheEngineTheVerificationCoreBudget`;
12. the gate's Verification gets `Cores` 0 →
    `TestTheRepositoryGateCommandReceivesTheCoreBudget`;
13. `runLinuxLane` never adds the two variables →
    `TestLinuxLaneRunsTheChangedPackagesInAContainer`;
14. doctor never warns →
    `TestDoctorReportsTheVerificationCoreBudget`, row
    `undivided with three Runs`.

Every other existing test of `./internal/daemon`, `./internal/agent`, and the
Implement, Doctor, Deliver and Linux Lane tests of `./internal/cli` passes
unchanged; the prototype confirmed it.

## Build Order

1. task_01: `internal/corebudget`, `VerifyRequest.Cores` and their tests;
   depends on nothing in this Spec but lands after Spec 0374 (`requires`).
2. task_02: the Daemon and Agent Session wiring, the Implement Command, the
   configuration guide section and the `implement` reference; depends on step
   1 and on Spec 0374's `Runs.MaxVerifying` and configuration section.
3. task_03: the Delivery Queue gate, the Linux Lane, the doctor line, the
   doctor and deliver guides and the `deliver` and `setup` references;
   depends on step 1, and runs after step 2 because both change the Roundfix
   skill version and the Behavior Surface Record, and it calls task_02's
   `verificationCoreBudget` and test helper; it places its doctor line after
   Spec 0374's `verification-capacity:` line.
4. task_04: the glossary entry; depends on step 3, so the term lands once
   every guide that uses it exists, after Spec 0374's
   Machine Verification Capacity entry.
5. task_05: the QA gate.

Spec 0374 is required: its Verification needs `runs.max_verifying` in the
loaded config and its sections and lines to place this Spec's text after. It
and Spec 0373 change other functions and sections of `internal/daemon/engine.go`,
`internal/daemon/task_engine.go`, `internal/cli/implement.go`,
`internal/cli/deliver_workflow.go`, `internal/cli/doctor.go`,
`internal/cli/health.go`, the three command and configuration guides, the
Roundfix `implement`, `deliver` and `setup` references, the `SKILL.md` pair,
the owned-skill version record, the Behavior Surface Record and `CONTEXT.md`;
each Task regenerates the records from the tree it runs on.

## Risks & Considerations

- A shell or queue owner that exports `GOMAXPROCS` keeps deciding the number:
  the operator's interim `GOMAXPROCS` on the queue owner and a Supervisor
  shell started with `GOMAXPROCS=2` both win. The operator drops the interim
  measure after the release.
- With the default, a lone Verification on an idle machine uses half of it;
  `GOMAXPROCS` set higher by the owner restores more.
- A test runner that reads neither variable keeps its own parallelism.
- The budget does not bound the processes a test starts that ignore
  `GOMAXPROCS`, such as `git` or shell children.
- Roundfix's own tests compute budgets from command dependencies; every new
  test sets them, because its own Verification will carry a budget.

## Validation

- Secondbrain, read only: `wiki/index.md`, then `qmd query` for a core budget
  shared by test verifications, for test-runner workers under machine load
  and for limiting cores across parallel projects. Relevant results:
  `projects/conexus/mirror/docs/references/test-execution-capacity.md`
  (Conexus limited Vitest to `maxWorkers: 2` after joint runs timed out, and
  only after "Diego autorizou essas duas linhas", a grant committed before the
  change), which is why the Baseline does not tell adopters' runners to read
  the variable and the user guide shows the opt-in instead;
  `projects/vortex/mirror/docs/references/2026-11-10-autorizacao-limite-de-tempo-dos-testes-do-backend.md`
  (the 30 s `testTimeout` Vortex adopted the same day, also an expressly
  authorized test-runner change, keeping "o teto de workers"); and this
  repository's mirrored decision that a test bounds by real time only what it
  asserts, the rule the new tests follow.
- Exa: the Go runtime documentation and the Go blog's container-aware
  `GOMAXPROCS` post (`GOMAXPROCS` is a parallelism limit and an explicit value
  disables the automatic default), Jest's `getMaxWorkers` (half the CPUs in
  watch mode, "unobtrusive"), Vitest's `maxWorkers` page (half of the
  available parallelism in watch mode) and Jest pull request 15944 (a worker
  count read from `JEST_MAX_WORKERS`). They confirm that Go obeys the
  variable alone, that half is the share runners take beside other work, and
  that other runners read a count only when configured to.

## Glossary

- adds: **Verification Core Budget**

## Decisions

Autonomous route: no person present; decisions taken from `_prd.md` → Decisions, the operator's cycle BB batch 2 brief for Spec 0375, Spec 0374's plan and the characterization above.

- A new package holds the computation, because three packages (`daemon`
  through its caller, `cli` twice) read it and none of them owns cores.
- The budget travels in the existing request types (`VerifyRequest`,
  `ExecuteRequest`) instead of a wrapper verifier or runner, so the fakes the
  tests already use observe it.
- The pre-work probe shares the budget; `spec check --run-verification`,
  which reaches the same function, passes 0.
- The Agent Session path reuses the Run Marker's overrides, which already
  replace inherited copies.
- The doctor line counts active Implement Runs from the Run Database, the
  number the vortex request named, and warns only on oversubscription.
- No Baseline change, no key, no Makefile change.
