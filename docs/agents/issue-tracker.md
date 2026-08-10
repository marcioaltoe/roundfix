<!-- setup-context-driven:begin id=guide.issue-tracker version=0.0.1 -->

# Issue tracker

- **mandatory**: Use the local Spec folder, Task Graph, and Task files as the implementation issue tracker. Do not introduce external triage labels or external issue status as Task state.

- **mandatory**: Keep Task status only in the assigned Task file frontmatter. The Task Graph records topology and dependencies, not progress.

- **mandatory**: Keep each Spec under `docs/specs/<feature-slug>/`. Dependencies live only in `_tasks.md`; status lives only in each Task file, and the local Task files are the published planning issues.

<!-- setup-context-driven:end id=guide.issue-tracker -->
