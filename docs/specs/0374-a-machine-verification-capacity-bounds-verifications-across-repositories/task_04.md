---
task: task_04
spec: 0374-a-machine-verification-capacity-bounds-verifications-across-repositories
status: completed
type: docs
complexity: low
---

# Task 04: The glossary defines the Machine Verification Capacity, and ADR-0056 names its amendment

## Overview

`CONTEXT.md` adds **Machine Verification Capacity** and changes
**Verification Capacity**, **Waiting for Verification** and **Run Budget**,
so the vocabulary every guide written by tasks 02 and 03 uses is defined.
ADR-0056 carries the notice that ADR-0384 wraps its per-Run capacity in a
machine capacity.

## Requirements

1. MUST, following `domain-modeling`, add to `CONTEXT.md` right after the
   **Verification Capacity** entry the entry **Machine Verification Capacity**
   with the definition and `_Avoid_` line of `_techspec.md` → API Contract 10
   word for word, its closing `(ADR-0384)` written as
   `([ADR-0384](docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md))`.
2. MUST change, exactly as API Contract 10 says: **Verification Capacity**
   keeps its first sentence, gains "An attempt holds it while it waits for
   the Machine Verification Capacity and gives it back last." and its
   `_Avoid_` line becomes
   `_Avoid_: Task Capacity, worktree concurrency, Machine Verification Capacity`;
   in **Waiting for Verification**, "before the Task acquires Verification
   Capacity" becomes "before the Task holds its Run's Verification Capacity
   and the Machine Verification Capacity"; **Run Budget** gains "It does not
   run while a Verification of the Run waits for the Machine Verification
   Capacity." Every other entry keeps its bytes.
3. MUST write into `docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md`,
   after the paragraph that ends "projects must choose a safe Task Capacity
   for mutations outside the Daemon Verification boundary.", the paragraph
   `**Amended by [ADR-0384](0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md) (2026-11-10).** Every Verification attempt also takes one unit of the Machine Verification Capacity after the Run's own unit and gives it back first, so the Daemon Verifications of every Run and the Delivery Queue's repository gate on one machine no longer run together beyond runs.max_verifying. The per-Run capacity, its fairness and the exclusive retry still hold; the retry takes one machine unit.`
   and set its `updated_at` to `2026-11-10T23:30:00Z`. Its other bytes stay.
4. MUST run the glossary contracts named in the Verification and keep them
   passing. The Task changes no Governed Path.

## Subtasks

- [ ] Add **Machine Verification Capacity** and change the three entries.
- [ ] Write ADR-0056's amendment notice.
- [ ] Run the glossary contracts.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Machine Verification Capacity** and carries the three changed entries.
- [ ] ADR-0056 carries the amendment notice naming ADR-0384.
- [ ] The glossary contracts pass.

## Context

- instruction: `docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md`
- interface: `CONTEXT.md`
- interface: `docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md`

## Consumers

- The glossary entries and ADR-0056's notice — consumers: `./internal/docscontract`

## Verification

- `grep -q '^\*\*Machine Verification Capacity\*\*:$' CONTEXT.md || { echo 'missing term **Machine Verification Capacity**' >&2; exit 1; }; grep -q '^\*\*Verification Capacity\*\*:$' CONTEXT.md && grep -q '^\*\*Waiting for Verification\*\*:$' CONTEXT.md && grep -q '^\*\*Run Budget\*\*:$' CONTEXT.md || { echo 'missing changed terms **Verification Capacity**, **Waiting for Verification** or **Run Budget**' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'An attempt holds it while it waits for the Machine Verification Capacity and gives it back last.' || { echo 'missing: Verification Capacity change' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qE -- 'before the Task holds its Run.s Verification Capacity and the Machine Verification Capacity' || { echo 'missing: Waiting for Verification change' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'It does not run while a Verification of the Run waits for the Machine Verification Capacity.' || { echo 'missing: Run Budget change' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md | grep -qF -- 'Amended by [ADR-0384]' || { echo 'missing: ADR-0056 notice' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryAvoidLinesArePlainLists|TestGlossaryStatesEachRuleOnce)$'` — expected: exit 0. Before this Task **Machine Verification Capacity** is not in the glossary and ADR-0056 has no notice; afterwards the new and changed entries are present and the four glossary contracts pass.

## References

- `_prd.md` → Success Metric 7; Glossary
- `_techspec.md` → API Contract 10; Build Order 4
- ADR-0384
- ADR-0056

## Result

Skills loaded: domain-modeling

### What changed

- `CONTEXT.md`: added **Machine Verification Capacity** immediately after the
  **Verification Capacity** entry, with API Contract 10's definition and
  `_Avoid_: machine-wide test lock, core budget, Run ceiling`, closing with the
  glossary's accepted-ADR link `([ADR-0384](docs/adr/0384-a-machine-verification-capacity-bounds-the-verifications-of-every-repository.md))`.
  Changed exactly three entries: **Verification Capacity** appends "An attempt
  holds it while it waits for the Machine Verification Capacity and gives it
  back last." and its `_Avoid_` line becomes `_Avoid_: Task Capacity, worktree
  concurrency, Machine Verification Capacity`; **Waiting for Verification**
  replaces "before the Task acquires Verification Capacity" with "before the
  Task holds its Run's Verification Capacity and the Machine Verification
  Capacity"; **Run Budget** appends "It does not run while a Verification of
  the Run waits for the Machine Verification Capacity." No other entry's bytes
  changed.
- `docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md`: added
  the amendment paragraph after the paragraph ending "projects must choose a
  safe Task Capacity for mutations outside the Daemon Verification boundary.",
  naming ADR-0384, and set `updated_at` to `2026-11-10T23:30:00Z`. Its other
  bytes stay.

### Acceptance criteria evidence

- **`CONTEXT.md` defines Machine Verification Capacity and carries the three
  changed entries** — the first grep block names `**Machine Verification
  Capacity**:` and the three changed terms, and the normalized-whitespace greps
  find each added/replaced sentence: output `term ok` / `vc ok` / `wfv ok` /
  `rb ok`.
- **ADR-0056 carries the amendment notice naming ADR-0384** — normalized grep
  finds `Amended by [ADR-0384]`: output `adr ok`; the added paragraph sits at
  `docs/adr/0056-spec-runs-separate-task-and-verification-capacity.md:29`.
- **The glossary contracts pass** — focused check
  `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryAvoidLinesArePlainLists|TestGlossaryStatesEachRuleOnce)$'`
  → `ok roundfix/internal/docscontract 0.415s`. The whole package was also run
  green: `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract`
  → `ok ... 27.089s`.

### Notes

- No Governed Path changed; only `CONTEXT.md` and ADR-0056 were edited.
- The Daemon owns the Verification run and Task status; this Result does not
  claim the Task completed or failed.
