# ADR-0039: Model availability preflight uses a disposable Agent Session

Status: Accepted

Before an operational command creates a Run, Roundfix opens a uniquely named
disposable Agent Session in the Git root, assigns the effective Agent Model and
runtime-specific reasoning option through acpx, and closes the session on every
path. This validates the installed adapter, authentication, model metadata, and
reasoning vocabulary together; the extra adapter startup is accepted because a
static catalog or runtime-owned config inspection cannot prove that the exact
selection will start. The real Agent Session repeats the same assignment order,
and any rejection fails Preflight Validation without fallback.
