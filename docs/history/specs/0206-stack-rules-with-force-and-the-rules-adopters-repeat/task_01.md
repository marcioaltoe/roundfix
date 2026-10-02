---
task: task_01
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
status: completed
type: backend
complexity: high
---

# Task 01: Recurring rules become core clauses, and lint warnings block

## Overview

Eight rules that adopters keep writing by hand become core clauses with force: lint warnings block, flaky tests block, no production hook exists only for tests, generated files are regenerated, vendored skills are not edited, and three database mutation rules. The lint clause replaces the Bun clause that blocked warnings only when Verification already treated them as errors, and declares the replacement so a Source Baseline adopter's update records it `replaced`. Each new clause gains its Source Baseline row. The slice is verified through the rendered guides of two profiles and a Managed Refresh plan for a Source Baseline adopter.

This is an authorized tooling Task. It may change only the files in its Context, the derived files the sanctioned regeneration rewrites, and this Task file.

## Requirements

1. MUST add, in `internal/baseline/assets/modules/core.json`, the eight clauses the TechSpec's "Exact texts" gives for task_01, each with its exact identifier, force and guidance, the lint clause with its `replaces` list, and MUST keep every existing clause byte-identical.
2. MUST remove `rule.bun.warning-free-verification` and its one clause from the `bun` module, from `guide.bun`, from the Standard TypeScript Monorepo profile's `requiredRules` and from the composed `go-cli-typescript-monorepo` profile's `requiredRules`, changing no other line of either profile by hand, and MUST retarget the legacy retention transition's `clause.legacy.block-warnings` mapping as the TechSpec states, keeping its `replaced` disposition.
3. MUST add the eight task_01 Source Baseline rows the TechSpec's "Source Baseline rows" table gives: corpus entry, manifest row and `entryIds` identifier, each after its named anchor. MUST NOT change or remove any existing row, including the Bun clause's row.
4. MUST raise by one, from the value on the starting main, the versions the TechSpec's "Version changes" lists for task_01, keeping each module's existing formatting.
5. MUST update the existing tests the TechSpec's "Existing tests that change" names for task_01: the force record gains the eight clauses and loses the Bun clause; the maintained Source Baseline entry count rises by eight; the TypeScript and Bun wording test drops the conditional warnings sentence from its required list and the Bun warnings row from its force table. No other existing test line changes, except the whitespace `gofmt` realigns in the force record's map once its longest key, the Bun clause, is gone; run `gofmt -w` on that file and change nothing else in it by hand.
6. MUST create `internal/baseline/promoted_core_clauses_test.go` with `expectedClause`, `promotedCoreClauses`, `clauseForceFindings` and the five tests the TechSpec's Testing Approach 1 names. The replacement test MUST assert a ready plan, the `replaced` disposition and the single target; the undeclared-replacement test MUST assert `action_required` and `unaccounted`.
7. MUST run `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`;
   a second refresh MUST report `File changes: 0`. MUST NOT hand-edit a pin, a golden, a snapshot or a generated guide.
8. MUST change no exported function signature and rename or remove no top-level test.

## Subtasks

- [ ] Add the eight core clauses and remove the conditional Bun clause with its rule.
- [ ] Retarget the legacy transition and remove the rule from both profiles.
- [ ] Add the eight Source Baseline rows.
- [ ] Update the three existing tests and create the new test file.
- [ ] Regenerate and refresh twice.

## Acceptance Criteria

- [ ] The rendered instruction guide of the Standard TypeScript Monorepo and Go CLI/TUI profiles states the seven instruction clauses with their force, and the skill guide states the vendored-skill clause.
- [ ] The rendered TypeScript and Bun guide no longer states the conditional warnings clause.
- [ ] A Source Baseline adopter's Managed Refresh plan is ready and records the Bun clause `replaced` by `clause.core.lint-warnings-block`; without `replaces` the plan is refused with the clause `unaccounted`.
- [ ] No built-in profile requires the removed rule: the catalog loads, and the composed profile's other fields are unchanged.
- [ ] The force record, the duplicate-text check, the adopter-neutral check and the catalog validation pass, and a second Managed Refresh is a no-op.

## Context

