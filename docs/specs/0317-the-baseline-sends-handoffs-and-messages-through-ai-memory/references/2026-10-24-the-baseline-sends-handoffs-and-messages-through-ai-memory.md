---
type: refactor
status: promoted
created: 2026-10-24
spec: 0317-the-baseline-sends-handoffs-and-messages-through-ai-memory
reason: null
---

# The Baseline sends handoffs and messages through ai-memory

## Opportunity

On 2026-10-24 at 10:26 BRT the maintainer's personal assistant project sent
`default/marcioaltoe-roundfix` an ai-memory message (`01a16e17`) with the
maintainer's request for the Baseline:

- remove the `gws-*`, `handoff` and `brand-guidelines` skills from the
  Baseline and from the setups it generates;
- move session handoff and cross-session and cross-project messaging in the
  Baseline to ai-memory: the `memory_handoff_*` and `memory_message_*` tools
  and the `ai-memory-handoff` and `ai-memory-messaging` skills, in place of
  the Secondbrain inbox;
- keep no Anthropic brand in the Baseline;
- agree skill names and versions with the skills project, work on a
  `refactor/` branch, and answer with a message titled "Retorno: baseline
  Roundfix".

The same morning the maintainer added, in his own words: "Lembrando que
solicitei em algum momento para matar o inbox de secondbrain e usar somente
as mensagens do ai-memory".

## Evidence

- Upstream `marcioaltoe/skills` pull request #131, commit
  `65d5f365cace27c4b468db1723213723b10c870b`, replaced `handoff` with
  `ai-memory-durable-pages`, `ai-memory-handoff`,
  `ai-memory-learning-maintenance`, `ai-memory-messaging`,
  `ai-memory-retrieval` and `ai-memory-routing-install` in all six presets,
  replaced `gws-*` with `gog` in the content preset and removed
  `brand-guidelines`.
- Roundfix at `b5e45c63`: the `core` Baseline module requires and dispatches
  `handoff`; the four upstream Setup Snapshots name commit `b3c45a45`; the
  `secondbrain` module renders the Inbox Entry contract, the capture triggers
  and the `inbox/**` write permission into `docs/agents/secondbrain.md`;
  `roundfix follow-up` reads a `secondbrain-inbox` source from
  `secondbrain.root`. No Roundfix snapshot holds a `gws-*` or
  `brand-guidelines` skill.

## Expected outcome

1. Every upstream Setup Snapshot follows `65d5f36`, with the six ai-memory
   skills and no `handoff`.
2. The Baseline's core guidance sends handoffs and cross-project messages
   through ai-memory when the session exposes it, and says what is missing
   when it does not.
3. The Secondbrain guidance keeps the knowledge base read-only and drops the
   inbox.
4. `roundfix follow-up` stops reading the Secondbrain inbox, and a User
   Config that still names its keys keeps loading.
