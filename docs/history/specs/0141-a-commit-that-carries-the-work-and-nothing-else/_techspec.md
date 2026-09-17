---
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
prd: _prd.md
created: 2026-09-17
---

# A commit that carries the work and nothing else — Technical Spec

## Executive Summary

Three changes at one boundary. The stageable filter asks Git whether a path is
tracked before it refuses an executable file. The commit step stops treating a
refusal as advisory: a path the work produced and the Daemon could not stage
fails the Task, with the paths and cause in the outcome. And the QA step takes a
second snapshot around the repository Verification, so the QA Report commit
excludes what the verifier wrote.

The trade-off this design accepts is that a refusal now has two classes. A path
that is absent, or that lives outside the repository, is refused by design and
stays advisory; a path that exists in the worktree and was not staged is a lost
output and fails the Task. Collapsing the two would either break the supported
external Spec Root or keep the silence this Spec exists to end.

## Project Constraints

- Identifier strategy: applicable — the drop reasons, the Run and Task outcome vocabulary and the Run Event kinds keep their spelling and meaning; this Spec adds no identifier and renames none. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the commit boundary reads and writes the local worktree and index only, and opens no credential store or network transport. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the settlement path is governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0057 applies: the Daemon owns Task status, so a dropped path changes the status the Daemon writes and never asks an Agent to record it.
  ADR-0096 applies: the QA gate proves machine facts before it spends an Agent turn, so the commit baseline moves without changing what the mechanical stage decides.
  ADR-0117 applies: a defect is checked by the stage that can produce it, so the drop is refused at the commit boundary that creates it.
  ADR-0035 applies: a Spec Root outside the repository keeps its artifacts uncommitted, so a path external to the repository is refused by design and is not a dropped work output.
  ADR-0036 applies: review artifacts are committed in their own docs commit, which this Spec neither creates nor changes.
  ADR-0029 applies: review artifacts live with the Spec, and this Spec moves no artifact between locations.
  ADR-0142 applies: head-bound Review Source Evidence decides the watch outcome, and this Spec changes no watch, review or evidence binding.
  ADR-0080 applies: a QA row that no environment can run is recorded as environment-blocked and never as a product failure.
  ADR-0091 applies: the QA gate is a Task node of its own type, which is the terminal Task this Spec authors.
  ADR-0093 applies: Spec consistency is checked by citation and never by inference, and this Spec changes no consistency rule.
  ADR-0097 applies: a QA row carries forward only on declared unmoved evidence, and this Spec moves no row's evidence.
  ADR-0104 applies: a Spec accepts on evidence it did not author, which here is this repository's three tracked executable files.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0157 applies: tracked source is work output whatever its mode, and a drop stops the Task from settling completed.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/daemon`, which is not a Governed Path, and its test files are ordinary too. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Stageable filter | `internal/daemon` commit boundary | Ask the index whether a path is tracked before refusing it for its mode. |
| Drop classification | `internal/daemon` commit boundary | Separate a refusal by design from a lost output, and carry the lost ones to the Task result. |
| Task settlement | `internal/daemon` task cycle | Refuse to settle a Task completed when its commit lost an output, naming paths and cause. |
| QA commit baseline | `internal/daemon` QA step | Take a second snapshot around the repository Verification and exclude what appeared inside that window. |

No new package, file or directory is proposed. Every seam exists and is covered
by tests today.

## Implementation Design

### A tracked path is stageable

The filter keeps its order — external, symlink, mode, absence — and the mode
refusal gains one question: does the index already track this path? A tracked
path is kept, whatever its mode. An untracked executable is refused exactly as
today, with the same reason and mode in its event.

The index query uses the same shape the filter already uses to decide absence,
so the boundary gains no new dependency and no new failure mode. A query that
cannot be answered is treated as untracked, which preserves today's behavior
rather than inventing a permissive default.

### A refusal has two classes

- **Refused by design.** A path external to the repository, and a path absent
  from both worktree and index. The first is what a Spec Root outside the
  repository produces; the second names a path that no longer exists. Both keep
  today's console line and Run Event, and neither fails the Task.
- **Lost output.** A path that exists in the worktree and was not staged: an
  untracked executable, and a path that crosses a symbolic link. The Task does
  not settle completed, and the result carries every lost path with its reason.

The classification lives beside the filter, so a future reason declares its
class where it is written rather than in a caller.

### The QA Report commit excludes the Verification's writes

The QA step already snapshots the worktree before it begins. It now takes a
second snapshot immediately before the repository Verification starts and a
third immediately after it ends. Paths that appear between those two are what
the Verification wrote, and the report commit's path set subtracts them.

Everything else is unchanged: the report the mechanical stage seeded is already
dirty at the second snapshot and stays in the commit, and whatever the QA Agent
writes afterwards appears after the third snapshot and stays in the commit. A
Spec whose Verification writes nothing sees no difference at all.

