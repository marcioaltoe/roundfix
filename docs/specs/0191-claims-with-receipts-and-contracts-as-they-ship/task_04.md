---
task: task_04
spec: 0191-claims-with-receipts-and-contracts-as-they-ship
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes the Spec Consistency Check
in `internal/speccheck`, four authoring skills, the Roundfix skill, the
glossary and one user guide. Every behavior row runs the built `roundfix`
binary, or executes the named tests against the built tree. Each binary
command runs in a disposable clone of the candidate with a disposable Roundfix
Home. No command reaches GitHub, a provider or a model.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST reproduce Surface Transcript 1 through the built binary: in a
   disposable clone, a Spec committed after the guide whose PRD quotes a
   decision record with one changed word. The finding lines and the exit code
   MUST match the transcript.
3. MUST reproduce Surface Transcript 2 through the built binary: the same
   held Spec with an attribution and no receipt, without `--strict`. It MUST
   also show that `--strict` exits `1` for the same artifact, and that adding
   the verbatim receipt clears the finding.
4. MUST reproduce Surface Transcript 3 through the built binary: a held Spec
   whose first transcript has no exit line and whose second transcript the QA
   Task names in no Requirement.
5. MUST reproduce Surface Transcript 4 through the built binary: the same
   Spec in a disposable clone from which the guide was removed. It MUST
   report no declaration gap and list both skips with the guide path.
6. MUST verify, by executing the horizon tests task_01 names against the
   built tree, that a PRD committed before the guide is not held, that a PRD
   committed in or after the guide's commit and an uncommitted PRD are held,
   and that an unreadable history holds every Spec.
7. MUST verify the non-regression promise:
   - by executing task_01's characterization test against the built tree,
     the claims recorded for the four fixture artifacts are unchanged;
   - the corpus golden holds the five new codes at `0`, and every other count
     equals the count in the golden at the Delivery Base, read with
     `git show`;
   - `SC-CITATION-UNSUPPORTED` still reports the false attribution of the
     existing citation characterization.
8. MUST verify that this Spec obeys its own rules. In a disposable clone,
   copy this Spec's PRD and TechSpec as a new Spec committed after the guide,
   and run the built binary's `spec check --stage techspec --strict` on it.
   It MUST report no `SC-RECEIPT-UNPROVEN`, `SC-RECEIPT-MISSING`,
   `SC-TRANSCRIPT-UNDECLARED` or `SC-TRANSCRIPT-MALFORMED` finding.
9. MUST record, as evidence this Spec did not author:
   - the "Reduce hallucinations" page of the Claude documentation
     (<https://docs.anthropic.com/en/docs/test-and-evaluate/strengthen-guardrails/reduce-hallucinations>)
     and the Cram documentation (<https://pypi.org/project/cram/>), each with
     the passage that supports the rule;
   - a replay on the archived PRD and TechSpec of Spec 0182, copied into a
     disposable clone as a Spec committed after the guide. The built binary
     MUST report one `SC-RECEIPT-MISSING` per resolved attribution and no
     `SC-RECEIPT-UNPROVEN`, and the row MUST record the count.

   When a page or the archived Spec is unavailable, it MUST record that part
   of the row as blocked with its reason.
10. MUST verify that the guide ships in the embedded bundle at the path the
    horizon reads, that the four authoring skills and their mirrors carry the
    sections task_03 names, that `make skills-sync-check` exits `0`, and that
    `TestSettlementGuidanceIsOneTable` passes.
11. MUST verify that this Spec's own artifacts satisfy the promise rule, and
    check whether it introduced, changed or retired a glossary term that
    `CONTEXT.md` does not carry.
12. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, and that every
    Governed Path is bounded in `_authorization.md`. A path rewritten by a
    sanctioned regeneration command that `_authorization.md` names also
    counts as declared, and the scope row MUST name each one it relies on.
13. MUST record the non-waivable Pull Request row as
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
14. MUST NOT accept a row whose only evidence is that a file was read.
15. MUST NOT write to the live Run Database under `~/.roundfix` or to any path
    under `~/dev/secondbrain`, and MUST NOT call a model or a paid API.

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

- `newest="$(find 'docs/specs/0191-claims-with-receipts-and-contracts-as-they-ship/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-6; Success Metrics 1-4; Acceptance
evidence; `_techspec.md` → Testing Approach 1-7; API Contracts 1-3; Surface
Transcripts 1-4; ADR-0080; ADR-0088; ADR-0091; ADR-0093; ADR-0094; ADR-0096;
ADR-0104; ADR-0116; ADR-0117; ADR-0130; ADR-0155; ADR-0156; ADR-0166;
ADR-0167; ADR-0168; ADR-0176; ADR-0183; ADR-0184.

## Result
