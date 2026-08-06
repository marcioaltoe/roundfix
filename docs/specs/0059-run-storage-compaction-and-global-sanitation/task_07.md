---
task: task_07
spec: 0059-run-storage-compaction-and-global-sanitation
status: completed
type: backend
complexity: medium
---

# Task 07: Ship the compaction command and document all three surfaces

## Overview

Corrective Task from the QA gate's F-001 (`Blocks-Completion`) and F-002.

`internal/store` carries `PreviewCompaction` and `Compact` with every guard the
Spec asked for — Active Run, other writer, insufficient temporary capacity —
and none of it is reachable. `roundfix gc compact` exits `2` with
`unexpected argument "compact"`. The TechSpec declared the surface in its API
Contracts table and no Task requirement ever asked for it, so the feature
shipped as an unreachable library.

task_05 then refused to document it, correctly: writing the Skill contract for
a command that does not exist would have been a false contract. Its Result
names this exact follow-up. But task_05 still settled `completed` having
changed nothing, because all four of its Verification commands pass most easily
when no work happened — the shape `SC-VERIFY-WORK-INDEPENDENT` now refuses at
authoring time.

This Task ships the route and then the documentation, in that order, and its
Verification is written so that doing nothing fails it.

## Requirements

1. MUST add `roundfix gc compact` with a preview-first contract: the bare form
   previews bytes before, reclaimable, and projected after; `--apply` performs
   the guarded compaction and reports the same three numbers as measured.
2. MUST route every guard `internal/store` already implements: an Active Run,
   another writer, and insufficient temporary capacity each refuse by name,
   before any mutation.
3. MUST list the subcommand in `roundfix gc --help`.
4. MUST leave per-repository GC and `gc sanitize` behaving exactly as they do
   today.
5. MUST then document all three operator surfaces in the canonical Roundfix
   Skill and its mirror: `gc compact`, `gc sanitize`, and `storage report`,
   including compaction's three refusals and that it is explicit rather than an
   automatic side effect of a retention sweep.
6. MUST regenerate the mirror with `make skills-sync` and run the ADR-0081
   chain, including the two characterization corpora that
   `make baseline-digests` does not reach.
7. MUST change only `internal/cli/**`, `.agents/skills/roundfix/**`,
   `skills/roundfix/**`, this Task file, and the ADR-0081 digest fallout.

## Subtasks

- [ ] Add the `gc compact` route over the existing store API.
- [ ] Document all three surfaces in the Skill and regenerate the mirror.

## Acceptance Criteria

- [ ] `roundfix gc compact` previews three numbers and exits 0.
- [ ] `roundfix gc compact --apply` compacts and reports measured numbers.
- [ ] Each of the three refusals is reachable from the command and names its
      cause.
- [ ] `roundfix gc --help` lists the subcommand.
- [ ] Per-repository GC and `gc sanitize` are unchanged.
- [ ] The Skill documents `gc compact`, `gc sanitize`, and `storage report`.
- [ ] The mirror is byte-identical to the canonical Skill.

## Context

- interface: `internal/cli/gc.go`
- interface: `internal/store/journal.go`
- instruction: `.agents/skills/roundfix/SKILL.md`

## Verification

Every command below fails if this Task does nothing. That is deliberate:
task_05 settled `completed` having changed nothing because all four of its
checks passed most easily on an untouched repository.

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix gc --help | grep -q "gc compact"`
  — expected: exit 0; the shipped help lists the subcommand.
- `go run -buildvcs=false ./cmd/roundfix gc compact > /dev/null` — expected:
  exit 0; the preview runs rather than rejecting the argument.
- `grep -q "gc compact" .agents/skills/roundfix/SKILL.md && grep -q "gc sanitize" .agents/skills/roundfix/SKILL.md && grep -q "storage report" .agents/skills/roundfix/SKILL.md`
  — expected: exit 0; the Skill names all three surfaces.
- `output="$(go test ./internal/cli -count=1 -run 'GCCompact|Compact' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the command tests ran and passed.
- `make skills-sync-check` — expected: exit 0; the mirror matches.
- `make verify` — expected: exit 0.

## References

- `qa/qa-report-2026-08-06.md` → F-001, F-002.
- `_prd.md` → Core Feature 1; User Stories 1 and 2.
- `_techspec.md` → API Contracts.
- `docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`.
- ADR-0081.

## Result

Implemented the missing `gc compact` command route over the existing store
API. The bare command opens the existing Run Database read-only, prints bytes
before, reclaimable, and projected after, and leaves a missing database
missing while reporting zero measurements. When a concurrent writer advances
the database during the exact snapshot measurement, the bare command falls
back to the existing immutable storage-report measurements so preview remains
read-only and available. `--apply` opens the writer, previews first, invokes
guarded compaction with that preview, and prints the three measured result
values. It never uses the reporting fallback, so runtime failures remain on
stderr and preserve the store's named Active Run, competing-writer, and
temporary-capacity causes.

Updated `gc --help` and the canonical Roundfix Skill after the route existed.
The Skill now documents `gc compact`, `gc sanitize`, and `storage report`, the
three compaction refusals, the preview-before-apply flow, and the rule that
retention sweeps never compact automatically. `make skills-sync` regenerated
the mirror, and the ADR-0081 generator chain regenerated its derived digest
and characterization fallout.

The pre-change focused signal was:

- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/cli -run '^TestRunGCCompact' -count=1`
  — failed to compile because `gcDependencies` had no compaction route,
  proving the new command tests were red before the production edit.

Focused checks run before Verification feedback attempt 1:

- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/cli -run '^TestRunGCCompact' -count=1 -v`
  — passed the real SQLite preview/apply journey, real Active Run and
  competing-writer refusals, isolated capacity refusal, and unsupported-input
  case.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/cli -count=1`
  — exited `0` for the complete CLI package, including the existing
  per-repository GC, sanitation, storage-report, and help suites.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go vet ./internal/cli`
  — exited `0`.
- `rtk make skills-sync` — exited `0` and regenerated the Roundfix Skill
  mirror.
- `rtk make baseline-digests` — exited `0` and reported
  `baseline-digests: regenerated` for the catalog, parity, setup, and plan
  characterization digest fallout.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/baseline -count=1 -run '^TestBaselinePlanCharacterization$' -update-baseline-plan-characterization`
  — exited `0`.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/baseline -count=1 -run '^TestCatalogDiagnosticCharacterization$' -update-catalog-diagnostics`
  — exited `0`.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test ./internal/baseline -run '^Test(BaselinePlanCharacterization|CatalogDiagnosticCharacterization)$' -count=1`
  — exited `0` against the regenerated corpora.
- `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — exited `0`; the canonical Skill and mirror are byte-identical.
- `rtk git -c core.fsmonitor=false diff --check` — exited `0`.

Acceptance evidence:

- Preview: `TestRunGCCompactPreviewsAndAppliesMeasuredBytes` invokes the public
  CLI runner against a real temporary Run Database, reconciles `after = before
  - reclaimable`, and proves the preview leaves database bytes unchanged.
- Apply: the same test invokes `gc compact --apply`, reconciles `after = before
  - reclaimed`, and compares the reported after value with the resulting
  database file size.
- Refusals: `TestRunGCCompactRefusalsNameCause` reaches the real store guard for
  an Active Run and another writer, routes the capacity guard through the
  command's dependency boundary, asserts each named cause on stderr, and
  compares database bytes before and after every refusal.
- Help: `TestRunGCHelp` asserts that `roundfix gc compact [--apply]` appears in
  `gc --help`; it passes in the complete CLI package check.
- Existing behavior: the complete CLI package check retains passing
  `TestRunGCDryRunListsEligibleRunsAndChangesNothing`,
  `TestRunGCPrunesEligibleJournalsArtifactsAndOrphans`,
  `TestRunGCSkipsWhenJournalRetentionIsZero`, and
  `TestRunGCSanitizeClassifiesEveryRecordedRootAndMutatesOnlyProvenDirectories`.
- Skill coverage: the reviewed canonical diff names `gc compact`, `gc
  sanitize`, and `storage report`, including all requested guard and
  explicit-only compaction language; the sanctioned generator chain exited
  `0` after that edit.
- Mirror equality: the byte comparison above exited `0` after the final Skill
  edit and sync.

Follow-up: Task 07 and its authorization note say the two characterization
corpora sit outside `BASELINE_DIGEST_STEPS`, but the current Makefile includes
both. Consequently `make baseline-digests` itself regenerated four
`internal/baseline/testdata/plan-characterization/*.golden.json` files. Their
diffs contain only catalog/manifest digest and derived identity changes from
the authorized Skill edit, so they remain ADR-0081 deterministic fallout; no
golden value was hand-edited.

### Verification feedback attempt 1

The retained diagnostic identified a competing-writer fingerprint change
during the bare preview's exact snapshot measurement. That failure happened
after the read-only snapshot was measured but before its private fingerprint
could be accepted for a later apply. The repair keeps the exact store preview
as the first path, uses immutable `storage report` measurements only for a bare
preview invalidated by that typed writer condition, and leaves `--apply` on
the original refusal path.

Fresh feedback-repair evidence:

- Before the repair,
  `rtk env GOCACHE=/private/tmp/roundfix-task07-feedback-gocache go test ./internal/cli -run '^TestRunGCCompactPreviewReportsStorageMeasurementWhenWriterAdvances$' -count=1 -v`
  reproduced exit `1` for the bare preview under an injected typed writer
  invalidation.
- After the repair, the same focused regression test exited `0`; it reconciles
  all three fallback measurements, proves database bytes are unchanged, and
  separately proves `--apply` still refuses and names the competing writer.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-feedback-gocache go test ./internal/cli -run '^(TestRunGCHelp|TestRunGCCompactPreviewsAndAppliesMeasuredBytes|TestRunGCCompactPreviewReportsStorageMeasurementWhenWriterAdvances|TestRunGCCompactRefusalsNameCause)$' -count=1 -v`
  exited `0` for stable preview/apply, live-writer preview, all three guarded
  apply refusals, and help discovery.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-feedback-gocache go vet ./internal/cli`
  exited `0`.
- `rtk make skills-sync` and `rtk make baseline-digests` both exited `0` after
  the canonical Skill documented the fallback-versus-refusal distinction; the
  latter ran its catalog-diagnostic and plan-characterization update steps.
- `rtk cmp -s .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` and
  `rtk git -c core.fsmonitor=false diff --check` both exited `0` after the
  feedback repair.

The acceptance evidence above remains applicable, with the new live-writer
regression adding direct evidence that the required bare preview exits `0`
while the competing-writer refusal remains reachable from `--apply` before
mutation.

The commands under `## Verification` were not run; the Daemon owns those
commands, Task status, settlement, and the Task commit.
