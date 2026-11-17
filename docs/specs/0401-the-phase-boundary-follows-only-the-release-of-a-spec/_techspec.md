---
spec: 0401-the-phase-boundary-follows-only-the-release-of-a-spec
prd: _prd.md
created: 2026-11-17
---

# The Phase Boundary follows only the release of a Spec — Technical Spec

## Executive Summary

The Phase Boundary moves from every Authoring Phase to the delivery and release of a Spec. One condition in the Phase Record's boundary rule changes: an entry is marked, and the Boundary Notice printed, only when the same session's latest phase entry is `release` and the new entry is not a `release` of the same Spec. Every text that states the boundary changes with it: the `autonomous-work` Baseline clause, the `Phase boundary` sections of the four authoring skills, the Roundfix skill's `--follow-up` and `--unattended` modes and its phase reference, the glossary, ADR-0407 and the user guide. The trade-off accepted is that a Spec which releases nothing reaches its boundary at its merge with no notice, because the record holds no merge entry; adding one would need a new phase and a new writer, for a reminder the instruction already gives. The notice keeps its text, so a reader that matches it keeps working.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the record keeps its schema and fields. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; `roundfix phase start` reads local day files only. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0413 (new) records the decision and amends ADR-0407; the design keeps ADR-0401 and ADR-0408 (where an Authoring Phase runs), ADR-0410 (the status view reads phase entries, not the boundary mark), ADR-0186 and ADR-0257 (a reworded Baseline clause keeps its identity and its adopter-neutral words); ADR-0258, ADR-0342, ADR-0353 and ADR-0405 are listed only to record why they do not apply. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-16, quoted in the authorization record. It covers the `autonomous-work` Baseline module, its derived Standard TypeScript golden and profile pin, this repository's managed `autonomous-work` guide and Setup Manifest, and the owned skills this Spec changes. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0401-the-phase-boundary-follows-only-the-release-of-a-spec/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`, `internal/baseline/assets/modules/autonomous-work.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## System Architecture

Three parts change, and none gains a component.

- **The Phase Record's boundary rule** in `internal/phaserecord/record.go`. `Recorder.boundary` already finds the same session's latest phase entry of the last seven days; only the condition that decides a crossing changes. `internal/cli/phase.go` prints the notice from the returned `Boundary` and does not change.
- **The Baseline clause** `clause.autonomous.loop-07-outlive-the-turn` of `internal/baseline/assets/modules/autonomous-work.json`. Its guidance is reworded; the module record command, `make baseline-digests` and the sanctioned `baseline update` regenerate the derived files and this repository's `docs/agents/autonomous-work.md`.
- **The skills and documents** that state the boundary: the `## Phase boundary` sections of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`; step 8 of `--follow-up` and the `--unattended` section of the Roundfix skill's `references/modes.md`; its `references/phase.md`; `CONTEXT.md`; ADR-0407; `docs/user-guide/commands/phase.md` and the Phase Record section of `docs/user-guide/usage.md`. Each owned skill's canonical copy under `.agents/skills` is edited, `make skills-sync` writes its mirror, and the version record command raises its version.

## Implementation Design

### Interfaces

```go
// internal/phaserecord/record.go — unchanged signatures
func (r *Recorder) StartPhase(ctx context.Context, entry PhaseEntry) (Boundary, error)

type Boundary struct {
	Crossed       bool
	PreviousPhase Phase  // always "release" when Crossed
	PreviousSpec  string // the Spec of that release, possibly empty
}
```

Numbered invariants:

1. Let `latest` be the same non-empty session's phase entry with the latest `time` in the current UTC day file and the six before it. The new entry crosses only when `latest` exists, `latest.Phase == "release"`, and not (`entry.Phase == "release"` and `entry.Spec == latest.Spec`).
2. Every other earlier phase (`follow-up`, `idea`, `prd`, `techspec`, `tasks`, `delivery`), of the same Spec or another, crosses nothing.
3. An entry without a session, another session's entry, and an entry older than the seventh day never match, as before.
4. A crossing writes `"boundary":true`, `previous_phase` and `previous_spec`, and prints the notice line unchanged; standard output and the exit code never depend on it.
5. The reworded clause keeps its id, its `mandatory` enforcement and no `replaces`, and the module gains or loses no clause or rule.

### Data Models

No stored data changes. The `roundfix/phase-record/v1` schema, its key order and the conditional `previous_phase` and `previous_spec` keys stay as they are. Entries written before this Spec keep the mark the earlier rule gave them.

### API Contracts

1. API Contract: `roundfix phase start` boundary — the command marks the entry and prints the Boundary Notice only under invariant 1. Its usage text, flags, standard output, refusals and exit codes are unchanged.
2. API Contract: the Boundary Notice — unchanged text: `roundfix: notice: this session already recorded <phase>[ of <spec>]; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.` After this Spec `<phase>` is always `release`.
3. API Contract: the `autonomous-work` clause `clause.autonomous.loop-07-outlive-the-turn` — its guidance becomes, word for word: "A supervising session that intends to keep working unattended must arm something that outlives its own turn before that turn ends: a detached Run, a scheduled wake-up, or a daemon. Never end a turn stating that work will continue, because a turn that has ended is indistinguishable, to whoever is waiting, from a turn still working. When no such mechanism is armed, say plainly that work is stopping and what remains. The phase boundary of a Spec comes after the Spec is delivered and released, or after it merges when it releases nothing; there, write the session handoff with the Spec, the phase that ended, the next phase, each open decision and what stays armed, tell the user what stays armed, and start no further phase in that session. An authoring phase is not a phase boundary: the session goes on to the next phase, and each authoring skill says when that phase runs in a subagent with its own context. A scheduled wake-up returns to the same session with its whole context, so it keeps a phase going and is never the way across a boundary."
4. API Contract: the authoring skills' `## Phase boundary` — the first paragraph, which records the phase, is unchanged. The rest becomes, with the same words in all four skills:

   "When the phase ends, no Phase Boundary follows: a Spec's Phase Boundary follows its delivery and release, or its merge when it releases nothing.

   - **You are the phase's Author Subagent:** return your report and end no session; the session that started you goes on to the next phase.
   - **Otherwise** go on to the next phase in this session, handing an Authoring Phase to its Author Subagent as its skill's `Declared model` section says, and write no Session Handoff for this phase. At the Spec's Phase Boundary, write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed, tell the user what stays armed, and start no further phase in this session. The next session starts with the Roundfix skill's `follow-up` Skill Mode, which claims the handoff."
