---
task: task_01
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 01: Classify a retired decision apart from a pending one

## Overview

One pure function per family answers whether a decision record or a typed intent
entry has retired. It is deliberately narrower than the existing inactive-status
predicate: that one answers whether a decision is in force, and reusing it would
file a pending proposal as history. Verifiable on its own through table-driven
tests over frontmatter fixtures, with no filesystem and no caller yet.

## Requirements

1. MUST report a decision record retired when its lifecycle status is
   `rejected`, `deprecated`, or `superseded`.
2. MUST report a decision record active when its lifecycle status is `proposed`
   or `accepted`.
3. MUST report a legacy decision record active when it carries no lifecycle
   status and its body does not mark it inactive.
4. MUST report a typed intent entry retired when its status is `declined`, and
   active when its status is `open` or `promoted`.
5. MUST name the status that retired the record, so a caller can report the
   reason rather than only the verdict.
6. MUST NOT read or touch the filesystem; the input is document content.

## Subtasks

- [ ] Add the retirement classification for decision records.
- [ ] Add the retirement classification for typed intent entries.
- [ ] Cover every lifecycle status of both families with table-driven tests.
- [ ] Cover the legacy no-status decision record and the body-marked case.

## Acceptance Criteria

- [ ] A `proposed` decision record classifies as active, exercised as a named
      `proposed` subtest that fails if it is ever classified retired.
- [ ] Each of `rejected`, `deprecated`, and `superseded` classifies as retired
      and names itself as the reason.
- [ ] A decision record with no lifecycle status classifies as active.
- [ ] A `declined` intent entry classifies as retired; `open` and `promoted`
      classify as active.
- [ ] The classification compiles and is exercised without any filesystem access
      in its tests.

## Verification

- `grep -q 'func ClassifyADR' internal/spec/*.go && grep -q 'func ClassifyBacklogEntry' internal/spec/*.go` — expected: exits 0, proving both classifications exist. Fails on a tree where no work has happened, because neither function is declared today.
- `go test -count=1 ./internal/spec -run 'TestClassifyADRRetirement|TestClassifyBacklogEntryRetirement' -v > /tmp/0094-task-01.log 2>&1; s=$?; grep -q '^--- PASS: TestClassifyADRRetirement' /tmp/0094-task-01.log && grep -q '^--- PASS: TestClassifyBacklogEntryRetirement' /tmp/0094-task-01.log || { cat /tmp/0094-task-01.log; exit 1; }; exit $s` — expected: exits 0 and the log names both passing tests. The names are exact rather than substring patterns: `-run 'Retire'` matches the existing `…EveryRetiredFamily` and `…EveryRetiredKind` tests, which is how an earlier draft of this gate passed before any work was done.
- `grep -q '^[[:space:]]*--- PASS: TestClassifyADRRetirement/proposed' /tmp/0094-task-01.log` — expected: exits 0, proving the pending-proposal case named by ADR-0122 is exercised as its own named subtest rather than only declared. The leading-whitespace allowance matches Go's standard indentation for nested subtest results. Reads the log the previous command wrote, so no pipeline can hide a test's exit status.

## References

`_techspec.md` → Interfaces: `ClassifyADR`, `ClassifyBacklogEntry`; Build Order 1;
Testing Approach: retirement classification. `_prd.md` → Core Feature 2. ADR-0122.

## Result

Implemented a pure `Retirement` classifier for decision records and typed intent
entries. Decision records retire only for `rejected`, `deprecated`, and
`superseded`; typed intent entries retire only for `declined`. A retired result
names its lifecycle status in `Reason`, while active and pending records return an
empty result. Legacy decision records remain active unless a body status line
marks one with a retiring decision status.

Focused checks and acceptance evidence:

- `rtk go test ./internal/spec -run '^TestClassify(ADR|BacklogEntry)Retirement$' -count=1`
  passed with 12 tests. The table exercises all five decision lifecycle statuses,
  all three typed intent statuses, the legacy no-status record, and the legacy
  body-marked record.
- `rtk go test ./internal/spec -run '^TestClassifyADRRetirement$/^proposed$' -count=1 -v`
  passed with the parent and named `proposed` subtest selected, proving the pending
  proposal remains active.
- Source inspection confirms the retired decision rows expect their own status as
  `Reason`, the `declined` row expects `Reason: "declined"`, and active rows expect
  an empty `Retirement`.
- The implementation and tests accept in-memory `[]byte`/string content and import
  no filesystem packages. The focused Go test compiles and exercises both public
  classifiers without creating fixtures on disk.
- `rtk git diff --check` exited 0 after the implementation edit.
- `rtk make verify-incremental` exited 0, covering formatting, the Go suite,
  skill sync and validation, and the build with reusable caches.

Verification feedback repair:

- Attempt 1's classifier tests passed, including the named `proposed` subtest,
  but the final detector required a nested Go subtest result at column 1. Go
  indents nested `--- PASS` lines, so the detector could not match genuine test
  output.
- The detector now permits Go's leading subtest indentation while retaining the
  exact `TestClassifyADRRetirement/proposed` name. No production or test behavior
  changed.
- `rtk go test ./internal/spec -run '^TestClassifyADRRetirement$/^proposed$' -count=1`
  passed with 2 tests after the detector repair.
- `rtk make verify-incremental` exited 0 after the detector repair; the Go suite,
  skill checks, and build passed.

The Daemon-owned Verification commands were not run.
