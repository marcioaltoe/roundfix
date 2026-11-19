---
spec: 0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release
prd: _prd.md
created: 2026-11-19
---

# The status counts a delivering Spec's Tasks, and a Roundfix update needs no release — Technical Spec

## Executive Summary

Two independent corrections ship together. The status view of the Follow-up Command reads a Delivery Queue item's Task Graph from its Item Branch through Git objects, with the Task Graph loader given a reader instead of the file system, and falls back to the Spec Root exactly where it reads today. The Release Plan gains an optional source of Adopter Maintenance Paths, read from the Setup Manifest and `skills-lock.json` at the range's base and target, and counts those paths as maintenance-only when the target's lock records the `roundfix` skill (ADR-0447). The trade-off accepted for the count is staleness by at most the commit that follows a Task's settlement: the branch, not the worktree, is read, so the view keeps its promise to read no file of another checkout. The trade-off accepted for the release is that an adopter whose lock lost the `roundfix` entry keeps the conservative classification until its next Baseline update.

## Project Constraints

- Identifier strategy: not applicable — no identifier is minted; the reader uses the Item Branch name and the Archive Record the Delivery Queue and the branch already hold, and the Release Plan reads existing skill names and paths. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or network call; both readers run `git` on the local repository. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0447 (new); ADR-0412, ADR-0349 and ADR-0359 constrain the release rule; ADR-0410, ADR-0434, ADR-0426, ADR-0057 and ADR-0071 bear on the status count. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: "Spec fix 0414 logo após a 0411 (Recommended)" with "Mesmas do BE", of 2026-11-19, quoted in the record; it bounds the three Roundfix skill files below and grants no Makefile, workflow, Go module, lint or formatter change. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0414-the-status-counts-a-delivering-specs-tasks-and-a-roundfix-update-needs-no-release/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/release.md`, `skills/roundfix/SKILL.md`.

## System Architecture

Four components change; nothing new is wired between them.

1. The Task Graph loader, `spec.Load` in `internal/spec/spec.go`. Its three file reads (`_prd.md` through `requireActive`, `_tasks.md` through `loadManifestNodes`, each Task file through `loadTask`) take a reader. `spec.Load` and `spec.LoadForRecovery` pass `os.ReadFile` and behave as before; the new `spec.LoadAt` takes the reader from its caller.
2. The status view's counts, in `internal/cli/follow_up_queue.go` (`followUpQueueRows`, the Spec queue) and `internal/cli/follow_up_status.go` (`followUpPendingSpecs`, the `Pending Specs` lines). Both call one new function, `followUpItemTaskGraph`, for a Delivery Queue item; a planned Spec keeps `spec.Load` on the Spec Root. The Daemon commits each Task's file with its settled status on the Item Branch (characterization, section 1), so the branch is the record of progress.
3. The Release Plan classifier, `internal/releaseplan`. `ClassifyRequest` gains a `Managed` value, the boundary test asks it about each path that is not planning evidence, and `Build` fills it from the Git source when the source implements the new `ManagedPathSource`, the same optional-interface seam `ReleaseGrantSource` already uses. `ClassifyCommit` keeps its signature and classifies as before.
4. The Release Plan's Git source, `releasePlanGitSource` in `internal/cli`, implements `ManagedPaths` in a new file. `roundfix release plan` and the `Next release` line of `roundfix follow-up --status` both build their plan through this source, so both get the rule.

The help of `release plan`, the release runbook, the usage guide and the Roundfix skill's release reference state the release rule; the `follow-up` guide states the count rule. The skill reference reaches adopters through the owned-skill installation.

## Implementation Design

### Interfaces

```go
// internal/spec/spec.go
// LoadAt behaves like Load but reads every file through read, which returns
// an error wrapping os.ErrNotExist for an absent path.
func LoadAt(specsRoot string, slug string, read func(string) ([]byte, error)) (*Graph, error)

// internal/cli/follow_up_queue.go
// followUpItemTaskGraph reads a Delivery Queue item's Task Graph from its
// Item Branch as Git holds it, and from the Spec Root when it cannot.
func followUpItemTaskGraph(ctx context.Context, git preflight.GitRunner,
	gitRoot, specRoot, archiveRoot string, item store.DeliveryQueueItem) (*spec.Graph, error)
