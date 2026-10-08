---
task: task_03
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: completed
type: backend
complexity: high
---

# Task 03: Each review, tracker and local-research rule is stated once

## Overview

The review policy appears in two core clauses and again in the autonomous
loop. The tracker rule appears in three clauses, and the local-research
prohibition appears in core and again in the Secondbrain guide. The
repository-owned skill-sync rule repeats a CLI clause (findings B04, B05, B06
and B19 of the Baseline audit of 2026-10-08,
`docs/references/2026-10-08-baseline-audit.md`). This Task merges, trims and
removes clauses with `_techspec.md` → Exact clause texts. Each removed clause
is named in the `replaces` list of the clause that absorbs it (ADR-0257). It
deletes the repository duplicate. It is verifiable on its own: each rule
renders once, and an adopter's update records every removed clause
`replaced`.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer findings B04, B05, B06 and B19 of the Baseline audit of
   2026-10-08 by applying the task_03 edits of `_techspec.md` → Exact clause
   texts to `core.json`, `autonomous-work.json`, `spec-workflow.json` and
   `secondbrain.json`. That means removing
   `clause.core.request-pull-request-review`, `clause.spec.status-only-in-task`
   and `clause.secondbrain.prohibit-external-local-discovery`, adding the three
   `replaces` lists, and replacing five `guidance` strings in place. Every
   other byte of each module stays as it is. Every kept clause keeps its `id`
   and enforcement.
2. MUST raise by one the task_03 versions of `_techspec.md` → Version changes,
   then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   `make baseline-digests` and
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. No pin, golden or guide is
   hand-edited, and the Result MUST name every file the commands rewrote.
3. MUST delete the "HARD RULE — roundfix skill sync" bullet of
   `docs/agents/specific-repository.md` whole, as `_techspec.md` → Skill and
   document texts says, and change no other byte of that file.
4. MUST create `internal/baseline/deduplicated_clauses_test.go` with
   `TestTheMergedClausesCarryTheirText`,
   `TestEveryRemovedClauseHasOneSuccessor` and
   `TestAnAdopterRecordsTheRemovedClausesReplaced`, as `_techspec.md` →
   Testing Approach describes. The texts and the three removed identities are
   literals in the test. The adopter test MUST use
   `newClauseReplacementAdopter`, require a ready Managed Refresh, record each
   removed identity `replaced` with its named successor, and find no clause
   `unaccounted`.
5. MUST change only these declared breaks:
   - drop the three removed identities from the force record in
     `internal/baseline/clause_characterization_test.go`;
   - change `TestStandardTypeScriptStructuralClauseRetention` in
     `internal/baseline/plan_test.go` so that it requires
     `clause.spec.status-only-in-task` `replaced` by
     `clause.spec.tracker-artifacts`;
   - change the structural-clause fleet fixture in
     `internal/cli/baseline_update_test.go` so that its `issue-tracker.md`
     rows name the clauses the guide now renders;
   - change the two ownership phrases that `TestLegacySpecConstraintExemption`
     in `skills/baseline_skill_contract_test.go` requires to the kept text
     "Keep completed or archived legacy Specs byte-identical.".
6. MUST prove each new gate can fail. The Result MUST record one sabotage that
   removes a `replaces` list, one that restores a removed clause's text
   elsewhere, and one that changes a merged text, each with the test that
   failed, and that the source was restored and regenerated.
7. MUST NOT change `clause.autonomous.loop-01-qa-once`,
   `clause.autonomous.loop-04-verify-the-class`, any `go` module clause, the
   Source Baseline assets, the retention transition or any production Go file.

## Subtasks

