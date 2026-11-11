---
spec: 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
prd: _prd.md
created: 2026-11-10
---

# A stale user-level skill copy is reported with its refresh command — Technical Spec

## Executive Summary

A new reader in the `skills` package lists five user-level skill directories,
follows symbolic links, and returns each copy of a Roundfix-owned skill whose
declared version is below the binary's, by the comparison ADR-0189 defines.
`roundfix skills check`, `roundfix doctor` and the Update Check call it
through one command dependency and print what it returns, each on its own
surface, without changing an exit code and without printing anything when no
copy is stale. The primary trade-off is reporting instead of refreshing: a
session sees the stale copy and the exact command, but a person (or the
`--update` Skill Mode, after asking) runs it, because those directories belong
to the user's harnesses and, on the maintainer's machine, lead through links
into a synchronized vault (ADR-0386).

## Project Constraints

- Identifier strategy: not applicable — no identifier or record schema is
  minted; `roundfix/update-check/v1` gains one optional field. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the reader is local and
  offline; no endpoint, credential or network call is added. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0386 (new) records the decision.
  ADR-0189 applies unchanged: the binary's version of a skill is the minimum,
  read from the embedded `SKILL.md`. ADR-0364 applies unchanged: an owned copy
  passes at or above the binary's version, offline. ADR-0371 applies
  unchanged: the Update Check stays read-only and every report exits `0`.
  ADR-0221, ADR-0224, ADR-0233, ADR-0241, ADR-0250, ADR-0256, ADR-0264 and
  ADR-0359 do not apply; `## ADR Alternatives` says why for each.
  Source: `docs/agents/domain.md`, `docs/agents/spec-routing.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle BB batch 2 answers of 2026-11-10 quoted in the authorization record;
  the Governed Paths are the Roundfix skill files task_02 changes and the CLI
  test helper file task_01 changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/modes.md`, `.agents/skills/roundfix/references/setup.md`, `internal/cli/cli_test.go`, `skills/roundfix/SKILL.md`.

## Current behavior

Measured at `2b561a38`:

1. `runSkillsCommand` `check` (`internal/cli/cli.go`) prints
   `Roundfix skill check passed: <names>` from the embedded bundle and the
   repository's `.agents/skills`; it reads no home directory.
2. `runDoctorCommand` (`internal/cli/doctor.go`) prints `skills:` from
   `repositorySkillsCheck`, which the release plan's `skills` check also
   calls (`collectReleasePlanChecks`, `internal/cli/releaseplan_checks.go`).
3. `runUpdateCommand` and `printUpdateCheckReport`
   (`internal/cli/update_check.go`) print the release line, the Baseline line
   and, when `next` is `apply`, the `Next:` line; the JSON pin of
   `TestUpdateCheckTranscripts` lists every field of the report.
4. `skills.DefaultTargetDirs` (`skills/skills.go`) maps `claude`, `codex` and
   `opencode` to `$CLAUDE_HOME`, `$CODEX_HOME` and `$OPENCODE_HOME` plus
   `/skills`, else `~/.claude/skills`, `~/.codex/skills` and
   `~/.opencode/skills`; `--dir` replaces one target's directory.
5. `commandEnvironmentForTest` (`internal/cli/cli_test.go`) keeps the host's
   home and `os.Getenv` unless a test sets them; `TestRunSkillsCheck` requires
   an empty stderr.
6. Step 1 of `## --update` in `references/modes.md` reports the check and
   stops when nothing can be applied; `TestTheUpdateModeAppliesOnlyThePlanItShowed`
   pins five steps and phrases of each.

## System Architecture

- `skills/user_copies.go` (new, package `skills`): the directory list, the
  reader and the shared phrases. It sits beside `Readiness` and
  `parseSkillFrontmatter`, which it reuses, and beside `DefaultTargetDirs`,
  whose directories it mirrors.
- `internal/cli/user_skills.go` (new): one helper that calls the reader
  through the command dependency `userSkillDirectories`, and the three
  renderers.
- `internal/cli/cli.go`, `doctor.go`, `health.go` and `update_check.go`: one
  call each.
