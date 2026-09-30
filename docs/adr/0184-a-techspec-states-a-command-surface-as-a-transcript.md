---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A TechSpec states a command surface as a transcript

A TechSpec described a command's new behavior in prose, and each reader rebuilt
the exact output from it: the Task author in a test, the QA gate in a row and
the reviewer in a finding. The three versions drifted, and the drift surfaced
as corrective Tasks. API Contracts are numbered and traced to Tasks (ADR-0156),
but a contract sentence does not fix one byte of output.

A TechSpec now declares numbered Surface Transcripts, or `None.` with the
reason. Each is a `transcript` block holding one command line, its standard
output, its standard error and its exit code, written as the command will
behave once the Spec ships. Task tests and QA Requirements name the transcript
by number and reuse its text. The Spec Consistency Check verifies the block's
shape, that a non-QA Task names each transcript, and that the authored QA Task
names it in a Requirement. It never runs the command.

The declaration gap binds a Spec from the same guide commit as ADR-0183's
receipt gap. Shape and reference findings apply to every declared transcript.
Interfaces written as code signatures and numbered invariants are taught by
the same guide and are not checked mechanically.

## Consequences

- One text is the expected output for the test, the gate and the review.
- A transcript is a promise, so it is declared and traced like the promises
  ADR-0156 already names.
- The check proves shape and traceability. The QA gate proves the behavior by
  running the command.
