---
task: task_04
spec: 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks
status: completed
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec rewords the Baseline's three
message clauses so a Cross-project Message expects a reply only when its
sender's work is blocked, adds the optional `blocked_work` key to the Follow-up
Command's messages file with the `declared` basis that ranks such a message
first without the judge, changes the Roundfix skill's follow-up mode to fill
the key, plan a declared message first and reply only to it, and records the
decision in ADR-0383, its notices in ADR-0319 and ADR-0330, and the glossary.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts an Agent
session or an ACP adapter, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh`, reaches the
network, reads another repository, sends a request to a provider, the judge
or a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept
at their `go env` values, puts a fake `gh` script that prints `[]` first on
`PATH`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and every Judge
Key Variable for that command, and removes the temporary directories
afterwards. The QA never writes into the checkout while a test suite runs. The
gate keeps its evidence under names the repository does not ignore, `.txt` for
command output, and links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter (`rustfmt` through rustup) that needs the
session's toolchain configuration under `HOME`. Only rows that run
`bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 6 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestFollowUpRanksADeclaredBlockingMessageFirstWithoutTheJudge`,
   `TestFollowUpRefusesAMalformedBlockedWork`,
   `TestFollowUpOnlyDeclaredItemsAskNoJudge` and
   `TestFollowUpDeclaredBlockedWorkTranscripts` over `./internal/cli`,
   `TestTheMessageClausesAskForAReplyOnlyForBlockedWork`,
   `TestEveryInboxClauseAnAdopterHoldsIsReplaced` and
   `TestTheGuidesSendCrossProjectWorkThroughAiMemory` over
   `./internal/baseline`, and
   `TestTheFollowUpModeRepliesOnlyToDeclaredBlockedWork` and
   `TestTheFollowUpModeTakesTheHandoffAndSettlesEveryMessage` over
   `./skills`, recording every row by name; then the whole `./internal/cli`,
   `./internal/baseline` and `./skills` packages,
   `-tags docscontract -run 'TestTheFollowUpGuidesNameTheJSONFieldsAndMessageTypes' ./internal/cli`
   and `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce Surface Transcript 1 and Surface Transcript 2 of
   `_techspec.md` through the built binary, in a temporary Git repository
   built as that section's fixture says (one commit of `README.md`, an
   `origin` remote `git@github.com:owner/widget.git`, the messages file in
   the temporary directory outside the repository), comparing standard
   output, standard error and exit code with the transcript's conventions.
   Then run the first fixture with `--format json` and record that items 1
   and 2 carry `priority` `1`, `basis` `declared`, `blocks` equal to
   `requested_by`, their `blocked_work` and null `answers`, that item 3 has
   no `blocked_work` key, and that the same file without the two
   `blocked_work` keys lists all three items with no `blocked_work` key and
   no `declared` basis.
4. MUST reproduce the managed refresh in a temporary copy of the Delivery
   Base (the merge base of the audited head with the default branch) made
   with `git archive` and committed: record that
   `bin/roundfix baseline update --repo . --no-skills --no-history --yes`
   writes into `docs/agents/agent-instructions.md`, exactly once each, "A
   message expects no reply unless the sending project's own work depends on
   the answer", "say so in the message's subject and first lines" and "after
   the work already running, which is never interrupted for it", and that a
   second run reports the repository current.
5. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails on
   the named case: sabotages 1 to 7 as listed there. Sabotage 6 also runs
   the module record command and `make baseline-digests` in that copy before
   the test.
6. MUST execute the tests tasks 01 to 03 name in their Verification.
7. MUST search and record:
   - the command guide and both `follow-up` reference copies carry the
     `## Declared blocked work` section of API Contract 7;
   - both `modes` reference copies carry the steps of API Contract 9;
   - this repository's `docs/agents/agent-instructions.md`, the formatter
     golden of the Standard TypeScript profile and
     `internal/baseline/assets/modules/core.json` carry the three clauses of
     API Contract 8, and their ids are unchanged;
   - `CONTEXT.md` carries API Contract 10 and ADR-0319 and ADR-0330 carry
     API Contract 11; ADR-0383 is `accepted`;
   - `.agents/skills/roundfix` and `skills/roundfix` hold the same
     `follow-up` and `modes` bytes;
   - the phrases "when a reply is due" and "when the message asks for a
     reply" appear in none of `.agents/skills/roundfix`, `skills/roundfix`,
     `docs/agents`, `docs/user-guide`, `internal/baseline/assets/modules` and
     `internal/baseline/assets/formatter-fixtures`, searching only those
     paths (this Spec's folder and `docs/history/` quote the old words and
     are excluded).
8. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked, through the passages
   `_prd.md` → Acceptance evidence quotes:
   - Hohpe and Woolf, *Enterprise Integration Patterns*, "Request-Reply"
     (read 2026-11-10): compare "When two applications communicate via
     Messaging, the communication is one-way." with the sender rule of
     requirement 4, which makes a message one-way unless the sender needs
     the answer;
   - the Azure Architecture Center "Priority Queue pattern" (read
     2026-11-10): compare "The application or producer that sends a request
     assigns a priority value to the message, and consumers process the
     messages by priority." with requirement 3's report, where the sender's
     `blocked_work` decides the first place.
9. MUST record the PRD's one Unreachable Acceptance entry, a real Blocking
   Request in the operator's follow-up after the release, as one row counted
   in `rows_blocked_declared`, with its `satisfied-by`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `adds:` for **Blocking Request** and
    `changes:` for **Cross-project Message** and **Follow-up Priority**, that
    `CONTEXT.md` defines them with the text of API Contract 10, and that
    `TestGlossaryPhrasesAreDefinedTerms` passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, any file under
      `internal/judge/`, any Vendored Skill under `.agents/skills` (the
      `ai-memory-messaging` skill included), any file under
      `internal/baseline/assets/modules/` other than `core.json`, any file
      under `internal/baseline/assets/source-baselines/`, any file under
      `.claude/`, `.codex/` or `.opencode/`, any skill other than the
      Roundfix skill;
    - `internal/cli/follow_up_status.go` keeps its bytes from the Delivery
      Base;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the transcripts, the JSON report and the managed refresh.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A message entry that declares blocked work lists first with basis `declared` and never reaches the judge; a malformed declaration is refused; a file without it keeps its report.
- [ ] An adopter's managed refresh carries the reworded clauses, and the follow-up mode replies only to declared blocked work.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0373-a-message-expects-a-reply-only-when-it-declares-the-work-it-blocks/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-11; Surface Transcript 1; Surface Transcript 2; Invariants 1-5; Declared breaks; Testing Approach; Build Order 4
