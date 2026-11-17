---
task: task_04
spec: 0401-the-phase-boundary-follows-only-the-release-of-a-spec
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec moves the
**Phase Boundary** from every Authoring Phase to a Spec's delivery and
release: `roundfix phase start` marks an entry and prints the
**Boundary Notice** only after a `release` of the same session, and the
`autonomous-work` Baseline clause, the four authoring skills, the Roundfix
skill's modes and phase reference, the glossary, ADR-0407 and the user
guide state the same boundary.

Every row runs on this repository's own files, temporary repositories,
temporary homes and temporary copies of the audited head. No row starts an
Agent harness, an Agent session or an ACP adapter, starts a Run, a Daemon
or a Delivery Owner, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh`, reaches the
network, sends a request to a provider, the judge or a reviewer, reads
another repository, or asserts a wall-clock duration. Every `go test` it
runs carries an explicit `-timeout`.

Every row that runs the product runs `bin/roundfix`, the binary
`make build` made from the audited head, and never a `roundfix` found on
`PATH`; requirement 6 alone also runs a binary the gate builds from the
Delivery Base. Every such row takes `m="$(go env GOMODCACHE)"` and
`c="$(go env GOCACHE)"` first, then sets `HOME` and `TMPDIR` to fresh
temporary directories outside the checkout with `GOMODCACHE="$m"` and
`GOCACHE="$c"`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and
every Judge Key Variable for that command, sets `CLAUDE_CODE_SESSION_ID`
only where a row names a session, reads standard input from `/dev/null`,
and removes the temporary directories afterwards. The QA never writes into
the checkout while a test suite runs. The gate keeps its evidence under
names the repository does not ignore, `.txt` for command output, and
links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 by running, with `-v`, the tests task_01
   names in its Verification in `./internal/phaserecord` and
   `./internal/cli`, recording every test and subtest by name, and both
   whole packages.
3. MUST reproduce Surface Transcripts 1 and 2 of `_techspec.md` through
   `bin/roundfix`, each in its own temporary Git repository and temporary
   `HOME`, with `CLAUDE_CODE_SESSION_ID` set to one fixed identifier: first
   the earlier entry the transcript names, then the transcript's command.
   Record standard output, standard error and the exit code of the
   transcript's command, compared line by line with the transcript, where
   `<session>` stands for that identifier, and the last line of the day
   file under that `HOME`, whose `boundary` is `false` for transcript 1 and
   `true` with `previous_phase` `release` and `previous_spec`
   `0412-example` for transcript 2.
4. MUST verify Success Metric 2 by searching, outside `docs/history/`,
   `CHANGELOG.md`, `.git` and this Spec's folder, every tracked file for
   "after each authoring phase", "after each Authoring Phase", "When the
   phase ends, a Phase Boundary follows", "a different phase or Spec" and
   "A different phase or Spec crosses", recording each match; the row
   passes only with none. Then record that both copies of the four
   authoring skills, of the Roundfix modes and of its phase reference, the
   Standard TypeScript `autonomous-work` golden, `docs/agents/autonomous-work.md`,
   `CONTEXT.md`, ADR-0407 and both user-guide pages carry the text tasks
   01, 02 and 03 prescribe; that each copy under `skills/` equals its
   `.agents/skills/` source; that the `## Declared model` section of each
   authoring skill, the first paragraph of each `## Phase boundary` section
   and every `### QA settlement` section are byte-identical to the Delivery
   Base's (the merge base of the audited head with the default branch).
5. MUST verify Success Metric 3 by running, with `-v`,
   `TestTheOutliveTheTurnClauseStatesThePhaseBoundary` and the whole
   `./internal/baseline` package; recording, from the embedded catalog,
   that `clause.autonomous.loop-07-outlive-the-turn` keeps its id,
   `mandatory` enforcement and no `replaces`, and that the module has the
   same clause ids as at the Delivery Base; and running
   `bin/roundfix baseline update --repo . --no-skills --no-history --format json`
   without confirmation, recording `"state":"current"` and the checkout's
   `git status --porcelain` unchanged afterwards.
6. MUST compare the audited head with the Delivery Base through the
   product. Build the Delivery Base's binary in a temporary copy made with
   `git archive` and committed in a fresh Git repository, outside the
   checkout. In one temporary Git repository and one temporary `HOME` per
   binary, record the same sequence with one session: `follow-up`,
   `prd --spec 0412-example`, `release --spec 0412-example`, `release
   --spec 0412-example`, `follow-up`. Record that both binaries print the
   same standard output and exit with the same code for every command, that
   only the head omits the notice after `prd` and both print it after the
   second `follow-up`, and that the written lines differ only in their
   `boundary`, `previous_phase` and `previous_spec` keys.
