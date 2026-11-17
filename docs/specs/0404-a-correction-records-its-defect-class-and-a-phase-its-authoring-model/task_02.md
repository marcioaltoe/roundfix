---
task: task_02
spec: 0404-a-correction-records-its-defect-class-and-a-phase-its-authoring-model
status: pending
type: backend
complexity: high
---

# Task 02: A phase entry of an Authoring Phase records its model, effort, runner and repository

## Overview

An entry of `idea`, `prd`, `techspec` or `tasks` in the Phase Record gains
six keys after the nine Spec 0401 leaves: the Repository Identity, the model
the session reports with `--model`, the reasoning effort from
`CLAUDE_EFFORT` or `--effort`, how each was learned, and whether the phase's
Author Subagent ran it. A value the command cannot record becomes `unknown`
with one warning and never costs the entry or changes the exit code. An
entry of `follow-up`, `delivery` or `release` keeps its nine keys byte for
byte. The four authoring skills tell the session what to pass, and the
phase reference, the guides, the glossary and ADR-0407 describe it.

## Requirements

1. MUST add to `internal/phaserecord` the six keys of `_techspec.md` → Data
   Models, after the nine keys Spec 0401 leaves and in the order
   `repository_id`, `model`, `model_source`, `effort`, `effort_source`,
   `runner`, on an entry of an Authoring Phase only, written with or without
   a session; `PhaseEntry` and `RecordedPhase` each gain `RepositoryID`,
   `Model`, `ModelSource`, `Effort`, `EffortSource` and `Runner`. An entry
   without the six keys reads back with the three sources `unknown` and the
   rest empty. An entry of another phase keeps the nine keys
   (`_techspec.md` → invariant 6). The schema keeps its version and the
   command reads no earlier entry (ADR-0414).
2. MUST add to `roundfix phase start` the flags `--model`, `--effort` and
   `--author-subagent` and the rules of `_techspec.md` → API Contract 7: the
   model only from `--model` and only when it matches the pattern the
   Authoring Model Declaration accepts for a model name, at most 200
   characters; the effort from `CLAUDE_EFFORT` first, then `--effort`, each
   when it matches `^[a-z0-9][a-z0-9_-]{0,31}$`; the runner
   `author-subagent` or `session`; the three `unknown` for an Authoring
   Phase without a session, whose standard output stays that of the
   Delivery Base; `repository_id` from `baseline.InspectRepositoryIdentity`
   for every Authoring Phase entry, empty when it cannot be computed; the
   result line with a session; and the three warnings, each one line on
   standard error. The command still loads no configuration, and its
   refusals and exit codes stay those of the Delivery Base (invariant 7).
   The environment is read through the command's injected lookup.
3. MUST make the help of `roundfix phase` Surface Transcript 6 and the root
   help's `phase start` line the usage line it shows; no other root help
   line changes in this Task.
4. MUST write these new tests, each in a new file, each calling
   `t.Parallel()` first:
   - `internal/phaserecord`: `TestAnAuthoringPhaseEntryHoldsTheSixKeysInOrder`
     (the keys of an Authoring Phase entry with every source, and the nine
     of a `release` entry) and
     `TestAnEntryWithoutTheSixKeysReadsBackAsUnknown`.
   - `internal/cli`: `TestPhaseSurfaceTranscript6` and
     `TestPhaseSurfaceTranscript7`, each asserting its whole transcript;
     `TestPhaseStartRecordsTheModelAndEffortOfEachAuthoringPhase`, with a
     subtest for each of the four Authoring Phases and the cases of
     `_prd.md` → Success Metric 8: a reported model with an effort in the
     environment, an effort by flag only, no reported model, no session, a
     model of 300 characters, an effort that does not match, and the
     Author Subagent; each asserts the entry's six values, the result line,
     the warnings and exit 0; `TestPhaseStartIgnoresTheModelFlagsOnAnotherPhase`
     (the warning on `release`, `follow-up` and `delivery`, and their nine
     keys); and `TestPhaseStartRecordsTheRepositoryIdentityOfItsClone`
     (the same identity from a clone and a linked worktree, and none from a
     remote address with a user name and a password).
   - `skills`: `TestTheAuthoringSkillsPassTheirModelToThePhaseCommand`, in
     a new file, over both copies of the four authoring skills.
   Every test sets or clears `CLAUDE_CODE_SESSION_ID` and `CLAUDE_EFFORT`
   through the injected lookup or `t.Setenv`.
