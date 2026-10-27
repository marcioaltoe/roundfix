---
status: accepted
created_at: 2026-10-27T01:30:00Z
updated_at: 2026-10-27T01:30:00Z
deprecated_at: null
superseded_by: null
---

# A profile change carries the recorded decisions, and a restore keeps the Claude skill links

ADR-0263 keeps the managed refresh on the recorded Baseline Profile and
rejected an automation path that "asks an adopter to restate every recorded
decision to change one". A Baseline Profile change still goes through
`roundfix baseline plan --profile <id>` and `roundfix baseline apply`, and on
2026-10-26 an adopter changing to the `python` profile restated eleven
recorded decisions by hand across three refusals. Separately, the skills CLI
that installs upstream skills keeps the canonical copy under
`.agents/skills/<name>` and gives Claude Code a symbolic link
`.claude/skills/<name>` to it, while `roundfix baseline skills restore`
writes only the canonical copy, so a skill the restore adds is invisible to
Claude Code in such a repository.

**The plan carries what the Setup Manifest records.** `baseline plan` reads
the repository's Setup Manifest when one exists and takes, as a Carried
Decision, each recorded decision the selected profile selects, does not fix,
and that no `--decision` or `--decision-file` answers, when the catalog still
accepts its value. A flag always wins. A recorded decision the profile does
not select, fixes, or whose value the catalog refuses is not carried and is
named in a plan warning. Every answer still missing, the decisions a missing
answer's catalog effect would require, and `preservation.mode` are named in
one refusal. `baseline update` keeps refusing a profile change (ADR-0263) and
names the plan command that performs it.

**A restore keeps the Claude Skill Links.** When `.claude/skills` is a real
directory that already holds at least one symbolic link to its
`.agents/skills/<name>` counterpart, the restore plans one Claude Skill Link
for each selected skill that has no `.claude/skills/<name>` entry, binds it
into the Plan Digest, and writes it inside the same recoverable transaction
as the skill tree (ADR-0073). An existing entry of any kind is the
repository's and is left as it is; a `.claude/skills` that is absent, is a
symbolic link, or holds no such link gets nothing.

## Considered options

Letting `baseline update --decision profile=<id>` change the profile was
rejected: the update is the Managed Refresh, which takes no root backup and
accounts no Source Baseline (ADR-0100), while a profile change regenerates
the root carriers and needs a preservation mode. Carrying recorded decisions
only in the interactive command was rejected because Agents and scripts
cannot use it. Writing the links after the transaction committed was
rejected for the same reason ADR-0073 rejected sequential best-effort
writes. Linking whenever `.claude/skills` is a directory, as `roundfix skills
install` does for the Roundfix skill, was rejected for external skills: a
directory of copies, kept where symbolic links do not work, would gain links
it cannot follow.

## Consequences

A profile change in an adopted repository needs only `--profile`,
`preservation.mode` and any decision the new profile introduces. A plan's
digest now depends on the Setup Manifest's recorded decisions when flags
leave them unanswered. A repository that keeps per-skill Claude links gets
one for every skill a restore installs, and an installed skill that lacks its
link gets it from a restore of that skill.
