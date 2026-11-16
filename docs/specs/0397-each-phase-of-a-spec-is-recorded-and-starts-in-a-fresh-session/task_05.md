---
task: task_05
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Phase Record under
Roundfix Home and `roundfix phase start`, which records the phase a session
starts and prints the Boundary Notice without refusing anything; records
each prompt of a Run and of a pre-PR review; removes old day files in the
Run Retention Sweep; installs a published authoring skill without the
Declared Skill Header and declares this repository's authoring models with
their generated subagents; and states the Phase Boundary and the fix route
in the four authoring skills, the Roundfix skill and the autonomous-work
guide of the Baseline, with ADR-0407, ADR-0408, the notice in ADR-0401 and
the glossary.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head or of the Delivery
Base. No row starts Claude Code, another Agent harness, an Agent session or
an ACP adapter, calls ai-memory, reads a harness transcript, reads or writes
the real `~/.roundfix`, the real home directory or a user-level skill
directory of the host, runs the real `gh`, reaches the network, reads
another repository, or sends a request to a provider, the judge or a
reviewer. Every `go test` it runs carries an explicit `-timeout`.

Every row that runs the product runs `bin/roundfix`, the binary `make build`
made from the audited head, and never a `roundfix` found on `PATH`: an older
binary would write the Declared Skill Header into this repository's
canonical skills. No row runs `skills install` or `baseline update` with the
checkout as its repository; those commands run only in temporary
repositories and copies. Every row that runs `bin/roundfix` takes
`m="$(go env GOMODCACHE)"` and `c="$(go env GOCACHE)"` first, then sets
`HOME` and `TMPDIR` to fresh temporary directories outside the checkout with
`GOMODCACHE="$m"` and `GOCACHE="$c"`, puts a fake `gh` script that prints
`[]` first on `PATH`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`
and every Judge Key Variable for that command, reads standard input from
`/dev/null`, and removes the temporary directories afterwards. The only
environment value a row sets and records is the session identifier of the
Surface Transcripts' fixtures. The QA never writes into the checkout while a
test suite runs. The gate keeps its evidence under names the repository does
not ignore, `.txt` for command output, and links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter that needs the session's toolchain configuration
under `HOME`. Only rows that run `bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 3 and Invariants 1, 2 and 4 by running
   `go test -count=1 -timeout 30m -v` with the three `internal/phaserecord`
   tests and the six `internal/cli` tests that task_01 names in its
   Verification, recording every test and subtest by name; then the whole
   `./internal/phaserecord` package and
   `GOOS=windows go build -buildvcs=false ./internal/phaserecord`, recording
   each pass.
3. MUST verify Success Metrics 4 and 5 and Invariant 3 by running, with
   `-v`, the two `internal/phaserecord` tests, the five `internal/daemon`
   tests and the five `internal/cli` tests that task_02 names in its
   Verification, recording every test and subtest by name; then
   `-run 'Retention' ./internal/cli`, recording the pass.
4. MUST verify Success Metrics 6 and 7 and Invariants 6 and 7 by running,
   with `-v`, the `skills` tests and the four `internal/cli` tests that
   task_03 names in its Verification, recording every test and subtest by
   name; then the whole `./skills` package and
   `-run 'Skill|Doctor|UpdateCheck|Authoring' ./internal/cli`, recording
   each pass.
5. MUST verify Success Metric 8 and Invariant 8 by running, with `-v`, the
   seven `skills` tests and the two `internal/baseline` tests that task_04
   names in its Verification, recording every test and subtest by name; then
   the whole `./internal/baseline` package,
   `-run 'Baseline|Guidance|SetupManifest' ./internal/cli` and
   `-tags docscontract ./internal/docscontract`, recording each pass.
6. MUST reproduce Surface Transcript 1, Surface Transcript 2, Surface
   Transcript 3, Surface Transcript 4 and Surface Transcript 5 of
   `_techspec.md` through `bin/roundfix`, each with the fixture the TechSpec
   describes: a temporary home and, as the working directory, a temporary
   Git repository; `CLAUDE_CODE_SESSION_ID` set to
   `11111111-2222-4333-8444-555555555555` for fixtures A and B and removed
   from the command's environment with `env -u` for fixture C. Record
   standard output, standard error and the exit code of each, compared line
   by line with the transcript, and the line each of transcripts 2 to 4
   added to the day file: its keys in the order of `_techspec.md` → Data
   Models, `boundary` false in transcripts 2 and 4, and `boundary` true with
   `previous_phase` `prd` and `previous_spec` `0412-example` in transcript
   3. Record that the record's directory has mode `0700` and the day file
   `0600`, and that transcript 5 left no file and no directory.
