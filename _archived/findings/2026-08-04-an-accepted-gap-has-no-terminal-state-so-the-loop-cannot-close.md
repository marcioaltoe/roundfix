---
status: done
created_at: 2026-08-04
updated_at: 2026-08-06
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
---

# An accepted gap has no terminal state, so the autonomous loop cannot close (2026-08-04)

A Vortex session drove Specs 0013, 0014 and a Baseline upgrade from a paused queue to
merged `main`. The implementation work was almost fully autonomous. The **closing** of a
Pull Request was not: it required six maintainer interventions, and **four of the six were
mechanical or encoding problems, not decisions of authority**.

That ratio is the finding. The operator's stated goal is a loop where they author the PRD
and everything downstream runs unattended. Measured against that goal, Roundfix's
implementation half is close; its *merge* half hands control back for reasons that carry no
judgement.

The proximate cause of the worst stall is narrow and fixable: **there is no Review Issue
status meaning "valid finding, deliberately not fixed, accepted by the maintainer."** Every
available terminal status either lies about the finding or blocks the loop forever.

This is distinct from
[gate and review rounds need a convergence rule](2026-08-03-gate-and-review-rounds-need-a-convergence-rule.md),
which describes the gate *racing* the reviewer. Here the loop cannot converge **by
construction**, no matter how patiently it waits.

## Symptom / evidence

PR #116 (Spec `0013-sync-cortes-de-dados`) took **five review rounds**. Rounds 001 and 002
each found a real defect that a hollow `Clean` would have shipped, so the rounds were
earning their cost. Cumulative outcome:

```text
Pull Request cumulative: 20 resolved, 6 invalid, 0 duplicated, 3 failed, 0 unresolved
```

All three `failed` are **the same finding**, re-raised each round: the new
`captureFinancialTransactionSourceVersion` query is covered only by a test that stubs
`Pool.query`, so it passes even if MySQL would reject the SQL. The finding is correct. The
resolving Agent triaged it `VALID` and then could not act, for a reason it verified
independently each round:

```text
Exact blocker: no hermetic MySQL test service or credential is available, and the
Batch is not authorized to add the required tooling surface.
```

Closing it means touching `docker-compose.yml`, `packages/backend/package.json`, `bun.lock`
and `Makefile` — protected tooling the Spec authorizes for no bounded path. The maintainer
accepted the gap and recorded it in a repository finding. CodeRabbit itself agreed, in the
thread:

> The MySQL-boundary finding remains valid and unresolved. The current PR does not
> authorize the required hermetic MySQL test infrastructure. This thread will remain open
> for a later Round.

### The status vocabulary has no cell for this

Per the Assigned Review Issue Batch contract, an Agent may set `resolved`, `invalid`, or
`failed`; `duplicated` is daemon-owned. Applied to an accepted gap:

| Status | What it would assert | Why it is wrong here |
| --- | --- | --- |
| `resolved` | the Batch fixed it | it did not; the test boundary is unchanged |
| `invalid` | false positive or does not apply | contract wording is explicit, and it *does* apply |
| `failed` | could not be safely completed | true, but permanently blocks push and Clean |

The Supervisor first invented `status: deferred`. Preflight rejected it, correctly:

```text
Review Issue artifact ".../round-003/issue_006.md" has unsupported status "deferred"
```

So `failed` is the only honest cell — and `failed` is exactly what stops the loop:

```text
Final Push blocked: 3 Unresolved Review Issue(s) remain.
Resolve Run ... reached Unresolved.
```

`roundfix watch --until-clean` requires zero unresolved Review Issues. With a permanently
`failed` issue, **`--until-clean` can never terminate**. Running it as instructed would loop
until it exhausted rounds or budget. The Supervisor had to abandon the documented command
and drive bounded single rounds by hand.

### The merge block that Roundfix never sees

Even after every actionable finding was closed, the PR would not merge. Two independent
gates, neither visible to Roundfix:

```text
mergeStateStatus: BLOCKED   reviewDecision: CHANGES_REQUESTED
```

