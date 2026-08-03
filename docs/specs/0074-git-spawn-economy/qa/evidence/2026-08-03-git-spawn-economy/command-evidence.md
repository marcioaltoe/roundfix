# Command evidence — Git spawn economy QA

Build: `4fd7539f20ed9a25c526a99d15ea9e9aa1ebbdd0`

## QA-01 — Task Graph and Project Constraints

- `_tasks.md` declares `qa: task_07`; `task_07` is the sole `type: qa` node,
  has no dependent node, and needs `task_06`. `task_06` transitively covers
  `task_01` through `task_05`.
- `rtk proxy rg -n '^status:|^type:'
  docs/specs/0074-git-spawn-economy/task_*.md` showed `task_01` through
  `task_06` as `completed`; only the Daemon-owned `task_07` remains `pending`.
- `rtk git log --reverse --format=oneline dbdad8a..HEAD` showed the six Task
  commits in Task order: `dd7d472`, `83a33e7`, `41ca137`, `73bd134`,
  `a9d0097`, and `4fd7539`.
- `rtk git diff-tree --no-commit-id --name-status -r <commit>` ran for all six
  commits. Task 01 and Task 06 changed only their Task file and Spec-local
  measurement artifacts; Tasks 02 through 04 changed only their Task file and
  `internal/baseline/`; Task 05 changed only its Task file and
  `internal/agent/`.
- `rtk git diff --name-only dbdad8a..HEAD` independently listed the same 27
  paths. It contained no linter, formatter, typechecker, test-runner,
  architecture-checker, build-tool, package-manager, code-generator,
  repository-tooling configuration, script, ignore file, plugin declaration,
  skill, workflow, or version pin. The only executable script is the
  Spec-local measurement shim under `baseline/`; it is QA evidence, not
  repository tooling.
- `rtk git -c core.fsmonitor=false status --short` showed only the new
  Spec-local `qa/` directory after the report was opened.
- The PRD and TechSpec both account for all four Project Constraints:
  identifier strategy and authentication/HTTP are not applicable with reasons;
  active ADR obligations identify accepted ADR-0089 and ADR-0090 and explain
  why ADR-0081 is not triggered; tooling authority is not applicable because
  no protected tooling surface is changed. Each row cites its operative source
  under `docs/agents/`.

Result: pass. No graph defect, missing constraint, tooling authorization
problem, out-of-scope Task path, or current non-QA delta was found.

## QA-02 — Repository Verification

`rtk make verify` exited 0 on build
`4fd7539f20ed9a25c526a99d15ea9e9aa1ebbdd0`:

```text
Go test: 3125 passed in 24 packages
Go test: 4 passed in 1 packages
Roundfix skill check passed: roundfix, write-idea, write-prd, write-techspec,
write-tasks, setup-context-driven, implement-task, implement-spec,
brainstorming, council, business-analyst, archive-spec, qa-gate,
evidence-gate
go build ... -o bin/roundfix ./cmd/roundfix
```

The gate's final build command exited 0 and produced `bin/roundfix`. Result:
pass.

## QA-03 — Reproducible spawn census baseline

- `rtk proxy sh -n .../baseline/git-census-shim`: exit 0.
- `rtk proxy test -x .../baseline/git-census-shim`: exit 0.
- `rtk proxy awk -f .../baseline/attribution.awk /dev/null`: exit 0 and
  reported zero in all four buckets, including malformed records.
- `rtk git cat-file -e <revision>^{commit}`: exit 0 for both published
  measurement revisions, `dbdad8ac1b8a2335ab88c65a0a47f50d86ef6c4e`
  and `a9d0097c590412dafd173fbcb4deaf1923bcae3a`.
- `baseline/README.md` records the exact census/timing commands, environment,
  attribution limits, and baseline results. `baseline/after.md` reuses those
  commands and identifies its revision and same host class.

Result: pass. The artifact is executable and reproducible; the fresh current
timing result is tracked separately in QA-08.

## QA-04 — Skills restore batches object reads

The current `bin/roundfix` exercised an offline multi-file restore against a
disposable clone with `agentic-cli-design` removed and
`/Users/marcio/dev/skills` as the immutable source.

