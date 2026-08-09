---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An unobserved Verification is unknown, not a verdict

Task Verification answers in two values: the command exited zero, or it did not.
Three failures on 2026-08-09 were all failures of that binary. A gate passed
against an untouched tree, which should have been unknown rather than a pass. A
Batch that failed on one issue overwrote twenty recorded outcomes, treating
"this Batch did not finish" as a refutation of work already done. And `make
verify` returned exit 2 and then exit 0 on one unchanged tree, because a wait
budget expired — an unobserved result reported as a verdict.

Roundfix therefore separates a command's verdict from the runner's ability to
observe it. A command that ran and exited non-zero has a verdict, and that
verdict is failure. A command that timed out, ran partially, or could not be
executed has no verdict, and the outcome carries `unknown` with its cause. The
distinction reaches the Task's terminal reason and its Run Event, so a maintainer
reading the record can tell "the work is wrong" from "we did not find out".

This is deliberately a cause on the existing outcome, not a fourth `spec.Status`.
A Task whose Verification could not be observed still settles `failed` — it is
not complete, and the Run must not proceed as though it were. What changes is the
reason attached to it and the evidence published beside it. Introducing a fourth
Task status would change what a Run's outcome means, which is a larger contract
question and belongs to the Spec that owns Run terminal disposition.

Folding `unknown` into failure without recording the cause was rejected because
it is the state that most needs to be visible: a Task that failed for a real
defect wants a repair turn, while a Task that failed because the machine was
loaded wants a retry, and today both produce the same line. Folding it into
success was never on the table — a gate that reports pass when it did not observe
the surface it names is the defect this Spec exists to remove.
