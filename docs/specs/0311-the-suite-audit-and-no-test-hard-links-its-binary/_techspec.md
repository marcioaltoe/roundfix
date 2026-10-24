---
spec: 0311-the-suite-audit-and-no-test-hard-links-its-binary
prd: _prd.md
created: 2026-10-24
---

# The suite audit, and no test hard-links its binary — Technical Spec

## Executive Summary

The `internal/agent` `TestMain` hard-links the compiled test binary to three
paths under a temporary `node_modules` tree and removes them after `m.Run`.
Every re-executed child of the binary that reaches `TestMain` does the same,
and the fixture child of `TestSpawnedFixtureDeathFailsFast` leaves its links
behind because it ends through `os.Exit`. macOS checks each new process of an
ad-hoc signed binary through `syspolicyd`, which resolves the executable's
path from the process's file; when the name it resolved has just been
removed, it rejects the process and kills it with `SIGKILL`. Under load, a
fake acpx child started while one of those names disappears dies before
`TestMain`, writes no `args.json`, and fails whichever test started it. This
Spec makes the package adapters copies (API Contract 1), keeps the fixture
children from making any, and adds a repository rule that refuses a hard link
in any test (API Contract 2), as ADR-0324 decides.

The trade-off accepted: each `internal/agent` test process writes three
copies of its binary, about 10 MB each, before its first test, instead of
three directory entries.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0324 is the decision both
  implementation Tasks carry out, as `_prd.md` → Project Constraints states.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized; the Makefile and the CI workflows stay unchanged. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0311-the-suite-audit-and-no-test-hard-links-its-binary/_authorization.md`;
  bounded files: `docs/agents/specific-repository.md`.

## Current behavior

Read and characterized on 2026-10-24 at `56dd1a76` (v0.116.0) in scratch
clones; the audit reference holds every command and number.

1. `internal/agent/acpx_runner_test.go` `TestMain` calls
   `provisionPackageAdapterLinks`, which makes a root with
   `os.MkdirTemp("", "roundfix-agent-package-adapters-")` and calls
   `os.Link(os.Args[0], target)` for `ClaudeAdapterPackage`,
   `@example/claude-code-acp` and `@example/claude-agent-acp`; the returned
   function removes the root after `suiteguard.Main` returns.
   `installSymlinkedPackageAdapter` points a symbolic link at one target.
   It is the only `os.Link` call in the repository.
2. `TestSpawnedFixtureDeathFailsFast` (`acpx_integration_test.go`) re-executes
   the binary with `ROUNDFIX_TEST_SPAWNED_FIXTURE_DEATH=wait`, which
   re-executes it again with `=fixture`; both children run `TestMain` in
   full, so both provision the links, and the fixture child ends through
   `os.Exit(23)` inside the test, so its root is never removed. On
   2026-10-24, 300 such roots sat under `$TMPDIR`, pinning 26 deleted test
   binaries, 184 MiB, 20 of them from Runs of 2026-10-23.
3. `runFakeACPXPrompt` reads `args.json` through `readJSONArgs` before it
   returns the `RunPrompt` error, so a child that never ran fails the test
   with "no such file or directory" and hides the error that said why.
4. The two Run failures of 2026-10-23 are the same event. At 18:14:41 and
   23:10:10 BRT `syspolicyd` logged "Unable (errno: 2) to read file at
   <private> for pid: N" and "Terminating process due to Malware rejection:
   N"; the Run logs fail one second later on a missing `args.json`.
5. Reproduction: twelve concurrent runs of the package binary,
   `-test.count=3 -test.parallel 16`, failed in four of eight rounds. A
   marker written as the first statement of `TestMain` was absent for each
   lost child; a stdlib overlay that logged each `SIGKILL` sent through
   `os.Process` showed none aimed at it; the unified log held one
   `syspolicyd` rejection for each. A standalone program lost 48 of 101,438
   self-started children while four goroutines created and removed hard
   links to its executable, and none of 115,022 without them.
6. Prototype of API Contract 1 in a scratch clone: eight rounds of the same
   load, 96 runs at load averages from 143 to 171, all passed, with no
   `syspolicyd` rejection in their windows.

## System Architecture

- `internal/agent/acpx_runner_test.go`: `provisionPackageAdapterLinks` copies
  instead of linking; `TestMain` skips it in the fixture children;
  `runFakeACPXPrompt` reports the `RunPrompt` error when no arguments were
  written.
- `internal/agent/package_adapter_test.go` (new): the copy and leak tests.
- `internal/testfixture/hard_link_rule_test.go` (new): the Hard Link rule.
- `docs/agents/specific-repository.md`, `CONTEXT.md`: the rule and the
  **Single-Name Test Binary** term.

No production package, import edge, exported symbol, schema or command
changes.

## Implementation Design

### Interfaces

```go
// internal/agent/acpx_runner_test.go (test code only)
const packageAdapterMode = 0o755 // named, so the written-executable rule's literal-mode inventory is unchanged

