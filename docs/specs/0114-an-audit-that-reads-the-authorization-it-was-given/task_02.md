---
task: task_02
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: pending # pending | in_progress | completed | failed — only implement-task changes this
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
