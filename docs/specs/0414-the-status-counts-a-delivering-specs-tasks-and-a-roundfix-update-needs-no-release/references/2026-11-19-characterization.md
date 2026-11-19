# Characterization of Spec 0414 — 2026-11-19

Measured by the operator and the planning author on 2026-11-19, on the planning
clone of `docs/plan-0414` at `145d6c76` (v0.148.1) and on a read-only clone of
the Fiscus repository at `ad1381b`. A prototype of both fixes ran on a scratch
worktree of the same commit, outside this repository, and was removed. No
Fiscus file is copied here; only commit identifiers, subjects and changed path
names are quoted.

## 1. The status view counts a delivering Spec's Tasks on the default branch

`followUpQueueRows` in `internal/cli/follow_up_queue.go` fills a row's
`tasks_completed` and `tasks_total` with `spec.Load(in.SpecRoot, slug)` for
every row that is not authoring, Delivery Queue items included. `in.SpecRoot` is
the Spec Root of the checkout the command runs in. `followUpPendingSpecs` in
`internal/cli/follow_up_status.go` counts the `Pending Specs` lines the same
way. The Daemon settles each Task in the Run's work directory and commits the
Task file with its status (`settleTask`, then the Task commit, in
`internal/daemon/task_engine.go`), so the progress lives on the Item Branch
while the default branch keeps the plan with every Task `pending`.

Reproduction on a temporary repository built with `newFollowUpReadInput`: a
Spec with two pending Tasks committed on `main`; an Item Branch
`roundfix/deliver-0097-widget-0123456789abcdef` in a linked worktree whose
commit settles `task_01` completed; a Delivery Queue item at stage `running`
naming that branch and worktree.

| Tree | Row status |
| --- | --- |
| `145d6c76`, unchanged | `running 0/2 Tasks` |
| prototype, Item Branch read through Git | `running 1/2 Tasks` |

A second case archives the Spec on the Item Branch after three completed Tasks
(a corrective Task appended on the branch) and puts the item at `reviewing`:
the unchanged tree prints `reviewing 0/2 Tasks`; the prototype, reading the
folder at the Archive Record's `source_revision`, prints `reviewing 3/3 Tasks`.

The prototype added `spec.LoadAt(specsRoot, slug, read)`, the Task Graph loader
with its three file reads (`_prd.md`, `_tasks.md`, each Task file) taken from a
function, and kept `spec.Load` on `os.ReadFile`. With it, every test of
`./internal/spec` and every `./internal/cli` test matching
`FollowUp|Queue|SpecIsBeingAuthored|RenamedSpec` passed unchanged, including
`TestTheQueueRowsComeInTableOrder` and `TestAQueueItemNamesItsWait`, whose
items have no branch and keep `queued 1/2 Tasks` and `parked 0/0 Tasks`.

The `follow-up` user guide states: "it reads no file of another checkout: for
the Spec queue, the status view asks other directories only for their `origin`
remote." Reading the Item Branch through Git objects of the same repository
keeps that sentence true; reading the item worktree's files would not.

## 2. The release plan asks to classify a Roundfix update by hand in an adopter

`roundfix release plan --from v0.9.0 --to HEAD --format json`, with the binary
built from `145d6c76`, in the Fiscus clone, exit 3:

| Commit | Subject (shortened) | crosses the boundary |
| --- | --- | --- |
| `e48a461` | docs: propose Spec 0098 | false |
| `4b83839` | chore: update the Roundfix skills and Baseline to 0.148.0 | true |
| `ec8481d` | chore: adopt the Roundfix 0.148.0 recommendations | true |
| `f923d37` | docs: plan Spec 0097 | true |
| `81286a9` | fix: unpark Spec 0097 | false |
| `ad1381b` | docs: plan Spec 0099 | true |

State `manual_classification_required`, blocking `4b83839`, `ec8481d`,
`f923d37`, `ad1381b`. The paths that cross, by commit:

- `4b83839`: `.agents/skills/{roundfix,setup-context-driven,write-idea,write-prd,write-tasks,write-techspec}/**`, `AGENTS.md`, `docs/agents/agent-instructions.md`, `docs/agents/autonomous-work.md`, `docs/agents/setup-context.json`, `skills-lock.json`.
- `ec8481d`: four authoring skills under `.agents/skills/`, five generated `.claude/agents/roundfix-*.md`, `.roundfixrc.yml`, `AGENTS.md`, `docs/agents/{agent-instructions,autonomous-work,backend,secondbrain}.md` and `docs/agents/specific-repository.md`.
- `f923d37`: `packages/backend/src/contexts/tax/CONTEXT.md`.
- `ad1381b`: `CONTEXT-MAP.md`, `docs/agents/specific-repository.md`, six files under `docs/user-guide/` and `packages/backend/.env.example`.

