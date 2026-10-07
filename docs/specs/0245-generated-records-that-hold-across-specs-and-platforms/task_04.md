---
task: task_04
spec: 0245-generated-records-that-hold-across-specs-and-platforms
status: failed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the Spec's authored terminal gate. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes a Baseline module's
version the output of a record step that checks it against the module's
content, makes the Coverage Record the same bytes on any host with the
platforms of each platform-limited test, and declares both records as derived
paths. Every behavior row runs in this checkout's tests, in temporary
directories or in a disposable clone. No row opens a network connection or
reads or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify Success Metric 1, API Contract 1 and API Contract 2 by running
   task_01's tests, then, in a disposable `git clone --no-local` of the
   audited head outside the repository, change one clause of a module without
   touching its version: record that the check fails naming the module and
   the record command, that the record step raises the version to one above
   the highest recorded and changes only that line of the module, and that a
   second record step changes nothing.
3. MUST verify Success Metric 2: on the audited checkout, the module check
   passes, every catalog module has an entry, every module's recorded latest
   version equals its file's version, no module changed a version number
   from the Spec's starting head, and the record step plus
   `make baseline-digests` change no file.
4. MUST verify Success Metric 3, the merge ablation, in a disposable clone:
   from the audited head create two branches that each change a different
   clause of `context-workflow`, each run the record step, `make
   baseline-digests` and the Managed Refresh, and commit. Merge the first
   into the clone's main, then merge main into the second with
   `merge.conflictStyle=merge`. Resolve as the derived merge of ADR-0192 and
   ADR-0233 does with the clone's `.roundfixrc.yml`: take the default
   branch's bytes for each conflicted declared path and each conflict hunk
   confined to declared lines, then run each matched regeneration in
   declaration order. Record the conflicted paths, that no source conflict
   remains, that the merged module's version is one above the first branch's,
   and that `make verify` exits 0 on the result.
5. MUST verify Success Metric 4 and API Contract 3: the Coverage Record's
   `platforms`, the platform lists of the macOS-only, Linux-only and
   Windows-only tests, that re-recording reproduces the same bytes, and that
   deleting one Linux-only test in a disposable clone makes
   `TestCoverageEquivalence` fail on this host naming that test and `linux`.
6. MUST verify Success Metric 5 by running task_02's fixture-module and
   anchor tests and recording their output.
7. MUST record that the Linux side of Goal 3 is proven by the Pull Request's
   Linux CI run of `TestCoverageEquivalence`. Until a Pull Request exists,
   record that row as `blocked (environment: no open Pull Request)` with the
   macOS evidence of Requirements 5 and 6.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the intervention log entries 215 and
   224-225, the Go command documentation, and the authoring measurement. For
   each source the gate reaches, record what it read and whether it still
   supports the design; for each it cannot reach, record the row as blocked
   with the reason.
9. MUST verify the guidance: the repository rule in
   `docs/agents/specific-repository.md`, the four derived declarations in
   `.roundfixrc.yml` and their config test, and ADR-0250 accepted.
10. MUST perform the glossary check of `docs/agents/domain.md`. Record that
    **Module Version Record** and **Coverage Record** are in `CONTEXT.md`,
    match the behavior and cite no Spec.
11. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`. Every Governed
    Path changed must be bounded in `_authorization.md`, and `Makefile`,
    `go.mod`, `_ownership.yml` records and the CI workflows must not change.
    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance, and support it with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the
      audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
13. MUST write the QA Report's `## Outcome` paragraph, which this Spec's own
    Archive Record will carry.
14. MUST NOT accept a row whose only evidence is that a file was read, and
    MUST remove every disposable clone it created.

## Subtasks

- [ ] Run the repository Verification and task tests.
- [ ] Exercise the record step and the check in a disposable clone.
- [ ] Run the two-branch merge ablation.
- [ ] Exercise the Coverage Record rows.
- [ ] Run the guidance, glossary and scope checks, and write the report.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The merge ablation ends on the next free version with `make verify` at
      exit 0.
- [ ] The outside-evidence row records what each source says or why it was
      unreachable.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0245-generated-records-that-hold-across-specs-and-platforms/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Data Models; Build Order 4
- ADR-0250; ADR-0192; ADR-0233; ADR-0080; ADR-0091; ADR-0104
