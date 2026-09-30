---
task: task_01
spec: 0183-run-storage-that-says-what-it-holds
status: completed
type: backend
complexity: medium
---

# Task 01: A retention prune reports only what it reclaimed

## Overview

Journal Retention never deletes a `runs` row, so every terminal Run past the cutoff stays a prune candidate forever. `gc`, `gc --dry-run` and the operational sweep of `implement`, `resolve` and `watch` count and report emptied Runs again on every call, and the sweep takes the machine-wide write lock to delete nothing. This Task makes a Run count only while it still holds Run Event Journal rows or an artifact directory, and makes the prune take the write lock only when some candidate still has events. Candidates come from the machine-wide Run Database. The report goes to the operator's stdout (`gc`) or stderr (the sweep). Nothing but the eligible Runs' events and artifact directories is ever deleted.

## Requirements

1. MUST change `PruneTerminalRuns` in `internal/store/journal.go` to fill each eligible candidate's event count on the read connection before any write, return without opening a write transaction when no candidate has events, and otherwise, inside one `withWriteTx`, count the candidates' events again and delete only the events of the candidates that have them at that moment, so the count and the deletion come from the same transaction.
2. MUST make `PruneResult.RunIDs` name only the Runs whose events were deleted, keep `Events` as the rows the DELETE removed, and add `EligibleRunIDs` naming every Run past the cutoff at prune time; `TerminalRunPruneCandidates` and the eligibility rule (terminal, completed before the cutoff) MUST NOT change, and no `runs` row or Active Run lock may ever be deleted.
3. MUST add `retentionReclaimable(candidates []store.PruneCandidate, hasArtifactDir func(runID string) (bool, error)) ([]string, error)` in `internal/cli/gc.go`, returning in candidate order every candidate with events or an existing artifact directory, and MUST use it for the `gc --dry-run` report so `Runs eligible` and `Eligible Runs` list only reclaimable Runs.
4. MUST make `pruneRunRetention` filter the candidates' artifact directories by `EligibleRunIDs`, remove them, and report the union of `RunIDs` and the removed directories' Run IDs in candidate order, so a Run whose events are gone but whose directory survived is still removed and reported.
5. MUST keep `sweepRunRetention`'s silence rule, so the sweep prints no `pruned Run storage` line when nothing was reclaimed.
6. MUST update the `gc` section of `docs/user-guide/commands.md` and the Run storage retention section of `.agents/skills/roundfix/SKILL.md` to contain the phrase `count only Runs that still hold Run Event Journal rows or an artifact directory`, and the skill's sweep paragraph to contain `only when it reclaimed something`; then MUST run `make skills-sync` so `skills/roundfix/SKILL.md` matches.
7. MUST add the sentence `A Run it already emptied is not counted again.` to the **GC Command** entry of `CONTEXT.md`.
8. MUST put the new tests in `internal/store/prune_reclaimable_test.go` and `internal/cli/gc_reclaimable_test.go`, using temporary Roundfix Homes only, and MUST keep `TestPruneTerminalRunsDeletesOnlyEligibleJournalRows` and every existing `gc` and sweep test green without renaming any.

## Subtasks

- [ ] Count events before writing and lock only to delete.
- [ ] Report deleted, eligible and reclaimed Runs separately.
- [ ] Share one reclaimable predicate between the dry run and the prune report.
- [ ] Update the guide, the skill and its mirror, and the glossary entry.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A second `gc` after a complete one reports `Runs pruned: 0`, and a dry run then reports `Runs eligible: 0`.
- [ ] A Run with no events left but a surviving artifact directory is removed and reported by `gc`.
- [ ] A store holding only emptied candidates is pruned with no write transaction: the prune succeeds while another holder keeps the machine-wide write lock.
- [ ] A candidate with events is still pruned, and `EligibleRunIDs` names every candidate.
- [ ] The operational sweep prints no `pruned Run storage` line over emptied Runs.
- [ ] The guide, the skill and its mirror, and the glossary state the rule, and `make skills-sync-check` passes.

## Context

- interface: `internal/store/journal.go`
- interface: `internal/store/journal_test.go`
- creates: `internal/store/prune_reclaimable_test.go`
- interface: `internal/cli/gc.go`
- interface: `internal/cli/gc_test.go`
- creates: `internal/cli/gc_reclaimable_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`
- instruction: `docs/adr/0171-a-retention-prune-reports-only-what-it-reclaimed.md`
- instruction: `docs/adr/0033-the-run-event-journal-is-pruned-by-retention.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestPruneTerminalRunsSkipsCandidatesWithoutEvents|TestPruneTerminalRunsNeedsNoWriteLockWhenNoCandidateHasEvents)$" ./internal/store 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestPruneTerminalRunsSkipsCandidatesWithoutEvents TestPruneTerminalRunsNeedsNoWriteLockWhenNoCandidateHasEvents; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task neither test exists, so the command fails.
- `out="$(go test -count=1 -v -run "^(TestGCSecondRunReportsNothingPruned|TestGCReclaimsARunWhoseArtifactDirectoryOutlivedItsEvents|TestRetentionSweepIsSilentWhenNothingIsReclaimed)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestGCSecondRunReportsNothingPruned TestGCReclaimsARunWhoseArtifactDirectoryOutlivedItsEvents TestRetentionSweepIsSilentWhenNothingIsReclaimed; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the tests exists.
- `for file in docs/user-guide/commands.md .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "count only Runs that still hold Run Event Journal rows or an artifact directory" || { printf 'missing phrase in %s: %s\n' "$file" "count only Runs that still hold Run Event Journal rows or an artifact directory" >&2; exit 1; }; done; for file in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "only when it reclaimed something" || { printf 'missing phrase in %s: %s\n' "$file" "only when it reclaimed something" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task the phrases are absent.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "A Run it already emptied is not counted again." || { printf 'missing phrase in %s: %s\n' CONTEXT.md "A Run it already emptied is not counted again." >&2; exit 1; }` — expected: exit 0; before this Task the sentence is absent.

