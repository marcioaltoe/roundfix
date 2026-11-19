---
status: accepted
knowledge: true
created_at: 2026-11-19T17:30:00Z
updated_at: 2026-11-19T17:30:00Z
deprecated_at: null
superseded_by: null
---

# The Flow Command List names subcommands and holds every Git and GitHub read the skills name

The Flow Command List holds one entry for each `roundfix` subcommand the owned
skills name and one for each Git and GitHub read they name. An entry names a
command down to its subcommand and never constrains its arguments. What the
list leaves out is named, with the reason. A read whose options can write a
file, run a program or change GitHub is on the list all the same, by the
maintainer's choice, with a Guard Entry for those options where the harness
can express one.

## Context

A list of allowed commands is short only if something decides what belongs on
it, and it stays true only if something compares it with what the skills ask
for. Two questions decided the content, and the maintainer answered both on
2026-11-19 when he approved the Idea.

Which `roundfix` subcommands: "Os que as skills nomeiam (Recommended)", that
is, one per subcommand, except those the skill tells a session to ask the
user about first and the command that writes the list.

Which Git and GitHub reads: "Todas as leituras nomeadas". The author had
recommended only the reads a harness would ask about, in a form that cannot
write, and never `gh api`. The maintainer chose every read the skills name,
with the measured risk recorded and reduced inside that choice, not by
narrowing it.

The risk was measured on 2026-11-19. In Claude Code, an allow entry for `git
show` let `git show --output=<file>` write that file with nobody asked, and an
ask entry for the option restored the question, in a subagent and in a
session that skips every prompt too. In Codex, a rule matches the words at the
start of a command and nothing after them, so no rule tells `git show` from
`git show --output=<file>`; an allow rule means "Run the command outside the
sandbox without prompting", and inside the sandbox the same write is refused.
Pillar Security published the attack this enables: "OpenAI's Codex CLI trusts
git show as a safe command, but git show --output can write arbitrary content
to arbitrary files—including .git/config." `gh api` sends a change with one
option: "The default
HTTP request method is `GET` normally and `POST` if any parameters were
added."

## Decision

**One entry per subcommand.** A `roundfix` entry names a top-level subcommand
of the binary. A Git entry names the Git subcommand, and a GitHub entry names
the `gh` command path. No entry holds a wildcard in the place of a command,
and no entry tries to tell one option from another.

**A subcommand is on the list or off it.** Every top-level `roundfix`
subcommand that an owned skill names is on the list unless it is a Named
Exclusion. The Named Exclusions are the command that applies the list and the
subcommands the Roundfix skill tells a session to run only after it asked the
user or the user asked for it: `baseline`, `skills`, `profiles` and `release`
on 2026-11-19. A flag that is dangerous keeps the guard it has inside
Roundfix.

**Every named read.** Every Git and GitHub command an owned skill names is on
the list when, in the form the skill names it, it changes no ref, index,
working-tree file or configuration of the repository and changes nothing on
GitHub. Every other Git and GitHub command the skills name is a Named
Exclusion, and so is every commit, push, merge and Pull Request command,
named or not.

**Never on the list.** A shell, an interpreter, a package manager and a build
tool, with or without arguments.

**Guard Entries where the harness has them.** For each listed read with a
known option that writes a file, runs a program or changes GitHub, the Claude
Code rendering adds an entry that keeps that option asked about. The Codex
rendering has none, because its rules cannot express one, and it says so
where the user confirms.

**A contract keeps it true.** A test of the repository fails when an owned
skill names a `roundfix` subcommand, or a Git or GitHub command, that is
neither on the list nor a Named Exclusion.

## Considered Options

- **One entry for all of `roundfix`.** Rejected. It never falls behind, and
  it covers the subcommands the skill asks about first and the command that
  applies the list.
- **Exclusion by flag.** Rejected. A harness matches a command from its
  start, and the documentation of one says that "Bash permission patterns
  that try to constrain command arguments are fragile". The Guard Entries
  are such patterns; they add a question and never remove one, so a pattern
  that misses leaves the entry as the maintainer chose it.
- **Only the reads a harness would ask about, in a form that cannot write.**
  The author's recommendation, refused by the maintainer.
- **In Codex, leave the Git reads to the sandbox.** The measurement favors
  it: a Git read runs inside the sandbox with no rule, and the rule only
  takes it out. It is the refused option applied to one harness, so it is
  not taken; the maintainer can still choose it.
- **A partial guard in Codex**, a prompt rule for the option as a word of its
  own. Rejected: the measurement shows the same option passing in two other
  spellings, and a guard that looks whole and is not is worse than none.
- **Leave `gh api` out**, as the Idea first proposed. Overruled by the
  maintainer's answer, which names every read.

## Consequences

- In Claude Code a Guard Entry asks in every permission mode. A session that
  skips prompts is asked, or refused when it cannot ask, before `gh api` with
  an option that sends a change, in a project that applied the list.
- In Codex a project that applies the list lets every listed Git and GitHub
  read run outside the sandbox with any option. The preview, the guide and
  the confirmation say so.
- Most named Git reads already run with no rule in both harnesses. The
  entries change the outcome for the forms a harness does not class as a
  read by itself.
- A new subcommand, or a new Git or GitHub command in a skill, fails the
  contract until someone decides on which side it falls.
- The list of Named Exclusions follows the Roundfix skill's text. A skill
  change that tells a session to ask first about another subcommand moves
  that subcommand, and the Doctor Command reports the projects that are then
  behind.
