---
schema: roundfix/archive-record/v1
spec: 0136-a-rename-the-committer-can-stage
title: A rename the committer can stage
status: archived
created: "2026-09-14"
archived: "2026-09-14"
disposition: pass
source: docs/history/specs/0136-a-rename-the-committer-can-stage
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_04
qa_report: qa-report-2026-09-14-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "188"
delivery_commit: 30069a624aa4adb01062311b43e2109b4915eaa8
---

# A rename the committer can stage

Spec 0135 made the changed-path reader carry a rename's source, so the governed side of a rename reaches the classifier. That set is read by two stages, and only one of them wanted the extra path. The commit stage stages explicit paths with `git add -f`, and after `git mv` the source exists in neither the worktree nor the index, because the index already records the rename. Git answers `pathspec did not match any files` and exits non-zero, so a Task that renames a file with `git mv` now fails to commit. Measured directly against Git: the same command succeeds for an unstaged rename and fails for a staged one.
