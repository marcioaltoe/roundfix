---
type: fix
status: open
created: 2026-09-25
spec: null
reason: null
---

# The mandatory docs-layout rule still describes hand-stamping a QA Archive Override

## Symptom

`docs/agents/docs-layout.md` tells agents to "stamp `qa_override: true`" and
record the approval themselves. An agent following this mandatory guide can edit
`_prd.md` by hand and skip the command's refusals and provenance, contradicting
the archive-spec skill, which now says an override is performed only through
`roundfix archive <slug> --qa-override --approval <source> --reason <text>`.

## Where

The Baseline module `internal/baseline/assets/modules/spec-workflow.json`, its
generated `docs/agents/docs-layout.md` and the formatter golden fixture.

## Expected

The module says a runtime that supports overrides must use its override command
and never hand-edit the stamp, keeping the "unsupported runtime is reported as
unsupported" clause; guides and fixtures regenerated.

## Evidence

Recorded limit of Spec 0169 (second pre-PR review of 2026-09-25).
