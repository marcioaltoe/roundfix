---
status: accepted
created_at: 2026-07-06T21:05:00Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# The Run Event Journal is pruned by a retention window

The Run Event Journal (ADR-0008) is append-only and never pruned, so `run_events`
grows without bound — a dogfood database reached ~220 MB, almost all journal
rows. Roundfix therefore bounds it with a retention window: journal events (and
the on-disk artifact directory) of a **terminal** Run older than a configured
duration become eligible for pruning, applied on demand by a GC Command and
best-effort during the operational preflight sweep, while the `runs` row and its
active-run lock — the load-bearing state for one-active-run, stop, detach, and
recovery (ADR-0012/0016/0022/0028) — are never touched by retention. Active Runs
and their journals are never pruned; a retention of `0` keeps everything, so the
default is a bounded window, not deletion of live history.
