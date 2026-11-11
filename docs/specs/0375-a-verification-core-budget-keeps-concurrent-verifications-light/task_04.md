---
task: task_04
spec: 0375-a-verification-core-budget-keeps-concurrent-verifications-light
status: completed
type: docs
complexity: low
---

# Task 04: The glossary defines the Verification Core Budget

## Overview

`CONTEXT.md` gains **Verification Core Budget**, the term the guides and
references of tasks 02 and 03 use, right after the entry
**Machine Verification Capacity**, with the text of `_techspec.md` → API
Contract 7 and a link to ADR-0385.

## Requirements

1. MUST add the entry through `domain-modeling`, right after the entry
   **Machine Verification Capacity** in `CONTEXT.md`, as the line
   `**Verification Core Budget**:` followed by the definition of
   `_techspec.md` → API Contract 7 word for word, with `(ADR-0385)` written as
   the glossary's usual Markdown link to
   `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md`,
   and the final line
   `_Avoid_: core limit, CPU quota, Machine Verification Capacity`.
2. MUST change no other glossary entry.
3. MUST keep the glossary contracts passing:
   `TestGlossaryPhrasesAreDefinedTerms`, `TestGlossaryLinksNameAcceptedADRs`,
   `TestGlossaryAvoidLinesArePlainLists` and `TestGlossaryStatesEachRuleOnce`.

## Subtasks

- [ ] Write the entry through `domain-modeling`.
- [ ] Run the glossary contracts.

## Acceptance Criteria

- [ ] `CONTEXT.md` defines **Verification Core Budget** with the prescribed text and link, and the four glossary contracts pass.

## Context

- instruction: `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md`
- interface: `CONTEXT.md`

## Consumers

- The glossary entry — consumers: `./internal/docscontract`

## Verification

- `grep -F -- '**Verification Core Budget**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Verification Core Budget\*\*(:| —( |$))' || { echo 'missing term **Verification Core Budget**' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qE -- 'and each Agent Session of an Implement Run, as .GOMAXPROCS. and .ROUNDFIX_VERIFICATION_CORES.\.' || { echo 'missing: Verification Core Budget definition' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- "_Avoid_: core limit, CPU quota, Machine Verification Capacity" || { printf 'missing phrase in %s: %s\n' CONTEXT.md "_Avoid_: core limit, CPU quota, Machine Verification Capacity" >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryAvoidLinesArePlainLists|TestGlossaryStatesEachRuleOnce)$'` — expected: exit 0. Before this Task **Verification Core Budget** is not in the glossary; afterwards the entry is present and the four glossary contracts pass.

## References

- `_prd.md` → Glossary; Success Metric 5
- `_techspec.md` → API Contract 7; Glossary; Build Order 4
- ADR-0385

## Result

Skills loaded: `implement-task`, `domain-modeling`, `tech-writer`, `evidence-gate`

Implementation:
- `CONTEXT.md` now defines **Verification Core Budget** immediately after the **Machine Verification Capacity** entry, using the API Contract 7 text verbatim, with the `(ADR-0385)` link to `docs/adr/0385-each-verification-receives-a-verification-core-budget-in-its-environment.md` and the final `_Avoid_: core limit, CPU quota, Machine Verification Capacity` line. No other glossary entry changed.

Focused checks:
- Pre-change: `grep -n 'Verification Core Budget\|core limit\|CPU quota\|ROUNDFIX_VERIFICATION_CORES' CONTEXT.md` returned no matches.
- Post-change: `grep -n 'Verification Core Budget' CONTEXT.md` returns line 548; the definition line and the `_Avoid_` line match the prescribed text.
- `go test -count=1 -tags docscontract ./internal/docscontract -run 'TestGlossary'` → `ok roundfix/internal/docscontract 0.461s` (covers the four glossary contracts, including the new entry).

Acceptance evidence:
- The new entry is present at the prescribed position and carries the prescribed text, link, and `_Avoid_` line.
- The four glossary contracts pass in the focused run; the declared `## Verification` command remains for the Daemon to run.

## Carry-forward provenance

- Source Run: `run_20261111T064558Z_6d985452394229fa`
- Source commit: `763fb7f208209b282299d8a0faca96078093c864`
