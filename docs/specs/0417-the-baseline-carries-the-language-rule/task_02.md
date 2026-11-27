---
task: task_02
spec: 0417-the-baseline-carries-the-language-rule
status: pending
type: backend
complexity: high
---

# Task 02: A codebase that keeps its language declares a Code Language Exception with its source

## Overview

This Task delivers `_techspec.md` → Build Order 2: the optional decision `language.code_exception` of the new type `code-language-exception`, which holds a language and the source of the owner's decision and refuses a declaration without a source; the `present` form of a clause gate; the three code-language clauses gated on the decision's absence and one clause gated on its presence; the second token of the Reader Language line; and an interactive flow that never asks for the exception. It is verifiable on its own: a declared exception swaps the clauses and names its source, a declaration without a source is refused and writes nothing, and a repository that declares nothing keeps the English clauses.

## Requirements

1. MUST add `internal/baseline/code_language_exception.go` with the constants, the type `CodeLanguageException` and `normalizeCodeLanguageException(value any) (CodeLanguageException, error)` of `_techspec.md` → Interfaces. It refuses a value with the reasons of `_techspec.md` → API Contract 9, word for word and in that order of checks, and refuses a language or a source that `validateStructuredRenderText` refuses as a single line. `validateDecisionValue` calls it for the type, so `baseline plan`, `baseline update --decision` and a Decision Document refuse the same way and write nothing. A recorded exception always holds a non-empty language other than English, in any case, and a non-empty source (invariant 2).
2. MUST make the catalog accept the type `code-language-exception` in its type list, count it among the types with exactly one render binding, and, for an optional decision of this type, refuse a `default` or a `suggestion` instead of requiring a default (`_techspec.md` → API Contract 9).
3. MUST add the decision `language.code_exception` to `internal/baseline/assets/decisions.json` after `language.reader`, exactly as `_techspec.md` → Data Models gives it: version 1, optional, type `code-language-exception`, the summary written there, no default and no suggestion, and one effect `when {"present": true}` with one `renderBindings` entry for `guide.agent-instructions`, `template.guide.agent-instructions` and the token `language.code_exception`.
4. MUST implement the clause gate of `_techspec.md` → API Contract 10: `appliesWhen` holds `decision` and exactly one of `equals` and `present`; with `present` the value is a boolean, the decision is optional and of any type, and `clauseApplies` returns whether the decision is recorded compared with the gate, never reading a decision's default (invariant 5); with `equals`, validation and evaluation are unchanged. The Source Baseline transition names presence in the reason of a clause a `present` gate takes out.
5. MUST gate, in `internal/baseline/assets/modules/core.json`, `clause.core.code-names-in-english`, `clause.core.code-names-keep-untranslated-terms` and `clause.core.code-names-new-only` with `{"decision": "language.code_exception", "present": false}`, and add to `rule.core.code-and-reader-language` the clause `clause.core.code-names-follow-the-existing-codebase`, `mandatory`, gated with `present` true, whose guidance is the first block of `_techspec.md` → API Contract 11 word for word. A rendered guide holds either that clause or the three clauses gated on absence, never both and never neither (invariant 3). `clause.core.reader-language` stays ungated.
6. MUST append the token to the template line, so that it reads `The Reader Language is {{language.reader}}.{{language.code_exception}}`, list the token in the template's entry of `internal/baseline/assets/templates/index.json`, and render it as `_techspec.md` → API Contract 11 says: nothing while the decision is absent, and while it is recorded a line break and the second block of API Contract 11, with the language escaped as Markdown text and the source in a code span.
7. MUST list `language.code_exception` in the `entryDecisions` of each of the five built-in profiles directly after `language.reader`, and change nothing else in a profile by hand.
8. MUST make `promptBaselineDecisions` in `internal/cli/baseline_human.go` never prompt for `language.code_exception`: it carries a recorded value into the answers unchanged and otherwise skips the decision (`_techspec.md` → API Contract 12). No scripted answer of an existing interactive test changes for it.
9. MUST list the new clause with its force in `internal/baseline/clause_characterization_test.go` and make the helper `manifestInputDecisionsForProfile` of `internal/baseline/update_test.go` skip `language.code_exception`, which has no suggested value to record. No parity-corpus fixture changes.
10. MUST regenerate, in this order and by these commands only, as `_authorization.md` → Sanctioned regeneration declares them: `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`; `make baseline-digests`; `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text --decision language.reader=English`, which declares no exception for this repository; and last `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`. The record command writes the next version of `core`, never a number written by hand; no version line, digest, golden, guide, Setup Manifest field or record is hand-edited. Repeating a recorded answer is no change, so the refresh is expected to rewrite the Setup Manifest alone; whatever it writes is its output and is not edited.
11. MUST add, in a new file `internal/baseline/code_language_exception_test.go`:
    - `TestACodeLanguageExceptionSwapsTheCodeClauses` (Success Metric 3), with subtests `a_declared_exception_states_the_language_and_the_source`, which renders a guide holding the exception's sentence with the language and the source and `clause.core.code-names-follow-the-existing-codebase`, and none of the three clauses gated on absence; `a_declaration_without_a_source_is_refused`, which writes nothing; and `no_declaration_states_english`, which renders the three clauses and not the exception's clause;
    - `TestACodeLanguageExceptionValueIsRefused`, a table with subtests `not_an_object`, `an_unknown_field`, `no_source`, `a_blank_source`, `no_language`, `english_in_any_case` and `a_multiline_source`, each asserting its reason of API Contract 9;
    - `TestAClauseGateOnPresence`, with subtests `present_true_applies_while_recorded`, `present_false_applies_while_absent`, `a_default_is_not_presence`, `equals_is_unchanged`, `both_forms_are_refused` and `a_non_boolean_present_is_refused`;
    - `TestTheCatalogRefusesADefaultForACodeLanguageException`, with subtests `a_default` and `a_suggestion`.