7. MUST record the outside-evidence row of `_prd.md` → Success Metric 4
   and Acceptance evidence. Read the eight phase entries, in order, from the
   table of
   [`references/2026-11-17-phase-record-boundary-measurement.md`](references/2026-11-17-phase-record-boundary-measurement.md),
   never from the operator's home directory. Replay them through
   `bin/roundfix phase start <phase> [--spec <slug>]` in one temporary Git
   repository and one temporary `HOME` with one fixed
   `CLAUDE_CODE_SESSION_ID`, giving each Spec number a slug of the form
   `<number>-replayed`, and record each command's standard error and the
   written `boundary`. The row passes when none of the eight prints a
   notice or is marked, where the table records four marks by Roundfix
   0.147.0, and when a `release --spec 0401-replayed` followed by a
   `follow-up` in the same replay prints the notice naming that release.
8. MUST run each sabotage below on its own temporary copy of the audited
   head made with `git archive` and committed in a fresh Git repository
   (the suite guard refuses a tree outside Git), with `-timeout` on every
   `go test`, and record that the named test fails on the named subtest or
   check. No sabotage regenerates or synchronizes anything:
   1. The crossing condition returns to "a different phase or Spec":
      `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/after_prd_of_the_same_spec`
      and `TestPhaseBoundaryNoticeFollowsOnlyARelease/surface_transcript_1`.
   2. The condition drops the exemption of the same release:
      `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/release_of_the_same_spec_again`.
   3. The condition ignores the release:
      `TestStartPhaseMarksABoundaryOnlyAcrossPhasesOfOneSession/after_release_follow_up`
      and `TestPhaseBoundaryNoticeFollowsOnlyARelease/surface_transcript_2`.
   4. The old sentence "At a phase boundary of a Spec, after each authoring
      phase and after a Spec is delivered and released" is restored in the
      clause guidance of `autonomous-work.json`:
      `TestTheOutliveTheTurnClauseStatesThePhaseBoundary/phase_boundary_after_delivery_and_release`.
   5. "When the phase ends, a Phase Boundary follows:" is restored in the
      canonical `write-prd` skill: the `no_old_boundary` check of
      `TestTheAuthoringSkillsStateThePhaseBoundary`.
   6. "after each Authoring Phase and" is restored in step 8 of the
      canonical modes: the `no_boundary_after_each_phase` check of
      `TestTheRoundfixSkillRecordsEachPhaseAndStatesThePhaseBoundary`.
   7. "last recorded a different phase or Spec" is restored in the
      canonical phase reference: `TestThePhaseReferenceDocumentsTheCommandAndTheRecord`.
9. MUST execute the tests tasks 01, 02 and 03 name in their Verification.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and the TechSpec declare `changes:` for **Phase Boundary** and
    **Boundary Notice**, that task_01 revised the second and task_03 the
    first, that no other entry changed, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, `.roundfixrc.yml`,
      `AGENTS.md`, `CLAUDE.md`, `docs/agents/specific-repository.md`, any
      file under `.claude/agents/`, any Baseline module other than
      `autonomous-work.json`, any template, Source Baseline or Setup
      Snapshot, any profile other than `standard-typescript-monorepo.json`,
      any Vendored Skill under `.agents/skills`, any ADR other than
      ADR-0407, `internal/cli/phase.go`, and any owned skill other than
      `roundfix`, `write-idea`, `write-prd`, `write-techspec` and
      `write-tasks`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent
    of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0401-the-phase-boundary-follows-only-the-release-of-a-spec --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write
    the report's Outcome as one paragraph that does not end with a colon.
    It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a
      Daemon or a Delivery Owner;
    - send any message, or any request to a provider, the judge or a
      reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of the three Tasks and their packages.
- [ ] Reproduce both Surface Transcripts and the comparison with the Delivery Base.
- [ ] Replay the measured phase entries in a temporary home.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the searches, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An Authoring Phase is followed by no Boundary Notice, and a session that goes on after a release still gets one.
- [ ] No skill, mode, Baseline clause, guide, glossary entry or user-guide page puts a boundary after each Authoring Phase.
- [ ] The clause keeps its identity and this repository's managed guide is current.
- [ ] The measured entries replay without a notice, and each sabotage makes its named test fail.

## Verification

- `newest="$(find 'docs/specs/0401-the-phase-boundary-follows-only-the-release-of-a-spec/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1–4; Success Metrics 1–4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Interfaces (invariants 1–5); API Contracts 1–6; Surface Transcripts 1 and 2; Testing Approach (QA); Build Order 4; Glossary
- [`references/2026-11-17-phase-record-boundary-measurement.md`](references/2026-11-17-phase-record-boundary-measurement.md), the measurement of 2026-11-17
- [`references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md`](references/2026-11-16-the-phase-boundary-ends-the-session-after-every-authoring-phase.md), the Backlog Entry of 2026-11-16 this Spec adopts
- ADR-0413, ADR-0407, ADR-0401, ADR-0408, ADR-0410, ADR-0186, ADR-0257
