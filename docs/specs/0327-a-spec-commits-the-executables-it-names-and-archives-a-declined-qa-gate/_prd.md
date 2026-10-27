---
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
status: active
created: 2026-10-26
surfaces: [backend, cli]
---

# A Spec commits the executables it names, settle says what it removes, and archive honours a declined QA gate

Minimal PRD for a bug fix, minted by `write-techspec`. Quaestor lost a Run of
Roundfix 0.123.0 because the Task commit refused two launchd scripts the Spec
asked for, saw `roundfix settle` remove the Run Worktree and then an exit 128,
and had to archive two Specs that declare `qa: declined` with a QA override.
After this Spec a Task commit keeps an executable file at a path the Spec's
`_authorization.md` lists under `executables:` and still refuses every other
one, including one inside a new directory; `roundfix settle` names each
worktree it removes; and `roundfix archive` accepts a Declared QA Decline once
every Task is completed, without an override.

## Sources

One Backlog Entry, written and adopted at birth by this Spec's planning change
(see `references/_index.md`):

1. [`references/2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md`](references/2026-10-26-a-spec-cannot-commit-its-executables-settle-exits-128-and-archive-refuses-a-declined-qa.md),
   `type: fix`: part 2 of ai-memory message
   `01a17abd-20fe-7213-ab5c-f1a1bdb043fe` from
   `default/gesttione-solutions-quaestor` (2026-10-26 21:22 BRT), which
   consolidates the messages `01a17994` and `01a179a8`. Part 1 of the same
   message, the Baseline profile-change frictions, belongs to Spec 0328 and
   is out of scope here.

At the authoring base `docs/backlog/` and `docs/findings/` hold no entry. The
planning change writes two open Backlog Entries for neighbours this Spec does
not fix (`_techspec.md` → Current behavior 4 and 8).

## Prerequisites

