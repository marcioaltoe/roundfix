---
spec: 0081-a-journal-cheap-to-write-and-keep
status: active
created: 2026-08-06
surfaces: [backend, data, docs, cli]
---

# A journal cheap to write and cheap to keep

The Run Event Journal is 2.9 GB of a 3.1 GB Run Database, and every byte of it
is inside the fourteen-day retention window: 1,645,457 events, one Run holding
42,000, roughly 1.8 KB of payload each. Retention is working exactly as
configured — `roundfix gc --dry-run` finds 302 eligible Runs and zero prunable
rows, and no event predates the boundary.

The size is the symptom that made this visible, but the cost is not disk. Three
paths make the journal expensive in ways a smaller database would not fix:
every Run start takes the machine-wide write lock and holds it across a full
aggregate of the largest table, even when nothing is eligible; every agent
output line is its own immediate transaction and its own fsync; and the live
cockpit re-reads the entire journal, payloads included, on every append. The
`busy_timeout` raised from 5s to 30s on 2026-08-05 after cross-project
`SQLITE_BUSY` failures treated the symptom of the first of those.

## Project Constraints

- Identifier strategy: not applicable — Run identifiers and event cursors
  already exist and are unchanged; no project-owned Internal Identifier is
  minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the work is local SQLite access
  and file I/O; no endpoint, credential, or transport is created. Source:
  `docs/agents/go.md`.
- Active ADR obligations: applicable and, in one case, in tension. ADR-0104
  requires at least one acceptance row to rest on evidence this Spec did not
  author; the row is recorded blocked with its reason when the source cannot be
  obtained, and the QA gate accounts for it before delivery. ADR-0098 is
  the decision this Spec implements: appends batch, with durability per batch
  and `synchronous` unchanged. ADR-0008
  states that Run Event payload stores raw producer JSON, write-once and
  read-as-blob, and that producers must never re-serialize or prune payload
  JSON; any payload-shedding or compressing design must supersede or amend it
  explicitly, never quietly. ADR-0030 removed per-Batch agent log files
  *because* the journal already stores every raw payload durably, so the
  journal is the sole durable copy and dropping payloads removes a capability
  rather than a duplicate. ADR-0033 owns the current retention policy —
  age-based, terminal-only, `runs` rows and locks never touched, zero means
  keep everything. ADR-0009 makes the cockpit read the journal exclusively,
  so any read-path change affects live rendering and not only replay.
  ADR-0004 owns the single machine-wide Run Database that makes contention
  cross-project. ADR-0022 does not apply: Stop Requests travel through the Run Database as
  control state, and this Spec changes event persistence and read paths
  without touching Run state, Stop semantics, or their tables. ADR-0027 owns
  how a removed or renamed config key degrades to
  a warning. ADR-0080 owns QA verdict semantics and the typed blocked-cause
  counts, and ADR-0091 keeps the authored QA gate a terminal Task node, under
  which this Spec's own graph is authored. ADR-0093, ADR-0096, and ADR-0097 do
  not apply: this Spec adds no consistency check, no gate stage, and no
  carried report row — they belong to Spec 0080, which shares only the QA
  verdict contract with this one. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the work is Go code, SQL, a user-guide
  document, and Go-side config defaults. It creates, edits, renames, moves, or
  deletes no repository-tooling configuration, script, ignore file, plugin
  declaration, or version pin. This repository's own `.roundfixrc.yml` is
  explicitly out of scope; defaults ship in code. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- Starting a Run never blocks on the size of the journal.
- Writing an event costs an amortized fraction of a transaction, not a
  transaction and an fsync per output line.
- A reader pays for what it reads: nothing loads payload bytes it will
  discard, and nothing rescans from the beginning of a Run to render its
  latest line.
- Journal Retention keeps its terminal-only, age-based contract, while the
  eligibility query becomes bounded by the configured cutoff and stops holding
  the write lock across a full-journal scan. No payload is shed inside the
  retention window.
- Parallel Runs on one machine stop being a contention story; the 30-second
  timeout stops being what makes them survivable.

## User Stories

1. As a session starting a Run, I want Run start to be independent of journal
   size, so a fleet machine with two weeks of history starts as fast as a
   fresh one.
2. As a Daemon appending agent output, I want event writes amortized across a
   batch, so a talkative Run does not pay one lock acquisition and one fsync
   per line.
3. As an operator running six Runs on one machine, I want concurrent Runs to
   coexist without `SQLITE_BUSY`, because a longer timeout made contention
   survivable without making the design concurrent.
4. As an operator watching a long Run in the cockpit, I want rendering cost to
   track new events, not total events, so a 42,000-event Run stays responsive.
5. As a supervisor consuming `roundfix events`, I want the stream contract and
   every daemon payload field it depends on to keep working byte-for-byte.
