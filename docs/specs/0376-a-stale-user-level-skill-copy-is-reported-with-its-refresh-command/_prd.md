---
spec: 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
status: active
created: 2026-11-10
surfaces: [backend, cli, docs]
type: fix
---

# A stale user-level skill copy is reported with its refresh command

Recorded by the operator on 2026-11-10: a Claude Code session in this
repository loaded the Roundfix skill at version 0.1.140 from
`~/.claude/skills/roundfix` while the repository's copy and the installed
binary carried 0.1.208, seven more Roundfix-owned skills were behind in the
same way, and `roundfix skills check` passed. Claude Code loads a personal
skill in place of the project skill of the same name, and no Roundfix command
reads the user-level skill directories that `roundfix skills install --target`
writes. After this Spec, `roundfix skills check`, `roundfix doctor` and the
Update Check read those directories and name each Stale User Skill Copy with
the command that refreshes it, without writing to them and without changing
any exit code.

## Sources

One Backlog Entry, adopted by this Spec's planning change (see
`references/_index.md`):

1. [`references/2026-11-10-a-user-level-skill-copy-older-than-the-project-copy-is-loaded-unnoticed.md`](references/2026-11-10-a-user-level-skill-copy-older-than-the-project-copy-is-loaded-unnoticed.md),
   `type: fix`, recorded by the operator from the maintainer's MacBook.

At the authoring base `docs/findings/` is empty, and the other open Backlog
Entries (the push gate's suite-time budget, the catalog bump's stale digest
and the Portuguese not-applicable marker) belong to Specs 0377, 0378 and 0379.

## Prerequisites

None in `requires`. Authored against `2b561a38` (v0.137.9) and rebased on
`b233761e`, which adds the 0374 plan. Specs 0373, 0374,
0375 and 0377 deliver before this one and may change files this Spec also
changes, each in its own section:

- The Roundfix owned skill's version fields in `.agents/skills/roundfix/SKILL.md`
  and `skills/roundfix/SKILL.md` and `skills/testdata/owned-skill-versions.json`
  (0373 and 0374 raise it too). task_02 raises the version with the record
  command from the tree it runs on and never writes a number.
- `.agents/skills/roundfix/references/modes.md` and its mirror (0373 changes
  the `--follow-up` steps; this Spec only step 1 of `--update`), and
  `skills/roundfix_modes_test.go` (0373 adds a test; this Spec adds two
  phrases to `TestTheUpdateModeAppliesOnlyThePlanItShowed`).
- `.agents/skills/roundfix/references/setup.md` and its mirror,
  `docs/user-guide/commands/doctor.md`, `internal/cli/doctor.go` and
  `internal/cli/health.go` (0374 adds the `verification-capacity:` line after
  `disk-space:`; this Spec adds the `user-skills:` line after `skills:`; 0375
  may report its core budget in Doctor too).
- `internal/cli/cli.go` (`runSkillsCommand` and `commandDependencies`),
  `internal/cli/cli_test.go` (`commandEnvironmentForTest`) and
  `internal/cli/update_check.go`.
