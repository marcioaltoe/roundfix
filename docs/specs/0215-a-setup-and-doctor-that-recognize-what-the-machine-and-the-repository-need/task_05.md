---
task: task_05
spec: 0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds the `gh`, `git`, `remote`,
`toolchain` and `environment` readiness lines to Doctor and Setup, reports
upstream skills that trail their Setup Snapshot, makes `deliver start` refuse
when this machine cannot publish, adds `verification.tools`, and has this
repository declare its tools and derived paths. Every behavior row is
exercised through the built binary or by executing the named tests against
the built tree, with disposable repositories, a disposable Roundfix Home and
fake `gh` and `git` executables. No command reaches GitHub or the network,
except the reads of published pages the outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST reproduce Surface Transcripts 1 and 2 through the built binary, with
   `HOME` set to a disposable directory, a disposable Git repository, and
   fake `gh`, `git`-wrapping and `acpx` executables first on `PATH` that
   answer as an unauthenticated and as an unreachable forge; record standard
   output, standard error and exit code for each, and that no fake received
   `--show-token`.
3. MUST verify, by executing task_01's and task_02's tests against the built
   tree, every code of the TechSpec's "Finding codes" except
   `DR-SKILL-TRAILS-SNAPSHOT`, Invariants 1 to 3, the bounded probe and the
   key sentinel.
4. MUST reproduce Surface Transcript 4 through the built binary in a
   disposable repository holding one upstream skill that matches its lock and
   differs from its snapshot tree, and run `roundfix baseline update` there
   without confirming, recording that its preview lists the skill; and MUST
   execute task_03's tests against the built tree.
5. MUST reproduce Surface Transcript 3 through the built binary in a
   disposable repository and Roundfix Home whose Spec carries a delivery
   grant, with the fake `gh` unauthenticated, and record that no Delivery
   Queue was created; and MUST verify through task_04's tests that a `warn`
   proceeds and that Setup prints the five lines (API Contract 2).
6. MUST verify that this repository's `.roundfixrc.yml` loads with the four
   tools and three derived-path declarations task_02 names, and record the
   `toolchain` line the built binary prints in this repository with a fake
   forge.
7. MUST verify that the `doctor`, `setup`, `deliver`, `baseline` and
   configuration guides and the Roundfix Skill's `setup`, `deliver` and
   `baseline` references describe the behavior the Tasks name, that each
   mirror equals its canonical file, that the version was raised and
   recorded, and that the `### QA settlement` section of every skill is
   unchanged.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the maintainer's request and the adopter
   runbook in the Secondbrain inbox; the two measurements of 2026-10-02,
   re-measured read-only on the built tree; the GitHub CLI release notes for
   2.81.0 and `pkg/cmd/auth/status/status.go` in `cli/cli`; Git's 2.23.0
   release notes; and the intervention log entry 66. For each source it
   reaches it MUST record what it read and whether it still supports the
   design; for each it cannot reach it MUST record the row as blocked with the
   reason.
9. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether the
   readiness line, the finding code and the trailing skill need glossary
   terms beside **Doctor Command** and **Repository Skill Set**.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, and that
    `Makefile`, `go.mod`, `internal/cli/cli_test.go` and the CI workflows did
    not change. This row reads Task commits, so it declares a `commit_range`
    input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
12. MUST NOT accept a row whose only evidence is that a file was read.
13. MUST NOT log in or out of `gh`, change Git configuration, reach GitHub
    from a test or row other than the published pages of requirement 8, or
    write to the live Run Database under `~/.roundfix`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0215-a-setup-and-doctor-that-recognize-what-the-machine-and-the-repository-need/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; User Stories 1-7; Core Features 1-9; Success Metrics 1-6;
Acceptance evidence; `_techspec.md` → Testing Approach 1-6; API Contract 1;
API Contract 2; API Contract 3; API Contract 4; API Contract 5; Surface
Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript
4; Integration Points; ADR-0080; ADR-0088; ADR-0091; ADR-0104; ADR-0155;
ADR-0156; ADR-0167; ADR-0220; ADR-0221.
