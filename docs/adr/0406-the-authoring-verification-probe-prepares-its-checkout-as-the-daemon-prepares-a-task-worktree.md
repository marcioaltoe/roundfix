---
status: accepted
created_at: 2026-11-14T15:00:00Z
updated_at: 2026-11-14T15:00:00Z
deprecated_at: null
superseded_by: null
---

# The authoring Verification probe prepares its checkout as the Daemon prepares a Task Worktree

`roundfix spec check --run-verification` prepares its disposable checkout of
`HEAD` the way the Daemon prepares a Task Worktree before its pre-work probe:
it copies the configured `worktree.copy` files from the repository and runs
the configured Worktree Bootstrap there, with the same command, timeout and
environment, through the same provisioning function the item worktree uses.
It does so once per invocation, right before the first command it probes. A
bootstrap failure ends the check before any command runs. A repository that
configures neither `worktree.copy` nor `worktree.bootstrap` is probed and
reported exactly as before.

## Context

On 2026-11-14 indexus reported that `roundfix spec check --strict
--run-verification` (0.143.1) reported a `bunx vitest run ...` command
"honest — exited non-zero before work", and that minutes later the Daemon
refused the same command as vacuous, because it exited zero against the
unchanged tree in the Run's worktree. Indexus configures a Worktree Bootstrap
that installs its dependencies.

The authoring probe runs each command in `speccheck.DisposableCheckout`, a
bare `git worktree add --detach HEAD`. A Run Worktree and a Task Worktree
receive the `worktree.copy` files and the Worktree Bootstrap before any Agent
work or Verification (ADR-0034), so the Daemon's pre-work probe runs on a
prepared tree. A command that needs installed dependencies or a copied file
therefore fails in the authoring probe for a reason unrelated to the work and
is called honest, while the Daemon, on the prepared tree, sees it pass and
refuses the Task. ADR-0148 made one prober serve both callers so that "a
checker that approves what the probe later refuses" cannot exist, but it let
the authoring caller supply "a working directory and nothing else"; the
working directory is where the two callers still differ.

Characterized on 2026-11-14 at `4f8a5b7d` with a temporary repository whose
Project Config sets `worktree.bootstrap: "touch .deps-installed"`: the check
reported `test -f .deps-installed` as `honest`; the Daemon, which runs the
bootstrap in the Task Worktree before the probe, sees it exit zero.

## Decision

- **The same preparation, through the same code.** Before the first command
  it probes, the check calls the worktree package's provisioning function on
  its disposable checkout with the repository's `worktree.copy` list and its
  Worktree Bootstrap (`worktree.bootstrap` and `worktree.bootstrap_timeout`).
  It is the function that provisions an item worktree, built from the two
  steps that provision a Run and a Task Worktree, in their order: copy, then
  bootstrap. Commands keep running through the shared prober and the
  Daemon's command runner.
- **Once per invocation, only when a command will run.** One checkout serves
  every Spec of the invocation, as before, and the preparation runs at most
  once, right before the first probed command. A check that probes nothing
  (no Task Graph, only completed Tasks, only a QA gate) prepares nothing.
- **A failed bootstrap probes nothing.** The check exits `1`, prints no
  report, and names the bootstrap failure on standard error with the
  bootstrap's last output, as a Run stops at the same failure before any
  Agent work. A command is never called honest on a tree whose preparation
  failed.
- **The report names the preparation.** A Spec report whose commands ran
  after the bootstrap prints `Verification tree bootstrap: "<command>"`
  after `Verification tree: HEAD`, and its JSON carries
  `verification.bootstrap`. Without a bootstrap neither appears.

## Considered Options

- **Keep the bare checkout and say when a failure does not come from the
  command.** Rejected: the probe sees an exit status, not a cause, so it
  cannot tell a missing dependency from a missing feature without guessing,
  and the verdict would still differ from the Daemon's.
- **Point the checkout at the user's installed dependencies** (copy or link
  `node_modules`, a virtual environment). Rejected: it is not what the
  Daemon does, it ties the verdict to the state of the user's checkout, and
  it covers one ecosystem at a time.
- **Prepare the checkout as soon as it exists.** Rejected: every
  `--run-verification` call that probes nothing, such as a TechSpec-stage
  check of a Spec without Tasks, would pay a dependency install for nothing.
- **Prepare one checkout per Spec.** Rejected: the commands of every Spec run
  against the same `HEAD`, so a second preparation repeats the first at the
  same cost.
- **Report each command `unknown` with the bootstrap failure as its cause.**
  Rejected: the failure is one fact about the tree, repeated once per
  command it would be noise, and the Daemon stops before any probe on the
  same failure.

## Consequences

- A command that needs what the Worktree Bootstrap or `worktree.copy`
  provides gets the verdict the Daemon gives it: vacuous when it already
  passes on the prepared tree, honest only when it fails there.
- A `--run-verification` check that probes at least one command in a
  repository with a Worktree Bootstrap runs that bootstrap once, bounded by
  `worktree.bootstrap_timeout`: an install with a warm package cache costs
  seconds, a cold one can cost minutes. The bootstrap is not gated by the
  Machine Verification Capacity and receives no Verification Core Budget,
  as ADR-0384 and ADR-0385 already state for every Worktree Bootstrap and for
  this check.
- A bootstrap that needs a network the caller's environment denies, such as
  an Agent's sandbox, fails the check with the bootstrap named, where it used
  to produce verdicts the Daemon would contradict.
- A missing `worktree.copy` source prints the same `Run Worktree copy skipped
  missing <path>` line on standard error that the Daemon prints.

This decision amends ADR-0148, whose authoring caller now prepares its
working directory as the Daemon's Task Worktree is prepared, and extends
ADR-0034, whose Worktree Bootstrap now also runs in the authoring probe's
checkout.
