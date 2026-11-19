---
schema: spec-tasks/v1
spec: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
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

# Tasks — The status counts a delivering Spec's Tasks, and a Roundfix update needs no release

| id      | title                                                                              | type    | complexity | needs   |
| ------- | ---------------------------------------------------------------------------------- | ------- | ---------- | ------- |
| task_01 | The status view counts a Delivery Queue item's Tasks on its Item Branch            | backend | medium     | —       |
| task_02 | In an adopter, a change to the files Roundfix manages or to a CONTEXT.md needs no release | backend | high       | task_01 |
| task_03 | Run the final QA gate                                                              | qa      | high       | task_02 |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03

The graph follows `_techspec.md` → Build Order. The two implementation Tasks share no component, but both edit `CONTEXT.md` and re-record `docs/references/behavior-surfaces.json`, so they run in series under the Wave collision rule: task_01 changes the status view's counts and the Task Graph loader with the `follow-up` guide, the Spec Queue Table entry and the Coverage Review; task_02 changes the Release Plan's boundary with the help, the release documents, the Roundfix skill version and the Adopter Maintenance Path entry. Each runs the suites that consume its changed contracts in its own Verification, so no separate Integration Checkpoint is needed.

Sized by the grouping criterion, the Spec holds two defects of one adopter's status report that share no component, kept together by the maintainer's answer "Spec fix 0414 logo após a 0411 (Recommended)" (`_prd.md` → Decisions). Each half is one Task that earns its verification cycle: a loader seam, a reader with eight cases and two documents; and a classifier rule with seventeen path cases, five repository cases, a help paragraph and four documents pinned phrase by phrase. Splitting them would give each a one-Task Spec paying a whole Run, QA gate, review and release.

## What the Tasks stage settled beyond the TechSpec

Each item was measured on a prototype on a scratch worktree of the planning branch, based on `145d6c76`, removed afterwards; the measurements are in [`references/2026-11-19-characterization.md`](references/2026-11-19-characterization.md).

- **Existing tests.** With the prototype, every test of `./internal/spec` and `./internal/releaseplan`, and every `./internal/cli` test matching `FollowUp|Queue|SpecIsBeingAuthored|RenamedSpec|ReleasePlan|Release`, passed unchanged; items without a branch keep `queued 1/2 Tasks` and `parked 0/0 Tasks`.
- **The Behavior Surface Record.** task_01's guide edit and task_02's help and guide edits each move it, so each re-records it as its last edit. task_02 changes the Roundfix release reference, which covers `command release plan`, `guide docs/user-guide/release-runbook.md` and `guide docs/user-guide/usage.md`, so those surfaces read as described; task_01 records the Coverage Review that `_prd.md` → Skills declares for the two status surfaces.
- **Project Config is not read.** The prototype binary planned an adopter range whose target holds an invalid `.roundfixrc.yml` and ended `no_release`, so `project_config_is_not_read` holds by design.
- **Owned-skill version.** Spec 0411, delivered first, also raises the Roundfix skill; task_02 raises it to the next version with the record command on its own tree.
- **Exact text.** Both implementation Tasks copy their text word for word from `_techspec.md` → API Contracts 1 and 2, so the QA gate searches for fixed strings.

## Time in tests

No test of this Spec reads the clock or waits on a timer. The release-plan tests create their own tags in temporary repositories, and no test depends on this repository's tags.
