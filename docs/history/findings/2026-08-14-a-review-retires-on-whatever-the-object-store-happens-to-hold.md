---
status: done
absorbed_by: 0120-knowledge-lifecycle-and-durable-capture
created_at: 2026-08-14
updated_at: 2026-09-08
kind: finding
---

# A Review retires on whatever the object store happens to hold (2026-08-14)

Running the migration Spec 0094 shipped, against the repository that shipped it, produced three different answers for the same fifty orphan Review Artifacts within one session. Nothing about the pull requests changed between them. What changed was the local Git object store.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-14-a-review-retires-on-whatever-the-object-store-happens-to-hold.md`.
