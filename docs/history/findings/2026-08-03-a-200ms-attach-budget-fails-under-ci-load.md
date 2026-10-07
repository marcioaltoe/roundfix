---
date: 2026-08-03
surface: internal/cli
status: done
updated_at: 2026-09-08
absorbed_by: 0124-verification-capacity-and-measured-economics
---

# A 200 ms attach budget fails the Verification gate under CI load

`TestRunImplementDetachSurvivesCallerProcessGroupKill` failed the CI Verification gate on PR #98, a documentation-only change touching fourteen Markdown files and no Go source. Re-running the identical job on the identical commit passed. The same test passed locally in `make verify` (3,136 tests) both before and after the failure, and passed in isolation in 0.80 s.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-03-a-200ms-attach-budget-fails-under-ci-load.md`.
