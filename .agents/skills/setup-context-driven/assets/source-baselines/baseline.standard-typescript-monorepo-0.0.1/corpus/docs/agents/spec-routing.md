<!-- source-baseline-entry: contract.spec.route-matrix -->
# Spec route matrix

| Change | Required route |
| --- | --- |
| Large or fuzzy product initiative | `write-idea` → `write-prd` → `write-techspec` → `write-tasks` |
| Standard feature that changes product behavior | `write-prd` → `write-techspec` → `write-tasks` |
| Refactor or bug fix without product-behavior change | minimal `_prd.md` → `write-techspec` → `write-tasks` |
| Trivial one-line fix, typo, or configuration change with obvious intent and Verification | Direct implementation |

Use the smaller sufficient route when two routes fit. Use brainstorming before creative or feature work. A TechSpec can be omitted only when the feature has no architectural surface. Implementation always executes from the Task Graph.
<!-- /source-baseline-entry: contract.spec.route-matrix -->

<!-- source-baseline-entry: contract.spec.task-ownership -->
# Local Task ownership protocol

1. Keep every Spec under `docs/specs/<feature-slug>/` with its idea, PRD, TechSpec, Task Graph, Task files, and `qa/` evidence.
2. Dependencies live only in `_tasks.md`; Task status lives only in each Task file's frontmatter.
3. Local Task files are the planning issues. External issue trackers are unsupported by this baseline.
4. An Agent executing one Task changes only that Task's slice and records fresh acceptance evidence in that Task file.
5. Final QA begins only after every Task is completed.
6. Archive only a completed Spec with a passing QA verdict under `docs/specs/_archived/`.
<!-- /source-baseline-entry: contract.spec.task-ownership -->
