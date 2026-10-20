---
spec: 0288-the-archive-rewrites-relative-citations-from-any-worktree
status: active
created: 2026-10-20
surfaces: [backend, cli, docs]
---

# The archive rewrites relative citations from any worktree

Minimal PRD for a bug fix, minted by `write-techspec`. `roundfix archive` and
`roundfix cancel` build repository-relative paths by comparing two spellings
of the same checkout: the root Git reports, which resolves symbolic links,
and the Spec Root built from the work directory, which keeps them. When the
command runs from a path that crosses a symbolic link, such as a Delivery
Queue item worktree under a linked `~/.roundfix/worktrees`, the two spellings
differ. The archive then refuses a Spec whose PRD links outside its folder,
silently leaves Live Documents citing the removed folder, and the cancel
refuses outright. After this Spec both commands compute every
repository-relative path from one repository root, whichever spelling of the
checkout they run from, and the TechSpec holds the design.

## Sources

One Secondbrain Inbox Entry, copied verbatim by this Spec's planning change
into `references/` (see `references/_index.md`):

1. [`references/2026-10-20-o-archive-com-qa-override-nao-reescreve-link-relativo-no-item-worktree.md`](references/2026-10-20-o-archive-com-qa-override-nao-reescreve-link-relativo-no-item-worktree.md),
   origin vortex, type hint `fix`. Run at the root of a Delivery Queue item
   worktree, as the next action of `deliver status` asks for a
   `qa-environment-partial` item, `roundfix archive <slug> --qa-override
   --approval … --reason …` refused with `rebase archive outcome: rebase
   citation "../../backlog/<file>.md": Rel: can't make docs/backlog/<file>.md
   relative to ../../../../../../../Users/marcio/.roundfix/worktrees/…/docs/history/specs`,
   and the operator had to edit the Spec's text to archive it.

No open Backlog Entry or unresolved Finding shares this context. The only open
Backlog Entry, `docs/backlog/2026-10-20-gc-sanitize-preserves-every-artifact-root-a-run-recorded-empty.md`,
is about which Artifact Roots `roundfix gc sanitize` reaches; Spec 0289 owns
it. `docs/findings/` holds no entry.

## Prerequisites

