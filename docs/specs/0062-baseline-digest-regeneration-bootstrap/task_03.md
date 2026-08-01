---
task: task_03
spec: 0062-baseline-digest-regeneration-bootstrap
status: completed
type: backend
complexity: high
---

# Task 03: Break the regeneration cycle on the update path

## Overview

This is the slice that fixes the reported defect. The regeneration tests load
the catalog strictly, so a module edit that invalidates a derived pin makes the
load refuse and the refresh never runs — the command's remediation is itself.
This Task points the update path at the regeneration-mode entry point and adds
the fixture that reproduces the cycle, which fails before the change and passes
after it.

## Requirements

1. MUST make every regeneration step that runs under the update flag load the
   catalog through the regeneration-mode entry point, so a stale derived pin no
   longer blocks the run that rewrites it.
2. MUST leave every non-update test and every production path loading strictly.
3. MUST add a fixture that reproduces the reported cycle — a catalog whose
   generated guide changed, invalidating its formatter golden digest — and
   prove regeneration completes against it.
4. MUST prove the same fixture still fails strict validation when loaded
   outside regeneration mode, so the deferral is scoped and not a hole.
5. MUST leave the diagnostics of every other catalog unchanged.

## Subtasks

- [ ] Point the update-flag regeneration steps at the regeneration entry point.
- [ ] Add the cycle fixture reproducing the invalidated derived pin.
- [ ] Assert regeneration completes against the fixture.
- [ ] Assert the same fixture still fails strict validation.

## Acceptance Criteria

- [ ] The cycle fixture completes regeneration without the formatter digest
      mismatch stopping it.
- [ ] The same fixture, loaded outside regeneration mode, still produces the
      mismatch diagnostic.
- [ ] Regeneration steps that do not run under the update flag still load
      strictly.
- [ ] The characterization corpus from task 01 is unchanged.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/catalog_test.go`
- interface: `internal/baseline/catalog_load.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestRegenerationBreaksGoldenDigestCycle -count=1`
  — expected: exit 0; the fixture that reproduces the reported cycle now
  regenerates instead of refusing.
- `go test ./internal/baseline -run TestCatalogRegenerationMode -count=1` —
  expected: exit 0; scoping from task 02 still holds.
- `go test ./internal/baseline -run TestCatalogDiagnosticCharacterization -count=1`
  — expected: exit 0; no other catalog's diagnostics moved.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 2; Success Metrics (adding a clause and running the
  regeneration once succeeds).
- `_techspec.md` → Testing Approach: cycle regression fixture; Build Order 3.
- ADR-0085.

## Result

Implemented the update-path wiring without changing the strict public path.
The regeneration-only catalog loader now accepts any package-owned asset
filesystem, while `LoadCatalog` and `LoadEmbeddedCatalog` retain their strict
behavior. The maintained Source Baseline and catalog-compatibility refreshes
load their on-disk assets through regeneration mode. Formatter composition
loads the embedded catalog through regeneration mode only under `-update` and
passes that catalog through package-local planning and plan validation;
ordinary `BuildPlan` and `ValidatePlanDocument` still acquire the strict
embedded catalog themselves.

Added `TestRegenerationBreaksGoldenDigestCycle` with one catalog fixture whose
generated backend guide differs from the formatter digest pin. Separate
subtests prove strict loading reports
`catalog.profile.formatter.goldenDigest.mismatch` and regeneration loading
accepts the same fixture.

Focused checks:

- Red signal: before adding the regeneration filesystem entry point, the
  regeneration subtest failed to compile at
  `catalog_test.go:166:19: undefined: loadCatalogForRegeneration`.
- `rtk gofmt -w internal/baseline/catalog.go internal/baseline/catalog_test.go
  internal/baseline/plan.go internal/baseline/plan_test.go` — exit 0.
- Each `TestRegenerationBreaksGoldenDigestCycle` subtest was run separately
  with `-count=1` and a writable local `GOCACHE`; both exited 0.
- Focused `TestFormatterComposition` runs without and with `-update` each
  exited 0. The update run exercised the regeneration catalog through planning
  and its package-local self-validation; the ordinary run exercised public,
  strict `BuildPlan`.
- Focused `TestReadoptionCompatibilityMaintainedFixture -update` and
  `TestCatalogCompatibility -update` runs each exited 0, exercising the other
  catalog-loading regeneration steps.
- Focused ordinary `TestCatalogCompatibility` exited 0, and focused
  `TestFileChangesProjectionRejectsMismatch` exited 0 after the plan-validation
  split.
- The non-allowlisted-diagnostic subtest from
  `TestCatalogRegenerationMode` exited 0.
- `rtk git diff --exit-code --
  internal/baseline/testdata/catalog.diagnostics.golden.json` — exit 0.
- `rtk git diff --check` — exit 0.

Acceptance evidence:

1. The cycle fixture's regeneration subtest accepted the changed generated
   guide, and all three catalog-loading `-update` steps exited 0 through the
   regeneration entry points.
2. The strict subtest loaded the same fixture through `LoadCatalog` and
   observed `catalog.profile.formatter.goldenDigest.mismatch`.
3. The ordinary formatter-composition and catalog-compatibility focused runs
   exited 0 through public strict loaders. Source inspection also confirms
   public `BuildPlan` and public `ValidatePlanDocument` still call
   `LoadEmbeddedCatalog` before entering their package-local cores.
4. The task 01 diagnostic characterization golden is byte-unchanged.
5. Final changed-path inspection is limited to `internal/baseline/` and this
   task file; the task status line was the pre-existing Daemon-owned change.

The commands under `## Verification` were not run; the Daemon owns them.
