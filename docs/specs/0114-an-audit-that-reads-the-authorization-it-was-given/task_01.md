---
task: task_01
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 01: Stop refusing the row the template teaches

## Overview

The PRD template tells an author to write the Tooling authority row as
`applicable — no protected tooling mutation proposed or authorized`, and the
checker refuses exactly that. This Spec's own PRD was refused twice for it and its
TechSpec once more, each time resolved by writing the row in the wording ADR-0131
replaces. The row states whether the constraint governs the Spec; bounded files
become required only when the row declares a mutation.

## Requirements

1. MUST accept a Tooling authority row recorded as applicable whose reason
   declares no proposed or authorized mutation.
2. MUST still refuse a row that declares a proposed or authorized mutation
   without an exact bounded files list.
3. MUST still refuse a row citing an authorization record that does not name the
   Spec.
4. MUST still refuse a row whose cited record states its grant only in prose.
5. MUST NOT change the PRD template, which ADR-0131 keeps as the correct reading.

## Subtasks

- [ ] Require bounded files from the declared mutation rather than from the word
      applicable.
- [ ] Cover the template's exact wording and each surviving refusal.

## Acceptance Criteria

- [ ] The template's verbatim no-mutation wording passes `spec check --stage prd`.
- [ ] A declared mutation without bounded files still fails.
- [ ] A record that does not name the Spec still fails.
- [ ] A prose-only grant still fails.
- [ ] The PRD template is unchanged.

## Verification

- `go test -count=1 ./internal/speccheck -run 'TestToolingRowStatesApplicability' -v > /tmp/0114-t02.log 2>&1; s=$?; grep -q '^--- PASS: TestToolingRowStatesApplicability' /tmp/0114-t02.log || { cat /tmp/0114-t02.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0114-t02.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0114-t02.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^--- PASS\|^    --- PASS' /tmp/0114-t02.log > /tmp/0114-t02-n.txt; test "$(cat /tmp/0114-t02-n.txt)" -ge 4 || { echo "expected the accepted row and the three surviving refusals as their own cases, got $(cat /tmp/0114-t02-n.txt)"; cat /tmp/0114-t02.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each direction is its own case.
- `git diff --name-only HEAD -- .agents/skills/write-prd skills/write-prd > /tmp/0114-t02-tmpl.txt; test ! -s /tmp/0114-t02-tmpl.txt || { echo 'the template was changed, which this Task forbids:'; cat /tmp/0114-t02-tmpl.txt; exit 1; }; grep -rq 'no protected tooling mutation' internal/speccheck/constraints.go || { echo 'the template is untouched, but the detector does not read the declared mutation'; exit 1; }` — expected: exits 0, proving the repair landed in the checker and not in the template. Fails today on the second clause.

## Context

- interface: `internal/speccheck/constraints.go`
- instruction: `.agents/skills/write-prd/references/prd-template.md`

## References

`_techspec.md` → Build Order 1; API Contracts.
`_prd.md` → Core Feature 3; Core Feature 4; Goal 3; User Story 3.
ADR-0131, ADR-0117.
