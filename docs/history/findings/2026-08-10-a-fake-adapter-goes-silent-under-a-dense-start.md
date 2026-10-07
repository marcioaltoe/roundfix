---
status: done
absorbed_by: 0103-a-suite-that-leaks-nothing
created_at: 2026-08-10
updated_at: 2026-09-08
---

# A fake adapter goes silent under a dense start

The CI suite carries a flake family that predates any single branch: `main` shows 2 failures in its last 20 verification runs, one of them on a docs-only commit, and across 2026-08-10 seven different tests failed exactly once each — `TestCheckAdapterProvesOfficialClaudePackageAndVersion/version_only_with_command_package_identity`, `TestProveExactSelectionTimeoutCleanup`, `TestOwnerProcessControllerTerminateTreeProvesOutlivingGrandchildGone`, `TestProveExactSelectionOfficialFixturesNoPrompt/Sol_high`, `TestACPXRunSkipsEmptyReasoningEffort`, and `TestProveExactSelectionCleanupJoinedFailure`, and `TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow` (locally, under a fresh -count=1 suite). Every one exercises a spawned fake — a `#!/bin/sh` adapter or a child process — and none has failed twice with the same name.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-10-a-fake-adapter-goes-silent-under-a-dense-start.md`.
