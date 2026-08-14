---
task: task_11
spec: 0103-a-suite-that-leaks-nothing
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: medium
---

# Task 11: Let the declared command name itself

## Overview

The guard reads the sanctioned-regeneration declaration correctly, and then asks
the operating system which command is running by walking process ancestry with
`ps`. On a host that denies the process table — the same denial that blocked two
Force Stop rows in this Spec's own gate — the ancestry walk returns nothing, no
declaration is selected, and `make baseline-digests` is refused for the two
outputs its record declares. The command should say who it is rather than be
guessed at.

## Requirements

1. MUST identify the running sanctioned command from an explicit declaration made
   in process, not from reading the host process table.
2. MUST keep the exemption bound to the declared command: a declared output
   written outside its command is still a violation.
3. MUST NOT read `ps`, `kern.proc.all`, or any host-wide process enumeration to
   decide an exemption.
4. MUST NOT modify the Makefile or any other protected tooling; the declaration
   is made in Go by the test that performs the regeneration.
5. MUST leave every other guard behaviour unchanged.

## Subtasks

- [ ] Add the in-process declaration and bind it to the exemption.
- [ ] Declare it from the tests that run the sanctioned commands.
- [ ] Remove the process-table ancestry walk.
- [ ] Cover the declared, undeclared, and wrong-command cases.

## Acceptance Criteria

- [ ] `make baseline-digests` completes with the guard installed on a host whose
      process table cannot be read.
- [ ] The plan-characterization update flag completes the same way.
- [ ] A declared output written without the declaration is still refused.
- [ ] No guard code reads a host-wide process table.
- [ ] The Makefile is unchanged.

## Verification