- [ ] Merge, trim and remove the clauses, and raise the versions.
- [ ] Record, regenerate and refresh twice.
- [ ] Delete the repository-owned duplicate rule.
- [ ] Add the de-duplication tests, update the four declared breaks and record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/agent-instructions.md` states the pre-PR review policy in
      one clause.
- [ ] `docs/agents/issue-tracker.md`, `docs/agents/spec-routing.md` and
      `docs/agents/secondbrain.md` no longer repeat the tracker, ownership,
      query-order and local-discovery sentences.
- [ ] `docs/agents/specific-repository.md` no longer carries the skill-sync
      HARD RULE.
- [ ] A Source Baseline adopter's refresh records the three removed clauses
      `replaced` and none `unaccounted`, and no two clauses share text.

## Context

- instruction: `docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`
- instruction: `docs/adr/0222-a-baseline-guide-says-only-what-holds-for-the-repository-that-reads-it.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- instruction: `internal/baseline/clause_replacement_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/modules/secondbrain.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/issue-tracker.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/cli/baseline_update_test.go`
- interface: `skills/baseline_skill_contract_test.go`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/autonomous-work.md`
- interface: `docs/agents/issue-tracker.md`
- interface: `docs/agents/secondbrain.md`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/specific-repository.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/deduplicated_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheMergedClausesCarryTheirText|TestEveryRemovedClauseHasOneSuccessor|TestAnAdopterRecordsTheRemovedClausesReplaced|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestStandardTypeScriptStructuralClauseRetention|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization|TestLegacySpecConstraintExemption|TestBaselineUpdateFleetSweep)$" ./internal/baseline ./skills ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheMergedClausesCarryTheirText TestEveryRemovedClauseHasOneSuccessor TestAnAdopterRecordsTheRemovedClausesReplaced TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestStandardTypeScriptStructuralClauseRetention TestEveryBaselineModuleVersionIsRecorded TestLegacySpecConstraintExemption TestBaselineUpdateFleetSweep; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three de-duplication tests do not exist, so the command fails.
- `tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- "Resolve and follow the repository's pre-PR review policy before publication" || { printf 'the merged review clause is missing\n' >&2; exit 1; }; for pair in "docs/agents/agent-instructions.md|Follow the explicit configured review policy without re-enabling an opted-out provider" "docs/agents/issue-tracker.md|Keep Task status only in the assigned Task file frontmatter" "docs/agents/spec-routing.md|Dependencies remain owned only by the Task Graph" "docs/agents/secondbrain.md|Do not use Exa or another external research tool" "docs/agents/secondbrain.md|Follow the index-first and local query workflow defined in this guide" "docs/agents/autonomous-work.md|do not manufacture a reviewed verdict" "docs/agents/specific-repository.md|HARD RULE — roundfix skill sync"; do file="${pair%%|*}"; phrase="${pair#*|}"; if tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase"; then printf 'duplicate still in %s: %s\n' "$file" "$phrase" >&2; exit 1; fi; done; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the agent-instructions guide lacks the merged review clause, so the command fails at its first check.

## References

- [_prd.md](_prd.md) — Goal 4; User Story 4; Core Feature 4; Success Metric 4
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 4; Exact clause texts; Skill and document texts; Version changes; Retention; Derived files; Invariants 1 to 4; Testing Approach; Build Order 3
- ADR-0257; ADR-0222; ADR-0186; ADR-0250

## Result

Implemented only task_03's clause consolidation. The five authored texts
replace their previous guidance in place; the three removed identities have
one same-enforcement successor each. The tracker successor keeps its existing
text. The declared rule and guide versions each rose by one; the Module
Version Record command selected and recorded the four module versions.
Only the declared characterization, structural-retention, fleet fixture and
legacy-Spec assertions changed. No production Go, Source Baseline, retention
transition, Go-module clause, loop-01 or loop-04 changed.

### Acceptance evidence

- Pre-PR review: the generated `docs/agents/agent-instructions.md` carries
  the authored merged review text exactly once. The two previous review
  texts are absent. `TestTheMergedClausesCarryTheirText` checks all five
  literal texts and their enforcement.
