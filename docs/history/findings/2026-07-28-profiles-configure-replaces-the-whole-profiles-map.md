---
status: done
created_at: 2026-07-28
updated_at: 2026-09-08
absorbed_by: 0056-profiles-configure-merge-semantics
---

# profiles configure — a one-category fragment deletes every other configured profile (2026-07-28)

`roundfix profiles configure --scope project --file <fragment>` replaces the entire `profiles:` map with the fragment instead of merging the named categories into it. Configuring one category therefore **silently deletes every other configured profile**.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-28-profiles-configure-replaces-the-whole-profiles-map.md`.
