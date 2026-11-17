---
status: accepted
created_at: 2026-11-17T22:30:00Z
updated_at: 2026-11-17T22:30:00Z
deprecated_at: null
superseded_by: null
---

# Roundfix reads the limit from the error a session passes, and prints the model an Author Subagent starts on

Two subcommands of `roundfix phase` carry the Authoring Fallback of ADR-0427
and ADR-0428. `roundfix phase fallback` takes the error an Author Subagent
ended on, decides from its text and the Authoring Model Declaration whether
the phase falls back, records the fallback and prints the model and effort to
start on. `roundfix phase model` prints the model and effort a phase's Author
Subagent starts on, from the declaration and the Spec's recorded fallbacks.
The authoring skills run them and pass what they print. A refusal exits `3`.

## Context

The session sees the error; Roundfix does not. The rules ADR-0427 sets
compare the error with the declared model, and ADR-0428 asks whether an
earlier phase of the Spec fell back. Both are decisions a test can hold.
Written into four skills as prose, they would be decided by each session
again. ADR-0418 keeps `roundfix phase start` free of configuration, and
ADR-0414 keeps it from reading earlier entries.

Claude Code writes a usage limit as `You've hit your <name> limit · resets
<time>`. The name is `session` or `weekly` for the limits every model shares,
and the family name, such as `Opus` or `Fable`, for a model's own limit. Spend
limits read `You've hit your monthly spend limit` and similar, with more than
one word before `limit`. A Fable consent message begins `Fable limit reached`.

## Decision

**Two subcommands, no change to `phase start`.**
- `roundfix phase fallback <phase> --spec <slug> --error <message>` loads
  Project Config and records a fallback entry when one applies.
- `roundfix phase model <phase> --spec <slug>` loads Project Config, reads
  the Phase Record and writes nothing.
- Both take only `idea`, `prd`, `techspec` and `tasks`.
- `roundfix phase start` keeps its contract.

**How the error is read.** The first line of the message that matches
`You've hit your <name> limit` or begins `<name> limit reached`, with either
apostrophe and in any letter case, names the limit. A one-word name other
than `session` or `weekly` is a model's own limit. `session` and `weekly` are
shared limits. A name of more than one word is not a model's limit. Without
such a line the error is not a usage limit. The kept line is cut to 200
characters, and control characters become spaces.

**How a name matches a model.** Roundfix knows no model family. A limit is
the declared model's when the declared model's name contains the limit's name,
ignoring letter case. A fallback is in the limited family when its model's
name contains it too.

**The answer.** One line on standard output, written for the session to read.
`0` means start on the printed model and effort. `3` means no fallback
applies, and the line says why. `2` is a usage or configuration error, and `1`
means the record could not be written.

## Considered Options

- **The rules as skill prose.** Rejected: a session would decide them again
  each time, and no test would hold them.
- **New flags on `roundfix phase start`.** Rejected: it loads no
  configuration and reads no earlier entry (ADR-0414, ADR-0418), and it is run
  by the Author Subagent after it starts, which is too late.
- **A list of model families shipped in Roundfix.** Rejected: Roundfix ships
  no model name (ADR-0401).
- **JSON output.** Not taken: one session reads one line, and the exit code
  carries the decision.
- **Exit `1` for a refusal.** Rejected: `1` already means the record could not
  be written.

## Consequences

- A new message wording from the harness reads as "not a usage limit", and the
  phase stops as before.
- A declared model whose name does not contain its family's name, such as a
  custom deployment ID, never falls back.
- The fallback entry is written before the session restarts the Author
  Subagent, so a fallback that cannot start still holds for the Spec.
