---
spec: 0028-settlement-and-reporting
status: archived
created: 2026-07-14
surfaces: [cli, docs]
archived: "2026-07-15"
source_slug: 0028-settlement-and-reporting
---


# Settlement and Reporting Robustness

Field sessions running spec Runs (specs 0004/0005, 0009/0010 dogfood cycles) lost work or misled operators through settlement and reporting gaps that are each small but compound badly. A killed run process leaves an orphaned Active-Run lock that blocks every relaunch until a manual force stop. An Agent that writes an obvious status synonym (`done` instead of `completed`) voids an entire batch of finished work. A Task whose Verification was already green on the baseline settles `completed` with a commit containing nothing but the status flip. The Settle Command sweeps every change in a shared worktree into one task's commit without saying so. The final report names failed items but never says why they failed, forcing journal archaeology. And a missing ACP adapter binary surfaces late, as an opaque spawn failure with no install hint. This spec closes those gaps.

## Goals

- A dead run process never blocks a new Run: orphaned Active-Run locks are detected and reclaimed safely.
- An Agent's completed work is never voided by a trivially normalizable status synonym.
- The final report alone tells an operator why each failed Work Item failed.
- Silent no-op completions and settle commits that sweep unrelated work are visible at the moment they happen.
- A missing ACP adapter binary is diagnosed at Preflight Validation with the install action.

## User Stories

1. As a developer relaunching after a killed run, I want Roundfix to detect that the lock's owning process is dead and reclaim the orphaned lock with a warning, so that relaunches are never blocked by a Run that no longer exists.
2. As a Supervisor, I want the Daemon to normalize obvious task status synonyms and validate the frontmatter before closing the Batch, so that a wording slip never converts finished work into a failed Task.
3. As a Supervisor, I want a warning and a Run Event whenever a Task commit contains no change outside the Spec Root, so that a no-op completion is caught while the Run is still observable.
4. As a developer using the Settle Command in a worktree shared by several failed Tasks, I want the commit's contents listed and a warning that other failed Tasks' work may be included, so that history stays honest.
5. As a developer reading the final report, I want every failed Task or Review Issue to carry a one-line reason in the report, so that diagnosis starts from stdout, not from the journal.
6. As a developer whose ACP adapter binary is not installed, I want Preflight Validation to name the missing binary and the install command, so that setup failures are self-explanatory.

## Core Features

1. **Orphaned-lock detection and reclamation.** Active-Run locks carry enough identity to check whether the owning process is still alive. When a new Run, the Stop Command, or the review Branch Integrity Preflight encounters a lock whose owner is provably dead, Roundfix treats it as orphaned: it reclaims the lock, records the reclamation in the Run Event Journal, and warns on stderr. A lock whose owner is alive keeps blocking exactly as today, naming the run id and the stop command. See ADR-0044.
2. **Status synonym normalization.** When reloading a task file after Agent work, the Daemon normalizes a small documented set of unambiguous status synonyms to the canonical vocabulary and rewrites the frontmatter before settlement. Statuses outside the canonical set and the synonym set still fail, with diagnostics naming the allowed values.
3. **No-op completion warning.** When a Task settles `completed` and its Task commit contains no change outside the Spec Root, the Daemon settles it but emits a prominent warning and a Run Event. The authoring contract documents that a Task's Verification must prove the Task's effect, with executable checks preferred.
4. **Settle transparency.** The Settle Command reports every path included in its commit and warns when other failed Tasks share the same worktree, since their work is swept into this commit.
5. **Failure reasons in the final report.** The end-of-run report includes, for every failed Task and failed Review Issue, a one-line reason: the failed step and command with its exit status when Verification failed, or the Agent-reported cause otherwise, with a pointer to the full diagnostics.
6. **Early adapter diagnostics.** When the configured ACP Runtime's adapter binary cannot be spawned, Preflight Validation reports that the binary is not installed and names the install command, instead of surfacing a late spawn failure during model probing.
7. **Authoring guidance shipped together.** The spec-authoring guidance gains the portability note proven in the field (prefer portable shell forms over `wc`-style pipelines in Verification gates on macOS) and the effect-proving Verification requirement; the Roundfix Skill is re-checked against the shipped behavior per the skill-sync contract.

## User Experience

- The orphaned-lock reclamation is loud but non-interactive: one stderr warning naming the dead run id, then the new Run proceeds.
- Report reason lines are short, human-readable, and stable enough for scripts to grep; full diagnostics stay in the journal and verification logs.
- The settle path summary prints before the commit is created, so the operator sees what is being swept.

## Non-Goals / Out of Scope

- Review-loop behavior — branch guardrails, Clean semantics, GitHub propagation — is spec 0027.
- Upstream ACP runtime defects (acpx internal errors, subagent write loss, hallucinated context).
- Blocking no-op completions — warn-only by product decision; a block would false-positive on legitimate docs-only Tasks.
- A pathspec-restricted settle mode — transparency first; restriction can follow if the warning proves insufficient.
- Verification exit-code allow-lists or skip-verify settlement.
- Changes to lock scope or the one-Active-Run-per-target rule.

## Success Metrics

- A killed run no longer requires a manual force stop before relaunching, in dogfood sessions.
- Zero Tasks failed for status-vocabulary reasons after normalization ships.
- Operators diagnose failed items from the final report without opening the journal in routine cases.
- Every no-op completion in dogfood runs is accompanied by its warning event.

## Decisions

- Orphaned locks are reclaimed automatically only on proven owner death; anything less than proof keeps the block. See ADR-0044.
- No-op completions warn and journal; they do not block. Effect-proving Verification is an authoring obligation.
- Normalization covers only a small documented synonym set; the canonical vocabulary is unchanged and everything else still fails.
- Settle stays whole-worktree but becomes transparent about it.

## Open Questions

- The exact synonym set — default until the tech spec fixes it: `done` → `completed`, plus hyphen/space variants of the canonical statuses.
- Whether the per-item reason line also appears for `invalid` Review Issues in the report — default: yes, one line, same format.
