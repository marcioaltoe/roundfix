---
status: accepted
created_at: 2026-08-09T00:00:00Z
updated_at: 2026-08-09T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A Task gate is run against the unchanged tree before the Agent

Task 05 of Spec 0089 settled `completed` without doing its work. Its authored
gate was `grep -q 'reasoning_effort: xhigh' .roundfixrc.yml`, and the file
already carried that string on an unrelated `claude`/`sonnet` fallback, from
before the Spec began. The command exited zero against a tree nobody had
touched. The Agent changed only its own Task file, the Daemon ran the gate, the
gate passed, and the configuration the Spec existed to write was never written.

Roundfix therefore runs a Task's `## Verification` commands once against the
unchanged tree, immediately before the Agent Session owner is created, and
refuses the Task when any command already passes. A command that exits zero
before the work happened cannot be evidence that the work happened. The refusal
names the offending command and costs one execution of a Task-scoped command,
against the Agent turn it replaces.

The probe executes the commands rather than inspecting their text. A static rule
about `grep` shapes, needle length, or negation would have cleared Task 05's
command, which is well-formed and would be correct against a file that did not
already contain its needle. The defect was never in the command; it was in the
relationship between the command and the tree it ran against, and only running
it reveals that.

The seam is a sibling of one that already exists. `verifyRepositoryPrecondition`
runs a Verification before the Agent and settles the Task `failed` when the
repository is not green on entry. The probe sits beside it, on the same shape and
the same event channel, so no new execution path is introduced.

Refusing at dispatch rather than at authoring was accepted as the cheaper half of
a two-part answer. Authoring is where the refusal would be cheapest to act on,
but the `write-tasks` contract is protected tooling and its edit needs its own
authorization; the Daemon can enforce the rule today without one. A Task whose
Verification legitimately holds before and after is refused by this decision,
which is intended: the authoring contract already forbids a Verification that
does not prove its own effect, and this makes that rule executable instead of
advisory.
