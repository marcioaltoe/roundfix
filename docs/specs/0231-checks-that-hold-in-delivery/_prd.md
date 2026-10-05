---
spec: 0231-checks-that-hold-in-delivery
status: active
created: 2026-10-05
surfaces: [backend, docs]
---

# Checks that hold in delivery

This is a bug fix with two causes, both recorded on 2026-10-05 in the
operator's intervention log and adopted from two Backlog Entries:
[the survivor test meets EPERM from an exiting owner](references/2026-10-05-the-survivor-test-meets-eperm-from-an-exiting-owner.md)
and [a Pull Request check ran on a stale merge](references/2026-10-05-a-pull-request-check-ran-on-a-stale-merge.md).
Each is a check that gave a verdict about something other than the work it
was checking.

1. **A cleanup that reads an exiting owner as alive.** In a Run's
   `make verify-changed`, `TestRunImplementDetachSurvivesCallerProcessGroupKill`
   failed with `kill fixture process group 52154: operation not permitted;
   live members=[52154]` (entry 161). The detached owner exits on its own
   after its Run reaches Clean, and the test's cleanup kills its group during
   that exit. On macOS the kernel stops counting an exiting process as a
   signal target well before it becomes a zombie, so `kill(-pgid)` returns
   `EPERM`, while the process table still lists the member as running with
   its exiting flag set. Spec 0220's reading excludes only zombies, so it
   reports a live member and the cleanup fails. The authoring probe caught
   this window in every one of 230 launches; the Agent sandbox plays no part.
2. **A failed check judged on an old default branch.** Pull Request #391
   failed CI because of a defect on main. After the fix merged, a re-run kept
   the original merge commit, as GitHub documents, and a reopen four and a
   half minutes later still tested the old merge. The Delivery Queue parked
   the item `flaky-check` and then `checks-failed`; a second reopen passed
   (entries 164 and 165).

## Prerequisites

None. Spec 0230 is being delivered while this Spec is authored; the two
Specs declare no file in common, and no Verification here depends on its
artifacts.

## Project Constraints

