---
task: task_01
spec: 0405-authoring-model-fallback
status: completed
type: backend
complexity: medium
---

# Task 01: The Authoring Model Declaration accepts a fallback for each Authoring Phase, and the recommendation carries one

## Overview

An entry `idea`, `prd`, `techspec` or `tasks` of the Authoring Model
Declaration accepts `fallback` with a `model` and an optional
`reasoning_effort`. The research entry refuses it, and so does every case in
which the fallback escapes no limit. The configuration guide, the Roundfix
skill's setup reference and the glossary describe the key. Installation
ignores it, so no header, subagent file or check changes. The Recommended Authoring Declaration that the Doctor Command and the Update Check print gains the maintainer's fallbacks, and its repository contract compares a fallback only where this repository declares one.

## Requirements

1. MUST add `AuthoringFallback` and the `Fallback` field of `_techspec.md` → Interfaces to the entries of the four Authoring Phases. MUST validate them in the order and with the exact messages of `_techspec.md` → Data Models, Project Config.
2. MUST change the refusal of an unknown entry key to `expected model, reasoning_effort or fallback`. MUST update the existing case `TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/phase_field` to that text, and change no other expectation of that test.
3. MUST refuse `fallback` under `research` with `skills.authoring.claude.research.fallback is not a supported key; expected model, reasoning_effort or tools`.
4. MUST compare a fallback model with its entry's model ignoring letter case.
5. MUST leave the installation's declaration (`authoringDeclaration`) and the generated headers and subagent files unchanged. A repository whose declaration gains a fallback MUST install byte-identical skills and subagents.
6. MUST add tests in a new file, each calling `t.Parallel()`:
   - `TestAnAuthoringEntryAcceptsAFallback` covers a fallback with and without an effort on each of the four phases, the decoded values, and an install that writes the same files with and without the fallback.
   - `TestAFallbackIsRefusedWhereItEscapesNoLimit` is a table with one case per refusal of Data Models plus the research case, asserting each full message.
7. MUST describe the key in `docs/user-guide/configuration.md`, in the section on `skills.authoring.claude`. The description MUST say:
   - its two fields;
   - that the model is a name the harness accepts when a subagent is started, an alias in Claude Code;
   - the sentence `a fallback must name another model`;
   - that a Roundfix older than the release refuses a Project Config that holds it.
8. MUST describe the same in the paragraph of `.agents/skills/roundfix/references/setup.md` that describes the declaration, with the phrase `a fallback model and reasoning effort`.
9. MUST revise **Authoring Model Declaration** in `CONTEXT.md`, through domain-modeling, so that it says each Authoring Phase may declare a fallback model and reasoning effort, with the phrase `may declare a fallback`.
10. MUST regenerate, never hand-edit, with the commands of `_authorization.md` → Sanctioned regeneration, in this order:
    - `make skills-sync`;
    - the `-record-skill-versions` step (`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`);
    - the `-record-skill-coverage` step (`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`), which re-records the `config keys` and configuration-guide surfaces.
11. MUST add the fallbacks of `_techspec.md` → API Contract 9 to the Recommended Authoring Declaration and its rendering, and leave every model and effort Spec 0402 set unchanged.
    - MUST change Spec 0402's tests that hold the rendered block to the block of Surface Transcript 9, and no other expectation of theirs: `TestRecommendedAuthoringDeclarationRendersValidProjectConfig`, `TestDoctorRecommendsAnAuthoringDeclarationWhenNoneIsDeclared`, `TestUpdateCheckRecommendsAnAuthoringDeclarationWhenNoneIsDeclared`, `TestSurfaceTranscript1UpdateCheckRecommendsTheDeclaration` and `TestRunDoctorRealRepositoryCheckDoesNotMutateState`.
    - MUST add `TestSurfaceTranscript9UpdateCheckRecommendsTheFallbacks`, calling `t.Parallel()`, asserting Surface Transcript 9 line by line through `runCLI`.
    - MUST change `TestTheRecommendedAuthoringDeclarationMatchesThisRepository` so that it holds the contract of API Contract 9. It compares the guide's block in full, compares this repository's declared fallbacks only where declared, and has a subtest `a_declared_fallback_that_differs_fails` on a fixture copy of `.roundfixrc.yml`.
    - MUST make the configuration guide's block the new rendering.
    - MUST revise the glossary definition of the Recommended Authoring Declaration that Spec 0402 adds, so that it names the fallbacks.
    - MUST add to ADR-0415 the notice `Amended by ADR-0430` with the date and what changed.
