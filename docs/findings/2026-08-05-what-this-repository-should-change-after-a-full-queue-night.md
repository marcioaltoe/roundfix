---
date: 2026-08-05
surface: internal, Makefile, docs/agents
status: open
---

# What this repository should change, after a full queue night

One session drove five Specs from PRD to archive and stopped a sixth on a
maintainer decision. Thirty-three implementation Tasks settled on Verification
attempt 1; exactly one failed on the work. That ratio is the finding's premise:
**the implementation half of this repository is healthy, and its costs are
elsewhere.** This records where, with what it cost, and what would change it.

Two patterns already have their own findings and are referenced rather than
repeated:
[an absence grep rejects the work it was written to protect](2026-08-05-an-absence-grep-rejects-the-work-it-was-written-to-protect.md)
and
[a finding does not prevent the recurrence it describes](2026-08-04-a-finding-does-not-prevent-the-recurrence-it-describes.md).

## 1. Timing budgets written as bare literals

**Seven test failures, in seven different tests, across one night** — most of
them on changes that touched no Go source at all. Every one passed in isolation
and passed on rerun.

| Test | Where it burned |
| --- | --- |
| `TestRunImplementDetachSurvivesCallerProcessGroupKill` | CI on PR #98, a 14-file Markdown change; again on Spec 0076's gate |
| `TestCheckCorpusGoldenAndBudget` | Spec 0064's gate: 0.64 s alone, 3.0–4.0 s under `-parallel 16` |
| `TestBaselinePlanCharacterization` and four siblings | Spec 0070's task_04 |
| `TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree` | Spec 0068's planning gate |
| `TestProveExactSelectionMalformedEvidenceCleanup` | Spec 0068's archive PR, which branch protection then correctly refused to merge |

The principle is already written down, in the same file as two of the
offenders:

> `implementWaitBudget` — *A tight budget here measures how loaded the machine
> is, not whether the code works.*

Two literals four hundred lines below it violated it. They were named in #104.
The rest remain.

**Suggestion.** A wait budget that is never reached costs nothing, because the
timer is only read on failure — so generous is free and tight is a coin flip.
Add a check that fails on a bare duration literal in a wait, timeout, or
deadline position inside `_test.go` files, with an allow-list for the budgets
that are the *subject* of an assertion rather than scaffolding around one.
That distinction matters: Spec 0064's corpus budget was a real product
guarantee and was correctly repaired by moving it to a serial gate step, not by
widening it.

## 2. Derivation ownership modelled per directory

Spec 0067 built ownership records to end three years of inference about which
artifacts `make baseline-digests` owns. Its own gate then found the model too
coarse: `internal/baseline/testdata/parity-corpus/` holds **fifteen files**, of
which the sanctioned command rewrites exactly **two** — `v1/manifest.json` and
`v1/fixtures/asset-sync.json`. The other thirteen, including `blobs.json` and
`matrix.json`, it never touches.

A per-directory `owner` cannot express that, so whichever value it carries is
false about most of the directory. Both prior readings — "the PRD is wrong" and
"the command is wrong" — were themselves wrong; the directory is simply mixed.

**Suggestion.** Ownership needs per-path exceptions, not per-directory verdicts.
Spec 0067 is the natural home and its work is already done except for this.

## 3. A sanctioned command that cannot reach what it claims

`make baseline-digests` reports `changed: false` and exits 0 while two
characterization corpora it does not cover are stale, leaving `make verify`
red. Their regeneration flags **do not match their test names**:

| Test | Flag |
| --- | --- |
| `TestBaselinePlanCharacterization` | `-update-baseline-plan-characterization` |
| `TestCatalogDiagnosticCharacterization` | `-update-catalog-diagnostics` |

The second was guessed from the test name during this session and the guess was
wrong — the fifth occurrence of a class first recorded on 2026-08-01.

**Suggestion.** Beyond Spec 0067's coverage fix: make each flag derivable from
its test name, or make the failure message print the exact invocation. A name
that must be remembered will be guessed.

## 4. Pinned values duplicated as literals

`const wantUpstreamDigest` appeared at **three** call sites in
`skills/baseline_skill_contract_test.go`. The repository's own rule names this
exact defect:

> *an assertion reads the constant it means* — when a value must be duplicated,
> change every occurrence in the same commit; fixing one of three is the most
> repeated defect in this repository's history.

Three copies, and the rule's own example is "one of three". They were hoisted to
one package-level constant in #100.

**Suggestion.** The rule is right and was not enough. A check that fails on the
same 64-hex literal appearing more than once outside a single `const` block
would make it mechanical.

## 5. Environment-dependent probes inside the required gate

`TestProveExactSelectionMalformedEvidenceCleanup` shells out to `npx` to prove
adapter package lineage. On PR #118 — an archive PR moving Markdown — it
returned an `AdapterLineageError` where the test expected a
`CapabilityEvidenceError`, failing the gate. It passed locally 2/2 and had
passed three times on the same code minutes earlier.

**Suggestion.** A required gate should not depend on a package registry being
reachable. Either stub the lineage probe in the gate tier and prove the real
one in a dedicated step, or classify the network failure explicitly so it
reports as unavailable rather than as a wrong error type.

## 6. Authorization wording that invites re-litigation

The 2026-08-02 record said a Baseline module asset is one the grant "does not
obviously reach". That phrasing was read three separate times this session as
implying module assets are protected tooling, and it cost a maintainer question
each time. It was settled on 2026-08-05: they are product content.

**Suggestion.** An authorization record should state what a class **is**, not
what a grant does not obviously reach. The settled classification is now
recorded; the pattern is worth watching for in future records.

## Evidence

- Session of 2026-08-04 into 2026-08-05: Specs 0064, 0076, 0070, 0068 and 0066
  archived; 0067 stopped on a maintainer decision; 21 Pull Requests merged.
- Task settlement: 33 implementation Tasks, 1 failure attributable to the work.
- `docs/handoffs/2026-08-05-the-night-the-queue-moved.md`.
- Pull Requests #98, #100, #104, #113, #114, #118 for the specific cases above.