- Identifier strategy: not applicable — no identifier scheme changes; queue
  items, Runs, blockers and checks keep their names, and the annotation title
  `tested-base` is a fixed label, not an identifier scheme. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the Delivery Queue gains one GitHub
  read through the authenticated `gh` CLI it already uses: the failing job's
  check-run annotations, through `gh api`. It sends no credential of its own,
  and its destination is the repository's own GitHub API. No test or
  Verification command reaches GitHub; every call goes through a scripted
  command runner. The CI job's new merge step uses the checkout's fetched
  refs and no credential. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0236 (this Spec) decides the CI
  merge and the queue rule: "A failure whose tested base does not contain
  that tip is stale". ADR-0213: "A process a test starts, directly or through
  a detached Run, must end when the test binary that started it ends"; Core
  Feature 1 keeps that proof and stops reading an exiting member as alive.
  ADR-0028 owns detach, which keeps its behavior. ADR-0179 bounds the one
  Governed Path in `_authorization.md`. The authored QA gate follows ADR-0080:
  "QA verdicts distinguish environment-blocked rows", and ADR-0091: "required
  to be terminal and to depend on every leaf"; ADR-0096, ADR-0097 and
  ADR-0167 bind its machine stage, its row carry and its pre-PR Pull Request
  row, ADR-0104: "Every Spec therefore rests at least one named" acceptance
  row on outside evidence, ADR-0194, ADR-0195 and ADR-0210 bind what a QA row
  records, when it is observed again and its evidence snapshot, and ADR-0182
  settles each Task on the facts its gate checks. ADR-0093, ADR-0117,
  ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this Spec's consistency by
  citation and receipt. ADR-0184 does not apply, because no command's output
  changes; the queue adds one line to its owner log. ADR-0030 cites ADR-0028
  but decides that agent run logs are opt-in, ADR-0229 cites ADR-0167 but
  decides how an operator archive resumes a park, and ADR-0098 cites
  ADR-0030 but decides how Run Events append in batches; this Spec changes
  none of them, so none applies. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — `.github/workflows/ci-verify.yml` is a
  Governed Path. On 2026-10-05 the maintainer gave the named grant
  "Concedo" for the governed test files this Spec declares and for
  `.github/workflows`, only the paths the Spec declares; the probe found no
  governed test file among them. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0231-checks-that-hold-in-delivery/_authorization.md`; bounded
  files: `.github/workflows/ci-verify.yml`.

## Goals

- The detach death and survivor tests pass in every environment where they
  ran before, without a skip, and a running fixture is never counted as
  ended.
- Every CI attempt on a Pull Request, including a re-run, tests the head
  merged with the default branch as it stands when the job starts, and says
  which tip it tested.
- The Delivery Queue never parks a failed check as `flaky-check` or
  `checks-failed` when that check tested an older default branch; it
  re-tests and keeps polling.
- A failed check that tested the current default branch, or that says
  nothing about what it tested, is classified exactly as before.

## Core Features

1. **An exiting member is not a live member.** On macOS the fixture's group
   reading excludes members whose process-table entry carries the exiting
   flag, as it already excludes zombies, so the cleanup accepts the `EPERM`
   the kernel returns for a group that is only exiting. Linux and the other
   Unix reading are unchanged.
2. **CI tests the current merge and names it.** The pull request job checks
   out the head commit, merges the default branch tip it fetched, verifies
   the changed sets against that tip and records it as a `tested-base`
   notice annotation.
3. **The queue judges a failure on the current tip.** Before classifying a
   failed check, the queue reads the job's `tested-base` annotation and the
   delivery remote's default branch tip. A stale failure is re-run once per
   tip, logged, and polled; a run re-run as stale keeps its one outside-change
   re-run for its first failure on the current tip.

## Non-Goals / Out of Scope

- Changing detach, `internal/store` process inspection or any production
  process control.
- Re-running checks of other workflows or non-Actions checks, closing or
  reopening Pull Requests, or pushing to an item branch to re-trigger CI.
- Persisting stale re-runs across queue owner restarts.
- The `push` path of CI, the Makefile, `go.mod`, and any skill or agent guide.

## Success Metrics

1. Success Metric: a darwin test that catches a fixture group whose only
   member is exiting (`kill(-pgid, 0)` returns `EPERM`, the member is listed,
   not a zombie, with the exiting flag) reports the group ended, and the
   detach death, survivor and unreaped-member tests pass 10 consecutive
   iterations each (before: the reading reports the exiting member live).
2. Success Metric: the CI merge step, run in a disposable clone whose default
   branch moved after the head branched, leaves a merge whose parents are the
   head and the moved tip, exports that tip as `VERIFY_BASE` and prints the
   `tested-base` notice for it (before: the job tested the event's merge).
3. Success Metric: replaying #391 through the queue with scripted GitHub
   responses, a failure that tested the old tip is re-run, not parked, and the
   item reaches merging when the re-run passes (before: `checks-failed`).
4. Success Metric: every existing check re-run test passes, and a failure
   without a `tested-base` annotation is classified as before.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- GitHub's re-run documentation
  (<https://docs.github.com/en/actions/how-tos/manage-workflow-runs/re-run-workflows-and-jobs>),
  read 2026-10-05: a re-run uses "the same `GITHUB_SHA` (commit SHA) and
  `GITHUB_REF` (git ref) of the original event".
- The CI runs of Pull Request #391 (37306411237 attempts 1 to 3,
  37308972464, 37309659532), whose checkout logs record which merge each
  tested.
- GitHub's workflow commands and check-run annotation documentation
  (<https://docs.github.com/en/actions/using-workflows/workflow-commands-for-github-actions>,
  <https://docs.github.com/en/rest/checks/runs>): a `notice` with a `title`
  creates an annotation, and `GET /repos/{owner}/{repo}/check-runs/{id}/annotations`
  returns its `title` and `message`. On 2026-10-05 a job id from
  `gh pr checks` returned that job's annotations through this endpoint.
- `actions/checkout`'s README: `fetch-depth: 0` fetches "all history for all
  branches and tags".
- XNU's `bsd/kern/kern_exit.c` and `bsd/kern/kern_sig.c`
  (<https://github.com/apple-oss-distributions/xnu>): `proc_exit` sets
  `P_REF_DEAD` ("will not be visible via proc_find") before `p_stat = SZOMB`,
  `killpg1` returns `EPERM` when it counts no member, and `kern_sysctl.c`
  reports `P_LEXIT` as `P_WEXIT`.
- The Daemon's Verification log of Run
  `run_20261005T101255Z_f599a5cc8ea7f4c2` (batch 3, attempt 1), which shows
  the failure.

## Research basis

The Secondbrain was consulted through `wiki/index.md` and
`qmd query "GitHub Actions pull request CI stale merge ref rerun base branch"`
and `qmd query "macOS kill process group EPERM exiting process test"`
(`--all --files --min-score 0.3`); it returned only this repository's mirror
of Spec 0220 and an unrelated QA report, so it holds no prior answer. Exa
found GitHub's re-run documentation, the workflow commands and check-run
annotation references, actions/checkout issue 1036, and bitcoin/bitcoin
issue 33303, where a project met the same stale re-run and chose the dynamic
merge ref, which ADR-0236 rejects because GitHub recomputes that ref lazily.
The one open Backlog Entry, a judge-assigned model tier per Task, shares no
context with these sources, and there is no unresolved Finding.

## Decisions

- Count an exiting member as ended; it has passed the point of no return.
  Keep the start-time identity and zombie rules of Spec 0220.
- Merge the current tip in CI and record it, and let the queue judge a
  failure by the tip it tested. See ADR-0236.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
