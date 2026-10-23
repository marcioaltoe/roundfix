---
spec: 0305-the-documentation-contracts-cost-seconds-per-active-spec
status: active
created: 2026-10-22
surfaces: [backend]
---

# The documentation contracts cost seconds per active Spec

Minimal PRD for a performance fix, minted by `write-techspec`. The
`internal/docscontract` package costs about a minute per active Spec because
each of its corpus tests runs `speccheck.Check` on every active Spec, five
Checks per Spec in all, and every Check spends almost all of its time in one
repository-wide pass that compares each path-like token of every Citation
Source with every Removed History Path, one string concatenation per pair.
After this Spec that pass looks each token up in a set, which makes one
Check of this repository about ten times cheaper, and the corpus tests share
one Check per active Spec per test process. Every finding and every contract
assertion stays as it is. The TechSpec holds the design and the measurements.

## Sources

One open Backlog Entry, adopted by this Spec's planning change into
`references/` (see `references/_index.md`):

1. [`references/2026-10-22-the-documentation-contract-suite-costs-about-a-minute-per-active-spec-and-runs-twice-in-one-ci-job.md`](references/2026-10-22-the-documentation-contract-suite-costs-about-a-minute-per-active-spec-and-runs-twice-in-one-ci-job.md),
   type `perf`, written by the operator on 2026-10-22 from the CI runs of
   pull request #619: with five active Specs the package took about 764 s in
   the "Verify changed" step and about 541 s in "Verify docs", and the
   `ci-verify` job took 30m21s and 32m20s.

The entry's second cost, the package running twice in one CI job, needs a
Makefile or workflow change that the cycle does not authorize. This planning
change records it as the open Backlog Entry
`docs/backlog/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md`,
the follow-up named under Non-Goals. The other open entry at the authoring
base, the QA refusal report's audited commit, is Spec 0302's source and
shares no context with this one; `docs/findings/` holds no entry.

## Prerequisites

