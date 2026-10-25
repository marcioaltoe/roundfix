---
task: task_03
spec: 0313-internal-cli-tests-stop-reading-the-live-documentation
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec moves thirteen
`internal/cli` tests that read this repository's guides, skills and records
into two `docscontract` files declared `//verify:always` (ADR-0326), so a
documentation edit no longer re-runs the whole package, and classifies those
files in the repository contract declarations test. Its corrective Task,
task_04, answers F-01 of the first Run's gate: the upgrade characterization
test, a fourteenth reader, now runs in a fixture workspace. task_05, the
second corrective Task, makes the suite guard fingerprint the repository
through Git and keeps every leak class it reported. This gate runs after
task_05, so the cache, test input log and sabotage rows measure the tree
with all four implementation Tasks.

Every row runs on temporary worktrees or copies of the audited head and of
the Delivery Base and this repository's own documents. No row reads or
writes `~/.roundfix`, reaches the network, sends a request to a provider, the
judge or a reviewer, reads a credential, or reads another repository. Every
`go test` it runs carries an explicit `-timeout`, and every row that runs
`bin/roundfix` sets `HOME` to a fresh temporary directory, because the CLI
loads `$HOME/.roundfix/config.yml`. A row that runs `go test` with a fresh
temporary `HOME` keeps the existing Go module cache and Rust toolchain
reachable through `GOMODCACHE`, `GOPATH`, `RUSTUP_HOME` and `CARGO_HOME`,
and records none of their values; the first Run's gate lost a cache series
to a `rustfmt` shim without a toolchain.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Invariant 2: for the thirteen tests
   `TestBaselineUpdatePreviewListsADifferenceWithARepositoryProfile`, `TestBaselineUpdateReachesCurrentWithARepositoryProfile`, `TestBaselineUpdateSkillsStageRefreshesSkillsWithARepositoryProfile`, `TestDoctorComparesARepositoryProfileWithItsSnapshot`, `TestBaselineUpdatePreviewListsASetupSnapshotDifference`, `TestARepositoryMissingARequiredExternalSkillIsReported`, `TestThisRepositoryHoldsEveryRequiredExternalSkill`, `TestBaselineExamplesParse`, `TestProjectConstraintPlanningExamplesParse`, `TestEventsHelpDocumentsAgentSelectionFilter`, `TestEveryStatementOfTheHTTPDefaultNamesTheCatalogDefault`, `TestReleasePlanChecksOnAnAdoptedRepository` and `TestReleasePlanChecksWriteNothing`, record that `go test -list` of `./internal/cli` without tags
   lists none of them, that `go test -count=1 -timeout 20m -tags docscontract
   -v` runs and passes each, and, with `git diff` between the Delivery Base
   and the audited head, that each function body is byte-identical to its
   Delivery Base body.
3. MUST verify Success Metric 2 and Invariant 3: run
   `go run ./cmd/verify-select -contracts -all` and record that the thirteen
   are listed with `./internal/cli`; run
   `TestRepositoryContractTestsDeclareTheirRelevance` and record its pass.
4. MUST verify Success Metric 3 in two temporary worktrees, one at the
   audited head and one at the Delivery Base, each made with
   `git worktree add` outside the checkout: run `go test -parallel 16
   -timeout 20m ./internal/cli` twice in the same shell environment (the
   second run must print `(cached)`), append one line to
   `docs/user-guide/commands/gc.md` and to `.agents/skills/roundfix/SKILL.md`,
   run it again, and record the three result lines of each worktree. The row
   passes when the audited head's third run prints `(cached)`; record the
   Delivery Base's re-run as context. Every run of a pair uses the same
   environment: `GODEBUG` and the other variables are part of Go's cache key.
   The row runs on the head that includes task_04 and task_05; the first
   Run's head without them re-ran the package (F-01).
