---
status: accepted
created_at: 2026-11-10T00:00:00Z
updated_at: 2026-11-10T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A user-level copy of a Roundfix-owned skill older than the binary is reported, never changed

`roundfix skills check`, `roundfix doctor` and `roundfix update check` read
the user-level skill directories that Claude Code, Codex and OpenCode load
and name each copy of a Roundfix-owned skill that declares a version below the
one the installed binary carries, with the `roundfix skills install` command
that refreshes it. None of them writes to those directories.

## Context

On 2026-11-10 a Claude Code session in this repository loaded the Roundfix
skill at version 0.1.140 from `~/.claude/skills/roundfix` while the
repository's copy and the installed binary carried 0.1.208; seven more owned
skills were behind in the same way, and `roundfix skills check` passed. Claude
Code documents that a personal skill wins over a project skill of the same
name; OpenCode registers `~/.config/opencode/skills` after a project's
`.agents/skills`, so that copy wins too; Codex lists a user copy beside the
project copy without merging them. `roundfix skills install --target <claude|codex|opencode>`
writes those copies, and nothing that upgrades Roundfix or refreshes a
repository touches them afterwards. The checks Roundfix had compared only the
repository's copies with the binary (ADR-0189, ADR-0364). The maintainer chose
"Skills de usuário velhas (Recommended)" — option text: "O `doctor` e o
`skills check` passam a apontar uma cópia de usuário de skill do Roundfix mais
velha que a do binário, com o comando que a atualiza. Hoje cada release deixa
as do vault velhas sem aviso."

## Decision

**What is read.** Five directories, in this order: the three directories
`roundfix skills install` writes for `--target claude`, `codex` and
`opencode` (each `$CLAUDE_HOME`, `$CODEX_HOME` or `$OPENCODE_HOME` followed by
`/skills`, else `~/.claude/skills`, `~/.codex/skills` or
`~/.opencode/skills`), then `~/.agents/skills`, Codex's documented user
directory and an OpenCode source, and `~/.config/opencode/skills`, OpenCode's
global directory. Only the twelve Roundfix-owned names are read; a copy's
version is the `version` its `SKILL.md` frontmatter declares.

**Symbolic links are followed.** A skill folder or a skills directory that is
a link is read where it leads, and a folder that several links reach is named
once, under the first directory in the order above, with the folder it
reaches. A missing name, a dangling link and an unreadable file are skipped.

**Older is the one comparison.** A copy is stale when its declared version is
below the binary's version of that skill, by the comparison ADR-0189 defines
and ADR-0364 applies to repository copies. A copy at or above it, or without a
comparable version, is not named: a user may keep a deliberate variant of a
skill, and only a declared older version proves an install the binary has
outgrown.

**Reported, never changed.** `skills check` prints a warning on stderr and
keeps its exit code; Doctor prints a `user-skills: warn` line after `skills:`
and keeps its exit code; the Update Check prints a `User skills:` line and the
JSON field `userSkills` and keeps `next`. Each names the refresh command of the
first directory that reached the copy: `roundfix skills install --target
<claude|codex|opencode>` for the three install directories, with `--dir <path>`
for the other two. The `--update` Skill Mode asks before running one. Nothing
prints when no copy is stale.

## Considered Options

- Refresh the user-level copies on `roundfix upgrade` or in the Update Check:
  rejected because those directories belong to the user's harnesses and, on
  the maintainer's machine, lead through links into a synchronized vault; a
  silent rewrite outside the repository is what the maintainer's answer
  replaced with a report.
- Fail `roundfix skills check` or Doctor on a stale copy: rejected because
  `make verify` and the release plan run them against the repository, and a
  file in one developer's home must not fail a repository gate.
- Add the finding to Doctor's `skills:` line: rejected because the release
  plan reads that line as the repository's skill readiness.
- Read only the directories `skills install` writes: rejected because Codex
  documents `~/.agents/skills` and OpenCode `~/.config/opencode/skills`, and
  a copy there wins or appears beside the project's.
- Name every copy whose content differs from the binary's: rejected because a
  newer or a deliberately edited user copy is not stale, and ADR-0364 keeps the
  version comparison as the owned-skill rule.

## Consequences

After each release, the first session's Update Check names the user-level
copies the new binary has outgrown, and one command per directory refreshes
them. A harness directory that no install target writes and that this record
does not list, such as one named by `CLAUDE_CONFIG_DIR` or an `opencode.json`
`skills` array, is not read.
