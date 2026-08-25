---
status: completed
type: test
---

# Task: Acceptance Verification

Verify three measured hook refusal cases.

## Work
- Case 1: 82-line function over 80-char limit
- Case 2: 2462-line file over 500-line limit
- Case 3: `sort()` instead of `toSorted()`

## Verification
- `grep -q "TestHookRefusalRecovery\|hook.*refusal.*recovery" internal/daemon/task_engine_test.go && go test -count=1 ./internal/daemon 2>&1 | grep -q "ok.*internal/daemon"`


## References

- User Story 1: Three measured cases
- User Story 3: Acceptance verification

## Result

Three cases resolve via settle without losing work.

### What changed

The three measured Run deaths are now covered by executable acceptance tests
that compute each finding over real content rather than declaring it. Every case
carries the work its Run died on — a function of 82 lines under an 80-line rule,
a 2462-line generated file under a 500-line rule, and `Array#sort()` where the
rule requires `toSorted()` — and a repository hook that enforces that rule while
the Task's declared Verification passes.

- `TestHookRefusalRecovery` (`internal/daemon/task_engine_test.go`) drives each
  case through the Daemon's commit boundary against a real repository with a
  real refusing `pre-commit` hook: the refusal is classified, the Task keeps the
  completed status its Verification earned, the work stays staged in the surface,
  the Run Event and Progress name the finding and the recovery command, and the
  commit settle performs then lands the work byte for byte once the over-strict
  rule leaves the hook.
- `TestSettleRecoversMeasuredHookRefusedWork` (`internal/cli/settle_test.go`)
  settles the same three findings through the `settle` Command itself. The first
  settle meets the hook still in place: it re-verifies, stages, is refused, and
  reports the refusal without discarding anything. After the misconfiguration is
  repaired, the second settle commits the same work untouched and leaves a clean
  checkout.
- `initHookRepoForTest` was extracted from `newHookRepoForTest`
  (`internal/daemon/commit_hook_test.go`) so a fixture repository built elsewhere
  can take a repository-local hooks path without duplicating it.

### Evidence per acceptance criterion

| Criterion | Evidence |
| --- | --- |
| Case 1 — 82-line function over the 80-line limit resolves via settle | `go test -count=1 -v -run TestHookRefusalRecovery ./internal/daemon` → `--- PASS: TestHookRefusalRecovery/function_over_the_80-line_limit`; the hook's own finding `src/parser.go:3: function exceeds the 80-line limit` is asserted on the refusal, the Run Event and the Progress output, and the file is recovered byte for byte. The settle Command half is `TestSettleRecoversMeasuredHookRefusedWork`, which carries the same file. |
| Case 2 — 2462-line generated file over the 500-line limit resolves via settle | Same run → `--- PASS: TestHookRefusalRecovery/generated_file_over_the_500-line_limit`; finding `internal/api/schema.gen.ts: 2462 lines exceeds the 500-line limit`, recovered byte for byte; same file settled by `TestSettleRecoversMeasuredHookRefusedWork`. |
| Case 3 — `sort()` instead of `toSorted()` resolves via settle | Same run → `--- PASS: TestHookRefusalRecovery/sort()_where_the_rule_requires_toSorted()`; finding `src/list.ts:2: use toSorted() instead of sort()`, recovered byte for byte; same file settled by `TestSettleRecoversMeasuredHookRefusedWork`. |
| No work is lost in any case | Both tests assert the surface after the refusal: nothing committed (`rev-list --count HEAD` unchanged), the work still staged/present with identical bytes, the Task file still `completed`. After recovery, `git show HEAD:<path>` equals the written content exactly and `git status --porcelain=v1` is empty. |

### Focused checks

- `go test -count=1 -run TestHookRefusalRecovery ./internal/daemon` → ok (3 subtests).
- `go test -count=1 -run TestSettleRecoversMeasuredHookRefusedWork ./internal/cli` → ok.
- `go test -count=1 ./internal/daemon ./internal/cli` → `ok roundfix/internal/daemon 3.033s`, `ok roundfix/internal/cli 53.267s`.
- `gofmt -l internal/daemon internal/cli` → no output; `go vet ./internal/daemon ./internal/cli` → clean.
- Mutation checks, to prove the tests detect rather than assert: forcing
  `ClassifyCommitHookRefusal` to classify nothing failed all three daemon
  subtests (`expected the hook refusal classified`); making each case's content
  comply with its rule failed all three daemon subtests (`got <nil>`) and the
  settle test (`expected settle to stop on the hook refusal with exit 1, got 0`).
  Both mutations were reverted.

### Follow-ups

- None. The settle Command's surface resolution and re-verification stay covered
  by the Task 02–05 tests; this Task adds only the measured-case acceptance.
