---
task: task_02
spec: 0254-a-faster-make-test
status: pending
type: test
complexity: high
---

# Task 02: Six more packages run their tests in parallel

## Overview

After `internal/cli`, the sequential chains that hold the suite are in
`internal/daemon` (144.8 s on 2026-10-08), `internal/speccheck` (138.4 s),
`internal/spec` (119.0 s), `internal/worktree` (46.2 s), `internal/baseline`
(21.9 s) and `internal/store` (`_techspec.md` → Current behavior). This Task
lists the six as Parallel Test Packages with the ceilings of API Contract 2
and converts them the way task_01 converted `internal/cli` (ADR-0259). It is
verifiable on its own: the rule test passes with all seven packages within
their ceilings, and each package passes under the race detector and a
shuffled order.

This is an authorized tooling Task for three Governed Paths:
`internal/baseline/derived_ownership_test.go`,
`internal/spec/coverage_test.go` and `internal/speccheck/mechanical_test.go`.
It may change only the paths in its Context, the other test files of the six
packages that receive an inserted `t.Parallel()` or a `// Sequential:`
comment, and this Task file. The Daemon records those other files under
`## Recorded paths`.

## Requirements

1. MUST add `internal/daemon` 16, `internal/baseline` 12, `internal/store`
   3, `internal/spec` 2, `internal/speccheck` 2 and `internal/worktree` 2 to
   the ceiling map, and keep `internal/cli` at 48.
2. MUST convert every top-level test of the six packages under task_01's
   Requirements 3 to 6, including any test Spec 0253 added under
   `internal/baseline` before this Task starts. Files under a contract build
   constraint stay as they are.
3. MUST keep sequential, with a `// Sequential: <reason>`, at least these
   tests the authoring prototype found: in `internal/daemon`,
   `TestGitCommitterStagesSelectedTrackedPathMatchedByGlobalIgnore`,
   `TestRunGitForTestIgnoresForcedSigningConfig`,
   `TestQAFormatStartFailureRestoresFiles`,
   `TestPriorQAPassImportsTheNewestUnintegratedReportByteForByte`,
   `TestTaskCycleSettlesAQualifyingPartial` and
   `TestProbeParserUnavailableRunsVerifier`, which set the environment, and
   `TestUnnamedTaskStaysBlockedByRedPrecondition`,
   `TestTaskCycleRepositoryGatePreconditionFailureStartsNoAgentSession` and
   the `TestPreWorkProbe…` tests, whose shared helpers call `t.Parallel()`
   themselves, unless that call moves from the helper into each test that
   calls it, which makes them parallel; in `internal/spec`,
   `TestPreconditionRefusalReportNamesItsAuditor`, which swaps `app.Version`,
   `app.BuildCommit` and `app.BuildTime`; in `internal/store`, the two
   `TestOwnerProcessIdentity…` tests; and in `internal/baseline`, the four
   `skills_lock_read_test.go` tests that reach `t.Setenv`.
4. MUST NOT convert a test that calls `captureStderr` in `internal/worktree`,
   which swaps `os.Stderr`, and MUST NOT change that helper.
5. MUST NOT change any assertion, expected value, golden, fixture, helper
   behavior or production file (Invariant 1), and every package keeps
   `suiteguard.Main` (ADR-0126).
6. MUST prove the ceiling can fail. The Result MUST record lowering one
   ceiling below its package's count, the failing line the rule printed, and
   the restored value.

## Subtasks

- [ ] Add the six packages and their ceilings to the rule.
- [ ] Convert each package and mark its Sequential Tests.
- [ ] Search each package for swapped variables.
- [ ] Run the six packages under `-race` and `-shuffle=on`.
- [ ] Record the ceiling sabotage.

## Acceptance Criteria

- [ ] The rule test logs all seven packages within their ceilings.
- [ ] The six packages pass with `-race -short` and with `-shuffle=on`.
- [ ] No assertion or production file changed.

## Context

- instruction: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- instruction: `docs/adr/0089-code-under-test-takes-its-environment-explicitly.md`
- instruction: `docs/adr/0126-the-suite-proves-its-own-isolation-per-package.md`
- creates: `internal/testfixture/parallel_tests_test.go`
- interface: `internal/baseline/derived_ownership_test.go`
- interface: `internal/spec/coverage_test.go`
- interface: `internal/speccheck/mechanical_test.go`
- interface: `internal/daemon/daemon_test.go`
- interface: `internal/daemon/qa_format_test.go`
- interface: `internal/daemon/qa_prior_pass_test.go`
- interface: `internal/daemon/qa_prior_pass_skip_test.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/verification_probe_malformed_test.go`
- interface: `internal/spec/qa_test.go`
- interface: `internal/spec/qa_frontmatter_test.go`
- interface: `internal/speccheck/evidence_digest_test.go`
- interface: `internal/store/process_unix_test.go`
- interface: `internal/baseline/skills_lock_read_test.go`
- interface: `internal/worktree/worktree_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^TestEveryTestInAParallelTestPackageRunsInParallel$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestEveryTestInAParallelTestPackageRunsInParallel ' || { printf 'the Parallel Test Package rule did not run\n' >&2; exit 1; }; for entry in 'internal/daemon:16' 'internal/baseline:12' 'internal/store:3' 'internal/spec:2' 'internal/speccheck:2' 'internal/worktree:2'; do package="${entry%%:*}"; ceiling="${entry##*:}"; printf '%s\n' "$out" | grep -Eq "$package: [0-9]+ Sequential Tests, ceiling $ceiling" || { printf '%s is not a Parallel Test Package with ceiling %s\n' "$package" "$ceiling" >&2; exit 1; }; done` — expected: exit 0; before this Task the rule lists only `internal/cli`, so the first package check fails.
- `grep -q '"internal/worktree": *2,' internal/testfixture/parallel_tests_test.go || { printf 'internal/worktree is not listed\n' >&2; exit 1; }; go test -count=1 -race -short -timeout 30m ./internal/daemon ./internal/spec ./internal/speccheck ./internal/worktree ./internal/store ./internal/baseline || exit 1; go test -count=1 -shuffle=on -timeout 20m ./internal/daemon ./internal/spec ./internal/speccheck ./internal/worktree ./internal/store ./internal/baseline || exit 1` — expected: exit 0; before this Task `internal/worktree` is not listed, so the command fails first. After it the six packages pass once under the race detector and once in a shuffled order.

## References

- [_prd.md](_prd.md) — Goals 1–2; User Stories 1, 2 and 4; Core Feature 3; Success Metrics 3–4
- [_techspec.md](_techspec.md) — Current behavior; API Contract 1; API Contract 2; Invariants 1–4; Testing Approach; Risks & Considerations; Build Order 2
- ADR-0259; ADR-0089; ADR-0126; ADR-0213