5. MUST change these existing tests' expectations, as Spec 0401 leaves
   them, and no other, declaring each in the Result: in
   `internal/cli/phase_test.go`, the help text of
   `TestPhaseStartSurfaceTranscripts`, and the result line and the entry
   line of every case that records an Authoring Phase through
   `phaseExpectedLine`, in that test, in
   `TestPhaseStartRecordsEachOfTheSevenPhases`, which requires nine keys for
   every phase, and in `TestPhaseStartLoadsNoConfiguration`; in
   `internal/phaserecord/record_test.go`, the literal `prd` line of
   `TestAPhaseEntryIsOneLineOfItsDeclaredKeys` and every case that decodes
   an Authoring Phase entry against the nine-key list; in
   `internal/phaserecord/read_test.go`, the entries
   `TestPhasesReadsBackThePhaseEntriesOfTheWindow` expects for an Authoring
   Phase; and in `skills/phase_reference_test.go`, the command line it
   requires in the phase reference. `_tasks.md` → What the Tasks stage
   settled lists what the measurement found; the Task runs
   `go test ./internal/phaserecord ./internal/cli ./skills` on its base and
   on its change and declares every other test whose expectation it had to
   change. It MUST NOT touch `internal/cli/cli_test.go` and renames or
   removes no test.
6. MUST add, in each of the four authoring skills, after the sentence that
   runs `roundfix phase start <phase> --spec <slug>` as Spec 0401 leaves
   it, the paragraph of `_techspec.md` → API Contract 11, the same in the
   four but for the phase; describe the three flags, the six values, the
   warning and the identity of a shallow clone in
   `.agents/skills/roundfix/references/phase.md`; and describe the same in
   `docs/user-guide/commands/phase.md` and in the Phase Record section of
   `docs/user-guide/usage.md`.
7. MUST revise **Phase Record** in `CONTEXT.md`, through domain-modeling, to
   say that an entry of an Authoring Phase holds the model and reasoning
   effort that ran it and how each was learned, and MUST add to ADR-0407 the
   notice "Amended by ADR-0418" with the date and what changed
   (`_techspec.md` → API Contract 12).
