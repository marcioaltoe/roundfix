---
spec: 0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal
status: active
created: 2026-11-05
surfaces: [cli, docs]
type: feat
---

# runs list reads as a table, and the Run Browser draws, at a terminal

At an interactive terminal `roundfix runs list` prints the Agents' stable line,
about three hundred characters per Run with the full Run id, the UTC start, the
branch and the Run Worktree path, and bare `roundfix attach` or `roundfix runs`
opens a Run Browser that draws nothing and returns to the prompt without a
word. On 2026-11-05 the maintainer met both on his MacBook. The Run Browser,
the attach cockpit and the owning cockpit of a foreground Run draw nothing
because the command hands the terminal program a wrapper of its output
instead of the terminal itself, so the program cannot measure the screen.
After this Spec `runs list` prints a short aligned table to a person at a
terminal and keeps the stable line for every Agent and pipe, the table's short
Run id works in `attach` and `runs show`, the Run Browser and both cockpits
draw on the terminal, and a Run Browser closed without a selection always says
so on one line.

## Sources

One source, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-05-macbook-runs-list-is-confusing-and-attach-shows-nothing.md`](references/2026-11-05-macbook-runs-list-is-confusing-and-attach-shows-nothing.md),
   reported by the maintainer on 2026-11-05 at about 06:50 BRT, copied
   verbatim: the MacBook `runs list` output, two bare `attach` calls that
   printed nothing, and his remark.

No open Backlog Entry or unresolved Finding exists on `main`, and this
project's ai-memory inbox held no message, so no other source shares this
context.

## Prerequisites

None in `requires`. Authored against `38cc5854` (v0.132.0), where no Spec is
active. Spec 0352, authored in parallel and delivered first, owns the
configuration loader and may change the same shared files: the Roundfix
skill's `SKILL.md` pair and version record, the `runs` reference of the
Roundfix skill and the `runs` user guide (its paragraph on a detected Git
root), and the Behavior Surface Record. This Spec keeps its edits to those
files narrow and rebases on 0352 before its plan merges; its Verification does
not depend on 0352.

## Project Constraints

- Identifier strategy: applicable — no identifier is minted; the short id is
  a display of the existing Run id (its first eight random hexadecimal digits)
  and is resolved back to that id, never stored. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential changes; the commands read only the Run Database and the
  terminal. Source: `docs/agents/go.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0357 makes `runs list` list every
  repository outside a Git repository, ADR-0357: "Outside a Git repository `runs list` behaves as `runs list --all`, with the
  repository column and exit `0`", and makes `attach` and `runs show`
  Read-Only Run Commands; both hold, and the table shows the repository
  column in the same cases. ADR-0362 (the Git-root fix, delivered first) narrows
  which `.git` entry marks a Git root and amends ADR-0357; it changes only
  how the repository is found, never the listing, the table or the Run
  Browser, so this Spec carries no obligation from it. ADR-0363, written by this Spec, records the dual
  text surface and the short id. Source: `docs/agents/domain.md`,
  `docs/agents/cli.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AZ answers of 2026-11-05 quoted in the authorization record, which
  keep the cycle AX rules; the governed files are the Roundfix skill files the
  Tasks change. No Makefile, workflow, Go module, lint, formatter or Baseline
  module change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0353-runs-list-reads-as-a-table-and-the-run-browser-draws-at-a-terminal/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/runs.md`, `skills/roundfix/SKILL.md`.

## Goals

- At an interactive terminal `runs list` prints a header and aligned columns
  a person can read in one screen width; piped, or with `--format plain`, it
  prints the stable line byte for byte as today.
- The short Run id the table shows is accepted by `attach` and `runs show`.
- Bare `attach` and bare `runs` draw the Run Browser on the terminal, the
  attach cockpit and the owning cockpit draw too, and a Run Browser that ends
  without a selection prints one line naming how it ended.
- The Roundfix skill, the user guides and ADR-0363 tell people and Agents
  which format they get and how to choose it.

