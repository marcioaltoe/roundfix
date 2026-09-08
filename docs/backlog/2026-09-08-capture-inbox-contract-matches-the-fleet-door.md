---
type: fix
status: open
created: 2026-09-08
spec: null
reason: null
---

# The documentation guide names a capture door the fleet no longer uses

## Symptom

The generated layout guide declares `docs/_inbox/` for raw notes while fleet
capture and destination-owned triage use the Secondbrain Inbox. The same
workflow has two declared homes with no rule explaining their relationship.

## Where

`docs/agents/docs-layout.md` and the canonical
`internal/baseline/assets/modules/context-workflow.json` declare the local
directory; `docs/agents/secondbrain.md` declares the fleet door.

## Expected

Make the canonical capture and triage instructions agree on the destination
Inbox. Preserve one disposition per Inbox Entry and its provenance. The precise
wording belongs to provisional group P1, authorization and knowledge lifecycle.

## Evidence

Secondbrain `inbox/roundfix/_triaged/2026-08-27-o-guia-docs-layout-exige-um-inbox-local-que-a-frota-nao-usa-mais.md` records the original
observation. Triage on 2026-09-08 checked the current local sources named above.
The Inbox body remains the observation's provenance; this entry records intent.
