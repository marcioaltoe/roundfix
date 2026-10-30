---
task: task_03
spec: 0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the coverage record
writer name its regeneration command to the suite guard. It makes the
regeneration contract model an archived item sanctioned only by its Archive
Record, with an unsanctioned case. It proves through the queue's own conflict
recovery that an archived item's derived-record conflict merges with the
sanction and is refused without it.

Every row runs on this repository's own files and on temporary repositories
the tests or the QA create. No row:

- reads or writes the real `~/.roundfix`;
- reaches the network or pushes to a remote other than a local bare
  repository;
- sends a request to a provider, the judge or a reviewer;
- reads another repository;
- reads a credential or records an environment value.

Every `go test` the QA runs carries an explicit `-timeout`. Every row that
runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary directories
and removes them afterwards. Such a row keeps the module and build caches
with `GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)"`, so no
cache fills a temporary home. The QA never writes into the checkout while a
test suite runs, and no row compares one command's output with another's
byte for byte.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 4 and Invariants 1, 2 and 5. Run
   `go test -count=1 -timeout 30m -tags docscontract ./internal/config -run '^TestRegenerationIsDeclared$' -v`
   and record two facts:
   - the test passes, including its subtest
     `an_unsanctioned_coverage_write_is_refused`;
   - its log shows every `.roundfixrc.yml` regeneration ran.

   Also run
   `go test -count=1 -timeout 20m ./internal/spec -run '^TestCoverageEquivalence$'`
   and record the pass.
3. MUST verify Success Metric 3 and Invariants 3 and 4. Run
   `go test -count=1 -timeout 20m ./internal/cli -run '^(TestAnArchivedItemsDerivedRecordConflictIsMergedAndRegenerated|TestAnArchivedItemWithoutTheSanctionKeepsItsDerivedRecordConflict|TestADerivedConflictIsMergedAndRegenerated|TestConflictRecoveryAbortsAFailedRegeneration)$' -v`
   and record each pass.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`. Run every `go test` with `-timeout`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - Remove the `suiteguard.DeclareSanctionedRegeneration` call from
     `TestCoverageEquivalence`. `TestRegenerationIsDeclared` fails, because
     its stale coverage record makes the coverage command write and exit 1.
     `TestAnArchivedItemsDerivedRecordConflictIsMergedAndRegenerated` fails
     with `regenerate derived paths`.
   - Make `removeSanctionedRegeneration` (`internal/suiteguard/suiteguard.go`)
     return no remaining violation whenever any declaration exists.
     `an_unsanctioned_coverage_write_is_refused` fails because the command
     exits 0. Its record lists the module record command and
     `make baseline-digests`, so declarations exist.
     `TestAnArchivedItemWithoutTheSanctionKeepsItsDerivedRecordConflict`
     fails because the merge commits. Its record lists
     `make baseline-digests`.
   - Write the archived item's record with the coverage command even when the
     unsanctioned case asks for none. `an_unsanctioned_coverage_write_is_refused`
     fails because the command exits 0.
   - Stop the regeneration contract's fixture from removing the copied
     `docs/specs/` tree. In the audited head, this Spec's own
     `_authorization.md` lists the coverage command, so
     `an_unsanctioned_coverage_write_is_refused` fails because the command
     exits 0.
5. MUST execute the tests task_01 and task_02 name in their Verification.
6. MUST search the repository for `regeneration_fixture_test.go` and for
   `DeclareSanctionedRegeneration` under `internal/config/`, excluding this
   Spec's folder and `docs/history/`. Record that neither remains, so no
   fixture injects a declaration (Invariant 5).
7. MUST record the outside-evidence row from two sources this Spec did not
   write. The row is never blocked.
   - The Go command documentation for `go test` (`cmd/go`,
     https://pkg.go.dev/cmd/go/internal/test, retrieved 2026-10-30), quoting
     "Tests that open files within the package's module or that consult
     environment variables only match future runs in which the files and
     environment variables are unchanged."
   - The park of Pull Request #709 quoted in the adopted Backlog Entry under
     `references/` (2026-10-27): the queue's regeneration log ended with
     `suiteguard: modified: docs/references/coverage-record.json`.

   Compare, offline, the refusal line that a passing
   `TestAnArchivedItemWithoutTheSanctionKeepsItsDerivedRecordConflict` reads
   from its regeneration log with that park. Confirm that both fixtures copy
   the repository in-process, as the documentation's rule requires for the
   cache to see the writer. Name the Go documentation and the Pull Request
   number with its date as the sources.
8. MUST record Success Metric 5 as exactly one row,
   `blocked (declared unreachable)`, against the PRD's single Unreachable
   Acceptance entry, and count it once in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that
   **Archive Record** in `CONTEXT.md` cites ADR-0345 and no Spec, and that
   `docs/agents/specific-repository.md` names
   `suiteguard.DeclareSanctionedRegeneration` and ADR-0345.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - `docs/references/coverage-record.json` changed only in task_02's
      commit;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `internal/suiteguard/`, `internal/suiteguardcontract/`,
      `internal/cli/deliver_workflow.go`, anything under
      `internal/baseline/assets/`, and any skill;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the search, the outside source, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every `.roundfixrc.yml` regeneration passes under an archived item's record alone, and the coverage write is refused when that record does not list its command.
- [ ] The queue's conflict recovery merges an archived item's coverage-record conflict, and leaves the item unchanged without the sanction.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0332-the-delivery-queue-regenerates-the-derived-records-of-an-archived-spec/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-5; Testing Approach; Build Order 3
