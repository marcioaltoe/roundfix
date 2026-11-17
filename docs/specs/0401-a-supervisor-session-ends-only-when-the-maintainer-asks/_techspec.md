---
spec: 0401-a-supervisor-session-ends-only-when-the-maintainer-asks
prd: _prd.md
created: 2026-11-17
---

# A Supervisor session ends only when the maintainer asks — Technical Spec

## Executive Summary

No phase of a Spec ends the Supervisor session any more. The Phase Record's writer stops reading earlier entries: `roundfix phase start` appends its entry with `"boundary":false` and no `previous_phase` or `previous_spec`, and prints no Boundary Notice. Every text that says when a session ends changes with it: the `autonomous-work` Baseline clause, the `Phase boundary` sections of the four authoring skills, the Roundfix skill's `--follow-up` and `--unattended` modes and its phase reference, the glossary, ADR-0407 and the user guide. The trade-off accepted is that one session may carry several Specs and a growing context; the Author Subagents keep each feature phase's own context small, and the harness's compaction and the maintainer's request remain. The `boundary` key stays in the `roundfix/phase-record/v1` schema as a constant, so an outside reader keeps parsing every entry, old and new.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the record keeps its schema and its keys. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; `roundfix phase start` now writes one local day file and reads none. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0414 (new) records the decision, amends ADR-0407 and supersedes ADR-0413; the design keeps ADR-0401 and ADR-0408 (where an Authoring Phase runs), ADR-0410 (the status view reads phase entries, never the `boundary` key), ADR-0186 and ADR-0257 (a reworded Baseline clause keeps its identity and its adopter-neutral words); ADR-0258, ADR-0342, ADR-0353 and ADR-0405 are listed only to record why they do not apply. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: the answer "Mesmas do BE (Recommended)" of 2026-11-16, unchanged by the answer of 2026-11-17, quoted in the authorization record. It covers the `autonomous-work` Baseline module, its derived Standard TypeScript golden and profile pin, this repository's managed `autonomous-work` guide and Setup Manifest, and the owned skills this Spec changes. It grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0401-a-supervisor-session-ends-only-when-the-maintainer-asks/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/phase.md`, `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-techspec/SKILL.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`, `internal/baseline/assets/modules/autonomous-work.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`, `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## System Architecture

Three parts change, and none gains a component.

- **The Phase Record's writer** in `internal/phaserecord/record.go`. `Recorder.StartPhase` loses its lookup of the same session's latest entry (`Recorder.boundary` and the line reader `readDay` that only it uses), and with it the `Boundary` value it returned. `internal/cli/phase.go` stops printing the notice. The reader `Phases` in `internal/phaserecord/read.go`, which the status view uses, does not change.
- **The Baseline clause** `clause.autonomous.loop-07-outlive-the-turn` of `internal/baseline/assets/modules/autonomous-work.json`. Its guidance is reworded; the module record command, `make baseline-digests` and the sanctioned `baseline update` regenerate the derived files and this repository's `docs/agents/autonomous-work.md`.
- **The skills and documents** that say when a session ends: the `## Phase boundary` sections of `write-idea`, `write-prd`, `write-techspec` and `write-tasks`; step 8 of `--follow-up` and the `--unattended` section of the Roundfix skill's `references/modes.md`; its `references/phase.md`; `CONTEXT.md`; ADR-0407; `docs/user-guide/commands/phase.md` and the Phase Record section of `docs/user-guide/usage.md`. Each owned skill's canonical copy under `.agents/skills` is edited, `make skills-sync` writes its mirror, and the version record command raises its version.

## Implementation Design

### Interfaces

```go
// internal/phaserecord/record.go
func (r *Recorder) StartPhase(ctx context.Context, entry PhaseEntry) error

// phaseLine keeps its declaration order, which is the public file format;
// PreviousPhase and PreviousSpec are removed, Boundary is always false.
type phaseLine struct {
	Schema     string `json:"schema"`
	Kind       string `json:"kind"`
	Time       string `json:"time"`
	Phase      Phase  `json:"phase"`
	Spec       string `json:"spec"`
	Repository string `json:"repository"`
	Harness    string `json:"harness"`
	Session    string `json:"session"`
	Boundary   bool   `json:"boundary"`
}
```

The type `Boundary`, the method `Recorder.boundary`, `readDay` and `maxLineBytes` are removed. The callers of `StartPhase` are `internal/cli/phase.go`, `internal/cli/follow_up_queue_test.go`, `internal/phaserecord/read_test.go` and `internal/phaserecord/record_test.go`.

