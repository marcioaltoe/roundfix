---
task: task_05
spec: 0392-each-authoring-phase-runs-on-its-declared-model
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Authoring Model
Declaration to Project Config, makes `roundfix skills install --target
project` and the skills stage of `roundfix baseline update` write the Declared
Skill Header and the generated Author Subagents and Research Subagent, makes
`roundfix skills check`, Doctor and the Update Check report an unapplied
declaration and name each reported user-level copy with the command that
removes it, adds `roundfix skills remove`, and writes the delegation text into
the four authoring skills and the Roundfix skill's `modes`, `setup` and
`baseline` references, with the glossary and ADR-0386's notice.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts Claude
Code, another Agent harness, an Agent session or an ACP adapter, calls
ai-memory, reads or writes the real `~/.roundfix`, the real home directory or
a user-level skill directory of the host, runs the real `gh`, reaches the
network, reads another repository, sends a request to a provider, the judge
or a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept
at their `go env` values, unsets `CLAUDE_HOME`, `CODEX_HOME`, `OPENCODE_HOME`,
`GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and every Judge Key Variable
for that command, and removes the temporary directories afterwards. The gate
runs through the built binary only `roundfix skills --help`, `roundfix skills
install --target project`, `roundfix skills check` and `roundfix skills
remove`: `roundfix update check` reads the release listing from the network
when its cache is empty, Doctor probes the host's tools, and `roundfix
baseline update` needs an adopted repository, so those three are proved
through their named tests, which use the repository's fakes. The QA never
writes into the checkout while a test suite runs. The gate keeps its evidence
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter (`rustfmt` through rustup) that needs the
session's toolchain configuration under `HOME`. Only rows that run
`bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 7 and Invariants 2 to 8 by running, with
   `go test -count=1 -timeout 20m -v` and recording every test and subtest by
   name:
   - over `./internal/config`:
     `TestTheAuthoringModelDeclarationIsReadFromProjectConfig`,
     `TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow` and
     `TestUserConfigIgnoresTheAuthoringModelDeclaration`;
   - over `./skills`: `TestNoShippedSkillDeclaresAModelOrEffort`,
     `TestTheDeclaredHeaderIsTheShippedSkillPlusTwoLines`,
     `TestGeneratedSubagentsMatchTheirContract`,
     `TestTheResearchSubagentNeverHoldsAShellOrAWriteTool`,
     `TestInstallAppliesTheDeclarationAndASecondRunChangesNothing`,
     `TestInstallRefusesAFileItDidNotGenerate`,
     `TestInstallRemovesOnlyItsOwnLeftOverSubagents`,
     `TestUserLevelTargetsReceiveNoHeader`,
     `TestInspectAuthoringNamesEachProblem`,
     `TestReportedUserSkillCopiesNameACurrentCopyThatHidesAHeader`,
     `TestTheRoundfixSkillOffersTheRemovalOfAUserLevelCopy`,
     `TestTheAuthoringSkillsDelegateADeclaredPhase`,
     `TestTheRoundfixSkillGathersDecisionsBeforeADelegatedPhase`,
     `TestTheUpdateModeAppliesOnlyThePlanItShowed` and
     `TestEveryOwnedSkillVersionIsRecorded`;
   - over `./internal/cli`:
     `TestSkillsInstallAppliesTheAuthoringModelDeclaration`,
     `TestBaselineUpdateAppliesTheAuthoringModelDeclaration`,
     `TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted`,
     `TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit`,
     `TestDoctorReportsTheAuthoringModels`,
     `TestUpdateCheckReportsTheAuthoringModels`,
     `TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit`,
     `TestSkillsRemoveListsThenRemovesOnlyOwnedCopies` and
     `TestSkillsUsageNamesTheRemoval`.

   Then the whole `./internal/config` and `./skills` packages, the
   `./internal/cli` tests selected by
   `-run 'Skill|Doctor|UpdateCheck|BaselineUpdate|Usage'` with
   `-timeout 30m`, and `-tags docscontract ./internal/docscontract`,
   recording each pass.
3. MUST reproduce Surface Transcript 1, Surface Transcript 2, Surface
   Transcript 3, Surface Transcript 4 and Surface Transcript 5 of
   `_techspec.md` through `bin/roundfix`, each in the fixture the TechSpec
   defines for it (fixture A, then fixture B built from it; the user-level
   copy of fixture B is the installed `write-prd` folder copied to
   `<home>/.claude/skills/write-prd`), and record for each one the whole
   standard output, the whole standard error and the exit code, compared line
   by line with the transcript and differing only in its named placeholders
   (`<path>`, `<count>`, `<names>`, `<write-prd>` and `<home>`).
