---
status: done
created_at: 2026-07-27
updated_at: 2026-09-08
absorbed_by: 0055-owner-identity-without-fork
---

# Force Stop — owner identity forks `/usr/bin/ps`, so the escape hatch fails exactly when the host is loaded (2026-07-27)

Spec 0037 gave Force Stop a real ownership proof: Runs record an opaque owner start-time identity, and Force Stop refuses to signal a PID whose live identity does not match. The proof is correct, but it is obtained by forking `/usr/bin/ps` on every read. When the host cannot fork, the proof fails, and Force Stop fails closed — refusing to stop the Run.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-27-owner-identity-forks-ps-and-fails-closed-under-load.md`.