## User Stories

1. As the maintainer at a terminal, I want `runs list` to show a short table
   with the repository, state, Spec, Agent, local start time and duration, so
   that I can tell at a glance what is running.
2. As an Agent or a script, I want `runs list` piped or with `--format plain`
   to keep printing today's stable line, so that nothing that parses it
   breaks.
3. As the maintainer, I want to pass the short id from the table to `attach`
   or `runs show`, so that I do not have to copy a full Run id.
4. As the maintainer, I want bare `attach` and bare `runs` to draw the Run
   Browser, and to say so when it closes without a selection, so that the
   command never leaves me at the prompt with nothing.
5. As the maintainer, I want `--wide` to add the kind, branch and worktree and
   show the full id, so that I can still see everything when I need it.

## Core Features

1. Format selection. `runs list` gains `--format <table|plain>` and `--wide`.
   `--format plain` prints the stable line; `--format table` prints the
   table; without `--format` the table is printed when `--wide` is given or
   standard output is an interactive terminal, and the stable line
   otherwise. `--wide --format plain`, and any other `--format` value, are
   usage errors that exit `2`.
2. The table. A header and two-space aligned columns `RUN` (short id),
   `REPO` (repository name, only when the listing spans repositories),
   `STATE`, `TARGET` (the Spec slug, or `pr:<number>`), `AGENT`,
   `STARTED (<zone>)` (local time, `HH:MM` today and `Mon DD HH:MM` on another
   day) and `DURATION`. A row whose owner identity is unproven adds a `NOTE`
   column. Only `TARGET` is cut, with `…`, to fit the terminal width.
   `No Runs found.` and every stderr note are unchanged.
3. The wide table. `--wide` shows the full Run id and adds `KIND`, `BRANCH`
   and `WORKTREE`, and cuts nothing.
4. The short id. `attach` and `runs show` accept the exact Run id as today,
   and a reference of eight to sixteen lowercase hexadecimal digits that
   names the one Run whose random suffix starts with it. Two or more matches
   exit `2` listing the full ids; no match keeps today's unknown-Run message.
5. The Run Browser and the cockpits draw. The Run Browser, the attach cockpit
   and the owning cockpit receive the terminal itself as their output, and a
   terminal program given anything else refuses with an error that names it
   instead of drawing nothing.
6. Never silent. When the Run Browser ends without a selection, bare `attach`
   and bare `runs` print one stderr line naming the key or condition that
   closed it and exit `0`. The Run Browser's `REPO` column names the
   repository, not the Run Worktree.
7. The rules in writing. The Roundfix skill's `runs` reference, the `runs` and
   `attach` user guides, the command help and ADR-0363 state the formats, the
   flags, the short id and the closing line.

## User Experience

At a terminal, from a directory that is not a repository:

```text
RUN       REPO      STATE               TARGET                                AGENT  STARTED (-03)  DURATION
93cee2a3  fiscus    ResolvingWithAgent  0063-um-gtin-por-mestre-e-a-pendenc…  codex  06:44          3m
cf0dfe58  quaestor  ResolvingWithAgent  0063-sequencia-pela-emissao-propria…  codex  06:32          15m
(9 terminal Run Worktrees retained; run 'roundfix reconcile' to inspect)
```

`roundfix attach 93cee2a3` opens that Run's cockpit. Bare `roundfix attach`
draws the Run Browser; pressing `q` returns to the prompt with
`roundfix: Run Browser closed (q); no Run selected.` on stderr.

## Non-Goals / Out of Scope

- Changing the stable line, its columns, its notes, its exit codes or the
  `--all`, `--state` and `--limit` flags.
- A JSON format for `runs list`.
- Accepting the short id in `events`, `stop`, `reconcile` or any command
  other than `attach` and `runs show`.
