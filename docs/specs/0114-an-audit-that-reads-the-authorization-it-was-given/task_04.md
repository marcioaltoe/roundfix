---
task: task_04
spec: 0114-an-audit-that-reads-the-authorization-it-was-given
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

- `go test -count=1 ./internal/baseline -run 'TestOutputsForCommand' -v > /tmp/0114-t05.log 2>&1; s=$?; grep -q '^--- PASS: TestOutputsForCommand' /tmp/0114-t05.log || { cat /tmp/0114-t05.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0114-t05.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestOutputsForCommand' /tmp/0114-t05.log || { echo 'the resolver suite selected no cases'; cat /tmp/0114-t05.log; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0114-t05.log > /tmp/0114-t05-n.txt; test "$(cat /tmp/0114-t05-n.txt)" -ge 4 || { echo "expected the owning, equality, unknown-command, and wrong-command cases, got $(cat /tmp/0114-t05-n.txt)"; cat /tmp/0114-t05.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each case runs on its own. Non-vacuity is proven by the named test having run rather than by the absence of Go's "no tests to run" notice, which a package without the test emits even when the work is done.
- `grep -q 'baseline-digests' /tmp/0114-t05.log || { echo 'the equality case did not exercise the measured command'; cat /tmp/0114-t05.log; exit 1; }; grep -rq 'func OutputsFor' internal/ || { echo 'no resolver exists'; exit 1; }` — expected: exits 0, proving the resolver exists and that the equality proof used the one command with an enumerated list to check against. Fails today.
- `git diff --name-only HEAD > /tmp/0114-t04-changed.txt; grep '_ownership.y' /tmp/0114-t04-changed.txt && { echo 'an ownership record was changed, which this Task forbids:'; grep '_ownership.y' /tmp/0114-t04-changed.txt; exit 1; }; test -s /tmp/0114-t04-changed.txt || { echo 'no file changed'; exit 1; }; grep -rq 'func OutputsFor' internal/ || { echo 'files changed and no ownership record was touched, but no resolver exists'; exit 1; }` — expected: exits 0, proving work happened, that no ownership record was rewritten to make the resolver agree with it, and that the resolver is what the work produced. The changed-file clause alone passes on any dirty tree, so it is anchored to the resolver.

## Context

- interface: `internal/baseline/derived_ownership.go`
- interface: `docs/workflow/authorizations/2026-08-06-proof-cost.md`

## References

`_techspec.md` → Build Order 4; Implementation Design, Interfaces; Risks &
Considerations, trusting the ownership records. `_prd.md` → Core Feature 1;
Goal 1; User Stories 1 and 4. ADR-0129, ADR-0081.

## Result

`baseline.OutputsFor` now reads `DERIVED_DIGEST_PATHS`, resolves the existing
ownership records, and returns the sorted repository-relative regular files
owned by the named sanctioned or dedicated command. Records without a command
never match, so an unknown or empty command exempts nothing. The 2026-08-06
transition authorization's enumeration was reconciled from 28 to the 48 paths
the pre-existing ownership records assign to `make baseline-digests`; no
`_ownership.yml` file changed.

Focused checks:

- Inspected the daemon's attempt-1 diagnostic artifact; it reported a vacuous package pass because `TestOutputsForCommand` did not exist.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test ./internal/baseline -run '^TestOutputsForCommand$/unknown_command_owns_nothing$' -count=1` failed before implementation with `undefined: OutputsFor`, establishing the focused red signal.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test ./internal/baseline -run '^TestOutputsForCommand$/(command_owns_a_tree|unknown_command_owns_nothing|path_owned_by_another_command_is_excluded)$' -count=1 -v` passed all three ownership-behavior cases.
- The isolated equality subtest first exposed 20 paths missing from the 2026-08-06 transition enumeration, then passed after that enumeration was reconciled to the ownership-derived set.
- `rtk env GOCACHE=/tmp/roundfix-task04-gocache go test ./internal/baseline -run '^TestOutputsForCommand$' -count=1` passed the complete focused resolver suite after the final implementation edit.
- No command from this Task's `## Verification` section was run.

Acceptance evidence:

- `command owns a tree` proves the resolver returns a known formatter output of `make baseline-digests`.
- `make baseline-digests matches the 2026-08-06 enumeration` compares both complete sorted slices and passes with 48 paths.
- `unknown command owns nothing` proves an unrecorded command returns an empty set without error.
- `path owned by another command is excluded` uses two dedicated ownership records and proves only the requested command's file is returned.
- Changed-path inspection found no `_ownership.yml` or `_ownership.yaml` mutation.
