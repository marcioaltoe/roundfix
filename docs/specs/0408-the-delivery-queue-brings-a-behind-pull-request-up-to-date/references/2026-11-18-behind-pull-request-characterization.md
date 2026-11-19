# A delivery Pull Request behind a strict base branch: characterization and published evidence

Recorded by the author of Spec 0408 on 2026-11-18, on the planning branch `docs/plan-0408` based on `1a8e08b4`. Every probe ran in a scratch worktree of that branch, outside this repository's tree, and was removed afterwards. Nothing here reached the network or a real forge.

## The defect, reproduced

### At the Delivery Engine

A test built on the existing `newMergeStateFixture` of `internal/delivery/merge_state_test.go` gave the item a Pull Request whose every poll reports `mergeStateStatus: BEHIND` with one passing check, and a forge whose merge is refused twice with `the base branch policy prohibits the merge`. The item parked:

```text
stage=parked blocker="delivery-error: merge pull request: merge pull request: the base branch policy prohibits the merge" attempts=2 polls=2
roundfix: checks: Delivery Queue item example: merge state BEHIND
roundfix: checks: Delivery Queue item example: 1 check(s) completed for candidate-original
roundfix: merge refused by branch policy: Delivery Queue item example; checking again
roundfix: checks: Delivery Queue item example: merge state BEHIND
roundfix: checks: Delivery Queue item example: 1 check(s) completed for candidate-original
```

The checking stage waits only while the merge state is `BLOCKED` or `UNKNOWN` (`internal/delivery/engine.go`, `checkCandidate`), so `BEHIND` counts as settled and the item moves to `merging`. `mergeCandidate` returns a policy refusal to `checking` once per head (`policyRefusals`); the second refusal at the same head is returned as an error, which parks `delivery-error`. Nothing in the queue merges the base branch into a branch that is merely behind: the only merge of the default branch, `ResolveConflict`, runs when GitHub reports `mergeable: CONFLICTING`.

### Through the command engine with a strict local forge

A copy of `TestCommandDeliveryEngineRecordsThePullRequest` (`internal/cli/deliver_record_pull_request_test.go`) used the same local forge script, made strict: after the item branch was cut, a second commit was pushed to `main` of the temporary bare origin; `gh pr view` printed `mergeStateStatus` `BEHIND` while `main` was not an ancestor of the Pull Request head, `CLEAN` otherwise; `gh pr merge` failed with `X Pull request #7 is not mergeable: the base branch policy prohibits the merge.` while behind. On `1a8e08b4` the item parked:

```text
roundfix: merge refused by branch policy: Delivery Queue item 0001-widget-flow; checking again
item stage=parked blocker="delivery-error: merge pull request: merge pull request: X Pull request #7 is not mergeable: the base branch policy prohibits the merge."
```

This is the sequence the Quaestor reported on 2026-11-16 (the adopted Backlog Entry).

## The existing merge of the default branch handles a branch that is only behind

A test called `commandDeliveryWorkflow.ResolveConflict` on a temporary repository whose item branch added `item.txt` while `main`, pushed to a bare origin, added `other.txt`, with no derived-path declaration. The call merged the fetched default branch without conflict and committed it:

```text
res={Head:b11306f0… SourcePaths:[] Regenerated:[]} err=<nil>
b11306f0… 7102291… 8d14136… chore: merge default branch and regenerate derived paths
Roundfix-Delivery: derived-merge
```

The new commit has the candidate as first parent and the default branch head as second. A second call on the merged head returned `Head` empty, without error: the fetched default branch was already in the item branch.

## A prototype of the design

The prototype changed only `internal/delivery/engine.go` and `internal/delivery/park_class.go`. At `checking`, a `BEHIND` merge state called the same resolver, recorded its head as a new candidate and returned the item to `gating`; a fourth `BEHIND` parked `pull-request-behind`; a resolver that found nothing to merge kept the Check Wait polling; a conflict parked `pull-request-conflict: <path>`.

- The strict-forge test above then ended `merged`. The owner log read `merge state BEHIND`, then the base update line, then `merge state CLEAN` and `1 check(s) completed for <new head>`. The merge commit on the origin held `moved.txt`, the file the base branch had gained.
- The whole `./internal/delivery` package passed unchanged.
- The whole `./internal/cli` package passed with the prototype; its only failure line was the suite guard reporting the author's own edits to the scratch tree during the run, not a test.
- A check-timeout while the merge state stayed `BEHIND` parked `checks-timeout: current-head checks not settled within delivery.check_timeout 3s; merge state BEHIND`.

## Published evidence

- GitHub GraphQL reference, "Pull requests", enum `MergeStateStatus` (<https://docs.github.com/en/graphql/reference/pulls>): "`BEHIND`: The head ref is out of date." and "`BLOCKED`: The merge is blocked." The same page documents `PullRequestBranchUpdateMethod` with the values `MERGE` ("Update branch via merge.") and `REBASE`.
- GitHub Docs, "Troubleshooting required status checks" (<https://docs.github.com/en/pull-requests/how-tos/merge-and-close-pull-requests/troubleshooting-required-status-checks>): "If branch protection requires your branch to be up-to-date, merge or rebase the base branch into your branch." and "Required checks must pass on the latest commit SHA. Checks from earlier commits don't satisfy the requirement."
- GitHub GraphQL reference, "Branches" (<https://docs.github.com/en/graphql/reference/branches>): `requiresStrictStatusChecks` — "Are branches required to be up to date before merging."; `strictRequiredStatusChecksPolicy` — "Whether pull requests targeting a matching branch must be tested with the latest code. This setting will not take effect unless at least one status check is enabled."
- multiplex-ai/muggle-ai-works pull request #265, merged 2026-06-09 (<https://github.com/multiplex-ai/muggle-ai-works/pull/265>), a measurement on a real pull request: "GitHub collapses merge state into one value with precedence. `DIRTY` (conflict) and `BLOCKED` (missing required review, pending/failing required check) outrank and mask `BEHIND`." and "`BEHIND` surfaces at all only when the base enforces \"require branches up to date\"."

## Secondbrain

- `wiki/sources/digest-quaestor-ruleset-recusa-github-actions-no-bypass-2026-11-16.md` records the Quaestor's Spec 0115: its `main` requires an up-to-date branch and `make verify` runs only on `pull_request`. A push to the item branch is a `pull_request` `synchronize` event there, so the Quaestor's CI runs again on the updated head.
- `wiki/concepts/tags-protegidas-e-release-por-workflow-no-github.md` states the cost of a strict required check: the branch must be up to date with the base before the merge, which costs more builds.
- `qmd query` failed on 2026-11-18 with `SQLITE_BUSY: database is locked`; the two pages above were found by a text search of `wiki/` instead.
