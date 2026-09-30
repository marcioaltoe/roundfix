---
task: task_06
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: completed
type: backend
complexity: medium
---

# Task 06: A removed Source Baseline clause is replaced, not lost

## Overview

task_01 removed the duplicate entry `rule.backend.boundary-contracts` from `backend.json`. The Standard TypeScript Source Baseline records that entry, and `classifySourceClauseTransition` in `internal/baseline/plan.go` knows only `retained` and `unaccounted`. An adopter whose Setup Manifest declares that Source Baseline now gets the removed entry reported `unaccounted`, and its next Baseline update is refused ("retention transition has 1 unaccounted clause(s)"). The first delivery Run's QA gate failed on `TestStandardTypeScriptStructuralClauseRetention`, the test that guards this. The same gate failed on `TestADRLifecycleContract`, which still pins the ADR sentence task_04 scoped on purpose. This Task gives the removal a declared successor and moves both contracts to the facts this Spec establishes. The TechSpec section "Clause replacement for adopters (task_06)" fixes the design.

## Requirements

1. MUST accept on a module clause an optional `replaces` field: a list of clause IDs that clause takes over.
2. MUST make `classifySourceClauseTransition` report a Source Baseline normative clause that is absent from the selected catalog as `replaced` (the existing `ClauseReplaced`), with the replacing clause's ID as its only target and a reason naming it, when exactly one selected clause lists it in `replaces` and that clause's enforcement equals the Source Baseline entry's. Every other absence MUST stay `unaccounted`, with today's reason. Retained clauses MUST keep today's classification and reason.
3. MUST make catalog validation refuse, each with its own diagnostic code under `catalog.clause.replaces.`: a `replaces` value that is not a list of unique non-empty strings; an ID that is still a clause ID anywhere in the catalog; an ID claimed by two clauses.
4. MUST add `"replaces": ["rule.backend.boundary-contracts"]` to `clause.backend.boundary-contracts` in `internal/baseline/assets/modules/backend.json`, change nothing else in that clause, and render no guide differently.
5. MUST change `TestStandardTypeScriptStructuralClauseRetention` in `internal/baseline/plan_test.go` so that `rule.backend.boundary-contracts` is expected `replaced`, targeting `clause.backend.boundary-contracts`, whose enforcement and guidance equal the Source Baseline entry's. Every other clause the test lists MUST stay expected `retained`, with its current checks.
6. MUST change the fragment `Only \`accepted\` is active.` in `TestADRLifecycleContract` to `For an ADR that carries lifecycle frontmatter, only \`accepted\` is active.`, the sentence task_04 wrote. It MUST change no other fragment or fixture of that test.
7. MUST add `internal/baseline/clause_replacement_test.go` with:
   - `TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor`: a Standard TypeScript adopter manifest whose managed artifacts drifted produces a plan that is not refused, and whose delta records `rule.backend.boundary-contracts` as `replaced`;
   - `TestAnUndeclaredRemovalStaysUnaccounted`: the same catalog with the `replaces` declaration removed reports `unaccounted`, and the plan is refused;
   - `TestAReplacementWithOtherForceStaysUnaccounted`: a declaration whose clause has a different enforcement reports `unaccounted`;
   - `TestClauseReplacementDeclarationsAreValidated`: the three refusals of Requirement 3, each by its code.
