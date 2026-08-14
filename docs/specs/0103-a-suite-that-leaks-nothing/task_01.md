---
task: task_01
spec: 0103-a-suite-that-leaks-nothing
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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

## Result

### Implementation

- The ACPX adapter fixtures now hard-link the package's compiled Go test binary
  and read observable behavior from non-executable JSON sidecars. `TestMain`
  dispatches those re-executions before the fake ACPX mode, including bare
  commands resolved through the harness's explicit environment.
- `testfixture.FixtureBinary` compiles Go fixture source through `go build`; the
  release-gate formatter uses it instead of writing a mode-0755 shell script.
- `TestFixtureBinarySurvivesConcurrentExec` runs 16 workers with 32 fixture
  provision-and-exec cycles each and rejects both exec errors and empty version
  probes.
- The adapter fixture comment cites Go issue 22315 as the upstream account of
  the write-then-execute `ETXTBSY` hazard.

### Focused checks

- Pre-change inspection found the adapter and formatter helpers writing shell
  scripts with mode `0o755`; the stress test and `FixtureBinary` constructor did
  not exist.
- `env GOCACHE=/private/tmp/roundfix-0103-task01-gocache go test ./internal/agent -run '^TestFixtureBinarySurvivesConcurrentExec$' -count=1`
  passed; all 512 concurrent cycles returned `fixture-version`.
- `env GOCACHE=/private/tmp/roundfix-0103-task01-gocache go test ./internal/agent -run '^(TestACPXRunAppliesSelectionBeforePrompt|TestCheckAdapterProvesOfficialClaudePackageAndVersion)$' -count=1`
  passed.
- `env GOCACHE=/private/tmp/roundfix-0103-task01-gocache go test ./internal/cli -run '^TestGuidanceCompositionJourney$/^standard-typescript-monorepo$' -count=1`
  passed, exercising the compiled formatter's version probe, check invocation,
  and marker behavior.
- `env GOCACHE=/private/tmp/roundfix-0103-task01-gocache go test ./internal/agent -count=1`
  passed after the fixture change.
- `env GOCACHE=/private/tmp/roundfix-0103-task01-gocache rtk make verify-incremental`
  passed, including formatting, the full Go suite, skill checks, and the build.
- `git diff --check` passed. Added-line inspection found no retry, sleep,
  deadline change, or executable-mode `os.WriteFile` in either named interface.

### Acceptance evidence

1. The executed adapter path is a hard link to the already-compiled test binary;
   the formatter executable is emitted by `go build`. Tests write only source,
   configuration, markers, and other non-executable data.
2. The named adapter regression checks and the full `internal/agent` package
   passed with their existing assertions.
3. The concurrent-exec stress check passed without an exec error or empty probe.
4. Diff inspection found no added retry, sleep, or deadline change.
5. `internal/agent/acpx_runner_test.go` cites Go issue 22315 next to the compiled
   fixture implementation.

The slice introduces no new product-domain term, so `CONTEXT.md` needs no
glossary change. No follow-up work was discovered inside this Task's scope.

The Daemon-owned Verification commands were not run in this Agent turn.
