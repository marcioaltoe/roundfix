---
spec: 0018-external-spec-root
status: archived
created: 2026-07-07
surfaces: [cli, docs]
archived: "2026-07-07"
source_slug: 0018-external-spec-root
---


# External Spec Root

Roundfix assumes Spec artifacts live at `docs/specs/` inside the repository
working tree. Repositories that keep planning artifacts in a knowledge
workspace — a separate git repository nested in or beside the code repository,
reached through a symlink — break that assumption twice: the per-Task commit
stages the task file through a path that crosses a symbolic link, which git
refuses (`pathspec beyond a symbolic link`), and Run Worktrees materialize the
symlink without its target, so a relative link dangles. External Spec Root
makes the Spec artifact location a first-class configuration, so Roundfix
reads and writes Specs wherever they live and never stages external artifacts
into code-repository commits.

## Goals

- A repository whose Spec artifacts live outside its working tree (nested
  knowledge repository, symlinked directory) runs the full Spec lifecycle —
  implement, settle, QA, archive — without a failed commit or a dangling path.
- Spec resolution behaves identically from the user checkout and from every
  Run or Task Worktree.
- Repositories with the default layout see zero behavior change.

## User Stories

1. As a user whose `docs/specs` is a symlink into a knowledge workspace
   repository, I want Spec Runs to complete and commit, so that the daemon's
   per-Task commit no longer dies staging a path that crosses the symlink.
2. As a user configuring a repository, I want to point Roundfix at the Spec
   artifact directory explicitly, so that Runs resolve Specs from the
   knowledge workspace directly instead of depending on how a symlink resolves
   inside a Run Worktree.
3. As a user reading a per-Task commit in the code repository, I want it to
   carry only that repository's changes when Spec artifacts are external, so
   that the two repositories' histories stay separate and truthful.

## Core Features

1. A configuration value sets the Spec Root — the directory holding
   `<slug>/` Spec folders — with Project Config over User Config over the
   built-in `docs/specs` default. Relative values resolve against the
   repository root of the user's checkout; absolute values are used as-is.
2. Every Spec consumer resolves through the Spec Root: the Implement Command
   (Task Graph load, task status writes, QA Reports), the Settle Command, the
   Archive Command, Interactive Input's active-Spec listing, Attach's Task
   detail, and the review artifact location that follows the Spec.
3. The Spec Root is resolved once per command against the user's checkout and
   carried into the Run, so Run and Task Worktrees read and write the same
   directory the checkout does — relative roots never re-resolve against a
   worktree.
4. When the Spec Root resolves outside the repository working tree, Daemon
   commits carry only repository changes: task files and QA Reports are not
   staged, and the Run Event Journal records that the artifact was settled
   outside the repository. When every candidate path is external and the
   repository has no changes, the Task settles without a commit instead of
   failing.
5. Independent of configuration, commit staging never fails on a path that
   crosses a symbolic link: such paths are dropped from staging with a
   journaled warning naming the path.
6. Configuration validation rejects a Spec Root that does not exist or is not
   a directory, with an actionable message naming the resolved path.

## User Experience

- Default-layout repositories: no visible change.
- External-root repositories: Run startup reports the resolved Spec Root on
  stderr; per-Task commit reports name the artifacts kept outside the
  repository. Everything else — reports, outcomes, exit codes — is unchanged.

## Non-Goals / Out of Scope

- Committing or pushing anything in the external knowledge repository — its
  history belongs to the user.
- Multiple Spec Roots per repository, or per-Spec root overrides.
- Discovering the Spec Root automatically by following symlinks — resolution
  is by configuration; the symlink tolerance in staging is a safety net, not a
  discovery mechanism.
- Changing the Review Issue artifact contract beyond following the Spec Root.

## Success Metrics

- The conexus-style layout (Spec artifacts in a nested knowledge repository
  behind a `docs/specs` symlink) completes an Implement Run with per-Task
  commits, where today the first commit fails.
- A default-layout repository's test suite passes unchanged.

## Decisions

- Spec Root is configuration, not symlink-following — explicit beats inferred
  for a path that decides where writes land. See ADR-0035.
- External Spec artifacts never ride in code-repository commits; the task
  file's settled status lives in the external repository, and committing it
  there stays a user action. See ADR-0035.
- The symlink staging guard ships regardless of configuration, so
  un-configured symlink layouts fail no worse than "artifact not committed".

## Open Questions

None.
