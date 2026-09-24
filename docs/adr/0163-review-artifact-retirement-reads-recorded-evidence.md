---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-09-24T00:00:00Z
updated_at: 2026-09-24T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# Review Artifact retirement reads recorded evidence

An orphan Review Artifact retires on its recorded outcome. A recorded squash merge with its merge commit is a valid receipt; when no outcome is recorded, the result is an explicit unknown and the artifact is retained. No local Git object, ref, or ancestry decides liveness.

Relocating `docs/specs/_reviews/` is independent of liveness. The Review Artifact root resolver never resolves into the history root, so new artifacts remain in the live root while migration can relocate legacy artifacts separately.

Supersedes ADR-0123 and ADR-0152 (2026-09-24); both are archived under `docs/history/adr/`.