- `internal/cli/cli_test.go`: `commandEnvironmentForTest` points the
  dependency at the test's own home only.

## Implementation Design

### Interfaces

```go
// package skills
type UserSkillDirectory struct{ Path, Refresh string }
type StaleUserSkillCopy struct {
	Skill    string `json:"skill"`
	Path     string `json:"path"`
	Resolved string `json:"resolved,omitempty"`
	Found    string `json:"found"`
	Embedded string `json:"embedded"`
	Refresh  string `json:"refresh"`
}
func UserSkillDirectories(home string, getenv func(string) string) []UserSkillDirectory
func StaleUserSkillCopies(ctx context.Context, directories []UserSkillDirectory) []StaleUserSkillCopy
func (copy StaleUserSkillCopy) Describe() string
func StaleUserSkillRefreshes(copies []StaleUserSkillCopy) []string

// package cli
// commandDependencies gains: userSkillDirectories func(string) []roundskills.UserSkillDirectory
func staleUserSkillCopies(ctx context.Context, homeDir string) []roundskills.StaleUserSkillCopy
```

### Data Models

`updateCheckReport` gains `UserSkills []roundskills.StaleUserSkillCopy` with
JSON tag `userSkills,omitempty`, declared right after `Baseline`. No stored
record changes.

### API Contracts

1. API Contract: directories. `UserSkillDirectories(home, getenv)` returns
   nil for a blank home, else five entries in this order, each `Path` and
   `Refresh`:
   - `<$CLAUDE_HOME or home/.claude>/skills`, `roundfix skills install --target claude`;
   - `<$CODEX_HOME or home/.codex>/skills`, `roundfix skills install --target codex`;
   - `<$OPENCODE_HOME or home/.opencode>/skills`, `roundfix skills install --target opencode`;
   - `home/.agents/skills`, `roundfix skills install --target codex --dir <path>`;
   - `home/.config/opencode/skills`, `roundfix skills install --target opencode --dir <path>`.

   `<path>` is single-quoted, with each `'` written `'\''`, when it holds a
   character outside letters, digits and `@%_+=:,./-`. A nil `getenv` reads
   no variable.
2. API Contract: stale copies. `StaleUserSkillCopies` walks the directories
   in order and, inside each, the names of `Names()` sorted; for each
   `<dir>/<name>` it resolves symbolic links (`filepath.EvalSymlinks`), skips
   a missing or dangling entry, a non-directory, a folder already reached
   (by its resolved path), an unreadable `SKILL.md` and frontmatter it cannot
   parse; it compares the frontmatter's `version` with the binary's version
   of that skill through `Readiness` and returns the copy only when the state
   is `below`. `Resolved` is the resolved folder when `<dir>/<name>` itself is
   a symbolic link, else empty. It returns early, with what it has, when the
   context is done. It opens nothing for writing.
3. API Contract: shared phrases. `Describe` returns
   `<skill> <found> (binary <embedded>) at <path>`, followed by
   ` -> <resolved>` when `Resolved` is set. `StaleUserSkillRefreshes` returns
   each copy's `Refresh` once, in first-seen order.
4. API Contract: skills check. After the passed or unversioned line,
   `roundfix skills check` calls `staleUserSkillCopies(ctx, environment.homeDir)`
   and, when it returns copies, writes to stderr, word for word:
   `roundfix: warning: user-level Roundfix skill copies older than this binary; an Agent harness may load them in place of the project copy:`
   then `  <Describe>` per copy and `  refresh: <command>` per refresh
   command. Its stdout and exit code do not change. A failed check prints
   nothing more.
5. API Contract: Doctor. Right after the `skills:` result, `runDoctorCommand`
   appends, only when copies exist, a result named `user-skills`
   (`HealthCheckUserSkills` in `internal/cli/health.go`), status `warn`,
   detail `user-level Roundfix skill copies older than this binary: ` plus
   each `Describe` joined by `, `, next action the refresh commands joined by
   ` && `, read with `loaded.HomeDir`. `repositorySkillsCheck` and the
   release plan do not change.
