---
schema: spec-tasks/v1
spec: 0073-skill-versions-decoupled-from-the-binary
qa: task_06
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
    - id: task_04
      file: task_04.md
      needs: [task_02]
    - id: task_05
      file: task_05.md
      needs: [task_03, task_04]
    - id: task_06
      file: task_06.md
      needs: [task_05]
---

# Tasks — Skill versions decoupled from the binary

| id      | title                                          | type    | complexity | needs            |
| ------- | ---------------------------------------------- | ------- | ---------- | ---------------- |
| task_01 | Let every owned skill declare its version       | chore   | medium     | —                |
| task_02 | Compare a version instead of matching bytes     | backend | high       | task_01          |
| task_03 | Stop gating compatibility on content            | backend | high       | task_02          |
| task_04 | Make every surface use the one comparison       | backend | medium     | task_02          |
| task_05 | Synchronise the Roundfix Skill                  | chore   | low        | task_03, task_04 |
| task_06 | Run the final QA gate                           | qa      | medium     | task_05          |

Waves: 1 → task_01 · 2 → task_02 · 3 → task_03 · task_04 · 4 → task_05 ·
5 → task_06

task_01 leads because it is inert. Skills gain a frontmatter field nothing
reads yet, so the repository stays green while the identity exists for task_02
to compare against. Doing it the other way — comparing before declaring — would
leave every owned skill `unversioned` at once and make the transition
indistinguishable from a regression.

task_03 and task_04 are independent of each other and share task_02's
comparison. One removes content from the compatibility gate; the other makes
Doctor and the gating commands call the same code, so two surfaces cannot
disagree.

task_01 and task_05 carry authorized tooling scope under the 2026-08-04 record,
which names Spec 0073: `Makefile` and, for any member of `OWNED_SKILLS`,
`.agents/skills/<owned-skill>/**` with its mirror. That set is read from the
`Makefile` variable — no Task in this graph copies the member list, because a
copied list is the defect the repository's "an assertion reads the constant it
means" rule names, and this Spec's whole subject is a pin that drifted from
what it pinned.

No acceptance in this graph asserts a recorded digest or a recorded version.
Each asserts the comparison's outcome, which holds at any version.

The QA gate is authored as task_06 and depends on task_05, the graph's only
non-QA leaf.
