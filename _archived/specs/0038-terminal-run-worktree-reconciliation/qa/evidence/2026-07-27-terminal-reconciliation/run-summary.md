# Terminal Run Worktree reconciliation — execution evidence

Build: `01a16c8d9fe430a4d70d0f8db5c5a135eb1d803b`

Application: `bin/roundfix`, built by `rtk make verify`

Environment: macOS Darwin 25.5.0 arm64; Go 1.26.5; Git 2.55.0.

Fixture root: `/private/tmp/roundfix-qa-0038.DEmQ8W`

The fixture helper used project Store and Worktree APIs only to seed isolated
Git repositories and one isolated Roundfix Home. All behavior checks used the
built public CLI, then confirmed state through fresh Git commands, a second CLI
process, database hashes, or a read-only Store reopen.

## Static gate

- First sandboxed `rtk make verify`: stopped before compilation because the
  sandbox could not read `/Users/marcio/Library/Caches/go-build/...`.
- Exact host-access `rtk make verify`: passed — 2,563 tests in 23 packages,
  four Skill-sync tests, fourteen shipped Skill checks, and `bin/roundfix`
  build.
- Postflight `git diff --check`: passed.
- Postflight status: only the Spec-local QA artifact plus the temporary
  `.qa-tmp/` fixture source; no formatter side effect.

## Classification and dry-run

Repository:
`/private/tmp/roundfix-qa-0038.DEmQ8W/matrix-repo`

`bin/roundfix reconcile` and `bin/roundfix reconcile --format json` both exited
0 and returned, newest-first:

| Run | Stored outcome | Classification | Key evidence |
| --- | --- | --- | --- |
| `run_20260727T195402Z_9cdadf385b4c788e` | Clean | released | Run Worktree and Run Branch absent |
| `run_20260727T195402Z_23fbf5f2e1ca151b` | TimedOut | unknown | target branch could not resolve unambiguously |
| `run_20260727T195402Z_05f7efcfd5cba8d3` | Unresolved | unintegrated | Run head `25b043c…` is not an ancestor of target `55d89c6…` |
| `run_20260727T195402Z_a82861776b685a09` | Stopped | dirty | tracked and untracked changes present |
| `run_20260727T195402Z_ff0fc4bf6471bcb9` | Failed | safe | clean worktree; both heads `55d89c6…` |

Summary:
`total=5 safe=1 unintegrated=1 dirty=1 unknown=1 released=1 applied=0 preserved=3 operational-failures=0`.

The JSON envelope was `roundfix-reconcile/v1`, mode `dry-run`, with the same
five results and summary. `--format=json` before the Run ID also exited 0 and
selected exactly the requested Run.

Run Database SHA-256 before and after all dry-runs:
`9a76468b6a4eb2b79f6d11f6ff6c3f2bac0029be530d66a6ca1d140da717fda7`.
`git worktree list --porcelain` and `git show-ref` were also byte-equivalent
before and after.

## Mixed apply and preservation

`bin/roundfix reconcile --apply` exited 0:

`total=5 safe=1 unintegrated=1 dirty=1 unknown=1 released=1 applied=1 preserved=3 operational-failures=0`.

Fresh Git reads found only the safe Run Worktree and Run Branch absent. The
dirty, unintegrated, and unknown worktrees and branches remained. Independent
content reads returned `tracked dirt`, `untracked dirt`, and the unintegrated
tip `25b043c unique run work`.

## Durable outcomes and idempotency

Repository:
`/private/tmp/roundfix-qa-0038.DEmQ8W/integration-repo`

Repository-wide JSON apply exited 0 with five `safe` results and `applied=5`.
A fresh `runs list --state all --limit 0` showed:

- `run_20260727T195403Z_9897285a8e208732`: Clean, previously
  IntegrationPending.
- `run_20260727T195403Z_756779953cf61ec3`: Unresolved.
- `run_20260727T195403Z_020645731249b5c6`: Failed.
- `run_20260727T195403Z_0f4079f721002bbc`: Stopped.
- `run_20260727T195403Z_a0a2b2dcfe3f1ce1`: TimedOut.

A fresh read-only Store process found one event for the promoted Run:

