---
status: accepted
knowledge: true
created_at: 2026-11-17T21:30:00Z
updated_at: 2026-11-17T21:30:00Z
deprecated_at: null
superseded_by: null
---

# An Authoring Phase falls back to its declared fallback when the declared model's usage limit runs out

An entry of the Authoring Model Declaration may name a fallback model and
reasoning effort. When the phase's Author Subagent ends because the usage
limit of its declared model ran out, the session tells the user and starts the
same Author Subagent again on the fallback, passing that model and effort, with
the brief and the work the first run left. Only a model's own usage limit
switches; a limit every model shares, or any other error, stops the phase as
before. This amends ADR-0401, under which the session always passed no model
and no effort.

## Context

On 2026-11-17 the TechSpec Author Subagent of Spec 0404, declared `fable`,
ended with "You've hit your Fable limit · resets Nov 22 at 8pm" partway
through the phase. The authoring skills said to report the error and author
in the session only if the user said so, so the phase stopped until the
maintainer answered and the work was restarted on Opus by hand. The
maintainer asked for the switch to be automatic "para não perder qualidade".

Read on 2026-11-17 in the Claude Code documentation, with the installed
2.1.326:

- No frontmatter declares a fallback. A skill's frontmatter takes `model` and
  `effort`, and so does a subagent's.
- The harness's own fallback chain, `fallbackModel` or `--fallback-model`,
  switches a subagent too, but "Authentication, billing, rate-limit,
  request-size, and transport errors … never trigger a switch".
- A subagent that ends on a usage limit fails with the error, and its partial
  output reaches the session. A resumed subagent keeps the model it started
  with.
- The model and effort passed when a subagent is started override its file.
  The model takes an alias there, not a full model identifier.
- The Fable, Opus and Sonnet limits apply only to their own family, so
  another model keeps working. The session and weekly limits are shared by
  every model.

## Decision

**The declaration.** An entry `idea`, `prd`, `techspec` or `tasks` under
`skills.authoring.claude` may hold `fallback`, with `model`, and optionally
`reasoning_effort`. The research entry takes none. Roundfix still ships no
model name. A fallback that names the entry's own declared model is refused,
because it escapes no limit of that model.

**The trigger.** Only an error that says a named model's usage limit ran out,
as Claude Code writes it, switches. The model the limit names must be the
declared model of the phase, and the fallback must not be that model. A
session limit, a weekly limit, a spend limit and every other error stop the
phase and are reported as before.

**The switch.** The session tells the user which phase switched, from which
model to which, and the limit that ran out, and goes on without waiting. It
starts the phase's Author Subagent again, passing the fallback's model and
effort, with the original brief, the output the first run returned and the
files it wrote. No second subagent file is generated.

**One step.** A fallback has no fallback. When the fallback's own limit runs
out, the phase stops and is reported.

**What a typed skill does.** A phase the user starts by typing the skill runs
in the session, which cannot change its own model. It keeps the harness's own
remedy.

## Considered Options

- **The harness's fallback chain.** Rejected: it leaves rate limits out.
- **A fallback field in the skill's or the subagent's frontmatter.** Not
  available: the harness reads none.
- **A rule in the Baseline with a fixed fallback.** Rejected: it ships a model
  name, which ADR-0401 refuses, and leaves the choice outside the repository.
- **A generated fallback subagent for each phase.** Rejected by the
  maintainer: the model and effort passed at the start already override the
  file. It would have allowed a full model identifier as the fallback.
- **Ask the user before switching.** Rejected by the maintainer, who chose to
  switch and be told.
- **A chain of fallbacks.** Left for later: one fallback answers the case
  observed.

## Consequences

- A fallback must be a model alias the harness accepts when a subagent is
  started. A full identifier is refused by the start, and the phase stops as
  before.
- The fallback starts fresh: what the first run decided and did not write or
  return is lost.
- An error text the harness changes is not recognised, and the phase stops as
  before.
- A Roundfix older than the release that carries the key refuses a Project
  Config that holds it (ADR-0408), so a repository declares a fallback only
  after every machine that runs its Roundfix is upgraded.
