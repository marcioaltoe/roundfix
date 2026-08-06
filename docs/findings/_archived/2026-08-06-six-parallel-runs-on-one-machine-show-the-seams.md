---
status: done
created_at: 2026-08-06
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-run-lifecycle-and-branch-integrity.md
---

# Six parallel Runs on one machine show the seams

**Date:** 2026-08-06
**Session:** the Oraculum queue — six Specs (0018, 0021–0025) driven to merged
Pull Requests and two production releases in ~16 hours, with all six Specs
running **concurrently**, each in its own git worktree of the same repository.

Previous queue nights ran Specs serially. This one ran them in parallel, and
that is the finding: **almost every new failure was a shared-resource seam,
not a defect in anyone's work.** The Agents implemented five Specs in
~1h15–1h40 each, simultaneously; the supervisor's night went to the seams
below. Each class carries its evidence and one ask.

## 1. One Run Database is a serial point under parallel Runs

Six concurrent Runs share `~/.roundfix/roundfix.db`. Four of them died in
their first Batch with:

```text
Agent failed: Agent Batch failed: publish Run Events: begin Run Event append: database is locked (5) (SQLITE_BUSY)
```

The `busy_timeout` 5s→30s fix (`d7a38967`) landed mid-queue and no recurrence
was observed afterward — but the field workaround while Runs started with the
old binary was **HOME isolation** (`HOME=<dir>` per worktree, symlinking
`.codex`, `.npm`, `.local`…), and that workaround has a sharp edge: `gh`
credentials do not resolve under a foreign HOME ("The token in default is
invalid"), so `watch` preflight dies with HTTP 401. The queue ended up running
`implement` under isolated HOMEs and `watch` under the real one.

**Ask:** make parallel Runs a first-class story — either a per-repository Run
Database path or an official `ROUNDFIX_HOME` override, documented together
with the `gh`-credential constraint. The 30s timeout makes contention
survivable; it does not make the design concurrent.

## 2. Branch Integrity Preflight reads the whole refs namespace

Worktrees share refs. With six Specs in flight, `watch` for one Spec's Pull
Request refused because it saw **every other Spec's** Run Branches as pending
work — ten branches named, none related to the PR's head branch.
`--skip-branch-integrity` plus its audit comment became the routine path on
every single watch of the night.

**Ask:** scope the preflight to Run Branches whose recorded target is the PR
Head Branch. Other Specs' branches are unrelated by construction; naming them
turns a safety check into a mandatory bypass, and a bypass that is always
taken protects nothing.

## 3. `watch` on a fresh Pull Request waits for a review nobody requested

With `auto_review` disabled and `review_source.request_review: true`, the
request is published only **after a Round's Final Push**. A freshly opened PR
has no Round, so the Run sat in `WaitingForReview` against a 30-minute
deadline for a review that would never arrive; a manual `@coderabbitai review`
comment unblocked it. The preflight coherence check validates the config pair
but not this startup state.

**Ask:** when the Run starts with no current-head Evidence and no Round,
publish the initial review request at startup. The config already declares the
intent; only the trigger is missing.

## 4. Review Source ceilings deserve preflight, and `MaxRoundsReached` conflates two endings

Observed from one Review Source in one night: a skip for **191 files > 150
limit**, and rate-limit refusals after roughly three successive PR reviews.
Both are knowable before work starts. Separately, one watch ended
`MaxRoundsReached` when resolution was in fact **complete** (18/18 issues
handled) but the reconfirmation round was rate-limited — the outcome reads as
"cap exhausted with findings open", which is the escalation case, when the
truth was "work done, evidence unavailable", which is the documented-gap merge
case. The supervisor had to read the issue artifacts to tell them apart.

**Ask:** (a) a preflight warning when the PR's file count exceeds the source's
review limit; (b) a distinct terminal reason (or `evidence_kind`) separating
"rounds exhausted with unresolved issues" from "issues settled, Review Source
unavailable for confirmation".

## 5. Squash merges make every Run Branch invisible to ancestry proofs

The repository merges every PR by squash. A Run Branch whose content landed
via squash is **never** an ancestor of the target tip, so `reconcile`
classifies it `unintegrated`/preserved forever, and old Runs additionally
refused classification with `short ref is ambiguous` on a legacy branch name.
End-of-queue cleanup of ten dead worktrees and branches was done with manual
`git worktree remove` / `branch -D` — exactly what the documentation forbids —
because the supported tool could not prove what the operator could see.

**Ask:** teach `reconcile` a squash-equivalence proof (patch-id or tree
comparison against the target head) and resolve branches by full ref. The
Spec Audit Command's classification would benefit from the same proof.

## 6. The tooling audit should read the registry its own prescription cites

Oraculum's QA gate flagged `src/infra/adapters/consultor/gerar-doutrina.ts` —
an application-code generator, run directly as a module — as an unauthorized
protected-tooling mutation (Spec 0025, F-01). The repository's durable
authorization document explicitly names "executable module under src, no new
package.json alias" as the **sanctioned route around** protected tooling, with
precedent. One full gate cycle went to teaching the gate via Project
Constraints prose what the registry already said.

**Ask:** the qa-gate skill's governance step should consult the repository's
durable tooling-authorization record before classifying, and default-classify
`src/**` application modules as product code. The inverse also held and is
worth keeping verbatim: the gate's refusal of **retroactive** authorization
(0018) was correct, and its prescription — re-land through a bounded
corrective Task after the authorization commit — was executable as written.

## 7. Authorization prose drifts one corrective round behind the work

Three times on one Spec (0023), the gate re-raised the same governance
finding: each corrective slice exercised a newly authorized path (a
`package.json` script, then a new test fixture), and the constraints prose
written the round before no longer named everything. The declaration was
always true when written and stale when audited.

**Ask:** a write-tasks invariant for corrective Tasks: "if this slice touches
an authorized tooling path, update both Project Constraints blocks in the same
slice." One sentence in the template ends the recurrence.

## 8. Confirmed from a second repository

Two classes already filed from the Roundfix-repo queue reproduced here,
unchanged:

- **The entry precondition is incompatible with a Task that repairs red**
  (see `2026-08-06-every-run-that-failed-tonight-failed-on-a-contract.md`,
  class 1). Three Oraculum corrective Tasks had to decompose `make verify`
  into component commands in their `## Verification`, with an explanatory
  note, to avoid settling `failed` on entry — the workaround works but every
  author must rediscover it.
- **The QA gate and the Pull Request cannot both be current** (see
  `2026-08-06-the-qa-gate-and-the-pull-request-cannot-both-be-current.md`).
  The 0022 gate returned `partial` solely for an unresolvable PR row; the fix
  that stuck was a "no PR surface" declaration in every QA Task. The 0018
  gate had reached the same conclusion by itself — making it deterministic
  cost one cycle per Spec that lacked the block.

## What held under parallel load — keep it

- **ADR-0020 classification** and Daemon-owned Verification: no false
  settlements across ~45 Task executions.
- **`settle`** as the supervisor-repair path: verified, committed and
  integrated a hand-finished Task cleanly, twice.
- **Automatic orphan reclamation** never fired wrongly; no Run was lost to
  the six-way concurrency after the timeout fix.
- **Preflight error messages** (projection-table row missing, Context label
  restriction) were exact enough to fix on the first try.
- **`release plan`**: the generic-request→patch-only guard plus the manual
  minor classification with an explicit approval question is the right shape.
  The operational lesson belongs to the operator, not the tool: **always pin
  the version on workflow dispatch** — an unpinned `cog bump --auto` shipped
  `v0.4.0` where `v0.3.1` had been approved.

## Evidence

- Oraculum queue record:
  `oraculum/docs/findings/2026-08-06-o-que-a-entrega-paralela-da-fila-0018-0025-ensinou.md`
  (22 gate cycles, 57 review issues, 5 documented-gap merges, per-Spec cycle
  counts).
- SQLITE_BUSY console lines: Oraculum Runs `run_20260805T200812Z_*`,
  `run_20260805T201057Z_*`, `run_20260805T201149Z_*`, `run_20260805T201245Z_*`.
- Preflight refusal naming ten unrelated Run Branches: watch startup for
  PR #48, 2026-08-05.
- `WaitingForReview` to deadline with no request published: Run
  `run_20260805T211951Z_*` (PR #48).
- `MaxRoundsReached` with 18/18 issues settled: Run
  `run_20260806T050233Z_*` (PR #51).
