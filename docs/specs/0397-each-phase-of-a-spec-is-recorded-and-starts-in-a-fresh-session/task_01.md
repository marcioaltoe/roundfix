---
task: task_01
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
status: completed
type: backend
complexity: high
---

# Task 01: The Phase Record and `roundfix phase start` record the phase a session starts

## Overview

A new store package appends phase entries to day files under Roundfix Home
and decides whether a session crossed a Phase Boundary, and a new command,
`roundfix phase start <phase> [--spec <slug>]`, writes one entry and prints
the Boundary Notice when the same session already recorded another phase.
The slice is whole on its own: with a temporary home and a session identity
supplied to the command, each of the seven phases leaves one line a JSON
reader parses, the second phase of a session prints the notice and keeps its
exit code, and a refusal writes nothing. The command's guide, the usage
guide's Phase Record section, the Roundfix skill's `phase` reference and the
four glossary terms describe it.

## Requirements

1. MUST add the package `internal/phaserecord` with the phase half of
   `_techspec.md` → Interfaces and Data Models: `Schema`, `Recorder`, `New`,
   `StartPhase`, `ParsePhase`, `Phase`, `PhaseEntry` and `Boundary`, with
   those signatures. It imports the standard library only and receives the
   home directory and the clock from its caller. `StartPhase` creates
   `<home>/.roundfix/phase-record` with mode `0700` when it is missing, opens
   the day file of the clock's UTC date, `YYYY-MM-DD.jsonl`, with
   `O_CREATE|O_APPEND|O_WRONLY` and mode `0600`, writes the entry as one
   line in one write and closes the file. The line holds the keys of a phase
   entry in the order Data Models gives, `previous_phase` and
   `previous_spec` only when `boundary` is true, every time as RFC 3339 in
   UTC with seconds, an absent value as the empty string and never a missing
   key, and each string value cut at 1,024 bytes so that the line still
   parses as JSON. `RecordPrompt` and `RemoveBefore` belong to task 02.
2. MUST implement the boundary rule of `_techspec.md` → API Contract 2 inside
   `StartPhase`, before the append: it reads the day files of the current UTC
   date and of the six dates before it, takes the latest phase entry whose
   `session` equals the new entry's non-empty `session`, and reports a
   boundary as that contract defines it. A line that does not parse, a line
   of another kind and a file whose name is not a date are skipped, never an
   error.
3. MUST add `roundfix phase start <phase> [--spec <slug>]` as
   `_techspec.md` → API Contract 1 and API Contract 3 define it, in a file of
   its own beside the other commands, with the three edits of
   `internal/cli/cli.go` that `_techspec.md` → Current behavior names:
   - the root usage gains the line
     `  roundfix phase start <phase> [--spec <slug>]` right after the
     `follow-up` line, and the `Commands:` list gains the row
     `  phase      Record the Supervisor Phase a session starts` right after
     the `follow-up` row;
   - `roundfix phase --help`, `roundfix phase -h` and the same two flags
     after `phase start` print the usage of Surface Transcript 1 and exit
     `0`;
   - the session identity comes from `harnessSession`, which reads
     `CLAUDE_CODE_SESSION_ID` through the command environment's lookup and
     never through the process environment, and the time comes from a clock
     the command's dependencies hold, so that a test fixes both;
   - the result line on standard output is
     `Recorded <phase> of <spec> for claude session <session>`, with
     ` of <spec>` left out when the entry names no Spec, and with
     ` for claude session <session>` replaced by
     ` without a harness session: the environment names none` when the
     environment names none (Surface Transcripts 2 and 4);
   - the Boundary Notice of API Contract 3 follows on standard error, word
     for word, when `StartPhase` reports a boundary, and changes neither
     standard output nor the exit code;
   - every refusal exits `2`, prints nothing on standard output, prints
     `roundfix: phase start failed: <reason>` and then
     `Run 'roundfix phase --help' for usage.` on standard error as Surface
     Transcript 5 shows, and writes nothing, not even the directory: an
     unknown phase with the reason of Surface Transcript 5, a phase that
     requires `--spec` without one, a `--spec` that does not match
     `^[0-9]{4}(-[a-z0-9]+)+$`, an extra argument, a missing or unknown
     subcommand, and a working directory outside a Git repository;
   - a record that cannot be written exits `1` with
     `roundfix: phase start failed: <error>` on standard error;
   - the command loads no Project Config and no User Config, so a
     `.roundfixrc.yml` that does not parse never stops it.
