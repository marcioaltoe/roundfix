---
task: task_05
spec: 0132-a-grant-read-exactly-where-it-lives
status: completed
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

## Result

Implementation evidence:

- Spec-record discovery now matches the shared `authoriz` filename stem. The
  existing `authorization` names and the `authorized` name therefore remain a
  bounded candidate set rather than widening discovery to every Markdown file.
- `TestDiscoveryRecognizesAuthorizedNaming` uses the exact
  `docs/specs/0122-verified-content-and-terminal-settlement/references/2026-09-08-authorized-qa-archive-override.md`
  path and asserts the discovered `make baseline-digests` declaration by deep
  equality. Before the implementation edit, this test returned no declaration.
- `TestSanctionedRegenerationReadsOnlyCandidateRecords` remains unchanged. Its
  read count stays at two after adding `task_01.md`, proving a Markdown filename
  without the `authoriz` stem is not read.
- `TestSanctionedRegenerationRejectsNonOperativeRecords` remains unchanged and
  continues to cover proposed, null-dated, malformed and unrelated records.

Focused-check evidence:

- Pre-change `rtk go test ./internal/suiteguardcontract -run
  TestDiscoveryRecognizesAuthorizedNaming` first encountered the sandbox's Go
  build-cache denial. The unchanged retry with cache access exited 1 because
  discovery returned `nil` instead of the expected declaration, establishing
  the intended red signal.
- `rtk go test ./internal/suiteguardcontract -run
  'Test(DiscoveryRecognizesAuthorizedNaming|SanctionedRegenerationReadsOnlyCandidateRecords|SanctionedRegenerationRejectsNonOperativeRecords)$'`
  exited 0 and reported 7 passing tests.
- `rtk gofmt -d internal/suiteguardcontract/regeneration.go
  internal/suiteguardcontract/regeneration_test.go` and `rtk git diff --check`
  both exited 0 with no output.
- `rtk make verify-incremental` stopped at `fmt-check` because the unchanged
  `internal/cli/baseline_skills_restore_test.go` and
  `internal/cli/baseline_assets_sync_test.go` need formatting. Those paths are
  outside this Task's authorized set and were not changed.
- The changed-file postflight lists only this Task file and the authorized
  `regeneration.go` implementation and test paths.
- The Task's declared `## Verification` commands were not run; Daemon
  Verification owns them after this handoff.
