---
task: task_03
spec: 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
status: completed
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the delivery item binary
tests report how their probe ended, records a diagnosis of the F6 signal, and
makes every delivery test replace its item binary by copy and rename, so ten
full suites pass. No row reads or writes the real `~/.roundfix`, reaches a
provider, reads a credential or reads another repository.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1, the acceptance ADR-0266 decides: at the
   audited head run `go test -count=1 -parallel 16 ./...` ten consecutive
   times on one machine, with `HOME` set to a temporary directory and
   `RUSTUP_HOME` and `CARGO_HOME` set to the operator's existing rustup and
   cargo homes, which the suites only read. Each suite runs the command
   exactly, without `-v` or other added flags, writes its raw output only to a
   temporary directory outside the checkout and is copied into the evidence
   directory after it ends. For each suite record its start time, load
   average, wall time, exit status and the count of output lines matching
   `was killed by signal`, `signal: killed`, `exited -1` or
   `suiteguard: helper residue`. The row passes only when all ten exit 0 with
   every count zero. On the first failing suite it stops, keeps the output,
   records the failing test and the child's Signal Exit, and fails the row; it
   never reruns a failed suite to obtain a pass. From the start of the first
   suite until the tenth ends, the QA Agent creates, changes or removes no
   file under the checkout, the QA directory included: the suite guard
   fingerprints the whole checkout, so a QA file written mid-suite fails the
   suite with a guard violation, as the first gate of 2026-10-17 did. Notes,
   timings and outputs go to the temporary directory, and the evidence
   directory and the QA Report are written only after the last suite ends.
3. MUST, when a suite of Requirement 2 shows a Signal Exit, read
   `~/Library/Logs/DiagnosticReports` in the Agent turn for a
   `roundfix-*.ips` report written during that suite with parent process
   `cli.test`, and record its termination namespace and indicator. With the
   namespace `CODESIGNING` the row fails as a fix that missed a site. With no
   such report it records the PRD's Unreachable Acceptance as
   `blocked (sender: process signal)`, quoting its `eslogger` command
   verbatim, so the Supervisor sends it to the maintainer.
4. MUST verify Success Metric 2 by sabotage: in a temporary worktree at the
   audited head, change only `deliveryItemBuildCommand` back to copying onto
   `bin/roundfix`, run
   `go test -count=1 -run '^TestDeliveryStepProbesARebuiltItemBinaryWhileThePreviousBuildStillRuns$' ./internal/cli`
   and record that it fails with the probe notice "was killed by signal
   SIGKILL (9), no core dump" and that a `CODESIGNING` diagnostic report was
   written for it; restore the builder, run it again and record that it
   passes; then remove the worktree.
5. MUST verify Success Metric 3: record the four lines of
   `diagnosis/f6-sender.md` and that its `sender` follows API Contract 3 from
   the suites and reports it lists. When it reads `process signal`, Requirement
   3's Unreachable Acceptance applies even if Requirement 2 passed.
6. MUST verify Invariant 4 and API Contract 6 against the delivery base, the
   merge base of the audited head with `main`: record that
   `git diff --name-only <delivery base> <audited head>` lists only
   `internal/cli/deliver_item_binary_test.go`,
   `internal/cli/deliver_item_binary_signal_unix_test.go`,
   `internal/cli/deliver_item_probe_report_test.go`,
   `internal/cli/deliver_item_binary_rebuild_test.go` and files of this Spec
   folder, so every production file is byte-identical to the delivery base;
   and that no test file of the repository copies an item binary onto an
   existing `bin/roundfix`, by repeating the sweep of task_02 over every
   `_test.go` file.
7. MUST execute the tests task_01 and task_02 name in their Verification.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the Apple "Updating Mac Software" passage
   on the in-kernel signature cache and the replacement by a temporary file,
   the Open Radar FB8914243 passage on the kernel's `SIGKILL`, and the
   `eslogger` manual page passage on super-user and Full Disk Access, quoting
   each and stating whether it still supports the design; the manual page is
   read locally with `man eslogger`. For a web source it cannot reach it
   records `blocked (environment: network denied: <host>)`.
9. MUST perform the glossary check of `docs/agents/domain.md`: record that the
   Spec declares no term and whether any term it uses needs one.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`; that no Governed
    Path changed, as `paths: []` in `_authorization.md` requires; and that the
    Makefile, the CI workflows, `go.mod`, `go.sum`, every skill and every
    Baseline module did not change. This row reads Task commits, so it
    declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, run `sudo` or `eslogger`, record any
    environment value, send any request to a provider, or start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the ten full suites and read any diagnostic report.
- [ ] Prove the rebuild test fails with the in-place build and passes with the replacement.
- [ ] Record the diagnosis, the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Ten consecutive full suites pass with no Signal Exit and no Helper Residue.
- [ ] The rebuild test fails when the build writes in place and passes with the replacement.
- [ ] No production file differs from the delivery base.
- [ ] The scope audit finds no Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Story 1; User Story 2; Core Feature 1; Core Feature 2; Core Feature 3; Success Metric 1; Success Metric 2; Success Metric 3; Unreachable Acceptance; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Invariants 1-5; Current behavior; Build Order 3
- ADR-0281; ADR-0266; ADR-0213
