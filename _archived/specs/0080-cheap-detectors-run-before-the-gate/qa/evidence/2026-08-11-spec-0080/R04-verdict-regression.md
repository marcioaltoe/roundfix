# R04 — verdict regression and root cause

Focused lower layers pass:

- `TestQAVerdict*`: 24 passed in `internal/spec`.
- Mechanical corpus/detector selection: 18 passed in `internal/speccheck`.

The public CLI matrix fails consistently. The executable reproduction is
`R04-minimal-reproduction.sh`; it selects only
`TestRunImplementQAVerdictMatrix/pass` and exits 1.

Root-cause trace:

1. `runQAGate` creates the mechanical report before the Agent
   (`internal/daemon/task_engine.go:1798`).
2. A non-blocking seed intentionally has no verdict
   (`internal/daemon/task_engine.go:2017-2026`).
3. The CLI integration Agent fake still writes the fixed
   `qa-report-2026-01-01.md` (`internal/cli/implement_test.go:754-765,799`).
4. The new seed is `qa-report-2026-08-11.md`, and `NewestQAReport` correctly
   chooses the later embedded date (`internal/spec/qa.go:34-72`).
5. Settlement therefore reads the untouched seed and returns `unreadable`
   instead of the scripted pass, partial, fail, missing, or unreadable result.

Task 03 updated its Daemon-unit fake to complete the seeded path, but did not
update the higher-level CLI fake. No fix was attempted during report-only QA.
