---
spec: 0390-a-nested-carrier-git-does-not-list-raises-no-conflict
status: active
created: 2026-11-13
surfaces: [backend, cli, docs]
type: fix
---

# A nested carrier Git does not list raises no conflict

Reported by quaestor on 2026-11-13: `roundfix baseline update` (0.143.0) warned `baseline.inventory.nested-carrier-conflict` for the `AGENTS.md` and `CLAUDE.md` of a linked worktree that a local script keeps under `.storage/`, a directory the repository's `.gitignore` lists, and the warnings disappeared when the worktree was removed. Those files are not the repository's: no commit of it holds them and the Baseline never writes them. After this Spec the Baseline inventory keeps a nested `AGENTS.md` or `CLAUDE.md` only when Git lists it as a file of the repository, tracked or untracked and not ignored, so an ignored copy and a carrier inside a linked worktree, a nested repository or a submodule raise no warning, while a tracked nested carrier is inventoried exactly as before.

## Sources

One ai-memory message, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-13-quaestor-a-nested-carrier-conflict-inside-a-git-ignored-worktree.md`](references/2026-11-13-quaestor-a-nested-carrier-conflict-inside-a-git-ignored-worktree.md),
   ai-memory message `01a1d4de-8d39-7f60-b575-5703959df8da`, copied
   verbatim; reported by quaestor on 2026-11-13. Quaestor declared no blocked
   work, asked for no reply and already moved its worktree out of the
   repository.

## Prerequisites

None in `requires`. Authored against `2e7a2a7a` (v0.143.1) and rebased on
`4cb2e08b`, which only deletes a Backlog Entry (#863). Spec 0389,
planned on the default branch and delivered before this one, changes no file
under `internal/baseline/`. The two Specs share only these files, each in its
own part:

- The Roundfix owned skill's version fields in `.agents/skills/roundfix/SKILL.md`
  and `skills/roundfix/SKILL.md`, and `skills/testdata/owned-skill-versions.json`:
  0389 raises them for its `deliver` reference; task_02 of this Spec raises
  the next version with the record command on the tree it runs on and never
  writes a number.
- `docs/references/behavior-surfaces.json`: 0389 re-records it after its
  `deliver` guide edits; task_02 re-records it after its own last guide edit.
- `CONTEXT.md` is not changed by this Spec.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted or changed.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network
  call is added; the new Git command reads the local repository only. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0399 (new) records the decision:
  "The Baseline inventories a nested carrier only when Git lists it". It
  amends ADR-0070, whose audit of every bounded carrier now covers the
  carriers Git lists, and task_02 writes the amendment notice. ADR-0372
  applies unchanged to guides under `docs/agents/`. Not applicable: ADR-0351,
  only the precedent for the same Git command; ADR-0330, which ADR-0372
  cites; and ADR-0074, ADR-0078 and ADR-0100, which cite ADR-0070 (reasons
  in the TechSpec's ADR Alternatives). Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BE batch 2 answer of 2026-11-13 quoted in the authorization record;
  the Governed Paths are the Roundfix skill files task_02 changes and
  `internal/baseline/repository_test.go`, which task_01 changes. No
  Makefile, workflow, Go module, lint or formatter change, and no Baseline
  module, template or profile change. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0390-a-nested-carrier-git-does-not-list-raises-no-conflict/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/baseline.md`, `internal/baseline/repository_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- A nested `AGENTS.md` or `CLAUDE.md` in a path Git ignores raises no
  `baseline.inventory.nested-carrier-conflict` and does not enter the Plan
  Digest.
- A nested carrier inside a linked worktree, a nested repository or a
  submodule raises no warning, wherever that checkout sits.
- A tracked nested carrier, and an untracked one Git does not ignore, is
  inventoried and warns exactly as before.
- The `baseline` guide and the Roundfix skill's `baseline` reference state
  the rule, and the amended Baseline audit decision carries a notice that
  names the new decision.

## Core Features

1. After the walk, an inspection that met a nested `AGENTS.md` or `CLAUDE.md`
   asks Git once for the repository's carriers with `git ls-files -z
   --cached --others --exclude-standard`, reads the output byte for byte, and
   inventories only the nested carriers it lists.
2. A failed listing stops the inspection with an error that names it.
3. The guide, both copies of the Roundfix `baseline` reference and the
   ADR-0070 amendment notice carry the rule.

## User Experience

Quaestor's plan with the worktree present carried two warnings like:

```text
baseline.inventory.nested-carrier-conflict .storage/simulacao-deploy/wt-aee32c4/AGENTS.md: nested carrier remains unchanged: bytes inside setup-context-driven begin/end markers are managed; bytes outside those markers are preserved and require conflict review
```

With this Spec the same command with the worktree present prints what it
prints without it, the same Plan Digest included.

## Non-Goals / Out of Scope

- Root carriers and guides under `docs/agents/`: they are inventoried as
  before, ignored or not.
- Reporting the carriers left out, as a count or a list.
- Any change to the walk's existing skips (`.git`, `node_modules`, `vendor`,
  `skills/` and the others) or to the inventory limits.
- Writing into a nested checkout, or inventorying it as a repository of its
  own: running the Baseline inside that checkout already does.
- Messages to quaestor: the operator sends them.

## Success Metrics

1. Success Metric: in a temporary repository whose `.gitignore` lists
   `.storage/`, a linked worktree under `.storage/`, a plain ignored copy, a
   linked worktree under a path Git does not ignore, a nested repository and
   a submodule add no carrier and no warning, and the snapshot digest equals
   the digest before they existed; the tracked nested carrier, a tracked one
   under an ignored directory and two untracked ones Git does not ignore keep
   their carriers and warnings. At `2e7a2a7a` every foreign carrier warns.
2. Success Metric: an inspection asks Git for the carrier list once, with the
   exact arguments of the TechSpec, only when the walk met a nested carrier,
   and a failed listing returns `list repository instruction carriers: <git
   error>`.
3. Success Metric: `roundfix baseline update --format json` on a clone of
   this repository with an ignored linked worktree reports the Plan Digest and
   the nested-carrier warnings it reports without the worktree; the binary of
   the delivery base reports the worktree's carriers.
4. Success Metric: the `baseline` guide, both copies of the Roundfix
   `baseline` reference and the ADR-0070 notice carry the new text; the
   documentation and skills contracts pass.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The adopted quaestor message, a measurement on a repository this Spec did
  not build: "com a worktree presente, o plano trazia os dois warnings ...
  Depois de `git worktree remove .storage/simulacao-deploy/wt-aee32c4`, o
  mesmo comando deu \"the repository already matches the current Baseline
  catalog\", sem nenhum warning."
- git-ls-files documentation (https://git-scm.com/docs/git-ls-files, read
  through Exa on 2026-11-13): "`--exclude-standard`: Add the standard Git
  exclusions: .git/info/exclude, .gitignore in each directory, and the user’s
  global exclusion file."
- gitignore documentation (https://git-scm.com/docs/gitignore, read through
  Exa on 2026-11-13): "A `gitignore` file specifies intentionally untracked
  files that Git should ignore. Files already tracked by Git are not
  affected".
- The installed `git` (2.56.1 on the authoring host), which this Spec did not
  build: in a temporary repository it lists neither the files of a linked
  worktree, a nested repository nor a submodule, and `git check-ignore`
  exits 128 with `fatal: Pathspec ... is in submodule` for a path inside a
  submodule.

## Unreachable Acceptance

None. Every criterion is proved with temporary repositories, a local clone
of the checkout and the installed `git`; no criterion needs the network, an
open Pull Request or another repository.

## Validation

The decisions were checked on 2026-11-13 against these sources.

- Characterization at `2e7a2a7a` with `bin/roundfix baseline plan --profile
  go-cli-tui` on a temporary repository: `baseline.inventory.nested-carrier-conflict`
  for `.storage/copy/CLAUDE.md`, two carriers of a linked worktree under the
  ignored `.storage/`, a nested repository's `AGENTS.md`, two carriers of a
  linked worktree under a path Git does not ignore and the tracked
  `packages/api/AGENTS.md`; only the last belongs to the repository.
- Characterization on a clone of this repository with `.storage/` in
  `.git/info/exclude` and a linked worktree under it: `baseline update
  --format json` reported four extra warnings and Plan Digest
  `sha256:97afcc52…` against `sha256:4d7703c9…` without the worktree.
- A prototype of task_01 on a scratch copy at `2e7a2a7a`: the same clone gave
  `sha256:4d7703c9…` and the two tracked carriers' warnings with and without
  the worktree; the new tests passed, `./internal/baseline` passed except
  `TestOutputsForCommand`, which needs this repository's real history and
  failed only because the scratch copy has one commit, and every sabotage of
  the TechSpec made its named test fail on its named case.
- Secondbrain, read only: `wiki/index.md`, then `qmd query --all --files
  --min-score 0.3` for an ignored nested worktree and its carriers; it
  returned this repository's mirrored decision that the archive and cancel
  commands refuse ignored Spec files, which already asks Git which files it
  ignores, and the worktree guides, nothing that contradicts asking Git.
- Exa: the two Git documentation pages above.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BE batch 2 brief for Spec 0390 of 2026-11-13, the maintainer's answer below and the adopted source.

The maintainer's answer (AskUserQuestion via `/roundfix --follow-up`,
2026-11-13 ~10:30 BRT), verbatim:

- Question "O Quaestor relatou que o `roundfix baseline update` acusa
  conflito nos AGENTS.md e CLAUDE.md de um worktree guardado numa pasta que o
  Git ignora (.storage/). Esses arquivos não são do repositório. Viro isso na
  Spec 0390 e entrego depois da 0389?" → "Sim, mesmas regras do BE
  (Recommended)" — option text: "O inventário passa a pular caminhos
  ignorados pelo Git e raízes de worktrees ligados. Push, merge e patch logo
  após o merge; até 3 correções; Makefile, CI, go.mod e lint liberados se a
  correção exigir."

Taken from that answer, the source and the characterization, with no
further interview:

- `type: fix`: a false warning becomes no warning; the Release Ceiling is
  `patch`.
- Git decides which nested carriers belong to the repository, with one
  `git ls-files --cached --others --exclude-standard` call: it skips ignored
  paths and stops at linked worktrees, nested repositories and submodules in
  one rule, where `git check-ignore` fails inside a submodule (ADR-0399).
- A directory whose `.git` file no longer points to a repository stays
  inventoried, because Git and `git status` treat it as ordinary content.
- Nothing is reported for a carrier left out, as for the trees the walk
  already skips.
- Adopters receive the change through the binary; no Baseline module,
  template or profile text changes.
- One code Task, one documentation Task, and the QA gate.

## Open Questions

None.
