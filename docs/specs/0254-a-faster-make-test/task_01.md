---
task: task_01
spec: 0254-a-faster-make-test
status: completed
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

## Result

Implemented the Parallel Test Package rule and the `internal/cli` conversion.
Task status and the declared Verification remain Daemon-owned; no declared
Verification, repository Verification, commit, push or pull request was run.

### Implementation

- The AST rule lists only `internal/cli`, with ceiling 48. It checks each
  test's own `*testing.T` parameter and first statement, accepts leading or doc
  Sequential reasons with at least two words, skips contract build constraints,
  reports every violation with its file, line and test name, and logs the count.
- The package has 1,409 top-level tests: 1,366 now have first-statement
  `t.Parallel()`, and 43 have valid Sequential reasons. Added 549 parallel
  calls and 29 Sequential comments across 106 CLI test files. All 14 existing
  Sequential comments remain byte-identical.
- `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded` already called
  `t.Parallel()` after its platform guard. Moved that existing call first
  rather than leaving a duplicate. The existing single-line
  `TestArchiveRefusesASpecWithAGlossaryGap` declaration was expanded only to
  format its inserted parallel call. Assertions, fixture code and helper
  behavior are preserved; neither production code nor another Task changed.
- Added both glossary terms beside Repository Contract Test, with `_Avoid_`
  lines and ADR-0259 references. The guide differs by exactly the requested
  repository-rule bullet, outside setup markers.

### State audit

An AST helper-call graph and independent read-only audits traced direct and
indirect environment, working-directory, signal and variable changes.

- Six repository-profile tests stay sequential through
  `newBaselineRepositoryProfileFixture`, which calls `t.Chdir`.
- Environment callers stay sequential, including the carry-forward hook,
  delivery-item binary, archive advice, script fixture and forge probe tests.
- All six owner-staleness tests stay sequential through `setOwnerBuildCommit`,
  which swaps and restores `app.BuildCommit`. Searches covered selector
  assignments for all imported packages, not only `app`, `config` and `store`,
  and assignments to this package's globals. No other test-time global swap
  was found: the notifier/version hooks run in `init`, the script binary path
  is set in `TestMain`, command override maps use `sync.Map`, and the shared
  cold binary cache is initialized under `sync.Once`.
- `TestCLIForceStopOwnerProcessHelper` keeps its SIGTERM reason.
  `TestSupersedeAcceptsAnArchiveRecord` stays sequential because it invokes
  the already parallel `TestSupersedeAcceptsASupersessionArchivedDeliverer`
  with the same `testing.T`. Existing detached-child reasons are unchanged.
  No current shared helper owns a `t.Parallel()` call that needs moving.

### Focused checks

These checks used `GOCACHE=/private/tmp/roundfix-0254-go-cache` after a host
cache access was denied. Exact Go output was retained through `rtk proxy`.

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Rule passes and CLI stays within ceiling | `go test ./internal/testfixture -run 'TestEveryTestInAParallelTestPackageRunsInParallel/repository$' -count=1 -v` exited 0 and logged `internal/cli: 43 Sequential Tests, ceiling 48`. Before conversion the same focused scan failed with 579 missing-form violations. |
| Seeded violations name their sites | `go test ./internal/testfixture -run 'TestEveryTestInAParallelTestPackageRunsInParallel/' -race -count=1 -v` exited 0 (`1.868s`), including the repository scan and 16 seeded subtests. Negative cases assert file, line, test and cause; cases cover both comment positions, late calls, both forms, one-word reasons, both contract tags, legacy tags, ceilings, multiple violations, wrong receiver, late comments, helpers/TestMain and unlisted packages. |
| CLI race and shuffled execution | Focused selector of 593 converted/newly annotated tests and existing stateful neighbors passed with `go test ./internal/cli -run "$selector" -race -short -count=1 -timeout 10m` (`102.215s`) and `go test ./internal/cli -run "$selector" -shuffle=on -count=1 -timeout 10m` (`64.028s`). The selector was retained in `/private/tmp/roundfix-cli-converted-selector.txt`; outputs in `/private/tmp/roundfix-cli-converted-race.log` and `/private/tmp/roundfix-cli-converted-shuffle.log`. Full-package acceptance remains for the Daemon's declared Verification. |
| Glossary and guide state the rule | A focused Python check confirmed both terms have `_Avoid_` lines and ADR-0259 links, and removing the exact one-bullet addition restores the guide's HEAD bytes. |

Additional checks: `gofmt -l internal/testfixture/parallel_tests_test.go
internal/cli` printed no paths; `git -c core.fsmonitor=false diff --check`
exited 0. A source comparison removing only parallel annotations and new
Sequential comments, accounting for the single-line declaration's formatting,
matched the original formatted code across all 106 changed CLI files.
Changed-file postflight covered 110 paths, all within the Task's bound;
the Daemon records the additional CLI test paths under Recorded paths.

### Sabotage and restoration

Removed the first `t.Parallel()` from
`TestArchiveRefusesASpecAFileStillPins` without adding a reason. The focused
repository rule exited 1 and printed:

```text
internal/cli/archive_active_spec_path_test.go:31: TestArchiveRefusesASpecAFileStillPins: missing first-statement Parallel() or Sequential: reason
```

Restored the file byte-for-byte in a `finally` block. The subsequent repository
scan and race-detector run of the rule both passed; the sabotage output is in
`/private/tmp/roundfix-parallel-sabotage.log`.

### Diagnostics resolved during implementation