12. MUST add, in a new file `internal/cli/baseline_code_language_exception_test.go`:
    - `TestAnExceptionWithoutASourceIsRefused`, which asserts `_techspec.md` → Surface Transcript 2 whole: every line of standard output and standard error in order and the exit code `2`, with only `<repository>` and `<digest>` as controlled variations, and that no file of the repository changed;
    - `TestTheInteractiveFlowAsksNothingForTheCodeLanguageException`, with subtests `nothing_recorded` and `a_recorded_exception_is_carried`, asserting that no prompt names the decision and that a recorded value reaches the answers unchanged.
13. MUST add no top-level test function other than the six this Task names, and rename, merge or remove no existing test; a table case is a subtest. No new test file carries a build tag, and no file is created under `internal/docscontract/`. No new test reads the network, the real home directory, another repository or the clock.
14. MUST NOT change the clause texts of Task 01, the review prompt, any skill, the glossary, ADR 0454, any user-guide page, `.secondbrain-export`, any file under `internal/baseline/assets/source-baselines/` or `internal/baseline/assets/setups/`, the Makefile, the workflows, `go.mod`, `go.sum`, lint configuration, `.roundfixrc.yml`, `docs/references/gate-time-targets.json` or any file that measures the Gate Time Rule.

## Subtasks

- [ ] Add the exception's type, its normalization and refusals, and the catalog's acceptance of the type.
- [ ] Add the `present` form of the clause gate and the decision to the catalog.
- [ ] Gate the three clauses, add the existing-codebase clause, the template token and the profile entries.
- [ ] Keep the interactive flow from asking for the exception, and update the clause list and the manifest-input helper.
- [ ] Run the module record command, `make baseline-digests`, the Managed Refresh and the coverage record, and add the six named tests.

## Acceptance Criteria

- [ ] A repository that declares a Code Language Exception with a language and a source renders a guide that states the declared language for code with its source and the existing-codebase clause, and holds none of the three clauses gated on absence.
- [ ] A declaration without a source, with a field other than the two, with a blank or multiline value, or with English as its language is refused with its reason and writes nothing; Surface Transcript 2 is reproduced line by line.
- [ ] A repository that declares nothing renders the English clauses, and a `present` gate never reads a default.
- [ ] The interactive flow never asks for the exception and carries a recorded one unchanged.
- [ ] The module version, the digests, the goldens and this repository's Setup Manifest are the output of the sanctioned commands.

## Context

- instruction: `internal/cli/baseline_human_test.go`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/custom_profile.go`
- interface: `internal/baseline/project_decision_render.go`
- interface: `internal/baseline/assets/decisions.json`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/profiles/go-cli-tui.json`
- interface: `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`
- interface: `internal/baseline/assets/profiles/python.json`
- interface: `internal/baseline/assets/profiles/rust-cli.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/templates/guides/agent-instructions.md`
- interface: `internal/baseline/assets/templates/index.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/domain.md`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/unsatisfied-blocking-capabilities.golden.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/update_test.go`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/cli/baseline_human.go`
- interface: `docs/references/coverage-record.json`
- creates: `internal/baseline/code_language_exception.go`
- creates: `internal/baseline/code_language_exception_test.go`
- creates: `internal/cli/baseline_code_language_exception_test.go`

## Consumers

