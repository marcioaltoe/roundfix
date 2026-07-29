---
task: task_01
spec: 0054-tooling-task-and-verification-hygiene
status: pending
type: test
complexity: high
---

# Task 01: Regenerate every derived digest from its canonical source

## Overview

Give each validator suite that pins a derived digest an update mode that
rewrites its artifacts from the canonical sources, so a Skill or Baseline
module edit no longer requires transcribing hashes by hand. Verifiable on its
own: editing a canonical source, running the update modes, and rerunning the
suites leaves them green with no hand-edited value.

## Requirements

1. MUST add an update mode to the authorial Skill-sync validator that
   rewrites every catalog setup snapshot's repository-sourced content digest
   and each snapshot's top-level digest from the canonical Skill sources.
2. MUST add an update mode to the catalog compatibility validator that
   rewrites the normalized catalog snapshot and its digest from the
   catalog's own normalization and digest accessors.
3. MUST add an update mode to the parity-corpus validator that rewrites the
   fixture's digests and the manifest's size and hash rows, leaving the
   frozen inventory digest, frozen date, and counted totals untouched.
4. MUST cover the second derived chain: an update mode that rewrites the
   formatter golden fixtures from the plan's own generated postimages and
   re-pins the profile's golden digest, and an update mode that rewrites the
   maintained source-baseline manifest.
5. MUST compute each source-baseline entry's span from its delimiting
   entry markers in the corpus rather than by offset arithmetic, then
   rewrite the identity and index digests, and MUST self-validate every
   entry digest before writing.
6. MUST compute Skill content digests through the production folder-hash
   helper, so the validator and the runtime agree on one algorithm.
7. MUST leave every suite's default behavior unchanged: without the update
   mode they compare and fail exactly as today.
8. MUST make each mismatch diagnostic name the sanctioned regeneration
   command so a stale pin reads as a stale snapshot rather than a broken
   catalog.

## Subtasks

- [ ] Add the update mode to the Skill-sync and catalog compatibility
      validators on the production folder-hash algorithm.
- [ ] Add the update mode to the parity-corpus validator, preserving its
      frozen fields.
- [ ] Add the update modes for the formatter goldens with the profile's
      pinned golden digest and for the maintained source baseline.
- [ ] Implement marker-based span recomputation with self-validation for
      source-baseline entries.
- [ ] Point every mismatch diagnostic at the regeneration command.

## Acceptance Criteria

- [ ] Running every update mode against an unchanged repository rewrites no
      bytes: the suites stay green and the working tree stays clean.
- [ ] After editing a canonical Skill source, running the update modes
      leaves the full gate green with no hand-edited digest value.
- [ ] After editing a Baseline module clause, running the update modes
      rewrites the source-baseline manifest spans and digests, the formatter
      goldens, the profile golden digest, and the catalog fixtures, and the
      suites pass.
- [ ] A stale pin without the update mode fails with a diagnostic naming the
      regeneration command.
- [ ] Every source-baseline entry digest validates against its corpus span
      after regeneration; a corrupted span fails instead of being written.
- [ ] The frozen parity-corpus inventory digest, frozen date, and counted
      totals are byte-identical after regeneration.

## Context

- interface: `skills/baseline_skill_contract_test.go`
- interface: `internal/baseline/catalog_test.go`
- interface: `internal/baseline/compatibility_corpus_test.go`
- interface: `internal/baseline/plan_test.go`
- interface: `internal/baseline/preservation_test.go`
- interface: `internal/skillhash/hash.go`
- interface: `skills/skills.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/baseline/ ./skills/` — expected: pass with the default comparing behavior unchanged.
- `go test -count=1 ./internal/baseline/ -run 'TestCatalogCompatibility|TestBaselineCompatibilityCorpus|TestFormatterComposition' && git diff --quiet -- internal/baseline/testdata internal/baseline/assets` — expected: default runs mutate nothing.

## References

`_prd.md` → User Stories 1–2, Core Features 1–2; `_techspec.md` → Build
Order 1, Interfaces: regeneration flags; ADR-0081.
