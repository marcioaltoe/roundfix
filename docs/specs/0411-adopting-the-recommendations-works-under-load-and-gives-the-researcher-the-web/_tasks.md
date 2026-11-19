---
schema: spec-tasks/v1
spec: 0411-adopting-the-recommendations-works-under-load-and-gives-the-researcher-the-web
qa: task_03
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: [task_01]
    - id: task_03
      file: task_03.md
      needs: [task_02]
---

# Tasks — Adopting the recommendations works under load and gives the researcher the web

| id      | title                                                                                     | type    | complexity | needs   |
| ------- | ----------------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | A profile write proves under the configured setup deadline                                 | backend | medium     | —       |
| task_02 | The recommendation gives the Research Subagent the web, and the checks name the missing tools | backend | high       | task_01 |
| task_03 | Run the final QA gate                                                                     | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The graph follows `_techspec.md` → Build Order. The two implementation Tasks share no code, but both edit `docs/user-guide/configuration.md`, raise the Roundfix skill's version in both copies of its `SKILL.md` and `skills/testdata/owned-skill-versions.json`, and re-record `docs/references/behavior-surfaces.json`, so task_02 needs task_01 under the Wave collision rule. The QA gate needs task_02 and, through it, every non-QA leaf. Each Task runs the suites that consume its changed contracts in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds two defects that share no component, kept together by the maintainer's answer "Juntar à Spec 0411 (Recommended)" (`_prd.md` → Decisions). Each half is one Task that earns its verification cycle: a reader with an eight-case agreement test, the wiring with six adoption cases and a transcript, and three documents; and a recommendation change, a two-check report with five cases, a generator test, a repository contract, five documents, the glossary entry and an ADR amendment. Splitting either Task further would leave a slice too small for its own Run turn, Verification and settlement.

## What the Tasks stage settled beyond the TechSpec

Each item was measured on a prototype on a scratch worktree of the planning branch, based on `a6a4a15b`, removed afterwards; the measurements are in [`references/2026-11-18-characterization.md`](references/2026-11-18-characterization.md).

- **A full configuration load breaks the repair path.** Resolving the whole configuration, as Setup does, failed `TestProfilesConfigureProofScope` and `TestProfilesConfigureChangeSummary`; the narrow reader keeps both passing unchanged.
- **Existing tests that change.** Only the recommendation pins move: `recommendedAuthoringYAML` and two subtests in `internal/config`, `authoringRecommendationBlock` in `internal/cli` (read by five tests, `TestRunDoctorRealRepositoryCheckDoesNotMutateState` among them, which needs no edit of its own), and the configuration guide's block the documentation contract compares. `TestTheSkillCoverageMapIsCurrent` fails after each guide edit until the Task re-records it as its last step.
- **The report reads the loaded configuration.** The installation's declaration copies the research tools with `append([]string{}, …)`, which turns a missing key into an empty list; only `loaded.Config.Skills.Authoring.Claude` tells `tools: []` from no key.
- **Through the product.** With the prototype binary, the Update Check rows of task_03 requirement 5 printed exactly the prescribed JSON and line, needed the `.claude/skills/write-prd` link to reach the `applied` state, and `profiles configure` printed Surface Transcript 1 without starting any proof.
- **Owned-skill versions.** Each Task raises the Roundfix skill's version once on the tree it starts from; Spec 0408, delivered before this one, raises it too.
- **Exact text.** Both implementation Tasks copy their text word for word from `_techspec.md` → API Contracts 2, 4 and 5, so the QA gate searches for fixed strings.

## Time in tests

No test of this Spec reads the clock or waits on a timer. The proof deadline is asserted as the value a fake runner receives, never as an elapsed time.
