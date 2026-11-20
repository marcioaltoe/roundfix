---
task: task_01
spec: 0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd
status: completed
type: backend
complexity: high
---

# Task 01: The Secondbrain search and the duplicate check name the Indexus in the Baseline, in write-idea and in a decision record

## Overview

This Task delivers `_techspec.md` → Build Order 1 and 2 in one change. Two clauses of the Baseline's Secondbrain module are reworded under their ids to name the Indexus tool, the research step of the owned skill `write-idea` names the same tool, a new accepted ADR records the decision and ADR-0330 gains its amendment notice. Every derived file comes from the repository's own commands, and this repository's Secondbrain guide and Setup Manifest come from the Managed Refresh. It is verifiable on its own through seven new tests over the embedded catalog, a plan on a temporary repository, the embedded owned skills and this repository's decision records, and through the existing Baseline, skills and documentation contracts.

## Requirements

1. MUST replace the `guidance` of `clause.secondbrain.02-query-order` in `internal/baseline/assets/modules/secondbrain.json` with the text of `_techspec.md` → API Contract 1, and the sentence of `clause.secondbrain.research-to-the-knowledge-owner` that begins "Run the advisory" with the text of API Contract 2, every earlier sentence of that clause unchanged. Each clause keeps its id, its `mandatory` force and its `replaces` value; no rule version, guide version, other clause, other module, template, profile source, Source Baseline or Setup Snapshot changes by hand (`_authorization.md` → Limits).
2. MUST replace, in `.agents/skills/write-idea/SKILL.md`, the part of the research-step paragraph that begins "When the session performed substantive external research" after its first sentence with the two sentences of `_techspec.md` → API Contract 3, wrapped as the file wraps its lines. No other line of the skill changes by hand, and no `### QA settlement` section of any skill changes.
3. MUST then run, in this order, the commands `_authorization.md` → Sanctioned regeneration declares, and hand-edit none of their outputs:
   1. `make skills-sync`;
   2. `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, the `-record-skill-versions` step, which raises both version fields of both `SKILL.md` copies of `write-idea` to the next version and records it in `skills/testdata/owned-skill-versions.json`;
   3. `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`, which writes the next version into the module's version line and records it in `internal/baseline/module-versions.json`;
   4. `make baseline-digests`, which rewrites the formatter golden of the Secondbrain guide, the Standard TypeScript profile's pin, the three catalog files and the four plan goldens.
   A version is "the next version", never a number written by hand. When another Spec changed a Baseline module or an owned skill first, these commands regenerate on the tree this Task runs on.
4. MUST create `docs/adr/0449-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-with-no-second-path.md`, the decision record of `_techspec.md` → API Contract 5, with the title "The Secondbrain search and the duplicate check go through the Indexus, with no second path", the lifecycle header of `docs/agents/docs-layout.md` with `status: accepted` and, on the next line, `knowledge: true`, and the content in the order API Contract 5 gives. It names no client and no acervo and cites no Spec folder path. When 0449 is no longer the next free number on the tree this Task runs on, stop and report it; do not take another number.
5. MUST add to `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md` the notice of `_techspec.md` → API Contract 6: one paragraph directly after the first paragraph of `## Decision`, beginning with the bold lead "Amended by", the new decision record's identifier as a link to its file and the date of its `created_at` in parentheses, in the form of the three notices below it. The sentence that names the earlier command stays byte-identical. Set ADR-0330's `updated_at` to the new decision record's `created_at`; leave `amended_by` to the update.
6. MUST then run the Managed Refresh once, `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`, which writes `docs/agents/secondbrain.md`, `docs/agents/setup-context.json`, the new path in the Standard Export List of `.secondbrain-export` and the new decision record's identifier as the last entry of ADR-0330's `amended_by`. None of these four is edited by hand.
7. MUST add `internal/baseline/secondbrain_indexus_clauses_test.go` with the four tests of `_techspec.md` → Interfaces and Testing Approach, over invariants 1 to 4:
   - `TestTheSecondbrainClausesNameTheIndexus`, built on `embeddedBaselineClauses`, with these subtests, each named as written: `keeps_identity_and_force` (both clauses `mandatory`, ids unchanged, `research-to-the-knowledge-owner` still replaces `clause.secondbrain.research-capture`); `search_clause`, holding one subtest per search-clause phrase of API Contract 4; `duplicate_check_clause`, holding one subtest per duplicate-check phrase; and `no_retired_word`.
   - `TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten`, with the subtests `planned_guide` (the postimage of `docs/agents/secondbrain.md` in a plan of the Standard TypeScript profile with `secondbrain.enabled` true), `formatter_golden` (read with `catalog.Asset`) and `repository_guide` (this repository's `docs/agents/secondbrain.md`), each holding one subtest per phrase of both clauses.
   - `TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest`, which clones the embedded catalog, gives both clauses the sentences of the Delivery Base with `mutateCatalogClause` of the package's plan tests, adopts the Standard TypeScript profile with the Secondbrain turned on in a temporary repository, commits, plans a Managed Refresh with the embedded catalog and asserts `_techspec.md` → API Contract 7: exactly the two files `docs/agents/secondbrain.md` and `docs/agents/setup-context.json` change; `clause.secondbrain.02-query-order` is `retained` with `mandatory` force; `clause.secondbrain.research-capture` is `replaced` with the target `clause.secondbrain.research-to-the-knowledge-owner` and `mandatory` force; no clause is `unaccounted`; and after the plan is applied a second plan changes no file. The earlier sentences are written in the test from `references/2026-11-20-characterization.md`, section 2, and the test reads no file under the Spec root.
   - `TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand`, which walks `assets/` and `../../docs/agents` as `TestAuthorizationClausesDropTheRemovedPhrases` walks them, for the whole word `qmd` in any letter case. It holds one exempt path, the Source Baseline corpus copy of the Secondbrain guide, with its reason; it fails when that path no longer holds the word, and fails when it scanned no module or no formatter golden.
   The phrase subtests are named, in API Contract 4's order: for the search clause `search_with_the_tool`, `one_group_per_reached_acervo`, `abrir_conteudo_opens_a_page`, `wiki_path_relative_to_wiki`, `no_wiki_result_is_a_limitation`, `no_other_search_tool`, `index_first`, `mirrors_are_references`; for the duplicate-check clause `duplicate_check_with_the_tool`, `each_returned_page_opens`, `extend_existing_knowledge`, `no_label_or_rank_decides`. Each asserts its phrase exactly once in its text after whitespace is collapsed (invariant 3).
8. MUST add `skills/indexus_duplicate_check_test.go` with `TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus`, which asserts the four `write-idea` phrases of API Contract 4 in the canonical copy and the mirror with `assertIdeaStageSkillPhrases`, in subtests named `duplicate_check_with_the_tool`, `each_returned_page_opens`, `no_label_or_rank_decides` and `absent_tool_substitutes_nothing`; and `TestNoOwnedSkillNamesTheRetiredSearchCommand`, which sweeps every file `skills.Files()` returns for the whole word `qmd` in any letter case and fails when it swept no file.
9. MUST add `internal/docscontract/indexus_decision_record_test.go`, with the `docscontract` build tag and the `verify:always` marker the package's other files carry, holding `TestTheIndexusDecisionAmendsTheConsultationRule`. It finds the decision record by its title, never by its number (invariant 5), and asserts with `spec.ReadADRHeader`: it is accepted and its Knowledge Mark is true; ADR-0330's `amended_by` holds its identifier; the notice of API Contract 6 follows the sentence of ADR-0330 that names the earlier command and carries the date of the record's `created_at`; and `.secondbrain-export` holds the record's path once, between the Standard Export List markers, with its ADR lines in ascending order.
10. MUST keep these existing tests unchanged and passing: `TestBaselineClauseForceIsCharacterized`, `TestEveryInboxClauseAnAdopterHoldsIsReplaced`, `TestFormatterComposition`, `TestTheSetupManifestIsCurrent`, `TestEveryBaselineModuleVersionIsRecorded`, `TestAuthorialSkillSync`, `TestEveryOwnedSkillVersionIsRecorded` and `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent`. No new test reaches the network, calls an MCP server, runs `qmd`, reads the real home directory, or reads a file under the Spec root or outside the repository (invariant 4).
11. MUST NOT change, by hand or by command, the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `CONTEXT.md`, any Vendored Skill (`knowledge-workspace` included), any other owned skill or file of `write-idea`, any guide under `docs/agents/` other than `secondbrain.md` and `setup-context.json`, the Source Baseline corpus, manifests and index, `docs/references/skill-coverage.json`, `docs/references/behavior-surfaces.json`, `docs/references/coverage-record.json`, any user-guide page, `CHANGELOG.md` or any version file, and any accepted ADR other than ADR-0330.

## Subtasks

- [ ] Reword the two clauses and the `write-idea` step.
- [ ] Run the skill sync, the skill record command, the module record command and `make baseline-digests`, in that order.
- [ ] Write the new decision record and the notice in ADR-0330, then run the Managed Refresh.
- [ ] Add the four Baseline tests, the two skills tests and the decision-record test.

## Acceptance Criteria

- [ ] Both clauses keep id, `mandatory` force and `replaces`, and each phrase of API Contract 4 occurs exactly once in the module, in the planned guide, in the formatter golden and in this repository's guide; no Baseline asset outside the exempt Source Baseline copy and no guide under `docs/agents/` names `qmd` in any letter case.
- [ ] An adopter that holds the earlier wording with the Secondbrain turned on receives a plan that changes only its Secondbrain guide and its Setup Manifest, keeps the clause accounting of API Contract 7 with nothing unaccounted, and plans no change after the apply.
- [ ] Both copies of `write-idea` hold the four phrases, are byte-identical, carry the next recorded version, and no owned skill names `qmd`.
- [ ] The new decision record is accepted with a true Knowledge Mark, ADR-0330 lists it in `amended_by` and carries its dated notice after the sentence that names the earlier command, and the Standard Export List holds its path once, in order.
- [ ] Every derived file is the output of the repository's commands, and the Baseline, skills and documentation contracts pass.

## Context

- instruction: `internal/baseline/clause_characterization_test.go`
- instruction: `internal/baseline/authorization_clauses_test.go`
- instruction: `internal/baseline/cross_project_continuity_test.go`
- instruction: `internal/baseline/plan.go`
- instruction: `skills/idea_stage_judgment_test.go`
- instruction: `skills/skills.go`
- instruction: `internal/spec/adr_header.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- interface: `internal/baseline/assets/modules/secondbrain.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/secondbrain.md`
- interface: `docs/agents/setup-context.json`
- interface: `.secondbrain-export`
- interface: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
- interface: `.agents/skills/write-idea/SKILL.md`
- interface: `skills/write-idea/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- creates: `docs/adr/0449-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-with-no-second-path.md`
- creates: `internal/baseline/secondbrain_indexus_clauses_test.go`
- creates: `skills/indexus_duplicate_check_test.go`
- creates: `internal/docscontract/indexus_decision_record_test.go`

## Consumers

- the embedded Baseline catalog, its digest and this repository's Setup Manifest — consumers: `./internal/baseline`, `./internal/cli`
- the owned skill `write-idea` and the owned-skill version record — consumers: `./skills`
- this repository's decision records and the Standard Export List — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheSecondbrainClausesNameTheIndexus|TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten|TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest|TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSecondbrainClausesNameTheIndexus/keeps_identity_and_force TestTheSecondbrainClausesNameTheIndexus/search_clause/search_with_the_tool TestTheSecondbrainClausesNameTheIndexus/search_clause/no_other_search_tool TestTheSecondbrainClausesNameTheIndexus/duplicate_check_clause/duplicate_check_with_the_tool TestTheSecondbrainClausesNameTheIndexus/no_retired_word TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten/planned_guide TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten/formatter_golden TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten/repository_guide TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the four Baseline tests and the named subtests pass.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus|TestNoOwnedSkillNamesTheRetiredSearchCommand)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus/duplicate_check_with_the_tool TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus/absent_tool_substitutes_nothing TestNoOwnedSkillNamesTheRetiredSearchCommand; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, both skills tests pass.
- `out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestTheIndexusDecisionAmendsTheConsultationRule|TestTheSetupManifestIsCurrent)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheIndexusDecisionAmendsTheConsultationRule TestTheSetupManifestIsCurrent; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the decision record amends ADR-0330 and this repository's Setup Manifest is current.
- `f=docs/adr/0449-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-with-no-second-path.md; test -f "$f" && grep -qx 'status:.accepted' "$f" && grep -qx 'knowledge:.true' "$f" && grep -q '^amended_by:.*ADR.0449' docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md && grep -qxF -- "$f" .secondbrain-export && grep -q 'pesquisar_acervos' internal/baseline/assets/modules/secondbrain.json && grep -q 'pesquisar_acervos' docs/agents/secondbrain.md && grep -q 'pesquisar_acervos' .agents/skills/write-idea/SKILL.md && ! grep -qiw 'qmd' internal/baseline/assets/modules/secondbrain.json docs/agents/secondbrain.md internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md .agents/skills/write-idea/SKILL.md skills/write-idea/SKILL.md && cmp -s .agents/skills/write-idea/SKILL.md skills/write-idea/SKILL.md` — expected: exit 0, the new decision record is accepted and marked, ADR-0330 and the export list name it, the module, the guide and both skill copies name the Indexus tool and not the retired command, and the two skill copies are byte-identical.
- `grep -q 'pesquisar_acervos' internal/baseline/assets/modules/secondbrain.json && grep -q 'pesquisar_acervos' skills/write-idea/SKILL.md && go test -count=1 -timeout 20m ./internal/baseline ./skills` — expected: exit 0, the whole Baseline and skills packages pass on the changed module and skill and their regenerated versions, digests and goldens.
- `grep -q 'pesquisar_acervos' docs/agents/secondbrain.md && go test -count=1 -timeout 20m -run 'Baseline|SetupManifest' ./internal/cli` — expected: exit 0, the Baseline commands' tests, `TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent` among them, pass after the Managed Refresh.

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-8; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6
- `_techspec.md` → System Architecture; Interfaces; Data Models; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Coverage Map; Testing Approach; Build Order 1; Build Order 2
- `_authorization.md` → Sanctioned regeneration; Limits
- [`references/2026-11-20-characterization.md`](references/2026-11-20-characterization.md), sections 2, 3, 5 and 7
- ADR-0330, ADR-0257, ADR-0186, ADR-0250, ADR-0438, ADR-0439


## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, testing-boss, tech-writer, setup-context-driven, roundfix, evidence-gate, systematic-debugging, no-workarounds

Implemented the Task 01 slice. Task status remains Daemon-owned; Subtask/Acceptance checkboxes and the Task Graph were left unchanged; no declared Verification command, repository Verification, commit, push, Pull Request or cross-project message was run.

### Implementation and acceptance evidence

1. **Clause wording and distributed guides.** Reworded only the two authorized guidance fields; both ids, mandatory forces and replacement lists are preserved, as are every rule and guide version and the earlier research-digest sentences. `TestTheSecondbrainClausesNameTheIndexus` pins all twelve phrases once after whitespace normalization and checks identity, force, replacements and the retired word. `TestTheSecondbrainGuideNamesTheIndexusWhereverItIsWritten` checks every phrase in a Standard TypeScript plan, the embedded formatter golden and this repository's guide. `TestNoBaselineAssetOrGuideNamesTheRetiredSearchCommand` sweeps all assets and repository guides, counts modules and formatter goldens, and requires its single historical Source Baseline exemption still to contain the retired word. All four new Baseline tests passed in the focused command below.
2. **Earlier adopter and convergence.** `TestAnAdopterOfTheEarlierWordingReceivesOnlyTheGuideAndTheManifest` uses literal Delivery Base sentences in a cloned catalog, adopts the Standard TypeScript profile with the Secondbrain enabled into a temporary repository, and commits that fixture. Its refresh changes exactly the Secondbrain guide and Setup Manifest, accounts the search clause as mandatory/retained and research-capture as mandatory/replaced by research-to-the-knowledge-owner, rejects any unaccounted clause, applies, and requires a second plan with no file changes. The focused Baseline run passed this test.
3. **Owned skill.** The research paragraph uses the four required Indexus phrases and the absent/silent-tool rule. The sanctioned recorder chose 0.0.18 for both version fields in both copies; the copies are byte-identical. `TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus`, all four named subtests, `TestNoOwnedSkillNamesTheRetiredSearchCommand`, `TestAuthorialSkillSync` and `TestEveryOwnedSkillVersionIsRecorded` passed in the focused skills run. The sweep reads every embedded owned-skill file and rejects an empty scan.
4. **Decision record and amendment.** Confirmed the highest existing ADR was 0448 before exclusively creating ADR-0449, accepted with `knowledge: true`. The record follows API Contract 5, including the maintainer's words, limitations, release condition, alternatives and amendment of ADR-0330; it names no client or acervo and cites no Spec path. Added only ADR-0330's dated notice and matching `updated_at` by hand, leaving its historical paragraph unchanged. The Managed Refresh added the new identifier as its last `amended_by` entry and inserted the new path once in the Standard Export List. `TestTheIndexusDecisionAmendsTheConsultationRule` finds the new record by title and passed its header, notice, timestamp and ordered-export checks.
5. **Derived artifacts and unchanged contracts.** All derived writes came from the sanctioned sequence below. The module recorder chose version 16; no version or digest was transcribed. The Managed Refresh ran once, reported four changed paths (guide, Setup Manifest, export list and ADR-0330) and applied its approved postimages. The Setup Manifest's recorded Roundfix version also moved from 0.150.0 to the running 0.152.0 through that command. It reported the two existing nested fixture/corpus-carrier warnings and left those carriers unchanged. Existing clause-force, inbox-replacement, formatter, module-version, skill-sync, skill-version, Setup Manifest and repository-update-transcript contracts passed unchanged in focused runs. Full Baseline/skills/CLI package Verification and the declared sequences remain for the Daemon; the focused results do not claim those broader gates.

### Sanctioned regeneration, in execution order

All commands exited 0; the shell used `rtk proxy` to preserve their output.

1. `make skills-sync`
2. `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
3. `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
4. `make baseline-digests` — reported `ok: true`, `changed: true`; regenerated the formatter golden, Standard TypeScript pin, three catalog records and four plan goldens.
5. After writing ADR-0449 and the notice: `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text` — applied exactly the four authorized outputs, with skills and history skipped.

### Focused implementation checks

- Red starting evidence, before source edits: `go test ./internal/baseline -run '^TestTheSecondbrainClausesNameTheIndexus$/search_clause$' -count=1` failed on the six missing Indexus-specific search phrases; `go test ./skills -run '^TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus$/duplicate_check_with_the_tool$' -count=1` failed on the missing phrase in both copies; `go test -tags docscontract ./internal/docscontract -run '^TestTheIndexusDecisionAmendsTheConsultationRule$' -count=1` failed because no ADR had the required title.
- `go test ./internal/baseline -run '^(TestTheSecondbrain|TestAnAdopterOfTheEarlierWording|TestNoBaselineAssetOrGuide)' -count=1` — passed all four new Baseline tests.
- `go test ./skills -run '^(TestWriteIdeaRunsTheDuplicateCheckThroughTheIndexus|TestNoOwnedSkillNamesTheRetiredSearchCommand|TestAuthorialSkillSync|TestEveryOwnedSkillVersionIsRecorded)$' -count=1 -v` — passed both new tests, all four phrase subtests and both existing contracts.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheIndexusDecisionAmendsTheConsultationRule$' -count=1` — passed after the Managed Refresh.
- `go test ./internal/baseline -run '^(TestBaselineClauseForceIsCharacterized|TestEveryInboxClauseAnAdopterHoldsIsReplaced|TestFormatterComposition|TestEveryBaselineModuleVersionIsRecorded)$' -count=1` — passed unchanged.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSetupManifestIsCurrent$' -count=1` — passed.
- `go test ./internal/cli -run '^TestBaselineUpdateSurfaceTranscriptThisRepositoryIsCurrent$' -count=1` — passed.
- `gofmt` applied only to the three new test files; `git diff --check` exited 0. A structural source audit confirmed the two authorized clause fields plus recorded module version are the only module changes, the preceding research sentences are preserved, the skill copies are byte-identical, and the new tests contain no Spec-root read or direct process call.

### Scope and execution notes

The only initially changed path was this Task file, whose `in_progress` status was already set by the Daemon. It was preserved. The executor refused to start under the symlinked writable root, including when given the physical working directory; approved escalated shell execution in the physical checkout was used for this Task. No repository or user configuration was changed to bypass that error. No live Indexus, network, Secondbrain, home-directory fixture, Agent session or MCP call was added to a test. No follow-up outside the Task slice was implemented.

Initial handback scope postflight: all 23 changed/new paths are declared interfaces/creates or this assigned Task file; every existing contract file named in the context is unchanged. Comparing the Task authored projection to HEAD confirmed no change outside `## Result` and the pre-existing Daemon status. `git diff --check` exited 0 after the Result edit.


### Verification Feedback repair — attempt 1

Inspected both Daemon diagnostic artifacts under the Run's `verification/` directory: `batch-001-attempt-1.log` and `batch-001-attempt-1-settlement-authorization.log`. The repository gate reached `internal/verifyselect` and exposed an exact file-count expectation invalidated by this Task's required new docscontract file; this was not a missing relevance declaration.

- Updated only the expected docscontract file count from 10 to 11 in `internal/verifyselect/contracts_repository_test.go`, leaving all declaration, relevance, selection and gate assertions intact. This additional ordinary test dependency is directly required by the new `indexus_decision_record_test.go`; it is not a governed path in `internal/speccheck/governed.go`. The Daemon may record that additional path under its ordinary-path contract. No test-runner configuration, Makefile, authorization record or other Task was edited.
- Focused check: `go test ./internal/verifyselect -run '^TestRepositoryContractTestsDeclareTheirRelevance$' -count=1` exited 0. The one-line diff was reviewed. No full repository gate or declared Task Verification sequence was rerun.
- Remaining integration blocker: the authorization audit uses delivery-target revision `2d8a2abda98ea907bb5a126c875711f2ee3189de`, where this Spec's `_authorization.md` is absent. A read-only `git ls-tree` of that exact revision and path confirmed the absence. Changing this worktree's implementation or rewriting the grant cannot place the grant in the delivery target's ancestry. The Supervisor must land the already approved bounded authorization there before the consuming commit; this child session created no commit, changed no ref and widened no grant.

The implementation repair is handed back with the focused check above; the authorization/ancestry blocker remains explicit. Task status remains `in_progress` and Daemon-owned. The earlier 23-path handback audit is historical; this repair adds only the ordinary test dependency named above, for 24 changed/new paths in total.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/verifyselect/contracts_repository_test.go`
