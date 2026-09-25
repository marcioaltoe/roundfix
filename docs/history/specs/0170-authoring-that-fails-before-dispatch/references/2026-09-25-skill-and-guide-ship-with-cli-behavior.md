---
type: fix
status: promoted
created: 2026-09-25
spec: 0170-authoring-that-fails-before-dispatch
reason: null
---

# A CLI behavior change lands without the skill or guide that describes it

## Symptom

A Task changes a command's flags, output, or exit codes, and the skill or agent guide that tells agents how to use it stays behind. The pre-PR review then finds the drift and a corrective Task follows. Two of the sixteen Unresolved Runs in Specs 0155–0169 came from this lag.

## Where

Task authoring in `write-tasks` and `spec check`; `.agents/skills/**`, `skills/*/SKILL.md`, and `docs/agents/*.md` against `internal/cli` help text.

## Expected

When a Task's bounded files include a CLI surface, the same Task (or a Task it depends on inside the Spec) also names the skill or guide that documents that surface, and `spec check` warns when neither does.

## Evidence

Efficiency diagnosis in secondbrain `raw/roundfix/2026-09-25-sequencia-de-eficiencia.md` (docs-lag Unresolved Runs).