None. This Spec is authored against `7a678ddb` (v0.97.0). Spec 0287, still in
the Delivery Queue, changes `internal/spec/archive_record.go` and reads
`internal/spec/cancel.go`; this Spec changes neither the summary reader nor
anything 0287 declares, and its plan merges after 0287, so it is rebased and
checked again before its Run. Spec 0289 is authored beside it and delivers
after it; they share no file.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; the Archive Record keeps `roundfix/archive-record/v1` and its
  fields. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; both commands read and write only the
  checkout they run in, and every test uses temporary repositories and
  homes. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the fix keeps every recorded archive
  and cancel decision and applies the resolved-path rule two earlier
  decisions set. ADR-0247: "An archive leaves an Archive Record, and the
  Spec folder stays in Git". ADR-0260: "An archive relinks live citations,
  and a Baseline update lists the citations it breaks". ADR-0278: "A
  retirement deletes, and the History Root holds only Archive Records".
  ADR-0262: "A cancelled Spec leaves an Archive Record and returns its
  sources". ADR-0154: "A QA archive override records user authority, not a
  pass". ADR-0269: "Archive and cancel refuse ignored Spec files and finish a
  failed removal forward". ADR-0268: "A recorded worktree path is resolved
  within its Worktree Root". ADR-0272: "A review record survives a fix, and a
  correction inside the archived Spec's declared paths returns to review".
  ADR-0299: "The history conversion rewrites every citation it lists, and the
  check reports a removed history file cited without a Git Locator".
  ADR-0259: "A test runs in parallel unless it names why it cannot".
  ADR-0271: "A test bounds by real time only what it asserts". ADR-0285:
  "Every test file compiles for Windows, and a repository contract proves
  it". ADR-0104: "A Spec accepts on evidence it did not author". ADR-0240:
  "One QA partial policy, and rows a Run sandbox cannot reach". ADR-0276: "A
  row the authorization keeps offline qualifies, and a byte invariant
  compares with the Delivery Base". ADR-0265: "A Spec proves what it claims,
  and every active Spec is held to it". ADR-0183: "An attributed claim
  carries a receipt the check proves". ADR-0184: "A TechSpec states a
  command surface as a transcript". ADR-0280: "A guide Task records its own
  surface fingerprint, a PRD records who decided, and the glossary is held to
  its own rule". ADR-0256: "A release waits for the skills that describe a
  changed surface". Set no obligation here, because this Spec changes none of
  their rules: ADR-0092 (triage intent lives in a typed backlog); ADR-0120
  (retired documentation lives under one history root); ADR-0167 (the pre-PR
  Pull Request row never decides a qualifying partial); ADR-0177 (the
  relocation citation scan reads through a repository root); ADR-0229 (an
  operator archive resumes any park from the Run start); ADR-0230 (an
  archived Spec keeps its relative links); ADR-0248 (existing history is
  sanitized in batches after a history-full tag); ADR-0250 (a module version
  is chosen when it is recorded, and the coverage record lists every
  platform); ADR-0251 (a Legacy Archive Folder is read leniently); ADR-0254
  (baseline update sanitizes pending history in the change it plans);
  ADR-0258 (a Spec is authored against the QA rerun classes); ADR-0261 (an
  Archive Record names its pull request); ADR-0277 (a Legacy Archive Folder
  conversion keeps the summary, links and dates its folder held); ADR-0279
  (status reports a Queue Owner proven gone); ADR-0281 (a test replaces an
  executable it rebuilds); ADR-0282 (the history cleanup names a superseded
  ADR in its successor); ADR-0286 (a fixed round-2 finding gets one
  Correction Review); ADR-0295 (reachability reads the fetched default
  branch); ADR-0298 (a timed-out profile proof gets a third attempt);
  ADR-0165 (a blocking review after archive parks publication for a
  corrective Spec); ADR-0197 (a pre-PR reviewer lineage spans at most two
  rounds); ADR-0237 (a Delivery Retry records a merge made outside the
  queue); ADR-0290 (merged items give back their Item Binaries); ADR-0169
  (the pre-PR review diffs the candidate from its merge base); ADR-0233 (a
  skill version raise is regenerated at merge); ADR-0196 (a pre-PR review
  finding parks only after validation). Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AL approval of 2026-10-20 through AskUserQuestion, option "Mesmas regras do
  AK (Recomendado)", which carries the cycle AJ grant "Mexer em Governed
  Paths declarados" ("'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica)."). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0288-the-archive-rewrites-relative-citations-from-any-worktree/_authorization.md`;
  bounded files: `internal/spec/archive.go`. No Makefile, workflow, Go
  module, lint, formatter, skill, Baseline module or Project Config file
  change.

## Goals

- `roundfix archive`, with or without `--qa-override`, archives a Spec whose
  files link outside its folder from any spelling of the checkout: the main
  checkout or a linked worktree, its root or a subdirectory, through a
  symbolic link or not. The Archive Record, the relinked Live Documents and
  the output are the same bytes whichever spelling it ran from.
- `roundfix cancel` does the same for a cancelled Spec.
- A path either command names in a refusal is repository-relative from every
  spelling.
- A path that really lies outside the repository is refused with a message
  that names it and the repository root, never rewritten as `../…` or
  silently skipped.

## Core Features

1. One resolver for repository-relative paths. A path is relativized to the
   repository root lexically; when that leaves the root, both are compared
   in their resolved spellings; when that also leaves the root, the path is
   refused. Every path that resolves lexically today keeps its bytes.
2. The Archive Command and the Cancel Command use that resolver for every
   repository-relative path they build: the Spec folder, the Spec Root, the
   Archive Root, the Archive Record path, promoted and returned sources, and
   the paths their refusals name.
3. The Delivery Queue's archive step runs the same Archive Command in the
   item worktree, so it is covered without a change of its own.
4. Documentation. The `archive` and `cancel` user guides say that each
   command reads every path relative to the checkout it runs in, whichever
   spelling of it the work directory uses.

## Non-Goals / Out of Scope