12. MUST change no Governed Path outside `_authorization.md` `paths`. MUST NOT change `.roundfixrc.yml`, any file under `.claude/agents/`, `internal/cli/cli_test.go` or any authoring skill.
13. MUST NOT write a test that depends on a release tag, a remote ref, the real home or the network.

## Subtasks

- [ ] Add the type, the field and the validation, with the two new tests.
- [ ] Update the existing refusal case.
- [ ] Describe the key in the configuration guide and the setup reference.
- [ ] Add the fallbacks to the recommendation, its tests, its transcript and the contract.
- [ ] Revise the glossary terms and regenerate the mirrors and records.

## Acceptance Criteria

- [ ] Every entry of the four phases accepts a fallback with a model, with or without an effort, and the research entry refuses one.
- [ ] Each refusal of Data Models prints its exact message, and a fallback that names its entry's model in any letter case is refused.
- [ ] Installing with and without a fallback writes byte-identical skills and subagents.
- [ ] The configuration guide, the setup reference and the glossary describe the key.
- [ ] The Doctor Command and the Update Check recommend the fallbacks, the guide shows the same block, and the contract fails on a declared fallback that differs and passes on this repository without fallbacks.

## Context

- interface: `internal/config/authoring.go`
- interface: `internal/config/authoring_test.go`
- creates: `internal/config/authoring_fallback_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- creates: `internal/config/authoring_recommendation.go`
- creates: `internal/config/authoring_recommendation_test.go`
- creates: `internal/cli/authoring_recommendation_test.go`
- interface: `internal/cli/doctor_test.go`
- creates: `internal/docscontract/authoring_recommendation_test.go`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md`
- instruction: `internal/cli/authoring_models.go`

## Consumers

- The Project Config schema of the Authoring Model Declaration — consumers: `./internal/config`, `./internal/cli`, `./skills`
- The Recommended Authoring Declaration and its rendered block — consumers: `./internal/config`, `./internal/cli`, `./internal/docscontract`
- The configuration guide, the setup reference and the glossary — consumers: `./internal/docscontract`, `./skills`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAnAuthoringEntryAcceptsAFallback|TestAFallbackIsRefusedWhereItEscapesNoLimit|TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow|TestSurfaceTranscript9UpdateCheckRecommendsTheFallbacks)$' ./internal/config ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAnAuthoringEntryAcceptsAFallback TestAFallbackIsRefusedWhereItEscapesNoLimit TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow TestSurfaceTranscript9UpdateCheckRecommendsTheFallbacks; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/config ./internal/cli ./skills` — expected: exit 0. Before this Task the two new tests do not exist and the refusal names no fallback, so the command fails. Afterwards the named tests pass, and so do the three packages that read the declaration.
- `tr -s '[:space:]' ' ' < docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md | grep -qF -- "Amended by ADR-0430" || { printf 'missing notice in ADR-0415\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "a fallback must name another model" || { printf 'missing phrase in %s\n' docs/user-guide/configuration.md >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- "a fallback model and reasoning effort" || { printf 'missing phrase in %s\n' .agents/skills/roundfix/references/setup.md >&2; exit 1; }; grep -F -- '**Authoring Model Declaration**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Authoring Model Declaration\*\*(:| —( |$))' && tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "may declare a fallback" && out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^TestTheRecommendedAuthoringDeclarationMatchesThisRepository$' ./internal/docscontract 2>&1)" && printf '%s\n' "$out" | grep -qF -- "--- PASS: TestTheRecommendedAuthoringDeclarationMatchesThisRepository/a_declared_fallback_that_differs_fails " && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract && make skills-sync-check` — expected: exit 0. Before this Task neither the guide, the reference nor the glossary names the fallback, and the contract has no fallback subtest. Afterwards they do, and the documentation contracts, the mirror and the records pass.

