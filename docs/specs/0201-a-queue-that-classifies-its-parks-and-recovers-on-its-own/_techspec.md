---
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
prd: _prd.md
created: 2026-09-30
---

# A queue that classifies its parks and recovers on its own — Technical Spec

## Executive Summary

Every change is a narrow extension of the Delivery Engine in
`internal/delivery` and of its command workflow in `internal/cli`. The engine
gains four optional dependencies: a prerequisite reader, a conflict resolver,
a failed-check inspector and an item history reader. When one is nil the
engine behaves exactly as today, so the existing engine tests and fakes keep
their meaning. A pure Park Class table turns every blocker into a class and a
next command, and `deliver status` and the Pending Question both read it.

The primary trade-off is in conflict recovery. The owner merges the default
branch and regenerates derived paths without a second Pre-PR Review, and it
trusts a Project Config declaration of which paths are derived. The accepted
risk is a declared path that also carries hand-edited source; the repository
gate that follows is the net, and nothing is pushed before it passes
(ADR-0192). The alternative, parking every conflict, keeps the manual
intervention this Spec exists to remove.

## Project Constraints

- Identifier strategy: applicable — the four new blockers, the nine Park
  Classes, the manifest key `requires` and the Project Config key
  `delivery.derived_paths` follow the existing kebab-case and snake_case
  forms; no generated identifier is added. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Git, the Run Database and the
  existing `gh` session only. The added `gh` calls are `pr view --json
  mergeable`, `run view <id> --json attempt`, `run view <id> --log-failed` and
  `run rerun <id> --failed`, on the repository the queue already publishes to.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0192 and ADR-0193 (this Spec),
  ADR-0149, ADR-0178, ADR-0154, ADR-0165, ADR-0170, ADR-0053, ADR-0158,
  ADR-0164, ADR-0161, ADR-0090, ADR-0153, ADR-0169, ADR-0174, ADR-0196,
  ADR-0197, ADR-0160, ADR-0081, ADR-0166, ADR-0167, ADR-0176 and ADR-0179 hold
  as the PRD states. The gate is bound by ADR-0080, ADR-0088, ADR-0091,
  ADR-0096, ADR-0104, ADR-0117, ADR-0155 and ADR-0156, and by ADR-0093 and
  ADR-0094. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the maintainer's 2026-09-30 statement that
  every skill may be adjusted covers the Roundfix skill's delivery reference
  and version, recorded in [_authorization.md](_authorization.md); bounded
  files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`. Sanctioned regeneration: `make skills-sync`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## System Architecture

| Concern | Reader or actor | File |
| --- | --- | --- |
| Park Classes | `ClassifyPark`, used by `PendingQuestionFor` and `deliver status` | `internal/delivery/park_class.go` (new), `internal/delivery/question.go`, `internal/cli/deliver.go` |
| Failed check re-run | `CheckRecovery`, implemented by `GitHubCLI` | `internal/delivery/check_rerun.go` (new), `internal/delivery/engine.go` |
| Prerequisites | `requires` in the Task Graph manifest; `PrerequisiteReader` | `internal/spec/spec.go`, `internal/delivery/engine.go`, `internal/cli/deliver_workflow.go`, `internal/cli/deliver.go`, `internal/store/delivery.go` |
| Environment-only partial | `RunSpec`, `InspectItem`, `Archive`, `Authorization`; `ItemHistory` | `internal/cli/deliver_workflow.go`, `internal/delivery/engine.go` |
| Conflicts | `CheckReport.Mergeable`; `ConflictResolver`; `delivery.derived_paths` | `internal/delivery/github.go`, `internal/delivery/engine.go`, `internal/cli/deliver_workflow.go`, `internal/config/config.go` |

No package or command is added. Every new owner action writes one line to the
owner's console log, in the `roundfix: <action>: Delivery Queue item <slug>: …`
form the premise warning already uses.

## Implementation Design

### Interfaces

```go
// internal/delivery/engine.go — optional; nil keeps today's behavior.
type PrerequisiteReader interface {
	UnmetPrerequisites(ctx context.Context, gitRoot, specSlug string) ([]string, error)
}
type ConflictResolver interface {
	ResolveConflict(ctx context.Context, workDir, specSlug, head string) (ConflictResolution, error)
}
type CheckRecovery interface {
	InspectFailedCheck(ctx context.Context, workDir, head string, check PullRequestCheck) (CheckFailure, error)
	RerunFailedCheck(ctx context.Context, workDir string, failure CheckFailure) error
}
type ItemHistory interface {
	Descends(ctx context.Context, workDir, ancestor, head string) (bool, error)
	RunStart(ctx context.Context, gitRoot, runID string) (string, error)
}
type ConflictResolution struct{ Head string; SourcePaths, Regenerated []string }
type CheckFailure struct {
	RunID string; Attempt int; Packages []string; OutsideChange bool
}
```

```go
// internal/delivery/park_class.go
type ParkClass string // dependency, conflict, environment, flaky-check,
                      // finding, budget, review, authorization, unclassified
