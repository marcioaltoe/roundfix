---
task: task_03
spec: 0183-run-storage-that-says-what-it-holds
status: pending
type: backend
complexity: medium
---

# Task 03: Doctor reports reclaimable Run storage

## Overview

Nothing tells the operator that Run storage is reclaimable until someone runs `roundfix gc`. This Task adds a read-only `storage` check to the Doctor Command. It opens the machine-wide Run Database through the read-only reader, reads the header's free-page count and page size, and counts the Runs past the Journal Retention cutoff that `gc` could still reclaim, using the predicate Task 01 added and one `lstat` per artifact directory. It prints one `storage:` line to stdout pointing to `roundfix gc` or `roundfix gc compact`. The check never writes, never creates a missing database, never runs `dbstat` or walks a directory, and never reports `failed`, so Doctor's exit code is unchanged.

## Requirements

1. MUST add `HealthCheckStorage = "storage"` in `internal/cli/health.go`, a `storage func(context.Context, roundconfig.Loaded) []CheckResult` field in `doctorDependencies` defaulting to `defaultDoctorStorageResults` in the new `internal/cli/doctor_storage.go`, and append its results in `runDoctorCommand` after the residue results and before `codex`.
2. MUST add `(*store.Store).FreeBytes(ctx) (int64, error)` in the new `internal/store/free_bytes.go`, reading `PRAGMA page_size` and `PRAGMA freelist_count` without a transaction or a pragma write.
3. MUST report `ok (no Run Database)` without creating anything when the Run Database path does not exist, and `partial (could not read Run Database: <err>)` when `store.OpenReader` fails, including on a schema-version mismatch, or when closing the reader fails.
4. MUST, with a non-zero Journal Retention, count the reclaimable Runs with `retentionReclaimable` over `TerminalRunPruneCandidates` at `now - retention`, where `hasArtifactDir` only `lstat`s `gcRunArtifactPath` under the Artifact Root that `ResolveArtifactDirectory` gives the loaded config and answers `false` when no root resolves; the detail then says `artifact directories not inspected outside a Git repository`. With a zero retention, the Runs clause MUST read `Journal Retention 0 keeps every Run`.
5. MUST report `found` with detail `Runs reclaimable: <n>; Run Database free bytes: <b>; next: <commands>` when the count is above zero or free bytes reach `doctorStorageFreeBytesThreshold` (64 MiB), naming `roundfix gc`, `roundfix gc compact` or both for whichever condition held, and otherwise `ok (nothing to reclaim; Runs reclaimable: 0; Run Database free bytes: <b>)`. It MUST NOT return `failed`.
6. MUST update the `doctor` section of `docs/user-guide/commands.md` with a `storage:` entry containing the phrase `never reports failed`, and state that `residue:` and `storage:` also report `found` or `partial`. MUST update the Doctor paragraph of `.agents/skills/roundfix/SKILL.md` to contain the phrase `storage check`, then run `make skills-sync` so `skills/roundfix/SKILL.md` matches. The Doctor help text MUST NOT change.
7. MUST add the phrase `whether Run storage is reclaimable` to the **Doctor Command** entry of `CONTEXT.md`.
8. MUST put the tests in `internal/cli/doctor_storage_test.go`, through `runDoctorCommand` with the other checks faked and a real Run Database in a temporary Roundfix Home, plus the pure result function for the free-bytes threshold, and MUST keep every existing Doctor test green without renaming any.

## Subtasks

- [ ] Read free bytes from the header and count reclaimable Runs with the shared predicate.
- [ ] Map the facts to one `storage:` line that never fails Doctor.
- [ ] Prove Doctor writes nothing and never creates a database.
- [ ] Update the guide, the skill and its mirror, and the glossary entry.
- [ ] Add a test for each acceptance criterion, each negative case separate.

## Acceptance Criteria

- [ ] A home with a Run past retention that still has events prints `storage: found` naming `roundfix gc`, and Doctor exits `0` when every other check passes.
- [ ] Free bytes at the threshold produce `found` naming `roundfix gc compact`.
- [ ] A home with nothing left to reclaim prints `storage: ok (nothing to reclaim; …)`.
- [ ] A missing Run Database prints `storage: ok (no Run Database)` and no database file exists afterwards.
- [ ] An unreadable Run Database prints `storage: partial`, and Doctor's exit code is unaffected.
- [ ] The Run Database file bytes are identical before and after Doctor.

## Context

- interface: `internal/cli/doctor.go`
- creates: `internal/cli/doctor_storage.go`
- creates: `internal/cli/doctor_storage_test.go`
- interface: `internal/cli/health.go`
- creates: `internal/store/free_bytes.go`
- interface: `docs/user-guide/commands.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `CONTEXT.md`
- instruction: `internal/cli/gc.go`
- instruction: `docs/adr/0172-doctor-reports-reclaimable-run-storage-from-cheap-reads.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestDoctorStorageReportsReclaimableRunsAndNamesGC|TestDoctorStorageResultNamesCompactAtTheFreeBytesThreshold|TestDoctorStorageIsOKWhenNothingIsLeftToReclaim|TestDoctorStorageDoesNotCreateAMissingRunDatabase|TestDoctorStorageIsPartialWhenTheRunDatabaseCannotBeRead|TestDoctorStorageLeavesTheRunDatabaseBytesUnchanged)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestDoctorStorageReportsReclaimableRunsAndNamesGC TestDoctorStorageResultNamesCompactAtTheFreeBytesThreshold TestDoctorStorageIsOKWhenNothingIsLeftToReclaim TestDoctorStorageDoesNotCreateAMissingRunDatabase TestDoctorStorageIsPartialWhenTheRunDatabaseCannotBeRead TestDoctorStorageLeavesTheRunDatabaseBytesUnchanged; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the tests exists, so the command fails.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands.md | grep -qF -- "never reports failed" || { printf 'missing phrase in %s: %s\n' docs/user-guide/commands.md "never reports failed" >&2; exit 1; }; for file in .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- "storage check" || { printf 'missing phrase in %s: %s\n' "$file" "storage check" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task the phrases are absent.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "whether Run storage is reclaimable" || { printf 'missing phrase in %s: %s\n' CONTEXT.md "whether Run storage is reclaimable" >&2; exit 1; }` — expected: exit 0; before this Task the phrase is absent.

## References

- [_prd.md](_prd.md) — Goal 3; Core Feature 2; Success Metric 3
- [_techspec.md](_techspec.md) — Doctor reports reclaimable storage; Interfaces; API Contract 2; Testing Approach 2; Build Order 3
- [references/2026-09-25-storage-reclaim-notice.md](references/2026-09-25-storage-reclaim-notice.md)
- ADR-0172; ADR-0171; ADR-0032; ADR-0107
