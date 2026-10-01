---
status: accepted
created_at: 2026-10-01T00:00:00Z
updated_at: 2026-10-01T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An evidence snapshot records one digest per declared input

ADR-0194 has the Daemon record, for each carriable passing row, the path and
SHA-256 digest of every file a declared input matched. A glob input such as
`internal/**` or `**` matches thousands of files, so the record grew with the
repository rather than with the matrix. Spec 0203's QA Report reached 157,212
lines, almost all of them per-file digests, and the pre-PR review's agent
session failed on the resulting diff.

The record now holds one entry per declared input: its ref, the number of
files it matched at the audited head, and one SHA-256 digest of a summary
holding, for each matched file in ascending byte order of its
repository-relative path, the hexadecimal SHA-256 of its content, two spaces,
the path and a newline. This is the summary Go's `h1:` module hash uses. A
file added to, removed from or changed under an input changes that digest, so
the carry compares one digest per input and names the input that moved. A
report recorded in the per-file form is still read: the stage derives the same
digest from its file list.

## Consequences

This refines ADR-0194 without superseding it, and ADR-0097's carry conditions
are unchanged. The block's size is bounded by rows times inputs. A carried row
proves the same thing as before, but the record no longer says which file
under a moved input changed. A disposition names the declared input, not the
changed files. A matched path containing a newline cannot be summarized, so
its input is unresolved and the row is executed again.
