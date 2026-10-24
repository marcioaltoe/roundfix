---
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
status: active
created: 2026-10-24
surfaces: [backend]
---

# internal/cli tests stop reading the live documentation

Minimal PRD for a test-suite performance fix, minted by `write-techspec`.
Thirteen `internal/cli` tests read this repository's own guides, skills and
records, so a one-line edit of a command guide or a skill makes Go run the
whole `internal/cli` package again, the longest of the suite at 160 to 240
s, although no code changed. After this Spec those thirteen tests live in
two `internal/cli` files under the `docscontract` build tag, declared
`//verify:always`, as ADR-0326 decides: the selective gate still runs them
on every change and the Full Contract Run on main and before a release,
while the package's cached result survives a documentation edit. Every
assertion they make stays as it is.

## Sources

1. The Backlog Entry adopted by this planning change,
   [`references/2026-10-24-internal-cli-tests-read-the-live-documentation-trees.md`](references/2026-10-24-internal-cli-tests-read-the-live-documentation-trees.md),
   type `perf`, written from the test suite audit
   [`docs/references/2026-10-24-test-suite-audit.md`](../../references/2026-10-24-test-suite-audit.md):
   with the documentation trees removed or unreadable exactly thirteen
   `internal/cli` tests fail; a guide edit re-ran the package (162.6 s and
   237.0 s) unless those thirteen were skipped.

## Prerequisites

None in `requires`. This Spec is authored against `56dd1a76` (v0.116.0) and
delivers after Specs 0311 and 0312 and before Specs 0314, 0307, 0309 and
0310, by the maintainer's answer of 2026-10-24 quoted in `_authorization.md`.
Its sources cite the audit reference, which Spec 0311's planning change adds,
so the operator merges the plans together. None of the pending plans moves
the thirteen tests; a plan that edits one of their bodies must follow the
test to its new file.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; the tests move between files unchanged.
  Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0326 (a code test never reads the
  live documentation trees), recorded by this planning change, is the
  decision both implementation Tasks carry out. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — protected paths limited to the
  repository guide line and the two governed test files the tests leave; the Makefile, the CI workflows and the build tags'
  wiring stay unchanged, and the cycle AP batch 4 approval of 2026-10-24
  grants the delivery operations. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0313-internal-cli-tests-stop-reading-the-live-documentation/_authorization.md`;
  bounded files: `docs/agents/specific-repository.md`,
  `internal/cli/cli_test.go`,
  `internal/cli/baseline_documentation_contract_test.go`.

## Goals

- A guide or skill edit leaves the `internal/cli` test result cached.
- The thirteen tests keep every assertion and still run in
  `make verify-changed` and `make verify-contracts`.
- A new `internal/cli` documentation contract file cannot land without
  `//verify:always`.

## Core Features

1. Repository-profile documentation contracts. The four repository-profile
   Baseline tests of `baseline_update_repository_profile_test.go` and
   `TestBaselineUpdatePreviewListsASetupSnapshotDifference` move, unchanged,
   to `internal/cli/repository_profile_docscontract_test.go`.
2. The other documentation contracts. The skill-set, guide-example, HTTP
   default statement, events help and release plan checks tests move,
   unchanged, to `internal/cli/live_documents_docscontract_test.go`, and the
   repository guide states the rule.
3. The repository contract declarations test classifies every
   `internal/cli` `docscontract` file as always selected and pins their
   number.

## Non-Goals / Out of Scope

- Rewriting any of the thirteen tests onto fixture copies.
- Other packages whose tests read source files of other packages
  (`internal/spec`, `internal/testfixture`, `internal/verifyselect`,
  `internal/suiteguardcontract`); their rules read code by design.
- The Makefile, the CI workflows and Spec 0312's documentation contract
  skip, which concerns only `internal/docscontract`.
- The shared build cache: Spec 0314.

## Success Metrics

1. Success Metric: `go test -list` of `./internal/cli` without tags lists
   none of the thirteen tests, and `go test -tags docscontract` runs and
   passes all thirteen.
2. Success Metric: `cmd/verify-select -contracts -all` lists the thirteen
   with `internal/cli`, and `TestRepositoryContractTestsDeclareTheirRelevance`
   classifies both new files as always selected.
3. Success Metric: in a temporary worktree of the audited head, after one
   uncached run of `go test -parallel 16 ./internal/cli`, appending a line to
   `docs/user-guide/commands/gc.md` and to `.agents/skills/roundfix/SKILL.md`
   leaves the next run `(cached)`; at the Delivery Base the same edit makes
   it run again.
4. Success Metric: the untagged `internal/cli` package passes with the five
   documentation trees unreadable.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The Go command documentation, `go help test` in the audited checkout:
  "Tests that open files within the package's module or that consult
  environment variables only match future runs in which the files and
  environment variables are unchanged."
- The test suite audit committed by Spec 0311's planning change before any
  Task of this Spec ran: the 162.6 s and 237.0 s re-runs after a guide edit,
  and the cached result with the thirteen tests skipped.

## Validation

The decisions were checked on 2026-10-24 against these sources.

- A prototype in a scratch worktree of `56dd1a76`: the thirteen tests moved
  with their doc comments into the two files; `go vet` passed with and
  without the tag and for Windows with both contract tags; the thirteen
  passed under the tag; `internal/testfixture` and `internal/verifyselect`
  passed once the declarations test classified the files; the selector
  listed the thirteen with `internal/cli`.
- Secondbrain, read only: `qmd query "roundfix test suite flaky slow
  performance" --all --files --min-score 0.3` returned the Roundfix mirror's
  records of Specs 0265 and 0103 and nothing on test caching. No Inbox Entry
  was written.
- Exa: none needed beyond the Go command's own documentation, which the QA
  reads offline.

## Glossary

None.

## Decisions

- 2026-10-24 (about 00:05 BRT): the cycle AP batch 4 answers quoted in
  `_authorization.md`.
- 2026-10-24 (about 01:50 BRT): "Escrever agora (Recomendado)" for this Spec.

Autonomous route: no person present; decisions taken from the operator's cycle AP batch 4 briefing of 2026-10-24, the maintainer's answers of 01:50 BRT and the audit.

Taken from those sources, with no further interview:

- Move the tests behind the tag rather than rewrite them on fixtures
  (ADR-0326): they assert what this repository ships.
- Two files, one per subject, so each Task moves one group.

## Open Questions

None.