- `CONTEXT.md` (each Spec its own terms) and
  `docs/references/behavior-surfaces.json` (each Task re-records it after its
  own last guide edit).

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted; the Update Check keeps `roundfix/update-check/v1` and gains one
  optional field. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or
  network call is added; the reading is local and offline. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0386 (new) records the decision:
  "A user-level copy of a Roundfix-owned skill older than the binary is
  reported, never changed". ADR-0189 applies unchanged: "The minimum version
  of an owned skill is now the version the binary carries for it", the one
  comparison a user-level copy is held to. ADR-0364 applies unchanged: "An
  installed owned copy passes when the version it declares is equal to or
  higher than the version the running binary carries", now also read for
  user-level copies, and still "Roundfix reads no catalog over the network".
  ADR-0371 applies unchanged: the Update Check "changes nothing", and its
  report gains a line. ADR-0221, ADR-0224, ADR-0233, ADR-0241, ADR-0250,
  ADR-0256, ADR-0264 and ADR-0359, which cite ADR-0189 or ADR-0364, do not
  apply: they decide how external skills meet their lock and Setup Snapshot,
  this repository's upstream skills, a version raise at merge, a repository
  profile's skill contracts, module versions, the release's skills check and
  an owned skill's lock entry, and this Spec changes none of them. Source:
  `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BB batch 2 answers of 2026-11-10 quoted in the authorization record,
  which grant the owned skills to this Spec; the Governed Paths are the
  Roundfix skill files task_02 changes and the CLI test helper file task_01
  changes. No Makefile, workflow, Go module,
  lint or formatter change. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Goals

- A user-level copy of a Roundfix-owned skill that the installed binary has
  outgrown is named, with the command that refreshes it, the first time a
  session runs the Update Check after a release.
- `roundfix skills check` and `roundfix doctor` name the same copies when a
  person asks.
- No Roundfix command writes to a user-level skill directory unasked, and no
  repository gate fails because of a file in one developer's home.

## User Stories

1. As an Agent session starting work, I want the Update Check to tell me that
   a user-level copy of a Roundfix skill is older than the installed binary,
   so that I do not follow instructions an older release wrote.
2. As the maintainer running `roundfix skills check` or `roundfix doctor`, I
   want each stale copy named with its path, both versions and the command
   that refreshes it, so that one command fixes it.
3. As the maintainer whose skill directories link into a synchronized vault,
   I want the report to say where a refresh writes, and nothing written
   without my approval.

## Core Features

1. The reader: five user-level directories in a fixed order (the three
   `roundfix skills install` targets for `claude`, `codex` and `opencode`,
   then `~/.agents/skills` and `~/.config/opencode/skills`), the twelve
   Roundfix-owned names, symbolic links followed, a folder reached by several
   links named once under the first directory, and a copy named only when its
   declared version is below the binary's version of that skill. It writes
   nothing.
2. `roundfix skills check` prints a warning block on stderr after its usual
   output, one line per copy and one refresh line per directory, and keeps
   its standard output and exit code.
3. `roundfix doctor` prints a `user-skills: warn` line right after `skills:`
   naming the copies and their refresh commands, and keeps its exit code; the
   release plan's `skills` check does not read it.
4. The Update Check prints a `User skills:` line after the Baseline line and
   the JSON field `userSkills`, and keeps `next`.
5. Nothing prints, and the JSON field is absent, when no copy is stale.
6. The command guides, the Roundfix `setup` reference and the `--update`
   Skill Mode, which asks before running a refresh, document it; the glossary
   adds **Stale User Skill Copy**.

## User Experience

On the maintainer's MacBook on 2026-11-10 at about 21:55 BRT, a prototype of
this Spec built from `2b561a38` printed, after its usual line:

```text
roundfix: warning: user-level Roundfix skill copies older than this binary; an Agent harness may load them in place of the project copy:
  qa-gate 0.0.25 (binary 0.0.26) at /Users/<user>/.claude/skills/qa-gate -> /Users/<user>/.vault/skills/qa-gate
  roundfix 0.1.208 (binary 0.1.210) at /Users/<user>/.claude/skills/roundfix -> /Users/<user>/.vault/skills/roundfix
  roundfix 0.1.140 (binary 0.1.210) at /Users/<user>/.opencode/skills/roundfix
  refresh: roundfix skills install --target claude
  refresh: roundfix skills install --target opencode
