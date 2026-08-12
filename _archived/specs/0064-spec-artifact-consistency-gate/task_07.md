---
task: task_07
spec: 0064-spec-artifact-consistency-gate
status: completed
type: docs
complexity: medium
---

# Task 07: Bring this repository's own Specs to a clean report

## Overview

Run the check across this repository's active Specs, correct every reported
`error`, and add the test that holds them at zero from here on. This is the
step that makes a fail-closed gate possible: wiring the gate before the Specs
are clean would turn `make verify` red for every contributor.

The output is a declared break list, not a discovered one. Every correction is
a real artifact contradiction the check located, with both sides named — never
a rewrite for style.

## Requirements

1. MUST run the check across every active Spec and correct each reported
   `error` in the Spec artifacts, so the active corpus reports zero errors.
2. MUST leave archived Specs byte-identical. An archived Spec's findings are
   recorded by the corpus golden, never fixed.
3. MUST add a test that runs the check across every active Spec and fails,
   naming each Spec and code, when any `error` remains. The test lands with the
   corrections it holds in place.
4. MUST record each correction in the Task's Result as a declared break: the
   Spec, the code, and what changed on each side.
5. MUST NOT suppress, downgrade, or narrow a detector to make a Spec pass. A
   detector that over-reaches is a defect in the detector, and correcting it
   there is in scope; silencing it is not.
6. MUST leave every reported `gap` visible, dismissing one only by writing its
   reason into the Spec that carries it.
7. SHOULD keep each correction minimal — the smallest edit that removes the
   contradiction the check located.

## Subtasks

- [ ] Run the check across the active corpus and capture the report.
- [ ] Correct each reported error in the owning Spec artifact.
- [ ] Resolve or reason-dismiss each reported gap.
- [ ] Add the active-corpus zero-error test.
- [ ] Update the corpus golden's active-Spec counts to zero.
- [ ] Record every correction as a declared break in the Result.

## Acceptance Criteria

- [ ] The check reports zero `error` findings across every active Spec.
- [ ] A test asserts that and fails naming the Spec and code when it regresses.
- [ ] Every archived Spec file is byte-identical to its pre-Task content.
- [ ] No detector was disabled, narrowed, or made advisory to reach zero,
      asserted by the unchanged detector tests from tasks 01 through 03.
- [ ] Each reported gap is either resolved or carries a written reason in its
      Spec.
- [ ] The Result lists every correction with its Spec, code, and both sides.

## Context

- instruction: `docs/agents/spec-routing.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix spec check` — expected: exit 0; the
  active corpus reports no error.
- `go test ./internal/speccheck -count=1 -run 'ActiveCorpus' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the active-corpus test ran and passed.
- `go test ./internal/speccheck -count=1` — expected: exit 0; the detector
  tests from tasks 01 through 03 still pass unchanged.
- `git diff --name-only HEAD -- docs/specs/_archived | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no archived Spec file changed.

## References

- `_prd.md` → Success Metrics; Decisions (non-regression).
- `_techspec.md` → Build Order 7; Risks & Considerations.
- ADR-0094.

## Result

### Implementation

- Added `TestCheckActiveCorpusHasNoErrors` to the existing repository-corpus
  suite. It discovers active Specs through `spec.ListActive`, checks each
  through the public `speccheck.Check` API, and reports every remaining error
  as `<Spec>: <SC-* code>: <summary>`. Gap findings remain visible in verbose
  output.
- Corrected every error and accounted for every gap from the pre-change active
  corpus report without changing detector code or archived Spec artifacts.
- Updated the active half of `corpus-golden.json` to zero for every diagnostic
  code. Archived counts remain unchanged.

### Declared error breaks

| Spec | Code | Reported side | Corrected side |
| --- | --- | --- | --- |
| `0064-spec-artifact-consistency-gate` | `SC-ADR-UNLISTED` (`ADR-0039`) | The TechSpec characterization describes the model-availability replay relationship. | The PRD and TechSpec Active ADR rows now classify it as a characterization-corpus input, not a product-behavior obligation. |
| `0064-spec-artifact-consistency-gate` | `SC-ADR-UNLISTED` (`ADR-0049`) | The TechSpec characterization describes the atomic-profile replay relationship. | The PRD and TechSpec Active ADR rows now classify it as a characterization-corpus input, not a product-behavior obligation. |
| `0064-spec-artifact-consistency-gate` | `SC-ADR-UNLISTED` (`ADR-0055`) | The TechSpec characterization names the exact-capability relation candidate. | The PRD and TechSpec Active ADR rows now classify it as a characterization-corpus input, not a product-behavior obligation. |
| `0064-spec-artifact-consistency-gate` | `SC-ADR-UNLISTED` (`ADR-0084`) | The TechSpec characterization names the trusted-publishing replay decision. | The PRD and TechSpec Active ADR rows now classify it as a characterization-corpus input, not a product-behavior obligation. |
| `0064-spec-artifact-consistency-gate` | `SC-ADR-UNLISTED` (`ADR-0086`) | The TechSpec characterization names the declared-category-removal replay decision. | The PRD and TechSpec Active ADR rows now classify it as a characterization-corpus input, not a product-behavior obligation. |
| `0064-spec-artifact-consistency-gate` | `SC-ADR-UNLISTED` (ADR `0081`) | The PRD Tooling authority row claimed digest fallout although this Spec's authorized tooling side is only the `Makefile`. | The stale digest-fallout citation was removed; the Active ADR row correctly continues to omit an obligation this Spec does not exercise. |
| `0065-loop-order-and-verification-honesty` | `SC-ADR-UNLISTED` (ADR `0081`) | The Tooling authority row cites sanctioned digest fallout for the owned-skill edit. | The Active ADR row now records ADR 0081 as governing that fallout. |
| `0073-skill-versions-decoupled-from-the-binary` | `SC-TOOLING-UNAUTHORIZED` | The PRD cited the queued-Spec authorization, whose scope does not name Spec 0073. | The PRD now states that authorization is not recorded, retains the exact bounded files, and requires a naming grant before implementation; the authorization record was not widened. |
| `0075-typed-docs-backlog` | `SC-CONSTRAINT-UNREASONED` | The TechSpec declared Authentication and HTTP not applicable without a reason. | The row now states that the work is local documentation and Baseline asset content, opens no transport, and handles no credential. |

