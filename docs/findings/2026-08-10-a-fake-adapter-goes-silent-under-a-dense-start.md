---
status: pending
created_at: 2026-08-10
updated_at: 2026-08-14
---

# A fake adapter goes silent under a dense start

The CI suite carries a flake family that predates any single branch: `main`
shows 2 failures in its last 20 verification runs, one of them on a docs-only
commit, and across 2026-08-10 seven different tests failed exactly once each —
`TestCheckAdapterProvesOfficialClaudePackageAndVersion/version_only_with_command_package_identity`,
`TestProveExactSelectionTimeoutCleanup`,
`TestOwnerProcessControllerTerminateTreeProvesOutlivingGrandchildGone`,
`TestProveExactSelectionOfficialFixturesNoPrompt/Sol_high`,
`TestACPXRunSkipsEmptyReasoningEffort`, and
`TestProveExactSelectionCleanupJoinedFailure`, and
`TestRunImplementVerificationCapacityAndDaemonStatusIntegratedFlow` (locally,
under a fresh -count=1 suite). Every one exercises a spawned
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

## 2026-08-14 — the silent-probe signature again, on a docs-only branch

`TestACPXRunAppliesSelectionBeforePrompt/codex_reasoning_effort` failed on
PR #161, a branch that changes nothing but Markdown. The message is the
silent-probe signature exactly: the fake `npx` written by `installFakeAdapter`
produced no parseable `--version` output, so `inspectAdapter` returned an
`AdapterLineageError` naming the package it could not prove. The same test
passes locally under `-count=1`, and the run's other twenty-four packages were
green.

This is the fourth month-to-date occurrence and the first on a branch that
cannot have caused it, which settles the question of whether the family is
branch-specific: it is not. It also raises the cost of the flake beyond wasted
CI minutes — a docs-only Pull Request cannot be merged without a rerun, so the
loop pays for it at every delivery, not only at every code change.

## 2026-08-14 — the kernel names the race: `text file busy`

PR #162's verification gate failed at
`TestProjectDecisionJourney/affected_Profile_apply_audit_and_empty_reapply`
with a message this record has been circling for four days:

    run fixture tool oxfmt --version: fork/exec /tmp/.../004/bin/oxfmt: text file busy

`ETXTBSY`. The kernel refuses to execute a file that some process still holds
open for writing. That converts the write-then-exec window from a hypothesis
into a named errno, and it settles the "silent probe" signature above: an
adapter that returns empty output and a fixture that refuses to start are the
same race caught at two different instants.

It also narrows the remedy. Retrying the probe treats the symptom. Closing the
file before exec, or replacing the shell fixtures with an `os.Args[0]` re-exec
the way the ACPX harness already does elsewhere, removes the window itself —
there is no descriptor left open to conflict with. The `os.WriteFile` call in
`installFakeAdapter` and its siblings do close their handle, so the open writer
is more likely a `cp`, an archive extraction, or a sibling test writing the same
fixture path concurrently; whichever it is, the fix is to stop executing a file
that was written moments earlier in a process tree that is still forking.

Two failures on two consecutive Pull Requests, in two different tests, on
branches whose changes cannot produce either signature. Both were green on
rerun.

