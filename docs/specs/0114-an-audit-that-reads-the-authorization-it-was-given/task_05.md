---
task: task_05
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Result

The changed-path audit now ignores ordinary paths, resolves every declared
regeneration command through `baseline.OutputsFor`, and unions those paths with
any outputs the authorization record still enumerates. A command-only
`Sanctioned regeneration` declaration is retained by the shared parser so the
audit can resolve it from the repository tree. The authorization record itself
remains governed even if its path does not match `GovernedPath`, which keeps a
Task from landing its grant in the commit that consumes it. Git diagnostics are
also kept off the `diff-tree` stdout stream so a warning cannot become a
fabricated offending path.

Focused checks:

- Before implementation,
  `rtk env GOCACHE=/tmp/roundfix-0114-task05-gocache go test ./internal/speccheck -count=1 -run '^TestAuditJudgesTheGrant$/historical_authorized_asset_and_ordinary_Go_split_now_share_one_audit$'`
  failed with the ordinary `internal/cli/archive.go` path refused. The same run
  also exposed an fsmonitor diagnostic parsed as a path.
- Before implementation,
  `rtk env GOCACHE=/tmp/roundfix-0114-task05-gocache go test ./internal/speccheck -count=1 -run '^TestAuditJudgesTheGrant$/historical_regeneration_resolves_outputs_absent_from_the_grant$'`
  failed with nine regenerated paths refused when the record named only
  `make baseline-digests`.
- `rtk env GOCACHE=/tmp/roundfix-0114-task05-gocache go test -json ./internal/speccheck -count=1 -run '^TestAuditJudgesTheGrant$'`
  passed the two historical acceptance cases and the four refusal cases as six
  separately reported subtests.
- `rtk env GOCACHE=/tmp/roundfix-0114-task05-gocache go test ./internal/speccheck -count=1`
  passed the complete package after the final production edit.
- `rtk env GOCACHE=/tmp/roundfix-0114-task05-gocache go test ./internal/suiteguard ./internal/suiteguardcontract -count=1`
  passed; `suiteguardcontract` has no direct test files.
- `rtk env GOCACHE=/tmp/roundfix-0114-task05-gocache go test ./internal/baseline -count=1 -run '^TestOutputsForCommand$'`
  passed the resolver suite.
- `rtk git diff --check` exited 0, and source inspection found both
  `GovernedPath` and `baseline.OutputsFor` in `internal/speccheck/mechanical.go`.
- No command from this Task's `## Verification` section was run.

Acceptance evidence:

- The ordinary-path case replays the actual split commits
  `419a4661ac769ff7ee6ce5423bd795185c859d01` and
  `65c51ebf2e19220ff50d25fe03be809fcdf353f0`; the same subtest also replays
  Spec 0095 Task commit `28acf39cc193ad490646cb5a1d23500e0c08c273`
  against its authorization. All three commits resolved from this repository's
  Git object database, so the outside-evidence row is not blocked.
- The command-resolution case replays the 2026-08-13 corrective commit
  `c80e1266658929f68e8046af82f88e13392dc56d` against a command-only grant.
  The commit resolved from Git and its nine derived paths passed without an
  enumerated `outputs` list, so this outside-evidence row is not blocked.
- `governed path outside the grant is refused by name` changes a bounded
  `Makefile` beside unbounded `.golangci.yml` and asserts exactly one
  `QA-AUTH-PATHS` finding whose detail names the path and whose file names the
  escaped grant.
- `hand edited derived value without a command is refused` asserts exactly one
  finding for `internal/baseline/assets/source-baselines/index.json` when the
  grant declares no regeneration command.
- `record that does not name the Spec is refused` preserves the existing
  `SC-TOOLING-UNAUTHORIZED` refusal and requires its summary to name the
  omitted Spec.
- `Task commit cannot carry its own authorization` changes a bounded path and
  its authorization record together, then asserts exactly one finding naming
  the record and grant.
