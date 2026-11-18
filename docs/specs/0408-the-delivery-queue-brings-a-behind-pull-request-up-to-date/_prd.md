---
spec: 0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date
status: active
created: 2026-11-18
surfaces: [backend, cli, docs]
type: fix
---

# The Delivery Queue brings a Pull Request that fell behind its base branch up to date

A delivery Pull Request whose base branch requires branches to be up to date before merging is refused at the merge when the base branch moves, and the Delivery Queue parks it `delivery-error` on the second refusal; the recovery runs the review, the gate and CI again. After this Spec the Queue Owner reads GitHub's `BEHIND` merge state at `checking` and makes a Base Update: it merges the refreshed default branch into the item branch, runs the repository gate on the merged tree, pushes it and waits for the checks of the new head before the merge. The pre-PR review, the QA gate and the Run are not repeated, a conflict parks as it does today, and three Base Updates per item bound a base branch that keeps moving.

## Sources

Adopted: the Backlog Entry [`references/2026-11-16-the-delivery-queue-and-the-release-guide-assume-loose-branch-rules-and-tag-releases.md`](references/2026-11-16-the-delivery-queue-and-the-release-guide-assume-loose-branch-rules-and-tag-releases.md), reported by quaestor (ai-memory message `01a1e522`, 2026-11-16). Its second observation, about tag-triggered releases, was split from it on 2026-11-16 and delivered by Spec 0400.

The author's reproduction, the reading of the existing merge of the default branch, a prototype of the design and the published GitHub evidence are in [`references/2026-11-18-behind-pull-request-characterization.md`](references/2026-11-18-behind-pull-request-characterization.md).

No other open Backlog Entry and no unresolved Finding shares this context: the open entries of 2026-11-11 (the Conexus status and authorization wait) and 2026-11-16 (the session context size) change other components, and the three entries of 2026-11-16 and 2026-11-18 are adopted by Specs 0407, 0409 and 0410.

## Prerequisites

None of this Spec's Verification depends on another active Spec. Spec 0407, authored beside it, delivers first and may also change `CONTEXT.md`, `docs/references/behavior-surfaces.json`, `skills/testdata/owned-skill-versions.json` and the version fields of the Roundfix skill; task_02 edits them narrowly and re-records the derived files on the tree it runs on. Authored against `1a8e08b4`.

## Project Constraints

- Identifier strategy: not applicable — no persistent identifier is minted. The new names are the blocker `pull-request-behind`, the owner log line `roundfix: base update: …` and the merge state `BEHIND` the queue now reads, all public CLI contract under `docs/agents/cli.md`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no endpoint, credential or new forge call. The owner keeps reading the merge state through the `gh pr view --json` call it already makes, and the Base Update uses the `git fetch`, `git merge`, commit and push it already makes for the derived merge. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0435 (new) records the design; it keeps ADR-0192, ADR-0369, ADR-0311, ADR-0366, ADR-0422, ADR-0302, ADR-0236 and ADR-0398; ADR-0115, ADR-0116, ADR-0197, ADR-0233, ADR-0250, ADR-0256, ADR-0282, ADR-0286, ADR-0309, ADR-0345, ADR-0352, ADR-0365 and ADR-0391 are listed only to record why they do not apply; the TechSpec relates each one. Source: `docs/agents/domain.md`, `docs/agents/specific-repository.md`.
- Tooling authority: applicable — express maintainer authorization: "Mesmas do BE" (2026-11-18), quoted in the authorization record. It covers the Roundfix skill files listed below. It grants no Makefile, workflow, Go module, lint or formatter change, and this design needs none. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`; Spec-contained authorization record: `docs/specs/0408-the-delivery-queue-brings-a-behind-pull-request-up-to-date/_authorization.md`; bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/deliver.md`, `skills/roundfix/SKILL.md`.

## Goals

- A delivery Pull Request that GitHub reports `BEHIND` never reaches the merge at that head; the Queue Owner merges the refreshed default branch into the item branch and takes the merged head through the repository gate, the push and its checks.
- The Base Update keeps the item branch the owner's: a merge, never a rebase or a forced push, committed in the item worktree with the trailer the Correction Sign check and the hook-refusal recovery already recognize.
- The pre-PR review, the QA gate and the Run are not repeated after a Base Update.
- A Base Update that conflicts parks as a pull request conflict does today; one that finds nothing to merge keeps the Check Wait polling.
- At most three Base Updates per item and Queue Owner; behind a fourth time, the item parks `pull-request-behind` with a next action.
- A repository whose base branch does not require up-to-date branches sees no change.

## Core Features

