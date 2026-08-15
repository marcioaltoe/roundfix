---
task: task_03
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 03: Hold the governed set to what was actually protected

## Overview

A declared set is a list, and lists rot. The repository already holds the ground
truth: every path any maintainer has ever bounded in an authorization record is,
by definition, a path they considered governed. This slice makes that record the
test, so the set can grow deliberately and cannot silently narrow.

## Requirements

1. MUST fail when any path bounded by any authorization record under
   `docs/workflow/authorizations/` is not matched by the declared governed set.
2. MUST name every unmatched path and the record that bounded it.
3. MUST read the records rather than a copied list, so a new record is covered
   the day it lands.
4. MUST skip rather than fail where no authorization record exists, so a fresh
   repository is not refused.

## Subtasks

- [ ] Read every authorization record's bounded paths.
- [ ] Fail on any that the set does not match, naming both.
- [ ] Cover the matched, unmatched, and no-records cases.

## Acceptance Criteria

- [ ] Every bounded path in every record today is matched by the set.
- [ ] A synthetic record bounding an unmatched path fails, naming the path and
      the record.
- [ ] A repository with no authorization records skips rather than fails.

## Verification

- `go test -count=1 -tags repocontract -run 'TestEveryBoundedPathIsGoverned' ./internal/speccheck -v > /tmp/0114-t04.log 2>&1; s=$?; grep -q '^--- PASS: TestEveryBoundedPathIsGoverned' /tmp/0114-t04.log || { cat /tmp/0114-t04.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing contract test; fails today, where it does not exist.
- `test -s /tmp/0114-t04.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0114-t04.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0114-t04.log > /tmp/0114-t04-n.txt; test "$(cat /tmp/0114-t04-n.txt)" -ge 3 || { echo "expected the matched, unmatched, and no-records cases, got $(cat /tmp/0114-t04-n.txt)"; cat /tmp/0114-t04.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving all three directions run.
- `grep -rq 'func TestEveryBoundedPathIsGoverned' internal/speccheck || { echo 'the contract test does not exist'; exit 1; }; grep -rl 'docs/workflow/authorizations' internal/speccheck | grep -q 'repocontract' || { echo 'the contract does not read the record directory itself'; exit 1; }; n=$(grep -c '^  - ' docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md); test "$n" -ge 20 || { echo "expected the umbrella record to still bound its paths, found $n"; exit 1; }` — expected: exits 0, proving the contract exists, reads the record directory rather than a copied list, and that the largest record still carries the paths it is checked against. Fails today, where the contract does not exist; the record-count clause alone passes on an untouched tree, so it is anchored to the two above it.

## Context

- interface: `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`

## References

`_techspec.md` → Build Order 3; Risks & Considerations, the declared set.
`_prd.md` → Core Feature 2; Open Questions. ADR-0130.
