---
task: task_02
spec: 0254-a-faster-make-test
status: completed
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

## Result

Added the six package ceilings while keeping `internal/cli` at 48. Audited
all six packages, including the current baseline tests, for direct and
helper-mediated environment changes, directory changes, package variable
swaps, stream swaps and existing parallel calls. Eligible top-level tests
now call `t.Parallel()` first; the exceptions carry reasons. Contract-tagged
files and existing Sequential comments are unchanged.

The daemon's eight wrappers keep their helper-owned parallel calls and name
those helpers rather than calling Parallel twice. Baseline's citation-hook
swap and module-record writer remain sequential alongside the four PATH
rewrite tests and existing regeneration exceptions. The worktree stderr
capture test was already parallel: its existing call was replaced with the
required Sequential reason; `captureStderr` is unchanged. The spec metadata
swap and both store environment tests remain sequential.

### Focused evidence by acceptance criterion

All Go checks used `GOCACHE=/tmp/roundfix-task02-gocache`.

1. Rule coverage: `go test -count=1 -v -run
   '^TestEveryTestInAParallelTestPackageRunsInParallel$/repository$'
   ./internal/testfixture` exited 0 after the last implementation edit. This
   is the focused repository subtest, not the declared Verification command.

   | Package | Sequential Tests | Ceiling |
   | --- | ---: | ---: |
   | internal/cli | 43 | 48 |
   | internal/daemon | 15 | 16 |
   | internal/baseline | 12 | 12 |
   | internal/store | 2 | 3 |
   | internal/spec | 1 | 2 |
   | internal/speccheck | 0 | 2 |
   | internal/worktree | 1 | 2 |

2. Race and shuffle implementation checks: a temporary Python launcher
   compared each edited test function with HEAD, built an anchored alternation
   of the 474 distinct edited top-level test names, and passed that selector
   to `go test -count=1 -race -short -shuffle=on -timeout 15m -run <selector>`
   for the six packages. Baseline (21.669 s), store (15.048 s), speccheck
   (22.162 s), and worktree (20.870 s) passed. The first daemon/spec check
   exposed two issues described below. After their repair,
   `python3 /tmp/roundfix-task02-focused.py daemon spec` reran the same edited
   name selector with `-race -short -shuffle=1791498359085405000 -timeout 15m`:
   daemon (20.321 s) and spec (16.854 s) passed, exit 0. The launcher and exact
   expanded command are `/tmp/roundfix-task02-focused.py` and
   `/tmp/roundfix-task02-command.json`; logs are
   `/tmp/roundfix-task02-focused-first.log` and
   `/tmp/roundfix-task02-focused.log`. These are focused subsets; the complete
   six-package race/short and shuffle checks remain Daemon-owned and unrun
   in this Agent turn.

3. Assertion and scope preservation: comparison against HEAD, removing only
   Parallel-call lines and Sequential-comment lines from both versions,
   produced identical source for every edited test file in the six packages.
   No assertion, expected value, golden, fixture, helper behavior or production
   file changed. All six `suiteguard.Main` installations remain unchanged.
   `git -c core.fsmonitor=false diff --check` exited 0. Changed-path review
   found only this Task file, the ceiling-map file, and test files of the six
   allowed packages. `GOOS=linux go build ./internal/store` exited 0 for the
   touched Linux-constrained test surface.

### Focused-check repairs and follow-up

- `TestCollisionsTreatsAnAbsentRepositoryAsAbsentEvidence` already called
  Parallel after a leading comment. The first insertion duplicated that call
  and produced `testing: t.Parallel called multiple times`. Moved the existing
  call ahead of the comment; a sweep found no other duplicate top-level calls.
- `TestBudgetExceededRunRecordsItsReason` failed under parallel load with
  `completed before budget = 0 with 0 commits, want 1 each`. Its existing
  500 ms real-time budget includes Git work, so process-wide scheduling
  contention consumed the budget before the first commit. The focused
  isolated command `go test -count=1 -race -short -run
  '^TestBudgetExceededRunRecordsItsReason$' ./internal/daemon` passed
  (2.451 s). Restored sequential execution and documented that dependency;
  its budget and assertions are unchanged. The subsequent focused shuffled
  daemon run passed. A future Task can remove this wall-clock dependency
  using the engine's controllable clock; that fixture change is outside this
  Task's slice.

### Ceiling sabotage

Temporarily lowered `internal/store` from 3 to 1, below its count of 2, and
ran the focused repository rule subtest. It exited 1 and printed:

```text
parallel_tests_test.go:225: internal/store/process_unix_test.go:603: TestOwnerProcessIdentityDoesNotSpawnPS: internal/store has 2 Sequential Tests, ceiling 1
```

Restored the value to 3 and confirmed the focused rule subtest exited 0.
The failure log is `/tmp/roundfix-task02-sabotage.log`; final counts are in
`/tmp/roundfix-task02-rule-final.log`.

