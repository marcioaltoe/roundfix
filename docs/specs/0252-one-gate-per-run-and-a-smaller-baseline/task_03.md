---
task: task_03
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: pending
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
