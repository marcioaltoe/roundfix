---
task: task_05
spec: 0258-authoring-that-proves-what-it-claims
status: failed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds four structural proof gaps to
the Spec Consistency Check for every active Spec, states the four proof rules
in `write-techspec` and `write-tasks`, and documents the codes for adopters.
Every behavior row runs against the built tree and temporary repositories. No
row reaches a provider, starts a real Agent Session, reads a credential, reads
or writes the real `~/.roundfix`, or reads another repository.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST execute every test that task_01, task_02 and task_03 created, the
   declared breaks of `_techspec.md` → Testing Approach, and task_04's phrase
   checks, and record each result.
3. MUST measure Success Metric 1 through the built `bin/roundfix`. In a
   temporary repository, it checks one fixture Spec per code twice, once with
   the section or line missing and once with it present, with and without
   `--strict`, and records each exit code and finding line. It covers API
   Contracts 1 to 4, a deviation that names no PRD unit, an Absence Evidence
   item without its result, a suite no Verification runs and a checkpoint
   outside the path to the gate.
4. MUST measure Success Metric 2. It records each rule of `_techspec.md` →
   Skill and document texts in both skills and both templates, the two raised
   versions in the owned skill version record, and that every mirror matches.
5. MUST measure Success Metric 3. It records the fix line of each finding from
   requirement 3 and that it names the section or line and the skill rule, and
   the four codes and the adopter paragraph in the command reference and the
   user guide, against `_prd.md` → Adopter notes.
6. MUST measure Success Metric 4. It records `bin/roundfix spec check --strict`
   over the active Spec corpus at the audited head with zero findings for the
   four codes, this Spec included, the results of `TestCheckCorpusGolden` and
   `TestCheckActiveCorpusHasNoErrors`, and that no file under `docs/history/`
   changed in the Spec's range.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - the adopted Backlog Entry in this Spec's references, comparing each of its
     four gaps with the rule and code that now answers it;
   - Robinson 2006, Juergens 2011, Zimmermann 2022 and Gotel and Finkelstein
     1994, through the passages `_prd.md` quotes.

   For each source it reaches, it records what it read and whether the source
   still supports the design. For a source it cannot reach, it records the row
   as blocked with the reason, using
   `blocked (environment: network denied: <host>)` when the Run sandbox denied
   the host. The adopted Backlog Entry needs no network.
8. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` defines **Absence Evidence**, **ADR Alternatives**,
   **Coverage Cross-Read**, **Consumer Declaration**, **Shared Contract** and
   **Integration Checkpoint**, and whether any other term the Spec introduced
   needs one.
9. MUST repeat the wording sweep. For the replaced sentence of `write-tasks`
   on whole-package suites, it searches `docs/`, `skills/`, `.agents/skills/`,
   `CONTEXT.md` and `AGENTS.md`, outside `docs/history/` and this Spec's
   folder, and records every remaining copy with whether it is allowed.
10. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, `.roundfixrc.yml`, the Makefile, the CI workflows, every
      Baseline module, every guide under `docs/agents/` and every production
      Go file outside `internal/speccheck` did not change;
    - that the `### QA settlement` section of every skill is byte-identical.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the Unreachable Acceptance criterion of `_prd.md` as not
    measured by this gate, with its `satisfied-by` follow-up.
12. MUST record the non-waivable Pull Request row as
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
13. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT send any request to a provider or write outside temporary directories
    and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests, the declared breaks and the phrase checks.
- [ ] Measure Success Metrics 1 to 4 and repeat the wording sweep.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] Each of the four codes fails a fixture missing its section or line and
      passes the same fixture with it present.
- [ ] The active corpus reports zero findings for the four codes.
- [ ] The scope audit finds no undeclared Governed Path change.

## Consumers

None. The gate changes no Shared Contract; it writes only its own QA report.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0258-authoring-that-proves-what-it-claims/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- [_prd.md](_prd.md) — Goals; Core Features 1-6; User Stories 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Adopter notes
- [_techspec.md](_techspec.md) — Current behavior; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Skill and document texts; Invariants 1-7; Testing Approach; Risks & Considerations; Build Order 5
- ADR-0265; ADR-0093; ADR-0094; ADR-0156; ADR-0168; ADR-0183; ADR-0176; ADR-0148; ADR-0252; ADR-0257; ADR-0258; ADR-0244; ADR-0256; ADR-0184; ADR-0253; ADR-0259
