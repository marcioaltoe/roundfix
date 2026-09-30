---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# The Roundfix skill and the command reference are read one command at a time

The Roundfix Skill was one 142 KB `SKILL.md`, and the command reference one
79 KB `docs/user-guide/commands.md`. An Agent that needed one command loaded
the whole file. Almost every CLI Task declared both files, so the Wave
collision rule forced those Tasks into a serial chain, and three QA gates in
two days refused because a Task wrote into a section all commands share.

Both documents are now an entry file plus one file per command family. The
skill's `SKILL.md` keeps what every reader needs: the contract of an assigned
Batch, the shared QA settlement table, and an index that says which file under
`references/` to read for which command. The command reference keeps its
global contract and an index of one file per command under
`docs/user-guide/commands/`. There is no generator: the files are the source,
and the existing mirror copy carries the skill's references unchanged.

The split moves bytes and rewords nothing. A contract that pinned a phrase in
either document now reads the entry file together with its companion files as
one text, so the phrase may live in any of them.

## Consequences

- A Task declares the one command file it changes. Two Tasks that change
  different commands no longer collide.
- The published guidance for agent skills asks for a short entry file with
  references loaded on demand, one level deep. A single-file skill stays
  easier to paste and print, and that is given up here because no reader
  needs every command at once.
- An adopting repository receives the new layout the next time its Roundfix
  skill is installed or restored. The old single file stays valid until then.
- The reference file names are fixed by command family. A new command joins
  the file of its family or adds one file and one index row.
