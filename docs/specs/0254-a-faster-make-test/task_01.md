---
task: task_01
spec: 0254-a-faster-make-test
status: pending
type: test
complexity: high
---

# Task 01: internal/cli runs its tests in parallel under the Parallel Test Package rule

## Overview

`internal/cli` runs 592 of its 1,409 top-level tests one at a time, and their
summed durations, 279 s to 324 s on 2026-10-08, are most of the package's wall
time (`_techspec.md` → Current behavior). This Task adds the repository test
that defines a Parallel Test Package (ADR-0259), lists `internal/cli` as the
first one, and converts the package: every top-level test calls `t.Parallel()`
first, or it is a Sequential Test that names why. It is verifiable on its own:
the rule test passes on the repository and fails on seeded violations, and the
package passes under the race detector and a shuffled order.

This is an authorized tooling Task for two Governed Paths,
`internal/cli/cli_test.go` and `docs/agents/specific-repository.md`. It may
change only the paths in its Context, the
other `internal/cli/*_test.go` files that receive an inserted `t.Parallel()`
or a `// Sequential:` comment, and this Task file. The Daemon records those
other test files under `## Recorded paths`, because they exceed the Context
bound.

## Requirements

1. MUST create `internal/testfixture/parallel_tests_test.go` with
   `TestEveryTestInAParallelTestPackageRunsInParallel` and the ceiling map of
   `_techspec.md` → Interfaces, holding `"internal/cli": 48`. The rule MUST be
   API Contract 1 exactly: it skips files under a `docscontract` or
   `repocontract` build constraint, treats a test as parallel only when its
   first statement is `<param>.Parallel()` for its own `*testing.T`
   parameter, accepts a `Sequential:` reason of at least two words in the doc
   comment or between the opening brace and the first statement, and reports
   every violation as `<file>:<line>: <Test>` with its cause. It MUST log one
   line per package in the form `<package>: <n> Sequential Tests, ceiling <m>`.
2. MUST test the rule on seeded trees under `t.TempDir()`, each as its own
   subtest: a parallel test passes; a test with neither form fails naming its
   `file:line`; a one-word reason fails; a test with both forms fails; a
   reason in the doc comment and one inside the body both pass; `t.Parallel()`
   after the first statement fails; a `repocontract` file is skipped; a
   package over its ceiling fails naming its count; an unlisted package is
   ignored. The rule test itself MUST call `t.Parallel()` first.
3. MUST insert `t.Parallel()` as the first statement of every top-level test
   in `internal/cli` that does not change process-wide state, directly or
   through a helper it calls. The authoring prototype's method is in
   `_techspec.md` → Current behavior.
4. MUST keep sequential, with a `// Sequential: <reason>` naming the state,
   every test that calls `t.Setenv` or `t.Chdir` directly or through a helper,
   swaps a package-level variable (for example `setOwnerBuildCommit`, which
   writes `app.BuildCommit`), installs a signal handler, or calls another test
   function that is already parallel (`TestSupersedeAcceptsAnArchiveRecord`).
   Existing `// Sequential:` comments stay byte-identical.
5. MUST search the package's test files for assignments to another package's
   variables (`app.`, `config.`, `store.` and the like) and to this package's
   package-level variables, and keep every test that makes one sequential.
6. MUST NOT change any assertion, expected value, golden, fixture, helper
   behavior or production file. A converted test differs from its starting
   bytes only by the inserted line, or by a `t.Parallel()` call moved from a
   shared helper into each test that calls it (Invariant 1).
7. MUST add **Sequential Test** and **Parallel Test Package** to `CONTEXT.md`
   through `domain-modeling`, beside **Repository Contract Test**, with the
   definitions of `_techspec.md` → API Contract 1 and ADR-0259, each with an
   `_Avoid_` line.
8. MUST add one bullet to the repository rules of
   `docs/agents/specific-repository.md`, outside every setup marker, that
   reads: "A top-level test in a Parallel Test Package calls `t.Parallel()`
   as its first statement or carries a `// Sequential: <reason>` comment
   naming the process-wide state it changes; the Parallel Test Package rule in
   `internal/testfixture` enforces it." No other line of that file changes.
