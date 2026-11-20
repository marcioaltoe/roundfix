---
schema: spec-tasks/v1
spec: 0416-the-secondbrain-search-and-the-duplicate-check-go-through-the-indexus-without-qmd
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

# Tasks — The Secondbrain search and the duplicate check go through the Indexus, without qmd

| id      | title                                                                                                        | type    | complexity | needs   |
| ------- | ------------------------------------------------------------------------------------------------------------ | ------- | ---------- | ------- |
| task_01 | The Secondbrain search and the duplicate check name the Indexus in the Baseline, in write-idea and in a decision record | backend | high       | —       |
| task_02 | Run the final QA gate                                                                                        | qa      | high       | task_01 |

Waves: 1 → task_01 · 2 → task_02

The graph follows `_techspec.md` → Build Order: steps 1 and 2 are task_01 and step 3, the gate, is task_02. Step 2 joins step 1 under the sizing rule of `docs/agents/spec-routing.md`: it is one paragraph of one owned skill, its mirror, its version record and one test file, the same replacement in the same contract, since the skill step repeats the guide's clause. As its own Task it would pay an Agent turn and the repository Verification for five lines. The joined Task runs the order the prototype proved: both edits, `make skills-sync`, the skill record command, the module record command, `make baseline-digests`, the decision record and the notice, then the Managed Refresh. The clause edit and the decision record stay in one Task because the PRD requires that the record and the guide never disagree on the default branch.

The Task's type is `backend`: its outcome is the text the Baseline catalog ships and `roundfix baseline update` writes into every adopter, proven by tests over the catalog, the plan and the embedded skills.

## What the Tasks stage settled beyond the TechSpec

- **The ADR number.** 0449. The default branch ends at decision record 0447, and the open planning Pull Request of Spec 0415 holds `docs/adr/0448-the-owner-of-a-skill-activation-requires-every-skill-of-its-bundle.md` (checked on 2026-11-20), which settles `_prd.md` → Open Questions, item 4. When 0449 is taken before the Run, the operator renames the file and its link and amends task_01; invariant 5 keeps the test independent of the number.
- **Subtest names.** API Contract 4 asks for one subtest per phrase; task_01 names each one, so the Verification and the gate check them by name.
- **Two PRD units corrected.** Core Feature 4 and Success Metric 2 said the update records both clauses as retained. `clause.secondbrain.research-to-the-knowledge-owner` has no disposition of its own: it is the target of the `replaced` record of `clause.secondbrain.research-capture`, and only `clause.secondbrain.02-query-order` is `retained`. Both units now state that accounting, so the gate reads them literally, and the TechSpec's deviation for them was removed.
- **The Governed Paths.** Measured again over every file task_01 declares; the seven paths of `_authorization.md` are exactly the governed ones, so the record is unchanged.

## Shared files

Spec 0415 is planned on another branch and its plan adds decision record 0448. If its implementation changes a Baseline module or an owned skill first, the two Specs share the module version record, the catalog files, the plan goldens, the Standard TypeScript profile's pin, this repository's Setup Manifest and the owned-skill version record; task_01 regenerates them on the tree it runs on and writes "the next version", never a number. No Verification of this Spec reads another Spec's artifact, and the graph needs no `requires`.

## Delivery

The authorization record grants `implement`, `commit`, `push` and `pull_request`, and neither `merge` nor `release`, so the implementation stops at an open Pull Request (`_prd.md` → Delivery limits). No Task sends a message.

## Time in tests

No test of this Spec reads the clock or waits on a timer; every new test reads the embedded catalog, a mutated in-memory copy of it, the embedded skills, this repository's files or files in a temporary repository.
