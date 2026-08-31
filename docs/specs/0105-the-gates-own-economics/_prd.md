---
spec: 0105-the-gates-own-economics
status: active
created: 2026-08-12
surfaces: [backend, docs]
---

# The gate's own economics

Of 201 failed Tasks measured across five repositories, 123 are the QA gate
returning a verdict rather than code breaking. The gate is not too strict — it
finds defects no suite catches — but its expensive tail has one cause: the Spec
assumes the world behaves like its fakes, so a design premise survives authoring,
implementation and every unit test and dies at the gate several Runs later. Two
structural costs compound it. The Pull Request row is unreachable by design in
every Spec, because the authored gate runs before a Pull Request exists, and one
Spec paid six of its eight gate executions for that row alone. And the gate's own
Verification is authored by hand over a verdict the Daemon already derived, which
in one measured case produced a gate that passed itself having failed.

## Project Constraints

- Identifier strategy: applicable — QA Report, verdict, blocked-row causes,
  Unreachable Acceptance, and Characterization are glossary terms this Spec
  changes the obligations of. The closing node checks whether the work introduced
  or changed a term the glossary should carry. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read by this Spec. Rows that exercise a hosting provider
  do so through the existing command surface. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0080 makes QA verdicts distinguish
  environment-blocked rows, which is the mechanism the Pull Request row must use
  rather than each Spec rediscovering it; ADR-0091 makes the gate a Task node of
  its own type; ADR-0088 authors the gate into the graph rather than requesting it
  per run, which is why the Pull Request row is unreachable by construction;
  ADR-0097 lets a row carry forward only on declared, unmoved evidence, which
  bounds what an equivalent-evidence path may accept; ADR-0096 proves machine
  facts before spending an Agent turn and ADR-0117 places a check with the stage
  that can produce its defect, which together are this Spec's cost thesis;
  ADR-0093 checks consistency by citation rather than inference; ADR-0104 makes a
  Spec accept on evidence it did not author, which is the rule the
  characterization change serves. ADR-0081 makes the generated skill copies
  deterministic fallout of an authorized edit, which is why the three skills and
  their mirrors are one authorization rather than six. ADR-0149 would normally
  spare a grant from enumerating a command's outputs by resolving them from an
  `_ownership.yml` declaration; `skills/` carries none, and the resolver reads
  only under `internal/baseline`, so this Spec's grant names the generated
  copies explicitly. That gap is recorded for Triage and is not this Spec's to
  close. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the QA gate skill, the task-authoring skill, and
  the task-execution skill are edited. Express maintainer authorization:
  `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
  granted 2026-08-12. Bounded files: `.agents/skills/qa-gate/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/implement-task/SKILL.md`.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. A Spec crossing an external boundary meets that boundary before its gate.
2. The Pull Request row costs one Spec's discovery, not every Spec's.
3. The gate's own Verification stops being hand-authored over a derived verdict.
4. A governance failure stops blinding the rest of the matrix.

## User Stories

1. As a Supervisor authoring a Spec that crosses an adapter, a contract, or a
   database, I want its characterization to record what the real thing does, so
   that a false premise dies on day one rather than after four Runs.
2. As a Supervisor whose gate runs before a Pull Request exists, I want the
   equivalent-evidence path applied to that row by default, so that every Spec
   does not pay to rediscover it.
3. As a Supervisor whose gate failed on governance, I want the flow rows reported
   too, so that a round produces signal about function rather than only a block.
4. As a maintainer, I want the gate's Verification owned by the tool that derives
   the verdict, so that a hand-written predicate cannot approve a failure.
5. As a Supervisor authoring in this repository's languages, I want the citation
   parser to accept the forms I actually write, so that a correction round is not
   spent on punctuation.

## Core Features

1. **Characterization touches the real boundary.** A Spec crossing an external
   surface records what the real thing does rather than what a fake does, so a
   premise the boundary does not support fails at authoring.
2. **The Pull Request row carries its own path.** The gate applies the
   equivalent-evidence path to that row by default, requiring the evidence to be
   recorded, so the row is neither a per-Spec discovery nor a silent pass.
3. **Roundfix owns the QA Task's Verification.** The gate Task's Verification is
   derived rather than authored, so an author cannot write a predicate that
   accepts a verdict outside the domain or selects an older report.
4. **Independent static findings report together.** A governance failure does not
   stop the matrix before the flow rows run, so a round reports both.
5. **The citation parser reads the forms Specs are written in.** A conjunction as
   a list separator and a decision number without its prefix are recognised, or
   the failure message names the form that is.

## User Experience

A gate round reports governance findings and flow findings in one pass. The Pull
Request row reports the equivalent evidence it accepted — the head, its ancestry,
the changed files, the local checks — rather than passing silently or blocking.

## Non-Goals / Out of Scope

- Loosening the gate. Eleven non-passing verdicts in one measured session all
  failed on contract rather than business logic, and the same gate found four real
  defects no suite would catch. Any change must be read against both numbers.
- Moving the gate after the Pull Request exists; that is settled by an accepted
  decision this Spec does not reopen.
- Changing verdict semantics or the typed blocked-cause counts.
- Building prepared data infrastructure for data-shaped repositories, which is
  infrastructure rather than method.
- The review-resolution Agent's guardrails: not rewriting an authored
  Verification, and not resolving a finding that needs absent infrastructure.
  Split out on 2026-08-31 under this Spec's own Open Question. The split is by
  subsystem: those two live in `internal/rounds`, which no other feature here
  touches and which carries no reference to Verification today, so the guardrail
  is new behaviour in a subsystem rather than an adjustment to this one. Carried
  as `docs/backlog/2026-08-31-the-review-agent-rewrites-the-contract-it-was-asked-to-satisfy.md`,
  with its undecided forbid-versus-report question travelling with it.
- Measuring the gate's own cost distribution, which this Spec's changes make
  measurable but does not itself deliver.

## Success Metrics

- A Spec whose premise the real boundary does not support fails at its
  characterization Task rather than at the gate. Source: a Spec measured on
  2026-08-10 whose premise survived four Runs before the gate refuted it.
- The Pull Request row consumes no gate rerun across the next Specs that carry it,
  against a measured baseline of six of eight executions in one Spec in a
  repository this Spec did not build.
- A hand-authored gate predicate can no longer accept a verdict outside the
  domain, proven by exercising the measured case that did.
- A gate round with a governance failure still reports flow rows, against a
  measured round where a single governance finding blocked fifteen of nineteen.
- The two measured citation forms parse without a correction round.

## Decisions

- The equivalent-evidence path becomes the gate's default for the Pull Request row
  rather than a template clause, because a template is per-Spec discovery again.

## Open Questions

- Resolved 2026-08-31: this Spec splits, and the review-Agent guardrails travel
  separately in the backlog entry named under Non-Goals. The undecided
  forbid-versus-report question travels with them rather than being settled
  here, because settling it would commit this Spec to a subsystem it does not
  otherwise touch.
- None outstanding.
