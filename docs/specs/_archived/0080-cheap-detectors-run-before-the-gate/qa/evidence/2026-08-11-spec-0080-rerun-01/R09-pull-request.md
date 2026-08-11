# R09 — Pull Request journeys

The Roundfix QA prompt states `Pull Request: none open; Pull Request journeys
are environment-blocked.` This is a proved absence, not an unresolved lookup.
Per the QA contract, this gate did not try to resolve the per-Run worktree
branch as a Pull Request.

`rtk env GOCACHE=/private/tmp/roundfix-spec0080-rerun-gocache go test
./internal/agent -run
'^TestBuildQAPrompt(StatesPullRequestJourneysAreEnvironmentBlockedWhenNoneIsOpen|SeparatesUnresolvedPullRequestFromProvenAbsence)$'
-count=1 -v` exited 0. The current builder preserves the distinction between a
confirmed absence and an unresolved lookup and gives each the required typed
environment-block wording.

Equivalent repository evidence: `_tasks.md` keeps `task_08` terminal, and
`docs/agents/autonomous-work.md` requires the order implement graph including
QA, archive, open Pull Request, watch until Clean, merge. The no-PR state is
therefore the intended pre-PR gate boundary, not a missing lookup of the Run
Worktree branch.

Approval, checks/status, unresolved review threads, Merge-Ready, and
review-artifact ancestry remain unobservable until a passing QA gate permits a
Pull Request. Unblock by rerunning those read-only journeys against the future
open Pull Request on `ma/cheap-detectors-run-before-the-gate`.
