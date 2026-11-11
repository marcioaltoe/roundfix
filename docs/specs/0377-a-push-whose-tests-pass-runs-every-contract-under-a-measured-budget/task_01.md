---
task: task_01
spec: 0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget
status: completed
type: infra
complexity: medium
---

# Task 01: The push gate records its suite-time verdict and judges it after the contracts

## Overview

On 2026-11-10 the operator recorded that `CI | Verify` fails most pushes to
`main` on its 240 s suite-time budget while every test passes, and skips
`Verify docs` and the Full Contract Run
([`references/2026-11-10-the-push-gate-on-main-fails-its-suite-time-budget-on-most-runs.md`](references/2026-11-10-the-push-gate-on-main-fails-its-suite-time-budget-on-most-runs.md),
recorded by the operator; the measurement of 33 passing pushes is in
[`references/2026-11-10-ci-verify-push-suite-times.md`](references/2026-11-10-ci-verify-push-suite-times.md)).
This Task makes the budget one measured number, 420 s, in the `Makefile`,
makes `make test-budget` record its verdict in a file when asked and a new
`make suite-budget-check` judge it, wires the push path of the workflow so the
verdict is judged in a last step after the docs and contract steps, adds the
test that proves both files, writes the amendment notice of ADR-0387 into
ADR-0259 and regenerates the setup manifest. The authorized tooling paths are
exactly `.github/workflows/ci-verify.yml` and `Makefile`, plus the derived
`docs/agents/setup-context.json` through its sanctioned command.

## Requirements

1. MUST, in the `Makefile`, replace the comment block above
   `SUITE_BUDGET_SECONDS` and its default with exactly these lines:

   ```make
   # The full suite's wall clock on a fresh CI run, in seconds (ADR-0387). The
   # pushes to main of 2026-11-06 to 2026-11-10 passed every test in 209 s to
   # 343 s, because GitHub assigns runners of different speeds: the same tree took
   # 215 s on one and 311 s on another. 420 s is 22 % above the slowest run and
   # fails a suite that slows by about 1.95 times on the fast runner.
   # Change this deliberately, with a new measurement, and say why.
   SUITE_BUDGET_SECONDS ?= 420
   ```

2. MUST replace the `test-budget` target (comment above it, recipe) with
   exactly this block, followed by the new `suite-budget-check` target, right
   after `SUITE_BUDGET_SECONDS`; recipe lines start with a tab. It follows
   `_techspec.md` → API Contracts 1 and 2: without `SUITE_BUDGET_VERDICT` the
   output and exit status of `test-budget` are byte-identical to today's.

   ```make
   # With SUITE_BUDGET_VERDICT naming a file, test-budget writes its verdict there
   # and suite-budget-check judges it later, so a slow suite does not skip the
   # steps that follow it; without it, an overrun fails at once.
   test-budget: gocache-bound ## Run Go tests once and fail if the suite exceeds SUITE_BUDGET_SECONDS
   	@if [ -n "$(SUITE_BUDGET_VERDICT)" ]; then rm -f "$(SUITE_BUDGET_VERDICT)" || exit $$?; fi; \
   	start=$$(date +%s); \
   	$(GO) test -count=1 -parallel $(GO_TEST_PARALLEL) $(PKGS) || exit $$?; \
   	elapsed=$$(( $$(date +%s) - start )); \
   	if [ "$$elapsed" -gt "$(SUITE_BUDGET_SECONDS)" ]; then \
   		verdict="suite-time budget exceeded: the suite took $${elapsed}s against a budget of $(SUITE_BUDGET_SECONDS)s"; \
   	else \
   		verdict="suite-time budget: $${elapsed}s of $(SUITE_BUDGET_SECONDS)s"; \
   	fi; \
   	echo "$$verdict"; \
   	if [ -n "$(SUITE_BUDGET_VERDICT)" ]; then \
   		printf '%s\n' "$$verdict" > "$(SUITE_BUDGET_VERDICT)" || exit $$?; \
   		echo "suite-time verdict recorded in $(SUITE_BUDGET_VERDICT); make suite-budget-check judges it"; \
   		exit 0; \
   	fi; \
   	case "$$verdict" in "suite-time budget exceeded:"*) \
   		echo "the budget is SUITE_BUDGET_SECONDS in the Makefile; raise it deliberately or find the regression"; \
   		exit 1;; \
   	esac

   suite-budget-check: ## Fail if the verdict test-budget recorded in SUITE_BUDGET_VERDICT exceeds the budget
   	@if [ -z "$(SUITE_BUDGET_VERDICT)" ] || [ ! -f "$(SUITE_BUDGET_VERDICT)" ]; then \
   		echo "suite-time budget: no verdict recorded at '$(SUITE_BUDGET_VERDICT)'; run make test-budget with SUITE_BUDGET_VERDICT set first"; \
   		exit 1; \
   	fi; \
   	verdict=$$(cat "$(SUITE_BUDGET_VERDICT)") || exit $$?; \
   	echo "$$verdict"; \
   	case "$$verdict" in \
   	"suite-time budget exceeded:"*) \
   		echo "the budget is SUITE_BUDGET_SECONDS in the Makefile; raise it deliberately or find the regression"; \
   		exit 1;; \
   	"suite-time budget: "*) ;; \
   	*) \
   		echo "suite-time budget: unreadable verdict at '$(SUITE_BUDGET_VERDICT)'"; \
   		exit 1;; \
   	esac
   ```

