### watch

```bash
roundfix watch --source coderabbit --pr <number> --until-clean [--max-rounds N]
```

Waits for Review Source Evidence on the current PR HEAD, observes the
configured quiet period, fetches, resolves Batches, and repeats. Evidence is
always bound to the expected head:

- `pending` means no usable expected-head signal exists. A stale check or
  review remains visible as detail but cannot verify the expected head.
- `reviewing` means a current-head CodeRabbit check or status is pending or in
  progress.
- `reviewed` means CodeRabbit produced a current-head result that does not
  prove Merge-Ready. A successful check, status, or approval stays `reviewed`
  while an unresolved CodeRabbit thread exists; a non-approved review is also
  only `reviewed`.
- `verified` accepts a successful current-head CodeRabbit check or commit
  status, or a current-head CodeRabbit `APPROVED` review, only with zero
  unresolved CodeRabbit threads.
- `skipped` requires an explicit structured CodeRabbit skip for the expected
  head. It ends Review Skipped with exit `3`, prints the Review Source reason
  and next action, fetches no Review Issues, and cannot mean Clean, Clean
  Unverified, or a zero-issue Round.
- `failed` records an explicit current-head Review Source failure.

`WaitingForReview` is the pre-fetch phase.
`WaitingForReviewCheck` is the Merge-Ready phase after Final Push. Each wait
records its expected head, start time, deadline, Evidence state and kind, and
retry status. Non-TTY progress prints on phase entry and when Evidence or retry
state changes; the Live Run View derives remaining time from the deadline.

Roundfix retries only a typed transient Review Source failure: a context
deadline not caused by Run cancellation, a temporary DNS failure, a connection
reset, HTTP `429`, or a GitHub `5xx` response. One episode records `started`,
then `recovered` or `exhausted`. Retry sleeps use the existing poll interval
and remain bounded by the existing Review Source timeout and Run Budget; there
is no new retry setting. The watch loop does not infer retryability from its
Console Log, progress lines, or Run Event summaries.

After Final Push, a proven Daemon-created review-artifact commit can inherit
its parent's verified Evidence. Roundfix requires the recorded commit to be
the current head, to have exactly the recorded verified parent, and to retain
the exact Daemon-generated artifact-commit subject. Its stageable review root
must be inside the repository without a symbolic-link crossing, its diff must
be non-empty, and every changed path must be below that root. Roundfix then
refreshes the parent's Evidence, which must still be verified with no
unresolved CodeRabbit threads. Missing or changed identity, a wrong or multiple
parent, a non-current head, a changed subject, an external or symbolic-link
root, an empty diff, any mixed or out-of-root path, stale parent Evidence, or
an unresolved thread refuses inheritance and returns to ordinary current-head
Evidence polling. User-authored documentation commits and all other
non-Daemon descendants never inherit.

If accepted Evidence never appears within `watch.check_grace_period` (default
`5m`), watch ends Clean Unverified with exit `3` and names the next action.
Other terminal outcomes are `MaxRoundsReached`, `BudgetExceeded`, `TimedOut`,
`Failed`, and `Stopped`.

