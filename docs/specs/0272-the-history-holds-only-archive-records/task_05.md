---
task: task_05
spec: 0272-the-history-holds-only-archive-records
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes a retirement delete, gives
every citation of a removed file a Git Locator in the Archive Command, the
Legacy Link Rebase and both history commands, holds the History Root to
Archive Records through the Spec Consistency Check once the Baseline carries
the rule, states that rule once in the Baseline, and removes this
repository's retired history through its own Managed Refresh. No row reads or
writes the real `~/.roundfix`, reaches a provider or the network, reads a
credential or reads another repository.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 on the audited head: `docs/history/` holds
   exactly 243 files, every one an Archive Record directly under
   `docs/history/specs/`, none of the ten `superseded` records 0120 to 0129,
   and no file under the Findings, Backlog, ADR, Review or Handoff family.
3. MUST verify Success Metric 2 and record the sabotage proof against the
   delivery base, the merge base of the audited head with `main`: build a
   second binary from a temporary worktree at the delivery base, and in a
   temporary repository whose docs-layout guide carries
   `RetirementRuleSentinel` plant, one at a time, a file under each retired
   family, a Reduced History Entry, a `done` Finding in `docs/findings/`, a
   relative link from an Archive Record to a missing file and a record
   linking to itself. Record that `roundfix spec check` exits 1 with
   `SC-HISTORY-RETAINED` and a Git Locator or deletion in its Fix for each
   with the audited binary and passes each with the delivery-base binary, and
   that the unplanted tree passes with both. Record that the same plants pass
   the audited binary when the guide lacks the sentence.
4. MUST verify Success Metric 3: every Git Locator in an Archive Record, a
   Live Document or a page under `docs/user-guide/` of the audited head names
   a commit that holds its path and that `main` or a tag reaches, or reads the
   pull request or `not reachable` form, and no Archive Record links to
   itself. Record that a planted Git Locator naming a commit no ref reaches
   fails the check, and that `ResolveGitLocator` gave such a commit to no
   writer in the task tests (API Contract 1, Invariant 1).
5. MUST verify Success Metric 4: the rendered `guide.docs-layout` region of
   `docs/agents/docs-layout.md` is at most 11,000 bytes, and no clause of the
   changed modules other than `clause.context.findings-09-archive` states the
   retirement rule (API Contract 15).
6. MUST verify Success Metric 5 with an adopter dry run in the QA scratch: a
   synthetic adopter repository built there the way
   `internal/cli/baseline_update_history_test.go` builds one, holding an
   annotated `history-full` tag, Reduced History Entries, a retired ADR in
   `docs/adr/`, a terminal Backlog Entry in `docs/backlog/`, an
   `_archived/findings/` tree, a cited and an uncited `superseded` record, a
   record with a broken relative link, and a user guide page citing a reduced
   entry. Record that `roundfix baseline update --format json` lists every
   removal, every citation with its `rewrite` and the `record-links` unit
   before anything changes and exits 3; that `--yes` applies it; and that a
   second run exits 0 with history `current`. No real repository is cloned.
7. MUST reproduce Surface Transcript 1 and Surface Transcript 2 through the
   built `bin/roundfix` on their declared fixture, comparing every stdout and
   stderr line and the exit code under the transcript conventions.
8. MUST verify API Contracts 4 and 6 through the built binary: archive a
   synthetic Spec whose outcome links an adopted reference and a QA evidence
   file and whose promoted copy links `../_prd.md` and an ADR, and record that
   the outcome and copy carry Git Locators or rebased destinations, that no
   destination resolves to a missing file, and that `roundfix deliver`'s
   exact-archive proof accepts the commit (the tests task_01 names).
9. MUST execute the tests task_01, task_02, task_03 and task_04 name in their
   Verification.
10. MUST record, as evidence this Spec did not author, each source of
    `_prd.md` → Acceptance evidence: the quereus documentation conventions
    passage "History is deleted, not archived.", the RFC 2119 section 6
    passage on imperatives, and the GitHub Docs passage that a commit
    identifier "will permanently link to the exact version of the file in that
    commit", quoting each from the PRD and stating whether it still supports
    the design. The sources are not fetched, because this Spec's
    authorization keeps the gate offline; a row that needs more than the
    quoted passage records
    `blocked (environment: network prohibited by Spec authorization)`.
11. MUST perform the glossary check of `docs/agents/domain.md`: record that
    `CONTEXT.md` defines **Git Locator** and that each changed term the
    Glossary sections declare cites ADR-0278, and whether any other term the
    Spec introduced needs one.
12. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`; that every
    Governed Path changed is bounded in `_authorization.md`; that the
    Makefile, the CI workflows, `go.mod`, `go.sum` and the lint configuration
    did not change; that the module and skill versions were raised and
    recorded; that `make baseline-digests` rewrites nothing on the audited
    head; and that no `### QA settlement` section changed. This row reads Task
    commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
14. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, record any environment value, send any
    request to a provider, or start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Verify the tree, the locators, the region and the transcripts.
- [ ] Prove the check against the delivery base and run the adopter dry run.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The History Root of the audited head holds only the 243 Archive Records.
- [ ] Each planted retained file or broken link fails the audited binary and passes the delivery-base binary.
- [ ] The adopter dry run lists every removal before applying and a second run reports the history current.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0272-the-history-holds-only-archive-records/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contracts 1-15; Surface Transcript 1; Surface Transcript 2; Invariants 1-6; Coverage Map deviations; Build Order 5
- ADR-0278; ADR-0248; ADR-0254; ADR-0260; ADR-0261; ADR-0247
