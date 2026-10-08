---
task: task_01
spec: 0253-authoring-rules-that-stop-qa-reruns
status: completed
type: backend
complexity: medium
---

# Task 01: An outside-evidence source is reachable, a transcript is asserted whole, and a prerequisite is declared

## Overview

Four of the last eleven Specs needed a QA Archive Override because an
outside-evidence row rested on a source the QA gate could not reach, and five
QA reruns came from transcript lines the implementing test never asserted
(findings B08 and B09 and briefing rules 8 and 13 of the Baseline audit of
2026-10-08, `docs/references/2026-10-08-baseline-audit.md`). This Task extends
three `spec-workflow` clauses in place with `_techspec.md` → Clause changes
(ADR-0258). It is verifiable on its own: this repository's spec-routing guide
carries the three sentences, and an adopter's update retains the three
clauses.

This is an authorized tooling Task. It may change only the files in its
Context, the derived files the sanctioned regeneration rewrites, and this Task
file. It starts from the merge of Spec 0252, which changes the same module.

## Requirements

1. MUST answer findings B08 and B09 of the Baseline audit of 2026-10-08 by
   editing the `guidance` of `clause.spec.project-constraints-06-outside-evidence`,
   `clause.spec.routing-05-task-graph` and
   `clause.spec.sources-02-bound-the-group` in
   `internal/baseline/assets/modules/spec-workflow.json` exactly as
   `_techspec.md` → Clause changes (task_01) says. It edits the strings in
   place, so every other byte stays as it is. Each clause MUST keep its `id`
   and its `mandatory` enforcement and MUST NOT gain a `replaces` list. The
   network-denied sentence of the outside-evidence clause MUST stay
   byte-identical.