3. MUST add `suite-budget-check` to the `.PHONY` line, right after
   `test-race`, and change no other target or variable of the `Makefile`.
4. MUST, in `.github/workflows/ci-verify.yml`, per `_techspec.md` → API
   Contract 3, make the push `Verify` step and its comment read exactly:

   ```yaml
         # Pushes to main owe a fresh verdict on the whole tree. -count=1 disables
         # the result cache. test-budget times the suite and records its verdict;
         # the last step judges it, so a slow runner never skips the docs and
         # contract runs (ADR-0387). The budget is the Makefile's.
         - name: Verify
           id: verify
           if: github.event_name == 'push'
           env:
             GOFLAGS: -count=1
             VERIFY_TEST_TARGET: test-budget
             SUITE_BUDGET_VERDICT: ${{ runner.temp }}/suite-budget-verdict
           run: make verify
   ```

   and add, as the job's last step after `Verify every contract`, separated
   by one blank line:

   ```yaml
         # The suite-time budget is judged last, whenever the timed suite passed,
         # even when the docs or contract runs failed (ADR-0387).
         - name: Check the suite-time budget
           if: ${{ !cancelled() && github.event_name == 'push' && steps.verify.outcome == 'success' }}
           env:
             SUITE_BUDGET_VERDICT: ${{ runner.temp }}/suite-budget-verdict
           run: make suite-budget-check
   ```

   Change no other step, trigger, permission, or the pull request path.
5. MUST add `internal/verifyselect/suite_budget_test.go` with the two tests of
   `_techspec.md` → Testing Approach, every test and subtest calling
   `t.Parallel()` first:
   - `TestTheSuiteTimeBudgetIsRecordedThenJudged` runs the real `Makefile`
     with `make --no-print-directory -f <Makefile>` in a `t.TempDir()`, with
     `GO=` an executable `go` stub script in that directory that exits 1 for
     `test` only in the failing subtest and 0 otherwise, and an environment
     of only `PATH`, `TMPDIR` and `HOME` (both set to that directory), plus the
     variable the subtest names. Its eight subtests and their assertions:
     `an overrun is recorded without failing` (`SUITE_BUDGET_VERDICT` from the
     environment, `SUITE_BUDGET_SECONDS=-1`: exit 0, the file holds exactly one
     line matching `^suite-time budget exceeded: the suite took [0-9]+s against a budget of -1s$`,
     and the output lacks the hint line); `without a verdict file an overrun
     fails at once` (budget -1, no variable: non-zero, the output holds the
     exceeded line and the hint, no file); `a failing suite records no
     verdict` (a stale file holding `suite-time budget: 1s of 420s`, the stub
     failing `test`: non-zero and the file is gone); `the default budget is
     420 seconds` (no `SUITE_BUDGET_SECONDS`: exit 0 and the file holds
     exactly one line matching `^suite-time budget: [0-9]+s of 420s$`); `the
     check fails an exceeded verdict` (file holding
     `suite-time budget exceeded: the suite took 500s against a budget of 420s`:
     non-zero, the output holds that line followed by the hint line); `the
     check passes a verdict within the budget` (file holding
     `suite-time budget: 300s of 420s`: exit 0 and the output is exactly that
     line); `the check fails without a verdict` (no file: non-zero, output
     holds `suite-time budget: no verdict recorded at `); `the check fails an
     unreadable verdict` (file holding `garbage`: non-zero, output holds
     `suite-time budget: unreadable verdict at `).
   - `TestTheCIPushJudgesTheSuiteTimeBudgetLast` parses the workflow with
     `gopkg.in/yaml.v3` and asserts API Contract 3: no step sets
     `SUITE_BUDGET_SECONDS`; exactly one step sets
     `VERIFY_TEST_TARGET: test-budget`, the `Verify` step, whose `id` is
     `verify`, `if` is `github.event_name == 'push'`, `run` is `make verify`
     and `SUITE_BUDGET_VERDICT` is `${{ runner.temp }}/suite-budget-verdict`;
     the last step of the `verify` job is `Check the suite-time budget`, after
     `Verify every contract`, with that `if`, that `run` and that verdict
     path.