5. API Contract: the Roundfix skill's modes — in step 8 of `--follow-up`, "At a Phase Boundary, after each Authoring Phase and after a Spec is delivered and released," becomes "At a Phase Boundary, after a Spec is delivered and released, or after it merges when it releases nothing,", and after "which claims the handoff." comes "An Authoring Phase ends at no Phase Boundary: the session goes on to the next phase and hands each Authoring Phase to its Author Subagent as the authoring skill's `Declared model` section says." In `--unattended`, "With nobody present the session still ends at a Phase Boundary and says" becomes "With nobody present the session still ends at a Phase Boundary, which follows a Spec's delivery and release and never an Authoring Phase, and says".
6. API Contract: the Roundfix skill's phase reference — its boundary paragraph becomes: "When the same session's latest phase entry of the last seven days is `release` and the new entry is not a `release` of the same Spec, the command prints the Boundary Notice on standard error and marks the entry; no other earlier phase prints it. The Boundary Notice refuses nothing: standard output and the exit code stay as they are. A Phase Boundary follows a Spec's delivery and release, never an Authoring Phase: there, write the Session Handoff and continue in a new session."

### Surface Transcripts

1. Surface Transcript: in a Git repository, with `CLAUDE_CODE_SESSION_ID` set and an earlier `roundfix phase start prd --spec 0412-example` of the same session, the next Authoring Phase prints no notice.

   ```transcript
   $ roundfix phase start techspec --spec 0412-example
   stdout:
   Recorded techspec of 0412-example for claude session <session>
   stderr:
   exit: 0
   ```

2. Surface Transcript: in the same conditions, after an earlier `roundfix phase start release --spec 0412-example` of the same session, a session that goes on receives the notice.

   ```transcript
   $ roundfix phase start follow-up
   stdout:
   Recorded follow-up for claude session <session>
   stderr:
   roundfix: notice: this session already recorded release of 0412-example; at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.
   exit: 0
   ```

## Coverage Map

- Goal 1 (no boundary after an Authoring Phase) → the boundary rule (invariants 1–2), API Contracts 3–5, Tasks 1–3.
- Goal 2 (the boundary after delivery and release) → API Contracts 3–6, Tasks 2 and 3.
- Goal 3 (the notice only after a release) → the boundary rule, API Contracts 1, 2 and 6, Task 1.
- Goal 4 (one boundary everywhere) → Tasks 1–3 and the QA gate's search for the old wording.
- Core Feature 1 → Task 1. Core Feature 2 → Task 2. Core Feature 3 → Task 3. Core Feature 4 → Task 1 (phase reference, user guide, Boundary Notice) and Task 3 (Phase Boundary, ADR-0407).
- Success Metric 1 → Task 1, Surface Transcripts 1 and 2. Success Metric 2 → Tasks 1–3 and the QA gate. Success Metric 3 → Task 2. Success Metric 4 → QA gate.

Deviations from the PRD: None.

## Absence Evidence

None. The design adds no component: it changes one condition of `Recorder.boundary` and rewords existing texts. The searches that located every text stating the per-phase boundary were `rtk proxy grep -rli "phase boundary" . --exclude-dir=.git --exclude-dir=history` and `rtk proxy grep -rn "continue in a new session" . --exclude-dir=.git --exclude-dir=history`; their results are the files the Build Order names, plus `CHANGELOG.md`, which a release writes, and the two Backlog Entries this Spec adopts or revises.

## ADR Alternatives

