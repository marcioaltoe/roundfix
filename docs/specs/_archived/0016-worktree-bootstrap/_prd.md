---
spec: 0016-worktree-bootstrap
status: archived
created: 2026-07-06
surfaces: [cli, infra, docs]
archived: "2026-07-06"
source_slug: 0016-worktree-bootstrap
---


# Worktree Bootstrap

Roundfix executes each agent Run in an isolated worktree created from committed
Git state. That worktree has no installed dependencies, no migrated or seeded
database, and no warmed caches — so for a stateful project (a TypeScript
monorepo with a package manager, a database, and a build cache) the Daemon's
Verification fails not because the work is wrong but because the environment was
never prepared. `worktree.copy` already places untracked files like `.env` into
the worktree, but it only copies files. This Spec adds a configured
`worktree.bootstrap` command that Roundfix runs once in each new worktree — after
copying files, before any Agent work and before Verification — so the worktree
is a working environment. Together with `worktree.copy` and sequential
execution, this unblocks Roundfix for stateful monorepos without giving up
worktree isolation.

## Goals

- A new Run or Task Worktree can be prepared into a working environment (install
  dependencies, migrate and seed a database, warm caches) before Agent work and
  Verification run. See ADR-0034.
- Environment preparation is configured by the repository, not hardcoded:
  Roundfix runs the command; the command owns dependency and database strategy.
- A bootstrap failure is unambiguous and non-silent: the Run or Task ends with a
  distinct bootstrap-failed outcome rather than a confusing Verification failure.
- Untracked environment files (`.env` and siblings) reliably reach the worktree
  through the existing `worktree.copy`, documented as part of this recipe.

## User Stories

1. As a developer with a TypeScript monorepo, I want Roundfix to run
   `bun install && bun run db:migrate && bun run db:seed` in each new worktree
   before Verification, so that `make verify` passes on a prepared environment
   instead of failing on missing dependencies or an unmigrated database.
2. As a developer, I want the bootstrap command to be repository configuration,
   so that each project encodes its own setup without changing Roundfix.
3. As a developer, I want a bootstrap failure to end the Run or Task with a clear
   "worktree bootstrap failed" message naming the command, so that I can tell an
   environment problem from a real Verification failure.
4. As a developer running a shared-database project, I want sequential execution
   to bootstrap once on the reused Run Worktree, so that concurrent bootstraps do
   not fight over one database. See ADR-0034.
5. As a developer, I want `.env` and similar untracked files copied into the
   worktree through `worktree.copy`, documented with the bootstrap recipe, so
   that the whole monorepo setup is one coherent configuration.

## Core Features

1. **Bootstrap command.** A `worktree.bootstrap` config value: a command Roundfix
   runs once in each newly created worktree, in the worktree root, after
   `worktree.copy` places files and before the Agent starts and before
   Verification. Empty means no bootstrap (today's behavior). See ADR-0034.
2. **Bootstrap failure handling.** A non-zero bootstrap exit ends the owning Run
   (for the Run Worktree) or settles the owning Task failed (for a Task
   Worktree) with a distinct message naming the bootstrap command and its
   failure; the environment is treated as not ready.
3. **Once-per-worktree scope.** Bootstrap runs once per created worktree that
   hosts Agent work: the Run Worktree for sequential Runs and review Runs, and
   each Task Worktree when Tasks run concurrently. The command is expected to be
   idempotent.
4. **Bounded execution.** Bootstrap runs under a configurable timeout so a hung
   setup cannot stall a Run indefinitely; its output goes to stderr and the Run
   Event Journal.
5. **Env-file recipe.** Documentation ties `worktree.copy` (for `.env` and
   untracked env files, gitignore-safe) together with `worktree.bootstrap` and
   `worktree.concurrency: 1` into a documented recipe for stateful monorepos.

## User Experience

A developer adds a `worktree` block — `copy` for env files, `bootstrap` for the
setup command, `concurrency: 1` for a shared database — and Roundfix Runs then
prepare each worktree before working. Bootstrap output streams to stderr like
other Run diagnostics; on success nothing new appears in the deterministic
stdout report. On failure, the Run or Task ends with a `worktree bootstrap
failed: <command>` message. Projects with no `worktree.bootstrap` set behave
exactly as today.

## Non-Goals / Out of Scope

- Provisioning or owning databases, containers, or services — the bootstrap
  command does that; Roundfix only runs it. Per-worktree database isolation is a
  project concern encoded in the command.
- Installing or symlinking `node_modules` on Roundfix's behalf — the bootstrap
  command installs dependencies (using the package manager's own cache).
- Caching or reusing a bootstrapped environment across Runs — each new worktree
  bootstraps fresh; caching is the package manager's job.
- A general pre/post hook system — this Spec adds one bootstrap step, not an
  arbitrary lifecycle hook framework.
- Changing `worktree.copy` behavior beyond documenting it (it already copies
  untracked files into worktrees).

## Success Metrics

- With `worktree.bootstrap` set, a new worktree has dependencies installed and
  the database migrated/seeded before Verification runs, and a monorepo spec's
  `make verify` passes in the worktree (asserted with a fake bootstrap in tests,
  validated by a real monorepo Run).
- A failing bootstrap command ends the Run or Task with a `worktree bootstrap
  failed: <command>` message and does not run Agent work or Verification for that
  worktree.
- With no `worktree.bootstrap` configured, Run behavior is byte-stable versus
  today.
- A shared-database monorepo at `worktree.concurrency: 1` bootstraps exactly once
  per Run.

## Decisions

- Roundfix runs a configured `worktree.bootstrap` command once per new worktree,
  after `worktree.copy`, before Agent work and Verification; it owns running,
  not what it does. See ADR-0034.
- Bootstrap failure is a distinct terminal signal (Run failed / Task failed with
  a bootstrap-failed message), never a silent or Verification-shaped failure.
- `.env` and untracked env files are handled by the existing `worktree.copy`
  (gitignore-safe); this Spec documents the combined recipe. The glossary gains
  Worktree Bootstrap.

## Open Questions

None.
