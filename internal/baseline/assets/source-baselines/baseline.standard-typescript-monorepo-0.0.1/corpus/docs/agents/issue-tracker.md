<!-- source-baseline-entry: clause.spec.local-task-tracker-only -->
- MUST use the local Spec folder, Task Graph, and Task files as the implementation issue tracker. Do not introduce external triage labels or external issue status as Task state.
<!-- /source-baseline-entry: clause.spec.local-task-tracker-only -->

<!-- source-baseline-entry: clause.spec.status-only-in-task -->
- MUST keep Task status only in the assigned Task file frontmatter. The Task Graph records topology and dependencies, not progress.
<!-- /source-baseline-entry: clause.spec.status-only-in-task -->

<!-- source-baseline-entry: clause.spec.tracker-artifacts -->
Keep each Spec under `docs/specs/<feature-slug>/`. Dependencies live only in `_tasks.md`; status lives only in each Task file, and the local Task files are the published planning issues.
<!-- /source-baseline-entry: clause.spec.tracker-artifacts -->
