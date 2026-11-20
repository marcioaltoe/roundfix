### permissions

Show the Flow Command List for an Agent harness. Preview first, then write or
remove exactly that plan by passing the digest it printed. There is no flag that writes without a digest.

```bash
roundfix permissions [--harness <claude|codex|opencode>] [--remove] [--apply <digest>] [--format <text|json>]
```

The command needs a Git repository and resolves its root from the current
directory. It loads neither Project Config nor User Config. With no
`--harness`, it previews each harness whose files the repository holds, in
Claude Code, Codex, OpenCode order. If none is present, it prints
`No harness files in this repository. Name one with --harness.` and exits `0`.
An explicitly named harness can be previewed before its files exist.

| Flag | Meaning |
| --- | --- |
| `--harness <claude\|codex\|opencode>` | Select one harness; required with `--apply`. OpenCode is reported, never written. |
| `--remove` | Preview removal; with `--apply`, remove what the preview showed. |
| `--apply <digest>` | Write exactly what the preview with this digest showed. |
| `--format <text\|json>` | Select text (default) or JSON. |

| Exit code | Meaning |
| --- | --- |
| `0` | A preview was printed, the plan was written, or nothing needed writing. |
| `1` | The permission file could not be read as the harness documents it, or could not be written. |
| `2` | Invalid input, including running outside a Git repository. |
| `3` | The preview changed since the digest was printed; nothing was written. |

Results go to standard output; refusals go to standard error.

#### Claude Code

Roundfix writes `.claude/settings.json`, adding the entries under
`permissions.allow` and the Guard Entries under `permissions.ask`. Other keys,
user rules and bytes outside the changed list elements stay as they were.
It opens no Harness Permission File under your home directory and never
writes `.claude/settings.local.json`.

The rules apply only after you trust this folder in Claude Code.
Guard Entries keep options that write output, run another program or change
GitHub asked about. These ask in every permission mode, a session that skips prompts included.

A fresh preview has these blocks, with every entry printed in full. This
example abbreviates the two rule lists; `<version>`, `<n>`, `<m>` and
`<digest>` stand for the installed version, actual counts and preview digest.

```text
Flow Command List for Claude Code, from Roundfix <version>
File: .claude/settings.json (absent, it would be created)
State: not applied

Would add under permissions.allow (<n>):
  Bash(roundfix archive *)
  ...
  Bash(gh pr view *)

Would add under permissions.ask, to keep these options asked about (<m>):
  Bash(git archive *--output*)
  ...
  Bash(gh api *--input*)
  These ask in every permission mode, a session that skips prompts included.

Left out, and still asked about:
  roundfix permissions: the command that applies this list
  roundfix baseline, profiles, release, skills: the Roundfix skill tells a session to ask you first
  git config, init, merge, mv, push, remote, repack, switch, worktree: they change the repository or its configuration
  gh auth: it changes your credentials
  A commit, a push, a merge and a Pull Request stay asked about, and a delivery's push and merge stay behind the Spec's authorization.

What this does not do:
  It is not a security boundary: a roundfix command runs what this repository declares.
  A listed read accepts any option; only the options listed above stay asked about.
  The rules apply only after you trust this folder in Claude Code.

Nothing was written. To write exactly this: roundfix permissions --harness claude --apply <digest>
```

Guard Entries cover Git archive output, remote and execution options; Git diff,
diff-tree and show output and external-diff options; Git grep's open-in-pager
options; and GitHub API method, field, raw-field and input options, including
their short forms where they exist. The preview prints their exact patterns.

#### Codex

Roundfix writes the whole `.codex/rules/roundfix-flow.rules` file. It contains
prefix rules for the entries, naming a command down to its subcommand, with no
argument constraint. Other files in `.codex/rules/` and `.codex/config.toml`
are left untouched.

In Codex an allowed command runs outside the sandbox with nobody asked.
No option of a listed read can be kept asked about in Codex.
For example, `git show --output=<file>` can write that file. Codex has no
Guard Entry because its rules match only the start of a command.
The rules apply only after you trust this project in Codex.

The Codex preview uses the same Named Exclusions block as Claude Code. Its
entry list is abbreviated here:

```text
Flow Command List for Codex, from Roundfix <version>
File: .codex/rules/roundfix-flow.rules (absent, it would be created)
State: not applied

Would allow (<n>):
  roundfix archive
  ...
  gh pr view
No option of a listed read can be kept asked about in Codex.

Left out, and still asked about:
  ...

What this does not do:
  It is not a security boundary: a roundfix command runs what this repository declares.
  In Codex an allowed command runs outside the sandbox with nobody asked.
  Codex rules match the start of a command, so a listed read runs with any option: git show --output=<file> would write that file.
  The rules apply only after you trust this project in Codex.

Nothing was written. To write exactly this: roundfix permissions --harness codex --apply <digest>
```