7. MUST reproduce Surface Transcript 6 and Surface Transcript 7 through
   `bin/roundfix`. Fixture D is a temporary Git repository with no `.claude`
   directory whose `.roundfixrc.yml` declares `prd` and `techspec` with
   `model: opus` and `reasoning_effort: xhigh`, and which holds
   `skills/write-prd/SKILL.md`, a copy of this checkout's
   `skills/write-prd/SKILL.md`. Fixture E is fixture D after Surface
   Transcript 6, with the line `model: opus` added to the frontmatter of its
   `.agents/skills/write-prd/SKILL.md` and with `.claude/skills` made a
   symbolic link to `../.agents/skills`; task_03 states why the link is part
   of the fixture. Record both streams and the exit code of each, compared
   line by line with the transcript, where `<path>`, `<count>` and `<names>`
   are whatever the command prints. Also record, in fixture D after
   transcript 6: that the installed `write-prd` equals the shipped file and
   the installed `write-techspec` holds the two header lines; and that a
   second `bin/roundfix skills install --target project` changes no byte of
   the repository.
8. MUST reproduce Surface Transcript 8 through `bin/roundfix`: a temporary
   Git repository and a temporary home with no Run Database, whose
   `.roundfix/phase-record` directory holds two day files dated more than 30
   days before the day the gate runs and one of that day. Record that
   `bin/roundfix gc --dry-run` prints `  Phase Record day files removable: 2`
   right after `  Runs removable: 0`, an empty standard error and exit `0`,
   and that the three files remain; then that `bin/roundfix gc` prints
   `  Phase Record day files removed: 2` right after `  Runs removed: 0`
   and leaves only the file of the day.
9. MUST verify Invariant 5 through `bin/roundfix`: in a temporary Git
   repository with a temporary home, run `follow-up`, `follow-up --status`
   and `runs list`, and record that the home holds no `phase-record`
   directory afterwards.
10. MUST verify this repository's adoption in the checkout, changing no
    file: record that `bin/roundfix skills check`, with a fresh temporary
    home, exits `0` and prints no line that names the Authoring Model
    Declaration; that the five files `roundfix-idea-author.md`,
    `roundfix-prd-author.md`, `roundfix-techspec-author.md`,
    `roundfix-tasks-author.md` and `roundfix-researcher.md` exist under
    `.claude/agents/`, each with its
    `<!-- roundfix:generated skills.authoring.claude.<key> -->` marker, the
    model and effort `.roundfixrc.yml` declares for it, and, for the
    researcher, the tool list `Read, Grep, Glob, WebSearch, WebFetch,
    mcp__exa__web_search_exa, mcp__exa__web_fetch_exa`; that no `SKILL.md`
    under `.agents/skills/` or `skills/` of an owned skill holds a top-level
    `model` or `effort`; and that `.agents/skills/<name>` and
    `skills/<name>` hold the same bytes for every owned skill.
11. MUST reproduce the managed refresh in a temporary copy of the Delivery
    Base (the merge base of the audited head with the default branch) made
    with `git archive` and committed: record that
    `bin/roundfix baseline update --repo . --no-skills --no-history --yes`
    writes into `docs/agents/autonomous-work.md`, exactly once each after
    whitespace normalization, "At a phase boundary of a Spec, after each
    authoring phase and after a Spec is delivered and released, write the
    session handoff with the Spec, the phase that ended, the next phase,
    each open decision and what stays armed, tell the user what stays armed,
    and start no further phase in that session." and "A scheduled wake-up
    returns to the same session with its whole context, so it keeps a phase
    going and is never the way across a boundary.", beside "When no such
    mechanism is armed, say plainly that work is stopping and what
    remains."; and that a second run of the same command changes no file.
12. MUST run each sabotage below on its own temporary copy of the audited
    head made with `git archive` and committed in a fresh Git repository
    (the suite guard refuses a tree outside Git), with `-timeout` on every
    `go test`, and record that the named test fails on the named subtest.
    No sabotage regenerates or synchronizes anything:
    1. `StartPhase` takes the latest phase entry of any session:
       `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/another_session`.
    2. The boundary rule loses its `follow-up` exception:
       `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/after_follow_up`.
    3. `RemoveBefore` also removes the day file of the cutoff's date:
       `TestRemoveBeforeRemovesOnlyDayFilesBeforeTheCutoff/cutoff_date_is_kept`.
    4. `roundfix phase start` exits `1` when a boundary is reported:
       `TestPhaseStartSurfaceTranscripts/transcript_3`.
    5. `recordPromptUsage` no longer writes the prompt entry:
       `TestEachPromptLeavesOnePromptEntryWithItsScope/task`.
    6. One `daemon.Dependencies` literal of `internal/cli/implement.go`
       loses its `PhaseRecord`:
       `TestEveryEngineTheCLIBuildsRecordsItsPrompts`.
    7. The installation applies the header to every declared skill:
       `TestInstallationWithholdsTheHeaderOfAPublishedSkill/one_published`.
    8. The inspector treats a published skill as any other:
       `TestInspectAuthoringAcceptsAPublishedSkillWithoutAHeader/no_header_is_no_problem`.
    9. The words "start no further phase in this session" leave the
       `## Phase boundary` section of the canonical `write-prd` only:
       `TestTheAuthoringSkillsStateThePhaseBoundary/canonical/write-prd/no_further_phase`.
    10. The two appended sentences leave the clause in the module file
        only: `TestTheOutliveTheTurnClauseStatesThePhaseBoundary`.
13. MUST execute the tests tasks 01, 02, 03 and 04 name in their
    Verification.
