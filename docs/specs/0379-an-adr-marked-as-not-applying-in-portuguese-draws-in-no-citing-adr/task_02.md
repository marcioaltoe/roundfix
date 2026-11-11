---
task: task_02
spec: 0379-an-adr-marked-as-not-applying-in-portuguese-draws-in-no-citing-adr
status: pending
type: backend
complexity: medium
---

# Task 02: The strict check reads the Portuguese not-applicable token

## Overview

Fiscus's Specs mark ADRs that do not apply with "não se aplica"
([`references/2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md`](references/2026-11-10-an-adr-marked-as-not-applying-in-portuguese-still-seeds-the-related-adrs.md),
recorded by the operator), and the related-ADR seed rule reads only
`not applicable`. This Task makes `notApplicableRelationPattern` read the
closed list ADR-0389 decides and names both tokens in the `SC-ADR-RELATED`
fix text, with tests for every accepted form and for the forms that stay
outside the list.

## Requirements

1. MUST set `notApplicableRelationPattern` in
   `internal/speccheck/authoring_proof.go` to
   `regexp.MustCompile(`(?i)^(?:not applicable|n(?:\x{00E3}|a\x{0303}|a)o se aplica)\b`)`
   (`_techspec.md` → Interfaces). `appliedADRNumbers` and every other caller
   stay unchanged.
2. MUST change the `SC-ADR-RELATED` fix in `internal/speccheck/citations.go`
   so that its token clause reads
   `begin its ADR Alternatives line in the TechSpec with "not applicable" or "não se aplica", so the ADRs that cite only it are not drawn in.`
   The summary, severity, code and locations keep their bytes.
3. MUST update `wantFix` in `TestAnAppliedADRStillDrawsInItsCitingADR`
   (`internal/speccheck/adr_related_seed_test.go`) to the new fix text; this
   is the declared break of API Contract 2. No other assertion there changes.
4. MUST add `internal/speccheck/adr_related_marker_test.go` (package
   `speccheck_test`, reusing `writeADRRelatedSeedFixture`,
   `checkADRRelatedSeeds` and `requireADRRelatedSeedFinding`) with two tests,
   each calling `t.Parallel()` first and in each subtest:
   - `TestAPortugueseNotApplicableADRDrawsInNoCitingADR`, subtests and
     relations of the second fixture ADR: `accented` `não se aplica; a decisão segue igual.`;
     `capitalized` `Não se aplica — a decisão segue igual.`; `upper case`
     `NÃO SE APLICA: a decisão segue igual.`; `unaccented`
     `nao se aplica; a decisão segue igual.`; `decomposed`, written as the Go
     literal `"não se aplica; a decisão segue igual."`; `line end`
     `não se aplica`. Each expects no `CodeADRRelated`, `CodeADRUnlisted` or
     `CodeADRAlternativesUnrelated` finding.
   - `TestAMarkerOutsideTheClosedListStillDrawsIn`, subtests and relations:
     `plural` `não se aplicam; a decisão segue igual.`; `not leading`
     `mantida; a alternativa rejeitada não se aplica.`; `abbreviation`
     `n/a; unchanged.`; `other wording` `não aplicável; a decisão segue igual.`.
     Each calls `requireADRRelatedSeedFinding(t, result, "0003", "0002")`.
5. MUST keep `TestANotApplicableADRDrawsInNoCitingADR`,
   `TestWithoutATechSpecEveryListedADRSeeds` and
   `TestAnADRCitingAnAppliedADRIsStillRelated` unchanged and passing, and
   leave `internal/speccheck/constraints.go` untouched.

## Subtasks

- [ ] Widen the not-applicable pattern to the closed list.
- [ ] Name both tokens in the `SC-ADR-RELATED` fix text and update its pinned test.
- [ ] Add the accepted-forms and outside-the-list tests.

## Acceptance Criteria

- [ ] Every accepted Portuguese form leaves the citing ADR undrawn, and every form outside the list still reports the gap.
- [ ] The English token keeps its behavior and the fix text names both tokens.
- [ ] Each mutation of the pattern (no Portuguese alternative, no case folding, no decomposed form, no unaccented form, no word boundary) fails a named subtest.

## Context

- interface: `internal/speccheck/authoring_proof.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/adr_related_seed_test.go`
- creates: `internal/speccheck/adr_related_marker_test.go`
- instruction: `docs/adr/0389-an-adr-alternatives-line-marks-an-adr-as-not-applying-in-english-or-portuguese.md`

## Consumers

- The `SC-ADR-RELATED` seed rule and fix text — consumers: `./internal/speccheck`

## Verification

- `out="$(GOMAXPROCS=2 go test -count=1 -timeout 10m ./internal/speccheck -run '^(TestAPortugueseNotApplicableADRDrawsInNoCitingADR|TestAMarkerOutsideTheClosedListStillDrawsIn|TestAnAppliedADRStillDrawsInItsCitingADR|TestANotApplicableADRDrawsInNoCitingADR|TestWithoutATechSpecEveryListedADRSeeds|TestAnADRCitingAnAppliedADRIsStillRelated)$' -v 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestAPortugueseNotApplicableADRDrawsInNoCitingADR' && printf '%s\n' "$out" | grep -q -- '--- PASS: TestAMarkerOutsideTheClosedListStillDrawsIn'` — expected: every accepted form and every guard passes, and the English and fix-text tests still pass; fails before the work because the new tests do not exist.
- `GOMAXPROCS=2 go test -count=1 -timeout 10m ./internal/speccheck && grep -qF -- 'n(?:\x{00E3}|a\x{0303}|a)o se aplica' internal/speccheck/authoring_proof.go` — expected: the whole package passes with the widened pattern; fails before the work because the pattern names only the English token.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-3; Success Metric 1; Success Metric 2
- `_techspec.md` → Interfaces; API Contract 1; API Contract 2; Testing Approach; Build Order 2
- ADR-0389
