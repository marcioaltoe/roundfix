---
type: fix
status: promoted
created: 2026-10-26
spec: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
reason: null
---

# The follow-up mode leaves the handoff and the handed-in messages pending

## Symptom

A session that runs the Roundfix skill's `--follow-up` Skill Mode reads the
project's ai-memory messages and hands them to `roundfix follow-up
--messages`, but consumes them only in the mode's fourth step, after the plan
and the answers. When the mode stops at the first step because there is
nothing to do, or before the fourth because the user stops it or answers
nothing, every message it read stays pending. On 2026-10-26 a session of this
repository found eight pending messages at its start, four of them
informational and from the previous day. The mode never looks at the
project's session handoff either, so the handoff the previous session wrote
stays open and its next steps never reach the plan.

## Where

`.agents/skills/roundfix/references/modes.md` (Spec 0307, ADR-0320,
ADR-0330): its `Optional sources` section says to list the messages "without
consuming them", step 4 of `--follow-up` is the only place that calls
`memory_message_pop`, and no section mentions `memory_handoff_list` or
`memory_handoff_accept`. The `Modes` summary of `.agents/skills/roundfix/SKILL.md`
says the same in one line per mode.

## Expected

The mode accepts the project's most recent open handoff, after reading a
newer cut-off session first, or uses the one the SessionStart hook already
claimed, and turns its next steps and open questions into plan items named in
the report. No message the mode handed in stays pending when the turn ends,
also when the mode stops early: a recorded message becomes a Spec, a Backlog
Entry or a Finding and is then consumed, an informational message is consumed
and cited in the report, and a message not yet handled becomes a Backlog Entry
and is then consumed.

## Evidence

The maintainer's request in chat, 2026-10-26 ~17:30 BRT, verbatim: "Precisamos
corrigir uma ação para o follow-up de roundfix relacionada ao ai-memory, é
necessário consumir o handoff e mensagens e não só ler." The operator's
session of the same day found the eight pending messages at its start. The
mode text at `feeffa28` reads "list this project's messages without consuming
them" and "consume each message the mode handed in (`memory_message_pop`)
only after its content is recorded", in step 4 only.