1. Preview exited 3 with `plan.confirmation.required`, five planned file
   creates, one lock edit, and Plan Digest
   `4ca9974fcd7c66de92eed1a01ceb13cfc8d254487ee88116d55eb70930220f8c`.
2. Confirming that exact digest exited 0 with `ok: true`, `applied: true`, the
   same plan digest, and `restore.completed`.
3. A fresh-process repeat exited 0 with `ok: true`, `applied: false`, and no
   planned changes.
4. `test -f` confirmed the restored fifth file at
   `.agents/skills/agentic-cli-design/references/templates.md`; target status
   showed only the expected lock update because the restored skill files were
   byte-identical to tracked content.
5. The committed PATH shim recorded 26 Git invocations across the three
   public commands. Each command recorded exactly one `ls-tree` and one
   `cat-file`, so each five-file read scope cost one object-reader process,
   not five.
6. Omitting `--profile` exited 2 with the public
   `restore.arguments-invalid` finding and an actionable rerun instruction.

Result: pass. Confirmation safety, apply, independent filesystem confirmation,
restart idempotence, invalid input, and operation-proportional object reads
were all observed through the public CLI.

## QA-05 — Assets sync batches provenance reads

- Two first probes against the live `/Users/marcio/dev/skills` checkout
  consistently exited 2 because that checkout had moved beyond this binary's
  pinned catalog (`bubbletea` and `tui-design` were outside the expected
  setup). The public error identified `skills.setup-snapshot.drift` and the
  exact corrective action. This was retained as stale-source recovery
  evidence, not treated as current-build success evidence.
- A disposable clone was detached at the immutable source revision reported
  by the restore plan,
  `14fdf46befa9a07203fde20b69b885dab4961844`, and its origin was normalized
  to `https://github.com/marcioaltoe/skills.git`.
- Three fresh-process public commands against that exact source exited 0 with
  byte-identical JSON: `ok: true`, zero errors/decisions/warnings/info, no
  findings, and no planned changes.
- `rtk git -C <source> status --short` remained empty after the checks,
  independently confirming `--check` was read-only.
- A dedicated census file for one successful command recorded: 2
  `rev-parse`, 1 `rev-list`, 1 `status`, 1 `remote`, 3 `show`, 132 `ls-tree`,
  and 132 `cat-file` processes. The production loop therefore opened one
  batch reader per immutable skill-tree scope; it did not spawn per file.
- The fresh focused seam check
  `rtk go test ./internal/baseline -count=1 -run
  '^TestAssetsSyncCommittedTreeDigestReadsManyFilesThroughOneBatchProcess$'
  -v` passed and directly counted one batch process for a multi-file tree.
- A non-Git archive source exited 2 with the public
  `skills.setup-snapshot.drift` finding, the failed combined `rev-parse`
  command, and an action to fix the canonical source.

Result: pass. Read-only success, restart stability, stale-source recovery,
non-repository failure, and multi-file batching were observed.

## QA-06 — Repository resolution combines fresh facts

- The public skills-restore preview ran through `git-args-shim`. Its first
  target-repository command was exactly one combined invocation:
  `git -C <target> -c core.fsmonitor=false rev-parse --show-toplevel
  --show-object-format --verify HEAD^{commit}`. The only separate identity
  command was Git's incompatible `rev-list --max-parents=0 HEAD` operation.
- The source acquisition separately logged one `rev-parse --verify
  FETCH_HEAD^{commit}`, one `ls-tree`, and one `cat-file --batch`.
- A new disposable clone first resolved at
  `4fd7539f20ed9a25c526a99d15ea9e9aa1ebbdd0`; after a checkout mutation it
  resolved at `a9d0097c590412dafd173fbcb4deaf1923bcae3a`. A fresh public process logged
  the combined query again after mutation. The stable Plan Digest remained
  `a19aba9e...`, which is correct because repository identity is stable across
  commits; the fresh `HEAD^{commit}` read validates usability rather than
  redefining repository identity.
- An empty Git repository exited 2 with `restore.repo-invalid` and the
  distinct `Baseline repository requires at least one commit` message.
- A non-repository directory exited 2 with `restore.repo-invalid` and the
  distinct `detect Git worktree root` message. Both supplied the same public
  next action to pass an existing Git worktree.
