---
schema: spec-tasks/v1
spec: 0419-round-2-of-a-pre-pr-review-opens-a-new-session-when-the-configured-model-changed
requires: [0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks]
qa: task_02
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
---

# Tasks — Round 2 of a pre-PR review opens a new session when the configured model changed

| id      | title                                                                                              | type    | complexity | needs   |
| ------- | -------------------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | Round 2 continues the reviewer session only on the recorded model, and the rule is documented      | backend | medium     | —       |
| task_02 | Run the final QA gate                                                                              | qa      | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02

The graph follows `_techspec.md` → Build Order one step per Task: step 1 is task_01 and step 2, the gate, is task_02. Step 1 stays one Task, as the TechSpec decides: the code is three lines, the documents describe only those lines, and the guide page, the skill reference and the Behavior Surface Record cannot be split from one another, because a changed guide page fails the surface contract until its record is re-recorded and fails the release check until its covering skill changes in the same range. A second Task for about ten lines of prose would pay an Agent turn and a repository Verification of its own.

task_01 is `backend`: its outcome is the behavior of the Pre-PR Review Command and the record it writes; the documents follow that behavior.

## What the Tasks stage settled beyond the TechSpec

- **`requires`.** The manifest names Spec 0418 in `requires`, by the decision of the Supervisor session of 2026-11-23 that started this phase. Both Specs change the Roundfix skill's version fields, `skills/testdata/owned-skill-versions.json`, `docs/references/behavior-surfaces.json`, `CONTEXT.md` and `.secondbrain-export`, and the maintainer's third decision delivers this Spec after Spec 0418 (`_prd.md` → Decisions). With the entry, the Delivery Queue starts this Spec only once Spec 0418 has merged, so task_01 runs on a tree that holds it. No Verification of this Spec reads an artifact of Spec 0418. `_prd.md` → Prerequisites and `_techspec.md` → Build Order now say so.
- **Test and subtest names.** The TechSpec names the six new tests; task_01 names their subtests so the Verification and the gate check them by name: `the_classifier_still_refuses_the_reported_answer` under the first test of `internal/cli`, `preferred_then_continued` and `fallback` under `TestReviewRecordsTheModelOfTheSessionOfEachRound`, and `guide`, `canonical_reference` and `mirror_reference` under `TestTheReviewGuideAndReferenceDescribeTheModelRule`. The test of `internal/cli/review_lineage_selection_test.go` that pins a `ceiling-closed` lineage is `TestACeilingRecordNamesNoSelection`.
- **The ADR number.** 0451. 0450 belongs to Spec 0418, and the default branch ends at decision record 0449 (checked on 2026-11-23). When 0451 is taken before task_01 runs, the operator renames the file, its link in the notice of ADR-0197, its link in the glossary and its export line, and amends task_01; invariant 7 keeps the decision record test independent of the number.
- **The Governed Paths.** Measured again on 2026-11-23 at `50143e74` with `speccheck.GovernedPath`, in a `go test -overlay` probe whose file lives outside the repository, over every file the two Tasks declare: the three paths of `_authorization.md` are the only governed ones, and only task_01 declares them. The record is unchanged.

## Shared files

Spec 0418 changes the Roundfix skill's version fields, the owned-skill version record, the Behavior Surface Record, `CONTEXT.md` and `.secondbrain-export` before this Spec is delivered. task_01 regenerates the shared records on the tree it runs on and writes "the next version", never a number, and no Verification or test of this Spec pins a version or a digest of a shared file.

## Delivery

The authorization record grants `implement`, `commit`, `push`, `pull_request`, `merge` and `release` at a `patch` Release Ceiling. No Task performs the release, edits `CHANGELOG.md`, a version file or the npm manifest, or sends a message; the operator cuts the release and tells the Fiscus, as the maintainer's first decision asks.

## Time in tests

No new test waits on elapsed time. The tests of `internal/cli` end on the command's return over a fake runner; the tests of `internal/docscontract` read files.
