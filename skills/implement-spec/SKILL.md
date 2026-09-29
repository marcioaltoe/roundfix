---
name: implement-spec
description: Prepare a spec with Roundfix's Delivery Plan, delegate implementation to Roundfix, monitor the Run or delivery queue, and ask the maintainer only its Pending Question.
disable-model-invocation: true
argument-hint: "<spec slug or path under docs/specs/> [--from task_NN]"
metadata:
  category: implementation
  tags: [workflow, agents, coding, issues]
  version: 0.1.0
  author: Marcio Altoé
  source: https://github.com/marcioaltoe/skills
version: 0.1.0
---

# Implement Spec

This entry point prepares and delegates a Spec through Roundfix. The Supervisor
never writes code or tests, never runs a Task or the **implement-task** cycle
itself, and never runs the **qa-gate** itself. The QA gate is the Spec's
terminal Task, which the Daemon runs.

## 1. Prepare

`$ARGUMENTS` names the Spec (slug or path under `docs/specs/`). Run
`roundfix deliver plan <slug>...` before choosing an implementation route. If
the plan reports a blocked Spec, stop and report its reasons. Do not start a
blocked Spec or ask the maintainer to approve an ordinary implementation step.

The Roundfix skill is the source of truth for Delivery Plan output, blockers,
and command details.

## 2. Delegate

For a single Spec on the current branch, hand implementation to
`roundfix implement --spec <slug>`. Add `--detach` when the session may end
before the Run does.

For a merge-through sequence, hand the queue to
`roundfix deliver start [--max-duration <duration>] [--max-retries <n>] <slug>...`.
Do not plan waves, execute Tasks, or perform implementation or Verification in
this skill; Roundfix and its Daemon own that work.

The Roundfix skill is the source of truth for the implement and delivery
commands, options, lifecycle, and recovery behavior.

## 3. Monitor and ask

Monitor a delivery queue with `roundfix deliver status` or monitor a Run
through its events. Report the observed Run or queue outcome and follow the
Roundfix skill's recovery path when it provides one.

When the queue presents a **Pending Question**, ask the maintainer only that
question. Preserve its wording and named action; do not invent a confirmation,
approval, recommendation, or second question. Continue only after the
maintainer answers through the explicit action the question names.

## Guardrails

- Do not edit `_prd.md`, `_techspec.md`, or `_tasks.md`; Spec amendments are a
  maintainer decision outside this handoff.
- Do not run destructive Git commands without explicit permission.
- Keep command semantics and delivery lifecycle details in the Roundfix skill;
  do not copy them into this entry point.
