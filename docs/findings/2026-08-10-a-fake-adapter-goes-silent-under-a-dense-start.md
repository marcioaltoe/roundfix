# A fake adapter goes silent under a dense start

The CI suite carries a flake family that predates any single branch: `main`
shows 2 failures in its last 20 verification runs, one of them on a docs-only
commit, and across 2026-08-10 six different tests failed exactly once each —
`TestCheckAdapterProvesOfficialClaudePackageAndVersion/version_only_with_command_package_identity`,
`TestProveExactSelectionTimeoutCleanup`,
`TestOwnerProcessControllerTerminateTreeProvesOutlivingGrandchildGone`,
`TestProveExactSelectionOfficialFixturesNoPrompt/Sol_high`,
`TestACPXRunSkipsEmptyReasoningEffort`, and
`TestProveExactSelectionCleanupJoinedFailure`. Every one exercises a spawned
fake — a `#!/bin/sh` adapter or a child process — and none has failed twice
with the same name.

Two signatures, one family:

- **Starved waits.** Before the cache work, failures were 90-second
  `agentWaitBudget` timeouts waiting for a file a shell script creates in
  milliseconds — the fork-storm of nested cold compiles kept the script from
  running at all.
- **Silent probes.** After the cache work made CI start dense (no compile
  phase staggering test starts), two tests failed at t=6s of a warm run with
  an `AdapterLineageError` whose message carries no package: the fake's
  `--version` probe returned empty output. `installFakeAdapter` writes the
  script with `os.WriteFile` and the probe executes it moments later while
  sibling parallel tests fork their own children.

The common condition is spawn density, not any specific test. Both reruns of
the same commit went green, and the local suite has never reproduced either
signature on 12 cores.

Worth settling when this is picked up: whether the probe should retry once on
empty output from a just-written script, whether `waitForFile` should watch the
child process so a dead fake fails in milliseconds instead of 90 seconds, and
whether the fakes should be compiled test binaries (`os.Args[0]` re-exec, which
the harness already uses elsewhere) instead of shell scripts, removing the
write-then-exec window entirely.
