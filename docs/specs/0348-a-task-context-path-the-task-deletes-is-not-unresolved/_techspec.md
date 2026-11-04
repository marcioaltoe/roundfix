---
spec: 0348-a-task-context-path-the-task-deletes-is-not-unresolved
prd: _prd.md
created: 2026-11-04
---

# A Task Context path the Task deletes is not unresolved — Technical Spec

## Executive Summary

`detectTaskContextReferences` in `internal/speccheck/citations.go` reports
`SC-REF-UNRESOLVED` for an `interface:` or `instruction:` Task Context path
that is missing from the checked tree and not found through the archive
root. The fix adds one more way for such a path to resolve: the path is
tracked at `HEAD` (the Task's own uncommitted removal, which is what the
Daemon's settlement check sees) or at the Delivery Base (a removal committed
on the branch, which is what the QA gate precondition, the Delivery Queue and
a local `roundfix spec check` see). The base is read only when the checked
repository root is the top level of its Git work tree, and any Git failure
leaves the finding standing. The finding's fix stops suggesting "Create" and
tells the reader to remove the entry or declare the deletion; the summary
keeps its bytes so recorded settlement baselines keep matching. The
`write-tasks` skill and the guides state that a removed file is listed only
under `deletes:`.

The trade-off accepted: a path created and removed on the same branch stays
unresolved after its removal is committed, because neither `HEAD` nor the
Delivery Base tracks it; reading Task settlement commits would cover it, but
`deletes:` already does, and the fix names `deletes:`.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added; `git` reads the local repository. Source:
  `docs/agents/go.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0295: "One Default Branch Head
  serves every reachability and merge-base read"; the Task Context base
  resolves that head with `preflight.ResolveDefaultBranchHead`, takes the
  Delivery Base as its merge base with `HEAD`, never fetches and never
  writes a ref. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the
  cycle AX answers of 2026-11-04 quoted in the authorization record; the
  Governed Paths are the `write-tasks` and Roundfix skill files task_02
  changes; no Makefile, workflow, module or lint configuration changes.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0348-a-task-context-path-the-task-deletes-is-not-unresolved/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-tasks/SKILL.md`, `.agents/skills/write-tasks/references/task-template.md`, `skills/roundfix/SKILL.md`, `skills/write-tasks/SKILL.md`.

## Current behavior

Read on 2026-11-04 at `6a0ed904`; `_prd.md` → Validation records the runs.

1. `detectTaskContextReferences` skips `creates:` and `deletes:` entries.
   For every other entry it accepts a path that exists in the checked tree
   (`repositoryPathExists`), a path found through an archived Spec
   (`archivedTaskContextExists`), or an active Spec path whose relative path
   exists under the archive root. Otherwise it appends one error finding:
   code `SC-REF-UNRESOLVED`, summary
   `<task> declares unresolved Task Context path <path>`, locations at the
   Task's Context line and at `<path>:1`, and the fix
   `Create <path> or update the Task Context entry in <task>.`
2. The caller, `detectTaskCoverageAndContextReferences`, runs it for every
   Task, completed Tasks included.
3. The Daemon reads the same check twice per settled Task: before the
   Agent turn it records the refusing findings by `RefusalReason` (code and
   summary) in `specConsistencyBaseline`; the `settlement check: spec
   consistency` refuses each finding not in that baseline and prints its
   `Fix:` line into the diagnostics the repair turn reads
   (`internal/daemon/settlement_spec_consistency.go`). The Task's change is
   still uncommitted at that point. The QA gate precondition
   (`specGatePrecondition`) runs the check with no baseline.
4. So a Task that removes a file it lists under `interface:` fails its own
   settlement, and if it settled, the QA gate precondition would refuse the
   Spec for the same path.
5. The `write-tasks` skill says "every path it removes is declared under
   `deletes:`, never `instruction:`"; it does not name `interface:`. The
   spec command guide and the Roundfix `spec` reference describe only the
   archive-root resolution of a missing Task Context path.

## System Architecture

No package, command or flag is added. Changes by component:

- `internal/speccheck/citations.go`: `detectTaskContextReferences` consults
  the Task Context base before appending the finding, and takes its fix from
  one helper.
- `internal/speccheck/task_context_base.go` (new): the base reader and the
  fix text, both unexported.
- `internal/speccheck/task_context_base_test.go` (new): the case table.
- `internal/speccheck/archive_record_test.go`: one fixture line (Declared
  breaks).
