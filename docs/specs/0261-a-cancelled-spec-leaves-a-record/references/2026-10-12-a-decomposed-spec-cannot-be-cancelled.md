---
type: feat
status: promoted
created: 2026-10-12
spec: 0261-a-cancelled-spec-leaves-a-record
---

# A decomposed Spec cannot be cancelled

## Opportunity

A maintainer who drops a Spec that already has a Task Graph can record why,
what it delivered, and hand its adopted sources back.

## Value

On 2026-10-10 the vortex maintainer dropped eight open Specs that already had
Task Graphs. One of them, 0030, had 5 Tasks completed. No command could record
that:

- `roundfix supersede` requires a Spec that was never decomposed and a
  delivering Spec;
- `roundfix archive --qa-override` requires every non-QA Task to be completed;
- Task `status` belongs to the Daemon.

The vortex deleted the folders. No Archive Record remains, and what 0030
delivered has no history.

## Shape

Maintainer decision, 2026-10-11 ("Comando com registro"):
`roundfix cancel <slug> --reason <text>` writes an Archive Record with the
`cancelled` disposition. The record keeps the reason, the date, and the Tasks
that were completed or pending at that moment. The command returns the
adopted sources to the backlog under the adoption contract, and refuses while
a Run or a Delivery Queue item is active for the Spec. The `archive-spec`
skill and the docs guide describe it.

Source:
- Secondbrain `inbox/roundfix/_triaged/2026-10-10-nao-ha-caminho-para-cancelar-uma-spec-ja-decomposta.md`
  (vortex, 2026-10-10).
