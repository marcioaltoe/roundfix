---
task: task_05
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the owner's
`roundfix archive` promote cited files and relink Live Documents, makes the
Delivery Queue commit and prove that change, and makes
`roundfix baseline update` list its History Citations, promote on request and
name Stale History Citations. Every behavior row runs against the built tree in
temporary repositories with synthetic Specs and history and a temporary home.
No row reaches a provider, starts a real Agent Session, reads a credential,
pushes, reads or writes the real `~/.roundfix`, or touches an adopter
repository.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST verify API Contracts 1 to 5 and Invariants 1 to 5 and 10 by executing
   task_02's ten tests and every existing archive and Archive Record test.
3. MUST verify API Contract 6 and Invariant 6 by executing task_03's four
   tests and the existing delivery archive tests.
4. MUST verify API Contracts 7 to 10 and Invariants 7 to 9 by executing
   task_04's eight tests and every existing `baseline update` and History
   Sanitize Command test.
5. MUST run Surface Transcript 1 with the rebuilt `bin/roundfix` from the root
   of a temporary repository built as the TechSpec's Surface Transcripts
   section describes, and record its output and exit code, the promoted copy's
   blob next to the removed one, a diff of each relinked file showing only the
   destination or path changed, and the record's `promoted` and `relinked`
   lists. That covers Success Metric 1. It also records the two refusals of
   Success Metric 2 with the tree byte-identical afterwards, and one archive
   without inbound citations for Success Metric 3.
6. MUST record Success Metric 4 through task_03's tests, because a real
   Delivery Queue needs a remote and a Pull Request; it records the commit's
   changed paths they assert and the tampered-file refusal.
7. MUST run Surface Transcript 2, Surface Transcript 3 and Surface Transcript 4
   through `go run -buildvcs=false ./cmd/roundfix` in that order against one
   temporary adopted repository built as the TechSpec's Surface Transcripts
   section describes, under a temporary home. Around Surface Transcript 3 it
   records the citing files' bytes before and after and the promoted copy, and
   it records one moved citing line changing the preview's digest. That covers
   Success Metric 5 and Success Metric 6, with a repository without Pending
   History compared against `roundfix baseline plan` or task_04's digest test.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - the characterization committed as
     `docs/references/2026-10-12-live-citations-characterization.md`;
   - the adopter reports it names;
   - the MD Links passage on rewriting only the URL token and warning on
     delete.

   For each source it reaches, it records what it read and whether the source
   still supports the design: the classes of broken citations, the share from
   Live Documents, the secondary-Spec links, and destination-only rewrites.
   For each source it cannot reach, it records the row as blocked with the
   reason, using `blocked (environment: network denied: <host>)` when the Run
   sandbox denied the host.
9. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` defines **Live Document**, **Archive Relink**,
   **History Citation** and **Stale History Citation** and carries the revised
   **Archive Record**, **Archive Command** and **Managed Refresh**, and whether
   any other term the Spec introduced needs one.
10. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, `.roundfixrc.yml`, the Makefile, the CI workflows, the
      Baseline modules and the Baseline guides did not change;
    - that ADR-0083, ADR-0247, ADR-0254 and every Archive Record are
      byte-identical;
    - that the `### QA settlement` section of every skill is byte-identical;
    - that the old wording task_01 replaced survives in no live document.

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
12. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT send any request to a provider, push, or write outside temporary
    directories and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests and the four Surface Transcripts.
- [ ] Record the refusals, the queue proof and the digest binding.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The owner's archive leaves every Live Document citation resolving and the queue proves it.
- [ ] The update lists History Citations, promotes only on request and names Stale History Citations.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0255-an-archive-that-relinks-and-an-update-that-lists-citations/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-7; User Stories 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; API Contract 10; Invariants 1-10; Surface Transcripts; Current behavior; Risks & Considerations; Build Order 5
- ADR-0260; ADR-0083; ADR-0247; ADR-0230; ADR-0173; ADR-0254; ADR-0248; ADR-0184