Numbered invariants:

1. `StartPhase` opens no day file for reading. It appends exactly one line, whose keys are, in order, `schema`, `kind`, `time`, `phase`, `spec`, `repository`, `harness`, `session` and `boundary`, with `"boundary":false`, whatever entries the record holds.
2. No new line holds `previous_phase` or `previous_spec`.
3. `Phases` reads a line that an earlier version marked, with `"boundary":true` and both `previous_` keys, as the same `RecordedPhase` as an unmarked one, and does not count it as skipped.
4. `roundfix phase start` keeps its usage text, flags, standard output, refusals and exit codes; after a successful append its standard error is empty.
5. The reworded clause keeps its id, its `mandatory` enforcement, no `replaces` and its first three sentences, and the module gains or loses no clause or rule.

### Data Models

No stored data changes shape. The `roundfix/phase-record/v1` schema and its key order stay; `boundary` becomes the constant `false` in new entries, and the conditional `previous_phase` and `previous_spec` keys are no longer written. Entries written before this Spec keep the marks the earlier rule gave them.

### API Contracts

1. API Contract: `roundfix phase start` — appends one entry under invariants 1 and 2 and prints only its result line on standard output; it prints nothing on standard error after any earlier entry. Its usage text, flags, refusals and exit codes are unchanged.
2. API Contract: the Phase Record line — the nine keys of invariant 1 in that order, `boundary` always `false`; an outside reader of the schema reads old and new lines alike.
3. API Contract: the `autonomous-work` clause `clause.autonomous.loop-07-outlive-the-turn` — its guidance becomes, word for word: "A supervising session that intends to keep working unattended must arm something that outlives its own turn before that turn ends: a detached Run, a scheduled wake-up, or a daemon. Never end a turn stating that work will continue, because a turn that has ended is indistinguishable, to whoever is waiting, from a turn still working. When no such mechanism is armed, say plainly that work is stopping and what remains. No phase of a Spec ends the session, whether a person is present or not: after an authoring phase, a delivery or a release, the session goes on with the next queued work, and each authoring skill says when a phase runs in a subagent with its own context. The session ends only when the user asks or the work stops because nothing remains or what remains waits on the user; then write the session handoff with the Spec, the phase that ended, the next phase, each open decision and what stays armed, and tell the user what stays armed."
4. API Contract: the authoring skills' `## Phase boundary` — the first paragraph, which records the phase, is unchanged. The rest becomes, with the same words in all four skills:

   "When the phase ends, the session goes on: no phase of a Spec ends the session, whether a person is present or not.

   - **You are the phase's Author Subagent:** return your report and end no session; the session that started you goes on with the next phase.
   - **Otherwise** go on in this session with the next phase or the next queued work, handing an Authoring Phase to its Author Subagent as its skill's `Declared model` section says, and write no Session Handoff for this phase. End the session only when the user asks or the work stops because nothing remains or what remains waits on the user; then write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed, and tell the user what stays armed. The next session starts with the Roundfix skill's `follow-up` Skill Mode, which claims the handoff."
5. API Contract: the Roundfix skill's modes — in step 8 of `--follow-up`, the two sentences from "At a Phase Boundary, after each Authoring Phase" to "never crosses a Phase Boundary." become: "No phase of a Spec ends the session, whether a person is present or not: after an Authoring Phase, a delivery or a release, go on with the next queued work, handing each Authoring Phase to its Author Subagent as the authoring skill's `Declared model` section says. End the session only when the user asks or the work stops because nothing remains or what remains waits on the user; then write the Session Handoff with the Spec, the phase that ended, the next phase, each decision still open and what stays armed, and tell the user what stays armed; the next session starts with this mode, which claims the handoff. A scheduled wake-up returns to the same session with its whole context." The last paragraph of `--unattended`, "With nobody present the session still ends at a Phase Boundary and says what stays armed and what waits. A scheduled wake-up keeps a phase going and never crosses a Phase Boundary, because it returns to the same session.", becomes: "With nobody present no phase of a Spec ends the session either: it goes on with the queued work and ends only when nothing remains or what remains waits on the user, saying what stays armed and what waits."
6. API Contract: the Roundfix skill's phase reference — its paragraph that begins "When the same session last recorded a different phase or Spec" becomes: "The command never reads the record back: whatever the session recorded before, it writes `boundary` as `false` and prints nothing on standard error. An entry an earlier Roundfix marked keeps its mark and stays readable. A phase entry ends no session: the session goes on with the next phase and ends only when the user asks or the work stops."

