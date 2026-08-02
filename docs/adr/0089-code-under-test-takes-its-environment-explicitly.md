---
status: accepted
created_at: 2026-08-02T00:00:00Z
updated_at: 2026-08-02T00:00:00Z
deprecated_at: null
superseded_by: null
---

# Code under test takes its environment explicitly

The verification suite's floor is `internal/cli` at 113 seconds of essentially
sequential work — 488 test functions against one `t.Parallel()` call — and the
reason it cannot parallelise is process-global state: 20 `t.Setenv` and 18
working-directory changes, each of which Go refuses to combine with
`t.Parallel()` because both mutate state the whole process shares. The
resolution is to make the code under test receive its environment and working
directory as arguments rather than reading them from the process, so a test
supplies them per call and needs no global mutation; where a test exists
precisely to verify that the process-level default is read correctly, it keeps
the global and states in one line why it stays sequential. Leaving the tests
sequential was rejected because the suite can never finish faster than its
slowest package and every Task pays that floor; running each test in its own
process was rejected because it trades a shared-state problem for a
process-spawn cost across 488 tests and hides the coupling instead of removing
it.
