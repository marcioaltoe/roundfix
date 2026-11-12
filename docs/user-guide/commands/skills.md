### skills

```bash
roundfix skills list
roundfix skills check
roundfix skills install [--target codex|claude|opencode|all]
```

The binary ships 12 Roundfix-owned skills: the operational `roundfix` skill
plus the authorial workflow skills (`write-idea`, `write-prd`,
`write-techspec`, `write-tasks`, `setup-context-driven`, `implement-task`,
`brainstorming`, `business-analyst`,
`archive-spec`, `qa-gate`, `evidence-gate`). `skills list` also prints
recommended external skills, which install through your own tooling and are
never shipped. By default `skills install` writes to `<repo>/.agents/skills`;
`--target` selects user-scoped Agent skill directories instead.

Inside an adopted Git repository, `skills check` also inspects the skills the
repository's Baseline requires and fails with exit `1` for each Missing
Required Skill, printing
`.agents/skills/<name>: required by the Baseline and not installed; restore with roundfix baseline skills restore --profile <id> --skill <name>`
on stderr. Outside a Git repository, or without a Setup Manifest, it checks
only the shipped bundle.

It also names each **Stale User Skill Copy**, a user-level Roundfix skill copy older than the binary, which an Agent harness can load in place of the project copy. It reads the copies from these directories, in order: `$CLAUDE_HOME/skills` or `~/.claude/skills`, `$CODEX_HOME/skills` or `~/.codex/skills`, `$OPENCODE_HOME/skills` or `~/.opencode/skills`, `~/.agents/skills` and `~/.config/opencode/skills`. Symbolic links are followed; when several links reach one folder, that folder is named once under the first directory. A copy at or above the binary's version, one without a comparable version and an unreadable one are not named. Nothing is written, and the exit code does not change.

```text
roundfix: warning: user-level Roundfix skill copies older than this binary; an Agent harness may load them in place of the project copy:
  <skill> <found> (binary <version>) at <path>[ -> <link target>]
  refresh: <command>
```

Refresh the first three directories with `roundfix skills install --target <claude|codex|opencode>`, `~/.agents/skills` with `roundfix skills install --target codex --dir <path>`, and `~/.config/opencode/skills` with `roundfix skills install --target opencode --dir <path>`. A refresh writes through a link to the folder the ` -> <link target>` names.

After bundle installation, an existing `skills-lock.json` entry naming an
owned skill receives the installed folder's `computedHash` and loses its
`ref`; every other field is preserved. No entry or lock is created when no
matching lock entry exists. Each name whose hash changed is reported as
`Updated skills-lock.json computedHash: <name>`, followed by
`Removed skills-lock.json ref: <name>` for each entry whose `ref` was removed.

Roundfix reads a `skills-lock.json` of version 1 or 2, the version the skills CLI 1.8 writes; every `computedHash` it writes into the lock or compares with an installed folder follows the hash algorithm of the version the file declares, and it never changes that version.
