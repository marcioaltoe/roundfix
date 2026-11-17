---
task: task_02
spec: 0402-doctor-and-the-update-check-recommend-an-authoring-model-declaration
status: completed
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
8. MUST change this repository's `.roundfixrc.yml` only in
   `skills.authoring.claude`: `prd` and `techspec` to `fable`/`high` and
   `tasks` to `opus`/`high`, the values of `_techspec.md` → Recommendation;
   `idea` and `research`, its `tools` included, and every other line stay
   as they are. Then run the sanctioned
   `go run -buildvcs=false ./cmd/roundfix skills install --target project`,
   which regenerates `.claude/agents/roundfix-prd-author.md`,
   `.claude/agents/roundfix-techspec-author.md` and
   `.claude/agents/roundfix-tasks-author.md`, one line each, and writes no
   header into the four authoring skills, because this repository publishes
   them (ADR-0408). Keep any other file the installation would change out
   of the commit and report it. Hand-edit no generated subagent.
9. MUST change no Governed Path outside `_authorization.md` `paths` and no
   file outside the ones this Task declares.

## Subtasks

- [ ] Set this repository's authoring declaration to the recommendation and regenerate its three Author Subagents.
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
- [ ] This repository's declaration and its generated Author Subagents carry the recommended values, and `roundfix skills check` reports them applied.

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
- interface: `.roundfixrc.yml`
- interface: `.claude/agents/roundfix-prd-author.md`
- interface: `.claude/agents/roundfix-techspec-author.md`
- interface: `.claude/agents/roundfix-tasks-author.md`
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

- `tr -s '[:space:]' ' ' < .claude/agents/roundfix-tasks-author.md | grep -qF -- "model: opus" || { printf 'missing phrase in %s: %s\n' .claude/agents/roundfix-tasks-author.md "model: opus" >&2; exit 1; }; tr -s '[:space:]' ' ' < .claude/agents/roundfix-prd-author.md | grep -qF -- "effort: high" || { printf 'missing phrase in %s: %s\n' .claude/agents/roundfix-prd-author.md "effort: high" >&2; exit 1; }; tr -s '[:space:]' ' ' < .claude/agents/roundfix-techspec-author.md | grep -qF -- "effort: high" || { printf 'missing phrase in %s: %s\n' .claude/agents/roundfix-techspec-author.md "effort: high" >&2; exit 1; }` — expected: exit 0. Before this Task the tasks author runs on `fable` and the PRD and TechSpec authors on `xhigh`, so the command fails; afterwards the regenerated subagents carry the recommended values.

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, tech-writer, domain-modeling, roundfix, evidence-gate

Implementation:

- Added the `docscontract` drift test with `repository_config` and `configuration_guide` subtests. The repository comparison names each missing or mismatched model and reasoning-effort key and intentionally ignores `research.tools`; the guide comparison requires the first YAML fence under `## Authoring models` to equal `config.RecommendedAuthoringYAML()`.
- Changed this repository's `skills.authoring.claude` declaration to the recommendation, updated the guide, setup reference, `--update` mode, glossary and ADR-0401, and updated the declared skills phrase test.
- Ran the sanctioned generators. The three Author Subagents now carry `fable/high`, `fable/high` and `opus/high`; the canonical and mirror Roundfix skills are synchronized at version `0.1.257`, the owned-skill version record and Behavior Surface Record are regenerated, and no other installer output changed in the worktree.
- Repaired the repository-contract expectation exposed by the Daemon gate: the required new docscontract file makes the always-selected docscontract count 10 instead of 9. The existing contract's expected count was updated without changing gate selection logic.

Focused checks:

- `gofmt -w internal/docscontract/authoring_recommendation_test.go`: passed.
- `git diff --check`: passed.
- `go test -count=1 -tags docscontract -v ./internal/docscontract -run 'TestTheRecommendedAuthoringDeclarationMatchesThisRepository/(repository_config|configuration_guide)'`: passed, 3 tests including both named subtests.
- `go test -count=1 ./internal/config -run '^TestRecommendedAuthoringDeclarationRendersValidProjectConfig$'`: passed, 4 tests.
- `go test -count=1 ./skills -run 'TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase|TestTheAuthoringSkillsDelegateADeclaredPhase'`: passed, 172 tests.
- `make skills-sync`: exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed, 1 test.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed, 1 test and regenerated `docs/references/behavior-surfaces.json`.
- `go run -buildvcs=false ./cmd/roundfix skills install --target project`: exited 0 and regenerated only the three declared Author Subagents among changed generated files.
- `go run -buildvcs=false ./cmd/roundfix skills check`: exited 0 with `Roundfix skill check passed`.
- `go test -count=1 ./internal/verifyselect -run '^TestRepositoryContractTestsDeclareTheirRelevance$'`: passed, 4 tests.

Acceptance evidence:

- The contract test observes a changed declaration or guide line through the respective named subtest and reports the affected key or exact block mismatch.
- The repository-contract declaration inventory now accounts for the required new docscontract file; its focused gate test passes.
- The configuration guide, setup reference, `--update` mode and glossary name the Recommended Authoring Declaration and the `undeclared` report; the optional research tools are shown separately so the governed recommendation remains tool-free.
- `make skills-sync`, the version recorder, mirror comparisons and the focused skills tests provide current evidence for synchronized skills and records.
- ADR-0401 contains the required ADR-0415 amendment immediately after the ADR-0408 amendment, without other ADR-0401 edits.
- The declaration, generated Author Subagents and local `skills check` output carry the recommended values and applied state.

The Task's declared Verification commands and the repository Verification were not run; they remain Daemon-owned. No commit, push or pull request was made.

## References

- `_prd.md` → Goals 4 and 5; Core Features 3 and 4; Success Metric 4; Glossary
- `_techspec.md` → Contract and guides; Interfaces (invariant 6); API Contract 4; Testing Approach (Task 2, contract and skill); ADR Alternatives (ADR-0401); Build Order 2
- ADR-0415, ADR-0401, ADR-0402

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/verifyselect/contracts_repository_test.go`

## Carry-forward provenance

- Source Run: `run_20261117T222450Z_a00f2c8f9217d395`
- Source commit: `f2ea52308c584184440d174fd9236250bd2e00e4`
