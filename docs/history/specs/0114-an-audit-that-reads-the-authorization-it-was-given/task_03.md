---
task: task_03
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Result

### Implementation

- Added a repository-contract test that scans every Markdown authorization
  record under `docs/workflow/authorizations/` and reads bounded paths through
  the checker’s existing typed-frontmatter and legacy `Bounded files` parser.
- The contract reports every unmatched path with the authorization record that
  bounded it. A synthetic record proves the diagnostic names both values.
- Added the exact historical exceptions exposed by the live authorization
  corpus to the compiled governed set. The entry cites ADR-0130, and exact
  matching keeps neighboring ordinary Go and documentation paths ungoverned.
- Added an explicit no-records subtest that records a Go test skip when the
  authorization directory is absent.

### Focused-check evidence

- Before the governed set changed, `rtk env
  GOCACHE=/tmp/roundfix-task-03-go-cache go test -tags repocontract
  ./internal/speccheck -run
  '^TestEveryBoundedPathIsGoverned/repository_records_are_governed$'` exited 1
  and named every historical path missing from the declared set together with
  its authorization record.
- After the governed set changed, the same focused repository-record command
  exited 0 and reported `ok roundfix/internal/speccheck`.
- `rtk env GOCACHE=/tmp/roundfix-task-03-go-cache go test -tags repocontract
  ./internal/speccheck -run
  '^TestEveryBoundedPathIsGoverned/(unmatched_path_names_path_and_record|no_authorization_records_skips)$'
  -v` exited 0. The unmatched case passed after observing one finding naming
  `README.md` and `docs/workflow/authorizations/synthetic.md`; the no-records
  case reported `SKIP`.
- `rtk env GOCACHE=/tmp/roundfix-task-03-go-cache go test
  ./internal/speccheck -run '^TestGovernedPath$'` exited 0, retaining the
  original kind-based and ordinary-path behavior.
- `rtk env GOCACHE=/tmp/roundfix-task-03-go-cache go test
  ./internal/speccheck` exited 0 for the complete package suite.
- The sandboxed `rtk env GOCACHE=/tmp/roundfix-task-03-go-cache make
  verify-incremental` reached `ok roundfix/internal/speccheck` but exited 2
  because two existing force-stop integration tests could not read the host
  process table. The rerun with host process-table access exited 0; all Go
  packages, skill checks, and the build succeeded.

### Acceptance evidence

- The live-corpus subtest reads the authorization directory and exited 0 only
  after every currently parsed bounded path matched `GovernedPath`.
- The synthetic unmatched-path subtest observed exactly one finding and checked
  that it named both the unmatched path and its source record.
- The empty-repository subtest exercised the absent-directory branch and
  reported a skip rather than a failure.

### Not run

- The three commands under `## Verification` are reserved for the Roundfix
  Daemon and were not run in this Agent turn.
