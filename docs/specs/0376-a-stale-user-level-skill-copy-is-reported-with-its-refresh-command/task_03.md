---
task: task_03
spec: 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
status: completed
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds a reader of five
user-level skill directories that names each Stale User Skill Copy, prints it
in `roundfix skills check`, `roundfix doctor` and the Update Check without
changing an exit code or writing a file, documents it in three command
guides, the Roundfix `setup` reference and the `--update` Skill Mode, and
records ADR-0386 and the glossary term.

Every row runs on this repository's own files, temporary homes and temporary
copies of the audited head. No row starts an Agent session or an ACP adapter,
calls ai-memory, reads or writes the real `~/.roundfix`, the real home
directory or any real `~/.claude`, `~/.codex`, `~/.opencode`, `~/.agents` or
`~/.config/opencode` directory, runs the real `gh`, reaches the network,
reads another repository, sends a request to a provider, the judge or a
reviewer, or reads a credential. Every `go test` it runs carries an explicit
`-timeout`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their
`go env` values, unsets `CLAUDE_HOME`, `CODEX_HOME`, `OPENCODE_HOME`,
`GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and every Judge Key Variable
for that command, puts a fake `gh` script that prints `[]` first on `PATH`,
and removes the temporary directories afterwards. The QA never writes into
the checkout while a test suite runs. The gate keeps its evidence under names
the repository does not ignore, `.txt` for command output, and links only
files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter that needs the session's toolchain configuration
under `HOME`. Only rows that run `bin/roundfix`, and the isolation row of
requirement 3, use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 4 by running
   `go test -count=1 -timeout 20m -v` with
   `TestStaleUserSkillCopiesReportsAnOlderCopyInEveryUserDirectory`,
   `TestStaleUserSkillCopiesFollowsLinksAndReportsAFolderOnce`,
   `TestStaleUserSkillCopiesWritesNothing`,
   `TestUserSkillDirectoriesFollowTheInstallTargetVariables` and
   `TestTheUpdateModeAppliesOnlyThePlanItShowed` over `./skills`, and
   `TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit`,
   `TestDoctorReportsStaleUserSkillCopiesWithoutFailing`,
   `TestUpdateCheckNamesStaleUserSkillCopies`, `TestRunSkillsCheck` and
   `TestUpdateCheckTranscripts` over `./internal/cli`, recording every row by
   name; then the whole `./internal/cli` and `./skills` packages and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST verify Success Metric 6 and Invariant 5: in a fresh temporary home
   holding `roundfix/SKILL.md` declaring `version: 0.0.1` under each of
   `.claude/skills`, `.codex/skills`, `.opencode/skills`, `.agents/skills`
   and `.config/opencode/skills`, run
   `go test -count=1 -timeout 40m -run 'Doctor|SkillsCheck|UpdateCheck|Skill' ./internal/cli`
   with `HOME` set to that home and `GOCACHE` and `GOMODCACHE` kept at their
   `go env` values, record that it passes, and remove the home afterwards.
4. MUST reproduce Surface Transcript 1 and Surface Transcript 2 of
   `_techspec.md` through the built binary, with `HOME` set to a temporary
   home built as that section's fixture says and the working directory a
   temporary directory outside any Git repository, comparing standard
   output, standard error and exit code with the transcript's conventions.
   Then run `bin/roundfix update check --format json` with the same home and
   record that `userSkills` holds the two copies with their `skill`, `path`,
   `resolved` (only for the linked `roundfix`), `found`, `embedded` and
   `refresh`, and that `next` is the value the same command gives with an
   empty temporary home. Then run `bin/roundfix skills check` and
   `bin/roundfix update check` with an empty temporary home and record that
   neither prints a warning or a `User skills:` line and that the JSON has no
   `userSkills` key.
5. MUST record that the fixture home's files, modes and modification times
   are byte-identical before and after the runs of requirement 4.
6. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails on
   the named case: sabotages 1 to 12 as listed there. Sabotage 12 runs
   `TestRunSkillsCheck` with `HOME` set to the five-directory home of
   requirement 3.
7. MUST execute the tests tasks 01 and 02 name in their Verification.
8. MUST search and record:
   - the `skills`, `doctor` and `update` command guides carry the `text`
     blocks of task_02's requirements 1 to 3 and the `userSkills` JSON row;
   - both `setup` reference copies carry the paragraph of API Contract 9,
     and both `modes` copies its `--update` sentence;
   - `.agents/skills/roundfix` and `skills/roundfix` hold the same `setup`
     and `modes` bytes;
   - `CONTEXT.md` carries API Contract 10, and ADR-0386 is `accepted`;
   - `repositorySkillsCheck` and `collectReleasePlanChecks` read no
     user-level directory, so the release plan's `skills` check is unchanged.
9. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date (read 2026-11-10) and none blocked, through the
   passages `_prd.md` → Acceptance evidence quotes:
   - Claude Code, "Extend Claude with skills": compare "personal over
     project. With `deploy` in both `~/.claude/skills/` and the project's
     `.claude/skills/`, `/deploy` runs the personal one" with the reader's
     first directory;
   - OpenCode, "Skills": compare its Discovery table's
     `~/.config/opencode/skills`, `~/.claude/skills` and `~/.agents/skills`
     with the reader's directories;
   - Codex, "Build skills": compare "`USER` | `$HOME/.agents/skills`" and
     the source comment "Deprecated user skills location
     (`$CODEX_HOME/skills`), kept for backward compatibility." with the
     reader's second and fourth directories.
10. MUST record the PRD's one Unreachable Acceptance entry, the operator's
    run on the maintainer's MacBook after the release, as one row counted in
    `rows_blocked_declared`, with its `satisfied-by`.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `adds:` for **Stale User Skill Copy**, that
    `CONTEXT.md` defines it with the text of API Contract 10, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, any file under
      `internal/baseline/`, any Vendored Skill under `.agents/skills`, any
      skill other than the Roundfix skill, any file under `.claude/`,
      `.codex/` or `.opencode/`, `internal/cli/releaseplan_checks.go`;
    - `skills/skills.go` and `skills/readiness.go` keep their bytes from the
      Delivery Base;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run a refresh command the report prints, or `roundfix skills install`
      with any user-level target;
    - run the real `gh`, call ai-memory or reach the network;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the isolation row, the transcripts and the JSON report.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A stale user-level copy is named on all three surfaces with its refresh command, and no exit code, `next` or file changes.
- [ ] No test reads the host's skill directories, and nothing prints without a stale copy.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-10; Surface Transcript 1; Surface Transcript 2; Invariants 1-5; Declared breaks; Testing Approach; Build Order 3
