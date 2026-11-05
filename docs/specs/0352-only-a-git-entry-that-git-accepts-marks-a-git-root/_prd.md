---
spec: 0352-only-a-git-entry-that-git-accepts-marks-a-git-root
status: active
created: 2026-11-05
surfaces: [backend, docs]
type: fix
---

# Only a .git entry that Git accepts marks a Git root

Minimal PRD for a bug fix, minted by `write-techspec`. Reported by the
maintainer on 2026-11-05: on his MacBook, Roundfix took the home directory for
a Git repository because `~/.git` existed, although that directory held only
`gk/` and `info/`, with no `HEAD` and no `objects/`, and `git -C ~ status`
answered `fatal: not a git repository`. Every command that creates a Run
Worktree then refused the default `~/.roundfix/worktrees` as inside the
repository tree, the Read-Only Run Commands printed that refusal as a warning,
and `roundfix init --scope project` would write Project Config into the home
directory. He renamed the directory, but the tool that owns `gk/` may create it
again. After this Spec Roundfix finds a Git root only where Git finds one: a
`.git` directory that passes Git's repository test, or a `gitdir:` file of a
linked worktree or submodule.

## Sources

One maintainer report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md`](references/2026-11-05-a-stray-home-git-directory-makes-the-home-a-git-root.md),
   `type: fix`, reported by the maintainer and relayed by the operator.

At the authoring base `docs/backlog/` and `docs/findings/` hold no open entry,
so no other source shares this context.

## Prerequisites

None in `requires`. Authored against `38cc5854` (v0.132.0), with no other
Spec active. This Spec owns `internal/config` Git root detection and the
`gittest` repository helper; it also changes the Roundfix `review-runs` and
`runs` references, the configuration and `runs` guides, and ADR-0357's amendment
notice.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change reads existing `.git` entries. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; detection reads the local file system. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0362 (new) records the decision:
  "A `.git` that Git rejects is skipped and the walk goes on to the parent, as
  Git's own discovery does." It amends ADR-0357, whose Decision kept any
  `.git` entry as a marker; task_03 writes that amendment notice. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AZ answer of 2026-11-05 and the cycle AY rules quoted in the authorization
  record; the governed files are the Roundfix skill files task_03 changes and
  the CLI test file and the `gittest` helper task_01 changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0352-only-a-git-entry-that-git-accepts-marks-a-git-root/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/review-runs.md`, `.agents/skills/roundfix/references/runs.md`, `internal/cli/cli_test.go`, `internal/gittest/gittest.go`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A `.git` directory that Git rejects, such as one holding only `gk/` and
  `info/`, marks no Git root; the walk continues to the parent, as Git does.
- A repository made by `git init`, a linked worktree, a submodule and a real
  repository nested under a rejected `.git` resolve to the same Git root as
  before.
- Under a rejected `.git` in the home directory, `runs list` prints no
  warning, `reconcile` no longer refuses `worktree.location` and `init --scope
  project` writes no Project Config into the home directory.
- The project scope of `skills install` and configuration resolution use the
  same rule.
- The Roundfix skill, the configuration and `runs` guides and ADR-0357 state
  the rule.

## Core Features

1. Git's repository test. One walk accepts a `.git` directory whose `HEAD`
   names a ref or an object and whose common directory holds `objects/` and
   `refs/`, or a regular file starting with `gitdir:`; any other `.git` entry
   is skipped.
2. One walk for every caller. `config.Load`, `roundfix init --scope project`,
   `roundfix profiles configure --scope project` and the project scope of
   `skills install` call it.
3. Fixtures that Git would accept. Tests that faked a repository with an empty
   `.git` directory create the smallest `.git` directory Git accepts.
4. The rule in writing: the Roundfix `review-runs` and `runs` references, the
   configuration and `runs` guides, ADR-0362 and the amendment notice in
   ADR-0357.

## Non-Goals / Out of Scope

- `GIT_DIR`, `GIT_COMMON_DIR`, `GIT_OBJECT_DIRECTORY`,
  `GIT_CEILING_DIRECTORIES`, bare repositories and `safe.directory`: the walk
  never read them and still does not.
- The other `.git` readers: `baseline profile`'s repository walk
  (`findBaselineRepositoryRoot`), the Run Database's repository backfill and
  the commands that ask `git rev-parse` for the root; each either asks Git or
  is followed by a Git command that refuses a rejected `.git`.
- The text of the `worktree.location` refusal and the Read-Only Run Command
  warning (ADR-0357), which keep their bytes when a Git root holds the
  location.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, any Baseline module or
  guide under `docs/agents/`, `CHANGELOG.md`.

## Success Metrics

1. Success Metric: with a temporary home whose `.git` directory holds only
   `gk/` and `info/` and the working directory `<home>/drive`, `runs list`
   prints `No Runs found.` with an empty stderr and exits `0`, `reconcile`
   exits `2` with `reconcile requires running inside a Git repository`, and
   `init --scope project` exits `2` with `project init requires a Git root;
   use --scope user outside a repository` and writes no file; at `38cc5854`
   the first prints the worktree-location warning, the second refuses
   `worktree.location`, and the third writes `<home>/.roundfixrc.yml`.
2. Success Metric: for each `.git` shape of `_techspec.md` → Testing Approach
   where `git rev-parse --show-toplevel` prints a top level or answers `not a
   git repository`, the Git root Roundfix finds equals that top level, or is
   empty where Git answers `not a git repository`.
