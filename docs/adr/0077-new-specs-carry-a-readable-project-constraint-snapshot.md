---
status: accepted
created_at: 2026-07-24T20:34:11Z
updated_at: 2026-07-24T20:34:11Z
deprecated_at: null
superseded_by: null
---

# ADR-0077: New Specs carry a readable Project Constraint snapshot

Every new PRD and TechSpec contains a `Project Constraints` section with a
concise snapshot of each applicable value and its source in `docs/agents/`.
The section is the human-readable contract; no separate authorization
frontmatter is added.

Repository-owned authorial skills must refuse to conclude a new artifact when
the section is absent or incomplete. A tooling mutation remains unauthorized
until the section records express maintainer approval and the bounded files it
covers.

Snapshots can become older than the live Setup Manifest. That is intentional:
they record the contract under which a Spec was approved, while a later
material decision change requires the Spec author to revise and re-approve the
affected artifact.
