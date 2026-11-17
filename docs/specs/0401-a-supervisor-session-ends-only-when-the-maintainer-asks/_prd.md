---
spec: 0401-a-supervisor-session-ends-only-when-the-maintainer-asks
status: active
created: 2026-11-17
surfaces: [backend, cli, docs]
type: fix
---

# A Supervisor session ends only when the maintainer asks

Spec 0397 put a Phase Boundary after every Authoring Phase and after a Spec's delivery and release: the four authoring skills, the Roundfix skill's `--follow-up` and `--unattended` modes and the `autonomous-work` Baseline clause tell the Supervisor to write the Session Handoff and end the session there, and `roundfix phase start` prints the Boundary Notice whenever the same session records a different phase or Spec. The maintainer found first that the stop after each phase slowed the work, and then that the stop at the end of each Spec also harms autonomous work, because an unattended session halts and asks to be restarted. He decided that no phase of a Spec ends the session. After this Spec the Supervisor follows its queue from phase to phase and from Spec to Spec in one session, and writes the Session Handoff and ends only when the maintainer asks or the work stops; `roundfix phase start` keeps recording every phase and prints no notice.

## Sources

One Backlog Entry, adopted by this Spec's planning change (see `references/_index.md`):

1. [`references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md`](references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md), recorded by the operator from the maintainer's message and answer of 2026-11-16. Whole; its decision is superseded by the answer of 2026-11-17 below.

It is not a Notice Request. The maintainer's answers of 2026-11-16 and 2026-11-17 are in [`references/2026-11-16-maintainer-answers.md`](references/2026-11-16-maintainer-answers.md). The measurement of the operator's Phase Record is in [`references/2026-11-17-phase-record-boundary-measurement.md`](references/2026-11-17-phase-record-boundary-measurement.md).

The only other open source that touches authoring sessions is `docs/backlog/2026-11-16-a-new-session-and-each-authoring-phase-start-with-a-large-context.md`. It stays apart: it measures and trims what a session and an Author Subagent read when they start, waits on a week of per-phase numbers, and changes no stop. The planning change revises it in place to say that a Supervisor session now starts only when the previous one ended at the maintainer's request or because the work stopped.

## Prerequisites

None in `requires`. Authored against `bcd716df` (`origin/main` at Roundfix 0.147.1, with Spec 0399 delivered and the plan of Spec 0400). This planning change rewrites the plan of this Spec, committed on `origin/main` before any Run, and supersedes ADR-0413, which that plan added and nothing implemented. Spec 0400 also raises the Roundfix skill's version. The files the two Specs share are derived records (`skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json` and the `SKILL.md` version fields), each with a derived-path declaration in `.roundfixrc.yml`. Whichever Spec delivers later records its versions on its merge base, as the record commands do, and the Delivery Queue resolves a conflict confined to those paths by regeneration.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the record keeps its schema and its keys. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; `roundfix phase start` now writes one local day file and reads none. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0414 (new) records the decision, amends ADR-0407 and supersedes ADR-0413; the design keeps ADR-0401 and ADR-0408 (where an Authoring Phase runs), ADR-0410 (the status view reads phase entries, never the `boundary` key), ADR-0186 and ADR-0257 (a reworded Baseline clause keeps its identity and its adopter-neutral words); ADR-0258, ADR-0342, ADR-0353 and ADR-0405 are listed only to record why they do not apply; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-16, unchanged by the answer of 2026-11-17, quoted in the authorization record. It covers the `autonomous-work` Baseline module, its derived Standard TypeScript golden and profile pin, this repository's managed `autonomous-work` guide and Setup Manifest, and the owned skills this Spec changes. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0401-a-supervisor-session-ends-only-when-the-maintainer-asks/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`, `internal/baseline/assets/modules/autonomous-work.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Goals

- No phase of a Spec ends the Supervisor session, with or without a person present: after an Authoring Phase, a delivery or a release the session goes on with the next queued work, and no skill, mode or guide asks for a Session Handoff there.
- The session writes the Session Handoff and ends only when the maintainer asks or the work stops because nothing remains or what remains waits on the maintainer; a turn that ends while work goes on still arms something that outlives it.
- `roundfix phase start` keeps recording every phase, prints nothing on standard error at any phase change, writes `"boundary":false` and no `previous_phase` or `previous_spec`, and changes nothing else in its output, its exit codes or the record's schema.
- Every place that says when a session ends says the same thing.

## Core Features

