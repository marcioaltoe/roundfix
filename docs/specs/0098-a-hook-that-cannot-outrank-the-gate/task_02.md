---
status: completed
type: backend
---

# Task: Extend Settle to Accept Completed Status

Update settle command contract to accept `completed` status.

## Work
- Modify settle preflight to accept `completed`
- Resolve settle surface in priority order
- Load Task file from selected surface

## Verification
- `grep -q "completed" internal/cli/settle.go && grep -q "taskStatus == \"completed\"" internal/cli/settle.go`


## References

- User Story 2: Settle accepts completed
- Core Feature 2: Settle Recovery

## Result

Settle now recovers the state a refused commit leaves behind: a Task the Daemon
settled `completed` after its Verification passed, whose verified work never
left the surface that produced it. The status alone does not open the door —
the surface must still hold that work uncommitted, so a `completed` Task whose
work is already committed keeps refusing exactly as before.

### Implementation

- `internal/cli/settle.go` gains `settleSurfaceSettles`, the single place that
  decides whether a candidate surface settles. `failed` is the ordinary target;
  `completed` qualifies only when the surface reports uncommitted work, read
  through the same `worktree.Snapshot` seam the settle commit path already uses,
  so "uncommitted work" means exactly what settle would commit.
  `resolveSettleSurface` takes a `context.Context` and calls it in place of the
  old `task.Status == spec.StatusFailed` test; the candidate order — Task
  Worktree, then Run Worktree, then the checkout — is untouched.
- `settleSurfaceReport` gains `detail`, so a refusal names *why* a candidate was
  passed over rather than only its status: `<path>: status completed (no
  uncommitted work)`.
- The refusal contract follows the new acceptance: `settleNoFailedSurfaceError`
  becomes `settleNoSurfaceError` ("has no settleable surface"), and the guidance
  closes on the shared `settleStatusRequirement` — `settle requires failed, or
  completed with uncommitted work`. The completed-surface message becomes
  "status is completed with no uncommitted work; nothing to settle", which is
  the condition actually measured, where the old "nothing to do" was now false
  whenever a surface was dirty.
- Two places assumed the settled Task arrived `failed`. The verification-failure
  line prints the Task's real status (`task_01 stays completed — verification
  failed`), and `settleTaskAndCommit` adds the task file to the commit set only
  when there is a status flip to commit — an already-`completed` task file is
  rewritten byte-identically, so forcing it into the commit printed a `commit
  <task file>` line for a path the commit did not contain.
- `settleUsage` states the two accepted states and drops "status stays failed"
  from exit `1`. `docs/user-guide/commands.md` and `docs/user-guide/usage.md`
  carry the same contract change.

### Acceptance criteria

1. **Settle accepts `completed` (it previously rejected).**
   `TestSettleAcceptsCompletedTaskWithUncommittedWorkInCheckout` seeds a
   `completed` Task with uncommitted work in the checkout and asserts exit `0`,
   the full stdout (`verify … — ok`, `commit done.txt`, `settled task_01
   completed — <sha>`), that the commit contains only `done.txt`, that the task
   file stays `completed`, and that the checkout ends clean.
   `TestSettleRefusesCompletedTaskWhenNoSurfaceHoldsUncommittedWork` is its
   negative control: same status, clean surfaces, exit `2`, HEAD unchanged, and
   the refusal names `no uncommitted work` for both candidates.
2. **Surface resolved in priority order.**
   `TestSettleAcceptsCompletedTaskFromKeptTaskWorktreeAfterHookRefusal` builds
   the refusal shape directly — work `git add`ed but never committed in a kept
   Task Worktree, its task file `completed`, the checkout's copy still
   `in_progress` — and asserts `Settle surface:` is the Task Worktree, that the
   work reaches the user checkout, and that both worktrees are cleaned up.
   `TestRunSettleRefusalEnumeratesCandidateStatuses` still proves the order is
   reported candidate by candidate.
3. **Task file loaded from the selected surface.** The same test gives the two
   surfaces different Verification commands: the Task Worktree's file declares
   `test -f done.txt` and the checkout's declares `test -f never-written.txt`.
   The Run passes, and stdout carries `verify test -f done.txt — ok` — the
   Verification could only have come from the selected surface.

### Focused checks

- `go test ./internal/cli/ -run 'TestSettleAccepts|TestSettleRefuses' -count=1 -v`
  — 3 new tests, all pass.
- Negative control: with `internal/cli/settle.go` restored to `HEAD` and the new
  tests unchanged, all three fail; restored and they pass. The tests discriminate
  the change rather than the fixture.
- `go test ./internal/cli/ -count=1` — ok (55.1s); `go test ./internal/daemon/
  ./internal/spec/ -count=1` — ok. The two existing expectations that moved are
  the ones this Task's contract change owns: the refusal headline
  (`no failed settle surface` → `no settleable surface`) and the completed-task
  message.
- `go test -race ./internal/cli/ -run 'TestSettleAccepts|TestSettleRefuses|TestRunSettle' -count=1` — ok.
- `go build ./...`, `go vet ./internal/cli/`, `make fmt-check` — clean.
- `make docs-test` fails on `SC-VERIFY-NON-HERMETIC` in
  `docs/specs/0113-a-refused-gate-writes-its-refusal-once/task_01.md`. Measured
  pre-existing: the same two tests fail identically with this Task's two
  documentation edits stashed.

### Follow-up for another Task

- `.agents/skills/roundfix/SKILL.md` still tells the Supervisor to settle "each
  failed Task". The skill is outside this Spec's bounded authorization
  (`internal/daemon`, `internal/cli/settle.go`, `docs/agents/autonomous-work.md`)
  and its edit would require `make baseline-digests`, so the skill-sync rule is
  left to the Spec's closing Task rather than partially applied here.
- `settleStatusError` in `internal/cli/settle.go` has no caller. It predates this
  change and is untouched by it.
