---
task: task_04
spec: 0223-adjustments-the-adopters-asked-for
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes `branch.prefix` an optional
Baseline decision, adds the re-execution rule for a row without inputs to the
`qa-gate` skill, and lets `roundfix release plan` report the skills and
baseline checks. Every behavior row is exercised through the built binary or by
executing the named tests against the built tree, with disposable repositories
and a disposable home. No command reaches GitHub, a provider or the network,
except the reads of published pages the outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST reproduce Surface Transcripts 1 and 2 through the built binary with
   `HOME` set to a disposable directory and `NODE_OPTIONS` unset, in
   disposable repositories tagged `v0.4.0`, one not adopted and one adopted
   with this repository's skills copied in, recording stdout, stderr, the exit
   code, the JSON `checks` object of `--format json`, and `git status
   --porcelain` and a checksum of the tree before and after each run. It MUST
   also show that `roundfix release plan --reset-to` and a refused plan print
   no check line, and that `roundfix doctor`'s `skills:` line reads the same as
   the plan's in the adopted repository.
3. MUST verify, by executing task_02's tests against the built tree, Fixed
   texts 1 and 2 and Success Metrics 1 and 2, and MUST adopt a disposable
   repository with the binary built from this Spec's base commit, recording
   `branch.prefix`, then run `roundfix baseline update` with the built binary
   and record that it reports no new decision and leaves
   `docs/agents/agent-instructions.md` byte-identical. It MUST also run the
   same update on a copy whose Setup Manifest has no `branch.prefix` and record
   the rendered paragraph.
4. MUST verify that this repository's `docs/agents/agent-instructions.md` is
   byte-identical to the base commit's, that a second
   `baseline update --repo . --no-skills` reports `current`, and that no
   Baseline Normative Clause was removed or renamed.
5. MUST verify that the `qa-gate` skill carries the re-execution paragraph of
   the TechSpec verbatim, that the Roundfix Skill release reference, the
   release runbook and the Context-Driven Development guide describe the
   behavior task_01 names, that each mirror equals its canonical file, that both
   versions were raised and recorded, and that the `### QA settlement` section
   of every skill is byte-identical.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the oraculum Inbox Entry and oraculum's
   rendered guide in the Secondbrain mirror, the fluxus Inbox Entry on the
   `qa-gate` skill, and the release runbook's check run against this
   repository at `v0.33.0`. For each source it reaches it MUST record what it
   read and whether it still supports the design; for each it cannot reach it
   MUST record the row as blocked with the reason.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Baseline Profile**, **Setup Manifest** and **QA Report** still describe
   the behavior, and whether any new term is needed.
8. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that `Makefile`,
   `.roundfixrc.yml`, `go.mod` and the CI workflows did not change. This row
   reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
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
10. MUST NOT accept a row whose only evidence is that a file was read.
11. MUST NOT start a Delivery Queue, merge, re-run or otherwise change
    anything on GitHub, and MUST NOT write to the live Run Database under
    `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0223-adjustments-the-adopters-asked-for/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; User Stories 1-4; Core Features 1-6; Success Metrics 1-4;
Acceptance evidence; `_techspec.md` → Testing Approach 1-3; API Contract 1;
API Contract 2; API Contract 3; API Contract 4; API Contract 5; Surface
Transcript 1; Surface Transcript 2; Integration Points; ADR-0080; ADR-0088;
ADR-0091; ADR-0104; ADR-0155; ADR-0167; ADR-0228.
