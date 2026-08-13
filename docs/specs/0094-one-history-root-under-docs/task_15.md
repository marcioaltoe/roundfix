---
task: task_15
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: docs
complexity: low
---

# Task 15: Correct the Archive Command's help text

## Overview

The QA gate found Core Feature 10's every-carrier promise unmet by two files
Task 07's sweep did not reach: the Archive Command's own help text, and the
embedded Source Baseline corpus that adopting repositories consume. A maintainer
running `archive --help` is told the built-in root archives to `_archived/specs/`
by the same binary that writes `docs/history/specs/`.

## Requirements

1. MUST make the Archive Command's help text name the destination the binary
   actually writes for the built-in Spec Root.
2. MUST NOT change the destination itself, the archive contract, or any behaviour;
   this slice corrects what two carriers say.
4. MUST NOT change any repository path outside the bounded scope below plus this
   Task file; stop and fail the Task if a changed-file check finds another path.

## Subtasks

- [ ] Correct the Archive Command help text.

## Acceptance Criteria

- [ ] `archive --help` names `docs/history/specs/` for the built-in Spec Root.
- [ ] The help text no longer names `_archived/specs`.
- [ ] No behaviour changed; the destination the binary writes is unchanged.

## Verification

- `f=internal/cli/archive.go; if grep -n '_archived/specs' "$f"; then echo "FAIL: $f still names the old destination on the line above"; exit 1; fi; grep -q 'docs/history/specs' "$f" || { echo "FAIL: $f never names docs/history/specs"; exit 1; }` — expected: exits 0, proving the help text lost the old destination and names the new one. It prints the offending line or the missing string on failure: an earlier draft used a silent `! grep -q … && grep -q …` pair, and its empty diagnostic left the repair turn with nothing to act on.
- `go build -buildvcs=false -o /tmp/0094-task-15-roundfix ./cmd/roundfix && /tmp/0094-task-15-roundfix archive --help 2>&1 | grep -q 'docs/history/specs'` — expected: exits 0, proving the corrected text reaches the built binary's own output rather than only the source file.
- `git diff --name-only HEAD > /tmp/0094-task-15-all.txt; test -s /tmp/0094-task-15-all.txt || { echo 'no file changed'; exit 1; }; grep -v -e '^internal/cli/archive\.go$' -e '^docs/specs/0094-one-history-root-under-docs/task_15\.md$' /tmp/0094-task-15-all.txt > /tmp/0094-task-15-scope.txt; test ! -s /tmp/0094-task-15-scope.txt || { cat /tmp/0094-task-15-scope.txt; exit 1; }` — expected: exits 0, proving work happened and every changed path is in bounds.

## Context

- interface: `internal/cli/archive.go`
- interface: `internal/spec/archive.go`

## References

`_prd.md` → Core Feature 10; User Story 1. The Source Baseline half of this
Spec's carrier gap is task_16, split out because a tooling Task's whole commit is
audited against the grant and this file needs none. QA report
`qa/qa-report-2026-08-13.md` → F-001. ADR-0120.
