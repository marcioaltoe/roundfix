---
type: fix
status: promoted
created: 2026-10-14
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
reason: null
---

# The implement reference states the QA secret refusal in another form

## Symptom

The `roundfix` skill's `implement` reference says a QA Task that carries a
secret settles with `QA evidence carries a secret: <path>:<line> (<kind>)`.
The Daemon prints each Secret Match as `<path>:<line>: <kind>`, and
`TestQAReportWithASecretMatchSettlesFailedWithoutACommit` asserts the colon
form.

## Where

`.agents/skills/roundfix/references/implement.md` (and its `skills/roundfix`
mirror), against `internal/daemon/task_engine.go` and
`internal/secretscan`'s `Match.String`.

## Expected

The reference states the message the Daemon prints, and a docs contract keeps
the two in step.

## Evidence

Found while writing the v0.78.0 release notes for Spec 0267 on 2026-10-14.
