---
task: task_05
spec: 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the Project Config key
`delivery.numbered_files`, makes Delivery Revalidation park an item whose
added Numbered File takes a number the default branch or an earlier unmerged
queue item already holds as `numbered-file-collision` in the conflict Park
Class, makes a retry of an item that has not run refuse while it collides,
and records the rule in `CONTEXT.md`, ADR-0367, the configuration and
`deliver` guides and the Roundfix `deliver` reference.

Every row runs on this repository's own files, temporary repositories with a
bare remote on disk, temporary Run Databases and temporary copies of the
audited head. No row starts an Agent session or an ACP adapter, reads or
writes the real `~/.roundfix`, reads another repository, reaches the network,
sends a request to a provider, the judge or a reviewer, or reads a credential.
Every `go test` it runs carries an explicit `-timeout`. Every row that runs
`bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary directories,
because the CLI loads `$HOME/.roundfix/config.yml`, passes
`GOMODCACHE="$(go env GOMODCACHE)"` and `GOCACHE="$(go env GOCACHE)"`, runs
with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and removes the
directories afterwards. The QA never writes into the checkout while a test
suite runs, and records environment variable names only, never values.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter (`rustfmt` through rustup) that needs the
session's toolchain configuration under `HOME`. Only rows that run
`bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 by running
   `go test -count=1 -timeout 20m -v` over `./internal/config` with
   `TestNumberedFileDeclarationsAreReadAndValidated`,
   `TestIsNumberedFileNeedsAMatchAndALeadingDigit` and every
   `TestDerivedLine` test, recording each pass and each subtest.
3. MUST verify Success Metrics 2 and 3 and Invariants 2, 3 and 5 by running
   `go test -count=1 -timeout 20m -v` over `./internal/delivery` with
   `TestFindNumberedFileCollisions`, `TestNumberedFileReportRendersTheBlocker`,
   `TestRevalidationParksANumberedFileCollisionBeforeTheRun`,
   `TestStrictFindingsParkBeforeANumberedFileCollision`,
   `TestRetryOfAnItemWithoutARunRefusesWhileANumberedFileCollides`,
   `TestEveryBlockerHasAParkClass`, `TestParkClassesAndNextActions` and every
   `TestRevalidation` and `TestRetryOfAnItemWithoutARun` test, recording each
   pass and each subtest.
4. MUST verify Success Metrics 2 and 4 and Invariants 1 and 4 by running
   `go test -count=1 -timeout 20m -v` over `./internal/cli` with
   `TestRevalidateReportsANumberedFileCollision`,
   `TestDeliverStatusNamesTheNumberedFileCollisionPark` and every
   `TestRevalidate` and `TestDeliverStatus` test, recording each pass and each
   subtest.
5. MUST verify Invariant 1 from the repository history: between the Delivery
   Base (the merge base of the audited head with `main`) and the audited head,
   no existing top-level test function of
   `internal/cli/deliver_revalidate_test.go`,
   `internal/cli/deliver_owner_staleness_test.go` or
   `internal/delivery/revalidate_test.go` changes its body, and
   `TestParkClassesAndNextActions` changes only by the new row and its next
   action check.
6. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails, with
   the case whose input makes the mutation observable.
7. MUST execute the tests task_01, task_02, task_03 and task_04 name in their
   Verification, and record that `GOOS=windows go build -buildvcs=false ./internal/cli ./internal/delivery ./internal/config`
   and `GOOS=linux go build -buildvcs=false ./internal/cli ./internal/delivery ./internal/config`
   succeed.
8. MUST search `docs/user-guide/`, `docs/adr/`, `CONTEXT.md` and both copies
   of the Roundfix skill, and record that: both copies of the `deliver`
   reference carry the paragraph, the conflict row and the action sentence
   task_03 prescribes and are byte-identical; the `deliver` guide carries the
   same paragraph, its conflict row and the action sentence; the
   configuration guide documents `delivery.numbered_files`; `CONTEXT.md`
   defines **Numbered File**, **Numbered File Collision** and the changed
   **Delivery Revalidation**; and ADR-0367 exists with status `accepted`.
9. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the fiscus report quoted in
     `references/2026-11-05-fiscus-delivery-queue-report.md` (2026-11-05):
     compare offline its five Specs that "começaram todas na 0052" with
     `TestRevalidateReportsANumberedFileCollision/ahead_branch` and
     `/ahead_plan`;
   - the drizzle-orm discussion 1104 passage `_prd.md` → Acceptance evidence
     quotes (both branches "created migrations numbered 0007", with the
     snapshot `both added`): compare it with the directory and number rule and
     with a pattern that covers the snapshot directory;
   - the drizzle-orm discussion 5522 passage `_prd.md` → Acceptance evidence
     quotes ("both create a migration with the same sequence number"):
     compare it with detection before the Run rather than at merge.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD and TechSpec declare adds **Numbered File** and
    **Numbered File Collision** and changes **Delivery Revalidation**, that
    the Task commit range writes exactly those entries into `CONTEXT.md`, and
    that `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'`
    passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/`,
      `docs/references/skill-coverage.json`, any existing Archive Record, any
      accepted ADR other than the new ADR-0367, any Baseline asset,
      `docs/agents/`, any Roundfix reference other than `deliver`, any skill
      other than the Roundfix files task_03 declares;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run --plan --offline; rm -rf "$h" "$t"`
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
- [ ] Run the named configuration, engine and revalidation tests.
- [ ] Check the unchanged tests in the history and run each sabotage on a temporary copy.
- [ ] Record the documents, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The key is read and validated, a collision parks before the Run with the exact blocker and next action, a retry refuses while it collides, and nothing changes without the key, in the tests.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0357-a-numbered-file-collision-between-queued-specs-parks-before-the-run/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-7; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-8; Invariants 1-5; Declared breaks; Testing Approach; Build Order 5
