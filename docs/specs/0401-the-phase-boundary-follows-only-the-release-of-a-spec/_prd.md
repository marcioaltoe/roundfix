---
spec: 0401-the-phase-boundary-follows-only-the-release-of-a-spec
status: active
created: 2026-11-17
surfaces: [backend, cli, docs]
type: fix
---

# The Phase Boundary follows only the release of a Spec

Spec 0397 put a Phase Boundary after every Authoring Phase: the four authoring skills, the Roundfix skill's `--follow-up` mode and the `autonomous-work` Baseline clause tell the Supervisor to write the Session Handoff and end the session after the Idea, PRD, TechSpec and Tasks phases, and `roundfix phase start` prints the Boundary Notice whenever the same session records a different phase or Spec. The maintainer found that it slows the work down, and decided that the session ends only after the Spec is delivered and released, with the Authoring Phases in agents that have their own context. After this Spec the Supervisor goes on from phase to phase in one session, every guide and skill states the boundary once, after delivery and release, and the notice prints only when a session goes on after a release.

## Sources

One Backlog Entry, adopted by this Spec's planning change (see `references/_index.md`):

1. [`references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md`](references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md), recorded by the operator from the maintainer's message and answer of 2026-11-16. Whole.

It is not a Notice Request. The maintainer's answers are in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md). The measurement of the operator's Phase Record is in [`references/2026-11-17-phase-record-boundary-measurement.md`](references/2026-11-17-phase-record-boundary-measurement.md).

The only other open source that touches authoring sessions is `docs/backlog/2026-11-16-a-new-session-and-each-authoring-phase-start-with-a-large-context.md`. It stays apart: it measures and trims what a session and an Author Subagent read when they start, waits on a week of per-phase numbers, and changes no boundary. The planning change revises it in place to say that a session now starts once per Spec.

## Prerequisites

None in `requires`. Authored against `379803c8` (`docs/plan-0401`, based on `origin/main` with the plans of Specs 0399 and 0400). Spec 0399 also edits Baseline modules and raises the Roundfix skill's version, and Spec 0400 raises the Roundfix skill's version. The files they share with this Spec are derived records (`internal/baseline/module-versions.json`, the catalog snapshots and plan goldens, the Standard TypeScript profile pin, `docs/agents/setup-context.json`, `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json` and the `SKILL.md` version fields), each with a derived-path declaration in `.roundfixrc.yml`. Whichever Spec delivers later records its versions on its merge base, as the record commands do, and the Delivery Queue resolves a conflict confined to those paths by regeneration.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the record keeps its schema and fields. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; `roundfix phase start` reads local day files only. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0413 (new) records the decision and amends ADR-0407; the design keeps ADR-0401 and ADR-0408 (where an Authoring Phase runs), ADR-0410 (the status view reads phase entries, not the boundary mark), ADR-0186 and ADR-0257 (a reworded Baseline clause keeps its identity and its adopter-neutral words); ADR-0258, ADR-0342, ADR-0353 and ADR-0405 are listed only to record why they do not apply; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-16, quoted in the authorization record. It covers the `autonomous-work` Baseline module, its derived Standard TypeScript golden and profile pin, this repository's managed `autonomous-work` guide and Setup Manifest, and the owned skills this Spec changes. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0401-the-phase-boundary-follows-only-the-release-of-a-spec/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`, `internal/baseline/assets/modules/autonomous-work.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- An Authoring Phase ends at no Phase Boundary: the session goes on to the next phase, and no skill, mode or guide asks for a Session Handoff after it.
- The Phase Boundary of a Spec follows its delivery and release, or its merge when it releases nothing, and keeps what ADR-0407 asks there: the Session Handoff, what stays armed, and no further phase.
- `roundfix phase start` prints the Boundary Notice only when the same session's latest phase entry is `release`, and changes nothing else in its output, its exit codes or the record's format.
- Every place that states the boundary states the same one.

## Core Features

