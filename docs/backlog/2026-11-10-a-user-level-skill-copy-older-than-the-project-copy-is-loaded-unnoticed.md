---
type: fix
status: open
created: 2026-11-10
spec: null
reason: null
---

# A user-level skill copy older than the project copy is loaded unnoticed

## Symptom

On 2026-11-10 a Claude Code session in this repository ran `/roundfix --follow-up` and received the Roundfix skill at version 0.1.140 from the user-level skills directory. The repository copy in `.agents/skills/roundfix`, which `.claude/skills` links to, was at 0.1.208, the version Roundfix 0.137.7 ships. The loaded text had no `--update` Skill Mode, no handoff claim in the follow-up mode and no Waiting Items. Seven more owned skills were behind in the same way: `qa-gate` 0.0.19 against 0.0.25, `write-tasks` 0.0.19 against 0.0.31, `write-prd` 0.0.13 against 0.0.18, `write-techspec` 0.0.14 against 0.0.18, `archive-spec` 0.0.14 against 0.0.19, `implement-task` 0.0.5 against 0.0.6 and `setup-context-driven` 0.0.6 against 0.0.8. The user-level copies dated from 2026-10-26, and `roundfix skills check` passed.

## Where

`roundfix skills install --target <claude|codex|opencode|all>` writes the user-level copies. `roundfix upgrade`, `roundfix update check`, the `--update` Skill Mode and `roundfix baseline update` refresh the binary and the project's `.agents/skills` and never those copies. `roundfix skills check` compares the project copies. Claude Code loads a user-level skill in place of the project skill of the same name, so one stale copy wins in every repository on the machine.

## Expected

A user-level copy of a Roundfix-owned skill that is older than the copy the installed binary embeds is reported with the command that refreshes it, by `roundfix skills check`, by `roundfix doctor` and by the Update Check, or the path that upgrades Roundfix refreshes the copies it installed. Whether to report or to refresh, and which harness directories to read, are choices for the Spec.

## Evidence

Measured by the operator on the maintainer's MacBook on 2026-11-10 at about 15:50 BRT: the session header named `~/.claude/skills/roundfix` as the skill's base directory, the `version:` fields of `~/.claude/skills/<name>/SKILL.md` and `.agents/skills/<name>/SKILL.md` differed as listed above, and `roundfix skills check` printed `Roundfix skill check passed`. The maintainer approved a refresh; after `roundfix skills install --target claude` the twelve owned skills matched the repository copies.
