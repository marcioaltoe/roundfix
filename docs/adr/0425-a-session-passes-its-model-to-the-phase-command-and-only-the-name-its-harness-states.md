---
status: accepted
created_at: 2026-11-17T19:30:00Z
updated_at: 2026-11-17T19:30:00Z
deprecated_at: null
superseded_by: null
---

# A session passes its model to the phase command, and only the name its harness states

A session reports its model to `roundfix phase start` with the flag
`--model`, and it passes only the exact identifier its harness states in the
session's own instructions; a session that is told none passes no flag, and
the entry reads `unknown`. The reasoning effort is read from `CLAUDE_EFFORT`
and taken from `--effort` only when that variable is unset, and
`--author-subagent` says who runs the phase. A reported value never costs the
entry.

## Context

ADR-0418 decides that the session reports the model and leaves the means to
the TechSpec. Measured a third time on 2026-11-17, in the Author Subagent of
this TechSpec on Claude Code 2.1.325: `CLAUDE_CODE_SESSION_ID` was set,
`CLAUDE_EFFORT` held `xhigh`, the effort this repository declares for the
phase, no variable's value named a model, and the subagent's instructions
stated its model by its exact identifier.

The harness documents the same. Its reference of environment variables, read
through Exa on 2026-11-17, says of `CLAUDE_EFFORT`: "Set automatically in Bash
tool subprocesses and hook commands to the effort level setting when the
subprocess starts", and "Only set when the current model supports the effort
parameter". It lists `ANTHROPIC_MODEL` as the "Name of the model setting to
use", which `--model` and `/model` override, and no variable that names the
model that is running.

## Decision

**Three flags.** `roundfix phase start` takes `--model <identifier>`,
`--effort <level>` and `--author-subagent`. They apply to the `idea`, `prd`,
`techspec` and `tasks` phases; on another phase they are ignored with one
warning, and the entry is written as before.

**The model.** The value of `--model`, recorded with the source `reported`,
when the environment names a Claude Code session and the value is one model
identifier of at most 200 characters. Otherwise the model is `unknown`.

**The effort.** The value of `CLAUDE_EFFORT`, with the source `environment`.
When the variable is unset, the value of `--effort`, with the source
`reported`. Otherwise `unknown`.

**Who ran the phase.** `author-subagent` with the flag, `session` without it,
and `unknown` when the environment names no session.

**What the skills say.** Each authoring skill tells the session to pass the
exact model identifier its harness states for it, never an alias and never a
guess, to leave the flag out when its instructions state none, and to pass
`--author-subagent` when it is the phase's Author Subagent. When the
installed Roundfix refuses the flags, the skill runs the command without
them and says so in its report.

**Never a refusal.** A value that is too long or malformed is recorded as
`unknown` with one warning on standard error. The entry is written, and the
command exits as it would have without the value.

## Considered Options

- **`ANTHROPIC_MODEL` or another variable as the source.** Rejected: it
  configures a model, a session may run on another, and a subagent's model
  comes from its definition.
- **The model as a positional argument, or in a second command.** Rejected: a
  flag is refused by an older Roundfix with its usage error, which the skill
  can answer by recording the phase as before.
- **A closed list of model names.** Rejected: Roundfix ships no model name
  (ADR-0401), and a new model would read `unknown`.
- **A refusal when an Authoring Phase is recorded without a model.**
  Rejected by ADR-0418.
- **Asking the session to guess from the Authoring Model Declaration.**
  Rejected: the declared model is not always the one that ran.

## Consequences

- A session can report a wrong identifier. The entry marks the value
  `reported`, and the report groups by the text it was given.
- An alias and an exact identifier of one model are two groups, so the skills
  ask for the exact identifier.
- The result line of an Authoring Phase names the model, the effort, how each
  was learned and who ran the phase, so a session that left the model out
  reads `model unknown` at once.
- A harness that stops setting `CLAUDE_EFFORT` leaves the effort to the flag.
