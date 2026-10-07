---
task: task_05
spec: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds the Project Config Format
Command `verification.format`. The QA step runs it over the Spec's `qa/`
directory after an import and before the QA Report commit, restoring the
original bytes when it fails. `roundfix reopen` also returns a completed gate
to `pending` over a Late Dependency proven from Git. Every behavior row runs
against the built tree with temporary repositories, a temporary home, fake
runners and a shell-script formatter. No row reaches a provider, starts a real
Agent Session, reads a credential, or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 by executing task_02's four config tests and by
   reading the default configuration rendered by the built binary's
   `roundfix init` in a temporary home and repository.
3. MUST verify Success Metric 1, Success Metric 2, Success Metric 3 and API
   Contract 2 by executing task_03's six Daemon tests and its CLI test. It MUST
   also confirm from the diff that the import (`copyPriorQAPass`) and the carry
   proof are unchanged, that no path outside the Spec's `qa/` directory reaches
   the formatter, and that paths are passed as arguments, never interpolated.
4. MUST re-measure the adopter's failure shape in a temporary repository. Set
   a repository Verification that fails on any file under `qa/` a shell-script
   formatter would rewrite, record an unformatted failed pass on a Run Branch,
   and run the built binary's gate with `verification.format` set and then
   empty, recording both outcomes. When a fake runner cannot drive the built
   binary to the QA step, MUST record the row on task_03's Daemon test evidence
   and say why.
5. MUST verify Success Metric 4 and API Contract 3 by executing task_04's spec
   and CLI tests, and by running the built `roundfix reopen` in a temporary
   repository with a Late Dependency, recording stdout, exit code and the QA
   Task's `## Invalidation` record.
6. MUST verify Success Metric 5: every existing config, Daemon prior-pass and
   carry, reopen, implement and spec test passes, and no existing test
   expectation changed.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the adopter's triaged report in the
   Secondbrain, the adopter's mirrored repository records named there, read
   only, and the lint-staged README. For each source it reaches it MUST record
   what it read and whether it still supports the design. For each source it
   cannot reach it MUST record the row as blocked with the reason.
8. MUST verify that the configuration guide, the Spec workflow guide, the
   reopen reference and the Roundfix Skill's `implement` and `settle`
   references describe the behavior, that each skill mirror equals its
   canonical file, and that the raised owned-skill version is recorded.
9. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` defines **Format Command** and **Late Dependency**, and whether
   any other term the Spec introduced needs one.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, and that
    `Makefile`, `go.mod`, `.roundfixrc.yml` and the CI workflows did not
    change. This row reads Task commits, so it declares a `commit_range` input.
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
12. MUST NOT accept a row whose only evidence is that a file was read, MUST NOT
    modify any adopter repository, and MUST NOT send any request to a provider.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Exercise the config, QA format and reopen rows.
- [ ] Re-measure the adopter's failure shape in a temporary repository.
- [ ] Record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The re-measurement shows the precondition passing with the Format
      Command configured, or the row rests on recorded Daemon test evidence
      with its reason.
- [ ] The outside-evidence row records what each source says or why it was
      unreachable.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Build Order 5
- ADR-0249; ADR-0194; ADR-0080; ADR-0091; ADR-0104
