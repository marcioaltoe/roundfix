---
task: task_03
spec: 0253-authoring-rules-that-stop-qa-reruns
status: completed
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


## Result

Implemented the Task 03 slice for Daemon Verification. Task status, checkboxes,
the Task Graph and other Task files are unchanged by this Agent. No commit,
push or Pull Request was made; authored Verification and repository
Verification were not run.

### Implementation and acceptance evidence

- The two autonomous loop clauses retain their identities, mandatory force
  and absence of `replaces`. Their additions match the TechSpec literals:
  parallel-Spec revalidation, the corrective `needs` edge before reopening,
  the wording sweep and characterization before change. A structural comparison
  against `HEAD` confirms loop-01's first sentence and network-denied sentence
  are byte-identical and every unrelated clause is unchanged.
- The new `clause.go.keep-tests-hermetic` is mandatory and follows
  `clause.go.test-observable-behavior` on one JSON line. The force
  characterization file gains exactly its one declared row.
- Version changes from this Task's starting commit: `rule.autonomous.loop`
  6 → 7, `guide.autonomous-work` 14 → 15, `rule.go.observable-tests` 3 → 4,
  and `guide.go` 4 → 5. The Module Version Record command chose module
  versions autonomous-work 14 → 15 and go 5 → 6; those module version lines
  and the record were not edited by hand.
- `TestTheLoopClausesCarryTheirText` checks each authored addition exactly
  once, mandatory enforcement, no `replaces`, and absence of the obsolete
  corrective instruction. `TestTheLoopClausesRenderInTheGuides` checks each
  literal exactly once after whitespace normalization in the Standard
  TypeScript golden and this repository's autonomous-work guide. Together
  they provide focused evidence for the first acceptance criterion.
- `TestTheHermeticGoClauseRendersInTheGoGuide` checks the catalog's mandatory
  clause, exact text and placement, and its mandatory bullet exactly once
  both in the Go adopter plan postimage and this repository's Go guide. It
  provides focused evidence for the second acceptance criterion.
- `TestAnAdopterRetainsTheLoopClauses` uses a temporary Source Baseline
  adopter with autonomous work enabled, applies that guide and ages its
  managed artifact digests. It requires a ready Managed Refresh, both loop
  clauses `retained` in the delta and retention evidence, and no
  `unaccounted` clause. The shared adopter starts with autonomous work
  disabled; the test explicitly enables it and supplies its two required
  runtime decisions before exercising retention.
- The second managed refresh and the refresh after sabotage restoration both
  exited 0 with `File changes: 0` and idempotence verified. Combined with the
  adopter test, this provides focused evidence for the third acceptance
  criterion.

### Commands and focused checks

All Go commands and regeneration used
`GOCACHE=/tmp/roundfix-task03-go-cache` with `rtk proxy`. The initial default
cache attempt was denied by the sandbox before compilation; the task-scoped
cache removed that environment blocker.

1. Before changing the modules,
   `go test ./internal/baseline -run '^(TestTheLoopClausesCarryTheirText|TestTheHermeticGoClauseRendersInTheGoGuide)$' -count=1`
   exited 1: the three loop additions were absent, the obsolete corrective
   instruction remained, and the Go clause was missing.
2. `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`
   exited 0, then `make baseline-digests` exited 0 and reported regeneration.
3. `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`
   initially could not create its Git-private transaction journal because
   that directory is outside the writable workspace. The same command with
   approved sandbox escalation exited 0, applied and verified three files.
   Its second invocation exited 0 with `File changes: 0`.
4. After the sabotages below, `make baseline-digests` exited 0 and reported
   `changed:false`; the restored sources already matched every derived file.
   The same managed-refresh command then exited 0 with `File changes: 0`.
5. Final focused command:
   `go test ./internal/baseline -run '^(TestTheLoopClausesCarryTheirText|TestTheLoopClausesRenderInTheGuides|TestTheHermeticGoClauseRendersInTheGoGuide|TestAnAdopterRetainsTheLoopClauses|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText)$' -count=1 -v`
   exited 0 with all six named tests passing, after restoration and regeneration.
6. Changed-path postflight inspected tracked and untracked paths: 18 paths,
   all in this Task's declared Context or this Task file. A comparison of the
   module objects against `HEAD` confirmed each required rule/guide version
   rose by one, the unrelated clauses are unchanged, and the Go clause is
   the only added identity. Source Baseline assets, the retention transition
   and all production Go files remain outside the diff.

### Sabotage evidence

Each sabotage changed only a declared module source, ran its focused test
with `-count=1 -v`, and restored the source byte-identically in a `finally`
block. No sabotage was recorded as a module version or left in the diff.

| Sabotage | Focused command | Observed failure |
| --- | --- | --- |
| Replaced loop-04's `Characterize current behavior before changing it, and declare each break:` with `SABOTAGED characterization:` | `go test ./internal/baseline -run '^TestTheLoopClausesCarryTheirText$' -count=1 -v` | Exit 1; loop-04 force/text differs. |
| Changed the Go hermetic clause's enforcement from mandatory to prohibited | `go test ./internal/baseline -run '^TestTheHermeticGoClauseRendersInTheGoGuide$' -count=1 -v` | Exit 1; catalog force/text differs and the Go adopter postimage lacks the mandatory bullet. |
| Changed loop-01's enforcement from mandatory to prohibited | `go test ./internal/baseline -run '^TestAnAdopterRetainsTheLoopClauses$' -count=1 -v` | Exit 1; refresh is `action_required`, with one unaccounted clause: loop-01. |

All three sources were restored, then regenerated and refreshed as recorded
above. Logs were inspected from `/tmp/task03-sabotage-loop-text.log`,
`/tmp/task03-sabotage-go-force.log` and `/tmp/task03-sabotage-retention.log`;
these are local working evidence, not durable dependencies of any test.

### Files rewritten by the sanctioned commands

Module Version Record:

- `internal/baseline/assets/modules/autonomous-work.json` (module version line)
- `internal/baseline/assets/modules/go.json` (module version line)
- `internal/baseline/module-versions.json`

`make baseline-digests`:

- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/autonomous-work.md`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`

Managed Refresh:

- `docs/agents/autonomous-work.md`
- `docs/agents/go.md`
- `docs/agents/setup-context.json`

No module version, pin, golden, snapshot or generated guide was hand-edited.
The refresh reported its existing nested-carrier warnings for the formatter
fixture and Source Baseline corpus AGENTS files; it preserved those carriers.
No follow-up work was added to this slice.
