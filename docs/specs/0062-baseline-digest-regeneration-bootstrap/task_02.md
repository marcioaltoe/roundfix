---
task: task_02
spec: 0062-baseline-digest-regeneration-bootstrap
status: completed
type: backend
complexity: high
---

# Task 02: Give the catalog loader a regeneration mode

## Overview

Catalog validation is unconditionally strict, which is correct everywhere
except inside the run whose purpose is to rewrite the very pins it checks. This
Task adds an explicit regeneration mode to the loader and an enumerable
allowlist of the diagnostic codes that mode defers. Nothing consumes the mode
yet — it is verifiable on its own by proving that the mode defers exactly the
allowlisted codes and that the default path is unchanged.

## Requirements

1. MUST add an explicit regeneration mode to the catalog loader, defaulting to
   off, so every existing caller keeps today's strict behavior with no change
   at the call site.
2. MUST express the deferred set as an enumerable allowlist keyed by diagnostic
   code, containing exactly the formatter golden-digest mismatch code, so the
   blast radius can be read off the list rather than inferred from a severity.
3. MUST defer an allowlisted code only when regeneration mode is on; a code
   outside the allowlist MUST still be recorded in regeneration mode.
4. MUST expose the regeneration-mode load through an unexported entry point, so
   no package outside the Baseline package and no production, CLI, or CI caller
   can reach it.
5. MUST leave the public catalog-loading API's signature and behavior
   unchanged.
6. MUST NOT wire the update path to the new entry point in this Task; that is
   the next slice.

## Subtasks

- [ ] Add the regeneration flag to the loader, defaulting to off.
- [ ] Add the deferred-code allowlist with the single formatter digest code.
- [ ] Apply the deferral at the diagnostic sink, gated on the flag.
- [ ] Add the unexported regeneration-mode load entry point.

## Acceptance Criteria

- [ ] With regeneration mode off, a catalog whose formatter golden digest
      disagrees with its fixtures still records the mismatch diagnostic.
- [ ] With regeneration mode on, the same catalog records no mismatch
      diagnostic.
- [ ] With regeneration mode on, a diagnostic whose code is absent from the
      allowlist is still recorded.
- [ ] The regeneration entry point is unexported and unreachable from outside
      the package.
- [ ] The public catalog-loading API's signature is unchanged.
- [ ] The characterization corpus from task 01 is unchanged, proving the
      default path did not move.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/catalog_load.go`
- interface: `internal/baseline/catalog_validate.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestCatalogRegenerationMode -count=1` —
  expected: exit 0; deferral applies only in regeneration mode and only to the
  allowlisted code.
- `go test ./internal/baseline -run TestCatalogDiagnosticCharacterization -count=1`
  — expected: exit 0; the strict default path produces the same diagnostics as
  before this Task.
- `grep -q "goldenDigest.mismatch" internal/baseline/catalog_validate.go` —
  expected: exit 0; the code is still emitted, not deleted.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `go vet ./internal/baseline` — expected: exit 0.

## References

- `_prd.md` → Core Features 1; Goals (outside regeneration still fails closed).
- `_techspec.md` → Implementation Design: Interfaces; Build Order 2.
- ADR-0085.

## Result

Implemented the loader-only regeneration slice. `catalogLoader` now carries an
off-by-default regeneration flag, and its diagnostic sink defers only codes in
the explicit `deferredDuringRegeneration` allowlist. The allowlist contains only
`catalog.profile.formatter.goldenDigest.mismatch`. `LoadCatalog` keeps its
existing signature and strict path, while the unexported
`loadEmbeddedCatalogForRegeneration` entry point enables the mode. No update,
CLI, CI, or production caller was wired to the entry point in this Task.

Focused checks:

- Red signal: the focused strict-mode subtest with a repository-local
  `GOCACHE` failed to compile before implementation because
  `deferredDuringRegeneration`, `loadCatalog`, and
  `loadEmbeddedCatalogForRegeneration` did not exist.
- `rtk gofmt -w internal/baseline/catalog.go internal/baseline/catalog_load.go internal/baseline/catalog_test.go`
  — exit 0.
- Each of the three focused `TestCatalogRegenerationMode` subtests was run
  separately with `-count=1` and a repository-local `GOCACHE`; all exited 0:
  strict mismatch recording, regeneration mismatch deferral, and regeneration
  recording of a non-allowlisted diagnostic.
- Focused `TestEmbeddedCatalog` run with `-count=1` and a repository-local
  `GOCACHE` — exit 0.
- `rtk git diff --exit-code -- internal/baseline/testdata/catalog.diagnostics.golden.json`
  — exit 0; the characterization corpus is byte-unchanged in this Task.
- `rtk git diff --check` — exit 0.

Acceptance evidence:

1. `strict load records allowlisted mismatch` passes through public
   `LoadCatalog` and observes the formatter mismatch diagnostic.
2. `regeneration load defers allowlisted mismatch` passes with the same fixture
   drift and no validation error.
3. `regeneration load records non-allowlisted diagnostic` observes
   `catalog.profile.module.unknown` while regeneration mode is enabled.
4. The regeneration entry point begins with a lowercase identifier; symbol
   inspection found only its definition and the package-local compile-time
   signature assertion, with no caller wired.
5. The compile-time assertion keeps `LoadCatalog` at
   `func(fs.FS) (*Catalog, error)`, and the strict focused cases pass through
   that public function.
6. The characterization golden file has no diff. The Daemon retains ownership
   of the declared characterization Verification command.
7. The final changed-path inspection is limited to `internal/baseline/` and
   this task file; the task status line was the pre-existing Daemon-owned
   change.

The commands under `## Verification` were not run; the Daemon owns them.