3. Success Metric: a `git init` repository, a linked worktree, a `.git` file
   holding `gitdir:` and a real repository nested under a rejected `.git`
   resolve to the same Git root as at `38cc5854`, and their Project Config is
   read from that root.
4. Success Metric: `skills install --scope project` run under a rejected
   `.git` outside any repository resolves no project root from it.
5. Success Metric: both copies of the Roundfix `review-runs` and `runs`
   references, the configuration guide and the `runs` guide say that only a
   `.git` Git accepts marks a Git root, none of them still says that any `.git` entry is found,
   and ADR-0357 carries the amendment notice naming ADR-0362.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The maintainer's report quoted in the adopted reference, written before this
  Spec existed: the `~/.git` contents (`gk/`, `info/`) and Git's answer
  `fatal: not a git repository (or any of the parent directories): .git`.
- Git's own answer, through the `git` binary the tests already use: Success
  Metric 2 compares Roundfix's root with `git rev-parse --show-toplevel` in
  each shape, so the oracle is Git's implementation, not this Spec's reading
  of it.
- The Git documentation, gitrepository-layout(5)
  (https://git-scm.com/docs/gitrepository-layout, read through Exa on
  2026-11-05): "a valid Git repository must have the HEAD file", and "you can
  have a plain text file `.git` at the root of your working tree, containing
  `gitdir: <path>` to point at the real directory that has the repository".

## Unreachable Acceptance

None. Every criterion is checked with this repository's files, temporary
homes and temporary repositories.

## Validation

The decisions were checked on 2026-11-05 against these sources.

- Characterization with the installed `roundfix` 0.132.0 (`38cc5854`), a
  temporary `HOME` and `TMPDIR` and no network, from `<home>/drive` with
  `<home>/.git/gk/` and `<home>/.git/info/`:
  - `runs list` exits `0` with `roundfix: warning: worktree.location must not
    be inside the repository tree: "<home>/.roundfix/worktrees" is inside the
    Git root "<home>" (found "<home>/.git"); commands that create a Run
    Worktree refuse it` and `No Runs found.`;
  - `reconcile` exits `2` with `roundfix: reconcile failed: worktree.location
    must not be inside the repository tree: ...`;
  - `init --scope project` exits `0` and writes `<home>/.roundfixrc.yml`;
  - `git rev-parse --show-toplevel` there answers `fatal: not a git
    repository (or any of the parent directories): .git`.
- Git 2.56.0 on the same machine, in a scratch repository with a subdirectory
  per shape: a `.git` directory holding only `gk/` and `info/`, one with
  `HEAD` and `objects/` but no `refs/`, and one whose `HEAD` holds `junk` are
  skipped, and Git reports the outer repository; a `.git` directory with a
  `ref: refs/heads/main` `HEAD`, `objects/` and `refs/` is a root; a `.git`
  file holding `garbage` fails with `invalid gitfile format`, and one holding
  `gitdir: /nonexistent/x` fails with `gitfile does not point to a valid
  repository`, both at that directory.
- Git's `setup.c` (`is_git_directory`, read from github.com/git/git at
  `master` on 2026-11-05): a Git directory needs a valid `HEAD`
  (`validate_headref`), then `objects` and `refs` under its common directory.
- A prototype of the design on a scratch copy at `38cc5854` passed the whole
  module's tests with the fixtures changed as `_techspec.md` → Declared breaks
  lists, and the three transcripts of `_techspec.md` were reproduced with its
  binary.
- Secondbrain, read only: `wiki/index.md`, then `qmd query "git repository
  discovery stray .git directory home roundfix findGitRoot worktree location"
  --all --files --min-score 0.3`. It returned the mirror of ADR-0357 (which
  rejected this change for its own scope), the `runs`, `gc` and `reconcile`
  guides and the digest
  `wiki/sources/digest-roundfix-worktree-git-sob-symlink-e-teste-de-lock-sem-prazo-2026-10-12.md`;
  none decides what marks a Git root.
- Exa: gitrepository-layout(5) at git-scm.com and gitcore-tutorial(7) at
  kernel.org, on the files a Git directory holds.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ brief for Spec 0352 of 2026-11-05 and the maintainer's answer below.

The maintainer's answer (AskUserQuestion, 2026-11-05 ~07:00 BRT), verbatim:

- "Spec agora + v0.132.1 (Recommended)" — option text: "findGitRoot passa a
  aceitar só um .git que o Git reconhece: pasta com HEAD e objects/, ou
  arquivo 'gitdir:'. Um teste cobre um ~/.git falso. Mesmas regras do AY
  (merge pela fila, notas aos 8 adotantes)."
- Question context, verbatim: "O ~/.git do MacBook só tinha as pastas gk/ (de
  10/06) e info/ (de 07/10), sem repositório. Se a ferramenta que usa gk/
  recriar ~/.git/gk, o Roundfix volta a recusar no MacBook, porque o
  findGitRoot aceita qualquer .git. Como tratar a correção no Roundfix?"

Taken from that answer, the operator's brief and the characterization, with
no further interview:

- The test is Git's own (`HEAD`, then `objects/` and `refs/` in the common
  directory), slightly stricter than the option text, so Roundfix never takes
  for a repository a directory Git rejects.
- A `gitdir:` file marks the root even when its target is missing, because
  Git stops there too; a regular file without `gitdir:` is skipped.
- Detection reads the file system and starts no `git` process.
- The project scope of `skills install` uses the same walk, because it held a
  copy of the same rule.
- ADR-0362 records the decision and amends ADR-0357. Three implementation
  Tasks in series and a QA gate.

## Open Questions

None.
