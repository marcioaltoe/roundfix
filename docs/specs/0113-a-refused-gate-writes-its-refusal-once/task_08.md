---
status: pending
type: backend
---

# Task: Give the refusal report a writer in production

QA finding F-1 (`Blocks-Completion`): `spec.WritePreconditionRefusalReport` and
`speccheck.PreconditionRefusal` are complete and correct, and nothing calls
either one. The only non-test reference to the writer in the tree is its own
definition. The Spec's stated goal is not delivered on any path a Supervisor
reaches.

## Work

- Wire the deriver into the Daemon's mechanical request so `Precondition`
  carries the refusing check and its reason instead of its zero value, and the
  report the Daemon writes emits `rows_blocked_precondition`,
  `precondition_check`, and `precondition_reason` when a precondition refused.
  The three pre-existing counts keep their current meaning.
- Give the refusing gate its contract: `.agents/skills/qa-gate/SKILL.md` states
  what a gate writes when it stops at a precondition — one terminal row, its
  frontmatter, and the check and reason that caused the refusal — so a gate that
  refuses leaves an artifact the mechanical stage accepts. This path is the one
  the PRD describes, and it is the reason the Tooling authority row was amended
  on 2026-08-26.
- Regenerate the skill mirror with the declared `make skills-sync`; never edit
  the mirror by hand.
- A gate that did not refuse at a precondition writes exactly what it writes
  today. The refusal path is additive.

## References

- User Story 1: Gate writes valid report
- User Story 2: Refusal recorded and auditable
- Core Feature 1: Terminal Row Writing
- Core Feature 3: Precondition recorded in the report

## Verification
- `grep -rq "WritePreconditionRefusalReport" internal/daemon/ && grep -q "rows_blocked_precondition" internal/daemon/task_engine.go && grep -qi "precondition" .agents/skills/qa-gate/SKILL.md && diff -q .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && go test -count=1 ./internal/daemon ./internal/spec 2>&1 | grep -q "^ok"`

## Result
A gate that refuses at a precondition writes a report its own contract accepts,
on the path a Supervisor actually reaches.
