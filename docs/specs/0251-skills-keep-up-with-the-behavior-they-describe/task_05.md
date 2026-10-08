---
task: task_05
spec: 0251-skills-keep-up-with-the-behavior-they-describe
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds the Skill Coverage Map and
the Behavior Surface Record with their repository contract, makes the release
plan refuse a release with a Lagging Surface, and adds the skills-Task rule to
the Spec Consistency Check. Every behavior row runs against the built binary
in temporary repositories and temporary homes with fixture Specs. No row
reaches a provider, starts a real Agent Session, reads a credential, creates a
tag, triggers a GitHub workflow, or reads or writes the real `~/.roundfix`.

## Requirements

1. MUST build the binary and run the repository Verification, recording its
   result as a gate fact, and run `make verify-docs`, recording its exit code
   and wall time.
2. MUST verify API Contract 1 by executing task_01's four tests, and API
   Contract 2 and API Contract 3 by executing task_02's four tests and the
   `./internal/verifyselect` test that the record's declaration selects.
3. MUST replay Success Metric 4 in a temporary clone of the audited head: one
   edited sentence in one command's help, without re-recording, makes
   `TestTheSkillCoverageMapIsCurrent` fail naming that surface; running the
   record command then makes it pass, and a second run of the record command
   changes no byte.
4. MUST run Surface Transcript 1, Surface Transcript 2 and Surface Transcript
   3 through the built `roundfix release plan` in temporary repositories built
   as each transcript describes, and record each output and exit code. It also
   records the JSON of Surface Transcript 1 and confirms that the decision
   state and proposed version equal those of the same range without the map
   and record files. That covers API Contracts 4 to 6 and Success Metrics 1 to
   3.
5. MUST replay the Coverage Review in a temporary repository: a range whose
   only change to the surface's entry is a new `review` text exits 0 with the
   surface counted `reviewed`, and the next range, with that review unchanged
   and the surface changed again, exits 3.
6. MUST run Surface Transcript 4 through the built `roundfix spec check` on a
   fixture Spec committed after the map in a temporary repository, record the
   output and exit code, then add a skills Task and record that the finding
   clears; and record a malformed `## Skills` line reporting
   `SC-SKILLS-MALFORMED`. That covers API Contracts 7 and 8 and Success Metric
   5.
7. MUST run, at the audited head, `roundfix release plan` and record its
   `skill-coverage:` line, which is expected to be `introduced` because the
   base release predates the map.
8. MUST verify Success Metric 6: every pre-existing release plan test passes,
   and the active corpus golden differs from the starting tree only by the two
   new codes at 0.
9. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - the operator's queue log entries 193, 231 and 232;
   - the Drift extension's published documentation on reviewed findings;
   - the study by Tan, Wagner and Treude on outdated code element references.

   For each source it reaches, it records what it read and whether the source
   still supports the design: drift between documentation and behavior is
   common and recurs here, and a review holds until the reviewed behavior
   changes again. For each source it cannot reach, it records the row as
   blocked with the reason.
10. MUST perform the glossary check of `docs/agents/domain.md`. It records
    that `CONTEXT.md` defines **Behavior Surface**, **Skill Coverage Map**,
    **Behavior Surface Record**, **Coverage Review**, **Lagging Surface** and
    **Skills Declaration** and carries the revised **Release Plan**, and
    whether any other term the Spec introduced needs one.
11. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, the `Makefile`, the CI workflows and every Baseline
      module did not change;
    - that no text inside a `### QA settlement` section changed;
    - that every owned skill whose content changed raised its version and is
      recorded, and that its mirror is byte-identical.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification and
      `make verify-docs` at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
13. MUST NOT accept a row whose only evidence is that a file was read. It
    MUST NOT send any request to a provider or start a workflow.

## Subtasks

- [ ] Build the binary, run the repository Verification and `make verify-docs`.
- [ ] Exercise the package, contract, release plan and spec check rows and the four Surface Transcripts.
- [ ] Replay the stale record and the Coverage Review.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] A Lagging Surface blocks a release plan, and a skill edit or a Coverage Review clears it.
- [ ] A stale Behavior Surface Record fails the repository contract.
- [ ] A Spec without its skills Task fails the Spec Consistency Check.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0251-skills-keep-up-with-the-behavior-they-describe/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; Invariants 1-8; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Build Order 5
- ADR-0256; ADR-0189; ADR-0192; ADR-0253
