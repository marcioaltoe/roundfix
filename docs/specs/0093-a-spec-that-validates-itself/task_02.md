---
task: task_02
spec: 0093-a-spec-that-validates-itself
status: pending
type: backend
complexity: high
---

# Task 02: Read a cited decision against the claim made about it

## Overview

The Spec's centre. When an artifact attributes a subject to a decision record —
"ADR-XXXX makes Y authoritative" — the detector resolves that record and matches
the claim against its text. An unsupported claim is reported with both texts, so
a maintainer settles it by reading two quotes rather than trusting the checker's
judgement.

## Requirements

1. MUST detect claims that attribute a subject to a decision record in a Spec's
   PRD and TechSpec.
2. MUST resolve each claim's target record and report
   `SC-CITATION-UNSUPPORTED` when the record's text does not carry the claimed
   subject.
3. MUST include, in the finding, the claiming sentence and the cited record's
   own subject, so the maintainer can settle it without opening either file.
4. MUST report the number of claims resolved, so an artifact making no claims is
   distinguishable from one whose claims were not parsed; a detector that
   silently matches nothing would pass everything.
5. MUST NOT report a finding for a citation that merely lists a record without
   claiming what it establishes; that case belongs to the existing checks.
6. MUST break the characterization case Task 01 declared, and update it in the
   same commit.

## Subtasks

- [ ] Parse attributions from PRD and TechSpec text.
- [ ] Resolve and match against the cited record.
- [ ] Report with both texts and the resolved-claim count.

## Acceptance Criteria

- [ ] Spec 0090's original claim about ADR-0083 is reported unsupported.
- [ ] The corrected wording is not reported.
- [ ] A bare listing of an ADR is not reported.
- [ ] The resolved-claim count is observable and distinguishes zero claims from
      zero parses.

## Rehearsal Cases

- Case: the exact sentence "ADR-0083 makes `make verify` the only authoritative
  gate" against the real ADR-0083; Observation: `SC-CITATION-UNSUPPORTED`, with
  the sentence and ADR-0083's own title in the finding.
- Case: "ADR-0096 has the QA gate prove machine facts before it spends an Agent
  turn" against the real ADR-0096; Observation: no finding.
- Case: an Active ADR obligations row that lists a decision record without
  claiming what it establishes; Observation: no finding, and the resolved-claim
  count excludes it.

## Bounded scope

This Task may create or modify only:

- `internal/speccheck/citations.go`
- `internal/speccheck/citations_test.go`
- `internal/speccheck/citation_characterization_test.go`
- `internal/speccheck/testdata/citation/**`
- `docs/specs/0093-a-spec-that-validates-itself/task_02.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestCitation' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCitationReportsAnUnsupportedClaim'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestCitation' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCitationAcceptsASupportedClaim'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestCitation' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCitationReportsHowManyClaimsItResolved'` — expected: exits 0.
- `grep -q 'SC-CITATION-UNSUPPORTED' internal/speccheck/citations.go` — expected: exits 0. This string does not exist before this Task.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestCitationCharacterization' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestCitationCharacterization'` — expected: exits 0, proving the declared break was updated to the new behaviour rather than left failing. A whole-package sweep would pass with the work absent; this names the case that must change.

## References

- `_prd.md` → Goals 1 and 4.
- `_techspec.md` → Build Order 2; Interfaces.
- ADR-0116.