- ADR-0413: the design is that ADR's decision. It rejects a boundary on a change of Spec, no notice at all, a boundary after `delivery`, and delegating the fix route's phases.
- ADR-0407: amended. Its Phase Record, phase and prompt entries, notice text and refusal of nothing stay; its boundary after each Authoring Phase goes. Its rejected alternatives (a refusal at a boundary, Roundfix opening the next session, labels set where the session starts) stay rejected: the notice still refuses nothing and Roundfix still ends no session.
- ADR-0401: kept. Each Authoring Phase of a feature still goes to its Author Subagent; the design relies on it for the context a phase starts with and takes none of its rejected alternatives.
- ADR-0408: kept. A fix or refactor is still authored in the session unless the user asks for the Author Subagent; the design does not take the per-fix question that ADR rejected.
- ADR-0410: kept. The status view reads phase entries, never the `boundary` mark, so the changed mark does not move what it shows. ADR-0410: "From the Phase Record it reads phase entries only, and it reports no usage."
- ADR-0186: kept. The reworded clause names no Roundfix-only artifact and keeps its `mandatory` enforcement. ADR-0186: "Every clause keeps its enforcement level when its wording changes".
- ADR-0257: kept. The clause keeps its identity, so no Source Baseline row and no `replaces` entry is needed. ADR-0257: "A clause whose wording changes keeps its identity".

- ADR-0258: not applicable — it cites ADR-0257 for authoring against the QA rerun classes and for reduced retirement history, and this design changes neither.
- ADR-0342: not applicable — it cites ADR-0257 for Task settlement on the repository Verification, and this design changes no settlement.
- ADR-0353: not applicable — it cites ADR-0410 for the follow-up list's Waiting Items, and this design changes neither the list nor what it holds back.
- ADR-0405: not applicable — it cites ADR-0401 for reporting a user-level skill copy, and this design changes no installation or report of a skill copy.

## Integration Points

None outside the repository. The Phase Record lives under Roundfix Home, and its outside reader, the operator's usage script, reads identifiers and times; the `boundary` mark it may read only marks fewer entries.

## Testing Approach

- **Task 1.** The table test `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession` in `internal/phaserecord/record_test.go` is rewritten to the new rule: every one of the seven phases as the earlier entry, of the same and of another Spec; `release` again of the same Spec and with no Spec; another session; no session; the sixth and the eighth day; and the latest entry deciding. `TestStartPhaseSkipsUnrelatedAndMalformedHistory` and `TestStartPhaseReadsAnOversizedLineWithoutBufferingIt` keep their purpose with a `release` as the valid earlier entry. Every time comes from the injected clock. In `internal/cli/phase_test.go` a new test asserts Surface Transcripts 1 and 2 whole, and `transcript_3` of `TestPhaseStartSurfaceTranscripts`, whose contract moves, takes a `release` as its earlier entry. `skills/phase_reference_test.go` asserts the new paragraph and the absence of the old sentence.
- **Task 2.** `internal/baseline/phase_boundary_clause_test.go` asserts the new phrases in the clause, the Standard TypeScript golden and this repository's guide, exactly once each, and the absence of "after each authoring phase".
- **Task 3.** `skills/phase_boundary_test.go` asserts the new sentences in each authoring skill's section and in the modes, and the absence of "When the phase ends, a Phase Boundary follows:" and "after each Authoring Phase".
- **QA.** The outside-evidence replay of the measured entries through the built binary, the two Surface Transcripts, and a search of the repository, outside `docs/history/` and `CHANGELOG.md`, for the old wording.

## Build Order

1. The Boundary Notice follows only a release: the boundary rule, its tests, the CLI transcripts, the phase reference, the user guide and the glossary's Boundary Notice.
2. The `autonomous-work` Baseline clause states the boundary after delivery and release, with its derived files and this repository's guide. Independent of step 1.
3. The authoring skills and the Roundfix skill's modes state the boundary after delivery and release, with the glossary's Phase Boundary and the amendment of ADR-0407 (depends on: 1, because both raise the Roundfix skill's version and both edit `CONTEXT.md`).
4. The QA gate (depends on: 2, 3).

## Risks & Considerations

- **Two Specs raise the same skill version.** Specs 0399 and 0400 also raise the Roundfix skill's version, and Spec 0399 edits Baseline modules. The record commands write the next version from the tree they run on, and the Delivery Queue regenerates the derived paths of a conflict.
- **A test that keeps passing.** The current skill and clause tests still pass after the reword, because they assert phrases the new text keeps. Each Task therefore adds an assertion of the new wording and of the absence of the old one, so its Verification fails before the work.
- **The record's earlier marks.** Entries written by 0.146.0 and 0.147.0 keep `"boundary":true` after Authoring Phases. A reader that counts marks across that release should split by date.

## Glossary

- changes: **Phase Boundary**
- changes: **Boundary Notice**

## Decisions

- The crossing condition is invariant 1, and the notice keeps its text. See ADR-0413.
- The boundary of a Spec that releases nothing is its merge, stated as an instruction, with no notice. See ADR-0413.
- The fix route stays in the session; the four skills' `Declared model` sections are not changed. See ADR-0408.
- The clause keeps its id and enforcement, so no Source Baseline row is added. See ADR-0257.
