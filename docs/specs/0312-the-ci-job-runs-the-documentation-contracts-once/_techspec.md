---
spec: 0312-the-ci-job-runs-the-documentation-contracts-once
prd: _prd.md
created: 2026-10-24
---

# The CI job runs the documentation contracts once — Technical Spec

## Executive Summary

`verify-docs` takes its documentation contract target from
`VERIFY_DOCS_TEST`, `docs-test` unless the caller says otherwise, and the
`ci-verify` "Verify docs" step passes `docs-test-skipped` on pull requests,
where "Verify changed" has just run every `internal/docscontract` test
because each of their files is `//verify:always`. The change is proved with
`make` runs over stub tools in temporary directories, the way
`internal/verifyselect` already proves the Makefile's cache and Full
Contract Run wiring, and a workflow parse.

The trade-off accepted: on a pull request the documentation contracts run
only under the selective step's command line, so a future
`internal/docscontract` file that is not `//verify:always` would escape the
pull request; a repository test refuses such a file.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — the Spec records no decision and
  changes no accepted one, as `_prd.md` → Project Constraints states.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — protected tooling mutation authorized for
  the Makefile and `.github/workflows/ci-verify.yml`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0312-the-ci-job-runs-the-documentation-contracts-once/_authorization.md`;
  bounded files: `Makefile`, `.github/workflows/ci-verify.yml`.

## Current behavior

Read on 2026-10-24 at `56dd1a76`.

1. `Makefile`: `verify-docs: build docs-test repo-test spec-budget
   spec-check`; `docs-test` runs `go test -count=1 -timeout 20m -tags
   docscontract ./internal/docscontract`; `spec-budget` runs the same package
   with `-parallel=1 -run '^TestCheckCorpusBudget$'`.
   `verify-changed-contracts` runs every contract `cmd/verify-select
   -contracts` selects, and every `internal/docscontract` test file carries
   `//verify:always` (nine files; the tenth, `suiteguard_repocontract_test.go`,
   is a `repocontract` file).
2. `.github/workflows/ci-verify.yml`: on a pull request, "Verify changed"
   (`make verify-changed`, `GOFLAGS=-count=1`), then "Verify docs"
   (`make verify-docs`); on a push, "Verify" (`make verify`), "Verify docs"
   and "Verify every contract" (`make verify-contracts`), which
   `TestMainAndReleaseRunEveryContract` requires directly after "Verify
   docs".
3. On PR #654 "Verify changed" took 271 s and "Verify docs" 173 s; locally
   `docs-test` took 14 s, `repo-test` 63 s and `spec-budget` 9 s.

## System Architecture

- `Makefile`: `VERIFY_DOCS_TEST`, `docs-test-skipped`, and the new
  prerequisite of `verify-docs`.
- `.github/workflows/ci-verify.yml`: the "Verify docs" step's `env`.
- `internal/verifyselect/verify_docs_makefile_test.go` (new).

## Implementation Design

### Interfaces

```make
VERIFY_DOCS_TEST ?= docs-test
verify-docs: build $(VERIFY_DOCS_TEST) repo-test spec-budget spec-check
docs-test-skipped:   # prints why, runs nothing
```

### Data Models

None. No persisted format changes.

### API Contracts

1. API Contract: `make verify-docs` runs `docs-test` unless
   `VERIFY_DOCS_TEST` names another target; `docs-test-skipped` prints
   `docs-test: skipped; make verify-changed ran internal/docscontract earlier
   in this job` and runs no command. The "Verify docs" step sets
   `VERIFY_DOCS_TEST: ${{ github.event_name == 'pull_request' &&
   'docs-test-skipped' || 'docs-test' }}`; both values are non-empty, so the
   expression never falls through on an empty string.

### Surface Transcripts

None. No command surface changes.

### Invariants

1. On a pull request every `internal/docscontract` test runs once in the
   `ci-verify` job; on a push, and in `make verify-docs` run alone, the
   package runs as today.
2. Every `docscontract`-tagged test file of `internal/docscontract`
   declares `//verify:always`.

## Coverage Map

- Goal 1, Core Feature 1, Success Metrics 1 and 2 → API Contract 1,
  Invariants 1 and 2, task_01.
- Goal 2 → Invariant 1, `TestMainAndReleaseRunEveryContract`, task_01.
- Goal 3 → Invariant 2, `TestEveryDocumentationContractFileIsAlwaysSelected`,
  task_01.
- Success Metric 3 → Unreachable Acceptance, operator after merge.

Deviations from the PRD: None.

## Absence Evidence

1. A skippable documentation contract step — command:
   `grep -n "VERIFY_DOCS_TEST\|docs-test-skipped" Makefile .github/workflows/ci-verify.yml`;
   result: no match.

## ADR Alternatives

None. The Active ADR obligations row lists no ADR: the Spec records no
decision and changes no accepted one.

## Integration Points

GitHub Actions reads the new `env` of one step; no other external system.

## Testing Approach

`internal/verifyselect/verify_docs_makefile_test.go` (new), in the shape of
`TestVerifyContractsRunsEveryContract`:

- `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped` runs
  `make -f <repo Makefile> verify-docs GO=<stub> BIN=<stub>` in a temporary
  directory with a `go` stub that appends its arguments to a file and a
  `roundfix-stub`, once by default and once with
  `VERIFY_DOCS_TEST=docs-test-skipped`; it counts the recorded `go test`
  lines that name `./internal/docscontract` without `TestCheckCorpusBudget`
  (1, then 0), and requires the build, `repo-test`, `spec-budget` and
  `spec-check` invocations both times and the skip message the second time.
- `TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts` parses
  `ci-verify.yml` and requires the "Verify docs" step's `env` to hold the
  API Contract 1 expression, and still `run: make verify-docs`.
- `TestEveryDocumentationContractFileIsAlwaysSelected` reads every
  `_test.go` file of `internal/docscontract` whose source holds
  `//go:build docscontract` and requires a `//verify:always` line.

## Build Order

1. The documentation contract skip, the workflow line and their tests
   (task_01).
2. The QA gate (task_02) (depends on: 1).

No other Spec is required; the audit reference comes from Spec 0311's
planning change, merged with this one.

## Risks & Considerations

- A change that adds a documentation contract file without
  `//verify:always` now fails a repository test instead of silently running
  only in "Verify docs".
- Spec 0313 adds documentation contracts in `internal/cli`; they run through
  the selective step's Contract Relevance, not through `docs-test`, so this
  skip does not touch them.

## Glossary

None.

## Decisions

- The skip lives in "Verify docs" and is selected by the workflow, so
  `make verify-docs` alone is unchanged.
- No ADR: the change moves no contract a later reader needs beyond the
  Makefile comment.
