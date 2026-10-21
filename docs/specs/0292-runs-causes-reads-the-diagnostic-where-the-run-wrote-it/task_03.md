---
task: task_03
spec: 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes `roundfix runs
causes` read each failed attempt's diagnostic under the Artifact Root the Run
wrote to, through the one resolver the storage report, GC sanitation and Run
Retention use, keeps every guard on the diagnostic, and describes the change
in the `runs` guide, the Roundfix `runs` reference and the glossary.

Every row runs on temporary repositories and homes. No row reads or writes the
real `~/.roundfix`, reaches the network, sends a request to a provider, the
judge or a reviewer, reads a credential, or reads another repository. The gate
never runs `roundfix spec judge`, never runs `roundfix review` against a real
provider, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run `TestRunsCausesReadsTheDiagnosticOfARunWithoutAnArtifactRootTranscript`
     and the `blank root` subtest of
     `TestBuildReadsTheDiagnosticUnderTheResolvedArtifactRoot`, and record the
     class and signature of each attempt (Invariant 1).
   - Through the built product, reproduce Surface Transcript 1 of
     `_techspec.md`. In a temporary copy of the audited head made with `git
     archive`, seed a temporary home outside every checkout, with its path's
     symbolic links resolved, using a throwaway test that opens the Run
     Database there, creates one Implement Run of Spec `0300-example` with no
     Artifact Directory for a temporary Git repository with one commit whose
     `docs/specs/0300-example/` holds the fixture's `_tasks.md` and
     `task_04.md`, appends the fixture's events with absolute
     `diagnostic_path` values under that repository's default root, writes the
     three diagnostics there, completes the Run `clean` and sets its creation
     time to 2026-10-01T12:00:00Z. Then, with `HOME` set to that home and the
     temporary repository as the working directory, run `bin/roundfix runs
     causes --since 2026-10-01` and record that standard output, standard
     error and the exit code match the transcript. Compare the file lists and
     hashes under the home's `.roundfix/artifacts` and under the repository
     before and after the run; the Run Database and its `-wal` and `-shm`
     files are not compared.
   - Build the Delivery Base's binary from a `git archive` of the Delivery
     Base, run it the same way against the same home, and record that it
     prints those three attempts `unclassified`, which is the defect this Spec
     fixes.
3. MUST verify Success Metric 2: run
   `TestRunsCausesReadsTheDiagnosticUnderARelativeArtifactRoot` and the
   `relative root` subtest, and record each attempt's class.
4. MUST verify Success Metric 3:
   - Run `TestBuildReadsTheDiagnosticUnderTheResolvedArtifactRoot` whole and
     `TestBuildIgnoresALogOutsideTheArtifactDirectory`, and record for each
     guard subtest that task_01 attempt 1 is `unclassified` with nil evidence,
     that `blank root without home` gives no error, and that `inspection
     error` returns an error naming the Run (Invariants 2-4).
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy, each keeping the package compiling, and record
     that the named test fails:
     - make `diagnosticTail` use the recorded Artifact Directory literally
       again: `TestRunsCausesReadsTheDiagnosticOfARunWithoutAnArtifactRootTranscript`
       and the `blank root` subtest fail;
     - drop the `os.Lstat` component walk: the `parent symlink` subtest fails;
     - drop the `filepath.Rel` containment: the `other run` subtest fails;
     - make an unresolved root fail the report instead of giving no evidence:
       the `blank root without home` subtest fails.
   - Run `TestBuildBoundsEvidenceAndExportsNoPrivateCommandText` and
     `TestRunsCausesPrintsJSON`, and record that no diagnostic text, home path
     or key reaches stdout or JSON (Invariant 5).
5. MUST verify Success Metric 4: record that `internal/runcause/report.go`
   calls `store.ResolveRunArtifactRoot`, that no production (non-test) file under
   `internal/runcause` calls `ResolveArtifactDirectory`,
   `DefaultArtifactDirectoryForPath` or `ValidateArtifactDirectory`, and that
   the storage report, `DiscoverArtifactRoots` and Run Retention still call
   the same resolver.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - quote from `docs/references/corrective-causes-measurement.md` its window,
     its table digest, the `unclassified` count 100 of 132, the `environment`
     count 0 and its five observed signatures, and record that the digest it
     names equals the digest `roundfix runs causes` prints at the audited head
     and that none of the seven diagnostic-only signatures of
     `internal/runcause/signatures.json` appears among the observed ones;
   - record that the bodies of `TestBuildIgnoresALogOutsideTheArtifactDirectory`,
     `TestBuildReportsAttemptsTasksAndCorrectives`,
     `TestBuildBoundsEvidenceAndExportsNoPrivateCommandText`,
     `TestRunsCausesPrintsAClassifiedWindow` and `TestRunsCausesPrintsJSON`
     are identical to the Delivery Base, and that all of them pass at the
     audited head;
   - because this Spec's authorization forbids network access, record the Go
     article as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
8. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Artifact Root** naming the run cause
     report among the resolver's callers, with links to ADR-0301, ADR-0255
     and ADR-0214;
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
     every file under `internal/store/`, `internal/runcause/signatures.json`,
     `docs/references/skill-coverage.json` and
     `docs/references/corrective-causes-measurement.md`;
   - the Roundfix skill version was raised and recorded;
   - no `### QA settlement` section changed;
   - the `runs causes` help text is unchanged from the Delivery Base's binary.

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
    `bin/roundfix archive 0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it --plan --offline`
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
- [ ] Reproduce Surface Transcript 1 through the built binary and compare it with the Delivery Base's binary on the same temporary home.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The report classifies the attempts of a Run that recorded no Artifact Directory from their diagnostics, exactly as Surface Transcript 1, where the Delivery Base's binary left them unclassified.
- [ ] Every guard holds on the resolved directory, and each removed guard makes its named test fail.
- [ ] One resolver serves the storage report, sanitation, Run Retention and the run cause report, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0292-runs-causes-reads-the-diagnostic-where-the-run-wrote-it/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Invariants 1-5; Current behavior; Build Order 3
- ADR-0301; ADR-0214; ADR-0104; ADR-0240; ADR-0276