6. API Contract: Update Check. `runUpdateCommand` sets `report.UserSkills`
   from `staleUserSkillCopies(ctx, loaded.HomeDir)` inside and outside a Git
   repository and for every release state. `printUpdateCheckReport` prints,
   only when it is not empty, right after the Baseline line and before the
   `Next:` line:
   `User skills: copies older than this binary: <Describe>, ...; refresh with <command> && <command>.`
   `next` does not change.
7. API Contract: dependency. `commandDependencies.userSkillDirectories`
   defaults to `func(home) { return roundskills.UserSkillDirectories(home, os.Getenv) }`.
   `staleUserSkillCopies` returns nil for a nil dependency or a blank home.
   `commandEnvironmentForTest` replaces it, after every override, with one
   that returns nil unless the test set a home through
   `setCommandHomeDirForTest` or `setCommandEnvironmentForTest`, and then
   reads that home with no environment variable.
8. API Contract: guides. `docs/user-guide/commands/skills.md` gains, after
   the paragraph that ends "it checks only the shipped bundle.", a paragraph
   naming the five directories in order, the link and reported-once rules,
   what is not named, that nothing is written, and the stderr block of API
   Contract 4 as a `text` block with `<skill> <found> (binary <version>) at
   <path>[ -> <link target>]` and `refresh: <command>`, followed by the refresh
   forms. `docs/user-guide/commands/doctor.md` gains, after the paragraph
   ending "The release plan's `skills` check reads this same line.", the
   `user-skills: warn (...)` line as a `text` block and the sentence that
   it follows `skills:`, never changes Doctor's exit code, is absent without
   a stale copy, and is not read by the release plan.
   `docs/user-guide/commands/update.md` gains, after the paragraph ending
   "The final outside-repository line skips the preview.", the `User skills:`
   line as a `text` block and the sentence that it does not change `next`,
   and the JSON table row `userSkills` before the `next` row. Each new
   sentence names a **Stale User Skill Copy** and points to `roundfix skills
   check` for the directories.
9. API Contract: Roundfix references. The canonical
   `.agents/skills/roundfix/references/setup.md` gains, right after the
   paragraph that ends "with `roundfix baseline skills restore --profile <id>
   --skill <name>`.", one paragraph: the three commands name each Stale User
   Skill Copy under the five directories read through links, `skills check`
   warns on stderr with its exit code kept, Doctor adds `user-skills: warn`
   after `skills:`, the Update Check adds `User skills:` and `userSkills`,
   each names the refresh command, none changes a copy, and a refresh runs
   only when the user asks because it writes outside the repository and
   through the link the report names. Step 1 of `## --update` in
   `.agents/skills/roundfix/references/modes.md` gains, right after "and the
   Baseline line.", word for word:
   "When the report carries `userSkills`, show each copy with its `refresh` command and ask once, recommending to run them, before going on: they write outside the repository, through any link the copy's `resolved` names; run only the ones the user approves."
10. API Contract: glossary. `CONTEXT.md` gains, right before **Missing
    Required Skill**, word for word, with `ADR-0386` in the parentheses
    written as a Markdown link to
    `docs/adr/0386-a-user-level-copy-of-a-roundfix-owned-skill-older-than-the-binary-is-reported-never-changed.md`,
    the form of the other glossary entries:

    ```text
    **Stale User Skill Copy**:
    A copy of a Roundfix-owned skill in a user-level skill directory of an Agent harness whose `SKILL.md` declares a version below the one the installed Roundfix carries. The harness can load it in place of the repository's copy, so `roundfix skills check`, the Doctor Command and the Update Check name it with the `roundfix skills install` command that refreshes it, and none of them changes it (ADR-0386).
    _Avoid_: Global skill, personal override, outdated install
    ```

### Surface Transcripts

The fixture is a temporary home: `<home>/vault/roundfix/SKILL.md` and
`<home>/.opencode/skills/qa-gate/SKILL.md` declare `version: 0.0.1`, and
`<home>/.claude/skills/roundfix` is a symbolic link to `<home>/vault/roundfix`;
the working directory is a temporary directory outside any Git repository.
`<names>` is the bundle's names joined by `, `, `<roundfix>` and `<qa-gate>`
the binary's versions of those skills, `<vault>` the resolved
`<home>/vault/roundfix`, and `<release line>` the Update Check's first line.