Facts about the records at the Fiscus target:

- `skills-lock.json` (version 1) records 100 skills, among them every Roundfix-owned skill (`roundfix`, `qa-gate`, `write-tasks`, `implement-task`, `archive-spec`, `setup-context-driven`), each with a `computedHash` and no `ref`.
- The Setup Manifest `docs/agents/setup-context.json` records 25 managed artifacts over 13 paths: `AGENTS.md` and 12 guides under `docs/agents/`. It does not record `docs/agents/specific-repository.md`, which the Baseline's repository extension writes at adoption and the repository then owns.
- A repository-local skill, `revisao-tributaria-cadastro`, lives under `.agents/skills/` with scripts and is not in `skills-lock.json`.

Facts about this repository at `145d6c76`:

- `skills-lock.json` (version 1) records 33 skills, all Vendored Skills of `marcioaltoe/skills`; no Roundfix-owned skill is recorded, because this repository authors them.
- The Setup Manifest records `AGENTS.md` and 11 guides under `docs/agents/`.

## 3. The rule measured on both histories

A script applied the rule of `_techspec.md` → API Contract 2 to the changed
paths of every commit, with the managed set read at the range's base and target
(the union) and the adopter test at the target.

- Fiscus, `v0.9.0..HEAD`: `4b83839`, `ec8481d` and `f923d37` stop crossing;
  `ad1381b` still crosses on `CONTEXT-MAP.md`, `docs/user-guide/**` and
  `packages/backend/.env.example`. Without `docs/agents/specific-repository.md`
  in the set, `ec8481d` would still cross.
- Fiscus, `v0.1.0..HEAD`: 18 of the 63 commits that cross today stop crossing,
  all of them Roundfix updates, Baseline applications and plans that touch a
  `CONTEXT.md`; 45 still cross.
- Roundfix, `v0.125.0..HEAD` (216 commits): with the adopter test, no commit
  changes. Applying the rule to every repository instead would have cleared one
  commit, `1a8e08b4` ("chore: turn on the declared defect class and the
  authoring fallbacks in this repository"), which changes only
  `.roundfixrc.yml`.

The prototype of the classifier (an optional `ManagedPathSource` read by
`releaseplan.Build`, reading `skills-lock.json` and the Setup Manifest at the
base and the target with `git ls-tree` and `git show`) was built with
`make build` and run:

| Repository and range | unchanged binary | prototype |
| --- | --- | --- |
| Fiscus `v0.9.0..HEAD` | manual, blocking 4 | manual, blocking `ad1381b` only |
| Roundfix `v0.148.0..HEAD` | manual, 1 blocking, 3 crossing | identical |
| Roundfix `v0.147.0..HEAD` | manual, 1 blocking, 11 crossing | identical |
| Roundfix `v0.140.0..HEAD` | manual, 2 blocking, 29 crossing | identical |

Every test of `./internal/releaseplan` and every `./internal/cli` test
matching `ReleasePlan|Release` passed unchanged on the prototype.

## 4. Knowledge sources consulted

- Secondbrain: `wiki/index.md` read; `qmd query "release plan manual
  classification adopter skills update" --all --files --min-score 0.3` failed
  twice with `SqliteError: database is locked`, so the wiki was searched
  directly. `wiki/entities/roundfix.md` lists the repositories that adopted
  the Baseline (Vortex, Fluxus, Oraculum, Conexus, Fiscus, Quaestor, Argus,
  Pantheon, Indexus, Onioncry and others): every one of them meets defect 2 at
  every Roundfix update. No page records a decision on the release boundary or
  on the status counts.
- Exa: release-please's `exclude-paths` option, merged in
  googleapis/release-please#1875
  (https://github.com/googleapis/release-please/pull/1875), states the same
  rule: "Idea is to exclude commits if all files in those commits are from
  specific paths." Its issue #2891
  (https://github.com/googleapis/release-please/issues/2891) shows the trap of
  that rule: an empty commit has no changed file, `[].every(...)` is true, and
  the commit was dropped. This repository's classifier already treats a commit
  without changed paths as crossing the boundary, and this Spec keeps that.
  release-please reads the excluded paths from its own configuration file;
  this Spec derives them from the records Roundfix writes, so the Release Plan
  Command still reads no Roundfix configuration.
