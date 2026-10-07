---
task: task_06
spec: 0242-an-archive-that-leaves-an-archive-record
status: completed
type: backend
complexity: low
---

# Task 06: The archive confirmation names promoted files, and a keyless plan says no advice

## Overview

The first QA gate of this Spec (`qa/qa-report-2026-10-06.md` in Run
`run_20261006T232220Z_7c2e4a372c1dd72e`) failed rows 03 and 08 on two output
lines that do not match the authored transcripts. Finding F-01: an archive
with `--promote` copies the file and records it, but the confirmation line
ends at the revision. Surface Transcript 2 requires the suffix
`; promoted 1 file(s) to docs/references/`. `runArchiveCommand` never reads
`result.Promoted`. Finding F-02: `roundfix archive <slug> --plan` without a
judge key prints the skip reason alone, as
`candidate references/<file>.md <bytes> bytes: <KEY_VARIABLE> is not set`.
Surface Transcript 4 requires `no advice (<KEY_VARIABLE> is not set)`.
`runArchivePlan` replaces the detail with the reason instead of wrapping it.

This Task corrects both lines in `internal/cli/archive.go` and proves them
with tests that compare the exact transcript lines. It changes no archive
behavior, no record field and no guidance. task_07 documents both lines in
the user guide.

## Requirements

1. MUST append `; promoted <n> file(s) to docs/references/` to the archive
   confirmation line when `result.Promoted` is not empty, where `<n>` is
   `len(result.Promoted)`. The suffix follows the
   `kept in Git at <12-hex>` part, for a normal archive and for an archive
   with `--qa-override`, `--approval` and `--reason` alike (Surface
   Transcript 2). Without `--promote`, the line stays byte-identical to
   Surface Transcript 1.
2. MUST print every candidate line of `--plan` whose advice has no choice
   as `candidate <path> <bytes> bytes: no advice (<reason>)`, wrapping
   whatever skip or failure reason the advice carries (Surface Transcript 4).
   An advice with neither a choice nor a reason keeps printing `no advice`.
   A candidate with a choice keeps its `<choice> (p=<probability>)` form.
3. MUST add `internal/cli/archive_output_test.go` with these tests:
   - `TestArchiveConfirmationNamesPromotedFiles`: in a committed temporary
     repository, archive an eligible Spec with
     `--promote references/knowledge.md`, once normally and once with
     `--qa-override --approval maintainer --reason <300 bytes>`. Each stdout
     equals the Surface Transcript 1 or override line built from the
     record, followed by `; promoted 1 file(s) to docs/references/` and a
     newline, and stderr is empty.
   - `TestArchivePlanWithoutAKeyPrintsNoAdvice`: with every judge key
     variable from `judge.Load().KeyVariables()` set to the empty string,
     `--plan` on a Spec with one candidate file exits 0 and stdout contains
     the exact line
     `candidate references/knowledge.md <bytes> bytes: no advice (<first key variable> is not set)`.
4. MUST keep `TestArchiveCommandLeavesTheArchiveRecord` and every other
   existing archive test unchanged and passing.
5. MUST NOT send a request to a provider, read a credential, or read or
   write the real `~/.roundfix` in any test.

## Subtasks

- [ ] Append the promotion suffix to the confirmation line.
- [ ] Wrap the plan's skip reason in `no advice (...)`.
- [ ] Add the two transcript tests.

## Acceptance Criteria

- [ ] A promoting archive prints Surface Transcript 2's line.
- [ ] A keyless plan prints Surface Transcript 4's candidate line.
- [ ] An archive without a promotion still prints Surface Transcript 1's
      line.

## Context

- creates: `internal/cli/archive_output_test.go`
- interface: `internal/cli/archive.go`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestArchiveConfirmationNamesPromotedFiles|TestArchivePlanWithoutAKeyPrintsNoAdvice|TestArchiveCommandLeavesTheArchiveRecord)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveConfirmationNamesPromotedFiles TestArchivePlanWithoutAKeyPrintsNoAdvice TestArchiveCommandLeavesTheArchiveRecord; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two new tests do not exist, so the command fails; with the tests and without the fix, both fail on the exact line.

## References

- `_prd.md` → Core Feature 3; Success Metric 3
- `_techspec.md` → API Contract 3; API Contract 4; Surface Transcripts 1, 2 and 4
- QA Report `qa-report-2026-10-06.md` → F-01; F-02; rows 03 and 08
- ADR-0247

## Result

Implemented the two output corrections in `internal/cli/archive.go`.
The confirmation appends the promotion count after the Git revision for
both normal and QA override archives. Plans retain `no advice` and wrap
the advice reason in parentheses. Choice/probability formatting and the
no-reason fallback remain unchanged.

Focused checks:

- Before the implementation change,
  `rtk proxy go test -count=1 -run 'TestArchive(ConfirmationNamesPromotedFiles|PlanWithoutAKeyPrintsNoAdvice)$' ./internal/cli`
  exited 1: both promotion cases lacked the suffix, and the keyless candidate
  printed the bare reason instead of `no advice (...)`.
- After the change and formatting,
  `rtk proxy go test -count=1 -run 'Test(Archive|RunArchive)' ./internal/cli`
  exited 0 (`ok roundfix/internal/cli`, 5.134s).

Acceptance evidence:

- Promoting archive: `TestArchiveConfirmationNamesPromotedFiles` compares
  complete stdout against the record-derived confirmation plus exactly
  `; promoted 1 file(s) to docs/references/` and a newline. Both normal and
  override cases assert empty stderr; the override uses approval `maintainer`
  and a 300-byte reason.
- Keyless plan: `TestArchivePlanWithoutAKeyPrintsNoAdvice` clears every
  variable returned by `judge.Load().KeyVariables()` and checks the exact
  candidate line with its byte count and the first key variable's skip reason.
  It asserts exit 0 and empty stderr.
- No promotion: the unchanged `TestArchiveCommandLeavesTheArchiveRecord`
  runs in the focused selection and still compares the complete normal and
  override confirmations with the original record-derived transcript.

The new tests use committed temporary repositories, temporary home directories,
explicit environments containing no credentials, and a fake HTTP transport
that rejects any request. Both tests assert zero provider requests. Existing
archive tests were left unchanged.

Declared Verification was not run; Task status and settlement remain
Daemon-owned. No commit, push, PR, Task Graph edit, or other Task edit was made.
The user-guide update remains task_07's slice.

## Carry-forward provenance

- Source Run: `run_20261007T012355Z_4430cb7122321779`
- Source commit: `ef9b4972a14500e984f489966fd3cc1956a471e6`
