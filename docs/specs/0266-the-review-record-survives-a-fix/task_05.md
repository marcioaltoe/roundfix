---
task: task_05
spec: 0266-the-review-record-survives-a-fix
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec keys the Pre-PR Review Record by
the resolved checkout path, compares every record repository through one
resolving helper, keeps a Reviewed Head Record per reviewed head, makes the
Delivery Retry's correction proof read the parked candidate's record, widens
the proof's path rule to the files the archived Spec's Tasks declared, and
documents all three in the guides, the Roundfix skill and `CONTEXT.md`. No row
reads or writes the real `~/.roundfix`, follows the maintainer's real worktree
link, reaches a provider or reads a credential.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify API Contract 1, API Contract 2 and Invariants 1-2 by executing
   task_02's tests by name, and by repeating task_02's sweep search for
   lexical `Repository` comparisons in `internal/cli/review.go`,
   `internal/cli/review_lineage.go`, `internal/cli/review_validation.go` and
   `internal/cli/deliver_review_correction.go`, recording that it finds none.
3. MUST verify API Contract 3, API Contract 4, Invariant 3 and Success Metric
   2 by executing task_03's tests by name, and the existing
   `TestAReviewOnlyCorrectionIsProvedFromTheReviewRecord`,
   `TestACorrectionOutsideTheArchivedSpecIsRefused`,
   `TestACorrectionWithAnUndisposedFindingIsRefused`,
   `TestACorrectionThatDoesNotDescendIsRefused` and
   `TestReviewCorrectionRequiresTheParkedReviewRecord`, which pin ADR-0233's
   descent and disposition proofs.
4. MUST verify API Contract 6, Invariants 4-5 and Success Metric 4 by
   executing task_04's tests by name, including both subtests of
   `TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview`, and MUST repeat
   the second characterization of `_techspec.md` → Current behavior on the
   built tree: a disposed fix to a file the archived Spec's Task declared is
   now accepted, and the same fix plus a change to an undeclared file is
   refused with a reason naming that file.
5. MUST verify Success Metric 1 and User Stories 1-3 through the built
   `bin/roundfix` with `HOME` set to a temporary home, `NODE_OPTIONS` unset,
   and a temporary Git repository whose linked worktree is added under a
   worktree root that is a symbolic link to a temporary directory: when a
   review can be driven without a real provider through the binary, record
   the review directory that `roundfix review` writes and that
   `roundfix review dispose` reads from both spellings of the worktree;
   otherwise the row records task_02's and task_03's tests as its evidence
   and says why. The Delivery Retry itself needs a parked Delivery Queue item,
   so its proof is evidenced by those tests.
6. MUST repeat the measurement of `_techspec.md` → Current behavior with the
   installed Git: add a worktree through a symbolic-linked directory in a
   temporary repository and record that `git worktree list --porcelain` and
   `git rev-parse --show-toplevel` run there name the resolved directory.
7. MUST record, as evidence this Spec did not author, each published source
   of `_prd.md` → Acceptance evidence: the Stack Overflow answer and the
   git-for-windows issue. When it reaches a source it records whether the
   quoted passage is still there and still supports the design; when it
   cannot, it records `blocked (environment: network denied: <host>)` and
   relies on the quoted passage and the measurement of Requirement 6.
8. MUST perform the glossary check of `docs/agents/domain.md`. It records
   that `CONTEXT.md` defines **Reviewed Head Record** and the revised
   **Pre-PR Review Record**, and whether any other term the Spec introduced
   needs one.
9. MUST search the user guides, `CONTEXT.md`, the agent guides and the
   Roundfix skill references for wording that says the review directory is
   named from the checkout path as spelled, that a later review replaces the
   record a retry reads, or that the retry admits only paths under the
   archived Spec, and record the result; and MUST record that both
   copies of the review and deliver skill references carry task_01's phrases
   and that `TestEveryOwnedSkillVersionIsRecorded` and
   `TestTheSkillCoverageMapIsCurrent` pass.
10. MUST verify from the repository history that each Task's changed files
   stay within its declarations or its `## Recorded paths`; that every
   Governed Path changed is one of the four the `_authorization.md` `paths:`
   bound; that no file under `internal/delivery` or `internal/config` changed;
   that no Baseline module, Makefile, CI workflow, `go.mod` or `go.sum`
   changed; that no skill other than `roundfix` changed;
   and that no `### QA settlement` section changed, compared with the
   delivery base (the merge base of the audited head with `main`). This row
   reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, use a broad `pkill`, send any request to
    a provider, start a Delivery Queue or a workflow, or persist an
    environment value.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests, the pinned proof tests and the sweep search.
- [ ] Exercise the review directory through the binary under a linked root.
- [ ] Record the Git measurement, the outside evidence, the glossary and wording checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The proof accepts a disposed docs-only fix under a linked root, also after round 2.
- [ ] A disposed code fix inside the declared paths returns to review; one outside them is refused naming the path.
- [ ] ADR-0233's descent and disposition refusals still hold.
- [ ] The scope audit finds no undeclared change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0266-the-review-record-survives-a-fix/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-6; User Stories 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Invariants 1-5; Current behavior; Build Order 5
- ADR-0272; ADR-0233; ADR-0165; ADR-0174; ADR-0197; ADR-0261; ADR-0259