- instruction: `docs/adr/0203-a-repository-rule-becomes-a-baseline-clause-only-when-it-recurs.md`
- instruction: `docs/adr/0186-baseline-guidance-states-what-the-product-does-in-adopter-neutral-words.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/modules/bun.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json`
- interface: `internal/baseline/assets/retention/transition.legacy-typescript-bun-to-portable-v3.json`
- interface: `internal/baseline/assets/source-baselines/index.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- interface: `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/skill-dispatch.md`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `internal/baseline/preservation_test.go`
- interface: `internal/baseline/stack_rule_wording_test.go`
- interface: `docs/agents/agent-instructions.md`
- interface: `docs/agents/skill-dispatch.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/promoted_core_clauses_test.go`

## Scope after Spec 0207

Spec 0207 merged after this Task was authored and shipped the composed profile `go-cli-typescript-monorepo`, which also lists `rule.bun.warning-free-verification`. A Run stopped at `catalog.profile.rule.unknown` because that path was outside this Task. A rehearsal on 2026-10-02 at `18ef15eb`, with that Run's change plus the one-line profile edit and `gofmt` on the force record, regenerated once, converged on the second `make baseline-digests` and the second Managed Refresh, passed this Task's Verification and passed `make verify`. The partial work in the stopped Run's worktree remains usable: only the composed profile line and the formatter pass were missing.

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheCoreGuidesStateThePromotedRules|TestThePromotedCoreClausesCarryTheirForce|TestAClauseForceCheckReportsAMissingOrChangedClause|TestTheWarningsClauseIsReplacedByTheLintClause|TestAnUndeclaredWarningsReplacementIsUnaccounted|TestTheTypeScriptAndBunGuideSaysWhatItGoverns|TestTheRewordedBunAndTypeScriptClausesKeepTheirForce|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestReadoptionCompatibilityMaintainedFixture|TestFormatterComposition|TestCatalogCompatibility|TestBaselinePlanCharacterization|TestBaselineCompatibilityCorpus|TestStandardTypeScriptStructuralClauseRetention)$" ./internal/baseline 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestTheCoreGuidesStateThePromotedRules TestThePromotedCoreClausesCarryTheirForce TestAClauseForceCheckReportsAMissingOrChangedClause TestTheWarningsClauseIsReplacedByTheLintClause TestAnUndeclaredWarningsReplacementIsUnaccounted TestTheTypeScriptAndBunGuideSaysWhatItGoverns TestBaselineClauseForceIsCharacterized TestReadoptionCompatibilityMaintainedFixture TestCatalogCompatibility; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && guide=docs/agents/agent-instructions.md && for phrase in "A lint warning fails Verification." "A flaky test is a blocking failure" "exists only for tests; test through the public entry points." "Do not hand-edit a generated file" "Stop and ask for express authorization before any statement that changes data or schema" "run a read with the same predicate and report its row count" "Do not route a database change through a migration, script, seed, or test"; do tr -s '[:space:]' ' ' < "$guide" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$guide" "$phrase" >&2; exit 1; }; done && tr -s '[:space:]' ' ' < docs/agents/skill-dispatch.md | grep -qF -- "Do not edit a skill the repository installs from an upstream source" || { printf 'missing vendored-skill clause\n' >&2; exit 1; }; golden=internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/typescript-bun.md && if tr -s '[:space:]' ' ' < "$golden" | grep -qF -- "When the repository's Verification treats warnings as errors"; then printf 'stale conditional warnings clause in %s\n' "$golden" >&2; exit 1; fi && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the five new tests do not exist and the guide lacks the promoted clauses, so the command fails.

## References

- `_prd.md` → Goal 4; Goal 5; Story 4; Story 5; Core Feature 3; Success Metric 4; Success Metric 5; Success Metric 7; Declared breaks
- `_techspec.md` → Candidate rules; Exact texts (task_01); Source Baseline rows; Version changes; Existing tests that change; API Contract 1; API Contract 5; Testing Approach 1; Testing Approach 5; Build Order 1
- ADR-0058, ADR-0059, ADR-0060, ADR-0081, ADR-0099, ADR-0149, ADR-0186, ADR-0203


## Result

Implemented task_01's eight core clauses with the TechSpec's exact identifiers,
guidance and force. The lint clause declares the removed Bun clause in
`replaces`; both profiles drop the removed rule, and the legacy warning
mapping remains `replaced` with the core lint clause as its target. Added
the eight Source Baseline entries after their named anchors while retaining
the historical Bun entry. Raised the nine specified module, rule and guide
versions once from the starting revision.

Updated only the three specified existing tests, including the maintained
entry count from 149 to 157 and the force map's sanctioned formatter
realignment. Added the five specified tests and test-only helpers in
`promoted_core_clauses_test.go`. No production Go code, exported signature,
top-level test name, skill or setup snapshot changed.

### Focused evidence

All Go and Make commands below used
`GOCACHE=/private/tmp/roundfix-task01-go-cache` because the sandbox denied
access to the default Go cache.

| Acceptance criterion | Implementation and fresh evidence |
| --- | --- |
| Two profiles render the seven instruction clauses and the vendored-skill clause with force | `TestTheCoreGuidesStateThePromotedRules` passed for `standard-typescript-monorepo` and `go-cli-tui`; it checks complete force-labelled guidance in Plan postimages. |
| TypeScript and Bun guidance loses conditional warnings | The same test asserts the conditional sentence is absent; focused existing wording and force tests passed. The golden was generated by `make baseline-digests`. |
| Source Baseline replacement is ready; undeclared removal is unaccounted | `TestTheWarningsClauseIsReplacedByTheLintClause` passed, asserting `ready`, delta and retention `replaced`, and exactly one target. `TestAnUndeclaredWarningsReplacementIsUnaccounted` passed after removing `replaces` in a cloned catalog, asserting no Plan, `action_required`, classification and `unaccounted`. |
| No profile requires the removed rule; composed fields stay unchanged | Catalog validation passed in sanctioned regeneration. A read-only JSON/Git comparison swept all profiles and proved the composed profile differs only by removal of that required rule. The composed-profile rendering and convergence tests passed. |
| Force, duplicate text, adopter-neutral guidance, catalog and no-op refresh | Focused force, duplicate-text and adopter-neutral tests passed. Both regeneration runs exited 0; the second reported `changed:false`. The first authorized Managed Refresh applied three files and verified postimages; the second exited 0 with `File changes: 0` and verified idempotence. |

Commands and outcomes:

- `go test -count=1 -v -run '^(TestTheCoreGuidesStateThePromotedRules|TestThePromotedCoreClausesCarryTheirForce|TestAClauseForceCheckReportsAMissingOrChangedClause|TestTheWarningsClauseIsReplacedByTheLintClause|TestAnUndeclaredWarningsReplacementIsUnaccounted)$' ./internal/baseline` — exit 0, all five tests passed.
- `go test -count=1 -run '^(TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestShippedGuidanceCitesNoRepositoryRecord|TestTheTypeScriptAndBunGuideSaysWhatItGoverns|TestTheRewordedBunAndTypeScriptClausesKeepTheirForce|TestTheComposedProfileRendersEveryGuideWithNoRepeatedClause|TestTheComposedProfilePlanConverges)$' ./internal/baseline` — exit 0.
- `make baseline-digests` — exit 0; second run exit 0 with `changed:false`.
- `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text` — first successful run exit 0 with three applied files; second exit 0 with zero file changes. The initial sandbox attempt could not create its Git-private journal; the authorized elevated retry applied the Plan.
- Read-only preservation/scope audit — exit 0: existing core and retained Bun clause object bytes are identical to HEAD, nine versions increment once, all previous Source Baseline entries retain identity/force/carrier/structure and order, old corpus bytes remain identical after removing the eight inserted entries, and all 27 changed paths fall within this Task's Context or created-file declaration.
- `git diff --check` — exit 0.
- `make verify-incremental` — initial sandbox run exited 2: two force-stop CLI integration tests could not enumerate the process table. The CLI suite guard also detected this Result edit while the suite was running. The Baseline package and other packages passed. The elevated rerun exited 0 with no concurrent repository edits: formatter check, vet, repository tests, skill synchronization/check and build passed.

An initial regeneration caught clause objects accidentally inserted into coverage
arrays by the edit script. The source arrays were repaired to their original
values before successful regeneration. The negative force test also caught a
test-helper decode mismatch; the helper now uses the existing force reader's
supported map representation.

The refresh reports the existing nested-carrier warnings for the formatter
fixture and Source Baseline corpus; those carriers were left unchanged.
Declared Task Verification remains unrun and Daemon-owned. Status, settlement
and commits remain Daemon-owned; no commit, push or Pull Request was made.
