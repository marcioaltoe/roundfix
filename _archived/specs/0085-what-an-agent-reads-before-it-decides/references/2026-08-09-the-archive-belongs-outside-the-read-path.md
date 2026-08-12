---
type: refactor # feat | fix | perf | refactor
status: promoted # open | promoted | declined
created: 2026-08-09
spec: 0085-what-an-agent-reads-before-it-decides # Spec slug when status: promoted
reason: null # required when status: declined
---

# The archive belongs outside the read path, not marked inside it

## Opportunity

Inactive knowledge lives inside the directories agents read. `docs/adr/` holds
105 decision records, `docs/specs/` holds every archived Spec under its own
`_archived/`, and `docs/findings/` does the same. An Agent loading decision
context reads the whole of `docs/adr/`, superseded records included.

The lifecycle markers the repository already uses — `status: superseded`,
`superseded_by`, `deprecated_at` — answer a different question. They tell a
reader which record still governs; they do not keep a retired record out of the
context an Agent pays for. Marking ADR-0106 superseded on 2026-08-09 left its
bytes exactly where they were.

The current shape also scatters the archive: `docs/specs/_archived/` and
`docs/findings/_archived/` are two trees today, and any future one adds a third.
Every consumer that wants to exclude history — `.coderabbit.yaml` path filters,
an Agent's context bundle, a grep — has to know each location separately.

## Value

Two costs, both measured on this repository.

**Context.** 105 ADRs sit in one directory with no structural separation between
the 31 that carry an accepted status, the 20 that carry only a legacy
`Status: Accepted` body line, and the 53 that carry no status at all. Nothing in
the layout tells a reader or an Agent which subset still governs.

**Review.** `.coderabbit.yaml` already excludes generated and tooling paths
through `path_filters`. Excluding history needs one filter per archive tree
today; a single top-level tree needs one line.

The proposed shape — `_archived/specs/`, `_archived/findings/`, `_archived/adr/`,
`_archived/backlog/` under one root — keeps active material in the read path and
history out of it, with one path for every consumer to know.

## What the Secondbrain says about it

`wiki/concepts/arquitetura-de-instrucoes-e-progressive-disclosure.md` treats
agent instructions as a context architecture rather than a store of every known
rule, and it makes two claims that bear directly on this.

The first raises the stakes past token cost: "documentação stale é mais perigosa
para agentes do que para humanos porque o agente pode tratá-la como instrução
atual." A superseded ADR left in `docs/adr/` is not merely paid for on every
task — it can be read as governing. Lifecycle markers do not fix that, because
they rely on the reader noticing them.

The second answers the obvious objection, that moving history hides it: "o
objetivo não é esconder conhecimento. É tornar o próximo documento descobrível
quando necessário, sem fazê-lo pagar tokens em todas as tarefas." That is the
test this work must pass — the archive stays reachable, and the forward pointer
question above is how it stays reachable rather than an optional nicety.

## Shape

Non-binding, but the blast radius is known: `_archived` appears 50 times in Go
code across `internal/speccheck`, `internal/spec`, `internal/specaudit`,
`internal/worktree`, and `internal/cli`, including hardcoded
`docs/specs/_archived` and `docs/findings/_archived` literals. `archive-spec`
and the Archive Command both write to the per-tree location, and
`docs/agents/docs-layout.md` is the contract that declares it.

Worth settling in the same work: whether an ADR moves to `_archived/adr/` when
it is superseded or deprecated, and if so whether the superseding record must
carry a forward pointer so a reader who starts from the active set can still
reach the history. A move without that pointer trades context cost for a broken
trail, which is the same failure in the other direction.
