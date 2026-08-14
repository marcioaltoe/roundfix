---
task: task_01
spec: 0103-a-suite-that-leaks-nothing
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: test
complexity: high
---

# Task 01: Compile the fixtures the suite executes

## Overview

Two places write an executable and run it moments later: the adapter harness in
`internal/agent` and the release-gate formatter in `internal/cli`. Under
concurrent forking the kernel refuses the exec with `text file busy`, or the
fixture answers before its bytes are readable. Both were measured here. This
slice replaces the written scripts with compiled binaries, which removes the
window rather than retrying inside it.

## Requirements

1. MUST replace every fixture that is written and then executed with a compiled
   binary — either a re-execution of the test binary behind an environment
   switch, or a binary built once for the package before its tests spawn.
2. MUST keep every existing assertion about what the fixtures do; a fixture's
   observable behaviour is unchanged.
3. MUST prove the new shape survives concurrent exec load, in a test that fails
   on the old shape and passes on the new one.
4. MUST NOT add a retry, a sleep, or a widened deadline anywhere in the repair.
5. MUST record where the failure mode is documented outside this repository.

## Subtasks

- [ ] Compile the adapter fixtures in the ACPX harness.
- [ ] Compile the release-gate formatter fixture.
- [ ] Add the concurrent-exec stress case.
- [ ] Cite the external documentation of the race.

## Acceptance Criteria

- [ ] No test writes a file and then executes that same file.
- [ ] The adapter harness's existing tests pass unchanged in behaviour.
- [ ] The stress case executes a fixture under concurrent forking and reports
      neither `text file busy` nor an empty version probe.
- [ ] No retry, sleep, or extended deadline was introduced.
- [ ] **Outside evidence.** The failure mode is documented outside this
      repository: Go's own issue tracker carries writing-then-executing from a
      forking process as a known `ETXTBSY` hazard (golang/go#22315). The
      repair cites it, so the fix rests on a diagnosis this Spec did not invent.
      Source: published upstream issue, obtained during authoring.

## Verification

- `! grep -rn 'os.WriteFile' internal/agent/acpx_runner_test.go internal/cli/baseline_release_gate_test.go > /tmp/0103-t01-writes.txt 2>&1; grep -E '0o755|0755' /tmp/0103-t01-writes.txt && { echo 'a fixture is still written with an executable mode:'; grep -E '0o755|0755' /tmp/0103-t01-writes.txt; exit 1; }; exit 0` — expected: exits 0, proving no executable fixture is written into place. It prints the offending lines on failure. Fails today, where both files write mode-0755 fixtures.
- `go test -count=1 ./internal/agent -run 'TestFixtureBinarySurvivesConcurrentExec' -v > /tmp/0103-t01.log 2>&1; s=$?; grep -q '^--- PASS: TestFixtureBinarySurvivesConcurrentExec' /tmp/0103-t01.log || { cat /tmp/0103-t01.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing stress case; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0103-t01.log` — expected: exits 0, refusing a vacuous run.
- `go test -count=1 ./internal/agent -run 'TestACPXRunAppliesSelectionBeforePrompt|TestCheckAdapterProvesOfficialClaudePackageAndVersion' -v > /tmp/0103-t01-regress.log 2>&1; s=$?; test $s -eq 0 || { cat /tmp/0103-t01-regress.log; exit 1; }; for t in TestACPXRunAppliesSelectionBeforePrompt TestCheckAdapterProvesOfficialClaudePackageAndVersion; do grep -q "^--- PASS: $t" /tmp/0103-t01-regress.log || { echo "FAIL: $t no longer passes"; cat /tmp/0103-t01-regress.log; exit 1; }; done; grep -rq 'func FixtureBinary' internal/ || { echo 'the existing assertions pass, but no compiled-fixture constructor exists'; exit 1; }` — expected: exits 0, proving the fixtures still behave as their existing tests require *and* that they became compiled binaries. Both halves are one command because the behaviour half alone passes on an untouched tree.
- `grep -q '22315' docs/specs/0103-a-suite-that-leaks-nothing/task_01.md && grep -rq '22315' internal/agent/ || { echo 'the repair does not cite the upstream diagnosis'; exit 1; }` — expected: exits 0, proving the outside evidence reached the code and not only the Task file. Fails today.

## Context

- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`

## References

`_techspec.md` → Build Order 1; System Architecture, the compiled fixture;
Testing Approach. `_prd.md` → Core Feature 4; Goal 3; Decisions. ADR-0125.
Evidence: `docs/findings/2026-08-10-a-fake-adapter-goes-silent-under-a-dense-start.md`,
both 2026-08-14 addenda.