4. MUST add these tests, each top-level test calling `t.Parallel()` first,
   with temporary homes and repositories through the packages' existing
   seams, never the real home, and without editing
   `internal/cli/cli_test.go`:
   - `internal/phaserecord`: `TestAPhaseEntryIsOneLineOfItsDeclaredKeys`,
     with the subtests `keys_in_order` (the decoded line holds exactly the
     declared keys in their order, so an added key fails it),
     `boundary_keys_only_at_a_boundary`, `modes` (directory `0700`, file
     `0600`, asserted where the platform keeps owner-only modes),
     `one_line_per_entry` and `oversized_field`;
     `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession`, with the
     subtests `another_phase`, `another_spec`, `after_follow_up`,
     `delivery_then_release_of_the_spec`,
     `delivery_then_release_without_a_spec`,
     `delivery_then_release_of_another_spec`, `same_phase_again`,
     `another_session`, `no_session`, `six_days_old`, `eight_days_old` and
     `latest_entry_decides`, each asserting the returned `Boundary` and the
     written line; and `TestParsePhaseKnowsTheSevenPhases`.
   - `internal/cli`: `TestPhaseStartSurfaceTranscripts`, with the subtests
     `transcript_1` to `transcript_5`, each asserting every line of standard
     output and standard error and the exit code of that Surface Transcript
     with its fixture, and, for transcripts 2 to 4, the whole line the day
     file gained; `TestPhaseStartRecordsEachOfTheSevenPhases` (one line per
     call, the schema, the kind `phase` and only the declared keys);
     `TestPhaseStartRefusalsWriteNothing`, with the subtests `missing_spec`,
     `malformed_spec`, `extra_argument`, `unknown_subcommand` and
     `outside_a_repository`; `TestPhaseStartExitsOneWhenTheRecordCannotBeWritten`,
     which puts a regular file where the record's directory belongs;
     and `TestPhaseStartLoadsNoConfiguration`.
   - `skills`: `TestThePhaseReferenceDocumentsTheCommandAndTheRecord`, which
     reads both copies through `roundfixLayoutSources` and requires, in one
     named subtest each, every phrase requirement 6 lists as pinned,
     exactly once, whitespace-normalized.
   - MUST change the reference count that
     `TestRoundfixSkillIndexNamesEveryReference` requires from 20 to 21, and
     nothing else in that test. Every test that pins the root usage on the
     default branch when the Task runs moves to the new text. No recorded
     test is removed or renamed, and every test file compiles for Windows.
5. MUST write the guides:
   - `docs/user-guide/commands/phase.md`, a new command guide in the shape
     of its neighbours: the usage line, the seven phases, when `--spec` is
     required and its form, the two result lines, the Boundary Notice word
     for word in a `text` block with its three exceptions, the exit codes
     `0`, `1` and `2`, and that the command loads no configuration and
     refuses nothing at a boundary;
   - the row for `phase` in the command index of
     `docs/user-guide/commands.md`, in the index's order;
   - in `docs/user-guide/usage.md`, a section `## Phase Record` right before
     `## Command reference`, and the row
     ``| `phase start` | Record in the Phase Record the phase a session starts |``
     in the command table. The section gives the record's place
     (`phase-record` under Roundfix Home, one `YYYY-MM-DD.jsonl` file per
     UTC date, readable by the user alone), the schema name, both entry
     kinds field by field in the order of `_techspec.md` → Data Models, the
     time format, the empty-string rule, the 1,024-byte bound, the seven
     phases, the boundary rule with its three exceptions and the notice, and
     the sentence "The record is the contract an outside reader uses:
     Roundfix reads no harness transcript, prices no token and reports no
     limit percentage." The section gives the prompt entry's format only:
     who writes one, and the retention, are added there by task 02.