14. MUST search and record:
    - both copies of the four authoring skills carry the `## Phase boundary`
      section of task_04, each with its own phase, and only `write-techspec`
      and `write-tasks` carry the fix-route sentence;
    - both copies of the Roundfix `modes`, `phase`, `setup` and `storage`
      references carry the text tasks 01 to 04 prescribe, and the reference
      index names 21 references;
    - this repository's `docs/agents/autonomous-work.md`, the formatter
      golden of the Standard TypeScript profile and the autonomous-work
      module carry the clause of API Contract 11; its id, its `mandatory`
      force and its absence of `replaces` are unchanged; and every file
      under `internal/baseline/assets/source-baselines/` is byte-identical
      to the Delivery Base;
    - `CONTEXT.md` defines **Phase Record**, **Supervisor Phase**, **Phase
      Boundary** and **Boundary Notice** and carries the changed **Declared
      Skill Header**, **Author Subagent** and **Session Handoff** with the
      text tasks 01, 03 and 04 prescribe; ADR-0401 carries one notice that
      names ADR-0408; and ADR-0407 and ADR-0408 are `accepted`;
    - the usage guide documents the record's place, both entry kinds field
      by field, the phases, the notice, the writers and the retention; the
      `phase`, `gc`, `skills`, `doctor` and `update` command guides, the
      Run Database guide and the configuration guide carry the text tasks
      01 to 03 prescribe;
    - no Go file outside `internal/cli/phase.go` and its tests reads
      `CLAUDE_CODE_SESSION_ID`, and no non-test Go file of
      `internal/phaserecord` imports a package outside the standard
      library.
15. MUST record, as outside-evidence rows, the published sources this Spec
    did not write, each with its date (read 2026-11-15) and none blocked,
    through the passages `_prd.md` → Acceptance evidence quotes, looking
    nothing up:
    - Claude Code documentation, *Manage costs effectively*: compare
      "Claude Code sends your full conversation with every request" with the
      rule of requirement 11, which ends the session at a boundary instead
      of carrying its conversation into the next phase;
    - the same page: compare "a session-scoped scheduled task fires on its
      interval even while the session is idle, sending your full context
      each time" with the wake-up sentence of requirement 11;
    - Claude Code documentation, *Monitoring*: compare "Claude Code's own
      session ID, which `/clear` replaces" with the `session` key of the
      lines of requirement 6 and with the notice of Surface Transcript 3,
      which the same session receives and another session does not.
16. MUST record the six criteria of the one entry of `_prd.md` → Unreachable
    Acceptance, one row each, as `blocked` with `declared: ` followed by
    that criterion's text, never a paraphrase.
17. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `adds:` for **Phase Record**, **Supervisor
    Phase**, **Phase Boundary** and **Boundary Notice** and `changes:` for
    **Declared Skill Header**, **Author Subagent** and **Session Handoff**,
    that tasks 01, 03 and 04 wrote those seven entries and no other, and
    that `TestGlossaryPhrasesAreDefinedTerms` passes.
18. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - `.roundfixrc.yml` changed only by gaining the `skills` section;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, `internal/cli/cli_test.go`,
      any Vendored Skill under `.agents/skills`, any file under
      `internal/baseline/assets/modules/` other than `autonomous-work.json`,
      any file under `internal/baseline/assets/source-baselines/`, any file
      under `docs/agents/` other than `autonomous-work.md` and
      `setup-context.json`, any ADR other than ADR-0401, any owned skill
      other than `roundfix`, `write-idea`, `write-prd`, `write-techspec` and
      `write-tasks`, and any Roundfix reference other than `modes`, `phase`,
      `setup` and `storage`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
19. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
20. MUST, for the Outcome for the Archive Record, run only
    `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`, or run `skills install` or
      `baseline update` with the checkout as its repository;
    - read or write `~/.roundfix`, the real home directory or a user-level
      skill directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session or an ACP adapter, or read a
      harness transcript;
    - send any message, or any request to a provider, the judge or a
      reviewer;
    - record any environment value other than the session identifier the
      fixtures set;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the four Tasks and the packages.
- [ ] Reproduce the eight Surface Transcripts, the read-only commands, the adoption and the managed refresh.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the outside sources, the six declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] `roundfix phase start` records each phase in one line a JSON reader parses, tells a session that crossed a boundary and refuses nothing; a Run and a pre-PR review leave one prompt entry per prompt and lose nothing when the record cannot be written; the Run Retention Sweep removes old day files; the eight Surface Transcripts reproduce through the built binary.
- [ ] A published authoring skill receives no header and the checks agree; this repository's declaration is applied, with its five generated subagents, and no shipped skill holds a model or an effort.
- [ ] An adopter's managed refresh carries the boundary rule; both copies of the skills, the guides, the glossary and ADR-0401 agree with ADR-0407 and ADR-0408.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-9; Core Features 1-13; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Success Metric 8; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-12; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Surface Transcript 7; Surface Transcript 8; Invariants 1-8; Declared breaks; Testing Approach; Build Order 5
- ADR-0407, ADR-0408, ADR-0401