- Changing how Roundfix records worktree paths, or the Worktree Root rule of
  ADR-0268.
- Changing the citation rewrite itself: which destinations it rewrites, Git
  Locator text, Legacy Link Rebase and Archive Relink rules stay as they are.
- `roundfix history sanitize` and the history conversion of `roundfix
  baseline update`: each already builds both sides of every comparison from
  one root, so they need no change.
- Rewriting any Archive Record already written, or repairing Live Documents
  an earlier symlinked archive left citing a removed folder.
- Changing the Roundfix skill: its archive reference says nothing about where
  the commands run, so it keeps its text and the Coverage Review of each
  surface is recorded instead (see Skills).

## Success Metrics

1. Success Metric: in a temporary repository with a linked worktree added
   outside the main checkout and reached through a symbolic link, `roundfix
   archive <slug> --qa-override --approval <a> --reason <r>` of a Spec whose
   PRD summary links `../../backlog/<file>.md`, while that Backlog Entry
   links the Spec's `_prd.md`, exits 0 from the worktree root and from a
   subdirectory spelled through the link, and from the main checkout spelled
   through a link. Each run prints the same standard output and standard
   error, and writes the same Archive Record and relinked Backlog Entry
   bytes, as the run from the resolved spelling.
2. Success Metric: `roundfix cancel <slug> --reason <r>` from the worktree
   root spelled through the link exits 0 and writes the same record and
   relinked bytes as from the resolved spelling.
3. Success Metric: from the worktree spelled through the link, the
   uncommitted-change refusal names `docs/specs/<slug>`, and the record-only
   `--pull-request` refusal for a record already on the default branch names
   `docs/history/specs/<slug>.md`.
4. Success Metric: `spec.Archive` and `spec.Cancel` given a repository root
   and a Spec Root spelled differently write the same bytes as given one
   spelling, and a path outside the repository is refused with a message
   naming both paths.
5. Success Metric: on the maintainer's machine after release, the next
   archive run from a real Delivery Queue item worktree under the linked
   `~/.roundfix/worktrees` archives a Spec whose PRD links outside its
   folder without editing the Spec.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write, read on
2026-10-20.

- The adopted Inbox Entry, written by a vortex session on 2026-10-20 and
  copied verbatim into `references/`: the refusal text, the item worktree
  path under `~/.roundfix/worktrees`, and that the same kind of link archived
  without trouble when the archive ran inside the Run.
- The Go standard library documentation of `os.Getwd`, readable offline with
  `go doc os.Getwd`: "If the current directory can be reached via multiple
  paths (due to symbolic links), Getwd may return any one of them." and "On
  Unix platforms, if the environment variable PWD provides an absolute name,
  and it is a name of the current directory, it is returned."
- The Go standard library documentation of `filepath.Rel`, readable offline
  with `go doc path/filepath.Rel`: "Rel returns a relative path that is
  lexically equivalent to targPath when joined to basePath".
