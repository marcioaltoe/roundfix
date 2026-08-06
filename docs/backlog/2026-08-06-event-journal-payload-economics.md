---
type: perf # feat | fix | perf | refactor
status: promoted # open | promoted | declined
created: 2026-08-06
spec: 0081-a-journal-cheap-to-write-and-keep
reason: null # required when status: declined
---

# Retain the event journal by who is watching, not by how old it is

## Opportunity

The Run Event Journal is 2.9 GB of a 3.1 GB database, and every byte of it is
inside the fourteen-day retention window — 1,645,457 events at roughly 1.8 KB
of payload each, with one Run alone holding 42,000. Retention is working
exactly as configured: `roundfix gc --dry-run` finds 302 eligible Runs with
zero prunable rows, and no event carries a `run_id` older than the boundary.

The policy is what deserves the change. Journal Retention asks how old a Run
is; the useful question is whether anyone is still watching it. The Run Event
Stream exists for live supervision — `roundfix events --follow` — and once a
Run reaches a terminal outcome its console log already persists as a file
artifact, leaving the journal's payloads forensic at best.

Candidate shapes, cheapest first: give terminal and non-terminal Runs
different retention, which changes no payload at all. Anything that drops or
rewrites a payload — including compression — is blocked by ADR-0008, which
makes the payload raw producer JSON, write-once and read-as-blob, and by
ADR-0030, which removed the per-Batch agent logs on the grounds that the
journal is the durable copy. Such a shape may only arrive as an explicit
amendment naming the lost capability; Spec 0081 defers that decision behind a
measurement that may conclude it is unnecessary.

## Why it is not only about disk

This is the database whose `busy_timeout` was raised from 5s to 30s on
2026-08-06 after cross-project `SQLITE_BUSY` reports. Write transactions
against a three-gigabyte table with a hot WAL hold locks longer than against a
small one, and the Daemon writes events continuously while Tasks commit and
settle. The measurement to take before designing anything is event-write
latency and `SQLITE_BUSY` frequency as a function of journal size.

## Relation to the gate work

This entry and `2026-08-06-two-stage-qa-gate-economics.md` are two ends of one
loop. Events are produced by agent turns, so the gate's cost and the journal's
growth share a cause: a gate round that re-audits every row with an agent
writes tens of thousands of events for work a grep could have settled. Cutting
gate rounds cuts journal bytes, and moving mechanical checks out of agent turns
cuts both directly.

The maintainer directed on 2026-08-06 that both be attacked next. Two measured
goals belong together in whichever Spec promotes them: **time to verdict** and
**bytes per Run**.

## Evidence

`docs/findings/2026-08-06-three-gigabytes-of-event-journal-inside-the-retention-window.md`
carries the measurements, the `dbstat` breakdown, and the GC dry-run output.