1. The checking stage of the Delivery Engine treats `BEHIND` as a merge state that never merges, stops the Check Wait on it, and makes a Base Update through the queue's existing merge of the default branch.
2. The bound, the `pull-request-behind` blocker, its Park Class `environment` and its next action, and the `; merge state BEHIND` suffix of a `checks-timeout` park.
3. The `roundfix deliver` reference of the Roundfix skill, the user guide page for `deliver`, the documented park actions and the glossary describe the Base Update.

## Non-Goals / Out of Scope

- GitHub's update-branch endpoint, a rebase, or a forced push (ADR-0435).
- Detecting a stale branch whose base does not require up-to-date branches.
- A merge queue, auto-merge, or a change to how the release guide treats tag-triggered releases (Spec 0400).
- Persisting the Base Update count across Queue Owner processes.
- Any change to the Pull Request policy refusal at a head that is not behind: its second refusal still parks `delivery-error`.

## Success Metrics

1. Success Metric: at `checking`, a Pull Request reported `BEHIND` sends the item back to `gating` with the merged head as its newest candidate, and no merge is attempted at the behind head.
2. Success Metric: through the command engine, a strict local forge that refuses the merge while the Pull Request head lacks the newest commit of `main` ends with the item `merged`; the merge commit holds the base branch's new file and the owner log names the Base Update.
3. Success Metric: a Base Update with a source conflict parks `pull-request-conflict: <path>`; one that finds nothing to merge keeps polling and merges once the merge state is `CLEAN`, or parks `checks-timeout: …; merge state BEHIND` at the end of the Check Wait.
4. Success Metric: behind a fourth time, the item parks `pull-request-behind: pull request #<n> is still behind its base branch after 3 Base Updates`, of Park Class `environment`, and after that park the next `BEHIND` makes a Base Update again.
5. Success Metric: a Pull Request reported `CLEAN`, `BLOCKED` or `UNKNOWN` keeps today's behavior, and the whole `internal/delivery` package passes.
6. Success Metric, outside evidence: GitHub's documentation and a published measurement agree that `BEHIND` means the head is out of date and surfaces only when the base branch requires up-to-date branches.

## Acceptance evidence

The outside-evidence row rests on published sources this Spec did not write: the GitHub GraphQL reference for `MergeStateStatus` and for the strict status-check settings, GitHub's "Troubleshooting required status checks", and the measurement in multiplex-ai/muggle-ai-works pull request #265. The planning change commits their quoted passages in [`references/2026-11-18-behind-pull-request-characterization.md`](references/2026-11-18-behind-pull-request-characterization.md), so the QA gate checks the row without the network.

## Unreachable Acceptance

- criterion: a delivery Pull Request on a real GitHub repository whose base branch requires up-to-date branches merges after a Base Update
  reason: the QA gate may not reach the network, run the real `gh` or write to a forge repository; the strict local forge of Success Metric 2 stands in for GitHub
  satisfied-by: the first delivery of an adopter whose `main` requires up-to-date branches, such as quaestor, after the release that carries this Spec

## Validation

The design was checked on 2026-11-18 against the reproduction and the prototype recorded in the characterization reference, the code of `internal/delivery/engine.go`, `park_class.go` and `internal/cli/deliver_workflow.go` (`ResolveConflict`), ADR-0192, ADR-0366 and ADR-0369, the Secondbrain pages that record the Quaestor's strict `main`, and the published GitHub sources above.

## Glossary

- adds: **Base Update**
- changes: **Check Wait**

## Decisions

- 2026-11-18: Backlog 2026-11-16 Delivery Queue com branch estrita (Quaestor)? — "Spec fix após a corretiva (Recommended)" — "Mesmas do BE".
- 2026-11-18: Ordem? — "Corretiva → regras → Indexus (Recommended)"; o fix da fila entra logo após a corretiva.

Autonomous route: no person present during authoring; the remaining decisions were taken from the operator's brief for Spec 0408 of 2026-11-18 (`BRIEFING-authoring.md`, section "Cycle BG"), the two answers above and the characterization reference.

- The Base Update reuses the queue's merge of the default branch and its commit, and never GitHub's update-branch endpoint (ADR-0435).
- The repository gate, the push and the checks run again; the pre-PR review, the QA gate and the Run do not (ADR-0435, ADR-0192).
- Three Base Updates per item and Queue Owner, counted in the owner's memory and dropped when the item parks (ADR-0435).
- `type: fix` with a `patch` Release Ceiling: a refusal that parked becomes a merge, and the new blocker names only a bound the old behavior never reached.

## Open Questions

None.