### Interfaces

The filter's signature keeps its shape; its result gains the class of each
dropped path.

```go
// DroppedStagePath gains the class that decides whether the Task can settle.
type DroppedStagePath struct {
    Path   string
    Reason string
    Mode   string
    Lost   bool // the path exists in the worktree and was not staged
}
```

### Data Models

No entity, schema, stored record or Run Event payload shape changes. The Task
result carries the lost paths it already collects as dropped paths.

### API Contracts

1. A Task whose commit lost an output does not settle `completed`. Its Task
   result and the Run outcome name each lost path and its reason, in the same
   vocabulary the existing console line uses.
2. `FilterStageablePaths` keeps a tracked path whatever its mode, and refuses an
   untracked executable with the reason and mode it reports today.
3. The QA Report commit carries the report, its evidence and the Agent's writes,
   and never a path that appeared while the repository Verification ran.

## Coverage Map

- Goal 1 → Stageable filter.
- Goal 2 → Drop classification, Task settlement.
- Goal 3 → QA commit baseline.
- Goal 4 → Drop classification (refused-by-design class).
- User Story 1 → Stageable filter.
- User Story 2 → Task settlement; API Contract 1.
- User Story 3 → Task settlement.
- User Story 4 → QA commit baseline; API Contract 3.
- Core Feature 1 → Stageable filter; API Contract 2.
- Core Feature 2 → Drop classification; Task settlement.
- Core Feature 3 → QA commit baseline.
- Core Feature 4 → Drop classification.
- Success Metric 1 → Testing Approach 1 and 4.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → Task settlement, Stageable filter, QA commit baseline.

## Integration Points

- **External Spec Root.** A Spec Root outside the repository produces paths
  refused as external. They stay in the refused-by-design class, so that
  supported configuration keeps settling as it does today.
- **QA gate.** The mechanical stage, the withholding rule and the verdict
  settlement are untouched; only the commit's path set moves.
- **Spec 0122.** This Spec is the first slice carved from it. Running the
  Verification against the integrated candidate, the acceptance eligibility
  policy, the archive override, the known-red repair and independent
  Verification groups stay there.

## Testing Approach

1. **Filter, at its existing unit seam.** Table-driven cases over a temporary
   repository: a tracked executable is kept; an untracked executable is refused
   with its mode; a tracked non-executable is unchanged; an unreadable index
   answer refuses. The tracked case fails on the tree as it stands today.
2. **Settlement, at the task-cycle seam.** A fixture Task whose work produces an
   untracked executable does not settle completed, and its result names the path
   and reason. A fixture whose only refusal is an external path still settles.
3. **QA baseline, at the QA-step seam.** A fixture whose configured repository
   Verification writes a tracked file produces a report commit without it, and a
   fixture whose Verification writes nothing produces the same commit as today.
4. **Outside evidence.** The replay uses the three executable files this
   repository already tracks, which this Spec did not author: a fixture Run
   edits one of them and the settlement commit carries the change, where the
   unchanged tree drops it. Where the file cannot be read, the row records that
   reason and does not block.
5. **Repository gate.** The terminal QA Task records the Daemon's `make verify`
   result as a fact.

## Build Order

1. Tracked-path question in the stageable filter, with its unit tests (depends
   on: none).
2. Drop classification and the Task settlement that reads it, with task-cycle
   tests (depends on: 1).
3. QA commit baseline around the repository Verification, with QA-step tests
   (depends on: none).
4. Terminal QA (depends on: 1, 2, 3).

Steps 1-2 and step 3 touch different parts of the same package and may run in
the same wave only if the graph proves they change different files; the Task
Graph decides that.

## Risks & Considerations

- **A drop that used to be tolerated now fails a Run.** That is the intended
  break. The refused-by-design class keeps the supported external Spec Root and
  absent-path cases settling as before, so the blast radius is untracked
  executables and symlink crossings.
- **The index query at the commit boundary.** It runs once per candidate path in
  a set that is already bounded by the changed-path reader, and reuses the
  existing query shape, so it adds no new failure mode. An unanswerable query
  refuses, which is today's behavior.
- **Three snapshots in the QA step.** Each is the same operation the step
  already performs once. A Spec whose Verification writes nothing pays two
  extra snapshots and changes no commit.

## Decisions

- **Tracked beats mode, and a drop fails the Task.** See ADR-0157.
- **Two classes of refusal.** Collapsing them would break the external Spec Root
  that ADR-0035 supports, or keep the silence this Spec ends.
- **A symlink crossing is a lost output.** The path exists and the work produced
  it; refusing to stage it loses work even though the refusal is safe.
- **Subtract the Verification window rather than re-snapshot.** Re-snapshotting
  after the Verification would drop the report the mechanical stage seeded,
  which must ride in the commit.

## Vocabulary Contract

No token is coined. "Lost output" describes an existing dropped path in prose
and names no new field, state, code or command.
