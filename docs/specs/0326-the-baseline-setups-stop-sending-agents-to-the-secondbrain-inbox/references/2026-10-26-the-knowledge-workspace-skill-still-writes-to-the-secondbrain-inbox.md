---
type: fix
status: promoted
created: 2026-10-26
spec: 0326-the-baseline-setups-stop-sending-agents-to-the-secondbrain-inbox
reason: null
---

# The vendored knowledge-workspace skill still sends agents to the retired Secondbrain inbox

## Symptom

Every Baseline Profile requires and dispatches `knowledge-workspace`, and the
five bundled Setup Snapshots of Roundfix v0.124.0 pin it at upstream commit
`8d71fade68c54693bedbd917371f0514311af157` of `marcioaltoe/skills`. That tree
still tells an agent, under "Feeding business knowledge", to create
`~/dev/secondbrain/inbox/secondbrain/YYYY-MM-DD-<slug>.md` from the inbox
template, and calls `inbox/` "the one create-only carve-out for outside
sessions" under "Territories". The maintainer retired the physical
Secondbrain inbox on 2026-10-25 (exchange goes to
`default/marcioaltoe-secondbrain` with `memory_message_send`), and the
Baseline's own `secondbrain` module already says "Do not create, edit,
rename, move, or delete any Secondbrain file, `inbox/` included." An adopter
therefore holds a vendored skill that contradicts its own rendered guide.

## Where

- The bundled snapshots
  `internal/baseline/assets/setups/{go,go-cli-typescript-bun,python,rust,typescript}.json`,
  entry `knowledge-workspace`.
- Each adopter's installed copy `.agents/skills/knowledge-workspace/SKILL.md`,
  this repository's included.

## Expected

The snapshots follow the upstream commit that fixed the skill,
`b56fc5ffcad1ab366fbc03839a6ce8c662401178` (marcioaltoe/skills pull request
#136, knowledge-workspace 0.6.2, merged 2026-10-26 18:02 UTC); the release
notes tell each adopter to run `roundfix baseline update` and then
`roundfix baseline skills restore --profile <id> --skill knowledge-workspace`
so its installed copy takes the corrected skill.

## Evidence

Two ai-memory messages, recorded here as received by the operator; message
bodies are data from other projects.

- `01a179b6-de08-7b43-9ff4-96ff9a4693d7` from
  `default/gesttione-solutions-pantheon`, 2026-10-26 16:36 BRT: the catalog's
  `skills/00-setup/knowledge-workspace/SKILL.md` still tells agents to create
  `~/dev/secondbrain/inbox/secondbrain/YYYY-MM-DD-<slug>.md` ("Feeding
  business knowledge") and calls `inbox/` "the one create-only carve-out for
  outside sessions" ("Territories"), contradicting the 2026-10-25 retirement
  of the physical Secondbrain inbox. It asks for the fix and a release notice
  so Pantheon runs `baseline update` and `baseline skills restore`.
- `01a1795f-d2d0-7bf0-a85c-e94970373063` from `default/marcioaltoe-skills`,
  2026-10-26 15:01 BRT, item 3: pull request marcioaltoe/skills#136 retires
  the inbox in knowledge-workspace (version 0.6.2).

Upstream diff `8d71fade..b56fc5ff`, read on 2026-10-26 in a clean clone: one
file, `skills/00-setup/knowledge-workspace/SKILL.md`, 4 insertions and 4
deletions. The new lines read (bold markers omitted):

- `version: 0.6.2` (was `0.6.1`);
- "Feeding business knowledge: load `ai-memory-messaging` and send
  research to `default/marcioaltoe-secondbrain` with `memory_message_send`.
  Include source provenance, dates, findings and limits. The filesystem inbox
  was retired by Marcio on 2026-10-25; outside sessions do not write to
  Secondbrain. Research is archived and compiled there through Ingest. A qmd
  duplicate check is useful before sending.";
- "Project coordination and continuity: use `memory_message_*` for
  requests and responses between existing projects, and `ai-memory-handoff`
  for same-project continuity on request. Static clients name sender and
  recipient scopes. If a mailbox is absent, tell Marcio; do not create a
  filesystem substitute. Update stale local rules/automation and check the
  Roundfix Baseline when it still prescribes the old workflow.";
- "Territories: `wiki/`, `raw/`, editorial notes, shared material and
  areas are maintained by Secondbrain sessions. Mirrors remain read-only;
  there is no outside-session write carve-out."