#### OpenCode

Roundfix writes nothing for OpenCode.
It reads repository configuration only to report whether shell commands are
already allowed by the base policy or the project holds rules of its own:

```text
Flow Command List for OpenCode, from Roundfix <version>
State: not needed
OpenCode's base policy allows every shell command, and this repository holds no OpenCode configuration that asks about or denies one.
Roundfix writes nothing for OpenCode.
```

When rules exist, it prints `State: own rules`, followed by
`This repository's OpenCode configuration (<file>) holds rules that ask about or deny shell commands. Compare them with the list:`,
the entries in plain words, and `Roundfix writes nothing for OpenCode.`
An apply for OpenCode exits `2`.

#### States, writes and removal

A preview writes nothing. An existing target omits the absent-file suffix.
An applied list prints `State: current` and
`Nothing to write: the file holds the installed list.` When entries differ,
it prints `State: behind, <a> to add, <r> to remove` (the add count includes
Guard Entries). Retired entries appear under
`Would remove, no longer on the list (<n>):`. A target that cannot be read
prints `State: not readable: <reason>` and offers no write digest.

After reviewing the preview, use its exact digest:

```bash
roundfix permissions --harness claude --apply <digest>
```

Claude Code reports `Wrote <n> entries and <m> Guard Entries to <path>.`;
Codex reports `Wrote <path> with <n> entries.` Each then prints its trust
sentence and `To remove them: roundfix permissions --harness <name> --remove`.

Removal also requires a separate preview and its digest:

```bash
roundfix permissions --harness claude --remove
roundfix permissions --harness claude --remove --apply <digest>
```

The removal preview prints the header and file, `Would remove (<n>):`, each
rule, and `An entry with the same text that you wrote yourself is removed with them.`
For Codex it also names the file with `Would delete .codex/rules/roundfix-flow.rules.`
It ends with `Nothing was written. To remove exactly this: roundfix permissions --harness <name> --remove --apply <digest>`.
With no entries to remove it prints
`Nothing to remove: the file holds no entry of the list.`

A confirmed Claude Code removal reports `Removed <n> entries from <path>.`;
Codex deletes only Roundfix's file and reports `Deleted <path>.`
Removal recognizes current and retired entries by their text, so a user rule
with identical text is removed too.

A stale or empty digest exits `3`, with these two lines on standard error
(the second command includes `--remove` for a removal):

```text
The preview changed since digest <digest> was printed. Nothing was written.
Run the preview again: roundfix permissions --harness <name>
```

An unreadable target on a write exits `1` regardless of the digest:
`<path> cannot be read as <Harness> documents it: <reason>. Nothing was written.`
A target that leaves the repository through a link exits `1`:
`<path> leaves the repository through a link. Nothing was written.`

#### JSON

`--format json` reports `roundfix/flow-permissions/v1`. Each plan includes
the harness, action, repository-relative file, its existence and state, the
entries and Guard Entries to add, entries to remove, Named Exclusions with
their reasons, the digest and the exact next command. Lists below are
abbreviated to one item for illustration:

```json
{
  "schemaVersion": "roundfix/flow-permissions/v1",
  "roundfix": "<version>",
  "harnesses": [
    {
      "harness": "claude",
      "action": "apply",
      "file": ".claude/settings.json",
      "fileExists": false,
      "state": "not_applied",
      "add": ["Bash(roundfix archive *)"],
      "addGuards": ["Bash(git archive *--output*)"],
      "remove": [],
      "exclusions": [{"command": "roundfix permissions", "reason": "the command that applies this list"}],
      "digest": "<digest>",
      "next": "roundfix permissions --harness claude --apply <digest>"
    }
  ]
}
```

`digest` and `next` are empty strings when nothing would be written. A write
prints the same object for its one harness with `"written": true`. Empty
rule lists are arrays, not null. The text for no detected harness is the
same in either format.

#### Scope of the list

It is not a security boundary.
A listed Roundfix command runs what this repository declares, and every
listed read accepts arbitrary arguments. Named Exclusions stay off the
list: the applying command, subcommands the skill says to ask about first,
and Git and GitHub commands that change the repository, its configuration,
GitHub or credentials. The preview groups them by reason. Shells,
interpreters, package managers and build tools receive no rule.

A session must show the preview and ask the user before applying or removing
it, in every permission mode. A recommendation, silence or a missing answer
is never that yes. Never apply or remove the list from a Run.