- Tracker, ownership and research: direct inspection of the five regenerated
  guides found each authored replacement once and the retired tracker,
  ownership, query-order, local-discovery and autonomous review sentences
  absent. `TestEveryRemovedClauseHasOneSuccessor` sweeps every module for
  removed identities and successor declarations, checks same enforcement,
  and rejects retired text restored under another identity.
- Repository rule: a byte comparison against the starting revision with only
  the whole skill-sync bullet removed matched `specific-repository.md`.
- Adopter: `TestAnAdopterRecordsTheRemovedClausesReplaced` uses
  `newClauseReplacementAdopter`, enables and records its optional Secondbrain
  guide in a temporary repository, ages managed digests, and requires a ready
  Managed Refresh. All three removed identities are `replaced` with exactly
  their named successor and explanatory evidence; none is `unaccounted`.
  `TestNoTwoBaselineClausesShareText` also passed.

### Focused checks and regeneration

The three new tests were run before the module edits and failed on all five
old texts, the three present removed identities, missing successors and
missing replacement accounting. After implementation and the sabotage
restorations, this focused command exited zero with all six tests passing:

```sh
GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/baseline -run '^(TestTheMergedClausesCarryTheirText|TestEveryRemovedClauseHasOneSuccessor|TestAnAdopterRecordsTheRemovedClausesReplaced|TestStandardTypeScriptStructuralClauseRetention|TestNoTwoBaselineClausesShareText|TestBaselineClauseForceIsCharacterized)$' -count=1 -v
```

This separate focused command also exited zero:

```sh
GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/cli ./skills -run '^(TestBaselineUpdateFleetSweep|TestLegacySpecConstraintExemption)$' -count=1
```

The required Module Version Record command, `make baseline-digests`, and
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
each exited zero. The second confirmed refresh reported `File changes: 0`
and idempotence verified. Regeneration after each sabotage exited zero.
The shared Go cache denied access during focused work; subsequent commands
used the task-local `GOCACHE` above. Baseline apply initially could not write
its Git-private transaction lock; the authorized refresh succeeded with
sandbox escalation. Existing nested-carrier warnings remained advisory.

Every file rewritten by those commands:

- `internal/baseline/assets/modules/core.json`
- `internal/baseline/assets/modules/autonomous-work.json`
- `internal/baseline/assets/modules/spec-workflow.json`
- `internal/baseline/assets/modules/secondbrain.json`
- `internal/baseline/module-versions.json`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/issue-tracker.md`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- `docs/agents/agent-instructions.md`
- `docs/agents/autonomous-work.md`
- `docs/agents/issue-tracker.md`
- `docs/agents/secondbrain.md`
- `docs/agents/spec-routing.md`
- `docs/agents/setup-context.json`

### Sabotage evidence

1. Removed the review successor's `replaces` list in `core.json`.
   `TestEveryRemovedClauseHasOneSuccessor` failed with successor count zero.
   `TestAnAdopterRecordsTheRemovedClausesReplaced` independently failed when
   planning reported `action_required`, classification, with the removed
   review identity `unaccounted`. Restored the exact source and regenerated.
2. Restored the retired status clause's literal text under the unrelated
   `clause.core.fix-root-causes` identity. `TestEveryRemovedClauseHasOneSuccessor`
   failed with `removed text restored in clause.core.fix-root-causes`.
   Restored the exact source and regenerated.
3. Changed the merged review guidance prefix in `core.json`.
   `TestTheMergedClausesCarryTheirText` failed on the literal text assertion
   for `clause.core.request-review-explicitly`. Restored the exact source and
   regenerated.

`git diff --check` passed. The changed-file postflight accounted for all
31 changed or new paths against this Task's Context and the Task file,
including the pre-existing Daemon status edit. No follow-up or implementation
blocker was found. Declared Task Verification, repository Verification and
incremental Verification were not run; the Daemon owns them. Task status and
other Task files remain untouched, and no commit, push or PR was made.

## Carry-forward provenance

- Source Run: `run_20261008T180729Z_1c5cbdf685b778d0`
- Source commit: `43213026fbebe2037374a8888876758aba47e000`
