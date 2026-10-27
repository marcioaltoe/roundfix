---
task: task_03
spec: 0325-the-follow-up-mode-consumes-the-handoff-and-every-message
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Roundfix skill's
`--follow-up` Skill Mode take the project's Session Handoff before its
inventory and hand its next steps and open questions to `roundfix follow-up`
as `handoff:<handoff id>:<n>` message entries, leave no handed-in message
pending when the turn ends, also when the mode stops early, and keep
`--status` and `--unattended` from accepting or consuming anything. A
contract test pins the text in the canonical skill and its mirror, and
`CONTEXT.md` gains **Session Handoff**.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row calls an ai-memory tool or
reaches an ai-memory server, reads or writes the real `~/.roundfix`, starts
an Agent session, reaches the network, sends a request to a provider, the
judge or a reviewer, reads a credential, or reads another repository. Every
`go test` it runs carries an explicit `-timeout`, and every row that runs a
built `roundfix` unsets every Judge Key Variable and `TYPESAFE_API_KEY`,
sets `HOME` and `TMPDIR` to fresh temporary directories, because the CLI
loads `$HOME/.roundfix/config.yml`, and puts first on `PATH` a temporary
directory holding a `gh` stub that prints `gh stub` to standard error and
exits `1`, so the command's Pull Request source never reaches GitHub.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1-4 and Invariants 1, 2 and 4 by running
   `TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage`,
   `TestTheFollowUpModeRunsItsStepsInOrder`,
   `TestTheUnattendedModeAsksWithinTheBaselineLimit`,
   `TestTheRoundfixSkillReadsItsModeFromTheRequest` and
   `TestEveryOwnedSkillVersionIsRecorded` in `./skills` with `-v`,
   recording each named test's and subtest's pass, for both the
   `canonical` and `embedded` layouts where the test has them.
3. MUST prove Invariant 3 and `_techspec.md` → Data Models through the built
   binary. In a fresh temporary Git repository named `app`, on branch
   `main` with one empty commit `feat: seed widget` and no remote, write a
   `roundfix/follow-up-messages/v1` file with two messages: id
   `handoff:h-42:1`, `from` `app`, `created_at`
   `2026-10-26T20:10:00Z`, subject `Finish the export retry`, text
   `Finish the export retry and rerun its QA.`, no `type`; and id `m-7`,
   `from` `billing`, `created_at` `2026-10-25`, text
   `FYI: billing moved its release to Friday.`. Then:
   - with the audited `bin/roundfix`, record that
     `follow-up --messages <file> --format json` exits `0`, reports
     `state` `work_found`, the `ai-memory-messages` source `checked`
     with count `2`, and an item with `ref` `message:handoff:h-42:1`,
     `kind` `message`, `requested_by` `app` and `title`
     `Finish the export retry`;
   - with a binary built from the Delivery Base in a temporary `git archive`
     copy (`go build -buildvcs=false -o <tmp>/roundfix-base ./cmd/roundfix`),
     record that the same command, in the same repository with the same
     file, prints byte-identical JSON, so the command and its schema are
     unchanged.
4. MUST record, from the audited checkout and from the Delivery Base read
   with `git show`, that no file under `.agents/skills/roundfix/` or
   `skills/roundfix/` of the Delivery Base names `memory_handoff_accept`,
   `newer_unended_session` or `not yet handled`, and that both
   `references/modes.md` copies of the audited head name all three, the
   defect the adopted Backlog Entry reports and its correction.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, run
   `go test -timeout 10m -count=1 -run '^TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage$' ./skills`
   there, and record that the named layout's subtest fails with the named
   message:
   - `skills/roundfix/references/modes.md` loses the words "read that
     session's newest observations with `memory_read_session_observations`
     before anything else, then ": `embedded` fails on step 1;
   - `.agents/skills/roundfix/references/modes.md` loses ", also when the
     mode stops early: at step 1, when the user stops it, or when step 3
     brings no answer" from the sentence that starts "No message the mode
     handed in": `canonical` fails on step 4;
   - `skills/roundfix/references/modes.md` loses "accepts no handoff, " in
     `## --status`: `embedded` fails with "--status section does not say it
     accepts no handoff";
   - `skills/roundfix/references/modes.md` loses the sentence "A message not
     yet handled becomes a Backlog Entry that records its sender, date and
     text, and is then consumed.": `embedded` fails on step 4;
   - `skills/roundfix/SKILL.md` loses ", and leave no handed-in message
     pending": `embedded` fails on the `## Modes` section;
   - `skills/roundfix/references/modes.md` replaces "List this project's
     messages with `memory_message_list`" with "List this project's messages
     without consuming them": `embedded` fails with "Optional sources still
     lists messages only to leave them pending".
6. MUST execute the commands task_01 and task_02 name in their
   Verification, and run `go test -timeout 15m -count=1 ./skills` and
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract`.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read `docs/adr/0330-the-secondbrain-is-read-only-and-ai-memory-carries-handoffs-and-messages.md`
     in the audited checkout and quote "consumes one only after recording its
     content in a Spec, a Backlog Entry or a Finding, or naming it when it
     only informs";
   - read the vendored `.agents/skills/ai-memory-handoff/SKILL.md` and quote
     "If the current context already contains a pending handoff block, answer
     from that block directly." and the sentence naming
     `newer_unended_session`;
   - read the vendored `.agents/skills/ai-memory-messaging/SKILL.md` and
     quote the sentence saying `memory_message_pop` claims "the oldest, or a
     specific `message_id`";
   - because this Spec's authorization forbids network access, record the
     published copies the PRD cites as `blocked (environment: network
     prohibited by Spec authorization)`, quoting the passages the PRD cites.
     Such a row never decides a qualifying `partial`.
8. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row, `blocked (declared unreachable)`, counted once in
   `rows_blocked_declared`: no row may reach an ai-memory server or start an
   interactive Agent session.
9. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` holds the **Session Handoff** entry of
     `_techspec.md` → Interfaces directly after the **Cross-project Message**
     entry, with one `_Avoid_` line, and the revised **Follow-up Item**
     (`docs/agents/domain.md`);
   - search `.agents/skills/roundfix`, `skills/roundfix` and
     `CONTEXT.md` for `without consuming them`, excluding this Spec's
     folder and `docs/history/`, and record that nothing matches;
   - record that the `roundfix` skill's two version fields are equal in the
     canonical copy and the mirror and recorded in
     `skills/testdata/owned-skill-versions.json`.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Go file
      outside `skills/roundfix_modes_test.go`, every file under
      `docs/user-guide/`, `docs/agents/` and `internal/baseline/`,
      `docs/references/behavior-surfaces.json`, and every skill other than
      `roundfix`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0325-the-follow-up-mode-consumes-the-handoff-and-every-message --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - call an ai-memory tool;
    - start an Agent session;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the end-to-end binary rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Both skill layouts take the Session Handoff in step 1 and settle every handed-in message, where the Delivery Base names neither.
- [ ] The audited and the Delivery Base binaries print the same JSON for a messages file with a handoff entry.
- [ ] Each sabotage makes its named subtest fail, and the scope audit finds no change outside the declared files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0325-the-follow-up-mode-consumes-the-handoff-and-every-message/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; Interfaces; Data Models; API Contracts 1-3; Invariants 1-4; Testing Approach; Build Order 3
