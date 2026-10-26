---
status: accepted
created_at: 2026-10-26T21:00:00Z
updated_at: 2026-10-26T21:00:00Z
deprecated_at: null
superseded_by: null
---

# The follow-up mode claims the handoff and leaves no handed-in message pending

ADR-0330 made ai-memory the channel for session handoffs and cross-project
messages, and ADR-0320 gave the Roundfix skill its `--follow-up` Skill Mode.
The mode listed this project's messages without consuming them, handed them to
`roundfix follow-up --messages`, and consumed them only in its fourth step,
after the plan and the answers. It never mentioned handoffs. A session that
stopped at the first step, because nothing was to be done, or before the
fourth, because the user stopped it or answered nothing, left every message it
had read pending, and the handoff the previous session wrote stayed open. On
2026-10-26 a session found eight pending messages at its start, four of them
informational and from the previous day. The maintainer asked the same day for
the mode to consume the handoff and the messages, not only read them, and
chose the two recommended answers quoted in Spec 0325's `_authorization.md`.

## Decision

**The mode claims the project's handoff before its inventory.** It uses the
handoff block the SessionStart hook already claimed when the context holds
one, and claims no other. Otherwise it lists the open handoffs, reads a newer
cut-off session that the listing names before anything else, and accepts the
project's most recent open handoff by its id. Each next step and open question
of that handoff reaches the Follow-up Command as one entry of the existing
`roundfix/follow-up-messages/v1` file, with id `handoff:<handoff id>:<n>` and
this project as its sender, so the Follow-up Priority orders it with every
other item. The plan and the report name each one.

**No message the mode handed in stays pending when the turn ends**, also when
the mode stops early. A message the plan records becomes a Spec, a Backlog
Entry or a Finding and is then consumed. An informational message is consumed
and named in the report. A message not yet handled becomes a Backlog Entry
that records it and is then consumed. Each pop names its message id and
follows the record, so the rule of ADR-0330, consume only after recording,
holds unchanged. `--status` and `--unattended` claim no handoff and consume no
message.

## Considered Options

- Keeping consumption in the fourth step alone was rejected: every early stop
  leaves the inbox as it found it, and the next session reads the same
  messages again.
- Leaving a message nobody handled pending for a later session was rejected
  by the maintainer: a pending message is listed again by every session and
  owned by none, while a Backlog Entry is owned by the repository and ordered
  by the Follow-up Priority.
- Listing the handoff without accepting it was rejected by the maintainer:
  an open handoff stays offered to the next session, which then resumes work
  the follow-up already planned.
- A new handoff field in the messages file, or a handoff source read by the
  command, was rejected: ADR-0330 keeps Roundfix's code away from ai-memory,
  and the existing schema already accepts any id, an empty type and the
  sender's name, so the command and its guide need no change.

## Consequences

- The Roundfix skill's `modes` reference and its `Modes` summary change, and
  the owned-skill version rises; no CLI flag, schema or Baseline clause
  changes.
- The Baseline's message clause already asks to consume only after recording
  in a Spec, a Backlog Entry or a Finding, or after naming an informational
  message, so adopters keep their guides unchanged.
- A Backlog Entry written for an unhandled message is a working-tree file
  until a session with commit authority commits it; the report names it.