```

(seven more `~/.opencode/skills` lines omitted here). The Claude copies the
operator refreshed at 15:50 BRT were already behind again after the two
releases of the same evening; the `~/.codex/skills`, `~/.agents/skills` and
`~/.config/opencode/skills` entries link to the same vault folders and are
not repeated.

## Non-Goals / Out of Scope

- Refreshing a user-level copy from `roundfix upgrade`, the Update Check,
  `roundfix baseline update` or any other command: the report names the
  command, and the `--update` Skill Mode runs it only when the user approves.
- Changing `roundfix skills install`, including its `--target opencode`
  directory `~/.opencode/skills`, which OpenCode does not list among its
  skill sources; that is recorded as a follow-up in the TechSpec Risks.
- Reading a harness directory that no install target writes and ADR-0386
  does not list, such as one named by `CLAUDE_CONFIG_DIR` or an
  `opencode.json` `skills` array, and any project-level directory.
- Naming a user copy at or above the binary's version, one without a
  comparable version, or one whose content differs at the same version.
- Changing the exit code of `roundfix skills check` or `roundfix doctor`, the
  `skills:` line, the release plan or the Update Check's `next`.
- A new schema version for the Update Check, or a message to adopters: they
  receive the report through their next release and the Update Check.

## Success Metrics

1. Success Metric: with a temporary home holding an older `roundfix` copy in
   each of the five directories, the reader returns five copies in the fixed
   order with the five refresh commands; an equal, a newer, an unversioned
   and a non-owned copy are not returned. At `2b561a38` no reader exists.
2. Success Metric: a skill folder reached through links from three
   directories is returned once, under the first directory, with the folder
   it reaches; the reader changes no file, mode or modification time in the
   home.
3. Success Metric: `roundfix skills check` with such a home exits `0`, prints
   its usual standard output byte for byte, and prints the warning block of
   Surface Transcript 1 on stderr; at `2b561a38` its stderr is empty.
4. Success Metric: `roundfix doctor` prints the `user-skills: warn` line
   right after `skills:` and exits as it would without it.
5. Success Metric: `roundfix update check` prints the `User skills:` line of
   Surface Transcript 2 and its JSON carries `userSkills` with `next` `none`;
   without a stale copy every existing Update Check transcript keeps its
   bytes.
6. Success Metric: every existing `skills check`, Doctor and Update Check
   test passes unchanged while the process's `HOME` holds stale copies in all
   five directories, so no test reads the host's skill directories.
7. Success Metric: the `skills`, `doctor` and `update` command guides, both
   copies of the Roundfix `setup` reference and step 1 of the `--update`
   Skill Mode document the report; `CONTEXT.md` defines **Stale User Skill
   Copy**; the skills, documentation and glossary contracts pass.

## Acceptance evidence

The outside-evidence row rests on published documentation this Spec did not
write, read through Exa on 2026-11-10.

- Claude Code, "Extend Claude with skills" (https://code.claude.com/docs/en/skills),
  section "Resolve skills that share a name": "Enterprise over personal, and
  personal over project. With `deploy` in both `~/.claude/skills/` and the
  project's `.claude/skills/`, `/deploy` runs the personal one". This is why
  a stale `~/.claude/skills` copy replaces the repository's.
- OpenCode, "Skills" (https://opencode.ai/docs/skills/): the Discovery table
  lists "Global | `~/.config/opencode/skills`" and "Global compatibility |
  `~/.claude/skills`, `~/.agents/skills`", and Precedence registers
  "`.agents/skills`, global first and then from the farthest ancestor toward
  the current directory" before "`~/.config/opencode/skills`". This is why
  the reader includes `~/.agents/skills` and `~/.config/opencode/skills`.
- Codex, "Build skills" (https://developers.openai.com/codex/skills): "`USER`
  | `$HOME/.agents/skills`" and "If two skills share the same `name`, Codex
  doesn't merge them; both can appear in skill selectors"; the Codex source
  `codex-rs/core-skills/src/loader.rs` keeps "Deprecated user skills location
  (`$CODEX_HOME/skills`), kept for backward compatibility."

## Unreachable Acceptance

- criterion: after the release that carries this Spec, `roundfix skills
  check` on the maintainer's MacBook names the user-level Roundfix skill
  copies older than the installed binary, and the refresh commands it prints
  bring them to the binary's versions.
  reason: it needs the maintainer's real home directory and its linked
  vault, which no Verification or QA row may read or write.
  satisfied-by: the operator's run on the maintainer's MacBook after the
  release, recorded session evidence.

## Validation

The decisions were checked on 2026-11-10 against these sources.

- Characterization at `2b561a38`:
  - `runSkillsCommand` (`internal/cli/cli.go`) checks only the embedded
    bundle and, inside an adopted repository, `.agents/skills`;
    `repositorySkillsCheck` (`internal/cli/doctor.go`) and
    `buildUpdateCheckReport` (`internal/cli/update_check.go`) read no
    user-level directory;
  - `skills.DefaultTargetDirs` writes `--target claude`, `codex` and
    `opencode` to `$CLAUDE_HOME`, `$CODEX_HOME` or `$OPENCODE_HOME` plus
    `/skills`, else `~/.claude/skills`, `~/.codex/skills` and
    `~/.opencode/skills`, and `--dir` overrides one target's directory;
  - `commandEnvironmentForTest` gives a test the host's home unless it sets
    one, so a reader of the real home would make `TestRunSkillsCheck` fail on
    a machine with a stale copy (it requires an empty stderr);
  - `make verify` runs `roundfix skills check`, and the release plan's
    `skills` check reads Doctor's `skills:` line through
    `repositorySkillsCheck`;
  - the maintainer's MacBook: `~/.claude/skills/<name>` link into
    `~/.vault/skills`, `~/.codex/skills`, `~/.agents/skills` and
    `~/.config/opencode/skills` link to `~/.claude/skills/<name>`, and
    `~/.opencode/skills` holds plain copies from 2026-10-26 at the versions
    the Backlog Entry lists.
- A prototype of both Tasks on a scratch copy at `2b561a38`: the four reader
  tests and three command tests passed; every `Doctor`, `SkillsCheck`,
  `UpdateCheck` and `Skill` test of `./internal/cli` and the `./skills`
  suite passed with `HOME` holding stale copies in all five directories;
  the glossary and documentation contracts passed on the prescribed text,
  and the Behavior Surface Record moved only for the three command guides.
- Secondbrain, read only: `wiki/index.md`, then `qmd query` for user-level
  skill copies and `roundfix skills install`. The top results were
  `wiki/entities/roundfix.md`, `wiki/entities/skills.md` and the mirror of
  ADR-0364; none records user-level copies or a decision about them, so the
  decision rests on the maintainer's answer and the harness documentation.
- Exa: the three published sources above, plus the issue
  "gh skill install --agent codex --scope user installs to the deprecated
  ~/.codex/skills" (https://github.com/cli/cli/issues/14153), which confirms
  that `$CODEX_HOME/skills` is still read.

## Glossary

- adds: **Stale User Skill Copy**

## Decisions

Autonomous route: no person present; decisions taken from the operator's cycle BB batch 2 brief for Spec 0376 of 2026-11-10, the maintainer's answers below and the adopted Backlog Entry.

The maintainer's answers (AskUserQuestion via `/roundfix --follow-up`,
2026-11-10 ~20:55–21:00 BRT), verbatim:

- Question "Quais Backlog Entries viram Spec agora, depois da 0374? As que
  você não marcar ficam abertas no backlog." → selected, among four, "Skills
  de usuário velhas (Recommended)" — option text: "O `doctor` e o `skills
  check` passam a apontar uma cópia de usuário de skill do Roundfix mais
  velha que a do binário, com o comando que a atualiza. Hoje cada release
  deixa as do vault velhas sem aviso."
- Question "Com que autorizações eu toco sem supervisão as seis Specs novas:
  0374 (capacidade da máquina), 0375 (orçamento de núcleos), 0376 (skills de
  usuário velhas), 0377 (portão do CI), 0378 (digest do bump de skills) e
  0379 (marcador \"não se aplica\")?" → "Mesmas do BB, teto minor
  (Recommended)" — option text: "Escrever, entregar pela fila, merge e
  release de cada Spec assim que pronta (patch nas correções, minor nas
  features; 0374 e 0375 saem juntas), teto de 2 Tasks corretivas, CI na 0377
  e Makefile na 0378 só nos arquivos que elas exigirem, eventos ao Jarvis e
  resposta ao Vortex quando a 0374 e a 0375 saírem."

Taken from those answers, the Backlog Entry and the characterization, with
no further interview:

- Report, never refresh: the answer asks for the copy and "o comando que a
  atualiza"; the directories lead through links into the maintainer's vault.
- The Update Check reports too, so the first session after a release sees it
  (the Backlog Entry's Expected names it); the `--update` Skill Mode asks
  before running a refresh.
- Five directories: the three install targets, and the two user directories
  Codex and OpenCode document that no install target writes.
- Links followed, one report per reached folder, under the first directory.
- Only a declared older version is stale, by ADR-0189's comparison.
- No exit code changes, a separate `user-skills:` Doctor line, and nothing
  printed without a stale copy, so `make verify`, the release plan and every
  existing transcript keep their behavior.
- ADR-0386 records the decision. One implementation Task, one docs Task and
  the QA gate.

## Open Questions

None.