6. MUST add the Roundfix skill's `phase` reference,
   `.agents/skills/roundfix/references/phase.md`, with the heading
   `## Phase Record`, a `bash` block holding
   `roundfix phase start <phase> [--spec <slug>]`, and these three
   paragraphs, word for word:

   > Run `roundfix phase start` when a session starts a Supervisor Phase: `follow-up`, `idea`, `prd`, `techspec`, `tasks`, `delivery` or `release`. `--spec` takes the Spec folder slug and is required for every phase but `follow-up` and `release`. The command appends one phase entry to the Phase Record and prints one line. It loads no Project Config and no User Config. It exits `2` and writes nothing for an unknown phase, a missing or malformed `--spec`, an extra argument or a directory outside a Git repository, and `1` when the record cannot be written. When the installed Roundfix does not know the command, say so and continue.

   > The Phase Record is the day files in the directory `phase-record` under Roundfix Home, one JSON object per line with the schema `roundfix/phase-record/v1`. A phase entry names the phase, the Spec, the repository, the harness and its session, read from the environment; without a session the entry is written with none. No command lists the record: read the files.

   > When the same session last recorded a different phase or Spec in the last seven days, the command prints the Boundary Notice on standard error and marks the entry. The Boundary Notice refuses nothing: standard output and the exit code stay as they are. At a Phase Boundary, write the Session Handoff and continue in a new session.

   The test of requirement 4 pins, each in its own subtest: the usage line
   of the `bash` block (`the_command`); "`follow-up`, `idea`, `prd`,
   `techspec`, `tasks`, `delivery` or `release`" (`the_seven_phases`);
   "required for every phase but `follow-up` and `release`"
   (`the_spec_option`); "It loads no Project Config and no User Config."
   (`no_configuration`); "When the installed Roundfix does not know the
   command, say so and continue." (`an_unknown_command_continues`); "the
   directory `phase-record` under Roundfix Home, one JSON object per line
   with the schema `roundfix/phase-record/v1`" (`the_record`); and "The
   Boundary Notice refuses nothing: standard output and the exit code stay
   as they are." (`the_notice_refuses_nothing`).

   MUST add to the reference index of `.agents/skills/roundfix/SKILL.md`,
   between the `modes` row and the `profiles` row, exactly this row:

   ```text
   | [phase](references/phase.md) | `phase` | Recording the phase a session starts. |
   ```

   and no other text: the entry file is held to 20,000 bytes by
   `TestRoundfixSkillEntryFileStaysWithinItsBudget` and measured 19,891
   before this row of 83. MUST then run `make skills-sync` and the version
   record
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the version of `roundfix` from the tree the Task runs on;
   never write a number or a mirror by hand. MUST NOT add or change text
   inside any `### QA settlement` section.
7. MUST add to `docs/references/skill-coverage.json`, in the map's order,
   the entry `command phase start`, with the skill
   `.agents/skills/roundfix/references/phase.md` and the sources
   `docs/user-guide/commands/phase.md`, `internal/cli/phase.go` and
   `internal/phaserecord/`, and the entry
   `guide docs/user-guide/commands/phase.md`, with the same skill and the
   guide as its source. MUST then re-record the Behavior Surface Record as
   the Task's last step, after the last guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
