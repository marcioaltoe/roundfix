---
task: task_01
spec: 0054-tooling-task-and-verification-hygiene
status: completed
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
- `before="$(git -c core.fsmonitor=false diff --no-ext-diff --binary -- internal/baseline/testdata internal/baseline/assets | git hash-object --stdin)" && go test -count=1 ./internal/baseline/ -run 'TestCatalogCompatibility|TestBaselineCompatibilityCorpus|TestFormatterComposition' && test "$before" = "$(git -c core.fsmonitor=false diff --no-ext-diff --binary -- internal/baseline/testdata internal/baseline/assets | git hash-object --stdin)"` — expected: default runs preserve the pre-existing artifact diff byte-for-byte.

## References

`_prd.md` → User Stories 1–2, Core Features 1–2; `_techspec.md` → Build
Order 1, Interfaces: regeneration flags; ADR-0081.

## Result

### Implementation

- Added package-owned `-update` modes for the authorial Skill-sync, catalog
  compatibility, parity-corpus, formatter-composition, and maintained
  source-baseline validators. Default invocations remain compare-only.
- The Skill updater rewrites every setup snapshot's Roundfix-owned
  `contentDigest` through `skills.SkillFolderHash`, then recomputes the setup
  digest from the catalog's canonical JSON payload. Moving from the former
  private lexicographic hash exposed and regenerated the existing pins in all
  three setup snapshots.
- The catalog updater writes `Catalog.Normalized()` and `Catalog.Digest()`.
  Formatter and source-baseline updates also refresh these compatibility pins
  from the filesystem catalog after changing assets, so their derived chain
  does not leave a stale catalog snapshot.
- The parity updater recomputes repository-owned Skill digests, fixture setup
  digests, the managed-entry ledger, and the fixture artifact's byte/hash row.
  It asserts that the frozen date, source-suite inventory digest and test
  count, profiles, adoption states, fixture inventory, and retired-behavior
  fields remain byte-identical before writing the manifest.
- The formatter updater first requires every declared golden postimage, then
  writes those generated bytes, recomputes the portable-file digest, and
  re-pins the Profile.
- The source-baseline updater derives every span from its unique
  `source-baseline-entry` markers, hashes the exact recomputed span,
  self-validates every marker/span/digest tuple before any manifest write,
  and recomputes the manifest, identity, corpus, and index digests.
- Derived mismatch diagnostics name `make baseline-digests`.

### Focused checks

- Verification feedback attempt 1 showed the owning Go tests passed; the
  chained `git diff --quiet` failed because it compared this Task's intentional
  tracked artifact updates against `HEAD`, rather than comparing repository
  state before and after the default validators. The declared non-mutation
  check now hashes the binary artifact diff on both sides of the test run and
  disables the worktree fsmonitor for those read-only comparisons.
- The three owning tests were rerun separately (not through the declared
  Verification command) with anchored selectors:
  `TestCatalogCompatibility`, `TestBaselineCompatibilityCorpus`, and
  `TestFormatterComposition`; each passed. The binary artifact diff's Git
  object hash was `dfe3d2226a4b90c78fb2e8b96c01374ab3715338` both before
  and after those focused runs.
- Red signal: the two owning packages initially rejected `-update` with
  `flag provided but not defined: -update`.
- `rtk env GOCACHE=/private/tmp/roundfix-0054-task01-gocache go test ./skills -run '^TestAuthorialSkillSync$' -update -count=1`
  passed.
- Each Baseline update selector passed independently with the same task-local
  cache: `TestCatalogCompatibility`, `TestBaselineCompatibilityCorpus`,
  `TestFormatterComposition`, and
  `TestReadoptionCompatibilityMaintainedFixture`.
- Re-running all five update modes in the TechSpec sequence changed no
  artifact bytes: the binary diff hash for
  `internal/baseline/assets` plus `internal/baseline/testdata` was
  `016f21d8882148ff3c600caf03a9e79884e5a7ac0a88e088392ccca81ea2c5e1`
  both before and after.
- `rtk env GOCACHE=/private/tmp/roundfix-0054-task01-gocache go test ./skills -run '^TestAuthorialSkillSync(UpdateModeRoundTrip)?$' -count=1`
  passed, including a disposable canonical Skill edit and automatic
  content/setup re-pinning.
- `rtk env GOCACHE=/private/tmp/roundfix-0054-task01-gocache go test ./internal/baseline -run '^(TestCatalogCompatibility|TestBaselineCompatibilityCorpus|TestFormatterComposition|TestReadoptionCompatibilityMaintainedFixture|TestSourceBaselineRegenerationRejectsCorruptedSpan)$' -count=1`
  passed.
- A reversible Baseline-clause probe changed one module clause and its corpus
  entry, ran source-baseline then formatter regeneration, and left the
  catalog, formatter, maintained-source, and corrupted-span validators green.
  Reversing the two probe edits and regenerating restored the exact artifact
  diff hash above.
- A deliberately stale `catalog.digest` failed in default mode with
  `run 'make baseline-digests'`; the catalog update mode restored it.
- `rtk git diff --check` passed.

### Acceptance evidence

- **Idempotent unchanged regeneration:** evidenced by all five passing update
  modes and the identical before/after artifact diff hash.
- **Canonical Skill edit:** the disposable round-trip test proved production
  folder hashing plus automatic setup re-pinning; the focused Skill
  validators passed afterward. The Daemon-owned full gate remains unrun.
- **Baseline module clause edit:** the reversible probe regenerated
  marker-shifted source-baseline spans/digests, formatter goldens, the Profile
  golden digest, and catalog compatibility pins; focused default validators
  passed and the probe was fully reversed.
- **Stale diagnostic:** the deliberate catalog pin corruption failed with the
  sanctioned command in its diagnostic, then regenerated successfully.
- **Source-baseline self-validation:** the corrupted-span regression passed by
  observing rejection, while the reversible probe exercised marker-based
  offset shifts before successful writes.
- **Frozen parity fields:** the update helper byte-compares every frozen
  manifest field before writing; parity update and default validation passed,
  and only the asset-sync artifact row changed in the manifest diff.

### Follow-up

- The later Makefile task must run maintained source-baseline regeneration
  before formatter regeneration when corpus bytes changed; otherwise the
  formatter test binary cannot load the stale embedded Source Baseline.
  This ordering belongs to the Make target task, not this diff.
- The Task's declared `## Verification` commands were not run; the Daemon
  owns them and terminal settlement.
