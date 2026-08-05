---
date: 2026-08-05
supersedes: 2026-08-03-after-the-0.3.1-release.md
---

# Handoff — the night the queue moved

Read `docs/workflow/spec-queue.md` for the ordered list. This records what
shipped across 2026-08-04 into 2026-08-05, what it cost, and the one pattern
worth changing before the next session.

## Where things stand

**Four Specs archived**: 0064 (Spec artifact consistency gate), 0076 (Force
Stop exit proof), 0070 (declared unreachable acceptance), 0068 (Spec close
audit). The queue went from eleven active Specs to seven.

**In flight when this was written**: Spec 0066 passed its gate and is in review
on PR #116; Spec 0067's planning PR is #117, awaiting checks.

**Authorizations granted this session**, all recorded in
`docs/workflow/authorizations/2026-08-04-queue-tail-tooling.md`:

- Specs 0059, 0065 and 0073 — the boundaries their PRDs cited but no record
  named. Spec 0059's cited a 2026-07-28 record that does not exist.
- Spec 0070 — the `qa-gate` skill, because Core Features 2, 3, 7 and 8 all live
  in that skill and no product code path can deliver them.
- **Standing**: `.agents/skills/roundfix/**` for CLI-contract synchronisation,
  bounded by purpose rather than by Spec. Five of the remaining Specs change
  CLI behaviour and would each have stopped at close without it.

**Branch protection is live on `main`**: `Verification gate` and
`Validate PR title` required, no required reviews, `enforce_admins` false,
force-push and deletion blocked. It earned its keep the same night — Spec
0068's archive PR failed CI and was refused a merge.

## What shipped

**0064 — `roundfix spec check`.** Eleven `SC-*` detectors over a Spec folder,
each error naming both sides of a contradiction with file and line. Detection is
by citation and declaration, never by subject-matter inference (ADR-0093), which
was settled by probe: only 1 of 89 accepted ADRs names a repository path, so an
ADR-to-code index cannot work, while the ADR citation graph can. Detectors are
presence-aware (ADR-0094). Wired into `make verify` fail-closed.

**0076 — the Force Stop exit proof.** Found *by* 0064's gate, outside its scope:
`TestOwnerProcessControllerForceKillExitProof` had been green for months while
proving nothing. Its `ignore` helper blocked on `select {}`, which Go's runtime
kills as a deadlock, so the process the test believed it was force-killing was
already dead. The repaired proof asserts **causation** — that the controller's
escalation ended the process — because "the process is gone" is satisfied by a
process that crashed on its own.

**0070 — declared unreachable acceptance.** A Spec declares, up front, an
acceptance no hermetic Verification can reach; the gate matches rows against
declarations and never declares unreachability itself. `rows_blocked_declared`
joins ADR-0080's typed counts as a third cause. Archive gains exactly one
accepting case and stamps the remaining debt under `unproven`, so `qa_override`
goes back to meaning genuinely failed evidence.

**0068 — `roundfix spec audit`.** Classifies every branch and worktree a Spec's
cycle produced, with the evidence that classified it, and verifies claimed
artifacts are present on the synced default branch. A Run whose target branch
its own squash merge deleted now resolves by content — proving integration only,
never disproving it, because a false `safe` is the one error that loses work.

**0066 — Run teardown** (in review). Termination reaches the process tree
including a child that outlives its parent, and an unprovable termination is
reported as unproven rather than as success — ADR-0044 reclaims orphaned locks
by reading that distinction.

## What it cost, and where

**Twenty-eight implementation Tasks settled on Verification attempt 1.** One
failed on the work. Every other failure was a defect in what the Supervisor
wrote.

Spec 0070's Core Feature 7 is the measurable win. One finding used to suppress
most of a matrix; after it shipped, a single finding blocked exactly one row:

| Gate | One finding blocked |
| --- | --- |
| Spec 0072 (2026-08-03) | 15 of 24 rows |
| Spec 0064 (2026-08-04) | 13 of 15 rows |
| Spec 0070, first run | 1 of 17 |
| Spec 0068, first run | 1 of 18 |

## The pattern worth changing

Recorded in full at
`docs/findings/2026-08-05-an-absence-grep-rejects-the-work-it-was-written-to-protect.md`.

Three Task Verifications failed in one night. **None failed on the work.** All
three were hand-written `grep` invocations with exclusion filters, authored by
the Supervisor, that rejected correct implementations:

- Spec 0070 task_04's changed-file allow-list omitted the ADR-0081 digest
  fallout that the Spec's own authorization record blesses.
- Spec 0068 task_02's Requirement contradicted a References note added to
  silence a coverage detector — prose instead of design.
- Spec 0068 task_08's negative grep rejected a correctly named constant,
  because its filter matched lowercase `reclaim` and the constant carried
  `Reclaim`.

The rule this suggests: **a Verification command asserts a behaviour the code
must exhibit, not a shape the source must have.** Where only a source-shape
check is possible, it belongs in review, not in a gate that stops a Task. The
third case was fixed that way — a fixture whose Git state is byte-identical
before and after — and passed immediately.

## Traps this session hit

- **Nothing watches a Pull Request after it is opened.** Spec 0070's review sat
  unaddressed for 1 h 32 min because the Supervisor was monitoring the Run, not
  the Pull Request, and no signal arrives when a review lands. Open the Pull
  Request and start `roundfix watch --until-clean` as one step, not two.
- **`git stash pop` with stderr silenced** put conflict markers into Spec
  0059's PRD and merged them to `main` in PR #102. The repository's own HARD
  RULE about pipes hiding exit status is the same defect one tool over.
- **Load-sensitive test failures appeared seven times**, in seven different
  tests, mostly on documentation-only changes. Two were repaired in #104 by
  naming their budgets; the rest were reruns. `implementWaitBudget` already
  states the principle: *a tight budget here measures how loaded the machine
  is, not whether the code works.*
- **Two characterization corpora sit outside `BASELINE_DIGEST_STEPS`**, so
  `make baseline-digests` reports "no changes" while the gate stays red. The
  flags do not match their test names and were guessed wrong again this
  session. **Spec 0067 owns this** and its planning is done.

## Waiting on a decision

1. **Spec 0073's close wants a release.** A release is irreversible and
   outward-facing; it was deliberately left for the maintainer. 0073 can be
   authored and run up to that point.
2. **Whether `spec check` should read `docs/workflow/` and `docs/handoffs/`.**
   This session began by trusting a queue file and a handoff that both said
   Spec 0064 had a TechSpec. It did not. The consistency gate has no reach into
   the documents that route the work.
3. **Whether the detector families should run on a PRD-only Spec.**
   `SC-TOOLING-UNAUTHORIZED` — the check for the defect that put 0064 first in
   the queue — is skipped entirely when `_techspec.md` is absent, which is the
   majority of the queue. Presence-awareness was implemented as "skip the
   detector if any input is missing" rather than "check what is present".

## Starting the next Spec

0066 and 0067 are in flight. After them the queue reads: 0075, 0069, 0059,
0073, 0065.

The mechanics per Spec, with this session's corrections folded in: branch,
author the TechSpec and graph, **correct the Tooling authority row during
authoring** — five Specs in a row had it wrong, and it is cheaper to fix before
a gate than at close — run `spec check` until clean, `make verify`, planning
Pull Request, merge, branch, `roundfix implement --detach`, integrate, open the
Pull Request **and start the watch together**, merge, archive.