### Surface Transcripts

1. Surface Transcript: in a Git repository, with `CLAUDE_CODE_SESSION_ID` set and an earlier `roundfix phase start prd --spec 0412-example` of the same session, the next Authoring Phase prints only its result line.

   ```transcript
   $ roundfix phase start techspec --spec 0412-example
   stdout:
   Recorded techspec of 0412-example for claude session <session>
   stderr:
   exit: 0
   ```

2. Surface Transcript: in the same conditions, after an earlier `roundfix phase start release --spec 0412-example` of the same session, the session goes on and prints only its result line.

   ```transcript
   $ roundfix phase start follow-up
   stdout:
   Recorded follow-up for claude session <session>
   stderr:
   exit: 0
   ```

## Coverage Map

- Goal 1 (no phase ends the session) → API Contracts 3–6, Tasks 1–3.
- Goal 2 (the session ends when the maintainer asks or the work stops; the turn still arms) → API Contracts 3–5, invariant 5, Tasks 2 and 3.
- Goal 3 (the record keeps measuring, no notice) → invariants 1–4, API Contracts 1, 2 and 6, Task 1.
- Goal 4 (one rule everywhere) → Tasks 1–3 and the QA gate's search for the old wording.
- Core Feature 1 → Task 1. Core Feature 2 → Task 2. Core Feature 3 → Task 3. Core Feature 4 → Task 1 (phase reference, user guide) and Task 3 (glossary, ADR-0407).
- Success Metric 1 → Task 1, Surface Transcripts 1 and 2. Success Metric 2 → Tasks 1–3 and the QA gate. Success Metric 3 → Task 2. Success Metric 4 → QA gate.

Deviations from the PRD: None.

## Absence Evidence

None. The design adds no component: it removes the writer's lookup and the notice and rewords existing texts. The searches that located every text stating when a session ends were `rtk proxy grep -rli -e "phase boundary" -e "boundary notice" -e "continue in a new session" . --exclude-dir=.git --exclude-dir=history` and `rtk proxy grep -rn -e "session handoff" -e "new session" -e "ends the session" .agents/skills docs/agents docs/user-guide`; their results are the files the Build Order names, plus `CHANGELOG.md`, which a release writes, the Vendored Skill `ai-memory-handoff`, which already writes a handoff only when the user wraps up, and the two Backlog Entries this Spec adopts or revises. `roundfix phase start` is the only caller that reads a `Boundary`; `rtk proxy grep -rn "StartPhase(" internal cmd` lists the four callers named in Interfaces.

## ADR Alternatives

- ADR-0414: the design is that ADR's decision. It rejects keeping a boundary after a Spec's release, keeping the notice after a release, dropping the `boundary` key, and a stop after a count of Specs or a context size.
- ADR-0413: superseded by ADR-0414 and deleted in the planning change. Its boundary after delivery and release and its notice after a release are reversed before any Run implemented them.
- ADR-0407: amended. Its Phase Record, phase and prompt entries, schema and refusal of nothing stay; its Phase Boundary that ends the session and its Boundary Notice go. Its rejected alternatives stay rejected: Roundfix still refuses nothing and opens or ends no session.
- ADR-0401: kept. Each Authoring Phase of a feature still goes to its Author Subagent; the design relies on it for the context a phase starts with and takes none of its rejected alternatives.
- ADR-0408: kept. A fix or refactor is still authored in the session unless the user asks for the Author Subagent; the design does not take the per-fix question that ADR rejected.
- ADR-0410: kept. The status view reads phase entries through `Phases`, never the `boundary` key, so a constant `false` does not move what it shows. ADR-0410: "From the Phase Record it reads phase entries only, and it reports no usage."
- ADR-0186: kept. The reworded clause names no Roundfix-only artifact and keeps its `mandatory` enforcement. ADR-0186: "Every clause keeps its enforcement level when its wording changes".
- ADR-0257: kept. The clause keeps its identity, so no Source Baseline row and no `replaces` entry is needed. ADR-0257: "A clause whose wording changes keeps its identity".

- ADR-0258: not applicable — it cites ADR-0257 for authoring against the QA rerun classes and for reduced retirement history, and this design changes neither.
- ADR-0342: not applicable — it cites ADR-0257 for Task settlement on the repository Verification, and this design changes no settlement.
- ADR-0353: not applicable — it cites ADR-0410 for the follow-up list's Waiting Items, and this design changes neither the list nor what it holds back.
- ADR-0405: not applicable — it cites ADR-0401 for reporting a user-level skill copy, and this design changes no installation or report of a skill copy.