- Changing the Run Browser's keys, layout, short id or column drop order.
- Hard-coding a time zone, or a configuration key for it.
- A new domain term; the short id and the table are described, not named.
- The configuration loader and Git root detection (Spec 0352), the Makefile,
  the CI workflows, `go.mod`, `go.sum`, lint or formatter configuration, any
  Baseline module and any `### QA settlement` section.

## Success Metrics

1. Success Metric: over a seeded Run Database, `runs list` with standard
   output not a terminal, and `runs list --format plain`, print exactly the
   bytes today's binary prints for the same Runs and flags, including `--all`
   and the owner-identity marker.
2. Success Metric: `runs list --format table` over Runs of two repositories
   prints the header `RUN  REPO  STATE  TARGET  AGENT  STARTED (<zone>)  DURATION`
   aligned, the eight-digit short ids, the repository names, local times and
   bare durations; inside one repository it has no `REPO` column; at an
   80-column budget no line is wider than 80 and only `TARGET` is cut.
3. Success Metric: `runs list --wide` prints the full ids and the `KIND`,
   `BRANCH` and `WORKTREE` columns, and `--wide --format plain` and
   `--format yaml` exit `2` with their messages.
4. Success Metric: `attach <short id>` and `runs show <short id>` act on the
   one matching Run; an ambiguous short id exits `2` listing every match; an
   exact id still wins over a prefix.
5. Success Metric: the Run Browser, the attach cockpit and the owning cockpit
   given the command's wrapped output reach the terminal program with the
   underlying writer, and a terminal program given a writer that is not a
   file refuses naming its type.
6. Success Metric: a cancelled Run Browser prints
   `roundfix: Run Browser closed (<reason>); no Run selected.` once on stderr
   for bare `attach` and bare `runs`, with no Run Database too, and stdout
   stays empty.
7. Success Metric: the `runs` reference of the Roundfix skill and its mirror,
   the `runs` and `attach` user guides, the help text and ADR-0363 state the
   formats, `--wide`, the short id and the closing line.

## Acceptance evidence

The outside-evidence rows rest on sources this Spec did not write.

- The maintainer's report, quoted in the adopted reference and written before
  this Spec existed: on roundfix 0.132.0 `runs list` printed two lines with
  the full id, the UTC start, the branch and the worktree path, and bare
  `roundfix attach` printed nothing twice; "A lista de runs está confusa,
  talves com informação em excesso e o roundfix attach não mostra nada."
- The Command Line Interface Guidelines (https://clig.dev/, read through Exa
  on 2026-11-05): "The most simple and straightforward heuristic for whether a
  particular output stream (`stdout` or `stderr`) is being read by a human is
  whether or not it’s a TTY", and "If human-readable output breaks
  machine-readable output, use `--plain` to display output in plain, tabular
  text format for integration with tools like `grep` or `awk`."
- Git's revision syntax (https://git-scm.com/docs/gitrevisions, read through
  Exa on 2026-11-05) names an object by "The full SHA-1 object name (40-byte
  hexadecimal string), or a leading substring that is unique within the
  repository", the published practice the short id follows.

## Unreachable Acceptance

- criterion: the Run Browser draws and closes with its line at the
  maintainer's MacBook terminal
  reason: a Run has no interactive terminal and cannot reach the MacBook; the
  QA gate proves the writer the terminal program receives and the closing
  line through the command, not a real screen
  satisfied-by: the model and command tests this Spec adds, and the
  maintainer running bare `roundfix attach` on the MacBook after v0.133.0

## Validation

The decisions were checked on 2026-11-05 against these sources.

- Reading at `38cc5854`: every command's standard output and standard error
  are wrapped with the command environment before dispatch, and that wrapper
  reaches the Run Browser, the attach cockpit and the owning cockpit as their
  terminal output. A bare `attach` without an interactive standard input or
  output already exits `2` with the non-interactive message. A cancelled Run
  Browser returns exit `0` with no output. The Run Browser's repository
  column is the base name of the Run's checkout, which is the Run Worktree
  for a queued Run. No Run id prefix is accepted anywhere.