// copyExecutable writes source to a new file at target and closes it.
func copyExecutable(source, target string) error
```

### Data Models

None. No persisted format changes.

### API Contracts

1. API Contract: in a process whose `ROUNDFIX_TEST_SPAWNED_FIXTURE_DEATH`
   is empty, `TestMain` makes the package adapter root and writes each of
   its three targets as a full copy of `os.Args[0]`, closed before `m.Run`
   starts, so no fork can inherit an open write descriptor; it removes the
   root after the run, as today. With the variable set it makes no root, and
   `packageAdapterLinks` stays nil. `runFakeACPXPrompt` fails with a message
   naming the missing `args.json` path and the `RunPrompt` error and type
   when the file does not exist.
2. API Contract: `TestNoTestCreatesAHardLink` parses every `_test.go` file
   under `internal/` and reports `path:line` for each call of `Link` on a
   package name that resolves, through the file's imports, to `os`,
   `syscall` or `golang.org/x/sys/unix`, and of `Linkat` on `syscall` or
   `golang.org/x/sys/unix`; a method or variable of another type named
   `Link`, `Symlink` and `Readlink` are not reported.

### Surface Transcripts

None. No command surface changes.

### Invariants

1. No `internal/agent` package adapter is the same file as the test binary.
2. A re-executed fixture child leaves no package adapter root under its
   temporary directory.
3. No `_test.go` file under `internal/` creates a hard link.
4. Every existing `internal/agent` test keeps its assertions.

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → API Contract 1, Invariants 1
  and 4, task_01, the QA load rows of task_03.
- Goal 2, Core Feature 2, Success Metric 3 → API Contract 2, Invariant 3,
  task_02.
- Goal 3, Success Metric 2 → API Contract 1, Invariant 2, task_01.
- Goal 4, Success Metrics 4 and 5 → the audit reference's targets;
  Unreachable Acceptance; Specs 0312 to 0314.

Deviations from the PRD: None.

## Absence Evidence

1. Copied package adapters — command:
   `grep -n "os.Link(\|copyExecutable" internal/agent/acpx_runner_test.go`;
   result: one `os.Link(os.Args[0], target)` in
   `provisionPackageAdapterLinks`; no copy.
2. A hard-link rule — command:
   `grep -rln "Link(" internal/testfixture`; result: no match.
3. The glossary term — command:
   `grep -n "Single-Name Test Binary" CONTEXT.md docs/agents/specific-repository.md`;
   result: no match.

## ADR Alternatives

- ADR-0324: this Spec implements it. Its rejected options: keeping the hard
  links and removing them only in the top-level process (a removal still
  races a start across concurrent packages and repeated suites), retrying a
  killed fake child (it hides an operating-system kill the test caused and
  breaks the flaky-test rule), and changing the host's security settings
  (the repository cannot own a machine setting).

## Integration Points

None. No external system, network or provider is reached.

## Testing Approach

- `internal/agent/package_adapter_test.go` (new, parallel tests):
  `TestPackageAdaptersAreCopiesOfTheTestBinary` asserts three targets, none
  `os.SameFile` with `os.Args[0]`, each the binary's size, executable on
  non-Windows hosts; the `os.Link` sabotage fails it on `os.SameFile`.
  `TestASpawnedFixtureChildLeavesNoPackageAdapter` runs the fixture child
  with `TMPDIR`, `TMP` and `TEMP` at a fresh temporary directory, requires
  exit status 23 and `fixture started`, and finds no root; removing the
  skip of API Contract 1 makes it find the leaked root.
- `internal/testfixture/hard_link_rule_test.go` (new):
  `TestNoTestCreatesAHardLink`, with a repository subtest and seeded
  subtests in the shape of `TestNoTestWritesAnExecutableOutsideTheResidue`.
  `internal/testfixture` is not a Parallel Test Package, so the test follows
  the shape of its neighbours.
- The load harness is a QA row (task_03): it needs about ten minutes of a
  loaded machine and reads the macOS unified log, which no Task Verification
  may do.

## Build Order

1. Copied package adapters, the fixture skip and the helper message, with
   their two tests (task_01).
2. The Hard Link rule, the guide line and the glossary term (task_02)
   (depends on: 1, because the rule's repository subtest fails while the
   `os.Link` call of task_01's file remains).
3. The QA gate (task_03) (depends on: 1, 2).

No other Spec is required.

## Risks & Considerations

- On Linux, writing an executable while another goroutine forks can make a
  later exec fail with "text file busy" (go.dev/issue/22315); API Contract 1
  writes and closes the copies before `m.Run`, when no test goroutine runs.
- The copies cost about 30 MB of temporary disk per `internal/agent` test
  process while it runs; the fixture children make none.
- `syspolicyd` may reject a process for other reasons (an unreadable file, an
  exhausted daemon); the load row records every rejection in its windows,
  whatever the process, so an unrelated one is visible and not counted as
  this Spec's.
- The audit's other findings (a second CI run of the documentation
  contracts, a 31 GB shared Go build cache, packages whose test results
  never hit the cache, the slowest tests) are not changed here.

## Glossary

- adds: **Single-Name Test Binary**

## Decisions

- Remove the hard links at their cause, copy instead (ADR-0324).
- Skip the copies in the fixture children by the variable that selects
  their mode, the only re-executed children that reach `TestMain`'s
  provisioning.
- Keep the load harness a QA row with a recorded baseline, not a test.
