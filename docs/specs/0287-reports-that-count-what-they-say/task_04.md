---
task: task_04
spec: 0287-reports-that-count-what-they-say
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers three fixes: the
Linux Lane's pass line counts the packages `go list` selected in the
container, `roundfix storage report` measures each Run under the Artifact Root
it wrote to, and an Archive Record summary carries the list or paragraph its
colon introduces.

Every row runs on temporary repositories and homes with compiled Docker and Go
stand-ins, except the one real lane row, which runs only when the Run sandbox
reaches a Docker daemon. No row reads or writes the real `~/.roundfix`,
reaches the network, pulls an image, sends a request to a provider, the judge
or a reviewer, reads a credential, or reads another repository. The gate
never runs `roundfix spec judge`, never runs `roundfix review` against a real
provider, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run `TestLinuxLaneReportsThePackagesTheContainerSelected`,
     `TestDeliverLanePackageCountTranscript`,
     `TestDeliverLaneNoPackageCountTranscript` and
     `TestLinuxLaneRunsTheChangedPackagesInAContainer`, and record the
     counted, zero, no-count and later-line cases (Invariants 1 and 2).
   - On a temporary copy of the audited head made with `git archive`, make the
     pass lines print `len(result.Packages)` again and record that
     `TestLinuxLaneReportsThePackagesTheContainerSelected` fails; on another,
     remove the script's count line and record that it fails; on a third,
     make the host read the last count line instead of the first and record
     that it fails.
   - Through the built product, reproduce Surface Transcripts 1 and 2 of
     `_techspec.md`: in a temporary repository whose `go.mod` declares `go
     1.26.1` and whose `.roundfixrc.yml` sets `delivery.linux_lane: false`,
     with `HOME` set to a temporary directory and compiled Docker and Go
     stand-ins first on `PATH`, run `bin/roundfix deliver lane --all` once
     with a Docker stand-in that writes `roundfix-linux-lane: packages: 43`
     on `run` and once with one that writes no count line, and record that
     standard output, standard error and the exit code match each transcript.
3. MUST verify Success Metric 2:
   - Run `TestStorageReportMeasuresTheArtifactRootARunWroteTo` and
     `TestRunStorageReportListsTheDefaultRootOfARunWithoutOne`, and record
     each subtest's root, status and bytes (Invariants 3 and 4).
   - On a temporary copy of the audited head, make `storageArtifactRoot`
     return the recorded value cleaned, as before, and record that both tests
     fail.
   - Through the built product: in a temporary copy of the audited head, seed
     a temporary home outside every checkout with a throwaway test that opens
     the Run Database there, creates an Implement Run with no Artifact
     Directory for a temporary repository and writes bytes under that Run's
     `runs/<run-id>` in the repository's default root; then run `bin/roundfix
     storage report` with `HOME` set to that home and record the Artifact
     Roots lines, that the root is `present` with those bytes, that no `.`
     root is listed, and that the home's file list and bytes are identical
     before and after the report.
4. MUST verify Success Metric 3:
   - Run `TestArchiveSummaryCarriesWhatItsColonIntroduces`,
     `TestAnArchiveRecordEndingWithAColonRoundTrips`,
     `TestACancelledRecordSummaryCarriesItsList` and
     `TestLegacySummaryCarriesTheListItsColonIntroduces`, and record each
     case's summary (Invariants 5 and 6).
   - On a temporary copy of the audited head, remove the carry from
     `archiveParagraph` and record that the first, third and fourth fail.
   - Through the built product: in a temporary Git repository with one commit
     holding `docs/specs/0001-demo/_prd.md` (status `active`, title paragraph
     "Two gaps failed it:" followed by a two-item list), a `_tasks.md` with
     `qa: declined` and one `task_01.md` with a `## Verification` bullet, and
     `HOME` set to a temporary directory outside that repository, run
     `bin/roundfix cancel 0001-demo --reason qa` and record the summary line
     of `docs/history/specs/0001-demo.md`, which holds both items.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - count the Archive Records under `docs/history/specs/` whose last line
     ends with a colon at the Delivery Base, and record the number and that
     no file under `docs/history/` differs between the Delivery Base and the
     audited head;
   - record from the adopted Backlog Entries the operator's lane run (43
     packages, line said 1) and the maintainer's storage report (`.`
     `missing`, 508 Run rows, 0 bytes);
   - because this Spec's authorization forbids network access, record each
     web source as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
7. MUST verify Success Metric 4 when the Run sandbox can, and otherwise record
   the row as declared unreachable:
   - When `docker info` exits `0` and `docker image inspect golang:1.26`
     finds the image, run `bin/roundfix deliver lane --all` in the audited
     head's worktree with `HOME` set to a temporary directory and
     `GOMODCACHE` set to the value `go env GOMODCACHE` printed before `HOME`
     changed. Record its outcome line, its exit code, the lane log's
     `roundfix-linux-lane: packages:` line and the number of `ok` and `?`
     result lines in the log. Record environment variable names only, never
     their values. Write that evidence only after the lane ends.
   - Otherwise record that the lane half is blocked (environment: Docker
     unreachable from the Run sandbox).
   - The storage report half never runs against the real `~/.roundfix`.
   - Record Success Metric 4 as ONE row, declared unreachable, holding both
     halves: the lane evidence or its block, and the storage report block.
     It counts once in `rows_blocked_declared`, because the PRD's
     Unreachable Acceptance declares one criterion; two declared rows exceed
     that declaration and the Daemon refuses the `partial`.
8. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` carries the revised **Archive Record** with the
     summary rule;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, `internal/store/run_retention.go`,
     `internal/cli/run_retention.go`, and `DiscoverArtifactRoots` and
     `classifyGCSanitationRoot` by content;
   - the Roundfix skill version was raised and recorded;
   - no `### QA settlement` section changed;
   - no guide, agent guide or skill copy still says the pass line counts
     package arguments.

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
    `bin/roundfix archive 0287-reports-that-count-what-they-say --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - start a workflow or pull an image.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Reproduce Surface Transcripts 1 and 2, the seeded storage report and the cancelled record through the built binary.
- [ ] Run the real lane when Docker is reachable, record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The lane's pass line prints the container's count or says it has none, and each removed guard makes its named test fail.
- [ ] The storage report lists the default root of a Run that recorded none as present with its bytes and changes nothing.
- [ ] A summary carries the list its colon introduces in every writer, and no Archive Record already written changed.
- [ ] The real lane run passes or is declared unreachable, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0287-reports-that-count-what-they-say/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Surface Transcript 1; Surface Transcript 2; Invariants 1-6; Current behavior; Build Order 4
- ADR-0294; ADR-0297; ADR-0296; ADR-0255; ADR-0247; ADR-0277; ADR-0262; ADR-0104; ADR-0240; ADR-0276
