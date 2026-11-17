---
task: task_01
spec: 0405-authoring-model-fallback
status: pending
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