4. MUST verify Success Metrics 1 and 3 through `bin/roundfix`, in temporary
   Git repositories with temporary homes:
   - with the example of `_prd.md` → User Experience as the Project Config,
     `roundfix skills install --target project` exits `0`; each of the four
     installed authoring `SKILL.md` differs from the same file of a control
     installation made in a second temporary repository without the section
     only by its added `model:` and `effort:` lines, right before the line
     that closes the frontmatter; the five files under `.claude/agents/`
     exist and each carries its `roundfix:generated` marker as the first
     non-blank line after the frontmatter; and a second run leaves every
     file under `.agents/skills` and `.claude/agents` with the same bytes;
   - in the control repository nothing exists under `.claude`, and every
     file under its `.agents/skills` equals, byte for byte, the file at the
     same path under the audited head's `skills/` directory;
   - with a hand-written `.claude/agents/roundfix-prd-author.md` that has no
     marker, in a fresh repository that declares `prd`, the command exits
     `2`, prints the refusal line of API Contract 5 and leaves no
     `.agents/skills` directory and the hand-written file unchanged;
   - after `research` is removed from the section, a new run prints
     `Removed generated subagent: .claude/agents/roundfix-researcher.md`,
     deletes that file, and leaves a hand-written
     `.claude/agents/my-reviewer.md` unchanged.
5. MUST verify Success Metric 2 through `bin/roundfix`: with
   `research: { tools: [Bash] }` in the section,
   `roundfix skills install --target project` exits `2` naming
   `skills.authoring.claude.research.tools names "Bash"; expected WebSearch, WebFetch or a name that begins with mcp__`
   and writes nothing; with `tools: [WebSearch, mcp__docs__search]` the
   generated `roundfix-researcher.md` holds exactly the line
   `tools: Read, Grep, Glob, WebSearch, mcp__docs__search`.
6. MUST verify Success Metric 5 through `bin/roundfix` in fixture B's home,
   after Surface Transcript 4: `roundfix skills remove --target claude
   --apply` prints `Removed <home>/.claude/skills/write-prd` and the folder
   is gone; with `<home>/.codex/skills/roundfix` made a symbolic link to a
   folder outside it and a folder `<home>/.codex/skills/not-roundfix` beside
   it, `roundfix skills remove --target codex --apply` prints the `Unlinked`
   line, the link is gone, and the link's target folder and `not-roundfix`
   still hold their files. The Doctor and Update Check parts of Success
   Metrics 4 and 5 are recorded from requirement 2's tests.
7. MUST verify Invariant 1 with two binaries: `bin/roundfix` and a binary
   built with `go build -buildvcs=false` from a temporary copy of the
   Delivery Base (the merge base of the audited head with the default
   branch) made with `git archive`. In one temporary repository without the
   section and with an empty temporary home, built once per binary in the
   same paths, `roundfix skills install --target project` and then `roundfix
   skills check` print the same standard output, the same standard error and
   the same exit code with both binaries. Before comparing, run the Delivery
   Base binary twice and confirm it equals itself; a line that differs
   between those two runs is exempt and named. The installed trees are not
   compared between the binaries, because task 04 changes the shipped
   skills: the head's tree is compared with the head's bundle in requirement
   4.
8. MUST run each sabotage of `_techspec.md` → Testing Approach, 1 to 12, on
   its own temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails in the
   named case:
   - 1, `roundfix skills install` records the lock entries before the header
     is written (it installs without the declaration, records, then installs
     with it) →
     `TestSkillsInstallAppliesTheAuthoringModelDeclaration/lock_records_the_installed_tree`;
   - 2, the installer applies the header for a user-level target →
     `TestUserLevelTargetsReceiveNoHeader/claude`;
   - 3, the installer overwrites a file at a generated name without testing
     its marker → `TestInstallRefusesAFileItDidNotGenerate/hand_written_file`;
   - 4, the installer deletes every file under `.claude/agents` whose entry
     is not declared → `TestInstallRemovesOnlyItsOwnLeftOverSubagents` (the
     unmarked `my-reviewer.md` is gone);
   - 5, the configuration loader accepts `Bash` in `research.tools` →
     `TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/tool_Bash`
     and, with that one change alone,
     `TestTheResearchSubagentNeverHoldsAShellOrAWriteTool`;
   - 6, the loader accepts any effort string →
     `TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/effort_hihg`;
   - 7, `roundfix skills check` exits non-zero when it prints the
     declaration warning →
     `TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit/declaration`;
   - 8, a copy that hides a header is named in every directory, not only in
     Claude Code's →
     `TestReportedUserSkillCopiesNameACurrentCopyThatHidesAHeader` (the
     current `write-prd` copy in the Codex directory is named);
   - 9, `roundfix skills remove --apply` removes the folder a link leads to →
     `TestSkillsRemoveListsThenRemovesOnlyOwnedCopies/apply` (the link's
     target is gone);
   - 10, `roundfix skills remove` deletes while it lists, without `--apply`
     → `TestSkillsRemoveListsThenRemovesOnlyOwnedCopies/surface_transcript_4`
     (the copy is gone after the listing);
   - 11, `model: x` is added to the frontmatter of the copy's
     `skills/write-prd/SKILL.md` →
     `TestNoShippedSkillDeclaresAModelOrEffort`;
   - 12, the inspector reports a file Roundfix did not generate at a
     declared entry's path as `differs from the declaration`, fixable by the
     installation → `TestInspectAuthoringNamesEachProblem/conflict`.