- `internal/daemon/task_deletes_settlement_test.go`: the settlement cases.
- `skills/task_context_deletion_test.go` (new): the guide wording.
- The `write-tasks` skill and Task template, the Roundfix `spec` reference,
  the spec command guide and `CONTEXT.md`.

## Implementation Design

### Interfaces

```go
// internal/speccheck/task_context_base.go

// taskContextTrackedAtBase reports whether path is tracked at HEAD or at the
// Delivery Base of the Git work tree whose top level is repoRoot. Any Git
// failure, a repoRoot that is not that top level, or an unresolved Default
// Branch Head drops the revision it concerns; false keeps the finding.
func taskContextTrackedAtBase(repoRoot, path string) bool

// taskContextUnresolvedFix is the fix of the Task Context SC-REF-UNRESOLVED.
func taskContextUnresolvedFix(path, taskDisplayPath string) string
```

### Data Models

No schema changes.

### API Contracts

1. API Contract: the Task Context base. After the three existing
   resolutions fail, `detectTaskContextReferences` calls
   `taskContextTrackedAtBase(repoRoot, ref.Path)` and skips the entry when it
   returns true. The function:
   - runs `git -C <repoRoot> rev-parse --show-toplevel` and continues only
     when the result and `repoRoot` name the same directory after
     `filepath.EvalSymlinks` (a repository root nested below another work
     tree's top level, or outside Git, reads no base);
   - reads `HEAD`, and the Delivery Base: `preflight.ResolveDefaultBranchHead`
     with the current branch from `git symbolic-ref --quiet --short HEAD`,
     then `git merge-base <head ref> HEAD`; an undetermined default or a
     failed merge base drops only the Delivery Base;
   - returns true when `git cat-file -e <revision>:<path>` exits 0 for
     either revision, with the path in slash form; every `git` call disables
     the filesystem monitor (`-c core.fsmonitor=false`), as
     `preflight.ExecGitRunner` does;
   - runs only for an entry that is still unresolved, so a Spec whose paths
     all exist starts no `git` process.
2. API Contract: the fix. The Task Context finding's fix becomes
   `taskContextUnresolvedFix(ref.Path, taskDisplayPath)`, which returns
   `Remove the Task Context entry for <path> from <task>, or declare it under deletes: when the Task removes it; replace it only with a path the Task itself reads or changes.`
   Code, severity, summary and locations keep their bytes; the other
   `SC-REF-UNRESOLVED` findings (reference index rows, adopted sources) keep
   theirs.
3. API Contract: the rules in writing.
   - `.agents/skills/write-tasks/SKILL.md`, bullet "Declared edits and
     governed paths are explicit": "every path it removes is declared under
     `deletes:`, never `instruction:`" becomes "every path it removes is
     declared only under `deletes:`, never `interface:` and never
     `instruction:`", with "never `instruction:`" kept on one line, followed
     by the sentence "A Task that reads a file before removing it still lists
     it once, under `deletes:`." The
     `references/task-template.md` Context comment gains
     "A file the Task removes is listed only as deletes:, never as interface: or instruction:."
   - `.agents/skills/roundfix/references/spec.md` and
     `docs/user-guide/commands/spec.md`, the paragraph on a missing Task
     Context entry: an `interface:` or `instruction:` path missing from the
     checked tree resolves when `HEAD` or the Delivery Base tracks it and the
     checked root is the top level of its Git work tree; otherwise
     `SC-REF-UNRESOLVED` reports it with the fix of API Contract 2.
   - `CONTEXT.md`, **Delivery Base**: one sentence that the Spec
     Consistency Check resolves a Task Context path missing from the tree
     at it.

### Surface Transcripts

None. The changed finding text is asserted by named tests; no whole-command
transcript is declared.

### Invariants

```text
1. A Task Context path that exists in the checked tree, or resolves through the archive root, resolves as at 6a0ed904.
2. A missing interface: or instruction: path resolves through the base only when HEAD or the Delivery Base tracks it and repoRoot is its work tree's top level.
3. A path neither revision tracks reports one SC-REF-UNRESOLVED whose code, severity, summary and locations match 6a0ed904 byte for byte.
4. The fix of the Task Context finding never starts with "Create".
5. creates: and deletes: entries, SC-TOOLING-UNDECLARED, the mechanical authorization audit and the deletes: settlement refusal keep their behavior.
```

### Declared breaks

- A missing `interface:` or `instruction:` path that `HEAD` or the
  Delivery Base tracks no longer reports `SC-REF-UNRESOLVED`, in Runs and in
  a local `roundfix spec check` alike, including an uncommitted removal in an
  author's checkout.
- The Task Context fix text changes. No test or golden pins it:
  `grep -rn 'or update the Task Context entry' internal skills docs` finds
  only `internal/speccheck/citations.go`.
- `TestTaskContextResolvesThroughTheArchiveRecord`
  (`internal/speccheck/archive_record_test.go`), subtest `active`: its
  fixture commits the archived Spec's `_techspec.md` on `main` and removes it
  from the working tree without committing the removal, so with this design
  `HEAD` still tracks the path and the "unavailable source" and "missing
  record" assertions see no finding. The fixture commits the removal (`git
  rm` and a commit), as `roundfix archive` does; its three assertions keep
  their expectations, and the `archived` subtest is unchanged.

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → `taskContextTrackedAtBase`
  (API Contract 1), Invariants 1 and 2.
- Goal 2, Core Feature 1, Success Metric 2 → API Contract 1's top-level and
  revision rules, Invariant 3.
- Goal 3, Core Feature 2, Success Metric 3 → `taskContextUnresolvedFix`
  (API Contract 2), Invariant 4.
- Core Feature 3, Success Metric 4 → the settlement cases of
  `internal/daemon/task_deletes_settlement_test.go`, Invariant 5.
- Goal 4, Core Feature 4, Success Metric 5 → API Contract 3 and
  `TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes`.

Deviations from the PRD: None.

## Absence Evidence

1. `internal/speccheck/task_context_base.go` and a Git base read in the Task Context check — command: `grep -n 'cat-file\|merge-base\|ResolveDefaultBranchHead\|preflight' internal/speccheck/citations.go; ls internal/speccheck/task_context_base.go`; result: no match in `citations.go`, and no such file.
2. `internal/speccheck/task_context_base_test.go` — command: `ls internal/speccheck/task_context_base_test.go; grep -rn 'TrackedAtBase\|UnresolvedFix' internal`; result: no such file, no match.
3. `skills/task_context_deletion_test.go` and the removal rule in the guides — command: `ls skills/task_context_deletion_test.go; grep -n 'never .interface:. and never\|listed only as deletes:' .agents/skills/write-tasks/SKILL.md .agents/skills/write-tasks/references/task-template.md`; result: no such file, no match.

## ADR Alternatives

- ADR-0295: kept as recorded; it rejected `git remote set-head --auto` (a
  fetch and a ref write) and a CI environment variable as the source of the
  default branch; the Task Context base reads local refs through the same
  resolver, writes nothing and gives the same answer on every machine from
  the same refs.

## Integration Points

Git, through the `git` executable already required by the check; no remote
is contacted.

## Testing Approach

All tests are hermetic: temporary Git repositories through `gittest`
(`InitRepo`, `Run`), the default branch named `main` with
`refs/remotes/origin/main` written by `git update-ref`, and no network.

- `internal/speccheck/task_context_base_test.go`:
  - `TestTaskContextPathTrackedAtTheBaseResolves`, a table over a pending
    Task with one Context entry, built on the `newUndeclaredFixture`
    helpers plus a minimal `_prd.md`, run through `speccheck.Check`. Rows,
    with today's count first (`_prd.md` → Validation) and the count after:
    `present` 0 → 0; `working tree removal` 1 → 0; `staged removal` 1 → 0;
    `committed removal` 1 → 0; `instruction committed removal` 1 → 0;
    `branch file removed in the working tree` 1 → 0; `never tracked` 1 → 1;
    `created and removed on the branch` 1 → 1; `removed on main before the
    branch` 1 → 1; `no Git repository` 1 → 1; `nested below another work
    tree` 1 → 1 (the outer repository's `HEAD` tracks the path at its own
    root; the checked root is a plain subdirectory); `no default branch
    committed removal` 1 → 1 (work branch only, no `main`, no origin ref).
    Each reporting row asserts code, severity, the exact summary, both
    locations and the exact fix of API Contract 2.
  - Every new top-level test calls `t.Parallel()` first.
- `internal/daemon/task_deletes_settlement_test.go`:
  - `TestAGatedTaskWhoseInterfacePathItDeletesSettles`: the fixture of
    `retiredHistoryDeletionFixture` with the Context entry labeled
    `interface:` (committed before the turn), for the modes `disk-only` and
    `staged`; `TaskCycle` completes one Task. At `6a0ed904` it fails with
    `SC-REF-UNRESOLVED`.
  - `TestSettlementStillRefusesAContextPathNoRevisionTracks`: the turn
    appends `- interface: \`packages/backend/test/never.test.ts\`` to the
    Task's Context; the Task fails, and its spec consistency diagnostics
    file carries `SC-REF-UNRESOLVED`, the path and the fix of API Contract 2.
- `skills/task_context_deletion_test.go`,
  `TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes`: both copies of
  the `write-tasks` `SKILL.md` and of `references/task-template.md` contain
  the API Contract 3 sentences on whitespace-normalized text, and neither
  `SKILL.md` copy still contains "declared under `deletes:`, never
  `instruction:`".
- `TestWriteTasksSkillStatesTheDeclaredPathRules`
  (`skills/baseline_skill_contract_test.go`) stays unchanged: it requires
  the raw text "never `instruction:`" in the canonical `write-tasks` skill,
  which is why API Contract 3 keeps that phrase on one line.
- Existing suites: `./internal/speccheck`, `./internal/daemon`,
  `./internal/spec`, `./internal/cli`, `./skills`, and
  `-tags docscontract ./internal/docscontract`.
- Sabotages, each run by the author on a scratch copy at `6a0ed904` with the
  prototype, with the row whose input makes it observable:
  - the base reader removed (today's behavior):
    `TestTaskContextPathTrackedAtTheBaseResolves` fails on
    `working tree removal` and `committed removal`;
  - the Delivery Base dropped (HEAD only): it fails on `committed removal`;
  - `HEAD` dropped (Delivery Base only): it fails on
    `branch file removed in the working tree`;
  - the top-level guard removed: it fails on `nested below another work
    tree`;
  - the whole history searched instead (`git log --all -- <path>` non-empty):
    it fails on `removed on main before the branch` and `created and removed
    on the branch`;
  - the old fix text restored: it fails on `never tracked`, and
    `TestSettlementStillRefusesAContextPathNoRevisionTracks` fails;
  - the `write-tasks` sentence removed from the mirror:
    `TestTheWriteTasksGuidesListARemovedFileOnlyUnderDeletes` fails.

## Build Order

1. The Task Context base, the fix text, the case table, the archive record
   fixture change and the two settlement tests (task_01).
2. The `write-tasks` skill and Task template, the Roundfix `spec`
   reference, the spec command guide, `CONTEXT.md`, the skills test, with
   the mirrors synced, the owned-skill versions recorded and the Behavior
   Surface Record re-recorded last (task_02). It needs no step of task_01:
   its test reads only guide text, and the two Tasks share no file.
3. QA gate (task_03) (depends on: 1, 2).

No other Spec's code is needed; see `_prd.md` → Prerequisites for the files
shared with Specs 0345, 0346 and 0347.

## Risks & Considerations

- A path created and removed on one branch and still listed under
  `interface:` reports the finding at the QA gate precondition after its
  removal is committed. The fix names `deletes:`, and the guides tell the
  author to use it.
- In a repository without a resolvable Default Branch Head (no `origin`
  default and a work branch not named `main` or `master`), a committed
  removal reports the finding as today; Run Worktrees share the user
  checkout's refs, so a Run of a cloned repository resolves it.
- The base adds at most four `git` processes per unresolved entry and none
  for a Spec whose paths exist.

## Glossary

- changes: **Delivery Base** — the Spec Consistency Check also reads a Task Context path at it

## Decisions

- Base revisions: `HEAD` and the Delivery Base, read only from local refs.
  The Run start head was not chosen, because `roundfix spec check`, the QA
  gate and the Delivery Queue cannot all name it, while all of them resolve
  the Delivery Base the same way, through one Default Branch Head resolver.
- Settlement commit trailers were not read: they would cover a file created
  and removed on one branch, at the cost of a history walk per entry, and
  `deletes:` already covers that case.
- A Git failure keeps the finding (fail closed); the check never errors
  because the base is unreadable.
- The summary is unchanged so that settlement baselines and refusal reasons
  recorded before the change keep matching; only the fix changes.
- No new ADR; the rule lives in the guides and `CONTEXT.md`.
- Secondbrain: `wiki/index.md` read; the `qmd query` named in `_prd.md` →
  Validation returned this repository's own mirrored documents
  (an ADR on undeclared Verification paths, the spec and settle command
  guides, `CONTEXT.md`), the mirrored Roundfix skill references and
  `wiki/entities/roundfix.md`; none decides how a check treats a deleted
  Task Context path, so the design rests on the code reading and the Git
  manual. Exa: the Git manual passages in `_prd.md` → Acceptance
  evidence.
- Two implementation Tasks in parallel and a QA gate.
