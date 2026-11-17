---
task: task_02
spec: 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
status: pending
type: docs
complexity: medium
---

# Task 02: This repository's declaration and its guides follow the Recommended Authoring Declaration

## Overview

The recommendation task_01 adds gets its drift gate and its words. A
repository contract test holds this repository's `.roundfixrc.yml` and the
configuration guide's block equal to the Recommended Authoring
Declaration. The configuration guide shows that block instead of an
unmarked example; the Roundfix skill's setup reference and `--update` mode
describe the `undeclared` report; the glossary adds the term and revises
**Authoring Model Declaration**; and ADR-0401 records that ADR-0415 amends
it.

## Requirements

1. MUST add `internal/docscontract/authoring_recommendation_test.go`,
   under the `docscontract` build tag and `//verify:always`, with
   `TestTheRecommendedAuthoringDeclarationMatchesThisRepository` and the
   subtests `repository_config` and `configuration_guide`, holding
   invariant 6 of `_techspec.md` → Interfaces:
   - `repository_config` reads this repository's `.roundfixrc.yml` and
     fails, naming the key, when the model or the reasoning effort of
     `idea`, `prd`, `techspec`, `tasks` or `research` differs from
     `config.RecommendedAuthoringDeclaration()`; it ignores
     `research.tools`;
   - `configuration_guide` reads the first fenced `yaml` block under
     `## Authoring models` in `docs/user-guide/configuration.md` and fails
     unless its content equals `config.RecommendedAuthoringYAML()`.
   The test follows the Parallel Test Package rule of `internal/testfixture`.
2. MUST replace the example block of `## Authoring models` in
   `docs/user-guide/configuration.md` with the Recommended Authoring
   Declaration, replace "Roundfix ships no model name" there with the
   statement that the block is the recommendation Doctor and the Update
   Check print and that Roundfix applies none by itself, show
   `tools: [WebSearch, WebFetch]` under `research` as an option the
   repository adds, and say that a file that already holds `skills:` takes
   the lines under it. Then run the `-record-skill-coverage` step to
   re-record `docs/references/behavior-surfaces.json`.
3. MUST change `.agents/skills/roundfix/references/setup.md` so the
   authoring paragraph says that Roundfix applies no model name by itself,
   and that with installed authoring skills the repository does not publish
   and no declaration, `roundfix doctor` (`authoring-models: warn`) and
   `roundfix update check` (`authoringModels.state` `undeclared`) print the
   Recommended Authoring Declaration to add to `.roundfixrc.yml` before
   `roundfix skills install --target project`, and that neither fails or
   writes. Change `.agents/skills/roundfix/references/modes.md` so the
   `--update` mode, when `authoringModels.state` is `undeclared`, shows the
   check's `recommendation` and leaves `.roundfixrc.yml` to the user,
   never writing it (`_techspec.md` → API Contract 4). Keep every other
   sentence of both files, and change no `### QA settlement` section.
4. MUST run `make skills-sync` and the `-record-skill-versions` step, so the
   mirrors under `skills/roundfix/`, the `SKILL.md` pair and
   `skills/testdata/owned-skill-versions.json` follow. Hand-edit no mirror,
   version field or record.
5. MUST change the `setup_declares_models` phrase in
   `skills/authoring_delegation_test.go` to the new sentence of the setup
   reference, a declared expectation change of `_techspec.md` → Testing
   Approach.
6. MUST write the glossary through `domain-modeling`: add
   **Recommended Authoring Declaration** to `CONTEXT.md` beside
   **Authoring Model Declaration**, citing ADR-0415, and revise
   **Authoring Model Declaration** so it no longer says that Roundfix ships
   no model name and names the `undeclared` warning of Doctor and the
   Update Check.
7. MUST add to ADR-0401, after its ADR-0408 amendment paragraph, the
   paragraph `**Amended by [ADR-0415](0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md) (2026-11-17).**`
   followed by one or two sentences: Doctor and the Update Check print a
   Recommended Authoring Declaration when the authoring skills are
   installed without one, and Roundfix still applies no model name by
   itself. Change nothing else in ADR-0401.
8. MUST change no Governed Path outside `_authorization.md` `paths` and no
   file outside the ones this Task declares; `.roundfixrc.yml` stays
   unchanged.

## Subtasks

- [ ] Add the repository contract test for the configuration and the guide.
- [ ] Show the recommendation in the configuration guide and re-record the Behavior Surface Record.
- [ ] Describe the report in the Roundfix skill's setup reference and `--update` mode, synchronize and record the version.
- [ ] Update the skills test phrase.
- [ ] Write the glossary entries and ADR-0401's amendment line.

## Acceptance Criteria

- [ ] Changing one model or effort of this repository's declaration, or one line of the configuration guide's block, fails the named subtest.
- [ ] The configuration guide, the setup reference, the `--update` mode and the glossary describe the Recommended Authoring Declaration and the `undeclared` report.
- [ ] The mirrors, the `SKILL.md` pair and the version record follow the canonical copies, and every skills test passes.
- [ ] ADR-0401 names its amendment by ADR-0415 and keeps the rest of its text.

## Context

- instruction: `docs/adr/0415-doctor-and-the-update-check-recommend-an-authoring-model-declaration-when-none-is-declared.md`
- instruction: `internal/docscontract/model_selection_test.go`
- creates: `internal/docscontract/authoring_recommendation_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/authoring_delegation_test.go`
- interface: `CONTEXT.md`
- interface: `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`

## Consumers

- The Roundfix skill's setup reference and `--update` mode text — consumers: `./skills`
- The configuration guide's authoring block — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 30m -tags docscontract -v -run '^TestTheRecommendedAuthoringDeclarationMatchesThisRepository$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheRecommendedAuthoringDeclarationMatchesThisRepository/repository_config TestTheRecommendedAuthoringDeclarationMatchesThisRepository/configuration_guide; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the test does not exist, so the command fails on the first missing pass; afterwards both subtests pass and the documentation contracts pass whole.
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- "Recommended Authoring Declaration" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/setup.md "Recommended Authoring Declaration" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/modes.md | grep -qF -- "undeclared" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/modes.md "undeclared" >&2; exit 1; }; go test -count=1 -timeout 30m ./skills` — expected: exit 0. Before this Task neither reference names the report, so the command fails; afterwards both do and the skills package, the mirror, version and phrase tests included, passes.
- `grep -F -- '**Recommended Authoring Declaration**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Recommended Authoring Declaration\*\*(:| —( |$))' && grep -F -- '**Authoring Model Declaration**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Authoring Model Declaration\*\*(:| —( |$))'` — expected: exit 0. Before this Task the glossary has no **Recommended Authoring Declaration** line, so the command fails; afterwards both entries are definitions.
- `tr -s '[:space:]' ' ' < docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md | grep -qF -- "Amended by [ADR-0415]" || { printf 'missing phrase in %s: %s\n' docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md "Amended by [ADR-0415]" >&2; exit 1; }` — expected: exit 0. Before this Task ADR-0401 names no amendment by ADR-0415.

## References

- `_prd.md` → Goal 4; Core Features 3 and 4; Success Metric 4; Glossary
- `_techspec.md` → Contract and guides; Interfaces (invariant 6); API Contract 4; Testing Approach (Task 2, contract and skill); ADR Alternatives (ADR-0401); Build Order 2
- ADR-0415, ADR-0401, ADR-0402
