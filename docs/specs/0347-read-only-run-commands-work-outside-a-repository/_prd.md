---
spec: 0347-read-only-run-commands-work-outside-a-repository
status: active
created: 2026-11-04
surfaces: [backend, docs]
type: fix
---

# Read-only Run commands work outside a repository

Minimal PRD for a bug fix, minted by `write-techspec`. Reported by the
maintainer on 2026-11-03: from `~/Google Drive/My Drive`, a directory that is
not a Git repository of his, `roundfix runs list` and `roundfix attach` both
refused with `worktree.location must not be inside the repository tree:
"/Users/marcio/.roundfix/worktrees"`. The likely cause is a stray `.git` in an
ancestor such as the home directory, which makes every directory beneath it a
"repository" whose tree holds the default worktree location; neither command
creates a Run Worktree, and the refusal names neither the Git root it found nor
the `.git` that made it one. After this Spec the commands that only read the
Run Database run from any directory, `runs list` outside a repository lists
every repository, and the refusal names the Git root and the `.git` path so
the user can find the stray one.

## Sources

One maintainer report, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md`](references/2026-11-03-runs-list-and-attach-refuse-outside-a-repository.md),
   `type: fix`, reported by the maintainer and relayed by the operator.

At the authoring base `docs/backlog/` holds one open entry, on the Waiting
Declaration in the Baseline contracts, which Spec 0346 adopts; it shares no
context with this one. `docs/findings/` holds no open entry on Run discovery or
`worktree.location`.

## Prerequisites

None in `requires`. Authored against `6a0ed904` (v0.130.1). Cycle AX delivers
0345 → 0348 → 0347 → 0346 in one serial queue, so 0345 and 0348 merge before
this Spec. Files this Spec may share with them, each changed narrowly so a
later plan rebases cleanly:

- the Roundfix owned skill version (`.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, `skills/testdata/owned-skill-versions.json`),
  which 0345 also raises for `references/implement.md`; the record command
  writes the next free version, so no number is named here;
- `docs/references/behavior-surfaces.json` and
  `docs/references/coverage-record.json`, re-recorded only with their
  commands after the last edit;
- `CONTEXT.md`, where this Spec adds one entry.

This Spec owns `.agents/skills/roundfix/references/runs.md` and
`internal/config`; it also changes `references/events.md` and
`references/review-runs.md` of the Roundfix skill, which no other Cycle AX
Spec owns. It edits no file owned by 0345, 0346 or 0348.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the change reads existing Run rows and configuration. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential; the commands read local files and the local Run Database.
  Source: `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0357 (new) records the decision:
  "The commands that only read the Run Database run from any directory." The
  location keeps its refusal for every command that creates, inspects or
  removes a Run Worktree. Source: `docs/agents/domain.md`,
  `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AX answers of 2026-11-04 quoted in the authorization record; the governed
  files are the Roundfix skill files task_04 changes, the CLI test file
  task_03 removes a test from, and the coverage record it re-records. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0347-read-only-run-commands-work-outside-a-repository/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/events.md`, `.agents/skills/roundfix/references/review-runs.md`, `.agents/skills/roundfix/references/runs.md`, `docs/references/coverage-record.json`, `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- `roundfix runs list` run outside a Git repository lists the Runs of every
  repository, as `runs list --all` does, and exits `0`.
- Bare `runs`, `runs list`, `runs show`, `attach` and `events` never fail on a
  `worktree.location` refusal; they print it once on stderr as a warning and
  keep their output, including under a stray `.git` file or directory.
- The `worktree.location` refusal keeps its text as a prefix and names the
  detected Git root and the `.git` path it found; every command that needs the
  location keeps refusing.
- Inside a Git repository the five commands keep their output byte for byte.
- The Roundfix skill, the command guides, the configuration guide and
  `CONTEXT.md` state the contract.

## Core Features

1. The refusal names its Git root. `ResolveWorktreeLocation` adds the root and
   its `.git` path after the existing text.
2. Optional worktree location. `config.Load` can return the refusal as data
   instead of failing, for the commands that never create a Run Worktree.
3. Read-Only Run Commands. The five commands load with that option and print
   the refusal as one warning.
4. Runs List outside a repository. With no Git root, `runs list` lists every
   repository.
5. The rules in writing: the Roundfix `runs`, `events` and `review-runs`
   references, the `runs`, `attach`, `events` and configuration guides, the
   `runs` help, and `CONTEXT.md`, which gains **Read-Only Run Command**.

## Non-Goals / Out of Scope

- How the Git root is detected (`findGitRoot`): any `.git` entry, directory or
  file, still marks a repository, so a linked worktree or a submodule keeps
  working; recognizing a `.git` only when Git would is left out (see
  `_techspec.md` → Decisions).
- `runs causes`, which reads the repository's Archive Records and keeps
  requiring a repository, and `stop`, `reconcile`, `gc`, `implement`,
  `deliver` and `doctor`, which keep refusing a location inside the root.
- A JSON output for `runs list`; it has none.
- The Makefile, the CI workflows, `go.mod`, `go.sum`, any Baseline module or
  guide under `docs/agents/`.

## Success Metrics

1. Success Metric: from a directory outside any Git repository, with a Run
   Database holding Runs of two repositories, `roundfix runs list` prints the
   same stdout and stderr as `roundfix runs list --all` run inside one of them
   and exits `0`; at `6a0ed904` it exits `2` naming `--all`.
2. Success Metric: with a `.git` directory or a `.git` file in the home
   directory and the working directory beneath it, `runs list --all`,
   `runs show <id>`, `attach <id>`, `events <id>`, and the Run Browser through
   bare `runs` and `attach`, exit as they do inside a repository with the same
   stdout, and their stderr is the warning line followed by the stderr they
   print inside a repository; `runs list` there prints `No Runs found.`; at
   `6a0ed904` each of them exits `2` with the refusal.
