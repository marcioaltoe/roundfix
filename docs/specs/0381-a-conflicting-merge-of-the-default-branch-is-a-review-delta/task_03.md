---
task: task_03
spec: 0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta
status: pending
type: qa
complexity: low
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes the pre-PR
review measure across a conflicting merge of the default branch from its
Conflicted Carry, keeps the Reviewer Lineage across that merge, retires the
conflict restart reason, says so in the round-2, Correction Review and
segment prompts, adds the absence-claim sentence to every review prompt, and
documents it in the `review` and `deliver` guides, the Roundfix skill
references, the glossary and the ADR-0366 and ADR-0197 notices.

Every row runs in this repository or in temporary directories it creates
and removes. No row creates, removes or changes any file in `/tmp`,
`/private/tmp` or `$HOME` outside a directory made by `mktemp -d`, reaches
the network, reads another repository, starts an Agent session or an ACP
adapter, reads or writes the real `~/.roundfix`, sends a request to a
provider, the judge or a reviewer, records an environment value, or sends a
signal to any process. Every `go test` carries an explicit `-timeout` and
keeps the session's own `HOME` and `TMPDIR`. Every row that runs
`bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary directories, with
`GOMODCACHE` and `GOCACHE` kept at their `go env` values, and removes them
afterwards. The QA never writes into the checkout while a test suite runs.
The gate keeps its evidence under names the repository does not ignore,
`.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v ./internal/cli`, the
   tests `TestAConflictingMergeOfTheDefaultBranchKeepsTheReviewerLineage`,
   `TestAReviewSegmentAcrossAConflictingMergeLeavesOutTheDefaultBranch`,
   `TestTheReviewPromptAsksToReadTheFileBeforeAnAbsenceClaim`,
   `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange` and
   `TestReviewLineageRestartNamesItsReason`, recording the pass line of each
   subtest; then `-run 'Review|Lineage|Dispos|Dismiss|Correction|Deliver|Segment'`
   over `./internal/cli` with `-timeout 30m`, and
   `GOOS=windows go vet ./internal/cli`.
3. MUST record Success Metrics 1–3 against the starting behavior: unpack a
   `git archive` of the Delivery Base (the merge base of the audited head
   with `main`) in a `mktemp -d` directory, make it a Git repository there
   (`git init`, `git add -A`, one commit with a throwaway identity passed
   with `-c`), copy the audited head's
   `internal/cli/review_conflicted_carry_test.go` over, and run its three
   tests there with `-timeout 20m`. Record that
   `TestAConflictingMergeOfTheDefaultBranchKeepsTheReviewerLineage/round_two`
   fails there with a round-1 record after the conflicting merge (the
   reported restart), `…/correction_review` fails there,
   `TestAReviewSegmentAcrossAConflictingMergeLeavesOutTheDefaultBranch` fails
   there with an empty `measuredFrom` for the second segment, and
   `TestTheReviewPromptAsksToReadTheFileBeforeAnAbsenceClaim` fails there;
   and that each passes on the audited head (row 2). Remove the directory
   afterwards.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive` and turned into a Git repository as in row 3, with
   `-timeout` on every `go test`, and record that the named test fails:
   - `decideReviewLineage` restarts again when the measurement is
     conflicted: `TestReviewLineageRestartNamesItsReason/conflict` and
     `TestAConflictingMergeOfTheDefaultBranchKeepsTheReviewerLineage/round_two`;
   - `reviewMeasuredFrom` returns the earlier commit on exit 1 instead of
     calling `reviewConflictedCarry`:
     `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange/conflicting_carry`;
   - `reviewConflictedCarry` drops `-c merge.conflictStyle=merge`:
     `TestReviewMeasuredFromCarriesOnlyTheCandidatesOwnChange/conflicting_carry`;
   - `buildRoundTwoPrompt` drops the Conflicted Carry note:
     `TestAConflictingMergeOfTheDefaultBranchKeepsTheReviewerLineage/round_two`;
   - `planReviewSegments` drops the segment note from the sent prompt:
     `TestAReviewSegmentAcrossAConflictingMergeLeavesOutTheDefaultBranch`;
   - `buildReviewPrompt` drops the absence-claim sentence:
     `TestTheReviewPromptAsksToReadTheFileBeforeAnAbsenceClaim`.
5. MUST execute the tests every non-QA Task names in its Verification.
6. MUST record, as the outside-evidence row, sources this Spec did not
   write, with their dates and not blocked: the adopted Conexus message
   (`references/`, reported by conexus, measured on Roundfix 0.138.0 in a
   repository this Spec did not build), whose round-1 review right after the
   conflict park's merge matches row 3's starting behavior; and the Git
   manual passages quoted in `_prd.md` → Acceptance evidence (`git
   merge-tree`: "If there were conflicts, then files within this tree may
   have embedded conflict markers" and "When the merge has conflicts, the
   exit status is 1"), compared with the `conflicting_carry` and `round_two`
   rows, whose tree holds the markers and whose delta removes them. The row
   reads the files in this repository, never the network.
7. MUST repeat the wording search of task_01: no file under
   `docs/user-guide/`, `.agents/skills/roundfix/` or `skills/roundfix/`
   contains `conflicts with it` or `one of the eight the`, both `review`
   texts list seven restart reasons and describe the Conflicted Carry, and
   no file under `internal/cli` that is not a test contains `conflicts with
   it`.
8. MUST perform the glossary check (`docs/agents/domain.md`): record that
   the PRD and TechSpec declare `adds: **Conflicted Carry**`, `changes:
   **Effective Base**` and `changes: **Reviewer Lineage**`, that `CONTEXT.md`
   defines all three as the Task commit range left them, and that ADR-0366
   and ADR-0197 carry the ADR-0391 notice with `updated_at` raised.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is listed in `_authorization.md` `paths`;
   - the Task commit range changes none of these: the Makefile, any
     workflow, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
     `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/**`,
     `internal/preflight/**`, `internal/delivery/**`, any ADR other than
     ADR-0366 and ADR-0197, any Baseline module or derived Baseline file, or
     the `### QA settlement` section of any skill.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository, or run `git fetch`, `git push` or `gh`;
    - read or write `~/.roundfix`;
    - create a file in `/tmp`, `/private/tmp` or `$HOME` outside a
      `mktemp -d` directory it removes;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification and the named suites.
- [ ] Characterize the Delivery Base with the new tests and run each sabotage on a temporary copy.
- [ ] Record the outside sources, the wording search, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A review after a conflicting merge of the default branch continues the lineage from the Conflicted Carry on the audited head and restarts at round 1 on the Delivery Base.
- [ ] Each sabotage makes its named test fail, and every Git failure still blocks.
- [ ] The guides, references and glossary describe the Conflicted Carry, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0381-a-conflicting-merge-of-the-default-branch-is-a-review-delta/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Testing Approach; Build Order 3