The initial focused CLI race run exposed a duplicate `t.Parallel()` call in
`TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded`. Its existing late call
was moved first, and the same focused run then passed (`7.143s`), followed by
the 593-test race and shuffled checks above. The panic occurred before this
test's `cmd.Start()`, so it had started no fixture PID to reap. No broad process
kill was used.

### Remaining Daemon checks

The two authored Verification commands, including full-package `-race -short`
and `-shuffle=on`, were not run in this turn. No terminal Task verdict is
claimed. No follow-up implementation was added to this diff.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/archive_active_spec_path_test.go`
- `internal/cli/archive_glossary_test.go`
- `internal/cli/archive_links_test.go`
- `internal/cli/archive_plan_test.go`
- `internal/cli/archive_record_readers_test.go`
- `internal/cli/archive_test.go`
- `internal/cli/baseline_branch_prefix_test.go`
- `internal/cli/baseline_history_citation_test.go`
- `internal/cli/baseline_update_history_test.go`
- `internal/cli/baseline_update_outdated_skills_test.go`
- `internal/cli/baseline_update_retired_skills_test.go`
- `internal/cli/baseline_update_test.go`
- `internal/cli/built_in_review_provider_test.go`
- `internal/cli/carryforward_completed_target_test.go`
- `internal/cli/carryforward_integration_order_test.go`
- `internal/cli/carryforward_test.go`
- `internal/cli/deliver_archive_links_test.go`
- `internal/cli/deliver_archived_retry_test.go`
- `internal/cli/deliver_conflict_test.go`
- `internal/cli/deliver_derived_lines_test.go`
- `internal/cli/deliver_derived_skill_layout_test.go`
- `internal/cli/deliver_item_branch_test.go`
- `internal/cli/deliver_limits_test.go`
- `internal/cli/deliver_merged_outside_test.go`
- `internal/cli/deliver_operator_archive_test.go`
- `internal/cli/deliver_park_status_test.go`
- `internal/cli/deliver_prerequisite_test.go`
- `internal/cli/deliver_publication_scope_test.go`
- `internal/cli/deliver_recovery_test.go`
- `internal/cli/deliver_retry_amendment_test.go`
- `internal/cli/deliver_retry_runs_test.go`
- `internal/cli/deliver_revalidate_test.go`
- `internal/cli/deliver_review_correction_test.go`
- `internal/cli/deliver_runtime_infrastructure_test.go`
- `internal/cli/deliver_start_readiness_test.go`
- `internal/cli/deliver_test.go`
- `internal/cli/deliver_token_ceiling_test.go`
- `internal/cli/detach_test.go`
- `internal/cli/doctor_characterization_test.go`
- `internal/cli/doctor_storage_test.go`
- `internal/cli/doctor_test.go`
- `internal/cli/doctor_trailing_skills_test.go`
- `internal/cli/events_usage_test.go`
- `internal/cli/gc_run_retention_test.go`
- `internal/cli/gc_sanitize_pre_key_root_test.go`
- `internal/cli/history_refusal_test.go`
- `internal/cli/history_test.go`
- `internal/cli/implement_budget_renewal_test.go`
- `internal/cli/implement_qa_format_test.go`
- `internal/cli/implement_test.go`
- `internal/cli/qa_partial_policy_test.go`
- `internal/cli/qa_report_test.go`
- `internal/cli/readiness_toolchain_test.go`
- `internal/cli/reconcile_legacy_key_test.go`
- `internal/cli/reconcile_staging_test.go`
- `internal/cli/reconcile_test.go`
- `internal/cli/releaseplan_checks_test.go`
- `internal/cli/releaseplan_skill_coverage_test.go`
- `internal/cli/reopen_late_dependency_test.go`
- `internal/cli/review_archived_spec_test.go`
- `internal/cli/review_convention_validator_test.go`
- `internal/cli/review_dispose_lock_test.go`
- `internal/cli/review_dispose_revision_test.go`
- `internal/cli/review_disposition_test.go`
- `internal/cli/review_final_message_test.go`
- `internal/cli/review_head_bound_test.go`
- `internal/cli/review_lineage_selection_test.go`
- `internal/cli/review_lineage_test.go`
- `internal/cli/review_merge_base_test.go`
- `internal/cli/review_override_convention_test.go`
- `internal/cli/review_permission_test.go`
- `internal/cli/review_prompt_bound_test.go`
- `internal/cli/review_provider_scope_test.go`
- `internal/cli/review_record_checkout_test.go`
- `internal/cli/review_scope_test.go`
- `internal/cli/review_selection_retry_test.go`
- `internal/cli/review_session_test.go`
- `internal/cli/review_test.go`
- `internal/cli/review_validation_test.go`
- `internal/cli/run_retention_start_test.go`
- `internal/cli/runs_list_vanished_checkout_test.go`
- `internal/cli/runs_show_test.go`
- `internal/cli/settle_test.go`
- `internal/cli/setup_readiness_test.go`
- `internal/cli/spec_audit_record_test.go`
- `internal/cli/spec_check_provenance_test.go`
- `internal/cli/spec_check_test.go`
- `internal/cli/spec_judge_tier_test.go`
- `internal/cli/supersede_test.go`
- `internal/cli/this_repository_skill_set_test.go`
- `internal/cli/version_freshness_isolation_test.go`
- `internal/cli/window_test.go`

## Carry-forward provenance

- Source Run: `run_20261008T215927Z_c26d3579de5392c8`
- Source commit: `7ee2c53ddd783583fa7528e6808950ab7ec67d45`