type ParkClassification struct{ Class ParkClass; Next string }
func ClassifyPark(queue store.DeliveryQueue, item store.DeliveryQueueItem) ParkClassification
```

Existing types gain fields: `CheckReport.Mergeable string`,
`RunResult.QAEnvironmentPartial bool`, `ItemState.QAOverride bool`,
`ArchiveResult.AlreadyArchived bool`, and `spec.Graph.Requires []string`.
`EngineDependencies` gains `Prerequisites`, `Conflicts`, `Checks` and
`History`. `newCommandDeliveryEngine` wires the workflow for the first, second
and fourth, and `delivery.NewGitHubCLI` for `Checks`.

### Data Models

No Run Database schema change. Blockers stay one string:

| Blocker | Written when |
| --- | --- |
| `prerequisite-unmerged: <slug>, …` | a queued item has an unmet prerequisite that is parked or not queued |
| `pull-request-conflict: <path>, …` | GitHub reports `CONFLICTING` and a conflicted path matches no declaration |
| `qa-environment-partial` | a Run ends unresolved on an environment-only partial QA Report |
| `flaky-check: <package>, …` | a re-run check fails again outside the item's change |

The retry transition in `internal/store/delivery.go` also accepts `queued` as
a target stage. Project Config gains:

```yaml
delivery:
  derived_paths:
    - paths: [internal/baseline/testdata/catalog.digest, internal/baseline/testdata/plan-characterization/]
      regenerate: make baseline-digests
