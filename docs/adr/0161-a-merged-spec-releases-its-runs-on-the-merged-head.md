---
status: accepted
created_at: 2026-09-28T16:02:46Z
updated_at: 2026-09-28T16:02:46Z
deprecated_at: null
superseded_by: null
---

# A merged Spec releases its Runs on the merged head

A terminal Run of a merged Spec is released when every commit its Run Branch holds and the merged head lacks is represented there. The Delivery Queue merge record is the first source of that head; the default branch carrying the archived Spec is the fallback.

The delivery owner releases the Spec's Runs after it records the merge. `--discard-superseded` remains the separate recorded act ADR-0115 requires; cleanup does not absorb that command or bypass its Branch Disposition.

This decision extends ADR-0053's proof set without superseding ADR-0053 or ADR-0115. The accepted cost is preservation whenever the merge record, archived-Spec fallback, Task status, QA Report supersession, or changed-file representation cannot be proven.
