---
task: task_05
spec: 0243-a-history-that-holds-only-records
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the Spec's authored terminal gate. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the History Sanitize
Command: a plan that writes nothing, and a batch that converts the next
Legacy Archive Folders into Archive Records, reduces retired Findings and
Backlog Entries, and removes retired Review Artifacts and handoffs, only on
a clean tree covered by the annotated `history-full` tag. Every behavior row
runs the built binary in temporary repositories or in a disposable clone,
with a temporary home and a fake judge transport. No row reaches a provider,
reads a credential, or reads or writes the real `~/.roundfix`. No row
applies a batch to this repository, creates a tag in it, or changes its
`docs/history` or `.secondbrain-export`.

## Requirements

1. MUST run the repository Verification and record its result as a gate
   fact.
2. MUST verify Success Metric 1 and API Contract 1. Run the built binary's
   plan in a temporary repository with three Legacy Archive Folders and
   every kind, reproducing Surface Transcript 1. Record
   `git status --porcelain --untracked-files=all` before and after, which
   must match. Also run the plan once against this repository's own
   checkout and record its unit count, which must be 227 at the authoring
   head (223 folders and four kinds) or name the difference, and that the
   checkout's status is unchanged.
3. MUST verify Success Metric 2 and API Contract 3 in that temporary
   repository: create the annotated `history-full` tag, apply
   `--batch 2` with one `--promote`, reproducing Surface Transcript 2, then
   record each record's size and round trip and `git show` of three removed
   files at `source_revision` and at the tag.
4. MUST verify Success Metric 3: reproduce Surface Transcript 3 and record
   exit 2 and an unchanged tree for a missing tag, a lightweight tag, a tag
   that is not an ancestor of `HEAD`, a tag lacking a batch path and a dirty
   tree. Record each usage error of Invariant 14.
5. MUST verify API Contract 2 by running `--batch 1 --advise` without a key,
   reproducing Surface Transcript 4, and with a fake judge transport, and
   confirm neither writes under the repository.
6. MUST verify Success Metric 4, the ablation row. In a disposable
   `git clone --no-local` of the audited head outside the repository, create
   the annotated `history-full` tag, run the built binary's
   `history sanitize --apply --batch 300`, drop the history exclusions from
   the clone's `.secondbrain-export`, commit, and run `make verify` and
   `make verify-docs` there. Both must exit 0. Record the units applied, the
   bytes under the clone's `docs/history` before and after, each exit code,
   and a `no-qa` record for `0003-dogfood-polish`. Remove the clone
   afterwards.
7. MUST verify Success Metric 5 by running task_04's export tests and by
   recording, in the ablated clone, that the export contract fails before
   the exclusions are dropped and passes after.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the Jev judgment of 2026-10-06, the Git
   documentation for `git tag` and `git merge-base`, the Git book, and the
   authoring ablation. For each source the gate reaches, record what it read
   and whether it still supports the design. For each it cannot reach,
   record the row as blocked with the reason.
9. MUST verify the guidance: the Roundfix skill section and index row, each
   mirror equal to its canonical file and the raised version recorded; the
   Baseline sentence in the module and the generated guide; `make
   baseline-digests` leaving no diff; the plan and confirmation phrases in
   the command guide; and the two repointed `CHANGELOG.md` paths resolving.
10. MUST perform the glossary check of `docs/agents/domain.md`. Record that
    **History Sanitize Command**, **Legacy Archive Folder**, **Sanitize
    Batch**, **Reduced History Entry** and **History Full Tag** are in
    `CONTEXT.md` and match the behavior, and that **Archive Record** and
    **History Root** were revised.
11. MUST verify that this Spec's own artifacts satisfy the promise rule, and
    that `git diff --name-status` of the audited range shows no deleted,
    renamed or modified path under `docs/history` and no change to
    `.secondbrain-export`, except this Spec's own archive at the end.
12. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`. Every Governed
    Path changed must be bounded in `_authorization.md`, and `Makefile`,
    `go.mod` and the CI workflows must not change. This row reads Task
    commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
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
14. MUST write the QA Report's `## Outcome` paragraph, which this Spec's own
    Archive Record will carry.
15. MUST NOT accept a row whose only evidence is that a file was read. MUST
    NOT send any request to a provider, and MUST NOT apply a batch to, tag,
    or change the `docs/history` or `.secondbrain-export` of this
    repository.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Exercise the plan, batch, refusal and advice rows in temporary repositories.
- [ ] Run the whole sanitize in a disposable clone and both gates there.
- [ ] Run the guidance, glossary and scope checks.
- [ ] Write the QA report with the Pull Request row blocked and an Outcome.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] Both gates exit 0 in the fully sanitized clone.
- [ ] The outside-evidence row records what each source says or why it was
      unreachable.
- [ ] The scope audit finds no undeclared Governed Path change and no
      history change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0243-a-history-that-holds-only-records/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → Measured inventory; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Exact texts; Operator batch procedure; Build Order 5
- ADR-0248; ADR-0247; ADR-0215; ADR-0080; ADR-0091; ADR-0104
