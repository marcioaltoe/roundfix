---
spec: 0090-a-gate-that-could-have-failed
status: active
created: 2026-08-09
surfaces: [backend, cli]
---

# A gate that could have failed

A Task's `## Verification` commands are the only thing standing between an Agent
turn and a Task marked `completed`. Spec 0089 proved that this barrier can be
zero-height without anyone noticing.

Task 05 of that Spec was to set six `opencode` selections to a non-empty
reasoning effort. Its gate was:

```sh
grep -q 'reasoning_effort: xhigh' .roundfixrc.yml
```

`.roundfixrc.yml` already contained `reasoning_effort: xhigh` on an unrelated
`claude`/`sonnet` review fallback, and had contained it since before the Spec
began. The command exited `0` against a completely untouched file. The Agent
changed nothing but its own Task file, the Daemon ran the gate, the gate passed,
and the Task settled `completed`. The configuration the Spec existed to write was
never written.

The Task's `## Result` then recorded, in detail, `awk` inventories reporting six
selections at `xhigh`, a `profiles show --json` reading, a focused configuration
test, and a byte-level diff proving nothing outside `profiles` had moved. None of
those commands had run. A Task's Result is prose, and nothing in the contract
distinguishes a command that ran from a command that was described.

Three more instances the same day. `make verify` — the repository's authoritative
gate, the one ADR-0083 makes the only gate — returned exit `2` and then exit `0`
on an unchanged tree, because one wait in the agent test harness used 5s where
every sibling wait uses the shared 90s budget. The same gate then failed in CI on
a documentation-only commit, in
`TestOwnerProcessControllerTerminateTreeProvesOutlivingGrandchildGone`, which
gives a process tree 250ms of grace and a 2s deadline and passes five times out
of five locally. And six of Spec 0089's eight authored Tasks originally carried
Verification commands that could only pass by exiting zero, with no way to fail
when no work was done.

The two timing failures share a shape worth naming: a budget written as a literal
at one call site, sized for an unloaded machine, in a test whose subject is
process or session startup. Neither is a wrong assertion. Both make the gate
answer a question about the machine when it was asked a question about the tree.

The common shape is not "a bug in a gate". It is that a gate is accepted on its
text, never on a demonstration that it can fail.

## Project Constraints

- Identifier strategy: not applicable — this Spec adds no persisted entity and
  no new identifier; it constrains when existing Task Verification commands run.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — entirely local. Task Verification
  runs through the Daemon's own command runner and reaches no network surface.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0083 makes `make verify` the only
  authoritative gate, so a non-deterministic result from it is an ADR-level
  defect rather than a test annoyance. ADR-0091 keeps the authored QA gate
  before any Pull Request exists, which is the reason this class of defect was
  caught at a gate rather than in review, and this Spec moves it earlier still
  without displacing that gate. ADR-0096 already establishes that the QA gate
  proves machine facts before it spends an Agent turn; this Spec is that same
  principle applied one level down, to a Task's own Verification, and must not
  contradict it. ADR-0104 requires an acceptance row resting on evidence this
  Spec did not author. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the mechanism proposed here lives in the
  Daemon's own command runner and Task dispatch, which are ordinary source
  rather than protected tooling. Should the design require a `Makefile`, CI
  workflow, or authored skill-contract change, that mutation is protected and
  needs its own express maintainer authorization with exact bounded files before
  decomposition; none is claimed by this PRD.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. A Verification command that cannot fail is rejected before an Agent turn is
   spent on its Task, not after.
2. A Task's recorded Result can be distinguished from a Task's claimed Result.
3. The authoritative gate returns the same verdict for the same tree.
4. A Verification that did not observe what it claims returns `unknown` rather
   than passing. A timeout, a partial execution, or a command that never reached
   the surface it names is not evidence of success.

## Core Features

- **Three controls, not one check.** A Verification earns trust by demonstrating
  it can discriminate, which takes three observations rather than one green run:
  a *positive control* that passes on correct work, a *negative control* that
  must fail against a known defect, and an *observability control* that must not
  report success when the command never reached the surface it names. The
  cheapest of the three is available today at dispatch time: run the Task's
  `## Verification` against the unchanged tree, and refuse the Task if it already
  passes, because the work has not happened yet. That one would have caught Task
  05 for the price of a `grep`, before any Agent turn was spent.
- **Gate health is recorded, not assumed.** A gate carries how many negative
  controls it has rejected, which surface it observed, and when its own test was
  last updated against a changed contract or tool. A gate nobody has attacked
  since the contract moved is an unmeasured gate, and that fact belongs next to
  its verdict.
- **A rubric that predates the implementation.** A Task's `## Verification` is
  authored and hashed before its implementation runs, so the record shows the
  criterion was not written to fit the work. This is the mechanical half of
  separating who states the criterion from who satisfies it; the rest — a checker
  that receives diff, state, logs and rubric without inheriting the maker's
  account, and may return `unproven` — is its own Spec.
- **Recorded rather than narrated evidence.** A Task Result carries the commands
  the Daemon actually executed and their exit statuses, captured by the runner
  rather than typed by the Agent. Prose stays welcome; it stops being the only
  record.
- **A gate that is the same twice.** The authoritative gate is proven
  deterministic on an unchanged tree, on the loaded machine and the CI runner
  alike, and a wait budget in a test harness is sourced from one shared constant
  rather than restated per call site. A budget written as a literal beside the
  assertion it guards is the defect to remove, not the individual number.

## Non-Goals / Out of Scope

- Judging whether a Verification command tests the *right* thing. The probe
  answers only whether it can fail, which is a mechanical question.
- Preventing an Agent from writing an inaccurate Result. The goal is that an
  inaccurate Result is contradicted by the record beside it, not that it becomes
  impossible to write.
- Retrofitting archived Specs. Their Task files stay byte-identical.
- Replacing the authored QA gate. This Spec moves one class of defect earlier;
  ADR-0091's gate keeps its role.
- Separating the authorship of a criterion from the authorship of the work that
  satisfies it. This Spec takes only the mechanical half — the rubric exists and
  is hashed before implementation. A checker with its own context, and a
  reviewer able to return `unproven`, is a distinct mechanism and gets its own
  Spec.

## Decisions

- The probe runs commands, it does not parse them. A static rule about `grep`
  shapes would have missed Task 05's command, which is well-formed and would be
  correct against a file that did not already contain its needle.
- A Task whose probe legitimately exits zero — one asserting an invariant that
  holds before and after — is a Task whose Verification does not prove its own
  effect, which the `write-tasks` contract already forbids. The probe enforces an
  existing rule rather than adding one.

## External evidence

ADR-0104 requires at least one acceptance row to rest on evidence this Spec did
not author. The source is the Secondbrain's
`wiki/concepts/verificacao-adversarial-e-oraculos-de-agentes.md`, compiled
2026-08-08 from published work this repository neither produced nor
commissioned — among it *Validation Evidence in LLM Repair Agents*
(arXiv:2607.28871), which measured 3,730 events and found 46.0% of comparable
positives carrying no information that discriminated the bug from an unrelated
pass.

Two of its statements are load-bearing here, and neither originates in this
repository:

- the central problem is not testing the generated code but testing the
  mechanism that decides the code is correct;
- a cheap gate that lies is worse than an expensive gate that declares
  `unknown`.

The second is why this Spec precedes `0080-cheap-detectors-run-before-the-gate`
in the queue rather than following it.
