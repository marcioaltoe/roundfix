---
task: task_05
spec: 0132-a-grant-read-exactly-where-it-lives
status: pending
type: backend
complexity: low
---

# Task 05: Recognize the record naming already in use

## Overview

Discovery filters candidate records on filenames containing `authorization`, so
the existing valid record `2026-09-08-authorized-qa-archive-override.md` is
ignored and a repository whose only sanctioned regeneration declaration uses
that naming reports legitimate generated output as an unexpected write. This
slice widens the filter to the naming already in use without widening it to
everything.

This is an authorized tooling Task. It may change only
`internal/suiteguardcontract/regeneration.go`, its test, and this Task file.
Stop before any other mutation. The bounded set comes from
[_authorization.md](_authorization.md).

## Requirements

1. MUST recognize a record whose filename carries the `authoriz` stem, so both
   `authorization` and `authorized` names are candidates.
2. MUST keep the filter meaningful: a file that cannot carry a record is still
   not read, so the walk's cost stays proportional to candidates.
3. MUST keep every current discovery outcome, including the refusal of
   proposed, null-dated, malformed and unrelated records.
4. MUST NOT read every markdown file under the Spec roots, which is the cost the
   filter exists to avoid.
5. MUST update only the characterization row this Task intentionally moves.

## Subtasks

- [ ] Widen the filename filter to the shared stem.
- [ ] Prove the existing `authorized`-named record is discovered.
- [ ] Prove non-candidate files are still not read.

## Acceptance Criteria

- [ ] `docs/specs/0122-verified-content-and-terminal-settlement/references/2026-09-08-authorized-qa-archive-override.md`
      is discovered as a candidate record; it was ignored before this Task.
- [ ] A markdown file whose name carries no authorization stem is still not
      read, proven by a read-count assertion rather than by inspection.
- [ ] Proposed, null-dated, malformed and unrelated records still contribute
      nothing.

## Context

- interface: `internal/suiteguardcontract/regeneration.go`

## Verification

- `grep -q 'func TestDiscoveryRecognizesAuthorizedNaming' internal/suiteguardcontract/regeneration_test.go && go test -count=1 ./internal/suiteguardcontract -run '^TestDiscoveryRecognizesAuthorizedNaming$'` — the `authorized`-named record is discovered; this fails today.
- `grep -q 'func TestDiscoveryRecognizesAuthorizedNaming' internal/suiteguardcontract/regeneration_test.go || exit 1; go test -count=1 ./internal/suiteguardcontract -run '^TestSanctionedRegenerationReadsOnlyCandidateRecords$'` — the read count still tracks candidates rather than every markdown file.
- `grep -q 'func TestDiscoveryRecognizesAuthorizedNaming' internal/suiteguardcontract/regeneration_test.go || exit 1; go test -count=1 ./internal/suiteguardcontract` — every preserved discovery outcome still holds.

## References

- `_prd.md` → Core Features 5; Success Metrics.
- `_techspec.md` → Testing Approach observation 3; Build Order 5.