1. **A stale `CHANGES_REQUESTED`.** CodeRabbit is incremental and refuses to re-review
   already-reviewed commits: *"CodeRabbit is an incremental review system and does not
   re-review already reviewed commits."* It therefore never issues a superseding `APPROVED`.
   The `CHANGES_REQUESTED` from round 004 stands forever. Dismissing it required a raw
   GitHub API call.

2. **`required_review_thread_resolution: true`.** The repository ruleset:

   ```json
   {"allowed_merge_methods":["squash"], "required_approving_review_count":0,
    "required_review_thread_resolution":true}
   ```

   Zero approvals needed — but every thread must be resolved, including the deliberately
   open one. Resolving it required a GraphQL mutation, and the Roundfix skill forbids manual
   thread resolution.

Both are ordinary GitHub configurations. Roundfix models neither, so a loop that is
otherwise complete stops one step from done.

### Three smaller stalls, same shape

**Preflight validates task files only when a Run starts.** Spec 0014's task files, authored
2026-08-01, carried four `pattern:` Context labels and one path with a trailing slash and a
parenthetical. Both are schema violations. They surfaced three days later, one at a time,
each aborting Run creation:

```text
Task Context entry: kind "pattern": expected "instruction" or "interface" label
Task Context entry: kind "interface": path "packages/backend/drizzle/migrations/": path must be clean
```

Authoring and validation are separated by however long the queue is. Worse, preflight
reports the **first** violation only, so a graph with N defects costs N start attempts.

**A pre-contract graph silently loses its QA gate.** Baseline 0.3.1 made QA an authored
terminal node and **removed `--qa` from `implement`**. Spec 0014's graph predates the
contract and declares no `qa:`. The contract says a legacy graph stays byte-identical — right
for a graph that already executed, wrong for one that has never run and is about to. Nothing
warned; the Supervisor noticed only by reading `implement --help` and finding the flag gone.
Left alone, the Spec would have run to completion with no gate, and the loop's "archive on a
QA pass" would have had no verdict to read.

**`roundfix baseline` deleted repository-authored rule blocks.** Both
`roundfix:repository-rule` blocks — the construct whose purpose is surviving baseline
writes — were removed. Most content was absorbed into the setup-owned block, but four rules
survived nowhere, including the absence-assertion form that a prior finding had already paid
for. No diff or report named what was dropped; the loss was found by reading `git diff`.

**A commit hook stricter than CI stalled a resolve Run with no diagnosis.** The batch commit
failed because `lint-staged` runs `oxlint --deny-warnings` while CI runs plain `oxlint`, and
a file sat at 509 counted lines against a 500 warning threshold. Roundfix surfaced it as a
raw dump ending in `Reverting to original state`, reported `watch failed after Run start`,
and left the batch work staged but uncommitted. The Run's own diagnosis did not distinguish
"the fix is wrong" from "the fix is fine and a hook rejected the commit" — opposite
remediations.

## Root cause

Roundfix's Review Issue lifecycle assumes every valid finding is eventually fixable inside
the Pull Request that received it. Real Specs violate that assumption routinely: a finding
can be correct, understood, and **out of the current Spec's authorization**. Protected-tooling
authorization is authored per Spec, up front, from the Spec's own scope — but review findings
arrive *after* authoring and can demand surface the author never anticipated. There is no
path from "reviewer asked for something this Spec may not do" to a settled, non-blocking
outcome.

Every downstream symptom follows: `failed` becomes a catch-all for two unrelated situations
(transient inability vs. accepted scope boundary); `--until-clean` inherits a termination
condition that cannot be met; and the Supervisor must leave the documented workflow exactly
when the workflow matters most.

The secondbrain frames this precisely. `wiki/concepts/agent-workflows-e-loop-engineering.md`
lists the seven parts of an operational loop and notes that a loop exists only with
**external verification, persistent state, and a stop condition** — *"sem stop condition,
consome tempo/tokens até quebrar ou drenar orçamento."* Roundfix has strong verification and
state; the stop condition is what is missing, and `--until-clean` is where the absence
concentrates. The same page's rule that *"uma verificação deve conseguir falhar quando nenhum
trabalho ocorreu"* has an unstated dual: **a verification must be able to pass when the
remaining work is legitimately out of scope.**

