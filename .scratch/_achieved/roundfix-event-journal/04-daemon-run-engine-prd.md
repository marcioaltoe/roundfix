---
title: "PRD: Daemon Run engine"
type: PRD
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
blocked_by:
  - 01-run-event-seam-prd.md
---

# PRD: Daemon Run engine

## Problem Statement

Run orchestration lives in the CLI module. The resolve command wires Agent execution, verification, Batch commits, Review Source thread resolution, and push behavior directly, and the watch path duplicates the same wiring. Test injection happens through package globals. ADR 0001 says the Daemon owns Batch commits and Final Push; today the CLI does. The daemon and Watch Run event stream PRD will make event publication non-optional for state transitions, which cannot be enforced while transitions are scattered across two command paths.

The current commit behavior also has three defects observed in real Runs:

1. The Batch commit stages everything in the worktree except Project Config. Unrelated user changes that slipped past Preflight Validation or appeared mid-Run are swept into a commit labeled as a review fix.
2. A triage-only Batch with no code changes makes the commit fail with a nothing-to-commit error, failing a Run whose Batch actually succeeded.
3. Git fsmonitor warnings pollute parsed git output, because Roundfix reads combined output instead of separating stdout from stderr.

## Solution

Move the review-resolution orchestration into a Run engine owned by the Daemon module. The engine executes one resolve cycle over a validated plan and exposes Final Push as a separate explicit operation, so the resolve command and the watch loop share one implementation instead of duplicating it. The engine receives all collaborators, including the Run Event sink, through an explicit dependencies struct, replacing CLI package globals.

The engine also replaces sweep-everything commits with snapshot-diff commits: it records the worktree state before the Agent starts and after it finishes, and stages only the paths the Agent touched. A Batch that changes nothing skips the commit and still succeeds. All Roundfix git invocations disable fsmonitor and parse stdout separately from stderr.

## User Stories

1. As a developer, I want Batch commits to contain only Agent-made changes, so that my unrelated work-in-progress never gets committed as a review fix.
2. As a developer, I want a triage-only Batch to succeed without a commit, so that a Run does not fail when the Agent correctly concludes no code change is needed.
3. As a developer, I want commit-skip outcomes reported clearly, so that I understand why a successful Batch produced no commit.
4. As a developer, I want git output parsing to ignore fsmonitor warnings, so that dirty-worktree and status checks stay reliable on machines with fsmonitor enabled.
5. As a developer, I want resolve and watch to share one orchestration implementation, so that fixes and instrumentation apply to both commands at once.
6. As a developer, I want the Daemon to own Batch commits and Final Push in code, so that the recorded architecture decision matches reality.
7. As a developer, I want Final Push to remain a separate explicit operation, so that push gating stays readable and never runs while Unresolved Review Issues remain.
8. As a developer, I want a Stop Request during a cycle to halt new Batches, verification, commits, pushes, and Review Source mutations, so that interrupted Runs stay safe and inspectable.
9. As a developer, I want Agent-created worktree changes preserved after a stop, so that I can inspect or reuse partial work.
10. As a developer, I want intermediate Run states updated during the cycle, so that the Run Database reflects what the Run is doing while it is active.
11. As a daemon, I want the engine to receive the Run Event sink as a dependency, so that later slices can make daemon event publication non-optional without restructuring orchestration again.
12. As a daemon, I want Review Source thread resolution to stay inside the cycle after verification, so that resolved and invalid findings do not reappear in later Rounds.
13. As a watch loop, I want to call the engine once per Round, so that polling, quiet periods, and round policy stay in the watch state machine.
14. As a resolve command, I want to call the engine once and receive per-Batch outcomes plus the remaining Unresolved Review Issue count, so that exit codes and output rendering stay in the CLI.
15. As a test author, I want the engine to be constructible with fake collaborators, so that the resolve-verify-commit-resolve-push contract is provable without a terminal or network.
16. As a test author, I want CLI package globals for orchestration collaborators removed, so that dependency flow is auditable and tests stop mutating shared state.

## Implementation Decisions

