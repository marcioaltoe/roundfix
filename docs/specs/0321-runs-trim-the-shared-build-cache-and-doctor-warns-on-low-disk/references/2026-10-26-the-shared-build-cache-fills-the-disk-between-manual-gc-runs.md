---
type: fix
status: promoted
created: 2026-10-26
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
reason: null
---

# The shared build cache fills the disk between manual GC runs

## Symptom

The Go build cache this repository's Makefile shares across worktrees,
`{tmpdir}/roundfix-go-build-{uid}`, filled the Mac mini's data volume twice
in four days and stopped Runs mid-build. Spec 0314 (ADR-0327) made
`roundfix gc` trim a declared Shared Build Cache, but the trim runs only when
someone runs `roundfix gc` by hand, so the cache grows past its declared
`max_size_mb` between those runs. `roundfix doctor` says nothing when the disk
that holds the worktrees or the cache is almost full.

## Where

- `internal/cli/gc_shared_cache.go` `sweepSharedCaches`, called only by
  `runGC` in `internal/cli/gc.go`.
- `internal/sharedcache` `Measure` and `Apply`, which take no
  `context.Context` and so cannot be bounded in time by a caller.
- The Implement Run start in `internal/cli/implement.go`
  `runImplementCommand` (direct Runs and the Delivery Queue's
  `implement --spec` child alike) and the queue's repository gate,
  `commandDeliveryWorkflow.Gate` in `internal/cli/deliver_workflow.go`, the
  two places where new build output enters the shared cache.
- `internal/cli/doctor.go` `runDoctorCommand` and
  `internal/cli/doctor_storage.go`, which read the Run Database's free pages
  but never the free space of a filesystem.

## Expected

The declared Shared Build Cache trim, with exactly ADR-0327's rules, also runs
by itself at a point every Implement Run passes, best-effort, bounded in time,
and recorded with its outcome; and `roundfix doctor` warns when the
filesystem that holds the worktree location or a declared Shared Build Cache
is low on free space.

## Evidence

Operator records of the Roundfix queue on the Mac mini (this repository's
own delivery machine):

- 2026-10-22 23:30 BRT: the Run of Spec 0300 ended Unresolved because
  `task_02` failed with `No space left on device`; `/System/Volumes/Data` was
  100 % used with 399 MB free. The shared cache held 43.7 GB. The operator
  emptied it with `go clean -cache` and retried.
- 2026-10-25 23:21 BRT: the disk filled again, 157 MiB free, with the shared
  cache at 51 GB. The operator ran the GC Command of Spec 0314, which removed
  27 GB (26 GB of it from the cache). Spec 0315's item had parked
  `gate-failed` and Spec 0307's `make build` had failed, both for lack of
  space; both were retried.
- 2026-10-26 15:30 BRT: the cache held 17 GB again; on 2026-10-26 at about
  16:40 BRT `du -sh` read 20 GB in 37,975 files, with 48 GiB free on the
  228 GiB data volume, and the worktree location `~/.roundfix/worktrees` on
  a different filesystem.
- ADR-0327 (2026-10-24) records 33 GB on 2026-10-24, none of it older than
  fifteen hours, and that the go command deletes an entry only after five
  days without use.
- `roundfix doctor` at v0.123.0 has no line that reads a filesystem's free
  space: `grep -rn "Statfs\|GetDiskFreeSpace" internal` finds only
  `internal/store/compaction_capacity_unix.go` and
  `internal/store/compaction_capacity_windows.go`, the unexported reader
  `gc compact` uses for the Run Database's directory.
