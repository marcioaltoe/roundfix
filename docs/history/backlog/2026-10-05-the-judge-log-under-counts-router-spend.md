---
type: fix
status: declined
created: 2026-10-05
spec: null
reason: "The Jev Router was retired on 2026-10-05 (ADR-0235), so no routed prompt writes to the Judge Log any more; the 2026-10-05 relay replay measured a gap of about 1 %, and the judge's own spend is cents."
---

# The Judge Log under-counts the Jev Router's spend

## Problem

On 2026-10-04 the two routed chore replays (0200 task_04 and 0195 task_06)
raised the OpenRouter key's usage by US$10.32, and the account's usage by the
same amount. The Judge Log recorded only US$7.76. The gate reads usage once,
right after the prompt, but OpenRouter kept reporting more for about a minute
after each prompt (about 25% more on these two). The month's Jev ceiling
(ADR-0231) is compared with the Judge Log sum, so it lets the real spend
exceed the configured value by that lag. Evidence:
`docs/history/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router-2026-10-04.md`,
section "Chore replays, 2026-10-04 evening".

## Expected

The spend the ceiling compares against matches what OpenRouter bills for the
month, within a stated tolerance. Two ways to get there: read the key's
`usage_monthly` (or the account usage) as the authority, or settle each
routed prompt's Judge Log line after usage stops changing.

## Notes

Non-binding. Spec 0229's relay sees each response's generation id, and
OpenRouter's generation lookup returns the final cost of that id, so the
relay could record a settled cost per response.

## Disposition — 2026-10-05

Declined. The maintainer retired the Jev Router on 2026-10-05 ("Aposentar o
router", then "Vamos seguir com a remoção"), recorded in ADR-0235 and carried
out by Spec 0230-retire-the-jev-router. The lag this entry measured belonged
to routed prompts only. The relay replay of 2026-10-05 (the "Relay
confirmation, 2026-10-05" section of the same measurement record) found the
Judge Log US$0.1173 below the key's usage change, about 1 %, and the judge's
own 456 calls from 2026-10-01 to 2026-10-05 cost US$0.014, so the judge's
Judge Log sum needs no settled cost.
