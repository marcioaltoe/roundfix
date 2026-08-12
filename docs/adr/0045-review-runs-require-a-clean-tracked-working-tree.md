---
status: accepted
created_at: 2026-07-15T16:53:06Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Review Runs require a clean tracked working tree

Because review Runs execute in the user's checkout (ADR-0042), resolve and watch Preflight Validation fails when tracked files are dirty, naming each dirty path and the stash/commit action; untracked files stay allowed because batch commits stage only paths changed since the batch's snapshot. The alternative — tolerating a dirty tree and excluding user paths from batch commits — was rejected: the Daemon cannot reliably distinguish user edits from Agent edits in the same file, and a swept-in user change would corrupt the batch commit contract (ADR-0001). Consequence: after a failed batch, everything dirty in the checkout is Agent work by construction, so the report can say so and leave recovery to the user without guessing.
