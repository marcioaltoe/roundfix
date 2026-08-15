---
task: task_05
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 05: Let the audit judge the grant it was given

## Overview

Where the three repairs meet. The audit stops reading every changed path as
governed, resolves a command's outputs from the tree instead of from an
enumerated list, and keeps refusing everything it should. The two measured
refusals this Spec exists to remove must both pass; the unauthorized cases must
all still fail.

## Requirements

1. MUST judge only governed paths, leaving an ordinary changed path unaudited.
2. MUST treat a path the named command owns as sanctioned fallout, resolved from
   the tree, with any enumerated list still honoured as a union.
3. MUST still refuse a governed path no grant reaches.
4. MUST still refuse a hand-edited derived value, which no command regenerated.
5. MUST still refuse a record that does not name the Spec, and a Task commit that
   folds an authorized change together with its own authorization.
6. MUST report one finding per offending path, naming the path and the grant it
   escaped.

## Subtasks

- [ ] Read the predicate and the resolver in the audit.
- [ ] Keep every existing refusal.
- [ ] Cover the two measured passes and the four refusals end to end.

## Acceptance Criteria

- [ ] A Task commit touching one authorized asset and one ordinary Go file passes
      without being split — the refusal measured during Spec 0095.
- [ ] A Task commit regenerating derived artifacts through the sanctioned command
      passes with those artifacts absent from the record — the refusal measured on
      2026-08-13.
- [ ] A governed path outside every grant still fails, named.
- [ ] A hand-edited derived value still fails.
- [ ] A record that does not name the Spec still fails.
- [ ] **Outside evidence.** The two passing cases are replayed from commits this
      Spec did not author: the Spec 0095 Task commit that was split to satisfy the
      audit, and the 2026-08-13 refusal recorded in the backlog. Source: this
      repository's own history, read with `git`, not fixtures written here. Record
      the row as blocked with its reason if either commit cannot be resolved.

## Verification

- `go test -count=1 ./internal/speccheck -run 'TestAuditJudgesTheGrant' -v > /tmp/0114-t06.log 2>&1; s=$?; grep -q '^--- PASS: TestAuditJudgesTheGrant' /tmp/0114-t06.log || { cat /tmp/0114-t06.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0114-t06.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0114-t06.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0114-t06.log > /tmp/0114-t06-n.txt; test "$(cat /tmp/0114-t06-n.txt)" -ge 6 || { echo "expected the two passes and the four refusals as their own cases, got $(cat /tmp/0114-t06-n.txt)"; cat /tmp/0114-t06.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving both directions are covered case by case.
- `grep -rq 'GovernedPath' internal/speccheck/mechanical.go && grep -rq 'OutputsFor' internal/speccheck/mechanical.go || { echo 'the audit does not read both the predicate and the resolver'; exit 1; }` — expected: exits 0, proving the composition happened where the audit lives. Fails today.
- `go test -count=1 ./internal/speccheck ./internal/suiteguard ./internal/suiteguardcontract ./internal/baseline > /tmp/0114-t06-regress.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0114-t06-regress.log && { echo 'a neighbouring package regressed:'; grep -B 3 -A 6 'FAIL' /tmp/0114-t06-regress.log | head -30; exit 1; }; grep -rq 'GovernedPath' internal/speccheck/mechanical.go || { echo 'the packages pass, but the audit never read the predicate'; exit 1; }; exit $s` — expected: exits 0, proving the audit's neighbours still pass and that they pass with the composition in place rather than without it.

## Context

- interface: `internal/speccheck/mechanical.go`
- instruction: `docs/backlog/2026-08-13-the-changed-path-audit-does-not-know-sanctioned-fallout.md`

## References

`_techspec.md` → Build Order 5; Coverage Map; Risks & Considerations, the
direction that loses safety. `_prd.md` → Core Features 1 and 2; Goals 1 and 2;
Success Metrics. ADR-0129, ADR-0130, ADR-0081.
