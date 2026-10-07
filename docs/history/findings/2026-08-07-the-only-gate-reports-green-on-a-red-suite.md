---
status: done
created_at: 2026-08-07
updated_at: 2026-09-08
kind: finding
absorbed_by: 0083-a-gate-that-can-say-no
---

# The only gate reports green on a red suite (2026-08-07)

`make verify` exits `0` on a working tree whose Go test suite exits `1`. The repository's own rule says the local gate is the only gate; that gate is currently capable of certifying a red build as verified.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-07-the-only-gate-reports-green-on-a-red-suite.md`.
