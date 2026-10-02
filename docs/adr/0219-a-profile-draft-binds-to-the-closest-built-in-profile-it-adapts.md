---
status: accepted
created_at: 2026-10-02T00:00:00Z
updated_at: 2026-10-02T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A profile draft binds to the closest built-in profile it adapts

A repository-owned Profile draft passed with `--profile-file` names no source.
Roundfix finds it by checking the draft against every built-in profile and
used to accept it only when exactly one profile admitted it (ADR-0075). A
composed built-in profile selects the modules of two others, so a draft cut
from the Standard TypeScript Monorepo Profile, or from the Go CLI/TUI profile
without its TUI surface, is now also a valid adaptation of
`go-cli-typescript-monorepo` (ADR-0204). Under the old rule both drafts were
refused as ambiguous, although nothing about them had changed.

A draft now binds to the closest built-in profile it adapts. That is the one
whose modules the draft removes fewest of. Fewest removed capabilities break a
tie. If two or more profiles are still equally close, the draft stays
ambiguous, and the refusal names only those profiles. A draft that adapts one
profile binds to it as before, and a draft that adapts none is still
unresolved.

Two alternatives were rejected. Leaving composed profiles out of the search
unless the draft names a module only they add would hard-code the current
catalog: the next composed profile would need its own exception. Leaving out
every profile whose modules contain another candidate's would not resolve the
Go CLI/TUI draft, because the composed profile lacks the TUI surface and so
contains neither candidate.

## Consequences

The source still chooses the remediation profile and the retention lineage of
the plan, so a TypeScript-only draft keeps the TypeScript profile's setup and
history instead of the composed one's. A draft that an earlier catalog refused
as ambiguous between two profiles at different distances now binds to the
closer one. A tie is never broken by catalog order or by name; it is refused,
so the source is always a property of the draft and not of the listing.
