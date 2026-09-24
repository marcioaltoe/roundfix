---
spec: 0165-a-measured-run-ceiling
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# A measured Run ceiling

On 2026-09-24 this machine ran five Implement Runs at once while the
maintainer's session ran the repository gates. Load average passed 15. The QA
gate of Spec 0157 then failed its precondition for a reason unrelated to the
Spec: `internal/cli` exceeded `go test`'s ten-minute default inside
`make verify`, and `TestRunImplementBudgetExceededPreservesRunWorktreeAndBranch`
failed three times out of three, on `main` as well. With three concurrent Runs
the same gates passed. Nothing in Roundfix bounds how many Runs a machine
starts: `worktree.concurrency` bounds Task Worktrees inside one Run, and the
Active Run lock is per checkout.

The budget test fails under load because it configures a 500 ms Run Budget and
then reads the Run id only from the `Implement Run: <id>` header, which a loaded
machine prints after the budget has already expired.

This Spec is delivery 7 of the restructured queue in the reduced scope the
maintainer chose on 2026-09-24: Spec 0124 Core Feature 3 as a measured ceiling,
and the budget-test flake that measurement exposed. Spec 0124 Core Feature 2 is
already met for Daemon Verification, which writes every attempt and exclusive
retry to its own diagnostic file, and Core Feature 6 is deferred.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; the ceiling counts
  existing Active Runs. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Run Database only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0056 separates Task and Verification
  Capacity inside one Run, which this machine-wide ceiling leaves unchanged;
  ADR-0093 checks Spec consistency by citation, ADR-0104 accepts on evidence a
  Spec did not author, ADR-0130 keeps a path governed once bounded, ADR-0155
  makes the `qa` Task declare the matrix and ADR-0156 makes a declared promise
  name a consuming Task. This Spec's gate is bound by ADR-0080, ADR-0091,
  ADR-0096, ADR-0097 and ADR-0117. All hold. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `GovernedPath` is empty; the Makefile and CI stay untouched as the
  maintainer required. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## Goals

- A machine never starts more Implement Runs than it was measured to carry.
- The budget test proves the budget, not the machine's speed.

## Core Features

1. **A machine-wide Active Run ceiling.** `runs.max_active` in User Config,
   default 3 as measured on 2026-09-24, bounds Active Implement Runs across every
   repository in the Run Database. `implement` refuses at the ceiling before
   creating a Run, naming the Active Runs that hold it; `0` disables the bound.
2. **A budget test that does not race the header.** The budget test identifies
   its Run without depending on the header line a loaded machine may not print
   before the budget expires.

## Non-Goals / Out of Scope

- A core-count formula or a raised deadline.
- Queueing a refused Run; the delivery loop decides when to retry.
- Spec 0124 Core Feature 6 and any Makefile, CI or cache change.

## Success Metrics

1. With `runs.max_active: 2` and two Active Runs in two repositories, a third
   `implement` refuses with both Runs named and creates nothing.
2. The budget test passes when the header line is absent from stderr.

## Decisions

- **Measure, then bound.** The default is the concurrency observed to keep the
  repository gate green on this machine, not a formula.
- **Refuse, don't queue.** A refusal is immediate and legible; waiting belongs to
  the caller that owns the schedule.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative case carries the weight: a ceiling counted per
repository would pass a single-repository test while letting a machine
overload.

## Research basis

Load, failures and the three-Run pass were observed on 2026-09-24 and captured
in the secondbrain inbox (`inbox/roundfix/2026-09-24-budget-test-flakes-under-load.md`).
`implementRunIDFromStderr` in `internal/cli/implement_test.go` reads only the
header line.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
