---
status: pending
created_at: 2026-08-14
updated_at: 2026-08-14
---

# A Spec that coins a term cannot pass its own gate

A Spec that declares a Vocabulary Contract promises that its coined term reaches
the glossary. The task-authoring contract puts that promise in the closing node:
the QA gate checks whether the work introduced, changed, or retired a term and
updates the domain context when it found something. The QA gate also refuses
before doing anything else when `spec check --strict` reports a finding — and
`SC-VOCABULARY-UNDOCUMENTED` is exactly the finding an undocumented coined term
produces.

So the gate is assigned the repair and forbidden from reaching it. Spec 0103
measured it on 2026-08-14: eight implementation Tasks settled `completed`, and
`task_09` failed on

```text
[error] SC-VOCABULARY-UNDOCUMENTED: internal/cli/doctor.go emits undocumented
token "residue" absent from CONTEXT.md
```

The only exits are for a human to write the glossary entry before the gate runs,
which makes the gate's own requirement dead text, or for the Spec to not declare
its vocabulary, which is worse.

Worth settling when this is picked up: whether the glossary update belongs to the
Task that coins the term rather than to the gate — the gate would then verify
rather than perform it, which is what every other row does; or whether the strict
precondition should exempt the vocabulary code when the Spec being gated is the
one that declared the term, so the gate can still do the job it was given.

## The same Run reproduced Spec 0113's defect

The refused gate wrote `qa-report-2026-08-14.md` with `verdict: fail` and an
empty Results table:

```markdown
## Results

| # | Status | Provenance |
| - | --- | --- |
```

That is the malformed shape Spec 0113 exists to repair, reproduced on demand
rather than recalled from Spec 0078. Spec 0095's delivery had produced a report
with a terminal row instead, which had suggested the defect was narrower than
0113 claims; this run shows the empty-table path is still live and is reached
whenever the gate stops at the static precondition rather than at a mechanical
finding.

## And a corrective Task after a failed gate deadlocks the Spec

Adding `task_10` to Spec 0103 after `task_09` had already run made
`make verify-docs` fail across three corpus tests with

```text
validate qa gate: QA gate result is invalidated for Task "task_09"
because these dependencies are not completed: task_10
```

The invalidation is correct in substance — a gate that ran before a new
dependency existed did not judge it. What is wrong is that it is reported as a
Spec-consistency error rather than as a pending state, so the repository's
documentation gate is red for as long as the graph has work left in it. Before
`task_09` ever ran, the same graph with every node pending passed the same gate.

The distinguishing fact is the gate node's own status: `pending` is fine,
`failed` with an incomplete dependency is an error. A graph that gained a
corrective Task is the ordinary case this should describe, not refuse.

It is worse than a red gate. `roundfix implement` refuses the same way, in
Preflight, before creating a Run:

```text
Preflight failed
Reason:
  validate qa gate: QA gate result is invalidated for Task "task_09"
  because these dependencies are not completed: task_10
No side effects:
  Roundfix did not create a Run, fetch Review Source issues, start an Agent,
  commit, or push.
```

So the tool refuses to execute the graph whose execution is the only thing that
would clear the refusal. `roundfix settle` does not help: it re-runs one failed
Task's own Verification, and the failed Task here is the gate, whose Verification
cannot pass until the corrective Task it does not yet judge has run.

There is no supported way out. Spec 0103 was unblocked on 2026-08-14 by editing
`task_09`'s `status` from `failed` back to `pending` by hand — a field whose own
frontmatter comment reserves it to the Implement Command. That edit claims no
work; it withdraws a claim, which is why it is defensible and why it should not
have been necessary. Whatever repairs this should offer the withdrawal as an act
the tool owns.

Worth settling when this is picked up: whether a gate whose recorded status is
`failed` should invalidate anything at all, since a failed gate has no result to
protect; and whether adding a node to a graph should reset a gate that already
ran, rather than refusing the graph that contains it.