None in `requires`. Authored against `6a2630e5`, with Specs 0321, 0322, 0324,
0325 and 0326 active under `docs/specs/`; this Spec delivers after 0326. Spec
0322 changes `docs/user-guide/commands/settle.md`, `archive.md`, the Roundfix
`settle` and `archive` references and the `archive-spec` skill; task_03 edits
the same files after it merges and pins none of their bytes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema is
  minted; the authorization record gains one optional list of paths and the
  Archive Record keeps its schema. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, credential or
  network access; the commands read and write a local repository. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0157 (an untracked executable is
  refused from a Task commit; tracked source is work output) is amended by the
  new ADR-0340, which lets the Spec's authorization name the executables a
  Task commit keeps; ADR-0154 (a declined gate keeps its own eligibility and
  acquires no invented override) is the contract the archive fix implements.
  ADR-0251, ADR-0262, ADR-0269 and ADR-0306 cite ADR-0154 or one another
  for other archive and cancel paths (a legacy folder's failed QA, the
  cancelled record, ignored Spec files, a link leaving the repository); this
  Spec leaves each of them unchanged.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AQ batch 5 answers of 2026-10-26 quoted in the authorization record, which
  keep the cycle AJ grant for the Governed Paths each Spec declares; the
  governed files are the owned skills `write-tasks`, `roundfix` and
  `archive-spec` and the historically bounded `internal/spec/archive.go`.
  Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/settle.md`, `.agents/skills/roundfix/references/spec-delivery.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `internal/spec/archive.go`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`. No Makefile, workflow, Go module, lint or formatter change.

## Goals

- A Spec that asks for an executable file delivers it: a Task commit keeps an
  untracked executable at a path its granted `_authorization.md` lists under
  `executables:`, with its mode.
- Every other untracked executable stays refused, also inside a new directory,
  and a Task commit's refusal says how to name one.
- `roundfix settle` exits with the settled outcome and names on stderr each
  worktree it removes.
- `roundfix archive` accepts a Spec whose Task Graph declares `qa: declined`
  with a `qa_reason` once every Task is completed, without an override, as the
  Baseline's docs-layout guide promises.
- The `write-tasks` skill, the Roundfix skill, the `archive-spec` skill, the
  user guides and `CONTEXT.md` say all three rules.

## Core Features

1. Authorized Executables. `_authorization.md` accepts `executables:`, a list
   of exact repository-relative file paths. A Task commit keeps a listed
   untracked executable; its refusal of any other one adds that the Spec's
   `_authorization.md` can list the path. A QA Report commit and a Batch
   commit never read the list and keep their refusal bytes.
2. One file at a time. A Task commit and a Batch commit expand an untracked
   directory into its files before the stage filter, so an executable inside a
   new directory is refused like one beside it.
3. Settle names what it removes. After a successful integration, settle prints
   `Task Worktree removed: <path>` and `Run Worktree removed: <path>` on stderr
   for each worktree it removed, and its exit code stays the settled outcome.
4. Archive of a Declared QA Decline. A Task Graph with `qa: declined` and a
   non-empty `qa_reason` archives when every Task is completed, reads no QA
   Report and records the disposition `no-qa`; `--qa-override` is refused for
   such a Spec, because it already qualifies.
5. The rules in writing. `write-tasks`, the Roundfix references `implement`,
   `settle`, `archive` and `spec-delivery`, the `archive-spec` skill, the
   `archive`, `settle` and `implement` guides, both command usages and
   `CONTEXT.md` describe them.

## Non-Goals / Out of Scope

- Part 1 of Quaestor's message (Spec 0328).
- Settle's own staging, which stages every change with `git add -A` and so
  commits an executable a Task commit refuses, and a tracked file whose mode a
  Task changes: recorded as an open Backlog Entry, not changed here.
- Settle run from inside the Run Worktree, which takes that worktree for the
  user checkout and reports `integration pending`: recorded as an open Backlog
  Entry.
- Globs or directories under `executables:`, executables in QA Report or Batch
  commits, and a repository whose `core.fileMode` is `false`, where Git records
  `100644` whatever the file's mode.
- A Baseline clause for `executables:`; the Makefile, the CI workflows,
  `go.mod`, `go.sum`.

## Success Metrics

1. Success Metric: a Task Cycle with a real Git commit, whose granted
   authorization lists `scripts/launch.sh`, completes the Task and commits
   `scripts/launch.sh` with mode `100755`; with the list removed from the
   filter, the same test fails.
2. Success Metric: an unlisted executable beside a listed one is refused and
   its progress line ends `not valid Work Item output unless the Spec's
   _authorization.md lists the path under executables`; a QA Report commit
   and a Batch commit refusal keep their bytes, even for a listed path.
3. Success Metric: a Task Cycle whose Agent writes `newdir/helper.sh` with
   mode `0755` and no authorization fails the Task naming
   `newdir/helper.sh (executable file)` and commits nothing under `newdir/`;
   at `6a2630e5` the same test completes the Task and commits the file.
4. Success Metric: `executables:` with an absolute path, a `..` segment, a
   glob, a trailing slash, surrounding whitespace, a duplicate or a non-list
   value makes the record malformed, naming the field and the value.
5. Success Metric: `roundfix archive` of a fixture Spec with `qa: declined`,
   a `qa_reason`, every Task completed and no QA Report exits `0` and writes a
   record with `disposition: no-qa` and no `qa_override`; at `6a2630e5` it
   exits `2` with `no QA Report found`. With `--qa-override` it exits `2`
   with `the Spec already qualifies for normal archive`; with a pending Task it
   still exits `2` naming the Task.
6. Success Metric: settling a kept Run Worktree exits `0` with stderr exactly
   `Settle surface: <run worktree>` and `Run Worktree removed: <run worktree>`;
   settling a kept Task Worktree also names `Task Worktree removed:`; an
   integration left pending prints neither and exits `1`.
7. Success Metric: every file `_prd.md` → Core Feature 5 names carries the
   three rules, and `CONTEXT.md` defines **Authorized Executable** and
   **Declared QA Decline** and revises **Run Worktree**, **Settle Command** and
   **Archive Command**.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- Git's own documentation, read through Exa on 2026-10-26.
  `Documentation/config/core.adoc` (https://raw.githubusercontent.com/git/git/master/Documentation/config/core.adoc):
  "`core.fileMode` Tells Git if the executable bit of files in the working
  tree is to be honored. … The default is true (when core.filemode is not
  specified in the config file)." Git's test `t/t3700-add.sh`
  (https://code.googlesource.com/git/+/HEAD/t/t3700-add.sh), "git add: Test
  that executable bit is not used if core.filemode=0": a file made `chmod 755`
  and added is checked to be in the index as `100644`. So a Task commit that
  keeps the file records `100755` under Git's default and `100644` where the
  repository turned the setting off.
- Git's behaviour on the QA host, a tool this Spec did not build: in a shell
  whose working directory was removed, `git rev-parse HEAD` prints `fatal:
  Unable to read current working directory: No such file or directory` and
  exits `128`, measured at authoring. That is the exit Quaestor's caller saw
  after settle removed the Run Worktree it stood in.
- Quaestor's report, quoted verbatim in the adopted Backlog Entry, written by
  another project before this Spec existed.

## Unreachable Acceptance

- criterion: Quaestor, on the release that carries this Spec, delivers a Spec
  that ships a launchd script with mode `0755` and archives a `qa: declined`
  Spec without `--qa-override`.
  reason: it needs another repository's Runs after the release, which no
  Verification or QA row may read.
  satisfied-by: Quaestor's reply through ai-memory to the release message the
  operator sends, recorded session evidence.

## Validation

The decisions were checked on 2026-10-26 against these sources.

- Characterization at `6a2630e5`, in a prototype clone with in-process tests
  and the built binary (`_techspec.md` → Current behavior): the executable
  refusal and its bypass through a new directory; archive's refusal of a
  declined Spec and its acceptance of an override for it; settle's exit `0`
  after removing the Run Worktree, and `integration pending` with exit `1`
  when run from inside it.
- A prototype of the three Tasks' code in that clone passed every new test,
  failed each new test under its sabotage, and, across the full
  `internal/daemon`, `internal/authorization` and `internal/cli` packages,
  moved only the existing tests `_techspec.md` → Declared breaks names.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix executable file refused task commit mode 0755" --all --files --min-score 0.3`
  (the Roundfix `deliver` reference mirror, the digest
  `wiki/sources/digest-roundfix-git-add-f-apaga-arquivo-2026-09-17.md` of an
  earlier exit 128 that came from Roundfix's own `git add`, the mirrored
  decision record on queue commit subjects),
  `qmd query "roundfix archive qa declined no QA Report" --all --files --min-score 0.3`
  (Archive Records 0300, 0070 and 0294, the mirrored decision record on
  legacy archive folders, the `archive-spec` and `spec-delivery` mirrors) and
  `qmd query "roundfix settle removes Run Worktree exit code" --all --files --min-score 0.3`
  (the mirrored decision records on per-Run worktrees, their bootstrap and
  their terminal reconciliation, the `implement` and `stop` mirrors). None
  records a decision on executables, a declined gate's archive, or settle's
  removal; the 2026-09-17 digest showed that a 128 can come from Roundfix, so
  the settle exit was measured rather than assumed.