1. The Phase Record's writer reads no earlier entry and marks none, and `roundfix phase start` prints no Boundary Notice; the reader still reads an entry an earlier version marked.
2. The `autonomous-work` Baseline clause states that no phase of a Spec ends the session and when it ends, keeping its identity, its force and its rule to arm something that outlives the turn; its derived files and this repository's guide are regenerated.
3. The `Phase boundary` sections of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`, and step 8 of `--follow-up` and the `--unattended` section of the Roundfix skill's modes, state the same rule.
4. The Roundfix skill's phase reference, the `phase` command page and the Phase Record section of the usage page describe the record without a boundary; the glossary redefines Phase Boundary and Session Handoff and drops Boundary Notice, and ADR-0407 is amended by ADR-0414.

## User Experience

Today, in one session:

```text
$ roundfix phase start release --spec 0412-example
Recorded release of 0412-example for claude session <session>
$ roundfix phase start follow-up
Recorded follow-up for claude session <session>
roundfix: notice: this session already recorded release of 0412-example; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.
```

With this Spec the second command prints only its result line, and so does every phase after any other:

```text
$ roundfix phase start release --spec 0412-example
Recorded release of 0412-example for claude session <session>
$ roundfix phase start follow-up
Recorded follow-up for claude session <session>
```

A Supervisor session that delivered and released a Spec reads the next queued work and goes on; it writes the Session Handoff only when the maintainer says to stop or nothing it can do remains.

## Non-Goals / Out of Scope

- Where an Authoring Phase runs: ADR-0401 and ADR-0408 stay as they are, so a feature's phases go to their Author Subagents and a fix or refactor is authored in the session unless the user asks.
- The Phase Record's schema, keys, file layout, retention and the prompt entries.
- Rewriting the `boundary` mark of entries already written.
- Any limit on how long a session runs or how much context it holds; the harness's compaction and the maintainer's request stay the only ends besides the work stopping.
- The maintainer's own instruction files outside the repository, which the operator adjusts.
- What a session or an Author Subagent reads when it starts, which the large-context Backlog Entry keeps.
- Any Source Baseline, Setup Snapshot, other Baseline module or other profile.

## Success Metrics

1. Success Metric: after any earlier phase entry of the same session — each of `follow-up`, `idea`, `prd`, `techspec`, `tasks`, `delivery` and `release`, of the same or another Spec, and an entry an earlier version marked — `roundfix phase start` exits `0`, prints only its result line, prints nothing on standard error and writes a line whose keys are exactly the nine of a phase entry with `"boundary":false`.
2. Success Metric: no authoring skill, Roundfix skill mode, Baseline clause, guide, glossary entry or user-guide page tells the Supervisor to end the session or write the Session Handoff after an Authoring Phase or after a Spec's delivery and release, and each of them says the session ends only when the maintainer asks or the work stops.
3. Success Metric: the reworded `autonomous-work` clause keeps its id, its `mandatory` force, no `replaces` and its sentence on arming something that outlives the turn; a second `roundfix baseline update` of this repository reports `current`.
4. Success Metric, outside evidence: replayed in order through the built binary in a temporary Roundfix Home, the eight phase entries of the operator's Phase Record measured on 2026-11-17, followed by a `release` and a `follow-up`, print no Boundary Notice and write no mark, where Roundfix 0.147.0 printed four notices for the eight.

## Acceptance evidence

The outside-evidence row rests on the operator's own Phase Record of 2026-11-16 to 2026-11-17, written by Roundfix 0.147.0 while the Supervisor planned Specs 0399, 0400 and this one, which this Spec did not design. The measurement in `references/` copies each entry's phase, Spec number and mark, without the session identifier. The QA gate replays them read from that reference, never from the operator's home directory.

## Unreachable Acceptance

None. The measurement is committed with the plan, and the replay runs in a temporary directory.

## Validation

The decisions were checked on 2026-11-17 against these sources:

- The measurement in `references/`.
- The code of `internal/phaserecord/record.go`, `internal/phaserecord/read.go` and `internal/cli/phase.go`, and the tests `internal/phaserecord/record_test.go`, `internal/phaserecord/read_test.go`, `internal/cli/phase_test.go`, `internal/cli/follow_up_queue_test.go`, `skills/phase_boundary_test.go`, `skills/phase_reference_test.go` and `internal/baseline/phase_boundary_clause_test.go`.
- ADR-0407, ADR-0413, ADR-0401, ADR-0408 and ADR-0410.
- A scratch worktree of `bcd716df` with the texts reworded and the writer changed, where the regeneration commands changed only the files `_tasks.md` names and the current tests named in the TechSpec failed.

## Skills

- unchanged: guide docs/user-guide/usage.md — only its Phase Record section changes, and the Roundfix skill's phase reference, which task 01 rewrites, covers that section; none of the guide's other covering references describes the Phase Record.

## Glossary

- changes: **Phase Boundary**
- changes: **Session Handoff**

## Decisions

- 2026-11-16: The maintainer's message — "Esse negocio de criar sessões a cada etapa está atrasando demais."
- 2026-11-16: Where does the Phase Boundary go? — "Só por Spec (Recommended)" (~21:00 BRT). Superseded by the answer of 2026-11-17.
- 2026-11-16: Which authorization applies to Spec 0401? — "Mesmas do BE (Recommended)" (~21:45 BRT): implement, commit, push, pull request, merge and release right after the Spec, up to three corrections. Makefile, CI, `go.mod` and lint only if a correction requires one.
- 2026-11-17: When does a Supervisor session end? — "Só quando você pedir (Recommended)" (~10:00 BRT), after the maintainer noticed that a session stopped at the end of a Spec's delivery asking to be restarted, which "vai prejudicar demais o autonomous work": no automatic stop, neither after a phase nor after a Spec, with or without a person present; the Session Handoff is written only when the session really ends; the Phase Record keeps measuring; the notice of `roundfix phase` goes. The authorization stays "Mesmas do BE".

Autonomous route: no person present during the amendment; the remaining decisions were taken from the operator's brief of 2026-11-17 and the measurement.

- `roundfix phase start` reads no earlier entry, writes `"boundary":false` and never `previous_phase` or `previous_spec`; the reader keeps reading old marked entries (ADR-0414).
- "The work stops" means nothing queued remains or what remains waits on the maintainer; the rule to arm something that outlives the turn stays.
- The term Phase Boundary stays, redefined as the change from one phase to the next that ends no session, because the four skills' `Phase boundary` sections still record the phase there; task 01 removes the glossary entry of the notice `roundfix phase start` no longer prints, because the term names nothing after this Spec.
- ADR-0413 is superseded by ADR-0414 and deleted in this planning change, as the retirement clause of `docs/agents/docs-layout.md` asks, because its decision is reversed and nothing implemented it.
- `type: fix` with a `patch` Release Ceiling: no key, field, command, flag or exit code is removed or added. The removed line was an advisory diagnostic on standard error that never changed the result, and `docs/agents/cli.md` makes stderr placement, not each diagnostic, public.

## Open Questions

None.