8. MUST write into `CONTEXT.md`, through `domain-modeling`, these four
   entries right after the entry **Research Subagent**, in this order, each
   followed by its `_Avoid_` line and one blank line as the neighbouring
   entries are, where `(ADR-0407)` is the glossary's Markdown link to
   `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`
   with the text `ADR-0407`:

   > **Phase Record**:
   > The append-only record under Roundfix Home, in the directory `phase-record`, of which Spec and Supervisor Phase each Supervisor session works on and which Work Item each prompt of a Run or of a pre-PR review served. It is a set of day files named by the UTC date, one JSON object per line under the schema `roundfix/phase-record/v1`. An entry holds identifiers and times, never prompt text, a token count, a price or an environment value; no line is rewritten, and the Run Retention Sweep removes a day file older than its cutoff (ADR-0407).
   > _Avoid_: Usage log, phase marker, session journal

   > **Supervisor Phase**:
   > One of the seven phases of a Supervisor's work that the Phase Record names: `follow-up`, the four Authoring Phases `idea`, `prd`, `techspec` and `tasks`, `delivery` and `release`. `roundfix phase start <phase>` records that the session that runs it starts one, with the Spec that `--spec` names; every phase but `follow-up` and `release` requires a Spec (ADR-0407).
   > _Avoid_: Stage, work mode, session type

   > **Phase Boundary**:
   > The point that follows each Authoring Phase and the delivery and release of a Spec. There the Supervisor writes the Session Handoff, tells the user what stays armed and starts no further phase in the session, and the next session starts with the follow-up Skill Mode; an Author Subagent ends no session and returns its report. Roundfix cannot end a session, so the boundary is an instruction that the Boundary Notice recalls (ADR-0407).
   > _Avoid_: Session limit, checkpoint, compaction point

   > **Boundary Notice**:
   > The one line `roundfix phase start` prints on standard error when the same session, in the last seven days of the Phase Record, last recorded a different Supervisor Phase or Spec. It names the earlier phase and Spec and marks the new entry. An earlier `follow-up`, `delivery` followed by `release` of the same Spec, and the same phase of the same Spec recorded again print none. It refuses nothing: standard output and the exit code are those of the command without it (ADR-0407).
   > _Avoid_: Boundary error, session warning, phase refusal

   MUST add no other entry and change no other entry.
9. Of the Governed Paths the Task may change only
   `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/phase.md` and
   `skills/roundfix/SKILL.md`, the last through the two commands of
   requirement 6. MUST NOT change `internal/cli/cli_test.go`,
   `docs/references/coverage-record.json`, the Daemon, `roundfix review`,
   the GC Command, any other skill, any Baseline file, `.roundfixrc.yml`,
   the Makefile, the CI workflows, `go.mod`, `go.sum` or the lint
   configuration. MUST NOT run a `roundfix` found on `PATH` for
   `skills install` or `baseline update`.

## Subtasks

- [ ] Add the store package with the phase entry and the boundary rule, and its tests.
- [ ] Add the command, its usage and its tests.
- [ ] Write the command guide, the index row and the usage guide's section.
- [ ] Add the `phase` reference and its index row, sync the mirror and record the skill version.
- [ ] Add the two Skill Coverage Map entries and the four glossary entries.
- [ ] Re-record the Behavior Surface Record after the last guide edit.

## Acceptance Criteria

- [ ] For each of the seven phases, `roundfix phase start` appends exactly one line to the day file of the current UTC date, with the schema name, the kind `phase` and only the declared keys in their order, in a directory and a file readable by their owner alone.
- [ ] A second phase recorded by the same session prints the Boundary Notice with the earlier phase and Spec on standard error, marks the entry and leaves standard output and the exit code as they are without it; an earlier `follow-up`, `delivery` followed by `release` of the same Spec or of none, the same phase of the same Spec, another session and an entry without a session print none; an entry eight days old is not matched and one six days old is.
- [ ] Without a session identity the entry is written with none and the command says so; every refusal exits `2` and writes nothing; an unwritable record exits `1`; a Project Config that does not parse changes nothing.
- [ ] The command guide, the index, the usage guide, the `phase` reference in both copies, the Skill Coverage Map and the Behavior Surface Record describe the command, and the Roundfix skill's entry file stays within its budget.
- [ ] `CONTEXT.md` defines **Phase Record**, **Supervisor Phase**, **Phase Boundary** and **Boundary Notice**.

## Context

