---
status: accepted
created_at: 2026-07-06T21:05:00Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Operational Runs can detach from their caller

`--detach` on resolve, watch, and implement re-executes roundfix as a session leader independent of the invoking process, hands the run id back through a startup handshake, prints it with the attach/stop hints, and exits 0 — because a Run's lifetime must not depend on the caller's: three Runs died mid-flight when an invoking session reaped its background tasks, and caller-side `nohup` gymnastics is not a contract anyone should depend on. Detached Runs write their console stream to a per-Run log under the Artifact Directory (their only record), and `roundfix attach`/`roundfix stop` are the follow and control surfaces. This is the stepping stone to the long-lived `roundfix serve` daemon (work-plan item 3), not a replacement for it.
