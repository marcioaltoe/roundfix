---
task: task_01
spec: 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
status: pending
type: backend
complexity: high
---

# Task 01: skills check, Doctor and the Update Check name a Stale User Skill Copy and write nothing

## Overview

A reader in the `skills` package lists five user-level skill directories,
follows symbolic links and returns each copy of a Roundfix-owned skill whose
declared version is below the binary's. `roundfix skills check` warns on
stderr, `roundfix doctor` adds a `user-skills: warn` line after `skills:`, and
the Update Check adds a `User skills:` line and the JSON field `userSkills`;
no exit code changes, nothing prints without a stale copy, and no test reads
the host's skill directories. This answers the Backlog Entry recorded on
2026-11-10, "A user-level skill copy older than the project copy is loaded
unnoticed".

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree
   `roundfix skills check`, `roundfix doctor` and `roundfix update check`
   read no user-level skill directory, and `TestRunSkillsCheck` requires an
   empty stderr while `commandEnvironmentForTest` keeps the host's home.
2. MUST create `skills/user_copies.go` with the types and functions of
   `_techspec.md` → Interfaces, behaving as API Contracts 1, 2 and 3 say:
   the five directories in their order with their refresh commands and the
   quoting rule; links resolved with `filepath.EvalSymlinks`, a reached
   folder skipped by its resolved path, the frontmatter `version` compared
   with the binary's version of that skill through `Readiness`, only state
   `below` returned, `Resolved` set only when `<dir>/<name>` itself is a
   symbolic link; `Describe` and `StaleUserSkillRefreshes` word for word. It
   MUST open no file for writing, create nothing and remove nothing, and it
   MUST reuse `parseSkillFrontmatter`, `Readiness` and the embedded versions,
   never a second version list.
3. MUST add the `userSkillDirectories` field to `commandDependencies` in
   `internal/cli/cli.go`, defaulting as API Contract 7 says, and create
   `internal/cli/user_skills.go` with `staleUserSkillCopies` and the three
   renderers, so that:
   - `runSkillsCommand` `check` writes, only when copies exist and after its
     passed or unversioned line, the stderr block of API Contract 4, whose
     first line is word for word
     `roundfix: warning: user-level Roundfix skill copies older than this binary; an Agent harness may load them in place of the project copy:`,
     then `  <Describe>` per copy and `  refresh: <command>` per command, and
     keeps its standard output and exit code;
   - `runDoctorCommand` appends right after the `skills:` result, only when
     copies exist, the `user-skills` result of API Contract 5
     (`HealthCheckUserSkills = "user-skills"` in `internal/cli/health.go`,
     status `warn`, detail beginning
     `user-level Roundfix skill copies older than this binary: `, next action
     the refresh commands joined by ` && `), read with `loaded.HomeDir`, and
     leaves `repositorySkillsCheck` unchanged;
   - `runUpdateCommand` sets `UserSkills` (JSON `userSkills,omitempty`,
     declared right after `Baseline` in `updateCheckReport`) from
     `loaded.HomeDir`, and `printUpdateCheckReport` prints, only when it is
     not empty, right after the Baseline line and before the `Next:` line,
     `User skills: copies older than this binary: <Describe>, ...; refresh with <command> && <command>.`
     without changing `next`.
4. MUST change `commandEnvironmentForTest` in `internal/cli/cli_test.go` as
   API Contract 7 says, after every override is applied: the dependency
   returns nil unless the test set a home with `setCommandHomeDirForTest` or
   `setCommandEnvironmentForTest`, and then reads the home it is given with no
   environment variable. No existing test changes otherwise.
5. MUST create `skills/user_copies_test.go` with the four tests of
   `_techspec.md` → Testing Approach, each calling `t.Parallel()` first:
   `TestStaleUserSkillCopiesReportsAnOlderCopyInEveryUserDirectory`,
   `TestStaleUserSkillCopiesFollowsLinksAndReportsAFolderOnce`,
   `TestStaleUserSkillCopiesWritesNothing` and
   `TestUserSkillDirectoriesFollowTheInstallTargetVariables`, with the inputs
   that section names (an equal, a one-patch-newer, an unversioned and a
   non-owned copy; a folder linked from two directories and a skills
   directory that is itself a link; a home path holding a space). The
   expected versions read the embedded ones, never a literal of a shipped
   version.
