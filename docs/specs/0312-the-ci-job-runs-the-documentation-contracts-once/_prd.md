---
spec: 0312-the-ci-job-runs-the-documentation-contracts-once
status: active
created: 2026-10-24
surfaces: [backend]
---

# The CI job runs the documentation contracts once

Minimal PRD for a verification-cost fix, minted by `write-techspec`. The
`ci-verify` job of every pull request runs `internal/docscontract` twice,
once in "Verify changed", where every one of its test files is
`//verify:always`, and again in "Verify docs". After this Spec a pull
request's "Verify docs" step skips the contract package the selective step
already ran, while `make verify-docs` run alone and every push to main still
run it, and a repository test keeps every documentation contract file
`//verify:always`, the condition that makes the skip safe. The maintainer
authorized the Makefile and workflow change on 2026-10-24.

## Sources

1. The open Backlog Entry adopted by this planning change,
   [`references/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md`](references/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md),
   type `perf`: the package took about 764 s in "Verify changed" and 541 s in
   "Verify docs" with five active Specs before Spec 0305; its target is one
   run per `ci-verify` job with `make verify-docs` alone still running it.
2. The test suite audit
   [`docs/references/2026-10-24-test-suite-audit.md`](../../references/2026-10-24-test-suite-audit.md),
   committed by Spec 0311's planning change: the redundant runs in the CI job
   and the CI timings this Spec starts from.

The audit's shared build cache finding is Spec 0314's, by the maintainer's
answer of 2026-10-24 quoted in `_authorization.md`.

## Prerequisites

None in `requires`. This Spec is authored against `56dd1a76` (v0.116.0) and
delivers after Spec 0311 and before Specs 0313, 0314, 0307, 0309 and 0310.
Its Verification reads the audit reference, which Spec 0311's planning change
adds, so the operator merges the plans together. Specs 0307, 0309 and 0310
compare the Makefile and the CI workflows with their own Delivery Base, which
will include this Spec's change.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; tests run `make` with stub tools in
  temporary directories. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — the Spec records no decision and
  changes no accepted one: which contracts the selective gate runs, and the
  Full Contract Run on main and before a release, stay as the glossary
  states them (Contract Relevance, Full Contract Run); the skip only stops a
  second run of tests the selective step already ran. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — protected tooling mutation authorized for
  the Makefile and `.github/workflows/ci-verify.yml`, by the maintainer's
  answer of 2026-10-24 "Sim, só Makefile e workflows (Recomendado)". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0312-the-ci-job-runs-the-documentation-contracts-once/_authorization.md`;
  bounded files: `Makefile`, `.github/workflows/ci-verify.yml`.

## Goals

- A pull request's `ci-verify` job runs every `internal/docscontract` test
  once.
- `make verify-docs` run alone, and every push to main, still run the
  package.
- A documentation contract file that stops being `//verify:always` fails a
  test before it can escape a pull request.

## Core Features

1. One documentation contract run per pull request. The Makefile gains
   `VERIFY_DOCS_TEST`, `docs-test` by default; `verify-docs` depends on it,
   and a `docs-test-skipped` target says why it runs nothing. The "Verify
   docs" step of `ci-verify` sets it to `docs-test-skipped` on pull requests
   and to `docs-test` on pushes. A repository test fails if any
   `internal/docscontract` test file stops being `//verify:always`.

## Non-Goals / Out of Scope

- The repository contracts that both "Verify changed" and the `repo-test`
  step of "Verify docs" run, and the `spec-budget` run of
  `TestCheckCorpusBudget`; the audit measures them, and they stay.
- `internal/cli` tests that read the live documentation trees: Spec 0313.
- The shared Go build cache: Spec 0314.
- Any change to which tests a pull request, a push to main or a release
  runs, other than the second documentation contract run.

## Success Metrics

1. Success Metric: with `VERIFY_DOCS_TEST=docs-test-skipped`,
   `make verify-docs` starts no `go test` of `./internal/docscontract` other
   than the `spec-budget` run of `TestCheckCorpusBudget`, and still runs the
   build, `repo-test`, `spec-budget` and `spec-check`; without it, it runs
   the package exactly once.
2. Success Metric: the `ci-verify` workflow passes `docs-test-skipped` on
   pull requests and `docs-test` on pushes, and every `internal/docscontract`
   test file under the `docscontract` tag declares `//verify:always`.
3. Success Metric: measured on 2026-10-24 on the maintainer's Mac, the
   `internal/docscontract` run that a pull request stops repeating took 14 s
   at a load average of 2.5 with four active Specs; on a GitHub runner it is
   part of the 173 s "Verify docs" step of PR #654. The first pull request
   after this Spec shows a "Verify docs" step shorter by that package's run.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted Backlog Entry under `references/`, written by the operator on
  2026-10-22 from the CI runs of pull request #619: the package ran in both
  steps, 764 s and 541 s with five active Specs.
- The GitHub Actions expression documentation
  (https://docs.github.com/en/actions/reference/workflows-and-actions/expressions):
  "falsy values (`false`, `0`, `-0`, `""`, `''`, `null`) are coerced to
  `false`", so `cond && '' || 'b'` would always give `'b'`; that is why both
  values of `VERIFY_DOCS_TEST` are non-empty target names.

## Unreachable Acceptance

- criterion: Success Metric 3, the first pull request after this Spec shows
  a "Verify docs" step shorter by the documentation contract package's run.
  reason: only GitHub Actions runs the job, and reading its step timings
  needs the network and an open pull request, which no Verification or QA
  row may reach.
  satisfied-by: the operator's reading of that pull request's `ci-verify`
  run after this Spec merges.

## Validation

The decisions were checked on 2026-10-24 against these sources.

- A prototype in a scratch worktree of `56dd1a76`: `make -n verify-docs
  VERIFY_DOCS_TEST=docs-test-skipped` printed no documentation contract run
  and the default printed one; the three tests of task_01 passed, and each
  sabotage the QA lists failed its named test.
- The repository contracts of `make verify-docs` measured locally:
  `docs-test` 14 s, `repo-test` 63 s, `spec-budget` 9 s.
- Secondbrain, read only: `qmd query "roundfix test suite flaky slow
  performance" --all --files --min-score 0.3` returned the Roundfix mirror's
  records of Specs 0265 and 0103 and nothing on the CI job's steps. No Inbox
  Entry was written.
- Exa: the GitHub Actions expression documentation cited above.

## Glossary

None.

## Decisions

- 2026-10-24 (about 00:05 BRT): the cycle AP batch 4 answers quoted in
  `_authorization.md`; "Sim, só Makefile e workflows (Recomendado)"
  authorizes exactly the Makefile and workflow change this Spec declares.
- 2026-10-24 (about 01:50 BRT): "roundfix gc recolhe" moved the shared build
  cache out of this Spec into Spec 0314.

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 4 briefing of 2026-10-24, the maintainer's answers of 01:50 BRT, the adopted Backlog Entry and the audit.

Taken from those sources, with no further interview:

- Skip the second run in the "Verify docs" step, not in "Verify changed":
  the selective step is the Run's own Verification, and `make verify-docs`
  alone must keep running the package before a pull request.

## Open Questions

None.
