---
task: task_01
spec: 0312-the-ci-job-runs-the-documentation-contracts-once
status: pending
type: infra
complexity: medium
---

# Task 01: A pull request runs the documentation contracts once

## Overview

The `ci-verify` job of a pull request runs every `internal/docscontract` test
in "Verify changed", because each of their files is `//verify:always`, and
runs the whole package again in "Verify docs" through `docs-test`. This Task
makes `verify-docs` take its documentation contract target from
`VERIFY_DOCS_TEST`, `docs-test` by default, adds a `docs-test-skipped` target
that says why it runs nothing, makes the "Verify docs" step choose it on pull
requests only, and adds the repository test that keeps every
`internal/docscontract` file `//verify:always`, the condition that makes the
skip safe.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-22,
   [`references/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md`](references/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md)
   (`_prd.md` → Sources 1), by implementing API Contract 1 of `_techspec.md`.
2. MUST change the `Makefile`:
   - add `VERIFY_DOCS_TEST ?= docs-test` above `verify-docs`, with a comment
     that a pull request's CI job already ran every `internal/docscontract`
     test in `make verify-changed` because each of their files is
     `//verify:always`;
   - make the prerequisites of `verify-docs` read
     `build $(VERIFY_DOCS_TEST) repo-test spec-budget spec-check`, keeping
     its `##` help text;
   - add the target `docs-test-skipped`, whose only recipe line is
     `@echo "docs-test: skipped; make verify-changed ran internal/docscontract earlier in this job"`,
     and add it to `.PHONY`.
   Every other target keeps its bytes.
3. MUST add to the "Verify docs" step of `.github/workflows/ci-verify.yml`,
   between its `name` and its `run`, an `env` with
   `VERIFY_DOCS_TEST: ${{ github.event_name == 'pull_request' && 'docs-test-skipped' || 'docs-test' }}`.
   The step keeps its name, its `run: make verify-docs` and its place
   directly before "Verify every contract"; every other step keeps its bytes.
4. MUST add the new file `internal/verifyselect/verify_docs_makefile_test.go`
   with a file header in the shape of `full_contract_run_test.go` and three
   tests that each call `t.Parallel()` first:
   - `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped`: in a
     temporary directory, a `go` stub and a `roundfix-stub` that append
     their arguments to a file; `make --no-print-directory -f <repository
     Makefile> verify-docs GO=<go stub> BIN=<roundfix stub>`, once as is and
     once with `VERIFY_DOCS_TEST=docs-test-skipped`; it counts the recorded
     `go test` lines that name `./internal/docscontract` without
     `TestCheckCorpusBudget` (1, then 0) and requires one `go build`, one
     `-tags repocontract` run, one `TestCheckCorpusBudget` run and one
     `spec check` in both, and the skip message only in the second;
   - `TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts`:
     parses `ci-verify.yml` with `gopkg.in/yaml.v3` and requires exactly one
     "Verify docs" step in the `verify` job, with `run: make verify-docs`
     and the Requirement 3 value of `VERIFY_DOCS_TEST`;
   - `TestEveryDocumentationContractFileIsAlwaysSelected`: every
     `_test.go` file of `internal/docscontract` whose source holds
     `//go:build docscontract` has a `//verify:always` line, and at least one
     such file exists.
5. MUST keep `TestMainAndReleaseRunEveryContract`, `TestMakeSharesOneGoBuildCache`,
   `TestVerifyContractsRunsEveryContract` and the `internal/suiteguardcontract`
   tests passing unchanged.
6. MUST NOT change `internal/docscontract`, any production file, any skill
   or any other workflow.
7. MUST, as the last step after the Makefile edit, regenerate this
   repository's Setup Manifest, whose verification entries record the
   Makefile digest, only with
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`,
   confirm that a second run reports no file change, and change no file
   under `docs/agents/` other than `docs/agents/setup-context.json`
   (`_authorization.md` → Sanctioned regeneration).

## Subtasks

- [ ] Make the documentation contract step of `verify-docs` selectable and add `docs-test-skipped`.
- [ ] Select it per event in the "Verify docs" step.
- [ ] Add the three tests.

## Acceptance Criteria

- [ ] `make verify-docs` runs the documentation contracts once, and none with `VERIFY_DOCS_TEST=docs-test-skipped`, with every other step unchanged.
- [ ] The workflow skips them on pull requests only.
- [ ] A documentation contract file that stops being `//verify:always` fails a test.

## Consumers

None. The change is confined to the Makefile, one workflow step and new tests; no Go package contract changes.

## Context

- instruction: `internal/verifyselect/full_contract_run_test.go`
- instruction: `internal/verifyselect/build_cache_makefile_test.go`
- instruction: `internal/verifyselect/verifyselect_test.go`
- instruction: `internal/suiteguardcontract/repository_gate_test.go`
- instruction: `internal/docscontract/corpus_test.go`
- instruction: `docs/references/2026-10-24-test-suite-audit.md`
- creates: `internal/verifyselect/verify_docs_makefile_test.go`
- interface: `Makefile`
- interface: `.github/workflows/ci-verify.yml`

## Verification

- `out="$(go test -count=1 -timeout 15m -v -run '^(TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped|TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts|TestEveryDocumentationContractFileIsAlwaysSelected|TestMainAndReleaseRunEveryContract|TestMakeSharesOneGoBuildCache|TestVerifyContractsRunsEveryContract)$' ./internal/verifyselect 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts TestEveryDocumentationContractFileIsAlwaysSelected TestMainAndReleaseRunEveryContract TestMakeSharesOneGoBuildCache; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -timeout 15m ./internal/suiteguardcontract` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check; with `verify-docs` depending on `docs-test` directly the skipped case fails, without the workflow `env` the workflow test fails, and with any `internal/docscontract` file marked other than `//verify:always` the declaration test fails. Afterwards the new tests, the Full Contract Run, shared-cache and repository-gate tests pass.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Success Metric 1; Success Metric 2
- `_techspec.md` → Current behavior; API Contract 1; Invariants 1-2; Testing Approach; Build Order 1
