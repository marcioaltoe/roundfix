---
task: task_03
spec: 0140-a-spec-traces-the-promises-it-makes
status: pending
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
