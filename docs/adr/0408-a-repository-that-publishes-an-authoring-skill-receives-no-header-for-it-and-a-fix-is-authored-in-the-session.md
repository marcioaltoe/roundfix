---
status: accepted
created_at: 2026-11-16T11:00:00Z
updated_at: 2026-11-16T11:00:00Z
deprecated_at: null
superseded_by: null
---

# A repository that publishes an authoring skill receives no header for it, and a fix is authored in the session

Installation writes no Declared Skill Header into an authoring skill whose
source the repository itself publishes under `skills/<name>/`, and still
generates the declared subagents. An authoring skill the agent activated on
the fix or refactor route writes in the session and hands the phase to its
Author Subagent only when the user asks. Both amend ADR-0401, which left this
repository without a declaration and handed every phase the agent starts to
the declared author.

## Context

ADR-0401 says: "This repository's `.agents/skills` are the source of the
shipped bundle, so it does not declare the section until a later decision says
how". Measured on 2026-11-16 in a throwaway clone of this repository with the
declaration added: `roundfix skills install --target project` wrote `model`
and `effort` into the four authoring skills, the mirror under `skills/` then
differed from them, and the five subagent files were generated. The
maintainer decided that this repository adopts the generated subagents alone,
with no header in its skills, and that a skill the user types there runs on
the session's model.

The same ADR says an authoring skill the agent activated "hands its phase to
that subagent". The maintainer keeps a small fix on the session's model and
chose the exception by route: the fix or refactor route writes in the
session and delegates only on request, a feature always delegates, and the
rule ships in the skills for every adopter.

The installed binary of the Delivery Queue reads Project Config from the
default branch and refuses a key it does not know, so this repository cannot
gain a new configuration key in the same change that teaches the binary that
key.

## Decision

**A published skill receives no header.** When the repository root holds
`skills/<name>/SKILL.md` for an authoring skill whose phase is declared,
`roundfix skills install --target project` and the skills stage of
`roundfix baseline update` write that skill without the Declared Skill Header
and say so. They generate every declared subagent as before. The rule is read
from the tree; no configuration key turns it on.

**The checks agree.** For such a skill `roundfix skills check`, the Doctor
Command and the Update Check count the declaration as applied when the
installed skill holds no `model` and no `effort`, and report a header found
there as the problem. A user-level copy of that skill hides no header.

**A fix or a refactor is authored in the session.** In `write-techspec` and
`write-tasks`, a skill the agent activated for a fix or a refactor runs the
phase in the session, on the session's model, and starts the phase's Author
Subagent only when the user's request or one of the user's answers asks for
it. A feature's phases are handed over as ADR-0401 decided. The report says
where the phase ran.

## Considered Options

- **A configuration key that withholds the header.** Rejected: the installed
  binary would refuse this repository's Project Config between the merge and
  its own upgrade, and one more key would describe what the tree already
  shows.
- **A marker file in the skills directory.** Rejected: a second thing to keep
  true beside the mirror that already says the repository publishes the
  skill.
- **Removing the header when the mirror is synchronized.** Rejected: it
  changes the build's own tooling, and the installed skills would differ from
  their mirror between two synchronizations.
- **Asking, for each fix, whether to delegate.** Rejected: one more question
  per fix; the user who wants the declared author says so once.
- **A judgment of size instead of the route.** Rejected by the maintainer:
  the route is read from the Spec, and a judgment would vary.

## Consequences

- In a repository that publishes an authoring skill, typing that skill runs
  the phase on the session's model; only a phase the agent hands over runs on
  the declared one.
- Until the installed binary is the release that carries this decision, an
  installation run with the older binary in this repository writes the
  headers; the test that keeps `model` and `effort` out of every shipped
  skill and the mirror check stop that tree before it is committed.
- An adopter who wants every TechSpec and every Task Graph of a fix on the
  declared model asks for the Author Subagent in the request.
