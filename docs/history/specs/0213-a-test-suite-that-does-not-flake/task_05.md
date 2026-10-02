---
task: task_05
spec: 0213-a-test-suite-that-does-not-flake
status: failed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes tests and test helpers in
`internal/store`, `internal/daemon`, `internal/cli`, `internal/testfixture`
and `internal/baseline`, and no production file. Every behavior row runs the
named tests against the built tree, and the flake rows repeat them under load.
No command reaches GitHub, a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST repeat the authoring stress shapes against the built tree, compile each
   package's test binary once, and record the iterations and failures of each
   row beside the PRD's before-rates:
   - the linger subtest: 8 parallel workers × 10 rounds ×
     `-count=50 -cpu 1,4` (before: 74 of 8,000 failed);
   - both Run Budget tests: 8 workers × 6 rounds × `-count=20 -cpu 1,4`
     beside ten self-terminating CPU burners (before: 9 and 5 of 1,920
     failed).

   A row passes only with zero failures. The suite guard fingerprints the
   checkout, so the stress MUST run on a checkout nothing else edits during
   the run.
3. MUST verify, by executing `TestScriptFixtureIsTheCompiledTestBinary` and
   `TestNoTestWritesAnExecutableOutsideTheResidue`, that no non-governed
   `internal/cli` test executes a file it wrote. The guard MUST pass on the
   tree, fail on a seeded written executable naming `file:line`, and hold the
   residual inventory of the TechSpec. The gate MUST also confirm, from the
   Git history, that `internal/cli/cli_test.go` is unchanged by this Spec.
4. MUST verify, by executing `TestImplementDetachChildEndsWhenItsTestBinaryDies`
   and `TestDetachSurvivorEndsWhenItsTestBinaryDies`, that a killed test
   binary leaves no fixture process. MUST repeat the authoring probe outside
   Go: compile the `internal/cli` test binary, run
   `TestRunImplementDetachSurvivesCallerProcessGroupKill` with a private
   `TMPDIR`, send `SIGKILL` to that binary once its `prompt-started` file
   exists, and confirm that no process whose command names that binary remains
   within 60 seconds (before: 3 of 3 runs leaked). Any process this probe
   leaves MUST be killed by its recorded PID.
5. MUST verify, by executing `TestAssetsSyncTemplateLeavesNoDirectoryBehind`
   and the Assets Sync tests under `-count=3 -cpu 1,4`, that the template
   leaves no directory behind and every copy is a valid repository.
6. MUST record, as evidence this Spec did not author:
   - the CI log of run 36903640400 attempt 1, read with
     `gh run view 36903640400 --attempt 1 --log-failed`, which shows the
     `adapter_floor_test.go:50` failure;
   - Go issue golang/go#22315 (<https://github.com/golang/go/issues/22315>),
     confirming the forked-child descriptor race behind `ETXTBSY`;
   - the CI Verification gate of this Spec's Pull Request as the Linux run of
     the converted fixtures, recorded as blocked until that Pull Request
     exists.

   When a source is unavailable, it MUST record that row as blocked with its
   reason.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   check whether it introduced, changed or retired a glossary term. Script
   fixture and residual inventory are test-suite terms, not domain terms;
   record whether `CONTEXT.md` needs either.
8. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, and that no Governed
   Path and no production file changed, as `_authorization.md` records
   `paths: []`. This row reads Task commits, so it declares a `commit_range`
   input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It MUST support that row with the pre-PR equivalent of
   each control:
   - approval: the maintainer's recorded delivery authority;
   - checks and status: the Daemon's repository Verification at the audited
     head;
   - unresolved threads: none, because no Pull Request exists;
   - Merge-Ready acceptance: none yet, because the pre-PR review follows
     archive;
   - review-artifact ancestry: the audited head named as the claimed
     candidate.
10. MUST NOT accept a row whose only evidence is that a file was read.
11. MUST NOT write to the live Run Database under `~/.roundfix`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree, including the stress rows, and
      record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.
- [ ] Each flake row records its after-rate beside its before-rate.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0213-a-test-suite-that-does-not-flake/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-5; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-4; The residual inventory;
Integration Points; Build Order 5; ADR-0080; ADR-0088; ADR-0091; ADR-0104;
ADR-0125; ADR-0126; ADR-0155; ADR-0156; ADR-0167; ADR-0213.
