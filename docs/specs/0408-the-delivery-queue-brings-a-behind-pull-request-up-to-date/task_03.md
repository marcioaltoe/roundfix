---
task: task_03
spec: 0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec makes the Delivery Engine's checking stage treat the merge state `BEHIND` as one that never merges and answer it with a Base Update: the Queue Owner merges the refreshed default branch into the item branch through the queue's existing merge, records the merged head and returns the item to `gating`. Three Base Updates per item and Queue Owner bound it, and a fourth parks `pull-request-behind` of Park Class `environment`. The `deliver` reference, the user guide and the glossary describe it.

Every row runs on this repository's own files, temporary homes, temporary repositories and temporary copies of the audited head. No row:

- starts an Agent harness, an Agent session or an ACP adapter;
- starts a Run, a Daemon or a Delivery Owner outside a test;
- calls ai-memory, runs the real `gh` or reaches the network;
- reads or writes the real `~/.roundfix`;
- sends a request to a provider, the judge or a reviewer;
- asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout` and keeps the session's own `HOME` and `TMPDIR`. When the gate isolates `HOME` for a row that runs `./internal/cli` or the repository Verification, it keeps the host Rust toolchain readable: `PATH`, `RUSTUP_HOME` and `CARGO_HOME` point at the host's toolchain, read only. Only rows that run `bin/roundfix` use fresh temporary homes. Each such row:

- takes `m="$(go env GOMODCACHE)"` and `c="$(go env GOCACHE)"` first;
- sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"` and `GOCACHE="$c"`;
- unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- reads standard input from `/dev/null` and removes the temporary directories afterwards.

The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1, Success Metric 3, Success Metric 4 and Success Metric 5 by running, with `-v`, the `internal/delivery` tests task_01 names in its Verification, recording every test and subtest by name, then the whole `./internal/delivery` package, recording its pass.
3. MUST verify Success Metric 2 by running, with `-v`, `TestCommandDeliveryEngineBringsABehindPullRequestUpToDate` in `./internal/cli` and recording from its output the `merge state BEHIND` line, the `roundfix: base update: Delivery Queue item …` line, the `merge state CLEAN` line for the merged head, and the pass. Then run the whole `./internal/cli` package and record its pass.
4. MUST record the exact texts of `_techspec.md` → API Contract 1 and API Contract 2 as the tests of requirement 2 assert them: the owner log line, the two `pull-request-behind` blockers, the next action and the `checks-timeout … ; merge state BEHIND` form. It MUST also record that `delivery.ClassifyPark` returns Park Class `environment` for `pull-request-behind`.
5. MUST run `-tags docscontract ./internal/docscontract` and `./skills` whole, recording each pass, and record:
   - the sentences of `_techspec.md` → API Contract 3 in `docs/user-guide/commands/deliver.md`, `.agents/skills/roundfix/references/deliver.md` and `skills/roundfix/references/deliver.md`, and that the mirror equals its canonical source (`make skills-sync-check`);
   - that the sentence "A second refusal for the same head parks `delivery-error`." is still present in each of the three documents;
   - that a search of `docs/user-guide`, `.agents/skills/roundfix` and `CONTEXT.md` finds no remaining text saying the item moves to `merging` when the merge state is neither `BLOCKED` nor `UNKNOWN` alone.
6. MUST run each sabotage below on its own temporary copy of the audited head made with `git archive` and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. Every `go test` carries `-timeout`. Record that the named test fails on the named case. No sabotage regenerates or synchronizes anything.
   1. The `BEHIND` path of `checkCandidate` is removed, so `BEHIND` counts as settled: `TestABehindPullRequestReturnsToTheGateWithTheMergedHead` fails, and `TestCommandDeliveryEngineBringsABehindPullRequestUpToDate` fails with the item parked `delivery-error`.
   2. The bound test always passes (the count is never compared with `BaseUpdateLimit`): `TestAFourthBehindPullRequestParksPullRequestBehind` fails.
   3. `Engine.park` no longer drops the item's count: `TestTheBaseUpdateCountStartsOverAfterAPark` fails.
   4. `BEHIND` is removed from the awaited merge states of `checksTimeoutBlocker`: `TestABaseUpdateWithNothingToMergeKeepsWaiting/times_out_with_the_merge_state` fails.
   5. The Base Update runs before the `CONFLICTING` test: `TestAConflictingAndBehindPullRequestTakesTheConflictPath` fails.
   6. `ClassifyPark` drops the `pull-request-behind` case: `TestParkClassesAndNextActions/pull-request-behind` and `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem` fail.
7. MUST record the outside-evidence row of `_prd.md` → Success Metric 6 and Acceptance evidence from [`references/2026-11-18-behind-pull-request-characterization.md`](references/2026-11-18-behind-pull-request-characterization.md), section "Published evidence", without the network: quote the passages of the GitHub GraphQL reference for `MergeStateStatus` ("`BEHIND`: The head ref is out of date."), of "Troubleshooting required status checks" ("If branch protection requires your branch to be up-to-date, merge or rebase the base branch into your branch."), and of multiplex-ai/muggle-ai-works pull request #265 ("`BEHIND` surfaces at all only when the base enforces \"require branches up to date\"."), each with its URL as the reference gives it, and record that the design reads `BEHIND` alone, so a base branch that does not require up-to-date branches keeps today's behavior.
8. MUST record the Unreachable Acceptance entry of `_prd.md` as one row `blocked (declared: a delivery Pull Request on a real GitHub repository whose base branch requires up-to-date branches merges after a Base Update)`, supported by the strict local forge evidence of requirement 3.
9. MUST execute the tests task_01 and task_02 name in their Verification.
10. MUST perform the glossary check (`docs/agents/domain.md`). Record that the PRD and the TechSpec declare `adds:` for **Base Update** and `changes:` for **Check Wait**, that task_02 wrote those entries and changed no other, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `skills-lock.json`, `.roundfixrc.yml`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, `docs/references/coverage-record.json`, any file under `internal/baseline/`, any Vendored Skill under `.agents/skills`, any owned skill other than `roundfix`, any ADR, and any Go source under `internal/cli` other than `internal/cli/deliver_base_update_test.go`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent harness, an Agent session, an ACP adapter, a Run, a Daemon or a Delivery Owner outside a test;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of both Tasks and the packages.
- [ ] Record the exact texts, the documents and the glossary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the outside evidence, the unreachable row, the scope audit and the Pull Request row, and write the QA report.

## Acceptance Criteria

- [ ] A Pull Request reported `BEHIND` is brought up to date and merged at its new head, never at the behind head, in the Delivery Engine and through the command engine with a strict local forge.
- [ ] A conflict, nothing to merge, the bound and the restart after a park each end as the TechSpec says, with the exact texts.
- [ ] The three documents and the glossary describe the Base Update, and the second policy refusal still parks `delivery-error`.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-4; Interfaces and invariants 1 to 7; Testing Approach; Build Order 3
- [`references/2026-11-18-behind-pull-request-characterization.md`](references/2026-11-18-behind-pull-request-characterization.md)
- ADR-0435, ADR-0192, ADR-0369, ADR-0311, ADR-0366, ADR-0422, ADR-0302, ADR-0236, ADR-0398
