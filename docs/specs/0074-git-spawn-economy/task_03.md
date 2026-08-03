---
task: task_03
spec: 0074-git-spawn-economy
status: pending
type: backend
complexity: medium
---

# Task 03: Batch object reads in assets-sync provenance

## Overview

The second per-file loop: `internal/baseline/assets_sync_git.go` reads
source files with one `git cat-file blob` spawn per tree entry while
computing the provenance digest. This Task applies the batch-reader shape
proven in task_02 to this loop, keeping the reader package-private and the
assets-sync error vocabulary intact.

## Requirements

1. MUST route the provenance loop through a batch reader with the same
   framing-exact semantics task_02 established.
2. MUST preserve the assets-sync error vocabulary: unsafe path, non-blob
   entries, and read failures report exactly as before.
3. MUST keep the provenance digest byte-identical for identical inputs —
   the existing sync and compatibility tests are the characterization and
   pass unmodified.
4. MUST cover the same framing edge cases where the shape differs from
   task_02's reader; shared cases need no duplication if the reader is
   shared within the package.

## Subtasks

- [ ] Route the loop through the batch reader.
- [ ] Preserve the error vocabulary.
- [ ] Confirm digests are byte-identical via the existing tests.

## Acceptance Criteria

- [ ] Provenance over N files spawns one `cat-file` process, not N.
- [ ] `TestAssetsSyncCompatibilityMatchesMaintainedPythonContract` and the
      whole assets-sync family pass unmodified.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/`
      and this task file.

## Verification

- `go test ./internal/baseline -count=1 -run 'AssetsSync' -v | grep -q -- "--- PASS"`
  — expected: exit 0.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `grep -n "cat-file" internal/baseline/assets_sync_git.go | grep -qv -- "--batch" && exit 1 || exit 0`
  — expected: exit 0; no per-file spawn remains.

## References

- `_prd.md` → Core Feature 1.
- `_techspec.md` → Build Order 3 (reuses the reader shape proven in
  task_02).
- ADR-0090.
