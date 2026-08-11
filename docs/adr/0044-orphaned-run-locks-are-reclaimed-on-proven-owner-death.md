---
status: accepted
created_at: 2026-07-15T16:53:06Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Orphaned Run locks are reclaimed only on proven owner death

Active-Run locks carry the owning process identity, and any surface that hits a lock (new Run preflight, Stop Command, Branch Integrity Preflight) checks the owner's liveness: a provably dead owner makes the lock orphaned, and Roundfix reclaims it automatically with a stderr warning and a Run Event Journal record; anything short of proof keeps the block exactly as ADR-0012/0016 define, naming the run id and the stop command. Observed motivation: a killed run process left its Stop Request unprocessed and its lock Active, blocking every relaunch until a manual force stop — but silent reclamation without proof of death would let two live Runs race on one target.
