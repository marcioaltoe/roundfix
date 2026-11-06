# Fluxus pull requests #466, #469 and #477: when each check appeared

Measured on 2026-11-06 by the operator for Spec 0359, read only, with `gh pr
view <n> -R gesttione-solutions/fluxus --json mergedAt,headRefOid,createdAt,statusCheckRollup`,
`gh api repos/gesttione-solutions/fluxus/actions/runs?head_sha=<head>` and
`gh api repos/gesttione-solutions/fluxus/commits/<head>/check-suites`. Times
are GitHub's, in UTC. Nothing in Fluxus was changed, and no Fluxus file is
copied here.

## Pull request #466 (Spec 0149), head `a960a60d`

- 12:46:40 pull request created.
- 12:46:45 a check suite of the `cursor` app is created; it stays `queued`
  with zero check runs and never completes.
- 12:46:48 three workflow runs are created together, all with event
  `pull_request`: `CI | Conventions`, `CI | Frontend`, `CI | Backend`.
- 12:46:53 the check run `Validate PR title` (`CI | Conventions`) starts;
  it ends `success` at 12:47:01.
- 12:47:11 the Delivery Queue merges the pull request.
- 12:48:18 and 12:48:20 the check runs `Lint, Type Check, Test & Build` of
  `CI | Frontend` and `CI | Backend` start, 67 and 69 seconds after the merge.
- 12:53:53 `CI | Backend` ends `failure`.

## Pull request #469 (Spec 0151), head `0e7ee646`

- 13:14:19 created; 13:14:27 the same three workflow runs are created
  together; 13:14:43 `Validate PR title` ends `success`; 13:14:49 merged;
  13:15:57 and 13:16:01 the two `Lint, Type Check, Test & Build` check runs
  start; 13:21:34 `CI | Backend` ends `failure`.

## Pull request #477 (Spec 0150), head `bb13eb8b`

- 15:28:52 created; 15:29:00 the same three workflow runs are created
  together; 15:29:21 `Validate PR title` ends `success`; 15:29:40 merged;
  15:30:24 and 15:30:37 the two `Lint, Type Check, Test & Build` check runs
  start; both end `success`, at 15:33:54 and 15:35:26.

## What the timing shows

In all three, every workflow run the pull request event triggered existed
from about eight seconds after the pull request was created, in one batch,
but the check run of a job appeared only when the job started, more than a
minute later. A reader of check runs alone (`gh pr checks`) saw one passing
check and nothing pending; a reader of the workflow runs of the head saw two
runs still in progress. The `cursor` suite shows that a check suite with no
check run can stay queued forever, so check suites are not a completion
signal.

The workflow triggers, read in the Fluxus checkout without copying it:
`CI | Conventions` runs on `pull_request` to `main` with the types `opened`,
`reopened`, `edited` and `synchronize` and no path filter; `CI | Backend` and
`CI | Frontend` run on `pull_request` to `main` with a `paths` filter, and
`CI | Backend` cancels an earlier run of the same ref.
