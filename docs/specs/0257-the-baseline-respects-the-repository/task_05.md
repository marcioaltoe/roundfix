---
task: task_05
spec: 0257-the-baseline-respects-the-repository
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec records the Better Auth base path
as a repository decision, lets the managed refresh record declared decisions,
keeps a skill that only differs from its Setup Snapshot instead of calling it
trailing, and states the vendored-skill divergence order. Every behavior row
runs against the built tree, the embedded catalog, temporary repositories and
fake restore dependencies. No row reaches a provider, the network or an
adopter repository, starts a real Agent Session, reads a credential, or reads
or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST execute every test that task_01 to task_04 created or renamed and
   every declared break of `_techspec.md` → Testing Approach, and record each
   result.
3. MUST measure Success Metric 1 through the built binary: adopt a temporary
   Standard TypeScript repository through `roundfix baseline plan` and
   `roundfix baseline apply`, run `roundfix baseline update --decision auth.base_path=/auth`
   without confirmation and record that nothing was written, then with
   `--yes`, and record the backend guide's Better Auth lines, the Setup
   Manifest's `auth.base_path`, `auth.provider` and `http.contract` values,
   the count of Better Auth exceptions, and a second update's state.
4. MUST measure Success Metric 2 on a second temporary Standard TypeScript
   adopter that records no base path: the backend guide's bytes before and
   after an update, and the absence of `auth.base_path` in its Setup Manifest.
   It MUST also record that the Standard TypeScript golden `backend.md` did
   not change in the Spec's commits.
5. MUST measure Success Metric 3 through the built binary on a temporary
   repository whose required external skill matches its lock but not its
   snapshot: Doctor's `skills` line and exit code, the update preview's state
   and `skills.differing`, and the skill tree's digest before and after
   `--yes` with `--skills-source-dir` pointing at a temporary object store. It
   MUST record the repository search for `DR-SKILL-TRAILS-SNAPSHOT`,
   "trails its Setup Snapshot", "trails the Setup Snapshot" and
   `TrailingSetupSkills` under `internal`, `cmd`, `docs/user-guide`,
   `.agents/skills` and `skills`, and record any hit.
6. MUST measure Success Metric 4: the three choices, in order, in
   `docs/agents/skill-dispatch.md` and in the `roundfix` skill, and for each
   row of `_techspec.md` → Divergence audit the clause text that overrides the
   skill, read from the embedded catalog.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the Better Auth options reference, the
   Microsoft APM drift-detection page, ingraft's editable-vendors page and
   the Safeguard Research article. For each source it reaches, it records what
   it read and whether the source still supports the design; the quoted
   passages in the PRD need no network. For a source it cannot reach, it
   records the row as `blocked (environment: network denied: <host>)`.
8. MUST perform the glossary check of `docs/agents/domain.md`: record that
   `CONTEXT.md` carries **Better Auth Base Path**, **Setup Snapshot
   Difference** and **Vendored Skill**, and whether any other term the Spec
   introduced needs one.
9. MUST verify from the repository history:
   - that each Task's changed files stay within its declarations or its
     `## Recorded paths`;
   - that every Governed Path changed is bounded in `_authorization.md`;
   - that `go.mod`, `.roundfixrc.yml`, the Makefile, the CI workflows, the
     Source Baseline assets, `internal/baseline/assets/retention/`, the Setup
     Snapshots and every module other than `core` and `rust` did not change;
   - that the `### QA settlement` section of every skill is byte-identical;
   - that `make verify-docs` passes at the audited head.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
11. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT send any request to a provider or write outside temporary directories
    and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests and the declared breaks.
- [ ] Measure Success Metrics 1 to 4 through the built binary.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] A declared base path moves the guide and the exception, and an
      unrecorded one changes nothing.
- [ ] A differing skill is kept and named without a warning.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0257-the-baseline-respects-the-repository/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- `_prd.md` → Goals; Core Features 1-5; User Stories 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence
- `_techspec.md` → Current behavior; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Exact clause texts; Divergence audit; Derived files; Invariants 1-10; Testing Approach; Build Order 5
- ADR-0263; ADR-0264; ADR-0221; ADR-0224; ADR-0076
