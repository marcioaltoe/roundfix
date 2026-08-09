---
type: fix # feat | fix | perf | refactor
status: open # open | promoted | declined
created: 2026-08-09
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# A failed Batch erases the issues it resolved

## Opportunity

When a Batch's Verification fails, Roundfix overwrites the status of **every**
Review Issue in that Batch to `failed`, including the ones the Agent had already
settled as `resolved`:

```go
func MarkBatchFailed(batch rounds.Batch, terminalReason string) error {
	for _, issue := range batch.Issues {
		if err := rounds.SetIssueStatus(issue.Path, rounds.StatusFailed, "", terminalReason); err != nil {
```

`internal/agent/agent.go:328`. The status answers "did this Batch pass its
gate?" and is written into a field that reads as "was this issue fixed?".

## Value

Measured on 2026-08-09, resolving CodeRabbit's review of pull request #143. The
Run reported `17 resolved, 2 invalid, 0 duplicated, 21 failed, 9 unresolved`
across 49 issues. Batch 001's twenty were all marked `failed` — and all twenty
edits were correct and are now committed. The Batch failed because it had edited
two Roundfix-owned authorial skills without regenerating their distributed
copies, a defect with no relationship to any of the twenty findings.

So one unrelated gate failure erased the record of twenty correct resolutions.
The cost lands twice: a later Round re-fetches issues that were already fixed and
spends an Agent turn rediscovering that, and a human reading the round's issue
files is told the work failed when it did not.

The blanket write also loses information the Run had. The Agent settled each
issue individually — `resolved`, `invalid`, `duplicated` — and that per-issue
judgement is discarded rather than kept alongside the Batch outcome.

## Shape

Non-binding. The distinction the data wants is between an issue's resolution
status and the Batch's gate outcome; today one field carries both and the
coarser one wins. Preserving the Agent's per-issue settlement and recording the
Batch failure separately would keep both facts, and would let a later Round skip
what was already resolved instead of re-deriving it.

Worth settling in the same work: what a re-fetched Round should do with an issue
whose fix is already in the tree, since that is the case this defect creates most
often. And whether a Batch whose Verification fails for a cause outside its own
issues should be distinguishable from one that failed because its edits were
wrong — the two currently look identical in the artifacts.
