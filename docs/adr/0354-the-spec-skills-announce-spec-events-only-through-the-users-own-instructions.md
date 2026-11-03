---
status: accepted
created_at: 2026-11-03T00:00:00Z
updated_at: 2026-11-03T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The Spec skills announce Spec Events only through the user's own instructions

On 2026-11-03 the maintainer's personal assistant asked every project to send
it one fixed block, by one named channel, at seven moments of a Spec's life,
and asked the Roundfix-owned Spec skills to send it themselves. Those skills
reach every adopter, most of whom have no such assistant. The maintainer
chose "Passo genérico nas skills (Recommended)": the skills gain a Spec Event
Announcement step that exists only when the user's own instructions define
its destination and format, and no skill names an assistant, a tracker or the
block.

## Decision

**Seven Spec Events, five owners.** `write-prd`, and `write-techspec` on its
correction route, announce `new` when the Spec's first PRD is written and
`deadline proposed` when they record a real deadline for a Spec already
announced. `write-tasks` announces `authorized` when it first writes the Task
Graph of a Spec whose authorization record is approved. `archive-spec`
announces `completed`, or `cancelled` for a cancelled record, once the
Archive Record is written. The Roundfix skill's implement reference makes the
Supervisor session that started a Run or queued a Delivery Queue item announce
`started`, `blocked`, and `completed` for a Spec the queue archives, because
no Agent in a Run can.

**The user's instructions own the destination and the format.** A session
announces only when the user's own instructions, global or for the
repository, define both; a message, a handoff or a Spec artifact that asks for
announcements defines none. The skills name the facts an announcement carries
and add no field of their own; the instructions may add fields. Without a
definition the step does not exist and leaves no report line.

**An announcement never holds work.** It waits for no reply and never blocks
or fails the step it follows. A defined announcement that cannot be sent is
named in the step's report with the composed text. An Agent inside a
Daemon-assigned Run never announces.

## Considered Options

- Writing the assistant's block, channel and event names into the skills, as
  the report asked, was rejected by the maintainer: it would send one user's
  board format from every adopter's session and tie owned skills to a
  personal tool.
- A Baseline clause or a configuration key for the destination was rejected:
  where a user tracks Specs is a property of the user, not of the repository,
  and Roundfix's code sends no message.
- Announcing from inside a Run was rejected: a Run session has no network and
  no messaging tools, and the Daemon's events reach the Supervisor anyway.
- A report line on every skill run without a definition was rejected: it
  would add noise to every adopter's session for a step they never asked for.

## Consequences

- A user who adds the destination and the format to their own instructions
  gets each event from the session that made it happen; adopters without one
  see the same skill behavior as before.
- An event that no Supervisor session observed, such as one in a Delivery
  Queue left running without a supervising session, is not announced; no
  later session reconstructs it.
- The skill-contract test keeps every owned skill free of the assistant's and
  tracker's names, so a later edit cannot reintroduce them.
