---
task: task_05
spec: 0342-the-claude-review-segments-and-omits-declared-paths
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Project Config key
`pre_pr_review.omit_paths`, which the pre-PR review reads at the base commit
and records as `omit-path` omissions; reviews a Claude round-1 candidate over
its token budget in up to four Review Segments; makes every size refusal name
the three largest files of its diff and the key; parks a blocked review with
its reason line; and documents them.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head, with the fake Agent runners of the
review tests. No row starts an Agent session or an ACP adapter, reads or
writes the real `~/.roundfix`, reads another repository (Vortex included),
reaches the network, sends a request to a provider, the judge or a reviewer,
or reads a credential. Every `go test` it runs carries an explicit
`-timeout`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` set to the
values `go env` prints, and removes them afterwards. The QA never writes into
the checkout while a test suite runs. The only byte comparison is the
omission summary row of requirement 3, between the summary line of a review
without the key and the literal sentence the TechSpec quotes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 6 and Invariants 1 to 6 by running
   `go test -count=1 -timeout 30m -v` with
   `TestPrePRReviewOmitPathsLoadsAndValidates` and
   `TestProjectConfigTemplateDescribesOmitPaths` over `./internal/config`;
   `TestReviewOmitsConfiguredPathsFromTheBaseCommit`,
   `TestReviewOmissionSummaryKeepsItsBytesWithoutOmitPaths`,
   `TestReviewSegmentsAClaudeCandidateOverItsBudget`,
   `TestReviewBlocksAClaudeCandidateNeedingMoreThanFourSegments`,
   `TestReviewKeepsClaudeRoundTwoAndOversizedCommitsUnsegmented`,
   `TestReviewKeepsCodexSegmentRecordsWithoutTokenEstimates`,
   `TestReviewSizeRefusalNamesTheLargestFiles`,
   `TestReviewSegmentsACandidateAboveTheBound` and
   `TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound` over
   `./internal/cli`; and `TestReviewBlockedParkCarriesTheReviewReason` over
   `./internal/delivery`, recording each pass.
3. MUST reproduce the Vortex shape through the built binary's config loader
   in a temporary repository the row creates and removes: a Project Config at
   the base commit with `pre_pr_review.omit_paths` listing
   `packages/backend/drizzle/migrations/meta/*_snapshot.json`, then
   `bin/roundfix profiles show` with fresh `HOME` and `TMPDIR`, recording exit
   0; the same repository with the entry `../outside.json` recording exit 2
   and an error naming `pre_pr_review.omit_paths` and the entry. (At
   `4e61bf6d` the key itself is refused with exit 2 and `field omit_paths not
   found in type config.prePRReviewOverlay`.) The review
   behavior itself is proven by requirement 2's tests, which use fake
   runners; no row runs `roundfix review` against a real provider.
4. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails, with
   the case whose input makes the mutation observable.
5. MUST execute the tests task_01 to task_04 name in their Verification.
6. MUST search `docs/user-guide/`, both copies of the Roundfix skill,
   `CONTEXT.md` and `docs/adr/` (excluding `docs/history/` and this Spec's
   folder) and record that: no live text still says Claude "is never
   segmented"; the review guide and both review reference copies name
   `omit-path`, the Claude segment sentence and the refusal suffix; the
   configuration guide has the `pre_pr_review.omit_paths` row; the deliver
   guide and both deliver reference copies state the park reason; and
   ADR-0227 and ADR-0282 carry their ADR-0352 amendment lines.
7. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the Vortex report quoted in
     `references/2026-11-03-the-claude-pre-pr-review-blocks-on-size-instead-of-segmenting.md`
     (2026-11-03): compare offline its shape (a Claude refusal at 541,455
     estimated tokens with one 669 KB generated snapshot, a `derived_paths`
     workaround whose regeneration exits 1) with requirement 2's omit and
     refusal tests and requirement 3;
   - the GitHub passage `_prd.md` → Acceptance evidence quotes ("To keep
     certain files from displaying in diffs by default"): compare it with the
     repository-declared omission the key provides;
   - the Claude Agent SDK passage `_prd.md` → Acceptance evidence quotes
     ("When the context window approaches its limit, the SDK automatically
     compacts the conversation"): compare it with the one-session Claude
     segments of requirement 2.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD declares `adds: **Review Omit Path**` and `changes: **Review
   Segment**`, that `CONTEXT.md` defines both, and that the Task commit range
   changes no other glossary entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/`,
      `docs/references/skill-coverage.json`, any Archive Record, any Baseline
      module or guide under `docs/agents/`, any skill file other than the
      Roundfix `SKILL.md` pair and its review and deliver references, and
      the Roundfix `modes.md`, `follow-up.md` and `implement.md` references;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0342-the-claude-review-segments-and-omits-declared-paths --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and the config rows through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A base-commit omit path leaves its file out with reason `omit-path`, and a review without the key keeps its prompt bytes.
- [ ] A Claude candidate over its budget is reviewed in segments within the budget, and every size refusal names its largest files and the key.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0342-the-claude-review-segments-and-omits-declared-paths/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-6; Declared breaks; Testing Approach; Build Order 5