## Integration Points

None outside the repository. The Phase Record lives under Roundfix Home, and its outside reader, the operator's usage script, reads identifiers and times; it finds the same keys, with `boundary` always `false` in new entries.

## Testing Approach

- **Task 1.** In `internal/phaserecord/record_test.go`, `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession` gives way to `TestStartPhaseMarksNoBoundaryWhateverCameBefore`, a table over every phase as the earlier entry of the same session, a `release` of another Spec, the same `release` again and a hand-written entry an earlier version marked; `TestStartPhaseSkipsUnrelatedAndMalformedHistory` and `TestStartPhaseReadsAnOversizedLineWithoutBufferingIt` give way to `TestStartPhaseReadsNoEarlierEntry`, which proves the writer never opens a day file to read: a directory stands where yesterday's day file would be, and today's file holds a malformed oversized line. `TestAPhaseEntryIsOneLineOfItsDeclaredKeys` keeps its subtests, with `boundary_keys_only_at_a_boundary` becoming `no_boundary_keys_after_a_release`. `internal/phaserecord/read_test.go` gains a case for an entry an earlier version marked. In `internal/cli/phase_test.go` a new test asserts Surface Transcripts 1 and 2 whole, `transcript_3` of `TestPhaseStartSurfaceTranscripts` takes no notice, and `TestPhaseBoundaryNoticeOmitsAnAbsentEarlierSpec` gives way to a case after a `release` without a Spec. `skills/phase_reference_test.go` asserts the new paragraph and the absence of the notice. Every time comes from an injected clock.
- **Task 2.** `internal/baseline/phase_boundary_clause_test.go` asserts the new sentences in the clause, the Standard TypeScript golden and this repository's guide, exactly once each, keeps the three sentences on arming the turn, and asserts the absence of the old boundary sentences.
- **Task 3.** `skills/phase_boundary_test.go` asserts the new sentences in each authoring skill's section and in the modes, and the absence of the old ones.
- **QA.** The outside-evidence replay of the measured entries through the built binary, the two Surface Transcripts, a comparison with the Delivery Base's binary, sabotages of each changed rule, and a search of the repository, outside `docs/history/` and `CHANGELOG.md`, for the old wording.

## Build Order

1. The Phase Record marks no boundary and `roundfix phase start` prints no notice: the writer, the command, their tests, the reader's test, the phase reference and the user guide.
2. The `autonomous-work` Baseline clause says no phase of a Spec ends the session, with its derived files and this repository's guide. Independent of step 1.
3. The authoring skills and the Roundfix skill's modes say the same, with the glossary's Phase Boundary, Session Handoff and the removal of the notice's entry, and the amendment of ADR-0407 (depends on: 1, because both raise the Roundfix skill's version).
4. The QA gate (depends on: 2, 3).

## Risks & Considerations

- **A session that grows.** One session may now plan, deliver and release several Specs. The Author Subagents keep each feature phase's context their own, the Phase Record still attributes usage by phase, and the harness's compaction and the maintainer's request stay available; no limit is added.
- **Two Specs raise the same skill version.** Spec 0400 also raises the Roundfix skill's version. The record commands write the next version from the tree they run on, and the Delivery Queue regenerates the derived paths of a conflict.
- **Tests that keep passing.** Measured on a scratch worktree of `bcd716df`: with the texts reworded and the writer changed, the current tests that assert the old sentences fail, but each Task still adds assertions of the new wording and of the absence of the old one, so its Verification proves the new contract and not only the removal.
- **The record's earlier marks.** Entries written by 0.146.0 and 0.147.x keep `"boundary":true` where the earlier rules marked them. A reader that counts marks across this release should split by date.

## Glossary

- changes: **Phase Boundary**
- changes: **Session Handoff**

## Decisions

- The writer reads no earlier entry and writes `"boundary":false`; the notice goes. See ADR-0414.
- No phase of a Spec ends the session; it ends when the user asks or the work stops, and the rule to arm what outlives the turn stays. See ADR-0414.
- ADR-0413 is superseded by ADR-0414 and deleted, as the retirement clause asks. See ADR-0414.
- The fix route stays in the session; the four skills' `Declared model` sections are not changed. See ADR-0408.
- The clause keeps its id and enforcement, so no Source Baseline row is added. See ADR-0257.
