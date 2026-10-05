---
spec: 0231-checks-that-hold-in-delivery
prd: _prd.md
created: 2026-10-05
---

# Checks that hold in delivery — Technical Spec

## Executive Summary

Two checks judge something other than the work. The darwin reading of a
detach fixture's process group counts a member that is already exiting as
live, while the kernel already refuses to count it as a signal target; the
reading gains the exiting flag the process table reports. CI tests GitHub's
event merge, which a re-run never refreshes, and the Delivery Queue
classifies the failure without asking which default branch it tested. CI's
pull request job instead merges the default branch tip it fetches and records
it as a `tested-base` annotation, and the queue compares that tip with the
delivery remote's default branch before it classifies a failure. The primary
trade-off is that CI tests a merge it computes rather than GitHub's test
merge, and the queue depends on a label the workflow writes; a failure
without the label falls back to the existing rules, so the dependency can
only make the queue more patient, never pass a check (ADR-0236).

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; the
  annotation title `tested-base` is a fixed label. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — `InspectFailedCheck` adds one read
  through the authenticated `gh` CLI: the failing job's check-run annotations
  from the repository's own GitHub API. No credential is read, stored or sent
  by Roundfix; no test reaches GitHub, because every command goes through a
  scripted `CommandRunner`. The CI merge step reads only refs the checkout
  fetched. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0236: "A failure whose tested base
  does not contain that tip is stale", and "A failure without the annotation
  keeps the rules it had before". ADR-0213: "A process a test starts,
  directly or through a detached Run, must end when the test binary that
  started it ends"; the reading still proves that end. ADR-0028 owns detach,
  unchanged. ADR-0179 bounds `.github/workflows/ci-verify.yml`. The QA gate
  follows ADR-0080: "QA verdicts distinguish environment-blocked rows", and
  ADR-0091: "required to be terminal and to depend on every leaf"; ADR-0096,
  ADR-0097, ADR-0104, ADR-0167, ADR-0182, ADR-0194, ADR-0195 and ADR-0210 bind
  its stage, row carry, outside evidence, Pull Request row, settlement, row
  record, re-observation and evidence snapshot. ADR-0093, ADR-0117, ADR-0156,
  ADR-0168, ADR-0176 and ADR-0183 check consistency by citation and receipt.
  ADR-0184 does not apply: no command's output changes. ADR-0030 cites
  ADR-0028 but decides that agent run logs are opt-in, ADR-0229 cites
  ADR-0167 but decides how an operator archive resumes a park, and ADR-0098
  cites ADR-0030 but decides how Run Events append in batches; this Spec
  changes none of them, so none applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — task_02 changes the Governed Path
  `.github/workflows/ci-verify.yml`, and only it; express maintainer
  authorization: the named grant "Concedo" of 2026-10-05 for the paths this
  Spec declares; bounded files: `.github/workflows/ci-verify.yml`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0231-checks-that-hold-in-delivery/_authorization.md`.

## System Architecture

No new package, command or flag.

| Component | Where | Change |
| --- | --- | --- |
| Fixture group reading (darwin) | `detachFixtureGroupLiveMembers` in `internal/cli/detach_fixture_group_darwin_test.go` | Also excludes members whose `kinfo_proc` `p_flag` has `P_WEXIT` |
| CI pull request job | `.github/workflows/ci-verify.yml`, job `verify` | Checks out the head, merges the fetched default branch tip, exports it as `VERIFY_BASE`, emits the `tested-base` notice |
| Failed-check inspection | `CheckFailure`, `InspectFailedCheck`, `RerunFailedCheck` in `internal/delivery/check_rerun.go` | Reads the tested base and the default branch tip, sets `Stale` |
| Check classification | `checkCandidate` in `internal/delivery/engine.go` | Re-runs a stale failure instead of classifying it |
| Record | `docs/user-guide/commands/deliver.md` | Describes the stale rule beside the re-run rule |

## Implementation Design

### Interfaces

```go
// internal/delivery/check_rerun.go
type CheckFailure struct {
	RunID         string
	Attempt       int
	Packages      []string
	OutsideChange bool
	TestedBase    string // from the job's tested-base annotation; "" when absent
	DefaultTip    string // delivery remote's default branch tip, set with TestedBase
	Stale         bool   // TestedBase set and not containing DefaultTip
}

