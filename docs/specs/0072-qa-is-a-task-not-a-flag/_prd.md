---
spec: 0072-qa-is-a-task-not-a-flag
status: active
created: 2026-08-02
surfaces: [backend, cli, docs]
---

# QA is a Task, not a flag

The QA gate is requested with `--qa` on the Implement Command. The Daemon
withholds it correctly — no gate begins while any Task is unsettled, and Spec
0057's first Run proves it by ending with one Task failed and no gate at all.
Nothing runs early.

The gap is the other direction. Because the gate belongs to the invocation
rather than to the Task Graph, a graph that grows *after* a gate has reported
leaves no structural trace. On Spec 0057 a corrective Task was appended after
each gate, and three gates ran against three different graphs at roughly twenty
to twenty-five minutes each. Read from outside, that is three normal cycles. It
was one decomposition that was wrong twice.

`docs/agents/autonomous-work.md` already warns about the serial chain and
already caps corrective Tasks at two. An advisory cap on a flag is what let
three closings pass unnoticed.

Making QA a Task the Spec authors — a terminal node depending on every leaf —
puts the gate inside the artifact that describes the work, where growing the
graph afterwards is visibly a node inserted before a terminal that already
reported.

## Project Constraints

- Identifier strategy: not applicable — Task identifiers, QA report names, and
  verdict values keep their existing contracts; no project-owned Internal
  Identifier is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the governing clause prohibits
  reading, printing, committing, or generating secrets and forbids inventing
  authentication, authorization, transport, or deployment policy; this Spec
  changes graph shape and a command surface. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0080 owns verdict semantics and the
  typed blocked-row counts, which a QA Task must produce unchanged; ADR-0088
  records removing the flag in favour of the authored node. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: on
  2026-08-02 the maintainer authorized tooling adjustment naming the `Makefile`
  and the owned skills, recorded at
  `docs/workflow/authorizations/2026-08-02-queued-spec-tooling.md`; bounded
  files: the owned skill pair, since `write-tasks` and the QA gate skill both
  change. Deterministic digest fallout is sanctioned by ADR-0081. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- A Spec that wants a QA gate declares it when its artifacts are authored, and
  the gate appears in the Task Graph as its terminal node.
- The gate cannot begin before the last Task settles, and cannot be requested
  apart from the graph it closes.
- Appending work after the gate has reported invalidates that result instead of
  producing a second cycle beside it.
- A Spec that legitimately needs no QA gate says so once, rather than by
  omitting a flag at every invocation.

## Core Features

1. Authoring decides the gate: when a Spec's artifacts are written, a QA gate
   is either included or explicitly declined with a reason, and the decision is
   recorded in the Spec rather than re-made at each invocation.
2. An included gate is emitted into the Task Graph as a terminal node that
   depends on every leaf Task, so the graph itself encodes that it runs last.
3. The Implement Command stops taking a QA parameter. What runs is what the
   graph says.
4. A Task added after the terminal node has reported invalidates that gate's
   result, and the graph shows the insertion rather than accumulating a second
   independent cycle.
5. The gate produces exactly the report, verdict, and typed blocked-row counts
   it produces today; ADR-0080's semantics are untouched.
6. Gate runs per Spec are countable from the graph and its history, so a Spec
   that needed four closings is legible as a decomposition problem.
7. Specs authored before this contract keep working unchanged, and their
   archived artifacts stay byte-identical.

## Non-Goals / Out of Scope

- Changing what the QA gate checks, how its matrix is derived, or what any
  verdict means.
- Changing the Daemon's withholding, which is already correct.
- Making the gate optional in the sense of skippable at run time; declining it
  is an authoring decision with a recorded reason.
- Test execution cost and the per-Task suite tax, owned by Spec 0071.

## Success Metrics

- A Spec authored with a gate produces a graph whose terminal node is the gate,
  depending on every leaf.
- Running the Implement Command with no QA parameter still closes such a Spec
  with a gate; the parameter no longer exists.
- A Spec authored without a gate closes without one and records why.
- Appending a Task to a graph whose gate already reported invalidates that
  result, proven by a fixture.
- Every archived Spec's artifacts remain byte-identical.
- Replaying Spec 0057's history, the three closings are countable as three.

## Decisions

- The gate belongs to the graph, not to the invocation. That is what makes
  "the graph grew after the gate ran" impossible to do quietly.
- Whether a Spec has a gate is an authoring decision, made once with a reason,
  not a flag chosen per run.
- The Daemon's current withholding is correct and stays. This Spec changes
  where the gate lives, not when it is allowed to begin.
- This Spec evolves the graph contract and never regresses it: no Spec authored
  before it changes behavior, and no archived artifact is rewritten.

## Open Questions

None.