6. MUST, in `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`,
   set `updated_at: 2026-11-10T23:00:00Z` and add, after the last
   `## Consequences` bullet ("The CI suite budget stays at 240 s until CI
   measures the converted suite."), this paragraph, unchanged:
   `**Amended by [ADR-0387](0387-the-suite-time-budget-is-measured-on-ci-and-judged-after-the-contract-runs.md) (2026-11-10).** CI measured the converted suite on the pushes to main from 2026-11-06 to 2026-11-10 at 209 s to 343 s, depending on the runner. The CI suite budget is now 420 s, kept in the Makefile, and it is judged after the docs and contract runs instead of stopping them. Every other decision of this record still holds.`
7. MUST, as the last step of the Task, after every `Makefile` edit,
   regenerate the setup manifest with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   and confirm that a second run reports `File changes: 0`; never hand-edit
   `docs/agents/setup-context.json`.
8. MUST NOT change any other `Makefile` target or variable, any other workflow,
   `release.yml`, `go.mod`, `go.sum`, `.roundfixrc.yml`, Go production code,
   any guide, skill, `CONTEXT.md`, any other ADR, or any existing test. No
   test reaches the network, starts a workflow, reads the real `~/.roundfix`
   or writes outside its temporary directories, and the test file compiles
   for Windows.

## Subtasks

- [ ] Set the measured budget, record the verdict in `test-budget` and add `suite-budget-check`.
- [ ] Wire the push `Verify` step and the last budget step in the workflow.
- [ ] Add the two tests, the ADR-0259 notice, and regenerate the setup manifest last.

## Acceptance Criteria

- [ ] With `SUITE_BUDGET_VERDICT` set, an overrun is recorded and does not fail `make test-budget`; `make suite-budget-check` fails it; without the variable an overrun fails at once with today's lines.
- [ ] The push path records the verdict, keeps `Verify docs` and `Verify every contract`, and judges the verdict last; no step sets `SUITE_BUDGET_SECONDS`, and the default is 420.
- [ ] ADR-0259 carries the notice, the setup manifest is current, and every `./internal/verifyselect` test passes.

## Context

- interface: `Makefile`
- interface: `.github/workflows/ci-verify.yml`
- interface: `docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/verifyselect/suite_budget_test.go`
- instruction: `docs/adr/0387-the-suite-time-budget-is-measured-on-ci-and-judged-after-the-contract-runs.md`
- instruction: `internal/verifyselect/full_contract_run_test.go`
- instruction: `internal/verifyselect/verify_docs_makefile_test.go`
- instruction: `internal/verifyselect/verifyselect_test.go`
- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/go.md`

## Consumers

- The `test-budget` and `suite-budget-check` targets and the push steps of `.github/workflows/ci-verify.yml` — consumers: `./internal/verifyselect`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheSuiteTimeBudgetIsRecordedThenJudged|TestTheCIPushJudgesTheSuiteTimeBudgetLast)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheCIPushJudgesTheSuiteTimeBudgetLast TestTheSuiteTimeBudgetIsRecordedThenJudged/an_overrun_is_recorded_without_failing TestTheSuiteTimeBudgetIsRecordedThenJudged/without_a_verdict_file_an_overrun_fails_at_once TestTheSuiteTimeBudgetIsRecordedThenJudged/a_failing_suite_records_no_verdict TestTheSuiteTimeBudgetIsRecordedThenJudged/the_default_budget_is_420_seconds TestTheSuiteTimeBudgetIsRecordedThenJudged/the_check_fails_an_exceeded_verdict TestTheSuiteTimeBudgetIsRecordedThenJudged/the_check_passes_a_verdict_within_the_budget TestTheSuiteTimeBudgetIsRecordedThenJudged/the_check_fails_without_a_verdict TestTheSuiteTimeBudgetIsRecordedThenJudged/the_check_fails_an_unreadable_verdict; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task neither test exists, so their pass lines are missing; afterwards both tests and all eight subtests pass.
- `tr -s '[:space:]' ' ' < docs/adr/0259-a-test-runs-in-parallel-unless-it-names-why-it-cannot.md | grep -qF -- 'The CI suite budget is now 420 s, kept in the Makefile' || { echo 'ADR-0259 lacks the amendment notice' >&2; exit 1; }; go test -count=1 -timeout 20m ./internal/verifyselect && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestTheSetupManifestIsCurrent|TestGlossaryPhrasesAreDefinedTerms)$' && GOOS=windows go vet ./internal/verifyselect` — expected: exit 0. Before this Task ADR-0259 lacks the notice, so the first check fails; afterwards every Makefile and workflow test of the package passes, the setup manifest is current after the `Makefile` edit, every capitalized phrase is a defined term, and the package vets for Windows.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; User Story 3; Core Feature 1; Core Feature 2; Core Feature 3; Core Feature 4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Testing Approach; Build Order 1
- ADR-0387; ADR-0259; ADR-0253