const testedBaseAnnotationTitle = "tested-base"
```

```text
1. InspectFailedCheck reads the attempt, then, when the check link names a job (".../actions/runs/<run>/job/<job>"), the job's annotations.
2. TestedBase is the message of the only annotation titled "tested-base" whose message is 40 lowercase hex characters; none, several, or a malformed one leave it "".
3. With a TestedBase, it resolves and fetches the remote default branch as it does today, sets DefaultTip, and sets Stale when the tip is not an ancestor of, or equal to, TestedBase, or when TestedBase is not a commit after the fetch. A Stale failure returns before the failed log is read.
4. Without a TestedBase, InspectFailedCheck behaves exactly as before.
5. RerunFailedCheck accepts a failure with a valid run id that is Stale or OutsideChange; the engine owns the attempt bound.
6. checkCandidate never parks a Stale failure as flaky-check or checks-failed. It re-runs it once per (run, DefaultTip), logs API Contract 4, restarts the check timeout at most once in total, and keeps polling; a Stale failure whose (run, DefaultTip) it already re-ran keeps polling.
7. A failure that is not Stale is classified by the existing rules, except that the first failure of a run re-run as Stale, on a newer attempt, is eligible for the one outside-change re-run, and a failure whose attempt is not newer than the attempt re-run for it keeps polling.
8. A re-run error or an inspection error parks checks-failed and is logged, as today; the check timeout still ends in checks-timeout.
9. The darwin fixture reading excludes a member with P_WEXIT (0x2000) or p_stat SZOMB; a member without both is live.
10. The CI job's pull request path tests the head merged with refs/remotes/origin/<base_ref> as fetched at job start; its push path is unchanged.
```

### Data Models

`CheckFailure` gains `TestedBase`, `DefaultTip` and `Stale`. The engine keeps
two in-memory sets for one `checkCandidate` call: the (run, tip) pairs it
re-ran as stale and the runs it re-ran as stale, plus the attempt it re-ran
for each run. Nothing is persisted; the Run Database and the queue item are
unchanged.

### API Contracts

1. API Contract: the CI pull request job prints exactly one
   `::notice title=tested-base::<sha>` line, where `<sha>` is the full commit
   id of `refs/remotes/origin/<base_ref>` it merged, before the `Verify
   changed` step, and writes `VERIFY_BASE=<sha>` to `$GITHUB_ENV`.
2. API Contract: the queue reads annotations with
   `gh api repos/{owner}/{repo}/check-runs/<job-id>/annotations?per_page=100`
   in the item's work directory, and decides staleness with
   `git cat-file -e <tested-base>^{commit}` and
   `git merge-base --is-ancestor refs/remotes/<remote>/<branch> <tested-base>`
   after the existing fetch of the default branch.
3. API Contract: a stale re-run uses the existing
   `gh run rerun <run-id> --failed`.
4. API Contract: the owner log line for a stale re-run is
   `roundfix: check stale: Delivery Queue item <slug>: <check> tested <tested-base>, default branch is at <tip>; re-run (run <run-id>)`.
5. API Contract: the existing log lines, blockers (`flaky-check`,
   `checks-failed`, `checks-timeout`), Park Classes and the
   `flaky-check: <check> passed on re-run` Warning keep their text; a pass
   after a stale re-run adds no Warning.

### Surface Transcripts

None. No command's standard output, standard error or exit code changes; the
queue adds one line to its owner log (API Contract 4), recorded in the
Vocabulary Contract.

## Coverage Map

- Goal 1 → darwin `detachFixtureGroupLiveMembers` (Invariant 9)
- Goal 2 → CI pull request job (Invariant 10, API Contract 1)
- Goal 3 → `InspectFailedCheck`, `checkCandidate` (Invariants 1 to 6)
- Goal 4 → `InspectFailedCheck` without a tested base, `checkCandidate` (Invariants 4, 7, 8; API Contract 5)
- Core Feature 1 → darwin fixture reading
- Core Feature 2 → CI pull request job
- Core Feature 3 → `CheckFailure`, `InspectFailedCheck`, `RerunFailedCheck`, `checkCandidate`, deliver guide
- Success Metric 1 → darwin reading and its exiting-member test
- Success Metric 2 → the CI merge step rehearsal
- Success Metric 3 → `checkCandidate` replay of #391
- Success Metric 4 → existing check re-run tests, Invariant 4

## Integration Points

- GitHub Actions: the workflow's checkout uses `actions/checkout` with
  `ref` set to the pull request's head commit on `pull_request` events and
  `fetch-depth: 0`, which fetches every branch at job start.
- GitHub REST through `gh api`: one annotation read per inspected failure.
- Git: the existing fetch of the delivery remote's default branch.

## Testing Approach

- The darwin reading is tested by a new darwin test in
  `detach_fixture_group_darwin_test.go`: it starts short-lived children in
  their own groups and polls until it observes, through `kern.proc.pgrp`
  read directly, a group whose only member is not a zombie and has `P_WEXIT`
  while `kill(-pgid, 0)` returns `EPERM`; at that moment the reading must
  report no live member. It records how many launches it needed and fails at
  its `testwait` deadline if it never sees the window. The existing
  unreaped-member test keeps proving that a running member is live.
- The CI step is rehearsed by its Verification in a disposable clone, with
  the step's `run` text read from the workflow file.
- The queue is tested through the existing seams: `fakeCheckRecovery` and
  `exerciseCheckRecovery` for `checkCandidate`, and `scriptedCommandRunner`
  for `InspectFailedCheck` and `RerunFailedCheck`. The scripted command
  sequences of the existing `InspectFailedCheck` tests change only where a
  job id and an annotation read are added.

## Build Order

1. Count an exiting member as ended in the darwin fixture reading.
2. Make CI's pull request job test and record the current default branch tip.
3. Judge a failed check on the tip it tested in the Delivery Queue, and
   describe the rule in the deliver guide (depends on: 2).
4. Final QA gate (depends on: 1, 2, 3).

## Risks & Considerations

- If the window test never observes an exiting member on some machine, it
  fails at its deadline rather than skipping; the probe saw the window in
  every launch, so a failure would mean the kernel changed.
- A default branch that keeps moving re-runs once per tip until the check
  timeout parks `checks-timeout`.
- The first Pull Request after task_02 merges is the first CI run of the new
  step; until it merges, a failure carries no annotation and is classified as
  before.

## Vocabulary Contract

- emits: `internal/delivery/engine.go`
  pattern: `check stale: Delivery Queue item`
  documented-in: `docs/user-guide/commands/deliver.md`
- emits: `internal/delivery/check_rerun.go`
  pattern: `tested-base`
  documented-in: `docs/user-guide/commands/deliver.md`

No glossary term is adopted; "tested base" is used in its plain sense beside
**Delivery Queue** and **Park Class**.

## Decisions

- Read the exiting flag instead of retrying the kill or waiting for the
  zombie: the flag is set at the point of no return, and a retry would
  depend on timing.
- Merge in CI and record the tip, and judge a failure by it. See ADR-0236.
- Keep the stale bookkeeping in memory for one `checkCandidate` call.
