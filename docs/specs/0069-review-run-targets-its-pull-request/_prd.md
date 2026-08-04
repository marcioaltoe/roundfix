---
spec: 0069-review-run-targets-its-pull-request
status: active
created: 2026-08-02
surfaces: [backend, cli, docs]
---

# A Review Run targets its Pull Request

`roundfix watch --source coderabbit --pr N` names the Pull Request it should
process, but resolves the branch it acts on from the main checkout instead. The
two are assumed to agree and nothing checks that they do, so a mismatch is
silent until its consequences appear.

Both failure modes cost a round in one session. Running the command for one
Pull Request while the checkout sat on a different branch pushed that Pull
Request's review-artifact commit onto the wrong branch, where it had to be
cherry-picked out and the branch force-pushed. Switching the checkout while a
Review Run was Active made every Review Issue fail with
`checkout branch mismatch: expected <branch> at <sha>, found <other> at <sha>` —
two legitimate security findings failed for an environmental reason and had to
be re-run from scratch.

The second case proves the check exists; it just runs too late, after the
Review Source has been queried and the Agent Session has started.

## Project Constraints

- Identifier strategy: not applicable — Pull Request numbers, branch names, and
  Run IDs keep their existing contracts; no project-owned Internal Identifier
  is created. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the governing clause forbids inventing
  authentication or transport policy and prohibits printing secrets; resolving
  a Pull Request's head branch uses the existing forge access path and adds no
  new credential or endpoint. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0052 protects terminal completion,
  so a refusal must never leave a Run in a state that cannot settle. ADR-0053
  is relation-only here: its Run Worktree reconciliation contract belongs to
  spec Runs, while this Spec changes review Run target validation. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized; the work is product code, CLI surface, and documentation.
  Source: `docs/agents/agent-instructions.md`.

## Goals

- A Review Run acts on the branch its Pull Request names, or refuses before it
  does anything observable.
- A checkout that disagrees with the target is caught at Preflight, not after
  the Review Source has been queried and an Agent Session started.
- A review artifact can never be committed or pushed to a branch other than the
  one under review.

## Core Features

1. The Run resolves its target branch from the Pull Request it names, and
   validates the checkout against that target before any Review Source query,
   Agent Session, or write.
2. A checkout that does not match refuses at Preflight with both branches and
   both revisions named, and creates no Run.
3. Every review-artifact commit and push targets the Pull Request's head
   branch, so no artifact can land on an unrelated branch.
4. A checkout that moves while a Review Run is Active is detected and reported
   as an environmental interruption, distinctly from a Review Issue that failed
   on its own merits, so a re-run is not confused with a real failure.
5. The refusal names the command that would succeed, so recovery is one step.

## Non-Goals / Out of Scope

- Checking out branches on the user's behalf, or moving the working tree.
- Changing how Review Issues are fetched, batched, resolved, or classified.
- Changing Branch Integrity Preflight's Run Branch rules, owned by Spec 0066.
- Supporting a Review Run against a Pull Request from a fork.

## Success Metrics

- Running a Review Run for a Pull Request while the checkout is on another
  branch refuses at Preflight, creates no Run, and names both branches.
- Replaying this session's two failures, neither produces a misplaced artifact
  or a Review Issue failed for an environmental reason.
- Every review-artifact commit in a Run's history lands on the Pull Request's
  head branch, asserted from Git rather than from the log.
- A checkout that moves mid-Run is reported as an interruption, not as a Review
  Issue failure.

## Decisions

- The Pull Request is the authority for the branch; the checkout is validated
  against it, never the other way round.
- The check moves to Preflight, where the existing contract already promises no
  side effects on refusal — the mismatch was already detected, just too late to
  be cheap.
- This Spec evolves the Review Run and never regresses it: a Run whose checkout
  already matches its Pull Request behaves exactly as it does today.

## Open Questions

None.
