---
task: task_05
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes the Delivery Engine, its
command workflow, the Task Graph manifest reader, the Project Config and the
documents that state them. Every behavior row is exercised through the built
binary or by executing the named tests against the built tree. Each binary
command runs in a disposable repository with a disposable Roundfix Home and a
bare local origin. No command reaches GitHub, runs `gh` against a real
repository, starts an ACP adapter or opens the network.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary that `roundfix deliver start` of two Specs that require each other reproduces Surface Transcript 2, exits `2` and records no queue, and that an unknown prerequisite is refused before any queue is recorded.
3. MUST verify, by executing `TestDeliverStatusReproducesSurfaceTranscriptOne` against the built tree, that `roundfix deliver status` reproduces Surface Transcript 1, and, through the built binary on a queue without a parked item, that the status output is unchanged.
4. MUST verify, by executing task_01's tests against the built tree, that every blocker constant has a Park Class other than `unclassified`, that a check failed outside the change is re-run exactly once, and that a failure in a changed package, a build failure or an unattributable log parks `checks-failed` with no re-run.
5. MUST verify, by executing task_02's tests against the built tree, that an item waits for a prerequisite ahead in the queue, parks `prerequisite-unmerged` without a worktree when the prerequisite is parked, and returns to `queued` with an unchanged retry count once the prerequisite's archive is on the refreshed default branch.
6. MUST verify, by executing task_03's tests against the built tree, that an environment-only partial parks `qa-environment-partial`, that a retry of the operator-archived item resumes at `reviewing` and reaches `gating` with no new archive commit, that the same retry without a QA override is refused with today's text, and that the authorization is read before the archive commit.
7. MUST verify, by executing task_04's tests against the built tree, that a `CONFLICTING` Pull Request stops the check wait on the first read, that a derived conflict is merged and regenerated into a commit with the `Roundfix-Delivery: derived-merge` trailer, and that a source conflict aborts the merge and names only the undeclared path.
8. MUST verify through the built binary that a Project Config with a malformed `delivery.derived_paths` entry is refused with the key named, and that a config without the key loads as before.
9. MUST record, as evidence this Spec did not author, the sources quoted in `_prd.md` → Acceptance evidence: the finding's recorded session evidence for PR #298, the two 0195 Runs and #296, checked against the finding's text; and the GitHub GraphQL `MergeableState` reference, the DeFlaker paper, the generated-file conflict guide and the GitHub CLI run manual. When a published page is reachable, it MUST record what it read there; when it is not, it MUST record the row as blocked with its reason.
10. MUST verify that the delivery command guide, the configuration guide, the Roundfix skill's delivery reference and its mirror describe the prerequisites, the Park Classes, the re-run, the operator-archived retry and the conflict handling; that the `### QA settlement` section of the skill is unchanged; that the skill's version is recorded; and that `make skills-sync-check` exits `0`.
11. MUST check whether this Spec introduced, changed or retired a glossary term (Park Class, prerequisite Spec, derived path declaration) and record the result for `docs/agents/domain.md`'s glossary rule, without editing `CONTEXT.md` inside this gate.
12. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that `internal/cli/cli_test.go` and `.roundfixrc.yml` did not change.
13. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows archive;
    - review-artifact ancestry: the audited head named as the claimed candidate.
14. MUST NOT accept a row whose only evidence is that a file was read.
15. MUST NOT write to the live Run Database under `~/.roundfix`, start a session on a real ACP adapter, or run `gh` or `git push` against a real remote.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; User Stories 1-6; Core Features 1-5; Success Metrics
1-5; Acceptance evidence; `_techspec.md` → Testing Approach 1-5; API Contracts
1-4; Surface Transcripts 1-2; ADR-0080; ADR-0088; ADR-0091; ADR-0093;
ADR-0094; ADR-0096; ADR-0104; ADR-0117; ADR-0155; ADR-0156; ADR-0166;
ADR-0167; ADR-0176; ADR-0192; ADR-0193.
