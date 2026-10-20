---
task: task_04
spec: 0289-gc-sanitize-reaches-the-default-artifact-roots
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the storage report's
resolver the one resolver of a Run's Artifact Root, moves Run Retention onto
it, and makes `roundfix gc sanitize` classify the root each Run wrote to with
its existing proofs, check a directory proven only by absence again before
removing it, and name each resolved root and what it would reclaim.

Every row runs on temporary repositories and homes. No row reads or writes the
real `~/.roundfix`, reaches the network, sends a request to a provider, the
judge or a reviewer, reads a credential, or reads another repository. The gate
never runs `roundfix spec judge`, never runs `roundfix review` against a real
provider, never runs an archive plan without `--offline`, and never runs
`roundfix gc sanitize --apply` outside a temporary home.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run `TestGCSanitizeDryRunReachesTheDefaultRootTranscript`,
     `TestGCSanitizeReachesTheDefaultRootAndKeepsEveryProof` and
     `TestDiscoverArtifactRootsGroupsRunsByTheRootTheyWroteTo`, and record
     each removed and kept path (Invariants 3, 4 and 6).
   - Through the built product, reproduce Surface Transcript 1 of
     `_techspec.md`. In a temporary copy of the audited head made with `git
     archive`, seed a temporary home outside every checkout with a throwaway
     test that opens the Run Database there, creates one Implement Run with no
     Artifact Directory for a temporary Git repository with one commit,
     completes it `clean`, sets its completion 400 hours before now, and
     writes `runs/<run-id>/agent/batch-001.log` under that repository's
     default root. Then, with `HOME` set to that home and the temporary
     repository as the working directory, run `bin/roundfix gc sanitize` and
     record that standard output, standard error and the exit code match the
     transcript. Then run `bin/roundfix gc sanitize --apply` twice and record
     that the first removes that one directory and the second reports
     `Directories removed: 0`. Compare the file lists under the home's
     `.roundfix/artifacts` and under the repository before and after each
     run; the Run Database and its `-wal` and `-shm` files are not compared.
3. MUST verify Success Metric 2:
   - Run `TestGCSanitizeKeepsADirectoryWhoseRunAppearsBeforeRemoval`,
     `TestGCSanitizePreservesARunWhoseArtifactRootCannotBeResolved`,
     `TestRunGCSanitizeClassifiesEveryRecordedRootAndMutatesOnlyProvenDirectories`,
     `TestGCSanitizeKeepsTheSharedRootAfterWorktreeRemoval` and the five tests
     in `internal/cli/gc_sanitize_pre_key_root_test.go`, and record each kept
     case.
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy and record that the named test fails:
     - make `DiscoverArtifactRoots` group Runs by their recorded value again:
       `TestGCSanitizeDryRunReachesTheDefaultRootTranscript` and
       `TestDiscoverArtifactRootsGroupsRunsByTheRootTheyWroteTo` fail;
     - make `classifyGCSanitationRoot` accept a root that is not every Run's
       default (skip the `overridden` result):
       `TestGCSanitizeReachesTheDefaultRootAndKeepsEveryProof` fails;
     - remove the Run Database read before an absence-proven removal:
       `TestGCSanitizeKeepsADirectoryWhoseRunAppearsBeforeRemoval` fails;
     - make an `Unresolved` group skip its proof:
       `TestGCSanitizePreservesARunWhoseArtifactRootCannotBeResolved` fails.
4. MUST verify Success Metric 3:
   - Run `TestResolveRunArtifactRootIsTheRootTheRunWroteTo`,
     `TestRunRetentionReadsTheRootTheSharedResolverReturns`,
     `TestStorageReportMeasuresTheArtifactRootARunWroteTo`,
     `TestRunStorageReportListsTheDefaultRootOfARunWithoutOne`,
     `TestGCRemovesRunsPastRunRetention` and
     `TestGCKeepsARunUnderAnUnprovenArtifactRoot`, and record each subtest's
     root (Invariants 1, 2 and 5).
   - On a temporary copy of the audited head, make `runRetentionArtifact` use
     the recorded root as is again, and record that
     `TestRunRetentionReadsTheRootTheSharedResolverReturns` fails.
   - Record that `internal/store/journal.go` defines `ResolveRunArtifactRoot`
     and no `storageArtifactRoot`, and that `internal/cli/run_retention.go`
     calls neither `ResolveArtifactDirectory` nor
     `DefaultArtifactDirectoryForPath`.
   - Through the built product: seed a second temporary home as in
     Requirement 2 with one Run recording no root, one recording a relative
     root and one recording the default root as an absolute path. Build the
     Delivery Base's binary from a `git archive` of the Delivery Base, run
     both binaries' `storage report` against that home, and record that the
     `Artifact Roots` sections are identical.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - quote from the adopted Backlog Entry and `_techspec.md` → Current
     behavior the maintainer's dry run of 2026-10-20 (root `""`, `unsafe`, 0
     directories, 0 bytes);
   - record that `internal/cli/gc_test.go` and
     `internal/cli/gc_sanitize_pre_key_root_test.go` are byte-identical to the
     Delivery Base, that the bodies of `TestGCRemovesRunsPastRunRetention` and
     `TestGCKeepsARunUnderAnUnprovenArtifactRoot` in
     `internal/cli/gc_run_retention_test.go` are identical to the Delivery
     Base, and that all of those tests pass at the audited head;
   - because this Spec's authorization forbids network access, record each
     web source as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
7. MUST record Success Metric 4 as exactly one row, `blocked (environment:
   the real Roundfix Home is out of the Spec's reach)`, naming the PRD's one
   Unreachable Acceptance entry: the maintainer's dry run and conditional
   apply on the real `~/.roundfix` after the release. Count that single row
   once in `rows_blocked_declared`; do not split the dry run and the apply
   into two rows. Never run either command against the real home.
8. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Artifact Root** with its resolution
     rules and links to ADR-0301 and ADR-0255;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, every file under `internal/config/`,
     and `docs/user-guide/commands/storage-report.md`;
   - the Roundfix skill version was raised and recorded;
   - no `### QA settlement` section changed;
   - no guide, agent guide or skill copy still says that `gc sanitize`
     discovers every recorded Artifact Root.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0289-gc-sanitize-reaches-the-default-artifact-roots --plan --offline`
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
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Reproduce Surface Transcript 1, the applies and the storage report comparison through the built binaries on temporary homes.
- [ ] Record the outside evidence, the single unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The dry run names the default root of a Run that recorded none with what it would reclaim, exactly as Surface Transcript 1, and the applies remove only what it named.
- [ ] Every existing proof keeps its paths on the resolved root, and each removed guard makes its named test fail.
- [ ] One resolver serves the storage report, sanitation and Run Retention, and the storage report's Artifact Roots section matches the Delivery Base's.
- [ ] The maintainer's real dry run is one declared unreachable row, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0289-gc-sanitize-reaches-the-default-artifact-roots/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; Surface Transcript 1; Invariants 1-6; Current behavior; Build Order 4
- ADR-0301; ADR-0255; ADR-0296; ADR-0033; ADR-0104; ADR-0240; ADR-0276