- A planned fixture commit was not used: inherited GPG signing tried to write
  under sandboxed `~/.gnupg` and failed. Per sandbox policy, it was not
  retried; the two-existing-commit checkout supplied equivalent mutation
  evidence without signing or user-repository mutation.

Result: pass. Combined positional reads, restart freshness, stable repository
identity, missing-HEAD failure, and non-repository failure were observed
through the public command.

## QA-07 — Agent runner takes its environment explicitly

- `roundfix doctor --help` exited 0 and confirmed the Doctor Command is
  offline, read-only, and mutates nothing.
- A controlled ACPX wrapper around the real `/opt/homebrew/bin/acpx` ran the
  public Doctor Command with `QA_ENV_MARKER=alpha`. The wrapper log recorded
  `marker=alpha --version` twice, proving the process-default environment
  reached each child boundary.
- A fresh Doctor process used `QA_ENV_MARKER=beta` and a controlled old-version
  ACPX response. Its independent wrapper log recorded `marker=beta --version`
  twice, with no stale alpha value. Doctor reported the expected minimum
  version failure and exact update action.
- The real Doctor process reported Node and ACPX ready, then exited 1 for
  machine state outside this Spec: unknown Claude adapter lineage, an invalid
  local Codex signature, and a missing external `knowledge-workspace` skill.
  Each failure included its next action. These readiness facts do not prevent
  observing the environment boundary; no Run or agent session was started.
- `rtk proxy rg -n 't\.Setenv\(' internal/agent -g '*_test.go'` found exactly
  three stated process-default tests. The parallel inventory found 134
  declarations across the package's test files.
- `rtk go test ./internal/agent -race -count=2 -parallel 16` exited 0: 530
  test executions passed with no race report.

Result: pass. Invocation-scoped process-default environments changed across
fresh public processes, public readiness failures stayed actionable, and the
exact Task race/repetition gate passed.

## QA-08 — Fresh suite completes under 60 seconds

The exact committed timing procedure ran with a newly created empty
`GOCACHE=/private/tmp/roundfix-qa-0074-timing.fExVI8/gocache`:

```text
rtk proxy env GOCACHE=<empty> /usr/bin/time -p \
  go test ./... -count=1 -parallel 16
```

The command exited 0, and every package passed. The slowest Go-reported
packages were `internal/cli` at 72.501s and `internal/baseline` at 71.872s.
`/usr/bin/time -p` reported:

```text
real 86.36
user 137.32
sys 243.25
```

Result: fail, F-001 (Friction). The fresh suite exceeded the strict target by
26.36 seconds. This confirms and worsens the committed Task 06 after result
of 83.38 seconds by 2.98 seconds on the same host class.

## QA-09 — Observable behavior remains byte-stable

Status: blocked (environment: exported base build inherited a sandboxed Go
cache).

- `rtk git archive -o <base.tar> dbdad8a...` and extraction both exited 0.
- The first and only base build attempt failed before compilation:
  `open /Users/marcio/Library/Caches/go-build/...: operation not permitted`.
  Per the sandboxed QA rule, it was not retried with another cache.
- Equivalent current-build evidence: `rtk make verify` passed 3,125 package
  tests plus the skills characterization; QA-04 observed public restore
  preview/apply/idempotence and exact restored file state; QA-05 observed
  repeated byte-identical assets-sync JSON and a clean source; QA-06 observed
  exact combined commands and distinct failure surfaces; the current
  `internal/agent` race/repetition gate passed 530 executions.
- Static equivalence evidence: `git diff --quiet dbdad8a..HEAD -- cmd/roundfix
  internal/cli docs/user-guide` exited 0. No public CLI implementation or user
  command documentation changed.
- Unblocking action: in a full-access session, build the already-exported base
  revision with a writable fresh Go cache and compare matched public success
  and failure inputs byte for byte.

## QA-11 — CLI contract did not change

Status: blocked (environment: same exported base build denial as QA-09).

- The source-level CLI contract comparison exited 0 with no diff across
  `cmd/roundfix`, `internal/cli`, or `docs/user-guide`.
- The current binary's root, Baseline, skills-restore, and assets-sync help all
  exited 0 and exposed the documented command names, flags, output contract,
  and exit codes.
