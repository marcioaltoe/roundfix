---
task: task_03
spec: 0417-the-baseline-carries-the-language-rule
status: pending
type: backend
complexity: high
---

# Task 03: The review and the QA gate report what breaks the rule, and the skills, guides and glossary describe it

## Overview

This Task delivers `_techspec.md` → Build Order 3: the Pre-PR Review prompt gains one instruction that asks for a finding on a new name that breaks the repository's language rule, the `qa-gate` skill gains a section that asks for a finding on text a user reads outside the Reader Language, and the Roundfix skill, the `setup-context-driven` skill, three user-guide pages and the glossary describe what Tasks 01 and 02 built. The skills are mirrored and versioned by their commands. It is verifiable on its own: the prompt carries the instruction in its place with its grammar unchanged, the documents quote it as the code writes it, the embedded QA gate skill carries its section, and the glossary defines the five terms.

## Requirements

1. MUST add to `internal/cli/review.go` the constant `reviewLanguageRuleInstruction`, holding exactly the text of `_techspec.md` → API Contract 13, and make `buildReviewPrompt` write it and a line break directly after the absence-claim instruction. The finding grammar, the Delivery Conventions, `deliveryConventionsVersion` and every other line of the prompt stay as they are (invariant 7). It is no Delivery Convention.
2. MUST add, in a new file `internal/cli/review_language_rule_test.go`, `TestTheReviewPromptCarriesTheLanguageRuleInstruction` (Success Metric 5): the prompt holds the constant once, directly after the absence-claim instruction, and still holds the finding grammar lines and the Delivery Conventions block with its version.
3. MUST quote the instruction, as the constant holds it, in `.agents/skills/roundfix/references/review.md` and `docs/user-guide/commands/review.md`, directly after the quotation of the absence-claim instruction, and add to `internal/docscontract/command_documentation_test.go` `TestTheLanguageRuleInstructionIsDocumentedAsTheReviewerPromptsIt`, which reads the constant from the Go source of `internal/cli/review.go` and holds the two files and the mirror `skills/roundfix/references/review.md` to it with whitespace normalized, and `TestTheLanguageRuleInstructionContractReportsAStaleQuotation`, which proves the same check fails on a quotation that differs by one word. No file is created under `internal/docscontract/`.
4. MUST add to `.agents/skills/qa-gate/SKILL.md`, directly after the section `### Frontend sweep`, the section `### Reader Language` with exactly the paragraph of `_techspec.md` → API Contract 14. The `### QA settlement` section and every other section stay byte-identical. Add, in a new file `skills/qa_gate_reader_language_test.go`, `TestTheQAGateReportsTextOutsideTheReaderLanguage` (Success Metric 6), which reads the embedded `qa-gate` skill and asserts the section heading after `### Frontend sweep` and its paragraph.
5. MUST describe the Baseline's two decisions as `_techspec.md` → API Contract 15 says, in `.agents/skills/roundfix/references/baseline.md` and `docs/user-guide/commands/baseline.md`, beside the `--decision` paragraph: what `language.reader` is; that an update without it stops with exit `3`; that `--adopt-suggested` records `English` and so is right only for a product whose users read English; the whole command `roundfix baseline update --decision 'language.code_exception={"language":"<language>","source":"<source>"}'`; its refusal with exit `2` when the source is missing; and that a repository-owned profile lists `language.reader`, and lists `language.code_exception` before it can declare one. Add two sentences to `## Project decisions and Spec constraints` of `.agents/skills/setup-context-driven/SKILL.md` that name both decisions. Add to `docs/user-guide/context-driven-development.md` a short subsection on the Language Rule, the `_In code_` line of the Name in Code and the Code Language Exception.
6. MUST write the five Glossary Definitions the PRD's `## Glossary` declares — **Language Rule**, **Reader Language**, **Untranslated Term**, **Name in Code** and **Code Language Exception** — into `CONTEXT.md` through the `domain-modeling` skill, as `docs/agents/domain.md` says, from the definitions of `_prd.md` → Core Features 1, 2, 3, 4 and 7, each citing ADR 0454 as a link the way the neighbouring entries cite theirs, with no `_In code_` line (`_techspec.md` → API Contract 16). No other glossary entry changes.
7. MUST regenerate, in this order and by these commands only, as `_authorization.md` → Sanctioned regeneration declares them: `make skills-sync`; `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, which writes the next version of `qa-gate`, `roundfix` and `setup-context-driven` into both copies of each `SKILL.md` and `skills/testdata/owned-skill-versions.json`, never a number written by hand; `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, because this Task edits pages under `docs/user-guide/`; and last `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`. The mirrors under `skills/` are the output of `make skills-sync` and are not edited by hand.
8. MUST add no top-level test function other than the four this Task names, and rename, merge or remove no existing test. No new test file carries a build tag. No new test reads the network, the real home directory, another repository or the clock.
9. MUST NOT change the Baseline catalog, any file under `internal/baseline/`, the guides under `docs/agents/`, `docs/agents/setup-context.json`, ADR 0454, `.secondbrain-export`, `docs/references/skill-coverage.json`, any other owned or Vendored Skill, the Makefile, the workflows, `go.mod`, `go.sum`, lint configuration, `.roundfixrc.yml`, `docs/references/gate-time-targets.json` or any file that measures the Gate Time Rule.

