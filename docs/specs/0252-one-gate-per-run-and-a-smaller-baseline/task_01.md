---
task: task_01
spec: 0252-one-gate-per-run-and-a-smaller-baseline
status: pending
type: backend
complexity: medium
---

# Task 01: A Daemon-assigned Task runs focused tests and leaves the full gate to the Daemon

## Overview

The core and Spec workflow clauses tell every Task agent to run the selected
incremental Verification before handoff, and this repository's incremental
command is the full suite. The Baseline audit of 2026-10-08 counted 74 such
agent runs (284 agent-minutes) in Specs 0225–0248, beside a Daemon gate that
already ran at every settlement (findings B01 and B15,
`docs/references/2026-10-08-baseline-audit.md`). This Task rewords the four
clauses in place with `_techspec.md` → Exact clause texts, aligns
`implement-task` §7 and revises **Incremental Verification** (ADR-0257). It is
verifiable on its own: this repository's guides scope both tiers, and an
adopter's update retains the four clauses.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer findings B01 and B15 of the Baseline audit of 2026-10-08 by
   replacing the `guidance` of `clause.core.run-selected-verification` and
   `clause.core.verification-two-tiers` in
   `internal/baseline/assets/modules/core.json`, and of
   `clause.spec.verification-two-tiers` and
   `clause.spec.verification-fails-before-the-change` in
   `internal/baseline/assets/modules/spec-workflow.json`, with the exact texts
   of `_techspec.md` → Exact clause texts. It edits the strings in place, so
   every other byte stays as it is. Each clause MUST keep its `id` and its
   `mandatory` enforcement and MUST NOT gain a `replaces` list.
2. MUST raise by one, from the value on this Task's starting commit, the
   task_01 versions of `_techspec.md` → Version changes. It MUST then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   then `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. It MUST NOT hand-edit a
   module version, pin, golden, snapshot or generated guide, and the Result
   MUST name every file the commands rewrote.
3. MUST replace step 3 of `.agents/skills/implement-task/SKILL.md` §7 with the
   text of `_techspec.md` → Skill and document texts. It MUST then run
   `make skills-sync` and
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and `make skills-sync` again, so that both version fields rise and the
   mirror matches. No other line of the skill and no `### QA settlement`
   section of any skill changes.
4. MUST revise the **Incremental Verification** entry of `CONTEXT.md` through
   `domain-modeling` to the text of `_techspec.md` → Skill and document texts,
   keeping its `_Avoid_` line.
5. MUST create `internal/baseline/verification_tier_clauses_test.go` with
   `TestTheVerificationTierClausesCarryTheirText`,
   `TestTheVerificationTierClausesRenderInTheGuides` and
   `TestAnAdopterRetainsTheVerificationTierClauses`, following
   `internal/baseline/glossary_clauses_test.go`. The four texts are literals in
   the test. The rendering test MUST require each one exactly once as a
   `mandatory` bullet in the Standard TypeScript golden guide that renders it
   (`agent-instructions.md` or `spec-routing.md`). The retention test MUST
   require a ready Managed Refresh of a Source Baseline adopter that records
   all four `retained` and no clause `unaccounted`.
6. MUST change the sentence that
   `internal/baseline/promoted_spec_and_typescript_clauses_test.go` pins for
   `clause.spec.verification-fails-before-the-change` to the new text and
   nothing else in that file. That is a declared break.
7. MUST prove each new gate can fail. The Result MUST record one sabotage of a
   clause text and one of the retention, for example a changed enforcement,
   with the test that failed for each, and that the source was restored and
   regenerated.
8. MUST NOT change any other clause, the Source Baseline assets, the retention
   transition, the force record or any production Go file.

## Subtasks

- [ ] Reword the four clauses and raise the versions.
- [ ] Record, regenerate and refresh twice.
- [ ] Change `implement-task` §7, record its version and sync the mirror.
- [ ] Revise **Incremental Verification** in `CONTEXT.md`.
- [ ] Add the clause test file, update the pinned sentence and record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/agent-instructions.md` and `docs/agents/spec-routing.md`
      scope both tiers to work outside a Run and name focused tests for a
      Daemon-assigned turn.
- [ ] `docs/agents/spec-routing.md` names `--strict --run-verification`.
- [ ] `implement-task` §7 forbids the selected commands and the full suite, and
      its mirror matches.
- [ ] `CONTEXT.md` carries the revised **Incremental Verification**.
- [ ] A Source Baseline adopter's refresh retains the four clauses, and a
      second refresh of this repository is a no-op.

## Context

- instruction: `docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- instruction: `internal/baseline/glossary_clauses_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/promoted_spec_and_typescript_clauses_test.go`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `skills/implement-task/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- creates: `internal/baseline/verification_tier_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheVerificationTierClausesCarryTheirText|TestTheVerificationTierClausesRenderInTheGuides|TestAnAdopterRetainsTheVerificationTierClauses|TestTwoTierClausesNameTheSelectedCommands|TestTheSpecAndTypeScriptGuidesStateThePromotedRules|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheVerificationTierClausesCarryTheirText TestTheVerificationTierClausesRenderInTheGuides TestAnAdopterRetainsTheVerificationTierClauses TestTwoTierClausesNameTheSelectedCommands TestTheSpecAndTypeScriptGuidesStateThePromotedRules TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestEveryBaselineModuleVersionIsRecorded TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three verification-tier tests do not exist, so the command fails.
- `for phrase in "Outside a Run, run the selected repository Verification before a completion claim" "In a Daemon-assigned turn, run focused tests of the changed packages and neither selected command"; do tr -s '[:space:]' ' ' < docs/agents/agent-instructions.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/agent-instructions.md: %s\n' "$phrase" >&2; exit 1; }; done; for phrase in "a Daemon-assigned Task hands back after focused tests" "--strict --run-verification"; do tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/spec-routing.md: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < .agents/skills/implement-task/SKILL.md | grep -qF -- "its incremental Verification, or the full test suite" || { printf 'implement-task lacks the handoff rule\n' >&2; exit 1; }; cmp -s .agents/skills/implement-task/SKILL.md skills/implement-task/SKILL.md || { printf 'the implement-task mirror differs\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "for validating the current change outside a Run while reusing safe local state" || { printf 'CONTEXT.md lacks the revised **Incremental Verification**\n' >&2; exit 1; }; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' ./skills || exit 1; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the agent-instructions guide asks every Task for the incremental tier, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goal 1; User Story 1; Core Features 1 and 6; Success Metric 1
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 5; Exact clause texts; Skill and document texts; Version changes; Retention; Derived files; Testing Approach; Build Order 1
- ADR-0257; ADR-0250; ADR-0244; ADR-0186