- `baseline skills restore --format yaml` exited 2, wrote the diagnostic to
  stderr, and emitted the stable text payload with
  `restore.arguments-invalid` and an actionable next step.
- Equivalent success/exit evidence is also recorded in QA-04 through QA-06.
- Unblocking action: use the QA-09 base binary build in a full-access session
  and compare the four help outputs plus the invalid-format stdout, stderr,
  and exit code.

## QA-10 — Tests were not deleted, skipped, or weakened

- `git diff --diff-filter=D --name-only dbdad8a..HEAD -- '*_test.go'` exited 0
  with no deleted test path.
- The added-line sweep for `t.Skip`, `Skipf`, `testify`, or `mockery` found no
  match (expected search exit 1).
- Static top-level test inventories increased from 1,515 at `dbdad8a` to 1,529
  at current HEAD.
- Test-file numstat showed ten touched test files with 1,012 additions and 152
  removals. The removals are environment-plumbing/parallelization rewrites;
  the inventory grew by 14 and the current full gate passed 3,125 tests.
- The exact changed-seam checks also passed: multi-file skills restore,
  multi-file assets provenance, and 530 race/count agent executions.

Result: pass. No test file or top-level test was deleted, no skip or forbidden
test dependency was added, and fresh verification exercised the expanded
inventory.

## QA-12 — Batching stays bounded; runners stay package-local

- Symbol inspection found the private `batchObjectReader` only in
  `internal/baseline/skills_restore_git.go`. Both skills restore and assets
  provenance close it on success and error paths.
- The changed-file list adds no shared Git package or client. Production
  changes remain inside the owning `internal/baseline` and `internal/agent`
  packages.
- A focused search found no `sync.Map`, `repositoryCache`, `gitClient`, or
  `GitClient` type in the three changed repository-read implementations.
- `rtk go test ./internal/baseline -count=1 -run
  '^Test(SkillsRestoreBatchFailuresKeepReadErrorSurface|AssetsSyncCommittedTreeDigestKeepsTreeAndReadErrors)$'
  -v` exited 0 with seven cases, covering close/read failures and the two
  package-specific error vocabularies.
- QA-06's fresh-process checkout mutation logged a new combined Git query
  after mutation, independently proving no in-process fact cache crossed the
  boundary.

Result: pass. Batching is scoped and closed, mutation triggers a fresh read,
and no shared client extraction landed.

## QA-13 — Pull Request review readiness on exact head

Status: blocked (environment: no Open Pull Request).

- The Roundfix QA prompt explicitly states: `Pull Request: none open; Pull
  Request journeys are environment-blocked.` This is the supervised evidence
  for the absence; the Run Worktree branch was not used to infer or resolve a
  Pull Request.
- Local ancestry evidence: `ma/spawn-economy` resolves to
  `dbdad8ac1b8a2335ab88c65a0a47f50d86ef6c4e`; the QA build resolves to
  `4fd7539f20ed9a25c526a99d15ea9e9aa1ebbdd0`; `git merge-base --is-ancestor
  ma/spawn-economy HEAD` exited 0.
- Equivalent local evidence covers integration ancestry and current
  verification, but it cannot substitute for a reviewer decision, check state,
  or unresolved-thread state on an Open Pull Request.
- Unblocking action: integrate the Run commits onto `ma/spawn-economy`, open
  its Pull Request, obtain review on the exact head, and rerun this row
  read-only.

## QA-14 — QA Report contract

- Frontmatter inspection showed `status: closed`, `verdict: fail`,
  `rows_blocked_environment: 3`, and `rows_blocked_finding: 0`.
- Table parsing counted 10 pass rows, 1 fail row, 2 base-build environment
  blocks, and 1 no-Pull-Request environment block: 14 terminal rows total.
- The pending/skipped table-row search returned no match (expected search
  exit 1).
- `rg --files` resolved the report and all ten evidence artifacts.
- All three QA helper scripts passed `sh -n`.
- `rtk git diff --check` exited 0; status showed only the new Spec-local
  `qa/` directory.

Result: pass. The report path is the day's first collision-safe name and its
closed verdict/count contract matches the result matrix.
