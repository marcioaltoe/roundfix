---
task: task_05
spec: 0254-a-faster-make-test
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the tests of seven packages
run in parallel under the Parallel Test Package rule, shrinks the sequential
residue of `internal/cli`, narrows two pinned-history fixtures and runs both
contract tags in one `go test`. Every row runs against the built tree,
temporary clones of this repository and the installed Go toolchain. No row
reaches a provider, starts a real Agent Session, reads a credential, or reads
or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST execute every test that task_01, task_03 and task_04 created or named,
   the declared breaks of `_techspec.md` → Testing Approach, and the rule test
   with its seeded cases, and record each result.
3. MUST measure Success Metrics 1 and 2 in two temporary clones of this
   repository, one at the starting commit `2ce5abe8` and one at the audited
   head, each with its own warmed `GOCACHE`. It runs
   `go test -count=1 -parallel 16 -json ./...` in the starting clone and then
   in the head clone, back to back, records each wall time, user and system
   CPU, the load average before and after, every package's elapsed time, and
   the sequential phase of `internal/cli` (the summed elapsed time of its
   top-level tests that never paused). Success Metric 1 holds when the head's
   wall time is at most 65 % of the start's; Success Metric 2 holds when the
   head's `internal/cli` sequential phase is at most 25 % of the start's. When
   the load average differs by more than half between the two runs, it
   repeats the pair and records both.
4. MUST measure Success Metric 3: the rule test passes on the audited head
   and logs the seven packages within the ceilings of `_techspec.md` → API
   Contract 2, and a temporary copy of one converted test file with its
   `t.Parallel()` removed makes the rule fail naming that file and line.
5. MUST measure Success Metric 4: each of the seven packages passes once with
   `-race -short` and once with `-shuffle=on` at the audited head, with the
   shuffle seed recorded.
6. MUST measure Success Metric 5: `go run -buildvcs=false ./cmd/verify-select -contracts -all`
   prints one line naming both tags, and `make verify-contracts` passes in
   each clone, with its wall time recorded for both.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - the installed Go toolchain's `go doc testing.T.Parallel` and
     `go help testflag`, quoting the two passages `_prd.md` quotes;
   - the three GitHub Actions runs, read with `gh run view <id> --log` when
     the network is reachable, comparing their `suite-time budget` lines and
     `internal/cli` times with `_prd.md`;
   - Candido et al., ASE 2017, through the passage `_prd.md` quotes;
   - the Secondbrain mirrors `_prd.md` names, when the Run can read them.

   For each source it reaches, it records what it read and whether the source
   still supports the design. For a source it cannot reach, it records the row
   as blocked with the reason, using
   `blocked (environment: network denied: <host>)` when the Run sandbox denied
   the host. The Go toolchain needs no network.
8. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` defines **Sequential Test** and **Parallel Test Package**, and
   whether any other term the Spec introduced needs one.
9. MUST verify from the repository history:
   - that each Task's changed files stay within its declarations or its
     `## Recorded paths`, and that every recorded path is a test file of a
     package that Task names;
   - that every Governed Path changed is bounded in `_authorization.md`;
   - that the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`
     and every production Go file other than `internal/verifyselect/contracts.go`
     did not change;
   - that in the test files task_01 and task_02 changed, every removed or
     altered line belongs to a file task_03 or task_04 declares, so the
     conversions only inserted `t.Parallel()` lines and `// Sequential:`
     comments (Invariant 1).

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the Unreachable Acceptance criterion of `_prd.md` as not
    measured by this gate, with its `satisfied-by` follow-up.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
12. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT send any request to a provider, and it MUST NOT write outside
    temporary directories and the Spec's `qa/` folder. It reaps any process a
    measurement leaves by the PID it started.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests, the declared breaks and the rule's seeded cases.
- [ ] Measure Success Metrics 1 to 5 in the two clones.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The head's suite and `internal/cli` sequential phase meet Success
      Metrics 1 and 2 against the starting commit on the same machine.
- [ ] The seven packages pass under `-race -short` and `-shuffle=on`.
- [ ] The scope audit finds no changed assertion and no undeclared Governed
      Path.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0254-a-faster-make-test/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- [_prd.md](_prd.md) — Goals; Core Features 1-6; User Stories 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance
- [_techspec.md](_techspec.md) — Current behavior; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Invariants 1-5; Testing Approach; Risks & Considerations; Build Order 5
- ADR-0259; ADR-0089; ADR-0125; ADR-0126; ADR-0213; ADR-0252; ADR-0253
