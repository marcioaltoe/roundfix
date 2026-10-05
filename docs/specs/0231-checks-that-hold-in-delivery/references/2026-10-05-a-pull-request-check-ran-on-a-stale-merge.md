---
type: fix
status: promoted
created: 2026-10-05
spec: 0231-checks-that-hold-in-delivery
reason: null
---

# A Pull Request check ran on a merge with an old default branch, and the queue parked it

## Symptom

On 2026-10-05 Spec 0228's Pull Request #391 failed CI's docscontract corpus
golden because of a defect on main, not in the item. The Delivery Queue
re-ran the failed job once and parked the item `flaky-check` (intervention log
entry 164). The fix merged to main as `8f71887f` at 12:17:00Z. Re-running the
failed run again, and then reopening the Pull Request, still tested the old
merge; the queue parked the item `checks-failed`, and only a second reopen
produced a passing run (entry 165).

## Where

- CI: `.github/workflows/ci-verify.yml` checks out the `pull_request` event's
  merge commit.
- Delivery Queue: `checkCandidate` in `internal/delivery/engine.go` and
  `InspectFailedCheck` / `RerunFailedCheck` in `internal/delivery/check_rerun.go`
  classify a failed check without asking which default branch it tested.

## Expected

A failed check that tested an older default branch is never parked as
`flaky-check` or `checks-failed`; the queue re-tests the item against the
current default branch and keeps polling, and a re-run of CI tests the current
default branch.

## Evidence

Read from GitHub with `gh` on 2026-10-05:

- Run 37306411237, attempts 1 to 3 (11:58, 12:02, 12:17:05Z): every attempt
  checked out `32e025ce`, "Merge 25c00261 into 7eabb720". Attempt 3 started
  five seconds after main moved to `8f71887f`.
- Run 37308972464 (first reopen, 12:21:29Z): again `32e025ce` into
  `7eabb720`, four and a half minutes after main moved.
- Run 37309659532 (second reopen, 12:27Z): `67fdede8`, "Merge 25c00261 into
  8f71887f"; it passed, and #391 merged at 12:34:47Z.
- GitHub's re-run documentation: a re-run uses "the same `GITHUB_SHA` (commit
  SHA) and `GITHUB_REF` (git ref) of the original event".
