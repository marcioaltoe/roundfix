---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A repository rule becomes a Baseline clause only when it recurs

Adopters write rules by hand in their Repository-Specific Normative Rules, and
an audit on 2026-09-30 found the same rules written again and again: database
mutation authority in five repositories, Verification that must fail first in
five, lint warnings that must block in three. A rule written once records one
repository's choice; a rule written in several repositories is a gap in the
Baseline.

A hand-written rule is promoted into the Baseline only when it appears in at
least two repositories, or when one repository tied it to measured rework. It
is also declined when it binds one machine or one domain, when the formatter,
the linter or a recorded decision already owns it, or when a shipped clause
already says it. A promoted rule becomes a clause with force in the module
that owns its concern, and its promotion is declared as a Baseline change.
When it supersedes an existing clause, the new clause names the old one in
`replaces` with the same enforcement, so an adopter whose Source Baseline
retains the old clause records it as `replaced`.

The first promotion follows this rule. Among the rules it promotes, lint
warnings now block in every repository: the clause that blocked warnings only
when a repository's Verification already treated them as errors is replaced by
a core clause that makes a lint warning fail Verification.

## Consequences

The Baseline grows only with evidence from outside Roundfix. Adopters whose
linters only warn must change their Verification to fail on a warning, which
is a tooling change they authorize in their own repository. A rule that a
single repository keeps without measured rework stays where its author wrote
it.