- A pseudo-terminal probe on the Mac mini with roundfix built at
  `38cc5854`, sized 120 by 40: bare `runs` entered the alternate screen,
  cleared it and drew nothing for 2.5 seconds, a resize drew nothing, and `q`
  or `Esc` returned to the prompt with no output. The same Run Browser started
  by a small program handed the terminal file drew its header, the empty
  state and the footer; handed a wrapper of that file, it drew nothing. An
  escape byte that reached the program alone, as a terminal reply split
  across reads would, closed it as `Esc`. A prototype on a scratch copy
  that passed the unwrapped writer drew the Run Browser in the same
  pseudo-terminal and printed a closing line after `q`.
- Secondbrain, read only: `wiki/index.md`, then
  `qmd query "roundfix runs list attach run browser terminal table" --all --files --min-score 0.3`,
  which did not finish in fifteen minutes, so
  `qmd search "roundfix runs list attach run browser" --all --files` was
  used: it returned the Roundfix and skills project mirrors (the archived
  Run Browser Spec 0020, the `runs`, `attach` and `usage` guides, the
  Roundfix skill references, `CONTEXT.md`, the CHANGELOG) and
  `wiki/concepts/execucao-remota-persistente-com-tmux-no-brain-macmini.md`,
  which records that the maintainer reaches the Mac mini from the MacBook
  through cmux, mosh and tmux. None treats the Run Browser drawing nothing;
  a multiplexer between the terminal and Roundfix is one way a terminal
  reply can arrive split and be read as `Esc`.
- Exa: the two pages quoted above.

## Glossary

None.

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle AZ brief for Spec 0353 of 2026-11-05 and the maintainer's answers below.

The maintainer's answers (AskUserQuestion, 2026-11-05 ~07:10 BRT), verbatim:

- Listing: "Tabela para pessoas (Recommended)" — option text: "No terminal
  interativo, runs list mostra cabeçalho, colunas alinhadas, id curto, nome
  do repositório, Spec, hora em BRT e duração. Branch e worktree só com
  --wide. Em pipe ou com --format plain, continua a linha estável de hoje
  para os agentes."
- Rules: "Mesmas do AX" — "Sem supervisão; merge pela fila; corrigir e
  continuar nas paradas; Makefile, CI, go.mod e lint liberados; mudanças de
  Baseline e skills que chegam aos adotantes liberadas."
- Release: minor v0.133.0, cut by the operator right after the merge, on
  2026-11-05. The maintainer, 2026-11-05: "Quero que faça os releases hoje
  assim que estiverem prontos. E se existe no baseline regra para um release
  por dia por algum motivo deve ser removida." No such rule exists in the
  repository or the Baseline (it was the operator's habit), so this Spec
  changes nothing for it.

Taken from those answers, ADR-0357 and the reading above, with no further
interview:

- The command stays in English like every other surface: the headers are
  `RUN`, `REPO`, `STATE`, `TARGET`, `AGENT`, `STARTED` and `DURATION`; the
  preview's Portuguese labels showed the layout. See ADR-0363.
- "Hora em BRT" is the machine's local zone, named once in the header; no
  zone is hard-coded. See ADR-0363.
- The short id is accepted as a unique prefix by `attach` and `runs show`,
  so the id shown is usable. See ADR-0363.
- `REPO` appears only when the listing spans repositories, as the stable
  line's repository column does.
- The silent end is fixed at its cause in Roundfix (the terminal programs get
  the terminal itself) and generically (a closing line on every end without a
  selection); a non-terminal standard input already makes bare `attach`
  non-interactive, and this Spec keeps it.
- Three implementation Tasks (the Run Browser and cockpits, the table, the
  short id with ADR-0363) and a QA gate.

## Open Questions

None. The MacBook diagnostics (exit code, `TERM`, `TERM_PROGRAM`) had not
arrived when the plan was written; if they show a second cause, such as a
terminal reply read as `Esc`, the closing line names it and a follow-up Spec
treats it.
