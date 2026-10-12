---
task: task_04
spec: 0257-the-baseline-respects-the-repository
status: pending
type: backend
complexity: medium
---

# Task 04: The Baseline states how to handle a diverging vendored skill

## Overview

The Baseline forbids editing a vendored skill but says nothing about what to
do when one diverges, and two divergences on the skills maintainer's list
still have no overriding rule. This Task rewords two core clauses and one Rust
clause in place with `_techspec.md` → Exact clause texts, so the
skill-dispatch guide states the three choices in order, and states the same
order in the `roundfix` skill (ADR-0264). It is verifiable on its own: this
repository's guide renders the order, and an adopter's refresh retains the
clauses. It answers the Backlog Entry of 2026-10-12,
[the Baseline does not say how to handle a diverging vendored skill](references/2026-10-12-the-baseline-does-not-say-how-to-handle-a-diverging-vendored-skill.md).

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-12 on diverging vendored skills
   by replacing the `guidance` of `clause.core.activate-matching-skills` and
   `clause.core.prohibit-editing-vendored-skills` in the `core` module and of
   `clause.rust.read-current-docs` in the `rust` module with the exact texts
   of `_techspec.md` → Exact clause texts, editing the strings in place. Each
   clause MUST keep its `id` and enforcement and MUST NOT gain `replaces`.
2. MUST raise `rule.core.skill-dispatch` from 6 to 7 and
   `rule.rust.current-docs` from 3 to 4, from the values on this Task's
   starting commit, then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   `make baseline-digests` and
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`;
   a second refresh MUST report `File changes: 0`. The Result MUST name every
   file the commands rewrote, which the prototype measured as the files listed
   under task_04 in `_techspec.md` → Derived files. No module version, pin,
   golden, snapshot or guide is hand-edited.
3. MUST create `internal/baseline/vendored_skill_clauses_test.go` with the
   three tests of `_techspec.md` → Testing Approach, following
   `internal/baseline/verification_tier_clauses_test.go`, and MUST change the
   vendored-skill sentence that
   `internal/baseline/promoted_core_clauses_test.go` pins to the new text and
   nothing else in that file, a declared break.
4. MUST add the task_04 section of `_techspec.md` → Skill and document texts
   to `.agents/skills/roundfix/references/baseline.md` under its own heading,
   then run `make skills-sync`,
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   and `make skills-sync` again. No `### QA settlement` section of any skill
   changes.
5. MUST add **Vendored Skill** to `CONTEXT.md` through `domain-modeling` with
   the text of `_techspec.md` → Skill and document texts.
6. MUST prove each new test can fail: the Result records one sabotage of a
   clause text and one of the retention, with the test that failed for each,
   and that the source was restored and regenerated.
7. MUST NOT change any other clause, any module other than `core` and `rust`,
   the Source Baseline assets, the clause force record or any production Go
   file.

## Subtasks

- [ ] Reword the three clauses and raise the two rule versions.
- [ ] Record, regenerate and refresh twice.
- [ ] Add the clause tests and update the pinned sentence.
- [ ] State the order in the `roundfix` skill and record its version.
- [ ] Add **Vendored Skill** to `CONTEXT.md`.

## Acceptance Criteria

- [ ] `docs/agents/skill-dispatch.md` states the three choices in order and
      the installed-directory rule.
- [ ] A Source Baseline adopter's refresh records both core clauses retained
      and no clause unaccounted.
- [ ] The `roundfix` skill states the order, and its mirror matches.
- [ ] `CONTEXT.md` defines **Vendored Skill**.

## Context

- instruction: `docs/adr/0264-roundfix-never-replaces-a-vendored-skill-it-cannot-prove-older.md`
- instruction: `docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`
- instruction: `internal/baseline/verification_tier_clauses_test.go`
- instruction: `internal/baseline/stack_force_rust_test.go`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/rust.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/promoted_core_clauses_test.go`
- interface: `docs/agents/setup-context.json`
- interface: `docs/agents/skill-dispatch.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`
- creates: `internal/baseline/vendored_skill_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTheVendoredSkillClausesCarryTheDivergenceOrder|TestTheVendoredSkillClausesRenderInTheGuides|TestAnAdopterRetainsTheVendoredSkillClauses|TestTheCoreGuidesStateThePromotedRules|TestNoTwoBaselineClausesShareText|TestBaselineClauseForceIsCharacterized|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility)$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheVendoredSkillClausesCarryTheDivergenceOrder TestTheVendoredSkillClausesRenderInTheGuides TestAnAdopterRetainsTheVendoredSkillClauses TestTheCoreGuidesStateThePromotedRules TestNoTwoBaselineClausesShareText TestBaselineClauseForceIsCharacterized TestEveryBaselineModuleVersionIsRecorded TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three vendored-skill tests do not exist, so the command fails.
- `for phrase in "When such a skill diverges from what the repository needs, choose in this order." "Read a skill's bundled files from the directory where the skill is installed"; do tr -s '[:space:]' ' ' < docs/agents/skill-dispatch.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/skill-dispatch.md: %s\n' "$phrase" >&2; exit 1; }; done; grep -qF -- '## Vendored skills' .agents/skills/roundfix/references/baseline.md || { printf 'the baseline reference lacks its Vendored skills section\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/baseline.md | grep -qF -- "only when the vendored skill is useless or almost useless to the repository" || { printf 'the baseline reference lacks the third choice\n' >&2; exit 1; }; cmp -s .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md || { printf 'the baseline reference mirror differs\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "**Vendored Skill**:" || { printf 'CONTEXT.md lacks **Vendored Skill**\n' >&2; exit 1; }; go test -count=1 -run '^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync)$' ./skills || exit 1; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the skill-dispatch guide states no order, so the command fails at its first phrase.

## References

- `_prd.md` → Goal 4; User Story 5; Core Features 4 and 5; Success Metric 4
- `_techspec.md` → API Contract 7; Exact clause texts; Divergence audit; Skill and document texts; Derived files; Invariants 6 to 9; Testing Approach; Build Order 4
- ADR-0264; ADR-0257; ADR-0186; ADR-0250; ADR-0244