## Result

Skills loaded: implement-task, coding-guidelines, golang-testing, golang-structs-interfaces, testing-boss, tech-writer, domain-modeling, no-workarounds, setup-context-driven, roundfix, evidence-gate.

The Makefile now owns the measured 420 s default. With
`SUITE_BUDGET_VERDICT` set, `test-budget` removes a stale verdict, records a
passing suite's verdict even on an overrun, and leaves the decision to
`suite-budget-check`. A failing suite leaves no verdict. Without the variable,
an overrun still fails immediately with the verdict and hint.

The push `Verify` step records the verdict under `runner.temp`, and the last
step checks it whenever that step succeeded, including after a docs or
contract failure. The pull request path and existing docs and contract steps
retain their bytes. ADR-0259 carries the exact ADR-0387 amendment and requested
`updated_at` value.

Focused evidence by acceptance criterion:

1. Recorded and immediate verdict behavior:
   `go test -count=1 -run 'TestThe(SuiteTimeBudgetIsRecordedThenJudged|CIPushJudgesTheSuiteTimeBudgetLast)' -v ./internal/verifyselect`
   exited 0 after the change: both tests and all eight subtests passed. Before
   the Makefile/workflow edits, the same focused command failed the workflow
   test and seven subtests; only the existing immediate-overrun case passed.
   The cases cover recording an overrun without its hint, immediate failure,
   stale-file removal after a test failure, the 420 s default, exceeded and
   passing recorded verdicts, and missing and unreadable verdicts. Each test
   and subtest starts with `t.Parallel()`; each make subprocess gets only
   `PATH`, disposable `TMPDIR` and `HOME`, and its optional verdict variable.
2. Push wiring and preserved consumers: the workflow test above passed.
   `go test -count=1 -run '^(TestVerifyContractsRunsEveryContract|TestMainAndReleaseRunEveryContract|TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped|TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts)$' ./internal/verifyselect`
   exited 0. A Python source comparison confirmed the authored Makefile and
   YAML blocks exactly, and confirmed all neighboring Makefile and existing
   workflow bytes are unchanged except the authorized `.PHONY` addition.
3. Documentation and generated manifest: the source comparison confirmed the
   exact ADR notice and timestamp. As the last implementation step,
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   exited 0 and changed only `docs/agents/setup-context.json` (`File changes:
   1`, approved plan digest
   `sha256:e8d59cf21f5970f7f2d058dc42eb2610cabde957279280dcd5926bb6d948dfad`).
   Its second run exited 0 with `File changes: 0`, `Idempotence: verified`,
   and plan digest
   `sha256:93bfc4f8108f35867522c9c7be74efe19ec38949b7d45f33773a1a2dc7d2dbc7`.
   Both runs reported unchanged nested instruction-carrier warnings in
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/AGENTS.md`
   and
   `internal/baseline/assets/source-baselines/baseline.standard-typescript-monorepo-0.0.1/corpus/AGENTS.md`.
   The generated manifest was never edited by hand.

Additional focused check: `GOOS=windows go test -c -o <temporary-directory>/verifyselect.test.exe ./internal/verifyselect`
exited 0; the executable was written to a disposable directory and removed
with it. `git diff --check` exited 0. The initial worktree contained only the
Daemon's task-status change; implementation touches only the five declared
implementation paths and this Result section.

The declared Verification commands, full package run, docscontract checks,
Windows vet, and repository Verification remain for the Daemon. No status,
Task Graph, other Task, commit, push, or pull request was authored here.

## Carry-forward provenance

- Source Run: `run_20261111T175059Z_b2ca8b179566f7a4`
- Source commit: `d8651e47b71fb74994202d77588942b0f82a92fa`