2. MUST raise by one, from the value on this Task's starting commit, the
   task_01 versions of `_techspec.md` → Version changes. It MUST then run
   `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1`,
   then `make baseline-digests`, then
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --yes --format text`.
   A second refresh MUST report `File changes: 0`. It MUST NOT hand-edit a
   module version, pin, golden, snapshot or generated guide, and the Result
   MUST name every file the commands rewrote.
3. MUST create `internal/baseline/authoring_evidence_clauses_test.go` with
   `TestTheAuthoringEvidenceClausesCarryTheirText`,
   `TestTheAuthoringEvidenceClausesRenderInTheGuide` and
   `TestAnAdopterRetainsTheAuthoringEvidenceClauses`, following
   `internal/baseline/glossary_clauses_test.go`. The three added texts are
   literals in the test. The text test MUST require each literal exactly once
   in its clause, with `mandatory` enforcement and no `replaces`. The
   rendering test MUST require each literal exactly once, on
   whitespace-normalized text, in the Standard TypeScript golden
   `docs/agents/spec-routing.md` and in this repository's
   `docs/agents/spec-routing.md`. The retention test MUST require a ready
   Managed Refresh of a Source Baseline adopter that records the three clauses
   `retained` and no clause `unaccounted`.
4. MUST change the `clause.spec.sources-02-bound-the-group` literal in
   `internal/baseline/grouped_sources_clauses_test.go` to the new text and
   nothing else in that file. That is a declared break.
5. MUST prove each new gate can fail. The Result MUST record one sabotage of a
   clause text and one of the retention, for example a changed enforcement,
   with the test that failed for each, and that the source was restored and
   regenerated.
6. MUST NOT change any other clause, the Source Baseline assets, the retention
   transition, the force record or any production Go file.

## Subtasks

- [ ] Extend the three clauses and raise the versions.
- [ ] Record, regenerate and refresh twice.
- [ ] Add the clause test file and update the grouping literal.
- [ ] Record each sabotage.

## Acceptance Criteria

- [ ] `docs/agents/spec-routing.md` carries the reachable-source sentence, the
      whole-transcript sentence and the `requires` sentence, each once.
- [ ] A Source Baseline adopter's refresh retains the three clauses with no
      clause `unaccounted`.
- [ ] A second refresh of this repository is a no-op.

## Context

- instruction: `docs/adr/0258-a-spec-is-authored-against-the-qa-rerun-classes-and-a-retirement-writes-reduced-history.md`
- instruction: `docs/adr/0257-inside-a-run-the-daemon-is-the-only-full-gate-and-the-baseline-states-each-rule-once.md`
- instruction: `docs/adr/0250-a-module-version-is-chosen-when-recorded-and-the-coverage-record-lists-every-platform.md`
- instruction: `docs/references/2026-10-08-baseline-audit.md`
- instruction: `internal/baseline/glossary_clauses_test.go`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/module-versions.json`
- interface: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- interface: `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- interface: `internal/baseline/testdata/catalog.diagnostics.golden.json`
- interface: `internal/baseline/testdata/catalog.digest`
- interface: `internal/baseline/testdata/catalog.normalized.json`
- interface: `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- interface: `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`
- interface: `internal/baseline/grouped_sources_clauses_test.go`
- interface: `docs/agents/spec-routing.md`
- interface: `docs/agents/setup-context.json`
- creates: `internal/baseline/authoring_evidence_clauses_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestTheAuthoringEvidenceClausesCarryTheirText|TestTheAuthoringEvidenceClausesRenderInTheGuide|TestAnAdopterRetainsTheAuthoringEvidenceClauses|TestTheGroupingClausesCarryTheirForceAndText|TestTheGroupingClausesRenderInTheGuides|TestAnAdopterRetainsTheGroupingClauses|TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow|TestBaselineClauseForceIsCharacterized|TestNoTwoBaselineClausesShareText|TestEveryBaselineModuleVersionIsRecorded|TestCatalogCompatibility|TestFormatterComposition|TestBaselinePlanCharacterization)$" ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheAuthoringEvidenceClausesCarryTheirText TestTheAuthoringEvidenceClausesRenderInTheGuide TestAnAdopterRetainsTheAuthoringEvidenceClauses TestTheGroupingClausesCarryTheirForceAndText TestTheGroupingClausesRenderInTheGuides TestTheGuidesExemptANetworkDeniedOutsideEvidenceRow TestBaselineClauseForceIsCharacterized TestNoTwoBaselineClausesShareText TestEveryBaselineModuleVersionIsRecorded TestCatalogCompatibility; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the three authoring-evidence tests do not exist, so the command fails.
- `for phrase in "Name a source the QA gate can check without the network before the Run starts" "asserts the whole transcript in one named test" "state the dependency in the TechSpec Build Order"; do tr -s '[:space:]' ' ' < docs/agents/spec-routing.md | grep -qF -- "$phrase" || { printf 'missing phrase in docs/agents/spec-routing.md: %s\n' "$phrase" >&2; exit 1; }; done; plan="$(go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --format json)" || { printf '%s\n' "$plan"; exit 1; }; printf '%s\n' "$plan" | grep -qF -- '"state":"current"' || { printf 'the guides are not refreshed: %s\n' "$plan" >&2; exit 1; }` — expected: exit 0; before this Task the spec-routing guide lacks the reachable-source sentence, so the command fails at its first phrase.

## References

- [_prd.md](_prd.md) — Goals 1-3; User Stories 1-3; Core Features 1 and 2; Success Metric 1; Acceptance evidence
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; Clause changes; Version changes; Retention; Derived files; Invariants 1-6 and 8; Testing Approach; Build Order 1
- ADR-0258; ADR-0257; ADR-0250; ADR-0186

## Result

Implemented the task_01 authoring guidance in place: outside evidence names a
source reachable without the network, an implementing Task asserts the whole
Surface Transcript in one named test, and an unmerged prerequisite is declared
in `requires` and the TechSpec Build Order. The three clause identities remain
mandatory with no `replaces`; the existing network-denied sentence is unchanged.
The only change to `grouped_sources_clauses_test.go` is the declared grouping
literal extension.

### Implementation and focused checks

- Starting worktree: only this Task file was modified, with the Daemon's
  `status: in_progress`; that status was preserved.
- Added `TestTheAuthoringEvidenceClausesCarryTheirText`,
  `TestTheAuthoringEvidenceClausesRenderInTheGuide`, and
  `TestAnAdopterRetainsTheAuthoringEvidenceClauses`. Their three literals are
  independent authored expectations. The text test requires each literal once,
  mandatory enforcement and no replacements; rendering checks both guides after
  whitespace normalization; retention checks a ready Source Baseline Managed
  Refresh, all three retained dispositions and evidence, and no unaccounted
  clause.
- Before the module edit, `go test ./internal/baseline -run
  '^TestTheAuthoringEvidenceClauses(CarryTheirText|RenderInTheGuide)$' -count=1`
  exited 1: all three clause literals were absent and each was absent in both
  guides. The first attempt could not access the default Go cache; all later
  Go and regeneration commands used
  `GOCACHE=/tmp/roundfix-0253-task01-gocache`.
- Raised `rule.spec.routing` 8 to 9, `rule.spec.project-constraints` 5 to 6,
  and `guide.spec-routing` 13 to 14 from the starting commit. Ran
  `go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$'
  -record-module-versions -count=1`: exit 0. The recorder selected module
  version 18 (starting version 17).
- Ran `make baseline-digests`: exit 0, `changed: true`. Ran it again after
  restoring the sabotage edits: exit 0, `changed: false`.
- Ran `go run -buildvcs=false ./cmd/roundfix baseline update --repo .
  --no-skills --yes --format text`. Its first sandboxed attempt failed before
  application because the Git-private transaction directory was not writable.
  The authorized rerun with sandbox escalation exited 0, reported
  `Baseline update: verified`, and applied two file changes. The next invocation
  exited 0, reported `File changes: 0` and `Idempotence: verified`.
- Final focused command (exit 0; all six named tests reported PASS):

  ```sh
  GOCACHE=/tmp/roundfix-0253-task01-gocache go test ./internal/baseline -run '^(TestTheAuthoringEvidenceClausesCarryTheirText|TestTheAuthoringEvidenceClausesRenderInTheGuide|TestAnAdopterRetainsTheAuthoringEvidenceClauses|TestTheGroupingClausesCarryTheirForceAndText|TestTheGroupingClausesRenderInTheGuides|TestAnAdopterRetainsTheGroupingClauses)$' -count=1 -v
  ```

### Acceptance evidence

1. Both the Standard TypeScript golden and this repository's
   `docs/agents/spec-routing.md` contain each added literal exactly once:
   `TestTheAuthoringEvidenceClausesRenderInTheGuide` passed all six guide
   subtests. The text/force test passed all three clause subtests.
2. A temporary Source Baseline adopter's Managed Refresh is ready, retains
   all three clauses with matching retention evidence, and has no unaccounted
   clause: `TestAnAdopterRetainsTheAuthoringEvidenceClauses` passed.
3. The second repository refresh reports zero file changes and verified
   idempotence, as recorded above.

### Negative evidence and restoration

Temporarily injected two subject mutations into the new test source, using
`mutateCatalogClause` on the loaded test catalog. Neither assertion was changed,
and no production source, Source Baseline asset or generated expectation was
mutated for these sabotages.

- Clause-text sabotage: changed `asserts the whole transcript` to
  `asserts part of the transcript` in the loaded
  `clause.spec.routing-05-task-graph`. The focused
  `TestTheAuthoringEvidenceClausesCarryTheirText` failed its Task Graph clause
  subtest with `force/text differs`.
- Retention sabotage: after constructing the adopter, changed that target
  clause's enforcement from `mandatory` to `stop-and-ask`.
  `TestAnAdopterRetainsTheAuthoringEvidenceClauses` failed with
  `refresh is not ready`, `State:action_required`, and
  `retention transition has 1 unaccounted clause(s):
  clause.spec.routing-05-task-graph`.
- Command for the two sabotages:
  `go test ./internal/baseline -run
  '^(TestTheAuthoringEvidenceClausesCarryTheirText|TestAnAdopterRetainsTheAuthoringEvidenceClauses)$'
  -count=1` exited 1, with both expected named failures.
- The rendering gate's negative evidence is the starting run above: every
  literal was missing in both rendered guides, and all six rendering subtests
  failed with `guide contains literal 0 times, want exactly once`.
- Restored the new test source byte-identically in a `finally` block, ran the
  sanctioned digest regeneration again, refreshed this repository, and then
  obtained the final six-test pass recorded above. No sabotage remains.

### Files rewritten by sanctioned commands

Module Version Record:

- `internal/baseline/assets/modules/spec-workflow.json` (module version line)
- `internal/baseline/module-versions.json`

`make baseline-digests`:

- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/spec-routing.md`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- `internal/baseline/testdata/plan-characterization/advisory-only-divergences.golden.json`
- `internal/baseline/testdata/plan-characterization/clean-adoption.golden.json`
- `internal/baseline/testdata/plan-characterization/idempotent-replan-after-verified-apply.golden.json`
- `internal/baseline/testdata/plan-characterization/same-baseline-changed-profile-and-catalog-digests.golden.json`

Repository Managed Refresh:

- `docs/agents/spec-routing.md`
- `docs/agents/setup-context.json`

Byte comparison against `HEAD` confirms the module has exactly the three
specified guidance edits, three inner version increments and recorder-written
module version, and the grouping test has only its specified literal change.
Changed-file postflight covers only Task Context paths plus this Task and its
new test. The Source Baseline assets, retention transition, force record,
production Go files, other Tasks and Task Graph manifest have no diff.

The Task's declared Verification and both repository Verification tiers were
not run; they remain Daemon-owned. No commit, push or Pull Request was made.
No follow-up work was found in this slice.
