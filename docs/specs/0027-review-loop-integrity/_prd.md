---
spec: 0027-review-loop-integrity
status: active
created: 2026-07-14
surfaces: [cli, docs]
---

# Review Loop Integrity

Field runs against real pull requests (PR #4 and PR #17 sessions, 2026-07-14) reproduced three integrity gaps in the review loop. First, a review Run can start from a HEAD that omits work stranded on Run Branches in kept worktrees — it then reviews stale code, and its Final Push can silently publish a HEAD missing completed work. Second, when the Review Source check never appears for the pushed head, watch declares Clean with only a stderr warning, so a script caller cannot tell a verified Merge-Ready outcome from an unverified one. Third, Review Issue outcomes are invisible or misleading on GitHub: invalid issues are resolved silently, failed and duplicated threads stay open with no explanation, and the final report mixes per-Run and cumulative counts. This spec makes the review loop anchored, verifiable, and auditable end to end.

## Goals

- A review Run can never silently review or push a HEAD that omits pending Run Branch work: it starts on a reconciled PR Head Branch or refuses with the exact recovery action.
- A script caller can distinguish a verified Clean (Merge-Ready) from an unverified one by terminal outcome and exit code.
- A reviewer following the pull request on GitHub can tell every Review Issue's outcome and its reason without access to local Roundfix artifacts.
- The final review report separates this Run's counts from the pull request's cumulative counts.
- A Review Issue artifact whose Agent-claimed resolution was overridden by the Daemon records the terminal reason.

## User Stories

1. As a developer running watch, I want the review Run to operate directly on my checkout on the PR Head Branch, so that review fixes are always a delta over the published HEAD and no work strands in a Run Worktree.
2. As a developer starting watch or fetch, I want a Branch Integrity Preflight that integrates pending Run Branch work automatically when fast-forward resolves it — and otherwise blocks naming each pending worktree and the exact integration command — so that I never review or fix stale code.
3. As a developer starting watch or fetch, I want the same preflight to block while any Run bound to the branch or working tree is active, naming the run id and the stop command, so that concurrent Runs never race on my branch.
4. As a developer in an exceptional situation, I want an explicit bypass flag that publishes an audit comment on the pull request recording what was skipped, so that shortcuts stay visible to everyone following the review.
5. As a script caller using until-clean watch, I want a distinct Clean Unverified terminal outcome and exit code when the Review Source check never appears for the pushed head, so that my automation never treats an unverified push as Merge-Ready.
6. As a human reviewer on GitHub, I want Review Issue threads updated as their fixes settle and an Outcome Comment explaining every non-resolved outcome, so that the pull request is self-auditable.
7. As a Supervisor inspecting a failed Review Issue artifact, I want the terminal reason recorded in the artifact — which step failed, the command, its exit status, and where the diagnostics live — so that recovery does not require replaying the journal.
8. As a script caller parsing the final report, I want this Run's counts reported separately from the pull request's cumulative counts, so that "how much did this Run resolve" is derivable.

## Core Features

1. **Review Runs execute in the user's checkout.** fetch, resolve, and watch operate directly on the checked-out PR Head Branch and create no Run Worktree. Worktree isolation remains the contract for spec Runs (implement), where Task concurrency needs it. Integration Pending ceases to exist as a review-Run outcome. See ADR-0042.
2. **Branch Integrity Preflight, on watch and on fetch.** Before any fetch, Agent Session, review comment, or code change, the preflight enumerates every Roundfix Run Branch and kept worktree whose base is the PR Head Branch. Pending commits are integrated automatically when fast-forward resolves them; otherwise the preflight fails naming each pending worktree, its divergent commits, and the exact integration command. The Run starts only with zero pending Run Branch work related to the branch.
3. **Active-Run guardrail.** The same preflight blocks while any Run bound to the branch, the pull request, or a related worktree is active, naming the run id and the stop command (including the force form when the owning process may be dead).
4. **Audited bypass.** A single explicit flag skips both guardrails. When used, Roundfix publishes a pull request comment recording who and when, the run id, which guardrails were skipped, and the ignored state (pending worktrees and active Runs enumerated). If the comment cannot be published, the bypass fails — the audit trail is part of the contract.
5. **Merge-Ready confirmation with a grace period.** After Final Push, watch polls for the Review Source check on the pushed head using the same settle-wait mechanism the fetch step already uses. Check success with no new Review Issues ends Clean; new Review Issues start the next Round as today; a check that never appears within the grace period ends the Run Clean Unverified — a distinct terminal outcome with its own exit code. See ADR-0043.
6. **Per-issue status propagation.** Review Issue thread status is updated on the Review Source as each issue's fix settles. If per-issue calls measurably slow the Round, propagation degrades to end-of-Batch — never later. The GitHub state tracks the Run's real progress, and no thread is resolved before the code and the mandatory verification gate confirm the result.
7. **Outcome Comments for every non-resolved outcome.** Invalid issues get a brief comment with the triage reason before their thread is resolved; failed issues get a comment naming the failed step and the needed action, and their thread stays open; duplicated issues get a comment pointing at the canonical thread before resolution; issues still unresolved at Run end get a comment saying why and when they will be revisited. Comments are idempotent across retries and are journaled with the Review Issue reference.
8. **Terminal reason in the issue artifact.** When the Daemon converts an Agent-claimed resolved issue to failed, the artifact records the failed step, the command, its exit status, and the diagnostics location alongside the status flip.
9. **Per-Run and cumulative counts.** The final report states this Run's issue counts and the pull request's cumulative counts as separate, labeled figures.
10. **Skill and vocabulary shipped together.** The Roundfix Skill sections and command resources for watch, fetch, resolve, and Batch resolution are updated in the same delivery: the no-worktree review contract, the Branch Integrity Preflight and its bypass, Clean Unverified, Outcome Comments, and the clarification that a verification retry never consumes a Round nor counts as a new Review Source review. The glossary entries for Run Worktree, Fetch Run, and related terms are aligned.

## User Experience

- Preflight refusal is deterministic and actionable: each pending worktree or active Run is listed with the one command that unblocks it. No partial work happens before the refusal.
- The bypass flag is explicit and never implied by other flags; its PR comment is visible to every reviewer.
- The final report reads the same for humans and scripts: outcome name, this-Run counts, cumulative counts, and — for every failed or unresolved issue — a one-line reason.
- Clean and Clean Unverified are visually and programmatically distinct; the Clean Unverified report names the next action (confirm the Review Source check before merging).

## Non-Goals / Out of Scope

- ACP adapter internals: the acpx internal-error crash at batch close, subagent writes not persisting, and context hallucination are upstream runtime issues; the mitigations (settle recovery, fallback to Codex) are already documented.
- Verification exit-code allow-lists or a skip-verify settle mode — the authoring contract owns hermetic, satisfiable Verification.
- Any change to the spec-Run (implement) worktree contract, Task Worktrees, or the cherry-pick integration queue.
- Orphaned-lock liveness detection — spec 0028; until it ships, the Active-Run guardrail names the force stop command for dead-owner cases.
- Review artifact commit policy — shipped (separate docs commit, auto-commit opt-out, round number in the message).
- The Reprocess Command.

## Success Metrics

- Zero review Runs in dogfood sessions start while unintegrated Run Branch commits exist for the branch (the preflight blocks or integrates them).
- Integration Pending no longer occurs for review Runs.
- A script can distinguish Clean from Clean Unverified by exit code alone.
- In dogfood pull requests, every thread whose outcome was not resolved carries an Outcome Comment, and no invalid thread is resolved without one.
- The final report answers "how many issues did this Run resolve" without consulting artifacts.

## Decisions

- Review Runs abandon Run Worktrees and run in the user's checkout; spec Runs keep worktree isolation. See ADR-0042.
- A missing Review Source check after Final Push ends Clean Unverified, never Clean; watch polls through a grace period first. See ADR-0043.
- Thread status propagates per issue, degrading to per-Batch only under measured latency pressure, never later than Batch end.
- Invalid and duplicated threads are resolved only after their Outcome Comment is published; failed and unresolved threads stay open with an Outcome Comment.
- The bypass is a dedicated explicit flag with a mandatory published audit comment; publish failure fails the bypass.
- The verification-retry-versus-Round distinction is documentation and skill contract, not a behavior change: a failed Batch already ends the Run Unresolved without consuming extra Rounds.

## Open Questions

- Grace-period duration for the post-push check poll — default until tuned: the same quiet-period setting the fetch wait already uses.
- Bypass flag name — default until the tech spec settles CLI surface: a self-describing dedicated flag (not an overloaded `--force`).
- The latency threshold that triggers per-issue → per-Batch degradation — default: degrade only when per-issue propagation is observed to dominate Round wall-clock; the tech spec defines the measurement.
