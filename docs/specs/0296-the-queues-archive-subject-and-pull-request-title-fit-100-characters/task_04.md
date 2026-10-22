---
task: task_04
spec: 0296-the-queues-archive-subject-and-pull-request-title-fit-100-characters
status: pending
type: test
complexity: low
---

# Task 04: The title test tells UTF-16 units from bytes

## Overview

Corrective Task 1 of 2. It answers finding F-01 of the QA Report of
2026-10-21 (Run `run_20261022T014206Z_8044b8cc31aec528`, row QA-04,
task_03 Requirement 3). The QA gate made `queueHeaderLength` in
`internal/cli/deliver_workflow.go` return `len(header)`, counting bytes, and
`TestQueuePullRequestTitleKeepsTheSlugWhileItFits` still passed: every slug
and prefix in its table is ASCII, where bytes and UTF-16 units agree, and the
title cut counts its runes with `utf16.RuneLen` without calling
`queueHeaderLength`. A byte count would move a fitting non-ASCII slug to the
Spec Number form, which is the defect class the Spec's UTF-16 rule exists to
prevent. The production code follows the TechSpec; only the title test is
too weak to prove it.

## Requirements

1. MUST add to `queuePullRequestTitleCases` in
   `internal/cli/deliver_queue_subject_test.go` a case named
   `77 units non ASCII` whose slug is `"0296-" + strings.Repeat("é", 72)`
   (77 UTF-16 units, 149 bytes) and whose expected title is
   `"feat: deliver " + slug`, so the slug form fits the 91-unit budget only
   when the length is counted in UTF-16 units.
2. MUST make `TestQueuePullRequestTitleKeepsTheSlugWhileItFits` treat the new
   case like `77 units`: the title reader is never called and the test fails
   if it is. Derive that from the case (for example, the expected title equals
   `"feat: deliver " + slug`) instead of comparing the case name with one
   literal, and keep every existing case's expectation and the `reader error`
   subtest unchanged.
3. MUST keep `TestQueueMessagesStayReadableBySquashReaders`, which ranges over
   the same table, passing with the new case: its squash subject still yields
   the number and the type `feat`.
4. MUST prove the negative control before handing back: in a temporary
   `git archive` copy of the tree outside the checkout, replace
   `return len(utf16.Encode([]rune(header)))` with `return len(header)` in
   `queueHeaderLength`, run
   `go test -count=1 -v -run '^TestQueuePullRequestTitleKeepsTheSlugWhileItFits$' ./internal/cli`
   there, and observe it fail on `77_units_non_ASCII`. Write nothing under the
   checkout for this proof, and leave no copy behind.
5. MUST NOT change production code, any other test, any skill, guide, record
   or planning file. The finding is a test sensitivity gap, not a production
   defect.

## Subtasks

- [ ] Add the `77 units non ASCII` case and derive the no-read expectation from the case.
- [ ] Prove in a temporary copy that byte counting fails `77_units_non_ASCII`.

## Acceptance Criteria

- [ ] Every title case, the new one included, passes on the intended code.
- [ ] A `queueHeaderLength` that counts bytes fails `TestQueuePullRequestTitleKeepsTheSlugWhileItFits/77_units_non_ASCII`.

## Consumers

None. The Task changes one test file and no contract another package reads.

## Context

- instruction: `docs/adr/0309-a-queue-archive-subject-and-pull-request-title-keep-the-slug-while-it-fits-100-characters.md`
- interface: `internal/cli/deliver_queue_subject_test.go`
- instruction: `internal/cli/deliver_workflow.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestQueuePullRequestTitleKeepsTheSlugWhileItFits|TestQueueMessagesStayReadableBySquashReaders)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestQueuePullRequestTitleKeepsTheSlugWhileItFits/77_units_non_ASCII TestQueuePullRequestTitleKeepsTheSlugWhileItFits/77_units TestQueueMessagesStayReadableBySquashReaders/title/77_units_non_ASCII; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run '^TestQueue' ./internal/cli || exit 1; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the table has no `77 units non ASCII` case, so the command fails at the first missing pass; afterwards the new subtests, the existing ones and the Windows vet pass.

## References

- `_techspec.md` → API Contract 2
- `task_03.md` → Requirement 3
