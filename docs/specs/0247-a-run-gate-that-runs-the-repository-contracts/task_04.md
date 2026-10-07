---
task: task_04
spec: 0247-a-run-gate-that-runs-the-repository-contracts
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes `make verify-changed`, the
Verification of Runs and of the QA gate, run the Repository Contract Tests
that their Contract Relevance selects. Every behavior row runs against the
built tree in temporary clones with fake or stub commands where the Task tests
use them. No row reaches a provider, starts a real Agent Session, reads a
credential, or reads or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 and Success Metric 4 by executing task_01's
   five `./internal/verifyselect` tests.
3. MUST verify API Contract 2 by executing task_02's two
   `./internal/verifyselect` tests. It also confirms from the diff that the
   `verify`, `verify-docs`, `docs-test`, `repo-test`, `spec-budget`,
   `verify-changed-core` and `verify-changed-baseline` recipes and
   `REPO_CONTRACT_TESTS` are byte-identical to `59548f13`.
4. MUST run Surface Transcript 1, Surface Transcript 2 and Surface
   Transcript 3 through `go run -buildvcs=false ./cmd/verify-select` in a
   temporary clone of the audited head. Each runs on its own committed branch
   built as the transcript describes, against the clone's base. It records
   the three outputs and exit codes, and that covers Success Metric 3.
5. MUST replay the CI failure of 2026-10-07 for Success Metric 1. In a
   temporary clone of the audited head it removes `"internal/config",` from
   `guardedSpawningPackages` in `internal/suiteguardcontract/contract.go`,
   commits that on a branch together with a one-line comment in
   `internal/config/config.go`, and runs `make verify-changed` against the
   clone's base. That run MUST exit non-zero and name `internal/config`. It
   records the same replay at `59548f13` exiting 0, or cites the TechSpec's
   Current behavior section when that tree is unavailable.
6. MUST measure Success Metric 2 with the procedure of the TechSpec's
   Current behavior section. That is a temporary clone, a branch with one
   comment line in `internal/config/config.go`, `GOFLAGS=-count=1`, a warm
   `GOCACHE` and `make verify-changed` against the clone's base. It records
   the wall time beside the 247 s baseline. A difference above 20 s fails
   the row unless the report shows that the extra time is outside
   `verify-changed-contracts`, by timing that target alone in the same
   clone.
7. MUST run `make verify-docs` at the audited head and record that `repo-test`
   and `docs-test` still run every Repository Contract Test.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - the operator's queue log entries 237 and 238;
   - the three published Test Impact Analysis sources.

   For each source it reaches, it records what it read and whether the
   source still supports the design: an unknown or module-wide change runs
   everything, and declared non-source inputs select tests. For each source
   it cannot reach, it records the row as blocked with the reason.
9. MUST perform the glossary check of `docs/agents/domain.md`. It records
   that `CONTEXT.md` defines **Repository Contract Test** and
   **Contract Relevance**, and whether any other term the Spec introduced
   needs one.
10. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, `.roundfixrc.yml`, the CI workflows, the skills and the
      Baseline modules did not change;
    - that no contract test body changed.

    This row reads Task commits, so it declares a `commit_range` input.
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
12. MUST NOT accept a row whose only evidence is that a file was read. It
    MUST NOT send any request to a provider. The Linux-only behavior stays a
    documented limit, and no row tries to reproduce it on this host.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Exercise the selector and Makefile rows and the three Surface Transcripts.
- [ ] Replay the CI failure and measure the gate's time.
- [ ] Run `make verify-docs`, and record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The replayed CI failure now fails `make verify-changed`.
- [ ] The typical change grows the gate by no more than 20 s.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0247-a-run-gate-that-runs-the-repository-contracts/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; Surface Transcripts; Current behavior; Build Order 4
- ADR-0252; ADR-0182; ADR-0184
