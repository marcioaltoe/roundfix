---
task: task_01
spec: 0307-the-roundfix-skill-modes-and-up-to-three-questions-per-call
status: pending
type: backend
complexity: medium
---

# Task 01: The Baseline clause allows up to three independent questions per call

## Overview

The Baseline's stop-and-ask clause `clause.core.ask-user-answerable-decisions`
requires exactly one question per call. On 2026-10-22 the maintainer answered
"Até 3 no geral" for every adopter. This Task rewords the clause to at most
three independent questions per call, keeps its identity and enforcement so an
adopter's refresh records it `retained`, regenerates every derived file through
the sanctioned commands, and holds the result with three tests (`_prd.md` →
Core Feature 1; ADR-0320).

## Requirements

1. MUST answer the Backlog Entry Spec 0306 adopted,
   [`../0306-the-follow-up-command-lists-pending-work-in-priority-order/references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md`](../0306-the-follow-up-command-lists-pending-work-in-priority-order/references/2026-10-22-follow-up-status-and-unattended-modes-for-the-roundfix-skill.md)
   of 2026-10-22, decision 3 ("Questions").
2. MUST set the `guidance` of `clause.core.ask-user-answerable-decisions` in
   `internal/baseline/assets/modules/core.json` to exactly the text of
   `_techspec.md` → API Contract 1, raise the `version` of its rule
   `rule.core.decision-authority` by one over the value the tree holds, and
   change nothing else in the module by hand: the clause keeps its `id`, its
   `enforcement` `stop-and-ask` and gains no `replaces`.
3. MUST then run, in this order, and hand-edit none of their outputs:
   - `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
     which writes the next module version into `core.json` and records it in
     `internal/baseline/module-versions.json`;
   - `make baseline-digests`, which rewrites the derived files this Task
     declares (measured on a prototype: the Standard TypeScript golden
     `docs/agents/agent-instructions.md`, the Standard TypeScript profile, the
     three catalog files and the four plan characterization goldens);
   - `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
     the Managed Refresh of this repository, which rewrites
     `docs/agents/agent-instructions.md` and `docs/agents/setup-context.json`,
     then the same command again, which must report no file change. Run both
     with `HOME` set to a temporary directory, so the refresh reads and writes
     no real Roundfix Home.
   A derived file the commands leave unchanged is not edited; a file they
   change that this Task does not declare is a stop-and-report condition.
4. MUST create `internal/baseline/question_batch_clause_test.go` with the
   constants of `_techspec.md` → Interfaces and the tests
   `TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall`,
   `TestTheAskUserClauseRendersInTheAgentInstructionsGuide` and
   `TestAnAdopterRetainsTheAskUserClause`, as `_techspec.md` → Testing
   Approach describes them, using the existing helpers
   `embeddedBaselineClauses` and `newClauseReplacementAdopter`, comparing on
   whitespace-normalized text through a helper whose name no other file of the
   package uses. Each top-level test and subtest calls `t.Parallel()` first.
5. MUST keep every existing test of `internal/baseline`, `internal/cli` and
   `internal/docscontract` passing without editing it; the clause
   characterization pins only the enforcement, which does not change.
6. MUST search the repository, outside `docs/history/`, for "exactly one
   question", "one question per call", "single question" and "Never bundle
   questions", and report every match left: the owned interview skills keep
   their one-question interviews by decision (ADR-0320), and any other match
   is a stop-and-report condition.
7. MUST NOT change the Source Baseline corpus, manifest, accounting or index,
   `internal/baseline/assets/retention/**`, any module other than `core`, any
   skill, any user guide page, `.roundfixrc.yml`, or any Go production file.

## Subtasks

- [ ] Reword the clause and raise its rule version.
- [ ] Record the module version, run `make baseline-digests` and the Managed Refresh twice.
- [ ] Add the three tests.
- [ ] Sweep for the retired wording and report.

## Acceptance Criteria

- [ ] The embedded clause carries the three-question limit and the independence sentence once, keeps `stop-and-ask` and no `replaces`, and no longer carries the retired sentence.
- [ ] The Standard TypeScript golden and this repository's agent instructions render the new text, and a second Managed Refresh changes nothing.
- [ ] A Standard TypeScript adopter's refresh records the clause `retained` with itself as target.

## Consumers

- the core module's clause text and the catalog digests — consumers: `./internal/baseline`

## Context

- instruction: `docs/adr/0320-a-question-call-carries-up-to-three-independent-questions.md`
- instruction: `internal/baseline/clause_characterization_test.go`
- instruction: `internal/baseline/clause_replacement_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/module-versions.json`
- creates: `internal/baseline/question_batch_clause_test.go`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/setup-context.json`

## Verification

- `for name in TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall TestTheAskUserClauseRendersInTheAgentInstructionsGuide TestAnAdopterRetainsTheAskUserClause; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run '^(TestTheAskUserClauseAllowsThreeIndependentQuestionsPerCall|TestTheAskUserClauseRendersInTheAgentInstructionsGuide|TestAnAdopterRetainsTheAskUserClause|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestCatalogDiagnosticCharacterization|TestBaselinePlanCharacterization|TestFormatterComposition|TestBaselineCompatibilityCorpus|TestTheCoreGuidesStateThePromotedRules|TestNoTwoBaselineClausesShareText|TestBaselineClauseForceIsCharacterized)$' ./internal/baseline` — expected: exit 0. Before this Task none of the three tests exists, so the command fails at its first check. Afterwards they exist and pass together with the module record, the catalog, the derived-file characterizations and the clause characterization of the package.
- `tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- 'Ask at most three questions per call or message, never more than the tool accepts, and keep only that call pending.' || { echo 'agent instructions lack the three-question limit' >&2; exit 1; }; if grep -qF 'Ask exactly one question per call' docs/agents/agent-instructions.md; then echo 'agent instructions keep the retired sentence' >&2; exit 1; fi; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the repository's agent instructions carry the one-question sentence, so the command fails at its first check. Afterwards the Managed Refresh has rendered the new clause and the documentation contract suite, which loads the embedded catalog, passes.

## References

- `_prd.md` → Goals; Core Feature 1; Success Metric 1; Success Metric 2; Acceptance evidence
- `_techspec.md` → Current behavior 1-4; API Contracts 1-2; Invariants 1-3; Testing Approach; Build Order 1
- ADR-0320