- instruction: `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`
- instruction: `internal/judge/log.go`
- instruction: `internal/cli/reopen.go`
- instruction: `internal/cli/command_help_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/phaserecord/record.go`
- creates: `internal/phaserecord/record_test.go`
- creates: `internal/cli/phase.go`
- creates: `internal/cli/phase_test.go`
- interface: `internal/cli/cli.go`
- creates: `docs/user-guide/commands/phase.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- creates: `.agents/skills/roundfix/references/phase.md`
- creates: `skills/roundfix/references/phase.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `skills/phase_reference_test.go`
- interface: `skills/roundfix_layout_test.go`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- The root usage and the command dispatch — consumers: `./internal/cli`
- The Roundfix skill's reference index and the owned-skill version record — consumers: `./skills`
- The command index, the usage guide, the Skill Coverage Map, the Behavior Surface Record and the glossary — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestAPhaseEntryIsOneLineOfItsDeclaredKeys|TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession|TestParsePhaseKnowsTheSevenPhases)$' ./internal/phaserecord 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAPhaseEntryIsOneLineOfItsDeclaredKeys/keys_in_order TestAPhaseEntryIsOneLineOfItsDeclaredKeys/modes TestAPhaseEntryIsOneLineOfItsDeclaredKeys/oversized_field TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/another_phase TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/after_follow_up TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/delivery_then_release_of_the_spec TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/delivery_then_release_of_another_spec TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/same_phase_again TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/no_session TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/eight_days_old TestParsePhaseKnowsTheSevenPhases; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 15m ./internal/phaserecord && GOOS=windows go build -buildvcs=false ./internal/phaserecord` — expected: exit 0. Before this Task the package does not exist, so the test command fails; afterwards a phase entry is one line of its declared keys, the boundary rule holds case by case, and the package builds for Windows.