1. The Phase Record's boundary rule marks an entry only after a `release` of the same session.
2. The `autonomous-work` Baseline clause states the new boundary, keeping its identity and force, and its derived files and this repository's guide are regenerated.
3. The `Phase boundary` sections of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`, and the `--follow-up` step 8 and `--unattended` section of the Roundfix skill's modes, state the new boundary.
4. The Roundfix skill's phase reference, the `phase` command page and the Phase Record section of the usage page describe the new notice rule; the glossary's Phase Boundary and Boundary Notice say the same, and ADR-0407 is amended by ADR-0413.

## User Experience

Today, in one session:

```text
$ roundfix phase start prd --spec 0412-example
Recorded prd of 0412-example for claude session <session>
$ roundfix phase start techspec --spec 0412-example
Recorded techspec of 0412-example for claude session <session>
roundfix: notice: this session already recorded prd of 0412-example; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.
```

With this Spec the second command prints no notice. Only a session that goes on after a release sees it:

```text
$ roundfix phase start release --spec 0412-example
Recorded release of 0412-example for claude session <session>
$ roundfix phase start follow-up
Recorded follow-up for claude session <session>
roundfix: notice: this session already recorded release of 0412-example; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.
```

## Non-Goals / Out of Scope

- Where an Authoring Phase runs: ADR-0401 and ADR-0408 stay as they are, so a feature's phases go to their Author Subagents and a fix or refactor is authored in the session unless the user asks.
- The Phase Record's schema, fields, file layout, retention and the prompt entries.
- The text, stream and exit behaviour of the Boundary Notice.
- Rewriting the `boundary` mark of entries already written.
- What a session or an Author Subagent reads when it starts, which the large-context Backlog Entry keeps.
- Any Source Baseline, Setup Snapshot, other Baseline module or other profile.

## Success Metrics

1. Success Metric: after an earlier `follow-up`, `idea`, `prd`, `techspec`, `tasks` or `delivery` entry of the same session, of the same or another Spec, `roundfix phase start` writes `"boundary":false` and prints nothing on standard error; after an earlier `release`, any other phase, or a `release` of another Spec, writes `"boundary":true` with `previous_phase` `release` and prints the Boundary Notice; a `release` of the same Spec again, or again without a Spec, prints none.
2. Success Metric: no authoring skill, Roundfix skill mode, Baseline clause, guide, glossary entry or user-guide page says that a Phase Boundary follows each Authoring Phase, and each of them states the boundary after a Spec's delivery and release.
3. Success Metric: the reworded `autonomous-work` clause keeps its id, its `mandatory` force and no `replaces`; a second `roundfix baseline update` of this repository reports `current`.
4. Success Metric, outside evidence: replayed in order through the built binary in a temporary Roundfix Home, the eight phase entries of the operator's Phase Record measured on 2026-11-17 print no Boundary Notice, where Roundfix 0.147.0 printed four; a `release` followed by a `follow-up` in the same replay prints one.

## Acceptance evidence

The outside-evidence row rests on the operator's own Phase Record of 2026-11-16 to 2026-11-17, written by Roundfix 0.147.0 while the Supervisor planned Specs 0399, 0400 and this one, which this Spec did not design. The measurement in `references/` copies each entry's phase, Spec number and mark, without the session identifier. The QA gate replays them read from that reference, never from the operator's home directory.

## Unreachable Acceptance

None. The measurement is committed with the plan, and the replay runs in a temporary directory.

## Validation

The decisions were checked on 2026-11-17 against these sources:

- The measurement in `references/`.
- The code of `internal/phaserecord/record.go` and `internal/cli/phase.go`, and the tests `internal/phaserecord/record_test.go`, `internal/cli/phase_test.go`, `skills/phase_boundary_test.go`, `skills/phase_reference_test.go` and `internal/baseline/phase_boundary_clause_test.go`.
- ADR-0407, ADR-0401, ADR-0408 and ADR-0410.

## Skills

- unchanged: guide docs/user-guide/usage.md — only its Phase Record section changes, and the Roundfix skill's phase reference, which task 01 rewrites, covers that section; none of the guide's other covering references describes the Phase Record.

## Glossary

- changes: **Phase Boundary**
- changes: **Boundary Notice**

## Decisions

- 2026-11-16: The maintainer's message — "Esse negocio de criar sessões a cada etapa está atrasando demais."
- 2026-11-16: Where does the Phase Boundary go? — "Só por Spec (Recommended)" (~21:00 BRT): the session ends only after the Spec is delivered and released; the Authoring Phases run in agents; a fix Spec adjusts the skills and the notice of `roundfix phase`.
- 2026-11-16: Which authorization applies to Spec 0401? — "Mesmas do BE (Recommended)" (~21:45 BRT): implement, commit, push, pull request, merge and release right after the Spec, up to three corrections. Makefile, CI, `go.mod` and lint only if a correction requires one.

Autonomous route: no person present during authoring; the remaining decisions were taken from the operator's brief for Spec 0401 of 2026-11-16 and the measurement.

- The notice prints only after a `release` of the same session, and keeps its text (ADR-0413).
- A Spec that releases nothing reaches its boundary at its merge, which the record does not see, so that boundary has no notice.
- The fix route keeps authoring in the session (ADR-0408); the answer's "fases de autoria em agentes" is read as the Author Subagents ADR-0401 already generates.
- `type: fix` with a `patch` Release Ceiling: no public output is added, and the notice prints in fewer cases.

## Open Questions

None.
