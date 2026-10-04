---
task: task_05
spec: 0222-a-review-through-claude-that-keeps-its-verdict
status: failed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec keeps the verdict of a read-only
review turn that ended after the session refused a permission, names a
runtime's `Prompt is too long` answer in the review reason, bounds a Claude
review prompt in estimated tokens against half of its context window, and
describes this in the Roundfix Skill and the `review` guide. Behavior rows are
exercised by executing the named tests against the built tree and through the
built binary with disposable repositories and homes. No command reaches GitHub
or the network, except the reads of published pages and the one optional real
review that the outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify API Contract 1 and Success Metric 1 by executing task_02's
   tests against the built tree, and record the cases they assert.
3. MUST verify API Contracts 2 and 3 and Success Metrics 2 and 3 by executing
   task_03's tests against the built tree, record the reasons and record
   fields they assert, and confirm that `TestReviewCommandBlocksOnTransportAnomaly`
   still passes.
4. MUST reproduce Surface Transcript 1 through the built binary in a
   disposable repository whose Project Config selects `claude` for
   `pre_pr_review` and the `review` profile, whose candidate adds a file of
   about 1,100,000 bytes, with `HOME` set to a disposable directory,
   `NODE_OPTIONS` unset and a `PATH` without acpx. It MUST record stdout,
   stderr and the exit code, and that the record's `estimatedPromptTokens`
   exceeds 500000; and MUST run the same candidate with `codex` selected and
   record the unchanged `review diff too large` reason. It MUST also verify
   API Contract 4 and Success Metric 4 by executing task_04's tests.
5. MUST verify that the Roundfix Skill's `review` reference and the `review`
   command guide describe the behavior task_01 names, that each mirror equals
   its canonical file, that the version was raised and recorded, and that the
   `### QA settlement` section of every skill is unchanged.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the two Secondbrain reports, the three
   authoring reviews, the installed acpx and Claude adapter, and the two
   Anthropic guides. For each source it reaches it MUST record what it read
   and whether it still supports the design; for each it cannot reach it MUST
   record the row as blocked with the reason. When a real `claude / opus`
   review is available, it MAY run at most one, through the built binary, in a
   disposable clone whose ignored `CLAUDE.local.md` asks the reviewer to run
   `go vet ./internal/cli/` before its verdict, and record whether acpx exited
   `5`, the outcome, `permissionRefused` and the finding ids; when it is not
   available, the row is blocked with that reason.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Pre-PR Review Command**, **ACPX Runner** and **Agent Session** still
   describe the behavior, and whether "refused permission" needs a glossary
   term.
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
    anything on GitHub, and MUST NOT write to the live Run Database or
    artifacts under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0222-a-review-through-claude-that-keeps-its-verdict/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; User Stories 1-4; Core Features 1-5; Success Metrics 1-4;
Acceptance evidence; `_techspec.md` → Testing Approach 1-3; API Contract 1;
API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1;
Integration Points; ADR-0080; ADR-0088; ADR-0091; ADR-0104; ADR-0155;
ADR-0167; ADR-0227.
