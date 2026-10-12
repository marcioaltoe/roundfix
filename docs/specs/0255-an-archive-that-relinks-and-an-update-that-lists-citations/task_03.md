---
task: task_03
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
status: pending
type: backend
complexity: medium
---

# Task 03: The Delivery Queue commits and proves an archive that relinked files

## Overview

After task_02 an archive can change Live Documents outside the Spec folder,
and today's exact-archive proof refuses any changed path besides the folder,
the record and the promoted copies. This Task makes the queue's archive stage
stage the relinked files the record names, and makes the proof accept each one
only when its bytes equal the archive's rewrite of the parent's bytes
(ADR-0260, as ADR-0230 accepted outward rewrites). The Task is verifiable on its
own through the queue's archive tests on temporary repositories.

## Requirements

1. MUST stage `record.Relinked` together with the record and the promoted
   copies in the Delivery Queue's archive stage, so the archive commit holds
   exactly the change the archive made.
2. MUST extend `archiveRetirementIsExact`, for both the committed and the
   working-tree form, as `_techspec.md` → API Contract 6 states: a relinked
   path is allowed only when it is outside the Spec folder and the History
   Root, exists at the parent, and its bytes equal `spec.RelinkContent` of its
   parent bytes under `spec.ArchiveCitationMoves`. The promotions map pairs
   each promoted destination with the one folder file at the parent that has
   the same basename and blob; zero or several candidates fail the proof.
3. MUST keep every other rule of the proof: the record's identity checks, the
   promoted blob equality, and the refusal of any changed path the proof
   cannot account for (`_techspec.md` → Invariant 6). A legacy exact move and
   a record retirement without relinked files are proven exactly as today.
4. MUST add these tests to `internal/cli/deliver_archive_relink_test.go`,
   reusing the fixtures of `deliver_archive_record_test.go`:
   - `TestDeliveryArchiveStageCommitsRelinkedFiles`: the archive stage over an
     owner Spec that an active Spec cites commits the record, the promoted
     copy and the relinked file, and reports an exact archive.
   - `TestDeliveryAcceptsARelinkedRetirement`: a commit whose relinked file
     equals the recomputed rewrite is proven exact, committed and in the
     working tree.
   - `TestDeliveryRefusesATamperedRelinkedFile`: the same commit with one
     extra byte in the relinked file, and a commit that changes a file the
     record does not name, are not exact.
   - `TestDeliveryRefusesAnAmbiguousPromotionPairing`: two folder files with
     the promoted basename and blob make the proof fail.
5. MUST NOT change any existing assertion; the existing delivery archive tests
   pass unchanged.

## Subtasks

- [ ] Stage the relinked files in the archive stage.
- [ ] Recompute and compare each relinked file in the proof.
- [ ] Pair promotions with their source by basename and blob.
- [ ] Write the queue tests.

## Acceptance Criteria

- [ ] An archive that relinked files is committed and proven exact by the queue.
- [ ] A tampered or unexplained change in the archive commit is refused.
- [ ] Existing exact-archive and legacy-move proofs behave as before.

## Context

- interface: `internal/cli/deliver_workflow.go`
- creates: `internal/cli/deliver_archive_relink_test.go`
- instruction: `internal/cli/deliver_archive_record_test.go`
- instruction: `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- instruction: `docs/adr/0230-an-archived-spec-keeps-its-relative-links.md`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `.agents/skills/roundfix/references/deliver.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliveryArchiveStageCommitsRelinkedFiles|TestDeliveryAcceptsARelinkedRetirement|TestDeliveryRefusesATamperedRelinkedFile|TestDeliveryRefusesAnAmbiguousPromotionPairing)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestDeliveryArchiveStageCommitsRelinkedFiles TestDeliveryAcceptsARelinkedRetirement TestDeliveryRefusesATamperedRelinkedFile TestDeliveryRefusesAnAmbiguousPromotionPairing; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n%s\n' "$t" "$out" >&2; exit 1; }; done; go test -count=1 -run '^TestDelivery(AcceptsAnExactRetirement|RefusesAnInexactRetirement|KeepsTheLegacyExactMove|ArchiveStageCommitsTheArchiveRecord)$' ./internal/cli` — expected: exit 0. Before this Task the four tests do not exist, so it fails; after it they pass and the existing delivery archive tests pass unchanged.

## References

- `_prd.md` → Core Feature 3; Goals; User Story 1; Success Metric 4
- `_techspec.md` → API Contract 6; Invariant 6; Build Order 3
- ADR-0260; ADR-0230; ADR-0247