- The Run engine lives in the existing Daemon module alongside the verifier, committer, and pusher abstractions it consumes. No new top-level module is created for it.
- Interface granularity is one resolve cycle, not a whole Run. The engine's first operation receives a validated plan — deduplicated Review Issues, assembled Batches, and an already-created Run — and, for each Batch: runs the Agent, verifies assigned issue statuses, creates the Batch commit when auto-commit is enabled, and resolves Review Source threads for resolved and invalid issues. It returns per-Batch outcomes and the remaining Unresolved Review Issue count.
- Final Push is the engine's second operation and is invoked explicitly by the caller, only when no Unresolved Review Issues remain and auto-push is enabled. Push gating policy stays in the caller, preserving ADR 0001 semantics: never push per Batch or Round.
- The watch state machine stays in its own module, owning polling, quiet period, review timeout, and round policy. It calls the engine once per Round. The resolve command calls the engine once.
- The engine updates intermediate Run states in the Run Database during the cycle. Run creation after Preflight Validation and terminal completion stay with the caller, because in watch a finished cycle does not end the Run.
- The engine constructor takes an explicit dependencies struct: Agent runner, verifier, committer, pusher, Review Source, rounds and store access, an injectable clock, and the Run Event sink fanout. The CLI package-global injection points are removed together with the orchestration they served.
- Daemon event kind taxonomy and full publication coverage belong to the daemon event stream PRD. This slice threads the sink dependency through the engine and keeps Agent-source events flowing unchanged.
- Commit strategy is snapshot-diff: the engine captures the worktree status before starting the Agent and again after it finishes, and the Batch commit stages only paths that changed between the snapshots. Agent-touched paths legitimately include assigned Review Issue files. Pre-existing user changes are never staged, even when they slipped past Preflight Validation or appeared mid-Run.
- When no paths changed between snapshots, the engine skips the commit and the Batch still succeeds. The skip decision becomes a journaled Run Event in the daemon event stream PRD.
- All Roundfix git invocations disable fsmonitor for the invocation and parse stdout separately from stderr. Combined-output parsing is removed.
- Stop semantics move with the orchestration: after context cancellation or a Stop Request, the engine starts no new Batch, verification, commit, push, or Review Source mutation; the stop is published to the sink before the engine returns; Agent-created worktree changes are preserved.
- Batch failure policy preserves current behavior: a failed verification or Agent failure fails the Run. No retry policy is added in this slice.
- The CLI keeps command parsing, Preflight Validation, Interactive Input, plan assembly, Run creation, terminal completion, exit code mapping, and output rendering.

## Testing Decisions

- Good tests prove the engine's externally visible contract through its interface — what was executed, committed, resolved, pushed, and published, in what order — never which private helpers ran.
- The engine is the highest seam: cycle tests use fake Agent runner, verifier, committer, pusher, and Review Source plus a temporary Run Database to prove the resolve, verify, commit, source-resolution, and push contract without a terminal or network.
- Snapshot-diff commit tests prove that pre-existing unrelated changes are never staged, that Agent-touched paths including assigned issue files are staged, and that mid-Run user edits stay out of Batch commits.
- Triage-only Batch tests prove the commit is skipped, the Batch succeeds, and the Run continues.
- Git invocation tests prove fsmonitor noise on stderr does not corrupt parsed status output.
- Stop tests prove that after cancellation no further Batch, verification, commit, push, or Review Source mutation occurs, the stop event reaches the sink, and worktree changes are preserved.
- Final Push tests prove the push operation is independent of the cycle and that the caller-side gating, not the engine, decides when it runs.
- CLI-level tests shrink to dispatch, preflight wiring, and exit code mapping. Existing CLI orchestration tests and their fakes migrate to engine tests rather than being deleted.
- Race tests are required wherever the engine owns goroutines.

## Out of Scope

- The SQLite Run Event Journal and its store APIs.
- Journal persistence of Agent stream events.
- Attach, replay, or cursor-based reads.
- The daemon event kind taxonomy and full state-transition publication coverage.
- Watch state machine changes, polling behavior, or round policy.
- Retry policies for failed Batches or Agents.
- Fetch and Round-artifact persistence orchestration.
- Changing deduplication rules or push policy.

## Further Notes

This slice makes ADR 0001 true in code and is the structural prerequisite for the daemon and Watch Run event stream PRD: once every state transition flows through the engine, making event publication non-optional becomes a local change. It depends on the Run Event seam PRD only for the sink type the engine accepts; it does not require the journal.

The snapshot-diff commit strategy also resolves the three commit defects observed in production Runs, which makes this slice valuable independently of the event-journal sequence.
