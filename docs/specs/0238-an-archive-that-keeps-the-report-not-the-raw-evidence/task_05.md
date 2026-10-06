---
task: task_05
spec: 0238-an-archive-that-keeps-the-report-not-the-raw-evidence
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes a normal archive drop a
Spec's raw QA evidence and keep its QA Reports with an Evidence Manifest,
lets the Delivery Queue accept that cut as an exact move, adds a
request-only cut of one archived Spec, removes the two dependencies of the
repository's gates on archived evidence, and states the cut in the skills
and the Spec workflow guides. Every behavior row runs against the built tree
in temporary repositories and homes or in a disposable clone. No row deletes
anything under this repository's `docs/history/`, reaches the network on
purpose, or reads or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify Success Metric 1 and API Contract 1 by executing task_02's
   `internal/spec` tests, and by reproducing Surface Transcript 1 through
   the built `roundfix archive` on a disposable repository whose Spec has
   three evidence files, one linked from its QA Report. The row recomputes
   each manifest row's size and SHA-256 from `git show` at the recorded
   revision and resolves the rewritten link. It also archives a Spec
   without evidence and compares the result byte-for-byte with the binary
   built from the delivery base.
3. MUST verify Success Metric 2 and API Contract 4 by executing
   `TestArchiveCommitWithAListedEvidenceCutIsAnExactMove` and
   `TestArchiveCommitWithAnUnlistedEvidenceCutIsNotAnExactMove`, and confirm
   from the diff that `archiveCommitIsExact` keeps every other exact-move
   rule.
4. MUST verify Success Metric 3 and API Contracts 2 and 3 by executing
   task_03's tests and by reproducing Surface Transcript 2, Surface
   Transcript 3, Surface Transcript 4, Surface Transcript 5 and Surface
   Transcript 6 through the built `roundfix archive` on a disposable
   repository.
5. MUST verify Success Metric 4 in a disposable clone outside the
   repository: remove every `docs/history/specs/*/qa/evidence/` directory,
   commit the removal there, and run `make verify` and `make verify-docs`,
   recording both exit codes. It also runs task_01's tests there with the
   removal left uncommitted.
6. MUST verify Success Metric 5 by running the built
   `roundfix archive <slug> --drop-evidence` (dry run) for every Spec under
   this repository's `docs/history/specs/`, summing the files and bytes it
   reports, counting the kept override and superseded Specs, and confirming
   that `git status` shows no change afterwards. The row compares the sums
   with the PRD's figures (about 1,698 files and 13.5 MB to drop, 37 kept
   override Specs) and explains any difference.
7. MUST verify Invariants 1 to 11 from the tests and the diff, including the
   refusal of a symbolic link and of a path with a pipe, the unchanged
   override and superseded archives, and the coverage record with no package
   under `roundfix/docs/`.
8. MUST verify that the archive-spec, qa-gate and Roundfix skills, the
   Roundfix archive reference, the archive user guide and both generated
   Baseline guides state the cut, the manifest and the request-only command;
   that the `### QA settlement` section is byte-identical to the delivery
   base in all three skills; that the Vocabulary Contract's pattern is
   documented where it says; that each skill mirror equals its canonical
   file; and that every raised owned-skill version is recorded.
9. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: Spec 0214's measurement in
   `docs/references/archived-evidence-measurement.md`, the operator's Jev
   judgment of 2026-10-06, the adopter repositories mirrored in the
   Secondbrain (re-measured with `du`), GitHub's published retention page
   and the Git book's maintenance chapter. For each source it reaches, it
   MUST record what it read and whether it still supports the design. A
   source the Run sandbox cannot reach because network access was denied is
   recorded as `blocked (environment: network denied: <host>)` with an
   `outside-evidence row` item in its provenance; any other unreachable
   source is an ordinary environment-blocked row with its reason.
10. MUST verify that this Spec's own artifacts satisfy the promise rule, and
    perform the glossary check of `docs/agents/domain.md`. It records
    whether `CONTEXT.md`'s **Archive Command** entry, which says the command
    moves the whole Spec, needs the cut, and whether **Evidence Manifest**
    needs a glossary term.
11. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, that no derived
    Baseline file, skill version record or coverage record was edited
    outside its regeneration command, that nothing under `docs/history/`
    changed, and that `Makefile`, `go.mod` and the CI workflows did not
    change. This row reads Task commits, so it declares a `commit_range`
    input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance, and support it with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
13. MUST NOT accept a row whose only evidence is that a file was read.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Reproduce the six Surface Transcripts and recompute a manifest from Git.
- [ ] Run both gates in a disposable clone without archived evidence.
- [ ] Dry-run the cut over the whole History Root and check it changed nothing.
- [ ] Record the outside evidence, glossary check and scope audit, and write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] `make verify` and `make verify-docs` exit 0 in the clone without
      archived evidence.
- [ ] The dry run over the History Root changes no file and matches the
      PRD's figures or explains the difference.
- [ ] The scope audit finds no undeclared Governed Path change and no change
      under `docs/history/`.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0238-an-archive-that-keeps-the-report-not-the-raw-evidence/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-8; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Vocabulary Contract; Build Order 5
- ADR-0243; ADR-0215; ADR-0230; ADR-0154
