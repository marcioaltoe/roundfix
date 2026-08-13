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
6. Archive only a completed Spec with a passing QA verdict under `docs/history/specs/`.
<!-- /source-baseline-entry: contract.spec.task-ownership -->

<!-- source-baseline-entry: clause.spec.routing-01-large-initiative -->
Large or fuzzy product initiative: run `write-idea` → `write-prd` → `write-techspec` → `write-tasks`.
<!-- /source-baseline-entry: clause.spec.routing-01-large-initiative -->

<!-- source-baseline-entry: clause.spec.routing-02-standard-feature -->
Standard feature that changes product behavior: run `write-prd` → `write-techspec` → `write-tasks`; skip the TechSpec only when the feature has no architectural surface.
<!-- /source-baseline-entry: clause.spec.routing-02-standard-feature -->

<!-- source-baseline-entry: clause.spec.routing-03-refactor-bugfix -->
Refactor or bug fix without product-behavior change: run `write-techspec` → `write-tasks` and create the minimal `_prd.md` required by the downstream artifact contract.
<!-- /source-baseline-entry: clause.spec.routing-03-refactor-bugfix -->

<!-- source-baseline-entry: clause.spec.routing-04-trivial-change -->
Trivial one-line fix, typo, or configuration tweak: implement directly without a Spec only when intent, acceptance criteria, and Verification are obvious.
<!-- /source-baseline-entry: clause.spec.routing-04-trivial-change -->

<!-- source-baseline-entry: clause.spec.routing-05-task-graph -->
Use `brainstorming` before creative or feature work, start with the smaller sufficient route when two routes fit, and execute implementation from the Task Graph.
<!-- /source-baseline-entry: clause.spec.routing-05-task-graph -->

<!-- source-baseline-entry: clause.spec.verification-two-tiers -->
- MUST use the active Baseline Profile's declared incremental verification command for each Task to answer whether the current slice remains valid before handoff. CI MUST use the Profile's declared complete verification command from a fresh run to answer whether the assembled tree satisfies the repository contract. A missing incremental command leaves the Profile's two-tier contract unmet and never authorizes skipping the local tier.
<!-- /source-baseline-entry: clause.spec.verification-two-tiers -->

<!-- source-baseline-entry: clause.spec.project-constraints-06-outside-evidence -->
- MUST rest a Spec's acceptance, in at least one named row, on evidence originating outside the Spec's own artifacts — a repository the Spec did not build, a measurement it did not design, or published literature — and record in that row where the evidence came from. A row whose outside source cannot be obtained is recorded as blocked with its reason; it never requires human interaction and never blocks the Spec.
<!-- /source-baseline-entry: clause.spec.project-constraints-06-outside-evidence -->