3. Success Metric: the refusal reads `worktree.location must not be inside the
   repository tree: "<location>" is inside the Git root "<root>" (found
   "<root>/.git")` for a `.git` directory, a `.git` file and a location equal to
   the root, and starts with the text it had at `6a0ed904`; `reconcile` under a
   stray `.git` still exits `2` with it.
4. Success Metric: inside a Git repository `runs list`, `runs show <id>`,
   `attach <id>` and `events <id>` print what they printed at `6a0ed904`, with
   no warning.
5. Success Metric: both copies of the Roundfix `runs` reference state the
   outside-repository listing and the warning, the `runs` guide and help state
   that outside a Git repository it lists every repository, the `attach`,
   `events` and configuration guides name the warning or the refusal's Git
   root, and `CONTEXT.md` defines **Read-Only Run Command**.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The maintainer's report quoted in the adopted reference, written before this
  Spec existed: the commands he ran, the directory, and the exact refusal.
- The Git documentation, git(1), environment variable
  `GIT_DISCOVERY_ACROSS_FILESYSTEM` (https://git-scm.com/docs/git, read through
  Exa on 2026-11-04): "When run in a directory that does not have ".git"
  repository directory, Git tries to find such a directory in the parent
  directories to find the top of the working tree". A `.git` in an ancestor is
  what makes a directory part of a working tree, so naming the `.git` found is
  what lets the user locate it.

## Unreachable Acceptance

None. Every criterion is checked with this repository's files, temporary
homes and temporary repositories.

## Validation

The decisions were checked on 2026-11-04 against these sources.

- Characterization with `bin/roundfix` built at `6a0ed904` (v0.130.1), with a
  temporary `HOME` and `TMPDIR` and no network:
  - outside any repository, no stray `.git`: `runs list` exits `2` with
    `runs list requires a Git repository; pass --all to list every
    repository`; `runs list --all` prints `No Runs found.` and exits `0`;
    `attach <id>`, `runs show <id>` and `events <id>` reach the Run Database;
  - inside a repository: `runs list` and `runs list --all` print `No Runs
    found.` and exit `0`;
  - with an empty `.git` directory in the temporary home and the working
    directory beneath it: `runs list`, `runs list --all`, `attach <id>`,
    `runs show <id>`, `events <id>` and `reconcile` each exit `2` with
    `worktree.location must not be inside the repository tree:
    "<home>/.roundfix/worktrees"`; an empty `.git` file does the same. The
    maintainer's report reproduces.
  - `git rev-parse --show-toplevel` in the same directory answers `fatal: not
    a git repository`, so Git does not take an empty `.git` for a repository
    while Roundfix does.
- A prototype of the design on a scratch copy at `6a0ed904` passed the new
  tests and the whole `./internal/cli` and `./internal/config` packages, and
  each sabotage named in `_techspec.md` → Testing Approach made its named test
  fail.
- Reading at `6a0ed904`: `findGitRoot` accepts the first ancestor holding any
  `.git` (`os.Stat`); `Load` ends with `ResolveWorktreeLocation`, whose error
  fails every command that loads configuration; `runRunsListCommand` refuses a
  missing Git root without `--all`; `runs`, `runs list`, `runs show`, `attach`
  and `events` load configuration through `loadCommandConfig` and never read
  `worktree.location`; `stop`, `reconcile`, `implement` and `deliver` read it.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix runs list attach outside repository stray .git home worktree location" --all --files --min-score 0.3`.
  It returned mirrors of the Roundfix `runs` reference, the command guides and
  two worktree decision records; none decides behavior outside a repository.
- Exa: the git(1) passage quoted above and the Pro Git chapter "Git Internals
  - Environment Variables" (https://git-scm.com/book/en/v2/Git-Internals-Environment-Variables).

## Glossary

- adds: **Read-Only Run Command** — a command that only reads the Run Database and runs from any directory

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AX brief for Spec 0347 of 2026-11-04 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-04), verbatim:

- 0347: "Spec neste ciclo (Recommended)" — option text: "runs list/attach
  funcionam fora de repositório (leem o store global), e a recusa do
  worktree.location nomeia a raiz Git detectada e o caminho do .git, para o
  usuário achar o .git perdido."
- Rules (all four): "Mesmas do AW (Recommended)" — option text: "Sem
  supervisão: push e PR dos planos, entrega e merge pela fila, consertar e
  seguir em falha técnica; Makefile, .github/workflows, go.mod/go.sum e lint
  liberados quando o desenho pedir; mudança de Baseline que chega aos
  adotantes liberada."
- Release (all four): "Uma v0.131.0 no fim (Recommended)" — option text: "Uma
  minor (a mudança da Baseline e o runs list fora de repo são feat); notas aos
  8 adotantes, resposta ao Pantheon com o defeito corrigido e o contorno que
  podem largar."

Taken from those answers, the operator's brief and the characterization, with
no further interview:

- The Read-Only Run Commands are bare `runs`, `runs list`, `runs show`,
  `attach` and `events`: the five that read only the Run Database. `runs
  show` and `events` join the two the maintainer named because they fail the
  same way and need the location no more.
- Outside a repository `runs list` lists every repository; inside one,
  including a stray one, its scope is unchanged.
- The refusal becomes a warning on those commands, never silence, so a stray
  `.git` stays visible.
- `findGitRoot` is unchanged.
- ADR-0357 records the decision. Four implementation Tasks in series and a QA
  gate.

## Open Questions

None.