- `! grep -rn "\"ps\"\|kern.proc.all\|pgrep" internal/suiteguard internal/suiteguardcontract > /tmp/0103-t11-ps.txt 2>&1; test ! -s /tmp/0103-t11-ps.txt || { echo 'the guard still reads the host process table:'; cat /tmp/0103-t11-ps.txt; exit 1; }; grep -rq 'Sanctioned regeneration' internal/suiteguard internal/suiteguardcontract || { echo 'the guard no longer reads the declaration'; exit 1; }` — expected: exits 0, proving the ancestry walk is gone and the declaration is still read. Fails today, where `activeCommands` shells out to `ps`.
- `go test -count=1 ./internal/suiteguard -run 'TestSanctionedRegeneration' -v > /tmp/0103-t11.log 2>&1; s=$?; grep -q '^--- PASS: TestSanctionedRegenerationIsDeclaredInProcess' /tmp/0103-t11.log || { cat /tmp/0103-t11.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing case; fails today, where no such test exists.
- `test -s /tmp/0103-t11.log || { echo 'the guard suite produced no output'; exit 1; }; grep -qi 'no tests to run' /tmp/0103-t11.log && { echo 'the guard suite selected no cases'; exit 1; }; grep -c '^--- PASS' /tmp/0103-t11.log > /tmp/0103-t11-n.txt; test "$(cat /tmp/0103-t11-n.txt)" -ge 3 || { echo "expected the declared, undeclared, and wrong-command cases, got $(cat /tmp/0103-t11-n.txt)"; cat /tmp/0103-t11.log; exit 1; }; grep -q 'TestSanctionedRegenerationIsDeclaredInProcess' /tmp/0103-t11.log || { echo 'the in-process case is absent; the count came from the cases task_10 already had'; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each direction is its own case. The count is anchored to the new case's name, because task_10 already left three passing cases in this package.
- `git diff --name-only HEAD -- Makefile > /tmp/0103-t11-make.txt; test ! -s /tmp/0103-t11-make.txt || { echo 'the Makefile was modified without authorization'; cat /tmp/0103-t11-make.txt; exit 1; }; grep -rq 'func Declare\|DeclareSanctioned' internal/suiteguard || { echo 'no in-process declaration exists'; exit 1; }` — expected: exits 0, proving the declaration is Go rather than a tooling edit. Fails today on the second clause, so an untouched Makefile alone cannot pass it.
- `make baseline-digests > /tmp/0103-t11-digests.log 2>&1; s=$?; grep -q 'repository boundary violated' /tmp/0103-t11-digests.log && { echo 'the guard still refuses the sanctioned regeneration:'; grep -A 6 'repository boundary violated' /tmp/0103-t11-digests.log; exit 1; }; test $s -eq 0 || { cat /tmp/0103-t11-digests.log; exit 1; }; grep -rq 'func DeclareSanctionedRegeneration' internal/suiteguard || { echo 'the command completed, but nothing declares itself in process'; exit 1; }` — expected: exits 0, proving the declared command completes through the in-process declaration. The command already completes on a host whose process table is readable, so the declaration is what this measures.

## Context

- interface: `internal/suiteguard/suiteguard.go`
- interface: `internal/suiteguardcontract/regeneration.go`

## References

`_techspec.md` → Build Order 9; Risks & Considerations, the sanctioned
regeneration. `_prd.md` → Core Feature 7. ADR-0128, ADR-0126.
Evidence: this Spec's QA report `qa/qa-report-2026-08-14-02.md`, finding F-02.

## Result

Implemented an in-process sanctioned-regeneration declaration in `suiteguard`.
The guard now matches the declared command against the existing authorization
records and exempts only that record's exact outputs. It no longer inspects its
process ancestry. The guarded Baseline tests declare `make baseline-digests`
before each update mode writes derived artifacts, including the plan and catalog
characterization flags.

Focused-check evidence:

- Red signal: `GOCACHE=/private/tmp/roundfix-task11-go-cache go test
  ./internal/suiteguard -run '^TestSanctionedRegenerationIsDeclaredInProcess$'`
  failed to compile because `suiteguard.DeclareSanctionedRegeneration` did not
  exist.
- `GOCACHE=/private/tmp/roundfix-task11-go-cache go test
  ./internal/suiteguard -run '^TestSanctionedRegenerationIsDeclaredInProcess$'`
  exited 0 after the implementation.
- `GOCACHE=/private/tmp/roundfix-task11-go-cache go test
  ./internal/suiteguard -run
  '^TestSanctionedRegenerationIsNotAViolation(WrongCommandIsRefused|UndeclaredCommandIsRefused)$'`
  exited 0, proving a wrong or absent in-process command does not exempt an
  authorization record's declared output.
- `GOCACHE=/private/tmp/roundfix-task11-go-cache go test
  ./internal/baseline -run
  '^TestBaselinePlanCharacterizationDiffNamesShapeAndField$'` exited 0, compiling
  the Baseline package with all regeneration declaration call sites.
- `GOCACHE=/private/tmp/roundfix-task11-go-cache go test
  ./internal/suiteguard` exited 0.
- `rg -n '"ps"|kern\.proc\.all|pgrep' internal/suiteguard
  internal/suiteguardcontract` found no matches (exit 1), so the guard no longer
  contains a host-wide process-table lookup.
- `git diff --quiet HEAD -- Makefile` exited 0; the Makefile is unchanged.
- `git diff --check` exited 0.

Acceptance evidence:

- Baseline digest regeneration: every guarded `-update` test invoked by
  `BASELINE_DIGEST_STEPS` now calls the shared in-process declaration. The
  Daemon-owned `make baseline-digests` Verification remains pending.
- Plan-characterization update: its update branch calls the same declaration,
  and the Baseline package compiled in the focused check. The Daemon owns the
  mutating acceptance run.
- Missing declaration: the focused undeclared-command subprocess test exited 0
  while asserting that `declared.txt` was refused and named as a violation.
- Process-table independence: the ancestry walk and its `ps` execution were
  removed; the focused source scan found no forbidden process-table token.
- Protected tooling: the focused Git inspection found no Makefile change.

The commands under `## Verification` were not run; the Roundfix Daemon owns
those commands and Task settlement.
