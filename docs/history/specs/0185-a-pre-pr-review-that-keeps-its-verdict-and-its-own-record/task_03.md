---
task: task_03
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
status: completed
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the Spec's authored terminal gate. It declares what the matrix covers
and settles the Spec on evidence. The Spec changes production code in
`internal/agent` and `internal/cli/review.go`, and the guides that describe
`roundfix review`. Every behavior row is exercised by running the named tests
against the built tree. The tests use fake ACP streams, fake Agent runners,
disposable repositories and disposable Artifact Directories. The only reads of
the live machine are the read-only outside-evidence rows below. No command
reaches a reviewer, a provider or the network.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing the message-log and runner tests against the
   built tree:
   - chunks without `messageId` split at a thought or a tool call and stay one
     message across a status update;
   - chunks with `messageId` split where the identifier changes and stay one
     message across a tool call when it is the same;
   - `Answer()` returns the last non-blank message;
   - a one-message stream keeps `Message` unchanged;
   - a sealed stream with commentary, a thought and JSON returns exactly the
     JSON, and the sealed cap still refuses.
3. MUST verify, by executing the review classification tests against the
   built tree:
   - a progress message then `Findings:` records `findings`;
   - a progress message then `No findings.` records `reviewed`;
   - a final message with both verdicts records `blocked`;
   - the answer file keeps every message;
   - `TestReviewClassifiesVerdictVariants` and
     `TestReviewPassesOnlyAWholeAnswerVerdict` still pass.
4. MUST verify, by executing the record-per-checkout tests against the built
   tree:
   - two checkouts sharing one Artifact Directory keep separate records;
   - a review in one leaves the other's record and answer bytes unchanged;
   - each checkout reuses and disposes only its own record;
   - a record at the old shared location is neither reused nor removed.
5. MUST record, as evidence this Spec did not author, the Agent Client
   Protocol's message-id RFD
   (https://github.com/agentclientprotocol/agent-client-protocol/blob/main/docs/rfds/message-id.mdx):
   `messageId` is stable across a message's chunks, and consecutive chunks
   without it are ambiguous. It MUST verify that the implemented boundary rule
   follows it.
6. MUST read, without writing, the installed codex-acp adapter at
   `/opt/homebrew/lib/node_modules/@agentclientprotocol/codex-acp/dist/index.js`.
   It MUST record whether it builds agent message chunks with a `messageId`
   (`createAgentTextMessageChunk(item.text, item.id, …)`), as evidence from a
   repository this Spec did not build. When the file is absent, it MUST record
   the row as blocked with its reason.
7. MUST read, without writing, the live shared record
   `~/.roundfix/artifacts/339f8dac2b687a04/pre-pr-review.json` and the ledger
   beside it, and record which checkout and head the record names and which
   heads the ledger's dispositions name. This is evidence written by other
   sessions on 2026-09-29. When either file is absent, it MUST record the row
   as blocked with its reason.
8. MUST verify that `docs/user-guide/commands.md`, the Roundfix skill and its
   mirror, and `CONTEXT.md` describe the final-message rule and the record per
   checkout, and that `make skills-sync-check` exits `0`.
9. MUST verify that this Spec's own artifacts satisfy the promise rule.
10. MUST verify from the repository history that the changed files stay within
    the paths the Tasks declare or the Daemon recorded, and that every governed
    path is bounded in `_authorization.md`.
11. MUST NOT accept a row whose only evidence is that a file was read.
12. MUST NOT run `roundfix review` against a real reviewer, and MUST NOT write
    to any Artifact Directory under `~/.roundfix`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Core Features 1-4; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-4; API Contracts 1-2; ADR-0017;
ADR-0020; ADR-0080; ADR-0091; ADR-0093; ADR-0094; ADR-0096; ADR-0104;
ADR-0117; ADR-0151; ADR-0153; ADR-0155; ADR-0156; ADR-0165; ADR-0169;
ADR-0174.
