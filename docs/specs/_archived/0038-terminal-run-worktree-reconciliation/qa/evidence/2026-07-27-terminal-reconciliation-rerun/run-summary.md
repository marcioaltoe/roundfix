# Terminal Run Worktree reconciliation — rerun evidence

Build: `3744d03669f75d05e8f8e728dc3e03e6300b61ce`
(`roundfix 0.0.1`, built as `3744d03-dirty` because the QA report was open).

Application: `bin/roundfix`, built by the current run's `rtk make verify`.

Environment: macOS Darwin 25.5.0 arm64; Go 1.26.5; Git 2.55.0.

Fixture root: `/private/tmp/roundfix-qa-0038-rerun.4GhgTo/state`.

The build-tagged `fixture.go` helper seeded only isolated Git repositories and
one isolated Roundfix Home through project Store and Worktree APIs. Every
verdict-bearing operation used the built public CLI. Confirmation used fresh
Git reads, database hashes, a second CLI process, or a read-only Store reopen.

## Prior-failure regression and scope audit

The prior QA run failed because Task 04 commit `0d56f42` added a 20.1 MiB
Mach-O executable at repository-root `roundfix`.

Current-build evidence:

- remediation commit `af03b5f` deletes `roundfix`;
- both `git ls-tree 8ec92ad -- roundfix` and
  `git ls-tree HEAD -- roundfix` return no entry;
- base and HEAD `.gitignore` both resolve to blob
  `498dd1c424271e8e612674bfcba3fed500eee252`;
- the post-failure commits add only the dated QA and findings artifacts after
  restoring the protected ignore file byte-for-byte;
- current status before report closure contains only the dated QA report and
  its evidence directory.

Every Task commit was audited with
`git diff-tree --no-commit-id --name-status -r <commit>`. Tasks 01-06 matched
their declared implementation and documentation slices except for the Task 04
binary now removed from HEAD.

Task 07 commit `16da6754efe30e28228d0498318eb01701d61a7d`
contains exactly:

- `.agents/skills/roundfix/SKILL.md`;
- `skills/roundfix/SKILL.md`;
- the five derived digest pins expressly authorized in the PRD and TechSpec;
- `task_07.md`.

Task 06 commit `c443d76` contains only `CONTEXT.md`, the two user guides, and
`task_06.md`; it touches no protected tooling path. `cmp` confirms the Skill
pair is byte-identical, and
`TestAuthorialSkillSync/typescript-bun.json` passes.

## Static gate

- The first exact `rtk make verify` attempt could not access the sandboxed
  default Go cache and stopped before compilation.
- The same exact command rerun with approved normal cache access passed:
  2,563 tests in 23 packages, four Skill-sync tests, all fourteen shipped
  Skill checks, and the production build.
- `rtk git -c core.fsmonitor=false diff --check` passed after the gate.
- Verification introduced no tracked side effect.

## Classification and dry-run

Repository:
`/private/tmp/roundfix-qa-0038-rerun.4GhgTo/state/matrix-repo`.

Text and `roundfix-reconcile/v1` JSON dry-runs both exited 0 and returned,
newest-first:

| Run | Stored outcome | Classification | Evidence |
| --- | --- | --- | --- |
| `run_20260727T204846Z_00c039441bf0c350` | Clean | released | Worktree and Run Branch absent |
| `run_20260727T204846Z_f5301d45ef12691c` | TimedOut | unknown | target branch cannot resolve |
| `run_20260727T204846Z_8f0ecb4080eff186` | Unresolved | unintegrated | Run head `3f7686a…` not ancestor of target `347f101…` |
| `run_20260727T204846Z_de3aa0216c7170e7` | Stopped | dirty | tracked and untracked changes |
| `run_20260727T204846Z_26931c4f50c1ca51` | Failed | safe | clean worktree; both heads `347f101…` |

Summary:
`total=5 safe=1 unintegrated=1 dirty=1 unknown=1 released=1 applied=0 preserved=3 operational-failures=0`.

The Run Database SHA-256 remained
`67bbd9f535dcea188f85fc91dae06a06f8f7df28c499174586d64f222b594cf7`
across the text, JSON, and flags-before/after-ID dry-runs. Fresh worktree and
ref reads also showed no mutation.

## Runs List discovery

Before apply:

- Active view: `No Runs found.` plus an exact retained count of four.
- Terminal and all-state repository views: five rows plus the same retained
  count of four.
- Machine-wide `--all --state all --limit 0`: ten retained terminal Run
  Worktrees across all fixture repositories.

After the matrix's safe entry was released, the Active view reported exactly
three retained terminal Run Worktrees. After the integration repository was
fully released, its all-state view preserved the five terminal history rows
and emitted no retained-worktree note. Focused stream tests separately confirm
that rows stay on stdout and guidance stays on stderr.

