---
date: 2026-07-31
supersedes: 2026-07-31-specs-0055-and-0060-shipped.md
---

# Handoff — paused with 0058 in remediation

The maintainer paused here to travel. Nothing is mid-flight: no Run is Active,
the working tree is clean, and every commit made in this session is on a remote.

## Where things stand

| Branch | Carries | Remote | PR |
| --- | --- | --- | --- |
| `ma/spec-non-regression-clauses` | PRD corrections for 0056, 0057, 0059 | yes | #59 open |
| `ma/npm-trusted-publishing-and-release-preflight` | 0058 PRD fix, ADR-0084, TechSpec, 5-task graph | yes | none yet |
| `roundfix/run-run_20260731T195234Z_e21abffeb5ebadcf` | 0058 implementation + QA report | yes (backup only) | not integrated |

`main` is unchanged at `99a0885`.

**The 0058 implementation is not on its `ma/` branch.** All five Tasks settled
`completed` inside the Run Worktree, but the Run reached `Unresolved` on a QA
`fail`, so nothing was integrated. The six commits were pushed to the run branch
purely as a backup before the pause — pushing was not integration and implies no
approval. The Run Worktree is retained at
`~/.roundfix/worktrees/roundfix-339f8dac/run_20260731T195234Z_e21abffeb5ebadcf`;
`roundfix reconcile` inspects it.

## What the queue actually is

0060 is **done** — archived under `_archived/` since PR #56/#57. The remaining
Specs are 0056, 0057, 0058, and **0059**. A request to implement "0060" in this
session was a slip for 0059.

The maintainer confirmed the execution shape: **0058 runs in parallel, and
0056 → 0057 → 0059 stay sequential.** The reason is mechanical — those three
each authorize edits to the same derived digest artifacts
(`catalog.digest`, `catalog.normalized.json`, `typescript-bun.json`, the
parity-corpus fixtures), so parallel branches would collide on every merge.
0058 touches only `.github/workflows/release.yml` and shares nothing.

Only 0058 is decomposed. 0056, 0057, and 0059 are still PRD-only and each needs
`write-techspec` → `write-tasks`.

## The standing acceptance criterion

The maintainer set a criterion that governs every remaining Spec: **a Spec must
evolve Roundfix, never regress it.** They came close to cancelling 0056 and 0057
over a suspicion that the two would break working behavior.

