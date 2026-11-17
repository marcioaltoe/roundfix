---
type: fix
status: promoted
created: 2026-11-16
spec: 0401-a-supervisor-session-ends-only-when-the-maintainer-asks
reason: null
---

# The Phase Boundary ends the Supervisor session after every Authoring Phase

## Symptom

Spec 0397 (ADR-0407) made each Authoring Phase end the session: the authoring skills' `## Phase boundary`, step 8 of the `--follow-up` Skill Mode and the notice of `roundfix phase start` tell the Supervisor to write the Session Handoff after the Idea, PRD, TechSpec and Tasks phases and continue in a new session. The maintainer said on 2026-11-16: "Esse negocio de criar sessões a cada etapa está atrasando demais."

## Decision

Marcio, 2026-11-16 ~21:00 BRT: "Só por Spec (Recommended)" — the session ends only after the Spec is delivered and released; the Authoring Phases run in agents with their own context. The Phase Record keeps measuring.

## Where

`.agents/skills/*/SKILL.md` `## Phase boundary` of write-idea, write-prd, write-techspec and write-tasks; `references/modes.md` step 8; the notice of `roundfix phase start`; ADR-0407.
