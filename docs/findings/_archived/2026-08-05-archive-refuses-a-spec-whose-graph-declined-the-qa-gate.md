---
status: done
created_at: 2026-08-05
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-qa-gates-and-verification-evidence.md
---

# 2026-08-05 — Archive refuses a Spec whose graph declined the QA gate

## What was observed

Spec 0016 declared `qa: declined` in its `_tasks.md` frontmatter, with the
reasoned `qa_reason` the contract requires. Its Run reached Clean with all four
Tasks completed. `roundfix archive` then refused it:

```
Preflight failed
Reason:
  no passing QA verdict: no QA Report found for Spec directory
  ".../docs/specs/0016-reconciliation-grain-clarity"; run the qa-gate workflow
  to produce one
```

There is no QA Report because the graph declined the gate, which is the
supported shape. The instruction to "run the qa-gate workflow" contradicts the
declaration the Spec authored and the Implement Command honoured — that Run
executed no gate, exactly as declined.

The Spec was archived by hand instead, stamping `qa_override: true` because
that is the field the verifier reads, with a comment recording that the gate
was declined by design rather than failed.

## Root cause

Two commands disagree about what a complete Spec looks like.

The Implement Command implements the authored-QA-decision contract fully: it
accepts `qa: task_NN` or `qa: declined` plus `qa_reason`, refuses a graph with
neither or both, and runs no gate for a decline.

The Archive Command predates that contract and still requires a passing QA
Report unconditionally. It never reads the manifest's declaration, so the only
way past it is the override intended for *failed or missing* verification —
which then records the wrong thing. An auditor reading `qa_override: true`
sees a Spec archived despite bad QA evidence, when in fact the Spec correctly
declared it needed none.

## What would settle it

Teach Archive the same declaration Implement already reads. When the manifest
declares `qa: declined` with a non-empty `qa_reason`, that is the QA precondition
satisfied — no report required, no override recorded, and the reason carried
into the archive stamp so the decision travels with the Spec.

Keep the current refusal for every other shape: an included gate still needs
its passing report, and a graph declaring neither form is still a defect.

Worth checking whether other post-contract commands share the assumption. The
QA decision became manifest state; anything still inferring it from the
presence of a report on disk will disagree with the graph the same way.

## Related

[[2026-08-04-pre-contract-spec-graphs-run-with-no-qa-gate-and-say-nothing]]
covers the opposite corner of the same contract: a graph that predates the
declaration. This one is a graph that makes the declaration correctly and is
refused for it.

## Spec pointer

None yet.