## Subtasks

- [ ] Add the review instruction to the prompt, its test and its documentation contract.
- [ ] Add the `### Reader Language` section to the QA gate skill and its test.
- [ ] Describe both decisions in the Roundfix skill's Baseline reference, the `setup-context-driven` skill and the three user-guide pages.
- [ ] Write the five Glossary Definitions through `domain-modeling`.
- [ ] Run `make skills-sync`, the owned-skill record command, the surface record and the coverage record.

## Acceptance Criteria

- [ ] Every review prompt carries the language-rule instruction once, directly after the absence-claim instruction, with its finding grammar and Delivery Conventions unchanged, and the review reference, its mirror and the user guide quote it as the code writes it.
- [ ] Both copies of the `qa-gate` skill carry the `### Reader Language` section after `### Frontend sweep`, are byte-identical, and carry a version the owned-skill record holds for those bytes; the settlement table is unchanged.
- [ ] The Roundfix skill, the `setup-context-driven` skill and the user guide name `language.reader` and `language.code_exception`, the stop with exit `3`, the refusal with exit `2` and the whole `--decision` command.
- [ ] The glossary defines the five terms, each citing ADR 0454.

## Context

- instruction: `docs/agents/domain.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `internal/cli/review_conflicted_carry_test.go`
- interface: `internal/cli/review.go`
- interface: `internal/docscontract/command_documentation_test.go`
- interface: `.agents/skills/roundfix/references/review.md`
- interface: `skills/roundfix/references/review.md`
- interface: `docs/user-guide/commands/review.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills/setup-context-driven/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/references/coverage-record.json`
- creates: `internal/cli/review_language_rule_test.go`
- creates: `skills/qa_gate_reader_language_test.go`

## Consumers

- The review prompt that `buildReviewPrompt` writes — consumers: `./internal/cli`
- The owned skills `qa-gate`, `roundfix` and `setup-context-driven` and their version record — consumers: `./skills`
- The review quotation, the user-guide pages, the surface record and the glossary — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheReviewPromptCarriesTheLanguageRuleInstruction|TestAConflictingMergeOfTheDefaultBranchKeepsTheReviewerLineage)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestTheReviewPromptCarriesTheLanguageRuleInstruction ' || { echo 'missing pass: TestTheReviewPromptCarriesTheLanguageRuleInstruction' >&2; exit 1; }; grep -q 'reviewLanguageRuleInstruction' internal/cli/review.go || exit 1; names="$(sed -n 's/^func \(Test[A-Za-z0-9_]*\)(t \*testing\.T) {$/\1/p' internal/cli/review_test.go internal/cli/review_validation_test.go internal/cli/review_data_segments_test.go internal/cli/review_claude_segment_test.go internal/cli/review_omit_paths_test.go internal/cli/review_override_convention_test.go internal/cli/review_convention_qa_task_test.go | paste -sd '|' -)"; test -n "$names" || exit 1; go test -count=1 -timeout 20m -run "^($names)$" ./internal/cli` — expected: exit 0. Before this Task the prompt test does not exist; afterwards the prompt carries the instruction in its place, and the tests of the files that build a review prompt pass with it.
- `out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestTheLanguageRuleInstructionIsDocumentedAsTheReviewerPromptsIt|TestTheLanguageRuleInstructionContractReportsAStaleQuotation|TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem|TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms|TestProjectConstraintDocumentation|TestBaselineDecisionExamples)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheLanguageRuleInstructionIsDocumentedAsTheReviewerPromptsIt TestTheLanguageRuleInstructionContractReportsAStaleQuotation TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem TestTheSkillCoverageMapIsCurrent TestGlossaryPhrasesAreDefinedTerms; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec` — expected: exit 0. Before this Task the two quotation tests do not exist; afterwards the three review documents quote the instruction as the code writes it, the surface record and the glossary contracts pass, and the coverage record holds the new tests.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheQAGateReportsTextOutsideTheReaderLanguage|TestAuthorialSkillSync|TestEveryOwnedSkillVersionIsRecorded|TestSettlementGuidanceIsOneTable)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheQAGateReportsTextOutsideTheReaderLanguage TestAuthorialSkillSync TestEveryOwnedSkillVersionIsRecorded TestSettlementGuidanceIsOneTable; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; cmp -s .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && cmp -s .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md && cmp -s .agents/skills/roundfix/references/review.md skills/roundfix/references/review.md && cmp -s .agents/skills/setup-context-driven/SKILL.md skills/setup-context-driven/SKILL.md` — expected: exit 0. Before this Task the QA gate skill test does not exist; afterwards the embedded skill carries the section, every mirror matches its source, each changed skill's version is recorded for its bytes, and the settlement table is one.
- `grep -qF 'language.reader' .agents/skills/roundfix/references/baseline.md && grep -qF 'language.code_exception' .agents/skills/roundfix/references/baseline.md && grep -qF 'language.reader' docs/user-guide/commands/baseline.md && grep -qF 'language.code_exception' docs/user-guide/commands/baseline.md && grep -qF 'language.reader' .agents/skills/setup-context-driven/SKILL.md && grep -qF 'language.code_exception' .agents/skills/setup-context-driven/SKILL.md && grep -qF '### Reader Language' .agents/skills/qa-gate/SKILL.md && grep -F -- '**Language Rule**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Language Rule\*\*(:| —( |$))' && grep -F -- '**Reader Language**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Reader Language\*\*(:| —( |$))' && grep -F -- '**Untranslated Term**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Untranslated Term\*\*(:| —( |$))' && grep -F -- '**Name in Code**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Name in Code\*\*(:| —( |$))' && grep -F -- '**Code Language Exception**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Code Language Exception\*\*(:| —( |$))'` — expected: exit 0. Before this Task neither decision is named in the skills or the user guide and none of the five terms is defined; afterwards both decisions are described and each term has its Glossary Definition line.

## References

- `_prd.md` → Goal 4; User Stories 6 and 7; Core Features 8, 9 and 12; Success Metric 5; Success Metric 6; Success Metric 7; Glossary
- `_techspec.md` → System Architecture (Review and QA, Record of the decision); Interfaces, invariant 7; API Contract 7; API Contract 13; API Contract 14; API Contract 15; API Contract 16; Testing Approach; Build Order 3
- `_authorization.md` → Sanctioned regeneration
- ADR-0189, ADR-0196, ADR-0233, ADR-0244, ADR-0256, ADR-0280, ADR-0360
