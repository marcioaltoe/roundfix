---
task: task_05
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes `TestRegenerationIsDeclared`
run when its inputs change, adds the Full Contract Run that CI runs on every
push to main and before every release, and names the selective gate's
exclusions. Every behavior row runs against the built tree in temporary
clones, with stub commands where the Task tests use them. No row reaches a
provider, starts a real Agent Session, reads a credential, triggers a GitHub
workflow, or reads or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 and API Contract 2 by executing task_01's three
   new tests and `TestRunPrintsContractInvocations`.
3. MUST verify Invariants 6 and 7 by executing task_02's two
   `./internal/verifyselect` repository tests.
4. MUST verify API Contract 3 and API Contract 4 by executing task_03's two
   tests. From the diff against `ae56aba0` it confirms that the `verify`,
   `verify-changed`, `verify-changed-contracts`, `verify-docs`, `docs-test`,
   `repo-test`, `spec-budget`, `verify-changed-core` and
   `verify-changed-baseline` recipes and `REPO_CONTRACT_TESTS` are
   byte-identical, and that each workflow diff only adds the one
   `Verify every contract` step and its comment.
5. MUST run Surface Transcript 1, Surface Transcript 2, Surface Transcript 3,
   Surface Transcript 4 and Surface Transcript 5 through
   `go run -buildvcs=false ./cmd/verify-select` in a temporary clone of the
   audited head. Transcripts 1 and 2 each run on their own committed branch
   built as the transcript describes, against the clone's base. It records
   each output and exit code; that covers Success Metric 2.
6. MUST replay Success Metric 1. In a temporary clone of the audited head it
   removes the `derived_paths` entry for
   `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/`
   from `.roundfixrc.yml`, commits that on a branch, and runs
   `make verify-changed-contracts VERIFY_BASE=<clone base>`. That run MUST
   exit non-zero with `TestRegenerationIsDeclared` failing. It records the
   same replay at `ae56aba0` exiting 0 without running that test, or cites
   the TechSpec's Current behavior section when that tree is unavailable.
7. MUST measure Success Metric 3: `make verify-contracts` at the audited head
   exits 0, and the number of tests it runs equals the number
   `verify-select -contracts -all` reports. It records the wall time.
8. MUST measure Success Metric 5: `make verify-docs` run twice at the audited
   head, recording both wall times and exit codes. The second, warm run MUST
   take no more than 60 s; the report records it beside the 45.5 s baseline
   in the TechSpec's Current behavior section.
9. MUST record, for the maintainer's 60 s target, the per-change cost of the
   selective gate: in a temporary clone, a one-line comment change in
   `internal/daemon/daemon.go` does not select `TestRegenerationIsDeclared`,
   and a one-line change to a Baseline module selects it. It lists the next
   contributors to the remaining suite time named in the TechSpec's Risks and
   Current behavior sections, without implementing any of them.
10. MUST record, as evidence this Spec did not author, each source of
    `_prd.md` → Acceptance evidence:
    - the operator's queue log entries 233 and 238;
    - the two published test impact analysis guides.

    For each source it reaches, it records what it read and whether the
    source still supports the design: selective tests on a pull request, and
    every test on the default branch. For each source it cannot reach, it
    records the row as blocked with the reason.
11. MUST perform the glossary check of `docs/agents/domain.md`. It records
    that `CONTEXT.md` defines **Full Contract Run** and carries the revised
    **Repository Contract Test** and **Contract Relevance**, and whether any
    other term the Spec introduced needs one.
12. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, `.roundfixrc.yml`, the skills and the Baseline modules
      did not change;
    - that `internal/config/regeneration_declared_test.go` changed only its
      directive line.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
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

    The push-to-main and release workflow runs happen after merge and are
    recorded as not yet observable, not as failures.
14. MUST NOT accept a row whose only evidence is that a file was read. It
    MUST NOT send any request to a provider or start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Exercise the selector, Makefile and workflow rows and the five Surface Transcripts.
- [ ] Replay the undeclared-path failure and measure the Full Contract Run and `make verify-docs`.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The replayed undeclared path now fails the selective gate.
- [ ] `make verify-contracts` runs every discovered contract and passes.
- [ ] `make verify-docs` stays within 60 s warm.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0248-the-regeneration-contract-runs-when-its-inputs-change/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Invariants 1-8; Surface Transcripts; Current behavior; Risks & Considerations; Build Order 5
- ADR-0253; ADR-0252; ADR-0182; ADR-0184