6. As a maintainer deciding retention, I want the trade-off stated in terms of
   what becomes unrecoverable, because the journal is the only durable copy of
   agent payloads.

## Core Features

1. **Measurement before design.** Event-write latency, lock-wait time, and
   `SQLITE_BUSY` frequency measured as functions of journal size and
   concurrent writer count, recorded as a baseline artifact this Spec's later
   claims are checked against. No performance change ships without a before
   and after from the same harness.
2. **Run start stops scanning the journal.** Retention eligibility is answered
   by a query bounded by the retention cutoff rather than by aggregating the
   whole event table, and the eligibility scan stops happening inside the
   write transaction. The reported event count is either derived cheaply or
   dropped from the hot path.
3. **Appends amortize.** Agent output events are written in batches within one
   transaction instead of one transaction per line, with durability behavior
   stated explicitly rather than inherited from defaults.
4. **Reads project what they need.** A payload-free read path exists for the
   consumers that only use headers, and the live cockpit advances a forward
   cursor instead of rescanning a Run from its first event.
5. **Retention semantics stay unchanged.** Only terminal Runs older than the
   configured Journal Retention window remain eligible, and pruning still
   removes their complete Run Event Journal. This Spec makes that eligibility
   query cheap; it adds no second window, payload side table, payload shedding,
   or new data-loss boundary.
6. **Concurrency is a declared story, not a timeout.** The Spec states what
   guarantees hold for parallel Runs on one machine and proves them, whether
   through write batching, connection discipline, or a scoped database path.
7. **The durable-table lifecycle document stays true.** Any table this Spec
   adds or changes lands in the machine-checked lifecycle table in the same
   change.

## User Experience

Nothing about the operator's vocabulary changes: `roundfix events`,
`roundfix attach`, the cockpit, and `roundfix gc` keep their commands, output
shapes, and stream schema. What changes is felt rather than read — Runs start
without a pause proportional to history, six parallel Runs stop dying in their
first Batch, and a long Run's cockpit stays responsive at hour three. The one
retention contract does not change: `roundfix gc`, the configured Journal
Retention window, and the lifecycle guide keep their current meaning. The
query becomes cheaper without creating a new point at which the journal's only
durable payload copy disappears.

## Non-Goals / Out of Scope

- Changing the Run Event Stream schema, its categories, or the `events`
  command contract that supervisors depend on.
- Changing what a Run is, how Runs are scheduled, or any Daemon behavior
  outside event persistence.
- Replacing SQLite, or moving to a second store.
- Reducing the number of events a Run produces — that is Spec 0080's effect,
  not this Spec's mechanism.
- Editing this repository's `.roundfixrc.yml`; defaults ship in code.
- Retro-compacting existing databases beyond what `roundfix gc compact`
  already does under its exclusive fence.
- Changing Journal Retention semantics, adding a second retention window, or
  shedding payloads from retained Runs.

## Success Metrics

- Run start wall-clock is independent of journal size, measured across an
  empty and a full-retention database.
- `SQLITE_BUSY` occurrences across a fixed parallel-Run scenario: zero, at the
  pre-change `busy_timeout` rather than the raised one.
- Event-write cost per 1,000 events improves against the baseline harness, and
  the improvement is attributable to batching rather than to fewer events.
- Bytes per Run, reported before and after, with the retention shape stated.
- Cockpit refresh cost tracks new events rather than total events, measured on
  a Run of at least ten thousand events.
- No consumer regression: `events`, `attach`, cockpit rendering, reconcile
  replay detection, and `gc` all behave identically on a corpus recorded
  before the change.

## Decisions

- Measurement precedes design, per the finding that opened this Spec — the
  three-gigabyte number is the symptom, and the costs worth fixing were only
  visible once the write path was read.
- Lock-holding and write amplification rank above disk size. A journal half
  the size that still aggregates the whole table under the write lock at every
  Run start would fix nothing that hurt.
- ADR-0008 is treated as binding until explicitly amended. A design that keeps
  every payload and still wins on the paths above is preferred to one that
  buys speed by discarding the only durable copy.
- The current terminal-only, age-based Journal Retention contract is preserved.
  Payload shedding and a second retention window are rejected because neither
  repairs the measured lock and write paths, and both create a new data-loss
  boundary for the journal's only durable payload copy.
- Daemon payloads are load-bearing and agent payloads are not, for the
  purposes of the `events` stream: the stream is daemon-only and requires
  those payload fields, while agent payloads serve human timeline rendering.
  This Spec records that asymmetry but does not use it to delete either payload
  class.

## Open Questions

None. The TechSpec owns the writer-concurrency mechanism; this PRD fixes the
product boundary at unchanged Journal Retention semantics and no payload
shedding.