- Git's `builtin/rev-parse.c`
  (https://github.com/git/git/blob/master/builtin/rev-parse.c): for
  `--show-toplevel` it prints `repo_get_work_tree(the_repository)`; a Visual
  Studio Code maintainer confirmed on 2022-12-13
  (https://github.com/microsoft/vscode/issues/116064): "At the moment `git
  rev-parse` resolves symbolic links".

## Unreachable Acceptance

- criterion: Success Metric 5, on the maintainer's machine after release, an
  archive run from a real Delivery Queue item worktree under the linked
  `~/.roundfix/worktrees` archives a Spec whose PRD links outside its folder
  without editing the Spec.
  reason: no Verification or QA row may read or write the real
  `~/.roundfix`, and only a real Delivery Queue item on that machine has the
  linked Worktree Root.
  satisfied-by: the operator's next `qa-environment-partial` recovery, or any
  archive run from an item worktree, after this Spec's release.

## Validation

The decisions were checked on 2026-10-20 against these sources.

- Characterization at `7a678ddb` with the installed `roundfix 0.97.0`, in
  temporary repositories under a temporary `HOME` and `GIT_CONFIG_GLOBAL`,
  with a linked worktree outside the main checkout and a symbolic link to the
  directory holding both. From the link spelling of the worktree root, of its
  `docs/` subdirectory and of the main checkout, the archive refused with
  `rebase archive outcome: rebase citation
  "../../backlog/2026-10-20-follow-up.md": Rel: can't make
  docs/backlog/2026-10-20-follow-up.md relative to
  ../../link/item/docs/history/specs`, the vortex error; from the resolved
  spelling it archived and relinked one file. Without a link in the PRD, the
  archive from the link spelling succeeded but left the Backlog Entry citing
  `../specs/0001-demo/_prd.md`, with no `relinked` count and no warning. The
  cancel from the link spelling refused with `plan cancellation citation
  advice: list archive citation targets: lstat <resolved root><link
  root>/docs/specs/0001-demo: no such file or directory`. The TechSpec
  records each.
- Secondbrain, read only. Read `wiki/index.md`, then
  `wiki/concepts/caminhos-de-worktree-e-links-simbolicos.md` and
  `wiki/sources/digest-roundfix-worktree-git-sob-symlink-e-teste-de-lock-sem-prazo-2026-10-12.md`,
  which record that Git writes and lists the resolved path and that a tool
  comparing a stored spelling with Git's must resolve both sides; ran `qmd
  query "archive relative link rewrite worktree symbolic link filepath.Rel
  repository root" --all --files --min-score 0.3`, which returned the mirror
  of ADR-0268, that digest, ADR-0177, ADR-0230 and ADR-0260, and `qmd query
  "git rev-parse show-toplevel resolved path symlink cwd PWD" --all --files
  --min-score 0.3`, which returned two unrelated Git talks. They confirmed the
  class and changed nothing in the design. No Inbox Entry was written,
  because the briefing forbids writing there.
- Exa: Git's `builtin/rev-parse.c`, the Visual Studio Code issue above, a
  Stack Overflow question (https://stackoverflow.com/questions/39746533)
  showing `--show-toplevel` printing the physical path of a linked project,
  and an editor issue (https://github.com/adriandeleon/Editora/issues/434)
  describing the same mismatch ("root and file share no prefix") with the fix
  direction of resolving both sides. Challenging: resolving both sides
  always would change the bytes of a Spec Root reached through a link inside
  the repository, so the resolver tries the lexical spelling first.

## Skills

- unchanged: command archive — task_02 changes how `internal/cli/archive.go` builds repository-relative paths; no flag, help text, output line or exit code changes, so the Roundfix skill's archive reference keeps its text.
- unchanged: command cancel — task_02 changes how `internal/cli/cancel.go` builds repository-relative paths; no flag, help text, output line or exit code changes, so the same reference keeps its text.
- unchanged: guide docs/user-guide/commands/archive.md — task_02 adds one sentence saying the archive reads every path relative to the checkout it runs in; the skill's archive reference names no work directory, so it keeps its text.
- unchanged: guide docs/user-guide/commands/cancel.md — task_02 adds the same sentence for the cancel; the same reference keeps its text.

## Glossary

None.

## Decisions

- 2026-10-20 (about 13:10 BRT): the cycle AL sequence "Aprovar como está
  (Recomendado)", P (0288) then Q (0289), one minor each; the authorization
  "Mesmas regras do AK (Recomendado)".
- 2026-10-19 (about 16:20 BRT): the cycle AJ grant "Mexer em Governed Paths
  declarados" (option text "'Concedo' para os caminhos que cada Spec declarar
  (módulos Baseline, registros, skills, testes governados; fora: Makefile,
  CI, go.mod, lint, que pedem autorização específica).").

Autonomous route: no person present; decisions taken from the operator's cycle AL briefing for Spec 0288 of 2026-10-20 and the characterization recorded in the TechSpec.

Taken from those sources, with no further interview:

- Fix the class in the two commands that build archive paths, with one
  resolver, rather than in the citation rewrite: the rewrite's contract
  already says its paths are repository-relative.
- Lexical first, resolved second, refused third, so every case that works
  today keeps its bytes.
- No new ADR: the fix applies ADR-0268's and ADR-0272's rule that a stored
  spelling and Git's spelling are compared resolved.

## Open Questions

None.
