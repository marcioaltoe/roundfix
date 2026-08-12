---
task: task_01
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
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

- [ ] A `proposed` decision record classifies as active, with a test that fails
      if it is ever classified retired.
- [ ] Each of `rejected`, `deprecated`, and `superseded` classifies as retired
      and names itself as the reason.
- [ ] A decision record with no lifecycle status classifies as active.
- [ ] A `declined` intent entry classifies as retired; `open` and `promoted`
      classify as active.
- [ ] The classification compiles and is exercised without any filesystem access
      in its tests.

## Verification

- `go test -count=1 ./internal/spec -run 'Retire|Classif' -v > /tmp/0094-task-01.log 2>&1; s=$?; grep -q '^--- PASS: .*Retire\|^--- PASS: .*Classif' /tmp/0094-task-01.log || { cat /tmp/0094-task-01.log; exit 1; }; exit $s` — expected: exits 0 and the log shows at least one passing classification test; fails when the named tests do not exist or do not run.
- `! grep -rn 'inactiveStatusPattern' internal/spec` — expected: exits 0, proving the narrower rule is its own predicate rather than the in-force one reused.
- `go build -buildvcs=false ./...` — expected: exits 0.

## References

`_techspec.md` → Interfaces: `ClassifyADR`, `ClassifyBacklogEntry`; Build Order 1;
Testing Approach: retirement classification. `_prd.md` → Core Feature 2. ADR-0122.