8. MUST regenerate the catalog snapshots, digest pins and plan goldens only through their sanctioned commands (`make baseline-digests` and the package's golden update path), never by hand, and MUST leave a second refresh reporting no change.
9. MUST NOT edit the Source Baseline corpus or index, any other governed test file, or any rendered guide.

## Subtasks

- [ ] Characterize: run the two failing tests and record the failure.
- [ ] Add the `replaces` producer and its validation.
- [ ] Declare the successor in `backend.json` and regenerate derived files.
- [ ] Move the two plan contracts and add the four new tests.

## Acceptance Criteria

- [ ] The whole `./internal/baseline` package passes, including the two tests the first Run failed.
- [ ] Removing the `replaces` declaration makes `TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor` fail.
- [ ] `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json` exits 0.

## Context

- instruction: `docs/specs/0193-baseline-guides-that-describe-the-product-as-it-is/_techspec.md`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/catalog_validate.go`
- interface: `internal/baseline/assets/modules/backend.json`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/clause_replacement_test.go`

## Verification

- `out="$(go test -count=1 -v ./internal/baseline 2>&1)" || { printf "%s\\n" "$out" | grep -E -- "--- FAIL|_test.go:" ; exit 1; }; for name in TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor TestAnUndeclaredRemovalStaysUnaccounted TestAReplacementWithOtherForceStaysUnaccounted TestClauseReplacementDeclarationsAreValidated TestStandardTypeScriptStructuralClauseRetention TestADRLifecycleContract TestBaselinePlanCharacterization TestCatalogCompatibility; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json >/dev/null` — expected: exit 0; before this Task the package fails on `TestStandardTypeScriptStructuralClauseRetention` and `TestADRLifecycleContract`, and the four new tests do not exist.

## References

- task_01, task_04
- `_techspec.md` → Clause texts; Build Order
- QA Report of Run `run_20260930T153457Z_b5736250c995b99a` (verification log `batch-005-attempt-1.log`)

## Result

Implemented the declared successor for `rule.backend.boundary-contracts`.
The Source Baseline classifier now emits `replaced` only for an absent clause
with exactly one selected successor carrying the same enforcement. Its sole
target and reason name that successor. Retained classifications and the
existing unaccounted reason remain unchanged.

Catalog validation rejects malformed or duplicate replacement lists with
`catalog.clause.replaces.invalid`, a target still present anywhere in the
catalog with `catalog.clause.replaces.present`, and competing claimants with
`catalog.clause.replaces.duplicate`. The backend clause gained only its
`replaces` field. The structural-retention contract checks the successor's
Source Baseline enforcement and guidance, and the ADR contract changed only
the requested sentence. The four new tests exercise adopter refresh,
undeclared removal, mismatched enforcement, and all three validation codes.

Focused evidence from this turn (Go commands used
`GOCACHE=/private/tmp/roundfix-task06-cache` and `rtk proxy`):

- Before implementation, `go test ./internal/baseline -run
  '^(TestStandardTypeScriptStructuralClauseRetention|TestADRLifecycleContract)$'
  -count=1` exited 1. The removed clause was `unaccounted` and missing from the
  selected catalog; the ADR contract required the old sentence.
- After the final source edits, `go test ./internal/baseline -run
  '^(TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor|TestAnUndeclaredRemovalStaysUnaccounted|TestAReplacementWithOtherForceStaysUnaccounted|TestClauseReplacementDeclarationsAreValidated|TestStandardTypeScriptStructuralClauseRetention|TestADRLifecycleContract)$'
  -count=1` exited 0.
- `make baseline-digests` regenerated catalog snapshots and plan goldens through
  the sanctioned update paths. A final second invocation exited 0 with
  `changed:false`. An intermediate JSON layout edit changed the raw catalog
  identity; regeneration from the final bytes repaired the stale expectations
  exposed by the first incremental run. No pin or golden was hand-edited.
- The authorized `go run -buildvcs=false ./cmd/roundfix baseline update --repo .
  --no-skills --yes --format text` exited 0 and updated only the Setup Manifest's
  catalog digest. Its second invocation exited 0 with `File changes: 0` and
  verified idempotence. The first sandboxed apply attempt could not create its
  Git-private transaction directory; the authorized retry with host access
  succeeded. No rendered guide, Source Baseline corpus/index, other task, or
  Task Graph was changed.
- `git -c core.fsmonitor=false diff --check` exited 0.

Acceptance evidence:

| Criterion | Evidence and remaining Daemon check |
| --- | --- |
| Whole Baseline package, including the two original failures | The final `make verify-incremental` run reported `ok roundfix/internal/baseline` (78.676s). The six focused tests also exited 0. Authored Verification remains Daemon-owned. |
| Removing the declaration breaks the positive adopter test | `go test -overlay /private/tmp/task06-mutation-overlay.json ./internal/baseline -run '^TestARemovedSourceClauseIsReplacedByItsDeclaredSuccessor$' -count=1` exited 1. The temporary overlay removed `replaces` from the test's catalog copy before planning; the positive assertion reported `action_required` and `1 unaccounted clause(s): rule.backend.boundary-contracts`. Repository source remained unchanged. |
| Public update exits 0 | The authorized apply and second refresh above both exited 0; the second proposed zero file changes. The exact unconfirmed JSON command in Verification was not run and remains for the Daemon. |

Follow-up outside this slice: the final `make verify-incremental` exited 2
because `TestBaselineUpdateFleetSweep/structural-clauses-missing` in
`internal/cli/baseline_update_test.go:686` still expects the backend boundary
paragraph twice (`occur 1 times ... want 2`). This contradicts task_01's
intentional duplicate removal. That CLI fixture was not edited. The initial
incremental run also encountered sandbox process-table refusals in two
force-stop tests; those failures cleared in the final run with host access.
The CLI fixture mismatch was the only final incremental failure.

Task status is unchanged. No authored Verification command, commit, push, or
pull request was performed. This Result records implementation evidence for
Daemon Verification and settlement, without assigning a terminal Task verdict.

## Carry-forward provenance

- Source Run: `run_20260930T184402Z_6aa3e4970df8def4`
- Source commit: `981aabafbc705f7c435daa84be36446fa8fe87c6`
