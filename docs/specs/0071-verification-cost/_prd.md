---
spec: 0071-verification-cost
status: active
created: 2026-08-02
surfaces: [test, backend, infra, docs]
---

# Verification cost

Spec 0057 spent roughly five hours of Run time across four cycles, and about a
third of that was one line repeated in every Task. Measured on a warm cache:

| Command | Wall clock |
| --- | --- |
| `go build -buildvcs=false ./...` | 0s |
| `go test ./... -count=1`, warm | **136.9s** |
| `go test ./... -count=1`, **cold** | **146.2s** |
| `go test ./internal/cli -count=1` alone | **113.2s** |
| `go test ./internal/baseline -count=1` alone | **83.9s** |

Compilation is not the cost: cold versus warm is 9.3s. Packages already overlap,
so the suite can never finish faster than its slowest package — and
`internal/cli` alone is 113.2s of essentially sequential work, with **488 test
functions and one `t.Parallel()` call** on a twelve-core machine.
`internal/cli` and `internal/baseline` together are 84% of the cost; the
fifteen smallest packages sum to about eight seconds. The full baseline is
recorded under `baseline/`.

Every one of that Spec's fourteen Tasks carried one of those whole-package
commands as the last line of its Verification — about 28 minutes per pass,
paid again on every retry, to prove something the Run-level gate already
proves.

The suite itself is the other half, and the critical path is `internal/cli`.
Adding `t.Parallel()` to the subtests of the slowest Baseline test — each of
which already builds its own repository — took it from **29s to 17s**. Forty-one
Baseline tests already isolate their filesystem with `t.TempDir()`, and
`internal/cli` has 253 subtest blocks against one parallel call. Binary builds
there are already shared through `TestMain`, so repeated compilation is not the
cost.

A third cost is how often the gate closes the graph, and Spec 0072 owns the
mechanism. In summary: The Daemon withholds QA
correctly — Spec 0057's first Run ended with one Task failed and ran no gate at
all — so nothing runs early. What is missing is the other direction: the gate
is a flag on the command rather than part of the graph, so a graph that grows
*after* a gate has reported leaves no structural trace. On Spec 0057 a
corrective Task was appended after each gate, and three gates ran against three
different graphs at roughly twenty to twenty-five minutes each. Read from the
outside those look like three normal cycles instead of one decomposition that
was wrong twice. `docs/agents/autonomous-work.md` already warns about the
serial chain and already caps corrective Tasks at two; an advisory cap on a
flag is what allowed it to pass unnoticed.

None of the three is about having too many tests. Coverage is the asset; the
cost is how the suite is executed, how often each Task pays for it, and how
many times the gate is requested.

## Project Constraints

- Identifier strategy: not applicable — no project-owned Internal Identifier is
  created; test names, package paths, and target names keep their existing
  contracts. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the governing clause prohibits
  reading, printing, committing, or generating secrets and forbids inventing
  authentication, authorization, transport, or deployment policy; this Spec
  changes test execution and authoring guidance. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0081 keeps sanctioned digest
  regeneration a fallout of the authorized edit, which any change to the
  verification target must preserve. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: on
  2026-08-02 the maintainer authorized tooling adjustment naming the `Makefile`
  and the owned skills, recorded at
  `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`; bounded
  files: `Makefile` and the owned skill pair for the Verification authoring
  rule. Deterministic digest fallout is sanctioned by ADR-0081. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- The verification a Task pays for is proportional to what that Task changed.
- The suite uses the machine it runs on instead of one core at a time.
- Coverage is identical afterwards: the same tests, the same assertions, the
  same behaviors exercised.
- A regression in suite time is visible rather than discovered four Runs later.

## Core Features

1. Tests that own their filesystem and mutate no shared state run in parallel;
   a test that cannot states why in one line, so sequential execution is a
   recorded decision rather than an omission.
2. Parallelising surfaces shared-state defects instead of hiding them: a test
   that only passes sequentially is reported as a defect to fix, never silenced
   by reverting it to sequential without a stated reason.
3. Task Verification stops carrying the whole-package suite. A Task proves its
   own effect with focused checks; the Run-level gate proves nothing else
   regressed, which is what it exists for.
4. The authoring rule for Task Verification is recorded, so the next fourteen
   Tasks do not reintroduce the same per-Task tax.
5. A measured suite-time budget is asserted, so a change that makes
   verification materially slower fails rather than accumulating.
6. Coverage equivalence is proven, not assumed: the set of test functions
   executed before and after is identical.
7. The Spec closes with a published before-and-after comparison: the same
   commands on the same machine, the headline and per-package tables side by
   side, and the delta stated. The baseline recorded under `baseline/` before
   any change is the "before"; it is not re-derived afterwards.

## Non-Goals / Out of Scope

- Deleting tests, skipping tests, or weakening any assertion.
- Reducing what is covered in exchange for speed.
- Changing the QA gate's discovery order or detector placement, owned by
  Spec 0063.
- Where the QA gate lives in the Task Graph, owned by Spec 0072.
- Rewriting the test suite's structure beyond what parallel execution requires.

## Success Metrics

- `go test ./... -count=1` completes materially faster than the recorded 136.9s
  baseline on the same machine, and `internal/cli` alone faster than 113.2s —
  the package that sets the floor.
- A published before-and-after comparison accompanies the Spec at close,
  measured with the baseline's own commands.
- The set of test functions that execute is identical before and after,
  compared mechanically rather than by inspection.
- No Task Verification in any active Spec carries a whole-package suite command.
- A deliberately introduced slow test trips the suite-time budget.
- Every test left sequential carries a stated reason.

## Decisions

- Coverage is the asset and is not negotiable; only execution changes.
- A test that fails under parallel execution has found a real shared-state
  defect. Fixing it is the work; reverting it to sequential without a reason is
  not.
- The Run-level gate is where "nothing else regressed" belongs. Asking every
  Task to prove it costs the same answer fourteen times.
- How often the gate closes a Spec is a graph-shape question, owned by Spec
  0072. This Spec owns what each closing costs.
- This Spec evolves verification cost and never regresses coverage: any change
  that reduces what is exercised is a defect, not a saving.

## Open Questions

None.