None in `requires`. This Spec is authored against `42f523c2` (v0.109.0),
with Specs 0299 and 0300 active, and delivers after them and after the cycle
AP Specs 0301 to 0304. None of those plans changes
`internal/spec/citation_sources.go` or `internal/docscontract/corpus_test.go`.
Every active Spec on the delivery tree is part of the corpus the shared sweep
checks, so the corpus golden counts and the active-corpus assertions read
that tree, not this authoring base.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or
  schema is minted; the change replaces a lookup inside one function and
  shares test results inside one test process. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses in-memory content,
  temporary fixture repositories or this repository's own documents. Source:
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: not applicable — the Spec records no decision and
  changes no behavior: every finding of `roundfix spec check`, every
  Contract Relevance and every contract assertion stays as it is, and no
  accepted ADR constrains how a Check computes its findings or how a test
  package shares its own results. The Contract Relevance of the
  `internal/docscontract` tests stays `//verify:always`, and the Sequential
  Test rule is followed as the glossary states it. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized. The Makefile and the CI workflows stay unchanged; the follow-up
  that would change them is a Backlog Entry. The cycle AP approval of
  2026-10-22 through AskUserQuestion ("Ruído C2 e suite docscontract",
  authorization the same as cycle AO) grants the delivery operations. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0305-the-documentation-contracts-cost-seconds-per-active-spec/_authorization.md`;
  bounded files: none (`paths: []`).

## Goals

- One `speccheck.Check` of this repository costs about a tenth of its CPU
  time today, with the same findings.
- The removed-history citation pass costs the same whatever the number of
  Removed History Paths, for a source with a fixed content.
- The four corpus tests of `internal/docscontract` run one Check per active
  Spec per test process, plus the one `CheckStage` call whose equality with
  that Check they assert.
- Every assertion of the corpus tests and of the `SC-HISTORY-RETAINED`
  citation rule keeps failing on the inputs it fails on today.

## Core Features

1. Set lookup of Removed History Paths. `spec.HistoryCitationsInSource`
   builds a set of the Removed History Paths once per call and answers each
   path-like token by looking up the token and each of its parent
   directories, instead of comparing the token with every Removed History
   Path. It reports the same History Citations, in the same order, for every
   input: a token equal to a removed path, a token below a removed path, a
   token with a trailing slash, an explicit `./` or `../` prose path, and a
   token that only shares a string prefix with a removed path.
2. One shared corpus sweep. `TestCheckCorpusGolden`, `TestCheckCorpusBudget`,
   `TestCheckActiveCorpusHasNoErrors` and
   `TestStageScopeDefaultSweepIsUnchanged` read one memoized sweep that runs
   `speccheck.Check` once per active Spec per test process.
   `TestStageScopeDefaultSweepIsUnchanged` keeps calling
   `speccheck.CheckStage(..., StageAll)` and compares it with the shared
   result. A new contract test proves the memo on a fixture corpus.

## Non-Goals / Out of Scope

- The Makefile targets `verify-changed-contracts` and `docs-test` and the
  `ci-verify` workflow, which run the package twice in one CI job; the open
  Backlog Entry
  `docs/backlog/2026-10-22-the-ci-verify-job-runs-the-documentation-contract-package-twice.md`
  keeps that follow-up.
- Any change to the findings, severities, summaries or fix lines of
  `roundfix spec check`, to the corpus golden counts, or to the Contract
  Relevance of any test.
- The Markdown pass of `HistoryCitationsInSource`
  (`archiveContentCitations`), which still reads each line once per Removed
  History Path; after this Spec it is about a fifth of a Check's CPU time.
- Caching results across `speccheck.Check` calls in production code; the
  Daemon and the CLI keep computing every detector on every call.

## Success Metrics

1. Success Metric: `TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths`
   counts, for a fixed 200-token line, at most 64 more allocations per call
   with 2,000 Removed History Paths than with one; measured on 2026-10-22, the
   current code makes 1,603,622 against 2,423 and the set lookup 1,633
   against 1,625.
2. Success Metric: `TestCheckLoopOrderRepositoryAgrees`, one Check of the
   repository, measured on the same machine at the audited head and at the
   Delivery Base, uses at most a fifth of the user CPU time at the Delivery
   Base; measured on 2026-10-22, 15.7 s to 22.5 s before and 1.3 s to 1.5 s
   with the set lookup.
3. Success Metric: with a fixture corpus of two Specs, the corpus sweep
   requested three times in one process runs exactly two Checks and returns
   the same sweep each time.
4. Success Metric: with the active Specs of the audited tree, one run of the
   package pays the corpus sweep in its first corpus test only, and
   `TestStageScopeDefaultSweepIsUnchanged` pays one `CheckStage` per active
   Spec. Measured on 2026-10-22 with five active Specs (TechSpec →
   Measurements), the package took 17.5 s of wall time and 16.9 s of user
   CPU, against 316.5 s to 631.1 s and 504.4 s to 635.1 s today.
5. Success Metric: the corpus golden counts, `TestCheckActiveCorpusHasNoErrors`,
   `TestStageScopeDefaultSweepIsUnchanged`, every existing `internal/spec`
   History Citation test and every `internal/speccheck` `SC-HISTORY-RETAINED`
   test pass unchanged, and each sabotage the QA Task lists makes its named
   test fail.
6. Success Metric: the `ci-verify` job of the first pull request delivered
   after this Spec spends, in its two `internal/docscontract` runs together,
   less than a fifth of the 1,305 s the adopted Backlog Entry measured with
   five active Specs, for the same number of active Specs or scaled per
   active Spec.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted Backlog Entry, committed under `references/` before this Spec
  existed: the operator's reading of the CI runs of pull request #619 on
  2026-10-21, about 764 s and 541 s for the package with five active Specs.
  The QA compares the per-Check cost it measures with that baseline only as
  a ratio, never as an absolute CI number.
- The Go standard library documentation of `testing.AllocsPerRun`
  (https://pkg.go.dev/testing#AllocsPerRun): "AllocsPerRun returns the
  average number of allocations during calls to f." and "AllocsPerRun sets
  runtime.GOMAXPROCS to 1 during its measurement and will restore it before
  returning." Its source in Go 1.26 panics with "testing: AllocsPerRun called
  during parallel test", which is why the new `internal/spec` test is a
  Sequential Test.

## Unreachable Acceptance

- criterion: Success Metric 6, the `ci-verify` job of the first pull request
  delivered after this Spec spends, in its two `internal/docscontract` runs
  together, less than a fifth of the 1,305 s the adopted Backlog Entry
  measured with five active Specs.
  reason: only GitHub Actions runs the `ci-verify` job, and reading its step
  timings needs the network and an open pull request, which no Verification
  or QA row may reach.
  satisfied-by: the operator's reading of the "Verify changed" and "Verify
  docs" step logs of that pull request's `ci-verify` run after this Spec
  merges.

## Validation

The decisions were checked on 2026-10-22 against these sources.

- Characterization at `42f523c2` in scratch clones, every `go test` with
  `-timeout`: the package took 104.9 s with two active Specs, of which the
  four corpus tests took 95.3 s and `TestCheckLoopOrderRepositoryAgrees`,
  one Check against a one-Spec carrier corpus, 9.0 s. A CPU profile put 16.2 s
  of the 16.8 s of CPU under `speccheck.Check` in `detectHistoryRetained`,
  11.1 s of it in `spec.HistoryCitationsInSource`, dominated by
  `runtime.concatstring2`. A logging `git` wrapper on `PATH` counted 31 Git
  processes per Check, 0.16 s in all. The TechSpec records the code read, the
  prototypes and the measurement matrix.
- Secondbrain, read only. Read `wiki/index.md`; ran `qmd query "roundfix
  docscontract spec check corpus performance slow test" --all --files
  --min-score 0.3`, which returned the Roundfix mirror's records of Spec 0080
  (cheap detectors before the gate), Spec 0166 (docs changes run the tests
  that read them) and Spec 0071 (Verification cost). None
  records a measurement of the corpus sweep or a decision against sharing
  it. No Inbox Entry was written, because the briefing forbids writing there.
- Exa: the Go `testing` documentation and the Go 1.26 source of
  `AllocsPerRun` above. Challenging: an allocation count is process-wide, so
  a parallel test could add noise, but `AllocsPerRun` refuses to run during a
  parallel test, which makes the measurement exact for a Sequential Test.

## Glossary

None.

## Decisions

- 2026-10-22 (questions of 15:25 BRT, answered by about 18:30 BRT): the cycle
  AP scope "Ruído C2 e suite docscontract", the order "Correções, follow-up,
  Python (Recomendado)", and the authorization "the same as cycle AO"; the
  operator's briefing for this Spec says to fix inside Go code only, to keep
  every assertion, and to record a Makefile or workflow fix as the follow-up.
- 2026-10-22 (about 00:55 BRT): the cycle AO answer "Backlog por enquanto
  (Recomendado)" that kept the adopted entry as a Backlog Entry until this
  cycle.

Autonomous route: no person present; decisions taken from the operator's cycle AP briefing for Spec 0305 of 2026-10-22 and the measurements recorded in the TechSpec.

Taken from those sources, with no further interview:

- The fix is the set lookup in `HistoryCitationsInSource`, where the profile
  put the cost, plus one shared corpus sweep in the test package; the
  quadratic `strings.TrimLeft` copy in `CitationRewrite.paths` measured no
  difference on this repository and is left alone.
- The cost proof counts allocations, which do not depend on machine load,
  instead of timing the call.
- No ADR: neither change alters a contract a later reader must know.

## Open Questions

None.
