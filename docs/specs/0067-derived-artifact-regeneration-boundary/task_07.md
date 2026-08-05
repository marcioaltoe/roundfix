---
task: task_07
spec: 0067-derived-artifact-regeneration-boundary
status: completed
type: backend
complexity: medium
---

# Task 07: Let a record carry per-path exceptions

## Overview

The gate's remaining F-002 is not a wrong premise on either side. It is a model
that is too coarse.

`internal/baseline/testdata/parity-corpus/` holds fifteen files. The sanctioned
command rewrites exactly two of them — `v1/manifest.json` and
`v1/fixtures/asset-sync.json`. The other thirteen, including `blobs.json` and
`matrix.json`, it never touches. A single directory-level `owner` is therefore
false about whichever part it does not describe: `frozen` lies about two files,
`sanctioned` lies about thirteen.

The maintainer settled this on 2026-08-05: the record gains per-path
exceptions. The PRD's frozen premise stands for the corpus, and the two derived
files are declared as what they measurably are.

## Requirements

1. MUST let an ownership record declare a directory-level `owner` plus an
   `exceptions` list, each entry naming a path and its own `owner`.
2. MUST restore the parity record to `owner: frozen` with `v1/manifest.json`
   and `v1/fixtures/asset-sync.json` declared `sanctioned`, so the thirteen
   frozen files are described truthfully again.
3. MUST resolve ownership for any path by the most specific matching record —
   an exception wins over its directory.
4. MUST extend the exhaustiveness test to exceptions: a path claimed by two
   exceptions, or by an exception outside its record's directory, fails.
5. MUST make the declared-step and frozen assertions evaluate per resolved path
   rather than per directory, so a frozen path rewritten by any step fails and
   a declared exception rewritten by its owner passes.
6. MUST NOT change which artifacts any command rewrites. This Task changes what
   the records can express, never what regeneration does.
7. MUST NOT change any digest value or artifact content.

## Subtasks

- [ ] Add `exceptions` to the record shape and its validation.
- [ ] Restore the parity record to frozen with its two declared exceptions.
- [ ] Resolve ownership by most-specific match.
- [ ] Extend exhaustiveness and the frozen assertion to per-path resolution.

## Acceptance Criteria

- [ ] The parity record reads `owner: frozen` with exactly two `sanctioned`
      exceptions, and every other file beneath it resolves `frozen`.
- [ ] A frozen-resolved path rewritten by any command fails the test, proven by
      a fixture that rewrites one.
- [ ] An exception path rewritten by its declared owner passes.
- [ ] A path claimed by two exceptions fails validation.
- [ ] An exception naming a path outside its record's directory fails
      validation.
- [ ] `make skills-sync && make baseline-digests && make verify` exits 0 for all
      three, and a second sanctioned run is a no-op with the gate still green.
- [ ] No digest value or artifact content changed by this Task.

## Context

- interface: `internal/baseline/testdata/parity-corpus/_ownership.yml`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -count=1 -run 'Ownership|Exception|Frozen|Exhaustive' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the per-path resolution tests ran and passed.
- `make skills-sync && make baseline-digests && make verify` — expected: exit 0;
  one sanctioned run leaves the gate green.
- `make baseline-digests && make verify` — expected: exit 0; the second run is a
  no-op.
- `grep -q "exceptions" internal/baseline/testdata/parity-corpus/_ownership.yml`
  — expected: exit 0; the parity record declares its exceptions.
- `grep -q "owner: frozen" internal/baseline/testdata/parity-corpus/_ownership.yml`
  — expected: exit 0; the corpus is frozen again at the directory level.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 1, 3 and 5; Decisions ("frozen stays frozen").
- `_techspec.md` → Interfaces; Risks & Considerations ("a manifest can lie").
- `qa/qa-report-2026-08-05-02.md` → F-002.

## Result

### Implementation

- Ownership records now accept strict `exceptions` entries with a relative
  path and owner. Validation rejects invalid owners, missing paths, escaping
  paths, duplicate claims across records, unused exceptions, and dedicated
  exceptions without a command.
- Ownership resolution applies an exact exception before a directory record,
  including when a nearer directory record exists. Remediation, measured
  sanctioned ownership, declared-step probes, and frozen probes use the
  resolved path owner.
- The parity record is frozen with exactly two sanctioned exceptions:
  `v1/manifest.json` and `v1/fixtures/asset-sync.json`. No regeneration command
  changed.

### Acceptance evidence

1. `TestDerivedOwnershipDeclaresKnownBoundaries` reads the real parity record,
   requires exactly the two sanctioned exceptions, and checks every other
   regular path under the corpus resolves frozen. It passed in the focused
   ownership command below.
2. `TestDeclaredStepRegenerationAndFrozenBoundaries/frozen_resolved_path_rejects_rewrite`
   resolves a fixture frozen, runs the command that rewrites it, and observes
   the required `rewrote frozen artifact` failure. The parent focused test
   passed.
3. The same declared-step test perturbs both parity exception paths and proves
   `make baseline-digests` restores them as sanctioned declarations. The
   parent focused test passed.
4. `TestDerivedOwnershipRejectsDuplicateExceptions` declares the same path in
   two records and observes the required multiple-exception rejection. It
   passed in the focused ownership command.
5. `TestDerivedOwnershipRejectsExceptionOutsideRecordDirectory` declares
   `../outside.json` and observes the required outside-directory rejection. It
   passed in the focused ownership command.
6. Daemon-owned Verification remains pending. This Agent did not run the
   declared `make skills-sync`, `make baseline-digests`, or `make verify`
   commands.
7. `rtk git diff --exit-code -- internal/baseline/testdata/parity-corpus/v1`
   exited 0. `rtk git diff --name-only -- internal/baseline/testdata` listed
   only `internal/baseline/testdata/parity-corpus/_ownership.yml`; no digest or
   artifact content changed.

### Focused checks

- Red signal: the exception regression command initially failed because the
  strict YAML decoder reported `field exceptions not found`.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1 -run '^TestDerivedOwnership'`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache GOFLAGS=-buildvcs=false go test ./internal/baseline -count=1 -run '^Test(DeclaredStepRegenerationAndFrozenBoundaries|MeasuredSanctionedOwnershipMatchesRecords)$'`
  — passed in 62.091s against isolated repository fixtures.
- `rtk git diff --check` — passed.

### Follow-ups

None discovered within this Task's slice.
