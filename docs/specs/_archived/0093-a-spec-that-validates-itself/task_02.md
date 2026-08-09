---
task: task_02
spec: 0093-a-spec-that-validates-itself
status: completed
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

## Result

Implemented semantic citation detection for PRD and TechSpec prose. The parser
recognizes subject attributions, ignores fenced examples and bare ADR listings,
resolves accepted decision records, and reports `SC-CITATION-UNSUPPORTED` with
the claiming sentence, the record title, and both source locations. Surface
matching uses meaningful-word coverage plus title or phrase anchors; short,
ambiguous claims fail safe rather than accusing an artifact from incidental
words.

`CitationClaims` and `ResolvedCitationClaimCount` expose the parsed and resolved
counts without changing the existing `Result` report model. A caller can
therefore distinguish one parsed-and-resolved attribution from an artifact with
zero attribution claims.

Focused checks:

- Red control: `rtk proxy env GOCACHE="$PWD/.gocache" go test
  -buildvcs=false ./internal/speccheck -run
  'TestCitationReportsAnUnsupportedClaim$' -count=1` failed before production
  implementation because `CodeCitationUnsupported` and the citation APIs did
  not exist.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test -buildvcs=false
  ./internal/speccheck -run
  'TestCitation(Characterization|CharacterizationReadsACitedRecordBody|AcceptsASupportedClaim|BareListingIsNotAClaim|ReportsHowManyClaimsItResolved|ReportsAnUnsupportedClaim)$'
  -count=1` — exited `0` after the final matcher edit.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test -buildvcs=false
  ./internal/speccheck -run
  'TestCheck(ADRUnlisted|ADRClosureDepthOne|ADRCorpusAbsentSkipsADRDetectors|CitationCoverageErrorLocations)$'
  -count=1` — exited `0`, preserving the existing ADR inventory, closure, skip,
  and location behavior.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test -buildvcs=false
  ./internal/speccheck -count=1` — exited `1`. The new detector intentionally
  exposes unsupported claims in active Specs 0080, 0090, 0091, 0092, and 0093;
  the existing corpus golden does not yet characterize the new code, and the
  vocabulary check reports that `CONTEXT.md` does not yet document it. Those
  files are outside this Task's bounded scope.

Acceptance evidence:

- Spec 0090's original claim is reported unsupported:
  `TestCitationCharacterization` parses and resolves the real ADR-0083
  attributions and requires the active-row finding to quote
  ``ADR-0083 makes `make verify` the only authoritative gate`` and ADR-0083's
  title, `Adopted sources move to their owning Spec`.
- Supported wording is not reported: `TestCitationAcceptsASupportedClaim`
  checks the exact ADR-0096 rehearsal sentence against the real decision
  subject and observes no `SC-CITATION-UNSUPPORTED` finding. The updated
  characterization also removes the finding after ADR-0083's fixture body is
  changed to carry the claimed subject.
- A bare ADR listing is not reported: `TestCitationBareListingIsNotAClaim`
  lists ADR-0083 and ADR-0096 and observes no semantic finding.
- The resolved count distinguishes work from silence:
  `TestCitationReportsHowManyClaimsItResolved` observes one parsed claim and one
  resolved record for the attribution, then zero parsed and zero resolved
  claims for the bare listing.

Follow-ups outside this Task's slice:

- The stage registry added by Task 03 landed concurrently and does not yet
  register `SC-CITATION-UNSUPPORTED`; the authoring-stage integration owner
  must add this detector at the PRD stage before the authoring skills rely on
  stage-scoped checks.
- The new code and newly exposed active-Spec findings need an explicitly
  bounded owner for `CONTEXT.md`, the corpus characterization/golden, and any
  active artifact corrections. This Task cannot edit those paths.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