## Mixed apply and preservation

`roundfix reconcile --apply --format json` in the matrix repository exited 0:

`total=5 safe=1 unintegrated=1 dirty=1 unknown=1 released=1 applied=1 preserved=3 operational-failures=0`.

Fresh Git reads found only the safe worktree and Run Branch absent. The dirty,
unintegrated, and unknown worktrees and branches remained. Independent reads
returned `# tracked dirt`, `untracked dirt`, and the unique unintegrated commit
`3f7686a`.

## Durable outcomes and idempotency

Repository:
`/private/tmp/roundfix-qa-0038-rerun.4GhgTo/state/integration-repo`.

Repository-wide apply returned five `safe` results and `applied=5`. A fresh
`runs list --state all --limit 0` showed:

- `run_20260727T204846Z_0df4aca31b4cd727`: Clean, previously
  IntegrationPending;
- `run_20260727T204846Z_e848687119f05bc8`: Unresolved;
- `run_20260727T204846Z_ffa8eacb96affb88`: Failed;
- `run_20260727T204846Z_9285f861aef47ffe`: Stopped;
- `run_20260727T204846Z_2f9abf528397aaa5`: TimedOut.

A fresh read-only Store process found exactly one event for the promoted Run.
The event records source `daemon`, kind `daemon.outcome`, previous outcome
IntegrationPending, current outcome Clean, classification safe, both branches,
both `347f101…` heads, the worktree, and action `cleanup`.

The second apply returned five `released` results and `applied=0`. The Run
Database SHA-256 remained
`f53305c7bbd3bd656a6b3727d4120565e78c453087892fa966ebcbe3275f6767`,
and fresh worktree/ref reads were unchanged.

The sibling Task Worktree and branch ending in `-task_01` remained after the
parent Run Worktree and Run Branch were released.

## Selector refusals and operational recovery

From the selector repository, missing, Active, review, and cross-repository
Run selectors each exited 2 with a specific Preflight Validation message.
The no-ID JSON form returned an empty repository-scoped result set.

Invalid format, extra argument, `--force`, and missing `--format` value each
exited 2. The database SHA-256 remained
`f53305c7bbd3bd656a6b3727d4120565e78c453087892fa966ebcbe3275f6767`,
and selector-repository refs/worktrees remained unchanged.

A locked clean worktree produced exit 1, `operational-failures=1`, the exact
rerun action, and both remaining recoverable surfaces. Fresh Git reads found
the worktree still registered and locked and the Run Branch still present.
A subsequent dry-run continued to classify the retained entry as safe.

## Non-Goal probes

- Mixed apply did not integrate the divergent commit or repair either dirty
  file.
- The sibling Task Worktree and Task Branch remained.
- The Integration Pending reconciliation event remained readable after
  cleanup; Reconcile did not prune journal or artifact state.
- The selector repository's no-ID form excluded every other repository.
- Help and user guides expose no age-, outcome-, missing-path-, or force-based
  deletion path and keep GC ownership separate.

## Focused current-build evidence

- `go test ./internal/worktree -run
  'Test(InspectTerminalRun|ApplyTerminalRun|PruneTerminal)' -count=1` passed.
- `go test ./internal/store -run 'TestReconcileIntegration' -count=1`
  passed.
- The focused CLI family for Reconcile, retained Runs List, Implement
  Preflight, help, and documentation contracts passed.
- `go test -race ./internal/cli ./internal/worktree -run
  'Test.*Reconcile' -count=1` passed.
- `go test -race ./internal/worktree -run
  'Test(ApplyTerminalRun|PruneTerminal)' -count=1` passed.
- `go test -race ./internal/worktree -run 'TestInspectTerminalRun' -count=1`
  passed.
- The Active, terminal, and all-state Runs List regression family passed.

These focused checks cover the stale-head, changed-metadata, newly-dirty,
persistence-failure-before-Git, unsafe/symlink path, command-ordering,
non-force removal, stream separation, and automatic reaper seams that cannot
be paused safely through the public process.

## Docs and shipped Skill

- Built `roundfix reconcile --help` documents dry-run default, the five
  states, safe-only apply, text/JSON, and no force switch.
- The glossary and user guides define the same states, stdout/stderr boundary,
  outcome behavior, GC separation, and supported commands.
- The shipped Skill directs Agents to dry-run first, preserve unintegrated,
  dirty, and unknown work, apply only fresh safe proof, and never substitute
  manual Git deletion.
- `cmp`, `make skills-sync-check`, the authorial digest test, and the full
  shipped Skill check all passed.
