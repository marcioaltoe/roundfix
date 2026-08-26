---
status: superseded # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-07-10T20:01:58Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null
superseded_by: ADR-0140
---

# ADR-0039: Model availability preflight uses a disposable Agent Session

Before an operational command creates a Run, Roundfix opens a uniquely named
disposable Agent Session in the Git root, assigns the effective Agent Model and
runtime-specific reasoning option through acpx, and closes the session on every
path. This validates the installed adapter, authentication, model metadata, and
reasoning vocabulary together; the extra adapter startup is accepted because a
static catalog or runtime-owned config inspection cannot prove that the exact
selection will start. The real Agent Session repeats the same assignment order,
and any rejection fails Preflight Validation without fallback.
