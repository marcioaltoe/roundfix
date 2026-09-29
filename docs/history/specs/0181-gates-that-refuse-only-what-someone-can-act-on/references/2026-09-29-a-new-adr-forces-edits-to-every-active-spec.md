---
type: fix
status: promoted
created: 2026-09-29
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
reason: null
---

# A new ADR forces edits to every other active Spec's PRD

## Symptom

`SC-ADR-RELATED` reports every accepted ADR that cites an ADR an active Spec lists, when that Spec does not list the citing ADR. A Spec that mints a new ADR citing an older one therefore opens a gap in every other active Spec listing the older one. Wherever gaps are promoted, that gap refuses: `spec check --strict`, `deliver plan`, queue revalidation, and the QA gate. The corpus golden then fails `make verify`. PRs #268 (ADR-0161), #272 (ADR-0164) and #274 (ADR-0165) each carried hand edits to other Specs' Project Constraints rows before they could merge. None of those Specs could act on an ADR that did not exist when it was written.

## Where

`detectADRConsistency` and `firstListedCitation` in `internal/speccheck/citations.go`; `PromoteGaps` in `internal/speccheck/report.go`; the pinned count in `internal/docscontract/testdata/corpus-golden.json`.

## Expected

An ADR that entered the repository after a Spec's PRD was committed does not open `SC-ADR-RELATED` on that Spec. A Spec being written, or one whose PRD history cannot be read, keeps the full check. `SC-ADR-UNLISTED` and `SC-CITATION-UNSUPPORTED` are unchanged.

## Evidence

- Squash commits `06afa835` (#268), `bc339356` (#272) and `f9c5fcce` (#274).
- ADR-0116 states the check's purpose.
- `docs/history/findings/2026-08-06-minting-an-adr-opens-gaps-no-one-can-ever-close.md` states that "a check must fail only where someone can act on it".
