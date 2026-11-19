---
status: accepted
knowledge: true
created_at: 2026-11-19T17:30:00Z
updated_at: 2026-11-19T17:30:00Z
deprecated_at: null
superseded_by: null
---

# Roundfix writes a harness permission file only as a list the user saw and confirmed

Roundfix writes allow entries into a project's permission file of an Agent
harness in one case only: a person ran the command that shows the closed Flow
Command List for that harness, and then confirmed that exact preview. It
writes the project's file, never a file of the user's own, and never a
permission mode. Nothing else in Roundfix writes such a file, and these files
become Governed Paths.

## Context

A supervised session of Claude Code or Codex asks the person before it runs a
command that no rule allows. The flow is made of `roundfix` commands: the
owned skills name 35 subcommands in 447 places. The maintainer described the
result on 2026-11-19: "Outra coisa que normalmente trava o fluxo é o uso de
comandos bash, se houver algo que possa ser feito para resolver isso tb, uma
inclusão de autorização permanente seria muito bom. Mas depende do harness."

Until now Roundfix wrote no settings file of a harness, and two decisions
explain why. For a Run, the full-access mode is set through the protocol and
"not by editing the user's runtime configuration files" (ADR-0011), and
"Editing `~/.codex/config.toml` was rejected again" (ADR-0293). For a session,
a hook merged into a harness's project file was refused as the Baseline's own
mechanism, because it "means merging repository-owned JSON
(`.claude/settings.json`, `.codex/hooks.json`) and an OpenCode plugin, which
is tooling configuration the core clauses forbid changing without express
authorization" (ADR-0371).

The maintainer chose the form on 2026-11-19: "Mostra a lista e grava com seu
sim (Recommended)", with the target "Sessão Supervisor e autores
(Recommended)" and the file ".claude/settings.json (Recommended)".

A published figure limits what a prompt is worth, "in practice users accept
93% of them anyway" (Anthropic, 2026-03-25), and published attacks limit what
a list is worth: "Command-name allowlists don't work for commands with complex
option semantics" (Pillar Security, 2026-07-20).

## Decision

**One command shows and writes.** `roundfix permissions` prints, for a
harness, the entries it would add, the file it would write, what the list
leaves out and what the list does not protect against, and it writes nothing.
The same command writes only when it is given the digest that a preview
printed, and only when that digest still describes the list and the file.

**The yes is the user's.** The command is not on the list it applies, so a
harness that asks will ask the person before it runs. The Roundfix skill tells
a session to show the preview and to put one question to the user before it
applies, in every permission mode. `roundfix baseline update`, `roundfix
skills install`, `roundfix setup`, the Doctor Command, the Update Check and
every Agent of a Run never apply it.

**Only the project's file.** Claude Code: the shared project settings file.
Codex: a rules file of Roundfix's own in the project's rules folder. OpenCode:
a report and no write. No file under the user's home directory is read for
its rules or written, and no write follows a link that leaves the repository.

**Entries, and nothing else.** The command adds and removes its own entries.
It sets no default mode, approval policy or sandbox mode, and it changes no
other byte of a file it did not create.

**The files are tooling configuration.** The project permission files of the
three harnesses join the governed set, so a Task changes one only under a
grant that names it.

**It is not a boundary.** The list removes prompts. It isolates nothing, and
every surface that shows it says so.

## Considered Options

- **Print the entries and let the person paste them.** Rejected by the
  maintainer's answer. Nothing would tell a project that its pasted list fell
  behind, and a paste is where a wildcard gets added.
- **Write at installation or at a Baseline update.** Rejected. That is the
  write nobody asked for which ADR-0371 refused; a permission is asked for by
  name.
- **Write the user's own file.** Rejected. It reaches every repository of the
  machine, no teammate can review it, and ADR-0011 and ADR-0293 already
  decided that Roundfix does not edit it.
- **A hook that approves the commands.** Rejected. A hook is code where a list
  is readable, and ADR-0371 decided that subject.
- **Recommend a session with every prompt off inside a sandbox.** Kept as the
  dissent. It may be the stronger answer where a person can run one; this
  decision serves the session that keeps prompts on.
- **Record the yes in Project Config.** Rejected. The harness's file is the
  record, and a new key would make every older binary refuse the repository.

## Consequences

- A project that applies the list commits a file that grants commands to
  every session of that repository whose user trusted the folder. One
  person's yes becomes the team's through review of that commit.
- A rule that is written is not always in force: each harness applies a
  project's allow rules only after the person trusts the folder, and Roundfix
  cannot read that state.
- A repository that already edits a harness permission file in a Spec needs a
  grant for it from the release on.
- This decision does not contradict ADR-0371, ADR-0011 or ADR-0293: no hook
  is written, nothing is written without a person's confirmation, and the
  user's own files stay untouched.
