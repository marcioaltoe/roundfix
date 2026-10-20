---
type: fix
status: promoted
created: 2026-10-20
spec: 0287-reports-that-count-what-they-say
reason: null
---

# An Archive Record summary ends at the colon that introduces its list

## Symptom

The Archive Record of Spec 0282 ends mid-sentence: "Two environment gaps
failed it:". The record's summary is the PRD's title paragraph, and that
paragraph introduces a list of the two gaps that the summary leaves out. On
2026-10-20, 38 of the 256 Archive Records under `docs/history/specs/` ended
with a colon.

## Where

`internal/spec/archive_record.go`, `archiveParagraph`: it collects the lines
of one paragraph and stops at the first blank line or heading. Every summary
goes through it:

- `BuildArchiveRecord` reads the PRD's title paragraph, or the QA Report's
  `## Outcome` paragraph when one exists;
- `legacyArchiveParagraph` strips block quote markers and calls it for a
  Legacy Archive Folder;
- `internal/spec/cancel.go` reads a cancelled Spec's summary with it;
- `ParseArchiveRecord` reads a written record back with it.

## Expected

A summary paragraph that ends with a colon carries what it introduces, so the
record reads as a whole. Records already written stay byte-identical.

## Evidence

`docs/history/specs/0282-the-linux-lane-runs-like-ci.md` and the PRD at its
`source_revision` `b026a989f1fb3fc600aa71dd5c060b27eb022fdc`; the count of
records ending with a colon on 2026-10-20; no Finding.
