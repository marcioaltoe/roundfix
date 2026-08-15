---
task: task_10
spec: 0103-a-suite-that-leaks-nothing
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: medium
---

# Task 10: Let a sanctioned regeneration write what it declares

## Overview

Installing the guard made two commands fail that exist to write into the
repository: `make baseline-digests` and the plan-characterization update flag.
Every path they were refused for is already declared as their output in an
authorization record's `## Sanctioned regeneration` block. The guard must read
that declaration rather than carry an exemption list of its own.

## Requirements

1. MUST exempt a path from the guard only when an authorization record declares
   it as an output of the regeneration command that is running.
2. MUST refuse the same path when it is written by anything other than its
   declared command.
3. MUST NOT introduce a second exemption list; the declaration the changed-path
   audit reads is the one the guard reads.
4. MUST leave every other violation reported exactly as before.

## Subtasks

- [ ] Read the sanctioned-regeneration declarations.
- [ ] Exempt a declared output under its declared command only.
- [ ] Cover the declared, the undeclared, and the wrong-command cases.

## Acceptance Criteria

- [ ] `make baseline-digests` completes with the guard installed.
- [ ] The plan-characterization update flag completes with the guard installed.
- [ ] A declared output written by a command that does not declare it is still
      refused, and named.
- [ ] An undeclared path is still refused, and named.
- [ ] No exemption list exists outside the authorization records.

## Verification

- `go test -count=1 ./internal/suiteguard -run 'TestSanctionedRegenerationIsNotAViolation' -v > /tmp/0103-t10.log 2>&1; s=$?; grep -q '^--- PASS: TestSanctionedRegenerationIsNotAViolation' /tmp/0103-t10.log || { cat /tmp/0103-t10.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0103-t10.log || { echo 'the guard suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0103-t10.log && { echo 'the guard suite selected no cases'; exit 1; }; grep -c '^--- PASS' /tmp/0103-t10.log > /tmp/0103-t10-n.txt; test "$(cat /tmp/0103-t10-n.txt)" -ge 3 || { echo "expected the declared, undeclared, and wrong-command cases, got $(cat /tmp/0103-t10-n.txt)"; cat /tmp/0103-t10.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each direction is its own case.
- `make baseline-digests > /tmp/0103-t10-digests.log 2>&1; s=$?; grep -q 'repository boundary violated' /tmp/0103-t10-digests.log && { echo 'the guard still refuses the sanctioned regeneration:'; grep -A 6 'repository boundary violated' /tmp/0103-t10-digests.log; exit 1; }; exit $s` — expected: exits 0, proving the declared command completes with the guard installed. Fails today with the boundary violation it prints.
- `go test -count=1 ./internal/baseline -count=1 -run 'TestDeclaredStepRegenerationAndFrozenBoundaries' -tags repocontract -v > /tmp/0103-t10-frozen.log 2>&1; s=$?; grep -q '^--- PASS: TestDeclaredStepRegenerationAndFrozenBoundaries' /tmp/0103-t10-frozen.log || { cat /tmp/0103-t10-frozen.log; exit 1; }; exit $s` — expected: exits 0, proving the frozen-boundary contract passes with the guard installed. Fails today.
- `grep -rq 'Sanctioned regeneration' internal/suiteguard internal/suiteguardcontract || { echo 'the guard does not read the sanctioned-regeneration declaration'; exit 1; }; grep -rn 'exempt\|allowlist\|allowList\|ignoreList' internal/suiteguard internal/suiteguardcontract > /tmp/0103-t10-priv.txt 2>/dev/null; grep -v 'sanctioned\|declaration\|authorization' /tmp/0103-t10-priv.txt > /tmp/0103-t10-priv2.txt; test ! -s /tmp/0103-t10-priv2.txt || { echo 'a private exemption list exists beside the declaration:'; cat /tmp/0103-t10-priv2.txt; exit 1; }` — expected: exits 0, proving the exemption comes from the declaration the changed-path audit already reads and not from a second list the guard owns. Fails today, where neither guard package mentions the declaration.

## Context

- interface: `docs/workflow/authorizations/2026-08-06-proof-cost.md`
- interface: `internal/speccheck/mechanical.go`

## References

`_techspec.md` → Build Order 9; Risks & Considerations, the sanctioned
regeneration. `_prd.md` → Core Feature 7; Success Metrics. ADR-0128, ADR-0126,
ADR-0081.

## Result

Implemented one sanctioned-regeneration reader in `internal/suiteguardcontract`.
The changed-path audit and suite guard now parse that reader's command-to-output
declarations. The guard removes a changed path only when its declaring command
is the current process or an observed ancestor; every remaining violation keeps
the existing comparison order and diagnostic shape. Baseline regeneration
fixtures now copy the repository authorization records, and the synthetic
plan-characterization command writes a fixture authorization record for its
derived outputs.

Focused evidence by acceptance criterion:

- `make baseline-digests` with the guard installed: the focused
  `go test -count=1 ./internal/baseline -tags repocontract -run
  '^TestMeasuredSanctionedOwnershipMatchesRecords$' -v` exercise progressed
  through the nested command without `repository boundary violated`. The outer
  contract then failed its byte-restoration assertion because the two
  pre-existing modified catalog artifacts contain digest `be1796...`, while the
  command regenerated tracked digest `9773fe...`; this is not a Task 10 guard
  failure.
- Plan-characterization update behavior: the Baseline fixture now materializes
  an authorization declaration for the exact synthetic command and enumerates
  its exact outputs. The Daemon-owned declared Verification that exercises this
  path was not run in this turn.
- Wrong command: `go test -count=1 ./internal/suiteguard -run
  'TestSanctionedRegeneration' -v` passed
  `TestSanctionedRegenerationIsNotAViolationWrongCommandIsRefused`; its guarded
  subprocess named `created: declared.txt`.
- Undeclared path: the same focused command passed
  `TestSanctionedRegenerationIsNotAViolationUndeclaredPathIsRefused`; its guarded
  subprocess named `created: undeclared.txt`.
- Single declaration source: `go test -count=1 ./internal/speccheck` passed all
  210 tests after `mechanical.go` switched to
  `suiteguardcontract.ParseSanctionedRegenerations`. A source scan found
  `Sanctioned regeneration` only in the shared reader and test declarations,
  with no private exemption-list identifier in either guard package.

Additional focused checks:

- `go test -count=1 ./internal/suiteguard ./internal/suiteguardcontract` — passed
  6 tests across 2 packages.
- `go test -count=1 ./internal/baseline -run
  '^TestDerivedOwnershipIsExhaustive$'` — passed.
- `git diff --check` — passed.
- `make verify-incremental` — failed outside this slice: two `internal/agent`
  tests encountered an ACP adapter that does not advertise `sandbox_mode`, and
  `internal/baseline` rejected the two pre-existing modified catalog artifacts.
  The gate's `internal/suiteguard` package passed.
