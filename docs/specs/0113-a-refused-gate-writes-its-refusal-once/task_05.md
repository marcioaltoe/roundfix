---
status: completed
type: backend
---

# Task: Read Newest Report Only

Read only newest report, ignore superseded.

## Work
- List all `qa-report-*.md`
- Sort by filename (newest first)
- Select first (newest)
- Ignore all older

## Verification
- `grep -q "NewestQAReport" internal/speccheck/mechanical.go && go test -count=1 ./internal/speccheck ./internal/spec 2>&1 | grep -q "^ok"`


## References

- Core Feature 3: Newest Report Only

## Result

The mechanical stage now decides for itself which QA Report is current. It lists
the `qa-report-*.md` family in the directory its caller pointed at, orders that
family by the shared recency key, reads the newest member, and ignores every
older one — so a refusal or a `fail` row a previous run wrote is read by nothing
and cannot block the run that supersedes it.

### What changed

Before, `RunMechanicalStage` validated exactly the path its caller handed it
(`MechanicalRequest.ReportPath`). Which report got validated was therefore the
caller's answer, not the stage's: the Daemon resolves that path once at Run start
(`internal/daemon/task_context.go:45`) and passes it as `previousReportPath`
(`internal/daemon/task_engine.go:2073`), so a superseded report named there was
validated and its findings blocked the gate.

`newestMechanicalReportPath` (`internal/speccheck/mechanical.go:880`) now sits
between the request and `loadMechanicalReport`. Given a report path in the
`qa-report-*.md` family it globs that family in the report's own directory and
returns the newest by `spec.NewestQAReportFromPaths` — the same
date-then-run-sequence key `spec.NewestQAReport` applies to the filesystem and
`internal/worktree/worktree.go:1202` applies to a Git tree, so the stage and the
callers that hand it a path cannot come to different answers about which report
is current.

Three boundaries are deliberate:

- **Recency is the shared key, not filename order.** Raw byte order puts
  `qa-report-2026-08-25.md` above `qa-report-2026-08-25-02.md`, because `.`
  sorts above `-` — it would read a date's first report as newer than every
  rerun that supersedes it. A subtest pins the ordering directly rather than
  trusting the key by reference.
- **A path outside the family is read as named.** A fixture or a draft under
  another name has no family to be the newest of, and silently redirecting that
  read would discard the report the caller meant. Every existing mechanical
  fixture (`report-green.md`, `report-red.md`, `report-raw.md`) is read exactly
  as before.
- **A carried row's establishing report is not resolved through here.**
  `resolveCarriedRows` still loads its cited path exactly, because a citation
  names one report on purpose; pointing it at the newest would move the evidence
  a carried row rests on.

### Evidence per Work item

Red first: with `internal/speccheck/mechanical.go` stashed and the new test in
place, `go test -count=1 -run TestMechanicalStageReadsTheNewestQAReportOnly
./internal/speccheck` failed on four of six subtests — `a superseded refusal does
not block the run that supersedes it` (the stage raised `QA-REPORT-SHAPE` against
`qa-report-2026-08-14.md`, the superseded report), `the newest report is the one
read, not merely the older one skipped`, `recency is the date and run sequence,
not raw filename order` (it read the defective `qa-report-2026-08-25.md` over its
`-10` rerun), and `a requested report that is gone yields to the newest one on
disk`. The two that passed unchanged are the regression guards: the
outside-the-family read and the missing-report skip.

Focused check after the last edit: `go test -count=1 -v -run
TestMechanicalStageReadsTheNewestQAReportOnly ./internal/speccheck` → `ok`, all
six subtests running and passing.

| Work item | Evidence |
| --- | --- |
| List all `qa-report-*.md` | `newestMechanicalReportPath` globs `qa-report-*.md` in the requested report's directory. Subtest `a requested report that is gone yields to the newest one on disk` proves the listing is of the directory and not of the requested name: the requested `qa-report-2026-08-13.md` was never written, and the finding comes from `qa-report-2026-08-15.md` |
| Sort by filename (newest first) | Ordering is `spec.NewestQAReportFromPaths`. Subtest `recency is the date and run sequence, not raw filename order` writes three same-date reports — an unsuffixed one and `-02` both defective, `-10` valid — and requires zero findings, which fails under raw name order and under a lexicographic sequence compare |
| Select first (newest) | Subtest `the newest report is the one read, not merely the older one skipped` requests `qa-report-2026-08-14.md` (valid) beside a defective `qa-report-2026-08-15.md` and requires exactly one finding, naming row `R02` with `File` = the newest path. Selection is proven by the finding raised, not by findings being absent |
| Ignore all older | Subtest `a superseded refusal does not block the run that supersedes it` requests the older defective report beside a newer valid one and requires no `QA-REPORT-SHAPE` and `Blocking = false` — the measured deadlock, with no artifact deleted |

Supporting checks, all after the last edit: `gofmt -l internal/speccheck/` → no
output; `go build ./...` → clean; `go vet ./internal/speccheck ./internal/spec` →
clean; `go test -count=1 ./internal/speccheck ./internal/spec` → `ok`, `ok`;
`go test -count=1 ./internal/daemon` → `ok` (the stage's production caller);
`go test -count=1 ./...` → every package `ok`, no failures.

Changed paths: `internal/speccheck/mechanical.go`,
`internal/speccheck/mechanical_test.go`, and this Task file. No tooling
configuration, derived artifact, or other Spec file was touched.

### Follow-ups for later Tasks

- The Daemon's `previousReportPath` is now belt-and-braces rather than
  load-bearing: it still resolves the newest report at Run start, and the stage
  re-resolves it at read time. Nothing needs it removed, but a later reader
  should know the stage no longer depends on the caller being right.
- Still nothing calls `spec.WritePreconditionRefusalReport` in production — the
  Daemon writes its own `QA-PRECONDITION` row instead of the refusal report whose
  shape task_01 minted and task_04 taught the stage to accept. Every Task in
  `_tasks.md` except the qa gate is now implemented, and no Task claims that
  wiring.
- Identifier strategy: this Task introduced no glossary term. It reuses the
  existing QA Report naming contract (`qa-report-YYYY-MM-DD[-NN].md`) and the
  recency key that already governs it. Whether `CONTEXT.md`'s `QA Report` entry
  should name the precondition count and its metadata keys remains the closing
  node's check, as task_01, task_03, and task_04 recorded.
- `golangci-lint` is not installed in this environment, so lint coverage comes
  from the Daemon's Verification run.