1. Surface Transcript: `skills check` warns and passes.

   ```transcript
   $ roundfix skills check
   stdout:
   Roundfix skill check passed: <names>
   stderr:
   roundfix: warning: user-level Roundfix skill copies older than this binary; an Agent harness may load them in place of the project copy:
     roundfix 0.0.1 (binary <roundfix>) at <home>/.claude/skills/roundfix -> <vault>
     qa-gate 0.0.1 (binary <qa-gate>) at <home>/.opencode/skills/qa-gate
     refresh: roundfix skills install --target claude
     refresh: roundfix skills install --target opencode
   exit: 0
   ```

2. Surface Transcript: the Update Check names the copies.

   ```transcript
   $ roundfix update check
   stdout:
   <release line>
   Baseline: not checked (outside a Git repository).
   User skills: copies older than this binary: roundfix 0.0.1 (binary <roundfix>) at <home>/.claude/skills/roundfix -> <vault>, qa-gate 0.0.1 (binary <qa-gate>) at <home>/.opencode/skills/qa-gate; refresh with roundfix skills install --target claude && roundfix skills install --target opencode.
   stderr:
   exit: 0
   ```

### Invariants

1. No surface writes, creates or removes a file under a user-level skill
   directory or the folder a link there reaches; only a person, or the
   `--update` Skill Mode after the user approves, runs a refresh.
2. No exit code of `roundfix skills check`, `roundfix doctor` or `roundfix
   update check` depends on a user-level copy, so no repository gate (`make
   verify`, which runs `roundfix skills check`, or the release plan) fails
   because of a file in a developer's home.
3. Without a stale copy every surface prints the bytes it printed at
   `2b561a38`, and the JSON has no `userSkills` key.
4. Doctor's `skills:` line and the release plan's `skills` check do not
   read the user-level directories.
5. No test in `internal/cli` reads the host's user-level skill directories.

### Declared breaks

- `roundfix skills check`, `roundfix doctor` and `roundfix update check` print
  more on a machine with a stale copy, as API Contracts 4 to 6 say.
- `skills/roundfix_modes_test.go`'s `TestTheUpdateModeAppliesOnlyThePlanItShowed`
  requires two more phrases in step 1.
- No existing test, golden or fixture changes otherwise.

## Coverage Map

- Goal 1, User Story 1, Core Feature 4, Success Metric 5 → API Contract 6,
  Surface Transcript 2, `staleUserSkillCopies`.
- Goal 2, User Story 2, Core Features 2 and 3, Success Metrics 3 and 4 → API
  Contracts 4 and 5, Surface Transcript 1.
- Goal 3 (nothing written unasked) → API Contract 2: `StaleUserSkillCopies`
  opens no file for writing, and `TestStaleUserSkillCopiesWritesNothing`
  proves it; the `--update` Skill Mode asks before a refresh (API Contract 9).
- Goal 3 (no repository gate fails) → API Contracts 4 to 6 keep every exit
  code, so `make verify`'s `roundfix skills check` still exits `0`, and
  `repositorySkillsCheck`, which the release plan reads, is unchanged
  (Invariants 2 and 4).
- User Story 3, Core Features 1 and 5, Success Metrics 1, 2 and 6 → API
  Contracts 1, 2, 3 and 7, Invariants 1, 3 and 5.
- Core Feature 6, Success Metric 7 → API Contracts 8, 9 and 10.

Deviations from the PRD: None.

## Absence Evidence

