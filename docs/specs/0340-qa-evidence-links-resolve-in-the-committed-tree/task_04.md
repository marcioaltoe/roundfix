---
task: task_04
spec: 0340-qa-evidence-links-resolve-in-the-committed-tree
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes every QA settlement
path, the Daemon's QA settlement, `roundfix qa-report accept` and the Settle
Command for a `qa` Task, refuse a QA Report whose local link the committed
tree would not hold, naming the link, its line and the reason `ignored` or
`missing`; the Daemon first gives the gate one repair turn; the QA check and
`SC-HISTORY-RETAINED` share one link extraction; and the `qa-gate` skill tells
the gate to keep evidence under names the repository does not ignore.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. No row compares one
command's output with another's byte for byte. The gate keeps its own
evidence under names the repository does not ignore, `.txt` for command
output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 6 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestUnresolvedEvidenceLinksNameTheReason`,
   `TestUnresolvedEvidenceLinksOutsideGitReadTheDisk`,
   `TestEvidenceLinksErrorNamesEachLink` and
   `TestMissingRelativeLinksReadsTheSharedExtraction` over `./internal/spec`,
   `TestQAReportAcceptRefusesAnUnresolvedEvidenceLink`,
   `TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink` and
   `TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape` over
   `./internal/cli`, the whole `./internal/daemon` and `./internal/speccheck`
   packages, recording each pass.
3. MUST reproduce the Argus case through the built binary in a temporary Git
   repository the row creates and removes: commit a `.gitignore` holding
   `*.log` and a fixture `docs/specs/9001-fixture/qa/qa-report-2026-11-01.md`
   with `verdict: pass` and a Results table of three rows linking
   `evidence/consumer-tests.log` (present, ignored), `evidence/journeys.txt`
   (absent) and `evidence/kept.txt` (committed). With
   `h="$(mktemp -d)"; t="$(mktemp -d)"`, run
   `HOME="$h" TMPDIR="$t" bin/roundfix qa-report accept docs/specs/9001-fixture/qa/qa-report-2026-11-01.md`
   from the repository and record exit 1 and a stderr line naming
   `evidence/consumer-tests.log` with `ignored` and `evidence/journeys.txt`
   with `missing`, both with their lines. Then save the log as
   `evidence/consumer-tests.txt`, create `evidence/journeys.txt`, update both
   links, rerun, and record exit 0 with empty output; then
   `rm -rf "$h" "$t"` and the repository.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - `UnresolvedEvidenceLinks` treats every existing target as resolved:
     `TestUnresolvedEvidenceLinksNameTheReason` fails on `evidence/a.log`;
   - `--cached` dropped from the `ls-files` call:
     `TestUnresolvedEvidenceLinksNameTheReason` fails on the forced-add
     `evidence/c.log`;
   - `RelativeLinks` keeps links inside code spans:
     `TestUnresolvedEvidenceLinksNameTheReason` fails on the code-span case;
   - the check removed from `requireSettleQAReportEligibility`:
     `TestSettleRefusesAQAReportWithAnIgnoredEvidenceLink` fails;
   - the repair turn removed from `runQAGate`:
     `TestQAEvidenceLinkRepairTurnKeepsTheEvidence` fails;
   - the commit-boundary check removed:
     `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed` fails;
   - the repair turn repeated until the links resolve:
     `TestQAEvidenceLinkUnresolvedAfterRepairSettlesFailed` fails on the
     repair prompt count.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
6. MUST search `docs/user-guide/`, every copy of the Roundfix skill and every
   copy of the `qa-gate` skill and record that the qa-report and settle
   guides and both copies of the Roundfix `settle` reference state that every
   local link of the report must resolve in the committed tree with both
   reasons; that the implement guide and both copies of the Roundfix
   `implement` reference state the one repair turn and the refusal without a
   QA Report commit; and that both copies of `qa-gate` carry the section
   `### Evidence links resolve in the committed tree` with "Never cite a file
   you did not keep."
7. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the Argus report quoted in
     `references/2026-11-01-qa-evidence-links-resolve-in-the-committed-tree.md`
     (2026-11-01): compare offline its shape (an accepted report linking
     ignored `.log` files absent from the QA commit) with the row of
     requirement 3;
   - the `git-check-ignore` manual passage `_prd.md` → Acceptance evidence
     quotes ("tracked files are not shown at all since they are not subject
     to exclude rules"): compare it with the forced-add case of
     `TestUnresolvedEvidenceLinksNameTheReason` and Invariant 2;
   - the `gitignore` manual passage `_prd.md` → Acceptance evidence quotes
     ("specifies intentionally untracked files to ignore"): compare it with
     the `ignored` reason of API Contract 2.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `adds: **Unresolved Evidence Link**` and
   `changes: **QA Report**`, that `CONTEXT.md` defines the first and carries
   the new sentence in the second, and that the Task commit range changes no
   other glossary entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `docs/references/skill-coverage.json`, `internal/spec/qa.go`,
      `internal/spec/archive.go`, anything under `internal/baseline/assets/`
      or `docs/agents/`, any skill other than the Roundfix and `qa-gate`
      files the Tasks declare;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0340-qa-evidence-links-resolve-in-the-committed-tree --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the Argus case through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An ignored or missing evidence link is refused with its line and reason by `qa-report accept`, the Settle Command and the Daemon, in the tests and through the binary.
- [ ] The Daemon runs exactly one repair turn and never commits a refused report.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0340-qa-evidence-links-resolve-in-the-committed-tree/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-8; Invariants 1-5; Declared breaks; Testing Approach; Build Order 4
