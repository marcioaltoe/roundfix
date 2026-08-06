---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-06
---

# Three gigabytes of event journal inside the retention window

**Date:** 2026-08-06
**Found by:** the maintainer noticing `~/.roundfix/roundfix.db` at 3.1 GB and
asking whether any of it is useful.

## Measurements

```text
roundfix.db      3.1G
roundfix.db-wal  4.4M
```

By table, from `dbstat`:

| object | MB |
| --- | --- |
| `run_events` | 2915 |
| `run_events` autoindex | 86 |
| every other table | ~0 |

The journal holds **1,645,457 events**: 736,318 stamped July and 909,139
stamped August. The largest single Run holds **42,000 events**; the next two
hold 23,692 and 17,539. Average payload works out to roughly 1.8 KB per event.

## The retention machinery is working

`roundfix gc --dry-run` reports Journal Retention of 336h, 302 eligible Runs,
**0 journal rows eligible**, and **0 artifact bytes reclaimable**. A direct
query agrees: zero events carry a `run_id` older than `20260723`, the
fourteen-day boundary.

So nothing is leaking past retention. The three gigabytes are **fourteen days
of current work**, at the granularity the Run Event Stream records: every tool
call, state transition, and payload of every Run.

## The question the size actually raises

The Run Event Stream exists for live supervision — `roundfix events --follow`
while a Run is in flight. Once a Run reaches a terminal outcome, its console
log already persists as a file artifact, and the journal's remaining value is
forensic.

That suggests the retention policy is asking the wrong question. Instead of
"how old is this Run", the useful question is "is anyone still watching it":

- payloads could be dropped at terminal settlement while event headers stay,
  keeping the timeline readable and shedding nearly all the bytes;
- Journal Retention could differ for terminal and non-terminal Runs;
- payloads could be compressed, which costs nothing in behavior.

## Why this is not only about disk

This is the same database whose `busy_timeout` was raised from 5s to 30s
today after cross-project `SQLITE_BUSY` reports. Write transactions against a
three-gigabyte table with a hot 4.4 MB WAL hold locks longer than against a
small one, and the Daemon writes events continuously while Tasks commit. Disk
is the cheap part of this observation; lock contention is the expensive part.

A measurement worth taking before designing anything: event write latency and
`SQLITE_BUSY` frequency as a function of journal size.
