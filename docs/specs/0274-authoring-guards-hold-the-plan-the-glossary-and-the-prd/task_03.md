---
task: task_03
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
status: pending
type: test
complexity: high
---

# Task 03: A contract test holds the glossary to the rule it applies to Specs

## Overview

Every Spec is held to the glossary, but nothing holds `CONTEXT.md` to itself
(`_techspec.md` → Current behavior). This Task adds a Repository Contract Test
with the five rules of API Contract 4, its allowlist and its constant mapping
in one place, and fixes the gaps it finds in `CONTEXT.md` so it starts green.
It is verifiable on its own by the six named tests, one of which plants a
violation of each rule.

## Requirements

1. MUST answer the adopted Backlog Entry of 2026-10-16,
   [the glossary is not held to the rule it applies to Specs](references/2026-10-16-the-glossary-is-not-held-to-the-rule-it-applies-to-specs.md),
   by adding its tests to the existing `internal/docscontract/user_guide_contract_test.go`
   (`docscontract` build tag, `//verify:always`) and to no new file, because
   `TestRepositoryContractTestsDeclareTheirRelevance` in
   `internal/verifyselect/contracts_repository_test.go` requires every
   `docscontract` file to be an always contract and pins their number at nine
   (ADR-0280).
2. MUST implement the five rules of API Contract 4 exactly as written, one
   top-level test each on the live `CONTEXT.md`:
   `TestGlossaryPhrasesAreDefinedTerms`, `TestGlossaryLinksNameAcceptedADRs`,
   `TestGlossaryClosedListsEqualTheirConstants`,
   `TestGlossaryAvoidLinesArePlainLists` and `TestGlossaryStatesEachRuleOnce`,
   plus `TestGlossaryContractReportsAPlantedViolationOfEachRule`, which runs
   the same rule functions on a synthetic glossary holding an undefined
   compound such as `Spec Runner Lane`, a link to a `proposed` ADR and one to
   a file under `docs/specs/`, a Task Type list missing `qa`, an `_Avoid_`
   line with a sentence, and one sentence of more than sixty characters in two
   entries, and requires each rule to fail with the entry it names. Every test
   calls `t.Parallel()` first and reads no file outside the repository.
3. MUST keep the allowlist and the closed-list mapping as two package-level
   tables at the top of the test file; each allowlist entry carries a
   non-blank reason, and a test fails on an allowlist entry rule 1 no longer
   needs. The mapping MUST read the exported constants, never copy their
   values (`docs/agents/agent-instructions.md`: "An assertion reads the
   constant it means.").
4. MUST resolve every finding the rules report on the tree this Task starts
   from, which includes Spec 0272's glossary changes when it merged first: for
   each of the uncovered runs listed in `_techspec.md` → Current behavior and
   any new one, reword the definition to a defined term without changing its
   meaning, or allowlist a proper name, a product name or an artifact section
   label with its reason; rewrite the Corrective Spec `_Avoid_` line as a
   plain list. MUST NOT add, rename or remove a glossary term, change an ADR
   link's target, or change the meaning of any definition.
5. MUST NOT change any file other than the new test and `CONTEXT.md`.

## Subtasks

- [ ] Write the entry parser and the five rule functions.
- [ ] Add the allowlist and the closed-list mapping tables.
- [ ] Add the six tests, including the planted violations.
- [ ] Fix the gaps the rules find in `CONTEXT.md`.

## Acceptance Criteria

- [ ] The five live-glossary tests pass on the delivered `CONTEXT.md`.
- [ ] Each planted violation makes its rule fail and names its entry.
- [ ] Every allowlist entry has a reason and is still needed.
- [ ] No glossary term was added, renamed or removed.

## Consumers

None. This Task adds a contract test and rewords glossary prose; no code consumes `CONTEXT.md`.

## Context

- instruction: `docs/adr/0280-a-guide-task-records-its-own-surface-fingerprint-a-prd-records-who-decided-and-the-glossary-is-held-to-its-own-rule.md`
- instruction: `docs/adr/0244-a-spec-declares-the-domain-terms-it-introduces-and-the-check-holds-them-to-the-glossary.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/speccheck/glossary.go`
- interface: `internal/docscontract/user_guide_contract_test.go`
- instruction: `internal/verifyselect/contracts_repository_test.go`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryClosedListsEqualTheirConstants|TestGlossaryAvoidLinesArePlainLists|TestGlossaryStatesEachRuleOnce|TestGlossaryContractReportsAPlantedViolationOfEachRule)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGlossaryPhrasesAreDefinedTerms TestGlossaryLinksNameAcceptedADRs TestGlossaryClosedListsEqualTheirConstants TestGlossaryAvoidLinesArePlainLists TestGlossaryStatesEachRuleOnce TestGlossaryContractReportsAPlantedViolationOfEachRule; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -q '^//verify:relevant .*CONTEXT.md' internal/docscontract/glossary_contract_test.go || { printf 'the glossary contract declares no relevance to CONTEXT.md\n' >&2; exit 1; }` — expected: exit 0. Before this Task the six tests do not exist, so the command fails at its first pass check.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 3; Success Metric 3
- `_techspec.md` → API Contract 4; Invariant 5; Current behavior; Testing Approach; Build Order 3
- ADR-0280; ADR-0244; ADR-0091
