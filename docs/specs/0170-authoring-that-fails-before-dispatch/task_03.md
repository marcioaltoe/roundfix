---
task: task_03
spec: 0170-authoring-that-fails-before-dispatch
status: completed
type: backend
complexity: medium
---

# Task 03: Close three checker honesty gaps and pin the closure cases

## Overview

Three existing Spec Consistency Check detectors accept what they should refuse
or refuse what they should accept: `SC-VERIFY-INVERTED-EXIT` approves a tool
piped into `grep` whose status the pipeline hides, `parsePromiseSection` refuses
a wrapped `None.` reason, and the reference index accepts an adopted source whose
original is still in place. The archived-Finding closure check is correct but
unpinned. This Task fixes the three and adds the missing tests.

## Requirements

1. MUST make `invertedExitVerificationCommand` in
   `internal/speccheck/verification.go` report `SC-VERIFY-INVERTED-EXIT` for a
   command substitution holding a pipeline whose last member is `grep` and whose
   first member is not `printf` or `echo`, when that substitution, or a variable
   assigned from it, is tested with `test -z` or `[ -z`; the finding MUST name the
   form and give the status-preserving replacement
   `out="$(tool 2>&1)" || exit 1; ! printf '%s\n' "$out" | grep -q pattern`.
2. MUST NOT report that form when `set -o pipefail` precedes it in the same
   command, and MUST NOT report the capture form the write-tasks skill
   recommends (`matches="$(find path -name '*.tmp' -print)" || exit 1; test -z "$matches"`)
   or the repository's named-pass pattern, which pipes `printf` into `grep -q`.
3. MUST make `parsePromiseSection` in `internal/speccheck/citations.go` accept,
   for both Success Metrics and API Contracts, one paragraph that begins with
   `None.` and carries a reason on the same or a following line.
4. MUST keep refusing a section that is missing, empty, holds `None.` without a
   reason, or follows the `None.` paragraph with a blank-line-separated block, a
   list line or a table line.
5. MUST make `detectReferenceIndex` in `citations.go` read each row's `source`
   cell and report `SC-REF-UNRESOLVED` when that source is under
   `docs/findings/` or `docs/backlog/` and still exists as a file, naming both
   the original and the indexed copy.
6. MUST add tests to `internal/speccheck/citations_test.go` for an archived
   Finding with only `closure_reason`, with only `closure_evidence`, and with a
   whitespace-only `closure_reason`, each expecting `SC-ARCHIVE-LICENSE`.
7. MUST leave every count of the active corpus golden unchanged.

## Subtasks

- [ ] Recognise the pipe-to-`grep` emptiness form and its safe neighbours.
- [ ] Accept a wrapped `None.` paragraph and keep the refusals.
- [ ] Refuse an adopted source left at its original path.
- [ ] Add the three closure tests.

## Acceptance Criteria

- [ ] The captured and inline pipe-to-`grep` forms are refused; the
      status-preserving, skill-recommended, named-pass and `pipefail` forms are
      not.
- [ ] A wrapped `None.` passes in both sections; a reasonless or mixed section
      is refused.
- [ ] A left-behind Backlog Entry and Finding are refused; a single move passes.
- [ ] Each closure case is refused with `SC-ARCHIVE-LICENSE`.

## Context

- interface: `internal/speccheck/verification.go`
- interface: `internal/speccheck/verification_test.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/citations_test.go`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestInvertedExitRefusesACapturedPipeToGrepTestedEmpty|TestInvertedExitRefusesAnInlinePipeToGrepTestedEmpty|TestInvertedExitAcceptsAStatusPreservingCapture|TestInvertedExitAcceptsTheSkillCaptureForm|TestInvertedExitAcceptsAPipefailPipeToGrep|TestInvertedExitAcceptsTheNamedPassPattern|TestWrappedNoneSuccessMetricsIsAccepted|TestWrappedNoneAPIContractsIsAccepted|TestNoneWithoutAReasonIsRefused|TestNoneFollowedByASecondBlockIsRefused|TestAdoptedBacklogEntryLeftInPlaceIsRefused|TestAdoptedFindingLeftInPlaceIsRefused|TestAdoptedSourceMovedOncePasses|TestArchivedFindingWithOnlyClosureReasonIsRefused|TestArchivedFindingWithOnlyClosureEvidenceIsRefused|TestArchivedFindingWithABlankClosureReasonIsRefused|TestCheckCorpusGolden|TestCheckActiveCorpusHasNoErrors)$" ./internal/speccheck ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestInvertedExitRefusesACapturedPipeToGrepTestedEmpty TestInvertedExitRefusesAnInlinePipeToGrepTestedEmpty TestInvertedExitAcceptsAStatusPreservingCapture TestInvertedExitAcceptsTheSkillCaptureForm TestInvertedExitAcceptsAPipefailPipeToGrep TestInvertedExitAcceptsTheNamedPassPattern TestWrappedNoneSuccessMetricsIsAccepted TestWrappedNoneAPIContractsIsAccepted TestNoneWithoutAReasonIsRefused TestNoneFollowedByASecondBlockIsRefused TestAdoptedBacklogEntryLeftInPlaceIsRefused TestAdoptedFindingLeftInPlaceIsRefused TestAdoptedSourceMovedOncePasses TestArchivedFindingWithOnlyClosureReasonIsRefused TestArchivedFindingWithOnlyClosureEvidenceIsRefused TestArchivedFindingWithABlankClosureReasonIsRefused TestCheckCorpusGolden TestCheckActiveCorpusHasNoErrors; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the sixteen new tests exists, so the command fails.

## References

- `_prd.md` → Goal 4; Core Feature 3; Success Metric 5.
- `_techspec.md` → Honesty gaps; API Contract 3; Testing Approach 3; ADR-0083;
  ADR-0092; ADR-0093.

## Result

Implemented the three checker corrections and pinned the existing archived-
Finding closure behavior:

- `SC-VERIFY-INVERTED-EXIT` now recognises captured and inline pipe-to-`grep`
  emptiness tests, names that form, and gives the status-preserving replacement.
  It leaves `pipefail`, the write-tasks capture form, and the named-pass pattern
  accepted.
- Promise declarations now treat wrapped lines as one `None.` paragraph while
  refusing a missing reason, a second paragraph, a list, a table, or numbered
  declarations mixed into that paragraph.
- Reference-index rows retain their `source` cell and report an adopted Finding
  or Backlog Entry that remains at its active-family source path, naming both
  that original and the indexed copy.
- Regression tests now pin reason-only, evidence-only, and blank-reason
  archived Findings as `SC-ARCHIVE-LICENSE` errors.

Focused-check evidence by acceptance criterion:

1. `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 ./internal/speccheck`
   passed, including the captured, inline, status-preserving, skill-capture,
   `pipefail`, and named-pass cases.
2. The same package check passed the wrapped Success Metrics and API Contracts
   cases and the reasonless, second-block, list, table, and mixed numbered
   refusals.
3. The same package check passed left-behind Backlog Entry, left-behind Finding,
   and single-move reference-index cases.
4. The same package check passed all three incomplete closure-field cases with
   `SC-ARCHIVE-LICENSE`.

Additional focused evidence:

- `GOCACHE=/tmp/roundfix-task03-gocache go test -count=1 -tags docscontract ./internal/docscontract`
  passed; the active corpus golden remains unchanged.
- `make verify-incremental` passed after rerunning outside the restricted
  network sandbox required by its GitHub-backed tests.
- The Task's authored `## Verification` command was not run; Daemon
  Verification remains the settlement authority.