1. User-level skill reader — command: `rg -n 'UserHomeDir|DefaultTargetDirs|\.claude/skills|\.agents/skills' skills internal/cli --type go -g '!*_test.go'`; result: at `2b561a38` the matches are the install path (`DefaultTargetDirs` and `targetDirs` in `skills/skills.go`), the command home (`os.UserHomeDir` in `internal/cli/cli.go`), the `skills` usage text, and readers of the repository's own `.agents/skills` (`skills/repository.go`, `internal/cli/review.go`, `internal/cli/baseline_update.go`, `internal/cli/required_skills.go`, `internal/cli/skills_check_repository.go`); none reads a user-level skill directory.
2. `user-skills` Doctor result — command: `rg -n 'user-skills|HealthCheckUserSkills|User skills:|userSkills' internal skills docs/user-guide .agents/skills/roundfix`; result: no match at `2b561a38`.

## ADR Alternatives

- ADR-0386: the design is the decision; it takes none of the five options the
  record rejects (refresh on upgrade or in the Update Check, failing a check,
  the `skills:` line, only the install directories, a content comparison).
- ADR-0189: it rejected a fixed literal minimum and a second list of
  versions; the reader compares with the version the embedded `SKILL.md`
  declares, through the same `Readiness`.
- ADR-0364: it rejected comparing an owned skill with the catalog over the
  network; the reader is offline and compares with the binary's copy only.
- ADR-0371: it rejected applying an update from the check itself; the Update
  Check only reports user-level copies, and the `--update` Skill Mode asks
  before running a refresh.

- ADR-0221: not applicable — it held a required external skill to its Setup
  Snapshot; the reader reads only Roundfix-owned names in user-level
  directories and no external skill.
- ADR-0224: not applicable — it decides how this repository holds its own
  upstream skills; no Vendored Skill or lock entry changes.
- ADR-0233: not applicable — it decides when a skill version raise is
  regenerated at merge; task_02 raises the version with the record command
  as that decision already requires.
- ADR-0241: not applicable — it gives a repository profile its skill
  contracts; no profile or Setup Snapshot is read or changed.
- ADR-0250: not applicable — it chooses Baseline module versions and the
  coverage record's platforms; no module or coverage record changes.
- ADR-0256: not applicable — it makes a release wait for the skills that
  describe a changed surface; the release plan and its checks are unchanged,
  and task_02 updates the skills that describe the changed surfaces.
- ADR-0264: not applicable — it decides when Roundfix may replace a vendored
  skill; the reader replaces nothing and reads no vendored skill.
- ADR-0359: not applicable — it names a stale owned lock entry in Doctor's
  `skills:` line; that line and the lock are unchanged, and the new
  `user-skills:` line reads no lock.

## Integration Points

- Claude Code, Codex and OpenCode read the same directories; Roundfix only
  reads them. A refresh is `roundfix skills install`, unchanged.

## Testing Approach

New tests (each calls `t.Parallel()` first):

- `skills/user_copies_test.go`:
  `TestStaleUserSkillCopiesReportsAnOlderCopyInEveryUserDirectory` (one
  `roundfix` copy at `0.0.1` in each of the five directories; in the Claude
  directory also `qa-gate` at the binary's version, `write-tasks` one patch
  above it, `write-prd` without a version and a non-owned `not-owned` at
  `0.0.1`; exactly the five `roundfix` copies with the five refresh commands
  return), `TestStaleUserSkillCopiesFollowsLinksAndReportsAFolderOnce` (a
  vault folder linked from the Claude and Codex directories, and
  `~/.agents/skills` itself a link to the Claude directory: one copy, under
  the Claude directory, with `Resolved` and the `->` description),
  `TestStaleUserSkillCopiesWritesNothing` (path, mode, modification time and
  bytes of every file under the home equal before and after) and
  `TestUserSkillDirectoriesFollowTheInstallTargetVariables` (the three
  variables set, a home whose path holds a space, the quoted `--dir`).
- `internal/cli/user_skills_test.go`, all through `runCLI` with
  `setCommandHomeDirForTest` and the fixture of Surface Transcripts:
  `TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit` (Surface
  Transcript 1 whole), `TestDoctorReportsStaleUserSkillCopiesWithoutFailing`
  (the line right after `skills:` word for word, exit `0`, through
  `withDoctorFakeLoadedAndReadiness` with `HomeDir` set to the home) and
  `TestUpdateCheckNamesStaleUserSkillCopies` (a development version, Surface
  Transcript 2 whole with its release line, and the JSON `userSkills` equal
  to the expected copies with `next` `none`).