- Exa: the Git documentation and test quoted above, and the commit that added
  `git add --chmod` (git/git `d64217b`), which states the executable bit is
  not detected where `core.filemode` is false.

## Glossary

- adds: **Authorized Executable**
- adds: **Declared QA Decline**
- changes: **Run Worktree**
- changes: **Settle Command**
- changes: **Archive Command**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AQ batch 5 brief for Spec 0327 of 2026-10-26 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-10-26, about 22:35 BRT), verbatim:

- Grouping: "Duas Specs (Recomendado)" — 0327 Run lifecycle (executable
  scripts, settle exit 128, archive with `qa: declined`); 0328 Baseline
  profile-change frictions.
- Executables: "_authorization.md nomeia os executáveis (Recomendado)" —
  option text: "O commit da Task aceita 0755 só nos caminhos que a
  autorização da Spec lista; todo o resto continua recusado. O write-tasks e a
  recusa passam a citar essa regra."
- Authorization: "Mesmas regras do ciclo AQ (Recomendado)" — option text:
  "implement, commit, push, PR e merge nos Governed Paths que cada Spec
  declarar; fora Makefile, .github/workflows, go.mod/go.sum e lint. A release
  continua a meu critério, pelos critérios gravados."

Taken from those answers and the characterization, with no further interview:

- The field is `executables:` in `_authorization.md` frontmatter: exact
  repository-relative file paths, read from the authorization the Run resolved
  at start, independent of `paths:` (ADR-0340).
- QA Report and Batch commits never read the list.
- The filter judges files: a new directory is expanded before it.
- Settle keeps removing the worktrees it integrated and names each one; its
  exit code was already the settled outcome, and the 128 came from Git in the
  caller's removed working directory.
- A Declared QA Decline archives with the existing `no-qa` disposition; the
  record schema does not change, and the reason stays in `_tasks.md` at the
  record's `source_revision`.

## Open Questions

None.
