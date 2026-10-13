---
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
status: active
created: 2026-10-12
surfaces: [backend, cli, docs]
---

# A Run Worktree under a symlinked root is cleaned

Since 2026-10-10 the maintainer's `~/.roundfix/worktrees` has been a symbolic
link to another volume, and no Run Worktree under it has been cleaned since.
Every Delivery Queue merge ends with
`warning: cleanup failed: keep Run "<run-id>": recorded Run Worktree path is unsafe`,
`roundfix reconcile` printed that reason 550 times on 2026-10-12, and the
worktree root held 24 GB that nothing reclaims
([the symlinked root Backlog Entry](references/2026-10-12-a-run-worktree-under-a-symlinked-root-is-never-cleaned.md)).
Roundfix records each worktree in the configured spelling, Git registers the
resolved path, and reconciliation refuses any recorded path that crosses a
symbolic link. The cycle's other worktree source, the timing-dependent
admin-lock test, is fixed by Spec 0259's corrective task_06 and closed in
history, so this Spec does not own it.

This is a bug fix: no command, flag, output format or exit code changes. A
symbolic-linked worktree location starts to work as a real directory does, and
`roundfix reconcile --apply` reclaims what earlier versions kept.

The maintainer decided on 2026-10-12, finally: a recorded Run Worktree path and
a registered worktree path that resolve to the same directory are the same
worktree; a path resolving outside the worktree root stays unsafe; every
comparison of recorded and registered worktree paths is swept for the same
class; and the command that reclaims the retained Run
Worktrees must exist, work on the symlinked root and be documented.

## Prerequisites

None. Spec 0259 merges first and fixes the disposal lock tests, which this
Spec no longer touches. Specs 0263 and 0264 of cycle AF change archive, cancel and the Baseline
group usage error, and touch no file this Spec declares. Spec 0258, not yet
merged, adds structural gaps to `spec check --strict` that hold every active
Spec; this Spec already carries the Absence Evidence, ADR Alternatives,
Coverage Cross-Read and Consumers sections it reads. No Verification here
depends on another Spec's artifacts.

## Project Constraints

