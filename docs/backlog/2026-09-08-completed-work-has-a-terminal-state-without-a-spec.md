---
type: fix
status: open
created: 2026-09-08
spec: null
reason: null
---

# Completed evidence and intent cannot truthfully close without a Spec

## Symptom

A directly completed or externally routed Finding has no matching archive
license, and a directly executed Backlog Entry has no terminal execution state.
Using `declined` for fulfilled intent would misdescribe what happened.

## Where

`docs/agents/docs-layout.md` and canonical
`internal/baseline/assets/modules/context-workflow.json` define the states.
`internal/speccheck/citations.go` resolves archive licenses only through an
active Rollup or a Spec. The Pantheon Inbox capture supplies concrete examples.

## Expected

Provide a truthful terminal path with explicit disposition evidence for work
completed or routed without a Spec. Preserve existing licenses and history
resolution; do not require a synthetic one-member Rollup. Route to provisional P1.

## Evidence

Secondbrain `inbox/roundfix/_triaged/2026-09-01-evidencia-e-intencao-resolvidas-fora-de-spec-nao-tem-estado-terminal.md` records the original
observation. Triage on 2026-09-08 checked the current local sources named above.
The Inbox body remains the observation's provenance; this entry records intent.