### Gap accounting

| Spec | Candidate | Disposition written into the Spec |
| --- | --- | --- |
| `0059-run-storage-compaction-and-global-sanitation` | ADR `0053` | Relation-only: terminal Run Worktree reconciliation remains owned by Spec 0038 and out of scope. |
| `0065-loop-order-and-verification-honesty` | `ADR-0091` | Applicable: it owns the authored QA gate Task whose position changes. |
| `0065-loop-order-and-verification-honesty` | `ADR-0093` | Relation-only: the Spec Consistency Check citation boundary does not govern loop order or Task Verification. |
| `0066-run-teardown-reclaims-what-it-created` | ADR `0053` | Applicable: its proof-based reconciliation contract extends to accumulated failed-cycle Run Branches. |
| `0068-spec-close-audit` | ADR `0053` | Applicable: its missing-target reconciliation case is resolved by content. |
| `0069-review-run-targets-its-pull-request` | ADR `0053` | Relation-only: Run Worktree reconciliation belongs to spec Runs; this Spec changes review Run target validation. |
| `0070-declared-unreachable-acceptance` | `ADR-0091` | Applicable: it owns the authored QA gate Task whose report carries the affected rows. |
| `0070-declared-unreachable-acceptance` | `ADR-0093` | Relation-only: the citation boundary does not govern QA evidence or archive semantics. |
| `0064-spec-artifact-consistency-gate` | `ADR-0040` | Relation-only candidate exposed by accounting for the characterization inputs; its reasoning policy does not govern this check. |
| `0064-spec-artifact-consistency-gate` | `ADR-0079` | Relation-only candidate exposed by accounting for the characterization inputs; its model-identifier policy does not govern this check. |

### Focused-check evidence

- Pre-change red signal:
  `GOCACHE=/private/tmp/roundfix-task07-gocache GOFLAGS=-buildvcs=false go test ./internal/speccheck -count=1 -run '^TestCheckActiveCorpusHasNoErrors$' -v`
  failed and named all nine errors by Spec and code, while logging all eight
  initial gaps.
- After the artifact corrections, the same focused test passed and emitted no
  remaining gap log. This exercises every active Spec through the public check
  API without running the Task's declared Verification command.
- `GOCACHE=/private/tmp/roundfix-task07-gocache GOFLAGS=-buildvcs=false go test ./internal/speccheck -count=1 -run '^TestCheckCorpusGoldenAndBudget$' -v`
  passed in 0.65 seconds. The checked-in active counts are all zero and the
  archived counts still match their pre-Task characterization.
- `git diff --name-only -- docs/specs/_archived` produced no paths.
- `git diff -- internal/speccheck/constraints.go internal/speccheck/citations.go internal/speccheck/vocabulary.go`
  produced no diff; no detector was disabled, narrowed, or downgraded.

### Acceptance-criterion evidence

- Active-corpus zero errors: the focused active-corpus test passed after the
  last Spec edit.
- Regression names the Spec and code: the pre-change run failed with messages
  such as `0064-spec-artifact-consistency-gate: SC-ADR-UNLISTED` and
  `0073-skill-versions-decoupled-from-the-binary: SC-TOOLING-UNAUTHORIZED`.
- Archived byte identity: the archived-path diff was empty and the archived
  golden counts did not change.
- Detector non-regression: detector source files are absent from this Task's
  diff; the Daemon retains ownership of the declared full detector suite.
- Gap accounting: every initial and correction-induced relation candidate is
  listed above with an applicable obligation or a written relation-only
  reason; the final focused run emitted no gap.
- Declared breaks: the table above names every error's Spec, code, reported
  side, and corrected side.

### Follow-ups

- Spec 0073 now records the real boundary: its protected `Makefile` and owned
  skill changes cannot start until a maintainer authorization record names the
  Spec and exact paths. This Task did not invent or widen that authority.
- The 0064 TechSpec's Build Order 7 still says archived Specs may be corrected,
  while its Testing Approach and this Task require archived Specs to remain
  byte-identical. That unreported planning contradiction was left unchanged
  because this Task permits only check-located artifact corrections.
