---
task: task_04
spec: 0257-the-baseline-respects-the-repository
status: completed
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

## Result

Implementation handed back for Daemon Verification. Task status, Subtasks,
Acceptance Criteria, the Task Graph and every other Task file are unchanged.

### Implementation

- Reworded the two core clauses and the Rust current-docs clause with the
  TechSpec's exact texts, keeping their identities and enforcement and adding
  no `replaces`. Only `rule.core.skill-dispatch` (6 → 7) and
  `rule.rust.current-docs` (3 → 4) changed rule versions. The sanctioned module
  record selected core module version 21 and Rust module version 5.
- Added `internal/baseline/vendored_skill_clauses_test.go`, following the
  verification-tier clause tests. The tests pin all three texts independently
  of the assets, verify each core text renders once with its enforcement, and
  exercise Source Baseline refresh accounting through
  `newClauseReplacementAdopter`.
- Changed only the declared vendored-skill sentence in
  `internal/baseline/promoted_core_clauses_test.go`.
- Added the exact `## Vendored skills` section to the canonical `roundfix`
  Baseline reference and synchronized its mirror. The owned-skill record
  selected version `0.1.66`; both version fields in each SKILL.md were written
  by the record command. The QA settlement section is byte-identical to HEAD.
- Added **Vendored Skill** beside the Repository Skill Set and Setup Snapshot
  Difference terms through `domain-modeling`, using the specified definition,
  avoid-list and a resolving Markdown link to ADR-0264.

### Acceptance evidence

| Acceptance criterion | Evidence observed in this turn |
| --- | --- |
| Skill-dispatch states the three choices and installed-directory rule | `TestTheVendoredSkillClausesCarryTheDivergenceOrder` and `TestTheVendoredSkillClausesRenderInTheGuides` passed. A direct inspection also found exactly one correctly forced bullet for each exact core text in both the repository guide and Standard TypeScript golden. |
| Source Baseline refresh retains both core clauses, none unaccounted | `TestAnAdopterRetainsTheVendoredSkillClauses` passed for both core clause identities. It checks every disposition for `ClauseUnaccounted`, each requested disposition for `ClauseRetained`, and the corresponding retention evidence. |
| Roundfix skill states the order and matches its mirror | The canonical reference contains the exact three-choice section. Its bytes equal `skills/roundfix/references/baseline.md`; the sanctioned regeneration's `TestAuthorialSkillSync` passed after synchronization. |
| CONTEXT defines Vendored Skill | Direct inspection confirmed the exact definition and avoid-list beside related skill terms, with the ADR-0264 link. |

### Commands and focused checks

- `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`: exit 0; recorded the module versions.
- `make skills-sync`, then
  `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
  then `make skills-sync` again: each exited 0.
- `make baseline-digests`: exit 0 after synchronization. The first invocation
  had stopped at `TestAuthorialSkillSync` because the canonical reference had
  been edited before its mirror was synchronized; the required sync/record/sync
  sequence corrected that cause. No expected bytes were edited by hand.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`:
  first refresh exited 0 and reported `File changes: 2`, updating the Setup
  Manifest and skill-dispatch guide. The second exited 0 and reported
  `File changes: 0` and `Idempotence: verified`. A final invocation after
  restoring the sabotage also exited 0 with the same zero-change result.
  Refresh reported the existing nested-carrier conflict warnings for the
  formatter-fixture and Source Baseline AGENTS.md carriers; neither was changed.
- `go test ./internal/baseline -run '^(TestTheVendoredSkillClausesCarryTheDivergenceOrder|TestTheVendoredSkillClausesRenderInTheGuides|TestAnAdopterRetainsTheVendoredSkillClauses|TestTheCoreGuidesStateThePromotedRules|TestTheRustGuideStatesTheErrorPolicy)$' -count=1 -v`:
  exit 0 on the restored final sources; all five named tests reported PASS.
- Direct structural comparison to HEAD confirmed the modules differ only in
  the three guidance strings, two rule versions and recorded module versions.
  The final changed-file audit found only Task Context paths plus this Result;
  `git diff --check` exited 0.

### Failure proofs and restoration

1. Before editing the modules, the focused text/rendering check
   `go test ./internal/baseline -run '^TestTheVendoredSkillClauses(CarryTheDivergenceOrder|RenderInTheGuides)$' -count=1`
   exited 1: all three old texts differed, and both core golden bullets lacked
   the new text.
2. Clause sabotage: temporarily changed `First, state the correction` to
   `First, ignore the correction` in the core source, then ran
   `make baseline-digests` successfully. The same focused text/rendering
   command exited 1: `TestTheVendoredSkillClausesCarryTheDivergenceOrder`
   reported a text mismatch and `TestTheVendoredSkillClausesRenderInTheGuides`
   reported the missing exact forced bullet, both for
   `clause.core.prohibit-editing-vendored-skills`.
3. Restored the exact core source bytes and ran `make baseline-digests` again,
   exit 0. No module-record command was run on the sabotaged content, so no
   temporary version entered the record.
4. Retention sabotage: a temporary Go test overlay removed that clause from
   the refresh catalog with `mutateCatalogClause` after the Source Baseline
   adopter was created. `go test -overlay /private/tmp/roundfix-task04-evidence/retention-overlay.json ./internal/baseline -run '^TestAnAdopterRetainsTheVendoredSkillClauses$' -count=1`
   exited 1. The test reported `action_required`, category `classification`,
   and `retention transition has 1 unaccounted clause(s): clause.core.prohibit-editing-vendored-skills`.
   The repository test source was unchanged by the overlay; both temporary
   overlay files were removed. The restored focused command above then
   reported PASS for all three new tests and both existing checks.

### Files changed by sanctioned commands

Module record command:

- `internal/baseline/assets/modules/core.json` (module version line)
- `internal/baseline/assets/modules/rust.json` (module version line)
- `internal/baseline/module-versions.json`

`make baseline-digests`:

- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`

Managed Refresh:

- `docs/agents/setup-context.json`
- `docs/agents/skill-dispatch.md`

Owned-skill record and synchronization (files whose bytes changed):

- `.agents/skills/roundfix/SKILL.md`
- `skills/roundfix/SKILL.md`
- `skills/roundfix/references/baseline.md`
- `skills/testdata/owned-skill-versions.json`

The remaining authored changes are the canonical reference
`.agents/skills/roundfix/references/baseline.md`, `CONTEXT.md`, the new test
file, the declared sentence in `promoted_core_clauses_test.go`, and this Result;
module guidance and rule versions were edited at their sources before recording.
No production Go file, Source Baseline asset, clause force record, other
module, other Task or Task Graph changed.

### Handoff boundary

The Task's two declared Verification commands, `make verify` and
`make verify-changed` were not run; they remain Daemon-owned. No commit, push
or pull request was made. No implementation blocker or follow-up was found.