8. MUST regenerate, never hand-edit, with the commands of
   `_authorization.md` → Sanctioned regeneration, in this order:
   `make skills-sync`; the `-record-skill-versions` step
   (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
   and the `-record-skill-coverage` step
   (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`).
   It records the Coverage Review of the Behavior Surface
   `guide docs/user-guide/usage.md` in `docs/references/skill-coverage.json`
   as a dated line naming this Task, saying that only the Phase Record
   section changed, which the phase reference covers.
9. MUST change no Governed Path outside `_authorization.md` `paths`, and
   MUST NOT change the `## Declared model` section of any authoring skill,
   any `### QA settlement` section, `.roundfixrc.yml` or any file under
   `.claude/agents/`.

## Subtasks

- [ ] Add the six keys to the writer and the reader, with their tests.
- [ ] Add the three flags, the effort, the identity and the warnings to the phase command, with Surface Transcripts 6 and 7.
- [ ] Move the declared existing tests.
- [ ] Add the paragraph to the four authoring skills and describe the values in the phase reference and the guides.
- [ ] Revise the glossary, amend ADR-0407 and run the sanctioned regeneration.

## Acceptance Criteria

- [ ] An Authoring Phase recorded with a session, a reported model and an effort in the environment holds both values with their sources and who ran the phase; without a reported model the model reads `unknown`; without a session both read `unknown`.
- [ ] A model of 300 characters or a malformed effort is recorded as `unknown` with one warning, and every case exits 0.
- [ ] An entry of `follow-up`, `delivery` or `release` holds the nine keys it holds at the Delivery Base, and an entry written at the Delivery Base reads back as `unknown`.
- [ ] Two clones of one repository record the same `repository_id`, and no entry holds a user name or a password.
- [ ] The four authoring skills, both copies, tell the session to pass the exact model identifier its harness states and `--author-subagent` when the Author Subagent runs the phase.

## Context

- instruction: `docs/adr/0418-a-phase-entry-holds-the-model-and-reasoning-effort-that-ran-an-authoring-phase.md`
- instruction: `docs/adr/0421-a-phase-entry-and-the-corrections-record-name-their-repository-by-its-git-lineage.md`
- instruction: `docs/adr/0425-a-session-passes-its-model-to-the-phase-command-and-only-the-name-its-harness-states.md`
- interface: `internal/phaserecord/record.go`
- interface: `internal/phaserecord/read.go`
- interface: `internal/phaserecord/record_test.go`
- interface: `internal/phaserecord/read_test.go`
- creates: `internal/phaserecord/model_keys_test.go`
- interface: `internal/cli/phase.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/phase_test.go`
- creates: `internal/cli/phase_model_test.go`
- interface: `skills/phase_reference_test.go`
- creates: `skills/phase_model_paragraph_test.go`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/write-tasks/SKILL.md`
- interface: `.agents/skills/roundfix/references/phase.md`
- interface: `skills/roundfix/references/phase.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/phase.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/references/skill-coverage.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`

## Consumers

- The phase entry's format and `phaserecord.RecordedPhase` — consumers: `./internal/phaserecord`, `./internal/cli`
- The help of `roundfix phase` and the root help — consumers: `./internal/cli`
- The authoring skills, the phase reference, the glossary and the Behavior Surface Record — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAnAuthoringPhaseEntryHoldsTheSixKeysInOrder|TestAnEntryWithoutTheSixKeysReadsBackAsUnknown|TestPhaseSurfaceTranscript6|TestPhaseSurfaceTranscript7|TestPhaseStartRecordsTheModelAndEffortOfEachAuthoringPhase|TestPhaseStartIgnoresTheModelFlagsOnAnotherPhase|TestPhaseStartRecordsTheRepositoryIdentityOfItsClone|TestTheAuthoringSkillsPassTheirModelToThePhaseCommand)$' ./internal/phaserecord ./internal/cli ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAnAuthoringPhaseEntryHoldsTheSixKeysInOrder TestAnEntryWithoutTheSixKeysReadsBackAsUnknown TestPhaseSurfaceTranscript6 TestPhaseSurfaceTranscript7 TestPhaseStartRecordsTheModelAndEffortOfEachAuthoringPhase TestPhaseStartIgnoresTheModelFlagsOnAnotherPhase TestPhaseStartRecordsTheRepositoryIdentityOfItsClone TestTheAuthoringSkillsPassTheirModelToThePhaseCommand; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/phaserecord ./internal/cli` — expected: exit 0. Before this Task the six keys, the three flags and the named tests do not exist, so the command fails; afterwards every named test passes and both packages that read and write the phase entry pass.
- `tr -s '[:space:]' ' ' < docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md | grep -qF -- "Amended by ADR-0418" || { printf 'missing phrase in %s: %s\n' docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md "Amended by ADR-0418" >&2; exit 1; }; grep -qF -- '--author-subagent' .agents/skills/write-prd/SKILL.md && grep -qF -- '--author-subagent' .agents/skills/roundfix/references/phase.md && grep -F -- '**Phase Record**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Phase Record\*\*(:| —( |$))' && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task ADR-0407 carries no notice of ADR-0418 and no skill names `--author-subagent`; afterwards both do, and the skill tests, the documentation contracts, the mirror and the records pass.

## References

- `_prd.md` → Goals (second bullet); User Story 8; Core Features 11, 12, 13 and 22; Success Metrics 8, 9 and 16
- `_techspec.md` → Phase entry values; Interfaces (invariants 6 and 7); Data Models; API Contracts 7, 11 and 12; Surface Transcripts 6 and 7; Testing Approach; Build Order 2
- ADR-0418, ADR-0421, ADR-0425, ADR-0407, ADR-0414, ADR-0410, ADR-0401, ADR-0415
