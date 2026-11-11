---
status: accepted
created_at: 2026-11-11T01:30:00Z
updated_at: 2026-11-11T01:30:00Z
deprecated_at: null
superseded_by: null
---

# Each Verification receives a Verification Core Budget in its environment

On 2026-11-10 the Delivery Queue of vortex stopped three times on
time-sensitive tests while a Run of another repository verified on the same
twelve-core machine, and `go test` used every core because `GOMAXPROCS` and
its package parallelism default to the number of CPUs. ADR-0384 lets only
`runs.max_verifying` Verifications run at once, but the Agents' own commands
still run beside them and an owner can raise or turn off that bound. The
maintainer asked, beside the capacity, for "o orçamento exportado para as
verificações que ainda rodarem juntas", read as half of the cores for one
repository's Verification.

## Decision

**Roundfix computes one number and exports it under two names.** The
Verification Core Budget is half of the machine's logical CPUs, at least 1,
divided by `runs.max_verifying` when it is above `0`, otherwise by
`runs.max_active` when that is above `0`, otherwise not divided, and never
below 1. Integer division rounds down. Roundfix sets `GOMAXPROCS`, which the
Go runtime and `go test -p` obey without any project change, and
`ROUNDFIX_VERIFICATION_CORES`, which any other test runner can read, to that
number, replacing any copy the environment carried.

**An explicit `GOMAXPROCS` wins.** When the environment of the Roundfix
process that starts the work holds `GOMAXPROCS` as a positive whole number,
that number is the budget and both variables carry it. Any other value is
ignored and replaced. The budget is computed once per process, when an
Implement Run or a Delivery Queue gate starts, from its loaded User Config.

**It reaches the Verifications and Agent sessions of an Implement Run and the
Delivery Queue's repository gate.** Every Verification command the Daemon
starts for an Implement Run, the Agent sessions of that Run through the same
path that carries the Run marker, and the Delivery Queue's repository
Verification and Linux Lane container receive it. `roundfix spec check
--run-verification`, `roundfix settle`, the queue's item binary build and
derived-file regeneration, Worktree Bootstrap and the review Runs receive
nothing new.

**`roundfix doctor` reports it against the active Runs.** One
`verification-cores:` line names the budget, where it came from and how many
Implement Runs are active, and warns when the Verifications the bounds allow
at once, each with the budget, would need more cores than the machine has.

**The Baseline does not tell adopters' test runners to read the variable.**
Changing a test runner's configuration is protected tooling in every adopter,
so the user guide documents the variable with a Vitest example and each
repository opts in by its own authorized change.

## Considered Options

- A fixed number such as the interim `GOMAXPROCS=6` the operator set on the
  queue owner: rejected, because it is right only on a twelve-core machine.
- Dividing by the Verifications holding the Machine Verification Capacity at
  the moment each one starts: rejected, because the first Verification would
  keep every core while a later one got the rest, and the number would change
  between attempts of one Task.
- A User Config key for the share: deferred; an explicit `GOMAXPROCS` already
  sets the number, and no measurement asks for another share yet.
- Overriding an explicit `GOMAXPROCS`: rejected, because it would take the
  number away from the person who chose it.
- Exporting nothing to Agent sessions: rejected, because the QA Agent's own
  suites are the one Verification load ADR-0384 cannot gate.
- A Baseline clause telling every adopter's test runner to read
  `ROUNDFIX_VERIFICATION_CORES`: rejected for now, because it would ask
  agents to change test-runner configuration that each repository's
  maintainer must authorize.

## Consequences

- After the update a twelve-core machine with the default capacity of `1`
  gives each Verification six cores, and `go test` runs at most six package
  binaries at once with six threads each.
- A test runner that reads neither variable keeps its own parallelism until
  its repository configures it to read `ROUNDFIX_VERIFICATION_CORES`.
- A shell or queue owner that exports `GOMAXPROCS` keeps deciding the number
  for the Runs and gates it starts.
- `roundfix doctor` prints one more line, `verification-cores:`.
