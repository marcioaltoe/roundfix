---
status: accepted
created_at: 2026-08-07T00:00:00Z
updated_at: 2026-08-07T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A managed refresh proves preservation instead of backing it up

ADR-0070 gives every mutated root carrier an immutable content-addressed backup
because adoption rewrites carriers whose pre-existing instructions it did not
author and cannot reconstruct. A managed refresh has no such exposure: its
defining invariant is that every byte outside a managed marker is identical
before and after, so the bytes a backup would protect are the bytes the mode
guarantees it never touches. Emitting one `AGENTS.<digest>.md` per refresh would
therefore accumulate a file per release in every adopted repository to insure
against a loss the mode makes impossible. A managed refresh consequently takes
no new root backup and instead carries a preimage-bound proof that non-managed
regions are unchanged, which the plan presents and apply verifies. Keeping the
backup anyway was rejected as noise the maintainer must later garbage-collect;
weakening the invariant to a promise rather than a verified postimage was
rejected because the backup is only redundant while the proof is mechanical.
This narrows ADR-0070 to the adoption path that motivated it and does not
supersede it: first adoption still backs up every root carrier it rewrites.
