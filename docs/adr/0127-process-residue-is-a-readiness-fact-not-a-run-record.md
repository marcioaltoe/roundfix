---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-14T21:06:00Z
updated_at: 2026-08-14T21:06:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Process residue is a readiness fact, not a Run record

A process Roundfix spawned that outlived its Run has, by definition, no live Run
record to hang from — the four survivors measured on 2026-08-06 burned two hours
and forty minutes of CPU while `runs list` correctly reported nothing. Residue is
therefore reported by the readiness diagnostic, which already answers "what is
true about this machine right now", rather than by a new command or by inventing
Run rows for processes the Run Database never owned. ADR-0014 keeps the Daemon the
owner of verification and settlement; the inventory reports and settles nothing.