9. MUST execute the tests tasks 01 to 04 name in their Verification.
10. MUST search and record:
    - no `SKILL.md` under `skills/` or under an owned skill of
      `.agents/skills/` holds a top-level `model` or `effort` key
      (Success Metric 6, Invariant 2);
    - `.agents/skills` and `skills` hold the same bytes for every file of
      `roundfix`, `write-idea`, `write-prd`, `write-techspec` and
      `write-tasks`;
    - both copies of the four authoring skills carry the section of API
      Contract 11 once, before `## Process`, and both copies of the Roundfix
      `modes`, `setup` and `baseline` references carry API Contract 12;
    - `CONTEXT.md` carries API Contract 13 and ADR-0386 carries API Contract
      14; ADR-0401, ADR-0402 and ADR-0405 are `accepted`;
    - `docs/user-guide/configuration.md` carries the section and the
      paragraph of API Contract 15, with its one example introduced by the
      sentence that says the names are an example;
    - `docs/references/skill-coverage.json` holds an entry for
      `command skills remove` and `docs/references/behavior-surfaces.json`
      holds its fingerprint.
11. MUST record, as outside-evidence rows, the published sources this Spec
    did not write, each with its date (read 2026-11-14) and none blocked,
    through the passages `_prd.md` → Acceptance evidence quotes, looking
    nothing up:
    - Claude Code, *Create custom subagents*, frontmatter reference: compare
      "Effort level when this subagent is active. Overrides the session
      effort level" with the `model:` and `effort:` lines requirement 4 found
      in the generated files;
    - the same page, *Available tools*: compare "The first filter removes
      these tools, even when listed in the `tools` field" with the second
      numbered item of the generated Author Subagent body, "You cannot ask
      the user.";
    - the same page, *Choose a model*: compare "The per-invocation `model`
      parameter, which Claude can pass when it starts the subagent." with
      "passing no model and no effort" in both copies of the authoring
      skills;
    - Claude Code, *Extend Claude with skills*, frontmatter reference:
      compare "Enterprise over personal, and personal over project." with
      the `hides the declared header` reason of Surface Transcript 3;
    - anthropics/claude-code issue 98898: compare "silently ignored when
      Claude invokes the same skill via the Skill tool" with the third
      bullet of the `## Declared model` section, which hands the phase to
      the Author Subagent;
    - OWASP, *LLM06:2025 Excessive Agency*: compare "Limit the extensions
      that LLM agents are allowed to call to only the minimum necessary."
      with the `tools:` line requirement 5 found.
12. MUST record the five criteria of the one entry under `_prd.md` →
    Unreachable Acceptance as five declared rows, one per criterion, each
    written `blocked (declared: <criterion>)` with the criterion's own text,
    and no other declared row.
13. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `adds:` for **Authoring Model Declaration**,
    **Authoring Phase**, **Declared Skill Header**, **Author Subagent** and
    **Research Subagent** and `changes:` for **Stale User Skill Copy**, that
    `CONTEXT.md` defines each with the text of API Contract 13, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
14. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, `internal/cli/cli_test.go`,
      any file under `internal/baseline/`, any file under `docs/agents/`,
      any Vendored Skill under `.agents/skills`, `.claude/`, `.codex/` or
      `.opencode/`, any skill other than `write-idea`, `write-prd`,
      `write-techspec`, `write-tasks` and `roundfix`, and any Roundfix
      reference other than `baseline`, `modes` and `setup`;
    - no `### QA settlement` section changed, and no recorded test was
      removed or renamed.

    This row reads Task commits, so it declares a `commit_range` input.
15. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
16. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0392-each-authoring-phase-runs-on-its-declared-model --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`, the real home directory or a user-level
      skill directory of the host;
    - run `roundfix skills remove --apply`, or `roundfix skills install`
      with a user-level target, outside a temporary home;
    - run `roundfix update check`, `roundfix doctor` or `roundfix baseline
      update` through the built binary;
    - start Claude Code or another Agent harness;
    - run the real `gh`, call ai-memory or reach the network;
    - send any message;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and the consuming suites.
- [ ] Reproduce the five Surface Transcripts and the installation rows through the built binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the outside sources, the declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A repository that declares its authoring models gets the Declared Skill Header and the generated subagents from one installation, a second installation changes nothing, a file Roundfix did not generate is never overwritten or deleted, and a repository that declares nothing prints and installs what it did before.
- [ ] The three checks warn of an unapplied declaration and keep their exit codes, each reported user-level copy is named with its removal command before its refresh command, and `roundfix skills remove` lists first, removes only owned names and never follows a link.
- [ ] Both copies of the authoring skills and of the Roundfix references carry the delegation text, the glossary and ADR-0386 agree with ADR-0401, ADR-0402 and ADR-0405, each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0392-each-authoring-phase-runs-on-its-declared-model/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-7; Core Features 1-8; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-15; Surface Transcripts 1-5; Invariants 1-8; Declared breaks; Testing Approach; Build Order 5