5. MUST verify Success Metrics 4 and 5 and Invariant 1. For Success Metric 5,
   in a temporary worktree of the audited head, build the untagged test
   binary with `go test -buildvcs=false -c -o <temporary path>
   ./internal/cli`, run it once from the worktree's `internal/cli` with a
   fresh temporary `HOME`, `-test.parallel 16`, `-test.timeout 20m` and
   `-test.testlogfile <temporary path>`, and record that it passes and that
   no `open` or `stat` line of the log names a path under the worktree's
   `docs/`, `.agents/`, `skills/`, `CONTEXT.md` or `.roundfixrc.yml`, with
   both the worktree path and its resolved form. Record the distinct
   in-module paths the log names, grouped by top-level directory. For
   Success Metric 4, in a temporary copy of the
   audited head made with `git archive`, initialized as a Git repository and
   committed, remove the five directories `docs/history/`, `docs/adr/`,
   `docs/specs/`, `docs/user-guide/` and `.agents/skills/` from disk without
   committing (Git keeps them tracked and reports them deleted; the suite
   guard accepts a tracked path that is missing from disk but cannot inspect
   a tracked unreadable one, Run run_20261025T193403Z_ade2bddbacb762cf F-02),
   run `go test -count=1 -parallel 16 -timeout 20m ./internal/cli`, record
   that it passes, and remove the copy.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, initialized as a Git repository and committed before
   the mutation (the suite guard refuses a root outside a Git work tree),
   and record that the named check fails, with the case
   whose input makes the mutation observable:
   - delete the `//go:build docscontract` line of
     `internal/cli/repository_profile_docscontract_test.go`: task_01's
     untagged `go test -list` check prints the five tests, and
     `TestRepositoryContractTestsDeclareTheirRelevance` fails because the
     file is no longer a contract;
   - replace `//verify:always` by `//verify:relevant` in
     `internal/cli/live_documents_docscontract_test.go`:
     `TestRepositoryContractTestsDeclareTheirRelevance` fails on the class
     of its eight tests;
   - move `TestReleasePlanChecksWriteNothing` back to
     `internal/cli/releaseplan_checks_test.go`: task_02's untagged list check
     prints it, and the Success Metric 3 row run on that copy re-runs the
     package after the guide edit;
   - delete the `withCLIWorkspace(t)` line that task_04 added to
     `TestUpgradeStdoutAndExitCodesAreCharacterized` in
     `internal/cli/upgrade_characterization_test.go`: task_04's Verification
     fails and names paths under `docs/`, and the Success Metric 5 row run on
     that copy names them too. The test's cases reach `roundfix upgrade`'s
     history notice, which inventories the copy's own `docs/history` and
     `docs/specs`, so the mutation is observable;
   - in `internal/suiteguard/suiteguard.go`, each of these four, on its own
     copy, makes task_05's Verification fail on the named leak cases of
     `TestGuardReportsEveryLeakClass`, whose fixture holds an untracked file,
     an empty untracked directory and tracked files with distinct content:
     drop `--others` from the first `git ls-files` call (the seven untracked
     cases, among them "untracked file created"); drop `--directory` from the
     second (the "empty directory created", "pre-existing empty directory
     removed" and "new directory holding only ignored files" cases); leave
     the mode of each state unset in the `os.Lstat` pass (the "tracked file
     permission narrowed" and "tracked directory permission narrowed"
     cases); drop `--remove` from `git update-index` (the "tracked file
     deleted" case);
   - restore the hand-written walk, so that `fingerprint` digests files with
     `os.Open` and accepts any directory:
     `TestGuardRefusesARootOutsideAGitWorkTree` fails.
7. MUST execute the tests task_01, task_02, task_04 and task_05 name in
   their Verification, record that the 22 leak cases pass at the audited
   head, and run every package that installs `suiteguard.Main` (the
   `guardedSpawningPackages` list of `internal/suiteguardcontract`) once at
   the audited head with `-count=1` and an explicit `-timeout`, recording
   that each passes and one `guard cost: measured` line.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - run `go help test` in the audited checkout and record the sentence the
     PRD quotes;
   - read the audit reference and record its 162.6 s and 237.0 s re-runs and
     the cached result with the thirteen skipped, as committed by Spec
     0311's planning change.
9. MUST perform the glossary check: record that the Spec declares no term,
   that `CONTEXT.md` is byte-identical to the Delivery Base, and whether any
   term the Spec introduced needs an entry (`docs/agents/domain.md`).
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every non-test `.go`
      file except `internal/suiteguard/suiteguard.go`, every file under
      `internal/docscontract/`, `internal/suiteguardcontract/`, `.agents/`,
      `skills/` and `docs/user-guide/`, every caller of `suiteguard.Main`,
      and `CONTEXT.md`;
    - task_05 changed only the four `internal/suiteguard` files in
      `_authorization.md` `paths`, besides its Task file;
    - task_04 changed only `internal/cli/upgrade_characterization_test.go`
      besides its Task file, by one added line;
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
    `HOME="$(mktemp -d)" bin/roundfix archive 0313-internal-cli-tests-stop-reading-the-live-documentation --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the list, tag, selector, cache, test input log and removed-trees rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The thirteen tests run only under the tag, unchanged, and the selector lists them.
- [ ] A guide and skill edit leaves the `internal/cli` result cached at the audited head, the untagged test binary's input log names no documentation path, and the untagged package passes with the trees unreadable.
- [ ] Each sabotage makes its named check fail, every guarded package passes with the Git fingerprint, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0313-internal-cli-tests-stop-reading-the-live-documentation/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Sources; Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-3; Invariants 1-4; Testing Approach; Build Order 5
