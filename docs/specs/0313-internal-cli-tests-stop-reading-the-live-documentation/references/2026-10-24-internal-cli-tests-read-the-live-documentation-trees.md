---
type: perf
status: promoted
created: 2026-10-24
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
reason: null
---

# internal/cli tests read the live documentation trees

## Slow

Thirteen `internal/cli` tests read this repository's own `docs/history`,
`docs/adr`, `docs/specs`, `docs/user-guide` or `.agents/skills`. Go reuses a
cached test result only when every file the test process read inside the
module is unchanged, so a change to a command guide or a skill makes the
whole package run again. `internal/cli` is the longest package of the suite,
and many Tasks edit a guide or a skill.

## Measured

On 2026-10-24, in a worktree of `56dd1a76` whose `internal/cli` result was
cached, one line appended to `docs/user-guide/commands/gc.md` made
`go test -parallel 16 ./internal/cli` run again: 162.6 s, and 237.0 s after a
second line. With the thirteen tests skipped (`-skip`), both edits left the
result cached; a change to a Task file's `status:` line left it cached in
both forms. With the five trees removed, and again with them unreadable, the
same thirteen tests fail and no other. Details in
`docs/references/2026-10-24-test-suite-audit.md`.

## Target

A guide or skill edit leaves the `internal/cli` test result cached: the tests
that check the repository's own documents run under the `docscontract` build
tag, where the selective gate's Contract Relevance decides when they run,
and every assertion they make today still fails on the inputs it fails on.
