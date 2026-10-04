---
schema: spec-tasks/v1
spec: 0222-a-review-through-claude-that-keeps-its-verdict
qa: task_05
graph:
  nodes:
    - id: task_01
      file: task_01.md
      needs: []
    - id: task_02
      file: task_02.md
      needs: []
    - id: task_03
      file: task_03.md
      needs: [task_01, task_02]
    - id: task_04
      file: task_04.md
      needs: [task_03]
    - id: task_05
      file: task_05.md
      needs: [task_04]
---

# Task Graph

| Task | Type | Title |
| --- | --- | --- |
| task_01 | docs | The review reference and guide describe a review through Claude that keeps its verdict |
| task_02 | backend | The ACPX Runner keeps a read-only turn that ended after a refused permission |
| task_03 | backend | The review classifies a verdict delivered after a refused permission and names a prompt that is too long |
| task_04 | backend | The review bounds a Claude prompt in estimated tokens against half of its window |
| task_05 | qa | Run the final QA gate |

Waves: 1 → task_01, task_02 · 2 → task_03 · 3 → task_04 · 4 → task_05

task_01 writes the `review` guide the CLI Tasks change and task_02 adds the
runner field task_03 reads, so task_03 follows both. task_04 shares
`internal/cli/review.go` with task_03, so it follows it. task_01 and task_02
share no declared file.
