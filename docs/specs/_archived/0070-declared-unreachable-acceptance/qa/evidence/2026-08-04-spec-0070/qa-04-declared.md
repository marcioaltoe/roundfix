# QA-04 — Declared-only archive

Status: pass

`roundfix archive qa-case` exited 0 and printed:

```text
archived qa-case -> docs/specs/_archived/qa-case
```

A fresh read found no active folder and found the archived PRD stamped with:

```yaml
status: archived
archived: "2026-08-04"
source_slug: qa-case
unproven:
    - a maintainer performs the live publication and records it
```

A second invocation exited 2 because the active Spec no longer existed,
confirming the move persisted and was not repeated.