```json
{
  "source": "daemon",
  "kind": "daemon.outcome",
  "summary": "Run reconciled IntegrationPending to Clean before terminal cleanup.",
  "payload": {
    "action": "cleanup",
    "classification": "safe",
    "current_outcome": "Clean",
    "event": "integration_reconciliation",
    "previous_outcome": "IntegrationPending",
    "run_branch": "roundfix/run-run_20260727T195403Z_9897285a8e208732",
    "run_head": "d0d4148e623582b25f7f702de3a155d456299d97",
    "target_branch": "ma/qa-target",
    "target_head": "d0d4148e623582b25f7f702de3a155d456299d97",
    "worktree": "/private/tmp/roundfix-qa-0038.DEmQ8W/worktrees/integration-repo-bdd31c45/run_20260727T195403Z_9897285a8e208732"
  }
}
```

The second repository-wide apply returned five `released` results,
`applied=0`, and no operational failure. Run Database SHA-256 stayed
`bf3b8a0a54d7b69c6a2c15fc35ec84af27bb000791ae0d4af084458cc97d0115`;
the Git worktree and ref surface also stayed unchanged.

## Runs List discovery

In the listing repository:

- Active view: `No Runs found.` plus
  `(3 terminal Run Worktrees retained; run 'roundfix reconcile' to inspect)`.
- Terminal and all-state repository views: the four existing terminal rows
  plus the same exact retained count of three.
- Machine-wide `--all --state all --limit 0`: six retained terminal Run
  Worktrees across the fixture repositories.
- Released-only integration repository: no retained-worktree note; the
  existing generic hidden-terminal-row note remained.

The full gate's stream assertions confirmed Runs List rows remain on stdout and
the retained-worktree note remains on stderr.

## Refusals and recovery

From the selector repository, these `--apply` selectors all exited 2:

- Active Run `run_20260727T195403Z_5f6793da655b49a3`.
- Review Run `run_20260727T195403Z_f66ef7d1d49059b4`.
- Missing Run `run_missing_qa_0038`.
- Other-repository Run `run_20260727T195402Z_a82861776b685a09`.

Database SHA-256 stayed
`bf3b8a0a54d7b69c6a2c15fc35ec84af27bb000791ae0d4af084458cc97d0115`;
refs and worktrees were unchanged. The no-ID JSON form in that repository
returned an empty current-repository result set.

Invalid format, extra argument, `--force`, and a missing `--format` value all
exited 2. Database SHA-256 stayed
`7e92dac68673c71408fe200efbc6f0a77c581c0684cc34815edd33584d1fa7cb`.

A locked clean worktree produced exit 1, `operational-failures=1`, the exact
next safe action, and the remaining worktree/branch in the refusal. Fresh Git
reads confirmed both remained and the worktree stayed locked.

## Non-Goal probes

- Mixed apply left the unique divergent commit and dirty tracked/untracked
  content unchanged: no integration or repair occurred.
- A sibling Task Worktree and Task Branch named
  `roundfix/run-run_20260727T195403Z_6200a0c37fe85dcd-task_01` remained after
  the safe parent Run Branch was released.
- The Integration Pending reconciliation event remained readable after
  cleanup; reconcile did not prune Run Event Journal or artifact storage.
- Repository-scoped no-ID reconciliation excluded every Run from the other
  fixture repositories.

## Focused current-build evidence

- Worktree classification/apply/reaper families: passed in 6.137s.
- CLI reconcile, retained Runs List, automatic Implement Preflight, and docs
  contract families: passed in 3.975s.
- Store reconciliation family: passed in 0.629s.
- `go test -race ./internal/cli ./internal/worktree -run 'Test.*Reconcile'`:
  passed.
- `go test -race ./internal/worktree -run
  'Test(ApplyTerminalRun|PruneTerminal)'`: passed.

These checks cover the timing-only stale-head/dirty/metadata boundary,
persistence-failure-before-Git boundary, unsafe/symlink path fixtures, and
automatic reaper seam that cannot be paused safely through the public process.

## Tooling and scope audit

Task 07 Daemon commit `16da6754efe30e28228d0498318eb01701d61a7d`
changed exactly:

- `.agents/skills/roundfix/SKILL.md`
- `skills/roundfix/SKILL.md`
- the five maintainer-authorized derived digest pins
- `task_07.md`

`cmp` passed, the `typescript-bun.json` authorial-sync test passed, the built
Skill check passed, and Task 06 changed no protected tooling path.

Task 04 Daemon commit `0d56f4281b06602b6108a8d304885237a897f4d6`
also created a repository-root file named `roundfix`. It is a 20.1 MiB Mach-O
arm64 executable. Base `8ec92ad` has no such path; HEAD stores it as executable
blob `01aaed49a73f1e33f4d2dc8bb9ca51208aaa4db9`. The repository's declared build
output is ignored `bin/roundfix` (`Makefile:68`, `README.md:126`), so this root
binary is an unexplained generated artifact.
