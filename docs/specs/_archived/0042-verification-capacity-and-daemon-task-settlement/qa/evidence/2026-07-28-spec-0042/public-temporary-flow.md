# Temporary Verification Failure flow

Current-build CLI/Daemon integration uses disposable Git repositories and real
shell Verification processes while faking only the external ACP boundary.

Passing checks covered:

- exit `75` once, followed by a passing exclusive retry and zero Agent repair;
- repeated exit `75`, one exclusive retry only, terminal Task failure, and
  preserved Task Worktree;
- deterministic non-`75` failure followed by exactly one Agent repair and
  numbered attempt 2;
- exit `1` containing timeout, listener, database, and port text remaining
  deterministic;
- distinct `attempt-1.log` and `attempt-1-retry-1.log` identities;
- bounded public event fields for classification, retry availability, mode,
  reason, verdict, and exhaustion.

The focused, repeated, race, full-suite, and full-race commands all exited 0.
