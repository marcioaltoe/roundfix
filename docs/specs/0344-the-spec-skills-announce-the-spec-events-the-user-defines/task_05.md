---
task: task_05
spec: 0344-the-spec-skills-announce-the-spec-events-the-user-defines
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec gives the `write-prd`,
`write-techspec`, `write-tasks` and `archive-spec` skills and the Roundfix
skill's implement reference a generic Spec Event Announcement step that
exists only when the user's own instructions define its destination and
format, adds **Spec Event** and **Spec Event Announcement** to the glossary,
describes the step in the context-driven development guide, and adds a
skill-contract test that holds the wording and keeps every owned skill free
of a personal assistant's or tracker's name.

Every row runs on this repository's own files and temporary copies of the
audited head. No row starts an Agent session or an ACP adapter, reads or
writes the real `~/.roundfix`, reads another repository or the user's own
instructions, reaches the network, sends a message or a request to a
provider, the judge or a reviewer, or reads a credential. Every `go test` it
runs carries an explicit `-timeout`. Every row that runs `bin/roundfix` sets
`HOME` and `TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. No row compares bytes with
a commit.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2 and Invariants 3 and 4 by running
   `go test -count=1 -timeout 20m -v` with
   `TestTheSpecSkillsAnnounceOnlyWhatTheUserDefines` and
   `TestNoOwnedSkillNamesAPersonalAssistant` over `./skills`, recording each
   pass and the ten subtests of the first, then `go test -count=1 -timeout 20m ./skills`
   and `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract`.
3. MUST run each of the four sabotages `_techspec.md` → Testing Approach
   lists on its own temporary copy of the audited head made with
   `git archive`, with `-timeout` on every `go test`, and record that the
   named test fails with the path and message the TechSpec names, from the
   input that makes the mutation observable (a sentence removed from one
   canonical file, a forbidden word in a non-`SKILL.md` canonical file, a
   mixed-case forbidden word in the mirror of a skill this Spec does not
   edit, a duplicated heading in a mirror).
4. MUST execute the tests task_01 to task_04 name in their Verification.
5. MUST search, case-insensitively, every file under `.agents/skills/` and
   `skills/` of each owned skill, `docs/user-guide/`, `docs/agents/` and
   `CONTEXT.md`, and record that none contains `jarvis`, `todoist`, `gt3`,
   or `spec para o quadro`; that the five surfaces and
   their mirrors each carry one `## Spec Event Announcement` section; that
   the guide's table names the owner of all seven Spec Events; and that the
   Roundfix skill's `references/modes.md` still says a release-announcement
   step exists only when the user's own instructions define one.
6. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the Jarvis report quoted in
     `references/2026-11-03-the-spec-skills-announce-each-spec-event.md`
     (2026-11-03): compare offline its seven events, the moment of each and
     the five skills it names with the owners and moments in API Contracts
     2 to 6 and the guide's table, and record each event's mapping;
   - the Claude Code memory documentation passage `_prd.md` → Acceptance
     evidence quotes ("User instructions | `~/.claude/CLAUDE.md` | Personal
     preferences for all projects"): compare it with the generic paragraph's
     "the user's own instructions, global or for this repository";
   - the Codex guide passage `_prd.md` → Acceptance evidence quotes ("Create
     persistent defaults in your Codex home directory so every repository
     inherits your working agreements."): the same comparison.
7. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
8. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `adds: **Spec Event**` and
   `adds: **Spec Event Announcement**`, that `CONTEXT.md` defines both with
   line-anchored definitions, and that the Task commit range changes no other
   glossary entry.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every changed Governed Path is listed in `_authorization.md` `paths`;
   - the Task commit range changes none of these: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
     `CHANGELOG.md`, `internal/app/version.go`, `dist/`,
     `docs/references/skill-coverage.json`, any Go file other than
     `skills/spec_event_announcement_test.go`, any Baseline module or guide
     under `docs/agents/`, any existing Archive Record, the Roundfix skill's
     `references/modes.md` and `references/follow-up.md`,
     `skills/roundfix_modes_test.go`, and any skill other than the
     `write-prd`, `write-techspec`, `write-tasks`, `archive-spec` and
     Roundfix files the Tasks declare;
   - no `### QA settlement` section changed.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0344-the-spec-skills-announce-the-spec-events-the-user-defines --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any message or request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and the suites.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The five surfaces and their mirrors carry the generic step and their events, in the tests and by search.
- [ ] No owned skill names the reported assistant, tracker, channel or block, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0344-the-spec-skills-announce-the-spec-events-the-user-defines/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-9; Invariants 1-5; Declared breaks; Testing Approach; Build Order 5