Task status, authored Verification, the Task Graph and all other Task files
remain Daemon-owned and unchanged by this Agent. No commit, push or pull
request was made.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/baseline/adopter_neutral_clauses_test.go`
- `internal/baseline/archive_record_clause_test.go`
- `internal/baseline/assets_sync_owned_membership_test.go`
- `internal/baseline/assets_sync_retired_test.go`
- `internal/baseline/authoring_evidence_clauses_test.go`
- `internal/baseline/authorization_clauses_test.go`
- `internal/baseline/branch_prefix_optional_test.go`
- `internal/baseline/clause_characterization_test.go`
- `internal/baseline/frontend_layout_decision_test.go`
- `internal/baseline/glossary_clauses_test.go`
- `internal/baseline/grouped_sources_clauses_test.go`
- `internal/baseline/history_citation_plan_test.go`
- `internal/baseline/history_citations_root_test.go`
- `internal/baseline/history_citations_test.go`
- `internal/baseline/history_sanitize_clause_test.go`
- `internal/baseline/loop_and_go_clauses_test.go`
- `internal/baseline/loop_clause_test.go`
- `internal/baseline/module_versions_test.go`
- `internal/baseline/network_denied_row_clause_test.go`
- `internal/baseline/qa_override_clause_test.go`
- `internal/baseline/reduced_retirement_clauses_test.go`
- `internal/baseline/release_clause_test.go`
- `internal/baseline/resolve_executable_test.go`
- `internal/baseline/retired_skills_installed_test.go`
- `internal/baseline/roundfix_skill_membership_test.go`
- `internal/baseline/runtime_and_verification_decisions_test.go`
- `internal/baseline/skills_restore_git_test.go`
- `internal/baseline/skills_snapshot_profile_test.go`
- `internal/baseline/skills_trailing_test.go`
- `internal/baseline/verification_tier_clauses_test.go`
- `internal/daemon/lost_rollout_recovery_test.go`
- `internal/daemon/qa_budget_renewal_test.go`
- `internal/daemon/qa_every_run_audit_test.go`
- `internal/daemon/qa_evidence_snapshot_test.go`
- `internal/daemon/qa_frontmatter_test.go`
- `internal/daemon/qa_row_carry_event_test.go`
- `internal/daemon/qa_two_pass_carry_test.go`
- `internal/daemon/settlement_active_spec_path_test.go`
- `internal/daemon/settlement_checks_test.go`
- `internal/daemon/task_budget_renewal_test.go`
- `internal/daemon/task_deletes_test.go`
- `internal/daemon/token_usage_test.go`
- `internal/spec/archive_links_test.go`
- `internal/spec/archive_record_test.go`
- `internal/spec/authorization_test.go`
- `internal/spec/cause_graph_test.go`
- `internal/spec/collision_test.go`
- `internal/spec/context_deletes_test.go`
- `internal/spec/gate_test.go`
- `internal/spec/history_sanitize_test.go`
- `internal/spec/history_sanitize_unproven_test.go`
- `internal/spec/requires_test.go`
- `internal/spec/retirement_deferred_test.go`
- `internal/spec/task_complexity_test.go`
- `internal/spec/task_test.go`
- `internal/speccheck/active_spec_paths_test.go`
- `internal/speccheck/archive_license_git_test.go`
- `internal/speccheck/archive_record_test.go`
- `internal/speccheck/authoring_horizon_test.go`
- `internal/speccheck/backlog_deferred_test.go`
- `internal/speccheck/checkout_test.go`
- `internal/speccheck/evidence_record_symlink_test.go`
- `internal/speccheck/evidence_record_test.go`
- `internal/speccheck/glossary_test.go`
- `internal/speccheck/qa_row_carry_test.go`
- `internal/speccheck/receipt_characterization_test.go`
- `internal/speccheck/receipts_test.go`
- `internal/speccheck/report_shape_findings_test.go`
- `internal/speccheck/skills_declaration_test.go`
- `internal/speccheck/transcripts_test.go`
- `internal/speccheck/wrap_fix_quoting_test.go`
- `internal/store/journal_baseline_test.go`
- `internal/store/journal_consumer_corpus_test.go`
- `internal/store/journal_parallel_runs_test.go`
- `internal/store/migrate_test.go`
- `internal/store/process_info_unix_test.go`
- `internal/store/process_linux_test.go`
- `internal/store/repository_key_backfill_test.go`
- `internal/store/token_usage_test.go`
- `internal/worktree/adminlock_guard_test.go`
- `internal/worktree/adminlock_test.go`
- `internal/worktree/archive_record_test.go`
- `internal/worktree/disposal_test.go`
- `internal/worktree/item_nested_cleanup_test.go`
- `internal/worktree/prune_repository_key_test.go`
- `internal/worktree/staging_test.go`

## Carry-forward provenance

- Source Run: `run_20261008T215927Z_c26d3579de5392c8`
- Source commit: `6710697fbfc0717e1c73a9b9cf4dc968fcbe6b80`
