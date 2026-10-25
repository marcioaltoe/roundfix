---
task: task_02
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
status: completed
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


## Result

Implemented this Task's catalog and guidance slice; status and declared
Verification remain Daemon-owned.

### Implementation

- Added `coverage.continuity` and the four mandatory clauses of
  `rule.core.cross-project-continuity`, with exactly the specified replacement
  lists and response fallback when ai-memory or its receiving project is missing.
- Removed the eight retired Secondbrain clauses. The research replacement lives
  in `rule.secondbrain.research-digests`; the kept prohibition now covers every
  Secondbrain file, including `inbox/`. Baseline guidance proposals are messages.
- Reworded the four workflow clauses without changing their identities or force.
  Raised each affected existing rule and guide version once; module versions
  were written by the Module Version Record command.
- Made only the three declared existing-test breaks and added the two specified
  parallel regression tests. Their boundary is the embedded catalog, actual
  Source Baseline transition, and rendered Plan postimages; no live service,
  Secondbrain filesystem, provider or network is exercised.
- Generated every derived artifact and rendered guide through the sanctioned
  commands. No skill, Setup Snapshot, Source Baseline, profile required rule,
  production Go file, other Task, Task Graph, user guide or glossary was changed.

### Acceptance evidence

1. `TestEveryInboxClauseAnAdopterHoldsIsReplaced` passed: all eight retired
   clauses have `replaced` disposition, one evidence entry and exactly their
   specified successor with mandatory force; no managed source clause is
   `unaccounted`. The kept prohibition, baseline-ownership clause and all four
   reworded workflow clauses are `retained`.
2. `TestTheGuidesSendCrossProjectWorkThroughAiMemory` passed: Plan postimages
   contain the optional-at-runtime rule, missing-server/project explanation,
   response fallback, owning-project message delivery and consume-after-recording
   contract, and dispatch both ai-memory skills. A separate inspection of the
   repository's refreshed guides confirmed the fallback, consumption contract
   and complete Secondbrain write prohibition.
3. The same guide test sweeps all five rendered Plan guides for `Inbox Entr`,
   `inbox/<destination>`, `_triaged` and the quoted `handoff` skill. A separate
   sweep of this repository's five actual guides found none of those strings.

### Commands and outcomes

- `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
  — exit 0; recorded the four changed modules.
- `make baseline-digests` — exit 0; regenerated the four formatter guides,
  profile golden digest, catalog records and four plan characterization goldens.
- `go build -o "$task0317_tmp/roundfix" ./cmd/roundfix` — exit 0, built under
  the Run environment in a private temporary directory.
- `HOME="$task0317_tmp/home-first" "$task0317_tmp/roundfix" baseline update --repo . --no-skills --no-history --yes --format text`
  — exit 0; approved postimages verified, `File changes: 5`.
- The identical update with a distinct fresh `HOME="$task0317_tmp/home-second"`
  — exit 0; approved postimages and idempotence verified, `File changes: 0`.
  Both refreshes reported the existing nested-carrier warnings for the formatter
  fixture and Source Baseline corpus; neither nested carrier was changed.
- `rtk proxy go test ./internal/baseline -run '^(TestEveryInboxClauseAnAdopterHoldsIsReplaced|TestTheGuidesSendCrossProjectWorkThroughAiMemory|TestTheGroupingClausesHaveSourceBaselineRows|TestAnAdopterRetainsTheGroupingClauses)$' -count=1`
  — exit 0 after regeneration.
- `rtk proxy go test ./internal/baseline -run '(BaselineClauseForceIsCharacterized|GroupingClausesCarryTheirForceAndText|GroupingClausesRenderInTheGuides|CaptureTriggerAsksForAConditionalLesson)' -count=1`
  — exit 0.
- Python scope/render audit — passed: all 28 staged/unstaged/untracked changed
  paths belong to this Task; Source Baseline paths and profile `requiredRules`
  are unchanged; the actual rendered guides satisfy the wording sweep.
- `rtk proxy git diff --check` — exit 0.

An earlier focused check overlapped formatter-golden regeneration and failed
compilation because embedded files changed length while the compiler read them.
After the generator exited successfully, the focused check above passed. The
initial new-test compilation also exposed that decisions are `[]DecisionValue`,
not a map; the tests now update the existing decision by ID.

### Remaining authority

The declared Verification commands, docscontract gate and repository Verification
were not run; the Daemon owns them. No terminal Task verdict, commit, push or
pull request is claimed. No follow-up outside this slice was implemented.

## Carry-forward provenance

- Source Run: `run_20261025T133232Z_0a2b2319e6c91868`
- Source commit: `d9704cb8cce31f291c7641a565fb760f65c2e8c5`