- Identifier strategy: not applicable — no identifier, schema or record format
  changes; recorded Run paths are read, never rewritten, and the Run Database
  is not migrated. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; tests and Verification use temporary homes, temporary
  repositories and the local `git` binary. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0268 (this Spec) decides the
  Worktree Root rule and the prefix-resolving comparison. ADR-0053: reconciliation "is read-only unless `--apply` is supplied"
  and "Ambiguous, dirty, and unintegrated work stays intact", which stays true
  under a symbolic-linked root. ADR-0023: each Run executes in a worktree
  "living under Roundfix Home, recorded on the Run row", the record this Spec
  reads without rewriting. ADR-0259: "Every top-level test in a Parallel Test Package calls
  `t.Parallel()` as its first statement", and each changed package runs once
  with `-race`. ADR-0256: "A release waits for the skills that describe a
  changed surface", answered by a Coverage Review for the reconcile surfaces.
  ADR-0244: "A Spec now carries a Glossary Declaration". ADR-0257: "Inside a
  Run the Daemon is the only full gate". ADR-0258: "A Spec is authored against
  the QA rerun classes". ADR-0161, ADR-0212 and ADR-0232 release a merged Spec's
  Runs and item branches on merge evidence; their proofs stay as they are and
  now also reach Runs under a symbolic-linked root. ADR-0170: carry-forward of
  a completed Task is untouched, because carry-forward reads no recorded
  worktree path through the changed rule. ADR-0141 does not apply: review Runs
  execute in the user checkout and record no Run Worktree. ADR-0237: the merge
  record a Delivery Retry writes is read as before. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or
  authorized; the cycle AF grant of 2026-10-12 ("Sim, como no AE") covers the
  Governed Paths each Spec declares, and this Spec's Tasks declare none, so
  the record carries `paths: []`. No Makefile, CI workflow, lint, formatter,
  module or skill file changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0262-a-run-worktree-under-a-symlinked-root-is-cleaned/_authorization.md`;
  bounded files: none (`paths: []`).

## Goals

- A terminal Run whose worktree lives under a symbolic-linked
  `worktree.location` is reconciled, released and pruned exactly as it would
  be under a real directory.
- A recorded path whose repository or worktree directory is itself a symbolic
  link stays unsafe and is preserved.
- Every comparison of a recorded worktree path with a registered one agrees
  when both resolve to the same directory, also when the directory is gone.
- `roundfix reconcile --apply` reclaims the Run Worktrees an earlier version
  kept as unsafe, and the user guide says so.

## User Stories

1. As the maintainer with a symbolic-linked worktree location, I want a
   Delivery Queue merge to remove the merged Spec's Run Worktrees, so that the
   disk stops filling.
2. As the maintainer, I want `roundfix reconcile --apply` to release the Run
   Worktrees that earlier versions kept as unsafe, so that I can reclaim the
   24 GB with Roundfix's own command.
3. As the maintainer, I want a recorded path that a symbolic link redirects
   away from its worktree location to stay preserved, so that cleanup never
   reaches another tree.
4. As the Delivery Queue, I want an item or Task Worktree whose directory is
   already gone under a symbolic-linked root to be pruned or recreated, so that
   a retried item does not fail on a stale registration.

## Core Features

1. **Worktree Root rule.** A recorded worktree path is accepted when its
   repository directory and worktree directory are real directories; symbolic
   links are followed only in the Worktree Root above them. It is then the same
   worktree as a registered path that resolves to the same directory.
2. **Outside the root stays unsafe.** A recorded path whose repository or
   worktree directory is a symbolic link keeps the reason
   `recorded Run Worktree path is unsafe` and is preserved.
3. **Prefix-resolving comparison.** Comparing a recorded and a registered
   worktree path resolves the longest existing prefix of each and keeps the
   missing remainder, so a worktree whose directory is gone still matches the
   path Git keeps for it.
4. **Every consumer covered.** Reconcile dry run and apply, the Delivery Queue
   release of merged Spec Runs, the terminal prune at `implement` and Force
   Stop, the superseded-branch disposition, and item worktree use and cleanup
   all work under a symbolic-linked root.
5. **Reclaim documented.** The reconcile user guide states that a
   symbolic-linked `worktree.location` is resolved within its Worktree Root and
   that `roundfix reconcile --apply` releases the Runs an earlier version kept.
6. **Glossary.** `CONTEXT.md` defines **Worktree Root**.

## User Experience

A repository whose worktree location is a real directory sees no change. One
whose location is a symbolic link sees `roundfix reconcile` classify its
terminal Runs as it would without the link, `safe` for a merged clean Run
instead of `unknown` with `recorded Run Worktree path is unsafe`, and
`roundfix reconcile --apply` release them. A Delivery Queue merge no longer
prints `warning: cleanup failed` for those Runs.

## Non-Goals / Out of Scope

- A recorded Git root that crosses a symbolic link; it stays refused, and its
  tests stay as they are.
- Threading the configured `worktree.location` into reconciliation, or
  rewriting the recorded paths in the Run Database.
- Changing the reconcile report shape, its reasons or its exit codes.
- Carry-forward staging worktrees, which already resolve their root before
  `git worktree add`.
- The disposal lock tests, which Spec 0259's corrective task_06 fixes.

## Success Metrics

1. Success Metric: in a temporary home whose worktree location is a symbolic
   link, a merged clean terminal Run is `safe` in `roundfix reconcile` and
   released by `roundfix reconcile --apply`, its Run Worktree and Run Branch
   gone.
2. Success Metric: a recorded path whose repository directory, worktree
   directory, or dangling worktree link is a symbolic link stays `unknown` with
   `recorded Run Worktree path is unsafe`, and nothing is removed.
3. Success Metric: under a symbolic-linked root, an item worktree whose
   directory is gone is cleaned up and recreated, and a Delivery Queue release
   and the superseded-branch disposition of a merged Run succeed.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- Git's `builtin/worktree.c`, `add_worktree`
  (https://github.com/git/git/blob/5d7573a151b918062087822a2f0d7661dc4bd707/builtin/worktree.c,
  found and read through Exa):
  `write_file(sb.buf, "%s", real_path(sb_git.buf));`. Git records a new
  worktree under its resolved path, which is why the recorded spelling and the
  registered path differ under a symbolic link. The installed Git 2.56.0,
  measured offline during authoring, lists a worktree added through a
  symbolic-linked directory under the resolved directory; the QA gate repeats
  that measurement.
- The Secondbrain, read during authoring through `wiki/index.md` and two
  `qmd query` searches: the mirror of ADR-0023
  (`projects/roundfix/mirror/docs/adr/0023-runs-execute-in-per-run-worktrees.md`)
  and ADR-0053 confirmed that Run Worktrees live under Roundfix Home and that
  reconciliation is proof-based, so the fix widens what counts as the same
  worktree without weakening any proof. No page covered symbolic-linked
  worktree roots.

## Glossary

- adds: **Worktree Root**
- not a term: **Worktree Root rule** — the name of this Spec's acceptance rule for recorded paths
- not a term: **Outside the root stays unsafe** — a Core Feature heading
- not a term: **Prefix-resolving comparison** — a Core Feature heading
- not a term: **Every consumer covered** — a Core Feature heading
- not a term: **Reclaim documented** — a Core Feature heading

## Skills

- unchanged: command reconcile — task_03 adds one paragraph to the reconcile user guide about a symbolic-linked `worktree.location`; the Roundfix skill's reconcile reference states no path rule and already tells agents to inspect the dry run and then apply, so its text stays.
- unchanged: guide docs/user-guide/commands/reconcile.md — the same paragraph; the skill's reconcile reference needs no new text.
- unchanged: guide docs/user-guide/configuration.md — task_03 adds one sentence to the `worktree.location` row's section saying a symbolic link is resolved within the Worktree Root; no key, default or validation changes.
- unchanged: config keys — the same sentence; `worktree.location` keeps its key, default and validation.

## Decisions

- A recorded worktree path is resolved within its Worktree Root; see ADR-0268.
- Path comparison resolves the longest existing prefix; see ADR-0268.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the sweep, the contracts and the
build order.