```

```go
// internal/releaseplan/managed.go
type ManagedPathSource interface {
	ManagedPaths(context.Context, Range) (ManagedPaths, error)
}
type ManagedPaths struct {
	Adopter bool     // the target's skills-lock.json records the roundfix skill
	Files   []string // every path the Setup Manifest records, at the base or the target
	Skills  []string // every skill skills-lock.json records, at the base or the target
}
func (managed ManagedPaths) Contains(changedPath string) bool

// internal/cli/releaseplan_managed_paths.go
func (source releasePlanGitSource) ManagedPaths(ctx context.Context, releaseRange releaseplan.Range) (releaseplan.ManagedPaths, error)
```

Numbered invariants:

1. `followUpItemTaskGraph` runs only `git rev-parse` and `git show` in `gitRoot`, opens no file under the item's worktree and writes nothing.
2. Its order: an item whose `Branch` is empty or names no local branch reads the Spec Root; otherwise the folder at the branch head; otherwise, when `spec.ReadArchivedSpecAt` finds the Spec's Archive Record at the branch head under the Archive Root's repository-relative path, the folder the record's `source` names at its `source_revision`; otherwise the Spec Root. A Spec Root or Archive Root outside the repository reads the Spec Root.
3. `ManagedPaths.Contains` returns false for every path when `Adopter` is false, and for a path `cleanChangedPath` rejects.
4. `ManagedPaths` reads `skills-lock.json` and `docs/agents/setup-context.json` at the target and at the base with `git ls-tree` and `git show`, never from the working tree and never `.roundfixrc.yml`. An absent file, a file that is not JSON, or a lock whose `skills` is not an object contributes nothing; a failing `git` command fails the plan with a `GitSourceError`. `Adopter` comes from the target's lock only.
5. A commit with no changed path still crosses the boundary.

### Data Models

None. The Delivery Queue, the Archive Records, the Setup Manifest and `skills-lock.json` are read as they are; nothing is written.

### API Contracts

1. API Contract: a Delivery Queue item's Task counts. Each Delivery Queue row of the Spec queue and each queued line under `Pending Specs` takes its completed and total Tasks from `followUpItemTaskGraph`, in the order of invariant 2. A planned Spec that is not queued and every other row keep today's reading. When no source loads, the counts are 0 and 0, as today. The text and JSON shape do not change. The `follow-up` guide inserts, right after the sentence of the paragraph beginning "The status of a row is" that ends "`not yet a Spec` for a Backlog Entry.", this text, word for word: "A Delivery Queue item's completed and total Tasks, in this table and under `Pending Specs`, are counted on its Item Branch as Git holds it: the Spec folder at the branch head or, once the branch holds the Spec's Archive Record, the folder at the `source_revision` that record names. An item with no branch yet, or whose branch holds neither, is counted on the Spec Root, and no file of the item's worktree is read." The glossary entry for the Spec Queue Table gains, after its sentence that ends "or a later counted entry" and before its ADR links, the sentence "A Delivery Queue item's Tasks are counted on its Item Branch as Git holds it, at the branch head or, once the branch holds the Spec's Archive Record, at the record's source revision, and on the Spec Root when the branch does not exist yet or holds neither". The Skill Coverage Map replaces the `review` of the entry `command follow-up` with "2026-11-19 — Spec 0414 task_01: a Delivery Queue item's Tasks are counted on its Item Branch instead of the Spec Root; no command, flag, field, output line or exit code changes, and the Roundfix follow-up reference relays the rows as the command gives them", and the `review` of the entry `guide docs/user-guide/commands/follow-up.md` with "2026-11-19 — Spec 0414 task_01: the guide states where a Delivery Queue item's Tasks are counted; the Roundfix follow-up reference does not restate how counts are read and needs no change".
2. API Contract: the adopter maintenance boundary. In an adopter, `ManagedPaths.Contains` is true for exactly these paths: `docs/agents/setup-context.json`; each path in `Files`; `skills-lock.json`; for each name in `Skills` that is one path segment other than `.` and `..`, `.agents/skills/<name>`, `.claude/skills/<name>` and every path under either; `CLAUDE.md`; `docs/agents/specific-repository.md`; a path in `.claude/agents/` whose name starts with `roundfix-` and ends with `.md`; `.roundfixrc.yml`; and a path whose last segment is `CONTEXT.md`. A commit crosses the maintenance-only boundary when one of its changed paths is neither `MaintenanceOnly` nor contained. The JSON field `crossesMaintenanceOnlyBoundary` and the text output keep their shape. The `release plan` help gains, right after its paragraph that begins "Planning evidence (docs/specs", this paragraph, word for word apart from line wrapping: "In an adopter, a repository whose skills-lock.json at the target records the roundfix skill, a commit that changes only planning evidence and Adopter Maintenance Paths needs no release either. Those paths are the Setup Manifest docs/agents/setup-context.json and every path it records, skills-lock.json and every skill directory it records under .agents/skills/ and .claude/skills/, CLAUDE.md, docs/agents/specific-repository.md, a generated .claude/agents/roundfix-*.md, .roundfixrc.yml and every CONTEXT.md; the two records are read at the base and at the target. Roundfix's own repository authors the roundfix skill and records none, so its skills and Baseline still set a release." Both copies of the Roundfix skill's `references/release.md`, `docs/user-guide/release-runbook.md` and `docs/user-guide/usage.md` append to their paragraph that begins "Backlog Entries under `docs/backlog/`", after "a commit to classify by hand.", this text, word for word apart from line wrapping: "In an adopter, a repository whose `skills-lock.json` at the target records the `roundfix` skill, an Adopter Maintenance Path is maintenance-only too: the Setup Manifest `docs/agents/setup-context.json` and every path it records, `skills-lock.json` and every skill directory it records under `.agents/skills/` and `.claude/skills/`, `CLAUDE.md`, `docs/agents/specific-repository.md`, a generated `.claude/agents/roundfix-*.md`, `.roundfixrc.yml` and every `CONTEXT.md`, the two records read at the base and at the target. Roundfix's own repository authors the `roundfix` skill and records none in its lock, so its skills and Baseline still set a release." The glossary gains the entry Adopter Maintenance Path, placed right after the entry for the Release Plan, with this definition: "A path whose change ships nothing in an adopter, a repository whose `skills-lock.json` at the target revision records the `roundfix` skill: the Setup Manifest and every path it records, `skills-lock.json` and every skill directory it records under `.agents/skills/` and `.claude/skills/`, `CLAUDE.md`, `docs/agents/specific-repository.md`, a generated `.claude/agents/roundfix-*.md`, `.roundfixrc.yml` and every `CONTEXT.md`. The Release Plan counts a commit confined to these paths and to planning evidence as maintenance-only, reading the two records at the range's base and target. Roundfix's own repository authors the `roundfix` skill and records none, so it has no Adopter Maintenance Path.", followed by a link to ADR-0447 whose target is `docs/adr/0447-in-an-adopter-a-change-to-the-files-roundfix-manages-or-to-a-context-md-needs-no-release.md`, in the form the other entries give their ADR links, and the line "_Avoid_: Managed file, adopter-owned file, no-release path".

### Surface Transcripts

None. No command line, field, output line or exit code changes shape: API Contract 1 changes the numbers a row carries, and API Contract 2 changes which commits cross the boundary and adds one paragraph to a help text, which a contract test pins phrase by phrase. The QA gate reproduces both through the built binary (`_prd.md` → Success Metrics 1 and 4).

## Coverage Map

- Goal 1 → `spec.LoadAt`, `followUpItemTaskGraph` in `followUpQueueRows` and `followUpPendingSpecs` (API Contract 1).
- Goal 2 → `ManagedPaths.Contains`, the `ManagedPathSource` seam of `releaseplan.Build`, `releasePlanGitSource.ManagedPaths` (API Contract 2).
- Goal 3 → the adopter test of invariant 3 and API Contract 2, through `releasePlanGitSource.ManagedPaths`.
- Core Feature 1 → `spec.LoadAt`. Core Feature 2 → `followUpItemTaskGraph`. Core Feature 3 → `ManagedPaths`, `ManagedPathSource`, `releasePlanGitSource.ManagedPaths`.
- Core Feature 4 → the `follow-up` guide sentence, the Spec Queue Table sentence and the Coverage Reviews (API Contract 1); the help paragraph, the four release documents, the Adopter Maintenance Path entry and ADR-0447 (API Contract 2).
- Success Metric 1 → `TestLoadAtReadsThroughTheGivenReader`, `TestADeliveryItemCountsTasksOnItsItemBranch`.
- Success Metric 2 → `TestAnAdopterManagedPathIsMaintenanceOnly`, `TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget`.
- Success Metric 3 → `TestTheReleaseDocumentsNameTheAdopterMaintenancePaths`, `TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths`.
- Success Metric 4 → the QA outside-evidence row through `bin/roundfix release plan`.

Deviations from the PRD: None.

## Absence Evidence

1. `spec.LoadAt` — command: `rg -n 'func LoadAt|read func\(string\)' internal/spec/spec.go`; result: no match; `spec.Load` reads with `os.ReadFile` in `requireActive`, `loadManifestNodes` and `loadTask`.
2. `followUpItemTaskGraph` — command: `rg -n 'item\.Branch|ReadArchivedSpecAt|LoadAt' internal/cli/follow_up_queue.go internal/cli/follow_up_status.go`; result: no match; both counts call `spec.Load` on the Spec Root.
3. `ManagedPaths`, `ManagedPathSource` and `releasePlanGitSource.ManagedPaths` — command: `rg -n -g '!*_test.go' 'ManagedPath|skills-lock|setup-context|CONTEXT.md|roundfixrc' internal/releaseplan internal/cli/releaseplan_git_source.go`; result: no match; `isPlanningEvidencePath` lists six `docs/` prefixes only.
4. The adopter paragraph of the help and the release documents — command: `rg -n 'Adopter Maintenance|adopter' internal/cli/cli.go docs/user-guide/release-runbook.md docs/user-guide/usage.md .agents/skills/roundfix/references/release.md`; result: no match.
5. The new tests — command: `rg -n 'TestLoadAtReadsThroughTheGivenReader|TestADeliveryItemCountsTasksOnItsItemBranch|TestAnAdopterManagedPathIsMaintenanceOnly|TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget|TestTheReleaseDocumentsNameTheAdopterMaintenancePaths|TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths' internal skills`; result: no match.

## ADR Alternatives

- ADR-0447: the design takes its decision; a rule in every repository, a Project Config key, a Roundfix source path, a fixed path list and per-commit reads stay rejected.
- ADR-0412: extended, not amended; its rule stays a rule on paths, the plan still reads committed records only, and its rejected options (reading file contents at every commit, a list of derived files) stay rejected: the two records are read twice per range, and `release_paths` keeps answering release-preparation commits.
- ADR-0349: not applicable — the Release Ceiling's authority and the approval answer are unchanged; a classification by hand is still never answered by a ceiling.
- ADR-0359: not applicable — the owned-skill lock entry is only read, as the evidence that an adopter installed the Roundfix skill; no entry is written.
- ADR-0410: not applicable — which Specs make rows of the Spec queue is unchanged; only a queued row's counts change source.
- ADR-0434: not applicable — the renamed-slug rule is unchanged and reads no Task count.
- ADR-0426: not applicable — the Corrections Report keeps counting the plan the default branch held; the status view counts progress, a different question.
- ADR-0057: not applicable — the Daemon keeps exclusive ownership of Task status; the view reads the statuses the Daemon committed and writes none.
- ADR-0071: not applicable — no symbolic link or file mode is read; Git paths use forward slashes on every platform.

## Integration Points

None external. Both readers run the local `git` through the runners the status view and the Release Plan already hold (`preflight.GitRunner`). The skill reference reaches adopters through the owned-skill bundle.

## Testing Approach

- `TestLoadAtReadsThroughTheGivenReader` in a new `internal/spec/load_at_test.go`, with a map-backed reader: `same_graph_as_load`, `an_absent_prd_is_not_found` (a `SpecNotFoundError`), `an_absent_task_file_is_named` (a `MissingTaskFileError`). Every existing test of `./internal/spec` keeps passing.
- `TestADeliveryItemCountsTasksOnItsItemBranch` in `internal/cli/follow_up_queue_test.go`, on `newFollowUpReadInput`, `writeImplementSpecAtRoot`, `gitImplement` and a linked worktree for the Item Branch, with the subtests `a_running_item_counts_its_branch`, `an_archived_item_counts_the_source_revision`, `an_item_without_a_branch_counts_the_spec_root`, `a_missing_branch_counts_the_spec_root`, `a_branch_without_the_spec_counts_the_spec_root`, `the_worktree_is_not_read`, `the_pending_specs_line_counts_the_same` and `the_rows_write_nothing`. The prototype measured the first two failing on the unchanged tree.
- `TestAnAdopterManagedPathIsMaintenanceOnly` in a new `internal/releaseplan/managed_test.go`, one subtest per kind of path of API Contract 2, a negative subtest for each near miss (`an_unlocked_local_skill_crosses`, `another_claude_agent_crosses`, `context_map_crosses`, `an_unrecorded_guide_crosses`, `a_mixed_commit_crosses`), `not_an_adopter_changes_nothing` and `an_empty_commit_still_crosses`, through `ClassifyChanges`.
- `TestTheReleasePlanReadsTheManagedRecordsAtBaseAndTarget` in a new `internal/cli/releaseplan_managed_paths_test.go`, on temporary repositories with a stable tag of their own, through `runReleasePlanCommand --format json`: `the_fiscus_shape`, `roundfix_itself`, `a_skill_removed_in_the_range`, `a_malformed_lock_makes_no_adopter`, `project_config_is_not_read` (an invalid `.roundfixrc.yml` at the target).
- `TestTheReleaseDocumentsNameTheAdopterMaintenancePaths` in a new `skills/release_adopter_maintenance_test.go`, reading both copies of `roundfix/references/release.md` through `roundfixLayoutSources` and the two guides from the repository, one subtest per document and phrase; `TestTheReleasePlanHelpNamesTheAdopterMaintenancePaths` in `internal/cli/releaseplan_command_test.go` pins the help paragraph phrase by phrase. `TestReleasePlanHelpDescribesFlagsDefaultsStatesAndReadOnlyBoundary` keeps passing, so "reads no Roundfix configuration" stays.
- No test reads the clock, the network, the real home or another repository. Tags exist only in the temporary repositories the tests create, so no test depends on this repository's tags.

## Build Order

1. The status view counts a Delivery Queue item's Tasks on its Item Branch: `spec.LoadAt`, `followUpItemTaskGraph` in both counts, their tests, the `follow-up` guide sentence, the Spec Queue Table sentence, the Coverage Reviews and the re-recorded Behavior Surface Record (API Contract 1).
2. The adopter maintenance boundary: `ManagedPaths`, the `ManagedPathSource` seam, `releasePlanGitSource.ManagedPaths`, their tests, the help paragraph, the four release documents with the skill version and mirror, the Adopter Maintenance Path entry and the re-recorded Behavior Surface Record (API Contract 2) (depends on: 1, because both steps edit `CONTEXT.md` and `docs/references/behavior-surfaces.json`, and Tasks that share a file run in series).
3. The QA gate (depends on: 1, 2).

Files shared with other Specs: Spec 0411, delivered before this one, also changes `CONTEXT.md`, `docs/references/behavior-surfaces.json`, the Roundfix skill's `SKILL.md` in both copies and `skills/testdata/owned-skill-versions.json`; Spec 0409, delivered after, is expected to touch the same records. Step 2 raises the Roundfix skill to the next version with the record command on the tree it runs on, and each step regenerates the Behavior Surface Record there. No Verification of this Spec depends on another Spec.

## Risks & Considerations

- A Task settled but not yet committed, or a Task of a parallel Wave not yet integrated, is not counted until its commit reaches the Item Branch; the delay is one commit.
- Each queued item costs at most `2 + N` `git` processes per status call for a Spec of N Tasks, and one `git show` more for an archived item; the Delivery Queue rarely holds more than a few unmerged items.
- An adopter that installed no skills, or lost the `roundfix` lock entry, keeps the conservative classification; ADR-0447 names `roundfix baseline update` as the way back.
- A `feat:` or `fix:` commit confined to Adopter Maintenance Paths sets no impact, as one confined to planning evidence already does; a skill a repository writes for itself stays outside the set unless the lock records it.
- The help paragraph changes the `release plan` digest in the Behavior Surface Record; the release reference changes in the same Task, so the Skill Coverage Map reads the surface as described.

## Glossary

- adds: **Adopter Maintenance Path**
- changes: **Spec Queue Table**

## Decisions

- The status counts read the Item Branch through Git, never the item worktree, and fall back to the Spec Root where they read today; no ADR, because the source is an implementation choice behind an unchanged contract (`_prd.md` → Decisions).
- `spec.LoadAt` takes a reader rather than writing the folder to a temporary directory, so the status view keeps its promise to write no file.
- The Release Plan's adopter test and managed set follow ADR-0447.
- The managed set is read through an optional source interface, as the release grants are, so `releaseplan.ClassifyCommit` and every existing caller keep their behavior.
- No Surface Transcript: no output shape changes; contract tests pin the new help and document text.