- `out="$(go test -count=1 -timeout 30m -v -run '^(TestPhaseStartSurfaceTranscripts|TestPhaseStartRecordsEachOfTheSevenPhases|TestPhaseStartRefusalsWriteNothing|TestPhaseStartExitsOneWhenTheRecordCannotBeWritten|TestPhaseStartLoadsNoConfiguration|TestEveryCommandPathPrintsItsHelp)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPhaseStartSurfaceTranscripts/transcript_1 TestPhaseStartSurfaceTranscripts/transcript_2 TestPhaseStartSurfaceTranscripts/transcript_3 TestPhaseStartSurfaceTranscripts/transcript_4 TestPhaseStartSurfaceTranscripts/transcript_5 TestPhaseStartRecordsEachOfTheSevenPhases TestPhaseStartRefusalsWriteNothing/missing_spec TestPhaseStartRefusalsWriteNothing/malformed_spec TestPhaseStartRefusalsWriteNothing/outside_a_repository TestPhaseStartExitsOneWhenTheRecordCannotBeWritten TestPhaseStartLoadsNoConfiguration TestEveryCommandPathPrintsItsHelp/phase_start_--help; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the command and its tests do not exist, so the pass lines are missing; afterwards the five Surface Transcripts hold line for line, every refusal writes nothing, and the root usage names a command that prints its help.

- `for f in docs/user-guide/commands/phase.md .agents/skills/roundfix/references/phase.md skills/roundfix/references/phase.md; do test -f "$f" || { echo "missing: $f" >&2; exit 1; }; done; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { echo 'the Roundfix skill differs from its mirror' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/phase.md | grep -qF -- 'at a phase boundary, write the Session Handoff and continue in a new session. Nothing was refused.' || { echo 'the phase guide lacks the Boundary Notice' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/usage.md | grep -qF -- 'The record is the contract an outside reader uses: Roundfix reads no harness transcript, prices no token and reports no limit percentage.' || { echo 'the usage guide lacks the Phase Record section' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^(TestThePhaseReferenceDocumentsTheCommandAndTheRecord|TestRoundfixSkillIndexNamesEveryReference|TestRoundfixSkillEntryFileStaysWithinItsBudget|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestThePhaseReferenceDocumentsTheCommandAndTheRecord TestRoundfixSkillIndexNamesEveryReference TestRoundfixSkillEntryFileStaysWithinItsBudget TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guide and the reference do not exist; afterwards both copies of the reference carry each pinned sentence once, the index names 21 references within the entry file's budget, the skill version is recorded, and the skills suite and the documentation contracts, the Skill Coverage Map and the Behavior Surface Record included, pass.

- `grep -F -- '**Phase Record**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Phase Record\*\*(:| —( |$))' || { echo 'missing term: Phase Record' >&2; exit 1; }; grep -F -- '**Supervisor Phase**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Supervisor Phase\*\*(:| —( |$))' || { echo 'missing term: Supervisor Phase' >&2; exit 1; }; grep -F -- '**Phase Boundary**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Phase Boundary\*\*(:| —( |$))' || { echo 'missing term: Phase Boundary' >&2; exit 1; }; grep -F -- '**Boundary Notice**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Boundary Notice\*\*(:| —( |$))' || { echo 'missing term: Boundary Notice' >&2; exit 1; }` — expected: exit 0. Before this Task `CONTEXT.md` defines none of the four terms; afterwards each has its Glossary Definition.

## References

- `_prd.md` → User Stories 1, 2 and 6; Core Features 1, 2, 3 and 13; Success Metrics 1, 2 and 3; Glossary
- `_techspec.md` → Current behavior; Interfaces; Data Models; API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Invariants 1, 2, 4 and 5; Declared breaks; Testing Approach; Build Order 1
- ADR-0407

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-context, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, domain-modeling, evidence-gate, systematic-debugging, no-workarounds.

Implemented the phase-only store and `phase start` command, with an injected
clock and command-environment session lookup. Prompt writing and retention
remain outside this slice. The initial checkout had only this Task file's
Daemon-owned status change; it is preserved.

### Acceptance evidence

1. The store appends one JSON line in one write to the UTC day file using
   private creation modes. The declared keys stay ordered, empty strings stay
   present, and boundary-only keys are conditional. Focused tests cover all
   seven phases through CLI dispatch, exact whole lines, both creation modes
   on platforms that preserve them, embedded newlines and quotes, and the
   1,024-byte limit including a multibyte rune cut.
2. The boundary tests exercise every named case: another phase or Spec,
   earlier follow-up, delivery then release of the same Spec, no Spec or
   another Spec, a repeated phase, another session, no session, six- and
   eight-day history, and selection by the latest recorded time. The CLI
   transcript check asserts the notice, unchanged result line and exit `0`;
   an additional CLI check covers an earlier entry without a Spec. Another
   store test skips malformed lines larger than 1 MiB, a prompt entry, an
   invalid time, another session and a file without a date name.
3. Surface Transcripts 1–5 match stdout, stderr and exit status exactly;
   transcripts 2–4 also assert the entire appended line. All refusal fixtures
   leave even `.roundfix` absent. The obstructing-file fixture exits `1`
   without changing the obstruction. Malformed Project Config and User
   Config at their actual paths do not stop recording; a nested working
   directory records its repository root. Tests supply session identity and
   clock through the existing command-dependency seams and use temporary
   homes and repositories.
4. Added the command guide and command-index row, the usage guide's Phase
   Record section, both skill-reference copies and its index row, and both
   sorted Skill Coverage Map entries. The required recorders generated skill
   version `0.1.241` from `0.1.240` and the Behavior Surface Record after the
   last guide edit. The skill trees are byte-identical; the entry file is
   19,974 bytes against its 20,000-byte budget. Focused skill, command-guide,
   link and glossary contracts passed.
5. Added exactly the four prescribed glossary entries immediately after
   Research Subagent, each with its Avoid line and accepted ADR-0407 link.
   Focused glossary contracts passed.

### Focused checks and regeneration

- Red starting signal: `go test ./internal/phaserecord -run '^TestParsePhaseKnowsTheSevenPhases$/prd$'`
  failed because the package had no non-test Go files before implementation.
- `go test ./internal/phaserecord -run 'Test(APhaseEntry|StartPhase|ParsePhase)' -count=1`:
  exit `0` after adding the store.
- `go test ./internal/cli -run '^TestPhaseStart' -count=1`: exit `0`.
- `go test ./internal/phaserecord ./internal/cli -run 'Test(APhaseEntry|StartPhase|ParsePhase|PhaseStart|PhaseBoundary|EveryCommandPathPrintsItsHelp)' -count=1`:
  exit `0`, including discovery of `phase` and `phase start` with both help flags.
- `go test ./internal/phaserecord ./internal/cli -run 'Test(APhaseEntry|StartPhase|ParsePhase|PhaseStart|PhaseBoundary)' -count=1`:
  exit `0` after the additional malformed-history and absent-earlier-Spec cases.
- `go test ./internal/cli -run '^TestPhaseStartLoadsNoConfiguration$' -count=1`:
  exit `0` after placing malformed User Config at `.roundfix/config.yml`.
- `make skills-sync`: exit `0`; run again after version recording.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`:
  exit `0`; raised both version fields and appended the generated digest.
- `go test ./skills -run 'Test(ThePhaseReference|RoundfixSkillIndexNamesEveryReference|RoundfixSkillEntryFileStaysWithinItsBudget)' -count=1`:
  exit `0`.
- `go test -tags docscontract ./internal/docscontract -run 'Test(EveryCommandIsNamedInTheUserGuide|UserGuideLinksResolve|CommandIndexNamesEveryCommandFile|CommandReferenceLinksResolve|EveryCommandIsNamedInTheRoundfixSkill|Glossary)' -count=1`:
  exit `0`.
- Windows test compilation: `GOOS=windows go test -c -o /tmp/roundfix-0397-phaserecord.test.exe ./internal/phaserecord`,
  `GOOS=windows go test -c -o /tmp/roundfix-0397-cli.test.exe ./internal/cli`,
  and `GOOS=windows go test -c -o /tmp/roundfix-0397-skills.test.exe ./skills`:
  each exited `0`; the CLI compilation was refreshed after its final fixture edit.
  These are compile checks, not Windows runtime evidence.
- Final implementation regeneration:
  `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`:
  exit `0`, adding the phase command and guide and refreshing both edited
  guide fingerprints.
- `git diff --check` and `gofmt -l` on the changed Go files: no diagnostics.
  Changed-file postflight found 19 paths, all declared by task_01; no Baseline,
  configuration, Daemon, other skill, other Task or Task Graph file changed.

The sandbox rejected the logical writable root because it contains a symlink;
commands ran with approved elevated execution from the physical checkout.
No declared Verification command, repository Verification, commit, push or PR
operation was run. Status and settlement remain with the Daemon.

### Verification Feedback — attempt 1

The Daemon's `make verify-changed` diagnostic identified
`TestRegenerationIsDeclared` as the failing contract. Its repository-copy
fixture removed the active Specs but retained the live backlog entry that
cited a reference under Spec 0397. The declared Baseline update reconciled
that newly dangling citation, producing an undeclared backlog change inside
the temporary fixture. An overlay probe reproduced the failure and confirmed
the citation rewrite without editing repository files.

Repaired the fixture in `internal/config/regeneration_declared_test.go` so
that replacing the host's Specs with synthetic archived history also removes
the copied live backlog. The actual backlog, production commands, configuration
and declared regeneration commands are unchanged. The contract still checks
all regenerated changes and its three declaration ablations; a new
`an_undeclared_backlog_write_is_refused` subtest explicitly proves that an
undeclared backlog write is still rejected. The existing unsanctioned
coverage-write refusal remains in place.

Focused evidence after the repair:

- `go test -tags docscontract ./internal/config -run '^TestRegenerationIsDeclared$' -count=1 -v`:
  exit `0`; the root test and all five named subtests passed.
- `GOOS=windows go test -tags docscontract -c -o /tmp/roundfix-0397-config.test.exe ./internal/config`:
  exit `0`; compilation evidence only.
- `git diff --check` and `gofmt -l internal/config/regeneration_declared_test.go`:
  no diagnostics.

The additional ordinary test path above repairs the repository contract
fixture exposed by this Task's Verification. It adds no governed mutation;
the original implementation and generated records stay unchanged. The Agent
did not rerun `make verify-changed` or the Task's declared Verification, did
not change Task status, and performed no commit, push or PR operation. The
Daemon owns the full Verification retry and settlement.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/config/regeneration_declared_test.go`
