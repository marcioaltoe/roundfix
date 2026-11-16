---
schema: spec-tasks/v1
spec: 0397-each-phase-of-a-spec-is-recorded-and-starts-in-a-fresh-session
qa: task_05
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
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_01, task_02, task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Tasks — Each phase of a Spec is recorded and starts in a fresh session

| id      | title                                                                                                      | type    | complexity | needs                     |
| ------- | ---------------------------------------------------------------------------------------------------------- | ------- | ---------- | ------------------------- |
| task_01 | The Phase Record and `roundfix phase start` record the phase a session starts                               | backend | high       | —                         |
| task_02 | Each prompt of a Run and of a pre-PR review is recorded, and the Run Retention Sweep removes old day files   | backend | high       | task_01                   |
| task_03 | A published authoring skill receives no Declared Skill Header, and this repository declares its models      | backend | high       | task_01, task_02          |
| task_04 | The skills and the autonomous-work guide state the Phase Boundary and keep a fix in the session             | docs    | medium     | task_01, task_02, task_03 |
| task_05 | Run the final QA gate                                                                                      | qa      | high       | task_04                   |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · 4 → task_04 · 5 → task_05

The Tasks run in series, in the order of `_techspec.md` → Build Order. Every
implementation Task raises the version of the Roundfix skill and writes
`skills/testdata/owned-skill-versions.json`; tasks 01 to 03 re-record
`docs/references/behavior-surfaces.json`; tasks 01 and 02 share the store
package, the usage guide and the Skill Coverage Map; tasks 01, 03 and 04
write `CONTEXT.md`; and tasks 03 and 04 both leave the four authoring skills
equal to their mirrors. Each implementation Task runs the suites that consume
its changed contracts in its own Verification, so no separate Integration
Checkpoint is needed.

Sized by the grouping criterion, the Spec holds one context, what a
Supervisor session works on and where an Authoring Phase runs, in four
implementation Tasks and the QA gate, and each Task earns its verification
cycle: a store package with its command, a Daemon and review change with the
retention, an installer rule with this repository's adoption, and a body of
pinned text with a Baseline regeneration.

What the Tasks stage settled beyond the TechSpec, each measured on the
planning branch with the binary built from it:

- The Roundfix skill's entry file is 19,891 bytes against a budget of 20,000
  and its reference count is pinned at 20, so task 01 prescribes the index
  row byte for byte and moves the count to 21.
- `roundfix gc` calls the Run Retention Sweep only when a Run Database
  exists, and Surface Transcript 8 has none, so task 02 counts and removes
  the day files in the GC Command with or without one.
- In fixture D, which has no `.claude` directory, the inspector also names
  the Claude Skill Link of `write-techspec`, a line Surface Transcript 7 does
  not show. Tasks 03 and 05 build fixture E with `.claude/skills` linked to
  `../.agents/skills` after Surface Transcript 6; the transcript's lines are
  unchanged.
- Task 02 changes the source files of the `review` and `implement` commands
  without changing those commands, and their only covering skill files are
  outside the authorization. `_prd.md` → Skills declares the three surfaces
  unchanged and task 02 records the Coverage Reviews.
- The usage guide's Phase Record section is written by task 01 and gains,
  in task 02, who writes a prompt entry and the retention, so that each Task
  documents what it built.

The hazard of `_techspec.md` → Risks & Considerations is planned for: no
Verification command and no QA row runs a `roundfix` found on `PATH`. The one
installation, in task 03, runs the product built from the Task's own tree
after `make skills-sync`, and every `roundfix baseline update` in a Task
carries `--no-skills`.

No other Spec is required.
