# Triage 2026-08-26 — the queue earns its tokens

Maintainer decision, 2026-08-26: the active queue keeps only work that reduces
Run errors, rework of finished Tasks, or token consumption in the implement
loop. Everything else is deferred with its evidence preserved. This record maps
every artifact the sweep touched; nothing was lost — deleted folders are
recoverable with `git log --diff-filter=D --oneline -- docs/specs/<slug>`.

Session evidence the criterion came from: 10 implement Runs on 2026-08-25/26
produced 1 Clean outcome directly, ~18 redundant Task re-executions (Unresolved
Runs leave completed Task commits unintegrated and the next Run redoes them),
5 QA gate reprovals of which 3 traced to defects in the auditing binary or the
authored artifacts rather than the code, and a profile configuration at
opus/xhigh that the built-in recommendation table prices at roughly twice
opus/high for about one point of result.

## Active queue, in order

1. **Unresolved-Run work reuse** (new, lean: techspec-direct, no product
   decision to record). An Unresolved Run leaves completed Task commits on its
   Run Branch; the next implement reads the checkout, sees `pending`, and redoes
   them. Either integrate completed work at Run end or refuse the next Run
   naming the integration command. Evidence: finding
   `2026-08-12-five-unresolved-runs-to-deliver-one-spec` (kept pending for
   adoption) plus this session's measurement above.
2. **0116-a-verdict-that-states-its-own-scope** (authored, source adopted,
   tooling authorization granted 2026-08-26). Stops vacuous Verification
   commands at authoring instead of one lost Run each.
3. **0097-a-wave-that-cannot-collide**. Two Tasks a graph declares independent
   editing the same file die at integration; measured first-hand here on
   2026-08-26 (0113 task_05 × task_07). Evidence kept pending:
   `2026-08-11-a-git-worktree-that-fails-only-under-load`.
4. **0105-the-gates-own-economics**. The gate is where the loop's cost
   concentrates (201 failed Tasks across five repositories in its PRD). Absorbs
   the whole of former Spec 0104 (the test-cache target, same authorization
   umbrella, Task-sized). Evidence kept pending/open:
   `2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks`,
   `2026-08-12-three-consecutive-specs-measure-the-loop`,
   `2026-08-07-the-only-gate-reports-green-on-a-red-suite`, backlog
   `2026-08-10-the-loop-is-measured-and-the-gate-is-where-it-costs`,
   `2026-08-08-go-clean-testcache-clears-a-cache-the-gate-does-not-use`,
   `2026-08-15-the-gate-runs-a-saturated-suite-inside-a-dense-run`.

Delivered and staying until archival: 0098 (Clean, 10/10) and 0113 (Clean, 9/9).

## Spec folders removed (PRD-only, no adopted sources, no inbound links)

Each maps to the evidence that stays on file; the PRD text lives in git history.

| Removed spec | Deferred evidence that carries the idea |
| --- | --- |
| 0099-a-tool-that-meets-the-repositorys-conventions | backlog 2026-08-12-archive-refuses-a-graph…, 2026-08-12-release-plan-requires-a-v-prefixed-tag |
| 0100-a-review-the-loop-always-asks-for | finding 2026-08-10-a-head-the-loop-did-not-push-is-a-head-nobody-reviews |
| 0101-a-terminal-branch-with-one-disposition | finding rollup 2026-08-06-rollup-run-lifecycle-and-branch-integrity |
| 0102-a-preflight-that-proves-what-the-run-needs | findings 2026-08-14-preflight-starves…, 2026-08-10-a-fake-adapter-goes-silent…, backlog 2026-08-08-a-failed-proof-appends-a-cleanup-error… |
| 0104-a-gate-that-cannot-certify-its-own-cache | absorbed into the 0105 queue entry above |
| 0106-a-decision-that-reaches-every-artifact | backlog 2026-08-10-an-excluded-artifact…, 2026-08-10-the-gate-accepts-a-manifest…, findings 2026-08-07-changing-the-http-contract…, 2026-08-07-greenfield-adoption… |
| 0107-the-authoring-rules-the-guides-do-not-carry | its one queue-relevant rule (wave independence by shared edit targets) rides with 0097; rest deferred |
| 0108-what-an-agent-loads-to-answer-one-question | backlog 2026-08-09-the-roundfix-skill-carries-nine-commands…, 2026-08-09-the-canonical-method-lives-in-rendered-guides… |
| 0109-what-a-session-consumed | backlog 2026-08-08-record-usage-per-agent-session |
| 0110-a-refresh-that-does-not-re-interview | finding 2026-08-07-the-setup-refresh-interviews-a-repository-that-already-answered |
| 0111-one-terminal-audit-across-a-runs-surfaces | deliberately last by its own design; no evidence lost |
| 0112-a-review-that-retires-on-its-own-facts | finding 2026-08-14-a-review-retires-on-whatever-the-object-store-happens-to-hold |
| 0115-an-archive-that-survives-its-own-move | backlog 2026-08-14-the-archive-destination-still-depends…, 2026-08-14-a-spec-archives-before-continuous-integration-answers |

## Closed as delivered

Backlog → done: a-verification-command-passes-only-by-exiting-zero (Spec 0095),
a-hook-failure-kills-a-run-that-already-verified-its-work (0098),
a-redirected-verification-hands-the-agent-an-empty-diagnostic (0096, ADR-0135),
a-refused-gate-writes-a-report-its-own-contract-rejects (0113),
the-changed-path-audit-does-not-know-sanctioned-fallout (ADR-0128/0129),
staging-a-task-commit-fails-when-the-task-deleted-a-file (0098),
the-vacuous-verification-event-lists-the-commands-that-passed (event carries
`commands` today, observed in this session's Run Events).

Findings → done: claude-agent-selections-are-never-proven (Preflight
exact-proves every configured tuple; observed on this session's
`profiles configure` and `doctor` runs).

## Deferred

Every backlog entry and finding not named above moves to `status: deferred`:
real observations, preserved verbatim, outside the active queue until the queue
above is delivered. The six 2026-08-06 rollups defer with their members.

## Secondbrain inbox

All 37 pending entries under `inbox/roundfix/` move to `_triaged/`. The
2026-08-06..08-19 entries were already minted into the backlog and findings
this document classifies — the inbox copies were residue of completed minting.
Of the five 2026-08-24/25 entries, two are minted as deferred findings needing
local reproduction (the QA Task Verification that is never executed; the Daemon
postcondition gap where a Task settles Clean while breaking the tree), and
three are recorded here as declined for the current cycle without a minted
artifact: the fluxus QA-override habit and the single-review-source question
are fleet-process discussions, not implement-loop defects; the
rollup-cannot-archive defect is real but is meta-tooling for closing findings,
which this sweep just did by hand at lower cost than fixing the mechanism.

## Profile change recorded

`.roundfixrc.yml` profiles now follow the binary's own recommendation table
(sol/high for backend-class work, opus/high only where design judgment leads,
luna/max for bounded docs/chore/review work). The prior configuration ran
opus/xhigh on every backend and frontend Task.
