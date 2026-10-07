---
status: done
created_at: 2026-07-25
updated_at: 2026-08-06
absorbed_by: 0058-npm-trusted-publishing-and-release-preflight
---

# npm publishing — token authentication and registry state can produce a partial release (2026-07-25)

The `v0.0.1` release reset exposed two independent publication risks. Roundfix still authenticates npm publication with a long-lived repository secret, and the release workflow does not prove that every package can accept the target version before it begins publishing the platform packages.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-25-npm-trusted-publishing-and-release-preflight.md`.
