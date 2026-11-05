---
status: accepted
created_at: 2026-11-05T12:00:00Z
updated_at: 2026-11-05T12:00:00Z
deprecated_at: null
superseded_by: null
---

# runs list prints a table at a terminal, and a short Run id names a Run

`roundfix runs list` prints an aligned table with a header when its standard
output is an interactive terminal, and keeps its stable line for every other
reader. The table's eight-digit short id is accepted by `attach` and
`runs show` as a unique prefix of a Run id's random suffix.

## Context

`runs list` was written for Agents: one line per Run with the full Run id,
the UTC start, the branch and, with `--all` or outside a repository, the Run
Worktree path, and no header. On 2026-11-05 the maintainer ran it on his
MacBook outside a repository and wrote: "A lista de runs está confusa, talves
com informação em excesso". The same morning he chose:

> No terminal interativo, runs list mostra cabeçalho, colunas alinhadas, id
> curto, nome do repositório, Spec, hora em BRT e duração. Branch e worktree só
> com --wide. Em pipe ou com --format plain, continua a linha estável de hoje
> para os agentes.

Outside a Git repository `runs list` lists every repository
([ADR-0357](0357-a-read-only-run-command-needs-no-repository-and-only-warns-on-worktree-location.md)).

## Decision

`--format plain` prints the stable line, byte for byte. `--format table`
prints the table. Without `--format`, the table is printed when `--wide` is
given or standard output is an interactive terminal by the test the Live Run
View uses, and the stable line otherwise. `--wide --format plain` is a usage
error.

The table has a header and the columns `RUN`, `REPO`, `STATE`, `TARGET`,
`AGENT`, `STARTED (<zone>)` and `DURATION`. `REPO` is the repository's base
name and appears only when the listing spans repositories. `STARTED` is the
machine's local time, its zone named once in the header. `--wide` shows the
full Run id and adds `KIND`, `BRANCH` and `WORKTREE`. Only `TARGET` is cut,
with `…`, to fit the terminal width. The stderr notes and `No Runs found.`
are unchanged.

The short id is the first eight hexadecimal digits of a Run id's random
suffix. `attach` and `runs show` match the exact id first; a reference of
eight to sixteen lowercase hexadecimal digits names the one Run whose suffix
starts with it, and two or more matches are a usage error listing the full
ids. Every other command keeps requiring the exact id.

## Considered Options

- One format for everyone: rejected, because Agents parse the stable line and
  a person reads the terminal, and neither shape serves the other.
- The table only behind a flag: rejected, because the person at the terminal
  is the one who would not know the flag.
- The full Run id in the table: rejected, because it is the widest column and
  the maintainer chose the short id; accepting the prefix keeps it usable.
- The prefix in every command that takes a Run id: rejected for now, because
  `stop` and `reconcile` act on the Run and `events` serves automation that
  holds the full id.
- A hard-coded America/Sao_Paulo zone: rejected, because the zone is the
  machine's, and the maintainer's machines already run in it.

## Consequences

A person at a terminal reads a short table. An Agent whose shell is a
pseudo-terminal would get the table too, so the Roundfix skill tells Agents
to pass `--format plain`. A short id can become ambiguous as Runs accumulate;
the refusal names every match.