`wiki/concepts/agentic-coding-e-unknowns.md` explains why this cannot be authored away.
Protected-tooling needs discovered by a reviewer are *unknown unknowns* relative to the
Spec — the operator did not know to authorize them. The prescription is to build machinery
that **reveals** unknowns and asks for the decision, not to demand the author predict them.

## What would settle it

**1. A terminal status for an accepted gap.** Add a fifth Agent-settable status — `accepted`
or `waived` — meaning *valid finding, deliberately not fixed, maintainer-accepted*. It
requires `terminal_reason` and a durable pointer (the finding or ADR recording the debt), it
counts as settled for Clean and Final Push, and it is reported separately so it never reads
as success. `failed` then narrows to its real meaning: the Batch could not complete and a
retry might. Without this, `--until-clean` is unsafe to run on any Spec carrying known debt.

**2. Authorization escalation as a first-class artifact.** When triage concludes a finding is
valid but needs protected-tooling surface the Spec does not authorize, emit a structured
request naming the exact paths and the blocked finding, and settle the issue `accepted`
pending that decision. This session produced that request three times in prose, re-derived
from scratch each round, because there was nowhere structured to put it.

**3. Model the merge gates, or state plainly that they are the operator's.** Read
`mergeStateStatus`, `reviewDecision`, and the ruleset before declaring a Pull Request
mergeable; distinguish "blocked by findings" from "blocked by a stale review decision" and
"blocked by unresolved threads". A reviewer that never re-approves is a permanent block, not
a wait, and the loop should say so rather than looking finished.

**4. Validate Task Graphs at authoring time, and report every violation at once.** Expose the
preflight schema check as a standalone command (`roundfix validate --spec <slug>`) that
`write-tasks` runs before reporting the breakdown, and make it enumerate all violations in
one pass. Three days and two failed Run starts were spent on defects a parser catches
instantly.

**5. Detect the pre-contract graph that has never run.** When the binary enforces the QA
contract and a graph declares no `qa:`, distinguish *already executed* (leave byte-identical)
from *never executed* (refuse to start, or warn loudly). Silently running a Spec with no gate
defeats the loop's own archive precondition.

**6. Report what `baseline` removes.** Print a summary of repository-authored content dropped
or absorbed, and refuse to delete a `roundfix:repository-rule` block without naming it. A
construct that exists to survive baseline writes must not vanish silently.

**7. Classify hook rejection as its own Run outcome.** When a batch commit fails because a
pre-commit hook rejected it, say so, name the rule, and preserve the work with an explicit
resume path — rather than a raw dump that reads like the fix was wrong.

Items 1 and 2 are the ones that block the operator's goal. With them, a Spec carrying known,
recorded debt converges to Clean unattended and stops for the maintainer exactly once: to
decide whether to authorize the surface or accept the gap — a real decision of authority,
which is the only kind the loop should ever hand back.

## Spec pointer

Vortex PR #116 (Spec `0013-sync-cortes-de-dados`), rounds 001–005; PR #118 (Baseline 0.3.1);
Spec `0014-sync-exclusao-mutua-e-steps` preflight and QA-contract gap. Runs
`run_20260803T233608Z_42202f2607d83d97` (watch, Failed),
`run_20260804T000123Z_deb53e48049c7335` (resolve, Unresolved),
`run_20260804T003838Z_47626d78ccee7b18` (watch, MaxRoundsReached),
`run_20260804T104109Z_bfdbf329b6bdab25` (implement, Clean).

Positive control worth preserving: the authored terminal QA node worked exactly as intended.
Spec 0014's gate returned 15 pass / 0 findings / 1 environment-blocked, confirmed
`singleton` support the TechSpec had asserted from a precedent that does not exist, and
produced executable probe scripts as evidence rather than prose. The QA half of the contract
is not what needs repair.
