---
spec: 0098-a-hook-that-cannot-outrank-the-gate
status: active
created: 2026-08-12
surfaces: [backend, docs]
---

# A hook that cannot outrank the gate

The Daemon runs the authoritative Verification and then commits. When the
repository's commit hook refuses, the work is reverted and the Run ends failed —
and the repair loop covers a Verification failure, not a hook failure. Three Runs
died this way in one Spec, every time with work that was correct and already
verified, left staged in a Task Worktree. Two gates where one cannot be satisfied
nor recovered is a design conflict, and the invariant that resolves it — a commit
hook may never be stricter than the authoritative Verification — is written
nowhere. The command built for recovery makes it worse: it refuses work that is
completed but uncommitted, because its contract assumes lost work is always
failed.

## Project Constraints

- Identifier strategy: applicable — Verification, Task Worktree, Run, and the Task
  status vocabulary are glossary terms this Spec changes the transitions of, and
  the state "verified, settled, not committed" may need a term the glossary does
  not carry. The closing node checks whether the work introduced or changed one.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is commit orchestration, recovery, and a
  Baseline module clause. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0014 gives the Daemon ownership of task
  verification and status settlement, which is the decision this Spec's conflict
  sits on: the Daemon is the verification authority and then defers to a second
  one at commit. ADR-0038 allows one Verification repair, which bounds what a
  reclassified hook failure may consume. The decisions extending that ownership
  are accounted and unchanged here: ADR-0020 ranks a parsed prompt result above
  the runtime's exit code, ADR-0057 gives the Daemon exclusive ownership of
  Implement Task status, ADR-0056 separates Task Capacity from Verification
  Capacity, and ADR-0096 with ADR-0117 place the gate's mechanical stage and its
  checks. ADR-0127 places process residue in the readiness diagnostic; this Spec governs the commit step and reports no residue, so it does not apply. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the autonomous-work module gains the
  hook-strictness invariant and the recovery ladder. Express maintainer
  authorization:
  `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
  granted 2026-08-12. Bounded files:
  `internal/baseline/assets/modules/autonomous-work.json` and
  `docs/agents/autonomous-work.md`. Source: `docs/agents/agent-instructions.md`.

## Goals

1. Verified work is never lost to a second gate the loop cannot satisfy.
2. Recovery covers the state where work is verified and settled but uncommitted.
3. The relationship between a commit hook and the authoritative Verification is
   written down, so a repository cannot configure the conflict by accident.

## User Stories

1. As a Supervisor whose repository has a strict commit hook, I want a Run to keep
   verified work rather than end failed, so that a legitimate lint finding costs a
   repair rather than a Run.
2. As a Supervisor recovering after a hook refusal, I want the recovery command to
   accept work that is completed but uncommitted, so that I do not extract a patch
   from a worktree and apply it by hand.
3. As a maintainer configuring a repository for the loop, I want the rule that a
   hook may not be stricter than the authoritative Verification stated in the
   guide, so that the conflict is designed out rather than discovered.

## Core Features

1. **A hook refusal has a defined outcome.** A commit refused by a repository hook
   resolves through one settled shape rather than ending the Run: the Daemon
   commits as the verification authority it already is, or the refusal is a
   repairable class that spends a repair round, or it is detected and names the
   recovery. Which of the three is the product decision this Spec settles.
2. **Recovery covers verified-but-uncommitted work.** The recovery contract
   accepts a Work Item that is completed and uncommitted, not only one that
   failed.
3. **Staging covers what a Task deleted.** A Task that correctly removes a file
   is committed rather than refused: the staging call records a deletion instead
   of failing on a pathspec that matches nothing. This is the same loss as a hook
   refusal — verified work discarded at the commit step — reached through a
   different call.
4. **The invariant is written where a repository reads it.** The autonomous-work
   guidance states that a commit hook may never be stricter than the authoritative
   Verification, so a repository adopting the loop configures the two
   consistently.

## User Experience

A hook refusal reports what the hook objected to and what the loop did about it,
in the Run's own output. A recovery run against verified-but-uncommitted work
states which surface holds the work and completes without a manual patch.

## Non-Goals / Out of Scope

- Changing any repository's commit hooks, or recommending hook content.
- Weakening or bypassing a repository's lint rules. Every finding measured in the
  triggering session was legitimate; the cost was where the check ran.
- Changing the authoritative Verification itself.
- Changing Run integration or push behavior.

## Success Metrics

- The three measured Run deaths — a function over a line limit, a generated file
  over a size limit, and a forbidden array method — each resolve without losing
  the Run. Source: a fifteen-Task Spec delivered in a repository this Spec did not
  build on 2026-08-07/08.
- Recovery of verified-but-uncommitted work completes with the tool, in the two
  measured cases that required a manual patch.
- The hook-strictness invariant appears in the rendered guidance of a repository
  that adopts the Baseline.

## Decisions

- The invariant is carried in the Baseline module rather than in this repository's
  own guide, because the conflict it prevents belongs to every repository that
  runs the loop and a local edit would be overwritten by the next update.

## Open Questions

- Which of the three shapes a hook refusal takes. Committing as the authority is
  the least surprising given that the Daemon already owns verification, and it is
  also the one that most changes what a repository's hooks mean, so it is the
  product decision rather than a derivation.
- Whether the recovery Spec that ships the hand-back path already covers the
  verified-but-uncommitted state. This is confirmed before decomposition, and the
  Spec narrows if so.