9. MUST prove the rule can fail. The Result MUST record one sabotage, a
   converted test whose `t.Parallel()` is removed without a reason, with the
   failing line the rule printed, and that the line was restored.
10. If a converted test panics and kills the test binary, MUST reap any
   fixture process it left by the PID it started, never with a broad `pkill`
   (`_techspec.md` → Risks & Considerations).

## Subtasks

- [ ] Write the rule test with its seeded cases.
- [ ] Convert `internal/cli` and mark the Sequential Tests.
- [ ] Search for swapped variables and keep those tests sequential.
- [ ] Run the package under `-race` and `-shuffle=on`.
- [ ] Add the two glossary terms and the repository rule, and record the sabotage.

## Acceptance Criteria

- [ ] The rule test passes on the repository and logs
      `internal/cli: <n> Sequential Tests, ceiling 48` with `n` at most 48.
- [ ] Each seeded violation fails with its `file:line`.
- [ ] `internal/cli` passes with `-race -short` and with `-shuffle=on`.
- [ ] `CONTEXT.md` defines **Sequential Test** and **Parallel Test Package**,
      and `docs/agents/specific-repository.md` states the rule.

## Context

- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `docs/adr/0089-code-under-test-takes-its-environment-explicitly.md`
- instruction: `docs/adr/0125-a-spawned-fixture-is-a-compiled-binary-never-a-written-script.md`
- instruction: `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`
- instruction: `internal/testfixture/written_executable_test.go`
- creates: `internal/testfixture/parallel_tests_test.go`
- interface: `CONTEXT.md`
- interface: `docs/agents/specific-repository.md`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/archive_output_test.go`
- interface: `internal/cli/archive_record_test.go`
- interface: `internal/cli/baseline_profile_test.go`
- interface: `internal/cli/baseline_update_repository_profile_test.go`
- interface: `internal/cli/carryforward_hooks_test.go`
- interface: `internal/cli/deliver_archive_record_test.go`
- interface: `internal/cli/deliver_item_binary_test.go`
- interface: `internal/cli/deliver_owner_staleness_test.go`
- interface: `internal/cli/deliver_retry_test.go`
- interface: `internal/cli/gc_test.go`
- interface: `internal/cli/implement_detach_teardown_test.go`
- interface: `internal/cli/orphan_unix_test.go`
- interface: `internal/cli/readiness_forge_test.go`
- interface: `internal/cli/script_fixture_test.go`
- interface: `internal/cli/spec_judge_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestEveryTestInAParallelTestPackageRunsInParallel$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestEveryTestInAParallelTestPackageRunsInParallel ' || { printf 'the Parallel Test Package rule did not run\n' >&2; exit 1; }; printf '%s\n' "$out" | grep -Eq 'internal/cli: ([0-9]|[1-3][0-9]|4[0-8]) Sequential Tests, ceiling 48' || { printf 'internal/cli is not within its ceiling of 48\n' >&2; exit 1; }; for term in "**Sequential Test**:" "**Parallel Test Package**:"; do tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "$term" || { printf 'CONTEXT.md lacks %s\n' "$term" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/agents/specific-repository.md | grep -qF -- "A top-level test in a Parallel Test Package calls" || { printf 'docs/agents/specific-repository.md lacks the rule\n' >&2; exit 1; }` — expected: exit 0; before this Task the rule test does not exist, so no PASS line is printed and the command fails.
- `test -f internal/testfixture/parallel_tests_test.go || { printf 'the rule test is missing\n' >&2; exit 1; }; go test -count=1 -race -short -timeout 30m ./internal/cli || exit 1; go test -count=1 -shuffle=on -timeout 20m ./internal/cli || exit 1` — expected: exit 0; before this Task the rule test file is missing, so the command fails first. After it the converted package passes once under the race detector and once in a shuffled order.

## References

- [_prd.md](_prd.md) — Goals 1–2; User Stories 3–4; Core Features 1, 2 and 6; Success Metrics 3–4
- [_techspec.md](_techspec.md) — Current behavior; Interfaces; API Contract 1; API Contract 2; Invariants 1–4; Testing Approach; Risks & Considerations; Build Order 1
- ADR-0259; ADR-0089; ADR-0125; ADR-0126; ADR-0213; ADR-0244