- The decision type `code-language-exception`, the `present` clause gate, the `core` module, the decision catalog, the five built-in profiles, the agent-instructions template and the derived catalog files — consumers: `./internal/baseline`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestACodeLanguageExceptionSwapsTheCodeClauses|TestACodeLanguageExceptionValueIsRefused|TestAClauseGateOnPresence|TestTheCatalogRefusesADefaultForACodeLanguageException|TestEveryBuiltInProfileRendersTheLanguageRule|TestBaselineClauseForceIsCharacterized|TestEveryBaselineModuleVersionIsRecorded|TestFormatterComposition|TestCatalogCompatibility)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestACodeLanguageExceptionSwapsTheCodeClauses/a_declared_exception_states_the_language_and_the_source TestACodeLanguageExceptionSwapsTheCodeClauses/a_declaration_without_a_source_is_refused TestACodeLanguageExceptionSwapsTheCodeClauses/no_declaration_states_english TestACodeLanguageExceptionValueIsRefused/not_an_object TestACodeLanguageExceptionValueIsRefused/an_unknown_field TestACodeLanguageExceptionValueIsRefused/no_source TestACodeLanguageExceptionValueIsRefused/a_blank_source TestACodeLanguageExceptionValueIsRefused/no_language TestACodeLanguageExceptionValueIsRefused/english_in_any_case TestACodeLanguageExceptionValueIsRefused/a_multiline_source TestAClauseGateOnPresence/present_true_applies_while_recorded TestAClauseGateOnPresence/present_false_applies_while_absent TestAClauseGateOnPresence/a_default_is_not_presence TestAClauseGateOnPresence/equals_is_unchanged TestAClauseGateOnPresence/both_forms_are_refused TestAClauseGateOnPresence/a_non_boolean_present_is_refused TestTheCatalogRefusesADefaultForACodeLanguageException/a_default TestTheCatalogRefusesADefaultForACodeLanguageException/a_suggestion TestEveryBuiltInProfileRendersTheLanguageRule TestBaselineClauseForceIsCharacterized TestEveryBaselineModuleVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; names="$(sed -n 's/^func \(Test[A-Za-z0-9_]*\)(t \*testing\.T) {$/\1/p' internal/baseline/plan_test.go internal/baseline/update_test.go internal/baseline/clause_characterization_test.go | paste -sd '|' -)"; test -n "$names" || exit 1; go test -count=1 -timeout 20m -run "^($names)$" ./internal/baseline` — expected: exit 0. Before this Task the four exception tests do not exist, so their pass lines are missing; afterwards the exception swaps the clauses, every refusal carries its reason, the gate never reads a default, the catalog refuses a default or a suggestion, every profile still renders the rule, and the tests of the changed Baseline test files pass.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAnExceptionWithoutASourceIsRefused|TestTheInteractiveFlowAsksNothingForTheCodeLanguageException|TestTheUpdateNamesTheReaderLanguageDecision|TestAnUpdateStopsForTheReaderLanguageAndThenChangesThreeFiles)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAnExceptionWithoutASourceIsRefused TestTheInteractiveFlowAsksNothingForTheCodeLanguageException/nothing_recorded TestTheInteractiveFlowAsksNothingForTheCodeLanguageException/a_recorded_exception_is_carried TestTheUpdateNamesTheReaderLanguageDecision TestAnUpdateStopsForTheReaderLanguageAndThenChangesThreeFiles; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; names="$(sed -n 's/^func \(Test[A-Za-z0-9_]*\)(t \*testing\.T) {$/\1/p' internal/cli/baseline_human_test.go internal/cli/baseline_release_gate_test.go internal/cli/baseline_profile_change_test.go internal/cli/baseline_update_test.go | paste -sd '|' -)"; test -n "$names" || exit 1; go test -count=1 -timeout 20m -run "^($names)$" ./internal/cli` — expected: exit 0. Before this Task the two CLI tests do not exist; afterwards Surface Transcript 2 is reproduced, the interactive flow asks nothing for the exception, Surface Transcript 1 still holds, and the interactive-flow and update tests pass unchanged.
- `grep -q 'func normalizeCodeLanguageException(value any) (CodeLanguageException, error)' internal/baseline/code_language_exception.go && grep -qF '{{language.reader}}.{{language.code_exception}}' internal/baseline/assets/templates/guides/agent-instructions.md && grep -q 'clause.core.code-names-follow-the-existing-codebase' internal/baseline/assets/modules/core.json && grep -q '"language.code_exception"' internal/baseline/assets/decisions.json && bin_dir="$(mktemp -d)" && home_dir="$(mktemp -d)" || exit 1; go build -buildvcs=false -o "$bin_dir/roundfix" ./cmd/roundfix || { rm -rf "$bin_dir" "$home_dir"; exit 1; }; out="$(HOME="$home_dir" "$bin_dir/roundfix" baseline update --repo . --no-skills --no-history --format json </dev/null 2>&1)"; status=$?; rm -rf "$bin_dir" "$home_dir"; test "$status" -eq 0 || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q '"state":"current"' || { printf '%s\n' "$out"; exit 1; }; ! { tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- "Code Language Exception"; } && go test -count=1 -timeout 20m -tags docscontract -run '^TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent$' ./internal/cli && go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec` — expected: exit 0. Before this Task the normalization, the template token, the gated clause and the decision do not exist; afterwards they do, a Managed Refresh planned under a fresh home reports this repository current, its guide states no exception, and the update transcript and coverage record contracts pass.

## References

- `_prd.md` → Goal 5; User Story 5; Core Feature 7; Core Feature 10; Success Metric 3
- `_techspec.md` → System Architecture (Code Language Exception); Interfaces, invariants 2, 3 and 5; Data Models; API Contract 2; API Contract 7; API Contract 9; API Contract 10; API Contract 11; API Contract 12; Surface Transcript 2; Testing Approach; Build Order 2; Risks & Considerations; Decisions
- `_authorization.md` → Sanctioned regeneration
- [`references/2026-11-27-characterization.md`](references/2026-11-27-characterization.md), section 8
- ADR-0118, ADR-0222, ADR-0250, ADR-0263
