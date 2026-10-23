---
status: accepted
created_at: 2026-10-22T22:00:00Z
updated_at: 2026-10-22T22:00:00Z
deprecated_at: null
superseded_by: null
---

# The follow-up order is code over four Jev yes/no answers

On 2026-10-22 the maintainer asked for one read-only command that lists a
repository's pending work and orders it, and chose the option "Comando CLI +
skill (Recomendado)": a command gathers the backlog, findings, Specs, Runs,
the Delivery Queue, Pull Requests, git and the Secondbrain inbox, "pergunta
ao Jev pelo internal/judge (model id no log) e ordena em código", while the
skill gathers the ai-memory messages, because only the agent sees that MCP.
The maintainer fixed the four questions (is it a fix, does it block another
project, is it small, does it affect production) and the four classes of the
order. ADR-0201 and ADR-0209 kept Inbox Entries and other Secondbrain content
out of every judge request, and no threshold was ever measured for these
questions.

## Decision

**The order is computed in code, never asked.** For each Follow-up Item the
command asks four Noul questions together in one request, each answered as a
probability that the answer is yes. A probability of 0.7 or more is yes, 0.3
or less is no, and anything between is unsure, which counts as not yes. The
code then places the item in the first class that holds: `unblocking` (it
blocks another project, and it is a fix or it is small), `health` (it is a
fix), `feature` (it affects production) and `other`. Inside a class, small
items come first, then the oldest, then the reference in byte order.

**Without an answer, the item's own type orders it.** When the judge is
unavailable for the whole run (no key, inside a Run, the monthly ceiling, an
unreadable Judge Log, a service failure) or skips one item, that item's class
comes from its declared type: `fix`, `perf` and a Finding are `health`,
`feat` is `feature`, and everything else is `other`. Nothing is `unblocking`
without an answer. The report names the reason, and the command still exits
`0`, as ADR-0200 requires of every advisory judgment.

**The judge reads Follow-up Items, and nothing else changes.** A request's
state is one item's kind, type, origin, title and an excerpt of its text: a
pending Inbox Entry of this repository's inbox namespace, a message the skill
handed in, an open Backlog Entry, an unresolved Finding, a note in
`docs/_inbox/` or an active Spec's PRD. An item whose excerpt matches a
secret pattern is not sent and is ordered by its type. Items are sent
whatever their language, because the inbox is written mostly in Portuguese
and the order is advisory. The recipients, the keys, the Judge Log, the
monthly ceiling and the Run Marker stay as ADR-0201, ADR-0231, ADR-0239 and
ADR-0283 decide, and the Judge Log records the versioned model each answer
came from.

**The command never reaches ai-memory or a remote Secondbrain.** Messages
arrive only as a file the caller passes. The Secondbrain is read from the
directory User Config names, and is never fetched or pulled.

## Considered Options

- Ask Jev to rank the whole list in one choice question: rejected. A ranking
  is a broad judgment that hides several factors, and TypeSafe's own guidance
  is to split it into one snap judgment per factor and combine the answers in
  code.
- Order by the declared type alone: rejected as the default, kept as the
  fallback. Only an answer can tell that an item blocks another project.
- Keep the English gate of ADR-0201: rejected. It would skip almost every
  Inbox Entry, which is where blocking requests from other projects arrive.
- Let the command call ai-memory: rejected. Only the Agent Session sees the
  MCP, and the command must behave the same in every harness.

## Consequences

- The thresholds are TypeSafe's suggested starting point, not a measurement
  of this judgment; moving them, or the pinned model, means measuring.
- A run asks one request per item, so its cost grows with the open work, under
  the same monthly ceiling.
- A Portuguese item is answered from English questions; the order is
  advisory, and the operator's first real run is the first check of it.
