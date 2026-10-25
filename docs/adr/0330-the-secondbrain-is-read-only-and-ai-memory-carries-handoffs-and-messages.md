---
status: accepted
created_at: 2026-10-24T17:30:00Z
updated_at: 2026-10-24T17:30:00Z
deprecated_at: null
superseded_by: null
---

# The Secondbrain is read-only and ai-memory carries handoffs and messages

Supersedes ADR-0095 (2026-10-24); it is retired, and Git keeps its text.

ADR-0095 made the Secondbrain's `inbox/<destination>/` the single door for
every fleeting observation: any session committed a file there, and only the
destination project triaged it. Since 2026-10-19 the maintainer's sessions
also run ai-memory, whose MCP server gives each project an inbox of directed,
claim-once messages (`memory_message_send`, `memory_message_list`,
`memory_message_pop`) and single-use session handoffs (`memory_handoff_*`).
Two doors for one job split the fleet's mail: on 2026-10-24 the maintainer
wrote "Lembrando que solicitei em algum momento para matar o inbox de
secondbrain e usar somente as mensagens do ai-memory". The same day the
upstream skills repository replaced its `handoff` skill with six `ai-memory-*`
skills in every preset (pull request #131, commit
`65d5f365cace27c4b468db1723213723b10c870b`).

## Decision

**The Secondbrain stays the business knowledge base, and it is read-only.**
The Baseline keeps its consultation rules: read `wiki/index.md` first, then
`qmd`, cite every file used, keep secrets out, and ask its owner to ingest a
durable change. A session never writes a Secondbrain file, `inbox/` included.
Substantive external research reaches the Secondbrain as a digest message to
the project that owns it, which ingests it by its own contract.

**ai-memory carries session handoffs and cross-project messages.** The core
module requires and dispatches `ai-memory-handoff` and `ai-memory-messaging`
in place of `handoff`, and its new rule says: an observation another project
owns becomes one self-contained message to that project's inbox; a project
reads its messages oldest first and consumes one only after recording its
content in a Spec, a Backlog Entry or a Finding, or naming it when it only
informs. Each retired inbox clause is replaced by one of these clauses with
the same force, so an adopter's next refresh accounts every clause it held.

**ai-memory is optional and detected at run time.** It is not a Baseline
decision: a repository renders the same guidance with or without it. A
session uses ai-memory only when it exposes the MCP tools; otherwise it tells
the user what is missing and puts the handoff or message text in its
response. Roundfix's own code never talks to ai-memory: the follow-up command
keeps reading messages only from a file the caller hands in, and drops its
`secondbrain-inbox` source. Its two keys, `secondbrain.root` and
`secondbrain.inbox_namespace`, become deprecated keys that are ignored with
one warning, so an existing User Config keeps loading.

**The upstream setups follow the commit that carries the ai-memory skills.**
Every upstream Setup Snapshot moves to `65d5f36` through the asset sync,
the only writer of a Setup Snapshot; the snapshots drop `handoff` with
upstream and gain the six `ai-memory-*` skills; `go` and `rust` also gain
`find-skills` and lose `council` and `the-fool`, which Roundfix had already
retired. The `gws-*` and
`brand-guidelines` skills upstream removed were never in a Roundfix snapshot,
because the Baseline imports no `content` or `devops` preset.

## Considered Options

- Keeping both doors, the Secondbrain inbox for captures and ai-memory for
  handoffs, was the 2026-10-19 arrangement; the maintainer retired it on
  2026-10-24 because the same observation could land in either.
- Making ai-memory a Baseline decision, with clauses that apply only when a
  repository answers yes, was rejected: whether the MCP server is reachable is
  a property of the session, not of the repository, and the fallback (say
  what is missing, put the text in the response) is the same everywhere.
- Having the follow-up command read ai-memory through its HTTP API was
  rejected: the CLI would need a credential and a network call, and the skill
  already hands messages in as a file.
- Removing the two configuration keys outright was rejected: strict
  validation would then refuse every User Config that still names them,
  including the maintainer's.

## Consequences

- Adopters with `secondbrain.enabled` lose the Inbox Entry contract, the
  capture triggers and the `inbox/**` write permission at their next managed
  refresh; each retired clause is recorded `replaced`.
- An adopter without ai-memory reads the same clauses and gets the message
  text in the response instead of a sent message.
- Adopters see the six `ai-memory-*` skills at their next
  refresh and keep any copy of `handoff` they hold until they remove it; the
  release note names the deliberate `baseline skills restore --skill` and
  `baseline skills reconcile` commands.
- Pending Inbox Entries are migrated to ai-memory by the operator and
  `inbox/**` is deleted by the Secondbrain's own project; no Roundfix code
  reads or writes either.