- Isolation: every existing test whose name matches `Doctor`, `SkillsCheck`,
  `UpdateCheck` or `Skill` in `./internal/cli` passes with the process's
  `HOME` pointing at a temporary home that holds a `roundfix` copy at
  `0.0.1` in all five directories.

Sabotages, each on a temporary copy of the audited head made with `git
archive`, each observed by the named test on the named case (all run by the
author on 2026-11-10 against the prototype):

1. Drop the version filter in `StaleUserSkillCopies` (return every readable
   copy) → `TestStaleUserSkillCopiesReportsAnOlderCopyInEveryUserDirectory`
   (the equal `qa-gate` and newer `write-tasks` copies appear).
2. Drop the reached-folder set → `TestStaleUserSkillCopiesFollowsLinksAndReportsAFolderOnce`
   (three copies instead of one).
3. Skip an entry that is a symbolic link instead of following it → the same
   test (no copy).
4. Point the fourth directory at `home/.codex/skills` instead of
   `home/.agents/skills` → `TestStaleUserSkillCopiesReportsAnOlderCopyInEveryUserDirectory`.
5. Write the `SKILL.md` before reading it → `TestStaleUserSkillCopiesWritesNothing`.
6. Drop the quoting of a `--dir` path → `TestUserSkillDirectoriesFollowTheInstallTargetVariables`.
7. Remove the warning call from `skills check` → `TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit`.
8. Make `skills check` exit `1` when copies exist → the same test.
9. Pass no copies to the Doctor result → `TestDoctorReportsStaleUserSkillCopiesWithoutFailing`.
10. Make the Doctor result `failed` → the same test (exit `1`).
11. Drop `report.UserSkills` in `runUpdateCommand` → `TestUpdateCheckNamesStaleUserSkillCopies`.
12. Make `commandEnvironmentForTest` keep the production reader when no home
    is set → `TestRunSkillsCheck` fails under the five-directory stale
    `HOME` (its stderr is not empty).

## Build Order

1. The reader, the three surfaces, the test isolation and their tests
   (task_01).
2. The command guides, the Roundfix `setup` reference, step 1 of the
   `--update` Skill Mode and its contract test, the glossary, the skill
   version and the Behavior Surface Record (task_02) (depends on: 1).
3. The QA gate (task_03) (depends on: 1, 2).

No other Spec is required. Specs 0373, 0374, 0375 and 0377 deliver first and
edit other sections of the shared files `_prd.md` → Prerequisites names; each
Task works on the tree it runs on, cites functions rather than lines, and
raises the skill version with the record command.

## Risks & Considerations

- A refresh writes through a link into whatever folder it reaches, possibly
  a synchronized vault or another checkout; the report names that folder,
  and no command runs a refresh unasked.
- A user directory that is a symbolic link to the repository's own
  `.agents/skills` is read like any other; a stale result then also shows in
  Doctor's `skills:` line.
- `roundfix skills install --target opencode` writes `~/.opencode/skills`,
  which OpenCode's documentation does not list; the maintainer's machine holds
  0.1.140 copies there that no harness loads. The reader still names them,
  because the install target wrote them and one command refreshes them;
  whether that target should move to `~/.config/opencode/skills` is a
  follow-up for a Backlog Entry, not this Spec.
- Reading twelve names in five directories costs at most sixty `stat` calls
  and twelve file reads per surface; the Update Check runs once per session.

## Glossary

- adds: **Stale User Skill Copy**

## Decisions

- Report, never refresh; ask in the `--update` Skill Mode. See ADR-0386.
- Five directories, links followed, one report per reached folder. See
  ADR-0386.
- A copy is stale only below the binary's declared version, through
  `Readiness`. See ADR-0189 and ADR-0364.
- No exit code changes; a separate `user-skills:` Doctor line; nothing
  printed without a stale copy.
- The reader lives in package `skills`, beside the comparison and the
  install targets it mirrors; the CLI reaches it through one command
  dependency that tests point at their own home.