## References

- `_prd.md` → Goals (first and fifth bullets); User Stories 1 and 7; Core Features 1, 7 and 8; Success Metric 4; Success Metric 5
- `_techspec.md` → Fallback declaration; Recommendation; Interfaces; Data Models (Project Config); API Contract 5; API Contract 9; Surface Transcript 9; Build Order 1; Testing Approach
- [`references/2026-11-17-harness-and-test-measurements.md`](references/2026-11-17-harness-and-test-measurements.md)
- ADR-0427, ADR-0430, ADR-0415, ADR-0401, ADR-0408, ADR-0392

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, systematic-debugging, no-workarounds.

Added the optional fallback declaration and its ordered validation. The four Authoring Phases accept a required fallback model and optional effort; research refuses it. The recommendation carries `opus`/`xhigh` for idea, prd and techspec and `sonnet`/`high` for tasks, preserving every existing primary model and effort. Recommendation copies own their fallback values. The configuration, Doctor and Update Check guides show the same block; the setup reference, glossary terms and ADR-0415 notice describe the change.

Acceptance evidence from focused checks:

- All four phases decode fallbacks with and without an effort in `internal/config`'s `TestAnAuthoringEntryAcceptsAFallback`. Its companion in `internal/cli` installs through `runCLI` before and after adding each fallback and compares the complete installed `.agents` and `.claude` snapshots. Both tests call `t.Parallel()`.
- `TestAFallbackIsRefusedWhereItEscapesNoLimit` asserts full messages for research, non-mappings, unknown keys and a nested fallback, malformed/non-string values, missing models, equal models including mixed case, and refusal precedence. The existing unknown-key test changes only its `phase_field` expectation. All selected cases passed.
- Installation's `authoringDeclaration` and production installation code remain unchanged. The CLI installation comparisons passed for all four phases with and without fallback effort, including the generated subagents and skill links. Existing skills-install and Baseline-install declaration tests also passed.
- The configuration guide and setup reference describe both fields, the harness-start alias requirement, another-model rule and older-binary refusal; the two glossary definitions and ADR-0415 amendment were updated. Focused glossary contracts passed. The setup mirror and version fields were regenerated, never hand-edited.
- The recommendation rendering loads as Project Config and keeps independent fallback copies. Doctor/Update Check recommendation tests, the existing transcript, the real-repository Doctor immutability test and the new line-by-line `TestSurfaceTranscript9UpdateCheckRecommendsTheFallbacks` passed. The shared expected block updates the existing Doctor test without changing its other expectations.
- `TestTheRecommendedAuthoringDeclarationMatchesThisRepository` compares the guide block in full, accepts undeclared or matching fallbacks, and detects model and effort drift for every phase through `a_declared_fallback_that_differs_fails` on fixture copies of `.roundfixrc.yml`. The actual repository config remains without fallbacks and passed the same comparison.

Focused commands and outcomes (each exited 0 unless noted):

