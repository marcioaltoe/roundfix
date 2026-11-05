---
task: task_05
spec: 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Project Config key
`pre_pr_review.data_paths`, takes the declared files out of every code diff
with the reason `data-path`, reviews their diff apart in at most four Data
Review Segments of the same session (ADR-0368), and states the rule in the
configuration and `review` guides, the Roundfix skill, the glossary and
ADR-0282's amendment.

Every row runs on this repository's own files, temporary homes, temporary
repositories, the fake reviewers of `internal/cli` and temporary copies of
the audited head. No row starts an Agent session or an ACP adapter, reads or
writes the real `~/.roundfix` or the real home directory, reads another
repository, reaches the network, sends a request to a provider, the judge or a
reviewer, or reads a credential. Every `go test` it runs carries an explicit
`-timeout`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR`
to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at
their `go env` values, and removes them afterwards. The QA never writes into
the checkout while a test suite runs. The gate keeps its evidence under names
the repository does not ignore, `.txt` for command output, and links only
files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestPrePRReviewDataPathsLoadsAndValidates` over `./internal/config`, and
   `TestReviewDataSegmentsPackCutAndBound`, `TestReviewSendsADataPathApart`,
   `TestReviewSplitsALargeDataFileAtLineBoundaries`,
   `TestReviewBlocksDataNeedingMoreThanFourDataSegments`,
   `TestReviewOmissionWinsOverADataPath` and
   `TestReviewReadsADataPathApartInRoundTwo` over `./internal/cli`,
   recording every subtest by name; then the whole `./internal/config`,
   `./internal/cli` and `./skills` packages,
   `go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec`
   and `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce through the built binary, with a fresh temporary `HOME`
   and `TMPDIR`, in a temporary Git repository whose committed
   `.roundfixrc.yml` holds `pre_pr_review:` with `provider: none` and
   `data_paths: [rules/]`, that `bin/roundfix review` no longer fails with
   `field data_paths not found in type config.prePRReviewOverlay` and exits
   `0` with outcome `omitted`; then, with `data_paths: [/abs]` instead,
   that it exits `2` and stderr names `pre_pr_review.data_paths has unsafe
   path "/abs"`. Remove every temporary directory afterwards.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails on the named case:
   - `ReviewsApart` always false: `TestReviewSendsADataPathApart` fails on
     `codex` and `claude`;
   - the `data-path` reason checked before `omit-path`:
     `TestReviewOmissionWinsOverADataPath` fails;
   - the new-line cursor not advanced between pieces:
     `TestReviewDataSegmentsPackCutAndBound` fails on `a new file is cut at
     lines`;
   - the four-segment limit unchecked: it fails on `a fifth segment is
     refused`, and `TestReviewBlocksDataNeedingMoreThanFourDataSegments`
     fails;
   - `appendReviewOmissions` counting `data-path` as omitted:
     `TestReviewSendsADataPathApart` fails.
5. MUST execute the tests tasks 01 to 04 name in their Verification.
6. MUST search `docs/user-guide/configuration.md`, both copies of the
   Roundfix `review` reference, `docs/user-guide/commands/review.md`,
   `CONTEXT.md` and ADR-0282, and record that each states the text of API
   Contract 6 for it, and that ADR-0282 carries the amendment naming ADR-0368
   with `updated_at` of 2026-11-05.
7. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the Fiscus report quoted in
     `references/2026-11-05-fiscus-the-delivery-queue-takes-days-and-a-large-data-catalog-is-omitted-from-review.md`
     (2026-11-05): compare offline its point 5 with
     `TestReviewSendsADataPathApart`, a catalog of about the size it names
     reviewed in a Data Review Segment instead of omitted;
   - Git's own content, through the `git` binary: record the two cut
     subtests of `TestReviewDataSegmentsPackCutAndBound`, whose oracle is
     `git show <head>:<path>`, and compare them with the git-diff(1) passage
     `_prd.md` → Acceptance evidence quotes.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `adds: **Review Data Path**` and `adds: **Data
   Review Segment**`, that `CONTEXT.md` defines both as API Contract 6
   states, and that the glossary test of `./internal/docscontract` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/skill-coverage.json`,
      `docs/references/coverage-record.json`, anything under
      `internal/baseline/`, `internal/delivery/`, `internal/store/` or
      `docs/agents/`, any ADR other than ADR-0282, any skill other than the
      Roundfix skill files task_04 declares;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0358-a-large-authored-data-file-is-reviewed-apart-not-omitted --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the config through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A declared authored data file is reviewed in Data Review Segments, never in the code diff, in every round, for `codex` and `claude`.
- [ ] Oversized data is cut with valid anchors up to four segments and refused beyond, before any session.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0358-a-large-authored-data-file-is-reviewed-apart-not-omitted/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; User Story 1; User Story 2; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-5; Declared breaks; Testing Approach; Build Order 5
