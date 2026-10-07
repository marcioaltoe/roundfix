---
type: fix
status: declined
created: 2026-10-05
spec: null
reason: "The Jev Router was retired on 2026-10-05 (ADR-0235), so no routed prompt writes to the Judge Log any more; the 2026-10-05 relay replay measured a gap of about 1 %, and the judge's own spend is cents."
---

# The Judge Log under-counts the Jev Router's spend

On 2026-10-04 the two routed chore replays (0200 task_04 and 0195 task_06) raised the OpenRouter key's usage by US$10.32, and the account's usage by the same amount. The Judge Log recorded only US$7.76. The gate reads usage once, right after the prompt, but OpenRouter kept reporting more for about a minute after each prompt (about 25% more on these two). The month's Jev ceiling (ADR-0231) is compared with the Judge Log sum, so it lets the real spend exceed the configured value by that lag. Evidence: `docs/history/specs/0218-the-jev-router-on-docs-and-chore-tasks/measurement/jev-router-2026-10-04.md`, section "Chore replays, 2026-10-04 evening".

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-10-05-the-judge-log-under-counts-router-spend.md`.