- Initial red signal: `rtk proxy go test ./internal/config -run '^TestAFallbackIsRefusedWhereItEscapesNoLimit$/scalar$' -count=1` exited 1 before implementation: `AuthoringEntry` had no `Fallback` field and `AuthoringFallback` was undefined.
- `rtk proxy go test ./internal/config ./internal/cli -run 'Test(AnAuthoringEntryAcceptsAFallback|AFallbackIsRefusedWhereItEscapesNoLimit|TheAuthoringModelDeclaration|UserConfigIgnoresTheAuthoringModelDeclaration|RecommendedAuthoring|SurfaceTranscript[19]UpdateCheck|DoctorRecommendsAnAuthoringDeclaration|UpdateCheckRecommendsAnAuthoringDeclaration|RunDoctorRealRepositoryCheckDoesNotMutateState|SkillsInstallAppliesTheAuthoringModelDeclaration|BaselineUpdateAppliesTheAuthoringModelDeclaration)' -count=1 -timeout=5m` — both packages passed against the regenerated bundle.
- `rtk proxy go test ./internal/docscontract -tags docscontract -run 'Test(TheRecommendedAuthoringDeclarationMatchesThisRepository|Glossary|TheSkillCoverageMapIsCurrent)' -count=1 -timeout=5m` — passed.
- After tightening the transcript's final-newline assertion: `rtk proxy go test ./internal/cli -run '^TestSurfaceTranscript9UpdateCheckRecommendsTheFallbacks$' -count=1` — passed.
- After final import formatting: `rtk proxy go test ./internal/docscontract -tags docscontract -run '^TestTheRecommendedAuthoringDeclarationMatchesThisRepository$' -count=1` — passed.

Sanctioned regeneration ran in the required order, all exiting 0:

1. `rtk proxy make skills-sync`.
2. `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` — recorded Roundfix `0.1.267` in both canonical and mirrored version fields and the owned-skill record.
3. `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — refreshed `config keys` and the changed configuration/Doctor/Update Check guide surfaces.

Scope inspection: only declared Task paths changed, plus this Result; the incoming `status: in_progress` remains untouched. `.roundfixrc.yml`, `.claude/agents/`, `internal/cli/cli_test.go`, the Task Graph and other Task files remain unchanged. `git diff --check` passed. The initial sandbox process error required escalation because the configured writable root contains a symlink; the work stayed in the same physical Run Worktree.

The declared Verification commands, full package suites, repository Verification and terminal Task settlement are left to the Daemon. No commit, push, pull request, memory message or Spec Event Announcement was made. Phase execution and authoring-skill behavior belong to the later Tasks.

### Verification feedback repair — attempt 1

The Daemon diagnostic artifact isolated the failure to the existing setup-declaration prose contract in `TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase`. The new fallback explanation had interrupted the original declaration sentence. Restored that sentence intact and placed the fallback explanation immediately after it, preserving all required fallback guidance. The owning test and its assertions remain unchanged; the repair changes only the canonical setup reference and sanctioned generated skill artifacts, plus this Result.

Focused reproduction and repair evidence:

- Before repair, `rtk proxy go test ./skills -run '^TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase$' -count=1` exited 1 for the canonical and embedded setup-declaration cases, reproducing the reported failure.
- After repair and regeneration, `rtk proxy go test ./skills -run 'Test(TheRoundfixSkillGathersDecisionsBeforeADelegatedPhase|EveryOwnedSkillVersionIsRecorded)' -count=1` exited 0. Both copies satisfy the original prose contract and the version record matches.
- `rtk proxy go test ./internal/cli -run 'Test(AnAuthoringEntryAcceptsAFallback|SurfaceTranscript9UpdateCheckRecommendsTheFallbacks)' -count=1` exited 0 against the regenerated bundle, retaining byte-identical installation behavior and the fallback recommendation transcript.
- Repeated the sanctioned generators in order: `rtk proxy make skills-sync`; `rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`; `rtk proxy go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`. Each exited 0. The version recorder now records Roundfix `0.1.268` in both version fields and the owned-skill history.
- Final scope inspection confirms only declared paths changed, the Roundfix mirror is byte-identical to its canonical source, and authored Task content and status remain unchanged. `git diff --check` exited 0.

The diagnostic artifact was read, not embedded. Declared Verification and terminal settlement remain with the Daemon; no commit, push or pull request was created.

## Carry-forward provenance

- Source Run: `run_20261118T111508Z_d52e529984ea1f74`
- Source commit: `151b7590d1e5f6f06a8b6c39ee1d1158eec2f752`
