---
task: task_04
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: high
---

# Task 04: Resolve a command's outputs from the tree

## Overview

A grant is asked to enumerate every output of its regeneration command, which is
the enumeration of consequences ADR-0081 rejected. The repository already records
that ownership where the artifacts live: `_ownership.yml` names each derived
tree's owner and the command that regenerates it. This slice lets a grant name the
command and the audit read the outputs from the tree.

## Requirements

1. MUST answer, for a named regeneration command, the repository-relative paths
   the ownership records say that command owns.
2. MUST prove the resolved set equals the enumerated set an existing record
   already carries for the same command, before anything relies on the resolver
   alone.
3. MUST return an empty set rather than an error for a command no record names,
   so an unknown command exempts nothing.
4. MUST NOT treat a path owned by a different command as belonging to this one.
5. MUST NOT change what `_ownership.yml` means or add a field to it.

## Subtasks

- [ ] Expose the resolver over the ownership records.
- [ ] Prove equality against the enumerated list it will replace.
- [ ] Cover the unknown-command and wrong-command cases.

## Acceptance Criteria

- [ ] The resolver returns the outputs of a command that owns a tree.
- [ ] For `make baseline-digests`, the resolved set equals the set the
      2026-08-06 record enumerates, proven by comparing them.
- [ ] An unknown command resolves to the empty set.
- [ ] A path owned by another command is not returned.
- [ ] `_ownership.yml` is unchanged.

## Verification

- `go test -count=1 ./internal/baseline ./internal/speccheck -run 'TestOutputsForCommand' -v > /tmp/0114-t05.log 2>&1; s=$?; grep -q '^--- PASS: TestOutputsForCommand' /tmp/0114-t05.log || { cat /tmp/0114-t05.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0114-t05.log || { echo 'the suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0114-t05.log && { echo 'the suite selected no cases'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0114-t05.log > /tmp/0114-t05-n.txt; test "$(cat /tmp/0114-t05-n.txt)" -ge 4 || { echo "expected the owning, equality, unknown-command, and wrong-command cases, got $(cat /tmp/0114-t05-n.txt)"; cat /tmp/0114-t05.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each case runs on its own.
- `grep -q 'baseline-digests' /tmp/0114-t05.log || { echo 'the equality case did not exercise the measured command'; cat /tmp/0114-t05.log; exit 1; }; grep -rq 'func OutputsFor' internal/ || { echo 'no resolver exists'; exit 1; }` — expected: exits 0, proving the resolver exists and that the equality proof used the one command with an enumerated list to check against. Fails today.
- `git diff --name-only HEAD > /tmp/0114-t04-changed.txt; grep '_ownership.y' /tmp/0114-t04-changed.txt && { echo 'an ownership record was changed, which this Task forbids:'; grep '_ownership.y' /tmp/0114-t04-changed.txt; exit 1; }; test -s /tmp/0114-t04-changed.txt || { echo 'no file changed'; exit 1; }; grep -rq 'func OutputsFor' internal/ || { echo 'files changed and no ownership record was touched, but no resolver exists'; exit 1; }` — expected: exits 0, proving work happened, that no ownership record was rewritten to make the resolver agree with it, and that the resolver is what the work produced. The changed-file clause alone passes on any dirty tree, so it is anchored to the resolver.

## Context

- interface: `internal/baseline/derived_ownership.go`
- interface: `docs/workflow/authorizations/2026-08-06-proof-cost.md`

## References

`_techspec.md` → Build Order 4; Implementation Design, Interfaces; Risks &
Considerations, trusting the ownership records. `_prd.md` → Core Feature 1;
Goal 1; User Stories 1 and 4. ADR-0129, ADR-0081.
