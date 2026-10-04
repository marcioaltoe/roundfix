---
status: accepted
created_at: 2026-10-04T00:00:00Z
updated_at: 2026-10-04T00:00:00Z
deprecated_at: null
superseded_by: null
---

# An unrecorded branch prefix follows the commit-type rule

The oraculum maintainer asked on 2026-10-01 to leave the branch prefix to the
rule the agent-instructions guide already states: new work branches are named
`<type>/<description>` from the work's Conventional Commit type. The Baseline
did not allow it. The core module required `branch.prefix`, and the guide
always rendered "The branch-prefix pattern is ...", so a repository could only
repeat the commit-type rule as a prefix value or record one of its own, which
the purpose-based branch policy forbids for personal prefixes.

`branch.prefix` is now an optional decision, as `frontend.layout` became one
in ADR-0205. The core module no longer requires it, and the built-in profiles
still offer it, so first adoption still asks it. While a repository records
none, the guide says that no branch prefix is recorded and names the
commit-type rule; a recorded value keeps today's sentence byte for byte. No
adopter is asked anything on its next update, and no recorded value is changed
or removed.

This narrows ADR-0118 for this one decision. ADR-0118 held that "a decision
nobody answered is not a default nobody chose", and that remains true: the
unrecorded wording names no prefix and applies no default, it states the rule
every repository already carries. Two alternatives were rejected. An explicit
`none` value would have kept the decision required and made a repository
record a non-answer. Rendering nothing while it is unrecorded would have
dropped the sentence that tells an Agent how branches are named.

## Consequences

A repository-owned profile may now leave the decision out, because no module
requires it. A repository that recorded a value keeps it; Roundfix offers no
command that clears a recorded decision.
