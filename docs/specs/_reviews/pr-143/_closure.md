# Closure — pull request 143

Pull request #143 was squash-merged into `main` on 2026-08-09 as `1a31c965`,
carrying Specs 0082, 0083, 0084, 0088 and 0089.

Round 001 of this review closes with work outstanding. The maintainer authorised
that closure explicitly on 2026-08-09:

> As 8 issues, 20 em failed do coderabbit e as quatro linhas do qa, se não faz
> parte do novo material considere finalizado e arquive com override autorizado
> por mim.

## What is being closed

| Outcome | Count | What it means |
| --- | ---: | --- |
| `resolved` | 17 | Addressed and recorded during the Run. |
| `failed` | 21 | Twenty carry triage notes describing the work they did; one was triaged `failed` on its merits. |
| `invalid` | 2 | Triaged as not applicable. |
| `pending` | 8 | Never assigned to a Batch. Seven minor or nitpick, one trivial. |

One further issue, 045, was a `major` functional-correctness finding. It was
confirmed against the code, fixed, and covered by a regression before the merge;
it is not part of this closure.

## Why the statuses were not rewritten

The twenty `failed` issues are not failures. `MarkBatchFailed` overwrote them
when one issue in their Batch failed, so finished work is filed as outstanding.
A spot check confirmed one of them landed in the tree.

They stay `failed` anyway. Restoring nineteen statuses on the strength of the
notes those same runs wrote is the class of evidence that proved false in Spec
0089's `task_05`, where a detailed Result described measurements that had never
run. Rewriting them would also destroy the record that the defect happened,
which the Spec that fixes it needs as its measurement.

So this file closes the round; the statuses stay as the accurate history of what
the loop did.

## What happens to the work

Nothing in round 001 is carried forward as a Task. Anything still true will be
raised again by a fresh review round against `main`, where it is re-fingerprinted
against the merged code rather than against a branch that no longer exists.

The defect that produced the twenty is filed at
`docs/backlog/2026-08-09-a-failed-batch-erases-the-issues-it-resolved.md`, with
the measurement and the reason the obvious one-function repair was tried and
reverted.
