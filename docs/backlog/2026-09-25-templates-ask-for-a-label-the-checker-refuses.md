---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# The PRD and TechSpec templates ask for `bounded paths:`, which the checker refuses

## Symptom

Template-faithful PRDs fail the first `spec check` with `SC-TOOLING-UNBOUNDED`; each template contradicts its own comment (`bounded files:`).

## Where

`.agents/skills/write-prd/references/prd-template.md`, `write-techspec/references/techspec-template.md` and mirrors; `internal/speccheck/constraints.go`.

## Expected

Templates say `bounded files:`, or the checker accepts `bounded paths:` as a synonym.

## Evidence

secondbrain `inbox/roundfix/2026-09-17-o-template-do-prd-pede-um-rotulo-que-o-checker-nao-aceita.md`.
