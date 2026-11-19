---
status: accepted
knowledge: true
created_at: 2026-11-14T11:00:00Z
updated_at: 2026-11-14T11:00:00Z
deprecated_at: null
superseded_by: null
---

# A user-level skill copy is reported with the command that removes it

`roundfix skills check`, Doctor and the Update Check name each user-level copy
they report with the command that removes it, before the command that
refreshes it, and `roundfix skills remove` is that command. They also name,
whatever its version, a copy in Claude Code's user-level skill directory of an
authoring skill whose phase the repository declares, because that copy hides
the Declared Skill Header. This amends ADR-0386.

## Context

ADR-0386 reported a Stale User Skill Copy with the `roundfix skills install`
command that refreshes it and compared nothing but the version. On 2026-11-13
the maintainer decided that the Roundfix skills live only in the project, and
pantheon reported that the refresh command is how the copies had come back:
"a ação certa é apagar a cópia de usuário, não atualizá-la". Claude Code
documents the precedence: "Enterprise over personal, and personal over
project". ADR-0401 writes a header into the project's copy of an authoring
skill, and a personal copy of that skill, older or not, makes Claude Code load
a file without it. The maintainer chose to join this change to the Spec that
declares the header.

## Decision

- **The command.** `roundfix skills remove --target <claude|codex|opencode|all>
  [--dir <path>] [--apply]` acts on the directories `roundfix skills install`
  writes for those targets and only on the twelve Roundfix-owned names.
  Without `--apply` it lists what it would remove and changes nothing. With
  it, a symbolic link is unlinked and the folder it leads to is kept, and a
  real folder is removed with its content. `--target project` is refused: the
  project's copies are the Repository Skill Set.
- **The reports.** Each copy the three reports name carries a `remove`
  command before its `refresh` command, and the Update Check's JSON entry
  gains `remove`. No field is renamed.
- **A copy that hides a declared header.** When the repository declares
  `skills.authoring.claude.<phase>`, a copy of that phase's skill under
  `$CLAUDE_HOME/skills`, else `~/.claude/skills`, is named whatever its
  version, with the words that it hides the declared header, and its JSON
  entry carries `hidesDeclaredHeader`. The other four directories and the
  other skills keep the older-than-the-binary comparison.
- **Still reported, never changed.** The three reports write nothing and keep
  their exit codes. Only a person, or the `--update` Skill Mode after the user
  approves, runs a removal or a refresh.

## Considered Options

- **Print `rm -rf <path>`.** Rejected: it differs per operating system and
  shell, removes whatever the path names, and follows a mistyped path.
- **Remove the copies during an installation or an upgrade.** Rejected, as in
  ADR-0386: those directories belong to the user's harnesses and can lead into
  a synchronized folder.
- **Name every user-level copy always.** Rejected: outside a repository that
  declares a header, a copy at or above the binary's version harms nothing
  Roundfix can prove, and ADR-0386 keeps a user's deliberate variant unnamed.
- **Remove by default and ask for a flag to list.** Rejected: the commands of
  this CLI that delete list first and delete with `--apply`.

## Consequences

- The `--update` Skill Mode asks once which command to run and recommends the
  removal, because the project's copy then loads and no later release leaves
  the user's copy stale.
- A removal inside a skills directory that is itself a link deletes the folder
  where the link leads; the listing names that folder first.
- A user who keeps user-level copies on purpose keeps the refresh command.
