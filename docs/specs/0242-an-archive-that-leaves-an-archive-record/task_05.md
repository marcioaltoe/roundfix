---
task: task_05
spec: 0242-an-archive-that-leaves-an-archive-record
status: failed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the Spec's authored terminal gate. It declares what the matrix
covers and settles the Spec on evidence. The Spec does three things. An
archive leaves one Archive Record and removes the Spec folder, whose bytes
stay in Git at the recorded `source_revision`. Every reader of archived
Specs works from the record or from Git. Jev advises on what to keep and
never gates. Every behavior row runs against the built tree with temporary
repositories, a temporary home and a fake judge transport. No row reaches a
provider, reads a credential, or reads or writes the real `~/.roundfix`. No
row deletes, moves or rewrites a file under this repository's
`docs/history`.

## Requirements

1. MUST run the repository Verification and record its result as a gate
   fact.
2. MUST verify Success Metric 1 and API Contracts 1 and 2 by executing
   task_01's tests and one archive of a disposable Spec in a temporary
   repository with the built binary. Record the confirmation line, the
   record's bytes and size, and `git show <source_revision>:<path>` of three
   removed files, reproducing Surface Transcript 1. Do the same for an
   override archive with a 300-byte reason and a promotion, reproducing
   Surface Transcript 2, and for a superseded Spec. Reproduce Surface
   Transcript 5 with an uncommitted change under the Spec folder.
3. MUST verify Success Metric 2 and API Contract 5 by executing task_01's
   exact-retirement tests and task_02's reader tests. Record a squash-merged
   record-only Spec that reconcile releases with merge evidence, and one
   whose `source_revision` is absent and whose worktree is kept.
4. MUST verify Success Metric 3 and API Contracts 3 and 4 by running
   `--plan` against a fake judge, reproducing Surface Transcript 3, and
   without a key, reproducing Surface Transcript 4, and by running one
   `--promote`. Confirm that `git status --porcelain` is empty after both
   plan runs.
5. MUST verify Success Metric 4, the ablation row. In a disposable
   `git clone --no-local` of the audited head outside the repository, remove
   every Spec folder under `docs/history/specs` in a clone commit, then run
   `make verify` and `make verify-docs` there. Both must exit 0. Record each
   exit code and the clone's removed file count, and remove the clone
   afterwards.
6. MUST verify Success Metric 5. Re-run the static reader scan from
   `_techspec.md` → Readers (`git grep -l -I -e docs/history -- ':!docs/history/**'`
   and the archive-root builder scan) at the audited head. Confirm that
   every reader it lists is either changed with a named test or recorded as
   unaffected with its reason. Any new reader is a finding.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence. These are the Jev judgment of
   2026-10-06, the Spec 0214 measurement, the authoring ablation, and the
   GitHub and Git documentation. For each source the gate reaches, record
   what it read and whether it still supports the design. For each it
   cannot reach, record the row as blocked with the reason.
8. MUST verify the guidance. Check these points:
   - the three skill sections exist, and `### QA settlement` is
     byte-identical across qa-gate, archive-spec and the Roundfix skill,
     carries the Archive Record wording and no longer says a passing archive
     keeps "The Spec and its QA report and evidence";
   - each mirror equals its canonical file, and the raised versions are
     recorded;
   - both Baseline sentences appear in the modules and the generated guide;
   - the Vocabulary Contract's patterns are documented where it says;
   - `make baseline-digests` leaves no diff.
9. MUST perform the glossary check of `docs/agents/domain.md`. Record that
   **Archive Record** and **Archive Advice** are in `CONTEXT.md` and match
   the behavior, that **Archive Command** no longer says it moves the whole
   Spec, and whether **History Root** or **QA Archive Override** needs a
   revision.
10. MUST verify that this Spec's own artifacts satisfy the promise rule, and
    that `git diff --name-status` of the audited range shows no deleted,
    renamed or modified path under `docs/history` except this Spec's own
    archive at the end.
11. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`. Every Governed
    Path changed must be bounded in `_authorization.md`, and `Makefile`,
    `go.mod` and the CI workflows must not change. This row reads Task
    commits, so it declares a `commit_range` input.
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
14. MUST NOT accept a row whose only evidence is that a file was read. MUST
    NOT send any request to a provider, and MUST NOT run `roundfix history`
    or migrate any archived Spec.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Exercise the archive, retirement, reader and plan rows.
- [ ] Run the ablation in a disposable clone.
- [ ] Re-run the reader scan, the guidance checks and the glossary check.
- [ ] Write the QA report with the Pull Request row blocked and an Outcome.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] Both gates exit 0 in the ablated clone.
- [ ] The outside-evidence row records what each source says or why it was
      unreachable.
- [ ] The scope audit finds no undeclared Governed Path change and no
      history change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0242-an-archive-that-leaves-an-archive-record/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → Readers; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Exact texts; Vocabulary Contract; Build Order 5
- ADR-0247; ADR-0215; ADR-0080; ADR-0091; ADR-0104
