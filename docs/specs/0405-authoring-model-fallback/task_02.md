---
task: task_02
spec: 0405-authoring-model-fallback
status: completed
type: backend
complexity: high
---

# Task 02: Roundfix decides and records an Authoring Fallback and prints the model a phase starts on

## Overview

Add the limit reader, the `fallback` kind of entry in the Phase Record, and
the two subcommands `roundfix phase fallback` and `roundfix phase model`, with
their help, guide, skill reference, coverage entries and glossary terms. After
this Task a session can pass the error an Author Subagent ended on and learn
whether to restart it, on which model and effort. A later phase can learn
whether its Spec already fell back on its declared model.

## Requirements

1. MUST add `ReadUsageLimit` with the types of `_techspec.md` → Interfaces. It holds invariants 1 to 3 and names no model.
2. MUST add `RecordFallback` and `Fallbacks` with the line format of `_techspec.md` → Data Models, Fallback entry, keys in that order, appended as `StartPhase` appends.
3. MUST keep `Phases` and every reader of phase entries skipping a `fallback` line without counting it as skipped (invariant 8).
4. MUST implement API Contract 1 (`roundfix phase fallback`) and API Contract 2 (`roundfix phase model`) in `internal/cli/phase.go`:
   - the usage errors, each decision and its order, and the exact lines and exit codes;
   - invariants 4 to 7;
   - repository matching by the identity Spec 0404 records in a phase entry, or by the Git root when either identity is empty.
5. MUST keep `roundfix phase start` as its base leaves it: arguments, output, record and exit codes. Invariant 6 holds.
6. MUST change the usage of `roundfix phase` and the root help as API Contract 3 says. MUST update `TestPhaseStartSurfaceTranscripts/transcript_1` and, when the base holds it, Spec 0404's test of the full `roundfix phase --help` to the new usage, and no other expectation of those tests.
7. MUST add tests, each calling `t.Parallel()`, with the seams of `_techspec.md` → Testing Approach. Each test creates its own temporary home, Git repository and `.roundfixrc.yml`. Tests:
   - `TestReadUsageLimitClassifiesTheHarnessMessages` takes every message quoted in `references/2026-11-17-harness-and-test-measurements.md` → Usage limit messages, the Spec 0404 observation, both apostrophes, upper and lower case, a multi-line error whose second line holds the limit, a line with control characters and a 300-character line.
   - `TestAFallbackEntryIsOneLineOfItsDeclaredKeys` and `TestFallbacksReadsBackTheFallbackEntries` cover the entry, including a line that does not decode.
   - `TestPhasesSkipsAFallbackEntry` covers invariant 8.
   - `TestPhaseFallbackSurfaceTranscripts` asserts Surface Transcript 1, Surface Transcript 2, Surface Transcript 3, Surface Transcript 4, Surface Transcript 5, Surface Transcript 6, Surface Transcript 7, Surface Transcript 8 of `_techspec.md`, each in a subtest named `transcript_<n>`, line by line with the transcripts' controlled variations, and with transcripts 3 to 5 run after transcript 2 in one home.
   - `TestPhaseFallbackRefusesWithoutWriting` asserts each refusal of API Contract 1 and each usage error: its line, its exit code and that no day file changes.
   - `TestPhaseModelKeepsTheFallbackForItsSpecAndModel` covers invariants 4 and 7: two Specs, two repositories, a phase whose declared model is not the limited one, and a phase that declares no fallback.
