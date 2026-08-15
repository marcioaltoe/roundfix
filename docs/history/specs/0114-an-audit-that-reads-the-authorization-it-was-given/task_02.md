---
task: task_02
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 02: Tell a governed path from an ordinary one

## Overview

The audit reads every path a Task commit touched as bound by the tooling rules,
so an ordinary Go file alongside an authorized asset fails a Task whose only fault
was being one commit. Spec 0095 split a Task for exactly that. The universal
clause names the governed class by kind and enumerates it nowhere; this slice
gives the checker a predicate that can answer.

## Requirements

1. MUST decide whether a repository path is governed by the tooling rules, from a
   declared set covering every kind the universal clause names: linter, formatter,
   typechecker, test-runner, architecture-checker, build-tool, package-manager and
   code-generator configuration and scripts, ignore files, plugin declarations,
   and version pins.
2. MUST answer false for ordinary source and test files, and for documentation
   that no tooling rule binds.
3. MUST record, beside the set, which clause each entry comes from, so a later
   reader can check the set against the rule rather than against intuition.
4. MUST NOT read the set from an editable data file; it is compiled in, so
   changing it requires the review that changes code.
5. MUST NOT yet change what the audit refuses; this slice only answers the
   question.

## Subtasks

- [ ] Declare the governed set with its per-entry clause citation.
- [ ] Implement the predicate.
- [ ] Cover one case per named kind and the ordinary-path cases.

## Acceptance Criteria

- [ ] Each kind the universal clause names has at least one governed example that
      answers true.
- [ ] An ordinary Go source file, a test file, and a Spec document answer false.
- [ ] Every entry in the set cites the clause it comes from.
- [ ] The audit's behaviour is unchanged by this Task alone.

## Verification

- `go test -count=1 ./internal/speccheck -run 'TestGovernedPath' -v > /tmp/0114-t03.log 2>&1; s=$?; grep -q '^--- PASS: TestGovernedPath' /tmp/0114-t03.log || { cat /tmp/0114-t03.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0114-t03.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0114-t03.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0114-t03.log > /tmp/0114-t03-n.txt; test "$(cat /tmp/0114-t03-n.txt)" -ge 12 || { echo "expected a case per governed kind plus the ordinary-path cases, got $(cat /tmp/0114-t03-n.txt)"; cat /tmp/0114-t03.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving the kinds are covered individually rather than as one assertion.
- `grep -rq 'func GovernedPath' internal/speccheck || { echo 'no predicate exists'; exit 1; }; for kind in linter formatter typechecker test-runner build-tool package-manager code-generator 'ignore file' 'plugin declaration' 'version pin'; do grep -qi "$kind" internal/speccheck/*.go || { echo "FAIL: the governed set does not cite the $kind clause"; exit 1; }; done` — expected: exits 0, proving the predicate exists and every named kind is accounted in the set's own text. Fails today.

## Context

- instruction: `docs/agents/agent-instructions.md`

## References

`_techspec.md` → Build Order 2; Implementation Design, Interfaces; Risks &
Considerations, the declared set. `_prd.md` → Core Feature 2; Goal 2; User
Story 2; Open Questions. ADR-0130.

## Result

### Implementation

- Added a compiled governed-path set with one cited entry for each of the eleven
  tooling kinds named by the universal Tooling Authority clause. The set covers
  conventional repository-relative configuration and script paths without
  reading an editable data file.
- Added `GovernedPath`, which applies the package's existing repository-path
  normalization and returns whether any fully cited set entry matches.
- Added fourteen named table cases: one positive case per governed kind and the
  required negative cases for ordinary Go source, an ordinary Go test, and a
  Spec document.
- Kept the changed-path audit untouched. This slice introduces the predicate but
  does not call it from the audit; Task 03 owns the historical-union contract and
  Task 05 owns audit integration.

### Focused-check evidence

- Before the implementation, `rtk env
  GOCACHE=/tmp/roundfix-task-02-go-cache go test ./internal/speccheck -run
  '^TestGovernedPath$'` exited 1 because `speccheck.GovernedPath` was undefined.
- After the implementation, the same focused command exited 0 and reported
  `ok roundfix/internal/speccheck`.
- `rtk env GOCACHE=/tmp/roundfix-task-02-go-cache make verify-incremental`
  initially reached the changed package successfully but exited 2 in the sandbox
  because two unrelated force-stop integration tests could not read the host
  process table. The unsandboxed rerun exited 0; all Go packages, skill checks,
  and the build passed.
- `rtk rg -c '^\s*kind:' internal/speccheck/governed.go && rtk rg -c
  '^\s*clause:' internal/speccheck/governed.go` exited 0 and printed `11` for
  both counts. This is direct evidence that every declared entry carries its
  clause citation.
- `rtk rg -l 'GovernedPath' internal/speccheck --glob '*.go'` exited 0 and named
  only `governed.go` and `governed_test.go`, proving this slice added no audit
  call site.

### Acceptance evidence

- Each named tooling kind has a true case: the eleven governed subtests passed in
  `TestGovernedPath` and in the incremental repository gate.
- The ordinary Go source, ordinary test, and Spec document subtests all passed
  with expected value `false`.
- The governed set has eleven entries and eleven adjacent clause fields; an entry
  without both its kind and clause cannot match.
- Audit behavior is unchanged because the predicate has no production caller;
  the incremental repository gate passed after the final code edit.

### Not run

- The three commands under `## Verification` are reserved for the Roundfix
  Daemon and were not run in this Agent turn.
