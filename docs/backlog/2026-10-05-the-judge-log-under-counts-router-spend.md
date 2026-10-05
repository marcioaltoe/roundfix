---
type: fix
status: open
created: 2026-10-05
spec: null
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