8. MUST describe both subcommands in `docs/user-guide/commands/phase.md` with their usage, lines and exit tables. MUST describe the fallback entry in the Phase Record section of `docs/user-guide/usage.md`.
9. MUST add to `.agents/skills/roundfix/references/phase.md` the `Authoring Fallback` section of API Contract 7, and make the Roundfix skill name `roundfix phase model` and `roundfix phase fallback` in its row for the phase command.
10. MUST add to `docs/references/skill-coverage.json` the entries `command phase model` and `command phase fallback`, each covered by `.agents/skills/roundfix/references/phase.md` with the sources of `command phase start`. MUST record a dated review line naming this Task on `command phase start`.
11. MUST add **Authoring Fallback** to `CONTEXT.md` and revise **Phase Record** so that it holds fallback entries, with the phrase `fallback entries`, both through domain-modeling. MUST add to ADR-0407 the notice `Amended by ADR-0428` with the date and what changed (`_techspec.md` → API Contract 8).
12. MUST regenerate, never hand-edit, with the commands of `_authorization.md` → Sanctioned regeneration, in this order:
    - `make skills-sync`;
    - the `-record-skill-versions` step (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
    - the `-record-skill-coverage` step (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`).
13. MUST change no Governed Path outside `_authorization.md` `paths`. MUST NOT change any authoring skill, `.roundfixrc.yml`, `internal/cli/cli_test.go` or any file under `.claude/agents/`.
14. MUST NOT write a test that depends on a release tag, a remote ref, the real home, a real harness or the network.

## Subtasks

- [ ] Add the limit reader with its tests.
- [ ] Add the fallback entry, its reader and the skip by `Phases`, with their tests.
- [ ] Add the two subcommands, the helps and the command tests, and update the help test.
- [ ] Describe them in the command guide, the Phase Record section and the phase reference.
- [ ] Add the coverage entries, the glossary terms and the ADR-0407 notice, and regenerate.

## Acceptance Criteria

- [ ] Every message of the harness reference is read as the TechSpec's invariants say, and none names a model in the code.
- [ ] `roundfix phase fallback` writes exactly one fallback entry and prints the start line when a fallback applies, and writes nothing on every refusal and usage error.
- [ ] `roundfix phase model` names the fallback only for a later phase of the same Spec, in the same repository, whose declared model is the limited one; it writes nothing.
- [ ] `roundfix phase start` behaves as on its base.
- [ ] Surface Transcripts 1 to 8 pass line by line.
- [ ] The guide, the reference, the coverage map, the glossary and ADR-0407 describe the change, and the documentation contracts pass.

## Context

- creates: `internal/phaserecord/limit.go`
- creates: `internal/phaserecord/limit_test.go`
- creates: `internal/phaserecord/fallback.go`
- creates: `internal/phaserecord/fallback_test.go`
- interface: `internal/phaserecord/read.go`
- interface: `internal/cli/phase.go`
- interface: `internal/cli/phase_test.go`
- creates: `internal/cli/phase_fallback_test.go`
- interface: `internal/cli/cli.go`
- interface: `docs/user-guide/commands/phase.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/references/phase.md`
- interface: `skills/roundfix/references/phase.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md`
- instruction: `internal/cli/authoring_models.go`
- instruction: `docs/specs/0405-authoring-model-fallback/references/2026-11-17-harness-and-test-measurements.md`

## Consumers

- The Phase Record line format and `phaserecord.Phases` — consumers: `./internal/phaserecord`, `./internal/cli`
- The help of `roundfix phase` and the root help — consumers: `./internal/cli`
- The command guide, the phase reference, the coverage map and the glossary — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReadUsageLimitClassifiesTheHarnessMessages|TestAFallbackEntryIsOneLineOfItsDeclaredKeys|TestFallbacksReadsBackTheFallbackEntries|TestPhasesSkipsAFallbackEntry|TestPhaseFallbackSurfaceTranscripts|TestPhaseFallbackRefusesWithoutWriting|TestPhaseModelKeepsTheFallbackForItsSpecAndModel)$' ./internal/phaserecord ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReadUsageLimitClassifiesTheHarnessMessages TestAFallbackEntryIsOneLineOfItsDeclaredKeys TestFallbacksReadsBackTheFallbackEntries TestPhasesSkipsAFallbackEntry TestPhaseFallbackSurfaceTranscripts TestPhaseFallbackRefusesWithoutWriting TestPhaseModelKeepsTheFallbackForItsSpecAndModel; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/phaserecord ./internal/cli` — expected: exit 0. Before this Task the reader, the entry, the subcommands and the named tests do not exist, so the command fails. Afterwards every named test passes, and so do both packages that read and write the Phase Record and print the help.
- `tr -s '[:space:]' ' ' < docs/adr/0407-a-phase-record-names-the-spec-and-phase-of-supervisor-work-and-a-phase-boundary-ends-the-session.md | grep -qF -- "Amended by ADR-0428" || { printf 'missing notice in ADR-0407\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/phase.md | grep -qF -- "roundfix phase fallback" || { printf 'missing command in the guide\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/phase.md | grep -qF -- "roundfix phase model" || { printf 'missing command in the phase reference\n' >&2; exit 1; }; grep -F -- '**Authoring Fallback**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Authoring Fallback\*\*(:| —( |$))' && grep -F -- '**Phase Record**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Phase Record\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "fallback entries" && go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task ADR-0407 carries no notice of ADR-0428, and neither the guide, the reference nor the glossary names the commands or the term. Afterwards they all do, and the skill tests, the documentation contracts, the mirror and the records pass.

## References

- `_prd.md` → Goals (first to fourth bullets); User Stories 2, 4, 5 and 6; Core Features 2, 4, 5, 6 and 7; Success Metrics 1, 2 and 3
- `_techspec.md` → Limit reader; Fallback entries; Phase subcommands; Interfaces (invariants 1 to 8); Data Models (Fallback entry); API Contracts 1, 2, 3, 4, 7 and 8; Surface Transcript 1, Surface Transcript 2, Surface Transcript 3, Surface Transcript 4, Surface Transcript 5, Surface Transcript 6, Surface Transcript 7, Surface Transcript 8; Testing Approach; Build Order 2
- [`references/2026-11-17-harness-and-test-measurements.md`](references/2026-11-17-harness-and-test-measurements.md)
- ADR-0429, ADR-0428, ADR-0427, ADR-0407, ADR-0414, ADR-0256

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, systematic-debugging, no-workarounds.

Implemented this Task's slice for Daemon Verification:

- Added the model-independent limit reader, with the first matching line retained, either apostrophe and letter case accepted, shared and multi-word limits classified separately, controls replaced by spaces and output bounded to 200 runes.
- Added the ordered `fallback` JSON line, private UTC day-file append and oldest-first reader. The existing `Phases` reader already skips non-phase kinds before decoding their payload; its behavior is preserved and covered for fallback entries, including an invalid phase payload.
- Added `phase fallback` with the seven ordered refusals, usage/configuration diagnostics, byte limit, append-before-output behavior and runtime write errors. Added read-only `phase model`, selecting the oldest applicable entry by Spec, repository identity (or Git root when either identity is empty), and the current phase's declared model and fallback.
- Preserved the `phase start` implementation and its output/record expectations. Updated only the help literals in its original transcript test and the Spec 0404 full-help test found on this base in `internal/cli/phase_model_test.go`, as Requirement 6 requests.
- Updated both helps, the command guide, the Phase Record guide, the Roundfix skill and its phase reference, the authored coverage entries and dated review of `command phase start`, the Authoring Fallback and Phase Record glossary definitions through domain-modeling, and the ADR-0407 amendment notice. The sanctioned recorder raised both Roundfix skill version fields from `0.1.268` to `0.1.269`, then to `0.1.270` for the feedback repair below, and recorded the digests; the mirror and behavior record were generated.

### Focused evidence by acceptance criterion

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Harness messages classified without model names in production reader code | `TestReadUsageLimitClassifiesTheHarnessMessages` covers all quoted usage messages, the Spec 0404 observation, straight/curly apostrophes, upper/lower case, multi-line order, control characters and a 300-character Unicode line. Initial focused compilation failed on the missing reader API; the implemented reader's focused run passed 21 tests. |
| One fallback entry on success; no day-file writes on refusal or usage error | `TestAFallbackEntryIsOneLineOfItsDeclaredKeys` checks every serialized key and its order, UTC date/instant, append preservation, empty values, JSON-escaped newlines, private modes and cancellation. `TestFallbacksReadsBackTheFallbackEntries` checks invalid JSON/schema/phase/time, other entry kinds, oldest-first ordering and read-only behavior. `TestPhaseFallbackRefusesWithoutWriting` checks all seven ordered refusals, usage/configuration errors, absent record creation and preservation of existing day files. Runtime cases cover write failure, a 4,096-byte accepted error, oversized multibyte errors, recorded repository identity, no effort/session, and bounded output. All are included in the final focused run below. |
| Later phases use only the same Spec/repository/model fallback; model writes nothing | `TestPhaseModelKeepsTheFallbackForItsSpecAndModel` covers two Specs, distinct committed repositories, a clone sharing identity, oldest applicable entry, different declared model, no fallback/model, fallback in the limited family, case-insensitive substring matching, and both empty-identity path rules. Read-only checks compare day files byte for byte, and an additional case checks use without a record directory. Included in the final focused run. |
| phase start retains its base behavior | The final focused run includes the existing `TestPhaseStart*` tests and `TestPhaseSurfaceTranscript6`; only the two help literals changed. `TestPhasesSkipsAFallbackEntry` preserves phase-reader results and skipped count. |
| Surface Transcripts 1–8 | `TestPhaseFallbackSurfaceTranscripts/transcript_1` through `transcript_8` compare complete stdout/stderr and exit codes. Transcripts 3–5 run after transcript 2 in the same temporary home, and every command except transcript 2 is checked for unchanged record files. Included in the final focused run. |
| Guides, reference, coverage, glossary and ADR describe the change | Sanctioned generation succeeded in the required order. The focused documentation contract selection passed 48 tests, including command coverage, guide/skill command names, glossary rules and accepted ADR links. The existing phase-reference selection passed 19 tests. Full documentation Verification remains with the Daemon. |

Focused commands and outcomes:

- `go test ./internal/phaserecord -run '^TestReadUsageLimitClassifiesTheHarnessMessages$/curly$' -count=1` — initial red signal: missing `LimitKind`, `UsageLimit` and `ReadUsageLimit`.
- `go test ./internal/phaserecord -run '^TestReadUsageLimitClassifiesTheHarnessMessages$' -count=1` — passed, 21 tests.
- `go test ./internal/phaserecord ./internal/cli -run 'Test(ReadUsageLimit|AFallbackEntry|FallbacksReads|PhasesSkips|PhaseFallback|PhaseModelKeeps|PhaseStart|PhaseSurfaceTranscript6)' -count=1 -timeout 5m` — final run passed, 224 tests across both packages.
- `make skills-sync` — succeeded.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — succeeded; generated the version/digest record and both version fields.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — succeeded; generated the behavior-surface record.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^(TestEveryCommandIsNamedInTheRoundfixSkill|TestEveryCommandIsNamedInTheUserGuide|TestAnUndocumentedCommandIsReported|TestCommandCoverageFollowsTheRoundfixReferenceIndex|TestSkillCoveragePathsRequireOwnedFilesAndMatchingSources|TestGlossary.*|TestCommandReferenceLinksResolve)$'` — passed, 48 tests.
- `go test ./skills -run '^TestThePhaseReferenceDocumentsTheCommandAndTheRecord$' -count=1` — passed, 19 tests.
- `git diff --check` — passed.

One new diagnostic expectation initially used `load config`; the saved failing transcript and existing config-loader tests showed the established diagnostic is `parse config`. Corrected the new expectation and reran it successfully, then reran the focused implementation selection. No existing config diagnostic was changed.

Scope: the Task file's `in_progress` status was pre-existing and remains untouched. No Task Graph, other Task file, authoring skill, `.roundfixrc.yml`, `internal/cli/cli_test.go`, or `.claude/agents/` file was edited. The Spec 0404 help-expectation file above is the only ordinary changed path beyond the listed interfaces and is required by Requirement 6. No commit, push, PR or Spec announcement was made. The Task's declared Verification commands and the repository Verification commands were not run; the Daemon owns them and Task settlement.

### Verification feedback repair — attempt 1

Inspected the Daemon diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261118T111508Z_d52e529984ea1f74/verification/batch-002-attempt-1.log`
and the owning `TestRoundfixSkillEntryFileStaysWithinItsBudget` test. The
added phase-command names put the Roundfix skill entry over its existing
20,000-byte budget. Shortened its introductory prose and the phase row's
purpose text while preserving both full command names and the reference's
fallback workflow. No budget, assertion or verification configuration changed.

Fresh focused evidence:

- `make skills-sync` — succeeded.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — succeeded; regenerated both version fields as `0.1.270` and recorded the new digest.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — succeeded after the version recorder.
- `go test ./skills -run '^(TestRoundfixSkillEntryFileStaysWithinItsBudget|TestRoundfixSkillIndexLinksResolve|TestThePhaseReferenceDocumentsTheCommandAndTheRecord)$' -count=1` — passed, 25 tests, including both canonical and embedded budget checks.
- A read-only byte comparison confirms the canonical and mirrored skill entry and phase reference match; each skill entry is 19,989 bytes. Both version fields are `0.1.270`.
- `git diff --check` — passed after the feedback repair and Result update.

This repair changes only the canonical Roundfix skill entry, its generated
mirror and version record, and this Result. The sanctioned coverage recorder
also ran and left its current record unchanged. Task status remains
Daemon-owned. The declared Verification commands were not rerun; the Daemon
owns the next full configured Verification sequence.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/phase_model_test.go`

## Carry-forward provenance

- Source Run: `run_20261118T111508Z_d52e529984ea1f74`
- Source commit: `c8b810e945b3989c8c352d40ea2a0f84272ac670`
