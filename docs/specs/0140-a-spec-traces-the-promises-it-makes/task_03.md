---
task: task_03
spec: 0140-a-spec-traces-the-promises-it-makes
status: completed
type: test
complexity: medium
---

# Task 03: Characterize the coined codes and the corpus

## Overview

Two coined codes need a glossary owner before they ship, and the repository's
corpus contract needs to accept them before it can count them. This slice
records the prediction this Spec made about the ten Specs in authoring, then
updates the corpus contract and its golden to what the sweep actually reports.

## Requirements

1. MUST give each coined code a glossary owner that states what it means and
   which artifact it reads.
2. MUST add both codes to the list the corpus sweep accepts, so that an emitted
   code is counted instead of failing the sweep as uncharacterized.
3. MUST record, in the Spec's QA evidence directory, the prediction this Spec
   authored before the sweep ran: ten metric findings and ten contract findings
   over the Specs in authoring, and no new error.
4. MUST update the golden map to the counts the sweep reports, and MUST NOT
   change any count this Spec did not cause.
5. MUST NOT edit any Spec in authoring to make its counts smaller.

## Subtasks

- [ ] Write the glossary entry for each coined code.
- [ ] Add both codes to the corpus contract's accepted-code list.
- [ ] Run the sweep, capture its counts, and compare them with the prediction.
- [ ] Update the golden map and record the comparison as evidence.

## Acceptance Criteria

- [ ] The glossary documents both codes, so the vocabulary detector runs over
      them and reports nothing.
- [ ] The corpus sweep accepts both codes instead of failing as uncharacterized.
- [ ] The golden map carries the two new counts, and the recorded evidence shows
      the prediction beside the observed counts.
- [ ] No Spec in authoring is edited by this Task.

## Context

- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `CONTEXT.md`

## Verification

- `grep -q "SC-METRIC-UNDECLARED" CONTEXT.md` — expected: exit 0; the coined code has a glossary owner.
- `grep -q "SC-CONTRACT-UNDECLARED" CONTEXT.md` — expected: exit 0.
- `grep -q "CodeMetricUndeclared" internal/docscontract/corpus_test.go` — expected: exit 0; the sweep accepts the code.
- `grep -q "CodeContractUndeclared" internal/docscontract/corpus_test.go` — expected: exit 0.
- `grep -q "SC-METRIC-UNDECLARED" internal/docscontract/testdata/corpus-golden.json && go test -count=1 -tags docscontract -run "^TestCheckCorpusGolden$" ./internal/docscontract` — expected: exit 0; the golden carries the new code and matches the sweep. Before this Task the key is absent, so the command fails.
- `grep -q "SC-METRIC-UNDECLARED" internal/docscontract/testdata/corpus-golden.json` — expected: exit 0; the golden carries the coined code as a key.
- `grep -q "SC-CONTRACT-UNDECLARED" internal/docscontract/testdata/corpus-golden.json` — expected: exit 0.
- `ls docs/specs/0140-a-spec-traces-the-promises-it-makes/qa/evidence` — expected: the recorded prediction and the observed counts exist as evidence.

## References

`_prd.md` → Goal 4; Acceptance evidence; Success Metric 1;
`_techspec.md` → Data Models; Testing Approach 3-4; Vocabulary Contract;
Build Order 3; ADR-0104.

## Result

Characterized the two promise-declaration codes and the active Spec corpus.
The glossary now defines each code and the artifact it reads, the corpus sweep
accepts both codes, and the golden records only the two observed counts added
by this Spec. The prediction and the later observation are recorded together
in `qa/evidence/task-03-corpus-sweep.md`.

Focused-check evidence:

- Before implementation, repository searches found neither code in
  `CONTEXT.md`, `corpusFindingCodes`, or the corpus golden.
- After the prediction was recorded, `rtk go run ./cmd/roundfix spec check
  --format json` exited 0 and returned 11 Spec documents. The ten pre-existing
  Specs 0120 through 0129 each emitted one `SC-METRIC-UNDECLARED` gap and one
  `SC-CONTRACT-UNDECLARED` gap; Spec 0140 emitted neither. The sweep reported
  20 gaps, zero errors, and no other finding code.
- `rtk go test -count=1 -tags docscontract -run
  '^(TestCheckCorpusBudget|TestCheckActiveCorpusHasNoErrors)$'
  ./internal/docscontract` passed both focused tests. This exercised the real
  corpus through the accepted-code list without invoking the Daemon-owned
  golden Verification command.
- `rtk jq empty internal/docscontract/testdata/corpus-golden.json` exited 0.

Acceptance evidence:

- `CONTEXT.md` defines `SC-METRIC-UNDECLARED` as the PRD declaration gap that
  reads `_prd.md`, and `SC-CONTRACT-UNDECLARED` as the TechSpec declaration gap
  that reads `_techspec.md`. The public-boundary sweep emitted no
  `SC-VOCABULARY-UNDOCUMENTED` finding.
- `corpusFindingCodes` includes `CodeMetricUndeclared` and
  `CodeContractUndeclared`; the focused corpus-budget sweep accepted and
  counted both codes.
- The golden adds only `SC-METRIC-UNDECLARED: 10` and
  `SC-CONTRACT-UNDECLARED: 10`; the evidence table places the predicted and
  observed values side by side with a zero difference and zero new errors.
- The changed-path review names no file under Specs 0120 through 0129. The
  sweep was read-only over those authoring Specs.

The first public-boundary attempt could not access the default Go build cache
inside the sandbox and did not reach the sweep. The permitted rerun reached the
boundary and produced the observation above. The Daemon-owned Verification
commands were not run in this turn.

Follow-up: `rtk make verify-incremental` reached the repository test suite but
exited 2. Synthetic QA fixtures in `internal/cli` and `internal/taskengine`
received one new promise-declaration finding, so tests that expected a clean QA
Run instead observed `QA verdict: fail`; examples include
`TestRunImplementUsesConfiguredExternalSpecRootEndToEnd` and
`TestTaskCycleQAPromptStaysUsableWithoutRecordedTargetBranch`. Those fixture
declarations belong to the detector/fixture implementation slice, not this
Task's glossary and active-corpus characterization, so this diff does not alter
them.

## Carry-forward provenance

- Source Run: `run_20260917T164016Z_ac67737203a1aa3d`
- Source commit: `94df5f4eb9f589ac3d43c6f4aa5ba1cdbb690b82`