## References

- [_prd.md](_prd.md) — Goals 1–2; Core Feature 1; Success Metrics 1–2
- [_techspec.md](_techspec.md) — A prune reports only what it reclaimed; Interfaces; API Contract 1; Testing Approach 1; Build Order 1
- [references/2026-09-29-gc-reports-already-emptied-runs-as-pruned.md](references/2026-09-29-gc-reports-already-emptied-runs-as-pruned.md)
- ADR-0171; ADR-0033; ADR-0090

## Result

### Implementation

- `PruneTerminalRuns` now counts candidate events before requesting the writer,
  returns every cutoff candidate in `EligibleRunIDs`, and recounts plus deletes
  only event-bearing candidates inside one write transaction. `RunIDs` and
  `Events` therefore describe only journal rows that transaction removed.
- GC dry runs and live prune reports use `retentionReclaimable` to include only
  Runs that still hold events or an artifact directory. Live cleanup filters
  directories by `EligibleRunIDs` and reports the candidate-ordered union of
  journal and directory reclamation.
- Added real temporary-store regression coverage for empty candidates, held
  write locks, repeated GC, directory-only reclamation, reclaimable ordering
  and callback errors, and silent operational sweeps.
- Updated the user guide, canonical Roundfix skill, generated skill mirror, and
  GC Command glossary entry. Ran the sanctioned skill sync and confirmed the
  Baseline digests required no derived update.

### Focused checks

- `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test ./internal/store` —
  passed, including the existing eligibility/lifecycle tests and the new
  reclaimable-prune tests.
- `rtk env GOCACHE=/tmp/roundfix-task01-gocache go test -run 'Test(RunGC|GCSecondRunReportsNothingPruned|GCReclaimsARunWhoseArtifactDirectoryOutlivedItsEvents|RetentionReclaimable|RetentionSweepIsSilentWhenNothingIsReclaimed|RunImplementPreflightRetention)' ./internal/cli`
  — passed, including existing GC and sweep coverage.
- `rtk make skills-sync` — passed; the shipped Roundfix skill mirror was
  regenerated from the canonical skill.
- `rtk env GOCACHE=/tmp/roundfix-task01-gocache make baseline-digests` —
  passed with `changed: false`.
- Canonical/mirror byte comparison and exact-phrase inspections for the guide,
  both skill paths, and `CONTEXT.md` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task01-gocache make verify-incremental` —
  passed on the host-permission rerun: `go vet`, the full package suite, the
  skill contract tests, `roundfix skills check`, and the build all exited 0.

### Acceptance evidence

1. `TestGCSecondRunReportsNothingPruned` observes `Runs pruned: 0` on the
   second live GC and `Runs eligible: 0` on the following dry run.
2. `TestGCReclaimsARunWhoseArtifactDirectoryOutlivedItsEvents` observes the
   directory-only Run in stdout and confirms its artifact directory is gone.
3. `TestPruneTerminalRunsNeedsNoWriteLockWhenNoCandidateHasEvents` holds the
   independent machine-wide write lock while an emptied-candidate prune
   succeeds and reports no reclamation.
4. `TestPruneTerminalRunsSkipsCandidatesWithoutEvents` confirms the
   event-bearing candidate is deleted and reported while `EligibleRunIDs`
   retains both cutoff candidates.
5. `TestRetentionSweepIsSilentWhenNothingIsReclaimed` observes empty stderr
   over an emptied terminal Run.
6. Exact wording is present in the guide, canonical skill, mirror, and glossary;
   the skill files compare byte-for-byte after `make skills-sync`. The authored
   `make skills-sync-check` remains part of Daemon-owned Verification.

### Non-gating environment observation

- An earlier broad `go test ./internal/store ./internal/cli` run passed the
  store package but the CLI package encountered the sandbox's process-table
  permission refusal in two unrelated force-stop tests. Its suite guard also
  observed documentation edits made while that long-running check was active.
  The stable retention/GC-focused CLI check above was rerun after edits and
  passed. The first sandboxed incremental gate saw the same process-table
  refusal; its host-permission rerun passed. The Task's authored Verification
  commands were not run.
