---
task: task_03
spec: 0253-authoring-rules-that-stop-qa-reruns
status: pending
type: backend
complexity: medium
---

# Task 03: The loop re-checks parallel Specs and sweeps old wording, and Go tests stay hermetic

## Overview

Specs authored in parallel needed hand fixes after a sibling merged an ADR, a
corrective Task needed a hand-added `needs` edge before `roundfix reopen` saw
it, stale wording in a second guide failed QA or CI four times, and tests that
read a host key or bound a long macOS socket path parked a Run or forced a
corrective Task (findings B10, B11, B13 and B14 and briefing rule 9 of the
Baseline audit of 2026-10-08, `docs/references/2026-10-08-baseline-audit.md`).
This Task extends the two loop clauses in place and adds one Go clause with
`_techspec.md` → Clause changes (ADR-0258). It is verifiable on its own: this
repository's autonomous-work and Go guides carry the new text, and an
adopter's update retains the loop clauses.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file.

## Requirements

1. MUST answer findings B10, B13 and B14 of the Baseline audit of 2026-10-08
   by editing the `guidance` of `clause.autonomous.loop-01-qa-once` and
   `clause.autonomous.loop-04-verify-the-class` in
   `internal/baseline/assets/modules/autonomous-work.json` exactly as
   `_techspec.md` → Clause changes (task_03) says. Each keeps its `id`, its
   `mandatory` enforcement and no `replaces` list. The first sentence of
   loop-01 and its network-denied sentence MUST stay byte-identical.
2. MUST answer finding B11 of the Baseline audit of 2026-10-08 by adding
   `clause.go.keep-tests-hermetic` (mandatory) with the text of
   `_techspec.md` → Clause changes to `rule.go.observable-tests` in
   `internal/baseline/assets/modules/go.json`, after
   `clause.go.test-observable-behavior`, on one line in the style of its
   neighbours.
3. MUST raise by one, from the value on this Task's starting commit, the
   task_03 versions of `_techspec.md` → Version changes. It MUST then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   then `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. It MUST NOT hand-edit a
   module version, pin, golden, snapshot or generated guide, and the Result
   MUST name every file the commands rewrote.
4. MUST create `internal/baseline/loop_and_go_clauses_test.go` with
   `TestTheLoopClausesCarryTheirText`, `TestTheLoopClausesRenderInTheGuides`,
   `TestTheHermeticGoClauseRendersInTheGoGuide` and
   `TestAnAdopterRetainsTheLoopClauses`. The added texts are literals in the
   test. The text test MUST also require that loop-01 no longer carries
   "reopen the gate with `roundfix reopen --spec <slug>`; never edit the QA
   Task by hand". The rendering test MUST require each literal exactly once,
   on whitespace-normalized text, in the Standard TypeScript golden and this
   repository's `docs/agents/autonomous-work.md`. The Go test MUST require the
   clause in the catalog as `mandatory`, and the line
   `- **mandatory**: <text>` exactly once in the `docs/agents/go.md`
   postimage of `buildTestPlan(t, newPlanRepository(t))` and in this
   repository's `docs/agents/go.md`. The retention test MUST require a ready
   Managed Refresh of a Source Baseline adopter that records both loop
   clauses `retained` and no clause `unaccounted`.
5. MUST add `"clause.go.keep-tests-hermetic": "mandatory"` to
   `characterizedBaselineForce` in
   `internal/baseline/clause_characterization_test.go` and nothing else in
   that file. That is a declared break.
6. MUST prove each new gate can fail. The Result MUST record one sabotage of a
   loop text, one of the Go clause's enforcement and one of the retention,
   with the test that failed for each, and that the source was restored and
   regenerated.
7. MUST NOT change any other clause, the Source Baseline assets, the retention
   transition or any production Go file.

## Subtasks

- [ ] Extend the two loop clauses, add the Go clause and raise the versions.
- [ ] Record, regenerate and refresh twice.
- [ ] Add the clause test file and the force row.
- [ ] Record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/autonomous-work.md` carries the parallel re-check, the
      corrective `needs` edge, the wording sweep and characterization before
      change.
- [ ] `docs/agents/go.md` carries the hermetic test clause as a `mandatory`
      bullet.
- [ ] A Source Baseline adopter's refresh retains both loop clauses, and a
      second refresh of this repository is a no-op.

## Context

- instruction: `docs/adr/0258-a-spec-is-authored-against-the-qa-rerun-classes-and-a-retirement-writes-reduced-history.md`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- instruction: `internal/baseline/glossary_clauses_test.go`
- instruction: `internal/baseline/stack_force_go_cli_tui_test.go`
- interface: `internal/baseline/assets/modules/autonomous-work.json`
- interface: `internal/baseline/assets/modules/go.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/clause_characterization_test.go`
- interface: `docs/agents/autonomous-work.md`
- interface: `docs/agents/go.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/loop_and_go_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheLoopClausesCarryTheirText|TestTheLoopClausesRenderInTheGuides|TestTheHermeticGoClauseRendersInTheGoGuide|TestAnAdopterRetainsTheLoopClauses|TestLoopClauseNamesTheDeliveryQueueAndItsRecoveryActs|TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow|TestTheGoCLIAndTUIGuidesStateTheirClauses|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheLoopClausesCarryTheirText TestTheLoopClausesRenderInTheGuides TestTheHermeticGoClauseRendersInTheGoGuide TestAnAdopterRetainsTheLoopClauses TestLoopClauseNamesTheDeliveryQueueAndItsRecoveryActs TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestEveryBaselineModuleVersionIsRecorded TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Loop' ./internal/delivery || exit 1` — expected: exit 0; before this Task the loop and Go clause tests do not exist, so the command fails.
- `for phrase in "Before a Spec authored beside another one enters a Run or the Delivery Queue" "add it to the QA Task's" "searches the repository's documentation, agent guides, and skill copies for the old wording" "Characterize current behavior before changing it, and declare each break"; do tr -s '[:space:]' ' ' < docs/agents/autonomous-work.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/autonomous-work.md: %s\n' "$phrase" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/agents/go.md | grep -qF -- "because macOS refuses a socket path longer than 104 bytes" || { printf 'docs/agents/go.md lacks the hermetic test clause\n' >&2; exit 1; }; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the autonomous-work guide lacks the parallel re-check, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goal 4; User Story 4; Core Feature 4; Success Metric 2; Acceptance evidence
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; Clause changes; Version changes; Retention; Derived files; Invariants 1-6 and 8; Testing Approach; Build Order 3
- ADR-0258; ADR-0257; ADR-0250; ADR-0186