6. MUST create `internal/cli/user_skills_test.go` with
   `TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit`,
   `TestDoctorReportsStaleUserSkillCopiesWithoutFailing` and
   `TestUpdateCheckNamesStaleUserSkillCopies`, each calling `t.Parallel()`
   first and using the fixture of `_techspec.md` → Surface Transcripts with
   `setCommandHomeDirForTest`. The first asserts Surface Transcript 1 whole
   (standard output, every stderr line in order, exit `0`); the second the
   line right after `skills:` word for word and exit `0`, through
   `withDoctorFakeLoadedAndReadiness` with `HomeDir` set; the third, with a
   development version, Surface Transcript 2 whole and the JSON `userSkills`
   equal to the two expected copies with `next` `none`.
7. MUST NOT change `repositorySkillsCheck`, the release plan, any exit code,
   `roundfix skills install`, `skills.DefaultTargetDirs`, any existing test
   or golden, or any file under a real home directory; no test reads
   `os.UserHomeDir`, `CLAUDE_HOME`, `CODEX_HOME` or `OPENCODE_HOME` of the
   host.

## Subtasks

- [ ] Record the characterization.
- [ ] Add the reader and its four tests.
- [ ] Wire the three surfaces through the command dependency, with the test isolation.
- [ ] Add the three command tests.

## Acceptance Criteria

- [ ] An older copy in any of the five directories is named once with its refresh command; equal, newer, unversioned and non-owned copies are not.
- [ ] `skills check`, Doctor and the Update Check print the prescribed lines and keep their exit codes and `next`.
- [ ] Every existing `Doctor`, `SkillsCheck`, `UpdateCheck` and `Skill` test passes while the process's `HOME` holds stale copies.

## Context

- instruction: `docs/adr/0386-a-user-level-copy-of-a-roundfix-owned-skill-older-than-the-binary-is-reported-never-changed.md`
- instruction: `skills/readiness.go`
- instruction: `skills/skills.go`
- instruction: `internal/cli/doctor_test.go`
- instruction: `internal/cli/update_check_test.go`
- creates: `skills/user_copies.go`
- creates: `skills/user_copies_test.go`
- creates: `internal/cli/user_skills.go`
- creates: `internal/cli/user_skills_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/update_check.go`

## Consumers

- The reader's exported types and functions in package `skills` — consumers: `./skills`, `./internal/cli`
- The `skills check` warning, the `user-skills:` Doctor line and the Update Check's `User skills:` line and `userSkills` field — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestStaleUserSkillCopiesReportsAnOlderCopyInEveryUserDirectory|TestStaleUserSkillCopiesFollowsLinksAndReportsAFolderOnce|TestStaleUserSkillCopiesWritesNothing|TestUserSkillDirectoriesFollowTheInstallTargetVariables)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestStaleUserSkillCopiesReportsAnOlderCopyInEveryUserDirectory TestStaleUserSkillCopiesFollowsLinksAndReportsAFolderOnce TestStaleUserSkillCopiesWritesNothing TestUserSkillDirectoriesFollowTheInstallTargetVariables; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the four tests do not exist, so their pass lines are missing; afterwards they pass.

- `h="$(mktemp -d)" || exit 1; trap 'rm -rf "$h"' EXIT; for d in .claude/skills .codex/skills .opencode/skills .agents/skills .config/opencode/skills; do mkdir -p "$h/$d/roundfix" && printf '%s\n' '---' 'name: roundfix' 'version: 0.0.1' '---' > "$h/$d/roundfix/SKILL.md" || exit 1; done; gocache="$(go env GOCACHE)"; gomodcache="$(go env GOMODCACHE)"; out="$(HOME="$h" GOCACHE="$gocache" GOMODCACHE="$gomodcache" go test -count=1 -timeout 40m -v -run 'Doctor|SkillsCheck|UpdateCheck|Skill' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit TestDoctorReportsStaleUserSkillCopiesWithoutFailing TestUpdateCheckNamesStaleUserSkillCopies TestRunSkillsCheck; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the three command tests do not exist; afterwards they pass, and every existing `Doctor`, `SkillsCheck`, `UpdateCheck` and `Skill` test of `./internal/cli`, `TestRunSkillsCheck` included, passes while the process's `HOME` holds a stale `roundfix` copy in all five directories.

## References

- `_prd.md` → Goals; User Story 1; User Story 2; User Story 3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6
- `_techspec.md` → Current behavior 1-5; Interfaces; Data Models; API Contracts 1-7; Surface Transcript 1; Surface Transcript 2; Invariants 1-5; Declared breaks; Testing Approach; Build Order 1
- ADR-0386
- [The adopted Backlog Entry](references/2026-11-10-a-user-level-skill-copy-older-than-the-project-copy-is-loaded-unnoticed.md)
