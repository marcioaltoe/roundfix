---
task: task_16
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: docs
complexity: low
---

# Task 16: Correct the Source Baseline carrier an adopting repository consumes

## Overview

The tooling half of the carrier gap the QA gate found. The embedded Source
Baseline's routing carrier still names `_archived/specs/`, so every repository
adopting this Baseline consumes the obsolete destination. It is its own Task
because a tooling Task's entire commit is audited against the authorization
record, and mixing it with production Go put paths in that commit the grant does
not govern.

## Requirements

1. MUST make the embedded Source Baseline's routing carrier name the destination
   the built-in Spec Root actually uses.
2. MUST regenerate the derived Source Baseline artifacts through their sanctioned
   command rather than by editing a digest or index by hand.
3. MUST NOT change what the routing carrier requires; only where it says the
   archive lives.
4. MUST NOT change any repository path outside the bounded scope below plus this
   Task file; stop and fail the Task if a changed-file check finds another path.

## Subtasks

- [ ] Correct the routing carrier's archive destination.
- [ ] Regenerate the derived Source Baseline artifacts.

## Acceptance Criteria

- [ ] The routing carrier names `docs/history/specs/` and no longer names
      `_archived/specs`.
- [ ] The derived Source Baseline artifacts match their canonical source.
- [ ] The carrier's other requirements are unchanged.
- [ ] The changed-file set is the bounded scope plus this Task file.

## Bounded scope

This Task may create or modify only:

- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md`
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/baseline.json`
- `internal/baseline/assets/source-baselines/index.json`
- `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/manifest.json`
- `docs/specs/0094-one-history-root-under-docs/task_16.md`

Express maintainer authorization:
`docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`. The
first path is granted there for the archive location. The three derived artifacts
are fallout of that edit under ADR-0081 and are listed because the changed-path
audit reads the record's exact paths and does not implement the fallout rule.

## Verification

- `f=internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/docs/agents/spec-routing.md; if grep -n '_archived/specs' "$f"; then echo "FAIL: $f still names the old destination on the line above"; exit 1; fi; grep -q 'docs/history/specs' "$f" || { echo "FAIL: $f never names docs/history/specs"; exit 1; }` — expected: exits 0, printing the offending line or the missing string on failure so a repair turn has a cause to act on.
- `go test -count=1 -tags repocontract ./internal/baseline -run 'TestDeclaredStepRegenerationAndFrozenBoundaries' -v > /tmp/0094-task-16.log 2>&1; s=$?; grep -q '^--- PASS: TestDeclaredStepRegenerationAndFrozenBoundaries' /tmp/0094-task-16.log || { cat /tmp/0094-task-16.log; exit 1; }; exit $s` — expected: exits 0, proving the derived artifacts match their canonical source after regeneration.
- `git diff --name-only HEAD > /tmp/0094-task-16-all.txt; test -s /tmp/0094-task-16-all.txt || { echo 'no file changed; this Task edits a carrier'; exit 1; }; grep -v -e '^internal/baseline/assets/source-baselines/' -e '^docs/specs/0094-one-history-root-under-docs/task_16\.md$' /tmp/0094-task-16-all.txt > /tmp/0094-task-16-scope.txt; test ! -s /tmp/0094-task-16-scope.txt || { echo 'out of bounds:'; cat /tmp/0094-task-16-scope.txt; exit 1; }` — expected: exits 0, proving work happened and every changed path is inside the grant. It names the offending paths on failure.

## Context

- instruction: `docs/workflow/authorizations/2026-08-12-the-archive-root-under-docs.md`
- instruction: `docs/workflow/baseline-digest-regeneration.md`

## References

`_prd.md` → Core Feature 10. QA report `qa/qa-report-2026-08-13.md` → F-001, the
Source Baseline half. ADR-0081, ADR-0120.
