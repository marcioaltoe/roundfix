---
task: task_02
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
status: pending
type: backend
complexity: high
---

# Task 02: The guidance sends cross-project work through ai-memory and retires the inbox clauses

## Overview

The Secondbrain guidance still grants `inbox/**` writes and describes the
Inbox Entry contract. This Task adds the core rule
`rule.core.cross-project-continuity`, whose four clauses replace the eight
inbox and capture clauses with the same force, moves the research digest to
a message in a new Secondbrain rule, makes the Secondbrain fully read-only
and rewords four workflow clauses. Every adopter's refresh then records each
retired clause `replaced` (`_prd.md` → Core Feature 2; ADR-0330).

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md`](references/2026-10-24-the-baseline-sends-handoffs-and-messages-through-ai-memory.md)
   of 2026-10-24, its expected outcomes 2 and 3.
2. MUST apply `_techspec.md` → Data Models 2 to 6 exactly: the coverage
   entry, the core rule with its four clauses and `replaces` lists, the
   Secondbrain module's removals, the new research rule and the two reworded
   clauses, the four workflow clauses, and the version rule. Every module and
   catalog file is edited as text, keeping every other byte and the file's
   layout (`TestCatalogDiagnosticCharacterization` reads it).
3. MUST then run, in order, the Module Version Record command and `make
   baseline-digests` declared in `_authorization.md`, and this repository's
   Managed Refresh twice with a binary built under the Run's environment and
   a fresh `HOME` each time (the second reports `File changes: 0`). No derived
   file, golden or rendered guide is hand-edited.
4. MUST NOT add a Source Baseline row, change any file under
   `internal/baseline/assets/source-baselines/`, or add a rule to any
   profile's `requiredRules`; the research clause lives in its new rule for
   that reason.
5. MUST make the declared breaks of `_techspec.md` → Testing Approach (Build
   Order 2) in `clause_characterization_test.go`,
   `grouped_sources_clauses_test.go` and `retirement_rule_clauses_test.go`,
   and no other change to existing tests.
6. MUST create `internal/baseline/cross_project_continuity_test.go` with
   `TestEveryInboxClauseAnAdopterHoldsIsReplaced` and
   `TestTheGuidesSendCrossProjectWorkThroughAiMemory` as `_techspec.md` →
   Testing Approach describes them. Each top-level test calls `t.Parallel()`
   first.
7. MUST NOT change the core module's skills, any Setup Snapshot, any profile
   source, any Go production file, any skill, any user guide page or
   `CONTEXT.md`.

## Subtasks

- [ ] Add the coverage entry, the core rule and the Secondbrain research rule; reword the kept clauses.
- [ ] Record the modules, regenerate the digests and refresh this repository twice.
- [ ] Make the declared test changes and add the two tests.

## Acceptance Criteria

- [ ] Each of the eight retired inbox clauses is `replaced` by exactly its named successor, none is unaccounted, and the kept clauses are `retained`.
- [ ] This repository's agent guides tell an Agent to send cross-project work through ai-memory, to fall back to the response when it is missing, and to consume a message only after recording it.
- [ ] No rendered guide names `Inbox Entr`, `inbox/<destination>`, `_triaged` or the `handoff` skill.

## Consumers

- the core continuity rule and the Secondbrain guide — consumers: `./internal/baseline`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
- instruction: `internal/baseline/plan.go`
- instruction: `internal/baseline/catalog_validate.go`
- instruction: `internal/baseline/backend_bucket_clause_test.go`
- instruction: `internal/baseline/catalog_test.go`
- creates: `internal/baseline/cross_project_continuity_test.go`
- interface: `internal/baseline/assets/coverage.json`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/secondbrain.json`
- interface: `internal/baseline/assets/modules/context-workflow.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/docs-layout.md`
- interface: `docs/agents/secondbrain.md`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/grouped_sources_clauses_test.go`
- interface: `internal/baseline/retirement_rule_clauses_test.go`

## Verification

- `for name in TestEveryInboxClauseAnAdopterHoldsIsReplaced TestTheGuidesSendCrossProjectWorkThroughAiMemory; do go test -list "^${name}\$" ./internal/baseline | grep -qx "$name" || { printf 'missing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 ./internal/baseline` — expected: exit 0. Before this Task neither test exists, so the command fails at its first check. Afterwards both exist and the package passes: the retention transition accounts every retired clause, the rendered guides carry the ai-memory clauses, the clause characterization, grouping and capture-lesson tests read the new text, and the derived files are current.
- `tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- 'ai-memory is optional and detected at run time' || { echo 'agent instructions lack the ai-memory rule' >&2; exit 1; }; for f in docs/agents/agent-instructions.md docs/agents/secondbrain.md docs/agents/docs-layout.md docs/agents/spec-routing.md docs/agents/skill-dispatch.md; do if grep -qF -e 'Inbox Entr' -e '_triaged' -e 'inbox/<destination>' -- "$f"; then printf '%s still names the inbox\n' "$f" >&2; exit 1; fi; done; go test -count=1 -run '^(TestBaselineClauseForceIsCharacterized|TestTheGroupingClausesCarryTheirForceAndText|TestTheGroupingClausesRenderInTheGuides|TestTheCaptureTriggerAsksForAConditionalLesson|TestCatalogDiagnosticCharacterization)$' ./internal/baseline && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task this repository's agent instructions carry no ai-memory rule, so the command fails at its first check. Afterwards the rule is rendered, no agent guide names the inbox, the named clause tests pass and the documentation contracts pass.

## References

- `_prd.md` → Goals 3-4; Core Feature 2; Success Metrics 2, 3 and 5; User Stories 2-3
- `_techspec.md` → Current behavior 3-7; Data Models 2-6; API Contracts 3-4; Invariants 2, 3 and 6; Testing Approach; Build Order 2
- ADR-0330