That was folded into both PRDs (PR #59) as three guarantees, dropped into the
existing nine sections rather than a new one, because the PRD section set is
fixed and QA verifies scope:

1. a characterization corpus captured **before** the change, as the regression
   gate — not a test written afterwards;
2. an explicit list of the only authorized behavior breaks (0056 declares
   exactly two: removal by omission ends, and a declined non-interactive write
   exits non-zero);
3. a false-positive bound on every new blocking gate — 0057's retention gate
   closes on evidence and never on doubt, and carrier classification may narrow
   a warning only on positive evidence.

**This has a scheduling consequence.** The characterization corpus must be
captured from the current binary before the first behavior-changing Task, so it
has to be `task_01` in both Specs. `write-tasks` will not infer that.

## 0058 — decomposed, implemented, QA failed

The Run implemented all five Tasks in ~55 minutes and the QA gate returned
`verdict: fail` — 21 pass, 3 fail, 3 environment-blocked, 0 skipped.
`make verify` passed with 2,941 tests across 24 packages.

The preflight was validated against the live registry: it classified all six
`0.0.2` coordinates as `used` and stopped before any mutation. The classifier
works.

Two of the three failures trace to the TechSpec and Task authoring, not to the
implementer:

- **QA-002 — failure attribution.** Any nonzero first `npm publish` enters the
  `identity:` branch and gets token-retried. A simulated network timeout was
  labelled `identity:`. That violates PRD Story 3, which requires the failure to
  distinguish identity from registry state. The TechSpec's interface sketch said
  "fails with ENEEDAUTH" loosely and `task_04`'s acceptance criteria let any
  failure through. Fix: only an authentication failure may enter the fallback.
- **QA-003 — runbook gap.** PRD Core Feature 4 requires token publication to be
  *disallowed* for the six packages after the window closes. `task_05` asked for
  the secret's removal but never for the registry-side shutdown, so the shipped
  runbook omits it.
- **QA-001 — a PRD promise reality cannot meet.** Core Feature 2 requires the
  preflight to detect missing ownership or trusted-publisher configuration.
  npm makes that impossible (below), ADR-0084 acknowledges it, but nobody
  amended the PRD, leaving the artifact set internally inconsistent. Fix is a
  PRD amendment, not code.

Three rows are environment-blocked and will stay that way: a live tagged
release, a remote `workflow_dispatch`, and a Pull Request journey. The PR row
unblocks as soon as a PR exists. **The end-to-end success metric — "one complete
release publishes all six packages through OIDC" — cannot be proven by any Task
or by QA.** It needs a live registry account and an irreversible publication, so
it stays a maintainer step.

### Facts about npm worth not rediscovering

- Trusted Publishing requires **npm 11.5.1+ on Node 22.14.0+**. The workflow
  pinned Node 20, whose bundled npm cannot perform an OIDC exchange at all.
- **npm never validates a trusted-publisher configuration until publish time.**
  There is no read-only endpoint, CLI command, or dry-run, and configuration is
  per-package. Identity therefore cannot be preflighted — which is the whole
  reason for ADR-0084's bounded per-coordinate token fallback.
- Publishing under Trusted Publishing generates provenance attestations
  automatically. Additive signed metadata; package contents are unchanged.

## A citation defect that was repo-wide

All four remaining PRDs cited `docs/agents/cli.md` for the Authentication and
HTTP constraint. That guide governs the CLI surface — command names, flags,
stream placement, exit codes — and says nothing about credentials or transport.
The operative rule is the secret-handling Normative Clause in
`docs/agents/agent-instructions.md`: never read, print, commit, or generate
secrets, and do not invent authentication, authorization, transport, or
deployment policy.

Fixed in all four. In 0058 the clause is *binding* — it constrains how the
retained token may be referenced, and became a requirement of `task_04`. In
0056, 0057, and 0059 the row stays "not applicable", now cited against the guide
that actually decides it.

Worth checking whether archived Specs carry the same wrong citation. They must
stay byte-identical, so this is an observation, not a backfill.

## Next, in order

1. **Merge PR #59.** 0056 and 0057 must be decomposed from the corrected PRDs.
2. **Remediate 0058**: amend PRD Core Feature 2 to match ADR-0084 (QA-001,
   Supervisor-owned), fix failure attribution in the workflow (QA-002), complete
   the runbook (QA-003), then re-run the QA gate and open the PR.
3. **Decompose 0056**, with the characterization corpus as `task_01`.
4. Then 0057, then 0059.

## Still-unrouted findings

Carried forward from the previous handoff and still open:

- `2026-07-31-a-rehearsal-task-can-settle-completed-without-rehearsing.md` — the
  Verification-that-passes-on-an-empty-tree defect. Fixing it in `write-tasks`
  before decomposing the remaining three would improve all three Task Graphs.
- `2026-07-30-baseline-digest-regeneration-cannot-bootstrap.md` — a direct
  prerequisite for 0057, which edits the Baseline catalog.
- `2026-07-30-failed-qa-runs-accumulate-unreleasable-run-branches.md` — now with
  one more instance, this session's retained Run Worktree.
- `2026-07-30-the-autonomous-loop-orders-qa-before-its-own-preconditions.md`
- `2026-07-29-qa-cycle-cost-is-cold-environments-and-agent-turns.md`
- The `archive-spec` scope defect (unscoped commit subject) remains unrouted
  with no finding.