```

A path entry is a clean repository-relative file, a directory ending in `/`
that covers everything below it, or a `path.Match` pattern over the whole
path. An empty list, an absolute or `..` entry, or an empty `regenerate` is a
config error that names the key. Project Config overrides User Config as for
every other key.

### Park Classes

`ClassifyPark` matches a blocker exactly or by its `<blocker>:` prefix:

| Class | Blockers | Next command |
| --- | --- | --- |
| `dependency` | `prerequisite-unmerged` | a prerequisite parked in the queue: `run roundfix deliver retry <prerequisite>; <slug> returns to the queue once <prerequisite> is merged`; otherwise `deliver or merge <prerequisite>, then run roundfix deliver retry <slug>` |
| `conflict` | `pull-request-conflict` | `merge the default branch into the item branch in <worktree>, resolve <paths>, commit, then run roundfix deliver retry <slug>` |
| `environment` | `qa-environment-partial`, `checks-timeout`, `item-worktree-missing`, `delivery-error` | for `qa-environment-partial`: `(cd <worktree> && roundfix reconcile <run-id> --carry-forward), satisfy the environment-blocked QA rows, run roundfix archive <slug> --qa-override --approval <source> --reason <text>, then run roundfix deliver retry <slug>`; otherwise the generic answer |
| `flaky-check` | `flaky-check` | `the check failed twice in <packages>, which <slug> did not change; fix or re-run it, then run roundfix deliver retry <slug>` |
| `finding` | `run-unresolved`, `review-findings`, `corrective-spec-required`, `gate-failed`, `checks-failed`, `revalidation-failed` | today's answer |
| `budget` | `run-budget-exceeded`, `queue-deadline` | today's answer |
| `review` | `review-blocked`, `review-stale` | today's answer |
| `authorization` | `unauthorized` | today's answer |
| `unclassified` | any other blocker | the generic answer |

The generic answer is today's `resolve the blocker, then run roundfix deliver
retry <slug>`. Every existing blocker keeps today's answer, so the Pending
Question changes only for the four new blockers. `deliver status` prints
`Park: <slug> <class>: <next>` for each parked item, in queue order, after the
Warning lines and before the Limits line.

### Failed check re-run

In `checkCandidate`, a check in bucket `fail` goes to `Checks` before the item
parks. `GitHubCLI.InspectFailedCheck`:

- takes the run ID from the check's `link`
  (`…/actions/runs/<id>/job/<job>`); without one, it returns an empty failure;
- reads `gh run view <id> --json attempt` and `gh run view <id> --log-failed`;
- collects each package from a Go test summary line `FAIL\t<import path>\t<n>s`.
  A `[build failed]` or `[setup failed]` line makes the failure unattributable;
- maps each import path to a directory through the `module` line of the item
  worktree's `go.mod`, and lists the item's changed paths with `git diff
  --name-only <merge-base> <head>` against the refreshed remote default branch;
- sets `OutsideChange` only when at least one package was found and no changed
  path lies under any package directory.

When `OutsideChange` holds and `Attempt` is `1`, the engine calls
`RerunFailedCheck` (`gh run rerun <id> --failed`), logs it, restarts the check
timeout once and keeps polling. A later pass that finds the check green joins
`flaky-check: <check> passed on re-run` to the item's Warning. A second failure
of a re-run check parks `flaky-check: <packages>`. Every other failure parks
`checks-failed` as today. The engine never re-runs a check twice for one head.

### Prerequisites

`spec.Load` reads `requires` from the manifest frontmatter. It must be a list
of distinct, non-empty slugs other than the Spec's own; anything else is a load
error naming `_tasks.md`. `deliver start` resolves every entry against the
Specs Root and its archive root, and refuses an unknown slug or a cycle among
the queued Specs, after every Spec loads and before the delivery authorization
check and the queue record (Surface Transcript 2).

`commandDeliveryWorkflow.UnmetPrerequisites` fetches the default branch from
the delivery remote, reads the item's `_tasks.md` at the refreshed
remote-tracking ref, and returns each prerequisite whose archived `_prd.md` is
absent at that ref. The engine calls it for a `queued` item before
`CreateItemBranch`:

- no unmet prerequisite: the item starts as today;
- every unmet prerequisite is a queue item that is neither parked nor merged:
  the item stays `queued`, is logged as waiting, and the pass moves on;
- otherwise the item parks `prerequisite-unmerged: <unmet>` with no worktree.

At the start of every pass, and after any item merges, the engine returns each
`prerequisite-unmerged` item whose prerequisites are all met to `queued`,
clears its blocker and leaves its retry count unchanged. `Retry` of such an
item skips the worktree, the revalidation and the carry-forward, and targets
`queued`.

### Environment-only partial and the operator-archived retry

After an unresolved Run, `RunSpec` reads the newest QA Report of the Spec at
the tip of the Run Branch (`roundfix/run-<run-id>`). It sets
`QAEnvironmentPartial` when the verdict is `partial`, `rows_blocked_finding` is
`0` and `rows_blocked_environment` exceeds the pre-PR Pull Request rows. The
engine then parks `qa-environment-partial`. A missing or unreadable report
keeps `run-unresolved`.

`InspectItem` sets `QAOverride` from the archived `_prd.md` frontmatter
(`qa_override: true`). In `Retry`, an archived item whose head differs from
its candidate is accepted when all of these hold, and refused with today's text
otherwise:

- the blocker is `qa-environment-partial` and `QAOverride` is true;
- `History` is set, and the head descends from the anchor: the last candidate
  commit, or the Run start of `item.RunID` when no candidate was recorded.

The accepted item appends the head to its candidate commits and resumes at
`reviewing`. `Archive` returns `AlreadyArchived` when the reviewed head holds
the archived Spec and not the active one, and the engine then moves to
`gating` without a new commit. `Authorization` reads the grant at the parent of
the newest first-parent commit that deleted the active `_prd.md`, instead of
`HEAD^`.

### Conflicts

`GitHubCLI` adds `mergeable` to the Pull Request fields and returns it in
`CheckReport`. In `checkCandidate`, `CONFLICTING` goes to `Conflicts` on the
first read; `UNKNOWN` counts as pending. `ResolveConflict` in the item
worktree:

1. requires a clean worktree whose `HEAD` is the candidate head;
2. fetches the default branch and runs `git merge --no-ff --no-commit
   <remote>/<default>`;
3. lists conflicted paths (`git diff --name-only --diff-filter=U`); when any
   matches no declaration, runs `git merge --abort` and returns them as
   `SourcePaths`;
4. checks out `--theirs` for each conflicted path, runs each matched
   declaration's `regenerate` command once in declaration order through the
   Verification executor, with its log under the artifact directory;
5. aborts and returns `regenerated <path> outside delivery.derived_paths` when
   the regeneration changed an undeclared path, and otherwise stages the
   declared paths and commits the merge with the trailer
   `Roundfix-Delivery: derived-merge`.

On `SourcePaths` the engine parks `pull-request-conflict: <paths>`. On a new
head it appends the head and moves to `gating`, so the existing gate, push and
check stages run again. `createPullRequest` re-reads an existing Pull Request
whose head is an earlier candidate commit of the item, every check interval up
to the check timeout, before it fails. A `Retry` of `pull-request-conflict`
accepts a head that descends from the candidate, like the operator-archived
retry, and resumes at `reviewing`.

### API Contracts

1. API Contract: `roundfix deliver status` — prints one `Park:` line per
   parked item with its Park Class and next command (Surface Transcript 1). A
   queue without a parked item prints exactly today's output.
2. API Contract: `roundfix deliver start <slug>...` — exits `2` before
   recording anything when a `requires` entry is unknown, names its own Spec,
   or closes a cycle among the queued Specs (Surface Transcript 2).
3. API Contract: `roundfix deliver retry <slug>` — returns a
   `prerequisite-unmerged` item to `queued`, and resumes an operator-archived
   `qa-environment-partial` item or an operator-resolved
   `pull-request-conflict` item at `reviewing`. Every other refusal keeps its
   text.
4. API Contract: Project Config `delivery.derived_paths` — a list of `paths`
   and `regenerate` pairs, validated at load.

### Surface Transcripts

1. Surface Transcript: `deliver status` with a conflict park and a dependency
   park behind it.

   ```transcript
   $ roundfix deliver status
   stdout:
   0204-first	parked	pull-request-conflict: internal/baseline/assets/profiles/standard-typescript-monorepo.json	/worktrees/0204-first
   0205-second	parked	prerequisite-unmerged: 0204-first	-
   Park: 0204-first conflict: merge the default branch into the item branch in /worktrees/0204-first, resolve internal/baseline/assets/profiles/standard-typescript-monorepo.json, commit, then run roundfix deliver retry 0204-first
   Park: 0205-second dependency: run roundfix deliver retry 0204-first; 0205-second returns to the queue once 0204-first is merged
   Limits: <limits>
   Pending question: 0204-first parked pull-request-conflict: internal/baseline/assets/profiles/standard-typescript-monorepo.json
   Answer: merge the default branch into the item branch in /worktrees/0204-first, resolve internal/baseline/assets/profiles/standard-typescript-monorepo.json, commit, then run roundfix deliver retry 0204-first
   Waiting behind it: 1 parked item(s)
   stderr:
   exit: 0
   ```

2. Surface Transcript: `deliver start` refuses two Specs that require each
   other.

   ```transcript
   $ roundfix deliver start 0204-first 0205-second
   stdout:
   stderr:
   Preflight failed

   Reason:
     Delivery Queue Specs require each other in a cycle: 0204-first -> 0205-second -> 0204-first

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix deliver start --help' for usage.
   exit: 2
   ```

## Coverage Map

- Goal 1 → Prerequisites; API Contracts 2, 3.
- Goal 2 → Conflicts; API Contracts 3, 4.
- Goal 3 → Environment-only partial and the operator-archived retry; API
  Contract 3.
- Goal 4 → Park Classes; Failed check re-run; API Contract 1.
- Story 1 → Prerequisites.
- Story 2 → Conflicts; Park Classes.
- Story 3 → Conflicts.
- Story 4 → Environment-only partial and the operator-archived retry.
- Story 5 → Park Classes.
- Story 6 → Failed check re-run.
- Core Feature 1 → Prerequisites.
- Core Feature 2 → Conflicts.
- Core Feature 3 → Environment-only partial and the operator-archived retry.
- Core Feature 4 → Park Classes.
- Core Feature 5 → Failed check re-run.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 4.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 1.
- Success Metric 5 → Testing Approach 1.

## Integration Points

- **GitHub through `gh`.** The same authenticated CLI and repository; every
  call goes through the `CommandRunner` seam, so tests use a fake runner.
- **Git.** One `git fetch` of the delivery remote before a prerequisite read
  and before a merge, resolved the way `CreateItemBranch` resolves it.
- **Verification executor.** Regeneration commands run through
  `daemon.ExecVerifier`, as the repository gate does.
- **QA Report reader.** `spec.ReadQAReportFile` semantics, applied to a report
  read from the Run Branch.

## Testing Approach

Every test uses temporary repositories, a temporary Roundfix Home and fake
command runners; none reaches GitHub or runs `gh`.

1. **Park Classes and the re-run.** `internal/delivery/park_class_test.go`
   parses the package's Go files and requires a class other than
   `unclassified` for every `Blocker…` constant, so a new blocker without a
   class fails. `internal/delivery/check_rerun_test.go` drives the engine with
   a fake check sequence and `GitHubCLI` with a fake runner: a failure outside
   the change is re-run once; a changed package, a build failure or a log
   without a Go package parks `checks-failed` with no re-run; a second failure
   parks `flaky-check`. `internal/cli/deliver_park_status_test.go` reproduces
   Surface Transcript 1. `TestPendingQuestionAnswersEachBlockerClass` and
   `TestDeliverStatusPrintsNoWarningLineWithoutAWarning` pass unchanged.
2. **Prerequisites.** `internal/spec/requires_test.go` covers the manifest
   rules. `internal/delivery/prerequisite_test.go` covers waiting, parking,
   release with an unchanged retry count, the Retry path and an item without
   `requires`. `internal/cli/deliver_prerequisite_test.go` uses a bare origin
   and a clone whose local default branch is stale, and reproduces Surface
   Transcript 2 and the unknown-slug refusal.
3. **Operator-archived retry.** `internal/delivery/operator_archive_retry_test.go`
   covers the park, the accepted retry, the refusal without an override, the
   refusal of a head that does not descend from the anchor, and the archive
   pass-through. `TestRetryRefusesAnArchivedItemWhoseHeadMoved` passes
   unchanged. `internal/cli/deliver_operator_archive_test.go` covers the Run
   Branch report read, the override read and the authorization read with a
   commit after the archive, over real repositories.
4. **Conflicts.** `internal/config/delivery_derived_paths_test.go` covers the
   key. `internal/delivery/conflict_test.go` covers `CONFLICTING`, `UNKNOWN`,
   both resolutions, the lagging Pull Request head and the retry.
   `internal/cli/deliver_conflict_test.go` builds a bare origin whose default
   branch and item branch change one derived file and one source file, with a
   regeneration command that writes a known digest.
5. **Docs.** Each Task phrase-checks its behavior in
   `docs/user-guide/commands/deliver.md`. task_04 also checks the Roundfix
   skill's delivery reference, the configuration guide, the recorded skill
   version and `make skills-sync-check`.

## Build Order

1. Park Classes, the `Park:` status lines and the failed check re-run, with
   the command guide, task_01 (depends on: none).
2. Prerequisites, with the command guide, task_02 (depends on: 1). It adds the
   `dependency` class and edits the engine and the guide.
3. Environment-only partial, the operator-archived retry and the pre-archive
   authorization read, with the command guide, task_03 (depends on: 2).
4. Conflicts, the derived-path declaration, the configuration guide and the
   Roundfix skill's delivery reference for all four behaviors, task_04
   (depends on: 3). It reuses task_03's retry acceptance and pre-archive
   authorization read.
5. Terminal QA, task_05 (depends on: 1, 2, 3, 4).

Every Task edits `internal/delivery/engine.go`, the Park Class table and the
command guide, so the chain is serial. The graph also depends on Spec 0194,
which creates the delivery reference files, and on Specs 0189, 0190 and 0196,
which edit `internal/config/config.go` and the Roundfix skill first. The queue
does not enforce that order; the operator does.

## Risks & Considerations

- **A mis-declared derived path.** A declared path that carries source loses
  the item's side. The gate runs before any push, and the ADR-0192 limit is in
  the guide.
- **Package attribution is coarse.** A change can break a test in a package it
  did not touch. The one re-run fails again, and the park says the failure is
  outside the change rather than calling it flaky.
- **GitHub lag.** `mergeable` starts `UNKNOWN` and the Pull Request head lags a
  push; both are read again within the existing check timeout.
- **Older owners.** A binary older than this Spec ignores `requires` and does
  not read `mergeable`. It never sees `delivery.derived_paths`, because this
  repository does not adopt the key in this Spec.

## Decisions

- Merge, never rebase; no second review for a derived merge. See ADR-0192.
- Prerequisites live in `_tasks.md`; the owner releases a dependency park.
  See ADR-0193.
- Optional engine dependencies keep every existing fake valid.
- A re-run is attributed by Go package and happens at most once per head.
- The authorization is read before the archive commit, so an owner merge or an
  operator commit after the archive does not hide the grant.
